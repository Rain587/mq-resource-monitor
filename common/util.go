package common

import (
	"fmt"
	"net"
)

func CheckPortInUse(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	defer listener.Close()
	return false
}

func FindNextFreePort(port int) int {
	for {
		if !CheckPortInUse(port) {
			return port
		} else {
			port++
		}
	}
}
