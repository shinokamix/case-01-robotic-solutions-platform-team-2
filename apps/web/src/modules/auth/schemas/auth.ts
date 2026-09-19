import { z } from "zod";

const nameSchema = (field: "имя" | "фамилия") =>
  z
    .string()
    .trim()
    .min(1, `Введите ${field}`)
    .max(100, `${field} не должно быть длиннее 100 символов`)
    .regex(/^[^\p{Cc}]*$/u, `${field} содержит недопустимые символы`);

const emailSchema = z
  .string()
  .trim()
  .min(1, "Введите email")
  .max(254, "Email не должен быть длиннее 254 символов")
  .pipe(z.email("Введите корректный email"));

const passwordSchema = z
  .string()
  .min(12, "Пароль должен содержать от 12 до 128 символов")
  .max(128, "Пароль должен содержать от 12 до 128 символов");

export const loginSchema = z.object({
  email: emailSchema,
  password: z.string().min(1, "Введите пароль"),
});

export const registerSchema = z.object({
  firstName: nameSchema("имя"),
  lastName: nameSchema("фамилия"),
  email: emailSchema,
  password: passwordSchema,
});

export type LoginFormValues = z.infer<typeof loginSchema>;
export type RegisterFormValues = z.infer<typeof registerSchema>;
