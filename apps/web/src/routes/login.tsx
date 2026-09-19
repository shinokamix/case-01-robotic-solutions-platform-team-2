import { createFileRoute } from "@tanstack/react-router";

import { AuthLayout, LoginForm } from "../modules/auth";

function LoginPage() {
  return (
    <AuthLayout
      title="С возвращением"
      description="Войдите, чтобы продолжить работу с вашими проектами."
    >
      <LoginForm />
    </AuthLayout>
  );
}

export const Route = createFileRoute("/login")({
  component: LoginPage,
});
