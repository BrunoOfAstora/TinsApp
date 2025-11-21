package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/BrunoOfAstora/internal/db/clients"
	"github.com/BrunoOfAstora/internal/db/generic"
	"github.com/BrunoOfAstora/internal/db/orders"
)

func StartUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TinsApp")

	dbPath := generic.DbFilePath()
	dbInit := clients.ClientDbInit(dbPath)

	var cli []string

	cli, _ = clients.ClientDbGetName(dbInit)

	// Pedidos simulados iniciais
	clientOrders := cli

	orderList := container.NewVBox()

	// Variáveis de estado da UI
	var selectedOrder string
	var selectedClientName string
	var mainScreen *fyne.Container
	var orderScreen *fyne.Container
	var orderDetailScreen *fyne.Container

	// CHANGE: Variáveis para manipular o conteúdo da tela de detalhes
	detailListContainer := container.NewVBox()
	totalLabel := widget.NewLabelWithStyle("Total: R$ 0,00", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	// CHANGE: Função auxiliar para criar linhas de produtos usando dados do banco de dados
	var populateOrderDetails func(string)
	populateOrderDetails = func(clientName string) {
		detailListContainer.Objects = nil

		products, _ := orders.OrdersDbGetInfoName(dbInit, clientName)

		if len(products) == 0 {
			detailListContainer.Add(widget.NewLabel("Nenhum pedido encontrado para este cliente"))
			detailListContainer.Refresh()
			totalLabel.SetText("Total do Pedido: R$ 0,00")
			return
		}

		var total float64 = 0

		header := container.NewHBox(
			widget.NewLabelWithStyle("Qtd", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Produto", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewLabelWithStyle("Valor", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)
		detailListContainer.Add(header)
		detailListContainer.Add(widget.NewSeparator())

		for _, p := range products {
			itemTotal := p.IPrice * float64(p.IQtd)
			total += itemTotal

			itemId := p.ItemId
			clientId := p.ClientId
			clientNameCaptured := p.CName

			btnRemove := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				err := orders.OrdersDbRemove(dbInit, clientId, itemId)
				if err != nil {
					fmt.Println("Erro ao remover item:", err)
					return
				}

				populateOrderDetails(clientNameCaptured)
			})
			btnRemove.Importance = widget.DangerImportance

			row := container.NewHBox(
				widget.NewLabel(fmt.Sprintf("%dx", p.IQtd)),
				widget.NewLabel(p.IName),
				layout.NewSpacer(),
				widget.NewLabel(fmt.Sprintf("R$ %.2f", itemTotal)),
				btnRemove,
			)
			detailListContainer.Add(row)
			detailListContainer.Add(widget.NewSeparator())
		}

		totalLabel.SetText(fmt.Sprintf("Total do Pedido: R$ %.2f", total))
		detailListContainer.Refresh()
	}

	reloadOrders := func() {
		orderList.Objects = nil
		for _, c := range clientOrders {
			clientName := c
			btn := widget.NewButton(clientName, func() {
				selectedClientName = clientName
				selectedOrder = clientName

				populateOrderDetails(selectedClientName)

				mainScreen.Hide()
				orderDetailScreen.Show()
			})
			btn.Resize(fyne.NewSize(360, 40))
			orderList.Add(btn)
		}
		orderList.Refresh()
	}

	// --- TELA PRINCIPAL ---
	title := widget.NewLabelWithStyle("TinsApp", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	btnOpenCash := widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {})
	btnAddOrder := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	btnAddOrderContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(70, 70)), btnAddOrder)

	topBar := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	topBar.SetMinSize(fyne.NewSize(1600, 60))
	bottomBar := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	bottomBar.SetMinSize(fyne.NewSize(1600, 100))

	topContainer := container.NewBorder(nil, nil, nil, btnOpenCash, topBar)
	rightAlignedAdd := container.NewHBox(layout.NewSpacer(), btnAddOrderContainer)
	bottomContainer := container.NewBorder(nil, nil, nil, rightAlignedAdd, bottomBar)

	scrollContainer := container.NewVScroll(orderList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 700))
	contentWithScroll := container.NewVBox(title, scrollContainer)
	leftAlignedContent := container.NewWithoutLayout(contentWithScroll)
	contentWithScroll.Move(fyne.NewPos(20, 20))
	contentWithScroll.Resize(fyne.NewSize(420, 740))

	mainScreen = container.NewBorder(topContainer, bottomContainer, nil, nil, leftAlignedContent)

	// --- TELA NOVO PEDIDO ---
	/*orderTitle :=*/
	widget.NewLabelWithStyle("Novo Pedido", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	btnBack := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), nil)
	topOrderBar := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	topOrderBar.SetMinSize(fyne.NewSize(1600, 60))
	topOrderContainer := container.NewBorder(nil, nil, btnBack, nil, topOrderBar)

	// NOVO: Conteúdo da tela de novo pedido
	// Barra de entrada para nome do cliente
	clientNameEntry := widget.NewEntry()
	clientNameEntry.SetPlaceHolder("Digite o nome do cliente...")

	// Dropdown de categorias
	categoryDropdown := widget.NewSelect([]string{
		"Hamburgueres",
		"Smash",
		"Refrigerantes",
		"Acompanhamentos",
		"Sobremesas",
	}, func(s string) {
		fmt.Println("Categoria selecionada:", s)
	})
	categoryDropdown.SetSelected("Hamburgueres") // Padrão

	// Container para nome e dropdown
	headerNewOrder := container.NewHBox(
		widget.NewLabel("Cliente:"),
		clientNameEntry,
		layout.NewSpacer(),
		widget.NewLabel("Categoria:"),
		categoryDropdown,
	)

	// Container para exibir itens disponíveis (será preenchido dinamicamente)
	itemsContainer := container.NewVBox()

	// Função para atualizar itens conforme a categoria
	updateItemsDisplay := func(category string) {
		itemsContainer.Objects = nil

		// DADOS DE EXEMPLO - Você deve buscar do banco de dados
		var items map[string][]struct {
			Name  string
			Price float64
		}

		items = map[string][]struct {
			Name  string
			Price float64
		}{
			"Hamburgueres": {
				{Name: "X-Burguer", Price: 25.00},
				{Name: "X-Tudo", Price: 35.50},
				{Name: "X-Frango", Price: 28.00},
			},
			"Smash": {
				{Name: "Smash Tradicional", Price: 22.00},
				{Name: "Smash Queijo Triplo", Price: 32.00},
				{Name: "Smash Bacon", Price: 28.00},
			},
			"Refrigerantes": {
				{Name: "Refrigerante Lata", Price: 6.00},
				{Name: "Refrigerante 2L", Price: 10.00},
				{Name: "Suco Natural", Price: 8.00},
			},
			"Porções": {
				{Name: "Batata Frita P", Price: 12.00},
				{Name: "Batata Frita G", Price: 18.00},
				{Name: "Cebola Frita", Price: 15.00},
			},
			"Sobremesas": {
				{Name: "Sorvete", Price: 8.00},
				{Name: "Brownie", Price: 12.00},
				{Name: "Pudim", Price: 10.00},
			},
		}

		// Cabeçalho da tabela de itens
		headerItems := container.NewHBox(
			widget.NewLabelWithStyle("Produto", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewLabelWithStyle("Preço", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Ação", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)
		itemsContainer.Add(headerItems)
		itemsContainer.Add(widget.NewSeparator())

		// Exibir itens da categoria selecionada
		if categoryItems, exists := items[category]; exists {
			for _, item := range categoryItems {
				itemName := item.Name
				itemPrice := item.Price
				clientName := clientNameEntry.Text

				// Botão para adicionar item ao pedido
				btnAddItem := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
					if clientName == "" {
						fmt.Println("Erro: Digite o nome do cliente")
						return
					}

					// TODO: Implementar no backend
					// Buscar clientId pelo nome
					// Buscar itemId pelo nome
					// Chamar: orders.InsertNewOrder(dbInit, clientId, itemId)

					fmt.Printf("Adicionando %s para cliente %s\n", itemName, clientName)
				})
				btnAddItem.Importance = widget.SuccessImportance

				// Linha do item
				row := container.NewHBox(
					widget.NewLabel(itemName),
					layout.NewSpacer(),
					widget.NewLabel(fmt.Sprintf("R$ %.2f", itemPrice)),
					btnAddItem,
				)
				itemsContainer.Add(row)
				itemsContainer.Add(widget.NewSeparator())
			}
		}

		itemsContainer.Refresh()
	}

	// Atualizar itens quando mudar a categoria
	categoryDropdown.OnChanged = func(s string) {
		updateItemsDisplay(s)
	}

	// Inicializar com a primeira categoria
	updateItemsDisplay("Hamburgueres")

	// Scroll para itens
	itemsScroll := container.NewVScroll(itemsContainer)
	itemsScroll.SetMinSize(fyne.NewSize(1400, 650))

	// Conteúdo completo da tela de novo pedido
	orderScreenContent := container.NewVBox(
		headerNewOrder,
		widget.NewSeparator(),
		itemsScroll,
	)

	orderScreen = container.NewBorder(
		topOrderContainer,
		nil,
		nil,
		nil,
		orderScreenContent,
	)

	// --- TELA DETALHES DO PEDIDO ---

	orderDetailTitle := widget.NewLabelWithStyle("Detalhes do Pedido", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	btnBackDetail := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), nil)

	topDetailBar := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	topDetailBar.SetMinSize(fyne.NewSize(1600, 60))
	topDetailContainer := container.NewBorder(nil, nil, btnBackDetail, nil, topDetailBar)

	headerDetail := container.NewVBox(orderDetailTitle, widget.NewSeparator())

	detailScroll := container.NewVScroll(detailListContainer)

	btnEditOrder := widget.NewButtonWithIcon("Editar Pedido", theme.DocumentCreateIcon(), func() {
		fmt.Println("Editar pedido clicado: ", selectedOrder)
	})
	btnEditOrder.Importance = widget.HighImportance

	bottomDetailInfo := container.NewHBox(
		layout.NewSpacer(),
		totalLabel,
		layout.NewSpacer(),
	)

	bottomDetailContainer := container.NewVBox(
		widget.NewSeparator(),
		bottomDetailInfo,
		btnEditOrder,
	)
	paddedBottomDetail := container.NewPadded(bottomDetailContainer)

	orderDetailScreen = container.NewBorder(
		container.NewVBox(topDetailContainer, headerDetail),
		paddedBottomDetail,
		nil, nil,
		detailScroll,
	)

	// Inicializa lógica
	reloadOrders()

	// Stack Navigation
	stack := container.NewStack(mainScreen, orderScreen, orderDetailScreen)
	orderScreen.Hide()
	orderDetailScreen.Hide()

	// Eventos
	btnAddOrder.OnTapped = func() {
		clientNameEntry.SetText("") // Limpa o campo de nome
		categoryDropdown.SetSelected("Hamburgueres")
		updateItemsDisplay("Hamburgueres")
		mainScreen.Hide()
		orderScreen.Show()
	}

	btnBack.OnTapped = func() {
		orderScreen.Hide()
		mainScreen.Show()
	}

	btnBackDetail.OnTapped = func() {
		orderDetailScreen.Hide()
		mainScreen.Show()
	}

	myWindow.SetContent(stack)
	myWindow.Resize(fyne.NewSize(1600, 900))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
