# Woodpecker CI sample extensions

This repo contains sample extensions for [Woodpecker CI](https://woodpecker-ci.org/).

## Run them locally

```shell
go run ./config

# or

go run ./secrets
```


## Sample configuration extension

This repository provides a very simplistic example of how to set up an
external configuration service for **Woodpecker CI** The external service
gets a **HTTP POST** request with information about the repo, current
build, and the configs that would normally be used. It can then decide to
acknowledge the current configs (By returning **HTTP 204**), or overriding
the configurations and returning new ones in the response

Use cases for this system are:

- Centralized configuration for multiple repositories at once
- Preprocessing steps in the pipeline like templating, macros or conversion from
  different pipeline formats to woodpeckers format

This service is written in go, to run it first copy the config example:
`cp .env.example .env`
Download the public key from your woodpecker instance from
`http(s)://your-woodpecker-server/api/signature/public-key` and save it
to file. Set `WOODPECKER_CONFIG_SERVICE_PUBLIC_KEY_FILE` to the path to
that file and add a filtering regex. The repositories that have a name
match the filtering regex will receive the config from `central-pipeline-config.yaml`,
while all other repositories will continue using their original configuration.

Then run using `go run .`.

Make sure to configure your woodpecker instance with the correct **endpoint** and
configure the same **secret**. See [Woodpeckers documentation here](https://woodpecker-ci.org/docs/administration/external-configuration-api)

eg:

```shell
# Server
# ...
WOODPECKER_CONFIG_SERVICE_ENDPOINT=http://<service>:8000/ciconfig
WOODPECKER_CONFIG_SERVICE_PUBLIC_KEY_FILE=public-key.pem
```

## Sample secrets extension
