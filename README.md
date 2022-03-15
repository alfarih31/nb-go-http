# Noob

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=alfarih31_nb-go-http&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=alfarih31_nb-go-http)
[![Go Reference](https://pkg.go.dev/badge/github.com/alfarih31/nb-go-http.svg)](https://pkg.go.dev/github.com/alfarih31/nb-go-http)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/alfarih31/nb-go-http?style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/alfarih31/nb-go-http?style=flat-square)

Noob is a REST API framework for faster development by providing several common functions, such as:
- Query Parsing
- Logging
- Try-Catch-Finally block
- Env loader
- Response Templating

**Currently, Noob based on [Gin v1.7.4](https://github.com/gin-gonic/gin). Thanks to [*gin-gonic*](https://github.com/gin-gonic/gin)**

## Contents

- [NOOB](#noob)
  - [TODO](#todo)

## TODO
- [ ] Write unit test

## Installation

To install this package, you need to install Go (**version 1.17+ is required**) & initiate your Go workspace first.

1. After you initiate your workspace then you can install this package with below command.

```shell
go get -u github.com/alfarih31/nb-go-http
```

2. Import it in your code

```go
import "github.com/alfarih31/nb-go-http"
```

## Quick Start & Usage

See the example: [sample_app](examples/sample_app.go)

## Contributors ##

- Alfarih Faza <alfarihfz@gmail.com>

## License

This project is licensed under the - see the [LICENSE.md](LICENSE.md) file for details