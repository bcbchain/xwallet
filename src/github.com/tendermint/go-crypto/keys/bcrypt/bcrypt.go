package bcrypt

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/blowfish"
)

const (
	MinCost		int	= 4
	MaxCost		int	= 31
	DefaultCost	int	= 10
)

var ErrMismatchedHashAndPassword = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")

var ErrHashTooShort = errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password")

type HashVersionTooNewError byte

func (hv HashVersionTooNewError) Error() string {
	return fmt.Sprintf("crypto/bcrypt: bcrypt algorithm version '%c' requested is newer than current version '%c'", byte(hv), majorVersion)
}

type InvalidHashPrefixError byte

func (ih InvalidHashPrefixError) Error() string {
	return fmt.Sprintf("crypto/bcrypt: bcrypt hashes must start with '$', but hashedSecret started with '%c'", byte(ih))
}

type InvalidCostError int

func (ic InvalidCostError) Error() string {
	return fmt.Sprintf("crypto/bcrypt: cost %d is outside allowed range (%d,%d)", int(ic), int(MinCost), int(MaxCost))
}

const (
	majorVersion		= '2'
	minorVersion		= 'a'
	maxSaltSize		= 16
	maxCryptedHashSize	= 23
	encodedSaltSize		= 22
	encodedHashSize		= 31
	minHashSize		= 59
)

var magicCipherData = []byte{
	0x4f, 0x72, 0x70, 0x68,
	0x65, 0x61, 0x6e, 0x42,
	0x65, 0x68, 0x6f, 0x6c,
	0x64, 0x65, 0x72, 0x53,
	0x63, 0x72, 0x79, 0x44,
	0x6f, 0x75, 0x62, 0x74,
}

type hashed struct {
	hash	[]byte
	salt	[]byte
	cost	int
	major	byte
	minor	byte
}

func GenerateFromPassword(salt []byte, password []byte, cost int) ([]byte, error) {
	if len(salt) != maxSaltSize {
		return nil, fmt.Errorf("Salt len must be %v", maxSaltSize)
	}
	p, err := newFromPassword(salt, password, cost)
	if err != nil {
		return nil, err
	}
	return p.Hash(), nil
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	p, err := newFromHash(hashedPassword)
	if err != nil {
		return err
	}

	otherHash, err := bcrypt(password, p.cost, p.salt)
	if err != nil {
		return err
	}

	otherP := &hashed{otherHash, p.salt, p.cost, p.major, p.minor}
	if subtle.ConstantTimeCompare(p.Hash(), otherP.Hash()) == 1 {
		return nil
	}

	return ErrMismatchedHashAndPassword
}

func Cost(hashedPassword []byte) (int, error) {
	p, err := newFromHash(hashedPassword)
	if err != nil {
		return 0, err
	}
	return p.cost, nil
}

func newFromPassword(salt []byte, password []byte, cost int) (*hashed, error) {
	if cost < MinCost {
		cost = DefaultCost
	}
	p := new(hashed)
	p.major = majorVersion
	p.minor = minorVersion

	err := checkCost(cost)
	if err != nil {
		return nil, err
	}
	p.cost = cost

	p.salt = base64Encode(salt)
	hash, err := bcrypt(password, p.cost, p.salt)
	if err != nil {
		return nil, err
	}
	p.hash = hash
	return p, err
}

func newFromHash(hashedSecret []byte) (*hashed, error) {
	if len(hashedSecret) < minHashSize {
		return nil, ErrHashTooShort
	}
	p := new(hashed)
	n, err := p.decodeVersion(hashedSecret)
	if err != nil {
		return nil, err
	}
	hashedSecret = hashedSecret[n:]
	n, err = p.decodeCost(hashedSecret)
	if err != nil {
		return nil, err
	}
	hashedSecret = hashedSecret[n:]

	p.salt = make([]byte, encodedSaltSize, encodedSaltSize+2)
	copy(p.salt, hashedSecret[:encodedSaltSize])

	hashedSecret = hashedSecret[encodedSaltSize:]
	p.hash = make([]byte, len(hashedSecret))
	copy(p.hash, hashedSecret)

	return p, nil
}

func bcrypt(password []byte, cost int, salt []byte) ([]byte, error) {
	cipherData := make([]byte, len(magicCipherData))
	copy(cipherData, magicCipherData)

	c, err := expensiveBlowfishSetup(password, uint32(cost), salt)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 24; i += 8 {
		for j := 0; j < 64; j++ {
			c.Encrypt(cipherData[i:i+8], cipherData[i:i+8])
		}
	}

	hsh := base64Encode(cipherData[:maxCryptedHashSize])
	return hsh, nil
}

func expensiveBlowfishSetup(key []byte, cost uint32, salt []byte) (*blowfish.Cipher, error) {

	csalt, err := base64Decode(salt)
	if err != nil {
		return nil, err
	}

	ckey := append(key, 0)

	c, err := blowfish.NewSaltedCipher(ckey, csalt)
	if err != nil {
		return nil, err
	}

	var i, rounds uint64
	rounds = 1 << cost
	for i = 0; i < rounds; i++ {
		blowfish.ExpandKey(ckey, c)
		blowfish.ExpandKey(csalt, c)
	}

	return c, nil
}

func (p *hashed) Hash() []byte {
	arr := make([]byte, 60)
	arr[0] = '$'
	arr[1] = p.major
	n := 2
	if p.minor != 0 {
		arr[2] = p.minor
		n = 3
	}
	arr[n] = '$'
	n += 1
	copy(arr[n:], []byte(fmt.Sprintf("%02d", p.cost)))
	n += 2
	arr[n] = '$'
	n += 1
	copy(arr[n:], p.salt)
	n += encodedSaltSize
	copy(arr[n:], p.hash)
	n += encodedHashSize
	return arr[:n]
}

func (p *hashed) decodeVersion(sbytes []byte) (int, error) {
	if sbytes[0] != '$' {
		return -1, InvalidHashPrefixError(sbytes[0])
	}
	if sbytes[1] > majorVersion {
		return -1, HashVersionTooNewError(sbytes[1])
	}
	p.major = sbytes[1]
	n := 3
	if sbytes[2] != '$' {
		p.minor = sbytes[2]
		n++
	}
	return n, nil
}

func (p *hashed) decodeCost(sbytes []byte) (int, error) {
	cost, err := strconv.Atoi(string(sbytes[0:2]))
	if err != nil {
		return -1, err
	}
	err = checkCost(cost)
	if err != nil {
		return -1, err
	}
	p.cost = cost
	return 3, nil
}

func (p *hashed) String() string {
	return fmt.Sprintf("&{hash: %#v, salt: %#v, cost: %d, major: %c, minor: %c}", string(p.hash), p.salt, p.cost, p.major, p.minor)
}

func checkCost(cost int) error {
	if cost < MinCost || cost > MaxCost {
		return InvalidCostError(cost)
	}
	return nil
}
