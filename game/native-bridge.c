#include "native-bridge.h"
#include "threadhelper.h"
#include "memory/ctmtypes.h"
#include "memory/addresses.h"
#include "memory/nativefunctions.h"
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <windows.h>

#include "minhook/MinHook.h"

extern int EnumerateVisibleObjectsCallback(int filter, uint64_t guid);
extern int WndProcGoCallback(int* hWnd, int Msg, int wParam, int lParam);

typedef struct {
    float X;
    float Y;
    float Z;
} vec3;

uint64_t GetPlayerGuid() {
    return GetPlayerGuidNative();
}

int GetObjectPtr(uint64_t guid) {
    return (int) GetObjectPtrNative(guid);
}

char* GetUnitName(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    int** ptr1 = (int**)(objectptr + UNIT_NAME_OFFSET);
    int* ptr2 = *ptr1;
    return (char*) *ptr2;
}

int* getDescriptorPtr(uint64_t guid) {
    int objectptr = GetObjectPtr(guid);
    return *(int*)(objectptr + DESCRIPTOR_OFFSET);
}

int32_t getIntFromDescriptorOffset(uint64_t guid, int offset) {
    int descriptor = getDescriptorPtr(guid);
    int value = *(int32_t*)(descriptor + offset);
    return value;
}

uint64_t getGuidFromDescriptorOffset(uint64_t guid, int offset) {
    int descriptor = getDescriptorPtr(guid);
    int value = *(uint64_t*)(descriptor + offset);
    return value;
}

int GetCurrentHealth(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, CURRENT_HEALTH_OFFSET);
}

int GetMaxHealth(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, MAX_HEALTH_OFFSET);
}

uint64_t GetTargetGuid(uint64_t guid) {
    return getGuidFromDescriptorOffset(guid, TARGET_GUID_OFFSET);
}

int32_t GetCurrentMana(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, CURRENT_MANA_OFFSET);
}

int32_t GetMaxMana(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, MAX_MANA_OFFSET);
}

int32_t GetCurrentRage(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, RAGE_OFFSET);
}

int32_t GetCurrentEnergy(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, ENERGY_OFFSET);
}

int32_t GetLevel(uint64_t guid) {
    return getIntFromDescriptorOffset(guid, LEVEL_OFFSET);
}

int32_t GetCurrentSpellCastId(uint64_t guid) {
    int ptr = GetObjectPtr(guid);
    int value = *(int*)(ptr + CURRENT_SPELLCAST_OFFSET);
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
        EnumerateVisibleObjectsNative(EnumerateVisibleObjectsCallbackInternal, filter);
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

void clickToMoveInternal(float x, float y, float z, uint32_t ctmType) {
    float destination[3] = {0.0, 0.0, 0.0};
    uint64_t playerGuid = GetPlayerGuid();
    int playerPtr = GetObjectPtr(playerGuid);
    destination[0] = x;
    destination[1] = y;
    destination[2] = z;
    ClickToMoveNative(playerPtr, ctmType, &interactGuid, destination, 2);
}

void ClickToMove(float x, float y, float z) {
    clickToMoveInternal(x, y, z, CTM_TYPE_MOVE);
}

void StopMovement() {
    clickToMoveInternal(0.f, 0.f, 0.f, CTM_TYPE_NONE);
}

void LuaCall(char* code) {
    void luaCalInternal() { // ignore this error: https://github.com/microsoft/vscode-cpptools/issues/1035
        LuaCallNative(code, "Superbot");
        return;
    }
    RunOnMainThread(luaCalInternal);
}

char* GetText(char* varName) {
    return GetTextNative(varName, 0xFFFFFFFF, 0);
}

void SetTarget(uint64_t guid) {
    void setTargetInternal() { // ignore this error: https://github.com/microsoft/vscode-cpptools/issues/1035
        SetTargetNative(guid);
        return;
    }
    RunOnMainThread(setTargetInternal);
}

__fastcall void (*SignalEventOriginal)(uint32_t) = 0;

void Log(char* log) {
    FILE* file = fopen("C:\\superbot\\native.log", "ab");
    fprintf(file, "%s\n", log);
    fclose(file);
}

__fastcall void SignalEventHook(uint32_t event) {
    char buffer[400];
    sprintf(buffer, "Event: %x", event);
    Log(buffer);
    SignalEventOriginal(event);
}

void HookEvents() {
    Log("Hooking events...");
    if (MH_Initialize() != MH_OK) {
        Log("Could not initialize MinHook");
        return;
    }
    if (MH_CreateHook(SIGNAL_EVENT_FUN_PTR, SignalEventHook, &SignalEventOriginal) != MH_OK) {
        Log("Could not create hook for SignalEvent");
        return;
    }
    if (MH_EnableHook(SIGNAL_EVENT_FUN_PTR) != MH_OK) {
        Log("Could not enable hook for SignalEvent");
        return;
    }
    Log("Hooking events...");
}
