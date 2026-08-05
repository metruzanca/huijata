# Huijata

Toolkit for inspecting and save editing in Noita (single-player game). Saves are
plain files on disk, so snapshot/restore and inventory editing are done by
reading/writing those files. Only affects the user's own copy of the game.

## What Noita's save looks like

Noita stores its data in a folder named `Nolla_Games_Noita`:

- Windows: `%LOCALAPPDATA%\Low\Nolla_Games_Noita\`
- Linux (Proton): `~/.local/share/Steam/steamapps/compatdata/881100/pfx/drive_c/users/steamuser/AppData/LocalLow/Nolla_Games_Noita/`
- macOS (Proton): `~/Library/Application Support/Steam/steamapps/compatdata/881100/pfx/drive_c/users/steamuser/AppData/LocalLow/Nolla_Games_Noita/`

Noita has no multi-save support, so the only save slot is `save00/`.

### What survives death vs. what gets wiped

When you die, run-specific files are deleted (`player.xml`,
`world_state.xml`, `session_numbers.salakieli`, the `world/` area chunks).
Meta-progression persists:

- `persistent/flags/` — one file per unlocked spell, named after the action
  (e.g. `action_accelerating_shot`). Presence of the file == unlocked. Note:
  these files are not actually empty; the game writes `why are you looking
  here..` into them.
- `persistent/orbs_new/` — numbered files marking orb unlocks.
- `persistent/bones_new/` — killed-by info / Kummitus wands (`item*.xml`).
- `persistent/sessions/` and `stats/` — long-term run stats.
- `world/` — biome generation data, persists across runs.

`save_files_examples/` (gitignored) holds three snapshots for reference:
`0_before_new_game`, `1_new_game`, `2_after_death`.

## Config

`internal/config/config.go` persists settings to `config.toml`:

- Config file lives in `os.UserConfigDir()/huijata/config.toml` unless
  `HUIJATA_CONFIG_PATH` is set, in which case it's
  `$HUIJATA_CONFIG_PATH/config.toml`.
- `game_path` stores the **`Nolla_Games_Noita` folder** (not `save00`).
  `Config.SavePath()` returns the actual save slot by appending `save00`.

`huijata init` interactively finds/validates the game folder (huh forms) and
writes the config.

## Development

- `.env` and `config.toml` are gitignored; in dev `.env` sets
  `HUIJATA_CONFIG_PATH="."` so the config lives in the repo root.
  `main.go` loads `.env` with `godotenv` before running the CLI.
- Stack: Go + `spf13/cobra` (CLI), `charmbracelet/huh` (interactive forms),
  `BurntSushi/toml` (config).
- Layout: `cmd/` for cobra commands, `internal/config/` for the config
  package, `internal/snapshots/` for snapshot file logic.
- Verify with `go build ./...`, `go test ./...`, `go vet ./...`.

## Planned

- Save snapshotting and restoring.
- Inventory (wand/card) editing.
