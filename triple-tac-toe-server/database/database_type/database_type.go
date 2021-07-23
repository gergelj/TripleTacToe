package database_type

type DatabaseType string

const (
	Transactional DatabaseType = "db_trx"
)

func (databaseType DatabaseType) String() string {
	return string(databaseType)
}
