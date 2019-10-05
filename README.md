[![Build Status](https://travis-ci.org/KCarretto/paragon.svg?branch=master)](https://travis-ci.org/KCarretto/paragon) [![Go Report Card](https://goreportcard.com/badge/github.com/kcarretto/paragon)](https://goreportcard.com/report/github.com/kcarretto/paragon) [![Coverage Status](https://coveralls.io/repos/github/KCarretto/paragon/badge.svg?branch=master)](https://coveralls.io/github/KCarretto/paragon?branch=master) ![GitHub release](https://img.shields.io/github/release-pre/kcarretto/paragon.svg) [![GoDoc](https://godoc.org/github.com/KCarretto/paragon?status.svg)](https://godoc.org/github.com/KCarretto/paragon)
# Paragon
![Logo](.github/images/logo.png)

## Installation
Ensure you have docker installed on your system. If you clone the repository, a `docker-compose.yml` file has been provided to enable easier setup. Running `docker run -e DEBUG_HTTP_ADDR=0.0.0.0:8080 -p 127.0.0.1:8080:8080 kcarretto/paragon:latest` will pull the latest agent image which will run with a debug HTTP transport, accessible at http://127.0.0.1:8080 from your browser.
