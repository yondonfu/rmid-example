package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var defaultCurve = elliptic.P256()

func GenerateAndStoreKey(filepath string) (*ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(defaultCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	// Store private key
	if err := StoreKey(filepath, priv); err != nil {
		return nil, err
	}

	pub := priv.Public().(*ecdsa.PublicKey)

	return pub, nil
}

func StoreKey(filepath string, priv *ecdsa.PrivateKey) error {
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	return pem.Encode(file, block)
}

func LoadKey(filepath string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block with EC private key")
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

func UnmarshalPublicKey(pub []byte) (*ecdsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	return key.(*ecdsa.PublicKey), nil
}
