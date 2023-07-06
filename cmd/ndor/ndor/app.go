package main

import (
	"context"
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/zrcoder/ndor/pkg"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GenImage(code string) map[string]string {
	runtime.LogInfo(a.ctx, code)
	src, err := pkg.Run(100, 100, code)
	errMsg := ""
	line := 0
	if err != nil {
		errMsg = err.Msg
		line = err.Number
	}
	return map[string]string{
		"err":  errMsg,
		"src":  src,
		"line": strconv.Itoa(line),
	}
}

func (a *App) AlertError(msg string) {
	a.alert(msg, runtime.ErrorDialog)
}

func (a *App) alert(msg string, dialogType runtime.DialogType) {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    dialogType,
		Message: msg,
	})
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
}
