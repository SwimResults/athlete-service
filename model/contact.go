package model

type Contact struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	EMail string `json:"e_mail,omitempty" bson:"e_mail,omitempty"`
	Phone string `json:"phone,omitempty" bson:"phone,omitempty"`
	Fax   string `json:"fax,omitempty" bson:"fax,omitempty"`
}
