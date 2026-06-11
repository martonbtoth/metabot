#ifndef _NATIVE_BRIDGE_H
#define _NATIVE_BRIDGE_H

#include "native-bridge.h"
#include <stdint.h>
#include <string.h>
#include <windows.h>

#endif

void RunOnMainThread(void (*action)());
int WndProcCallback(int* hWnd, int Msg, int wParam, int lParam);
void SetOldCallback(int callback);
int GetWndProcCallbackPtr();