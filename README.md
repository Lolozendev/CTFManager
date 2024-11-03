# CTFManager
CTFManager is my implementation of a manager of dockerised ctf environment.

It is writtent in go and stands at my first serious project in this language.

# How it works

The manager is a simple cli tool that allows you to manage a team,it's members and the challenges they are working on.

Each team is a docker compose project that contains a set of services used to create a ctf network:

## Network

The network is a simple bridge network that allows the team members to connect to the services.

Each team network is created with the following configuration: subnet: `10.0.[team_id].0/24`

## Services

### wireguard
A vpn service that allows the team to connect to the ctf network.

on the host a port between 50000 and 50254 is opened and mapped to the wireguard port. The team members can connect to the vpn using the host ip and the port.

### dnsmasq
A dns server that resolves the ctf domain to the wireguard network. each team member can access the services using the domain name once connected to the vpn.

## Challenges

Each challenges is stored in a separate folder and contains a dokcerfile and an .env file (will add more features like volume mounting and network configuration later).

Each challenge directory is in /challenges and must respect the following structure: `<number>-<name>`

with
- number: the number of the challenge in the network between 11 and 249
- name: the name of the challenge (must be alphanumeric)

# Installation
## Prerequisites
- docker
- docker-compose
- the binary

* (Optional) go [if you want to build the binary]

## Build the binary (optional)
```bash
go build -o ctfmanager
```

## Install the binary
```bash
sudo cp ctfmanager /usr/local/bin
```

## Setup the environment
```bash
sudo ctfmanager setup
```

# Usage

The CLI is still on my TODO list but here is a list of the commands that will be available:

## CTF management
### Setup the ctf
`ctfmanager setup`
### Start the ctf
`ctfmanager start`
### Stop the ctf
`ctfmanager stop`

## Team management
### List teams
`ctfmanager team list`
### Create a team
`ctfmanager team create <team_name>`
### Delete a team
`ctfmanager team delete <team_name>`
### Disables a team
`ctfmanager team disable <team_name>`
### Enables a team
`ctfmanager team enable <team_name>`

## Member management
### List members
`ctfmanager member list <team_name>`
### Add a member
`ctfmanager member add <team_name> <member_name>`
### Remove a member
`ctfmanager member remove <team_name> <member_name>`

## Challenge management
### List challenges
`ctfmanager challenge list`
### Enable a challenge
`ctfmanager challenge enable <challenge_name>`
### Disable a challenge
`ctfmanager challenge disable <challenge_name>`



