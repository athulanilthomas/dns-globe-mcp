MCP_URL ?= http://localhost:8080/mcp
DOMAIN ?= athulanilthomas.in
RECORD_TYPE ?= AAAA

export MCP_URL DOMAIN RECORD_TYPE

.PHONY: call-dns-propagation
call-dns-propagation:
	@set -e; \
	request=$$(python3 -c 'import json, os; print(json.dumps({"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "check_dns_propagation", "arguments": {"domain": os.environ["DOMAIN"], "recordType": os.environ["RECORD_TYPE"]}}}))'); \
	curl --fail-with-body --silent --show-error --request POST \
		--header 'Content-Type: application/json' \
		--header 'Accept: application/json, text/event-stream' \
		--header 'MCP-Protocol-Version: 2025-11-25' \
		--data-binary "$$request" \
		"$(MCP_URL)"
