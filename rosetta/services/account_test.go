package services

import (
	"reflect"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewAccountIdentifier(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	hexAddr := addr.Hex()
	metadata, err := types.MarshalMap(AccountMetadata{Address: addr.String()})
	if err != nil {
		t.Fatalf(err.Error())
	}

	referenceAccID := &types.AccountIdentifier{
		Address:  hexAddr,
		Metadata: metadata,
	}
	testAccID, rosettaError := newAccountIdentifier(addr)
	if rosettaError != nil {
		t.Fatalf("unexpected rosetta error: %v", rosettaError)
	}
	if !reflect.DeepEqual(referenceAccID, testAccID) {
		t.Errorf("reference ID %v != testID %v", referenceAccID, testAccID)
	}
}

func TestGetAddress(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf(err.Error())
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	hexAddr := addr.Hex()
	testAccID := &types.AccountIdentifier{
		Address: hexAddr,
	}

	testAddr, err := getAddress(testAccID)
	if err != nil {
		t.Fatal(err)
	}
	if testAddr != addr {
		t.Errorf("expected %v to be %v", testAddr.String(), addr.String())
	}

	defaultAddr := ethcommon.Address{}
	testAddr, err = getAddress(nil)
	if err == nil {
		t.Error("expected err for nil identifier")
	}
	if testAddr != defaultAddr {
		t.Errorf("expected errored addres to be %v not %v", defaultAddr.String(), testAddr.String())
	}
}
