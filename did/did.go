package did

import (
	"crypto/ed25519"
	"crypto/rand"
	"fabric-did/tools"
	"fmt"
)

const KeyType = "Ed25519VerificationKey2018"
const VCType = "VerifiableCredential"

type DID struct {
	ID             string         `json:"id"`
	Authentication Authentication `json:"authentication,omitempty"`
	PrivateKey     string         `json:"privateKey,omitempty"`
}

type Authentication struct {
	ID                 string `json:"id"`
	Type               string `json:"type,omitempty"`
	Controller         string `json:"controller,omitempty"`
	PublicKeyMultibase string `json:"publicKeyMultibase,omitempty"`
}

type VerifiableCredential struct {
	ID                string            `json:"id,omitempty"`
	Type              string            `json:"type"`
	Issuer            string            `json:"issuer"`
	IssuanceDate      int64             `json:"issuanceDate"`
	ExpirationDate    int64             `json:"expirationDate,omitempty"`
	CredentialSubject CredentialSubject `json:"credentialSubject"`
	Proof             Proof             `json:"proof"`
}

type CredentialSubject struct {
	ID    string                 `json:"id"`
	Claim map[string]interface{} `json:"claim"`
}

type Proof struct {
	Creator   string `json:"creator,omitempty"`
	Signature string `json:"signature"`
	Created   int64  `json:"created,omitempty"`
	Type      string `json:"type,omitempty"`
}

func CreateDid() DID {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	publicKeyBase58 := tools.Base58Encode(publicKey)
	privateKeyBase58 := tools.Base58Encode(privateKey)
	id := "did:example:" + tools.GetUUID()
	authentication := Authentication{
		ID:                 id + "#keys-1",
		Type:               KeyType,
		Controller:         id,
		PublicKeyMultibase: string(publicKeyBase58),
	}
	did := DID{
		ID:             id,
		Authentication: authentication,
		PrivateKey:     string(privateKeyBase58),
	}
	fmt.Printf("did:%#v \n", did)
	return did
}
