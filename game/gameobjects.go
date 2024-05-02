package game

import "fmt"

const (
	None          = 0
	Item          = 1
	Container     = 2
	Unit          = 3
	Player        = 4
	GameObject    = 5
	DynamicObject = 6
	Corpse        = 7
)

type WowObject struct {
	Guid          uint64
	Pointer       uintptr
	Type          uint8
	Name          string
	MaxHealth     int
	CurrentHealth int
	Position      Vec3
}

func (w WowObject) String() string {
	return fmt.Sprintf("WowObject { GUID: 0x%X, Type: %v}", w.Guid, ObjectTypeToString(w.Type))
}

func ObjectTypeToString(objectType uint8) string {
	switch objectType {
	case None:
		return "None"
	case Item:
		return "Item"
	case Container:
		return "Container"
	case Unit:
		return "Unit"
	case Player:
		return "Player"
	case GameObject:
		return "GameObject"
	case DynamicObject:
		return "DynamicObject"
	case Corpse:
		return "Corpse"
	}
	return "Unknown"
}
