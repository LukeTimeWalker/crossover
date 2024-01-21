![Crossover Logo](./assets/header_dark.svg#gh-dark-mode-only)
![Crossover Logo](./assets/header_light.svg#gh-light-mode-only)

Crossover is a lightweight and declarative webhook translator written in Go. So basically it receives HTTP-Requests on pre-configured endpoints and sends out HTTP-Request to different also configurable endpoints.

## Config

| Key | Description | Default |
| --- | ----------- | ------- |
| port | Specifies the port crossover is listening | 9999 |
| timeout | Sets the timeout for the specified endpoints | 15 |
| logLevel | LogLevel (0=Debug, 1=Info, 2=Warn, 3=Error) | 1 |
| envs | List of environment variables crossover should read | [] |
| endpoints | List of specified endpoints | [] |
| endpoints[].receiver | Endpoint that listens for incoming requests | |
| endpoints[].receiver.path |||
| endpoints[].receiver.method |||
| endpoints[].receiver.urlParams |  | [] |
| endpoints[].sender |  |  |
| endpoints[].sender[].url |  |  |
| endpoints[].sender[].method |  |  |
| endpoints[].sender[].header |  |  |
| endpoints[].sender[].body |  |  |
| endpoints[].sender[].conditions |  |  |

