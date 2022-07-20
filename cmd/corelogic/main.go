package main

import (
	"fmt"
	"github.com/jantytgat/corelogic/internal/controllers"
	"github.com/jantytgat/corelogic/internal/models"
)

func main() {
	fmt.Println("CoreLogic")

	element := models.Element{
		Name:        "Test",
		Tags:        nil,
		Fields:      nil,
		Expressions: models.Expression{},
	}

	fmt.Println(element)

	yamlController := controllers.YamlController{}
	frameworkController, err := yamlController.Load("11.0")
	if err != nil {
		panic(err)
	}
	fmt.Println(frameworkController.GetOutput("install", []string{}))
}
