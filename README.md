# caffeinate-mcp

An MCP (Model Context Protocol) server for macOS caffeinate command.

## Overview

caffeinate-mcp is a server that allows you to control the macOS caffeinate command via MCP. It can create and manage assertions to prevent the system from sleeping.

## Installation

### Using npx

```bash
npx -y @upamune/caffeinate-mcp
```

### Using npm

```bash
npm install -g @upamune/caffeinate-mcp
```

## Usage

### Available Tools

#### `caffeinate_start`

Start a caffeinate process to prevent system sleep.

Parameters:
- `display` (boolean): Prevent the display from sleeping (-d flag)
- `idle` (boolean): Prevent the system from idle sleeping (-i flag)
- `disk` (boolean): Prevent the disk from idle sleeping (-m flag)
- `system` (boolean): Prevent the system from sleeping when on AC power (-s flag)
- `user` (boolean): Declare that user is active (-u flag)
- `timeout` (number): Timeout value in seconds (-t flag)
- `pid` (number): Wait for the process with specified PID to exit (-w flag)

#### `caffeinate_stop`

Stop a running caffeinate process.

Parameters:
- `id` (string, required): ID of the caffeinate process to stop

#### `caffeinate_list`

List active caffeinate processes.

## MCP Client Configuration

### Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "caffeinate": {
      "command": "npx",
      "args": ["-y", "@upamune/caffeinate-mcp"],
      "env": {}
    }
  }
}
```

## Development

### Setup

```bash
npm install
```

### Building

```bash
npm run build
```

### Testing

```bash
npm test
```

### Linting

```bash
npm run lint
```

### Type Checking

```bash
npm run typecheck
```

## License

MIT