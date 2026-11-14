import React from "react";
import { Link } from "react-router";
import { useAuth } from "~src/lib/auth";
import { humanFileSize } from "~src/lib/humanizeFileSize";
import { routes } from "~src/routes";

// Keep in sync with httpserver/api_handlers.go and api_album_images.go.
const MAX_IMAGES_PER_REQUEST = 5;
const MAX_ALBUM_IMAGES_PER_REQUEST = 100;
const MIN_EXPIRES_IN_SECONDS = 60;
const MAX_EXPIRES_IN_SECONDS = 31536000;
// Keep in sync with domainmodels/album.go.
const MAX_CAPTION_LENGTH = 2000;
const MAX_TITLE_LENGTH = 255;

const ERROR_CODES: [string, number | string, string][] = [
  ["method_not_allowed", 405, "The request did not use POST."],
  [
    "invalid_api_key",
    401,
    "An Authorization header was sent, but it is malformed or the key is unknown or revoked. The request is not retried anonymously. For every call but createImages, also when no key was sent.",
  ],
  ["upload_disabled", 403, "This instance is not accepting new images."],
  [
    "invalid_request",
    400,
    "The body could not be parsed: not multipart/form-data for createImages; not JSON or form data with at most one title for createAlbum; not JSON in the documented shape, or with unknown fields, for addImagesToAlbum and removeImagesFromAlbum.",
  ],
  [
    "invalid_title",
    400,
    `The album title is longer than ${MAX_TITLE_LENGTH} characters.`,
  ],
  [
    "request_too_large",
    413,
    "The request body is larger than the instance allows.",
  ],
  [
    "invalid_expires_in",
    400,
    `expires_in is not a whole number of seconds between ${MIN_EXPIRES_IN_SECONDS} and ${MAX_EXPIRES_IN_SECONDS}, or was sent more than once.`,
  ],
  ["invalid_album", 400, "album was sent more than once."],
  [
    "album_not_found",
    404,
    "The album does not exist, or the API key does not belong to its creator.",
  ],
  [
    "invalid_caption",
    400,
    "A caption is too long, the number of captions is neither one nor one per image, or captions were sent without album.",
  ],
  [
    "invalid_images",
    400,
    `images is empty, has more than ${MAX_ALBUM_IMAGES_PER_REQUEST} entries, repeats an id, or has an id that is not an image id.`,
  ],
  [
    "image_not_found",
    404,
    "addImagesToAlbum: an image does not exist, was deleted, has expired, or does not belong to the API key's owner. Nothing was added.",
  ],
  ["no_images", 400, 'No files were sent in the "image" field.'],
  [
    "too_many_images",
    400,
    `More than ${MAX_IMAGES_PER_REQUEST} files were sent.`,
  ],
  ["file_too_large", "400 (per image)", "The image exceeds your size limit."],
  [
    "mime_type_mismatch",
    "400 (per image)",
    "The file extension does not match the file's content, e.g. a PNG named photo.jpg.",
  ],
  ["empty_file", "400 (per image)", "The file is empty."],
  [
    "invalid_image",
    "400 (per image)",
    "The file is corrupt, truncated, or not a supported image.",
  ],
  ["unreadable_file", "400 (per image)", "The file could not be read."],
  [
    "bad_image",
    "400 (per image)",
    "The image was rejected by the instance's content check.",
  ],
  [
    "rate_limited",
    429,
    "Too many requests: anonymous uploads from your IP address (reported per image), or album calls by your account. Slow down.",
  ],
  ["internal_error", "500 (per image)", "Something went wrong on the server."],
];

function CodeBlock({ children }: { children: string }) {
  return (
    <pre className="overflow-x-auto text-sm">
      <code>{children}</code>
    </pre>
  );
}

export function APIDocsPage() {
  const { data: authData } = useAuth();
  const isAuthenticated = !!authData?.viewer.organizationUser;
  const uploadLimit = authData?.viewer.uploadLimitBytes;
  const origin = window.location.origin;
  const endpoint = `${origin}/api/v1/createImages`;
  const albumEndpoint = `${origin}/api/v1/createAlbum`;
  const addEndpoint = `${origin}/api/v1/addImagesToAlbum`;

  return (
    <article className="prose prose-neutral dark:prose-invert prose-code:before:content-none prose-code:after:content-none max-w-3xl mx-auto px-4 py-8">
      <h1>HTTP API</h1>
      <p>
        Upload images to this site from scripts, apps, and other tools. The API
        works both anonymously and with an API key.
      </p>
      <p>
        <strong>For AI agents and other programs:</strong> the same reference is
        available as plain markdown at{" "}
        <a href="/docs/api.md">
          <code>/docs/api.md</code>
        </a>
        , with a quick start for sharing screenshots in a live album. It is also
        linked from{" "}
        <a href="/llms.txt">
          <code>/llms.txt</code>
        </a>
        .
      </p>

      <h2 id="design">Design: intention based, not resource based</h2>
      <p>
        This API is <strong>intentionally intention based</strong>, not resource
        based (it is not REST). Each endpoint is named after the one thing you
        want to do, as a verb plus a noun, such as <code>createImages</code>.
        That means:
      </p>
      <ul>
        <li>
          Every call is a <code>POST</code> to{" "}
          <code>/api/v1/&lt;intention&gt;</code>.
        </li>
        <li>
          There are no resource URLs like <code>/images/:id</code> whose meaning
          changes with the HTTP method. The URL alone tells you what a request
          does.
        </li>
        <li>
          New capabilities are added as new intentions rather than by
          overloading existing ones. An intention keeps its behavior for the
          lifetime of its API version.
        </li>
      </ul>

      <h2 id="authentication">Authentication</h2>
      <p>
        Authentication is optional for <code>createImages</code> without an
        album. Every album call requires an API key.
      </p>
      <ul>
        <li>
          <strong>Anonymous:</strong> send no <code>Authorization</code> header.
          Anonymous uploads use the guest size limit, are rate limited per IP
          address, and are not tied to any account, so you cannot list or delete
          them later.
        </li>
        <li>
          <strong>With an API key:</strong> send{" "}
          <code>Authorization: Bearer &lt;your API key&gt;</code>. Images are
          added to your account, use your size limit, and are not rate limited
          per IP.{" "}
          {isAuthenticated ? (
            <>
              Manage your keys on the{" "}
              <Link to={routes.profile.apiKeys}>API keys page</Link>.
            </>
          ) : (
            <>
              <Link to="/auth">Sign in</Link> to create API keys.
            </>
          )}
        </li>
      </ul>
      <p>
        If you send an <code>Authorization</code> header that is malformed or
        holds an unknown or revoked key, the request fails with <code>401</code>
        . It is never treated as anonymous. Browser session cookies and session
        tokens are not accepted by the API.
      </p>
      <p>
        Keep API keys secret. Anyone with a key can upload images to your
        account. Revoke a key as soon as you suspect it leaked.
      </p>

      <h2 id="create-album">createAlbum</h2>
      <p>
        <code>POST /api/v1/createAlbum</code>
      </p>
      <p>
        Creates an album: a page that lists images in the order they were added
        and updates live, so others can watch while you upload. Requires an API
        key.
      </p>
      <h3>Request</h3>
      <p>
        An optional <code>title</code>, sent as JSON (
        <code>{`{"title": "..."}`}</code>), as form data, or omitted. Titles are
        trimmed and can be at most {MAX_TITLE_LENGTH} characters.
      </p>
      <CodeBlock>{`curl -H "Authorization: Bearer $IMGDD_API_KEY" \\
  --data-urlencode "title=Fix login page layout" \\
  ${albumEndpoint}`}</CodeBlock>
      <h3>Response</h3>
      <CodeBlock>{`{
  "album": {
    "id": "O96EBCnmxZoNOUE4rJDw0A",
    "url": "${origin}/albums/O96EBCnmxZoNOUE4rJDw0A",
    "title": "Fix login page layout"
  }
}`}</CodeBlock>
      <ul>
        <li>
          Pass <code>id</code> as the <code>album</code> field of{" "}
          <code>createImages</code>, <code>addImagesToAlbum</code> and{" "}
          <code>removeImagesFromAlbum</code>.
        </li>
        <li>
          <code>url</code> is the album page. Anyone with the link can view it,
          so share it only with people who should see the images. Album pages
          are not indexed by search engines.
        </li>
        <li>
          Only the album&apos;s creator can add or remove images. An image can
          be in any number of albums, with a different caption in each. An album
          page never shows which other albums its images are in.
        </li>
      </ul>

      <h2 id="create-images">createImages</h2>
      <p>
        <code>POST /api/v1/createImages</code>
      </p>
      <p>Creates one or more images from uploaded files.</p>

      <h3>Request</h3>
      <p>
        A <code>multipart/form-data</code> body with one file per{" "}
        <code>image</code> field, up to {MAX_IMAGES_PER_REQUEST} files.
      </p>
      <ul>
        <li>
          The file name&apos;s extension must match the file&apos;s content. SVG
          is not supported.
        </li>
        <li>EXIF metadata is removed before the image is stored.</li>
        <li>
          Each image must fit within your size limit
          {uploadLimit ? (
            <>
              {" "}
              (currently <strong>{humanFileSize(uploadLimit)}</strong>
              {isAuthenticated ? " for your account" : " for anonymous uploads"}
              )
            </>
          ) : null}
          .
        </li>
      </ul>

      <h4>Fields</h4>
      <ul>
        <li>
          <code>image</code> (required, repeatable): a file to upload.
        </li>
        <li>
          <code>expires_in</code> (optional): the number of seconds until the
          images expire, from {MIN_EXPIRES_IN_SECONDS} (one minute) to{" "}
          {MAX_EXPIRES_IN_SECONDS} (365 days). It applies to every image in the
          request. Without it, images never expire.
        </li>
        <li>
          <code>album</code> (optional): the <code>id</code> of an album you
          created. The images are added to it as they are stored. Only the
          album&apos;s creator can add to it; anyone else gets{" "}
          <code>album_not_found</code>.
        </li>
        <li>
          <code>caption</code> (optional, needs <code>album</code>): what each
          image shows in the album. Either one value for every image, or exactly
          one value per image, in the same order as the files. Each is trimmed
          and at most {MAX_CAPTION_LENGTH} characters, and is shown as plain
          text.
        </li>
      </ul>
      <p>
        <code>expires_in</code>, <code>album</code> and <code>caption</code> are
        checked before any file is stored. If one is invalid, nothing is
        uploaded.
      </p>
      <p>
        Expired images are deleted by a periodic background task, so an image
        can stay available for a short time after it expires.
      </p>

      <h3>Examples</h3>
      <p>Anonymous upload:</p>
      <CodeBlock>{`curl -F "image=@photo.png" ${endpoint}`}</CodeBlock>
      <p>Upload two images to your account:</p>
      <CodeBlock>{`curl -H "Authorization: Bearer $IMGDD_API_KEY" \\
  -F "image=@photo.png" \\
  -F "image=@screenshot.jpg" \\
  ${endpoint}`}</CodeBlock>
      <p>Upload an image that expires in one hour:</p>
      <CodeBlock>{`curl -F "image=@photo.png" -F "expires_in=3600" ${endpoint}`}</CodeBlock>
      <p>
        Upload a screenshot into an album, with a caption. With curl, send text
        fields with <code>--form-string</code>: <code>-F</code> treats{" "}
        <code>;</code> in a value as a separator and would cut the caption
        short.
      </p>
      <CodeBlock>{`curl -H "Authorization: Bearer $IMGDD_API_KEY" \\
  --form-string "album=O96EBCnmxZoNOUE4rJDw0A" \\
  --form-string "caption=step 2: after the fix; the button no longer wraps" \\
  -F "image=@step-2.png" \\
  ${endpoint}`}</CodeBlock>
      <p>JavaScript:</p>
      <CodeBlock>{`const form = new FormData();
form.append("image", file, file.name);

const res = await fetch("${endpoint}", {
  method: "POST",
  headers: { Authorization: \`Bearer \${apiKey}\` }, // omit for anonymous
  body: form,
});
const { results } = await res.json();`}</CodeBlock>

      <h3>Response</h3>
      <p>
        A JSON object with one entry in <code>results</code> per uploaded file,
        in the order they were sent. Each entry has either an <code>image</code>{" "}
        or an <code>error</code>; one file failing does not stop the others.
      </p>
      <CodeBlock>{`{
  "results": [
    {
      "filename": "photo.png",
      "image": {
        "id": "4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11",
        "filename": "photo.png",
        "url": "${origin}/image/8f1c2d7e-0b4a-4c9e-a1f2-3d4e5f6a7b8c.png",
        "mimeType": "image/png",
        "width": 1280,
        "height": 720,
        "byteSize": 183204,
        "expiresAt": "2026-01-02T15:04:05Z",
        "pageUrl": "${origin}/images/4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11",
        "caption": "step 2: after the fix",
        "albumUrl": "${origin}/albums/O96EBCnmxZoNOUE4rJDw0A"
      }
    },
    {
      "filename": "screenshot.jpg",
      "error": {
        "code": "mime_type_mismatch",
        "message": "Declared MIME type does not match detected MIME type"
      }
    }
  ]
}`}</CodeBlock>
      <ul>
        <li>
          <code>id</code> uniquely identifies the image. It is not part of{" "}
          <code>url</code>. Pass it to <code>addImagesToAlbum</code> to show the
          image in another album.
        </li>
        <li>
          <code>url</code> is the public link to the image. Use it as returned
          rather than building links yourself.
        </li>
        <li>
          <code>byteSize</code> is the stored size, after EXIF removal.
        </li>
        <li>
          <code>expiresAt</code> is when the image expires (UTC, RFC 3339), or{" "}
          <code>null</code> if it never expires.
        </li>
        <li>
          <code>pageUrl</code> is the image&apos;s page in the web UI.
        </li>
        <li>
          <code>caption</code> and <code>albumUrl</code> are <code>null</code>{" "}
          when not set.
        </li>
      </ul>

      <h3>Status codes</h3>
      <ul>
        <li>
          <code>200</code> when at least one image was created. Check each entry
          in <code>results</code> for errors.
        </li>
        <li>
          When no image was created, the status of the first failure, such as{" "}
          <code>400</code> or <code>429</code>.
        </li>
        <li>
          When the request as a whole is rejected, there is no{" "}
          <code>results</code>; the body is{" "}
          <code>{`{"error": {"code": "...", "message": "..."}}`}</code>.
        </li>
      </ul>

      <h2 id="add-images-to-album">addImagesToAlbum</h2>
      <p>
        <code>POST /api/v1/addImagesToAlbum</code>
      </p>
      <p>
        Adds images you already uploaded to an album, each with its own caption
        there. Requires an API key that owns the album and every image, and the
        images must not have expired. New images go to the end of the album in
        the order given. The request is all or nothing, and images already in
        the album are left as they are, so it is safe to repeat.
      </p>
      <CodeBlock>{`curl -H "Authorization: Bearer $IMGDD_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"album": "O96EBCnmxZoNOUE4rJDw0A",
       "images": [{"id": "4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11",
                   "caption": "the wrapping bug, before the fix"}]}' \\
  ${addEndpoint}`}</CodeBlock>
      <p>
        <code>images</code> holds 1 to {MAX_ALBUM_IMAGES_PER_REQUEST} distinct
        image ids; <code>caption</code> is optional. The response is the album
        (as from <code>createAlbum</code>) and <code>added</code>, the number of
        images that were not in the album yet.
      </p>

      <h2 id="remove-images-from-album">removeImagesFromAlbum</h2>
      <p>
        <code>POST /api/v1/removeImagesFromAlbum</code>
      </p>
      <p>
        Removes images from an album you created. The images are not deleted and
        stay in any other albums. Ids not in the album are ignored.
      </p>
      <CodeBlock>{`{"album": "O96EBCnmxZoNOUE4rJDw0A", "images": ["4a0b6f0e-1c5d-4f3a-9a53-2d1f6b8e0c11"]}`}</CodeBlock>
      <p>
        The response is the album and <code>removed</code>, the number of images
        taken out of it.
      </p>

      <h2 id="errors">Error codes</h2>
      <p>
        Match on <code>code</code>, which is stable. <code>message</code> is for
        people and may change.
      </p>
      <table>
        <thead>
          <tr>
            <th>Code</th>
            <th>Status</th>
            <th>Meaning</th>
          </tr>
        </thead>
        <tbody>
          {ERROR_CODES.map(([code, status, meaning]) => (
            <tr key={code}>
              <td>
                <code>{code}</code>
              </td>
              <td className="whitespace-nowrap">{status}</td>
              <td>{meaning}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </article>
  );
}
