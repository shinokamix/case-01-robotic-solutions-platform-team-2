import { createFileRoute } from "@tanstack/react-router";

import { AuthLayout, RegisterForm } from "../modules/auth";

function RegisterPage() {
  return (
    <AuthLayout
      title="Создайте аккаунт"
      description="Сохраняйте проекты и возвращайтесь к расчетам в любое время."
    >
      <RegisterForm />
    </AuthLayout>
  );
}

export const Route = createFileRoute("/register")({
  component: RegisterPage,
});
