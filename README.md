# Shell Quest

A pirate-themed CLI learning game that teaches real shell commands to kids (ages 3–10) through treasure hunting adventures. Single Go binary, no setup required beyond `go build`.

```
        _,.
     ,` -.)
    ( _/-\\-._
   /,|`--._,-^|            ,
   \_| |`-._/||          ,'|
     |  `-, / |         /  /
     |     || |        /  /
      `r-._||/   __   /  /
  __,-<_     )`-/  `./  /
 '  \   `---'   \   /  /
     |           |./  /
     /           //  /
 \_/' \         |/  /
  |    |   _,^-'/  /
  |    , ``  (\/  /_
   \,.->._    \X-=/^
   (  /   `-._//^`
    `Y.-    b'
     |  \   |
     |  |  /|
```

## What You'll Learn

Navigate 10 missions across Skull Island, each teaching specific shell commands through objectives with pirate story flair:

| Mission | Location | Commands |
|---------|----------|----------|
| 1 | Awakening at the Docks | `ls`, `cat` |
| 2 | The Jungle Crossroads | `pwd`, `cd` |
| 3 | The Hidden Cove | `ls -a` |
| 4 | The Captain's Log | `cat` (long files) |
| 5 | The Shipwright's Yard | `mkdir`, `touch` |
| 6 | The Duplication Ritual | `cp` |
| 7 | The Scrambled Hold | `mv` |
| 8 | The Bilge | `rm` |
| 9 | The Spyglass Plateau | `find` |
| 10 | The Root Vault | Final treasure |

## Features

- **Split-pane TUI** — story/objectives on the left, interactive shell on the right
- **Real shell commands** — commands run against a sandboxed virtual filesystem
- **Player profiles** — multiple pirate profiles with persistent progress (SQLite)
- **Difficulty tiers** — Beginner (now), Explorer and Master (coming soon)
- **Hints** — after 3 failed attempts, a clue appears
- **Adventure log** — track mission progress (`Ctrl+L`)
- **Parent mode** — lock the game behind a quick math quiz (`Ctrl+P`)

## Install

Requires Go 1.23+.

```sh
git clone https://github.com/shanewilliams/shell-quest
cd shell-quest
go build -o shellquest ./cmd/shellquest
./shellquest
```

## Controls

| Key | Action |
|-----|--------|
| `Enter` | Run command |
| `Ctrl+L` | Adventure log |
| `Ctrl+P` | Parent mode |
| `Ctrl+C` | Quit |

## Tech

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [modernc SQLite](https://gitlab.com/cznic/sqlite) — pure Go, no CGo

Progress is saved to `~/.shellquest/progress.db`.
