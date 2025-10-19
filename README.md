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

### Available Actions

These are the available action names, their default key bindings, and descriptions. You can override any binding in `actions.json`.

| Action | Default Key | Description |
| --- | --- | --- |
| `toggle_ui` | `esc` | Shows User Interface |
| `clear_tracers` | `f1` | Clear Tracer(s) |
| `zoom_to_green` | `f3` | Zoom to aim point, then zoom to green |
| `free_float_camera` | `f5` | Free float camera; then use arrow keys to fly |
| `console_short` | `f8` | Console (short) / error log |
| `console_tall` | `f9` | Console (tall) |
| `full_screen` | `f11` | Full Screen |
| `cam_go_to_ball` | `5` | Camera go to ball |
| `cam_fly_to_ball` | `6` | Camera fly to ball |
| `goto_ball` | `8` | Goto ball |
| `prev_hole` | `9` | Previous hole |
| `next_hole` | `0` | Next hole |
| `mulligan` | `ctrl + m` | Mulligan |
| `pin_indicator` | `p` | Pin indicator |
| `fly_over` | `o` | Flyover of the hole |
| `line_of_sight_toggle` | `b` | Toggle objects in line of sight invisible |
| `putt_toggle` | `u` | Putt toggle |
| `heat_map` | `y` | Heat map on the green |
| `vertical_dots` | `d` | Vertical dots on the green |
| `shot_cam` | `j` | Shot cam |
| `lighting_menu` | `l` | Lighting/time-of-day adjustments |
| `fps` | `f` | Show FPS (frames per second) |
| `hud_toggle` | `h` | HUD toggle |
| `green_grid` | `g` | Green grid |
| `prev_club` | `k` | Previous club |
| `next_club` | `i` | Next club |
| `toggle_grass` | `z` | Hide/show 3D grass |
| `toggle_water` | `r` | Hide/show reflective water |
| `scorecard` | `t` | Scorecard display |
| `shadow_decrease` | `,` | Decrease shadow quality |
| `shadow_increase` | `.` | Increase shadow quality |
| `mute` | `-` | Mute game audio |
| `sound_on` | `shift + =` | Turn game audio on (often `+` on US keyboards) |
| `fast_forward_hold` | `space` | Hold to fast-forward ball roll |

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


