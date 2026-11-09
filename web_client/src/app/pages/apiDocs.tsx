import React from "react";
import { Link } from "react-router";
import { useAuth } from "~src/lib/auth";
import { humanFileSize } from "~src/lib/humanizeFileSize";
import { routes } from "~src/routes";

// Keep in sync with httpserver/api_handlers.go.
const MAX_IMAGES_PER_REQUEST = 5;
const MIN_EXPIRES_IN_SECONDS = 60;
const MAX_EXPIRES_IN_SECONDS = 31536000;

const ERROR_CODES: [string, number | string, string][] = [
  ["method_not_allowed", 405, "The request did not use POST."],
  [
    "invalid_api_key",
    401,
    "An Authorization header was sent, but it is malformed or the key is unknown or revoked. The request is not retried anonymously.",
  ],
  ["upload_disabled", 403, "This instance is not accepting new images."],
  ["invalid_request", 400, "The body is not valid multipart/form-data."],
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
    "429 (per image)",
    "Too many anonymous uploads from your IP address. Slow down, or use an API key.",
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

  return (
    <article className="prose prose-neutral dark:prose-invert prose-code:before:content-none prose-code:after:content-none max-w-3xl mx-auto px-4 py-8">
      <h1>HTTP API</h1>
      <p>
        Upload images to this site from scripts, apps, and other tools. The API
        works both anonymously and with an API key.
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
      <p>Authentication is optional.</p>
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
      </ul>
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
        "expiresAt": "2026-01-02T15:04:05Z"
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
          <code>url</code>.
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
