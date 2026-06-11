package gui

import (
	"fmt"
	"superbot/gui/objecttree"
	"superbot/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func luaPage(gui *gui) *fyne.Container {
	luaTextField := widget.NewEntry()
	luaTextField.MultiLine = true
	luaTextField.Text = "DefaultServerLogin(\"totlol\", \"a\")\nCharacterSelect_SelectCharacter(1)\nEnterWorld()"
	runLuaButton := widget.NewButton("Run", func() {
		lua := luaTextField.Text
		results := gui.g.RunLuaWithResults(lua)
		logger.GetLogger().Log(fmt.Sprintf("Lua results: %v", results))
	})
	return container.New(layout.NewBorderLayout(nil, runLuaButton, nil, nil), luaTextField, runLuaButton)
}

func objectTreePage(gui *gui) *fyne.Container {
	enumerateObjectsButton := widget.NewButton("Enumerate\nobjects", gui.EnumerateVisibleObjects)
	tree := objecttree.NewObjectTree(gui.g)
	enumerateButtonLayout := container.New(layout.NewCenterLayout(), enumerateObjectsButton)
	gui.treeContainer = container.New(layout.NewBorderLayout(nil, enumerateButtonLayout, nil, nil), tree, enumerateButtonLayout)
	return gui.treeContainer
}

func terminalPage(gui *gui) *fyne.Container {
	gui.terminal = widget.NewTextGrid()
	scroll := container.NewScroll(gui.terminal)
	gui.scroll = scroll
	scroll.SetMinSize(fyne.NewSize(600, 300))

	gui.tabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Terminal" {
			scroll.ScrollToBottom()
		}
	}

	terminalContainer := container.New(layout.NewStackLayout(), scroll)
	return terminalContainer
}

func coolButtonsPage(gui *gui) *fyne.Container {
	clickSomethingButton := widget.NewButton("Log something", func() {
		logger.GetLogger().Log("Button clicked")
	})
	listSpellsButton := widget.NewButton("List spells", func() {
		spells := gui.g.GetAvailableSpells()
		logger.GetLogger().Log(fmt.Sprintf("Spells: %v", spells))
	})
	return container.New(
		layout.NewGridWrapLayout(fyne.Size{Width: 120, Height: 60}),
		clickSomethingButton,
		listSpellsButton,
	)
}

func buildGui(gui *gui) {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	gui.window = a.NewWindow("Superbot")
	gui.window.Resize(fyne.NewSize(800, 400))

	gui.tabs = container.NewAppTabs()

	gui.tabs.SetItems(
		[]*container.TabItem{
			container.NewTabItem("Terminal", terminalPage(gui)),
			container.NewTabItem("Objects", objectTreePage(gui)),
			container.NewTabItem("Lua shell", luaPage(gui)),
			container.NewTabItem("Cool buttons", coolButtonsPage(gui)),
		},
	)

	gui.window.SetContent(gui.tabs)
}
