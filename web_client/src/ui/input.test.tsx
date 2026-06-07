import React from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

import { InputWithLabel } from "./input";

describe("InputWithLabel", () => {
  it("associates its label with the input and accepts text entry", async () => {
    const user = userEvent.setup();

    render(<InputWithLabel label="Display name" />);

    const input = screen.getByLabelText("Display name");
    await user.type(input, "Ada Lovelace");

    expect(input).toHaveValue("Ada Lovelace");
  });
});
