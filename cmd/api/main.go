package main

import (
	"fmt"
	"time"

	"github.com/icrxz/crm-api-core/internal/infra"
)

func main() {
	time.Local = time.UTC
	fmt.Println("Starting crm-core app!!")

	if err := infra.RunApp(); err != nil {
		panic(err)
	}
}
