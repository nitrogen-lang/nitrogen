#!/usr/local/bin/nitrogen

import 'std/string'
import 'std/collections'
import 'std/os'

const printMap = func(map) {
    const list = collections.reduce(map, func(acc, val, key) {
        acc + string.format('<li>{}: {}</li>', key, val)
    }, '<ul>')

    return list + '</ul>'
}

println(string.format('Content-Type: text/html

<!DOCTYPE html>
<html>
<head>
    <title>Nitrogen Webpage CGI Example</title>
</head>
<body>
    <h2>Hello from Nitrogen! CGI</h2>
    <h3>Script Environment:</h3>
    {}
    {}
</body>
</html>
', printMap(os.env()), printMap(_SERVER)))
