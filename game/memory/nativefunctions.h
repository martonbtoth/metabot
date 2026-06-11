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

typedef __stdcall (*SetTargetNative_t)(uint64_t guid);
SetTargetNative_t SetTargetNative = 0x00493540;

typedef char* __fastcall (*GetTextNative_t)(char* varName, unsigned int nonSense, int zero);
GetTextNative_t GetTextNative = 0x00703BF0;

#endif