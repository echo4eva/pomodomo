package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/rivo/tview"
)

var (
	app  *tview.Application
	view *tview.Modal
	db   *database.Database
)

func getTimer(end time.Time) string {
	difference := time.Until(end)
	difference = difference.Round(time.Second)
	countdown := formatTimer(difference.String())

	return fmt.Sprint(countdown)
}

func getElapsedTime(start time.Time) string {
	elapsed := time.Since(start)
	return formatTimer(elapsed.Round(time.Second).String())
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
		formattedCountdown = fmt.Sprintf("%s:%s:%s", hours, minutes, seconds)
	} else if minutesExists && secondsExist {
		minutes := duration[:minuteIndex]
		seconds := duration[minuteIndex+1 : secondIndex]
		formattedCountdown = fmt.Sprintf("%s:%s", minutes, seconds)
	} else {
		seconds := duration[:secondIndex]
		formattedCountdown = fmt.Sprint(seconds)
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
		if timer == "0" {
			break
		}
	}
}

func Exec(start, end time.Time, task string) {
	db, err := database.New()
	if err != nil {
		panic(err)
	}
	app = tview.NewApplication()
	view = tview.NewModal().
		AddButtons([]string{"give up :("}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "give up :(" {
				db.AddSession(database.Session{
					Date:     time.Now().Format(time.DateOnly),
					Duration: getElapsedTime(start),
					Task:     task,
				})
				app.Stop()
			}
		})
	go updateTimer(end)
	if err := app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
