package wire_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/monetas/bmd/wire"
)

// TestShaHash tests the ShaHash API.
func TestShaHash(t *testing.T) {

	shaHashStr := "14a0810ac680a3eb3f82edc878cea25ec41d6b790744e5daeef"
	shaHash, err := wire.NewShaHashFromStr(shaHashStr)
	if err != nil {
		t.Errorf("NewShaHashFromStr: %v", err)
	}

	buf := []byte{
		0x79, 0xa6, 0x1a, 0xdb, 0xc6, 0xe5, 0xa2, 0xe1,
		0x39, 0xd2, 0x71, 0x3a, 0x54, 0x6e, 0xc7, 0xc8,
		0x75, 0x63, 0x2e, 0x75, 0xf1, 0xdf, 0x9c, 0x3f,
		0xa6, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	hash, err := wire.NewShaHash(buf)
	if err != nil {
		t.Errorf("NewShaHash: unexpected error %v", err)
	}

	// Ensure proper size.
	if len(hash) != wire.HashSize {
		t.Errorf("NewShaHash: hash length mismatch - got: %v, want: %v",
			len(hash), wire.HashSize)
	}

	// Ensure contents match.
	if !bytes.Equal(hash[:], buf) {
		t.Errorf("NewShaHash: hash contents mismatch - got: %v, want: %v",
			hash[:], buf)
	}

	if hash.IsEqual(shaHash) {
		t.Errorf("IsEqual: hash contents should not match - got: %v, want: %v",
			hash, shaHash)
	}

	// Set hash from byte slice and ensure contents match.
	err = hash.SetBytes(shaHash.Bytes())
	if err != nil {
		t.Errorf("SetBytes: %v", err)
	}
	if !hash.IsEqual(shaHash) {
		t.Errorf("IsEqual: hash contents mismatch - got: %v, want: %v",
			hash, shaHash)
	}

	// Invalid size for SetBytes.
	err = hash.SetBytes([]byte{0x00})
	if err == nil {
		t.Errorf("SetBytes: failed to received expected err - got: nil")
	}

	// Invalid size for NewShaHash.
	invalidHash := make([]byte, wire.HashSize+1)
	_, err = wire.NewShaHash(invalidHash)
	if err == nil {
		t.Errorf("NewShaHash: failed to received expected err - got: nil")
	}
}

// TestShaHashString  tests the stringized output for sha hashes.
func TestShaHashString(t *testing.T) {
	wantStr := "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"
	hash := wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
		0x06, 0xe5, 0x33, 0xfd, 0x1a, 0xda, 0x86, 0x39,
		0x1f, 0x3f, 0x6c, 0x34, 0x32, 0x04, 0xb0, 0xd2,
		0x78, 0xd4, 0xaa, 0xec, 0x1c, 0x0b, 0x20, 0xaa,
		0x27, 0xba, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	hashStr := hash.String()
	if hashStr != wantStr {
		t.Errorf("String: wrong hash string - got %v, want %v",
			hashStr, wantStr)
	}
}

// TestNewShaHashFromStr executes tests against the NewShaHashFromStr function.
func TestNewShaHashFromStr(t *testing.T) {
	tests := []struct {
		in   string
		want wire.ShaHash
		err  error
	}{
		// Empty string.
		{
			"",
			wire.ShaHash{},
			nil,
		},

		// Single digit hash.
		{
			"1",
			wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}),
			nil,
		},

		{
			"3264bc2ac36a60840790ba1d475d01367e7c723da941069e9dc",
			wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
				0xdc, 0xe9, 0x69, 0x10, 0x94, 0xda, 0x23, 0xc7,
				0xe7, 0x67, 0x13, 0xd0, 0x75, 0xd4, 0xa1, 0x0b,
				0x79, 0x40, 0x08, 0xa6, 0x36, 0xac, 0xc2, 0x4b,
				0x26, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}),
			nil,
		},

		// Hash string that is too long.
		{
			"01234567890123456789012345678901234567890123456789012345678912345",
			wire.ShaHash{},
			wire.ErrHashStrSize,
		},

		// Hash string that is contains non-hex chars.
		{
			"abcdefg",
			wire.ShaHash{},
			hex.InvalidByteError('g'),
		},
	}

	unexpectedErrStr := "NewShaHashFromStr #%d failed to detect expected error - got: %v want: %v"
	unexpectedResultStr := "NewShaHashFromStr #%d got: %v want: %v"
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result, err := wire.NewShaHashFromStr(test.in)
		if err != test.err {
			t.Errorf(unexpectedErrStr, i, err, test.err)
			continue
		} else if err != nil {
			// Got expected error. Move on to the next test.
			continue
		}
		if !test.want.IsEqual(result) {
			t.Errorf(unexpectedResultStr, i, result, &test.want)
			continue
		}
	}
}
