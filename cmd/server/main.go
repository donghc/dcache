package main

import (
	"fmt"
	"github.com/donghc/dcache/pkg/cache"
)

func main() {
	c := cache.New("inMemoryCache")
	fmt.Println(c)
}
