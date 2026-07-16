---
name: postman-to-openapi
description: Use when asked to convert a Postman collection export (collection.json) into an OpenAPI YAML spec, or to update/sync an existing openapi.yaml so it matches a newer Postman collection. Triggers on "postman to openapi", ".postman_collection.json", "update openapi.yaml from postman", "sync swagger with postman".
---

# Postman to OpenAPI

## Overview

Merges the requests recorded in a Postman Collection v2.1 export into an OpenAPI 3.x YAML file. This is a **merge into existing hand-authored YAML**, not a blind regenerate — an existing `openapi.yaml` usually has curated `operationId`s, shared `components/schemas`, and descriptions that a naive converter would clobber or duplicate. The mechanical parts (walking folders, listing requests, diffing response bodies) are scripted; deciding operationIds, schema names/reuse, and path parameters is a judgment call you make.

## When to Use

- User gives you a Postman JSON path and an OpenAPI YAML path and asks to convert/update.
- Only one file is given: ask for the other (both a Postman collection AND a target OpenAPI file are required inputs — this skill updates an existing spec, it does not scaffold one from nothing beyond `openapi.yaml`'s minimal skeleton if the file doesn't exist yet).

## Process

1. **Extract mechanically.** Run `scripts/extract_requests.py <postman.json>` to get a flat JSON list of every leaf request: tag (top folder name), Postman item name, method, path segments, request body, and every recorded response (code + path segments + body). Postman folders map to OpenAPI tags; nested sub-folders are flattened to the top-level folder name.

2. **Read the existing OpenAPI file** (or start from a minimal skeleton — `openapi: 3.0.3`, `info`, empty `paths`/`components/schemas` — if it doesn't exist). Note what's already there: existing paths, existing schema names, existing `security`/`securitySchemes`.

3. **Detect path parameters.** A request's own `path_segments` is one example; its `responses[].path_segments` are more examples of the *same* endpoint recorded at different times. If a segment's value varies across these examples (e.g. `server`, `1`, `iphone` all appear in the same slot), that segment is a path parameter — replace it with `{name}` (pick a name from the resource/context, e.g. singular of the tag, or read the backend's own error messages in response bodies for a hint like `"category iphone not found"` → the segment is a category name). If a segment is constant across every example, it's a literal path segment.

4. **Decide operationId and summary.** Prefer the Postman item name (`Login` → `login`, `GetAll` → `listCategories`). Postman collections often have lazy/generic names like `New Request` or `Copy of Get` — when the name isn't usable, derive one from method + resource (`DELETE /api/category/{name}` → `deleteCategory`).

5. **Build/reuse schemas.**
   - For each distinct response body shape (grouped by status code across the endpoints), check whether an existing `components/schemas` entry already matches that shape — reuse it via `$ref` rather than duplicating (e.g. a generic `{"message": string}` body should reuse a shared `MessageResponse`-style schema if one already exists in the target file).
   - For genuinely new shapes, add a new schema named `<Operation><Request|Response>` (e.g. `CategoryCreateRequest`) and infer types from the example JSON values (`"id": 1` → `integer`, strings → `string`). Nested arrays of objects become their own named schema (e.g. `categories: Category[]`) so the item shape can be reused across list/detail endpoints.
   - Nest under the right response code (2xx from the success example(s), plus every documented error code with its own example body).

6. **Infer auth requirements from Postman scripts**, not just the request headers. Check the collection's collection-level `event` (`prerequest`) and per-item `event` scripts — a script that adds an `Authorization: Bearer <token>` header to every request except one (commonly the login request) means every *other* endpoint requires auth even though the raw `request.header` array is empty (Postman injects it at runtime, so the exported JSON won't show it as a static header). Add a `bearerAuth` (or matching) entry under `components/securitySchemes` if missing, and apply `security: [{bearerAuth: []}]` to the affected operations, leaving the login/public operations without it.

7. **Write the merged result** with the Edit tool, preserving formatting/ordering conventions already used in the file (e.g. blank line between path items, schema property ordering). Don't reformat parts of the file you didn't need to touch.

8. **Validate**: parse the resulting YAML (`python3 -c "import yaml; yaml.safe_load(open('openapi.yaml'))"`) to confirm it's syntactically valid before considering the task done.

## Quick Reference: Postman → OpenAPI mapping

| Postman | OpenAPI |
|---|---|
| Top-level folder | `tags` on each operation in that folder |
| `request.method` + `url.path` | HTTP method + `paths.<path>` |
| Path segment that varies across recorded examples | `{param}` path parameter |
| `request.body.raw` (mode `raw`, language `json`) | `requestBody.content.application/json.schema` |
| `response[].body` per status code | `responses.<code>.content.application/json.schema` (+ `example`) |
| Collection/folder/item `prerequest` script adding a header | `security` + `components.securitySchemes` |
| Collection `variable` (e.g. `{{BACKEND_URL}}`) | `servers` — only if it has a concrete non-empty value; skip if blank/templated |

## Common Mistakes

- **Regenerating the whole file instead of merging.** Always read the existing spec first; reuse its schema names and conventions.
- **Missing auth.** The exported JSON's `request.header` is often empty even for authenticated endpoints because Postman adds auth headers via a pre-request script at runtime, not as a static header. Always check `event[].script` at the collection and item level.
- **Treating every path segment as literal.** Compare a request's segments against all of its own response examples' segments before deciding a segment is fixed.
- **Duplicating a schema that already exists** (e.g. adding a second `{message: string}` schema instead of reusing the existing `MessageResponse`).
