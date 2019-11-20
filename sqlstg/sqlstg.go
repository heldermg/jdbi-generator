package sqlstg

import (
	"fmt"
	"os"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/pojo"
	"github.com/heldermg/jdbi-generator/util"
)

func MakeSqlStgFile(table pojo.Table, columns []pojo.Column, className string) {
	className += "Repository"
	os.Mkdir(constants.SQL_STG_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.SQL_STG_FOLDER+className+".sql.stg",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
	arquivo.WriteString(getInitialComments(className))
	if table.HasMultiplePK {
		arquivo.WriteString(getByComplexId(table.Name, columns))
		arquivo.WriteString(findAll(table.Name))
		arquivo.WriteString(save(table.Name, columns))
		arquivo.WriteString(deleteByComplexId(table.Name, columns))
	} else {
		arquivo.WriteString(getById(table.Name))
		arquivo.WriteString(findAll(table.Name))
		arquivo.WriteString(save(table.Name, columns))
		arquivo.WriteString(update(table.Name, columns))
		arquivo.WriteString(deleteById(table.Name))
	}
	arquivo.Close()
}

func getInitialComments(className string) string {
	return "<! " + className + " !>\n\n"
}

func getById(table string) string {
	sql := "findById(id) ::= <% \n" +
		"    select \n" +
		"        * \n" +
		"    from \n" +
		"        " + table + " \n" +
		"    where \n" +
		"        id = :id \n" +
		"%>\n\n"
	return sql
}

func getByComplexId(table string, columns []pojo.Column) string {
	sql := "findByComplexId(id) ::= <% \n" +
		"    select \n" +
		"        * \n" +
		"    from \n" +
		"        " + table + " \n" +
		"    where \n"

	connector := ""
	for i := 0; i < len(columns); i++ {
		column := columns[i]
		if column.IsPrimaryKey.Valid {
			sql += "        " + connector + column.Name + " = :id." +
				util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(column.Name)) + "\n"
		}
		if i != len(columns)-1 {
			connector += "and "
		}
	}

	sql += "%>\n\n"
	return sql
}

func findAll(table string) string {
	sql := "findAll() ::= <% \n" +
		"    select \n" +
		"        * \n" +
		"    from \n" +
		"        " + table + " \n" +
		"%>\n\n"
	return sql
}

func deleteById(table string) string {
	sql := "deleteById(id) ::= <% \n" +
		"    delete from " + table + " \n" +
		"    where \n" +
		"        id = :id \n" +
		"%>\n\n"
	return sql
}

func deleteByComplexId(table string, columns []pojo.Column) string {
	sql := "deleteByComplexId(id) ::= <% \n" +
		"    delete from " + table + " \n" +
		"    where \n"

	connector := ""
	for i := 0; i < len(columns); i++ {
		column := columns[i]
		if column.IsPrimaryKey.Valid {
			sql += "        " + connector + column.Name + " = :id." +
				util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(column.Name)) + "\n"
		}
		if i != len(columns)-1 {
			connector += "and "
		}
	}

	sql += "%>\n\n"
	return sql
}

func save(table string, columns []pojo.Column) string {
	sql := "save(entity) ::= <% \n" +
		"    insert into " + table + " ( \n"

	// Se o primeiro campo for 'id', o mesmo possui sequence então não precisa incluir no save
	posicaoInicial := 0
	if columns[0].Name == "id" {
		posicaoInicial = 1
	}

	for i := posicaoInicial; i < len(columns); i++ {
		sql += "        " + columns[i].Name
		if i != len(columns)-1 {
			sql += ", \n"
		} else {
			sql += " \n"
		}
	}
	sql += "    ) values ( \n"

	for i := posicaoInicial; i < len(columns); i++ {
		sql += getVersionColumnSql(columns[i].Name, false)
		sql += getAuditDatetimeColumnSql(columns[i].Name, false)
		sql += getAuditOperationColumnSql(columns[i].Name, false)
		sql += getOtherColumnsSql(columns[i].Name, false)

		if i != len(columns)-1 {
			sql += ", \n"
		} else {
			sql += " \n"
		}
	}

	sql += "    ) \n"
	sql += "%>\n\n"
	return sql
}

func update(table string, columns []pojo.Column) string {
	sql := "update(entity) ::= <% \n" +
		"    update " + table + " set \n"

	// Se o primeiro campo for 'id', ignora chave
	for i := 0; i < len(columns); i++ {
		if !columns[i].IsPrimaryKey.Valid {

			sql += getVersionColumnSql(columns[i].Name, true)
			sql += getAuditDatetimeColumnSql(columns[i].Name, true)
			sql += getAuditOperationColumnSql(columns[i].Name, true)
			sql += getOtherColumnsSql(columns[i].Name, true)

			if i != len(columns)-1 {
				sql += ", \n"
			} else {
				sql += " \n"
			}
		}
	}
	sql += "    where id = :entity.id \n"
	sql += "%>\n\n"
	return sql
}

func getVersionColumnSql(columnName string, isUpdateMathod bool) string {
	sql := ""
	if columnName == constants.VERSION_COLUMN_NAME {
		sql += "        "
		if isUpdateMathod {
			sql += columnName + " = (:entity." + constants.VERSION_COLUMN_NAME + " + 1)"
		} else {
			sql += "0"
		}
	}
	return sql
}

func getAuditDatetimeColumnSql(columnName string, isUpdateMathod bool) string {
	sql := ""
	if columnName == constants.AUDIT_DATETIME_COLUMN_NAME {
		sql += "        "
		if isUpdateMathod {
			sql += columnName + " = "
		}
		sql += "now()"
	}
	return sql
}

func getAuditOperationColumnSql(columnName string, isUpdateMethod bool) string {
	sql := ""
	if columnName == constants.AUDIT_OPERATION_COLUMN_NAME {
		sql += "        "
		if isUpdateMethod {
			sql += columnName + " = 'UPDATE'"
		} else {
			sql += "'INSERT'"
		}
	}
	return sql
}

func getOtherColumnsSql(columnName string, isUpdate bool) string {
	sql := ""
	if columnName != constants.VERSION_COLUMN_NAME &&
		columnName != constants.AUDIT_DATETIME_COLUMN_NAME &&
		columnName != constants.AUDIT_OPERATION_COLUMN_NAME {
		sql += "        "
		if isUpdate {
			sql += columnName + " = "
		}
		sql += ":entity." + util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(columnName))
	}
	return sql
}
