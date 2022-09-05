package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func GenKey(privFile string) {
	pub, err := GenerateAndStoreKey(privFile)
	if err != nil {
		log.Fatal(err)
	}

	derBytes, err := MarshalPublicKey(pub)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(derBytes))
}

func RootRMID(file string) {
	s, err := NewRawMediaNodeStream(file)
	if err != nil {
		log.Fatal(err)
	}

	rootHash, err := s.RootHash()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Root RMID: %s\n", rootHash)
}

func SignVideo(file string, privFile string) {
	s, err := NewRawMediaNodeStream(file)
	if err != nil {
		log.Fatal(err)
	}

	priv, err := LoadKey(privFile)
	if err != nil {
		log.Fatal(err)
	}

	sigs, err := s.Sign(priv)
	if err != nil {
		log.Fatal(err)
	}

	for _, sig := range sigs {
		fmt.Println(hex.EncodeToString(sig))
	}
}

func VerifyVideo(file string, pub string, sigFile string) {
	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		log.Fatal(err)
	}

	pubkey, err := UnmarshalPublicKey(pubBytes)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewRawMediaNodeStream(file)
	if err != nil {
		log.Fatal(err)
	}

	sigs, err := ReadHexLines(sigFile)
	if err != nil {
		log.Fatal(err)
	}

	if !s.Verify(pubkey, sigs) {
		log.Fatal("failed signature verification!")
	}

	fmt.Println("signature verification successful!")
}

func main() {
	args := os.Args[1:]

	switch opt := args[0]; opt {
	case "gen-key":
		privFile := args[1]
		GenKey(privFile)
	case "root-rmid":
		file := args[1]
		RootRMID(file)
	case "sign-video":
		file := args[1]
		privFile := args[2]
		SignVideo(file, privFile)
	case "verify-video":
		file := args[1]
		pub := args[2]
		sigFile := args[3]
		VerifyVideo(file, pub, sigFile)
	default:
		fmt.Println("invalid option")
	}
}
