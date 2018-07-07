#!/usr/local/bin/nitrogen

func printMap(map) {
    println('<ul>')

    const keys = sort(hashKeys(map))
    const keyLen = len(keys)

    for i = 0; i < keyLen; i += 1 {
        let key = keys[i]
        println('<li>', key, ': ', map[key], '</li>')
    }

    println(res, '</ul>')
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
