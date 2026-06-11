package game

/*
#include "native-bridge.h"
*/
import "C"

import (
	"fmt"
	"superbot/logger"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
	GWLP_WNDPROC       = int32(-4)
)

var oldCallback uintptr
var mainWindow syscall.Handle

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.SyscallN(procEnumWindows.Addr(), uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			logger.Log(fmt.Sprintf("EnumWindows failed: %v", e1))
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func NotifyMainThread() {
	if oldCallback == 0 {
		mainWindow, _ = FindWindow("World of Warcraft")
		// wndProcCallbackPtr := C.GetWndProcCallbackPtr()
		// cb := syscall.NewCallback(func(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
		// 	go func() {
		// 		logger.Log(fmt.Sprintf("oldCallback: %v, WndProc: %v %v %v %v", oldCallback, hwnd, msg, wParam, lParam))
		// 	}()
		// 	return win.CallWindowProc(oldCallback, hwnd, msg, wParam, lParam)
		// })
		oldCallback = win.SetWindowLongPtr(win.HWND(mainWindow), win.GWLP_WNDPROC, uintptr(C.GetWndProcCallbackPtr()))
		C.SetOldCallback(C.int(oldCallback))
	}
	win.SendMessage(win.HWND(mainWindow), 0x0400, 0, 0)
}
