![Crossover Logo](./assets/header_dark.svg#gh-dark-mode-only)
![Crossover Logo](./assets/header_light.svg#gh-light-mode-only)

Crossover is a lightweight and declarative webhook translator written in Go🩵  
It basically receives HTTP-Requests on pre-configured routes and sends out HTTP-Request to different also configurable routes.  
This project focuses on `minimalism` and `simplicity`.

![Crossover Architecture Overview](./assets/overview_light.svg#gh-light-mode-only)
![Crossover Architecture Overview](./assets/overview_dark.svg#gh-dark-mode-only)

> [!WARNING]
> This project is still under development. Do not use it in production.

## Installation

To install Crossover, you can either download the binary or use the Docker image.

## Usage

Crossover can be configured via a `config.yaml` file, which can be specified with the `-config` flag. If no config file is specified, the system searches for `./config.yaml` in the same directory as Crossover by default.

### Binary

```bash
./crossover -config "./config.yaml"
```

### Docker/Podman

```bash
docker run -v ./config.yaml:/crossover/config.yaml -p 9999:9999 crossover
```

## Config

| Key                                        | Description                                              | Default |
|--------------------------------------------|----------------------------------------------------------|---------|
| `port`                                     | Port on which the Crossover server should listen         | 9999    |
| `logLevel`                                 | The log level (0 = ERROR, 1 = WARN, 2 = INFO, 3 = DEBUG) | 1       |
| `routes`                                   | List of routes                                           | []      |
| `routes[*].name`                           | Name of the route. The name is used in log messages      |         |
| `routes[*].receiver`                       | Configuration for the receiving endpoint of a route      |         |
| `routes[*].receiver.endpoint`              | Endpoint the receiver should listen on                   |         |
| `routes[*].receiver.method`                | The HTTP method for the receiver                         |         |
| `routes[*].sender`                         | List of senders                                          | []      |
| `routes[*].sender[*].endpoint`             | The endpoint the request should be sent to               |         |
| `routes[*].sender[*].method`               | The HTTP method for the request                          |         |
| `routes[*].sender[*].body`                 | The body of that request                                 |         |
| `routes[*].sender[*].header`               | A map of HTTP headers                                    | {}      |
| `routes[*].sender[*].client`               | Configuration for the client used to send the request    |         |
| `routes[*].sender[*].client.timeout`       | Timeout for the client                                   | 10s     |
| `routes[*].sender[*].client.tlsSkipVerify` | Whether TLS verification should be skipped               |         |
| `routes[*].sender[*].conditions`           | List of conditions for that sender                       | []      |
| `routes[*].conditions`                     | List of conditions for that route                        | []      |


The following config configures Crossover to handle a route that opens a new GitLab Issue every time it receives a request on the `/gitlab` endpoint.

```yaml
port: 9999 # port on which Crossover listens
logLevel: 3 # log level
envs: # list of environment variable names that should be used for templating
  - TOKEN
  - PROJECT_ID
routes: # list of routes that Crossover should handle
  - receiver:
      endpoint: gitlab # endpoint on which Crossover listens for this route
      method: POST # HTTP method for which Crossover listens for this route
    sender:
      - endpoint: https://gitlab.com/api/v4/projects/{{ .env.PROJECT_ID }}/issues # endpoint to which a http request will be send
        method: POST # HTTP method
        header:
          Content-Type: application/json
          Authorization: "Bearer {{ .env.TOKEN }}"
        body: |
          {
            "title": "[Grafana Alert] {{ .body.title }}",
            "description": "Number of alerts: {{ len .body.alerts}}\nMessage: {{ .body.message}}"
          }
    conditions:
      - '{{ eq .body.status "firing" }}'
```

## Templating

Crossover leverages Go's `text/tempalte` module so the familiar go template syntax can be used to template the `body`, `endpoint` and `header` fields.
`Environment variables`, `URL values` and `body values` can all be used as variables.

## Conditions

It is possible to define conditions for a route as well as for individual senders, which must be fulfilled before an HTTP request is sent from Crossover. If one of the conditions fails, no request is sent.
To define conditions, the same variables can be used as for templating the `body`, `endpoint` and `header` fields.

> [!TIP]
> Use multiple conditions for logical conjuction and the `or` function within a condition for logical disjunction

