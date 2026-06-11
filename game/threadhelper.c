#include "threadhelper.h"
#include <stdint.h>
#include <string.h>
#include <windows.h>

extern void NotifyMainThread();

int oldCallback = 0;
char* luaCallCode = 0;

typedef struct ActionList_s {
    void (*action)();
    HANDLE semaphore;
    struct ActionList_s* next;
} ActionList;

ActionList* actions = 0;

void RunOnMainThread(void (*action)()) {
    ActionList* newAction = (ActionList*)malloc(sizeof(ActionList));
    newAction->action = action;
    newAction->next = actions;
    newAction->semaphore = CreateSemaphore(NULL, 0, 1, NULL);
    actions = newAction;
    NotifyMainThread();
    WaitForSingleObject(newAction->semaphore, 10000L);
}

ActionList* PopAction() {
    if (actions == 0) {
        return 0;
    }
    ActionList* action = actions;
    actions = actions->next;
    return action;
}

int WndProcCallback(int* hWnd, int Msg, int wParam, int lParam) {
    if (actions != 0) {
        ActionList* action = PopAction();
        action->action();
        ReleaseSemaphore(action->semaphore, 1, NULL);
        free(action);
    }
    return CallWindowProc(oldCallback, hWnd, Msg, wParam, lParam);
}

void SetOldCallback(int callback) {
    oldCallback = callback;
}

int GetWndProcCallbackPtr() {
    return WndProcCallback;
}
