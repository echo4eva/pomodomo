package ui

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TimerUI struct {
	BaseUI[*tview.Modal]
	start         *time.Time
	scheduled_end *time.Time
	task          *database.Task
	description   *string
}

func (tui TimerUI) getCountdown() string {
	difference := time.Until(*tui.scheduled_end)
	difference = difference.Round(time.Second)
	countdown := formatTime(difference.String())
	return fmt.Sprint(countdown)
}

func (tui TimerUI) getTimer() string {
	timeSince := time.Since(*tui.scheduled_end)
	timeSince = timeSince.Round(time.Second)
	countdown := formatTime(timeSince.String())
	return fmt.Sprint(countdown)
}

func (tui TimerUI) getElapsedTime() int {
	elapsed := time.Since(*tui.start)
	return int(elapsed.Seconds())
}

func (tui TimerUI) getCompletion() int {
	difference := time.Until(*tui.scheduled_end)
	difference = difference.Round(time.Second)
	if difference <= (time.Second * 0) {
		return 1
	}
	return 0
}

func formatTime(duration string) string {
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

func (tui *TimerUI) updateCountdown() {
	for {
		time.Sleep(1 * time.Second)
		timer := tui.getCountdown()
		tui.app.QueueUpdateDraw(func() {
			tui.view.SetText(fmt.Sprint(timer, " left"))
		})
	}
}

func (tui *TimerUI) updateTimer() {
	for {
		time.Sleep(1 * time.Second)
		timer := tui.getTimer()
		tui.app.QueueUpdateDraw(func() {
			tui.view.SetText(fmt.Sprint(timer, " elapsed"))
		})
	}
}

func (tui *TimerUI) captureInterruptSignal(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		tui.completeSession()
	}

	return event
}

func (tui *TimerUI) completeSession() {
	tui.db.AddSession(database.Session{
		Duration:     tui.getElapsedTime(),
		TaskID:       (*tui.task).Id,
		Start:        (*tui.start).Format(time.RFC3339),
		ScheduledEnd: (*tui.scheduled_end).Format(time.RFC3339),
		EndedAt:      time.Now().Format(time.RFC3339),
		Completed:    tui.getCompletion(),
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
	if err == sql.ErrNoRows {
		fmt.Printf("Invalid task name entered: %s\n", taskName)
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	tui := &TimerUI{
		BaseUI:        Initialize(app, view, db),
		start:         &start,
		scheduled_end: &scheduled_end,
		task:          task,
	}

	tui.app.SetInputCapture(tui.captureInterruptSignal)
	tui.view.
		AddButtons([]string{"give up :("}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "give up :(" {
				tui.completeSession()
				app.Stop()
			}
		})

	if tui.start.String() == tui.scheduled_end.String() {
		go tui.updateTimer()
	} else {
		go tui.updateCountdown()
	}

	if err := tui.app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
