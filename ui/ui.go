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
)

func StartUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TinsApp")

	dbPath := generic.DbFilePath()
	dbInit := clients.ClientDbInit(dbPath)

	var cli []string

	cli, _ = clients.ClientDbGetName(dbInit)

	// Pedidos simulados iniciais
	orders := cli

	orderList := container.NewVBox()

	// Variáveis de estado da UI
	var selectedOrder string
	var mainScreen *fyne.Container
	var orderScreen *fyne.Container
	var orderDetailScreen *fyne.Container

	// CHANGE: Variáveis para manipular o conteúdo da tela de detalhes
	detailListContainer := container.NewVBox() // Onde ficarão os produtos
	totalLabel := widget.NewLabelWithStyle("Total: R$ 0,00", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	// CHANGE: Função auxiliar para criar linhas de produtos (simulando banco de dados)
	populateOrderDetails := func(orderName string) {
		detailListContainer.Objects = nil // Limpa a lista anterior

		// Simula produtos diferentes
		mockProducts := []struct {
			Name  string
			Qtd   int
			Price float64
		}{
			{"X-Tudo Artesanal", 1, 35.50},
			{"Refrigerante Lata", 2, 6.00},
			{"Batata Frita Grande", 1, 18.00},
			{"Adicional de Bacon", 1, 4.50},
		}

		var total float64 = 0

		// Adiciona cabeçalho da tabela
		header := container.NewHBox(
			widget.NewLabelWithStyle("Qtd", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Produto", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewLabelWithStyle("Valor", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		)
		detailListContainer.Add(header)
		detailListContainer.Add(widget.NewSeparator())

		// Itera e cria as linhas
		for _, p := range mockProducts {
			itemTotal := p.Price * float64(p.Qtd)
			total += itemTotal

			// Layout da linha: Qtd | Nome ... | Preço
			row := container.NewHBox(
				widget.NewLabel(fmt.Sprintf("%dx", p.Qtd)),
				widget.NewLabel(p.Name),
				layout.NewSpacer(),
				widget.NewLabel(fmt.Sprintf("R$ %.2f", itemTotal)),
			)
			detailListContainer.Add(row)
			detailListContainer.Add(widget.NewSeparator()) // Linha divisória
		}

		// Atualiza o total lá embaixo
		totalLabel.SetText(fmt.Sprintf("Total do Pedido: R$ %.2f", total))
		detailListContainer.Refresh()
	}

	reloadOrders := func() {
		orderList.Objects = nil
		for _, c := range orders {
			orderName := c
			btn := widget.NewButton(orderName, func() {
				selectedOrder = orderName

				// CHANGE: Ao clicar, populamos a lista e calculamos o total
				populateOrderDetails(selectedOrder)

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
	orderTitle := widget.NewLabelWithStyle("Novo Pedido", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	btnBack := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), nil)
	topOrderBar := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	topOrderBar.SetMinSize(fyne.NewSize(1600, 60))
	topOrderContainer := container.NewBorder(nil, nil, btnBack, nil, topOrderBar)
	orderScreen = container.NewBorder(topOrderContainer, nil, nil, nil, container.NewCenter(orderTitle))

	// --- TELA DETALHES DO PEDIDO (MUDANÇA PRINCIPAL AQUI) ---

	orderDetailTitle := widget.NewLabelWithStyle("Detalhes do Pedido", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	btnBackDetail := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), nil)

	// Barra Topo Detalhes
	topDetailBar := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	topDetailBar.SetMinSize(fyne.NewSize(1600, 60))
	topDetailContainer := container.NewBorder(nil, nil, btnBackDetail, nil, topDetailBar)

	// Cabeçalho com o título dentro do container de conteúdo
	headerDetail := container.NewVBox(orderDetailTitle, widget.NewSeparator())

	// Scroll View para a lista de produtos
	detailScroll := container.NewVScroll(detailListContainer)

	// CHANGE: Construção da barra inferior da tela de detalhes (Total + Botão Editar)
	btnEditOrder := widget.NewButtonWithIcon("Editar Pedido", theme.DocumentCreateIcon(), func() {
		fmt.Println("Editar pedido clicado: ", selectedOrder)
	})
	btnEditOrder.Importance = widget.HighImportance // Deixa o botão azul/destacado

	// Container inferior: Total na esquerda (ou direita), Botão Editar na direita
	bottomDetailInfo := container.NewHBox(
		layout.NewSpacer(),
		totalLabel, // Label do total
		layout.NewSpacer(),
	)

	// Container que segura o total e o botão de editar
	bottomDetailContainer := container.NewVBox(
		widget.NewSeparator(),
		bottomDetailInfo,
		btnEditOrder,
	)
	// Padding para não ficar colado na borda
	paddedBottomDetail := container.NewPadded(bottomDetailContainer)

	// Montagem da tela de detalhes usando Border Layout
	// Top: Barra de voltar
	// Bottom: Total e Botão Editar
	// Center: Lista de produtos (com scroll)
	orderDetailScreen = container.NewBorder(
		container.NewVBox(topDetailContainer, headerDetail), // Topo
		paddedBottomDetail, // Fundo
		nil, nil,
		detailScroll, // Centro (Conteúdo)
	)

	// Inicializa lógica
	reloadOrders()

	// Stack Navigation
	stack := container.NewStack(mainScreen, orderScreen, orderDetailScreen)
	orderScreen.Hide()
	orderDetailScreen.Hide()

	// Eventos
	btnAddOrder.OnTapped = func() {
		orders = append(orders, "Cliente Novo")
		reloadOrders()
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
