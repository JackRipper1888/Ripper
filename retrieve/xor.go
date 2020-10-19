package retrieve

import (
	"bytes"
	"math/big"
	"math/bits"
	"errors"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/minio/sha256-simd"
)

// XORKeySpace is a KeySpace which:
// - normalizes identifiers using a cryptographic hash (sha256)
// - measures distance by XORing keys together
var XORKeySpace = &xorKeySpace{}
//var _ KeySpace = XORKeySpace // ensure it conforms

type xorKeySpace struct{}

// Key converts an identifier into a Key in this space.
func (s *xorKeySpace) Key(id []byte) Key {
	hash := sha256.Sum256(id)
	key := hash[:]
	return Key{
		Space:    s,
		Original: id,
		Bytes:    key,
	}
}

// Equal returns whether keys are equal in this key space
func (s *xorKeySpace) Equal(k1, k2 Key) bool {
	return bytes.Equal(k1.Bytes, k2.Bytes)
}

// Distance returns the distance metric in this key space
func (s *xorKeySpace) Distance(k1, k2 Key) *big.Int {
	// XOR the keys
	k3 := XOR(k1.Bytes, k2.Bytes)
	// interpret it as an integer
	dist := big.NewInt(0).SetBytes(k3)
	return dist
}

// Less returns whether the first key is smaller than the second.
func (s *xorKeySpace) Less(k1, k2 Key) bool {
	return bytes.Compare(k1.Bytes, k2.Bytes) < 0
}

// ZeroPrefixLen returns the number of consecutive zeroes in a byte slice.
func ZeroPrefixLen(id []byte) int {
	for i, b := range id {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(id) * 8
}


func XOR(a, b []byte) []byte {
	c := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}



// Returned if a routing table query returns no results. This is NOT expected
// behaviour
var ErrLookupFailure = errors.New("failed to find any peer in table")

// ID for IpfsDHT is in the XORKeySpace
//
// The type dht.ID signifies that its contents have been hashed from either a
// peer.ID or a util.Key. This unifies the keyspace
type ID []byte

func (id ID) equal(other ID) bool {
	return bytes.Equal(id, other)
}

func (id ID) less(other ID) bool {
	a := Key{Space: XORKeySpace, Bytes: id}
	b := Key{Space: XORKeySpace, Bytes: other}
	return a.Less(b)
}

func xor(a, b ID) ID {
	return ID(XOR(a, b))
}

func CommonPrefixLen(a, b ID) int {
	return ZeroPrefixLen(XOR(a, b))
}

// ConvertPeerID creates a DHT ID by hashing a Peer ID (Multihash)
func ConvertPeerID(id peer.ID) ID {
	hash := sha256.Sum256([]byte(id))
	return hash[:]
}

// ConvertKey creates a DHT ID by hashing a local key (String)
func ConvertKey(id string) ID {
	hash := sha256.Sum256([]byte(id))
	return hash[:]
}

// Closer returns true if a is closer to key than b is
func Closer(a, b peer.ID, key string) bool {
	aid := ConvertPeerID(a)
	bid := ConvertPeerID(b)
	tgt := ConvertKey(key)
	adist := xor(aid, tgt)
	bdist := xor(bid, tgt)

	return adist.less(bdist)
}
