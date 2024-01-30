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


### Templating

Crossover leverages the `text/tempalte` module. So you can use every feature of this module.
