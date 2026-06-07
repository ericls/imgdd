import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

import { notice, prompt, PromptContainer } from "./prompt";

describe("prompt", () => {
  it("renders the requested title and content", async () => {
    render(<PromptContainer />);

    const result = prompt({
      title: "Confirm delete",
      content: <p>Delete image.png?</p>,
    });

    expect(
      await screen.findByRole("dialog", { name: "Confirm delete" }),
    ).toBeInTheDocument();
    expect(screen.getByText("Delete image.png?")).toBeInTheDocument();

    await userEvent.click(screen.getByRole("button", { name: "OK" }));
    await expect(result).resolves.toBe(true);
  });

  it("resolves true when the yes button is clicked", async () => {
    const user = userEvent.setup();
    render(<PromptContainer />);

    const result = prompt({
      title: "Overwrite file",
      content: <p>Continue?</p>,
      yesText: "Overwrite",
    });

    await user.click(await screen.findByRole("button", { name: "Overwrite" }));

    await expect(result).resolves.toBe(true);
    await waitFor(() => {
      expect(screen.queryByRole("dialog")).toBeNull();
    });
  });

  it("resolves false when cancel is clicked", async () => {
    const user = userEvent.setup();
    render(<PromptContainer />);

    const result = prompt({
      title: "Discard changes",
      content: <p>Unsaved changes will be lost.</p>,
    });

    await user.click(await screen.findByRole("button", { name: "Cancel" }));

    await expect(result).resolves.toBe(false);
    await waitFor(() => {
      expect(screen.queryByRole("dialog")).toBeNull();
    });
  });

  it("resolves null when the dialog is dismissed", async () => {
    const user = userEvent.setup();
    render(<PromptContainer />);

    const result = prompt({
      title: "Dismissible prompt",
      content: <p>Press escape to dismiss.</p>,
    });

    expect(await screen.findByRole("dialog")).toBeInTheDocument();
    await user.keyboard("{Escape}");

    await expect(result).resolves.toBeNull();
    await waitFor(() => {
      expect(screen.queryByRole("dialog")).toBeNull();
    });
  });

  it("hides cancel for notices", async () => {
    render(<PromptContainer />);

    const result = notice("Upload complete", <p>Your image is ready.</p>);

    expect(
      await screen.findByRole("dialog", { name: "Upload complete" }),
    ).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Cancel" })).toBeNull();

    await userEvent.click(screen.getByRole("button", { name: "OK" }));
    await expect(result).resolves.toBe(true);
  });
});
