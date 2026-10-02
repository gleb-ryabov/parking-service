package main

import (
	"parking/internal/config"
)

func main() {
	_ = config.MustLoad()
}
