package model

import (
	"math"
	"time"
)

type City struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	Name		string			`gorm:"notNull" json:"name"`
}

type Center struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	Name		string			`gorm:"notNull" json:"name"`
	Info		string			`json:"info"`
}

type Machine struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	VersionID	int				`gorm:"notNull" json:"version_id"`
	Version		Version 		`gorm:"notNull;foreignKey:VersionID" json:"version"`
	Count		int				`gorm:"notNull" json:"machine_count"`
	Price		int				`json:"price"`
	LocationID	int				`gorm:"notNull" json:"location_id"`
	Location	Location		`gorm:"notNull;foreignKey:LocationID" json:"location"`
	Notes		string			`json:"notes"`
}

type Location struct {	
	ID 			int				`gorm:"primaryKey" json:"id"`
	Name 		string			`gorm:"notNull" json:"name"`
	Description	string			`json:"description"`
	Lat 		float64			`json:"lat"`
	Long		float64			`json:"long"`
	CenterID	int				`gorm:"notNull" json:"center_id"`
	Center		Center			`gorm:"notNull;foreignKey:CenterID" json:"center"`
	Machine		[]Machine		`gorm:"notNull;foreignKey:LocationID" json:"machine"`
	CityID		int				`gorm:"notNull" json:"city_id"`
	City 		City			`gorm:"notNull;foreignKey:CityID" json:"city"`
}

type Version struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	TitleID		int				`gorm:"notNull" json:"title_id"`
	Title		Title			`gorm:"notNull;foreignKey:TitleID" json:"title"`
	Name		string			`gorm:"notNull" json:"name"`
	Info		string			`gorm:"notNull" json:"info"`
}

type Title struct {
	ID 			int					`gorm:"primaryKey" json:"id"`
	Name		string				`gorm:"notNull" json:"name"`
	Version		[]Version			`gorm:"notNull; foreignKey:TitleID" json:"versions"`
}

type User struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	Name		string			`gorm:"notNull" json:"name"`	
	Email		string			`gorm:"notNull" json:"email"`
	Password	string			`gorm:"notNull" json:"password"`
	Role		string			`gorm:"notNull" json:"role"`
}

type Session struct {
	ID     int       			`gorm:"primaryKey" json:"id"`
	Token  string    			`json:"token"`
	Email  string    			`json:"email"`
	Expiry time.Time 			`json:"expiry"`
}

type Versions struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	Title		Title_compact	`gorm:"notNull;foreignKey:id" json:"title"`
	Name		string			`gorm:"notNull" json:"name"`
	Info		string			`gorm:"notNull" json:"info"`
}

type Titles struct {
	ID 			int					`gorm:"primaryKey" json:"id"`
	Name		string				`gorm:"notNull" json:"name"`
	Version		[]Version_compact	`gorm:"notNull;foreignKey:id" json:"version"`
}

type Location_range struct{
	ID 			int				
	Name 		string			
	Description	string			
	Lat 		float64			
	Long		float64			
	Distance	float64			
}

func (l *Location)CheckRange(lat float64, long float64) Location_range{
	return Location_range{
		Lat: l.Lat,
		Long: l.Long,
		Distance: math.Sqrt((l.Lat - lat) * (l.Lat - lat)) + ((l.Long - long) * (l.Long - long)),
	}
}