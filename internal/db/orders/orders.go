package orders

import (
	"database/sql"

	"github.com/BrunoOfAstora/internal"
	"github.com/BrunoOfAstora/internal/db/clients"
	"github.com/BrunoOfAstora/internal/db/items"

	_ "github.com/mattn/go-sqlite3"
)

type Orders struct {
	Id       int
	ItemId   int
	ClientId int
}

func OrdersDbInit(dbfilepath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbfilepath)
	internal.FatalErrChecking(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		itemid INTEGER,
		clientsid INTEGER NOT NULL,
		FOREIGN KEY (itemid) REFERENCES items(id),
		FOREIGN KEY (clientsid) REFERENCES clients(id)
		);
	`)

	internal.FatalErrChecking(err)
	return db
}

func OrdersDbInsert(db *sql.DB, orders *Orders) int {

	res, err := db.Exec(`
		INSERT INTO orders (itemid, clientsid) VALUES (?,?);
	`, orders.ItemId, orders.ClientId)

	internal.FatalErrChecking(err)

	id, err := res.LastInsertId()
	internal.FatalErrChecking(err)

	orders.Id = int(id)
	return orders.Id
}

func OrdersDbRemoveAll(db *sql.DB, orders *Orders) error {

	_, err := db.Exec(`
		TRUNCATE TABLE orders;
	`)

	if err != nil {
		return err
	}

	return nil
}

func OrdersDbRemove(db *sql.DB, orders *Orders) error {

	_, err := db.Exec(`
		DELETE FROM orders WHERE clientsid = ?;
	`, orders.Id)

	if err != nil {
		return err
	}

	return nil
}

func InsertNewOrder(db *sql.DB, clientId int, itemId int) int {
	orders := Orders{}

	orders.ClientId = clientId
	orders.ItemId = itemId

	return OrdersDbInsert(db, &orders)
}

func CreateNewOrder(db *sql.DB, cname string) int {
	clientId := clients.InsertNewClient(db, cname)

	InsertNewOrder(db, clientId, 0)

	return 0
}

func EditOrder(db *sql.DB, clientName string, itemName string) {

	itemId, _ := items.ItemsDbGet(db, itemName)
	clientId, _ := clients.ClientDbGet(db, clientName)

	InsertNewOrder(db, clientId, itemId)

}
