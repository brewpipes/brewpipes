# Browser Interaction & Dev Server Lifecycle

> **OpenCode frontmatter:** To enable Playwright MCP tools for an agent, add `playwright: true` to the agent's `tools:` section in the YAML frontmatter.

## Playwright MCP Browser Tools

You have access to browser automation via the Playwright MCP server. This gives you tools to interact with the running BrewPipes web UI directly:

### Key tools

- **`browser_navigate`** — Navigate to a URL (e.g., `http://localhost:3000`)
- **`browser_snapshot`** — Capture the page's accessibility tree (how you "see" the page structure)
- **`browser_click`** — Click an element by its accessibility snapshot reference
- **`browser_type`** — Type text into an input field
- **`browser_fill_form`** — Fill multiple form fields at once
- **`browser_select_option`** — Select a dropdown option
- **`browser_take_screenshot`** — Take a visual screenshot
- **`browser_console_messages`** — Read browser console output (useful for debugging)
- **`browser_tabs`** — List, create, close, or switch browser tabs
- **`browser_navigate_back`** — Go back in browser history
- **`browser_press_key`** — Press a keyboard key
- **`browser_wait_for`** — Wait for text to appear/disappear or a timeout

### Usage guidelines

- Use `browser_snapshot` (accessibility tree) for understanding page structure and finding interactive elements. This is more reliable than screenshots.
- Use `browser_take_screenshot` when you need to verify visual layout, styling, or responsive design.
- Element references (`ref`) come from the accessibility snapshot — take a snapshot first, then use the refs to interact.
- The browser launches headed (visible window) by default. It persists between tool calls within a session.

## Dev Server Lifecycle

BrewPipes has two servers that need to run for full local development:

### Backend (Go monolith on :8080)

```bash
# Start (runs in background)
DATABASE_URL="postgres://brewpipes:brewpipes@localhost:5432/brewpipes?sslmode=disable" \
BREWPIPES_SECRET_KEY="dummy" \
go run ./cmd/monolith &

# Or use make (foreground)
make run-server
```

Requires Postgres to be running first: `make start-postgres`

### Frontend (Vite dev server on :3000)

```bash
# Start (runs in background, from service/www/)
pnpm dev --force &

# Or use make (foreground)
make run-web
```

The Vite dev server proxies `/api` requests to `:8080` (the Go backend).

### CRITICAL: Always clean up

**You MUST stop any servers you start before completing your task.** The user needs these ports free for their own local development.

To stop background processes:
```bash
# Find and kill the Vite dev server
lsof -ti:3000 | xargs kill 2>/dev/null

# Find and kill the Go backend
lsof -ti:8080 | xargs kill 2>/dev/null
```

Always run these cleanup commands when you're done, even if your task failed or errored. Wrap your work in a try/finally pattern mentally: start servers → do work → **always** stop servers.

### Startup sequence

If you need both servers running:

1. Ensure Postgres is running: `make start-postgres`
2. Start the Go backend in background
3. Wait for it to be ready (check with `curl -s http://localhost:8080/api/batches` or similar)
4. Start the Vite dev server in background
5. Wait for it to be ready (check with `curl -s http://localhost:3000`)
6. Do your browser-based work
7. **Stop both servers when done**

### When to spin up servers

- **Do** spin up servers when you need to verify UI rendering, test user flows, or debug browser issues.
- **Do NOT** spin up servers for code-only tasks (writing code, running unit tests, type-checking). Unit tests and type-checks don't need running servers.
