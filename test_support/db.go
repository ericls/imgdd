package test_support

import "github.com/ericls/imgdd/db"

func ResetDatabase(dbConf *db.DBConfigDef) {
	conn := db.GetConnection(dbConf)
	conn.Exec(`DO $$ DECLARE
  r RECORD;
BEGIN
  FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname =current_schema()) LOOP
    EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
  END LOOP;
END $$;`)
	db.PopulateBuiltInRoles(dbConf)
}
