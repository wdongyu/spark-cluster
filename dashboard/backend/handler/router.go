package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(handler *APIHandler) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	var publicRoutes = Routes{
		{
			"ListSparkCluster",
			"GET",
			"/apis/sparkcluster",
			handler.ListSparkCluster,
		},
		{
			"CreateSparkCluster",
			"POST",
			"/apis/sparkcluster",
			handler.CreateSparkCluster,
		},
		{
			"CreateTerminal",
			"GET",
			"/apis/terminal",
			handler.CreateTerminal,
		},
		{
			"UploadFile",
			"GET",
			"/apis/file",
			handler.UploadFile,
		},
	}

	// The public route is always accessible
	for _, route := range publicRoutes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// The private route is only accessible if the user has a valid access_token.
	// We are chaining the middleware into the negroni handler function which will
	// check for a valid token.

	// Handle websocket routes with path prefix router
	// WebSocket need multi routes for a service.
	router.PathPrefix("/terminal/ws").Handler(NewTerminal(handler.kubeClient, handler.kubeConfig))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})

	return c.Handler(router)
}
