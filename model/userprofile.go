package model

type Expectation struct {
	Stock         string  `bson:"stock" json:"stock"`
	ExpectedRaise float64 `bson:"expectedRaise" json:"expectedRaise"`
}

type Userprofile struct {
	UserID             string        `bson:"_id" json:"userId"`
	Email              string        `bson:"email" json:"email"`
	Expectations       []Expectation `bson:"expectations" json:"expectations"`
	ExpectedReturn     float64       `bson:"expectedReturn" json:"expectedReturn"`
	DefaultExpectation float64       `bson:"defaultExpectations" json:"defaultExpectation"`
}

//GetExpectation return the expectation for the given stock if exists, returns
// the default otherwise
func (u *Userprofile) GetExpectation(stock string) float64 {
	expectations := u.Expectations

	for _, exp := range expectations {
		if exp.Stock == stock {
			return exp.ExpectedRaise
		}
	}

	return u.DefaultExpectation
}
