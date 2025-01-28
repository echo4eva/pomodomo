package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

var (
	app  *tview.Application
	view *tview.Modal
)

func getTimer(end time.Time) string {
	difference := time.Until(end)
	difference = difference.Round(time.Second)
	countdown := formatTimer(difference.String())

	return fmt.Sprint(countdown)
}

func formatTimer(duration string) string {
	var formattedCountdown string

	hoursExists := false
	minutesExists := false
	secondsExist := false

	hourIndex := strings.Index(duration, "h")
	if hourIndex != -1 {
		hoursExists = true
	}
	minuteIndex := strings.Index(duration, "m")
	if minuteIndex != -1 {
		minutesExists = true
	}
	secondIndex := strings.Index(duration, "s")
	if secondIndex != -1 {
		secondsExist = true
	}

	if hoursExists && minutesExists && secondsExist {
		hours := duration[:hourIndex]
		minutes := duration[hourIndex+1 : minuteIndex]
		seconds := duration[minuteIndex+1 : secondIndex]
		formattedCountdown = fmt.Sprintf("%s hours %s minutes %s seconds", hours, minutes, seconds)
	} else if minutesExists && secondsExist {
		minutes := duration[:minuteIndex]
		seconds := duration[minuteIndex+1 : secondIndex]
		formattedCountdown = fmt.Sprintf("%s minutes %s seconds", minutes, seconds)
	} else {
		seconds := duration[:secondIndex]
		formattedCountdown = fmt.Sprintf("%s seconds", seconds)
	}

	return formattedCountdown
}

func updateTimer(end time.Time) {
	for {
		time.Sleep(1 * time.Second)
		timer := getTimer(end)
		app.QueueUpdateDraw(func() {
			view.SetText(fmt.Sprint(timer, " left"))
		})
		if timer == "0 seconds" {
			break
		}
	}
}

func Exec(t0, t1 time.Time, task string) {
	app = tview.NewApplication()
	view = tview.NewModal().
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
