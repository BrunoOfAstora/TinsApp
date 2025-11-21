package orders

import (
	"database/sql"
	"fmt"

	"github.com/BrunoOfAstora/internal"
	"github.com/BrunoOfAstora/internal/db/clients"
	"github.com/BrunoOfAstora/internal/db/items"

	_ "github.com/mattn/go-sqlite3"
)

type Orders struct {
	Id       int
	ItemId   int
	ClientId int

	IName  string
	IPrice float64
	IQtd   int
	CName  string
}

func OrdersDbInit(dbfilepath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbfilepath)
	internal.FatalErrChecking(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		itemid INTEGER,
		clientsid INTEGER NOT NULL,

		iname TEXT,
		iprice FLOAT,
		iqtd INTEGER,
		cname TEXT,

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

func OrdersDbGetInfoName(db *sql.DB, cName string) ([]Orders, error) {
	rows, err := db.Query(`
		SELECT
			orders.id,
			orders.itemid,
			orders.clientsid,
			items.name, 
			items.price, 
			COUNT(orders.itemid) as qtd,
			clients.name
		FROM orders 
		INNER JOIN items ON orders.itemid = items.id
		INNER JOIN clients ON orders.clientsid = clients.id
		WHERE clients.name = ?
		GROUP BY items.id;
	`, cName)

	internal.FatalErrChecking(err)

	defer rows.Close()

	var orderDetails []Orders

	for rows.Next() {
		var details Orders
		if err := rows.Scan(
			&details.Id,
			&details.ItemId,
			&details.ClientId,
			&details.IName,
			&details.IPrice,
			&details.IQtd,
			&details.CName); err != nil { // ← ADICIONADO
			return orderDetails, err
		}
		orderDetails = append(orderDetails, details)
	}

	if err = rows.Err(); err != nil {
		return orderDetails, err
	}
	return orderDetails, err
}

func OrdersDbRemove(db *sql.DB, clientId int, itemId int) error {
	res, err := db.Exec(`
		 DELETE FROM orders WHERE id = (
            SELECT id FROM orders 
            WHERE clientsid = ? AND itemid = ? 
            LIMIT 1
        );
	`, clientId, itemId)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum pedido encontrado para remover (Client: %d, Item: %d)", clientId, itemId)
	}

	return nil
}

/*func OrdersGetItemsByCategory(db *sql.DB) (map[Orders][]struct{
	Name string
	Tpe string
	Price float64
}, error) {

res, _ := db.Query(`
	SELECT
		itemid,

	FROM
`)

}
*/
