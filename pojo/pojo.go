package pojo

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/util"
)

type Column struct {
	Name         string
	IsNullable   string
	DataType     string
	IsPrimaryKey sql.NullString
}

type Table struct {
	Name          string
	HasMultiplePK bool
}

func MakeAbstractClassFiles() {
	if _, err := os.Stat(constants.GENERATED_FOLDER); os.IsNotExist(err) {
		os.Mkdir(constants.GENERATED_FOLDER, os.ModePerm)
	}

	os.Mkdir(constants.DOMAIN_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.DOMAIN_FOLDER+"AbstractPojoAuditVersion.java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	packageImports := "package " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".domain;\n\n"
	packageImports += "import java.time.LocalDateTime;\n"
	packageImports += "import org.jdbi.v3.core.mapper.reflect.ColumnName;\n"

	pojoClass := "\npublic abstract class AbstractPojoAuditVersion implements Serializable {\n\n"
	pojoClass += "    private static final long serialVersionUID = 1L;\n\n"

	attName := util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(constants.AUDIT_LOGIN_COLUMN_NAME))
	pojoClass += "    @ColumnName(\"" + constants.AUDIT_LOGIN_COLUMN_NAME + "\")\n"
	pojoClass += "    private String " + attName + " = \"admim\";\n\n"
	getterSetter := getGettersSetters(attName, "String")

	attName = util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(constants.AUDIT_DATETIME_COLUMN_NAME))
	pojoClass += "    @ColumnName(\"" + constants.AUDIT_DATETIME_COLUMN_NAME + "\")\n"
	pojoClass += "    private LocalDateTime " + attName + " = LocalDateTime.now();\n\n"
	getterSetter += getGettersSetters(attName, "LocalDateTime")

	attName = util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(constants.AUDIT_OPERATION_COLUMN_NAME))
	pojoClass += "    @ColumnName(\"" + constants.AUDIT_OPERATION_COLUMN_NAME + "\")\n"
	pojoClass += "    private String " + attName + " = \"INSERCAO\";\n\n"
	getterSetter += getGettersSetters(attName, "String")

	attName = util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(constants.VERSION_COLUMN_NAME))
	pojoClass += "    @ColumnName(\"" + constants.VERSION_COLUMN_NAME + "\")\n"
	pojoClass += "    private Integer " + attName + " = 0;\n\n"
	getterSetter += getGettersSetters(attName, "Integer")

	pojoClass += getterSetter
	pojoClass += "}\n"

	arquivo.WriteString((packageImports + pojoClass))
	arquivo.Close()
}

func MakePojoFile(table Table, className string, columns []Column) {
	os.Mkdir(constants.DOMAIN_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.DOMAIN_FOLDER+className+".java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(getPojoClass(table, className, columns))
	arquivo.Close()
}

func getPojoClass(table Table, class string, columns []Column) string {
	packageImports := "package " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".domain;\n\n"
	packageImports += "import javax.validation.constraints.NotNull;\n"
	packageImports += "import org.jdbi.v3.core.mapper.reflect.ColumnName;\n"
	packageImports += "import org.springframework.data.annotation.Id;\n"

	pojoClass := "\npublic class " + class + " extends AbstractPojoAuditVersion {\n\n"
	pojoClass += "    private static final long serialVersionUID = 1L;\n\n"

	getterSetter := ""
	for i := 0; i < len(columns); i++ {
		column := columns[i]

		if i == 0 {
			pojoClass += "    @Id\n"

			if table.HasMultiplePK {
				pojoClass += "    private " + class + "Id id;\n\n"
				getterSetter += getGettersSetters("id", class+"Id")
				//gerar classe '<class>Id.java'
				getMultiplePkClass(class+"Id", columns)
				continue
			}
		}

		if (!column.IsPrimaryKey.Valid ||
			(column.IsPrimaryKey.Valid && !table.HasMultiplePK)) &&
			column.Name != constants.AUDIT_DATETIME_COLUMN_NAME &&
			column.Name != constants.AUDIT_OPERATION_COLUMN_NAME &&
			column.Name != constants.AUDIT_LOGIN_COLUMN_NAME &&
			column.Name != constants.VERSION_COLUMN_NAME {

			javaType := util.DbTypeToJavaType(column.DataType)
			attName := util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(column.Name))

			pojoClass += getNotNullAnnotation(column.IsNullable)
			pojoClass += "    @ColumnName(\"" + column.Name + "\")\n"
			pojoClass += "    private " + javaType.Type + " " + attName + ";\n\n"

			if !strings.Contains(packageImports, javaType.Imports) {
				packageImports += javaType.Imports
			}
			getterSetter += getGettersSetters(attName, javaType.Type)
		}
	}

	pojoClass += getterSetter
	pojoClass += "}\n"
	return (packageImports + pojoClass)
}

func getMultiplePkClass(class string, columns []Column) {
	packageImports := "package " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".domain;\n\n"
	packageImports += "import java.io.Serializable;\n"
	packageImports += "import javax.validation.constraints.NotNull;\n"
	packageImports += "import org.jdbi.v3.core.mapper.reflect.ColumnName;\n"
	packageImports += "import org.springframework.data.annotation.Id;\n"

	pojoClass := "\npublic class " + class + " implements Serializable {\n\n"
	pojoClass += "    private static final long serialVersionUID = 1L;\n\n"

	getterSetter := ""
	for i := 0; i < len(columns); i++ {
		column := columns[i]

		if column.IsPrimaryKey.Valid {
			javaType := util.DbTypeToJavaType(column.DataType)
			attName := util.MakeFirstLowerCase(util.SnakeCaseToCamelCase(column.Name))

			pojoClass += getNotNullAnnotation(column.IsNullable)
			pojoClass += "    @ColumnName(\"" + column.Name + "\")\n"
			pojoClass += "    private " + javaType.Type + " " + attName + ";\n\n"

			if !strings.Contains(packageImports, javaType.Imports) {
				packageImports += javaType.Imports
			}
			getterSetter += getGettersSetters(attName, javaType.Type)
		}
	}
	pojoClass += getterSetter
	pojoClass += "}\n"
	finalClass := (packageImports + pojoClass)

	arquivo, err := os.OpenFile(
		constants.DOMAIN_FOLDER+class+".java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(finalClass)
}

func getNotNullAnnotation(isNullable string) (annotation string) {
	annotation = ""
	if isNullable == "NO" {
		annotation = "    @NotNull\n"
	}
	return
}

func getGettersSetters(attName string, javaType string) (getterSetter string) {
	getterSetter = "    public " + javaType + " get" + util.MakeFirstUpperCase(attName) + "() {\n"
	getterSetter += "        return this." + attName + ";\n"
	getterSetter += "    }\n\n"
	getterSetter += "    public void set" + util.MakeFirstUpperCase(attName) +
		"(" + javaType + " " + attName + ") {\n"
	getterSetter += "        this." + attName + " = " + attName + ";\n"
	getterSetter += "    }\n\n"
	return
}
