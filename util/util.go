package util

import (
	"bytes"
	"fmt"
	"strings"
)

func SnakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func MakeFirstLowerCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]
	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func MakeFirstUpperCase(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}
	bts := []byte(s)
	lc := bytes.ToUpper([]byte{bts[0]})
	rest := bts[1:]
	return string(bytes.Join([][]byte{lc, rest}, nil))
}

type JavaType struct {
	Type, Imports string
}

func DbTypeToJavaType(dbType string) (javaType JavaType) {
	fmt.Println(dbType)
	if dbType == "bigint" {
		javaType.Type = "Long"
		javaType.Imports = ""
		return
	}
	if dbType == "character varying" || dbType == "text" || dbType == "character" {
		javaType.Type = "String"
		javaType.Imports = ""
		return
	}
	if dbType == "numeric" {
		javaType.Type = "BigDecimal"
		javaType.Imports = "import java.math.BigDecimal;\n"
		return
	}
	if dbType == "date" {
		javaType.Type = "LocalDate"
		javaType.Imports = "import java.time.LocalDate;\n"
		return
	}
	if dbType == "timestamp without time zone" {
		javaType.Type = "LocalDateTime"
		javaType.Imports = "import java.time.LocalDateTime;\n"
		return
	}
	if dbType == "integer" || dbType == "smallint" {
		javaType.Type = "Integer"
		javaType.Imports = ""
		return
	}
	if dbType == "double precision" {
		javaType.Type = "Double"
		javaType.Imports = ""
		return
	}
	if dbType == "real" {
		javaType.Type = "Float"
		javaType.Imports = ""
		return
	}
	if dbType == "boolean" {
		javaType.Type = "Boolean"
		javaType.Imports = ""
		return
	}
	if dbType == "bytea" {
		javaType.Type = "byte[]"
		javaType.Imports = ""
		return
	}
	return
}
