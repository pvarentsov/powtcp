package hashcash

import "errors"

// Errors
var (
	ErrIncorrectHeaderFormat        = errors.New("incorrect header format")
	ErrHashLengthLessThanZeroBits   = errors.New("hash length cannot be less than zero bits")
	ErrZeroBitsMustBeMoreThanZero   = errors.New("zero bits must be more than zero")
	ErrComputingMaxAttemptsExceeded = errors.New("max attempts to compute correct hash exceeded")
)
