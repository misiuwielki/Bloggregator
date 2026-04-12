package main

import (
	"fmt"

	"github.com/misiuwielki/Bloggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	cfg.SetUser("misiuwielki")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)

}
