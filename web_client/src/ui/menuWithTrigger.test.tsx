import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { MenuWithTrigger } from "./menuWithTrigger";

describe("MenuWithTrigger", () => {
  it("opens the menu from the trigger and closes when an item calls close", async () => {
    const user = userEvent.setup();
    const action = vi.fn(({ close }) => close?.());

    render(
      <MenuWithTrigger
        trigger={<button type="button">Open menu</button>}
        menuSections={{
          children: [
            {
              id: "primary",
              items: [{ id: "copy", children: "Copy link", action }],
            },
          ],
        }}
      />,
    );

    expect(screen.queryByRole("button", { name: "Copy link" })).toBeNull();

    await user.click(screen.getByRole("button", { name: "Open menu" }));
    await user.click(screen.getByRole("button", { name: "Copy link" }));

    expect(action).toHaveBeenCalledTimes(1);
    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "Copy link" })).toBeNull();
    });
  });

  it("disables the trigger when there are no menu items", () => {
    render(
      <MenuWithTrigger
        trigger={<button type="button">Open menu</button>}
        menuSections={{ children: [{ id: "empty", items: [] }] }}
      />,
    );

    expect(screen.getByRole("button", { name: "Open menu" })).toBeDisabled();
  });
});
