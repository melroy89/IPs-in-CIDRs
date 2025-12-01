# IPs in CIDR

IPs in CIDR is a simple Go tool which checks a list of IPs against a list of CIDRs. And report back the matches.

## Production

Create two text files, called:

- `cidrs.txt` (with your list of CIDRs, each new line a CIDR notation)
- `ips.txt` (with your list of IPs, each new line a IP address)

And then run the `./ips-in-cidrs` binary (see the latest binary in: [GitLab CICD artifacts](https://gitlab.melroy.org/melroy/ips-in-cidrs/-/jobs/artifacts/main/browse?job=build). You can download the pre-build binary from there.)

## Development

### Requirements

- [Golang](https://go.dev/doc/install)

### Start dev

Create two text files:

- `cidrs.txt`
- `ips.txt`

Either run: `go run .`

Or if you want to have watch mode, use: `gow run .`

### Build binary

Run: `go build .`

### Getting started

Assuming you already fulfilled the requirements above.

1. Clone the project: `git clone git@gitlab.melroy.org:melroy/ips-in-cidr.git`
2. To start the project executing: `go run .`

