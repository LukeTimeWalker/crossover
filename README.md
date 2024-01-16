![Crossover Logo](./assets/header_dark.svg#gh-dark-mode-only)
![Crossover Logo](./assets/header_light.svg#gh-light-mode-only)

Crossover is a declarative webhook translator written in Go. So basically it receives HTTP-Requests on preconfigured endpoints and sends out HTTP-Request to other endpoints based on conditions

You can receive Variables from the Request Body, Environment Variables & Url Parameter

You can template these Variables into Header, Body & Url

Crossover utilizes the golang template package and syntax

### Logging

Crossover knows 4 log levels (Debug, Info, Warn, Error).

| Level | Name | Description |
| ----- | ---- | ----------- |
| 0     | Error | Error just logs error notifications and causes the programm to crash |
| 1     | Warn | Warnings are errors that do not cause a crash of the application |
| 2 | Info provides usefull information about the current state of the appliation |
| 3 | Debug is usefull for debugging an error, as it provides a lot of information |