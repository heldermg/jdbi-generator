package repository

import (
	"fmt"
	"os"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/pojo"
	"github.com/heldermg/jdbi-generator/util"
)

func MakeAbstractClassFiles() {
	if _, err := os.Stat(constants.GENERATED_FOLDER); os.IsNotExist(err) {
		os.Mkdir(constants.GENERATED_FOLDER, os.ModePerm)
	}

	os.Mkdir(constants.REPOSITORY_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.REPOSITORY_FOLDER+"AbstractRepository.java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	packageImports := getPackages(false)
	pojoClass := "\npublic interface AbstractRepository<T> {\n\n"
	pojoClass += getFindAll()
	pojoClass += getSave()
	pojoClass += getUpdate()
	pojoClass += "}\n"

	arquivo.WriteString((packageImports + pojoClass))
	arquivo.Close()

	os.Mkdir(constants.REPOSITORY_FOLDER, os.ModePerm)
	arquivo, err = os.OpenFile(
		constants.REPOSITORY_FOLDER+"AbstractRepositoryWithSimpleId.java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	packageImports = getPackages(false)
	pojoClass = "\npublic interface AbstractRepositoryWithSimpleId<T, I> extends AbstractRepository<T> {\n\n"
	pojoClass += getFindById(false)
	pojoClass += getDelete(false)
	pojoClass += "}\n"

	arquivo.WriteString((packageImports + pojoClass))
	arquivo.Close()

	os.Mkdir(constants.REPOSITORY_FOLDER, os.ModePerm)
	arquivo, err = os.OpenFile(
		constants.REPOSITORY_FOLDER+"AbstractRepositoryWithComplexId.java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	packageImports = getPackages(true)
	pojoClass = "\npublic interface AbstractRepositoryWithComplexId<T, I> extends AbstractRepository<T> {\n\n"
	pojoClass += getFindById(true)
	pojoClass += getDelete(true)
	pojoClass += "}\n"

	arquivo.WriteString((packageImports + pojoClass))
	arquivo.Close()
}

func getFindById(hasMultiplePK bool) (findById string) {
	findById += "    @SqlQuery\n"
	findById += "    @UseStringTemplateSqlLocator\n"
	if hasMultiplePK {
		findById += "    Optional<T> findByComplexId(@BindBean(\"id\") I id);\n\n"
	} else {
		findById += "    Optional<T> findById(@Bind(\"id\") I id);\n\n"
	}
	return
}

func getFindAll() (findAll string) {
	findAll += "    @SqlQuery\n"
	findAll += "    @UseStringTemplateSqlLocator\n"
	findAll += "    List<T> findAll();\n\n"
	return
}

func getSave() (save string) {
	save += "    @SqlUpdate(\"save\")\n"
	save += "    @GetGeneratedKeys\n"
	save += "    @UseStringTemplateSqlLocator\n"
	save += "    T saveEntity(@BindBean(\"entity\") T entity);\n\n"

	save += "    default T save(T entity) {\n"
	save += "       if (entity instanceof AbstractPojoAuditVersion) {\n"
	save += "          AbstractPojoAuditVersion pojo = (AbstractPojoAuditVersion) entity;\n"
	save += "          pojo.set" + util.MakeFirstUpperCase(util.SnakeCaseToCamelCase(constants.AUDIT_DATETIME_COLUMN_NAME)) + "(LocalDateTime.now());\n"
	save += "          pojo.set" + util.MakeFirstUpperCase(util.SnakeCaseToCamelCase(constants.AUDIT_OPERATION_COLUMN_NAME)) + "(\"INSERT\");\n\n"

	save += "          Optional<String> userLogin = getCurrentUserLogin();\n"
	save += "          if (userLogin.isPresent()) {\n"
	save += "             pojo.set" + util.MakeFirstUpperCase(util.SnakeCaseToCamelCase(constants.AUDIT_USER_COLUMN_NAME)) + "(userLogin.get());\n"
	save += "          }\n"
	save += "       }\n"
	save += "       entity = saveEntity(entity);\n"
	save += "       return entity;\n"
	save += "    }\n"
	return
}

func getUpdate() (update string) {
	update += "    @SqlUpdate(\"update\")\n"
	update += "    @GetGeneratedKeys\n"
	update += "    @UseStringTemplateSqlLocator\n"
	update += "    T updateEntity(@BindBean(\"entity\") T entity);\n\n"

	update += "    default T update(T entity) {\n"
	update += "       if (entity instanceof AbstractPojoAuditVersion) {\n"
	update += "          AbstractPojoAuditVersion pojo = (AbstractPojoAuditVersion) entity;\n"
	update += "          pojo.set" + util.MakeFirstUpperCase(util.SnakeCaseToCamelCase(constants.AUDIT_DATETIME_COLUMN_NAME)) + "(LocalDateTime.now());\n"
	update += "          pojo.set" + util.MakeFirstUpperCase(util.SnakeCaseToCamelCase(constants.AUDIT_OPERATION_COLUMN_NAME)) + "(\"UPDATE\");\n\n"

	update += "          Optional<String> userLogin = getCurrentUserLogin();\n"
	update += "          if (userLogin.isPresent()) {\n"
	update += "             pojo.set" + util.MakeFirstUpperCase(util.SnakeCaseToCamelCase(constants.AUDIT_USER_COLUMN_NAME)) + "(userLogin.get());\n"
	update += "          }\n"
	update += "       }\n"
	update += "       entity = updateEntity(entity);\n"
	update += "       return entity;\n"
	update += "    }\n"
	return
}

func getDelete(hasMultiplePK bool) (delete string) {
	delete += "    @SqlUpdate\n"
	delete += "    @UseStringTemplateSqlLocator\n"
	if hasMultiplePK {
		delete += "    void deleteByComplexId(@BindBean(\"id\") I id);\n\n"
	} else {
		delete += "    void deleteById(@Bind(\"id\") I id);\n\n"
	}
	return
}

func getPackages(hasMultiplePK bool) (packageImports string) {
	packageImports = "package " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".repository;\n\n"
	packageImports += "import java.time.LocalDateTime;\n"
	packageImports += "import java.util.List;\n"
	packageImports += "import java.util.Optional;\n\n"
	if !hasMultiplePK {
		packageImports += "import org.jdbi.v3.sqlobject.customizer.Bind;\n"
	}
	packageImports += "import org.jdbi.v3.sqlobject.customizer.BindBean;\n"
	packageImports += "import org.jdbi.v3.sqlobject.statement.GetGeneratedKeys;\n"
	packageImports += "import org.jdbi.v3.sqlobject.statement.SqlQuery;\n"
	packageImports += "import org.jdbi.v3.sqlobject.statement.SqlUpdate;\n"
	packageImports += "import org.jdbi.v3.stringtemplate4.UseStringTemplateSqlLocator;\n"
	packageImports += "import " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".domain.AbstractPojoAuditVersion;\n\n"
	return
}

func MakeRepositoryFile(table pojo.Table, className string) {
	os.Mkdir(constants.REPOSITORY_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.REPOSITORY_FOLDER+className+"Repository.java",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(getRepositoryClass(table, className))
	arquivo.Close()
}

func getRepositoryClass(table pojo.Table, class string) string {
	repositoryClass := "package " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".repository;\n\n"
	repositoryClass += "import org.jdbi.v3.sqlobject.config.RegisterFieldMapper;\n"
	repositoryClass += "import org.jdbi.v3.sqlobject.statement.SqlQuery;\n"
	repositoryClass += "import org.jdbi.v3.stringtemplate4.UseStringTemplateSqlLocator;\n"
	repositoryClass += "import org.springframework.stereotype.Repository;\n\n"
	repositoryClass += "import " + constants.DEFAULT_PACKAGE + "." + constants.DB_SCHEMA + ".domain." + class + ";\n\n\n"
	repositoryClass += "@Repository\n"
	repositoryClass += "@RegisterFieldMapper(" + class + ".class)\n"

	extends := "AbstractRepository"
	idClass := "Long"
	if table.HasMultiplePK {
		extends = "AbstractRepositoryWithComplexId"
		idClass = class + "Id"
	}

	repositoryClass += "public interface " + class + "Repository extends " + extends + "<" + class + ", " + idClass + "> {\n\n"
	repositoryClass += "}\n"
	return repositoryClass
}
