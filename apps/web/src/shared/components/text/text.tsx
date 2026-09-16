import { cn } from "cn";
import type {
  ComponentPropsWithoutRef,
  CSSProperties,
  ElementType,
} from "react";

const textSizes = {
  sm: [
    "text-xs! leading-4! tracking-normal!",
    "lg:text-sm! lg:leading-5! lg:tracking-normal!",
  ],
  md: [
    "text-sm! leading-5! tracking-normal!",
    "lg:text-base! lg:leading-6! lg:tracking-normal!",
  ],
  lg: [
    "text-base! leading-6! tracking-normal!",
    "lg:text-lg! lg:leading-7! lg:tracking-normal!",
  ],
} as const;

type TextSize = keyof typeof textSizes;
type TextStyle = Omit<
  CSSProperties,
  | "font"
  | "fontFamily"
  | "fontSize"
  | "fontWeight"
  | "letterSpacing"
  | "lineHeight"
>;

export type TextProps<T extends ElementType = "p"> = {
  as?: T;
  size?: TextSize;
  bold?: boolean;
  className?: string;
  style?: TextStyle;
} & Omit<ComponentPropsWithoutRef<T>, "as" | "className" | "style">;

export function Text<T extends ElementType = "p">({
  as,
  size = "md",
  bold = false,
  className,
  ...props
}: TextProps<T>) {
  const Component = as ?? "p";

  return (
    <Component
      className={cn(
        className,
        "font-sans!",
        textSizes[size][0],
        textSizes[size][1],
        bold ? "font-semibold!" : "font-normal!",
      )}
      {...props}
    />
  );
}
