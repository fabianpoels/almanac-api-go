package models

type TranslatedField struct {
	En string `json:"en,omitempty" bson:"en,omitempty"`
	Ar string `json:"ar,omitempty" bson:"ar,omitempty"`
	Fr string `json:"fr,omitempty" bson:"fr,omitempty"`
}
