**基于区块链和隐私保护的分布式身份认证系统**

## 分布式身份认证系统：

**分为两个部分**：

![图片](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/640)

### 分布式身份标识 (DID)

[Decentralized Identifiers (DIDs) v1.0 (w3.org)](https://www.w3.org/TR/did-core/)

用`DID`来构建每一个**实体的标识符**，用`DID Document`文档来储存身份的元数据，用`DID私钥`来对Jason-LD格式的数字凭证进行**签名**，就可以加密安全，保护隐私和可由第三方进行机器验证的方式在Web上表达现实社会中各种各样的凭证。在区块链、密码学的支持下，这种凭证比物理凭证更加的可信赖。

#### DID身份的构成:

1. did开头表示分布式身份表示,
2. the identifier for the DID method
3. the DID method-specific identifier.

![imgpng](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/0d1d710da7c9b6e8a152dc4b6f7d4673-img.png)

上面的示例 DID 解析为 DID 文档。 DID 文档包含与 DID 关联的信息，例如以加密方式验证 DID 控制器的方法。

#### EXAMPLE 1: A simple DID document

```json
{
"@context": [
"https://www.w3.org/ns/did/v1",
"https://w3id.org/security/suites/ed25519-2020/v1"
],
"id": "did:example:123456789abcdefghi",
    "PK","SK"
"authentication": [{

    "id": "did:example:123456789abcdefghi#keys-1",
    "type": "Ed25519VerificationKey2020",
    "controller": "did:example:123456789abcdefghi",
    "publicKeyMultibase": "zH3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV"
}]
}
```

#### DID架构

<img src="https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/838f1e94d3c65297280d1d95257afef8-img_1.png" alt="img1png" style="zoom:80%;" />



### 可验证数字凭证 (Verifiable Credentials)

[Verifiable Credentials Data Model v1.1 (w3.org)](https://www.w3.org/TR/vc-data-model/)

在**物理世界**中，**凭证**可能包括：

- 与凭证持有者相关的信息（照片、名称或身份证号码）
- 与发行者有关的信息（市政府、国家机构或认证机构）
- 与凭据类型有关的信息（学历、护照）
- 与凭据约束有关的信息（日期或使用条款）

**可验证凭证**（Verifiable Credentials）需要能够表示物理凭据所代表的所有相同信息。所以数字签名被应用到了**可验证凭证**中。

持有**可验证凭证**的人可以与**验证者**分享这些凭证，以证明他们具有特定特征的**可验证凭据**。

<img src="https://www.w3.org/TR/vc-data-model/diagrams/ecosystem.svg" style="zoom:80%;" />

#### 持有者（holder）

**持有者**可以拥有一个或多个可验证凭据。示例：学生，员工和客户等。

#### 发行人（issuer）

**发行人**通过验证主体的一个或多个属性，根据这些属性创建**可验证凭证**，并将**可验证凭证**传输给持有者。示例：公司、非营利组织、行业协会、政府和个人。

#### 验证者（verifies）

验证者通过接收一个或多个可验证的凭据并进行验证。示例：雇主，安全人员和网站。

### DATA Model

<img src="https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/image-20230919161452375.png" alt="image-20230919161452375" style="zoom:80%;" />

EXAMPLE 1: A simple example of a verifiable credential

```
{
  // set the context, which establishes the special terms we will be using
  // such as 'issuer' and 'alumniOf'.
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://www.w3.org/2018/credentials/examples/v1"
  ],
  // specify the identifier for the credential
  "id": "http://example.edu/credentials/1872",
  // the credential types, which declare what data to expect in the credential
  "type": ["VerifiableCredential", "AlumniCredential"],
  // the entity that issued the credential
  "issuer": "https://example.edu/issuers/565049",
  // when the credential was issued
  "issuanceDate": "2010-01-01T19:23:24Z",
  // claims about the subjects of the credential
  "credentialSubject": {
    // identifier for the only subject of the credential
    "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
    // assertion about the only subject of the credential
    "alumniOf": {
      "id": "did:example:c276e12ec21ebfeb1f712ebc6f1",
      "name": [{
        "value": "Example University",
        "lang": "en"
      }, {
        "value": "Exemple d'Université",
        "lang": "fr"
      }]
    }
  },
  // digital proof that makes the credential tamper-evident
  // see the NOTE at end of this section for more detail
  "proof": {
    // the cryptographic signature suite that was used to generate the signature
    "type": "RsaSignature2018",
    // the date the signature was created
    "created": "2017-06-18T21:19:10Z",
    // purpose of this proof
    "proofPurpose": "assertionMethod",
    // the identifier of the public key that can verify the signature
    "verificationMethod": "https://example.edu/issuers/565049#key-1",
    // the digital signature value
    "jws": "eyJhbGciOiJSUzI1NiIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..TCYt5X
      sITJX1CxPCT8yAV-TVkIEq_PbChOMqsLfRoPsnsgw5WEuts01mq-pQy7UJiN5mgRxD-WUc
      X16dUEMGlv50aqzpqh4Qktb3rk-BuQy72IFLOqV0G_zS245-kronKb78cPN25DGlcTwLtj
      PAYuNzVBAh4vGHSrQyHUdBBPM"
  }
}
```
