import { Button } from "@base-ui/react/button";
import { Field } from "@base-ui/react/field";
import { Input } from "@base-ui/react/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { Eye, EyeOff } from "lucide-react";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";

import { useRegister } from "../../../shared/api";
import type { ErrorResponse } from "../../../shared/api/schemas";
import { Text } from "../../../shared/components/text";
import type { ApiError } from "../../../shared/fetcher";
import { registerSchema, type RegisterFormValues } from "../schemas/auth";
import { AuthFieldError, AuthFieldMotion } from "./auth-field-motion";

const inputClassName =
  "mt-2 h-11 w-full rounded-lg border border-white/15 bg-white/5 px-3 text-sm text-white outline-none transition-[border-color,box-shadow,background-color] placeholder:text-white/30 hover:bg-white/[0.07] focus:border-white/45 focus:ring-3 focus:ring-white/8 aria-invalid:border-red-400/80 aria-invalid:focus:ring-red-400/10";
const labelClassName = "text-white/90";
const errorClassName = "mt-1.5 text-red-300";

export function RegisterForm() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema, undefined, { mode: "sync" }),
    mode: "onSubmit",
    reValidateMode: "onChange",
    defaultValues: { firstName: "", lastName: "", email: "", password: "" },
  });
  const register = useRegister<ApiError<ErrorResponse>>({
    mutation: {
      onError: (error) => {
        const fields: Record<string, keyof RegisterFormValues | undefined> = {
          invalid_first_name: "firstName",
          invalid_last_name: "lastName",
          invalid_email: "email",
          email_already_registered: "email",
          invalid_password: "password",
        };
        const field = fields[error.data.code];

        form.setError(field ?? "root", {
          type: "server",
          message:
            error.data.message ||
            "Не удалось выполнить регистрацию. Попробуйте еще раз.",
        });
      },
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    register.mutate(
      { data: values },
      { onSuccess: () => void navigate({ to: "/" }) },
    );
  });

  return (
    <form
      className="space-y-5"
      onSubmit={(event) => void onSubmit(event)}
      noValidate
      aria-busy={register.isPending}
    >
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 sm:gap-3">
        <Controller
          name="firstName"
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
                  Имя
                </Text>
              </Field.Label>
              <AuthFieldMotion invalid={fieldState.invalid}>
                <Field.Control
                  render={<Input />}
                  {...field}
                  className={inputClassName}
                  autoComplete="given-name"
                  autoCapitalize="words"
                  maxLength={100}
                  aria-invalid={fieldState.invalid}
                />
              </AuthFieldMotion>
              {fieldState.error?.message && (
                <Field.Error match={true}>
                  <AuthFieldError className={errorClassName}>
                    {fieldState.error.message}
                  </AuthFieldError>
                </Field.Error>
              )}
            </Field.Root>
          )}
        />

        <Controller
          name="lastName"
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
                  Фамилия
                </Text>
              </Field.Label>
              <AuthFieldMotion invalid={fieldState.invalid}>
                <Field.Control
                  render={<Input />}
                  {...field}
                  className={inputClassName}
                  autoComplete="family-name"
                  autoCapitalize="words"
                  maxLength={100}
                  aria-invalid={fieldState.invalid}
                />
              </AuthFieldMotion>
              {fieldState.error?.message && (
                <Field.Error match={true}>
                  <AuthFieldError className={errorClassName}>
                    {fieldState.error.message}
                  </AuthFieldError>
                </Field.Error>
              )}
            </Field.Root>
          )}
        />
      </div>

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
              <Field.Error match={true}>
                <AuthFieldError className={errorClassName}>
                  {fieldState.error.message}
                </AuthFieldError>
              </Field.Error>
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
            <AuthFieldMotion className="relative" invalid={fieldState.invalid}>
              <Field.Control
                render={<Input />}
                {...field}
                className={`${inputClassName} pr-12`}
                type={showPassword ? "text" : "password"}
                autoComplete="new-password"
                maxLength={128}
                aria-invalid={fieldState.invalid}
              />
              <button
                type="button"
                className="absolute right-1 bottom-0 grid size-11 place-items-center rounded-md text-white/45 transition-colors outline-none hover:text-white focus-visible:ring-2 focus-visible:ring-white/35"
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
            {fieldState.error?.message ? (
              <Field.Error match={true}>
                <AuthFieldError className={errorClassName}>
                  {fieldState.error.message}
                </AuthFieldError>
              </Field.Error>
            ) : (
              <Text size="sm" className="mt-2 text-white/45">
                Не менее 12 символов
              </Text>
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
        disabled={register.isPending}
        className="h-11 w-full rounded-lg bg-white px-4 text-black transition-[background-color,transform] outline-none hover:bg-white/90 focus-visible:ring-3 focus-visible:ring-white/25 active:translate-y-px disabled:cursor-wait disabled:opacity-50"
      >
        <Text as="span" size="sm" bold>
          {register.isPending ? "Создаем аккаунт..." : "Создать аккаунт"}
        </Text>
      </Button>

      <Text size="sm" className="text-center text-white/50">
        Уже есть аккаунт?{" "}
        <Link
          to="/login"
          className="rounded-sm font-medium text-white underline decoration-white/35 underline-offset-4 outline-none hover:decoration-white focus-visible:ring-2 focus-visible:ring-white/35"
        >
          Войти
        </Link>
      </Text>
    </form>
  );
}
