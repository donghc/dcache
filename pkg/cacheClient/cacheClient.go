package cacheClient

// Cmd : command submitted by client to run on cache server
type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

// Client : client of the benchmark
type Client interface {
	Run(*Cmd)
}

// New : create new client by typ
func New(typ, server string) Client {
	if typ == "tcp" {
		return newTCPClient(server)
	}
	// Other type are not supported in my version
	panic("unknown client type: " + typ)
}
