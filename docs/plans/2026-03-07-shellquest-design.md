# Shell Quest — Design Document

**Date:** 2026-03-07
**Status:** Approved

---

## Overview

Shell Quest is a kids' command-line learning game for macOS (and later Linux), targeting ages 3–10. It runs as a single Go binary with no external service dependencies — SQLite handles all persistence.

The game is themed as a **treasure hunt**: children navigate a pirate world by typing real shell commands into a sandboxed fake terminal, discovering clues and finding treasure.

---

## Architecture

### Technology Stack

- **Language:** Go (single binary)
- **TUI Framework:** Bubble Tea (Charm) with Lip Gloss for styling
- **Persistence:** SQLite via `modernc.org/sqlite` (pure Go, no CGo required)
- **Content:** Embedded JSON mission files via `go:embed`
- **DB Location:** `~/.shellquest/progress.db`

### Project Structure

```
kids-cli/
  cmd/shellquest/
    main.go                  ← entry point
  internal/
    tui/                     ← Bubble Tea models, layout, key handling
      model.go               ← root model: profile select, game, menus
      panes.go               ← split-pane layout (story left, shell right)
    shell/                   ← fake shell engine
      executor.go            ← command dispatch
      filesystem.go          ← virtual FS tree (in-memory)
      commands/              ← one file per command (ls.go, cd.go, etc.)
    world/                   ← mission/narrative engine
      loader.go              ← loads missions from embedded JSON
      events.go              ← command→narrative event mapping
      missions.go            ← mission state machine
    db/                      ← SQLite layer
      schema.go              ← CREATE TABLE statements, migrations
      players.go             ← player CRUD
      progress.go            ← mission progress read/write
      worldstate.go          ← virtual FS snapshot save/load
  content/                   ← embedded game content
    worlds/
      skull_island.json      ← World 1: missions + filesystem layout
    avatars/
      pirate.txt             ← ASCII art avatars
  docs/plans/
    2026-03-07-shellquest-design.md
```

---

## Split-Pane TUI Layout

```
┌─────────────────────────────┬──────────────────────────────────┐
│  🗺  TREASURE MAP            │  🏴‍☠️  SHELL QUEST               │
│                             │                                  │
│  [ASCII pirate map]         │  > You are in: /island/cave      │
│                             │  > A chest glimmers nearby...    │
│  📍 You are here            │                                  │
│  🏆 Treasure: ???           │  pirate@quest:/island/cave$ _    │
│                             │                                  │
│  CLUE: "Look inside the     │  [scrollable command output]     │
│  hidden folder beneath      │                                  │
│  the old oak"               │                                  │
└─────────────────────────────┴──────────────────────────────────┘
```

- Left pane: story narrative, current clue, ASCII map, inventory
- Right pane: fake shell prompt + scrollable output history
- All keyboard input routes to the shell prompt

---

## Fake Shell & Virtual Filesystem

### Virtual Filesystem

An in-memory tree of nodes, persisted to SQLite as a JSON snapshot on save/quit.

Each node has: `id`, `parent_id`, `name`, `type` (file/dir), `content`, `hidden` (for `.dotfiles`), `permissions` (for chmod lessons).

Worlds ship with a pre-built filesystem defined in JSON content files. On first run, the filesystem is instantiated into memory from the content JSON; on save, it is serialized back to SQLite.

### Commands by Tier

| Tier | Commands |
|------|----------|
| **Beginner** (ages 3–6) | `ls`, `cd`, `pwd`, `cat`, `echo`, `clear`, `help` |
| **Explorer** (ages 6–8) | + `mkdir`, `touch`, `cp`, `mv`, `rm`, `find` |
| **Master** (ages 8–10) | + `grep`, `chmod`, `man`, `history`, pipe `\|`, glob `*`, redirect `>` |

Each command is a Go function implementing a `Commander` interface. Commands emit **narrative events** — finding a treasure file triggers a story beat in the left pane.

---

## Data Model (SQLite)

```sql
-- Player profiles
CREATE TABLE players (
  id         INTEGER PRIMARY KEY,
  name       TEXT NOT NULL,
  tier       TEXT NOT NULL CHECK(tier IN ('beginner','explorer','master')),
  avatar     TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Per-mission completion
CREATE TABLE progress (
  player_id    INTEGER REFERENCES players(id),
  world_id     TEXT NOT NULL,
  mission_id   TEXT NOT NULL,
  stars        INTEGER DEFAULT 0,
  completed_at DATETIME,
  PRIMARY KEY (player_id, world_id, mission_id)
);

-- Virtual filesystem snapshot per player per world
CREATE TABLE world_state (
  player_id  INTEGER REFERENCES players(id),
  world_id   TEXT NOT NULL,
  fs_json    TEXT NOT NULL,  -- JSON snapshot of virtual FS
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (player_id, world_id)
);

-- Shell history (for `history` command)
CREATE TABLE history (
  id         INTEGER PRIMARY KEY,
  player_id  INTEGER REFERENCES players(id),
  command    TEXT NOT NULL,
  timestamp  DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## Mission JSON Format

```json
{
  "id": "skull_island_1",
  "world": "skull_island",
  "title": "The Hidden Cave",
  "starting_clue": "Old Pete says the treasure is hidden deep in the cave. Start by listing what's around you.",
  "objectives": [
    { "event": "ls", "path": "/island" },
    { "event": "cd", "path": "/island/cave" },
    { "event": "cat", "path": "/island/cave/.treasure_map" }
  ],
  "success_message": "You found the treasure map! X marks the spot on Skull Mountain!",
  "treasure_reveal": "A worn map tumbles out. It shows a path to the mountain peak!",
  "unlocks": ["skull_island_2", "open_world:skull_island"],
  "filesystem": {
    "/island": { "type": "dir" },
    "/island/cave": { "type": "dir" },
    "/island/cave/.treasure_map": { "type": "file", "content": "The treasure is buried at /mountain/peak/chest" },
    "/island/cave/old_note.txt": { "type": "file", "content": "Arr, many have searched here before ye..." }
  }
}
```

---

## Game Loop

### Story Missions (Linear)

1. Player selects/creates a profile and chooses a tier
2. Current mission clue shown in left pane
3. Player types commands in the right pane
4. Shell executor emits events; world engine checks against mission objectives
5. When all objectives are met → success animation → stars awarded → next mission unlocked

### Open World (Post-Mission)

Completed areas remain explorable with no objectives. Hidden files, easter eggs, and funny responses to `cat`-ing random files reward curiosity. The "adventure log" tracks all completed missions and star counts.

### Progression & Rewards

- 1–3 stars per mission (based on number of commands used vs. optimal)
- Stars unlock cosmetic rewards: new ship name, ASCII pirate avatar styles, color themes
- No fail state — the game is encouraging, not punishing

### Hint System

- `help` or `?` always shows a tier-appropriate hint
- After 3 failed attempts at an objective, a gentle nudge appears in the story pane
- Parent mode (simple arithmetic gate) allows changing tier or resetting progress

---

## Testing Strategy

- Unit tests for each command implementation (virtual FS operations)
- Integration tests for mission objective detection
- Snapshot tests for TUI rendering (Bubble Tea test helpers)
- Manual playtesting with the target age group

---

## Out of Scope (v1)

- Real filesystem access
- Multiplayer / network features
- Linux packaging (added after macOS is stable)
- Sound / audio
- TypeScript / web dashboard
