# Huijata

Toolkit for inspecting and save editing in Noita (single-player game). Saves are
plain files on disk, so snapshot/restore and inventory editing are done by
reading/writing those files. Only affects the user's own copy of the game.

## What Noita's save looks like

Noita stores its data in a folder named `Nolla_Games_Noita`:

- Windows: `%USERPROFILE%\AppData\LocalLow\Nolla_Games_Noita\` (LocalLow is
  its own top-level AppData folder, not `Local\Low`)
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
- `persistent/sessions/` and `stats/` — long-term run stats. Exception:
  `stats/_streaks.salakieli` is run-specific (created on new game, deleted on
  death), so it belongs to the run state, not lifetime progression.
- `world/` — biome generation data, persists across runs.

`save_files_examples/` (gitignored) holds three snapshots for reference:
`0_before_new_game`, `1_new_game`, `2_after_death`.

## `.salakieli` encryption

The `.salakieli` files are AES-128 in CTR mode. The key and IV are the first
16 bytes of two fixed passphrase strings; which pair a file uses depends on
the **kind of data it holds**, not its name:

| File | Key | IV |
| --- | --- | --- |
| `stats/_stats.salakieli`, `stats/_streaks.salakieli` | `SecretsOfTheAllSeeing` | `ThreeEyesAreWatchingYou` |
| `session_numbers.salakieli`, `persistent/magic_numbers.salakieli` | `KnowledgeIsTheHighestOfTheHighest` | `WhoWouldntGiveEverythingForTrueKnowledge` |
| `player.salakieli` | `WeSeeATrueSeekerOfKnowledge` | `YouAreSoCloseToBeingEnlightened` |
| `world_state.salakieli` | `TheTruthIsThatThereIsNothing` | `MoreValuableThanKnowledge` |

`internal/salakieli.DecryptFile` handles these (CTR makes encrypt and decrypt
the same operation). CTR with a fixed key/IV means two files sharing a
passphrase also share their keystream, which is how the `_stats`/`_streaks`
and `session_numbers`/`magic_numbers` pairs were identified.

`data.wak` is **not** encrypted in current Noita builds (older builds used an
LCG-derived key), so `internal/noitadata` reads it directly.

## Stats

`huijata stats` decrypts `stats/_stats.salakieli` and prints the player's
lifetime, best-run and last-session stats plus top tracked kills.

The decrypted file is `<GameStats>` XML: a `KEY_VALUE_STATS` block of
`<E key="..." value="..."/>` pairs (enemy names and `action_*` spell casts)
followed by `session`, `highest`, `global`, `prev_best` blocks whose stats
are XML attributes (`death_count`, `enemies_killed`, `gold`, `playtime_str`,
etc.). `internal/stats.Parse` turns it into structs; `stats.Load` does the
decrypt + parse from a save path. `internal/player` already handles
`player.xml` the same way.

## Config

`internal/config/config.go` persists settings to `config.toml`:

- Config file lives in `os.UserConfigDir()/huijata/config.toml` unless
  `HUIJATA_CONFIG_PATH` is set, in which case it's
  `$HUIJATA_CONFIG_PATH/config.toml`.
- `game_path` stores the **`Nolla_Games_Noita` folder** (not `save00`).
  `Config.SavePath()` returns the actual save slot by appending `save00`.

`huijata init` interactively finds/validates the game folder (huh forms) and
writes the config.

## Snapshots

`huijata save <description>` / `huijata restore` / `huijata clear` snapshot and
roll back the *current run*. See `docs/snapshots.md`.

- A snapshot captures exactly the run-state files: `player.xml`,
  `session_numbers.salakieli`, `world_state.xml`, `world/` (recursive), and
  `stats/_streaks.salakieli` (see `internal/snapshots.RunFiles`). Everything
  else (lifetime progression) is never touched.
- Snapshots are stored in `snapshots/` next to the config file (resolved by
  `config.SnapshotsDir()`; in dev that's `./snapshots`).
- One folder per snapshot, named by a generated id
  (`20060102-150405-<hex>` = timestamp prefix + random suffix). Inside is
  `meta.json` with `id`, `description`, `created_at`. The id is the identity;
  the description is free text and duplicates are allowed — no overwrite
  prompt.
- `restore` **mirrors** run files: it copies each run-state path from the
  snapshot and removes paths the snapshot lacks, so restoring a dead-state
  snapshot correctly wipes a live run. `restore` and `clear` confirm via huh
  forms.
- Save file safety: only `save00` files listed in `RunFiles` are ever
  written/deleted.

## Spell names

`wands show` prints the player's spells by their in-game English name. Spell
display names are not stored in the save: `player.xml` only holds internal
action ids (e.g. `LIGHT_BULLET`, `MINE`). They are resolved at runtime from
the installed game's own data by `internal/noitadata`:

- `data.wak` (at `<install>/data/data.wak`) is a packed archive whose directory
  (from byte 24) is `[u32 name_len][name][u32 off][u32 a]` per entry; an
  entry's contents span `[previous entry's off, this entry's off)` (the first
  starts at the header's data offset). Contents are plain text, zlib if the
  first byte is `0x78`.
- `data/scripts/gun/gun_actions.lua` (inside the wak) maps each action
  `id = "X"` to a translation key `name = "$action_y"`. This is the
  authoritative id→key mapping — do NOT guess it from the id (e.g.
  `LIGHTNING_BOLT`→`$action_lightning`, `DISC_BULLET_BIGGER`→
  `$action_omega_disc_bullet`). Some entries have a `"???"` placeholder
  (`FUNKY_SPELL`); those simply don't resolve.
- `data/translations/common.csv` (unpacked on disk) maps the key to the
  English name (`action_light_bullet` → "Spark bolt").

Verified id → display-name examples (do not assume the display name from the
id — it is frequently unrelated):

- `LIGHT_BULLET` → "Spark bolt" (the starter wand's spell)
- `MINE` → "Unstable crystal"
- `RUBBER_BALL` → "Bouncing burst"
- `BOMB` → "Bomb"
- `LIGHTNING_BOLT` → "Lightning bolt" (via `$action_lightning`)

`gun_actions.lua` has ~491 actions; the `id`/`name` pattern matches all but a
few (`FUNKY_SPELL`'s name is a literal `"???"` placeholder), and `common.csv`
holds 514 `action_*` keys — some keys resolve to names not present in
`common.csv` at all, which is why the fallback to the raw id exists.

`internal/noitadata.SpellNames(installDir)` returns `id -> English name` for
every resolvable spell; pass `""` to auto-detect the install dir (Steam
default locations + `libraryfolders.vdf`, or `HUIJATA_NOITA_INSTALL`). Spells
that don't resolve fall back to the raw id in `wands show`.

## Development

- `.env` and `config.toml` are gitignored; in dev `.env` sets
  `HUIJATA_CONFIG_PATH="."` so the config lives in the repo root.
  `main.go` loads `.env` with `godotenv` before running the CLI.
- Stack: Go + `spf13/cobra` (CLI), `charmbracelet/huh` (interactive forms),
  `BurntSushi/toml` (config).
- Layout: `cmd/` for cobra commands, `internal/config/` for the config
  package, `internal/snapshots/` for snapshot file logic,
  `internal/noitadata/` for game-data reading (wak, spell names),
  `internal/player/` for player.xml parsing, `internal/salakieli/` for
  `.salakieli` decryption, `internal/stats/` for stats XML parsing.
- Verify with `go build ./...`, `go test ./...`, `go vet ./...`.

## Planned

- Inventory (wand/card) editing.
