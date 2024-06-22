package server

import (
	"GoSimProject/lib"
	"database/sql"
	"fmt"
	"time"
)

type DBConfig struct {
	DriverName string
	User       string
	Passwd     string
	Addr       string
	DBName     string
	Port       string
	QuerySql   string
	UpdateSql  string
}

type DB struct {
	db   DBConfig
	call DBConfig
}

// NewMySqlDB 配置数据库
func NewMySqlDB(db DBConfig) *sql.DB {
	var err error
	funcName := "NewMySqlDB"
	var dbClient *sql.DB
	dbClient, err = sql.Open(
		db.DriverName, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.User, db.Passwd, db.Addr, db.Port, db.DBName))

	if err != nil {
		fmt.Printf("%s err(%+v)\n", funcName, err)
		return nil
	}
	dbClient.SetMaxOpenConns(25)                 // 设置最大打开的连接数
	dbClient.SetMaxIdleConns(25)                 // 设置最大空闲连接数
	dbClient.SetConnMaxLifetime(5 * time.Minute) // 设置连接的最大生存时间
	return dbClient
}

func (db DBConfig) WriteBack(resp string) {
	var err error
	funcName := "db.WriteBack"
	updateSql := lib.JsHandler(config.Get("callJobScriptTemplate").(string), resp)
	dbClient := NewMySqlDB(db)
	//批量更新开票数据状态
	_, err = dbClient.Exec(updateSql)
	if err != nil {
		fmt.Printf("%s err(%+v) ID(%+v)\n", funcName, err)
	}
}

func (db DBConfig) FetchData() []map[string]interface{} {
	funName := "db.FetchData"
	var (
		err error
	)

	dbClient := NewMySqlDB(db)
	defer dbClient.Close()
	query := db.QuerySql
	var results []map[string]interface{}
	results, err = QueryAndProcess(dbClient, query)
	if err != nil {
		fmt.Printf("%s err(%+v)", funName, err.Error())
		return nil
	}
	return results
}

// QueryAndProcess 执行查询并处理结果
func QueryAndProcess(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	var err error
	var rows *sql.Rows
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	cols, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i, col := range cols {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			result[col] = v
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
