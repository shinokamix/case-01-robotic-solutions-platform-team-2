import type { ReactNode } from "react";

import { Heading } from "../../../shared/components/heading";
import { Text } from "../../../shared/components/text";

type AuthLayoutProps = {
  title: string;
  description: string;
  children: ReactNode;
};

export function AuthLayout({ title, description, children }: AuthLayoutProps) {
  return (
    <main className="flex min-h-dvh items-center justify-center bg-neutral-950 px-4 py-8 sm:py-12">
      <section className="w-full max-w-[26rem] rounded-2xl border border-white/10 bg-neutral-900/95 p-6 shadow-xl sm:p-8">
        <header className="mb-7">
          <Heading as="h1" size="sm">
            {title}
          </Heading>
          <Text size="md" className="mt-2 max-w-xs text-white/55">
            {description}
          </Text>
        </header>
        {children}
      </section>
    </main>
  );
}
