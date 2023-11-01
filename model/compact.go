package model

type Location_compact struct{
	ID		int				`gorm:"primaryKey" json:"location_id"`
	Name	string			`gorm:"notNull" json:"name"`
	Version	[]Version		`gorm:"notNull;foreignKey:version_id" json:"version"`
}

type Version_compact struct{
	ID		int				`gorm:"primaryKey" json:"version_id"`
	Name	string			`gorm:"notNull" json:"name"`
	Info	string			`gorm:"notNull" json:"info"`
}

type Title_compact struct{
	ID		int				`gorm:"primaryKey" json:"title_id"`
	Name	string			`gorm:"notNull" json:"name"`
}