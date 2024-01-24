![Crossover Logo](./assets/header_dark.svg#gh-dark-mode-only)
![Crossover Logo](./assets/header_light.svg#gh-light-mode-only)

Crossover is a lightweight and declarative webhook translator written in Go🩵
So basically it receives HTTP-Requests on pre-configured routes and sends out HTTP-Request to different also configurable routes.
This project focuses on minimalism and simplicity.

![Crossover Architecture Overview](./assets/overview_light.svg#gh-light-mode-only)
![Crossover Architecture Overview](./assets/overview_dark.svg#gh-dark-mode-only)

## Usage

### Binary

```bash
./crossover -config "./config.yml"
```

### Docker

```bash
docker run -p 9999:9999 LukeTimeWalker/crossover
```

## Config

| Key | Description | Default |
| --- | ----------- | ------- |
| port | Specifies the port crossover is listening | 9999 |
| timeout | Sets the timeout for the specified routes | 15 |
| logLevel | LogLevel (0=Debug, 1=Info, 2=Warn, 3=Error) | 1 |
| envs | List of environment variables crossover should read | [] |
| routes | List of specified routes | [] |
| routes[].receiver | Route that listens for incoming requests | |
| routes[].receiver.path |||
| routes[].receiver.method |||
| routes[].receiver.urlParams |  | [] |
| routes[].sender |  |  |
| routes[].sender[].url |  |  |
| routes[].sender[].method |  |  |
| routes[].sender[].header |  |  |
| routes[].sender[].body |  |  |
| routes[].sender[].conditions |  |  |
