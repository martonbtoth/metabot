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
	"slices"
	"strings"
	"superbot/logger"
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
	SetTarget(guid uint64)
	RunLua(lua string)
	RunLuaWithResults(lua string) []string
	GetAvailableSpells() []string
	Jump()
	ToggleAttack(attack bool)
	CastSpellByName(spellName string)
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
	})
}

func (g *game) getName(objectType uint8, guid uint64) string {
	name := "N/A"
	if objectType == Unit {
		name = C.GoString(C.GetUnitName(C.uint64_t(guid)))
	} else if objectType == Player {
		name = C.GoString(C.GetPlayerName(C.uint64_t(guid)))
	}
	return name
}

func (g *game) getCurrentHealth(objectType uint8, guid uint64) int32 {
	currentHealth := 0
	if objectType == Unit {
		currentHealth = int(C.GetCurrentHealth(C.uint64_t(guid)))
	} else if objectType == Player {
		currentHealth = int(C.GetCurrentHealth(C.uint64_t(guid)))
	}
	return int32(currentHealth)
}

func (g *game) getMaxHealth(objectType uint8, guid uint64) int32 {
	currentHealth := 0
	if objectType == Unit {
		currentHealth = int(C.GetMaxHealth(C.uint64_t(guid)))
	} else if objectType == Player {
		currentHealth = int(C.GetMaxHealth(C.uint64_t(guid)))
	}
	return int32(currentHealth)
}

func (g *game) getTargetGuid(objectType uint8, guid uint64) uint64 {
	targetGuid := uint64(0)
	if objectType == Unit || objectType == Player {
		targetGuid = uint64(C.GetTargetGuid(C.uint64_t(guid)))
	}
	return targetGuid
}

func (g *game) getCurrentMana(objectType uint8, guid uint64) int32 {
	currentMana := int32(0)
	if objectType == Unit || objectType == Player {
		currentMana = int32(C.GetCurrentMana(C.uint64_t(guid)))
	}
	return currentMana
}

func (g *game) getMaxMana(objectType uint8, guid uint64) int32 {
	currentMana := int32(0)
	if objectType == Unit || objectType == Player {
		currentMana = int32(C.GetMaxMana(C.uint64_t(guid)))
	}
	return currentMana
}

func (g *game) getCurrentRage(objectType uint8, guid uint64) int32 {
	currentRage := int32(0)
	if objectType == Unit || objectType == Player {
		currentRage = int32(C.GetCurrentRage(C.uint64_t(guid)))
	}
	return currentRage
}

func (g *game) getCurrentEnergy(objectType uint8, guid uint64) int32 {
	currentEnergy := int32(0)
	if objectType == Unit || objectType == Player {
		currentEnergy = int32(C.GetCurrentEnergy(C.uint64_t(guid)))
	}
	return currentEnergy
}

func (g *game) getCurrentSpellCastId(objectType uint8, guid uint64) int32 {
	currentSpellCastId := int32(0)
	if objectType == Unit || objectType == Player {
		currentSpellCastId = int32(C.GetCurrentSpellCastId(C.uint64_t(guid)))
	}
	return currentSpellCastId
}

func (g *game) getLevel(objectType uint8, guid uint64) int32 {
	level := int32(0)
	if objectType == Unit || objectType == Player {
		level = int32(C.GetLevel(C.uint64_t(guid)))
	}
	return level
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
	logger.Log(fmt.Sprintf("Moving to position %v", position))

	C.ClickToMove(C.float(position.X), C.float(position.Y), C.float(position.Z))
}

func (g *game) StopMovement() {
	logger.Log("Stopping movement")

	C.StopMovement()
}

func (g *game) SetTarget(guid uint64) {
	logger.Log(fmt.Sprintf("Setting target to %v", guid))
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
