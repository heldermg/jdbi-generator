package sqlaudit

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/pojo"
)

func MakeAuditSchemaSqlFile(tables []pojo.Table) {
	if _, err := os.Stat(constants.GENERATED_FOLDER); os.IsNotExist(err) {
		os.Mkdir(constants.GENERATED_FOLDER, os.ModePerm)
	}

	if _, err := os.Stat(constants.SQL_FOLDER); os.IsNotExist(err) {
		os.Mkdir(constants.SQL_FOLDER, os.ModePerm)
	}

	arquivo, err := os.OpenFile(
		constants.SQL_FOLDER+"create-schema-audit.sql",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	create := ""
	create += "create schema if not exists " + constants.DB_SCHEMA + "_audit authorization postgres;\n\n"
	for i := 0; i < len(tables); i++ {
		table := tables[i]
		create += "create table " + constants.DB_SCHEMA + "_audit." + table.Name + "_audit (\n"

		isAuditPkIncluded := false
		columns := getColumns(&table)
		for j := 0; j < len(columns); j++ {
			column := columns[j]

			if column.IsPrimaryKey.Valid {
				create += "    " + column.Name + " bigserial " + getNull(column.IsNullable) + ", \n"

			} else {
				if !isAuditPkIncluded {
					isAuditPkIncluded = true
					create += "    " + table.Name[strings.Index(table.Name, "_")+1:len(table.Name)] + "_fk bigint not null, \n"
				}
				create += "    " + column.Name + " " + getDataType(column) + " " + getNull(column.IsNullable) + ", \n"
			}

			if j == len(columns)-1 {
				create += "    constraint " + table.Name + "_audit_pk primary key (id)\n"
			}
		}
		create += ");\n\n"

		isAuditPkIncluded = false
		comment := ""
		for j := 0; j < len(columns); j++ {
			column := columns[j]
			comment = getColumnComment(table.Name, column.Name, true)

			create += "comment on column " + constants.DB_SCHEMA + "_audit." + table.Name + "_audit." + column.Name + " is '" + comment + "';\n"
			if !isAuditPkIncluded {
				isAuditPkIncluded = true
				create += "comment on column " + constants.DB_SCHEMA + "_audit." + table.Name + "_audit." +
					table.Name[strings.Index(table.Name, "_")+1:len(table.Name)] + "_fk is '" + getColumnComment(table.Name, "id", false) + "';\n"
			}
		}
		create += "\n"
	}

	arquivo.WriteString(create)
	arquivo.Close()
}

func getColumnComment(tableName string, columnName string, isAuditTable bool) (columnComment string) {
	columnComment = ""
	if columnName == constants.AUDIT_DATETIME_COLUMN_NAME {
		columnComment = constants.AUDIT_DATETIME_COLUMN_COMMENT

	} else if columnName == constants.AUDIT_OPERATION_COLUMN_NAME {
		columnComment = constants.AUDIT_OPERATION_COLUMN_COMMENT

	} else if columnName == constants.AUDIT_USER_COLUMN_NAME {
		columnComment = constants.AUDIT_USER_COLUMN_COMMENT

	} else if columnName == constants.VERSION_COLUMN_NAME {
		columnComment = constants.VERSION_COLUMN_COMMENT

	} else if columnName == "id" && isAuditTable {
		columnComment = constants.AUDIT_PK_COLUMN_COMMENT

	} else {
		columnComment = getColumnCommentFromOriginalTable(tableName, columnName)
	}
	return
}

func getNull(isNullable string) (sqlNullable string) {
	sqlNullable = "not null"
	if isNullable == "YES" {
		sqlNullable = "null"
	}
	return
}

func getDataType(column pojo.Column) (sqlDataType string) {
	sqlDataType = column.DataType
	if column.DataType == "character varying" {
		sqlDataType = column.DataType + "(" + strconv.FormatInt(column.DataLength.Int64, 10) + ")"
	}
	return
}

func getColumnCommentFromOriginalTable(tableName string, columnName string) (columnComment string) {
	columnComment = ""
	db, err := sql.Open("postgres", constants.DB_URL)
	defer db.Close()

	rows, err := db.Query(
		`select pd.description
		 from pg_description pd, pg_class pc, pg_attribute pa
     where pc.relname = '` + tableName + `'
     	and pa.attname = '` + columnName + `'
      and pa.attrelid = pc.oid
			and pd.objoid = pc.oid
			and pd.objsubid = pa.attnum`)
	if err != nil {
		fmt.Println("Failed to execute columns properties query!", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&columnComment)
	}
	return
}

func getColumns(table *pojo.Table) (columns []pojo.Column) {
	db, err := sql.Open("postgres", constants.DB_URL)
	defer db.Close()

	rows, err := db.Query(
		`select 
				column_name, is_nullable, data_type, character_maximum_length, 
				(select 'YES' from information_schema.table_constraints tco 
					join information_schema.key_column_usage kcu on kcu.constraint_name = tco.constraint_name 
					and kcu.constraint_schema = tco.constraint_schema and kcu.constraint_name = tco.constraint_name 
					where tco.constraint_type = 'PRIMARY KEY' and kcu.column_name = c.column_name and kcu.table_name = c.table_name) 
				as is_primary_key 
		from information_schema.columns c where c.table_name = '` + table.Name + `'`)
	if err != nil {
		fmt.Println("Failed to execute columns properties query!", err)
		return
	}
	defer rows.Close()

	pkCount := 0
	var column pojo.Column
	for rows.Next() {
		err := rows.Scan(&column.Name, &column.IsNullable, &column.DataType, &column.DataLength, &column.IsPrimaryKey)
		if err != nil {
			fmt.Println(err)
		}
		columns = append(columns, column)

		if column.IsPrimaryKey.Valid {
			pkCount++
		}
		if pkCount > 1 {
			table.HasMultiplePK = true
		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return
}
