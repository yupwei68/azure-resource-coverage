package sqlconnect

import (
	_ "github.com/denisenkom/go-mssqldb"
	"database/sql"
	"context"
	"log"
	"fmt"
	"errors"
)

var db *sql.DB

var server = "sqlserver.database.windows.net"
var port = 1433
var user = ""
var password = ""
var database = ""

func Connect() bool{
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
		return false
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
	return true
}

func CreateCoverage(nameSpace string,typeName string, providerName string, resourceName string, opsReqPath, versions, ops string, d string, support bool) (int64, error) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("CreateCoverage: db is null")
		return -1, err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tfsupport := 0
	if support{
		tfsupport =1
	}

	tsql:= "INSERT INTO TFCoverage.Coverage (Namespace,TypeName,ProviderName,ResourceName,OperationReqPath,Versions,Operations,TFSupport,UpdateDate) VALUES (@Namespace, @TypeName, @ProviderName,@ResourceName,@OperationReqPath,@Versions,@Operations,@TFSupport,@UpdateDate); select convert(bigint, SCOPE_IDENTITY());"

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Namespace", nameSpace),
		sql.Named("TypeName", typeName),
		sql.Named("ProviderName", providerName),
		sql.Named("ResourceName", resourceName),
		sql.Named("OperationReqPath", opsReqPath),
		sql.Named("Versions", versions),
		sql.Named("Operations", ops),
		sql.Named("TFSupport", tfsupport),
		sql.Named("UpdateDate", d))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}
