package main

import (
	"embed"
	"fmt"
	"github.com/jantytgat/corelogic/internal/controllers"
	"log"
)

//go:embed assets
var assets embed.FS

func main() {

	fmt.Println("Loading files")
	yamlController := controllers.YamlController{Assets: assets}
	frameworkController, err := yamlController.Load("0.1.0")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Get Output:")
	output1, err := frameworkController.GetOutput("0.1.0", "install", []string{})
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range output1 {
		fmt.Println(line)
	}
}
