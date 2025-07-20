package suicrypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// RFC 6979 deterministic k generation (SHA-256, secp256r1)
func deterministicK(priv *ecdsa.PrivateKey, hash []byte) *big.Int {
	curveOrder := priv.Curve.Params().N
	holen := sha256.Size
	rolen := (curveOrder.BitLen() + 7) >> 3
	// Step B: Process priv.D and hash as octets
	bx := append(int2octets(priv.D, rolen), bits2octets(hash, curveOrder, rolen)...)
	// Step C: Set V = 0x01 0x01 ... 0x01
	V := bytes.Repeat([]byte{0x01}, holen)
	// Step D: Set K = 0x00 0x00 ... 0x00
	K := bytes.Repeat([]byte{0x00}, holen)
	// Step E: K = HMAC_K(V || 0x00 || bx)
	K = hmacSHA256(K, append(append(V, 0x00), bx...))
	// Step F: V = HMAC_K(V)
	V = hmacSHA256(K, V)
	// Step G: K = HMAC_K(V || 0x01 || bx)
	K = hmacSHA256(K, append(append(V, 0x01), bx...))
	// Step H: V = HMAC_K(V)
	V = hmacSHA256(K, V)

	for {
		// Step K: Generate candidate k
		V = hmacSHA256(K, V)
		k := new(big.Int).SetBytes(V)
		if k.Sign() > 0 && k.Cmp(curveOrder) < 0 {
			return k
		}
		// Step H: Update K and V
		K = hmacSHA256(K, append(V, 0x00))
		V = hmacSHA256(K, V)
	}
}

func int2octets(x *big.Int, rolen int) []byte {
	out := x.Bytes()
	if len(out) < rolen {
		pad := make([]byte, rolen-len(out))
		out = append(pad, out...)
	}
	return out
}

func bits2octets(in []byte, curveOrder *big.Int, rolen int) []byte {
	z := new(big.Int).SetBytes(in)
	z.Mod(z, curveOrder)
	return int2octets(z, rolen)
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// Deterministic ECDSA sign for P-256, RFC 6979, SHA256
// The Golang's standard lib doesn't support deterministic Secp256r1 sign
func DeterministicSecp256r1Sign(priv *ecdsa.PrivateKey, msg []byte) (r, s *big.Int, err error) {
	hash := sha256.Sum256(msg)
	k := deterministicK(priv, hash[:])
	curve := priv.Curve
	N := curve.Params().N

	// (x, _) = k*G
	x, _ := curve.ScalarBaseMult(k.Bytes())
	r = new(big.Int).Mod(x, N)
	if r.Sign() == 0 {
		return nil, nil, fmt.Errorf("r is zero")
	}

	kInv := new(big.Int).ModInverse(k, N)
	z := new(big.Int).SetBytes(hash[:])
	s = new(big.Int).Mul(r, priv.D)
	s.Add(s, z)
	s.Mul(s, kInv)
	s.Mod(s, N)
	if s.Sign() == 0 {
		return nil, nil, fmt.Errorf("s is zero")
	}

	// Enforce low-S form
	halfN := new(big.Int).Rsh(N, 1)
	if s.Cmp(halfN) == 1 {
		s.Sub(N, s)
	}

	return r, s, nil
}

// Verifies a secp256r1 ECDSA signature (RFC 6979 compatible).
func VerifySecp256r1(pub *ecdsa.PublicKey, msg []byte, r, s *big.Int) bool {
	curve := pub.Curve
	N := curve.Params().N

	// Check that r and s are in the valid range
	if r.Sign() <= 0 || s.Sign() <= 0 {
		return false
	}
	if r.Cmp(N) >= 0 || s.Cmp(N) >= 0 {
		return false
	}

	// Hash message
	hash := sha256.Sum256(msg)
	z := new(big.Int).SetBytes(hash[:])

	// w = s^-1 mod n
	w := new(big.Int).ModInverse(s, N)
	if w == nil {
		return false
	}

	// u1 = z * w mod n
	u1 := new(big.Int).Mul(z, w)
	u1.Mod(u1, N)

	// u2 = r * w mod n
	u2 := new(big.Int).Mul(r, w)
	u2.Mod(u2, N)

	// (x1, y1) = u1*G + u2*Q
	x1, y1 := curve.ScalarBaseMult(u1.Bytes())
	x2, y2 := curve.ScalarMult(pub.X, pub.Y, u2.Bytes())
	X, Y := curve.Add(x1, y1, x2, y2)
	if X == nil || Y == nil {
		return false
	}

	// v = x1 mod n
	v := new(big.Int).Mod(X, N)

	return v.Cmp(r) == 0
}
