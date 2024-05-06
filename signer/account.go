package signer

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/move_types"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

const (
	ADDRESS_LENGTH = 64
	TEST_MNEMONIC  = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
)

var TEST_ADDRESS, _ = sui_types.NewAddressFromHex("0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5")

type Account struct {
	KeyPair *Signer
	Address string
}

type SignatureScheme []byte

func NewAccount(scheme SignatureScheme, seed []byte) *Account {
	// suiKeyPair := NewSuiKeyPair(scheme, seed)
	// tmp := []byte{scheme.Flag()}
	// tmp = append(tmp, suiKeyPair.PublicKey()...)
	// addrBytes := blake2b.Sum256(tmp)
	// address := "0x" + hex.EncodeToString(addrBytes[:])[:ADDRESS_LENGTH]

	s := NewSigner(seed)

	return &Account{
		KeyPair: s,
		Address: s.Address,
	}
}

func (a *Account) AccountAddress() *move_types.AccountAddress {
	addr := a.Address[2:]
	data, err := hex.DecodeString(addr)
	if err != nil {
		panic(err)
	}
	var accountAddress move_types.AccountAddress
	copy(accountAddress[32-len(data):], data[:])
	return &accountAddress
}

// func NewAccountWithKeystore(keystore string) (*Account, error) {
// 	ksByte, err := base64.StdEncoding.DecodeString(keystore)
// 	if err != nil {
// 		return nil, err
// 	}
// 	scheme, err := NewSignatureScheme(ksByte[0])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return NewAccount(scheme, ksByte[1:]), nil
// }

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := derivation.DeriveForPath("m/44'/784'/0'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	return NewAccount([]byte{}, key.Key), nil
}

func (a *Account) Sign(data []byte) Signature {
	// FIXME temporarily support only Ed25519
	ed25519Signature := Ed25519SuiSignature{}
	copy(ed25519Signature.Signature[:], ed25519.Sign(a.KeyPair.PriKey, data))
	return Signature{
		Ed25519SuiSignature: &ed25519Signature,
	}
}

// FIXME support only ed25519 now
func (a *Account) SignTransactionBlock(txnBytes []byte, intent Intent) (Signature, error) {
	message := MessageWithIntent(intent, bcsBytes(txnBytes))
	data, err := bcs.Marshal(message)
	if err != nil {
		return Signature{}, err
	}
	hash := blake2b.Sum256(data)
	return a.Sign(hash[:]), nil
}

type bcsBytes []byte

func (b bcsBytes) MarshalBCS() ([]byte, error) {
	return b, nil
}
