package game

import "fmt"

const (
	UnitType_None              = 0
	UnitType_Item              = 1
	UnitType_Container         = 2
	UnitType_Unit              = 3
	UnitType_Player            = 4
	UnitType_GameObject        = 5
	UnitType_DynamicObject     = 6
	UnitType_Corpse            = 7
	DynamicFlag_None           = 0x0
	DynamicFlag_CanBeLooted    = 0x1
	DynamicFlag_IsMarked       = 0x2
	DynamicFlag_Tapped         = 0x4 // Makes creature name tag appear grey
	DynamicFlag_TappedByMe     = 0x8
	UnitFlag_Unknown0          = 0x00000001 // Movement checks disabled, likely paired with loss of client control packet. We use it to add custom cliffwalking to GM mode until actual usecases will be known.
	UnitFlag_NonAttackable     = 0x00000002 // not attackable
	UnitFlag_ClientControlLost = 0x00000004 // Generic unspecified loss of control initiated by server script, movement checks disabled, paired with loss of client control packet.
	UnitFlag_PlayerController  = 0x00000008 // players, pets, totems, guardians, companions, charms, any units associated with players
	UnitFlag_Rename            = 0x00000010 // ??
	UnitFlag_Preparation       = 0x00000020 // don't take reagents for spells with SPELL_ATTR_EX5_NO_REAGENT_WHILE_PREP
	UnitFlag_Unknown6          = 0x00000040 // ??
	UnitFlag_NotAttackable_1   = 0x00000080 // ?? (UnitFlag_PvpAttackable | UnitFlag_NOT_ATTACKABLE_1) is NON_PVP_ATTACKABLE
	UnitFlag_ImmuneToPlayer    = 0x00000100 // Target is immune to players
	UnitFlag_ImmuneToNpc       = 0x00000200 // Target is immune to Non-Player Characters
	UnitFlag_Looting           = 0x00000400 // loot animation
	UnitFlag_PetInCombat       = 0x00000800 // in combat?, 2.0.8
	UnitFlag_Pvp               = 0x00001000 // is flagged for pvp
	UnitFlag_Silenced          = 0x00002000 // silenced, 2.1.1
	UnitFlag_Persuaded         = 0x00004000 // persuaded, 2.0.8
	UnitFlag_Swimming          = 0x00008000 // controls water swimming animation - TODO: confirm whether dynamic or static
	UnitFlag_NonAttackable2    = 0x00010000 // removes attackable icon, if on yourself, cannot assist self but can cast TARGET_UNIT_CASTER spells - added by SPELL_AURA_MOD_UNATTACKABLE
	UnitFlag_Pacified          = 0x00020000 // probably like the paladin's Repentance spell
	UnitFlag_Stunned           = 0x00040000 // Unit is a subject to stun, turn and strafe movement disabled
	UnitFlag_InCombat          = 0x00080000
	UnitFlag_TaxiFlight        = 0x00100000 // Unit is on taxi, paired with a duplicate loss of client control packet (likely a legacy serverside hack). Disables any spellcasts not allowed in taxi flight client-side.
	UnitFlag_Disarmed          = 0x00200000 // disable melee spells casting..., "Required melee weapon" added to melee spells tooltip.
	UnitFlag_Confused          = 0x00400000 // Unit is a subject to confused movement, movement checks disabled, paired with loss of client control packet.
	UnitFlag_Fleeing           = 0x00800000 // Unit is a subject to fleeing movement, movement checks disabled, paired with loss of client control packet.
	UnitFlag_Possessed         = 0x01000000 // Unit is under remote control by another unit, movement checks disabled, paired with loss of client control packet. New master is allowed to use melee attack and can't select this unit via mouse in the world (as if it was own character).
	UnitFlag_NotSelectable     = 0x02000000
	UnitFlag_Skinnable         = 0x04000000
	UnitFlag_Mount             = 0x08000000 // is mounted?
	UnitFlag_Unknown28         = 0x10000000 // ??
	UnitFlag_Unknown29         = 0x20000000 // used in Feing Death spell
	UnitFlag_Sheathe           = 0x40000000 // ??
	UnitFlag_Immune            = 0x80000000
	UnitReaction_Hated         = 0
	UnitReaction_Hostile       = 1
	UnitReaction_Unfriendly    = 2
	UnitReaction_Neutral       = 3
	UnitReaction_Friendly      = 4
	UnitReaction_Honored       = 5
	UnitReaction_Revered       = 6
	UnitReaction_Exalted       = 7
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
	DynamicFlags       []string
	UnitFlags          []string
	Reaction		   string
}

func (w WowObject) String() string {
	return fmt.Sprintf("WowObject { GUID: 0x%X, Type: %v}", w.Guid, ObjectTypeToString(w.Type))
}

func ObjectTypeToString(objectType uint8) string {
	switch objectType {
	case UnitType_None:
		return "None"
	case UnitType_Item:
		return "Item"
	case UnitType_Container:
		return "Container"
	case UnitType_Unit:
		return "Unit"
	case UnitType_Player:
		return "Player"
	case UnitType_GameObject:
		return "GameObject"
	case UnitType_DynamicObject:
		return "DynamicObject"
	case UnitType_Corpse:
		return "Corpse"
	}
	return "Unknown"
}

func HasFlag(byteValue uint32, flag uint32) bool {
	return byteValue&flag != 0
}
