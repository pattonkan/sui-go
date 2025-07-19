package suicrypto

// func ParsePublicKeyAndSignatureFromRawBytes(raw []byte) (KeySchemeFlag, []byte, error) {
// 	switch signature[0] {
// 	case KeySchemeFlagEd25519.Byte():
// 		if len(signature) != ed25519.PublicKeySize+ed25519.SignatureSize+1 {
// 			return errors.New("invalid ed25519 signature")
// 		}
// 		var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
// 		copy(signatureBytes[:], signature)
// 		s.Ed25519SuiSignature = &Ed25519SuiSignature{
// 			Signature: signatureBytes,
// 		}
// 	case KeySchemeFlagSecp256k1.Byte():
// 		if len(signature) != KeypairSecp256k1PublicKeySize+KeypairSecp256k1SignatureSize+1 {
// 			return errors.New("invalid secp256k1 signature")
// 		}
// 		var signatureBytes [KeypairSecp256k1PublicKeySize + KeypairSecp256k1SignatureSize + 1]byte
// 		copy(signatureBytes[:], signature)
// 		s.Secp256k1SuiSignature = &Secp256k1SuiSignature{
// 			Signature: signatureBytes,
// 		}
// 	default:
// 		return errors.New("not supported signature")
// 	}
// }
