package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID       primitive.ObjectID        `json:"_id,omitempty" bson:"_id,omitempty"`             
	Name     string                    `json:"name,omitempty"`
	Email    string                    `json:"email,omitempty"`
	Password string                    `json:"password,omitempty"`
	Token    string                     `json:"token,omitempty"`
}

type Claims struct {
	UserEmail string       `json:"user-email"`
	jwt.StandardClaims
}

type Token struct{
	UserEmail string     `json:"userEmail"` 
	Token string         `json:"token"`
}

