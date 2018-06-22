package bliss

import (
	"bytes"
	"crypto/rand"
	"testing"

	_ "github.com/coolsnady/hxd/chaincfg/chainec"
	_ "github.com/coolsnady/hxd/crypto"
)

func TestPublicKey(t *testing.T) {

	_, pk, err := Bliss.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal("Error in Generate keys")
	}

	pkBytes := pk.Serialize()
	restoredPK, err := Bliss.ParsePubKey(pkBytes)
	if err != nil {
		t.Fatal("Error in parsepubkey()")
	}
	pkBytes2 := restoredPK.Serialize()

	if !bytes.Equal(pkBytes, pkBytes2) {
		t.Fatal("Serialization() and ParsePubKey() do not match")
	}

	tp := pk.GetType()
	if tp != pqcTypeBliss {
		t.Fatal("GetType() result not matched")
	}

}
