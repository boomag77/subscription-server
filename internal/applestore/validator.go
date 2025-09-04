package appstore

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"subscription-server/internal/jws"
)

type appleJWSValidator struct {
	rootCA string
}

func NewAppleJWSValidator() jws.JWSValidator {
	return &appleJWSValidator{
		rootCA: appleRootCA,
	}
}

func (v *appleJWSValidator) Validate(header string, payload string, signature string) error {
	if header == "" || payload == "" || signature == "" {
		return fmt.Errorf("header, payload, and signature must all be non-empty")
	}

	hdrBytes, err := decodeBase64String(header)
	if err != nil {
		return fmt.Errorf("decode header: %w", err)
	}
	var hdr jwsHeader
	if err := json.Unmarshal(hdrBytes, &hdr); err != nil {
		return fmt.Errorf("unmarshal header: %w", err)
	}
	if hdr.Alg != "ES256" {
		return fmt.Errorf("unsupported algorithm: %s", hdr.Alg)
	}
	if len(hdr.X5c) == 0 {
		return fmt.Errorf("missing X5c field")
	}
	leaf := hdr.X5c[0]

	der, err := base64.StdEncoding.DecodeString(leaf)
	if err != nil {
		return fmt.Errorf("decode leaf: %w", err)
	}
	leafCert, err := x509.ParseCertificate(der)
	if err != nil {
		return fmt.Errorf("parse leaf cert: %w", err)
	}
	pubKey, ok := leafCert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an ECDSA public key")
	}

	signingInput := fmt.Sprintf("%s.%s", header, payload)
	digest := sha256.Sum256([]byte(signingInput)) //hash the signing input

	sigBytes, err := decodeBase64String(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	if !ecdsa.VerifyASN1(pubKey, digest[:], sigBytes) {
		return fmt.Errorf("invalid JWS signature")
	}
	if err := v.validateAppleChain(leafCert, hdr.X5c[1:]); err != nil {
		return fmt.Errorf("validate Apple chain: %w", err)
	}

	return nil
}

func (v *appleJWSValidator) validateAppleChain(leafCert *x509.Certificate, intermCerts []string) error {

	intermediates := x509.NewCertPool()
	roots := x509.NewCertPool()
	for _, certB64 := range intermCerts {
		der, err := base64.StdEncoding.DecodeString(certB64)
		if err != nil {
			return fmt.Errorf("decode cert: %w", err)
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			return fmt.Errorf("parse cert: %w", err)
		}
		intermediates.AddCert(cert)
	}

	block, _ := pem.Decode([]byte(v.rootCA))
	if block == nil {
		return fmt.Errorf("failed to parse root CA PEM")
	}

	rootCA, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse root Apple CA: %w", err)
	}
	roots.AddCert(rootCA)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
	}
	if _, err := leafCert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify Apple certificate chain: %w", err)
	}

	return nil
}
