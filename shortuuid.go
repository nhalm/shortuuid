// Package shortuuid provides functionality to encode strings and UUIDs into short, URL-safe identifiers using base62 encoding.
//
// The package offers two main approaches:
//   - Shorten/Expand: For encoding arbitrary strings
//   - ShortenUUID/ExpandUUID: For encoding UUID objects with optimized hex handling
//
// All functions use a base62 alphabet (0-9, A-Z, a-z) to create compact, readable identifiers.
package shortuuid

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// EncodeError represents an error that occurs during string or UUID encoding.
// It contains the original input and a description of what went wrong.
type EncodeError struct {
	Input  string // The input that failed to encode
	Reason string // Description of the error
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("encode error for input '%s': %s", e.Input, e.Reason)
}

// DecodeError represents an error that occurs during short ID decoding.
// It contains the short ID that failed to decode and a description of the error.
type DecodeError struct {
	ShortID string // The short ID that failed to decode
	Reason  string // Description of the error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode error for short ID '%s': %s", e.ShortID, e.Reason)
}

// defaultBase62 is the default alphabet for encoding
// Uses: '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
var defaultBase62 = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Shorten converts any string to a short, URL-safe identifier using base62 encoding.
// The input string is converted to bytes and then encoded using the base62 alphabet.
// Returns an error if the input string is empty.
func Shorten(input string) (string, error) {
	return encodeString(input)
}

// Expand converts a short ID back to the original string using base62 decoding.
// The short ID must contain only valid base62 characters (0-9, A-Z, a-z).
// Returns an error if the short ID contains invalid characters.
func Expand(shortID string) (string, error) {
	return decodeString(shortID)
}

// ShortenUUID converts a uuid.UUID to a short, URL-safe identifier using optimized hex encoding.
// This method is more efficient than Shorten for UUID objects as it works directly with
// the UUID's hex representation rather than converting to string first.
func ShortenUUID(u uuid.UUID) (string, error) {
	uuidStr := u.String()
	if uuidStr == "" {
		return "", &EncodeError{
			Input:  uuidStr,
			Reason: "UUID string cannot be empty",
		}
	}

	// Remove dashes and encode
	cleanUUID := strings.ReplaceAll(uuidStr, "-", "")
	return encodeHex(cleanUUID)
}

// ExpandUUID converts a short ID back to a uuid.UUID object.
// The short ID must have been created by ShortenUUID to ensure proper UUID format.
// Returns an error if the short ID is invalid or doesn't decode to a valid UUID.
func ExpandUUID(shortID string) (uuid.UUID, error) {
	// Decode the short ID to hex string
	hexStr, err := decodeHex(shortID)
	if err != nil {
		return uuid.UUID{}, err
	}

	// Add dashes to create proper UUID format
	if len(hexStr) != 32 {
		return uuid.UUID{}, &DecodeError{
			ShortID: shortID,
			Reason:  fmt.Sprintf("decoded to invalid length: expected 32 hex characters, got %d", len(hexStr)),
		}
	}

	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s",
		hexStr[0:8],
		hexStr[8:12],
		hexStr[12:16],
		hexStr[16:20],
		hexStr[20:32])

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, &DecodeError{
			ShortID: shortID,
			Reason:  "failed to parse UUID: " + err.Error(),
		}
	}

	return parsedUUID, nil
}

// encodeString converts any string to a short ID using base62
func encodeString(input string) (string, error) {
	if input == "" {
		return "", &EncodeError{
			Input:  input,
			Reason: "input string cannot be empty",
		}
	}

	// Convert string to bytes, then to big integer
	bytes := []byte(input)
	num := new(big.Int)
	num.SetBytes(bytes)

	// Convert to base62
	return intToBase(num), nil
}

// decodeString converts a short ID back to the original string
func decodeString(shortID string) (string, error) {
	// Convert from base to integer
	num, err := baseToInt(shortID)
	if err != nil {
		return "", err
	}

	// Convert big integer back to bytes, then to string
	bytes := num.Bytes()
	return string(bytes), nil
}

// encodeHex converts a hex string to a short ID using base62
func encodeHex(hexStr string) (string, error) {
	// Convert hex string to big integer
	num := new(big.Int)
	num.SetString(hexStr, 16)

	// Convert to base62
	return intToBase(num), nil
}

// decodeHex converts a short ID back to a hex string
func decodeHex(shortID string) (string, error) {
	// Convert from base to integer
	num, err := baseToInt(shortID)
	if err != nil {
		return "", err
	}

	// Convert to hex string with proper padding for UUID (32 chars)
	hexStr := fmt.Sprintf("%032s", num.Text(16))
	return hexStr, nil
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
				Reason:  fmt.Sprintf("invalid character '%c' in short ID (valid characters: 0-9, A-Z, a-z)", char),
			}
		}

		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(index)))
	}

	return result, nil
}
