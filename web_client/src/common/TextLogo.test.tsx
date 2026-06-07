import React from "react";
import { render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it } from "vitest";

import { TextLogo, TextLogoSmall } from "./TextLogo";

afterEach(() => {
  delete (window as Partial<Pick<Window, "SITE_NAME">>).SITE_NAME;
});

describe("TextLogo", () => {
  it("renders the default full logo text", () => {
    render(<TextLogo />);

    expect(screen.getByText("IMG")).toBeInTheDocument();
    expect(screen.getByText("DD")).toBeInTheDocument();
  });

  it("renders the configured site name", () => {
    window.SITE_NAME = "Example CDN";

    render(<TextLogo />);

    expect(screen.getByText("Example CDN")).toBeInTheDocument();
  });
});

describe("TextLogoSmall", () => {
  it("renders the default small logo text", () => {
    render(<TextLogoSmall />);

    expect(screen.getByText("I")).toBeInTheDocument();
    expect(screen.getByText("D")).toBeInTheDocument();
  });
});
