# How snapshots work

A snapshot is a point-in-time copy of a single Noita run. Take one before you
do something risky, and `restore` rolls your save back to it. Only the files
that make up the *current run* are captured, so the unlocks and stats that live
on after death come through the restore untouched. Nothing else moves.

## What gets captured

When you die, Noita deletes the files that describe the run. Anything that
survives death is lifetime progression and is **not** part of a snapshot. The
captured files are:

| Path | What it is |
| --- | --- |
| `player.xml` | The player entity (inventory, perks, spells). |
| `session_numbers.salakieli` | Run session counters. |
| `world_state.xml` | World state (seed, progress). |
| `world/` | Generated biome chunks. |
| `stats/_streaks.salakieli` | Death-streak counter for the current run. |

Everything else in the folder is left alone. `persistent/` (unlocked spells,
orbs, bones), the rest of `stats/`, `mod_config.xml`, `mod_settings.bin`, and
`steam_autocloud.vdf` all stay where they are.

## Where snapshots are stored

Snapshots live next to huijata's config file, in a `snapshots/` folder:

- Windows: `%AppData%\huijata\snapshots\`
- Linux: `~/.config/huijata/snapshots/`
- dev (when `HUIJATA_CONFIG_PATH="."` is set in `.env`): `./snapshots/`

Each snapshot is one folder, named by a generated id. Inside the folder sits a
`meta.json`, a small file that remembers what the snapshot is:

```json
{
  "id": "20260805-100000-4f3a9c",
  "description": "safe point before the tower",
  "created_at": "2026-08-05T10:00:00Z"
}
```

The id is the snapshot's identity. The description is whatever you typed when
saving, and it does not have to be unique. There is no overwrite prompt, so
saving twice just gives you a second folder.

## The commands

- `huijata save <description>` copies the run-state files into a new snapshot
  folder. Description must be non-empty and contain no path separators.
- `huijata restore` lists snapshots newest-first, lets you pick one, asks for
  confirmation, then rolls the run back.
- `huijata clear` asks for confirmation, then deletes every snapshot. The save
  file itself is never touched.

## How restore works

Restore **mirrors** the run-state files. For each captured path it copies the
snapshot's version into the save, and when the snapshot has no such path it
removes the path from the save instead. A snapshot taken after death wipes a
live run back to the dead state, and a live snapshot brings a dead run back to
life.

Any progress you made *after* the snapshot is gone. Lifetime progression
(unlocks, stats, orbs) is kept, because those files are never touched.
