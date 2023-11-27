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

type User_compact struct {
	ID 			int				`gorm:"primaryKey" json:"user_id"`
	Name		string			`gorm:"notNull" json:"name"`	
	Email		string			`gorm:"notNull" json:"email"`
	Password	string			`gorm:"notNull" json:"password"`
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
		ID: 	uf.ID,
		Name: 	uf.Name,
		Email:	uf.Email,
		Role:	uf.Role,
	}
}



// func MassCompactUser(uf []User) []User_compact{
// 	var result []model.User
// 	for _, data := range uf {
// 		UserToCompact(data)
// 	}
// }
