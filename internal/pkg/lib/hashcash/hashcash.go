package hashcash

import (
	"crypto/rand"
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
)

// New - returns new hashcash
func New(bits int, resource string) (*Hashcash, error) {
	salt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return nil, err
	}

	return &Hashcash{
		bits:     bits,
		date:     time.Now().UTC().Truncate(time.Second),
		resource: resource,
		salt:     salt.Bytes(),
	}, nil
}

// ParseHeader - parses hashcah from header
func ParseHeader(header string) (hashcash *Hashcash, err error) {
	parts := strings.Split(header, ":")
	if len(parts) != 6 {
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

	hashcash.salt, err = base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	counterStr, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	hashcash.counter, err = strconv.Atoi(string(counterStr))
	if err != nil {
		return nil, ErrIncorrectHeaderFormat
	}

	return
}

// Hashcash - hashcash structure
type Hashcash struct {
	bits     int
	date     time.Time
	resource string
	salt     []byte
	counter  int
}

// Key - returns string presentation of hashcash without counter
// Using to match original hashcash with solved hashcash
func (h *Hashcash) Key() string {
	return fmt.Sprintf("%d:%d:%s:%d", h.bits, h.date.Unix(), h.resource, binary.BigEndian.Uint32(h.salt))
}

// Header - returns string presentation of hashcash to share it
func (h *Hashcash) Header() string {
	return fmt.Sprintf("1:%d:%s:%s:%s:%s",
		h.bits,
		h.date.Format(dateLayout),
		h.resource,
		base64.StdEncoding.EncodeToString(h.salt),
		base64.StdEncoding.Strict().EncodeToString([]byte(strconv.Itoa(h.counter))),
	)
}
