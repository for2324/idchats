package biubiulib

import (
	"strings"

	"golang.org/x/net/idna"

	"golang.org/x/crypto/sha3"
)

var p = idna.New(idna.MapForLookup(), idna.StrictDomainName(false), idna.Transitional(false))
var pStrict = idna.New(idna.MapForLookup(), idna.StrictDomainName(true), idna.Transitional(false))

// Normalize normalizes a name according to the ENS rules
func Normalize(input string) (output string, err error) {
	output, err = p.ToUnicode(input)
	if err != nil {
		return
	}
	// If the name started with a period then ToUnicode() removes it, but we want to keep it
	if strings.HasPrefix(input, ".") && !strings.HasPrefix(output, ".") {
		output = "." + output
	}
	return
}

// LabelHash generates a simple hash for a piece of a name.
func LabelHash(label string) (hash [32]byte, err error) {
	normalizedLabel, err := Normalize(label)
	if err != nil {
		return
	}

	sha := sha3.NewLegacyKeccak256()
	if _, err = sha.Write([]byte(normalizedLabel)); err != nil {
		return
	}
	sha.Sum(hash[:0])
	return
}

// NameHash generates a hash from a name that can be used to
// look up the name in ENS
func NameHash(name string) (hash [32]byte, err error) {
	if name == "" {
		return
	}
	normalizedName, err := Normalize(name)
	if err != nil {
		return
	}
	parts := strings.Split(normalizedName, ".")
	for i := len(parts) - 1; i >= 0; i-- {
		if hash, err = nameHashPart(hash, parts[i]); err != nil {
			return
		}
	}
	return
}

func nameHashPart(currentHash [32]byte, name string) (hash [32]byte, err error) {
	sha := sha3.NewLegacyKeccak256()
	if _, err = sha.Write(currentHash[:]); err != nil {
		return
	}
	nameSha := sha3.NewLegacyKeccak256()
	if _, err = nameSha.Write([]byte(name)); err != nil {
		return
	}
	nameHash := nameSha.Sum(nil)
	if _, err = sha.Write(nameHash); err != nil {
		return
	}
	sha.Sum(hash[:0])
	return
}

func EncodeName(name string) []byte {
	n := strings.Trim(name, ".")
	if len(n) == 0 {
		return []byte{}
	}
	segments := strings.Split(n, ".")
	buf := make([]byte, 0, len(n)+2*len(segments))
	for _, seg := range segments {
		if len(seg) == 0 {
			continue
		}
		buf = append(buf, byte(len(seg)))
		buf = append(buf, []byte(seg)...)
	}
	buf = append(buf, 0x00)
	return buf
}
