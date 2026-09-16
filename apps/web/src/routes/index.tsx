import { createFileRoute } from "@tanstack/react-router";

import { Menu } from "../modules/menu";
import { Heading } from "../shared/components/heading";
import { Text } from "../shared/components/text";

function HomePage() {
  return (
    <main className="grid min-h-dvh place-items-center bg-black">
      <div className="flex w-full flex-col items-center px-6 text-center">
        <Heading as="h1" size="md" className="max-w-3xl text-white">
          Подбор роботизированных решений
        </Heading>
        <Text size="md" className="mt-6 max-w-xl text-white/70">
          Сравнивайте роботизированные решения по задачам, характеристикам и
          условиям внедрения.
        </Text>
      </div>
      <Menu />
    </main>
  );
}

export const Route = createFileRoute("/")({
  component: HomePage,
});
