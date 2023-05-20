package model

type Contact struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	EMail string `json:"email,omitempty" bson:"email,omitempty"`
	Phone string `json:"phone,omitempty" bson:"phone,omitempty"`
	Fax   string `json:"fax,omitempty" bson:"fax,omitempty"`
}
