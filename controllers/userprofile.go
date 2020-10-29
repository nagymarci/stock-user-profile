package controllers

import (
	log "github.com/sirupsen/logrus"

	"github.com/nagymarci/stock-user-profile/database"
	"github.com/nagymarci/stock-user-profile/model"
)

//UserprofileController represents a database connection to userprofile collection
type UserprofileController struct {
	userprofiles database.UserProfilesCollection
}

//NewUserprofileController creates a controller with the given db collection
func NewUserprofileController(u database.UserProfilesCollection) *UserprofileController {
	return &UserprofileController{
		userprofiles: u,
	}
}

//Create creates or updates
func (u *UserprofileController) Create(up model.Userprofile) error {
	err := u.userprofiles.Save(up)

	if err != nil {
		return model.NewInternalServerError(err.Error())
	}

	return nil
}

//Get returns the userprofile for the given ID
func (u *UserprofileController) Get(userID string) (model.Userprofile, error) {
	userprofile, err := u.userprofiles.Get(userID)

	if err != nil {
		message := "Cannot read userprofile " + err.Error()
		log.WithFields(log.Fields{"userId": userID}).Error(err)
		return model.Userprofile{}, model.NewBadRequestError(message)
	}

	return userprofile, nil
}
