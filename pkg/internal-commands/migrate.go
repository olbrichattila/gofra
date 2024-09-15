package internalcommand

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/olbrichattila/gofra/pkg/app/args"
	"github.com/olbrichattila/gofra/pkg/app/db"

	migrator "github.com/olbrichattila/godbmigrator"
)

const defaultMigrationFilePath = "./migrations"

func Migrate(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, MigrationProvider, migrationFilePath, step, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	err = migrator.Migrate(dbConn, MigrationProvider, migrationFilePath, step)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Rollback(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, MigrationProvider, migrationFilePath, step, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	err = migrator.Rollback(dbConn, MigrationProvider, migrationFilePath, step)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Refresh(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, MigrationProvider, migrationFilePath, _, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	err = migrator.Refresh(dbConn, MigrationProvider, migrationFilePath)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Report(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, MigrationProvider, migrationFilePath, _, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	report, err := migrator.Report(dbConn, MigrationProvider, migrationFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(report)
}

func Add(a args.CommandArger, dbConfig db.DBFactoryer) {
	migrationFilePath := getMigrationFilePath()
	customPrefix := ""
	params := a.GetAll()
	if len(params) > 0 {
		customPrefix = params[0]
	}

	err := migrator.AddNewMigrationFiles(migrationFilePath, customPrefix)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func constructMigratorOptions(a args.CommandArger, dbConfig db.DBFactoryer) (*sql.DB, migrator.MigrationProvider, string, int, error) {
	migrationFilePath := getMigrationFilePath()
	step := getStep(a)

	dbConf, err := dbConfig.GetConnectionConfig()
	if err != nil {
		return nil, nil, "", 0, err
	}

	dbConn, err := sql.Open(dbConf.GetConnectionName(), dbConf.GetConnectionString())
	if err != nil {
		return nil, nil, "", 0, err
	}

	MigrationProvider, err := migrator.NewMigrationProvider("db", "", dbConn)
	if err != nil {
		return nil, nil, "", 0, err
	}

	return dbConn, MigrationProvider, migrationFilePath, step, nil
}

func getStep(a args.CommandArger) int {
	stepStr, _ := a.GetFlagByName("step", "0")

	if nr, err := strconv.Atoi(stepStr); err == nil {
		return nr
	}

	return 0
}

func getMigrationFilePath() string {
	mPath := os.Getenv("MIGRATOR_MIGRATION_PATH")
	if mPath != "" {
		return mPath
	}

	return defaultMigrationFilePath
}
