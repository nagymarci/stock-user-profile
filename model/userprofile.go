package model

type Expectation struct {
	Stock         string  `bson:"stock"`
	ExpectedRaise float64 `bson:"expectedRaise"`
}

type Userprofile struct {
	UserID         string        `bson:"_id"`
	Email          string        `bson:"email"`
	Expectations   []Expectation `bson:"expectations"`
	ExpectedReturn float64       `bson:"expectedReturn"`
}

type UserprofileRequest struct {
	Expectations   []Expectation `json:"expectations"`
	ExpectedReturn float64       `json:"expectedReturn"`
}

//ToUserprofile creates a Userprofile from a request object
func (upr *UserprofileRequest) ToUserprofile(userID string, email string) Userprofile {
	return Userprofile{
		UserID:         userID,
		Email:          email,
		Expectations:   upr.Expectations,
		ExpectedReturn: upr.ExpectedReturn,
	}
}
