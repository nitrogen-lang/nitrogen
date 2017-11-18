# Nitrogen CGI Webserver Example

Requirements:

- Docker
- Docker Compose

This example demonstrates running Nitrogen scripts using CGI. fcgiwrap is used to "conver" FastCGI requests from Nginx
into CGI requests for the script.

## Setup

This example requires a container with fcgiwrap installed. The provided Dockerfile will build an Ubuntu image with
fcgiwrap. If you make your own image, make sure it has glibc if you want to use plugins with the Nitrogen interpreter.
The example index file expects the interpreter to live at `/usr/local/bin/nitrogen`.

Build Nitrogen using the make file at the root of the project. Simply run `make`.

## Running

While in the `webserver-cgi` directory, run `docker-compose up`. This will start two containers. One is for Nginx which exposes
the server on port 8080. The second is an app container running fcgiwrap. The second container should have the Nitrogen
binary. The Nginx container doesn't need Nitrogen.

## Testing

Once everything is up, open a web browser and go to `http://localhost:8080/cgi/index.ni`. You should see the environment
variables that are available to the script.
