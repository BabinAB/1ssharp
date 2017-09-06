package lib1ssarp

import (
	"fmt"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"os"
	"log"
)

var db *sql.DB
var err error
var isOpen bool = false

const DB_TYPE_MYSQL = "mysql"
const DB_TYPE_POSTGRES = "postgres"
const DB_TYPE_MSSQL = "mssql"
const DB_TYPE_SQLITE = "sqlite"

const FIELD_TYPE_UINT = "uint"
const FIELD_TYPE_INT = "int"
const FIELD_TYPE_STRING = "string"

func init() {
	fmt.Println("Init db...")
}

/**
Open connection
 */
func GetConn(d Database) *sql.DB {

	if isOpen {
		return db
	}

	var strConn string

	if d.Type != DB_TYPE_MYSQL {
		//TODO not yet
		panic(fmt.Errorf("While the system only supports - mysql driver"))
		os.Exit(1)
	}
	strConn = fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", d.Username, d.Password, d.Basename)

	fmt.Println("String connection: ", strConn)

	db, err = sql.Open(d.Type, strConn)
	detectSqlErr(err)

	fmt.Println("Connection DB has been open...")
	isOpen = true
	return db
}

/**
Close connection
 */
func CloseConn() {
	if isOpen {
		isOpen = false
		err = db.Close()
		detectSqlErr(err)
	}
}

/**

 */
func SafeName(name string) string  {
	return name
}

//Service
type Service struct {
	Database Database
	Model Model
}

func (s Service) FetchAll() []map[string]string {
	switch s.Database.Type {
	case DB_TYPE_MYSQL:
		return mysqlFetchAll(s.Model, s.Database)

	}

	return nil
}

func (s Service) FetchOne(id string) map[string]string  {
	switch s.Database.Type {
	case DB_TYPE_MYSQL:
		return mysqlFetchOne(id, s.Model, s.Database)

	}
	return nil
}


func (s Service) Create(data map[string]interface{}) uint {
	switch s.Database.Type {
	case DB_TYPE_MYSQL:
		return mysqlCreateModel(data, s.Model, s.Database)

	}
	return 0
}

//Service end


//Table
type Table struct {
	Database Database
	Model Model
}

func (t Table) Exists() bool {
	switch t.Database.Type {
	case DB_TYPE_MYSQL:
		return mysqlTableExists(t.Model.Name, t.Database)

	}
	return false
}

func (t Table) Create() {
	fmt.Println("Create table: ", t.Model.Name, ", db type: ",  t.Database.Type)
	switch t.Database.Type {
	case DB_TYPE_MYSQL:
		mysqlTableCreate(t.Model, t.Database)

	}
}

func (t Table) Update() {
	fmt.Println("Update table: ", t.Model.Name)
	var countUpdate uint
	switch t.Database.Type {
	case DB_TYPE_MYSQL:
		countUpdate = mysqlTableUpdate(t.Model, t.Database)

	}

	if countUpdate == 0 {
		fmt.Println("Nothing updated...")
	} else {
		fmt.Println("Has been updated: ", countUpdate)
	}
}

//Table end


//Column
type Column struct {
	Field

}

func (c Column) String() string {
	return fmt.Sprintf("Column{ Name: %s}", c.Name)
}

//Column end




//end mysql

func detectSqlErr(error error) {
	if err != nil {
		//TODO add to log
		switch {
		case err == sql.ErrNoRows:
			log.Printf("Empty row..." )
		default:
			log.Fatal(err)
		}
	}
}