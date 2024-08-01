package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fabric-did/did"
	"fabric-did/tools"
	"fmt"
	"testing"
	"time"
)

var Issuer = make([]string, 0)

func TestCreatDid(t *testing.T) {
	did.CreateDid()
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
	credentialSubject := did.CredentialSubject{
		ID:    "did:example:123",
		Claim: claim,
	}
	now := time.Now().UnixMilli()
	vc := did.VerifiableCredential{
		ID:                tools.GetUUID(),
		Type:              did.VCType,
		Issuer:            "did:example:061dce1b6e87f6ff110d50cf2ce5bd98",
		IssuanceDate:      now,
		ExpirationDate:    now + 60000,
		CredentialSubject: credentialSubject,
		Proof:             did.Proof{},
	}
	vcBytes, _ := json.Marshal(vc)
	//create issuer
	issuer := did.CreateDid()
	// 进行ed25519签名
	privateKey := tools.Base58Decode([]byte(issuer.PrivateKey))
	signature := ed25519.Sign(privateKey, vcBytes)
	signatureStr := string(tools.Base58Encode(signature))
	proof := did.Proof{
		Creator:   issuer.ID,
		Signature: signatureStr,
		Created:   now,
		Type:      did.KeyType,
	}
	vc.Proof = proof
	fmt.Printf("Create VC success:%#v \n", vc)

	//验证
	pk := tools.Base58Decode([]byte(issuer.Authentication.PublicKeyMultibase))
	vc.Proof = did.Proof{}
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
	id := tools.Base58Encode(publicKey)
	fmt.Println("公钥：", publicKeyStr)
	fmt.Println("id：", string(id))
	fmt.Println("私钥", privateKeyStr)
	fmt.Println("签名", signature)
	fmt.Println("验签结果：", verify)

}
