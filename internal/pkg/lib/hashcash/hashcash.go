package hashcash

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout = "20060102150405"
	zeroBit    = '0'
)

// New - returns new hashcash
func New(bits int, resource string) (*Hashcash, error) {
	rand, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return nil, err
	}

	if bits <= 0 {
		return nil, ErrZeroBitsMustBeMoreThanZero
	}

	return &Hashcash{
		bits:     bits,
		date:     time.Now().UTC().Truncate(time.Second),
		resource: resource,
		rand:     rand.Bytes(),
	}, nil
}

// Hashcash - hashcash structure
// Version 1
type Hashcash struct {
	bits      int       // number of zero bits in hashed code
	date      time.Time // time that the message was sent
	resource  string    // resource data string (IP address,  email address, etc)
	extension string    // extension, ignored in this version
	rand      []byte    // random characters
	counter   int       // computing counter
}

// Bits - returns number of zero bits
func (h *Hashcash) Bits() int {
	return h.bits
}

// Counter - returns counter
func (h *Hashcash) Counter() int {
	return h.counter
}

// EqualResource - check if input resource is equal with hashcash resource
func (h *Hashcash) EqualResource(resource string) bool {
	return h.resource == resource
}

// IsActual - check if hashcash expiration exceeded ttl
func (h *Hashcash) IsActual(ttl time.Duration) bool {
	return h.date.Add(ttl).After(time.Now().UTC())
}

// Compute - compute hash with enough zero bits in the begining
// Increase counter if hash does't have enough zero bits in the begining
func (h *Hashcash) Compute(maxAttempts int) error {
	if maxAttempts > 0 {
		h.counter = 0
		for h.counter <= maxAttempts {
			ok, err := h.Header().IsHashCorrect(h.bits)
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
			h.counter++
		}
	}

	return ErrComputingMaxAttemptsExceeded
}

// Key - returns string presentation of hashcash without counter
// Key is using to match original hashcash with solved hashcash
func (h *Hashcash) Key() string {
	return fmt.Sprintf("%d:%d:%s:%d", h.bits, h.date.Unix(), h.resource, binary.BigEndian.Uint32(h.rand))
}

// Header - returns string presentation of hashcash to share it
func (h *Hashcash) Header() Header {
	return Header(fmt.Sprintf("1:%d:%s:%s:%s:%s:%s",
		h.bits,
		h.date.Format(dateLayout),
		h.resource,
		h.extension,
		base64.StdEncoding.EncodeToString(h.rand),
		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))),
	))
}

// ParseHeader - parse hashcah from header
func ParseHeader(header string) (hashcash *Hashcash, err error) {
	parts := strings.Split(header, ":")

	if len(parts) < 7 {
		return nil, ErrIncorrectHeaderFormat
	}
	if len(parts) > 7 {
		for i := 0; i < len(parts)-7; i++ {
			parts[3] += ":" + parts[3+i+1]
		}
		parts[4] = parts[len(parts)-3]
		parts[5] = parts[len(parts)-2]
		parts[6] = parts[len(parts)-1]
		parts = parts[:7]
	}
	if parts[0] != "1" {
		return nil, ErrIncorrectHeaderFormat
	}

	hashcash = &Hashcash{}

	hashcash.bits, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	hashcash.date, err = time.ParseInLocation(dateLayout, parts[2], time.UTC)
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	hashcash.resource = parts[3]
	hashcash.extension = parts[4]

	hashcash.rand, err = base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	counterStr, err := base64.StdEncoding.DecodeString(parts[6])
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	hashcash.counter, err = strconv.Atoi(string(counterStr))
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	return
}

// Header - string presentation of hashcash
// Format - 1:bits:date:resource:externsion:rand:counter
type Header string

// IsHashCorrect - does header hash constain zero bits enough
func (header Header) IsHashCorrect(bits int) (ok bool, err error) {
	if bits <= 0 {
		return false, ErrZeroBitsMustBeMoreThanZero
	}

	hash, err := header.sha1()
	if err != nil {
		return ok, err
	}
	if len(hash) < bits {
		return false, ErrHashLengthLessThanZeroBits
	}

	ok = true
	for _, s := range hash[:bits] {
		if s != zeroBit {
			ok = false
			break
		}
	}
	return
}

func (header Header) sha1() (hash string, err error) {
	hasher := sha1.New()
	if _, err = hasher.Write([]byte(header)); err != nil {
		return
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
