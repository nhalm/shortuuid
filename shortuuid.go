package shortuuid

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

// EncodeError represents an error that occurs during UUID shortening
type EncodeError struct {
	UUID   string
	Reason string
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("encode error for UUID '%s': %s", e.UUID, e.Reason)
}

// DecodeError represents an error that occurs during short ID expansion
type DecodeError struct {
	ShortID string
	Reason  string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode error for short ID '%s': %s", e.ShortID, e.Reason)
}

// defaultBase62 is the default alphabet for encoding
// Uses: '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
var defaultBase62 = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Shorten converts any string to a short ID using the default alphabet
func Shorten(input string) (string, error) {
	return encode(input)
}

// Expand converts a short ID back to the original string using the default alphabet
func Expand(shortID string) (string, error) {
	return decode(shortID)
}

// ShortenUUID converts a uuid.UUID to a short ID using the default alphabet
func ShortenUUID(u uuid.UUID) (string, error) {
	return Shorten(u.String())
}

// ExpandUUID converts a short ID back to a uuid.UUID using the default alphabet
func ExpandUUID(shortID string) (uuid.UUID, error) {
	uuidStr, err := Expand(shortID)
	if err != nil {
		return uuid.UUID{}, err
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, &DecodeError{
			ShortID: shortID,
			Reason:  "failed to parse expanded UUID: " + err.Error(),
		}
	}

	return parsedUUID, nil
}

// encode converts any string to a short ID by treating it as bytes
func encode(input string) (string, error) {
	if input == "" {
		return "", &EncodeError{
			UUID:   input,
			Reason: "input string cannot be empty",
		}
	}

	// Convert string to bytes, then to big integer
	bytes := []byte(input)
	num := new(big.Int)
	num.SetBytes(bytes)

	// Convert to the target base
	return intToBase(num), nil
}

// decode converts a short ID back to the original string
func decode(shortID string) (string, error) {
	// Convert from base to integer
	num, err := baseToInt(shortID)
	if err != nil {
		return "", err
	}

	// Convert big integer back to bytes, then to string
	bytes := num.Bytes()
	return string(bytes), nil
}

// intToBase converts a big integer to the target base representation
func intToBase(num *big.Int) string {
	if num.Sign() == 0 {
		return string(defaultBase62[0])
	}

	var result []rune
	base := big.NewInt(int64(len(defaultBase62)))
	zero := big.NewInt(0)

	// Make a copy to avoid modifying the original
	n := new(big.Int).Set(num)

	for n.Cmp(zero) > 0 {
		remainder := new(big.Int)
		n.DivMod(n, base, remainder)
		result = append([]rune{defaultBase62[remainder.Int64()]}, result...)
	}

	return string(result)
}

// baseToInt converts a base representation back to a big integer
func baseToInt(encoded string) (*big.Int, error) {
	result := big.NewInt(0)
	base := big.NewInt(int64(len(defaultBase62)))

	for _, char := range encoded {
		// Find the character in the alphabet
		index := -1
		for i, alphabetChar := range defaultBase62 {
			if char == alphabetChar {
				index = i
				break
			}
		}

		if index == -1 {
			return nil, &DecodeError{
				ShortID: encoded,
				Reason:  fmt.Sprintf("contains invalid character '%c' not found in alphabet", char),
			}
		}

		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(index)))
	}

	return result, nil
}
