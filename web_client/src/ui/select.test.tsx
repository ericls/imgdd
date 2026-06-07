import React from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

import { SelectWithLabel } from "./select";

describe("SelectWithLabel", () => {
  it("associates its label with the select and accepts selection changes", async () => {
    const user = userEvent.setup();

    render(
      <SelectWithLabel label="Language" defaultValue="en">
        <option value="en">English</option>
        <option value="ko">Korean</option>
      </SelectWithLabel>,
    );

    const select = screen.getByLabelText("Language");
    await user.selectOptions(select, "ko");

    expect(select).toHaveValue("ko");
  });
});
