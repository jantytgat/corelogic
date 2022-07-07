package main

import (
	"fmt"
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
}
