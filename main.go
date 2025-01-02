package main

import (
	"github.com/per1Peteia/gator/internal/config"
	"fmt"
	"log"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	config.SetUser("justus")
	
	config, err = cfg.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("%s\n%s\n", config.DbURL, config.CurrentUserName)
	
}
