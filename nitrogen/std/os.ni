const env = fn native ()
const argv = fn native ()
const exec = fn native (cmd, args)
const system = fn native (cmd, args)

return {
	env,
	argv,
	exec,
	system,
}
