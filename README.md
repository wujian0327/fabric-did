# 基于区块链和隐私保护的分布式身份认证系统

## 测试环境
ubuntu 22.04

golang 1.18

fabric 2.4.3

**因为fabric-sdk-go v1.0.0只支持到fabric2.4，golang的版本也不能太高，只推荐1.18**

**后续可能会使用Fabric Gateway取代fabric-sdk-go以支持最新的fabric**

**但目前暂时继续使用fabric2.4**

## 1. 本地DID测试
执行
```shell
go mod tidy
```
运行本地DID测试
```shell
go test -v local_did_test.go
```

![image-20240801133603458](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/image-20240801133603458.png)

## 2. fabric2.4 环境测试搭建

### 2.1 克隆fabric-sample仓库

```bash
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-sample
git checkout v2.4.3
```

### 2.2 下载二进制文件

```bash
wget https://github.com/hyperledger/fabric/releases/download/v2.4.3/hyperledger-fabric-linux-amd64-2.4.3.tar.gz
wget https://github.com/hyperledger/fabric-ca/releases/download/v1.5.3/hyperledger-fabric-ca-linux-amd64-1.5.3.tar.gz
tar xvzf hyperledger-fabric-linux-amd64-2.4.3.tar.gz
tar xvzf hyperledger-fabric-ca-linux-amd64-1.5.3.tar.gz
cd bin
export PATH=${PWD}:$PATH
cd ../
```

### 2.3 下载docker镜像

```bash
docker pull hyperledger/fabric-peer:2.4.3
docker pull hyperledger/fabric-orderer:2.4.3
docker pull hyperledger/fabric-ccenv:2.4.3
docker pull hyperledger/fabric-tools:2.4.3
docker pull hyperledger/fabric-baseos:2.4.3
docker pull hyperledger/fabric-ca:1.5.3

docker tag hyperledger/fabric-peer:2.4.3 hyperledger/fabric-peer
docker tag hyperledger/fabric-peer:2.4.3 hyperledger/fabric-peer:2.4
docker tag hyperledger/fabric-orderer:2.4.3 hyperledger/fabric-orderer
docker tag hyperledger/fabric-orderer:2.4.3 hyperledger/fabric-orderer:2.4
docker tag hyperledger/fabric-ccenv:2.4.3 hyperledger/fabric-ccenv
docker tag hyperledger/fabric-ccenv:2.4.3 hyperledger/fabric-ccenv:2.4
docker tag hyperledger/fabric-tools:2.4.3 hyperledger/fabric-tools
docker tag hyperledger/fabric-tools:2.4.3 hyperledger/fabric-tools:2.4
docker tag hyperledger/fabric-baseos:2.4.3 hyperledger/fabric-baseos
docker tag hyperledger/fabric-baseos:2.4.3 hyperledger/fabric-baseos:2.4
docker tag hyperledger/fabric-ca:1.5.3 hyperledger/fabric-ca
docker tag hyperledger/fabric-ca:1.5.3 hyperledger/fabric-ca:1.5
```

### 2.3  启动fabric测试网络

```bash
cd test-network
./network.sh down
./network.sh up
```

### 2.4 创建通道

```bash
./network.sh createChannel
```

### 2.5 部署测试链码

```
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
```

## 3. 部署DID链码到fabric

### 3.1 复制fabric-did目录下的contract到test-network

```bash
cp -r {repo_path}/contract .
cd contract/
go mod tidy
go mod vendor
cd ..
```

### 3.2 安装DID链码

```bash
./network.sh deployCC -ccn did -ccp contract -ccl go
```

![image-20240801143741810](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/image-20240801143741810.png)

## 4. 测试fabric-sdk的功能

### 4.1 复制test-network/organizations/*到fabric-did/config目录下

```
cd fabric-did/config
cp -r {fabric-sample-path}/test-network/organizations/* .
cd ../
```

![image-20240801145017359](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/image-20240801145017359.png)

### 4.2 运行测试

```
go test -v  -run=TestGetBlockData
```

![image-20240801144608854](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/image-20240801144608854.png)

如果遇到网络问题，可能需要设置host

```
{fabric_ip}     peer0.org1.example.com
{fabric_ip}     peer0.org2.example.com
{fabric_ip}     orderer.example.com
```

## 5. 测试fabric_did的功能

```
go test -v  -run=TestCreateDID 
```

![image-20240801145235567](https://gitee.com/wujian2023/typora_images/raw/master/auto_upload/image-20240801145235567.png)



上述只是实现了did和vc的链上存储和验证，隐私保护、身份认证等功能敬请期待。