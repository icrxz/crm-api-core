package main

import (
	"fmt"

	"github.com/icrxz/crm-api-core/internal/infra"
)

func main() {
	fmt.Println("Starting crm-core app!!")

	if err := infra.RunApp(); err != nil {
		panic(err)
	}
}
