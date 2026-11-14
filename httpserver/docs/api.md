# imgdd HTTP API

Upload images to {{BASE_URL}} from scripts, tools and AI agents. This is the
plain-text version of {{WEB_UI_URL}}/docs/api, written to be read by programs as
well as people.

- Base URL: `{{BASE_URL}}`
- Every call is a `POST` to `{{BASE_URL}}/api/v1/<intention>`.
- Responses are JSON. Errors have a stable `code` to match on and a `message`
  for people.

## Contents

- [Quick start: screenshots from an AI agent](#quick-start-screenshots-from-an-ai-agent)
- [Design](#design)
- [Authentication](#authentication)
- [createAlbum](#createalbum)
- [createImages](#createimages)
- [addImagesToAlbum](#addimagestoalbum)
- [removeImagesFromAlbum](#removeimagesfromalbum)
- [Error codes](#error-codes)

## Quick start: screenshots from an AI agent

Use this to show a person the screenshots you take while you work, on one page
that updates live.

You need two environment variables from the person you work for:

- `IMGDD_URL`: `{{BASE_URL}}`
- `IMGDD_API_KEY`: an API key, created at `{{WEB_UI_URL}}/profile/api-keys`.
  Albums need a key. Plain uploads work without one.

1. Create an album once per task. Give its `url` to the person; it is the page
   they watch.

   ```sh
   curl -s -H "Authorization: Bearer $IMGDD_API_KEY" \
     --data-urlencode "title=Fix login page layout" \
     "$IMGDD_URL/api/v1/createAlbum"
   ```

   ```json
   {"album": {"id": "O96EBCnmxZoNOUE4rJDw0A", "url": "{{WEB_UI_URL}}/albums/O96EBCnmxZoNOUE4rJDw0A", "title": "Fix login page layout"}}
   ```

2. Upload each screenshot into the album with a caption that says what it shows.

   ```sh
   curl -s -H "Authorization: Bearer $IMGDD_API_KEY" \
     --form-string "album=O96EBCnmxZoNOUE4rJDw0A" \
     --form-string "caption=step 2: after the fix; the button no longer wraps at 375px" \
     -F "image=@step-2.png" \
     "$IMGDD_URL/api/v1/createImages"
   ```

   The album page shows new images within a few seconds, without a reload.
   Each result also has a `url` to the image itself and a `pageUrl` to share a
   single screenshot.

Tips:

- **Use `--form-string` for text fields** such as `caption` and `album`.
  curl's `-F` treats `;` in a value as a separator, so
  `-F "caption=before; after"` silently uploads only `before`. Use `-F` only
  for files (`image=@file.png`).
- **Look at each screenshot before uploading it.** A blank or wrong image with
  a good caption is worse than no image. Make sure each image shows what its
  caption says.
- **Wait for the page to settle before capturing.** One-shot captures such as
  `chrome --headless --screenshot` can fire before a page finishes rendering
  or scrolling, which gives a blank image. Prefer a real browser session
  (Playwright, Puppeteer or the Chrome DevTools protocol) and wait for content.
- **Set `expires_in`** (seconds) if the screenshots are only needed for a while.
- Upload up to {{MAX_IMAGES}} images per request. Give one `caption` for all of
  them, or one per image in order.
- **An image can be in several albums.** To also show a screenshot you already
  uploaded in another album, such as a summary of the task, use
  `addImagesToAlbum` with its `id` instead of uploading it again. It can have a
  different caption there.

## Design

The API is intention based, not resource based. Each endpoint is named after
the one thing you want to do, such as `createImages`, and is always called with
`POST` at `/api/v1/<intention>`. There are no resource URLs whose meaning
changes with the HTTP method. New capabilities are added as new intentions.

## Authentication

Send `Authorization: Bearer <api key>`, or no `Authorization` header at all for
an anonymous request.

- Anonymous requests can only call `createImages`, without `album` or `caption`. They use the guest size
  limit, are rate limited per IP address, and are not tied to any account.
- With an API key, images are added to the key owner's account and use their
  size limit. Keys are created and revoked at `{{WEB_UI_URL}}/profile/api-keys`.
- A malformed `Authorization` header, or an unknown or revoked key, fails with
  `401 invalid_api_key`. It is never treated as anonymous.
- Browser session cookies and session tokens are not accepted.

Keep keys secret: anyone with a key can upload to its owner's account and add
to their albums.

## createAlbum

`POST {{BASE_URL}}/api/v1/createAlbum`

Creates an album: a page that lists images in the order they were added and
updates live. Requires an API key.

### Request

An optional `title`, sent as JSON (`{"title": "..."}`), as form data, or
omitted entirely. Titles are trimmed and can be at most {{MAX_TITLE}}
characters. An album without a title is shown as "Untitled album".

### Response

```json
{
  "album": {
    "id": "O96EBCnmxZoNOUE4rJDw0A",
    "url": "{{WEB_UI_URL}}/albums/O96EBCnmxZoNOUE4rJDw0A",
    "title": "Fix login page layout"
  }
}
```

- `id` is what you pass as the `album` field of `createImages`,
  `addImagesToAlbum` and `removeImagesFromAlbum`.
- `url` is the album page. Anyone with the link can view it, so share it only
  with people who should see the images. Album pages are not indexed by search
  engines.
- Only the album's creator can add images to it or remove them. Images join an
  album when they are uploaded into it, or with `addImagesToAlbum`. An image
  can be in any number of albums, with a different caption in each. An album
  page never shows which other albums its images are in.

Album creation is rate limited per user.

## createImages

`POST {{BASE_URL}}/api/v1/createImages`

Creates one or more images from uploaded files. Works with or without an API
key.

### Request

A `multipart/form-data` body with these fields:

| Field | Required | Description |
| --- | --- | --- |
| `image` | yes, repeatable | A file to upload, up to {{MAX_IMAGES}} per request. The file name's extension must match its content; SVG is not supported. EXIF metadata is removed. |
| `expires_in` | no | Seconds until the images expire, from {{MIN_EXPIRES_IN}} to {{MAX_EXPIRES_IN}} (365 days). Applies to every image in the request. Without it, images never expire. Expired images are deleted by a periodic task, so they can stay available briefly after expiring. |
| `album` | no | The `id` of an album you created. The images are added to it as they are stored. Requires the API key of the album's creator; anyone else gets `album_not_found`. |
| `caption` | no, needs `album` | What each image shows, displayed under it on the album page. Either one value for every image, or exactly one value per image, in the same order as the files. Each is trimmed and at most {{MAX_CAPTION}} characters. Captions are shown as plain text and belong to the image's place in this album. |

The `expires_in`, `album` and `caption` fields are checked before any file is
stored. If one is invalid, nothing is uploaded.

### Response

One entry in `results` per uploaded file, in request order. Each entry has
either an `image` or an `error`; one file failing does not stop the others.

```json
{
  "results": [
    {
      "filename": "step-2.png",
      "image": {
        "id": "4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11",
        "filename": "step-2.png",
        "url": "{{EXAMPLE_IMAGE_URL}}",
        "mimeType": "image/png",
        "width": 1280,
        "height": 720,
        "byteSize": 183204,
        "expiresAt": null,
        "pageUrl": "{{WEB_UI_URL}}/images/4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11",
        "caption": "step 2: after the fix; the button no longer wraps at 375px",
        "albumUrl": "{{WEB_UI_URL}}/albums/O96EBCnmxZoNOUE4rJDw0A"
      }
    },
    {
      "filename": "notes.jpg",
      "error": {"code": "mime_type_mismatch", "message": "Declared MIME type does not match detected MIME type"}
    }
  ]
}
```

- `url` is the image file. `pageUrl` is its page in the web UI. Use both as
  returned rather than building links yourself.
- `byteSize` is the stored size, after EXIF removal.
- `id` is what you pass to `addImagesToAlbum` and `removeImagesFromAlbum`.
- `expiresAt` (UTC, RFC 3339), `caption` and `albumUrl` are `null` when not set.

### Status codes

- `200` when at least one image was created. Check each entry for errors.
- When no image was created, the status of the first failure, such as `400`
  or `429`.
- When the request as a whole is rejected, there is no `results`; the body is
  `{"error": {"code": "...", "message": "..."}}`.

## addImagesToAlbum

`POST {{BASE_URL}}/api/v1/addImagesToAlbum`

Adds images you already uploaded to an album, each with its own caption there.
Requires an API key that owns both the album and every image.

### Request

A JSON body:

```json
{
  "album": "O96EBCnmxZoNOUE4rJDw0A",
  "images": [
    {"id": "4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11", "caption": "the wrapping bug, before the fix"}
  ]
}
```

- `images` holds 1 to 100 distinct image `id`s, as returned by `createImages`.
  `caption` is optional, trimmed, and at most {{MAX_CAPTION}} characters.
- New images go to the end of the album, in the order given. Images already
  in the album are left as they are, caption included, so the call is safe to
  repeat.
- The request is all or nothing: if any image is not yours, or has expired,
  nothing is added.

```sh
curl -s -H "Authorization: Bearer $IMGDD_API_KEY" -H "Content-Type: application/json" \
  -d '{"album": "O96EBCnmxZoNOUE4rJDw0A", "images": [{"id": "4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11"}]}' \
  "$IMGDD_URL/api/v1/addImagesToAlbum"
```

### Response

```json
{"album": {"id": "O96EBCnmxZoNOUE4rJDw0A", "url": "{{WEB_UI_URL}}/albums/O96EBCnmxZoNOUE4rJDw0A", "title": "Fix login page layout"}, "added": 1}
```

`added` counts only images that were not in the album yet.

## removeImagesFromAlbum

`POST {{BASE_URL}}/api/v1/removeImagesFromAlbum`

Removes images from an album. The images are not deleted and stay in any other
albums. Requires the API key of the album's creator.

### Request

```json
{"album": "O96EBCnmxZoNOUE4rJDw0A", "images": ["4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11"]}
```

`images` holds 1 to 100 distinct image `id`s. Ids that are not in the album are
ignored.

### Response

```json
{"album": {"id": "O96EBCnmxZoNOUE4rJDw0A", "url": "{{WEB_UI_URL}}/albums/O96EBCnmxZoNOUE4rJDw0A", "title": "Fix login page layout"}, "removed": 1}
```

## Error codes

Match on `code`; `message` may change.

| Code | Status | Endpoint | Meaning |
| --- | --- | --- | --- |
| `method_not_allowed` | 405 | all | The request did not use `POST`. |
| `invalid_api_key` | 401 | all | The `Authorization` header is malformed or the key is unknown or revoked. For every call but `createImages`, also when no key was sent. |
| `invalid_request` | 400 | all | The body could not be parsed: not multipart for `createImages`; not JSON or form data with at most one `title` for `createAlbum`; not JSON in the documented shape, or with unknown fields, for `addImagesToAlbum` and `removeImagesFromAlbum`. |
| `rate_limited` | 429 | all | Too many requests: anonymous uploads per IP (reported per image), or album calls per user. Slow down. |
| `internal_error` | 500 | all | Something went wrong on the server. |
| `invalid_title` | 400 | createAlbum | The title is longer than {{MAX_TITLE}} characters. |
| `upload_disabled` | 403 | createImages | This instance is not accepting new images. |
| `request_too_large` | 413 | createImages | The request body is larger than the instance allows. |
| `no_images` | 400 | createImages | No files were sent in the `image` field. |
| `too_many_images` | 400 | createImages | More than {{MAX_IMAGES}} files were sent. |
| `invalid_expires_in` | 400 | createImages | `expires_in` is not a whole number between {{MIN_EXPIRES_IN}} and {{MAX_EXPIRES_IN}}, or was sent more than once. |
| `invalid_album` | 400 | createImages | `album` was sent more than once. |
| `album_not_found` | 404 | createImages, addImagesToAlbum, removeImagesFromAlbum | The album does not exist, or the API key does not belong to its creator. |
| `invalid_caption` | 400 | createImages, addImagesToAlbum | A caption is too long, the number of captions is neither one nor one per image, or captions were sent without `album`. |
| `invalid_images` | 400 | addImagesToAlbum, removeImagesFromAlbum | `images` is empty, has more than 100 entries, repeats an id, or has an id that is not an image id. |
| `image_not_found` | 404 | addImagesToAlbum | An image does not exist, was deleted, has expired, or does not belong to the API key's owner. Nothing was added. |
| `file_too_large` | 400, per image | createImages | The image exceeds your size limit. |
| `mime_type_mismatch` | 400, per image | createImages | The file extension does not match the content, e.g. a PNG named `photo.jpg`. |
| `empty_file` | 400, per image | createImages | The file is empty. |
| `invalid_image` | 400, per image | createImages | The file is corrupt, truncated or not a supported image. |
| `unreadable_file` | 400, per image | createImages | The file could not be read. |
| `bad_image` | 400, per image | createImages | The instance's content check rejected the image. |
