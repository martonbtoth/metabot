package game

/*
#include "native-bridge.h"
#cgo LDFLAGS: -L${SRCDIR}/minhook/ -l:libMinHook.a
*/
import "C"

import (
	_ "embed"
	"fmt"
	"math/rand"
	"metabot/logger"
	"slices"
	"strings"
	"unsafe"
)

const OBJECT_TYPE_OFFSET = 0x14
const DESCRIPTOR_OFFSET = 0x8

//go:embed lua/enumerate_spellbook.lua
var enumerateSpellbookLua string

//go:embed lua/attack_on.lua
var attackOnLua string

//go:embed lua/attack_off.lua
var attackOffLua string

type Game interface {
	GetPlayerGuid() uint64
	EnumerateVisibleObjects()
	EnumerateVisibleObjectsCallback(filter int, guid uint64)
	GetVisibleObjects() []WowObject
	GetVisibleObjectByGuid(guid uint64) *WowObject
	MoveToPosition(position Vec3)
	StopMovement()
	SetFacing(angle float32)
	FaceTarget()
	FacePosition(pos Vec3)
	Sit()
	SetTarget(guid uint64)
	RunLua(lua string)
	RunLuaWithResults(lua string) []string
	GetAvailableSpells() []string
	Jump()
	ToggleAttack(attack bool)
	CastSpellByName(spellName string)
	AutoLoot(guid uint64)
}

type game struct {
	objects []WowObject
}

var globalGame Game

func GetGame() Game {
	if globalGame == nil {
		globalGame = &game{}
	}
	FixClickToMove()
	UnlockProtectedLuaFunctions()
	logger.GetLogger().AddListener(func(logBuffer string, s string) {
		if globalGame.GetPlayerGuid() != 0 {
			split := strings.Split(s, "\n")
			for _, s := range split {
				globalGame.RunLua("DEFAULT_CHAT_FRAME:AddMessage('" + strings.ReplaceAll(s, "'", "\\'") + "')")
			}

		}
	})
	return globalGame
}

func (g *game) GetAvailableSpells() []string {

	if g.GetPlayerGuid() == 0 {
		return []string{}
	}

	rawSpellsString := g.RunLuaWithResults(enumerateSpellbookLua)[0]

	splitFn := func(c rune) bool {
		return c == '\n'
	}

	spells := strings.FieldsFunc(rawSpellsString, splitFn)

	return spells
}

func (g *game) GetPlayerGuid() uint64 {
	playerGuid := uint64(C.GetPlayerGuid())
	return playerGuid
}

func (g *game) EnumerateVisibleObjects() {
	g.objects = []WowObject{}
	if g.GetPlayerGuid() != 0 {
		C.EnumerateVisibleObjects(0)
	}
}

func (g *game) EnumerateVisibleObjectsCallback(filter int, guid uint64) {
	objectUintPtr := uintptr(C.GetObjectPtr(C.uint64_t(guid)))
	objectTypePtr := unsafe.Pointer(objectUintPtr + OBJECT_TYPE_OFFSET)
	objectType := *(*uint8)(objectTypePtr)

	position := Vec3{
		X: float32(C.GetObjectPositionX(C.uint64_t(guid))),
		Y: float32(C.GetObjectPositionY(C.uint64_t(guid))),
		Z: float32(C.GetObjectPositionZ(C.uint64_t(guid))),
	}

	g.objects = append(g.objects, WowObject{
		Guid:               guid,
		Pointer:            objectUintPtr,
		Type:               objectType,
		Name:               g.getName(objectType, guid),
		Level:              g.getLevel(objectType, guid),
		MaxHealth:          g.getMaxHealth(objectType, guid),
		CurrentHealth:      g.getCurrentHealth(objectType, guid),
		Position:           position,
		TargetGuid:         g.getTargetGuid(objectType, guid),
		CurrentMana:        g.getCurrentMana(objectType, guid),
		MaxMana:            g.getMaxMana(objectType, guid),
		Rage:               g.getCurrentRage(objectType, guid),
		Energy:             g.getCurrentEnergy(objectType, guid),
		CurrentSpellcastId: g.getCurrentSpellCastId(objectType, guid),
		DynamicFlags:       g.getDynamicFlags(objectType, guid),
		Reaction:           g.getReaction(objectType, guid),
	})
}

func (g *game) getName(objectType uint8, guid uint64) string {
	name := "N/A"
	if objectType == UnitType_Unit {
		name = C.GoString(C.GetUnitName(C.uint64_t(guid)))
	} else if objectType == UnitType_Player {
		name = C.GoString(C.GetPlayerName(C.uint64_t(guid)))
	}
	return name
}

func (g *game) getCurrentHealth(objectType uint8, guid uint64) int32 {
	currentHealth := 0
	if objectType == UnitType_Unit {
		currentHealth = int(C.GetCurrentHealth(C.uint64_t(guid)))
	} else if objectType == UnitType_Player {
		currentHealth = int(C.GetCurrentHealth(C.uint64_t(guid)))
	}
	return int32(currentHealth)
}

func (g *game) getMaxHealth(objectType uint8, guid uint64) int32 {
	currentHealth := 0
	if objectType == UnitType_Unit {
		currentHealth = int(C.GetMaxHealth(C.uint64_t(guid)))
	} else if objectType == UnitType_Player {
		currentHealth = int(C.GetMaxHealth(C.uint64_t(guid)))
	}
	return int32(currentHealth)
}

func (g *game) getTargetGuid(objectType uint8, guid uint64) uint64 {
	targetGuid := uint64(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		targetGuid = uint64(C.GetTargetGuid(C.uint64_t(guid)))
	}
	return targetGuid
}

func (g *game) getCurrentMana(objectType uint8, guid uint64) int32 {
	currentMana := int32(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		currentMana = int32(C.GetCurrentMana(C.uint64_t(guid)))
	}
	return currentMana
}

func (g *game) getMaxMana(objectType uint8, guid uint64) int32 {
	currentMana := int32(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		currentMana = int32(C.GetMaxMana(C.uint64_t(guid)))
	}
	return currentMana
}

func (g *game) getCurrentRage(objectType uint8, guid uint64) int32 {
	currentRage := int32(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		currentRage = int32(C.GetCurrentRage(C.uint64_t(guid)))
	}
	return currentRage
}

func (g *game) getCurrentEnergy(objectType uint8, guid uint64) int32 {
	currentEnergy := int32(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		currentEnergy = int32(C.GetCurrentEnergy(C.uint64_t(guid)))
	}
	return currentEnergy
}

func (g *game) getCurrentSpellCastId(objectType uint8, guid uint64) int32 {
	currentSpellCastId := int32(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		currentSpellCastId = int32(C.GetCurrentSpellCastId(C.uint64_t(guid)))
	}
	return currentSpellCastId
}

func (g *game) getLevel(objectType uint8, guid uint64) int32 {
	level := int32(0)
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		level = int32(C.GetLevel(C.uint64_t(guid)))
	}
	return level
}

func (g *game) getDynamicFlags(objectType uint8, guid uint64) []string {
	dynamicFlags := []string{}
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		flags := uint32(C.GetDynamicFlags(C.uint64_t(guid)))
		if HasFlag(flags, DynamicFlag_CanBeLooted) {
			dynamicFlags = append(dynamicFlags, "CAN_BE_LOOTED")
		}
		if HasFlag(flags, DynamicFlag_Tapped) {
			dynamicFlags = append(dynamicFlags, "TAPPED")
		}
		if HasFlag(flags, DynamicFlag_TappedByMe) {
			dynamicFlags = append(dynamicFlags, "TAPPED_BY_ME")
		}
		if HasFlag(flags, DynamicFlag_IsMarked) {
			dynamicFlags = append(dynamicFlags, "MARKED")
		}
	}
	return dynamicFlags
}

func (g *game) getReaction(objectType uint8, guid uint64) string {
	if objectType == UnitType_Unit || objectType == UnitType_Player {
		reactionId := int(C.GetUnitReaction(C.uint64_t(g.GetPlayerGuid()), C.uint64_t(guid)))
		switch reactionId {
		case UnitReaction_Exalted:
			return "EXALTED"
		case UnitReaction_Hated:
			return "HATED"
		case UnitReaction_Hostile:
			return "HOSTILE"
		case UnitReaction_Unfriendly:
			return "UNFRIENDLY"
		case UnitReaction_Neutral:
			return "NEUTRAL"
		case UnitReaction_Friendly:
			return "FRIENDLY"
		case UnitReaction_Honored:
			return "HONORED"
		case UnitReaction_Revered:
			return "REVERED"
		}
	}
	return ""
}

func (g *game) Sit() {
	g.RunLua("DoEmote('SIT')")
}

func (g *game) GetVisibleObjects() []WowObject {
	return g.objects
}

func (g *game) GetVisibleObjectByGuid(guid uint64) *WowObject {
	allObjects := g.GetVisibleObjects()
	index := slices.IndexFunc(allObjects, func(w WowObject) bool { return w.Guid == guid })
	return &allObjects[index]
}

func (g *game) MoveToPosition(position Vec3) {
	// logger.Log(fmt.Sprintf("Moving to position %v", position))

	C.ClickToMove(C.float(position.X), C.float(position.Y), C.float(position.Z))
}

func (g *game) StopMovement() {
	C.StopMovement()
}

func (g *game) SetFacing(angle float32) {
	playerGuid := g.GetPlayerGuid()
	playerPtr := uintptr(C.GetObjectPtr(C.uint64_t(playerGuid)))
	C.SetFacing(C.uint32_t(playerPtr), C.float(angle))
}

func (g *game) FaceTarget() {
	g.EnumerateVisibleObjects()
	playerGuid := g.GetPlayerGuid()
	player := g.GetVisibleObjectByGuid(playerGuid)
	targetGuid := player.TargetGuid
	if targetGuid == 0 {
		logger.Log("No target to face")
		return
	}
	target := g.GetVisibleObjectByGuid(player.TargetGuid)
	facing := player.Position.AngleTo(target.Position)
	g.SetFacing(facing)
}

func (g *game) FacePosition(pos Vec3) {
	g.EnumerateVisibleObjects()
	playerGuid := g.GetPlayerGuid()
	player := g.GetVisibleObjectByGuid(playerGuid)
	facing := player.Position.AngleTo(pos)
	g.SetFacing(facing)
}

func (g *game) SetTarget(guid uint64) {
	// logger.Log(fmt.Sprintf("Setting target to %v", guid))
	C.SetTarget(C.uint64_t(guid))
}

func (g *game) RunLua(lua string) {
	C.LuaCall(C.CString(lua))
	NotifyMainThread()
}

func (g *game) RunLuaWithResults(lua string) []string {

	luaVars := []string{}
	luaWithPlaceholders := lua

	for i := 0; i < 50; i++ {
		newLuaVar := randStringRunes(10)
		newLuaWithPlaceholders := strings.ReplaceAll(luaWithPlaceholders, "{"+fmt.Sprint(i)+"}", newLuaVar)
		if luaWithPlaceholders == newLuaWithPlaceholders {
			break
		}
		luaVars = append(luaVars, newLuaVar)
		luaWithPlaceholders = newLuaWithPlaceholders
	}

	g.RunLua(luaWithPlaceholders)

	results := []string{}

	for _, luaVar := range luaVars {
		results = append(results, C.GoString(C.GetText(C.CString(luaVar))))
	}

	return results
}

func (g *game) Jump() {
	g.RunLua("Jump()")
}

func (g *game) ToggleAttack(attack bool) {
	if attack {
		g.RunLua(attackOnLua)
	} else {
		g.RunLua(attackOffLua)
	}
}

func (g *game) CastSpellByName(spellName string) {
	g.RunLua("CastSpellByName('" + spellName + "')")
}

func (g *game) AutoLoot(guid uint64) {
	unitPtr := C.GetObjectPtr(C.uint64_t(guid))
	C.RightClick(unitPtr, C.int(1))
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//export EnumerateVisibleObjectsCallback
func EnumerateVisibleObjectsCallback(filter int, guid uint64) {
	globalGame.EnumerateVisibleObjectsCallback(filter, guid)
}

func HookEvents() {
	C.HookEvents()
}
