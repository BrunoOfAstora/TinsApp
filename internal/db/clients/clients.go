package clients

import (
	"database/sql"
	//"time"

	"github.com/BrunoOfAstora/internal"

	_ "github.com/mattn/go-sqlite3"
)

type Clients struct {
	Id   int
	Name string
	//CreateAt time.Time
}

func ClientDbInit(dbfilepath string) *sql.DB {

	db, err := sql.Open("sqlite3", dbfilepath)
	internal.FatalErrChecking(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		tmstmp TEXT NOT NULL
		);
	`)

	internal.FatalErrChecking(err)
	return db
}

func ClientDbInsert(db *sql.DB, clients *Clients) int {

	res, err := db.Exec(`
		INSERT INTO clients (name, tmstmp) VALUES (?, time('now', 'localtime'));
	`, clients.Name)

	internal.FatalErrChecking(err)

	id, err := res.LastInsertId()
	internal.FatalErrChecking(err)

	clients.Id = int(id)
	return clients.Id
}

func ClientDbGet(db *sql.DB, clientName string) (int, error) {
	c := Clients{}

	err := db.QueryRow(`
		SELECT id FROM clients WHERE name = ?;
	`, clientName).Scan(&c.Id)

	if err != nil {
		return 0, err
	}

	return c.Id, nil
}

func ClientDbGetName(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
	SELECT name FROM clients;
	`)
	internal.FatalErrChecking(err)

	defer rows.Close()

	var getAllClients []string

	for rows.Next() {
		var clientsInfo Clients
		if err := rows.Scan(&clientsInfo.Name); err != nil {
			return getAllClients, err
		}
		getAllClients = append(getAllClients, clientsInfo.Name)

	}
	if err = rows.Err(); err != nil {
		return getAllClients, err
	}
	return getAllClients, err
}

func ClientDbRemoveAll(db *sql.DB, clients *Clients) error {
	_, err := db.Exec(`
		DELETE FROM clients;
	`)

	if err != nil {
		return err
	}

	return nil
}

func ClientDbRemove(db *sql.DB, clients *Clients) error {
	_, err := db.Exec(`
		DELETE FROM clients WHERE id = ?;
	`, clients.Id)

	if err != nil {
		return err
	}

	return nil
}

func InsertNewClient(db *sql.DB, name string) int {
	client := Clients{}

	client.Name = name

	return ClientDbInsert(db, &client)
}
