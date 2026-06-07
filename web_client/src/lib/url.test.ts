import { describe, expect, it } from "vitest";

import { absoluteURL } from "./url";

describe("absoluteURL", () => {
  it("returns absolute HTTP URLs unchanged", () => {
    expect(absoluteURL("https://example.com/image.png")).toBe(
      "https://example.com/image.png",
    );
  });

  it("prefixes root-relative URLs with the current origin", () => {
    expect(absoluteURL("/uploads/image.png")).toBe(
      `${window.location.origin}/uploads/image.png`,
    );
  });
});
