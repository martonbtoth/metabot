#ifndef _NATIVE_BRIDGE_H
#define _NATIVE_BRIDGE_H

#include "native-bridge.h"
#include "threadhelper.h"
#include <stdint.h>
#include <string.h>
#include <windows.h>

#endif

const int GET_PLAYER_GUID_FUN_PTR = 0x00468550;
const int ENUMERATE_VISIBLE_OBJECTS_FUN_PTR = 0x00468380;
const int GET_OBJECT_PTR_FUN_PTR = 0x00464870;
const int CLICK_TO_MOVE_FUN_PTR = 0x00611130;
const int CLICK_TO_MOVE_FIX_PTR = 0x860A90;
const int LUA_CALL_FUN_PTR = 0x00704CD0;
const int GET_TEXT_FUN_PTR = 0x00703BF0;
const uint32_t UNIT_NAME_OFFSET = 0xB30;
const int* NAME_LIST_POINTER = 0xC0E230;
const int NAME_PTR_GUID_OFFSET = 0xC;
const int NAME_PTR_NAME_OFFSET = 0x14;
const int DESCRIPTOR_OFFSET = 0x8;
const int CURRENT_HEALTH_OFFSET = 0x58;
const int MAX_HEALTH_OFFSET = 0x70;
const int POS_X_OFFSET = 0x9B8;
const int POS_Y_OFFSET = 0x9BC;
const int POS_Z_OFFSET = 0x9C0;

extern int EnumerateVisibleObjectsCallback(int filter, uint64_t guid);
extern int WndProcGoCallback(int* hWnd, int Msg, int wParam, int lParam);

typedef struct {
    float X;
    float Y;
    float Z;
} vec3;

uint64_t GetPlayerGuid() {
    uint64_t (*GetPlayerGuidNative)(void) = (uint64_t (*)(void))GET_PLAYER_GUID_FUN_PTR;
    return GetPlayerGuidNative();
}

int GetObjectPtr(uint64_t guid) {
    __stdcall int* (*GetObjectPtrNative)(uint64_t) = (__stdcall int* (*)(uint64_t))GET_OBJECT_PTR_FUN_PTR;
    return (int) GetObjectPtrNative(guid);
}

char* GetUnitName(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    int** ptr1 = (int**)(objectptr + UNIT_NAME_OFFSET);
    int* ptr2 = *ptr1;
    return (char*) *ptr2;
}

int GetCurrentHealth(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    int descriptor = *(int*)(objectptr + DESCRIPTOR_OFFSET);
    int value = *(int*)(descriptor + CURRENT_HEALTH_OFFSET);
    return value;
}

int GetMaxHealth(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    int descriptor = *(int*)(objectptr + DESCRIPTOR_OFFSET);
    int value = *(int*)(descriptor + MAX_HEALTH_OFFSET);
    return value;
}

char* GetPlayerName(uint64_t guid) {
    int namePtr = *NAME_LIST_POINTER;
    for (uint32_t i = 0; i < 1000000; i++) {
    // while (1) {
        uint64_t nextGuid = *(uint64_t*)(namePtr + NAME_PTR_GUID_OFFSET);
        if (nextGuid == guid) {
            return namePtr + NAME_PTR_NAME_OFFSET;
        } else {
            namePtr = *(int*)namePtr;
        }
    }

    return "Player name not found!";
}


__stdcall int EnumerateVisibleObjectsCallbackInternal(uint64_t guid) {
    EnumerateVisibleObjectsCallback(0, guid);
    return 1;
}

__stdcall int EnumerateVisibleObjects(int filter) {
        // typedef void __fastcall func(__thiscall int (*callback)(void*, int, uint64_t), int filter);
        typedef void __fastcall func(__stdcall int (*callback)(uint64_t), int filter);
        func* enumerateFunction = (func*)ENUMERATE_VISIBLE_OBJECTS_FUN_PTR;
        enumerateFunction(EnumerateVisibleObjectsCallbackInternal, filter);
}

float GetObjectPositionX(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    return *(float*)(objectptr + POS_X_OFFSET);
}

float GetObjectPositionY(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    return *(float*)(objectptr + POS_Y_OFFSET);
}

float GetObjectPositionZ(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    return *(float*)(objectptr + POS_Z_OFFSET);
}

uint64_t interactGuid = 0;

void ClickToMove(float x, float y, float z) {
    float destination[3] = {0.0, 0.0, 0.0};
    // void __thiscall (*ClickToMoveNative)(uintptr_t, int, uint64_t*, vec3*, float) = (void __thiscall(*)(uintptr_t, int, uint64_t, vec3*, float))CLICK_TO_MOVE_FUN_PTR;
    typedef void (__thiscall* func)(uintptr_t, uint32_t, unsigned long long*, float*, float);
    func ClickToMoveNative = (func)CLICK_TO_MOVE_FUN_PTR;
    uint64_t playerGuid = GetPlayerGuid();
    int playerPtr = GetObjectPtr(playerGuid);
    destination[0] = x;
    destination[1] = y;
    destination[2] = z;
    ClickToMoveNative(playerPtr, 0x4, &interactGuid, destination, 2);
}

void LuaCall(char* code) {
    void luaCalInternal() {
        typedef void __fastcall func(char* code, const char* unused);
        func* f = (func*)LUA_CALL_FUN_PTR;
        f(code, "Superbot");
        return;
    }
    RunOnMainThread(luaCalInternal);
}

char* GetText(char* varName) {
    typedef char* __fastcall func(char* varName, unsigned int nonSense, int zero);
    func* f = (func*)GET_TEXT_FUN_PTR;
    return f(varName, 0xFFFFFFFF, 0);
}
