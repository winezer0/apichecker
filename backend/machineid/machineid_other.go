//go:build !windows

package machineid

import "errors"

// GetMachineID 获取当前设备的唯一机器码。
func GetMachineID() (string, error) {
	return "", errors.New("machine id is only implemented on windows")
}
