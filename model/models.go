package model

import (
	"time"
)

type City struct {
	ID 			int				`gorm:"primaryKey" json:"city_id"`
	Name		string			`gorm:"notNull" json:"name"`
}

type Center struct {
	ID 			int				`gorm:"primaryKey" json:"center_id"`
	Name		string			`gorm:"notNull" json:"name"`
	Info		string			`json:"info"`
}

type Machine struct {
	ID 			int				`gorm:"primaryKey" json:"machine_id"`
	VersionID	int				`gorm:"notNull" json:"version_id"`
	Version		Version 		`gorm:"notNull;foreignKey:VersionID" json:"version"`
	Count		int				`gorm:"notNull" json:"machine_count"`
	Price		int				`json:"price"`
	LocationID	int				`gorm:"notNull" json:"location_id"`
	Location	Location		`gorm:"notNull;foreignKey:LocationID" json:"center"`
}

type Location struct {	
	ID 			int				`gorm:"primaryKey" json:"location_id"`
	Name 		string			`gorm:"notNull" json:"name"`
	Description	string			`json:"description"`
	Lat 		float32			`json:"lat"`
	Long		float32			`json:"long"`
	CenterID	int				`gorm:"notNull" json:"center_id"`
	Center		Center			`gorm:"notNull;foreignKey:CenterID" json:"center"`
	Machine		[]Machine		`gorm:"notNull;foreignKey:LocationID" json:"machine"`
	CityID		int				`gorm:"notNull" json:"city_id"`
	City 		City			`gorm:"notNull;foreignKey:CityID" json:"city"`
}

type Version struct {
	ID 			int				`gorm:"primaryKey" json:"version_id"`
	TitleID		int				`gorm:"notNull" json:"title_id"`
	Title		Title			`gorm:"notNull;foreignKey:TitleID" json:"title"`
	Name		string			`gorm:"notNull" json:"name"`
	Info		string			`gorm:"notNull" json:"info"`
}

type Title struct {
	ID 			int					`gorm:"primaryKey" json:"title_id"`
	Name		string				`gorm:"notNull" json:"name"`
	Version		[]Version			`gorm:"notNull; foreignKey:TitleID" json:"version"`
}

type User struct {
	ID 			int				`gorm:"primaryKey" json:"user_id"`
	Name		string			`gorm:"notNull" json:"name"`	
	Email		string			`gorm:"notNull" json:"email"`
	Password	string			`gorm:"notNull" json:"password"`
	Role		string			`gorm:"notNull" json:"role"`
}

type Session struct {
	ID     int       			`gorm:"primaryKey" json:"session_id"`
	Token  string    			`json:"token"`
	Email  string    			`json:"email"`
	Expiry time.Time 			`json:"expiry"`
}
type Versions struct {
	ID 			int				`gorm:"primaryKey" json:"version_id"`
	Title		Title_compact	`gorm:"notNull;foreignKey:id" json:"title"`
	Name		string			`gorm:"notNull" json:"name"`
	Info		string			`gorm:"notNull" json:"info"`
}

type Titles struct {
	ID 			int					`gorm:"primaryKey" json:"title_id"`
	Name		string				`gorm:"notNull" json:"name"`
	Version		[]Version_compact	`gorm:"notNull;foreignKey:id" json:"version"`
}