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
	salt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
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
		salt:     salt.Bytes(),
	}, nil
}

// Hashcash - hashcash structure
type Hashcash struct {
	bits      int
	date      time.Time
	resource  string
	extension string
	salt      []byte
	counter   int
}

// Bits - returns zero bits count
func (h *Hashcash) Bits() int {
	return h.bits
}

// EqualResource - chech if input resource is equal with hashcash resource
func (h *Hashcash) EqualResource(resource string) bool {
	return h.resource == resource
}

// Compute - compute hash with enough zero bits in the begining
// Increase counter while hash is not correct
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
// Using to match original hashcash with solved hashcash
func (h *Hashcash) Key() string {
	return fmt.Sprintf("%d:%d:%s:%d", h.bits, h.date.Unix(), h.resource, binary.BigEndian.Uint32(h.salt))
}

// Header - returns string presentation of hashcash to share it
func (h *Hashcash) Header() Header {
	return Header(fmt.Sprintf("1:%d:%s:%s:%s:%s:%s",
		h.bits,
		h.date.Format(dateLayout),
		h.resource,
		h.extension,
		base64.StdEncoding.EncodeToString(h.salt),
		base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.counter))),
	))
}

// ParseHeader - parses hashcah from header
func ParseHeader(header string) (hashcash *Hashcash, err error) {
	parts := strings.Split(header, ":")
	if len(parts) != 7 {
		return nil, ErrIncorrectHeaderFormat
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

	hashcash.salt, err = base64.StdEncoding.DecodeString(parts[5])
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
// Format - bits:date:resource:externsion:salt:counter
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
