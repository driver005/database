module github.com/driver005/database

go 1.17

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jackc/pgx/v4 v4.16.1
	github.com/jinzhu/inflection v1.0.0
	github.com/jinzhu/now v1.1.5
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/text v0.3.7 // indirect
)

retract (
	v1.0.2
	v1.0.1
	v1.0.0
)
