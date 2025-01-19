package main

import (
	"fmt"
	"log"

	"github.com/englandrecoil/go-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading configuration")
	}

	fmt.Println(cfg)
}
