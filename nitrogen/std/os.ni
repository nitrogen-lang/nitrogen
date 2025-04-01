const env = fn native ()
const argv = fn native ()
const exec = fn native (cmd, args)
const system = fn native (cmd, args)

return {
	"env": env,
	"argv": argv,
	"exec": exec,
	"system": system,
}
