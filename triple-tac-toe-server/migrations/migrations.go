package migrations

import (
	"errors"

	"gregvader/triple-tac-toe/database"
	"gregvader/triple-tac-toe/domain"

	"log"
)

func Migrate(dbConn database.DBConn) {
	if err := dropAllTables(dbConn); err != nil {
		logger := log.Logger{}
		logger.Println(err)
	}
	dbConn.DB.AutoMigrate(&domain.User{})

}

func dropAllTables(dbConn database.DBConn) error {
	//Method used to remove all tables (INCLUDING MANY TO MANY (main reason why this method is used))
	err := dbConn.DB.Exec("DO $$ DECLARE\n    r RECORD;\nBEGIN\n    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP\n        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';\n    END LOOP;\nEND $$;")
	if err.Error != nil {
		return errors.New("migration error: tables could not be dropped")
	}
	return nil
}
