package gui

import (
	"superbot/gui/objecttree"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildGui(gui *gui) {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	gui.window = a.NewWindow("Superbot")
	gui.window.Resize(fyne.NewSize(800, 400))

	gui.terminal = widget.NewTextGrid()
	scroll := container.NewScroll(gui.terminal)
	gui.scroll = scroll
	scroll.SetMinSize(fyne.NewSize(600, 300))

	enumerateObjectsButton := widget.NewButton("Enumerate\nobjects", gui.EnumerateVisibleObjects)

	terminalContainer := container.New(layout.NewStackLayout(), scroll)

	tree := objecttree.NewObjectTree(gui.g)
	enumerateButtonLayout := container.New(layout.NewCenterLayout(), enumerateObjectsButton)
	gui.treeContainer = container.New(layout.NewBorderLayout(nil, enumerateButtonLayout, nil, nil), tree, enumerateButtonLayout)
	gui.tabs = container.NewAppTabs(
		container.NewTabItem("Terminal", terminalContainer),
		container.NewTabItem("Objects", gui.treeContainer),
	)

	gui.tabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Terminal" {
			scroll.ScrollToBottom()
		}
	}

	gui.window.SetContent(gui.tabs)
}
