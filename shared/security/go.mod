module github.com/pdh9523/go-hexarch/shared/security

go 1.24.5

require (
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/pdh9523/go-hexarch/shared/common v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.40.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/pdh9523/go-hexarch/shared/common => ../common
