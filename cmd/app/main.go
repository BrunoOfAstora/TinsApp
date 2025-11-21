package main

import (
	//"time"

	//"github.com/BrunoOfAstora/internal/db/clients"
	"github.com/BrunoOfAstora/internal/db/clients"
	"github.com/BrunoOfAstora/internal/db/generic"
	"github.com/BrunoOfAstora/internal/db/items"
	"github.com/BrunoOfAstora/internal/db/orders"

	//"github.com/BrunoOfAstora/internal/db/orders"
	"github.com/BrunoOfAstora/ui"
)

func main() {
	dbpath := generic.DbFilePath()
	/*	i := */ items.ItemsDbInit(dbpath)
	/*	o := */ orders.OrdersDbInit(dbpath)
	clients.ClientDbInit(dbpath)

	//items.InsertNewItem(i, "b", "burguer", 60)

	//id, _ := items.ItemsDbGet(i, "b")
	//cl, _ := clients.ClientDbGet(o, "José")

	//orders.InsertNewOrder(o, cl, id)

	//clients.ClientDbInit(dbpath)
	//orders.EditOrder(order_db, "Joseph", "b")
	ui.StartUI()
}
