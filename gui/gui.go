package gui

import (
	"fmt"

	"github.com/cfanatic/go-expense/account"
	"github.com/cfanatic/go-expense/database"
	"github.com/cfanatic/go-expense/datasheet"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Gui struct {
	widgets.QWidget

	qApp    *widgets.QApplication
	hlayout *widgets.QHBoxLayout
	vlayout *widgets.QVBoxLayout
	tview   *widgets.QTableView
	bload   *widgets.QPushButton
	bsave   *widgets.QPushButton
	bprint  *widgets.QPushButton
	bquit   *widgets.QPushButton
	sitem   *gui.QStandardItemModel

	ds *datasheet.Datasheet
	db *database.Database

	_ func() `constructor:"init"`

	_ func() `slot:"connect"`
}

func (w *Gui) init() {
	geometry := widgets.QApplication_Desktop().AvailableGeometry(0)
	w.SetWindowTitle("go-expensegui")
	w.SetGeometry2(0, 0, 850, 600)
	w.Move2((geometry.Width()-w.Width())/2, (geometry.Height()-w.Height())/2)

	w.ConnectKeyPressEvent(w.keypressevent)

	w.hlayout = widgets.NewQHBoxLayout()
	w.vlayout = widgets.NewQVBoxLayout()
	w.tview = widgets.NewQTableView(nil)
	w.bload = widgets.NewQPushButton2("Load", nil)
	w.bsave = widgets.NewQPushButton2("Save", nil)
	w.bprint = widgets.NewQPushButton2("Print", nil)
	w.bquit = widgets.NewQPushButton2("Quit", nil)
	w.sitem = gui.NewQStandardItemModel(nil)

	w.bsave.SetEnabled(false)
	w.bprint.SetEnabled(false)

	w.tview.SetModel(w.sitem)
	w.hlayout.AddWidget(w.bload, 0, 0)
	w.hlayout.AddWidget(w.bsave, 0, 0)
	w.hlayout.AddWidget(w.bprint, 0, 0)
	w.hlayout.AddWidget(w.bquit, 0, 0)
	w.vlayout.AddWidget(w.tview, 0, 0)
	w.vlayout.AddLayout(w.hlayout, 0)

	w.SetLayout(w.vlayout)

	w.bload.ConnectClicked(w.load)
	w.bquit.ConnectClicked(func(bool) { w.qApp.Exit(0) })
	w.tview.HorizontalHeader().ConnectSectionResized(
		func(idx, old, new int) { fmt.Printf("Index: %d, Size: %d\n", idx, new) },
	)
}

func (w *Gui) load(bool) {
	var export []datasheet.Content
	var err error

	if f := widgets.QFileDialog_GetOpenFileNames(nil, "Open datasheet", core.QDir_HomePath(), "*.xlsx", "", 0); len(f) > 0 {
		w.ds = datasheet.New(f[0], account.TAB)
		w.db = database.New(account.ADDRESS, account.NAME, account.COLLECT)
	} else {
		return
	}

	if export, err = w.ds.Content(); err == nil || w.ds.Err == nil || w.db.Err == nil {
		w.bsave.SetEnabled(true)
		w.bprint.SetEnabled(true)
	} else {
		panic("Error during datasheet import!")

	}

	for _, trans := range export {
		items := []*gui.QStandardItem{}
		if item, err := w.db.Document("payee", trans.Payee); err == nil {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(item.Label),
			)
		} else {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(""),
			)
		}
		items[1].SetToolTip(trans.Description)
		w.sitem.AppendRow(items)
	}

	w.sitem.SetHorizontalHeaderLabels([]string{"Date", "Payee", "Amount", "Label"})
	w.tview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Interactive)
	w.tview.HorizontalHeader().ResizeSection(0, 100)
	w.tview.HorizontalHeader().ResizeSection(1, 400)
	w.tview.HorizontalHeader().ResizeSection(2, 100)
	w.tview.HorizontalHeader().ResizeSection(3, 200)
	w.tview.VerticalHeader().SetSectionResizeMode(widgets.QHeaderView__Fixed)
}

func (w *Gui) keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Escape) {
		w.qApp.Exit(0)
	}
}

func (w *Gui) InitWith(qApp *widgets.QApplication) {
	w.qApp = qApp
}
