# dm-system Unix Tools

CLI tool for sending admin commands to the dm-system backend over its Unix domain socket.

## Build

```bash
make
```

This builds the `tool` binary via `go build -o tool`.

Run lint checks with:

```bash
make lint
```

## Usage

```bash
./tool -a <action> -s <socket file>
```

Flags:

- `-a, --action` (required): action to perform. Supported values:

    - `resetAccount`

- `-s, --socket` (default: `./dm.sock`): path to the backend's Unix socket file.

Example:

```bash
./tool -a resetAccount -s ./dm.sock
```

## Layout

- `main.go` — entrypoint, delegates to `cmd`
- `cmd/` — cobra command definition and action dispatch
- `action/` — implementation of each supported action
- `unix/` — Unix socket client used to talk to the backend
