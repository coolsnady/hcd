package bliss

import (
	"crypto/rand"
	"io"

	"github.com/coolsnady/bliss"
	"github.com/coolsnady/bliss/poly"
	"github.com/coolsnady/bliss/sampler"
	dcrcrypto "github.com/coolsnady/hxd/crypto"
)

var pqcTypeBliss = 4

type blissDSA struct {

	// Private keys
	newPrivateKey    func(s1, s2, a *poly.PolyArray) dcrcrypto.PrivateKey
	privKeyFromBytes func(pk []byte) (dcrcrypto.PrivateKey, dcrcrypto.PublicKey)
	privKeyBytesLen  func() int

	// Public keys
	newPublicKey   func(a *poly.PolyArray) dcrcrypto.PublicKey
	parsePubKey    func(pubKeyStr []byte) (dcrcrypto.PublicKey, error)
	pubKeyBytesLen func() int

	// Signatures
	newSignature      func(z1, z2 *poly.PolyArray, c []uint32) dcrcrypto.Signature
	parseDERSignature func(sigStr []byte) (dcrcrypto.Signature, error)
	parseSignature    func(sigStr []byte) (dcrcrypto.Signature, error)
	recoverCompact    func(signature, hash []byte) (dcrcrypto.PublicKey, bool, error)

	//
	generateKey func(rand io.Reader) (dcrcrypto.PrivateKey, dcrcrypto.PublicKey, error)
	sign        func(priv dcrcrypto.PrivateKey, hash []byte) (dcrcrypto.Signature, error)
	verify      func(pub dcrcrypto.PublicKey, hash []byte, sig dcrcrypto.Signature) bool

	// Symmetric cipher encryption
	//generateSharedSecret func(privkey []byte, x, y *big.Int) []byte
	//encrypt              func(x, y *big.Int, in []byte) ([]byte, error)
	//decrypt              func(privkey []byte, in []byte) ([]byte, error)
}

// Private keys
func (sp blissDSA) NewPrivateKey(s1, s2, a *poly.PolyArray) dcrcrypto.PrivateKey {
	return sp.newPrivateKey(s1, s2, a)
}
func (sp blissDSA) PrivKeyFromBytes(pk []byte) (dcrcrypto.PrivateKey, dcrcrypto.PublicKey) {
	return sp.privKeyFromBytes(pk)
}
func (sp blissDSA) PrivKeyBytesLen() int {
	return sp.privKeyBytesLen()
}

// Public keys
func (sp blissDSA) NewPublicKey(a *poly.PolyArray) dcrcrypto.PublicKey {
	return sp.newPublicKey(a)
}
func (sp blissDSA) ParsePubKey(pubKeyStr []byte) (dcrcrypto.PublicKey, error) {
	return sp.parsePubKey(pubKeyStr)
}
func (sp blissDSA) PubKeyBytesLen() int {
	return sp.pubKeyBytesLen()
}

// Signatures
func (sp blissDSA) NewSignature(z1, z2 *poly.PolyArray, c []uint32) dcrcrypto.Signature {
	return sp.newSignature(z1, z2, c)
}
func (sp blissDSA) ParseDERSignature(sigStr []byte) (dcrcrypto.Signature, error) {
	return sp.parseDERSignature(sigStr)
}
func (sp blissDSA) ParseSignature(sigStr []byte) (dcrcrypto.Signature, error) {
	return sp.parseSignature(sigStr)
}
func (sp blissDSA) RecoverCompact(signature, hash []byte) (dcrcrypto.PublicKey, bool,
	error) {
	return sp.recoverCompact(signature, hash)
}

// ECDSA
func (sp blissDSA) GenerateKey(rand io.Reader) (dcrcrypto.PrivateKey, dcrcrypto.PublicKey,
	error) {
	return sp.generateKey(rand)
}
func (sp blissDSA) Sign(priv dcrcrypto.PrivateKey, hash []byte) (dcrcrypto.Signature, error) {
	return sp.sign(priv, hash)
}
func (sp blissDSA) Verify(pub dcrcrypto.PublicKey, hash []byte, sig dcrcrypto.Signature) bool {
	return sp.verify(pub, hash, sig)
}

func newBlissDSA() DSA {
	var bliss DSA = &blissDSA{

		// Private keys
		newPrivateKey: func(s1, s2, a *poly.PolyArray) dcrcrypto.PrivateKey {
			if s1 == nil || s2 == nil || a == nil {
				return nil
			}

			n := s1.Param().N
			s1data := s1.GetData()
			s2data := s2.GetData()
			ret := make([]byte, n*2+1)
			ret[0] = byte(s1.Param().Version)
			s1part := ret[1 : 1+n]
			s2part := ret[1+n:]
			for i := 0; i < int(n); i++ {
				s1part[i] = byte(s1data[i] + 4)
				s2part[i] = byte(s2data[i] + 4)
			}

			blissPK, err := bliss.DecodePrivateKey(ret)
			if err != nil {
				return nil
			}

			return &PrivateKey{
				PrivateKey: *blissPK,
			}
		},
		privKeyFromBytes: func(pk []byte) (dcrcrypto.PrivateKey, dcrcrypto.PublicKey) {
			blissPK, err := bliss.DeserializePrivateKey(pk)
			if err != nil {
				return nil, nil
			}
			var privateKey PrivateKey
			var publicKey PublicKey
			privateKey.PrivateKey = *blissPK
			publicKey.PublicKey = *(blissPK.PublicKey())
			return privateKey, publicKey
		},

		privKeyBytesLen: func() int {
			return BlissPrivKeyLen
		},

		// Public keys
		newPublicKey: func(a *poly.PolyArray) dcrcrypto.PublicKey {
			if a == nil {
				return nil
			}

			n := a.Param().N
			data := a.GetData()
			ret := make([]byte, n*2+1)
			ret[0] = byte(a.Param().Version)
			for i := 0; i < int(n); i++ {
				ret[i*2+1] = byte(uint16(data[i]) >> 8)
				ret[i*2+2] = byte(uint16(data[i]) & 0xff)
			}
			blissPK, err := bliss.DecodePublicKey(ret)
			if err != nil {
				return nil
			}
			return &PublicKey{
				PublicKey: *blissPK,
			}
		},
		parsePubKey: func(pubKeyStr []byte) (dcrcrypto.PublicKey, error) {
			blissPK, err := bliss.DeserializePublicKey(pubKeyStr)
			if err != nil {
				return nil, err
			}
			return &PublicKey{
				PublicKey: *blissPK,
			}, nil

		},
		pubKeyBytesLen: func() int {
			return BlissPubKeyLen
		},

		// Signatures
		newSignature: func(z1, z2 *poly.PolyArray, c []uint32) dcrcrypto.Signature {
			if z1 == nil || z2 == nil || c == nil {
				return nil
			}
			n := z1.Param().N
			kappa := z1.Param().Kappa
			z1len := n * 2
			z2len := n + n/8
			clen := 2 * kappa

			z1data := z1.GetData()
			z2data := z2.GetData()
			cdata := c

			ret := make([]byte, 1+z1len+z2len+clen)
			ret[0] = byte(z1.Param().Version)

			z1part := ret[1 : 1+z1len]
			z2part := ret[1+z1len : 1+z1len+z2len]
			cpart := ret[1+z1len+z2len:]

			// It is easy to store z1. Take each element as
			// an uint16, although they are actually a littble
			// bit smaller than 16 bits.
			for i := 0; i < int(n); i++ {
				tmp := z1.NumModQ(z1data[i])
				z1part[i*2] = byte(uint16(tmp) >> 8)
				z1part[i*2+1] = byte(uint16(tmp) & 0xff)
			}

			// z2 is much smaller than z1, bounded by p/2
			// An additional bit array is used to store the signs
			z2left := z2part[:n]
			z2right := z2part[n:]
			for i := 0; i < int(n); i++ {
				z2left[i] = byte(uint16(bliss.Abs(z2data[i])) & 0xff)
			}
			for i := 0; i < int(n)/8; i++ {
				tmp := byte(0)
				for j := 0; j < 8; j++ {
					tmp <<= 1
					if z2data[i*8+j] > 0 {
						tmp += 1
					}
				}
				// Each extra bit takes a byte array of size n/8
				z2right[i] = tmp
			}

			// c is represented by a list of kappa integers in [0,n)
			// For simplicity, we use 2 bytes to store each index.
			for i := 0; i < int(kappa); i++ {
				cpart[i*2] = byte(uint16(cdata[i]) >> 8)
				cpart[i*2+1] = byte(uint16(cdata[i]) & 0xff)
			}

			sig, err := bliss.DecodeSignature(ret)
			if err != nil {
				return nil
			}
			return &Signature{
				Signature: *sig,
			}

		},
		parseDERSignature: func(sigStr []byte) (dcrcrypto.Signature, error) {
			sig, err := bliss.DeserializeBlissSignature(sigStr)
			if err != nil {
				return nil, err
			}

			return &Signature{
				Signature: *sig,
			}, nil
		},
		parseSignature: func(sigStr []byte) (dcrcrypto.Signature, error) {
			sig, err := bliss.DeserializeBlissSignature(sigStr)
			if err != nil {
				return nil, err
			}

			return &Signature{
				Signature: *sig,
			}, nil
		},
		recoverCompact: func(signature, hash []byte) (dcrcrypto.PublicKey, bool, error) {
			return nil, false, nil
		},

		generateKey: func(rand io.Reader) (dcrcrypto.PrivateKey, dcrcrypto.PublicKey, error) {
			seed := make([]byte, sampler.SHA_512_DIGEST_LENGTH)
			rand.Read(seed)
			entropy, err := sampler.NewEntropy(seed)
			if err != nil {
				return nil, nil, err
			}
			blissPK, err := bliss.GeneratePrivateKey(BlissVersion, entropy)
			if err != nil {
				return nil, nil, err
			}
			privateKey := &PrivateKey{
				PrivateKey: *blissPK,
			}
			publicKey := &PublicKey{
				PublicKey: *blissPK.PublicKey(),
			}
			return privateKey, publicKey, nil
		},

		sign: func(priv dcrcrypto.PrivateKey, hash []byte) (dcrcrypto.Signature, error) {
			seed := make([]byte, sampler.SHA_512_DIGEST_LENGTH)
			rand.Read(seed)
			entropy, err := sampler.NewEntropy(seed)
			if err != nil {
				return nil, err
			}
			priv1 := priv.(PrivateKey)
			sig, err := priv1.Sign(hash, entropy)
			if err != nil {
				return nil, err
			}
			return &Signature{
				Signature: *sig,
			}, nil
		},

		verify: func(pub dcrcrypto.PublicKey, hash []byte, sig dcrcrypto.Signature) bool {
			signature := sig.(*Signature)
			blissSig := signature.Signature
			result, _ := pub.(*PublicKey).Verify(hash, &blissSig)
			return result
		},
	}

	return bliss.(DSA)
}
