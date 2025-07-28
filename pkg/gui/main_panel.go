package gui

import (
	"math"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) scrollUpMainFast() error {
	err := gui.scrollUpMain()
	for i := 0; i < 9; i++ {
		gui.scrollUpMain()
	}
	return err
}

func (gui *Gui) scrullDownMainFast() error {
	err := gui.scrollDownMain()
	for i := 0; i < 9; i++ {
		gui.scrollDownMain()
	}
	return err
}

func (gui *Gui) scrollRightMainFast() error {
	err := gui.scrollRightMain(nil, nil)
	for i := 0; i < 9; i++ {
		gui.scrollRightMain(nil, nil)
	}
	return err
}

func (gui *Gui) scrollLeftMainFast() error {
	err := gui.scrollLeftMain(nil, nil)
	for i := 0; i < 9; i++ {
		gui.scrollLeftMain(nil, nil)
	}
	return err
}

func (gui *Gui) scrollUpMain() error {
	mainView := gui.Views.Main
	mainView.Autoscroll = false
	ox, oy := mainView.Origin()
	newOy := int(math.Max(0, float64(oy-gui.Config.UserConfig.Gui.ScrollHeight)))
	return mainView.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownMain() error {
	mainView := gui.Views.Main
	mainView.Autoscroll = false
	ox, oy := mainView.Origin()

	reservedLines := 0
	if !gui.Config.UserConfig.Gui.ScrollPastBottom {
		_, sizeY := mainView.Size()
		reservedLines = sizeY
	}

	totalLines := mainView.ViewLinesHeight()
	if oy+reservedLines >= totalLines {
		return nil
	}

	return mainView.SetOrigin(ox, oy+gui.Config.UserConfig.Gui.ScrollHeight)
}

func (gui *Gui) scrollLeftMain(g *gocui.Gui, v *gocui.View) error {
	mainView := gui.Views.Main
	ox, oy := mainView.Origin()
	newOx := int(math.Max(0, float64(ox-gui.Config.UserConfig.Gui.ScrollHeight)))

	return mainView.SetOrigin(newOx, oy)
}

func (gui *Gui) scrollRightMain(g *gocui.Gui, v *gocui.View) error {
	mainView := gui.Views.Main
	ox, oy := mainView.Origin()

	content := mainView.ViewBufferLines()
	var largestNumberOfCharacters int
	for _, txt := range content {
		if len(txt) > largestNumberOfCharacters {
			largestNumberOfCharacters = len(txt)
		}
	}

	sizeX, _ := mainView.Size()
	if ox+sizeX >= largestNumberOfCharacters {
		return nil
	}

	return mainView.SetOrigin(ox+gui.Config.UserConfig.Gui.ScrollHeight, oy)
}

func (gui *Gui) autoScrollMain(g *gocui.Gui, v *gocui.View) error {
	gui.Views.Main.Autoscroll = true
	return nil
}

func (gui *Gui) jumpToTopMain(g *gocui.Gui, v *gocui.View) error {
	gui.Views.Main.Autoscroll = false
	_ = gui.Views.Main.SetOrigin(0, 0)
	_ = gui.Views.Main.SetCursor(0, 0)
	return nil
}

func (gui *Gui) onMainTabClick(tabIndex int) error {
	gui.Log.Warn(tabIndex)

	currentSidePanel, ok := gui.currentSidePanel()

	if !ok {
		return nil
	}

	currentSidePanel.SetMainTabIndex(tabIndex)
	return currentSidePanel.HandleSelect()
}

func (gui *Gui) handleEnterMain(g *gocui.Gui, v *gocui.View) error {
	mainView := gui.Views.Main
	mainView.ParentView = v

	return gui.switchFocus(mainView)
}

func (gui *Gui) handleExitMain(g *gocui.Gui, v *gocui.View) error {
	v.ParentView = nil
	return gui.returnFocus()
}

func (gui *Gui) handleMainClick() error {
	if gui.popupPanelFocused() {
		return nil
	}

	currentView := gui.g.CurrentView()

	if currentView.Name() != "main" {
		gui.Views.Main.ParentView = currentView
	}

	return gui.switchFocus(gui.Views.Main)
}
