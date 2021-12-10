package main

import (
	"github.com/donghc/dcache/pkg/cache"
	"github.com/donghc/dcache/pkg/http"
)

func main() {
	c := cache.New("inMemoryCache")
	http.New(c).Listen()
}
