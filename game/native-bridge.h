#include <stdint.h>

uint64_t GetPlayerGuid();
__stdcall int EnumerateVisibleObjects(int filter);
int GetObjectPtr(uint64_t guid);
int GetCurrentHealth(uint64_t guid);
int GetMaxHealth(uint64_t guid);
char* GetUnitName(uint64_t guid);
char* GetPlayerName(uint64_t guid);
uint64_t GetTargetGuid(uint64_t guid);

int32_t GetLevel(uint64_t guid);
int32_t GetCurrentMana(uint64_t guid);
int32_t GetMaxMana(uint64_t guid);
int32_t GetCurrentRage(uint64_t guid);
int32_t GetCurrentEnergy(uint64_t guid);
int32_t GetCurrentSpellCastId(uint64_t guid);
uint32_t GetDynamicFlags(uint64_t guid);
uint32_t GetUnitReaction(uint64_t guid, uint64_t towardsUnitGuid);

float GetObjectPositionX(uint64_t guid);
float GetObjectPositionY(uint64_t guid);
float GetObjectPositionZ(uint64_t guid);
void ClickToMove(float x, float y, float z);
void StopMovement();
void SetFacing(uint32_t playerPtr, float angle);
void SetTarget(uint64_t guid);
void RightClick(int32_t unitPtr, int32_t autoLoot);

void LuaCall(char* code);
char* GetText(char* varName);

void HookEvents();
