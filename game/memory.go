package game

import (
	"golang.org/x/sys/windows"
)

func WriteMemory(addr uintptr, data []byte) {
	numberOfBytesWritten := uintptr(0)
	windows.WriteProcessMemory(windows.CurrentProcess(), addr, &data[0], uintptr(len(data)), &numberOfBytesWritten)
}

func UnlockProtectedLuaFunctions() {
	patch := []byte{0xB8, 0x01, 0x00, 0x00, 0x00, 0xc3}
	WriteMemory(0x494A50, patch)
}

func FixClickToMove() {
	patch := []byte{0x00, 0x00, 0x00, 0x00}
	WriteMemory(0x860A90, patch)
}
