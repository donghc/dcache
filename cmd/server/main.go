package main

import (
	"github.com/donghc/dcache/pkg/cache"
	"github.com/donghc/dcache/pkg/server"
)

func main() {
	c := cache.New("inMemoryCache")
	server.New(c).Listen()
}
