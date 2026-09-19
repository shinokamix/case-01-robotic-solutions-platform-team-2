import { motion } from "motion/react";
import type { ReactNode } from "react";

import { Text } from "../../../shared/components/text";

const fieldVariants = {
  idle: { x: 0 },
  invalid: { x: [0, -3, 3, -2, 0] },
};

type AuthFieldMotionProps = {
  invalid: boolean;
  className?: string;
  children: ReactNode;
};

export function AuthFieldMotion({
  invalid,
  className,
  children,
}: AuthFieldMotionProps) {
  return (
    <motion.div
      className={className}
      initial={false}
      animate={invalid ? "invalid" : "idle"}
      variants={fieldVariants}
      transition={{ duration: 0.24, ease: "easeOut" }}
    >
      {children}
    </motion.div>
  );
}

type AuthFieldErrorProps = {
  className: string;
  children: ReactNode;
};

export function AuthFieldError({ className, children }: AuthFieldErrorProps) {
  return (
    <motion.div
      className={className}
      role="alert"
      initial={{ opacity: 0, y: -3 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.18, ease: "easeOut" }}
    >
      <Text size="sm">{children}</Text>
    </motion.div>
  );
}
