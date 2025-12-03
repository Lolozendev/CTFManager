# CTFManager

> Simple CLI for managing isolated, Dockerized CTF environments

Deploy CTF challenges with automatic network isolation, WireGuard VPN, and DNS resolution for each team.

## Features

- 🎯 Isolated Docker networks per team (10.0.X.0/24)
- 🔒 WireGuard VPN configs auto-generated for team members
- 📡 Built-in DNS resolution for challenges
- ⚡ One command to deploy an entire team's infrastructure

## Install

```bash
go build -o ctfmanager ./cmd/ctfmanager
```

Or download from [releases](https://github.com/Lolozendev/CTFManager/releases).

## Quick Start

```bash
# Create a challenge (format: <network-id>-<name>)
mkdir -p challenges/11-webapp
echo "FROM nginx:alpine" > challenges/11-webapp/Dockerfile
touch challenges/11-webapp/.env

# Create a team
ctfmanager team create 1 redteam --members alice,bob

# Deploy
cd equipes/1-redteam
docker compose up -d
```

VPN configs are in `equipes/<team>/config/peer*/`.

## Commands

### Teams
```bash
ctfmanager team list
ctfmanager team create <id> <name> [--members user1,user2]
ctfmanager team delete <name>
```

### Challenges
```bash
ctfmanager challenge list
ctfmanager challenge validate
ctfmanager challenge enable <name> <network-id>
ctfmanager challenge disable <name>
```

Challenges are auto-loaded from `challenges/` directory:
- Enabled: `11-webapp`, `12-crypto` (numbers 11-249)
- Disabled: `x-oldchall` (prefix with `x-`)

## Network Layout

Each team gets:
- Subnet: `10.0.<team_id>.0/24`
- VPN: `.252`, DNS: `.253`, Gateway: `.254`
- Challenges: `.11-.249`
- VPN port: `50000 + team_id`

## Configuration

Edit `internal/config/config.go` to change default paths or network ranges.

## License

GPL-3.0 - See [LICENSE](LICENSE) for details.
