import { RouterProvider, createRouter } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import { routeTree } from "./routeTree.gen.ts";

import "./index.css";

const router = createRouter({ routeTree });
const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error('Element with id "root" was not found');
}

createRoot(rootElement).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);
