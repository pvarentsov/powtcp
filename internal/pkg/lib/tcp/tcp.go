package tcp

import "net"

// NewConnErrorChecker - create new connection error checker
func NewConnErrorChecker() *ConnErrorChecker {
	return &ConnErrorChecker{}
}

// ConnErrorChecker - tcp connection error checker
type ConnErrorChecker struct{}

// IsTimeout - define that tcp connection was timed out
func (ec *ConnErrorChecker) IsTimeout(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}
	return false
}
