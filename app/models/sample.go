package models

type Address struct {
	City    string `json:"city" bson:"city"`
	Pincode string `json:"pincode" bson:"pincode"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
}

type Sample struct {
	Id        int     `json:"id" bson:"id"`
	FirstName string  `json:"firstName" bson:"firstName"`
	LastName  string  `json:"lastName" bson:"lastName"`
	Age       uint16  `json:"age" bson:"age"`
	Address   Address `json:"address" bson:"address"`
}
