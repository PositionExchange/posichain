package services

import (
	"reflect"
	"strings"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestGetAddressFromPublicKey(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	refAddr := crypto.PubkeyToAddress(key.PublicKey)
	compressedPublicKey := crypto.CompressPubkey(&key.PublicKey)
	addr, rosettaError := getAddressFromPublicKey(&types.PublicKey{
		Bytes:     compressedPublicKey,
		CurveType: types.Secp256k1,
	})
	if rosettaError != nil {
		t.Fatal(rosettaError)
	}
	if !reflect.DeepEqual(refAddr, *addr) {
		t.Errorf("expected adder %v, got %v", refAddr, *addr)
	}

	_, rosettaError = getAddressFromPublicKey(&types.PublicKey{
		Bytes:     compressedPublicKey,
		CurveType: types.Edwards25519,
	})
	if rosettaError == nil {
		t.Error("expected error")
	}

	_, rosettaError = getAddressFromPublicKey(nil)
	if rosettaError == nil {
		t.Error("expected error")
	}
}

func TestGetAddressFromKnownPublicKey(t *testing.T) {
	refCompressKey := "0x033e4c030253cd932a73e24f1a852de98b67647e0e96c5a3aba905a26d1c09bd2a"
	compressedPublicKey, _ := hexutil.Decode(refCompressKey)
	addr, rosettaError := getAddressFromPublicKey(&types.PublicKey{
		Bytes:     compressedPublicKey,
		CurveType: types.Secp256k1,
	})
	if rosettaError != nil {
		t.Fatal(rosettaError)
	}
	addrID, rosettaError := newAccountIdentifier(*addr)
	if rosettaError != nil {
		t.Fatal(rosettaError)
	}
	refHexAddr := "0xef5c5de1a7be1d9ff7d538ff1dccb8f06f9f7242"
	if strings.ToLower(addrID.Address) != strings.ToLower(refHexAddr) {
		t.Error("account ID from key is incorrect")
	}
}
