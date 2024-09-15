appwizard:
	cd cmd/appwizard && go build -o appwizard
	./cmd/appwizard/appwizard
run:
	go run ./cmd/
build:
	go build -o gofra ./cmd/
build-and-run:
	go build -o gofra ./cmd/
	./gofra
db-sqlite-recreate:
	@if [ -f ./database/database.sqlite ]; then rm ./database/database.sqlite; fi
	csvimporter data.csv vehicles ";"
	migrator migrate
db-recreate:
	csvimporter data.csv vehicles ";"
	migrator refresh
db-recreate-firebird:
	csvimporter car_basemodel.csv car_basemodel
	csvimporter car_fuel_type.csv car_fuel_type
	csvimporter car_make.csv car_make
	csvimporter car_model.csv car_model
	csvimporter car_year.csv car_year
	csvimporter data.csv vehicles ";"
	migrator refresh
switch-sqlite:
	cp .env.sqlite.example .env.migrator
	cp .env.sqlite.example .env.csvimporter
	cp .env.sqlite.example .env.csvexporter
switch-mysql:
	cp .env.mysql.example .env.migrator
	cp .env.mysql.example .env.csvimporter
	cp .env.mysql.example .env.csvexporter
switch-pgsql:
	cp .env.pgsql.example .env.migrator
	cp .env.pgsql.example .env.csvimporter
	cp .env.pgsql.example .env.csvexporter
switch-firebird:
	cp .env.firebird.example .env.migrator
	cp .env.firebird.example .env.csvimporter
	cp .env.firebird.example .env.csvexporter
lint:
	gocritic check ./...
	revive ./...
	golint ./...
	go vet ./...
	staticcheck ./...
	golangci-lint run
	goconst ./...
test:
	go test ./...
