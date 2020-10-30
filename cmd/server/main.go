package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nagymarci/stock-user-profile/controllers"
	"github.com/nagymarci/stock-user-profile/routes"

	"github.com/nagymarci/stock-user-profile/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	db := database.New(os.Getenv("DB_CONNECTION_URI"))
	uDb := database.NewUserProfile(db)

	uC := controllers.NewUserprofileController(uDb)

	router := routes.Route(uC)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
