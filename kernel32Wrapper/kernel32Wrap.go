package kernel32Wrapper

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	openProcess              = kernel32.NewProc("OpenProcess")
	readProcessMemory        = kernel32.NewProc("ReadProcessMemory")
	createToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	process32FirstW          = kernel32.NewProc("Process32FirstW")
	process32NextW           = kernel32.NewProc("Process32NextW")
	module32First            = kernel32.NewProc("Module32First")
	module32Next             = kernel32.NewProc("Module32Next")
)

type Processentry32 struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [260]uint16
}

type MODULEENTRY32 struct {
	DwSize        uint32
	Th32ModuleID  uint32
	Th32ProcessID uint32
	GlblcntUsage  uint32
	ProccntUsage  uint32
	ModBaseAddr   uintptr
	ModBaseSize   uint32
	HModule       uintptr
	szModule      [256]byte
	szExePath     [260]byte
}

const (
	TH32CS_SNAPMODULE = 0x00000008
)

func GetAddressOfDll(pid uint32, dllName string) uintptr {
	moduleSnapshot, _, _ := createToolhelp32Snapshot.Call(uintptr(TH32CS_SNAPMODULE), uintptr(pid))
	defer syscall.CloseHandle(syscall.Handle(moduleSnapshot))

	// 初始化 MODULEENTRY32 结构体
	var moduleEntry MODULEENTRY32
	moduleEntry.DwSize = uint32(unsafe.Sizeof(moduleEntry))

	// 遍历目标进程的模块列表，输出 DLL 的地址
	ret, _, _ := module32First.Call(moduleSnapshot, uintptr(unsafe.Pointer(&moduleEntry)))
	for ret != 0 {
		if strings.HasPrefix(string(moduleEntry.szModule[:]), dllName) {
			fmt.Printf("Base Address: %X\n", moduleEntry.ModBaseAddr)
			break
		}
		ret, _, _ = module32Next.Call(moduleSnapshot, uintptr(unsafe.Pointer(&moduleEntry)))
	}
	return moduleEntry.ModBaseAddr
}

func OpenProcess(desiredAccess uint32, inheritHandle bool, processId uint32) (handle syscall.Handle, err error) {
	inherit := 0
	if inheritHandle {
		inherit = 1
	}
	r0, _, e1 := syscall.Syscall(openProcess.Addr(), 3, uintptr(desiredAccess), uintptr(inherit), uintptr(processId))
	handle = syscall.Handle(r0)
	if handle == syscall.InvalidHandle {
		err = e1
	}
	return
}

func ReadProcessMemory(hProcess syscall.Handle, lpBaseAddress uint64, lpBuffer *uint64, nSize uintptr, lpNumberOfBytesRead *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(readProcessMemory.Addr(), 5, uintptr(hProcess), uintptr(lpBaseAddress), uintptr(unsafe.Pointer(lpBuffer)), nSize, uintptr(unsafe.Pointer(lpNumberOfBytesRead)), 0)
	if r1 == 0 {
		err = e1
	}
	return
}

func CreateToolhelp32Snapshot(dwFlags int, th32ProcessID int) (uintptr, uintptr, error) {
	r1, r2, err := createToolhelp32Snapshot.Call(uintptr(dwFlags), uintptr(th32ProcessID))
	return r1, r2, err
}

func Process32First(hSnapshot uintptr, entry *Processentry32) (uintptr, uintptr, error) {
	(*entry).DwSize = uint32(unsafe.Sizeof(*entry))
	r1, r2, err := process32FirstW.Call(hSnapshot, uintptr(unsafe.Pointer(entry)))
	return r1, r2, err
}

func Process32NextW(hSnapshot uintptr, entry *Processentry32) (uintptr, uintptr, error) {
	r1, r2, err := process32NextW.Call(hSnapshot, uintptr(unsafe.Pointer(entry)))
	return r1, r2, err
}
