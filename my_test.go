package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const KeyType = "Ed25519VerificationKey2018"
const VCType = "VerifiableCredential"

type DID struct {
	ID             string         `json:"id"`
	Authentication Authentication `json:"authentication,omitempty"`
	privateKey     []byte
}

type Authentication struct {
	ID                 string `json:"id"`
	Type               string `json:"type,omitempty"`
	Controller         string `json:"controller,omitempty"`
	PublicKeyMultibase string `json:"publicKeyMultibase,omitempty"`
}

type VerifiableCredential struct {
	ID             string                 `json:"id,omitempty"`
	Type           string                 `json:"type"`
	Issuer         string                 `json:"issuer"`
	IssuanceDate   int64                  `json:"issuanceDate"`
	ExpirationDate int64                  `json:"expirationDate,omitempty"`
	Claim          map[string]interface{} `json:"claim"`
	Proof          Proof                  `json:"proof"`
}

type Proof struct {
	Creator   string `json:"creator,omitempty"`
	Signature string `json:"signature"`
	Created   int64  `json:"created,omitempty"`
	Type      string `json:"type,omitempty"`
}

var Issuer = make([]string, 0)

func TestCreatDid(t *testing.T) {
	CreateDid()
}

func CreateDid() DID {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	publicKeyBase58 := Base58Encode(publicKey)
	privateKeyBase58 := Base58Encode(privateKey)
	id := "did:example:" + GetUUID()
	authentication := Authentication{
		ID:                 id + "#keys-1",
		Type:               KeyType,
		Controller:         id,
		PublicKeyMultibase: string(publicKeyBase58),
	}
	did := DID{
		ID:             id,
		Authentication: authentication,
		privateKey:     privateKey,
	}
	fmt.Printf("did:%#v \n", did)
	println(privateKeyBase58)
	return did
}

func TestIssuerAdd(t *testing.T) {
	Issuer = append(Issuer, "did:example:061dce1b6e87f6ff110d50cf2ce5bd98")
	fmt.Printf("Issuer:%#v \n", Issuer)
}

func TestVC(t *testing.T) {
	claim := make(map[string]interface{})
	claim["name"] = "李白"
	claim["age"] = "1300"
	claim["poiet"] = "桃花潭水深千尺，不及汪伦送我情"
	now := time.Now().UnixMilli()
	vc := VerifiableCredential{
		ID:             GetUUID(),
		Type:           VCType,
		Issuer:         "did:example:061dce1b6e87f6ff110d50cf2ce5bd98",
		IssuanceDate:   now,
		ExpirationDate: now + 60000,
		Claim:          claim,
		Proof:          Proof{},
	}
	vcBytes, _ := json.Marshal(vc)
	//create issuer
	issuer := CreateDid()
	// 进行ed25519签名
	signature := ed25519.Sign(issuer.privateKey, vcBytes)
	signatureStr := string(Base58Encode(signature))
	proof := Proof{
		Creator:   issuer.ID,
		Signature: signatureStr,
		Created:   now,
		Type:      KeyType,
	}
	vc.Proof = proof
	fmt.Printf("Create VC success:%#v \n", vc)

	//验证
	pk := Base58Decode([]byte(issuer.Authentication.PublicKeyMultibase))
	vc.Proof = Proof{}
	vcBytes2, _ := json.Marshal(vc)
	verify := ed25519.Verify(pk, vcBytes2, signature)
	fmt.Println("验签结果：", verify)
}

func TestEd25519(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)

	msg := "abc123"

	msgByte := []byte(msg)

	// 进行ed25519签名
	signature := ed25519.Sign(privateKey, msgByte)

	// 使用公钥进行验签
	verify := ed25519.Verify(publicKey, msgByte, signature)

	publicKeyStr := base64.StdEncoding.EncodeToString(publicKey)
	privateKeyStr := base64.StdEncoding.EncodeToString(privateKey)
	id := Base58Encode(publicKey)
	fmt.Println("公钥：", publicKeyStr)
	fmt.Println("id：", string(id))
	fmt.Println("私钥", privateKeyStr)
	fmt.Println("签名", signature)
	fmt.Println("验签结果：", verify)

}
