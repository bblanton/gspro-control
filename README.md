# GSPro Control API

Small Go service exposing HTTP endpoints to trigger keyboard shortcuts on the host for GSPro actions.

## Endpoints

- `GET /health` → `{ "status": "ok" }`
- `GET /actions` → JSON map of actions to key combos
- `GET /command?action=<name>` → triggers combo mapped to `<name>`

Example:

```
GET http://localhost:8080/command?action=mulligan

200 OK
{
  "action": "mulligan",
  "result": "success"
}
```

## Actions

By default, the service ships with a built-in `mulligan` mapping to `["ctrl", "m"]`.

At startup, it attempts to read `actions.json` from the working directory and merge any overrides/additional actions.

See `actions.example.json` for format:

```json
{
  "mulligan": ["ctrl", "m"],
  "gimme": ["ctrl", "g"],
  "concede": ["ctrl", "q"]
}
```

Supported keys are those supported by `robotgo`. Modifiers like `ctrl`, `shift`, `alt`, `cmd` should be lowercase.

## Run

Prereqs: Go 1.21+

```
go run ./cmd/gspro-control
```

Change port with `PORT=3000`.

## macOS Accessibility Permissions

On first key press, macOS will require Accessibility permission for the terminal/app:

1. System Settings → Privacy & Security → Accessibility
2. Enable the binary (Terminal, Cursor, or the built app) running this server

## Windows Notes

- Key presses are sent using .NET Windows Forms `SendKeys` via PowerShell.
- Ensure the GSPro window has focus; SendKeys targets the active window.
- Modifiers map as: `ctrl` → `^`, `alt` → `%`, `shift` → `+`. The Windows key is not supported.
- Bracket keys like `[` and `]` are sent literally (no special escaping needed). Some symbols like `+`, `%`, `^`, `~`, `(`, `)`, `{`, `}` are escaped automatically.

## CORS

Open CORS (`*`) is enabled by default for quick testing from a browser-based UI. Restrict in production as needed.


