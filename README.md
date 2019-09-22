curl -XPOST -d'{"s":"hello, world"}' localhost:8080/import
{"v":"HELLO, WORLD","err":null}
curl -XPOST -d'{"s":"hello, world"}' localhost:8080/export
{"v":12}