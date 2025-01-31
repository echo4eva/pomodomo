package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app           *tview.Application
	view          *tview.Modal
	db            *database.Database
	start         *time.Time
	scheduled_end *time.Time
	task          *database.Task
	description   *string
}

func (ui UI) getTimer() string {
	difference := time.Until(*ui.scheduled_end)
	difference = difference.Round(time.Second)
	countdown := formatTimer(difference.String())
	return fmt.Sprint(countdown)
}

func (ui UI) getElapsedTime() string {
	elapsed := time.Since(*ui.start)
	return formatTimer(elapsed.Round(time.Second).String())
}

func (ui UI) getCompletion() int {
	difference := time.Until(*ui.scheduled_end)
	difference = difference.Round(time.Second)
	if difference <= (time.Second * 0) {
		return 1
	}
	return 0
}

func formatTimer(duration string) string {
	fillZeroes := func(s string) string {
		if len(s) < 2 {
			return "0" + s
		}
		return s
	}

	var formattedCountdown string

	hourIndex := strings.Index(duration, "h")
	minuteIndex := strings.Index(duration, "m")
	secondIndex := strings.Index(duration, "s")

	switch {
	case hourIndex != -1 && minuteIndex != -1 && secondIndex != -1:
		formattedCountdown = fmt.Sprintf("%s:%s:%s",
			fillZeroes(duration[:hourIndex]),
			fillZeroes(duration[hourIndex+1:minuteIndex]),
			fillZeroes(duration[minuteIndex+1:secondIndex]),
		)
	case minuteIndex != -1 && secondIndex != -1:
		formattedCountdown = fmt.Sprintf("00:%s:%s",
			fillZeroes(duration[:minuteIndex]),
			fillZeroes(duration[minuteIndex+1:secondIndex]),
		)
	default:
		formattedCountdown = fmt.Sprintf("00:00:%s",
			fillZeroes(duration[:secondIndex]),
		)
	}

	return formattedCountdown
}

func (ui *UI) updateTimer() {
	for {
		time.Sleep(1 * time.Second)
		timer := ui.getTimer()
		ui.app.QueueUpdateDraw(func() {
			ui.view.SetText(fmt.Sprint(timer, " left"))
		})
		if timer == "0" {
			ui.app.QueueUpdateDraw(func() {
				ui.view.SetText("You're done! Lock in if you're tapped in!")
			})
			break
		}
	}
}

func (ui *UI) captureInterruptSignal(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		ui.completeSession()
	}

	return event
}

func (ui *UI) completeSession() {
	ui.db.AddSession(database.Session{
		Duration:     ui.getElapsedTime(),
		TaskID:       (*ui.task).Id,
		Start:        (*ui.start).Format(time.RFC3339),
		ScheduledEnd: (*ui.scheduled_end).Format(time.RFC3339),
		EndedAt:      time.Now().Format(time.RFC3339),
		Completed:    ui.getCompletion(),
	})
}

func Exec(start, scheduled_end time.Time, taskName string) {
	db, err := database.New()
	if err != nil {
		panic(err)
	}
	app := tview.NewApplication()
	view := tview.NewModal()
	task, err := db.RetrieveTaskByName(taskName)
	if err != nil {
		panic(err)
	}

	ui := &UI{
		app:           app,
		view:          view,
		db:            db,
		start:         &start,
		scheduled_end: &scheduled_end,
		task:          task,
	}

	ui.app.SetInputCapture(ui.captureInterruptSignal)
	ui.view.
		AddButtons([]string{"give up :("}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "give up :(" {
				ui.completeSession()
				app.Stop()
			}
		})

	go ui.updateTimer()
	if err := ui.app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
