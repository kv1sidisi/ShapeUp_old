package main

import (
	"GatewayAPI/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("%+v\n", cfg)
}
