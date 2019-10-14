package repository

import (
	"fmt"
	"os"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/pojo"
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
	pojoClass := "\npublic interface AbstractRepository<T, I> {\n\n"
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
	save += "    @SqlUpdate\n"
	save += "    @GetGeneratedKeys\n"
	save += "    @UseStringTemplateSqlLocator\n"
	save += "    T save(@BindBean(\"entity\") T entity);\n\n"
	return
}

func getUpdate() (update string) {
	update += "    @SqlUpdate\n"
	update += "    @GetGeneratedKeys\n"
	update += "    @UseStringTemplateSqlLocator\n"
	update += "    T update(@BindBean(\"entity\") T entity);\n\n"
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
	packageImports = "package br.gov.economia.maisbrasil." + constants.DB_SCHEMA + ".repository;\n\n"
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
	repositoryClass := "package br.gov.economia.maisbrasil." + constants.DB_SCHEMA + ".repository;\n\n"
	repositoryClass += "import org.jdbi.v3.sqlobject.config.RegisterFieldMapper;\n"
	repositoryClass += "import org.jdbi.v3.sqlobject.statement.SqlQuery;\n"
	repositoryClass += "import org.jdbi.v3.stringtemplate4.UseStringTemplateSqlLocator;\n"
	repositoryClass += "import org.springframework.stereotype.Repository;\n\n"
	repositoryClass += "import br.gov.economia.maisbrasil." + constants.DB_SCHEMA + ".domain." + class + ";\n\n\n"
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
