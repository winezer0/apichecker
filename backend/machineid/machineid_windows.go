//go:build windows

package machineid

import (
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const machineGuidPath = `SOFTWARE\Microsoft\Cryptography`

// GetMachineID 获取当前 Windows 设备的唯一机器码。
func GetMachineID() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, machineGuidPath, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}

	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", errors.New("machine id is empty")
	}

	return normalized, nil
}
