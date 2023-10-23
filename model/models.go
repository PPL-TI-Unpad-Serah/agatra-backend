package model

import "time"

type Machine struct {
	ID 			int
	Version_ID	int
	Location_ID	int
	Count		int
	Price		int
}

type Location struct {
	ID 			int
	Centre_ID	int
	City_ID 	int
	Name 		string
}

type City struct {
	ID 			int
	Name		string
}

type Centre struct {
	ID 			int
	Name		string
	Price_Info	string
}

type Version struct {
	ID 			int
	Title_ID	int
	Name		string
	Price_Info	string
}

type Title struct {
	ID 			int
	Name		string
}

type User struct {
	ID 			int
	Name		string
	Email		string
	Password	string
	Role		string
}

type Session struct {
	ID     int       `gorm:"primaryKey" json:"id"`
	Token  string    `json:"token"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
}