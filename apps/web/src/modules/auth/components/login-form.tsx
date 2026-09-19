import { Button } from "@base-ui/react/button";
import { Field } from "@base-ui/react/field";
import { Input } from "@base-ui/react/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { Eye, EyeOff } from "lucide-react";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";

import { useLogin } from "../../../shared/api";
import type { ErrorResponse } from "../../../shared/api/schemas";
import { Text } from "../../../shared/components/text";
import type { ApiError } from "../../../shared/fetcher";
import { loginSchema, type LoginFormValues } from "../schemas/auth";
import { AuthFieldError, AuthFieldMotion } from "./auth-field-motion";

const inputClassName =
  "mt-2 h-11 w-full rounded-lg border border-white/15 bg-white/5 px-3 text-sm text-white outline-none transition-[border-color,box-shadow,background-color] placeholder:text-white/30 hover:bg-white/[0.07] focus:border-white/45 focus:ring-3 focus:ring-white/8 aria-invalid:border-red-400/80 aria-invalid:focus:ring-red-400/10";
const labelClassName = "text-white/90";
const errorClassName = "mt-1.5 text-red-300";

export function LoginForm() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema, undefined, { mode: "sync" }),
    mode: "onSubmit",
    reValidateMode: "onChange",
    defaultValues: { email: "", password: "" },
  });
  const login = useLogin<ApiError<ErrorResponse>>({
    mutation: {
      onError: (error) => {
        form.setError("root", {
          type: "server",
          message:
            error.data.message ||
            "Не удалось выполнить вход. Проверьте email и пароль.",
        });
      },
    },
  });

  const onSubmit = form.handleSubmit(async (values) => {
    await login.mutateAsync({
      data: { ...values, email: values.email.trim() },
    });
    await navigate({ to: "/" });
  });

  return (
    <form
      className="space-y-5"
      onSubmit={(event) => void onSubmit(event)}
      noValidate
      aria-busy={login.isPending}
    >
      <Controller
        name="email"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field.Root
            name={field.name}
            invalid={fieldState.invalid}
            touched={fieldState.isTouched}
            dirty={fieldState.isDirty}
          >
            <Field.Label className={labelClassName}>
              <Text as="span" size="sm" bold>
                Email
              </Text>
            </Field.Label>
            <AuthFieldMotion invalid={fieldState.invalid}>
              <Field.Control
                render={<Input />}
                {...field}
                className={inputClassName}
                type="email"
                autoComplete="email"
                inputMode="email"
                autoCapitalize="none"
                spellCheck={false}
                placeholder="mail@example.com"
                maxLength={254}
                aria-invalid={fieldState.invalid}
              />
            </AuthFieldMotion>
            {fieldState.error?.message && (
              <AuthFieldError className={errorClassName}>
                {fieldState.error.message}
              </AuthFieldError>
            )}
          </Field.Root>
        )}
      />

      <Controller
        name="password"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field.Root
            name={field.name}
            invalid={fieldState.invalid}
            touched={fieldState.isTouched}
            dirty={fieldState.isDirty}
          >
            <Field.Label className={labelClassName}>
              <Text as="span" size="sm" bold>
                Пароль
              </Text>
            </Field.Label>
            <AuthFieldMotion
              className="relative"
              invalid={fieldState.invalid}
            >
              <Field.Control
                render={<Input />}
                {...field}
                className={`${inputClassName} pr-12`}
                type={showPassword ? "text" : "password"}
                autoComplete="current-password"
                maxLength={128}
                aria-invalid={fieldState.invalid}
              />
              <button
                type="button"
                className="absolute bottom-0 right-1 grid size-11 place-items-center rounded-md text-white/45 outline-none transition-colors hover:text-white focus-visible:ring-2 focus-visible:ring-white/35"
                onClick={() => setShowPassword((visible) => !visible)}
                aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"}
                title={showPassword ? "Скрыть пароль" : "Показать пароль"}
              >
                {showPassword ? (
                  <EyeOff aria-hidden="true" size={18} strokeWidth={1.8} />
                ) : (
                  <Eye aria-hidden="true" size={18} strokeWidth={1.8} />
                )}
              </button>
            </AuthFieldMotion>
            {fieldState.error?.message && (
              <AuthFieldError className={errorClassName}>
                {fieldState.error.message}
              </AuthFieldError>
            )}
          </Field.Root>
        )}
      />

      {form.formState.errors.root?.message && (
        <AuthFieldError className="rounded-lg border border-red-400/20 bg-red-400/8 px-3 py-2.5 text-red-200">
          {form.formState.errors.root.message}
        </AuthFieldError>
      )}

      <Button
        type="submit"
        disabled={login.isPending}
        className="h-11 w-full rounded-lg bg-white px-4 text-black outline-none transition-[background-color,transform] hover:bg-white/90 active:translate-y-px focus-visible:ring-3 focus-visible:ring-white/25 disabled:cursor-wait disabled:opacity-50"
      >
        <Text as="span" size="sm" bold>
          {login.isPending ? "Входим..." : "Войти"}
        </Text>
      </Button>

      <Text size="sm" className="text-center text-white/50">
        Нет аккаунта?{" "}
        <Link
          to="/register"
          className="rounded-sm font-medium text-white outline-none underline decoration-white/35 underline-offset-4 hover:decoration-white focus-visible:ring-2 focus-visible:ring-white/35"
        >
          Зарегистрироваться
        </Link>
      </Text>
    </form>
  );
}
