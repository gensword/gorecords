package gorecords

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var dbUrl string = "user:password@tcp(localhost:3306)/dbname?tls=skip-verify&autocommit=true"

func TestMysqlDataBase_GetDBName(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	assert.Equal(t, "cron", mysqlDb.GetDBName())
	mysqlDb.Close()
}

func TestMysqlDataBase_GetConnection(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	sqlDb := mysqlDb.Open()
	assert.Equal(t, sqlDb, mysqlDb.GetConnection())
	mysqlDb.Close()
	assert.Empty(t, mysqlDb.GetConnection())
}

func TestMysqlDataBase_GetTableNames(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	table_names, err := mysqlDb.GetTableNames()
	assert.Nil(t, err)
	assert.Equal(t, []string{"logs", "users"}, table_names)
}

func TestMysqlDataBase_Exec(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	res, err := mysqlDb.Exec("insert into users (`user_name`, `password`) values ('gensword', 'gensword')")
	assert.Nil(t, err)
	rows, _ := res.RowsAffected()
	assert.EqualValues(t, 1, interface{}(rows))
	mysqlDb.Close()
}

func TestMysqlDataBase_Query(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	rows, err := mysqlDb.Query("select * from users limit 3")
	assert.Nil(t, nil, err)
	assert.Len(t, rows, 3)
	mysqlDb.Close()
}

func TestRecords_ToJson(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	rows, err := mysqlDb.Query("select created_at, deleted_at, id, updated_at, user_name from users limit 3")
	assert.Nil(t, nil, err)
	jsonStr, err := rows.ToJson()
	assert.Nil(t, nil, err)
	assert.JSONEq(t, `[{"created_at":"2019-08-16 18:07:10","deleted_at":null,"id":1,"updated_at":"2019-08-16 18:07:10","user_name":"xzy"},{"created_at":null,"deleted_at":null,"id":2,"updated_at":null,"user_name":"gensword"},{"created_at":null,"deleted_at":null,"id":3,"updated_at":null,"user_name":"gensword"}]`, jsonStr)
}

func TestRecords_ToMaps(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	rows, err := mysqlDb.Query("select user_name, deleted_at, created_at, updated_at, id from users limit 2")
	assert.Nil(t, nil, err)
	maps := rows.ToMaps()
	shouldEqual := []map[string]interface{}{map[string]interface{}{"deleted_at": nil, "user_name": "xzy", "created_at": "2019-08-16 18:07:10", "updated_at": "2019-08-16 18:07:10", "id": int64(1)}, map[string]interface{}{"deleted_at": nil, "user_name": "gensword", "created_at": nil, "updated_at": nil, "id": int64(2)}}
	assert.Equal(t, 2, len(rows))
	//fmt.Println(assert.ObjectsAreEqual(shouldEqual[0], maps[0]))
	assert.Equal(t, shouldEqual[0], maps[0])
	assert.Equal(t, shouldEqual[1], maps[1])
}

func TestRecords_ToSlices(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	rows, err := mysqlDb.Query("select user_name, deleted_at, id from users limit 2")
	assert.Nil(t, nil, err)
	slices := rows.ToSlices()
	assert.Equal(t, 2, len(slices))
	shouldEqual := make([][]interface{}, 0)
	shouldEqual = append(shouldEqual, []interface{}{"xzy", nil, int64(1)}, []interface{}{"gensword", nil, int64(2)})
	assert.Equal(t, shouldEqual, slices)
}

func TestRecords_ToCsv(t *testing.T) {
	mysqlDb := NewMysqlDataBase(dbUrl)
	mysqlDb.Open()
	rows, err := mysqlDb.Query("select user_name, deleted_at, id from users limit 2")
	assert.Nil(t, nil, err)
	data := rows.ToCsv(true)
	f, _ := os.Create("test.csv")
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(data)
	w.Flush()
}
