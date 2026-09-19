import { createRootRoute, Outlet } from "@tanstack/react-router";
import { MotionConfig } from "motion/react";

export const Route = createRootRoute({
  component: () => (
    <MotionConfig reducedMotion="user">
      <Outlet />
    </MotionConfig>
  ),
});
