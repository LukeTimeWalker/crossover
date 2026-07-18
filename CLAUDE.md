# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Crossover is a lightweight, declarative webhook translator written in Go. It receives HTTP requests on
configurable routes and forwards (translates) them into HTTP requests to other configurable
endpoints, using Go's `text/template` to map data between the incoming request and the outgoing
one(s). It's config-driven (see `config.yaml`) — most behavior changes come from route/sender/condition
config, not code. The project is under active development and explicitly not production-ready
(see README warning).

## Commands

Task (`Taskfile.yaml`) is the primary entry point; run `task` with no args to list tasks.

- `task go:test` — runs `go test ./...` (equivalent to running tests directly with `go test`)
- `task go:run` — runs Crossover locally via `go run ./cmd/crossover/main.go` (looks for `./config.yaml`)
- `task go:build PROJECT_VERSION=<version>` — cross-compiles the binary into `./target/` (requires `PROJECT_VERSION`)
- `task image:build PROJECT_VERSION=<version>` — builds a container image with Podman (requires `PROJECT_VERSION`)

Single test example: `go test ./pkg/utils/ -run TestName -v`

## Architecture

Two packages under `pkg/`:

- **`pkg/core`** — the domain logic: config loading, route registration, request handling, conditions, HTTP clients.
- **`pkg/utils`** — stateless helpers used by `core`: templating, body/URL parsing, env var loading, logging.

Entry point is `cmd/crossover/main.go`, which loads `Config` (`pkg/core/config.go`) and constructs a
`Crossover` (`pkg/core/crossover.go`), which wraps a single `http.ServeMux` + `http.Server`.

### Request flow (the core mental model)

1. **Config → Routes.** `Config` (`pkg/core/config.go`) is unmarshalled from YAML. Each `Route` has one
   `Receiver` (the inbound endpoint/method) and one or more `Sender`s (outbound requests), plus optional
   route-level and sender-level `Conditions`.
2. **Route registration.** `Crossover.RegisterRoute` (`pkg/core/crossover.go`) builds a `RouteHandler` via
   `CreateRouteHandler` (`pkg/core/route.go`) and returns an error if any template in the route fails to
   parse (`NewCrossover` exits on that; runtime callers can handle it). Because `http.ServeMux` panics on
   duplicate pattern registration, each pattern (`"METHOD /endpoint"`) is registered on the mux exactly
   once with a wrapper that looks up the current handler in a mutex-guarded `handlers` map; `RegisterRoute`
   and `UnregisterRoute` swap entries in that map (nil → 404), so routes can be added, removed, and
   re-added at runtime (see the dynamic-route example in `main.go`).
3. **Per-request handling** (`CreateRouteHandler`'s returned closure):
   - Parse the inbound JSON body (`utils.LoadBody`) and URL query params (`utils.LoadUrlParams`).
   - Assemble a `Variables` map (`pkg/core/route.go`): `{"body": ..., "url": ..., "env": ...}`. This is
     the data available to every template and condition.
   - Evaluate route-level conditions (`EvaluateConditions`, `pkg/core/condition.go`) — if any fails, stop
     with 200 OK and no outbound requests are sent.
   - Run each sender's `PathHandler` (built once at registration time by `CreatePathHandler`), each of
     which templates its own body/URL/headers, evaluates its own sender-level conditions, and sends the
     outbound request via a per-`Sender` cached `http.Client` (`GetClient`, `pkg/core/client.go`).
   - Response status is 200 if all senders succeeded, 400 if any failed.
4. **Templating.** All templated fields (`body`, `endpoint`, header values, conditions) use Go's
   `text/template` (`utils.CreateTemplate`/`utils.ExecuteTemplate`, `pkg/utils/template.go`) against the
   `Variables` map, so `{{ .body.foo }}`, `{{ .url.bar }}`, `{{ .env.TOKEN }}` are all valid inside config.
5. **Conditions.** A condition is a template string that must evaluate to the literal string `"true"` or
   `"false"` (parsed via `strconv.ParseBool`); anything else is an error. Conditions are AND'd together
   within a list; use the template `or` function for logical OR within a single condition string.

### Notable implementation details

- HTTP clients are cached per `*Sender` pointer (`pkg/core/client.go`), not recreated per request — this
  is why route/sender structs are built once at registration and reused across requests.
- Logging is a custom minimal `slog.Handler` (`pkg/utils/log.go`) with a numeric log level (0=ERROR..3=DEBUG)
  configured via `config.yaml`'s `logLevel`. There's a global `sync.Mutex`-guarded singleton logger
  (`GetLogger`) — don't introduce a second logger instance.
- `utils.LoadBody` escapes `\n` to `\\n` in the raw request body before JSON-unmarshalling — this is
  intentional handling for payloads containing literal newlines inside JSON string values (e.g. from
  webhook providers), not an oversight.

## Tests

- `pkg/utils/*_test.go` — unit tests for the utility helpers (uses `testify`).
- `tests/` — package `tests`, black-box integration tests that spin up a real `Crossover` instance
  (`httptest.NewServer(crossover.Router())`) plus mock downstream receivers to assert on end-to-end
  behavior. One file per capability; when adding a test, reuse the scaffolding and put it in the
  matching file:
  - `helpers_test.go` — shared scaffolding: `MockReceiver` (options `WithStatus`/`WithDelay`, methods
    `WaitForRequests`/`AssertNoRequests`) and the `startCrossover`/`singleRouteConfig`/`postSender`/
    `postJSON` helpers. Configs are built as Go structs, so defaults from `withDefaults` don't apply —
    always set `Client.Timeout` (via `defaultClient()`), or `GetClient` fails registration.
  - `routing_test.go` — method/endpoint matching (404/405) and dynamic register/unregister.
  - `templating_test.go` — `.body`/`.url`/`.env` templating in bodies, headers, and sender endpoints.
  - `conditions_test.go` — route- and sender-level condition semantics (false → 200 OK, nothing sent;
    non-boolean → 400).
  - `fanout_test.go` — multi-sender fan-out.
  - `failures_test.go` — downstream failures, bad inbound input, client timeouts, registration errors.
