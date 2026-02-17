#!/usr/bin/env python3
"""
Fetch QCY product database from the QCY server API.

Outputs:
  - products_raw.json  : Raw server response
  - products.json      : Merged product database (vendorId -> model info + features)

Usage:
  python3 scripts/dump_products.py [--output-dir DIR]
"""

import argparse
import io
import json
import os
import sys
import time
import zipfile
from pathlib import Path

try:
    import requests
except ImportError:
    print("ERROR: 'requests' package required. Install: pip install requests", file=sys.stderr)
    sys.exit(1)

BASE_URL = "https://api.watch.qcy.com/"
HEADERS = {
    "Content-Type": "application/json; charset=utf-8",
    "lang": "en",
    "country": "US",
    "sys_": "android",
    "app_version": "4.0.7_695",
}

# Layout type codes from control panel JSONs
LAYOUT_TYPES = {
    2: "eq",
    3: "find_earphone",
    4: "settings",
    8: "channel_balance",
    9: "anc",
    10: "key_function",
    100: "reset",
    101: "device_name",
    102: "auto_off_timer",
}

# Settings item type codes
SETTINGS_TYPES = {
    1: "firmware_update",
    3: "reset_default",
    100: "factory_reset",
    200: "sleep_mode",
    203: "game_mode",
}


def fetch_product_list():
    """Fetch full product list from QCY server.
    
    The server returns a ZIP URL containing a JSON file with all products.
    """
    url = BASE_URL + "product/findProductList"
    print(f"POST {url}")
    resp = requests.post(url, data="", headers=HEADERS, timeout=30)
    resp.raise_for_status()
    meta = resp.json()

    product_url = None
    if isinstance(meta, dict):
        data = meta.get("data", {})
        if isinstance(data, dict):
            product_url = data.get("product")
        elif isinstance(data, str):
            product_url = data

    if not product_url:
        print("WARN: No product ZIP URL in server response, returning raw response")
        return meta

    print(f"GET {product_url}")
    zip_resp = requests.get(product_url, timeout=60)
    zip_resp.raise_for_status()

    z = zipfile.ZipFile(io.BytesIO(zip_resp.content))
    for name in z.namelist():
        if name.endswith(".json"):
            print(f"  Extracted {name} ({z.getinfo(name).file_size} bytes)")
            return json.loads(z.read(name))

    print("WARN: No JSON found in ZIP, returning raw response")
    return meta


def fetch_control_panel(vendor_id, firmware_version=""):
    """Fetch control panel JSON for a specific product.
    
    Uses vendorId (not flageID) as the modelId parameter for the API.
    Response structure: data.controlPanel.layouts
    """
    url = BASE_URL + "product/findControlPanels"
    form_data = {
        "modelId": str(vendor_id),
        "firmwareVersion": firmware_version,
    }
    print(f"  POST {url} vendorId={vendor_id}")
    try:
        resp = requests.post(url, data=form_data, headers={
            "lang": "en",
            "country": "US",
            "sys_": "android",
            "app_version": "4.0.7_695",
        }, timeout=30)
        resp.raise_for_status()
        result = resp.json()
        
        if result.get("code") == 200:
            data = result.get("data", {})
            if isinstance(data, dict):
                control_panel = data.get("controlPanel", {})
                if isinstance(control_panel, dict):
                    return control_panel
        
        return None
    except Exception as e:
        print(f"    WARN: Failed to fetch control panel for vendorId={vendor_id}: {e}")
        return None


def extract_features_from_layouts(layouts):
    """Extract supported features from control panel layout definitions."""
    features = {}
    for layout in layouts:
        layout_type = layout.get("type")
        type_name = LAYOUT_TYPES.get(layout_type, f"unknown_{layout_type}")

        if layout_type == 9:  # ANC
            modes = layout.get("modes", [])
            anc_modes = []
            for mode in modes:
                anc_modes.append({
                    "name": mode.get("name", ""),
                    "startcmdid": mode.get("startcmdid"),
                    "endcmdid": mode.get("endcmdid"),
                    "defaultcmd": mode.get("defaultcmd"),
                    "viewtype": mode.get("viewtype"),
                    "items": [
                        {
                            "name": item.get("name", ""),
                            "startcmdid": item.get("startcmdid"),
                            "endcmdid": item.get("endcmdid"),
                        }
                        for item in mode.get("items", [])
                    ] if mode.get("items") else None,
                })
            features["anc"] = {"modes": anc_modes}

        elif layout_type == 2:  # EQ
            features["eq"] = {
                "bands": layout.get("count", 10),
                "mindb": layout.get("mindb"),
                "maxdb": layout.get("maxdb"),
                "freq": layout.get("freq", ""),
                "characteristic": layout.get("character", ""),
                "presets": [p.get("name", "") for p in layout.get("sys_eq", [])],
            }

        elif layout_type == 10:  # Key function
            music = layout.get("music", {})
            events = music.get("event", [])
            key_events = []
            for event in events:
                event_name = event.get("name", "")
                left = event.get("left", {})
                if isinstance(left, dict):
                    functions = [item.get("name", "") for item in left.get("list", [])]
                    key_events.append({"name": event_name, "functions": functions})
            features["key_function"] = {"events": key_events}

        elif layout_type == 8:  # Channel balance
            features["channel_balance"] = True

        elif layout_type == 3:  # Find earphone
            features["find_earphone"] = True

        elif layout_type == 4:  # Settings
            items = layout.get("items", [])
            settings = []
            for item in items:
                item_type = item.get("type")
                type_name = SETTINGS_TYPES.get(item_type, f"type_{item_type}")
                settings.append({
                    "name": item.get("title", ""),
                    "type": type_name,
                    "cmdid": item.get("cmdid"),
                    "cmd": item.get("cmd"),
                })
            features["settings"] = settings

        elif layout_type == 101:
            features["device_name"] = True

        elif layout_type == 102:
            features["auto_off_timer"] = {
                "cmdid": layout.get("cmdid"),
                "repeat": layout.get("repeat"),
            }

    return features


def build_product_database(raw_data, output_dir):
    """Build merged product database from raw server data + control panels."""
    products = {}
    all_items = []

    # Collect items from all categories
    if isinstance(raw_data, dict):
        data = raw_data.get("data", raw_data)
        if isinstance(data, dict):
            for key in ["earphones", "wactchInfos", "accessory", "speaker", "product"]:
                items = data.get(key, [])
                if isinstance(items, list):
                    for item in items:
                        item["_category"] = key
                        all_items.append(item)
        elif isinstance(data, list):
            all_items = data
    elif isinstance(raw_data, list):
        all_items = raw_data

    print(f"\nFound {len(all_items)} products total")

    # Fetch control panels and build database
    for i, item in enumerate(all_items):
        vendor_id = item.get("vendorID") or item.get("vendorId")
        if vendor_id is None:
            continue

        title = item.get("title", "Unknown")
        sub_title = item.get("subTitle", "")
        model_id = item.get("modelId") or item.get("flageID") or item.get("id")
        category = item.get("_category", "unknown")

        entry = {
            "vendorId": vendor_id,
            "title": title,
            "subTitle": sub_title,
            "category": category,
        }

        if model_id:
            entry["modelId"] = model_id

        # Only fetch control panels for earphones (not watches/accessories)
        if category in ("earphones", "unknown") and vendor_id:
            control_panel = fetch_control_panel(vendor_id)
            if control_panel and isinstance(control_panel, dict):
                layouts = control_panel.get("layouts", [])
                if layouts:
                    entry["features"] = extract_features_from_layouts(layouts)

            # Rate limit
            time.sleep(0.3)

        # Use vendorId as key; if duplicate, keep first
        key = str(vendor_id)
        if key not in products:
            products[key] = entry
        else:
            # Merge: keep existing but note duplicate
            existing = products[key]
            if "aliases" not in existing:
                existing["aliases"] = []
            existing["aliases"].append(title)

    return products


def main():
    parser = argparse.ArgumentParser(description="Dump QCY product database")
    parser.add_argument("--output-dir", default="scripts/output", help="Output directory")
    parser.add_argument("--skip-panels", action="store_true", help="Skip fetching control panels")
    parser.add_argument("--local-only", action="store_true",
                        help="Use bundled qcy_earphone.json instead of server")
    args = parser.parse_args()

    output_dir = Path(args.output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)

    if args.local_only:
        # Use bundled JSON from decompiled APK
        local_path = Path("tmp/resources/package_1/res/raw/qcy_earphone.json")
        if not local_path.exists():
            print(f"ERROR: {local_path} not found", file=sys.stderr)
            sys.exit(1)
        with open(local_path) as f:
            raw_data = json.load(f)
        print(f"Loaded {local_path}")
    else:
        # Fetch from server
        raw_data = fetch_product_list()
        raw_path = output_dir / "products_raw.json"
        with open(raw_path, "w") as f:
            json.dump(raw_data, f, indent=2, ensure_ascii=False)
        print(f"Saved raw response to {raw_path}")

    if args.skip_panels:
        # Just extract product list without control panels
        products = {}
        items = raw_data.get("earphones", []) if "earphones" in raw_data else []
        if isinstance(raw_data, dict) and "data" in raw_data:
            data = raw_data["data"]
            if isinstance(data, dict):
                for key in ["earphones", "wactchInfos", "accessory", "speaker"]:
                    items.extend(data.get(key, []))

        for item in items:
            vendor_id = item.get("vendorID") or item.get("vendorId")
            if vendor_id is None:
                continue
            key = str(vendor_id)
            if key not in products:
                products[key] = {
                    "vendorId": vendor_id,
                    "title": item.get("title", "Unknown"),
                    "subTitle": item.get("subTitle", ""),
                }
    else:
        products = build_product_database(raw_data, output_dir)

    # Save merged database
    db_path = output_dir / "products.json"
    with open(db_path, "w") as f:
        json.dump(products, f, indent=2, ensure_ascii=False)
    print(f"\nSaved {len(products)} products to {db_path}")


if __name__ == "__main__":
    main()
