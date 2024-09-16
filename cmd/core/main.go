package main

import (
	"fmt"
	"os"

	"avito2024/internal/app/core"
	"avito2024/internal/config"
)

func main() {

	host := os.Getenv("SERVER_ADDRESS")
	connStr := os.Getenv("POSTGRES_CONN")
	isTest := os.Getenv("IS_TEST_ENV")

	isT := func() bool { return isTest == "true" }()

	if isT {
		fmt.Println("server address: ", host)
		fmt.Println("postgress conn: ", connStr)
	}

	cfg := config.New(
		host,
		connStr,
		isT,
	)

	core.Run(cfg)
}
