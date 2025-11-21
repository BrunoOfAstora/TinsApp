package items

import (
	"database/sql"

	"github.com/BrunoOfAstora/internal"

	_ "github.com/mattn/go-sqlite3"
)

type Items struct {
	Id    int
	Name  string
	Type  string
	Price float64
}

func ItemsDbInit(dbfilepath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbfilepath)
	internal.FatalErrChecking(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		price FLOAT NOT NULL
		);
	`)

	internal.FatalErrChecking(err)
	return db
}

func ItemsDbInsert(db *sql.DB, items *Items) int {

	res, err := db.Exec(`
		INSERT INTO items (name, type, price) VALUES (?,?,?);
	`, items.Name, items.Type, items.Price)

	internal.FatalErrChecking(err)

	id, err := res.LastInsertId()
	internal.FatalErrChecking(err)

	items.Id = int(id)
	return items.Id
}

func ItemsDbGet(db *sql.DB, item string) (int, error) {
	i := Items{}
	err := db.QueryRow(`
		SELECT id FROM items WHERE name = ?;
	`, item).Scan(&i.Id)

	if err != nil {
		return 0, err
	}

	return i.Id, nil
}

func InsertNewItem(db *sql.DB, name string, tpe string, price float64) int {
	items := Items{}

	items.Name = name
	items.Type = tpe
	items.Price = price

	return ItemsDbInsert(db, &items)
}
