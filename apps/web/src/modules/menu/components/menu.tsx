import { Dialog } from "@base-ui/react/dialog";
import { Link } from "@tanstack/react-router";
import { cn } from "cn";

import { Text } from "../../../shared/components/text";

const menuButtonClassName =
  "fixed top-5 right-5 grid size-12 place-items-center sm:top-8 sm:right-8";

export function Menu() {
  return (
    <Dialog.Root>
      <Dialog.Trigger
        className={cn(menuButtonClassName, "group z-50 cursor-pointer")}
        aria-label="Открыть или закрыть меню"
      >
        <span className="relative block size-7" aria-hidden="true">
          <span className="absolute top-[calc(50%-4px)] left-0 h-0.5 w-full -translate-y-1/2 bg-current transition-[top,transform] duration-300 ease-in-out group-data-popup-open:top-1/2 group-data-popup-open:transform-[translateY(-50%)_rotate(45deg)]" />
          <span className="absolute top-[calc(50%+4px)] left-0 h-0.5 w-full -translate-y-1/2 bg-current transition-[top,transform] duration-300 ease-in-out group-data-popup-open:top-1/2 group-data-popup-open:transform-[translateY(-50%)_rotate(-45deg)]" />
        </span>
      </Dialog.Trigger>

      <Dialog.Portal>
        <Dialog.Popup className="fixed inset-0 z-40 flex min-h-dvh items-center justify-center bg-neutral-800 transition-opacity duration-200 ease-out data-ending-style:opacity-0 data-starting-style:opacity-0 motion-reduce:transition-none">
          <Dialog.Title className="sr-only">Меню</Dialog.Title>
          <nav className="flex flex-col items-center gap-4">
            <Link
              to="/login"
              className="transition-opacity hover:opacity-60"
            >
              <Text as="span" size="lg">
                Вход
              </Text>
            </Link>
            <Link
              to="/register"
              className="transition-opacity hover:opacity-60"
            >
              <Text as="span" size="lg">
                Регистрация
              </Text>
            </Link>
          </nav>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
