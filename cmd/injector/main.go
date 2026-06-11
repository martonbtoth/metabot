package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func setSeDebugPrivilege() {
	// Get current process (the one I wanna change)
	handle, err := windows.GetCurrentProcess()
	defer windows.CloseHandle(handle)
	if err != nil {
		log.Fatal(err)
	}

	// Get the current process token
	var token windows.Token
	err = windows.OpenProcessToken(handle, windows.TOKEN_ADJUST_PRIVILEGES, &token)
	if err != nil {
		log.Fatal(err)
	}

	// Check the LUID
	var luid windows.LUID
	seDebugName, err := windows.UTF16FromString("SeDebugPrivilege")
	if err != nil {
		fmt.Println(err)
	}
	err = windows.LookupPrivilegeValue(nil, &seDebugName[0], &luid)
	if err != nil {
		log.Fatal(err)
	}

	// Modify the token
	var tokenPriviledges windows.Tokenprivileges
	tokenPriviledges.PrivilegeCount = 1
	tokenPriviledges.Privileges[0].Luid = luid
	tokenPriviledges.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	// Adjust token privs
	tokPrivLen := uint32(unsafe.Sizeof(tokenPriviledges))
	fmt.Printf("Length is %d\n", tokPrivLen)
	err = windows.AdjustTokenPrivileges(token, false, &tokenPriviledges, tokPrivLen, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Debug Priviledge granted")
}

func main() {

	setSeDebugPrivilege()

	dPath := "C:\\superbot\\cmd\\dll\\superbot.dll" //Path to the DLL file to inject
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
