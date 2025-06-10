# caffeinate-mcp

An MCP (Model Context Protocol) server for macOS caffeinate command.

## Overview

caffeinate-mcp is a server that allows you to control the macOS caffeinate command via MCP. It can create and manage assertions to prevent the system from sleeping.

## Installation

### Using Go

```bash
go install github.com/upamune/caffeinate-mcp@latest
```

### Download Binary

Download the latest binary from the [releases page](https://github.com/upamune/caffeinate-mcp/releases).

## Usage

### Starting the Server

```bash
caffeinate-mcp
```

### Available Tools

#### `caffeinate_start`

Start a caffeinate process to prevent system sleep.

Parameters:
- `display` (boolean): Prevent the display from sleeping (-d flag)
- `idle` (boolean): Prevent the system from idle sleeping (-i flag)
- `disk` (boolean): Prevent the disk from idle sleeping (-m flag)
- `system` (boolean): Prevent the system from sleeping when on AC power (-s flag)
- `user` (boolean): Declare that user is active (-u flag)
- `timeout` (integer): Timeout value in seconds (-t flag)
- `pid` (integer): Wait for the process with specified PID to exit (-w flag)

#### `caffeinate_stop`

Stop a running caffeinate process.

Parameters:
- `id` (string, required): ID of the caffeinate process to stop

#### `caffeinate_list`

List active caffeinate processes.

## MCP Client Configuration

```json
{
  "mcpServers": {
    "caffeinate": {
      "command": "caffeinate-mcp",
      "args": [],
      "env": {}
    }
  }
}
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Create a Snapshot Release

```bash
make release-snapshot
```

### Available Make Targets

```bash
make help    # Show all available targets
```

## License

MIT