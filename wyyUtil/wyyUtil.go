package wyyUtil

import (
	"bemfa-demo/kernel32Wrapper"
	"bemfa-demo/systemUtil"
	"bemfa-demo/user32Wrapper"
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

const (
	PROCESS_VM_READ                  = 0x0010
	PROCESS_QUERY_INFORMATION        = 0x0400
	VK_CONTRO                        = 0x11
	VK_ALT                           = 0x12
	VK_P                             = 0x50
	VK_SHIFT                         = 0x10
	VK_UP                            = 0x26
	VK_DOWN                          = 0x28
	VK_LEFT                          = 0x25
	VK_RIGHT                         = 0x27
	KEYEVENTF_DOWN                   = 0
	KEYEVENTF_UP                     = 2
	CloudmusicName                   = "cloudmusic.exe"
	dllName                          = "cloudmusic.dll"
	BASE_ADDRESS_OFFSET       uint64 = 0x19B8EF0
)

var (
	memoryOffset = []uint64{0, 0xD0, 0x5D0, 0x860}
)

func ChangePlayStatus() {
	log.Println("ChangePlayStatus")
	user32Wrapper.Keybd_event(VK_CONTRO, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_ALT, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_P, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)

	user32Wrapper.Keybd_event(VK_P, KEYEVENTF_UP)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_CONTRO, KEYEVENTF_UP)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_ALT, KEYEVENTF_UP)
}

func getBaseAddress(pids []uint32) uint64 {
	return uint64(kernel32Wrapper.GetAddressOfDll(pids[0], dllName)) + BASE_ADDRESS_OFFSET
}

func GetPlayStatus() bool {
	pids, _ := systemUtil.GetProcessIDsByName(CloudmusicName)

	if pids == nil || len(pids) == 0 {
		return false
	}
	var value = getBaseAddress(pids)
	var bytesRead uintptr
	var address uint64
	var err error
	var hProcess syscall.Handle
	for _, pid := range pids {
		hProcess, err = kernel32Wrapper.OpenProcess(PROCESS_VM_READ|PROCESS_QUERY_INFORMATION, false, pid)
		fmt.Println("pid:", pid)
		defer syscall.CloseHandle(hProcess)
		if err != nil {
			log.Println("Error opening process:", err)
			continue
		}

		for i, offset := range memoryOffset {
			address, err = getMemoryValue(&offset, &value, &hProcess, &bytesRead)
			if err == nil {
				fmt.Printf("Level: %X, Address: %X, Value: %X\n", i, address, value)
			} else {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return value&1 == 1
}

func ChangeVoice(plus bool) {
	user32Wrapper.Keybd_event(VK_CONTRO, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_ALT, KEYEVENTF_DOWN)
	if plus {
		log.Println("plus voice")
		user32Wrapper.Keybd_event(VK_UP, KEYEVENTF_DOWN)
		time.Sleep(50 * time.Millisecond)
		user32Wrapper.Keybd_event(VK_UP, KEYEVENTF_UP)
	} else {
		log.Println("minus voice")
		user32Wrapper.Keybd_event(VK_DOWN, KEYEVENTF_DOWN)
		time.Sleep(50 * time.Millisecond)
		user32Wrapper.Keybd_event(VK_DOWN, KEYEVENTF_UP)
	}
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_CONTRO, KEYEVENTF_UP)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_ALT, KEYEVENTF_UP)
}

func ChangeMusic(next bool) {
	user32Wrapper.Keybd_event(VK_CONTRO, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_SHIFT, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_ALT, KEYEVENTF_DOWN)
	time.Sleep(50 * time.Millisecond)
	if next {
		log.Println("next music")
		user32Wrapper.Keybd_event(VK_RIGHT, KEYEVENTF_DOWN)
		time.Sleep(50 * time.Millisecond)
		user32Wrapper.Keybd_event(VK_RIGHT, KEYEVENTF_UP)
	} else {
		log.Println("last music")
		user32Wrapper.Keybd_event(VK_LEFT, KEYEVENTF_DOWN)
		time.Sleep(50 * time.Millisecond)
		user32Wrapper.Keybd_event(VK_LEFT, KEYEVENTF_UP)
	}
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_CONTRO, KEYEVENTF_UP)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_SHIFT, KEYEVENTF_UP)
	time.Sleep(50 * time.Millisecond)
	user32Wrapper.Keybd_event(VK_ALT, KEYEVENTF_UP)
}

func getMemoryValue(offset *uint64, value *uint64, hProcess *syscall.Handle, bytesRead *uintptr) (uint64, error) {
	address := *value + *offset
	err := kernel32Wrapper.ReadProcessMemory(*hProcess, address, value, unsafe.Sizeof(value), bytesRead)
	return address, err
}
