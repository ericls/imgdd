import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { Menu, type MenuSections } from "./menu";

function menuSections(overrides?: Partial<MenuSections>): MenuSections {
  return {
    children: [
      {
        id: "primary",
        items: [
          { id: "rename", children: "Rename" },
          { id: "disabled", children: "Disabled", disabled: true },
          { id: "copy", children: "Copy link" },
        ],
      },
      {
        id: "danger",
        items: [{ id: "delete", children: "Delete", variant: "danger" }],
      },
    ],
    ...overrides,
  };
}

describe("Menu", () => {
  it("renders sections and marks disabled and danger items", () => {
    render(<Menu menuSections={menuSections()} />);

    expect(screen.getByRole("button", { name: "Rename" })).toBeEnabled();
    expect(screen.getByRole("button", { name: "Disabled" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Delete" })).toHaveClass(
      "text-red-600",
    );
  });

  it("calls item actions with the click event and close handler", async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    const action = vi.fn(({ close }) => close?.());

    render(
      <Menu
        onClose={onClose}
        menuSections={{
          children: [
            {
              id: "primary",
              items: [{ id: "archive", children: "Archive", action }],
            },
          ],
        }}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Archive" }));

    expect(action).toHaveBeenCalledTimes(1);
    expect(action.mock.calls[0][0].event).toBeDefined();
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("moves keyboard focus through enabled items only", () => {
    render(<Menu menuSections={menuSections()} />);

    const menu = screen.getByRole("button", { name: "Rename" }).parentElement
      ?.parentElement;
    expect(menu).not.toBeNull();

    fireEvent.keyDown(menu!, { key: "ArrowDown" });
    expect(screen.getByRole("button", { name: "Rename" })).toHaveFocus();

    fireEvent.keyDown(menu!, { key: "ArrowDown" });
    expect(screen.getByRole("button", { name: "Copy link" })).toHaveFocus();

    fireEvent.keyDown(menu!, { key: "ArrowUp" });
    expect(screen.getByRole("button", { name: "Rename" })).toHaveFocus();

    fireEvent.keyDown(menu!, { key: "Tab", shiftKey: true });
    expect(screen.getByRole("button", { name: "Delete" })).toHaveFocus();
  });

  it("calls onClose when Escape is pressed", () => {
    const onClose = vi.fn();
    render(<Menu onClose={onClose} menuSections={menuSections()} />);

    const menu = screen.getByRole("button", { name: "Rename" }).parentElement
      ?.parentElement;
    fireEvent.keyDown(menu!, { key: "Escape" });

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("calls onClose for outside clicks after the menu is armed", async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();

    render(<Menu onClose={onClose} menuSections={menuSections()} />);

    await new Promise<void>((resolve) =>
      requestAnimationFrame(() => resolve()),
    );
    await user.click(document.body);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("does not close for clicks inside refs considered inside", async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    const insideRef = React.createRef<HTMLButtonElement>();

    render(
      <>
        <button ref={insideRef}>Trigger</button>
        <Menu
          onClose={onClose}
          refsConsideredInside={[insideRef]}
          menuSections={menuSections()}
        />
      </>,
    );

    await new Promise<void>((resolve) =>
      requestAnimationFrame(() => resolve()),
    );
    await user.click(screen.getByRole("button", { name: "Trigger" }));

    expect(onClose).not.toHaveBeenCalled();
  });
});
