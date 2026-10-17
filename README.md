# DNS Globe MCP

An [MCP](https://modelcontextprotocol.io) server that looks up DNS records across several public resolvers and shows the result on an interactive 3D globe, using the [MCP Apps](https://github.com/modelcontextprotocol/ext-apps) UI standard.

## What it does

The server exposes one tool, `lookup_dns_records`:

| Input        | Description                                           |
| ------------ | ----------------------------------------------------- |
| `domain`     | Domain to check, e.g. `example.com`                   |
| `recordType` | One of `A`, `AAAA`, `CNAME`, `MX`, `NS`, `TXT`, `SOA` |

It queries Cloudflare, Google, Quad9 and OpenDNS in parallel over DNS-over-HTTPS and returns one result per resolver:

- `status`: `resolved`, `pending` (empty answer), `nxdomain`, or `error`. The `stale` status exists in the widget but the server does not produce it yet.
- `records`: the answers returned
- `error`: the failure message, if any

The result is a short text summary plus `structuredContent` (`domain`, `recordType`, `results`).

## Usage

The MCP endpoint is `https://lab.athulanilthomas.in/mcp`. To run it locally, see [Run locally](#run-locally) and use `http://localhost:8080/mcp` instead.

### Cloudflare AI Playground

1. Open [playground.ai.cloudflare.com](https://playground.ai.cloudflare.com).
2. In the MCP servers section of the sidebar, enter the endpoint URL and click **Connect**.
3. Ask: "Look up the A record for example.com".

### VS Code

1. Run **MCP: Add Server** from the Command Palette, choose **HTTP**, and enter the endpoint URL. Or add it to `.vscode/mcp.json`:
   ```json
   {
     "servers": {
       "dns-globe": { "type": "http", "url": "https://lab.athulanilthomas.in/mcp" }
     }
   }
   ```
2. Click **Start** above the server entry in `mcp.json`.
3. Open Copilot Chat in Agent mode and ask: "Look up the A record for example.com". Allow the tool when prompted. The globe renders inline in the chat.

### Claude and ChatGPT

Add the endpoint URL as a custom connector in the app's connector settings. For a local server, expose it with a tunnel such as `cloudflared tunnel --url http://localhost:8080` and use the tunnel URL.

### WebMCP tools in Chrome

The globe also registers two [WebMCP](https://webmachinelearning.github.io/webmcp/) tools that let an agent control it:

- `focus_resolver`: rotate to a resolver (Cloudflare, Google, Quad9, OpenDNS) and highlight it
- `get_globe_state`: return the highlighted resolver and every resolver's status and records

To try them:

1. In Chrome, open `chrome://flags/#enable-webmcp-testing`, set it to **Enabled**, and relaunch.
2. Open a host that renders the globe (for example the Cloudflare AI Playground) in that Chrome, and run a DNS check.
3. Ask the agent to "focus on Quad9" or "which resolvers haven't resolved?".

## Run locally

Requires Go and [pnpm](https://pnpm.io).

```sh
cd widget && pnpm install && cd ..
make build-widget
go run ./cmd/server
```

The server listens on port 8080 by default; set `PORT` to change it. Test the tool from the command line:

```sh
DOMAIN=example.com RECORD_TYPE=A make call-lookup-dns-records
```