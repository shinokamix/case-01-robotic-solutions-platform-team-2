import { createFileRoute } from "@tanstack/react-router";
import { Menu } from "../modules/menu";

function HomePage() {
  return (
    <main className="relative min-h-dvh bg-black">
      <Menu />
    </main>
  );
}

export const Route = createFileRoute("/")({
  component: HomePage,
});
