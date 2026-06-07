import { describe, expect, it } from "vitest";

import { humanFileSize } from "./humanizeFileSize";

describe("humanFileSize", () => {
  it("formats bytes without a unit prefix", () => {
    expect(humanFileSize(512)).toBe("512 B");
  });

  it("formats binary units by default", () => {
    expect(humanFileSize(1536)).toBe("1.5 KiB");
    expect(humanFileSize(1024 * 1024)).toBe("1.0 MiB");
  });

  it("formats metric units when requested", () => {
    expect(humanFileSize(1500, true)).toBe("1.5 kB");
  });

  it("respects the requested decimal places", () => {
    expect(humanFileSize(1536, false, 0)).toBe("2 KiB");
  });
});
