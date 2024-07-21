package main

import (
	"fmt"

	ep "github.com/lanseg/homebox/endpoint"
)

func main() {
	fmt.Println("HELLO")

	ep.NewEndpoint(1234)
}
