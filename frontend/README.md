# TravelLog frontend

React + Vite UI for TravelLog. API calls use the `/api` prefix; in dev, Vite proxies that to the Go backend.

## Quick commands

```bash
npm install
npm run dev          # http://localhost:5173
npm run build
npm test             # Vitest (unit / component tests)
npm run cypress:open # E2E (needs backend + npm run dev)
```

## Backend URL (proxy)

Default proxy target is `http://localhost:8080` (same as backend `SERVER_PORT`).

If the API runs on another port, add `frontend/.env.local`:

```bash
VITE_API_PROXY_TARGET=http://localhost:8081
```

Restart the dev server after changing this.

## Full-stack instructions

See the main project **[README.md](../README.md)** — section **11. Local development and testing** (ports, seed user, Cypress, fixing “address already in use” on 8080).
