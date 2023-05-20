package model

type ColorSet struct {
	Primary   string `json:"primary,omitempty" bson:"primary,omitempty"`
	Secondary string `json:"secondary,omitempty" bson:"secondary,omitempty"`
	Contrast  string `json:"contrast,omitempty" bson:"contrast,omitempty"`
}
