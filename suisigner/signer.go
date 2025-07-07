package suisigner

import (
	"encoding/hex"
	"fmt"

	"github.com/pattonkan/sui-go/sui"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

var (
	TEST_MNEMONIC = "ordinary cry margin host traffic bulb start zone mimic wage fossil eight diagram clay say remove add atom"
	TEST_SEED     = []byte{4, 66, 186, 181, 112, 134, 111, 192, 149, 13, 68, 115, 67, 195, 58, 59, 33, 20, 200, 10, 150, 185, 145, 3, 106, 160, 105, 37, 4, 153, 172, 103, 69, 228, 114, 210, 176, 182, 208, 21, 252, 59, 50, 82, 135, 160, 1, 131, 156, 104, 159, 240, 183, 20, 250, 216, 26, 228, 91, 220, 15, 222, 75, 91}
	TEST_ADDRESS  = sui.MustAddressFromHex("0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5")
)

// FIXME support more than ed25519
type Signer struct {
	KeypairEd25519   *KeypairEd25519
	KeypairSecp256k1 *KeypairSecp256k1
	KeypairSecp256r1 *KeypairSecp256r1
	Address          *sui.Address
}

func NewSigner(seed []byte, flag KeySchemeFlag) *Signer {
	signer := Signer{}

	// IOTA_DIFF iota ignore flag when signature scheme is ed25519
	switch flag {
	case KeySchemeFlagEd25519:
		k := NewKeypairEd25519FromSeed(seed)
		signer.KeypairEd25519 = k
	case KeySchemeFlagSecp256k1:
		k := NewKeypairSecp256k1FromSeed(seed)
		signer.KeypairSecp256k1 = k
	case KeySchemeFlagSecp256r1:
		k := NewKeypairSecp256r1FromSeed(seed)
		signer.KeypairSecp256r1 = k
	default:
		panic("unrecognizable key scheme flag")
	}
	signer.Address = signer.calculateAddress(flag)
	return &signer
}

// there are only 256 different signers can be generated
func NewSignerByIndex(seed []byte, flag KeySchemeFlag, index int) *Signer {
	seed[0] = seed[0] + byte(index)
	return NewSigner(seed, flag)
}

// generate keypair (signer) with mnemonic which is referring the Sui monorepo in the following code
//
// let phrase = "asset pink record dawn hundred sure various crime client enforce carbon blossom";
// let mut keystore = Keystore::from(InMemKeystore::new_insecure_for_tests(0));
// let generated_address = keystore.import_from_mnemonic(&phrase, SignatureScheme::ED25519, None, None).unwrap();
func NewSignerWithMnemonic(mnemonic string, flag KeySchemeFlag) (*Signer, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	var derivePath string
	switch flag {
	case KeySchemeFlagEd25519:
		derivePath = DerivationPathEd25519
	case KeySchemeFlagSecp256k1:
		derivePath = DerivationPathSecp256k1
	default:
		return nil, fmt.Errorf("unsupported key scheme")
	}

	key, err := DeriveForPath(derivePath, seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive %s key for path: %w", flag.String(), err)
	}
	return NewSigner(key.Key, flag), nil
}

func (s *Signer) PrivateKey() []byte {
	switch {
	case s.KeypairEd25519 != nil:
		return s.KeypairEd25519.PriKey
	case s.KeypairSecp256k1 != nil:
		return s.KeypairSecp256k1.PriKey.Serialize()
	default:
		return nil
	}
}

func (s *Signer) PublicKey() []byte {
	switch {
	case s.KeypairEd25519 != nil:
		return s.KeypairEd25519.PubKey
	case s.KeypairSecp256k1 != nil:
		return s.KeypairSecp256k1.PubKey.SerializeCompressed()
	case s.KeypairSecp256r1 != nil:
		return compressSecp256r1PublicKey(s.KeypairSecp256r1.PubKey)
	default:
		return nil
	}
}

// Signer implements the UserSignature trait in Sui Rust SDK
// refer sui-rust-sdk/crates/sui-sdk-types/src/crypto/signature.rs
//
//	pub enum UserSignature {
//	    Simple(SimpleSignature),
//	    Multisig(MultisigAggregatedSignature),
//	    ZkLogin(Box<ZkLoginAuthenticator>),
//	    Passkey(PasskeyAuthenticator),
//	}
//
// SimpleSignature include Ed25519, Secp256k1, and Secp256r1 signatures
func (s *Signer) Sign(data []byte) Signature {
	var sig Signature
	switch {
	case s.KeypairEd25519 != nil:
		sig.Ed25519SuiSignature = NewEd25519SuiSignature(s, data)
	case s.KeypairSecp256k1 != nil:
		sig.Secp256k1SuiSignature = NewSecp256k1SuiSignature(s, data)
	case s.KeypairSecp256r1 != nil:
		sig.Secp256r1SuiSignature = NewSecp256r1SuiSignature(s, data)
	default:
		panic("signer does not have keypair")
	}
	return sig
}

// it is the signing_digest() interface in Sui Rust SDK
func (a *Signer) SignDigest(msg []byte, intent Intent) (Signature, error) {
	data := MessageWithIntent(intent, bcsBytes(msg))
	hash := blake2b.Sum256(data)
	return a.Sign(hash[:]), nil
}

func (a *Signer) calculateAddress(flag KeySchemeFlag) *sui.Address {
	var buf []byte
	switch flag {
	case KeySchemeFlagEd25519:
		buf = []byte{KeySchemeFlagEd25519.Byte()}
	case KeySchemeFlagSecp256k1:
		buf = []byte{KeySchemeFlagSecp256k1.Byte()}
	case KeySchemeFlagSecp256r1:
		buf = []byte{KeySchemeFlagSecp256r1.Byte()}
	default:
		panic("unrecognizable key scheme flag")
	}
	buf = append(buf, a.PublicKey()...)
	addrBytes := blake2b.Sum256(buf)
	addr := "0x" + hex.EncodeToString(addrBytes[:])
	return sui.MustAddressFromHex(addr)
}

type bcsBytes []byte

func (b bcsBytes) MarshalBCS() ([]byte, error) {
	return b, nil
}
