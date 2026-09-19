import { defineConfig } from "orval";

export default defineConfig({
  api: {
    input: {
      target: "../../packages/contracts/openapi.yaml",
    },
    output: {
      mode: "tags-split",
      target: "src/shared/api/generated.ts",
      schemas: "src/shared/api/schemas",
      client: "react-query",
      httpClient: "fetch",
      namingConvention: "kebab-case",
      clean: true,
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        query: {
          version: 5,
        },
        mutator: {
          path: "src/shared/fetcher.ts",
          name: "fetcher",
        },
      },
    },
  },
});
