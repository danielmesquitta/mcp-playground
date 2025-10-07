# Brasil CEP MCP Server 🇧🇷

An MCP (Model Context Protocol) server written in Go that provides Brazilian ZIP code (CEP) lookup functionality using the [Brasil API](https://brasilapi.com.br/).

## Features

- 🔍 Look up Brazilian addresses by ZIP code (CEP)
- ✅ Accepts CEP with or without hyphen (01310-100 or 01310100)
- 🚀 Fast and lightweight Go implementation
- 🎯 MCP-compliant for integration with AI assistants

## Quick Start

### Prerequisites

- Go 1.25 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/danielmesquitta/mcp-playground
cd mcp-playground

# Install dependencies
make install

# Build the server
make build
```

### Usage with Claude Desktop

Add to your Claude Desktop configuration file (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "brasil-cep": {
      "command": "/absolute/path/to/cep-server",
      "args": []
    }
  }
}
```

Restart Claude Desktop, and you can now use the `lookup_address` tool!

## Tool: lookup_address

Retrieves address information from a Brazilian ZIP code.

**Parameters:**

- `cep` (required): Brazilian ZIP code with or without hyphen (e.g., "01310-100" or "01310100")

**Example Response:**

```
Street: Avenida Paulista
Neighborhood: Bela Vista
City: São Paulo
State: SP
CEP: 01310-100
```

## API Reference

This server uses the Brasil API CEP endpoint:

```
GET https://brasilapi.com.br/api/cep/v1/{cep}
```

## Acknowledgments

- [Brasil API](https://brasilapi.com.br/) for providing free Brazilian public data
- [mcp-go](https://github.com/mark3labs/mcp-go) for the excellent MCP implementation
