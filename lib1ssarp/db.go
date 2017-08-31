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
	"strings"
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

func (s Service) FetchAll()  {
	switch s.Database.Type {
	case "mysql":
		//TODO fetch data

	}
}

//Service end


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
	var countUpdate uint
	switch t.Database.Type {
	case "mysql":
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


//mysql
func mysqlTableExists(name string, d Database) bool {
	conn := GetConn(d)
	var tn string
	e := conn.QueryRow(`SELECT TABLE_NAME FROM information_schema.tables
								WHERE table_schema = ?  AND table_name = ? LIMIT 1`, d.Basename, name).Scan(&tn)
	detectSqlErr(e)
	return  len(tn) > 0
}


func mysqlFieldType(f Field) string {
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

	return field
}


func mysqlTableCreate(m Model, d Database) {

	query := fmt.Sprintf("CREATE TABLE %s (", m.Name)
	for n, f := range m.Fields {
		if n > 0 {
			query += ", "
		}
		query += mysqlFieldType(f)
	}
	query += fmt.Sprintf(") CHARACTER SET utf8")

	mysqlExecPureQuery(d, query)

}

func mysqlExecPureQuery(d Database, query string ) {
	fmt.Println("Exec Query: ", query)

	conn := GetConn(d)
	r, e := conn.Exec(query)
	detectSqlErr(e)
	_, e = r.RowsAffected()
	detectSqlErr(e)
}

/**

 */
func mysqlTableUpdate(m Model, d Database) uint {

	columns := mysqlTableColumns(m.Name, d)
	columnsSee := map[string]bool{}
	var countUpdate uint = 0

	for _, f := range m.Fields {

		column := columns[f.Name]
		columnsSee[column.Name] = true
		if column.Name != f.Name {
			fmt.Println("Add new column: ", f.Name)
			query := `ALTER TABLE ` + d.Basename + `.` + m.Name + ` ADD COLUMN ` + mysqlFieldType(f)
			mysqlExecPureQuery(d, query)
			countUpdate ++
			continue
		}

		if column.Type != f.Type || column.Length != f.Length {

			fmt.Println("Modified column: ", f.Name, ", Type: ", f.Type, ", Length: ", f.Length)

			query := `ALTER TABLE ` + d.Basename + `.` + m.Name + ` MODIFY ` +  mysqlFieldType(f)
			mysqlExecPureQuery(d, query)
			countUpdate ++
			continue
		}
	}

	//TODO Perhaps, need safe delete columns
	for name, column := range columns {
		if !columnsSee[name] {
			fmt.Println("Drop column: ", column.Name)

			query := `ALTER TABLE ` + d.Basename + `.` + m.Name + ` DROP COLUMN ` +  column.Name
			mysqlExecPureQuery(d, query)
			countUpdate ++
		}
	}

	return countUpdate
}


func mysqlTableColumns(name string, d Database) map[string]Column {
	conn := GetConn(d)

	r, e := conn.Query(`SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, EXTRA, CHARACTER_MAXIMUM_LENGTH
								FROM information_schema.COLUMNS
								WHERE table_schema = ?  AND table_name = ?`, d.Basename, name)



	detectSqlErr(e)

	defer r.Close()

	var cname, dtype, ctype, extra string
	var length uint

	m := map[string]Column{}

	for r.Next() {

		e = r.Scan(&cname, &dtype, &ctype, &extra, &length)
		detectSqlErr(e)

		column := Column{}
		column.Name = cname
		column.Autoincrement = extra == "auto_increment"
		column.Length  = length

		//TODO move once func, see  mysqlTableCreate
		switch dtype {
		case "varchar":
			column.Type = "string"
		case "int":
			if strings.Contains(ctype, "unsigned") {
				column.Type = "uint"
			} else {
				column.Type = "int"
			}
		default:
			panic(fmt.Errorf("Error match type: %s", ctype))
		}

		m[column.Name] = column
	}

	return m
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