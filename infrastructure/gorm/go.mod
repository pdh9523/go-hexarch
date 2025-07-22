module github.com/pdh9523/go-hexarch/infrastructure/gorm

go 1.24.5

require (
	github.com/google/uuid v1.6.0
	github.com/pdh9523/go-hexarch/domains/user v0.0.0-00010101000000-000000000000
	github.com/pdh9523/go-hexarch/shared/common v0.0.0-00010101000000-000000000000
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)

replace github.com/pdh9523/go-hexarch/shared/common => ../../shared/common

replace github.com/pdh9523/go-hexarch/domains/user => ../../domains/user
