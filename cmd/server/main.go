package main

import (
	"github.com/donghc/dcache/pkg/cache"
	"github.com/donghc/dcache/pkg/http"
	"github.com/donghc/dcache/pkg/tcp"
)

func main() {
	c := cache.New("inMemoryCache")
	go tcp.New(c).Listen()
	http.New(c).Listen()
}
