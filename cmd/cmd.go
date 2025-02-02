package cmd

import "flag"

var srcType = flag.String("src", "", "Source database type (sqlite or mysql)")
var sqliteFile = flag.String("sqlite", "", "SQLite database file path")
var mysqlDSN = flag.String("mysql", "", "MySQL DSN (e.g. 'root:password@tcp(127.0.0.1:3306)/your_db')")
var migrate = flag.Bool("migrate", false, "Perform database migration")

func Cmd() {
	flag.Parse()
	if *migrate {
		migrateDB()
		return
	}
	run()
}
