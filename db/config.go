package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-jet/jet/v2/postgres"
)

type DBConfigDef struct {
	POSTGRES_DB       string
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	LOG_QUERIES       *bool
}

func (c *DBConfigDef) ConnectString() string {
	return "host=" + c.POSTGRES_HOST + " port=" + c.POSTGRES_PORT + " user=" + c.POSTGRES_USER + " password=" + c.POSTGRES_PASSWORD + " dbname=" + c.POSTGRES_DB + " sslmode=disable"
}

func (c *DBConfigDef) URI() string {
	return "postgres://" + c.POSTGRES_USER + ":" + c.POSTGRES_PASSWORD + "@" + c.POSTGRES_HOST + ":" + c.POSTGRES_PORT + "/" + c.POSTGRES_DB + "?sslmode=disable"
}

func (c *DBConfigDef) EnvLines() []string {
	return []string{
		"POSTGRES_DB=" + c.POSTGRES_DB,
		"POSTGRES_PASSWORD=" + c.POSTGRES_PASSWORD,
		"POSTGRES_USER=" + c.POSTGRES_USER,
		"POSTGRES_HOST=" + c.POSTGRES_HOST,
		"POSTGRES_PORT=" + c.POSTGRES_PORT,
	}
}

func (c *DBConfigDef) SetLogQueries() {
	println("Setting log queries")
	postgres.SetQueryLogger(func(ctx context.Context, queryInfo postgres.QueryInfo) {
		sql, args := queryInfo.Statement.Sql()
		fmt.Printf("- SQL: %s Args: %v \n", sql, args)
		fmt.Printf("- Debug SQL: %s \n", queryInfo.Statement.DebugSql())

		// Depending on how the statement is executed, RowsProcessed is:
		//   - Number of rows returned for Query() and QueryContext() methods
		//   - RowsAffected() for Exec() and ExecContext() methods
		//   - Always 0 for Rows() method.
		fmt.Printf("- Rows processed: %d\n", queryInfo.RowsProcessed)
		fmt.Printf("- Duration %s\n", queryInfo.Duration.String())
		// fmt.Printf("- Execution error: %v\n", queryInfo.Err)

		callerFile, callerLine, callerFunction := queryInfo.Caller()
		fmt.Printf("- Caller file: %s, line: %d, function: %s\n", callerFile, callerLine, callerFunction)
	})
}

func ReadConfigFromEnv() DBConfigDef {
	logQueries := os.Getenv("LOG_QUERIES") == "true"
	conf := DBConfigDef{
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		LOG_QUERIES:       &logQueries,
	}
	if conf.LOG_QUERIES != nil && *conf.LOG_QUERIES {
		conf.SetLogQueries()
	}
	return conf
}
