import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { configMocks, mockAnimationsApi } from "jsdom-testing-mocks";
import { afterAll, afterEach } from "vitest";

configMocks({ afterAll, afterEach });
mockAnimationsApi();

afterEach(() => {
  cleanup();
});
