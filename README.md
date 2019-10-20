# gorecords

## features
- exec and query raw sqls
- export records or a single record to json
- export records or a single record to map
- export records or a single record to slices (the same order as columns)
- export records or a single record to csv (the same order as columns)

## examples
```go
    mysqlDb := NewMysqlDataBase("user:password@tcp(localhost:3306)/dbname?tls=skip-verify&autocommit=true")
    mysqlDb.Open()
    table_names, err := mysqlDb.GetTableNames() // return table names under the database such as []string{"logs", "users"}
    res, err := mysqlDb.Exec("insert into users (`user_name`, `password`) values ('gensword', 'gensword')") // insert a reacord, res.RowsAffected = 1
    rows, err := mysqlDb.Query("select * from users")
    js, err := rows.ToJson() // get records json such as `[{"created_at":"2019-08-16 18:07:10","deleted_at":null,"id":1,"updated_at":"2019-08-16 18:07:10","user_name":"gensword", "password":"gensword"}]`
    singleRecordJs, err := rows[0].ToJson() // get single record json such as `{"created_at":"2019-08-16 18:07:10","deleted_at":null,"id":1,"updated_at":"2019-08-16 18:07:10","user_name":"gensword", "password":"gensword"}`
    RecordsMaps := rows.ToMaps() // get records maps such as `[]map[string]interface{}{map[string]interface{}{"deleted_at": nil, "user_name": "gensword", "created_at": "2019-08-16 18:07:10", "updated_at": "2019-08-16 18:07:10", "id": int64(1)}}`
    singleRecordMap := rows[0].ToMap() // get single record map such as `map[string]interface{}{"deleted_at": nil, "user_name": "gensword", "created_at": "2019-08-16 18:07:10", "updated_at": "2019-08-16 18:07:10", "id": int64(1)}`
    rows[0].ToSlice() // get a slice of single record coulmns' values ordered such as  `[]interface{}{"2019-08-16", nil, int64(1), "2019-08-16 18:07:10", "gensword", "gensword"}`
    rows.ToSlices() // get slices of records coulmns' values ordered such as `[][]interface{}{{"2019-08-16", nil, int64(1), "2019-08-16 18:07:10", "gensword", "gensword"}}`
    
    data := rows.ToCsv(true)  // the first line of the target csv file will be columns' name if true else only records of table will be written to the file.
    f, _ := os.Create("test.csv")
    defer f.Close()
    w := csv.NewWriter(f)
    w.WriteAll(data)
    w.Flush()
```