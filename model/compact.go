package model

type Tabler interface {
	TableName() string
}

type Location_compact struct{
	ID		int				`gorm:"primaryKey" json:"id"`
	Name	string			`gorm:"notNull" json:"name"`
	Version	[]Version		`gorm:"notNull;foreignKey:version_id" json:"version"`
}

type Version_compact struct{
	ID		int				`gorm:"primaryKey" json:"id"`
	Name	string			`gorm:"notNull" json:"name"`
	TitleID	int				`gorm:"notNull" json:"title_id"`
	Info	string			`gorm:"notNull" json:"info"`
}

func (Version_compact) TableName() string {
	return "versions"
}

type Title_compact struct{
	ID		int				`gorm:"primaryKey" json:"id"`
	Name	string			`gorm:"notNull" json:"name"`
}

func (Title_compact) TableName() string {
	return "titles"
}

type User_compact struct {
	ID 			int				`gorm:"primaryKey" json:"id"`
	Username	string			`gorm:"notNull" json:"username"`	
	Email		string			`gorm:"notNull" json:"email"`
	Role		string			`gorm:"notNull" json:"role"`
}

type User_login struct{
	Username 	string			`gorm:"notNull" json:"username"`
	Password	string			`gorm:"notNull" json:"password"`
}



func LocationToCompact(lf Location) Location_compact{
	compact := Location_compact{
		ID: lf.ID,
		Name: lf.Name,
	}

	for _, machine := range lf.Machine {
		compact.Version = append(compact.Version, machine.Version)
	}

	return compact;
}

func VersionToCompact(vf Version) Version_compact{
	return Version_compact{
		ID: 	vf.ID,
		Name: 	vf.Name,
		Info: 	vf.Info,
	}
}

func TitleToCompact(tf Title) Title_compact{
	return Title_compact{
		ID: 	tf.ID,
		Name: 	tf.Name,
	}
}

func UserToCompact(uf User) User_compact{
	return User_compact{
		ID: 		uf.ID,
		Username: 	uf.Username,
		Email:		uf.Email,
		Role:		uf.Role,
	}
}



// func MassCompactUser(uf []User) []User_compact{
// 	var result []model.User
// 	for _, data := range uf {
// 		UserToCompact(data)
// 	}
// }
