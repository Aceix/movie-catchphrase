package models

import "go.mongodb.org/mongo-driver/bson/primitive"

/**
Catchphrase model represents a catch phrase
*/
type Catchphrase struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MovieName string `json:"movieName,omitempty" bson:"movieName,omitempty"`
	CatchPhrase string `json:"catchPhrase,omitempty" bson:"catchPhrase,omitempty"`
	MovieContext string `json:"movieContext,omitempty" bson:"movieContext,omitempty"`
}
