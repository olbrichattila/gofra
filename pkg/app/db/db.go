package db

// Todo add sync, to avoid database locks as multiple go routines can use it the same time
import (
	"database/sql"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	// This needs to be blank imported as not directly referenced, but required
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

type DBer interface {
	// Construct(DBFactoryer, logger.Logger)
	Construct(DBFactoryer)
	LowerCaseResult()
	OriginalCaseResult()
	Open()
	Close()
	QueryAll(string, ...any) <-chan map[string]interface{}
	QueryOne(string, ...any) (map[string]interface{}, error)
	Execute(string, ...any) (int64, error)
	GetLastError() error
}

type DB struct {
	// l         logger.Logger
	lowerCaseResult bool
	db              *sql.DB
	dbConfig        DBConfiger
	lastError       error
}

func New() DBer {
	db := &DB{
		lowerCaseResult: true,
	}
	runtime.SetFinalizer(db, func(db *DB) {
		db.Cleanup()
	})
	return db
}

// func (d *DB) Construct(dbConfig DBFactoryer, l logger.Logger) {
// 	d.l = l
// 	var err error
// 	d.dbConfig, err = dbConfig.GetConnectionConfig()
// 	if err != nil {
// 		l.Error(fmt.Sprintf("Cannot get database config: %s", err.Error()))
// 		return
// 	}

// 	d.Open()
// }

func (d *DB) Construct(dbConfig DBFactoryer) {

	var err error
	d.dbConfig, err = dbConfig.GetConnectionConfig()
	if err != nil {
		return
	}

	d.Open()
}

func (d *DB) LowerCaseResult() {
	d.lowerCaseResult = true
}

func (d *DB) OriginalCaseResult() {
	d.lowerCaseResult = false
}

func (d *DB) Cleanup() {
	d.Close()
}

func (d *DB) Open() {
	var err error
	d.db, err = sql.Open(d.dbConfig.GetConnectionName(), d.dbConfig.GetConnectionString())
	if err != nil {
		d.logError(err.Error())
	}
}

func (d *DB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d *DB) QueryAll(sql string, pars ...any) <-chan map[string]interface{} {
	ch := make(chan map[string]interface{}, 1)
	d.lastError = nil
	if d.db == nil {
		d.lastError = fmt.Errorf("db not open")
		close(ch)
		return ch
	}

	stmt, err := d.db.Prepare(sql)
	if err != nil {
		d.lastError = err
		close(ch)
		return ch
	}

	rows, err := stmt.Query(pars...)
	if err != nil {
		d.lastError = err
		stmt.Close()
		close(ch)
		return ch
	}

	cols, err := rows.Columns()
	if err != nil {
		d.lastError = err
		rows.Close()
		stmt.Close()
		close(ch)
		return ch
	}

	colCount := len(cols)

	row := make([]interface{}, colCount)
	for i := range row {
		row[i] = new(interface{})
	}

	go func() {
		for rows.Next() {
			err := rows.Scan(row...)
			if err != nil {
				d.lastError = err
				break
			}
			result := make(map[string]interface{}, colCount)
			for i, colName := range cols {
				if d.lowerCaseResult {
					colName = strings.ToLower(colName)
				}
				value := *(row[i].(*interface{}))

				switch v := value.(type) {
				case string:
					result[colName] = v
				case []byte:
					result[colName] = string(v)
				default:
					result[colName] = v
				}
			}
			ch <- result
			result = nil
		}
		rows.Close()
		stmt.Close()
		close(ch)
	}()

	return ch
}

func (d *DB) QueryOne(sql string, pars ...any) (map[string]interface{}, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db not open")
	}

	stmt, err := d.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(pars...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colCount := len(cols)
	row := make([]interface{}, colCount)
	for i := range row {
		row[i] = new(interface{})
	}

	if !rows.Next() {
		return nil, fmt.Errorf("row cannot be found")
	}

	err = rows.Scan(row...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, colCount)
	for i, colName := range cols {
		if d.lowerCaseResult {
			colName = strings.ToLower(colName)
		}
		value := *(row[i].(*interface{}))
		switch v := value.(type) {
		case string:
			result[colName] = v
		case []byte:
			result[colName] = string(v)
		default:
			result[colName] = v
		}
	}

	return result, nil
}

func (d *DB) Execute(sql string, pars ...any) (int64, error) {
	if d.db == nil {
		return 0, fmt.Errorf("db not open")
	}

	stmt, err := d.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(pars...)
	if err != nil {
		return 0, err
	}

	// PG SQL and FirebirdSql does not support last insert ID, try to get it
	// It is really hacky, and could just use RETURNS id in the SQL, it is here to be compatible with MySql, SqLite
	// but way less performant
	switch d.dbConfig.GetConnectionName() {
	case DriverNamePostgres:
		return d.getPgSQLLastInsertId(sql)
	case DriverNameFirebird:
		return d.getFirebirdLastInsertId(sql)
	}

	return res.LastInsertId()
}

func (d *DB) getPgSQLLastInsertId(sql string) (int64, error) {
	tableName := d.extractTableName(sql)
	if tableName == "" {
		return 0, nil
	}

	primaryKey := d.getPrimaryKey(tableName)
	if primaryKey == "" {
		return 0, nil
	}

	r, err := d.QueryOne("SELECT currval(pg_get_serial_sequence('" + tableName + "', '" + primaryKey + "')) as last_id")
	if err != nil {
		return 0, err
	}
	if val, ok := r["last_id"]; ok {
		return val.(int64), nil
	}

	return 0, nil
}

func (d *DB) getFirebirdLastInsertId(sql string) (int64, error) {
	tableName := d.extractTableName(sql)
	if tableName == "" {
		return 0, nil
	}

	generatorName := d.getFirebirdGeneratorName(tableName)
	if generatorName == "" {
		return 0, nil
	}

	r, err := d.QueryOne("SELECT GEN_ID(" + generatorName + ", 0) AS \"last_id\" FROM RDB$DATABASE;")
	if err != nil {
		return 0, err
	}
	if val, ok := r["last_id"]; ok {
		return val.(int64), nil
	}

	return 0, nil
}

func (d *DB) getFirebirdGeneratorName(tableName string) string {
	sql := `SELECT FIRST 1
RDB$GENERATOR_NAME AS "generator_name"
FROM RDB$RELATION_FIELDS WHERE RDB$RELATION_NAME = ?
AND RDB$GENERATOR_NAME IS NOT NULL 
AND RDB$IDENTITY_TYPE = 1
ORDER BY RDB$FIELD_ID`

	res, err := d.QueryOne(sql, tableName)
	if err != nil {
		return ""
	}

	if id, ok := res["generator_name"]; ok {
		return id.(string)
	}

	return ""
}

func (d *DB) getPrimaryKey(tableName string) string {
	sql := `SELECT a.attname
        FROM   pg_index i
        JOIN   pg_attribute a ON a.attrelid = i.indrelid
                             AND a.attnum = ANY(i.indkey)
        WHERE  i.indrelid = $1::regclass
        AND    i.indisprimary`

	res, err := d.QueryOne(sql, tableName)
	if err != nil {
		return ""
	}

	if id, ok := res["attname"]; ok {
		return id.(string)
	}

	return ""
}

func (d *DB) extractTableName(sql string) string {
	query := strings.ToLower(sql)
	re := regexp.MustCompile(`insert\sinto\s+\"(\w+)\"`)
	match := re.FindStringSubmatch(query)
	if len(match) > 1 {
		return strings.ReplaceAll(match[1], `"`, "")
	}
	return ""
}

func (d *DB) GetLastError() error {
	return d.lastError
}

func (d *DB) logError(message string) {
	// if d.l != nil {
	// 	d.l.Error(message)
	// }
}
