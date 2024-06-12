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
intercept or modify the requests in transit. To ensure this, Woodpecker CI will
sign the requests by providing a http signature. You have to verify the signature
using the public key of your Woodpecker CI instance!

## More about Woodpecker extensions

For more information about Woodpecker CI extensions, please refer to the [official documentation](https://woodpecker-ci.org/docs/usage/extensions/).
