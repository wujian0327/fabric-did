package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fabric-did/did"
	"fabric-did/tools"
	"fmt"
	"testing"
	"time"
)

const chaincode = "did"

// issuer did:example:9eb940032acc12eb2a8637e3459aa605
//did:example:388fa399ee7acd6435c5ac62f5344755
func TestCreateDID(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	publicKeyBase58 := tools.Base58Encode(publicKey)
	privateKeyBase58 := tools.Base58Encode(privateKey)
	id := "did:example:" + tools.GetUUID()
	authentication := did.Authentication{
		ID:                 id + "#keys-1",
		Type:               did.KeyType,
		Controller:         id,
		PublicKeyMultibase: string(publicKeyBase58),
	}
	DID := did.DID{
		ID:             id,
		Authentication: authentication,
		PrivateKey:     string(privateKeyBase58),
	}
	didBytes, _ := json.Marshal(DID)
	_, _, err := ExecuteChaincode(chaincode, "CreateDID", string(didBytes))
	if err != nil {
		panic(err)
	}
	fmt.Printf("id: %s \n", id)
}

func TestSearchDID(t *testing.T) {
	result, err := QueryChaincode(chaincode, "SearchDID", "did:example:9eb940032acc12eb2a8637e3459aa605")
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %s \n", result)
}

func TestAddIssuer(t *testing.T) {
	_, _, err := ExecuteChaincode(chaincode, "AddIssuer", "did:example:9eb940032acc12eb2a8637e3459aa605")
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %v \n", true)
}

func TestSearchIssuer(t *testing.T) {
	result, err := QueryChaincode(chaincode, "SearchIssuer")
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %s \n", result)
}

func TestCreateVC(t *testing.T) {
	claim := make(map[string]interface{})
	claim["name"] = "李白"
	claim["age"] = "1300"
	claim["poiet"] = "桃花潭水深千尺，不及汪伦送我情"
	now := time.Now().UnixMilli()
	credentialSubject := did.CredentialSubject{
		ID:    "did:example:388fa399ee7acd6435c5ac62f5344755",
		Claim: claim,
	}
	issuerId := "did:example:9eb940032acc12eb2a8637e3459aa605"
	issuerBytes, err := QueryChaincode(chaincode, "SearchDID", issuerId)
	if err != nil {
		panic(err)
	}
	var issuer did.DID
	err = json.Unmarshal(issuerBytes, &issuer)
	if err != nil {
		panic(err)
	}
	vc := did.VerifiableCredential{
		ID:                tools.GetUUID(),
		Type:              did.VCType,
		Issuer:            issuer.ID,
		IssuanceDate:      now,
		ExpirationDate:    now + 60000,
		CredentialSubject: credentialSubject,
		Proof:             did.Proof{},
	}
	vcBytes, _ := json.Marshal(vc)
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
	vcBytes2, _ := json.Marshal(vc)
	_, _, err = ExecuteChaincode(chaincode, "CreateVC", "did:example:388fa399ee7acd6435c5ac62f5344755", "did:example:9eb940032acc12eb2a8637e3459aa605", string(vcBytes2))
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %v \n", true)
	fmt.Printf("vc: %v \n", string(vcBytes2))
	//e8b13fc0d78d4eecaabbc9e19bd52b1b
}

func TestSearchVC(t *testing.T) {
	result, err := QueryChaincode(chaincode, "SearchVC", "e8b13fc0d78d4eecaabbc9e19bd52b1b")
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %s \n", result)
}

func TestCheckVC(t *testing.T) {
	result, err := QueryChaincode(chaincode, "CheckVC", "e8b13fc0d78d4eecaabbc9e19bd52b1b")
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %s \n", result)
}

func TestCreateFakeVC(t *testing.T) {
	fakeVCJson := "{\"id\":\"e8b13fc0d78d4eecaabbc9e19bd52b1b\",\"type\":\"VerifiableCredential\",\"issuer\":\"did:example:9eb940032acc12eb2a8637e3459aa605\",\"issuanceDate\":1695636155254,\"expirationDate\":1695636215254,\"credentialSubject\":{\"id\":\"did:example:388fa399ee7acd6435c5ac62f5344755\",\"claim\":{\"age\":\"1300\",\"name\":\"李白\",\"poiet\":\"桃花潭水深千尺，不及汪伦送我情\"}},\"proof\":{\"creator\":\"did:example:9eb940032acc12eb2a8637e3459aa605\",\"signature\":\"CqpGwPKqF73G5berg5Ti6U9ThnPq3eonUFNannW4t6HHcmB8N8s4JQa7sSbFVkzxufWo9RQvSGPqjwEVzanNmki\",\"created\":1695636155254,\"type\":\"Ed25519VerificationKey2018\"}}"
	var vc did.VerifiableCredential
	_ = json.Unmarshal([]byte(fakeVCJson), &vc)
	vc.CredentialSubject.Claim["name"] = "杜甫"
	bytes, _ := json.Marshal(vc)
	_, _, err := ExecuteChaincode(chaincode, "PutValue", "e8b13fc0d78d4eecaabbc9e19bd52b1b", string(bytes))
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %v \n", true)
}
