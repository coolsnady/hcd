package crypto

import (
	"fmt"
	"testing"

	_ "github.com/coolsnady/hxd/chaincfg/chainec"
)

func TestCrypto(t *testing.T) {
	fmt.Println("test start")
	var pk PublicKey
	pk = new(PublicKeyAdapter)
	fmt.Println(pk.GetType())
}
