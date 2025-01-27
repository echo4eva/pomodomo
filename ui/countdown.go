package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

var (
	app  *tview.Application
	view *tview.Modal
)

func getTimer(end time.Time) string {
	now := time.Now()
	difference := end.Sub(now)
	return fmt.Sprint(difference.String())
}

func updateTimer(end time.Time) {
	for {
		time.Sleep(500 * time.Millisecond)
		app.QueueUpdateDraw(func() {
			view.SetText(getTimer(end))
		})
	}
}

func Exec(t0, t1 time.Time, task string) {
	app = tview.NewApplication()
	view = tview.NewModal().
		SetText(getTimer(t1)).
		AddButtons([]string{"give up :("}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "give up :(" {
				app.Stop()
			}
		})

	go updateTimer(t1)
	if err := app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
