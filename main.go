package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fabric-did/did"
	"fabric-did/ssi"
	"fabric-did/vc"
	"fmt"
)

func main() {
	didID, _ := did.ParseDID("did:example:123")

	// Empty did document:
	doc := &did.Document{
		Context: []ssi.URI{did.DIDContextV1URI()},
		ID:      *didID,
	}
	// Add an assertionMethod
	keyID, _ := did.ParseDIDURL("did:example:123#key-1")

	keyPair, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	verificationMethod, _ := did.NewVerificationMethod(*keyID, ssi.JsonWebKey2020, did.DID{}, keyPair.Public())

	// This adds the method to the VerificationMethod list and stores a reference to the assertion list
	doc.AddAssertionMethod(verificationMethod)
	//doc.AddAuthenticationMethod(verificationMethod)
	didJson, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println(string(didJson))

	// Unmarshalling of a json did document:
	parsedDIDDoc := did.Document{}
	_ = json.Unmarshal(didJson, &parsedDIDDoc)

	// It can return the key in the convenient lestrrat-go/jwx JWK
	jwk, _ := parsedDIDDoc.AssertionMethod[0].JWK()
	fmt.Println(jwk)
	// Or return a native crypto.PublicKey
	key, _ := parsedDIDDoc.AssertionMethod[0].PublicKey()
	fmt.Println(key)

	credential := vc.VerifiableCredential{}
	json.Unmarshal([]byte(`{
		  "id":"did:example:123#vc-1",
		  "type":["VerifiableCredential", "custom"],
		  "credentialSubject": {"name": "test"}
		}`), &credential)
}
