package routes

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/nagymarci/stock-commons/authorization"

	"github.com/nagymarci/stock-user-profile/controllers"
	"github.com/nagymarci/stock-user-profile/handlers"
)

func Route(userprofileController *controllers.UserprofileController) (http.Handler, http.Handler) {
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	userprofile := mux.NewRouter().PathPrefix("/userprofile").Subrouter()
	handlers.UserprofileCreateHandler(userprofile, userprofileController, handlers.DefaultExtractUserID)
	handlers.UserprofileGetHandler(userprofile, userprofileController, handlers.DefaultExtractUserID)

	audience := os.Getenv("USERPROFILE_AUDIENCE")
	authServer := os.Getenv("AUTHORIZATION_SERVER")
	userprofileScope := os.Getenv("USERPROFILE_SCOPE")

	auth := negroni.New(
		negroni.HandlerFunc(authorization.CreateAuthorizationMiddleware(audience, authServer).HandlerWithNext),
		negroni.HandlerFunc(authorization.CreateScopeMiddleware(userprofileScope, authServer, audience)))

	router.PathPrefix("/userprofile").Handler(auth.With(negroni.Wrap(userprofile)))

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false

	n := negroni.New(recovery, negroni.NewLogger())
	n.UseHandler(router)

	internalRouter := mux.NewRouter()
	internalUserProfile := mux.NewRouter().PathPrefix("/userprofile").Subrouter()
	handlers.UserprofileGetHandler(internalUserProfile, userprofileController, func(r *http.Request) string {
		return mux.Vars(r)["id"]
	})
	internalRouter.PathPrefix("/userprofile").Handler(internalUserProfile)
	internal := negroni.New(recovery, negroni.NewLogger())
	internal.UseHandler(internalRouter)

	return n, internal

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
