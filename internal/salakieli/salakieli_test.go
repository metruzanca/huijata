package salakieli

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var testPlain = []byte("<?xml version=\"1.0\"?>\n<GameStats TEST=\"1\"></GameStats>\n")

func TestDecryptStatsPassphrase(t *testing.T) {
	// Encrypting testPlain with the _stats key/IV produces testCipherStats.
	const testCipherStats = "d029674ece3cbddd124639c318831ddcc288673af03e89127a8877f47d60a5dd6718f3021b70d87277848900015320210a072b41a9deaf"
	ct, err := hex.DecodeString(testCipherStats)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decrypt("_stats", ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, testPlain) {
		t.Errorf("Decrypt(_stats) = %q, want %q", got, testPlain)
	}
}

func TestDecryptSessionNumbersPassphrase(t *testing.T) {
	// session_numbers.salakieli uses the magic_numbers passphrase.
	const testCipherSession = "d815b82fdd4a53cc983421814045748a8eba3b82bda1371b00efd47ee6860b88ad41a985b5f974923ce216d7348d7ee33b0eac1f0b444c"
	ct, err := hex.DecodeString(testCipherSession)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decrypt("session_numbers", ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, testPlain) {
		t.Errorf("Decrypt(session_numbers) = %q, want %q", got, testPlain)
	}
}

func TestDecryptUnknownName(t *testing.T) {
	if _, err := Decrypt("mystery", []byte("x")); err == nil {
		t.Fatal("expected an error for an unknown file name")
	}
}
