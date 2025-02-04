package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatsUI struct {
}

func StatsExec() {
	app := tview.NewApplication()
	view := tview.NewFlex()

	menuFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	infoFlex := tview.NewFlex()
	infoFlex.Box.SetBorder(true).SetTitle("Info")

	navFlex := tview.NewFlex()
	navFlex.Box.SetBorder(true).SetTitle("Navigation")
	navList := tview.NewList().ShowSecondaryText(false).
		AddItem("Day", "", 0, nil).
		AddItem("Week", "", 0, nil).
		AddItem("Month", "", 0, nil).
		AddItem("Year", "", 0, nil).
		AddItem("Alltime", "", 0, nil)
	navFlex.AddItem(navList, 0, 1, true)
	menuFlex.AddItem(navFlex, 0, 1, true)

	datesFlex := tview.NewFlex()
	datesFlex.Box.SetBorder(true).SetTitle("Dates")
	datesList := tview.NewList().ShowSecondaryText(false).
		AddItem("01-01-01", "", 0, nil).
		AddItem("02-02-02", "", 0, nil)
	datesFlex.AddItem(datesList, 0, 1, true)
	menuFlex.AddItem(datesFlex, 0, 1, false)

	tagsFlex := tview.NewFlex()
	tagsFlex.Box.SetBorder(true).SetTitle("tasks")
	tagsList := tview.NewList().ShowSecondaryText(false).
		AddItem("Study", "", 0, nil).
		AddItem("Work", "", 0, nil)
	tagsFlex.AddItem(tagsList, 0, 1, true)
	menuFlex.AddItem(tagsFlex, 0, 1, false)

	view.AddItem(menuFlex, 0, 1, true)
	view.AddItem(infoFlex, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			app.SetFocus(navList)
			navFlex.Box.SetBorderColor(tcell.ColorHotPink)
			datesFlex.Box.SetBorderColor(tcell.ColorWhite)
			tagsFlex.Box.SetBorderColor(tcell.ColorWhite)
		case tcell.KeyF2:
			app.SetFocus(datesList)
			navFlex.Box.SetBorderColor(tcell.ColorWhite)
			datesFlex.Box.SetBorderColor(tcell.ColorHotPink)
			tagsFlex.Box.SetBorderColor(tcell.ColorWhite)
		case tcell.KeyF3:
			app.SetFocus(tagsList)
			navFlex.Box.SetBorderColor(tcell.ColorWhite)
			datesFlex.Box.SetBorderColor(tcell.ColorWhite)
			tagsFlex.Box.SetBorderColor(tcell.ColorHotPink)
		}
		return event
	})

	if err := app.SetRoot(view, true).Run(); err != nil {
		panic(err)
	}
}
