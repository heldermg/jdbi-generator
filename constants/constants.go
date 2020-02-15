package constants

// DB constants
const DB_IP_PORT = "localhost:5432"
const DB_NAME = "postgres"
const DB_SCHEMA = "myschema"
const DB_URL = "postgres://postgres:postgres@" +
	DB_IP_PORT + "/" + DB_NAME + "?sslmode=disable"

// Files
const GENERATED_FOLDER = "_generated/"
const SQL_STG_FOLDER = GENERATED_FOLDER + "sql.stg/"
const REPOSITORY_FOLDER = GENERATED_FOLDER + "repository/"
const DOMAIN_FOLDER = GENERATED_FOLDER + "domain/"
const SQL_FOLDER = GENERATED_FOLDER + "sql/"

// Pojo and sqlstg
const DEFAULT_PACKAGE = "br.com.myapp"
const VERSION_COLUMN_NAME = "version"
const VERSION_COLUMN_COMMENT = "Current record version identifier"
const AUDIT_PK_COLUMN_COMMENT = "Audit table unique identifier"
const AUDIT_DATETIME_COLUMN_NAME = "audit_datetime"
const AUDIT_DATETIME_COLUMN_COMMENT = "Operation date/time"
const AUDIT_OPERATION_COLUMN_NAME = "audit_operation"
const AUDIT_OPERATION_COLUMN_COMMENT = "Operation action (INSERT, UPDATE, DELETE)"
const AUDIT_USER_COLUMN_NAME = "audit_user"
const AUDIT_USER_COLUMN_COMMENT = "Operation user"

// Triggers
const GENREIC_CONCURRENT_FN_NAME = "fn_generic_concurrent"
const CONCURRENT_TRIGGER_PREFIX = "tg_concurrent_"
const GENERIC_AUDIT_FN_NAME = "fn_generic_audit"
const AUDIT_TRIGGER_PREFIX = "tg_audit_"
