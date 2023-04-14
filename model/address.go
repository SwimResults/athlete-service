package model

type Address struct {
	Street     string `json:"street,omitempty" bson:"street,omitempty"`
	Number     string `json:"number,omitempty" bson:"number,omitempty"`
	City       string `json:"city,omitempty" bson:"city,omitempty"`
	PostalCode string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
}
