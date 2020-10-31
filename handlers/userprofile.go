package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/nagymarci/stock-user-profile/model"

	"github.com/gorilla/mux"
	"github.com/nagymarci/stock-user-profile/controllers"
)

func UserprofileCreateHandler(router *mux.Router, userprofileController *controllers.UserprofileController, extractUserIDFromToken func(*http.Request) string) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserIDFromToken(r)

		var userprofile model.Userprofile

		err := json.NewDecoder(r.Body).Decode(&userprofile)

		if err != nil {
			message := "Failed to deserialize payload."
			handleErrorResponse(message, w, http.StatusBadRequest)
			log.WithFields(log.Fields{"userId": userID}).Error(err)
			return
		}

		if userID != userprofile.UserID {
			message := "UserID in request doesn't match userID in token"
			handleErrorResponse(message, w, http.StatusUnauthorized)
			log.WithFields(log.Fields{"userId": userID, "request_userId": userprofile.UserID}).Error("Unauthorized")
			return
		}

		err = validateFields(userprofile)

		if err != nil {
			handleErrorResponse(err.Error(), w, http.StatusBadRequest)
			log.WithFields(log.Fields{"userId": userID}).Error(err)
			return
		}

		err = userprofileController.Create(userprofile)

		if err != nil {
			handleError(err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}).Methods(http.MethodPost, http.MethodOptions)
}

func UserprofileGetHandler(router *mux.Router, userprofileController *controllers.UserprofileController, extractUserIDFromToken func(*http.Request) string) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserIDFromToken(r)
		id := mux.Vars(r)["id"]

		if userID != id {
			message := "UserID in request doesn't match userID in token"
			handleErrorResponse(message, w, http.StatusUnauthorized)
			log.WithFields(log.Fields{"userId": userID, "request_userId": id}).Error("Unauthorized")
			return
		}

		result, err := userprofileController.Get(id)

		if err != nil {
			handleError(err, w)
			return
		}

		handleJSONResponse(result, w, http.StatusOK)
	}).Methods(http.MethodGet)
}

func validateFields(up model.Userprofile) error {
	if isEmpty(up.UserID) {
		return errors.New("Field \"userId\" is missing")
	}

	if isEmpty(up.Email) {
		return errors.New("Field \"email\" is missing")
	}

	if up.ExpectedReturn == nil {
		return errors.New("Field \"expectedReturn\" is missing")
	}

	if up.DefaultExpectation == nil {
		return errors.New("Field \"defaultExpectation\" is missing")
	}

	for _, exp := range up.Expectations {
		if isEmpty(exp.Stock) {
			return errors.New("Field \"stock\" from \"expectations\" is missing")
		}

		if exp.ExpectedRaise == nil {
			return errors.New("Field \"expectedRaise\" from \"expectations\" is missing")
		}
	}

	return nil
}

func handleError(err error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	if err, ok := interface{}(&err).(model.HttpError); ok {
		statusCode = err.Status()
	}
	message := "Failed to process request: " + err.Error()
	handleErrorResponse(message, w, statusCode)
	log.Println(message)
}

func handleErrorResponse(msg string, w http.ResponseWriter, status int) {
	response := model.ErrorResponse{Message: msg}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, model.UnknownError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)
}

func handleJSONResponse(object interface{}, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(object)
}

func DefaultExtractUserID(r *http.Request) string {
	user := r.Context().Value("user")
	email := user.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	return email
}

func isEmpty(str string) bool {
	return len(str) == 0 || str == " "
}
