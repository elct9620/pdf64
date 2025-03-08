package main

import (
	"net/http"

	"github.com/elct9620/pdf64/internal/app"
	v1 "github.com/elct9620/pdf64/internal/controller/v1"
)

func main() {
	apiV1 := v1.NewService()
	server := app.NewServer(apiV1)

	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
