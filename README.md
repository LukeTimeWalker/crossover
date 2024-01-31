![Crossover Logo](./assets/header_dark.svg#gh-dark-mode-only)
![Crossover Logo](./assets/header_light.svg#gh-light-mode-only)

Crossover is a lightweight and declarative webhook translator written in Go🩵  
So basically it receives HTTP-Requests on pre-configured routes and sends out HTTP-Request to different also configurable routes.  
This project focuses on `minimalism` and `simplicity`.

![Crossover Architecture Overview](./assets/overview_light.svg#gh-light-mode-only)
![Crossover Architecture Overview](./assets/overview_dark.svg#gh-dark-mode-only)

## Installation

To install crossover, you can either simply download the corresponding binary or use the docker image.

## Usage

Crossover can be configured via a `config.yml` file, which can be specified via the `-config` flag. If no config file is specified, the system searches for `./config.yml` in the same directory as crossover by default.

### Binary

```bash
./crossover -config "./config.yml"
```

### Docker/Podman

```bash
docker run -v ./config.yml:/crossover/config.yml -p 9999:9999 crossover
```

## Config

## Templating

Crossover leverages the `text/tempalte` module so the familiar go template syntax can be used to template `header`, `body` and `url`.

The variables for templating can be 


## Conditions

It is possible to define conditions for a route as well as for individual sender, which must be fulfilled before an HTTP request is sent from Crossover. If one of the conditions fails, no request is sent.
To define conditions, the same variables can be accessed as for templating the body, URL and header

