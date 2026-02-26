# trello-cli

A command-line interface for the [Trello API](https://developer.atlassian.com/cloud/trello/rest/), written in Go.

Manage boards, lists, cards, members, checklists, and labels directly from your terminal. Outputs human-readable tables in a terminal and JSON when piped — making it usable both interactively and as part of automated workflows.

---

## Installation

```bash
git clone https://github.com/the20100/trello-cli
cd trello-cli
go build -o trello .
mv trello /usr/local/bin/
```

---

## Authentication

You need a Trello **API key** and **API token**.

1. Go to https://trello.com/power-ups/admin and create (or use) a Power-Up to get your API key
2. Generate a token by visiting:
   ```
   https://trello.com/1/authorize?expiration=never&scope=read,write&response_type=token&key=YOUR_KEY
   ```
3. Save your credentials:
   ```bash
   trello auth setup YOUR_API_KEY YOUR_API_TOKEN
   ```

Credentials are stored at:
- macOS: `~/Library/Application Support/trello/config.json`
- Linux: `~/.config/trello/config.json`
- Windows: `%AppData%\trello\config.json`

You can also use environment variables (takes priority over config file):
```bash
export TRELLO_API_KEY=your_api_key
export TRELLO_API_TOKEN=your_api_token
```

---

## Usage

```
trello [command] [subcommand] [flags]
```

### Global flags

| Flag | Description |
|------|-------------|
| `--json` | Force JSON output |
| `--pretty` | Force pretty-printed JSON output |
| `--help` | Help for any command |

---

## Commands

### `auth`

```bash
trello auth setup <api-key> <api-token>   # Save credentials (validates against API)
trello auth status                         # Show current auth status
trello auth logout                         # Remove saved credentials
```

### `info`

```bash
trello info   # Show binary path, config location, credential source
```

---

### `boards`

```bash
trello boards list                                # List your boards (open by default)
trello boards list --filter all                   # Include closed boards
trello boards get <board-id>                      # Get board details
trello boards create "My Project"                 # Create a board
trello boards create "Q1" --workspace <ws-id>     # Create in a workspace
trello boards update <board-id> --name "New Name" # Rename a board
trello boards update <board-id> --closed          # Archive a board
trello boards delete <board-id>                   # Delete a board (permanent)
trello boards members <board-id>                  # List board members
trello boards labels <board-id>                   # List board labels
```

---

### `lists`

```bash
trello lists list --board <board-id>              # List all lists on a board
trello lists list --board <id> --filter all       # Include archived lists
trello lists get <list-id>                        # Get list details
trello lists create "To Do" --board <board-id>    # Create a list
trello lists create "Done" --board <id> --pos bottom
trello lists rename <list-id> "In Progress"       # Rename a list
trello lists archive <list-id>                    # Archive a list
trello lists unarchive <list-id>                  # Unarchive a list
trello lists cards <list-id>                      # List cards in a list
```

---

### `cards`

```bash
trello cards list --board <board-id>              # List cards on a board
trello cards list --list <list-id>                # List cards in a list
trello cards list --board <id> --filter all       # Include archived cards
trello cards get <card-id>                        # Get card details
trello cards create "Fix the bug" --list <list-id>
trello cards create "Deploy" --list <id> --desc "Deploy v2" --due 2024-12-31
trello cards update <card-id> --name "New title"
trello cards update <card-id> --due 2024-12-31 --due-complete
trello cards move <card-id> --list <target-list-id>          # Move to list
trello cards move <card-id> --list <list-id> --board <board-id>  # Cross-board move
trello cards archive <card-id>                    # Archive a card
trello cards delete <card-id>                     # Delete a card (permanent)
trello cards comment <card-id> "Looks good!"      # Add a comment
trello cards checklists <card-id>                 # Show checklists with items
trello cards attachments <card-id>                # List attachments
trello cards label <card-id> --add <label-id>     # Add a label
trello cards label <card-id> --remove <label-id>  # Remove a label
trello cards member <card-id> --add <member-id>   # Assign a member
trello cards member <card-id> --remove <member-id> # Unassign a member
```

---

### `members`

```bash
trello members me                         # Show your profile
trello members get <id-or-username>       # Get any member's profile
trello members boards                     # Your boards
trello members boards johndoe             # Another member's boards
trello members cards                      # Cards assigned to you
trello members cards johndoe              # Cards assigned to another member
trello members workspaces                 # Your workspaces
```

---

### `checklists`

```bash
trello checklists create "Criteria" --card <card-id>          # Create a checklist
trello checklists delete <checklist-id>                        # Delete a checklist
trello checklists add-item "Write tests" --checklist <cl-id>  # Add item
trello checklists check <item-id> --card <card-id> --checklist <cl-id>    # Check item
trello checklists uncheck <item-id> --card <card-id> --checklist <cl-id>  # Uncheck item
```

---

### `search`

```bash
trello search "deploy"                    # Search everything
trello search "bug" --type cards          # Cards only
trello search "John" --type members       # Members only
trello search "project" --limit 5         # Limit results per type
```

---

## JSON output & piping

When stdout is not a TTY (i.e. piped to another command), output is automatically JSON. You can also force it with `--json` or `--pretty`.

```bash
# Get all board IDs
trello boards list --json | jq '.[].id'

# Get open cards in a specific list
trello cards list --list <list-id> --json | jq '.[] | {id, name, due}'

# Find your first board's ID
BOARD_ID=$(trello boards list --json | jq -r '.[0].id')

# List all lists on that board
trello lists list --board "$BOARD_ID" --json | jq '.[] | {id, name}'

# Get all cards with a due date
trello cards list --board <id> --json | jq '[.[] | select(.due != null)]'
```

---

## Credential resolution order

1. `TRELLO_API_KEY` + `TRELLO_API_TOKEN` environment variables
2. Config file (`trello auth setup`)

---

## Project structure

```
trello-cli/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go          # Root command, auth resolution, info
│   ├── auth.go          # auth setup / status / logout
│   ├── boards.go        # boards subcommands
│   ├── cards.go         # cards subcommands
│   ├── lists.go         # lists subcommands
│   ├── members.go       # members subcommands
│   ├── checklists.go    # checklists subcommands
│   ├── search.go        # search command
│   └── helpers.go       # shared helpers (buildParams, printAPICardsTable)
└── internal/
    ├── api/
    │   ├── client.go    # HTTP client + all API methods
    │   └── types.go     # Board, Card, TrelloList, Member, etc.
    ├── config/
    │   └── config.go    # Config load/save/clear
    └── output/
        └── output.go    # Table, JSON, formatting helpers
```

---

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) — CLI framework
- [mattn/go-isatty](https://github.com/mattn/go-isatty) — TTY detection for auto JSON/table output
