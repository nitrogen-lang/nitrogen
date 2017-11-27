#!/usr/local/bin/nitrogen
print("Content-Type: text/html\n\n") // HTTP header section ends with two newlines

println('<!DOCTYPE html>
<html>
<head>
    <title>Nitrogen Webpage Example</title>
</head>
<body>')
println('<h2>Hello from Nitrogen!</h2>')

println('<h3>Script Environment:</h3>')
println('<ul>')

let keys = sort(hashKeys(_ENV))
for i = 0; i < len(keys); i += 1 {
    println('<li>', keys[i], ': ', _ENV[keys[i]], '</li>')
}

println('</ul>
</body>
</html>')
