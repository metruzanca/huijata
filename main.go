package main

import (
	"github.com/joho/godotenv"
	"github.com/metruzanca/huijata/cmd"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
