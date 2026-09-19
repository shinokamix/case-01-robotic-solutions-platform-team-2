import tailwindcss from "@tailwindcss/vite";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig, lazyPlugins } from "vite-plus";

// https://vite.dev/config/
export default defineConfig({
  fmt: {
    ignorePatterns: ["src/routeTree.gen.ts", "src/shared/api/**"],
    printWidth: 80,
    sortImports: true,
    sortTailwindcss: true,
  },
  lint: {
    ignorePatterns: ["src/routeTree.gen.ts", "src/shared/api/**"],
    categories: {
      correctness: "error",
      perf: "warn",
      suspicious: "error",
    },
    plugins: [
      "import",
      "jsx-a11y",
      "oxc",
      "promise",
      "react",
      "typescript",
      "unicorn",
    ],
    rules: {
      curly: "error",
      eqeqeq: ["error", "always", { null: "ignore" }],
      "import/first": "error",
      "import/newline-after-import": "error",
      "import/no-cycle": "error",
      "import/no-duplicates": "error",
      "import/no-self-import": "error",
      "import/no-unassigned-import": "off",
      "jsx-a11y/alt-text": "error",
      "jsx-a11y/anchor-is-valid": "error",
      "jsx-a11y/click-events-have-key-events": "error",
      "jsx-a11y/iframe-has-title": "error",
      "jsx-a11y/label-has-associated-control": "error",
      "jsx-a11y/no-autofocus": "warn",
      "jsx-a11y/no-static-element-interactions": "error",
      "no-alert": "error",
      "no-console": ["warn", { allow: ["warn", "error"] }],
      "no-var": "error",
      "object-shorthand": "error",
      "prefer-const": "error",
      "promise/no-return-wrap": "error",
      "promise/param-names": "error",
      "react/jsx-no-useless-fragment": "error",
      "react/no-array-index-key": "warn",
      "react/no-danger": "error",
      "react/only-export-components": [
        "warn",
        {
          allowConstantExport: true,
        },
      ],
      "react/react-in-jsx-scope": "off",
      "react/rules-of-hooks": "error",
      "react/self-closing-comp": "error",
      "stylistic/padding-line-between-statements": [
        "error",
        {
          blankLine: "always",
          prev: ["const", "let", "var"],
          next: "*",
        },
        {
          blankLine: "any",
          prev: ["const", "let", "var"],
          next: ["const", "let", "var"],
        },
        {
          blankLine: "always",
          prev: "*",
          next: ["return", "throw"],
        },
      ],
      "typescript/consistent-type-imports": "error",
      "typescript/no-explicit-any": "warn",
      "typescript/no-floating-promises": "error",
      "typescript/no-misused-promises": "error",
      "typescript/no-non-null-assertion": "warn",
      "typescript/switch-exhaustiveness-check": "error",
      "vite-plus/prefer-vite-plus-imports": "error",
    },
    overrides: [
      {
        files: ["src/routes/**/*.tsx"],
        rules: {
          "react/only-export-components": "off",
        },
      },
    ],
    options: {
      reportUnusedDisableDirectives: "error",
      typeAware: true,
      typeCheck: true,
    },
    jsPlugins: [
      {
        name: "stylistic",
        specifier: "@stylistic/eslint-plugin",
      },
      {
        name: "vite-plus",
        specifier: "vite-plus/oxlint-plugin",
      },
    ],
  },
  server: {
    proxy: {
      "/api": {
        target: process.env.API_PROXY_TARGET ?? "http://localhost:8080",
      },
    },
  },
  plugins: lazyPlugins(() => [TanStackRouterVite(), react(), tailwindcss()]),
});
