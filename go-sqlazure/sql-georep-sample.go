package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func syncSecondary(db *sql.DB, targetServer, targetDatabase string) error {

	query := fmt.Sprintf(`
		EXEC sys.sp_wait_for_database_copy_sync @target_server = N'%s', @target_database = N'%s';
		`,
		targetServer,
		targetDatabase,
	)

	fmt.Println(query)

	result, err := db.Exec(
		query,
	)

	if err != nil {
		log.Fatalf("syncSecondary: result %+v, error %+v", result, err)
	}
	return nil
}

func main() {

	primaryServer := "jiren-sql.database.windows.net"
	targetServer := "jiren-sql-secondary"
	//primaryDatabase := "jiren-sql"
	targetDatabase := "jiren-sql"

	connStr := fmt.Sprintf("server=%s;initial catalog=%s;database=%s;user id=jiren;password=Understand@Remote;port=1433;", primaryServer, targetDatabase, targetDatabase)
	fmt.Println(connStr)

	db, err := sql.Open("mssql", connStr)

	if err != nil {
		log.Fatal(db, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(tx, err)
	}
	// Create table users if it is not there yet
	result, err := db.Exec(
		`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'users' AND xtype = 'U')
		      CREATE TABLE users (
				  name varchar(64) not null, age int
			  )`,
	)

	if err != nil {
		log.Fatalf("Position 1 %+v, %+v", result, err)
	}

	// Add an entry in the table
	result, err = db.Exec(
		`
		INSERT INTO users (name, age) VALUES ($1, $2)`,
		"gopher",
		27,
	)

	if err != nil {
		log.Fatalf("Position 2 %+v, %+v", result, err)
	}

	// Remove duplicates
	result, err = db.Exec(
		`
		WITH CTE AS (
   			SELECT name, age,
       				RN = ROW_NUMBER() OVER (PARTITION BY name, age ORDER BY name, age)
   			FROM users
		) 
		DELETE FROM CTE WHERE RN > 1`,
	)

	if err != nil {
		log.Fatalf("Position 3 %+v, %+v", result, err)
	}

	// Wait for sync
	err = syncSecondary(db, targetServer, targetDatabase)

	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	age := 27
	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
