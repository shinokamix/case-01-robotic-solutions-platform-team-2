const apiUrl = import.meta.env.VITE_API_URL ?? "";

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
  const response = await fetch(`${apiUrl}${url}`, {
    ...options,
    credentials: "include",
  });

  const data: unknown = await response.json().catch(() => undefined);

  if (!response.ok) {
    throw new ApiError(response.status, data);
  }

  // The OpenAPI-generated endpoint type defines the expected response body.
  // eslint-disable-next-line typescript/no-unsafe-type-assertion
  return data as T;
}
