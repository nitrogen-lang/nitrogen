# Nitrogen SCGI Server

Like any other executable code, Nitrogen scripts can be used with CGI to generate dynamic web content. The Nitrogen interpreter
goes one step further and supports Simple CGI or SCGI. SCGI allows the interpreter to always run and accepts multiple requests
from another web server. Nitrogen does not currently support running as a stand alone web server. It needs something like
Nginx or Apache in front if it to proxy requests. Using SCGI is more efficient since the operating system doesn't need to
create a new process for every request. Instead, all requests are handled in a single process.

Nitrogen supports concurrent requests meaning it can fulfill multiple proxy requests at the same time. The number of workers is
configurable. By default, 5 workers are created to handle requests.

## Flags

These are the flag given to Nitrogen to configure SCGI.

- `-scgi`: Enable the SCGI server. With this flag, no script is executed and the interactive prompt is not shown. Without this,
none of the remaining flags have any effect.
- `-scgi-sock`: TCP or Unix socket to listen on. Ex: `tcp:127.0.0.1:9000` or `unix:/var/run/nitrogen-scgi.sock`. This defaults
to `tcp:0.0.0.0:9000`.
- `-scgi-workers`: The number of workers available to handle requests. Defaults to 5.
- `-scgi-worker-timeout`: The number of seconds the server will wait for an available worker. If all workers are busy, the server
will wait this long before closing the connection. If this timeout is reached, an error message will be printed to standard output
saying there weren't enough workers to handle incoming requests. You can use this to adjust the number of workers available.

## Scripts

The only change to normal script execution is any print statements will go to the client's browser, not the process's normal
standard output. The `_ENV` variable will contain any CGI variables provided by the upstream web server.
