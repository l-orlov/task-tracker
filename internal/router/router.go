package router

import (
	"net/http"
	"github.com/gorilla/mux"

	"github.com/LevOrlov5404/task-tracker/internal/controller"
)

const (
	apiV1Prefix = "/api/v1"
)

// NewRouter return new created and configured *mux.Router
func NewRouter(c controller.Controller) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Methods(http.MethodGet).
		Name("HomePage").
		PathPrefix(apiV1Prefix).
		Path("/").
		HandlerFunc(c.HomePage)

	r.Methods(http.MethodGet).
		Name("ContactsPage").
		PathPrefix(apiV1Prefix).
		Path("/contacts/").
		HandlerFunc(c.ContactsPage)

	return r
}
