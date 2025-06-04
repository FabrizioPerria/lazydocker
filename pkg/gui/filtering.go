package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleMainFilter() error {
	if gui.State.FilterMain.active {
		gui.State.FilterMain.active = false
		if err := gui.clearFilterMain(); err != nil {
			return err
		}
		return gui.returnFocus()
	} else {
		gui.State.FilterMain.active = true
		return gui.switchFocus(gui.Views.FilterMain)
	}
	return nil
}

func (gui *Gui) handleOpenFilter() error {
	panel, ok := gui.currentListPanel()
	if !ok {
		return nil
	}

	if panel.IsFilterDisabled() {
		return nil
	}

	gui.State.Filter.active = true
	gui.State.Filter.panel = panel

	return gui.switchFocus(gui.Views.Filter)
}

func (gui *Gui) onNewFilterNeedle(value string) error {
	gui.State.Filter.needle = value
	gui.ResetOrigin(gui.State.Filter.panel.GetView())
	return gui.State.Filter.panel.RerenderList()
}

func (gui *Gui) onNewFilterNeedleMain(value string) error {
	gui.State.FilterMain.needle = value
	return gui.Views.Main.Search(value)
}

func (gui *Gui) wrapEditor(f func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
		matched := f(v, key, ch, mod)
		if matched {
			if err := gui.onNewFilterNeedle(v.TextArea.GetContent()); err != nil {
				gui.Log.Error(err)
			}
		}
		return matched
	}
}
func (gui *Gui) wrapEditorMain(f func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
		matched := f(v, key, ch, mod)
		if matched {
			if err := gui.onNewFilterNeedleMain(v.TextArea.GetContent()); err != nil {
				gui.Log.Error(err)
			}
		}
		return matched
	}
}

func (gui *Gui) escapeFilterPrompt() error {
	gui.clearFilter()
	gui.clearFilterMain()
	// if gui.State.Filter.active {
	// 	if err := gui.clearFilter(); err != nil {
	// 		return err
	// 	}
	// }
	// if gui.State.FilterMain.active {
	// 	if err := gui.clearFilterMain(); err != nil {
	// 		return err
	// 	}
	// }

	return gui.returnFocus()
}

func (gui *Gui) clearFilter() error {
	gui.State.Filter.needle = ""
	gui.State.Filter.active = false
	panel := gui.State.Filter.panel
	gui.State.Filter.panel = nil
	gui.Views.Filter.ClearTextArea()

	if panel == nil {
		return nil
	}

	gui.ResetOrigin(panel.GetView())

	return panel.RerenderList()
}

func (gui *Gui) clearFilterMain() error {
	gui.State.FilterMain.needle = ""
	gui.State.FilterMain.active = false
	gui.Views.Main.Search("")
	// panel := gui.State.Filter.panel
	// gui.State.FilterMain.panel = nil
	gui.Views.FilterMain.ClearTextArea()

	return nil
}

// returns to the list view with the filter still applied
func (gui *Gui) commitFilter() error {
	if gui.State.Filter.needle == "" {
		if err := gui.clearFilter(); err != nil {
			return err
		}
	}

	return gui.returnFocus()
}
func (gui *Gui) commitFilterMain() error {
	if gui.State.FilterMain.needle == "" {
		if err := gui.clearFilterMain(); err != nil {
			return err
		}
	}

	return gui.returnFocus()
}

func (gui *Gui) filterPrompt() string {
	return fmt.Sprintf("%s: ", gui.Tr.FilterPrompt)
}
