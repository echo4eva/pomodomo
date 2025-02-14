package ui

import (
	"strconv"
	"time"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatsUI struct {
	BaseUI[*tview.Flex]
	nav           *tview.Flex
	info          *tview.Flex
	infoList      *tview.List
	timeframeFlex *tview.Flex
	timeframeList *tview.List
	dateFlex      *tview.Flex
	dateList      *tview.List
}

func createFlex(title string) *tview.Flex {
	flex := tview.NewFlex()
	if len(title) != 0 {
		flex.Box.
			SetTitle(title).
			SetBorder(true)
	}
	return flex
}

func (sui *StatsUI) initializeTUI() {
	sui.nav = createFlex("").SetDirection(tview.FlexRow)
	sui.info = createFlex("Info")
	sui.infoList = tview.NewList()
	sui.info.AddItem(sui.infoList, 0, 1, false)
	sui.view.AddItem(sui.nav, 0, 1, true)
	sui.view.AddItem(sui.info, 0, 1, false)

	sui.timeframeFlex = createFlex("Timeframes")
	sui.dateFlex = createFlex("Time periods")
	sui.nav.AddItem(sui.timeframeFlex, 0, 1, true)
	sui.nav.AddItem(sui.dateFlex, 0, 1, false)

	sui.timeframeList = tview.NewList().ShowSecondaryText(false).
		AddItem("Day", "", 0, nil).
		AddItem("Week", "", 0, nil).
		AddItem("Month", "", 0, nil).
		AddItem("Year", "", 0, nil).
		AddItem("Alltime", "", 0, nil)
	sui.timeframeFlex.AddItem(sui.timeframeList, 0, 1, true)

	sui.dateList = tview.NewList().ShowSecondaryText(false)
	sui.dateFlex.AddItem(sui.dateList, 0, 1, false)

	sui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			sui.app.SetFocus(sui.timeframeList)
			sui.timeframeFlex.Box.SetBorderColor(tcell.ColorHotPink)
			sui.dateFlex.Box.SetBorderColor(tcell.ColorWhite)
		case tcell.KeyF2:
			sui.app.SetFocus(sui.dateList)
			sui.timeframeFlex.Box.SetBorderColor(tcell.ColorWhite)
			sui.dateFlex.Box.SetBorderColor(tcell.ColorHotPink)
		}
		return event
	})

	sui.timeframeList.SetSelectedFunc(func(index int, main string, secondary string, shortcut rune) {
		sui.dateList.Clear()
		sui.app.SetFocus(sui.dateList)
		sui.dateFlex.Box.SetBorderColor(tcell.ColorHotPink)
		sui.timeframeFlex.Box.SetBorderColor(tcell.ColorWhite)
		switch main {
		case "Day":
			sui.displayDays()
		case "Week":
			sui.displayWeeks()
		case "Month":
			sui.displayMonths()
		case "Year":
			sui.displayYears()
		case "Alltime":
			sui.displayAlltime()
		}
	})
}

func (sui *StatsUI) displayDays() {
	rows, err := sui.db.RetrieveDailySummary()
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		sui.dateList.AddItem(row.Date, "", 0, func() {
			sui.displayStats(row)
			sui.displayTaskStats(row.Date, "Day")
		})
	}
}

func (sui *StatsUI) displayWeeks() {
	rows, err := sui.db.RetrieveWeeklySummary()
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		sui.dateList.AddItem(row.DateRange, "", 0, func() {
			sui.displayStats(row)
			sui.displayTaskStats(row.Date, "Week")
		})
	}
}

func (sui *StatsUI) displayMonths() {
	rows, err := sui.db.RetrieveMonthlySummary()
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		sui.dateList.AddItem(row.Date, "", 0, func() {
			sui.displayStats(row)
			sui.displayTaskStats(row.Date, "Month")
		})
	}
}

func (sui *StatsUI) displayYears() {
	rows, err := sui.db.RetrieveYearlySummary()
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		sui.dateList.AddItem(row.Date, "", 0, func() {
			sui.displayStats(row)
			sui.displayTaskStats(row.Date, "Year")
		})
	}
}

func (sui *StatsUI) displayAlltime() {
	rows, err := sui.db.RetrieveYearlySummary()
	if err != nil {
		panic(err)
	}
	for _, row := range rows {
		sui.dateList.AddItem("KAPPA PENIS", "", 0, func() {
			sui.displayStats(row)
			sui.displayTaskStats("", "Alltime")
		})
	}
}

func convertDuration(duration int) string {
	var t time.Time
	t = t.Add(time.Duration(duration) * time.Second)
	return t.Format(time.TimeOnly)
}

func (sui *StatsUI) displayStats(stats database.SessionSummary) {
	sui.infoList.Clear()
	if stats.Date != "" && stats.DateRange != "" {
		sui.infoList.AddItem("Date", stats.DateRange, 0, nil)
	} else if stats.Date != "" && stats.DateRange == "" {
		sui.infoList.AddItem("Date", stats.Date, 0, nil)
	}
	sui.infoList.AddItem("Total Duration", convertDuration(stats.TotalDuration), 0, nil)
	sui.infoList.AddItem("Completed Sessions", strconv.Itoa(stats.CompletedSessions), 0, nil)
	sui.infoList.AddItem("Total Sessions", strconv.Itoa(stats.TotalSessions), 0, nil)
}

func (sui *StatsUI) displayTaskStats(date string, timeframe string) {
	var tasks []database.TaskSummary
	var err error

	switch timeframe {
	case "Day":
		tasks, err = sui.db.RetrieveDailyTaskSummary(date)
	case "Week":
		tasks, err = sui.db.RetrieveWeeklyTaskSummary(date)
	case "Month":
		tasks, err = sui.db.RetrieveMonthlyTaskSummary(date)
	case "Year":
		tasks, err = sui.db.RetrieveYearlyTaskSummary(date)
	case "Alltime":
		tasks, err = sui.db.RetrieveAlltimeTaskSummary()
	}

	if err != nil {
		panic(err)
	}

	for _, task := range tasks {
		sui.infoList.AddItem(task.TaskName, strconv.Itoa(task.TotalDuration), 0, nil)
	}
}

func StatsExec() {
	app := tview.NewApplication()
	view := tview.NewFlex()
	db, err := database.New()
	if err != nil {
		panic(err)
	}

	sui := &StatsUI{
		BaseUI: Initialize(app, view, db),
	}

	sui.initializeTUI()
	if err := sui.app.SetRoot(sui.view, true).Run(); err != nil {
		panic(err)
	}
}
