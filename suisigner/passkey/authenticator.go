package passkey

import "github.com/pattonkan/sui-go/suisigner/suicrypto"

// pub struct PasskeyAuthenticator {
//     /// The secp256r1 public key for this passkey.
//     public_key: Secp256r1PublicKey,

//     /// The secp256r1 signature from the passkey.
//     signature: Secp256r1Signature,

//     /// Parsed base64url decoded challenge bytes from `client_data_json.challenge`.
//     challenge: Vec<u8>,

//     /// Opaque authenticator data for this passkey signature.
//     ///
//     /// See [Authenticator Data](https://www.w3.org/TR/webauthn-2/#sctn-authenticator-data) for
//     /// more information on this field.
//     authenticator_data: Vec<u8>,

//	    /// Structured, unparsed, JSON for this passkey signature.
//	    ///
//	    /// See [CollectedClientData](https://www.w3.org/TR/webauthn-2/#dictdef-collectedclientdata)
//	    /// for more information on this field.
//	    client_data_json: String,
//	}
type PasskeyAuthenticator struct {
	PublicKey         suicrypto.Secp256r1PubKey
	Signature         []byte
	Challenge         []byte
	AuthenticatorData []byte
	ClientDataJson    string
}
