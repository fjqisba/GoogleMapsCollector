package WalkApp

import (
	"github.com/lxn/walk"
)

type MainWnd struct {
	*walk.MainWindow
	ui mainWndUI
}

func RunMainWnd() (int, error) {
	mw := new(MainWnd)
	if err := mw.init(); err != nil {
		return 0, err
	}
	defer mw.Dispose()

	// TODO: Do further required setup, e.g. for event handling, here.

	mw.Show()

	return mw.Run(), nil
}
