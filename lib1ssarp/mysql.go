package lib1ssarp

import (
	"fmt"
	"strings"
)

//mysql
func mysqlTableExists(name string, d Database) bool {
	conn := GetConn(d)
	var tn string
	e := conn.QueryRow(`SELECT TABLE_NAME FROM information_schema.tables
								WHERE table_schema = ?  AND table_name = ? LIMIT 1`, d.Basename, name).Scan(&tn)
	detectSqlErr(e)
	return  len(tn) > 0
}


func mysqlFieldType(f  Field) string {
	field := fmt.Sprintf("%s ", f.Name)
	switch f.Type {
	case FIELD_TYPE_UINT:
		field += "INT(10) unsigned"
	case FIELD_TYPE_INT:
		field += "INT(10)"
	case FIELD_TYPE_STRING:
		field += fmt.Sprintf("VARCHAR(%d) ", f.Length)
	default:
		panic(fmt.Errorf("Error match type: %s", f.Type))
	}

	if f.Autoincrement {
		field += " auto_increment primary key"
	}

	return field
}


func mysqlTableCreate(m  Model, d  Database) {

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

func mysqlExecPureQuery(d  Database, query string ) {
	fmt.Println("Exec Query: ", query)

	conn := GetConn(d)
	r, e := conn.Exec(query)
	detectSqlErr(e)
	_, e = r.RowsAffected()
	detectSqlErr(e)
}

/**

 */
func mysqlTableUpdate(m  Model, d  Database) uint {

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

		if column.Type != f.Type ||
			(column.Type == FIELD_TYPE_STRING && column.Length != f.Length) {
			fmt.Println("Modified column: ", f.Name, ", Type: ", f.Type, ", Length: ", f.Length, ", (Type:", column.Type, ")")

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


func mysqlTableColumns(name string, d  Database) map[string]Column {
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
			column.Type = FIELD_TYPE_STRING
		case "int":
			if strings.Contains(ctype, "unsigned") {
				column.Type = FIELD_TYPE_UINT
			} else {
				column.Type = FIELD_TYPE_INT
			}
		default:
			panic(fmt.Errorf("Error match type: %s", ctype))
		}

		m[column.Name] = column
	}

	return m
}




func mysqlFetchAll(m  Model, d  Database) []map[string]string{
	conn := GetConn(d)

	r, e := conn.Query(`SELECT * FROM  ` + SafeName(m.Name) + ``)
	detectSqlErr(e)

	defer r.Close()

	var result []map[string]string

	for r.Next() {

		row := make(map[string]string, len(m.Fields))
		result = append(result, row)

		values := make([]string , len(m.Fields))
		tmp := make([]interface{}, len(m.Fields))
		for i, _ := range values {
			tmp[i] = &values[i]
		}
		e = r.Scan(tmp...)
		detectSqlErr(e)

		for i, f := range m.Fields{
			row[f.Name] = values[i]
		}
	}

	return result
}

func mysqlFetchOne(id string, m  Model, d   Database) map[string]string{
	conn := GetConn(d)

	r, e := conn.Query(`SELECT * FROM  ` + SafeName(m.Name) + ` WHERE id = ?`, id)
	detectSqlErr(e)

	defer r.Close()


	for r.Next() {

		row := make(map[string]string, len(m.Fields))

		values := make([]string , len(m.Fields))
		tmp := make([]interface{}, len(m.Fields))
		for i, _ := range values {
			tmp[i] = &values[i]
		}
		e = r.Scan(tmp...)
		detectSqlErr(e)

		for i, f := range m.Fields{
			row[f.Name] = values[i]
		}

		return row
	}

	return nil
}

/**

 */
func mysqlCreateModel(data map[string]interface{}, m  Model, d   Database) uint {

	fmt.Println(data)

	var fields, values string
	var args []interface{}
	num := 0

	for key, value := range data {
		if num > 0 {
			fields += ","
			values += ","
		}
		num ++
		fields += key
		values += "?"
		args = append(args, value)
	}

	fmt.Println(`INSERT INTO ` + SafeName(m.Name) + ` (` + fields + `) VALUES (` + values + `)`)
	conn := GetConn(d)
	r, e := conn.Exec(`INSERT INTO ` + SafeName(m.Name) + ` (` + fields + `) VALUES (` + values + `)`, args...)
	detectSqlErr(e)

	id, e := r.LastInsertId()
	detectSqlErr(e)

	return uint(id)
}

/**

 */
func mysqlUpdateModel(id string, data map[string]interface{}, m  Model, d   Database) bool {
	fmt.Println(data)

	var fields string
	var args []interface{}
	num := 0

	for key, value := range data {
		if num > 0 {
			fields += ","
		}
		num ++
		fields += key
		fields += "= ?"
		args = append(args, value)
	}

	fmt.Println(`UPDATE ` + SafeName(m.Name) + ` SET ` + fields + ` WHERE id = ?`)
	args = append(args, id)

	conn := GetConn(d)
	r, e := conn.Exec(`UPDATE ` + SafeName(m.Name) + ` SET ` + fields + ` WHERE id = ?`, args...)
	detectSqlErr(e)

	a, e := r.RowsAffected()
	detectSqlErr(e)

	return a > 0
}


func mysqlDeleteModel(id string,  m  Model, d   Database) bool {
	conn := GetConn(d)

	r, e := conn.Exec(`DELETE FROM  ` + SafeName(m.Name) + ` WHERE id = ?`, id)
	detectSqlErr(e)

	a, e := r.RowsAffected()
	detectSqlErr(e)

	return a > 0
}