package controller

import (
	"fmt"
	"net/http"
)

// Controller is responsible for http handlers
type Controller interface {
	HomePage(w http.ResponseWriter, r *http.Request)
	ContactsPage(w http.ResponseWriter, r *http.Request)
}

type controller struct {}

// NewController return new created Controller
func NewController() Controller {
	return &controller{}
}

func (c *controller) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
}

func (c *controller) ContactsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contacts page")
}