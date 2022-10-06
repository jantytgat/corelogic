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
	versions, err := yamlController.ListAvailableVersions()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range versions {
		fmt.Println(v)
	}
	frameworkController, err := yamlController.Load("0.1.8")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Get Output:")
	frameworkController.Parse()
	output1, err := frameworkController.GetOutput("0.1.8", "install", []string{})
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range output1 {
		fmt.Println(line)
	}
}
