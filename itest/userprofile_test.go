package itest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-user-profile/controllers"
	"github.com/nagymarci/stock-user-profile/handlers"
	"github.com/nagymarci/stock-user-profile/model"

	"github.com/nagymarci/stock-user-profile/database"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	dbConnectionURI := fmt.Sprintf(
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

func TestUserprofileGetHandler(t *testing.T) {
	t.Run("sends 200OK with data from db", func(t *testing.T) {
		upDb := database.NewUserProfile(db)
		upC := controllers.NewUserprofileController(upDb)

		expectedReturn := 9.0
		defaultExpectation := 9.0
		expectedRaise := 5.5

		testProfile := model.Userprofile{
			UserID:         "userId",
			Email:          "alic@example.com",
			ExpectedReturn: &expectedReturn,
			Expectations: []model.Expectation{
				model.Expectation{
					Stock:         "INTC",
					ExpectedRaise: &expectedRaise,
				},
			},
			DefaultExpectation: &defaultExpectation,
		}

		upDb.Save(testProfile)

		router := mux.NewRouter().PathPrefix("/userprofile").Subrouter()
		handlers.UserprofileGetHandler(router, upC, func(r *http.Request) string { return "userId" })

		req := httptest.NewRequest(http.MethodGet, "/userprofile/userId", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		var result model.Userprofile
		json.NewDecoder(res.Body).Decode(&result)

		assertEquals(t, &testProfile, &result)

	})
	t.Run("sends 404 when userprofile missing", func(t *testing.T) {
		upDb := database.NewUserProfile(db)
		upC := controllers.NewUserprofileController(upDb)

		router := mux.NewRouter().PathPrefix("/userprofile").Subrouter()
		handlers.UserprofileGetHandler(router, upC, func(r *http.Request) string { return "userId2" })

		req := httptest.NewRequest(http.MethodGet, "/userprofile/userId2", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected [%d], got [%d]", http.StatusNotFound, res.StatusCode)
		}
	})
}

func TestUserprofileCreateHandler(t *testing.T) {
	t.Run("sends 201OK and data in db", func(t *testing.T) {
		upDb := database.NewUserProfile(db)
		upC := controllers.NewUserprofileController(upDb)

		expectedReturn := 9.0
		defaultExpectation := 9.0
		expectedRaise := 5.5

		testProfile := model.Userprofile{
			UserID:         "userId",
			Email:          "alic@example.com",
			ExpectedReturn: &expectedReturn,
			Expectations: []model.Expectation{
				model.Expectation{
					Stock:         "INTC",
					ExpectedRaise: &expectedRaise,
				},
			},
			DefaultExpectation: &defaultExpectation,
		}

		router := mux.NewRouter().PathPrefix("/userprofile").Subrouter()
		handlers.UserprofileCreateHandler(router, upC, func(r *http.Request) string { return "userId" })

		body, _ := json.Marshal(testProfile)

		req := httptest.NewRequest(http.MethodPost, "/userprofile", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusCreated {
			t.Fatalf("expected [%d], got [%d]", http.StatusOK, res.StatusCode)
		}

		result, _ := upDb.Get("userId")

		assertEquals(t, &testProfile, &result)

	})
}

func assertEquals(t *testing.T, expected *model.Userprofile, got *model.Userprofile) {
	if expected.UserID != got.UserID {
		t.Fatalf("expected [%s], got [%s]", expected.UserID, got.UserID)
	}

	if expected.Email != got.Email {
		t.Fatalf("expected [%s], got [%s]", expected.Email, got.Email)
	}

	if *expected.ExpectedReturn != *got.ExpectedReturn {
		t.Fatalf("expected [%f], got [%f]", *expected.ExpectedReturn, *got.ExpectedReturn)
	}

	if len(expected.Expectations) != len(got.Expectations) {
		t.Fatalf("expected [%v], got [%v]", expected.Expectations, got.Expectations)
	}
}
