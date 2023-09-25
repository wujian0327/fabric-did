package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"math/big"
)

type SimpleChaincode struct {
	contractapi.Contract
}

const KeyType = "Ed25519VerificationKey2018"
const VCType = "VerifiableCredential"
const IssuerKey = "ISSUER_LIST"

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

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

func (t *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Init Ledger")
	var issuer = make([]string, 0)
	issuerBytes, err := json.Marshal(issuer)
	if err != nil {
		return err
	}
	//初始化issuerList列表
	err = ctx.GetStub().PutState(IssuerKey, issuerBytes)
	if err != nil {
		return err
	}
	return nil
}

func (t *SimpleChaincode) PutValue(ctx contractapi.TransactionContextInterface, key string, value string) error {
	err := ctx.GetStub().PutState(key, []byte(value))
	fmt.Printf("put value success,key:%v,value:%v", key, value)
	return err
}

func (t *SimpleChaincode) GetValue(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	b, err := ctx.GetStub().GetState(key)
	if b == nil {
		return "", fmt.Errorf("key doesn't exist")
	}
	return string(b), err
}

// CreateDID 如果不提供公私钥，那么CreateDID应该是链下处理，数据上链
func (t *SimpleChaincode) CreateDID(ctx contractapi.TransactionContextInterface, DIDJson string) error {
	var did DID
	err := json.Unmarshal([]byte(DIDJson), &did)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(did.ID, []byte(DIDJson))
	if err != nil {
		return err
	}
	return nil
}

func (t *SimpleChaincode) SearchDID(ctx contractapi.TransactionContextInterface, ID string) (string, error) {
	didBytes, err := ctx.GetStub().GetState(ID)
	if didBytes == nil || err != nil {
		return "", fmt.Errorf("DID doesn't exist")
	}
	return string(didBytes), nil
}

func (t *SimpleChaincode) AddIssuer(ctx contractapi.TransactionContextInterface, ID string) error {
	b, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return fmt.Errorf("DID doesn't exist")
	}
	var did DID
	err = json.Unmarshal(b, &did)
	if err != nil {
		return err
	}
	issuerList, err := GetIssuerList(ctx)
	issuerList = append(issuerList, ID)
	issuerBytes, err := json.Marshal(issuerList)
	if err != nil {
		return err
	}
	//更新issuer
	err = ctx.GetStub().PutState(IssuerKey, issuerBytes)
	if err != nil {
		return err
	}
	return nil
}

func (t *SimpleChaincode) SearchIssuer(ctx contractapi.TransactionContextInterface) (string, error) {
	issuerBytes, err := ctx.GetStub().GetState(IssuerKey)
	if err != nil {
		return "", fmt.Errorf("issuer doesn't exist")
	}
	return string(issuerBytes), nil
}

// CreateVC VC的创建需要链上和链下一起进行，链下进行issuer私钥的签名，链上判断issuer的权限
func (t *SimpleChaincode) CreateVC(ctx contractapi.TransactionContextInterface, ID string, issuerID string, vcJson string) error {
	//检查issuer是否存在
	_, err := CheckIssuer(ctx, issuerID)
	if err != nil {
		return err
	}
	//查询ID是否存在
	_, err = GetDID(ctx, ID)
	if err != nil {
		return err
	}
	var vc VerifiableCredential
	err = json.Unmarshal([]byte(vcJson), &vc)
	if err != nil {
		return fmt.Errorf("verifiableCredential format error")
	}
	if vc.Issuer != issuerID || issuerID != vc.Proof.Creator {
		return fmt.Errorf("issuerID error")
	}
	err = ctx.GetStub().PutState(vc.ID, []byte(vcJson))
	return nil
}

func (t *SimpleChaincode) SearchVC(ctx contractapi.TransactionContextInterface, vcID string) (string, error) {
	vcBytes, err := ctx.GetStub().GetState(vcID)
	if vcBytes == nil || err != nil {
		return "", fmt.Errorf("vc doesn't exist")
	}
	return string(vcBytes), nil
}

// CheckVC 检查VC可以链上也可以链下，我们这里链上处理
func (t *SimpleChaincode) CheckVC(ctx contractapi.TransactionContextInterface, vcID string) (bool, error) {
	vcBytes, err := ctx.GetStub().GetState(vcID)
	var vc VerifiableCredential
	err = json.Unmarshal(vcBytes, &vc)
	if err != nil {
		return false, err
	}
	//检查issuer是否存在
	_, err = CheckIssuer(ctx, vc.Issuer)
	if err != nil {
		return false, err
	}
	//获得issuer的公钥
	issuerDID, err := GetDID(ctx, vc.Issuer)
	if err != nil {
		return false, err
	}
	//验证
	signature := Base58Decode([]byte(vc.Proof.Signature))
	pk := Base58Decode([]byte(issuerDID.Authentication.PublicKeyMultibase))
	vc.Proof = Proof{}
	vcBytes2, _ := json.Marshal(vc)
	verify := ed25519.Verify(pk, vcBytes2, signature)
	fmt.Println("verify:", verify)
	return verify, nil
}

func CheckIssuer(ctx contractapi.TransactionContextInterface, issuerID string) (bool, error) {
	issuerList, err := GetIssuerList(ctx)
	if err != nil {
		return false, fmt.Errorf("issuer doesn't exist")
	}
	for _, eachItem := range issuerList {
		if eachItem == issuerID {
			return true, nil
		}
	}
	return false, fmt.Errorf("not issuer")
}

func GetIssuerList(ctx contractapi.TransactionContextInterface) ([]string, error) {
	issuerBytes, err := ctx.GetStub().GetState(IssuerKey)
	if issuerBytes == nil || err != nil {
		return nil, fmt.Errorf("issuer doesn't exist")
	}
	var issuerList = make([]string, 0)
	err = json.Unmarshal(issuerBytes, &issuerList)
	if err != nil {
		return nil, fmt.Errorf("unmarshal issuer error")
	}
	return issuerList, nil
}

func GetDID(ctx contractapi.TransactionContextInterface, ID string) (DID, error) {
	didBytes, err := ctx.GetStub().GetState(ID)
	if didBytes == nil || err != nil {
		return DID{}, fmt.Errorf("DID doesn't exist")
	}
	var did DID
	err = json.Unmarshal(didBytes, &did)
	if err != nil {
		return DID{}, err
	}
	return did, nil
}

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	if input[0] == 0x00 {
		result = append(result, b58Alphabet[0])
	}

	ReverseBytes(result)

	return result
}

// Base58Decode decodes Base58-encoded data
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)

	for _, b := range input {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()

	if input[0] == b58Alphabet[0] {
		decoded = append([]byte{0x00}, decoded...)
	}

	return decoded
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
