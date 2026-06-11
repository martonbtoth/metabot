package objecttree

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"metabot/game"
	"metabot/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mitchellh/mapstructure"
)

type ObjectTree struct {
	widget.Tree

	g game.Game
}

func NewObjectTree(g game.Game) *ObjectTree {
	logger := logger.GetLogger()
	objectTree := &ObjectTree{g: g}
	objectTree.ExtendBaseWidget(objectTree)
	objectTree.ChildUIDs = func(id widget.TreeNodeID) []widget.TreeNodeID {
		if g.GetPlayerGuid() == 0 {
			return []string{"Not logged in, please log in first"}
		}
		objects := g.GetVisibleObjects()
		if id == "" {
			nodes := []widget.TreeNodeID{}
			for _, wowObject := range objects {
				nodes = append(nodes, fmt.Sprintf("wowobject:%d", wowObject.Guid))
			}
			return nodes
		}

		if isWowObjectId(id) {
			wowObject := g.GetVisibleObjectByGuid(getWowObjectGuid(id))
			mappedAttrs := map[string]interface{}{}
			mapstructure.Decode(wowObject, &mappedAttrs)
			leaves := []string{}
			for k, v := range mappedAttrs {
				s := formatValue(k, v)
				leaves = append(leaves, "attr:"+s+id)
			}
			sort.Strings(leaves)
			return leaves
		}

		return []string{}
	}
	objectTree.IsBranch = func(id widget.TreeNodeID) bool {
		return id == "" || isWowObjectId(id)
	}
	objectTree.CreateNode = func(branch bool) fyne.CanvasObject {
		if branch {
			return newWowObjectNode(nil)
		}
		moveButton := widget.NewButton("Move", func() {})
		label := widget.NewLabel("Template")
		attrContainer := container.New(layout.NewBorderLayout(nil, nil, label, moveButton), label, moveButton)
		return attrContainer
	}
	objectTree.UpdateNode = func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
		if isWowObjectId(id) {
			wowObject := g.GetVisibleObjectByGuid(getWowObjectGuid(id))
			wowObjectNode := o.(*wowObjectNode)
			wowObjectNode.SetWowObject(wowObject)
		} else {
			container := o.(*fyne.Container)
			label := container.Objects[0].(*widget.Label)
			moveButton := container.Objects[1].(*widget.Button)
			cut, _ := strings.CutPrefix(id, "attr:")
			split := strings.Split(cut, "wowobject:")
			label.SetText(split[0])
			moveButton.OnTapped = func() {
				guidString := split[1]
				guid, _ := strconv.ParseUint(guidString, 10, 64)
				object := g.GetVisibleObjectByGuid(guid)
				logger.Log("Unit found: " + object.Name)
				g.MoveToPosition(object.Position)
			}
			moveButton.Hidden = !strings.Contains(id, "Position")
		}

	}
	return objectTree
}

func isWowObjectId(id widget.TreeNodeID) bool {
	return strings.HasPrefix(id, "wowobject:")
}

func getWowObjectGuid(w widget.TreeNodeID) uint64 {
	cut, _ := strings.CutPrefix(w, "wowobject:")
	guid, _ := strconv.ParseUint(cut, 10, 64)
	return guid
}

func formatValue(k string, v interface{}) string {
	_, isGuid := v.(uint64)
	if isGuid {
		return fmt.Sprintf("%s: 0x%X", k, v)
	}

	if k == "Pointer" {
		return fmt.Sprintf("%s: 0x%X", k, v)
	}

	if k == "Type" {
		return fmt.Sprintf("%s: %s", k, game.ObjectTypeToString(v.(uint8)))
	}
	return fmt.Sprintf("%s: %v", k, v)
}
