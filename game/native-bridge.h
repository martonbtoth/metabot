#include <stdint.h>

uint64_t GetPlayerGuid();
__stdcall int EnumerateVisibleObjects(int filter);
int GetObjectPtr(uint64_t guid);
int GetCurrentHealth(uint64_t guid);
int GetMaxHealth(uint64_t guid);
char* GetUnitName(uint64_t guid);
char* GetPlayerName(uint64_t guid);

float GetObjectPositionX(uint64_t guid);
float GetObjectPositionY(uint64_t guid);
float GetObjectPositionZ(uint64_t guid);
void ClickToMove(float x, float y, float z);

void LuaCall(char* code);
