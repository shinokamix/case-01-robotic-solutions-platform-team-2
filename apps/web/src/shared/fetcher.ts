import type { ErrorResponse } from "./api/schemas";

const apiUrl = import.meta.env.VITE_API_URL ?? "/api";

function isErrorResponse(data: unknown): data is ErrorResponse {
  return (
    typeof data === "object" &&
    data !== null &&
    "code" in data &&
    typeof data.code === "string" &&
    "message" in data &&
    typeof data.message === "string"
  );
}

export class ApiError<T = unknown> extends Error {
  readonly status: number;
  readonly data: T;

  constructor(status: number, data: T) {
    super(`API request failed with status ${status}`);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
  }
}

export async function fetcher<T>(
  url: string,
  options?: RequestInit,
): Promise<T> {
  try {
    const response = await fetch(`${apiUrl}${url}`, {
      ...options,
      credentials: "include",
    });

    const data: unknown = await response.json().catch(() => undefined);

    if (!response.ok) {
      throw new ApiError<ErrorResponse>(
        response.status,
        isErrorResponse(data)
          ? data
          : {
              code: "invalid_error_response",
              message: "Не удалось выполнить запрос.",
            },
      );
    }

    // The OpenAPI-generated endpoint type defines the expected response body.
    // eslint-disable-next-line typescript/no-unsafe-type-assertion
    return data as T;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    throw new ApiError(0, {
      code: "network_error",
      message: "Не удалось подключиться к серверу.",
    });
  }
}
