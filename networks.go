package derr

import (
	"net"
	"os"
	"syscall"
)

// IsClientSideNetworkError returns whether the error received is a network error caused by the client side
// that could not be possibly handled correctly on the server side anyway.
func IsClientSideNetworkError(err error) bool {
	netErr, isNetErr := err.(*net.OpError)
	if !isNetErr {
		return false
	}

	syscallErr, isSyscallErr := netErr.Err.(*os.SyscallError)
	if !isSyscallErr {
		return false
	}

	return syscallErr.Err == syscall.ECONNRESET || syscallErr.Err == syscall.EPIPE
}
