package game

/*
#include "native-bridge.h"
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

//go:embed enumerate_spellbook.lua
var enumerateSpellbookLua string

type Game interface {
	GetPlayerGuid() uint64
	EnumerateVisibleObjects()
	EnumerateVisibleObjectsCallback(filter int, guid uint64)
	GetVisibleObjects() []WowObject
	GetVisibleObjectByGuid(guid uint64) *WowObject
	MoveToPosition(position Vec3)
	RunLua(lua string)
	RunLuaWithResults(lua string) []string
	GetAvailableSpells() []string
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
	C.EnumerateVisibleObjects(0)
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
		Guid:          guid,
		Pointer:       objectUintPtr,
		Type:          objectType,
		Name:          g.getName(objectType, guid),
		MaxHealth:     g.getMaxHealth(objectType, guid),
		CurrentHealth: g.getCurrentHealth(objectType, guid),
		Position:      position,
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

func (g *game) getCurrentHealth(objectType uint8, guid uint64) int {
	currentHealth := 0
	if objectType == Unit {
		currentHealth = int(C.GetCurrentHealth(C.uint64_t(guid)))
	} else if objectType == Player {
		currentHealth = int(C.GetCurrentHealth(C.uint64_t(guid)))
	}
	return currentHealth
}

func (g *game) getMaxHealth(objectType uint8, guid uint64) int {
	currentHealth := 0
	if objectType == Unit {
		currentHealth = int(C.GetMaxHealth(C.uint64_t(guid)))
	} else if objectType == Player {
		currentHealth = int(C.GetMaxHealth(C.uint64_t(guid)))
	}
	return currentHealth
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
	logger.GetLogger().Log(fmt.Sprintf("Moving to position %v", position))

	C.ClickToMove(C.float(position.X), C.float(position.Y), C.float(position.Z))
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
