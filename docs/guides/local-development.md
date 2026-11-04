# Local development

This guide reflects the repository as it exists now. It is a web-only Next.js starter; the Go API, PostgreSQL service, voice adapter, and Compose configuration are planned and are not available yet.

## Prerequisites

- Node.js compatible with the repository's `package.json` and lockfile.
- pnpm 12.10.1, declared by the `packageManager` field.

## Start the web app

From the repository root:

```bash
pnpm install --frozen-lockfile
pnpm dev
```

Open <http://localhost:3000>. The current page is still the default Next.js starter screen.

## Available package scripts

| Command | Purpose |
| --- | --- |
| `pnpm dev` | Start the Next.js development server. |
| `pnpm build` | Build the web application. |
| `pnpm start` | Run the built application. |
| `pnpm lint` | Run ESLint. |

## Next.js documentation rule

Before changing application code, follow the repository's `AGENTS.md`: read the relevant guide under `node_modules/next/dist/docs/` for the installed Next.js version. The project has breaking changes relative to older Next.js conventions.

## Planned services

Do not use backend, database, voice, or Docker commands yet; corresponding services and configuration have not been implemented. Their setup instructions should be added here when those artifacts exist and can be verified.
