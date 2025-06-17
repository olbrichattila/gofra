package internalcommand

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/olbrichattila/godbmigrator/config"
	"github.com/olbrichattila/gofra/pkg/app/args"
	"github.com/olbrichattila/gofra/pkg/app/db"

	migrator "github.com/olbrichattila/godbmigrator"
)

const defaultMigrationFilePath = "./migrations"

func Migrate(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, m, step, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	err = m.Migrate(step)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Rollback(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, m, step, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	err = m.Rollback(step)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Refresh(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, m, _, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	err = m.Refresh()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Report(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, m, _, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	report, err := m.Report()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(report)
}

func Add(a args.CommandArger, dbConfig db.DBFactoryer) {
	dbConn, m, _, err := constructMigratorOptions(a, dbConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dbConn.Close()

	customPrefix := ""
	params := a.GetAll()
	if len(params) > 0 {
		customPrefix = params[0]
	}

	err = m.AddNewMigrationFiles(customPrefix)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func constructMigratorOptions(a args.CommandArger, dbConfig db.DBFactoryer) (*sql.DB, migrator.DBMigrator, int, error) {
	migrationFilePath := getMigrationFilePath()
	step := getStep(a)

	dbConf, err := dbConfig.GetConnectionConfig()
	if err != nil {
		return nil, nil, 0, err
	}

	dbConn, err := sql.Open(dbConf.GetConnectionName(), dbConf.GetConnectionString())
	if err != nil {
		return nil, nil, 0, err
	}

	newMigrator := migrator.New(dbConn, migrationFilePath, "gofra")
	newMigrator.SubscribeToMessages(func(et int, msg string) {
		fmt.Println(decorateMessage(et, msg))
	})

	return dbConn, newMigrator, step, nil
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

// DecorateMessage will return a message with the correct context from the event type and message
func decorateMessage(eventType int, message string) string {
	if formattedMsg, exists := getMessageFormat(eventType); exists {
		return fmt.Sprintf(formattedMsg, message)
	}

	return message
}

// getMessageFormat returns the message format string based on the event type.
func getMessageFormat(eventType int) (string, bool) {
	messages := map[int]string{
		config.MigratedItems:        "Migrated %s items",
		config.NothingToRollback:    "Nothing to roll back%s",
		config.RolledBack:           "Rolled back %s items",
		config.RunningMigrations:    "Running migration: %s",
		config.SkipRollback:         "Skip rollback as file '%s' not exists",
		config.RunningRollback:      "Running rollback %s",
		config.MigrationFileCreated: "Migration file created: %s",
	}

	format, exists := messages[eventType]
	return format, exists
}
