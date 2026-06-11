package objecttree

import (
	"fmt"

	"metabot/game"

	"fyne.io/fyne/v2/widget"
)

type wowObjectNode struct {
	widget.Label

	w *game.WowObject
}

func newWowObjectNode(w *game.WowObject) *wowObjectNode {
	objectNode := &wowObjectNode{}
	objectNode.ExtendBaseWidget(objectNode)
	objectNode.SetWowObject(w)
	return objectNode
}

func (node *wowObjectNode) SetWowObject(w *game.WowObject) {
	node.w = w
	if w != nil {
		node.SetText(fmt.Sprintf("0x%X - %s - %s", w.Guid, w.Name, game.ObjectTypeToString(w.Type)))
	} else {
		node.SetText("Empty")
	}
}
