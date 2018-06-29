#!/usr/local/bin/nitrogen

func printMap(map) {
    println('<ul>')

    let keys = sort(hashKeys(map))
    for i = 0; i < len(keys); i += 1 {
        println('<li>', keys[i], ': ', map[keys[i]], '</li>')
    }

    println('</ul>')
}

print("Content-Type: text/html\n")
print("\n") // HTTP header section ends with an empty line

println('<!DOCTYPE html>
<html>
<head>
    <title>Nitrogen Webpage SCGI Example</title>
</head>
<body>')
println('<h2>Hello from Nitrogen! SCGI</h2>')

println('<h3>Script Environment:</h3>')

printMap(_ENV)
printMap(_SERVER)

println('</body>
</html>')
