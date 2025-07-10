package shortuuid

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestCompatibility(t *testing.T) {
	// Test cases to ensure compatibility with expected encodings
	testCases := map[string]string{
		"53a8d1b9-4eca-4888-9b59-8fa91497857b": "2XrVqpuNYMfp5OSuawGnL1",
		"8658bb57-992d-4a4d-9292-a5b118d28c8b": "45VWNy74cXYBydTM0JO3rv",
		"d26abc73-a6bf-49c6-984d-e08c941fad4a": "6P3AMn3h7r9JJSeHECCjJS",
		"62263f37-9ef1-49ba-9790-d29edcb20f4a": "2zCj4Ne1tBbZaiZlHG7XOE",
		"9dff78bd-25bf-48a7-b962-47121e2f8abb": "4o8XraygEz74U5VcWAHfb9",
		"bd7da66f-0473-4f0b-b22b-13e50e091188": "5lYyGdkeeGhhagAO34mMO0",
		"eaf186de-b255-42bf-8046-288e286e5799": "79Ka7SqpMzbYobZdnAA1Cz",
		"4f77f8d9-1c48-46d7-b477-693c2168a63a": "2PxDiMQ6GcdX2gWH9uxQGg",
		"97ed0539-f38e-4851-9ea1-2c420b9e3c8e": "4cg9v4IOwgEgAjQXVl2aiM",
		"2e5b029c-bae4-4b4b-8265-699d911d2c42": "1PTEo1hMaUFibCywjsRLIg",
		"607a4a5e-28cb-48e4-892a-8f09d52c1a0c": "2w39JluOCmuC3CnE34ReI8",
		"5b109229-b3d8-4623-bf69-5d71ed611db4": "2lpsELdUhx0S8wiem4mMq4",
		"788e6389-ff1c-4b43-86b2-0ae92292c62c": "3fU9MsSs7HFv9LWe6BX50m",
		"1d152c86-c436-4d47-9269-006c7469b867": "ssS9A1oUhTFAbdjd6w93P",
		"e371489d-8cf9-4612-8bcb-ba6e7bcb6da2": "6vB1nTe4C1jJ3R7e5M0ETa",
		"e9360b9a-9b0c-442e-9f41-70a1cb2e92b1": "763uTm4g8wSb9ojuPwYwSH",
		"08f057f3-23e0-4b2a-8703-03f2dab8f628": "Grm6f7QVJVuVrufEOTgIC",
		"df540be5-e504-4c72-9ef9-17ecbfd4886d": "6nPhI89QMDrD1iQr7xEX9R",
		"54027bef-34b4-433e-92e8-bd250718d819": "2YWUQIU5tlKQqT5qbsw7ZJ",
		"a7dcc285-16ec-4b0a-af1a-3658a15ed6d1": "56kbZzRtV2g7yYf4G97VUv",
	}

	for originalUUID, expectedShort := range testCases {
		t.Run(originalUUID, func(t *testing.T) {
			// Test shortening
			actualShort, err := Shorten(originalUUID)
			if err != nil {
				t.Fatalf("Error shortening UUID %s: %v", originalUUID, err)
			}

			if actualShort != expectedShort {
				t.Errorf("Expected short ID %s, got %s", expectedShort, actualShort)
			}

			// Test expanding
			expandedUUID, err := Expand(actualShort)
			if err != nil {
				t.Fatalf("Error expanding short ID %s: %v", actualShort, err)
			}

			if expandedUUID != originalUUID {
				t.Errorf("Expected expanded UUID %s, got %s", originalUUID, expandedUUID)
			}

			t.Logf("✓ UUID: %s <-> Short: %s", originalUUID, actualShort)
		})
	}
}

func TestShorten(t *testing.T) {
	// Test basic shortening and expanding
	uuid := "4890586e-32a5-4f9c-a000-2a2bb68eb1ce"

	short, err := Shorten(uuid)
	if err != nil {
		t.Fatalf("Error shortening UUID: %v", err)
	}

	expanded, err := Expand(short)
	if err != nil {
		t.Fatalf("Error expanding short ID: %v", err)
	}

	if expanded != uuid {
		t.Errorf("Expected %s, got %s", uuid, expanded)
	}

	t.Logf("UUID: %s -> Short: %s -> UUID: %s", uuid, short, expanded)
}

func TestShortenUUID(t *testing.T) {
	// Test ShortenUUID and ExpandUUID with uuid.UUID types
	testUUID := uuid.New()

	short, err := ShortenUUID(testUUID)
	if err != nil {
		t.Fatalf("Error shortening UUID: %v", err)
	}

	expanded, err := ExpandUUID(short)
	if err != nil {
		t.Fatalf("Error expanding short ID: %v", err)
	}

	if expanded != testUUID {
		t.Errorf("Expected %s, got %s", testUUID, expanded)
	}

	t.Logf("UUID: %s -> Short: %s -> UUID: %s", testUUID, short, expanded)
}

func TestUUIDVersionPreservation(t *testing.T) {
	// Test that UUID versions are preserved
	testCases := []struct {
		name string
		uuid uuid.UUID
	}{
		{
			name: "UUIDv4",
			uuid: uuid.New(), // v4
		},
	}

	// Add UUIDv7 if available
	if uuidv7, err := uuid.NewV7(); err == nil {
		testCases = append(testCases, struct {
			name string
			uuid uuid.UUID
		}{
			name: "UUIDv7",
			uuid: uuidv7,
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalVersion := tc.uuid.Version()

			short, err := ShortenUUID(tc.uuid)
			if err != nil {
				t.Fatalf("Error shortening %s: %v", tc.name, err)
			}

			expanded, err := ExpandUUID(short)
			if err != nil {
				t.Fatalf("Error expanding %s: %v", tc.name, err)
			}

			if expanded.Version() != originalVersion {
				t.Errorf("Expected version %d, got %d", originalVersion, expanded.Version())
			}

			if expanded != tc.uuid {
				t.Errorf("Expected %s, got %s", tc.uuid, expanded)
			}

			t.Logf("%s: %s -> %s -> %s (version %d)", tc.name, tc.uuid, short, expanded, expanded.Version())
		})
	}
}

func TestBasicEncodeDecode(t *testing.T) {
	testCases := []string{
		"00000000-0000-0000-0000-000000000000",
		"ffffffff-ffff-ffff-ffff-ffffffffffff",
		"12345678-1234-5678-9abc-123456789abc",
		"550e8400-e29b-41d4-a716-446655440000",
	}

	for _, uuid := range testCases {
		t.Run(uuid, func(t *testing.T) {
			short, err := Shorten(uuid)
			if err != nil {
				t.Fatalf("Error shortening UUID %s: %v", uuid, err)
			}

			expanded, err := Expand(short)
			if err != nil {
				t.Fatalf("Error expanding short ID %s: %v", short, err)
			}

			if expanded != uuid {
				t.Errorf("Expected %s, got %s", uuid, expanded)
			}

			t.Logf("UUID: %s -> Short: %s", uuid, short)
		})
	}
}

func TestRandomUUIDs(t *testing.T) {
	// Test with random UUIDs
	for i := 0; i < 100; i++ {
		uuid := uuid.New().String()

		short, err := Shorten(uuid)
		if err != nil {
			t.Fatalf("Error shortening UUID %s: %v", uuid, err)
		}

		expanded, err := Expand(short)
		if err != nil {
			t.Fatalf("Error expanding short ID %s: %v", short, err)
		}

		if expanded != uuid {
			t.Errorf("Expected %s, got %s", uuid, expanded)
		}
	}
}

func TestErrorCases(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		isShortID bool
	}{
		{"invalid_not-a-uuid", "not-a-uuid", false},
		{"invalid_12345", "12345", false},
		{"invalid_12345678-1234-5678-9abc-123456789abcdef", "12345678-1234-5678-9abc-123456789abcdef", false},
		{"invalid_gggggggg-gggg-gggg-gggg-gggggggggggg", "gggggggg-gggg-gggg-gggg-gggggggggggg", false},
		{"invalid_short_@#$%", "@#$%", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isShortID {
				// Test invalid short ID
				_, err := Expand(tc.input)
				if err == nil {
					t.Errorf("Expected error for invalid short ID %s", tc.input)
				}
			} else {
				// Test invalid UUID
				_, err := Shorten(tc.input)
				if err == nil {
					t.Errorf("Expected error for invalid UUID %s", tc.input)
				}
			}
		})
	}
}

func TestError(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		isShortID      bool
		expectedInput  string
		expectedReason string
	}{
		{
			name:           "too short UUID",
			input:          "12345",
			isShortID:      false,
			expectedInput:  "12345",
			expectedReason: "invalid UUID format: expected 32 hex characters after removing hyphens, got 5",
		},
		{
			name:           "too long UUID",
			input:          "12345678-1234-5678-9abc-123456789abcdef",
			isShortID:      false,
			expectedInput:  "12345678-1234-5678-9abc-123456789abcdef",
			expectedReason: "invalid UUID format: expected 32 hex characters after removing hyphens, got 35",
		},
		{
			name:           "invalid hex UUID",
			input:          "gggggggg-gggg-gggg-gggg-gggggggggggg",
			isShortID:      false,
			expectedInput:  "gggggggg-gggg-gggg-gggg-gggggggggggg",
			expectedReason: "invalid UUID format: contains non-hex characters (valid characters: 0-9, a-f, A-F, hyphens)",
		},
		{
			name:           "invalid character in short ID",
			input:          "@#$%",
			isShortID:      true,
			expectedInput:  "@#$%",
			expectedReason: "invalid character '@' in short ID (valid characters: 0-9, A-Z, a-z)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			if tc.isShortID {
				_, err = Expand(tc.input)
			} else {
				_, err = Shorten(tc.input)
			}

			if err == nil {
				t.Fatalf("Expected error for input %s, but got none", tc.input)
			}

			if tc.isShortID {
				// Test DecodeError
				var decodeErr *DecodeError
				if !errors.As(err, &decodeErr) {
					t.Fatalf("Expected DecodeError, got %T: %v", err, err)
				}

				if decodeErr.ShortID != tc.expectedInput {
					t.Errorf("Expected ShortID %q in error, got %q", tc.expectedInput, decodeErr.ShortID)
				}

				if decodeErr.Reason != tc.expectedReason {
					t.Errorf("Expected reason %q in error, got %q", tc.expectedReason, decodeErr.Reason)
				}

				expectedMsg := fmt.Sprintf("decode error for short ID '%s': %s", tc.expectedInput, tc.expectedReason)
				if decodeErr.Error() != expectedMsg {
					t.Errorf("Expected error message %q, got %q", expectedMsg, decodeErr.Error())
				}
			} else {
				// Test EncodeError
				var encodeErr *EncodeError
				if !errors.As(err, &encodeErr) {
					t.Fatalf("Expected EncodeError, got %T: %v", err, err)
				}

				if encodeErr.UUID != tc.expectedInput {
					t.Errorf("Expected UUID %q in error, got %q", tc.expectedInput, encodeErr.UUID)
				}

				if encodeErr.Reason != tc.expectedReason {
					t.Errorf("Expected reason %q in error, got %q", tc.expectedReason, encodeErr.Reason)
				}

				expectedMsg := fmt.Sprintf("encode error for UUID '%s': %s", tc.expectedInput, tc.expectedReason)
				if encodeErr.Error() != expectedMsg {
					t.Errorf("Expected error message %q, got %q", expectedMsg, encodeErr.Error())
				}
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test that we can use errors.As with our error types
	_, err := Shorten("invalid")
	if err == nil {
		t.Fatal("Expected error")
	}

	// Test errors.As with EncodeError
	var encodeErr *EncodeError
	if !errors.As(err, &encodeErr) {
		t.Error("errors.As should work with EncodeError")
	}

	// Test with DecodeError
	_, err = Expand("@#$%")
	if err == nil {
		t.Fatal("Expected error")
	}

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Error("errors.As should work with DecodeError")
	}
}

// Benchmark tests
func BenchmarkShorten(b *testing.B) {
	uuid := "4890586e-32a5-4f9c-a000-2a2bb68eb1ce"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Shorten(uuid)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExpand(b *testing.B) {
	shortID := "2CvPdpytrcURpSLoPxYb30"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Expand(shortID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkShortenUUID(b *testing.B) {
	testUUID := uuid.New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ShortenUUID(testUUID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExpandUUID(b *testing.B) {
	shortID := "2CvPdpytrcURpSLoPxYb30"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExpandUUID(shortID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
