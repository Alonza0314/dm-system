#!/usr/bin/env python3
"""Flatten a Postman Collection v2.1 export into a list of raw request records.

This performs only the mechanical extraction (walk folders, pull out method,
path segments, headers, body, and response examples). It does NOT decide
operationIds, schema names, or which path segments are parameters -- that
requires judgment and is done by the agent following SKILL.md, not this
script.

Usage:
    python3 extract_requests.py <postman-collection.json>

Output: JSON array on stdout, one object per leaf request:
{
  "tag": "Category",                 // top-level folder name, or null
  "name": "Get",                     // Postman item name
  "method": "GET",
  "path_segments": ["api", "category", "server"],
  "request_body_raw": "{...}"|null,  // raw string body of the *live* request, if any
  "responses": [
    {
      "name": "200",                // Postman response label (often the status code)
      "code": 200,
      "path_segments": ["api", "category", "server"],
      "body_raw": "{...}" | ""      // raw response body string
    },
    ...
  ]
}
"""
import json
import sys


def walk(items, tag=None):
    records = []
    for item in items:
        if "item" in item:
            # Folder: recurse, tagging with the top-level folder name only.
            next_tag = tag or item.get("name")
            records.extend(walk(item["item"], next_tag))
            continue

        request = item.get("request")
        if request is None:
            continue

        url = request.get("url") or {}
        body = request.get("body") or {}
        record = {
            "tag": tag,
            "name": item.get("name"),
            "method": request.get("method"),
            "path_segments": url.get("path") or [],
            "request_body_raw": body.get("raw") or None,
            "responses": [],
        }

        for resp in item.get("response", []):
            orig = resp.get("originalRequest") or {}
            orig_url = orig.get("url") or {}
            record["responses"].append({
                "name": resp.get("name"),
                "code": resp.get("code"),
                "path_segments": orig_url.get("path") or record["path_segments"],
                "body_raw": resp.get("body") or "",
            })

        records.append(record)
    return records


def main():
    if len(sys.argv) != 2:
        print(__doc__, file=sys.stderr)
        sys.exit(1)

    with open(sys.argv[1], "r", encoding="utf-8") as f:
        collection = json.load(f)

    records = walk(collection.get("item", []))
    json.dump(records, sys.stdout, indent=2)
    print()


if __name__ == "__main__":
    main()
