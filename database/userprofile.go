package database

import (
	"context"

	"github.com/nagymarci/stock-user-profile/model"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserProfiles struct {
	collection *mongo.Collection
}

type UserProfilesCollection interface {
	Save(up model.Userprofile) error
	Get(userID string) (model.Userprofile, error)
}

//NewUserProfile creates a new object to manipulate the collection
func NewUserProfile(db *mongo.Database) UserProfilesCollection {
	return &UserProfiles{
		collection: db.Collection("userProfile"),
	}
}

//SaveUserProfile saves the userprofile in the parameter. If the userprofile for the given
//UserId already exists, then it updates the fields
func (u *UserProfiles) Save(up model.Userprofile) error {
	filter := bson.D{{Key: "_id", Value: up.UserID}}

	opts := options.Replace().SetUpsert(true)

	_, err := u.collection.ReplaceOne(context.TODO(), filter, up, opts)

	log.WithFields(log.Fields{"userId": up.UserID}).Error(err)
	return err
}

//Get returns the requested userprofile
func (u *UserProfiles) Get(userID string) (model.Userprofile, error) {
	filter := bson.D{{Key: "_id", Value: userID}}

	var result model.Userprofile
	err := u.collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.WithFields(log.Fields{"userId": userID}).Error(err)
	}
	return result, err
}
