package user32Wrapper

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	enumWindows              = user32.NewProc("EnumWindows")
	getWindowText            = user32.NewProc("GetWindowTextW")
	getClassName             = user32.NewProc("GetClassNameW")
	showWindow               = user32.NewProc("ShowWindow")
	isIconic                 = user32.NewProc("IsIconic")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procIsWindowVisible      = user32.NewProc("IsWindowVisible")
	findWindowW              = user32.NewProc("FindWindowW")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	keybd_event              = user32.NewProc("keybd_event")
)

func FindWindowW(className string, windowName string) syscall.Handle {
	var err error

	var classNameUTF16 *uint16
	if className == "" {
		classNameUTF16 = nil
	} else {
		classNameUTF16, err = syscall.UTF16PtrFromString(className)
		if err != nil {
			log.Fatal(err)
		}
	}

	var windowNameUTF16 *uint16
	if windowName == "" {
		windowNameUTF16 = nil
	} else {
		windowNameUTF16, err = syscall.UTF16PtrFromString(windowName)
		if err != nil {
			log.Fatal(err)
		}
	}

	r1, _, _ := findWindowW.Call(uintptr(unsafe.Pointer(classNameUTF16)), uintptr(unsafe.Pointer(windowNameUTF16)))
	return syscall.Handle(r1)
}

func GetWindowThreadProcessId(hwnd syscall.Handle) (uint32, uint32) {
	var pid uint32
	tid, _, _ := getWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return uint32(tid), pid
}

func IsWindowVisible(hwnd syscall.Handle) bool {
	ret, _, _ := procIsWindowVisible.Call(
		uintptr(hwnd),
	)
	return ret != 0
}

func SetForegroundWindow(hwnd syscall.Handle) {
	ret, _, err := procSetForegroundWindow.Call(
		uintptr(hwnd),
	)
	if ret == 0 {
		if err != nil {
			fmt.Println("Error calling SetForegroundWindow:", err)
		} else {
			fmt.Println("SetForegroundWindow failed")
		}
	}
}

func GetWindowClassName(hwnd syscall.Handle) (string, error) {
	const nMaxCount = 256
	className := make([]uint16, nMaxCount)
	_, _, err := getClassName.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&className[0])), uintptr(nMaxCount))
	if err != nil && err.Error() != "The operation completed successfully." {
		return "", err
	}
	return syscall.UTF16ToString(className), nil
}

func GetWindowText(hwnd syscall.Handle) string {
	buffer := make([]uint16, 256)
	getWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(256))
	return syscall.UTF16ToString(buffer)
}

func ShowWindow(hwnd syscall.Handle, nCmdShow int) {
	showWindow.Call(
		uintptr(hwnd),
		uintptr(nCmdShow))
}

func IsIconic(hwnd syscall.Handle) bool {
	ret, _, _ := isIconic.Call(uintptr(hwnd))
	return ret != 0
}

func EnumWindows(callback func(hwnd syscall.Handle, lParam uintptr) uintptr, lParam int) {
	enumWindows.Call(syscall.NewCallback(callback), uintptr(lParam))
}

func Keybd_event(bVk int, dwFlags int) {
	keybd_event.Call(uintptr(bVk), 0, uintptr(dwFlags), 0)
}
