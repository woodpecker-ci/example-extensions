# Woodpecker CI sample extensions

This repo contains sample extensions for [Woodpecker CI](https://woodpecker-ci.org/).

## Run locally

Copy `.env.example` to `.env` and adjust the values to your needs.

Then run the following command to start the server: `go run .`

Finally go to your repository settings, to the extensions tab and set
the URL of this server as the URL for the extension you want to test.

## Security

It is extremely important to secure your extension server. Only Woodpecker CI
should be able to send requests to your endpoints and no one should be able to
intercept or modify the requests in transit. To ensure this, Woodpecker CI
will sign the requests by providing a http signature. You can verify the signature using the
public key of your Woodpecker CI instance which you can find in the extension settings
of your repository or by accessing `<woodpecker-ci-url>/api/signature/public-key`.

## Background

This repository provides a very simplistic example of how to set up an
external server providing a config and secrets extensions for **Woodpecker CI**.

This service can be registered in the settings of a repository in Woodpecker CI. When
configured, Woodpecker CI will send a **HTTP** requests to predefined endpoints
on this server.

### Secrets extension

For example if you open the secrets list in the repository settings, the following will happen:

- Woodpecker CI will send a **HTTP GET** request to the `<this-server-url>/repo/<repo-id>/secrets`
  endpoint of this extension server
- The extension server will first verify the request is coming from the actual Woodpecker CI instance
  and wasn't modified in transit using a http signature. This is pretty important so no one can just
  send a request to this server and get all your secrets for example.
- The extension server will then return a list of secrets that are available for this repository.
- The secrets will be displayed in the UI.

### Configuration extension

The configuration extension allows you to provide a custom pipeline configuration
or modify the existing one. It will be called when the pipeline configuration is
compiled and before a new pipeline is executed.

Using such an extension can be useful if you want to:

- Implement custom pipeline config options
- Do something like templating or macro expansion
- Convert from a different pipeline format to the Woodpecker format
  - Think of a Gitlab CI to Woodpecker CI converter
  - Or a Jsonnet or starlark converter
- Centralized configuration for multiple repositories at once
- Or just some other special use case

The endpoint will receive a **HTTP POST** request with the following JSON data:

- `repo` - The repository object
- `pipeline` - The current pipeline configuration
- `netrc` - The netrc object

> [!TIP]
> The `netrc` data is pretty powerful as it contains credentials to access the repository.
> You can use this to clone the repository or even use the forge (Github or Gitlab, ...) api
> to get more information about the repository.

The endpoint can return a **HTTP 204** to tell Woodpecker CI to stick to the current configuration
or return a new configuration in the response with a **HTTP 200**.
