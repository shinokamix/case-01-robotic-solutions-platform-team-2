import { createFileRoute } from "@tanstack/react-router";

function HomePage() {
  return (
    <main>
      <h1>Платформа подбора роботизированных решений</h1>
    </main>
  );
}

export const Route = createFileRoute("/")({
  component: HomePage,
});
