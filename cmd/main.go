package main

import (
	"net/http"

	"github.com/elct9620/pdf64/internal/app"
	"github.com/elct9620/pdf64/internal/builder"
	v1 "github.com/elct9620/pdf64/internal/controller/v1"
	"github.com/elct9620/pdf64/internal/usecase"
)

func main() {
	// Initialize dependencies
	fileBuilder := builder.NewFileBuilder()
	convertUsecase := usecase.NewConvertUsecase(fileBuilder)
	
	// Initialize controllers
	apiV1 := v1.NewService(convertUsecase)
	
	// Initialize server
	server := app.NewServer(apiV1)

	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
