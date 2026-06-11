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
	Guid               uint64
	Pointer            uintptr
	Type               uint8
	Name               string
	Level              int32
	MaxHealth          int32
	CurrentHealth      int32
	Position           Vec3
	TargetGuid         uint64
	CurrentMana        int32
	MaxMana            int32
	Rage               int32
	Energy             int32
	CurrentSpellcastId int32
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
