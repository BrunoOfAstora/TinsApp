package main

import (
	//"time"

	//"github.com/BrunoOfAstora/internal/db/clients"
	//"github.com/BrunoOfAstora/internal/db/clients"
	"github.com/BrunoOfAstora/internal/db/generic"
	"github.com/BrunoOfAstora/internal/db/items"

	//"github.com/BrunoOfAstora/internal/db/orders"
	"github.com/BrunoOfAstora/ui"
)

func main() {
	dbpath := generic.DbFilePath()
	items.ItemsDbInit(dbpath)

	//clients.ClientDbInit(dbpath)
	//orders.EditOrder(order_db, "Joseph", "b")
	ui.StartUI()
}
