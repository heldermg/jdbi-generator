package trigger

import (
	"fmt"
	"os"

	"github.com/heldermg/jdbi-generator/constants"
	"github.com/heldermg/jdbi-generator/pojo"
)

func MakeGenericConcurrentTrigger() {
	if _, err := os.Stat(constants.GENERATED_FOLDER); os.IsNotExist(err) {
		os.Mkdir(constants.GENERATED_FOLDER, os.ModePerm)
	}

	os.Mkdir(constants.SQL_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.SQL_FOLDER+constants.GENREIC_CONCURRENT_FN_NAME+".sql",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	trigger := "create or replace function " + constants.DB_SCHEMA + "." + constants.GENREIC_CONCURRENT_FN_NAME + "() returns trigger \n"
	trigger += "language plpgsql as $function$ \n"
	trigger += "declare\n"
	trigger += "   _version_ integer;\n"
	trigger += "   _sql_ text;\n"
	trigger += "   _where_ text;\n"
	trigger += "   _table_pks_ record;\n"
	trigger += "   _pk_count_ int;\n"
	trigger += "   _table_name_ text := '" + constants.DB_SCHEMA + ".' || TG_TABLE_NAME;\n"
	trigger += "begin\n"
	trigger += "   _pk_count_ := 1;\n"
	trigger += "   _where_ := '';\n"
	trigger += "   for _table_pks_ in (\n"
	trigger += "      select a.attname as column_name from pg_index i\n"
	trigger += "      join pg_attribute a on a.attrelid = i.indrelid\n"
	trigger += "      and a.attnum = ANY(i.indkey)\n"
	trigger += "      where i.indrelid = _table_name_::regclass\n"
	trigger += "      and i.indisprimary)\n"
	trigger += "   loop\n"
	trigger += "      if (_pk_count_ > 1) then\n"
	trigger += "          _where_ := _where_ || ' and ';\n"
	trigger += "      end if;\n"
	trigger += "      _where_ := _where_ || _table_pks_.column_name || ' = ($1).' || _table_pks_.column_name;\n"
	trigger += "      _pk_count_ := _pk_count_ + 1;\n"
	trigger += "   end loop;\n\n"
	trigger += "   _sql_ := 'select ' || " + constants.VERSION_COLUMN_NAME + " || ' from ' || _table_name_ || ' where ' || _where_;\n"
	trigger += "   execute _sql_ into _version_ using old;\n\n"
	trigger += "   if (_version_ is null) then\n"
	trigger += "      _version_ := 0;\n"
	trigger += "   end if;\n\n"
	trigger += "   if (TG_OP = 'UPDATE') then\n"
	trigger += "      if (_version_ + 1) <> NEW." + constants.VERSION_COLUMN_NAME + " then\n"
	trigger += "         raise exception 'Invalid Version. Table %: version must be %, but got %', _table_name_, (_version_ + 1), NEW." + constants.VERSION_COLUMN_NAME + " using ERRCODE = '23501';\n"
	trigger += "      end if;\n"
	trigger += "      return new;\n"
	trigger += "   elsif (TG_OP = 'DELETE') then\n"
	trigger += "      if (_version_) <> OLD." + constants.VERSION_COLUMN_NAME + " then\n"
	trigger += "         raise exception 'Invalid Version. Table %: version must be %, but got %', _table_name_, _version_, OLD." + constants.VERSION_COLUMN_NAME + " using ERRCODE = '23501';\n"
	trigger += "      end if;\n"
	trigger += "      return old;\n"
	trigger += "   end if;\n"
	trigger += "   return null;\n"
	trigger += "end;\n"
	trigger += "$function$\n"
	trigger += ";"

	arquivo.WriteString(trigger)
	arquivo.Close()
}

func MakeTableConcurrentTrigger(table pojo.Table) {
	os.Mkdir(constants.SQL_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.SQL_FOLDER+constants.CONCURRENT_TRIGGER_PREFIX+"tables.sql",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	trigger := "-- TABLE: " + table.Name + "\n"
	trigger += "-- drop trigger if exists " + constants.CONCURRENT_TRIGGER_PREFIX + table.Name + " on " + constants.DB_SCHEMA + "." + table.Name + ";\n"
	trigger += "create trigger " + constants.CONCURRENT_TRIGGER_PREFIX + table.Name + "\n"
	trigger += "before delete or update on " + constants.DB_SCHEMA + "." + table.Name + "\n"
	trigger += "for each row execute procedure " + constants.DB_SCHEMA + "." + constants.GENREIC_CONCURRENT_FN_NAME + "();\n\n"

	arquivo.WriteString(trigger)
	arquivo.Close()
}

func MakeGenericAuditTrigger() {
	if _, err := os.Stat(constants.GENERATED_FOLDER); os.IsNotExist(err) {
		os.Mkdir(constants.GENERATED_FOLDER, os.ModePerm)
	}

	os.Mkdir(constants.SQL_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.SQL_FOLDER+constants.GENERIC_AUDIT_FN_NAME+".sql",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	trigger := "create or replace function " + constants.DB_SCHEMA + "." + constants.GENERIC_AUDIT_FN_NAME + "() returns trigger \n"
	trigger += "language plpgsql as $function$ \n"
	trigger += "declare\n"
	trigger += "   _sql_ text;\n"
	trigger += "   _audit_table_name_ text := TG_TABLE_NAME || '_audit';\n"
	trigger += "   _columns_ text := 'id';\n"
	trigger += "   _original_table_name_ text := '" + constants.DB_SCHEMA + ".' || TG_TABLE_NAME;\n"
	trigger += "   _table_columns_ record;\n"
	trigger += "begin\n"
	trigger += "   for _table_columns_ in (\n"
	trigger += "      select a.attname\n"
	trigger += "      from pg_catalog.pg_attribute a\n"
	trigger += "      where attrelid = _original_table_name_::regclass\n"
	trigger += "         and a.attnum > 0\n"
	trigger += "         and not a.attisdropped)\n"
	trigger += "   loop\n"
	trigger += "      _columns_ := _columns_ || ', ';\n"
	trigger += "      if (_table_columns_.attname = 'id') then\n"
	trigger += "         _columns_ := _columns_ || TG_TABLE_NAME || '_id';\n"
	trigger += "      else\n"
	trigger += "         _columns_ := _columns_ || _table_columns_.attname;\n"
	trigger += "      end if;\n"
	trigger += "   end loop;\n\n"
	trigger += "   _sql_ := FORMAT('INSERT INTO ' || TG_TABLE_SCHEMA || '_audit.%1$I\n"
	trigger += "      (' || _columns_ || ') VALUES\n"
	trigger += "      (NEXTVAL(pg_get_serial_sequence(''' || TG_TABLE_SCHEMA || '_audit.%1$I'', ''id'')),\n"
	trigger += "      ($1).*)', _audit_table_name_);\n\n"
	trigger += "   if (TG_OP = 'INSERT') or (TG_OP = 'UPDATE') then\n"
	trigger += "      execute _sql_ USING NEW;\n"
	trigger += "      return NEW;\n\n"
	trigger += "   elsif (TG_OP = 'DELETE') THEN\n"
	trigger += "      OLD.audit_operacao := 'DELETE';\n"
	trigger += "      execute _sql_ USING OLD;\n"
	trigger += "      return OLD;\n"
	trigger += "   end if;\n\n"
	trigger += "   return null;\n"
	trigger += "end;\n"
	trigger += "$function$\n"
	trigger += ";"

	arquivo.WriteString(trigger)
	arquivo.Close()
}

func MakeTableAuditTrigger(table pojo.Table) {
	os.Mkdir(constants.SQL_FOLDER, os.ModePerm)
	arquivo, err := os.OpenFile(
		constants.SQL_FOLDER+constants.AUDIT_TRIGGER_PREFIX+"tables.sql",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}

	trigger := "-- TABLE: " + table.Name + "\n"
	trigger += "-- drop trigger if exists " + constants.AUDIT_TRIGGER_PREFIX + table.Name + " on " + constants.DB_SCHEMA + "." + table.Name + ";\n"
	trigger += "create trigger " + constants.AUDIT_TRIGGER_PREFIX + table.Name + "\n"
	trigger += "after insert or delete or update on " + constants.DB_SCHEMA + "." + table.Name + "\n"
	trigger += "for each row execute procedure " + constants.DB_SCHEMA + "." + constants.GENERIC_AUDIT_FN_NAME + "();\n\n"

	arquivo.WriteString(trigger)
	arquivo.Close()
}
