# Nitrogen SCGI Webserver Example

Requirements:

- Docker
- Docker Compose

This example demonstrates Nitrogen running as an SCGI server. The Nitrogen executable
has a builtin SCGI server and doesn't need any extra software to handle requests.

## Setup

Build Nitrogen using the make file at the root of the project. Simply run `make`.

## Running

While in the `webserver-scgi` directory, run `docker-compose up`. This will start
two containers. One is for Nginx which exposes the server on port 8080. The second
is an app container running Nitrogen in SCGI mode.

## Testing

Once everything is up, open a web browser and go to `http://localhost:8081/cgi/index.ni`.
You should see the environment variables that are available to the script.
