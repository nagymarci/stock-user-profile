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
