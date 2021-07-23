package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gregvader/triple-tac-toe/globals"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConn struct {
	DB *gorm.DB
}

func (db *DBConn) ConnectToDatabase() error {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Belgrade",
			globals.DbHost, globals.DbUser, globals.DbPassword, globals.DbName, globals.DbPort),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: db.getLogger(),
	})
	db.DB = conn

	return err
}

func (db *DBConn) getLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)
	return newLogger
}
