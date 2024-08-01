package main

import (
	"fmt"

	ep "github.com/lanseg/homebox/endpoint"
)

func main() {
	fmt.Println("Starting data server")
	ep.NewEndpoint(12345)
}
