package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-user-profile/controllers"
	"github.com/nagymarci/stock-user-profile/model"

	"github.com/nagymarci/stock-user-profile/database"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbConnectionURI string
var db *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(time.Minute * 2),
		Env:          map[string]string{},
	}
	req.Env["MONGO_INITDB_ROOT_USERNAME"] = "mongodb"
	req.Env["MONGO_INITDB_ROOT_PASSWORD"] = "mongodb"
	req.Env["MONGO_INITDB_DATABASE"] = "stock-screener"

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalln(err)
	}
	defer mongoC.Terminate(ctx)
	ip, err := mongoC.Host(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		log.Fatalln(err)
	}

	dbConnectionURI = fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		"mongodb",
		"mongodb",
		ip,
		port.Int())

	db = database.New(dbConnectionURI)

	code := m.Run()

	mongoC.Terminate(ctx)

	os.Exit(code)
}

func TestSomething(t *testing.T) {
	t.Run("simply test something", func(t *testing.T) {
		upDb := database.NewUserProfile(db)
		upC := controllers.NewUserprofileController(upDb)

		testProfile := model.Userprofile{
			UserID:         "userId",
			Email:          "alic@example.com",
			ExpectedReturn: 9,
			Expectations: []model.Expectation{
				model.Expectation{
					Stock:         "INTC",
					ExpectedRaise: 5.5,
				},
			},
		}

		upDb.Save(testProfile)

		router := mux.NewRouter().PathPrefix("/userprofile").Subrouter()
		UserprofileGetHandler(router, upC, func(r *http.Request) string { return "userId" })

		req := httptest.NewRequest(http.MethodGet, "/userprofile/userId", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

	})
}
