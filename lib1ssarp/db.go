package lib1ssarp

/**
TODO add types database as constants - `mysql, pg, mssql and etc`
TODO add types database fields as constants - `uint, int, string and etc`
 */

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

	if d.Type != "mysql" {
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
		err = db.Close()
		detectSqlErr(err)
	}
}

/**

 */
func SafeName(name string) string  {
	return name
}


//Table
type Table struct {
	Database Database
	Model Model
}

func (t Table) Exists() bool {
	switch t.Database.Type {
	case "mysql":
		return mysqlTableExists(t.Model.Name, t.Database)

	}
	return false
}

func (t Table) Create() {
	fmt.Println("Create table: ", t.Model.Name, ", db type: ",  t.Database.Type)
	switch t.Database.Type {
	case "mysql":
		mysqlTableCreate(t.Model, t.Database)

	}
}

func (t Table) Update() {
	fmt.Println("Update table: ", t.Model.Name)
	switch t.Database.Type {
	case "mysql":
		mysqlTableUpdate(t.Model, t.Database)

	}
}

//Table end

//mysql
func mysqlTableExists(name string, d Database) bool {
	conn := GetConn(d)
	var tn string
	e := conn.QueryRow(`SELECT TABLE_NAME FROM information_schema.tables
								WHERE table_schema = ?  AND table_name = ? LIMIT 1`, d.Basename, name).Scan(&tn)
	detectSqlErr(e)
	return  len(tn) > 0
}


func mysqlTableCreate(m Model, d Database) {

	query := fmt.Sprintf("CREATE TABLE %s (", m.Name)
	for n, f := range m.Fields {

		//TODO set
		field := fmt.Sprintf("%s ", f.Name)
		switch f.Type {
		case "uint":
			field += "INT(10) unsigned"
		case "int":
			field += "INT(10)"
		case "string":
			field += fmt.Sprintf("VARCHAR(%d) ", f.Length)
		default:
			panic(fmt.Errorf("Error match type: %s", f.Type))
		}

		if f.Autoincrement {
			field += " auto_increment primary key"
		}

		if n > 0 {
			query += ", "
		}
		query += field
	}
	query += fmt.Sprintf(") CHARACTER SET utf8")

	fmt.Println("Exec Query: ", query)

	conn := GetConn(d)
	r, e := conn.Exec(query)
	detectSqlErr(e)
	_, e = r.RowsAffected()
	detectSqlErr(e)

}

func mysqlTableUpdate(m Model, d Database) {
	for _, f := range m.Fields {
		//TODO compare and modified
		fmt.Println(f)
	}
}


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