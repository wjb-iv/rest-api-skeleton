# REST API Template App
Another learning project in Go.

Inspired by some great ideas that have come up at my day job, this sample demonstrates the use of the excellent [HTTPRouter Project](https://github.com/julienschmidt/httprouter) to produce a simple, yet sophisticated and feature-rich API server.

Also uses the [Cobra CLI Library](https://github.com/spf13/cobra) to demonstrate command-line server startup with various configuration flags.

## What's Inside

### Apache Style Request logging
While you can find middleware for doing this from multiple sources, a simple demonstration of "rolling your own" is included.

### Basic Auth Demonstration
Simple basic auth is described using in-built methods of Go's net/http package.

### API Access Token (JWT) Demonstration
Methods for creating and parsing access tokens using
[dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) are shown; as well as one technique for protecting individual resources.
