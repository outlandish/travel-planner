# travelplanner

A holiday-planning assistant written in Go. Give it a budget, a holiday type, and a
preferred environment, and it builds a prompt, sends it to DeepSeek, and returns a list
of suggested countries. Responses are cached in SQLite so repeated requests are instant.

Two entry points are available, backed by the same core logic:

- **CLI** — a command-line tool.
- **HTTP API** — a small web server (stdlib `net/http`, no frameworks).

## Requirements

- Go 1.26+
- A [DeepSeek](https://www.deepseek.com/) API key

## Environment variables

| Variable           | Used by     | Description                                        |
| ------------------ | ----------- | -------------------------------------------------- |
| `DEEPSEEK_API_KEY` | CLI + API   | DeepSeek API key (required)                        |
| `SQLITE_DB_PATH`   | CLI + API   | Path to the SQLite cache database (required)       |
| `SERVER_PORT`      | API only    | Port the server listens on (required for the API)  |

The SQLite driver (`modernc.org/sqlite`) is pure Go, so no CGO is required.

## Parameters

| Parameter | Type   | Allowed values                       |
| --------- | ------ | ------------------------------------ |
| `budget`  | int    | positive integer (`> 0`)             |
| `type`    | string | `active`, `beach`, `attractions`     |
| `nature`  | string | `sea`, `mountains`, `city`           |

## CLI

```sh
export DEEPSEEK_API_KEY=your-key
export SQLITE_DB_PATH=./travelplanner.db

go run ./cmd/travelplanner --budget=1000 --type=beach --nature=sea
```

The result is printed to stdout.

## HTTP API

Start the server:

```sh
export DEEPSEEK_API_KEY=your-key
export SQLITE_DB_PATH=./travelplanner.db
export SERVER_PORT=8090

go run ./cmd/server
```

Request a plan:

```sh
curl "http://localhost:8090/plan?budget=1000&type=beach&nature=sea"
```

The endpoint is `GET /plan` with query parameters `budget`, `type`, and `nature`.

Success response (`200 OK`):

```json
{"response": "- Country: ...\nWhy it fits: ...", "status": "success"}
```

Error responses use the same shape with `"status": "error"` and the message in
`"response"`.

| Status | Meaning                            |
| ------ | ---------------------------------- |
| `200`  | success                            |
| `400`  | empty/invalid request parameters   |
| `500`  | AI or database error               |

## Project structure

```
cmd/
  travelplanner/   CLI entry point
  server/          HTTP API entry point
internal/
  app/             shared business flow (cache → AI → cache)
  domain/          types and validation (HolidayParams)
  prompt/          prompt construction
  ai/              DeepSeek HTTP client
  storage/         SQLite cache
```

## How it works

1. Parameters are parsed and validated (`domain.HolidayParams.Validate`).
2. A cache key is built from the parameters and looked up in SQLite.
3. On a cache miss, a prompt is built and sent to DeepSeek (`deepseek-chat`).
4. The answer is cached and returned.
