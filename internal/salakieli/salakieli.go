package salakieli

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// passphrases maps the basename (without the .salakieli extension) of each
// encrypted save file to the two passphrase strings whose first 16 bytes
// form the AES-128 key and CTR IV. Noita encrypts these files with AES-CTR
// using a fixed key/IV per kind of data: _stats and _streaks share one pair,
// session_numbers and magic_numbers share another.
var passphrases = map[string][2]string{
	"_stats":          {"SecretsOfTheAllSeeing", "ThreeEyesAreWatchingYou"},
	"_streaks":        {"SecretsOfTheAllSeeing", "ThreeEyesAreWatchingYou"},
	"session_numbers": {"KnowledgeIsTheHighestOfTheHighest", "WhoWouldntGiveEverythingForTrueKnowledge"},
	"magic_numbers":   {"KnowledgeIsTheHighestOfTheHighest", "WhoWouldntGiveEverythingForTrueKnowledge"},
	"player":          {"WeSeeATrueSeekerOfKnowledge", "YouAreSoCloseToBeingEnlightened"},
	"world_state":     {"TheTruthIsThatThereIsNothing", "MoreValuableThanKnowledge"},
}

// Decrypt decrypts the contents of a .salakieli file. name is the file's
// basename without the extension, e.g. "_stats" for "_stats.salakieli".
func Decrypt(name string, data []byte) ([]byte, error) {
	pp, ok := passphrases[name]
	if !ok {
		return nil, fmt.Errorf("no known passphrase for %q", name)
	}
	block, err := aes.NewCipher([]byte(pp[0][:16]))
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, []byte(pp[1][:16]))
	out := make([]byte, len(data))
	stream.XORKeyStream(out, data)
	return out, nil
}

// DecryptFile reads and decrypts a .salakieli file, selecting the passphrase
// from the file name.
func DecryptFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return Decrypt(name, data)
}
