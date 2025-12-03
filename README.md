# CTFManager

> A professional CLI tool for managing Dockerized CTF (Capture The Flag) environments

CTFManager simplifies the deployment and management of isolated, team-based CTF infrastructure using Docker Compose. Each team gets its own network environment with VPN access, DNS resolution, and challenge containers.

## Features

- 🎯 **Team Management**: Create, enable, disable, and delete CTF teams
- 🏆 **Challenge Management**: Organize and deploy Dockerized challenges
- 🔒 **Network Isolation**: Each team gets an isolated Docker network (10.0.X.0/24)
- 🌐 **VPN Access**: WireGuard VPN for secure team connectivity (ports 50000-50254)
- 📡 **DNS Resolution**: Built-in dnsmasq for challenge domain resolution
- ⚡ **Auto-Generation**: Automatic Docker Compose file generation per team

## Architecture

### Network Structure

Each team is assigned:
- **Subnet**: `10.0.<team_id>.0/24`
- **Gateway**: `10.0.<team_id>.254`
- **VPN**: `10.0.<team_id>.252` (WireGuard)
- **DNS**: `10.0.<team_id>.253` (dnsmasq)
- **Challenges**: `10.0.<team_id>.11-249`

### Directory Structure

```
CTFManager/
├── cmd/
│   └── ctfmanager/          # CLI entry point
│       └── main.go
├── internal/                # Private application code
│   ├── app/                 # Application logic
│   │   ├── challenge/       # Challenge management
│   │   ├── team/            # Team management
│   │   └── compose/         # Docker Compose generation
│   ├── model/               # Data models
│   │   ├── challenge.go
│   │   ├── team.go
│   │   ├── service.go
│   │   └── network.go
│   ├── config/              # Configuration
│   │   └── config.go
│   └── logger/              # Logging utilities
│       └── logger.go
├── challenges/              # Challenge definitions
│   ├── 11-web1/
│   │   ├── Dockerfile
│   │   └── .env
│   ├── 12-crypto1/
│   │   ├── Dockerfile
│   │   └── .env
│   └── x-disabled/          # Disabled challenges (prefix: x-)
├── equipes/                 # Team directories (auto-generated)
│   ├── 1-teamA/
│   │   └── compose.yml
│   └── 2-teamB/
│       └── compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- **Docker** (with Docker Compose v2+)
- **Go** 1.23+ (for building from source)
- **Linux/Unix** environment (recommended for production)

## Installation

### Option 1: Download Binary (Recommended)

```bash
# Download the latest release
curl -LO https://github.com/Lolozendev/CTFManager/releases/latest/download/ctfmanager

# Make it executable
chmod +x ctfmanager

# Move to PATH
sudo mv ctfmanager /usr/local/bin/
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/Lolozendev/CTFManager.git
cd CTFManager

# Build
go build -o ctfmanager ./cmd/ctfmanager

# Install (optional - installs to $GOPATH/bin)
go install ./cmd/ctfmanager
```

## Quick Start

### 1. Setup Environment

```bash
# Initialize CTF environment (validates configuration)
ctfmanager setup
```

### 2. Prepare Challenges

Create challenge directories in `/challenges` with the format: `<network_id>-<name>`

```bash
mkdir -p /challenges/11-webapp
cd /challenges/11-webapp

# Create Dockerfile
cat > Dockerfile <<EOF
FROM nginx:alpine
COPY index.html /usr/share/nginx/html/
EOF

# Create .env file (even if empty)
touch .env
```

**Challenge Naming Rules:**
- **Enabled**: `<number>-<name>` (e.g., `11-webapp`, `12-crypto`)
- **Disabled**: `x-<name>` (e.g., `x-beta-challenge`)
- **Network IDs**: 11-249 (reserved range)

### 3. Validate Challenges

```bash
# Check all challenges for correctness
ctfmanager challenge validate

# List all challenges
ctfmanager challenge list

# List including disabled
ctfmanager challenge list --all
```

### 4. Create a Team

```bash
# Create team with ID 1
ctfmanager team create 1 redteam --members alice,bob,charlie

# This generates:
#   - /equipes/1-redteam/ directory
#   - compose.yml with all enabled challenges
#   - WireGuard configs for 3 members
```

### 5. Deploy Team Infrastructure

```bash
cd /equipes/1-redteam
docker compose up -d

# Access WireGuard config
ls config/peer1/  # alice's VPN config
ls config/peer2/  # bob's VPN config
ls config/peer3/  # charlie's VPN config
```

## Usage

### Team Management

```bash
# List all teams
ctfmanager team list

# Create a new team
ctfmanager team create <id> <name> [--members user1,user2,...]

# Delete a team
ctfmanager team delete <name>

# Disable a team (stops but preserves data)
ctfmanager team disable <name>

# Re-enable a team
ctfmanager team enable <name> <id>
```

**Examples:**
```bash
ctfmanager team create 1 blueteam --members alice,bob
ctfmanager team create 2 redteam --members eve,mallory,trudy
ctfmanager team disable blueteam
ctfmanager team enable blueteam 1
```

### Challenge Management

```bash
# List enabled challenges
ctfmanager challenge list

# List all challenges (including disabled)
ctfmanager challenge list --all

# Validate challenge structure
ctfmanager challenge validate

# Enable a disabled challenge
ctfmanager challenge enable <name> <network-id>

# Disable a challenge
ctfmanager challenge disable <name>
```

**Examples:**
```bash
ctfmanager challenge enable beta-web 15
ctfmanager challenge disable old-crypto
ctfmanager challenge validate
```

## Configuration

Default configuration (can be customized in `internal/config/config.go`):

```go
Paths:
  Challenges:      "/challenges"
  Teams:           "/equipes"
  DnsmasqTemplate: "/dnsconf/dnsmasq.template"

Network:
  BaseSubnet: "10.0"

Challenges:
  MinNetworkID: 11
  MaxNetworkID: 249

Teams:
  MinID:       1
  MaxID:       254
  BaseVPNPort: 50000
```

## Challenge Development

### Minimal Challenge Template

```dockerfile
# Dockerfile
FROM ubuntu:22.04

# Install your challenge
COPY challenge.py /app/
COPY flag.txt /root/

WORKDIR /app
CMD ["python3", "challenge.py"]
```

### Environment Variables

Use `.env` files for challenge-specific configuration:

```bash
# .env
FLAG=CTF{this_is_a_flag}
PORT=8080
DIFFICULTY=medium
```

### Best Practices

1. **Isolation**: Never share secrets between challenges
2. **Resource Limits**: Set memory/CPU limits in Docker
3. **Logging**: Log user interactions for monitoring
4. **Health Checks**: Add Docker health checks
5. **Documentation**: Include README.md in each challenge

## Troubleshooting

### Build Issues

```bash
# Clean build
rm -rf build/
go build -o build/ctfmanager ./cmd/ctfmanager

# Verify Go modules
go mod tidy
go mod verify
```

### Network Conflicts

```bash
# Check for port conflicts
netstat -tuln | grep 50000

# Check Docker networks
docker network ls
docker network inspect <network-name>
```

### Challenge Not Loading

```bash
# Validate challenge structure
ctfmanager challenge validate

# Check Docker logs
cd /equipes/<team>
docker compose logs <challenge-name>
```

## Development

### Project Structure

- `cmd/`: Command-line entry points
- `internal/`: Private packages (not importable by other projects)
  - `app/`: Business logic (challenge, team, compose managers)
  - `model/`: Data structures and domain models
  - `config/`: Configuration management
  - `logger/`: Centralized logging

### Building

```bash
# Build binary
go build -o build/ctfmanager ./cmd/ctfmanager

# Run tests
go test -v ./...

# Format code
go fmt ./...

# Tidy modules
go mod tidy
```

### Adding New Commands

1. Create command in `cmd/ctfmanager/main.go`
2. Implement logic in `internal/app/`
3. Add models if needed in `internal/model/`
4. Update documentation

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Write tests for new features
- Update documentation
- Keep functions small and focused (KISS principle)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Uses [Charm Log](https://github.com/charmbracelet/log) for beautiful logging
- Inspired by the need for simple, scalable CTF infrastructure

## Support

- 📧 Email: [your-email@example.com]
- 🐛 Issues: [GitHub Issues](https://github.com/Lolozendev/CTFManager/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/Lolozendev/CTFManager/discussions)

---

**Made with ❤️ for the CTF community**
