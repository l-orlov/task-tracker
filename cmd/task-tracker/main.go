package main

import (
	"log"
	"net/http"

	"github.com/LevOrlov5404/task-tracker/internal/controller"
	"github.com/LevOrlov5404/task-tracker/internal/router"
)

func main() {
	c := controller.NewController()
	r := router.NewRouter(c)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err.Error())
	}
}
