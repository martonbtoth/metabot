package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

func main() {

	dPath := "..\\dll\\superbot.dll" //Path to the DLL file to inject
	gamePath := "C:\\wow\\WoW.exe"
	////

	processInfo := syscall.ProcessInformation{}
	startupInfo := syscall.StartupInfo{}

	gamePathUtf16, _ := syscall.UTF16PtrFromString(gamePath)
	syscall.CreateProcess(
		gamePathUtf16,
		nil,
		nil,
		nil,
		false,
		0x08000000,
		nil,
		nil,
		&startupInfo,
		&processInfo,
	)

	time.Sleep(500 * time.Millisecond)

	kernel32 := windows.NewLazyDLL("kernel32.dll")

	//Opens a handle to the target process with the needed permissions
	pHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, processInfo.ProcessId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Process opened")
	////

	//Allocates virtual memory for the file path
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	vAlloc, _, err := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(dPath)+1), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
	fmt.Println("Memory allocated")
	////

	//Converts the file path to type *byte
	bPtrDpath, err := windows.BytePtrFromString(dPath)
	if err != nil {
		log.Fatal(err)
	}
	////

	//Writes the filename to the previously allocated space
	Zero := uintptr(0)
	err = windows.WriteProcessMemory(pHandle, vAlloc, bPtrDpath, uintptr(len(dPath)+1), &Zero)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DLL path written")
	////

	//Gets a pointer to the LoadLibrary function
	LoadLibAddr, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
	if err != nil {
		log.Fatal(err)
	}
	////

	//Creates a remote thread that loads the DLL triggering it
	tHandle, _, _ := kernel32.NewProc("CreateRemoteThread").Call(uintptr(pHandle), 0, 0, LoadLibAddr, vAlloc, 0, 0)
	defer syscall.CloseHandle(syscall.Handle(tHandle))
	fmt.Println("DLL Injected")
	////

}
