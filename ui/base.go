package ui

import (
	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/rivo/tview"
)

type BaseUI[T tview.Primitive] struct {
	app  *tview.Application
	view T
	db   *database.Database
}

func Initialize[T tview.Primitive](app *tview.Application, view T, db *database.Database) BaseUI[T] {
	return BaseUI[T]{
		app:  app,
		view: view,
		db:   db,
	}
}
