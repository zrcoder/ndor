package main

import (
	"embed"

	"github.com/ELY-Open-Source/CrychicUI/navigation"
	"github.com/egoist/mygo/ui"
	"github.com/zrcoder/ndor/cmd/gui/internal/ndor"
	"github.com/zrcoder/ndor/cmd/gui/internal/turtle"
)

//go:embed resources
var resourcesFS embed.FS

type App struct {
	ndor        *ndor.App
	turtle      *turtle.App
	selectedTab int
}

func NewApp() *App {
	return &App{
		ndor:   ndor.New(resourcesFS),
		turtle: turtle.New(resourcesFS),
	}
}

func (a *App) view(c *ui.Context) {
	navigation.UnderlineTabs(c, "tabs", []string{"Ndor", "Turtle"}, &a.selectedTab)
	if a.selectedTab == 0 {
		a.ndor.View(c)
	} else {
		a.turtle.View(c)
	}
}
