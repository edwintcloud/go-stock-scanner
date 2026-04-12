package main

import (
	"github.com/edwintcloud/go-stock-scanner/cmd"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
