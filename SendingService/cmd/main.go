package cmd

import "SendingService/internal/config"

const (
	envLocal = "local"
	envDev   = "dev"
	endProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

}
