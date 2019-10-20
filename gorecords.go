package gorecords

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gensword/collections"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Export interface {
	ToMaps() []map[string]interface{}
	ToSlices() [][]interface{}
	ToJson() (string, error)
	ToCsv(withHeader bool) [][]string
}

type SingleExport interface {
	ToMap() map[string]interface{}
	ToSlice() []interface{}
	ToJson() (string, error)
	ToCsv(withHeader bool) [][]string
}

type Record struct {
	collections.OrderedMap
}

type Records []Record

type DB interface {
	Open() *sql.DB
	Close()
	GetTableNames() ([]string, error)
	GetDBName() string
	GetConnection() *sql.DB
	Query(sql string) (Records, error)
	Exec(sql string) (sql.Result, error)
}

type MysqlDataBase struct {
	Url        string
	dbName     string
	connection *sql.DB
}

func (mysqlDb *MysqlDataBase) Open() *sql.DB {
	parts := strings.Split(mysqlDb.Url, "/")
	dbNameAndParamParts := parts[len(parts)-1]
	dbNameEndPos := strings.Index(dbNameAndParamParts, "?")
	var dbName string
	if dbNameEndPos == -1 {
		dbName = dbNameAndParamParts[:]
	} else {
		dbName = dbNameAndParamParts[:dbNameEndPos]
	}
	mysqlDb.dbName = dbName
	db, err := sql.Open("mysql", mysqlDb.Url)
	if err != nil {
		panic(fmt.Sprintf("can not connect to DB %s", mysqlDb.dbName))
	}
	mysqlDb.connection = db
	return db
}

func (mysqlDb *MysqlDataBase) Close() {
	mysqlDb.connection.Close()
	mysqlDb.connection = nil
}

func (mysqlDb *MysqlDataBase) GetDBName() string {
	return mysqlDb.dbName
}

func (mysqlDb *MysqlDataBase) GetTableNames() ([]string, error) {
	rows, err := mysqlDb.connection.Query("show tables")
	tableNames := make([]string, 0)
	if err != nil {
		return tableNames, err
	}
	for rows.Next() {
		var tableName []byte
		rows.Scan(&tableName)
		tableNames = append(tableNames, string(tableName))
	}
	return tableNames, nil
}

func (mysqlDb *MysqlDataBase) GetConnection() *sql.DB {
	return mysqlDb.connection
}

func (mysqlDb *MysqlDataBase) Query(sql string) (Records, error) {
	return getRecords(mysqlDb.connection, sql)
}

func (mysqlDb *MysqlDataBase) Exec(sql string) (sql.Result, error) {
	return mysqlDb.connection.Exec(sql)
}

func (record *Record) ToMap() map[string]interface{} {
	colNamesValues := make(map[string]interface{})
	for item := range record.Iter() {
		for colName, colValue := range item {
			colNamesValues[colName.(string)] = colValue
		}
	}
	return colNamesValues
}

func (record *Record) ToSlice() []interface{} {
	colValues := make([]interface{}, 0)
	for item := range record.Iter() {
		for _, colValue := range item {
			colValues = append(colValues, colValue)
		}
	}
	return colValues
}

func (record *Record) ToJson() (string, error) {
	data, err := json.Marshal(record.ToMap())
	return string(data), err
}

func (record *Record) ToCsv(withHeader bool) [][]string {
	data := make([][]string, 0)
	recordSlice := record.ToSlice()
	recordStringSlice := make([]string, len(recordSlice))
	for i, v := range recordSlice {
		recordStringSlice[i] = fmt.Sprint(v)
	}
	if withHeader {
		data = append(data, record.getHeader(), recordStringSlice)
	} else {
		data = append(data, recordStringSlice)
	}
	return data
}

func (record *Record) getHeader() []string {
	colNames := make([]string, 0)
	for item := range record.Iter() {
		for col, _ := range item {
			colNames = append(colNames, col.(string))
		}
	}
	return colNames
}

func (records *Records) ToMaps() []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	for _, record := range *records {
		data = append(data, record.ToMap())
	}
	return data
}

func (records *Records) ToSlices() [][]interface{} {
	data := make([][]interface{}, 0)
	for _, record := range *records {
		data = append(data, record.ToSlice())
	}
	return data
}

func (records *Records) ToJson() (string, error) {
	data, err := json.Marshal(records.ToMaps())
	return string(data), err
}

func (records *Records) ToCsv(withHeader bool) [][]string {
	data := make([][]string, 0)
	recordSlices := records.ToSlices()
	if withHeader {
		data = append(data, (*records)[0].getHeader())
	}
	for _, recordSlice := range recordSlices {
		recordStringSlice := make([]string, len(recordSlice))
		for i, v := range recordSlice {
			recordStringSlice[i] = fmt.Sprint(v)
		}
		data = append(data, recordStringSlice)
	}

	return data
}

func getRecords(db *sql.DB, sqlString string) (Records, error) {
	var records Records
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return records, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return records, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return records, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return records, err
		}

		entry := collections.NewOederedMap()
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if (ok) {
				entry.Set(col, string(b), true)
			} else {
				entry.Set(col, v, true)
			}
		}

		records = append(records, Record{*entry})
	}
	return records, nil
}

func NewMysqlDataBase(url string) *MysqlDataBase {
	return &MysqlDataBase{Url:url}
}
