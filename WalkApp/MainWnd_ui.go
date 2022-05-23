// This file was created by ui2walk and may be regenerated.
// DO NOT EDIT OR YOUR MODIFICATIONS WILL BE LOST!
package WalkApp

import (
	"github.com/lxn/walk"
)

type mainWndActions struct {
	actionSetting *walk.Action
}

type mainWndUI struct {
	actions              mainWndActions
	centralWidget        *walk.Composite
	label_Country        *walk.Label
	combo_Country        *walk.ComboBox
	label_State          *walk.Label
	combo_State          *walk.ComboBox
	//list_City            *QListView
	pushButton_SelectAll *walk.PushButton
	label_City           *walk.Label
	pushButton_Start     *walk.PushButton
	//mainToolBar          *QToolBar
}

func (w *MainWnd) init() (err error) {
	if w.MainWindow, err = walk.NewMainWindow(); err != nil {
		return err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	var font *walk.Font
	if font == nil {
		font = nil
	}

	w.SetName("MainWnd")
	l := walk.NewVBoxLayout()
	if err := l.SetMargins(walk.Margins{0, 0, 0, 0}); err != nil {
		return err
	}
	if err := w.SetLayout(l); err != nil {
		return err
	}
	if err := w.SetClientSize(walk.Size{850, 609}); err != nil {
		return err
	}
	if err := w.SetTitle(`谷歌地图采集器`); err != nil {
		return err
	}


	// Actions

	// w.ui.actions.actionSetting
	w.ui.actions.actionSetting = walk.NewAction()
	if err := w.ui.actions.actionSetting.SetText(`设置`); err != nil {
		return err
	}

	// Menus

	// menu
	menu, err := walk.NewMenu()
	if err != nil {
		return err
	}
	menuAction, err := w.Menu().Actions().AddMenu(menu)
	if err != nil {
		return err
	}
	if err := menuAction.SetText(`程序`); err != nil {
		return err
	}

	if err := menu.Actions().Add(w.ui.actions.actionSetting); err != nil {
		return err
	}

	// centralWidget
	if w.ui.centralWidget, err = walk.NewComposite(w); err != nil {
		return err
	}
	w.ui.centralWidget.SetName("centralWidget")

	// label_Country
	if w.ui.label_Country, err = walk.NewLabel(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.label_Country.SetName("label_Country")
	if err := w.ui.label_Country.SetBounds(walk.Rectangle{40, 30, 81, 31}); err != nil {
		return err
	}
	if err := w.ui.label_Country.SetText(`选择国家:`); err != nil {
		return err
	}

	w.Children().Add(w.ui.label_Country)

	// combo_Country
	if w.ui.combo_Country, err = walk.NewComboBox(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.combo_Country.SetName("combo_Country")
	if err := w.ui.combo_Country.SetBounds(walk.Rectangle{130, 30, 131, 22}); err != nil {
		return err
	}
	w.Children().Add(w.ui.combo_Country)

	// label_State
	if w.ui.label_State, err = walk.NewLabel(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.label_State.SetName("label_State")
	if err := w.ui.label_State.SetBounds(walk.Rectangle{40, 100, 81, 31}); err != nil {
		return err
	}
	if err := w.ui.label_State.SetText(`选择省份:`); err != nil {
		return err
	}
	w.Children().Add(w.ui.label_State)

	// combo_State
	if w.ui.combo_State, err = walk.NewComboBox(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.combo_State.SetName("combo_State")
	if err := w.ui.combo_State.SetBounds(walk.Rectangle{130, 100, 131, 22}); err != nil {
		return err
	}
	w.Children().Add(w.ui.combo_State)
	// list_City
	//if w.ui.list_City, err = NewQListView(w.ui.centralWidget); err != nil {
	//	return err
	//}
	//w.ui.list_City.SetName("list_City")
	//if err := w.ui.list_City.SetBounds(walk.Rectangle{340, 80, 491, 311}); err != nil {
	//	return err
	//}

	// pushButton_SelectAll
	if w.ui.pushButton_SelectAll, err = walk.NewPushButton(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.pushButton_SelectAll.SetName("pushButton_SelectAll")
	if err := w.ui.pushButton_SelectAll.SetBounds(walk.Rectangle{490, 20, 181, 41}); err != nil {
		return err
	}
	if err := w.ui.pushButton_SelectAll.SetText(`全选`); err != nil {
		return err
	}

	// label_City
	if w.ui.label_City, err = walk.NewLabel(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.label_City.SetName("label_City")
	if err := w.ui.label_City.SetBounds(walk.Rectangle{390, 20, 81, 41}); err != nil {
		return err
	}
	if err := w.ui.label_City.SetText(`选择城市:`); err != nil {
		return err
	}

	// pushButton_Start
	if w.ui.pushButton_Start, err = walk.NewPushButton(w.ui.centralWidget); err != nil {
		return err
	}
	w.ui.pushButton_Start.SetName("pushButton_Start")
	if err := w.ui.pushButton_Start.SetBounds(walk.Rectangle{240, 440, 341, 91}); err != nil {
		return err
	}
	if err := w.ui.pushButton_Start.SetText(`开始任务`); err != nil {
		return err
	}

	// mainToolBar
	//if w.ui.mainToolBar, err = NewQToolBar(w); err != nil {
	//	return err
	//}
	//w.ui.mainToolBar.SetName("mainToolBar")

	// Tab order

	succeeded = true

	return nil
}
