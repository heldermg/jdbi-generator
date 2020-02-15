package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/pojo"
	"github.com/heldermg/jdbi-generator/repository"
	"github.com/heldermg/jdbi-generator/sqlaudit"
	"github.com/heldermg/jdbi-generator/sqlstg"
	"github.com/heldermg/jdbi-generator/trigger"
	"github.com/heldermg/jdbi-generator/util"
	_ "github.com/lib/pq"
)

var TABLES_FILTER []string

func main() {
	db, err := sql.Open("postgres", constants.DB_URL)
	defer db.Close()

	if err != nil {
		fmt.Println("Failed to connect database!", err)
		return
	}

	tables := getTables(db)
	for {
		showMenu()
		option := getOption()

		makeSqlAuditSchema := false
		makeSqls := false
		makeRepositories := false
		makePojos := false
		makeTriggers := false

		switch option {
		case 0:
			os.Exit(0)

		case 1:
			TABLES_FILTER = getTableFilter()
			tables = getTables(db)
			continue

		case 2:
			makeSqlAuditSchema = true
			makeSqls = true
			makeRepositories = true
			makePojos = true
			makeTriggers = true

		case 3:
			makeSqls = true

		case 4:
			makeRepositories = true

		case 5:
			makePojos = true

		case 6:
			makeTriggers = true

		case 7:
			makeSqlAuditSchema = true

		default:
			fmt.Println("Option not found!")
			fmt.Println("")
			continue
		}

		if makeSqlAuditSchema {
			sqlaudit.MakeAuditSchemaSqlFile(tables)
		}

		if makeRepositories {
			repository.MakeAbstractClassFiles()
		}

		if makePojos {
			pojo.MakeAbstractClassFiles()
		}

		if makeTriggers {
			trigger.MakeGenericAuditTrigger()
			trigger.MakeGenericConcurrentTrigger()
		}

		for i := 0; i < len(tables); i++ {
			table := tables[i]
			class := util.SnakeCaseToCamelCase(table.Name)

			columns := getColumns(db, &table)

			if _, err := os.Stat(constants.GENERATED_FOLDER); os.IsNotExist(err) {
				os.Mkdir(constants.GENERATED_FOLDER, os.ModePerm)
			}

			if makeTriggers {
				trigger.MakeTableAuditTrigger(table)
				trigger.MakeTableConcurrentTrigger(table)
			}

			// Generate <class>Repository.sql.stg files
			if makeSqls {
				sqlstg.MakeSqlStgFile(table, columns, class)
			}

			// Generate <class>Repository.java files
			if makeRepositories {
				repository.MakeRepositoryFile(table, class)
			}

			// Generate <class>.java (pojo) files
			if makePojos {
				pojo.MakePojoFile(table, class, columns)
			}
		}
	}
}

func getTableFilter() []string {
	var f string
	fmt.Print("Tables (comma-separated): ")
	fmt.Scan(&f)
	fmt.Println("")
	return strings.Split(f, ",")
}

func showMenu() {
	fmt.Println("0- Exit")
	fmt.Println("1- Filter by table (Current: " + strings.Join(TABLES_FILTER, ", ") + ")")
	fmt.Println("2- Generate all")
	fmt.Println("3- Generate only *Repository.sql.stg files")
	fmt.Println("4- Generate only *Repository.java files")
	fmt.Println("5- Generate only pojo *.java files")
	fmt.Println("6- Generate audit and concurrent triggers")
	fmt.Println("7- Generate audit schema script")
}

func getOption() int {
	var option int
	fmt.Print("Option chosen: ")
	fmt.Scan(&option)
	fmt.Println("")
	return option
}

func getTables(db *sql.DB) (tables []pojo.Table) {
	query := `select table_name from information_schema.tables where
		table_schema = '` + constants.DB_SCHEMA + `'
		and table_name not like 'jhi%'
		and table_name not like 'database%'`

	if len(TABLES_FILTER) > 0 {
		query += " and table_name in ("
		for i := 0; i < len(TABLES_FILTER); i++ {
			if i != 0 {
				query += ","
			}
			query += "'" + TABLES_FILTER[i] + "'"
		}
		query += ")"
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to execute table names query!", err)
		return
	}
	defer rows.Close()

	var table pojo.Table
	for rows.Next() {
		err := rows.Scan(&table.Name)
		if err != nil {
			fmt.Println(err)
		}
		tables = append(tables, table)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func getColumns(db *sql.DB, table *pojo.Table) (columns []pojo.Column) {
	rows, err := db.Query(
		`select 
				column_name, is_nullable, data_type, 
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
		err := rows.Scan(&column.Name, &column.IsNullable, &column.DataType, &column.IsPrimaryKey)
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
