package systemUtil

import (
	"bemfa-demo/kernel32Wrapper"
	"fmt"
	"syscall"
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
)

func GetProcessIDsByName(processName string) ([]uint32, error) {
	var pids []uint32

	snapshot, _, err := kernel32Wrapper.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if snapshot == uintptr(syscall.InvalidHandle) {
		return nil, err
	}
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	var entry kernel32Wrapper.Processentry32

	ret, _, _ := kernel32Wrapper.Process32First(snapshot, &entry)
	if ret == 0 {
		return nil, fmt.Errorf("failed to retrieve first process")
	}

	for {
		exeFile := syscall.UTF16ToString(entry.SzExeFile[:])
		if exeFile == processName {
			pids = append(pids, entry.Th32ProcessID)
		}
		ret, _, _ := kernel32Wrapper.Process32NextW(snapshot, &entry)
		if ret == 0 {
			break
		}
	}

	if len(pids) == 0 {
		return nil, fmt.Errorf("process %s not found", processName)
	}
	fmt.Println("returned pids:", pids)
	return pids, nil
}
