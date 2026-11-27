# mymquery

A terminal UI for [MQuery](https://github.com/czcorpus/mquery), a corpus querying HTTP API
built on Manatee-Open: concordances, frequencies, text types and collocations, plus an
AI chat that queries corpora for you (Gemini, Mistral or an OpenAI-compatible server such
as llama-server).

mymquery is pure Go (no CGO, no Manatee needed) – it just talks to a running MQuery server.

## Install

```
go install github.com/tomachalek/mymquery@latest
```

or from the sources:

```
make            # tests + build
./mymquery -api http://localhost:8989
```

## Configuration

Options can be set by flags (`mymquery -h`) or by variables in env files:

1. `$XDG_CONFIG_DIRS/mymquery/env` (e.g. `/etc/xdg/mymquery/env`)
2. `$XDG_CONFIG_HOME/mymquery/env` (e.g. `~/.config/mymquery/env`)
3. `./.env` – only when the binary is run from the current directory (development), or
   a file given by `-env-file`

A later file overrides an earlier one, real environment variables override the files.
Files containing secrets should be readable only by you (`chmod 600`).

| Variable | Meaning |
|---|---|
| `MQUERY_API_URL` (`API_URL`) | MQuery API URL (default `http://localhost:8989`) |
| `API_TOKEN`, `LOGIN_URL` | log in to a protected MQuery instance |
| `LLM_BACKEND` | `gemini`, `mistral` or `openai` (empty = the first one configured) |
| `LLM_MODEL` | comma separated models; the first is primary, the others are fallbacks |
| `GEMINI_API_KEY`, `MISTRAL_API_KEY` | API keys of the hosted backends |
| `OPENAI_API_URL` | OpenAI-compatible server (default `http://localhost:8080/v1`, llama-server) |
| `OPENAI_API_KEY` | optional token of that server (e.g. llama-server `--api-key`) |
| `CQLIZER_MCP` | CQLizer MCP server providing CQL helper tools to the AI chat |

The AI chat can also be set up from the UI (`ctrl+l` in the chat), optionally saving
the setup to `~/.config/mymquery/env`.

Query history and AI chats are stored in `$XDG_STATE_HOME/mymquery/mymquery.db`
(`~/.local/state/mymquery/mymquery.db`; `-db ""` disables it).

## Development

See [PLAN.md](PLAN.md) for the architecture, the implemented features and the backlog.
mymquery used to live in the MQuery repository (`cmd/mymquery`); its history was preserved.
