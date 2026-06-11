#ifndef NATIVEFUNCTIONS_H
#define NATIVEFUNCTIONS

#include <stdint.h>

typedef void __cdecl (*FrameScript_RegisterFunction_t)(char* functionName, int (__cdecl* function)(void*)); // last param is LuaState*
FrameScript_RegisterFunction_t FrameScript_RegisterFunction = 0x00704120;

typedef uint64_t (*GetPlayerGuidNative_t)(void);
GetPlayerGuidNative_t GetPlayerGuidNative = 0x00468550;

typedef __stdcall uint32_t (*GetObjectPtrNative_t)(uint64_t);
GetObjectPtrNative_t GetObjectPtrNative = 0x00464870;

typedef void __fastcall (*EnumerateVisibleObjectsNative_t)(__stdcall int (*callback)(uint64_t), int filter);
EnumerateVisibleObjectsNative_t EnumerateVisibleObjectsNative = 0x00468380;

typedef void __thiscall (*ClickToMoveNative_t)(uintptr_t, uint32_t, unsigned long long*, float*, float);
ClickToMoveNative_t ClickToMoveNative = 0x00611130;

typedef void __fastcall (*LuaCallNative_t)(char* code, const char* unused);
LuaCallNative_t LuaCallNative = 0x00704CD0;

typedef void __stdcall (*SetTargetNative_t)(uint64_t guid);
SetTargetNative_t SetTargetNative = 0x00493540;

typedef char* __fastcall (*GetTextNative_t)(char* varName, unsigned int nonSense, int zero);
GetTextNative_t GetTextNative = (char *(__fastcall *)(char *varName, unsigned int nonSense, int zero)) 0x00703BF0;

typedef void __fastcall (*LootSlotNative_t)(uint32_t slot, int32_t unused);
LootSlotNative_t LootSlotNative = (void (__fastcall *)(uint32_t slot, int32_t unused)) 0x004C2790;

typedef void __thiscall (*RightClickNative_t)(uint32_t unitPtr, int32_t autoLoot);
RightClickNative_t RightClickNative = (void (*)(uint32_t unitPtr, int32_t autoLoot)) 0x60BEA0;

typedef void __thiscall (*SetFacingNative_t)(uint32_t playerSetFacingPtr, float angle);
SetFacingNative_t SetFacingNative = (void (*)(uint32_t playerSetFacingPtr, float angle)) 0x007C6F30;

typedef void __thiscall (*SendMovementUpdateNative_t)(uint32_t playerPtr, uint32_t unknown, uint32_t opcode, uint32_t unknown2, uint32_t unknown3);
SendMovementUpdateNative_t SendMovementUpdateNative = (void (*)(uint32_t playerPtr, uint32_t unknown, uint32_t opcode, uint32_t unknown2, uint32_t unknown3)) 0x00600A30;

typedef uint32_t __thiscall (*GetUnitReactionNative_t)(int unitPtr, int towardsUnitPointer);
GetUnitReactionNative_t GetUnitReactionNative = (uint32_t (*)(int unitPtr, int towardsUnitPointer)) 0x006061E0;

#endif