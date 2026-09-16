import { cn } from "cn";
import type {
  ComponentPropsWithoutRef,
  CSSProperties,
  ElementType,
} from "react";

const headingSizes = {
  sm: [
    "text-2xl! leading-8! tracking-normal!",
    "lg:text-3xl! lg:leading-9! lg:tracking-tight!",
  ],
  md: [
    "text-4xl! leading-10! tracking-tight!",
    "lg:text-5xl! lg:leading-[3.25rem]! lg:tracking-tight!",
  ],
  lg: [
    "text-5xl! leading-none! tracking-tight!",
    "lg:text-8xl! lg:leading-none! lg:tracking-tight!",
  ],
} as const;

type HeadingSize = keyof typeof headingSizes;
type HeadingStyle = Omit<
  CSSProperties,
  | "font"
  | "fontFamily"
  | "fontSize"
  | "fontWeight"
  | "letterSpacing"
  | "lineHeight"
>;

export type HeadingProps<T extends ElementType = "h2"> = {
  as?: T;
  size?: HeadingSize;
  bold?: boolean;
  className?: string;
  style?: HeadingStyle;
} & Omit<ComponentPropsWithoutRef<T>, "as" | "className" | "style">;

export function Heading<T extends ElementType = "h2">({
  as,
  size = "md",
  bold = true,
  className,
  ...props
}: HeadingProps<T>) {
  const Component = as ?? "h2";

  return (
    <Component
      className={cn(
        className,
        "font-display!",
        headingSizes[size][0],
        headingSizes[size][1],
        bold ? "font-semibold!" : "font-normal!",
      )}
      {...props}
    />
  );
}
