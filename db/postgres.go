package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "agatra/model"
)

type Postgres struct{}

type Config struct{
	Host		string
	Port		string
	Password	string
	User		string
	DBName		string
	SSLMode		string
}

func (p *Postgres) Connect(creds *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", creds.Host, creds.User, creds.Password, creds.DBName, creds.Port)

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func NewDB() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Reset(db *gorm.DB, table string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE " + table).Error; err != nil {
			return err
		}

		if err := tx.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1").Error; err != nil {
			return err
		}

		return nil
	})
}