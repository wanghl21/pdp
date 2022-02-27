

###### fabric2.x运行

1.1 安装预置环境

```shell
cd first-network
```

1.2运行 利用docker镜像快速引导一个由4个代表2个不同组织的peer节点以及一个排序服务节点的Hyperledger fabric网络。它还将启动一个容器来运行一个peer节点加入channel、部署实例化链码服务以及驱动已经部署的链码执行交易的脚本。

```
./byfn.sh -h
Usage:
  byfn.sh -m up|down|restart|generate [-c <channel name>] [-t <timeout>]
  byfn.sh -h|--help (print this message)
    -m <mode> - one of 'up', 'down', 'restart' or 'generate'
      - 'up' - bring up the network with docker-compose up
      - 'down' - bring up the network with docker-compose up
      - 'restart' - bring up the network with docker-compose up
      - 'generate' - generate required certificates and genesis block
    -c <channel name> - config name to use (defaults to "mychannel")
    -t <timeout> - CLI timeout duration in microseconds (defaults to 10000)

Typically, one would first generate the required certificates and
genesis block, then bring up the network. e.g.:

  byfn.sh -m generate -c <channelname>
  byfn.sh -m up -c <channelname>
```


###### 启动网络

```shell
./network.sh up
```

###### 测试网络的组成部分

与Fabric网络交互的每个节点和用户都需要属于作为网络成员的组织即联盟。测试网络有两个联盟成员，Org1和Org2。网络还包括维护网络的排序服务的排序节点组织。

节点是任何Fabric网络的基本组成部分。节点存储区块链账本，并在提交到账本之前验证交易。节点运行智能合约，其中包含用于管理区块链账本上资产的业务逻辑。

排序服务允许节点集中精力验证交易并将其提交到账本。在排序节点从客户端接收到背书的交易后，它们就交易的顺序达成共识，然后将它们添加到块中。然后，区块被分发到普通节点，普通节点将区块添加到区块链账本中。排序节点还操作定义Fabric网络功能的系统通道，例如如何生成块以及节点可以使用的Fabric版本。系统通道定义了哪些组织是联盟的成员。

```
docker ps -a
```

###### 创建通道

通道是特定网络成员之间的专用通信层。通道只能由受邀加入通道的组织使用，并对网络的其他成员不可见。每个通道都有一个单独的区块链账本。创建默认名称为mychannel的通道

```
./network.sh createChannel
```

也可以使用通道标志创建具有自定义名称的通道。例如，以下命令将创建名为channel1的通道。

```
./network.sh createChannel -c channel1
```

如果要启动网络并在一个步骤中创建通道，可以同时使用up和createChannel模式：

```
./network.sh up createChannel
```

###### 在通道上启动链码

创建通道后，可以使用智能合约与通道账本进行交互。智能合约包含管理区块链账本上资产的业务逻辑。网络成员运行的应用程序可以调用智能合约在账本上创建资产，以及更改和转移这些资产。应用程序还可以查询智能合约以读取账本上的数据。

为确保交易有效，使用智能合约创建的交易通常需要由多个组织签署才能提交到通道账本。多重签名是Fabric信任模型不可或缺的一部分。要求对一个交易进行多重背书可以防止通道上的一个组织篡改其节点的账本或使用未经同意的业务逻辑。要签署一个交易，每个组织都需要在其节点上调用和执行智能合约，然后由节点对交易的输出进行签名。如果输出是一致的，并且已经有足够的组织签名，则可以将交易提交到账本。指定通道上需要执行智能合约的已设置组织的策略称为背书策略，它是作为链码定义的一部分为每个链码设置的。

在Fabric中，智能合约在网络上部署在称为链码的包中。链码安装在组织的节点上，然后部署到一个通道，在那里它可以被用来背书交易和与区块链账本交互。在将链码部署到通道之前，通道成员需要就建立链码治理的链码定义达成一致。当所需数量的组织达成一致时，可以将链码定义提交给通道，并且可以使用链码。

需把链码放在/home/blockchain/go/src/github.com/hyperledger/fabric/scripts/fabric-samples目录下

链码在import shim和peer时注意fabric2.x版本

```
"github.com/hyperledger/fabric-chaincode-go/shim"
sc "github.com/hyperledger/fabric-protos-go/peer"
```



```
./network.sh deployCC -ccn car_demo -ccp ../car_demo/chaincode-go -ccl go
./network.sh deployCC -ccn cbc -ccp ../cbc/chaincode-go -ccl go
```

- **-ccn**：为指定链码名称
- **-ccl**：为指定链码语言

###### 利用终端命令行与网络交互

在启动测试网络后，可以使用peer cli客户端与网络进行交互，通过peer cli客户端可以调用已部署的智能合约，更新通道，或者安装和部署新的智能合约。

首先确保操作目录为test-netwok目录

```
pwd
```

执行一下命令将cli客户端添加到环境变量中

```
export PATH=${PWD}/../bin:$PATH
```

还需要将fabric-samples代码库中的FABRIC_CFG_PATH设置为指向其中的core.yaml文件：

```
export FABRIC_CFG_PATH=$PWD/../config/
```

设置允许org1操作peer cli的环境变量：

```bash
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

`CORE_PEER_TLS_ROOTCERT_FILE` 和 `CORE_PEER_MSPCONFIGPATH` 环境变量指向Org1的 `organizations` 文件夹中的的加密材料。

执行以下命令用一些资产来初始化账本： 注意-n 后面参数改为链码名称

```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n car_demo --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"initLedger","Args":[]}'
```

执行以下指令来查询通道账本中的车辆列表：

```bash
peer chaincode query -C mychannel -n car_demo -c '{"Args":["queryAllCars"]}'
```

```
peer chaincode query -C mychannel -n cbc -c '{"Args":["RegisterUser","Alice"]}'

peer chaincode query -C mychannel -n abe_lightweight -c '{"Args":["RegisterUser","Alice"]}'
```

###### 利用application与网络进行交互

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/robfig/cron"
)

func main() {

	err := os.RemoveAll("./keystore")
	err = os.RemoveAll("./wallet")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("============ application-golang starts ============")

	err = os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")//链接到channel
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract := network.GetContract("test") //获得智能合约

	log.Println("--> Submit Transaction: RegisterUser, function creates the initial set of assets on the ledger")

	result, err := contract.SubmitTransaction("RegisterUser"，"Alice") //提交交易
	print(result)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
```

###### 关闭网络

```bash
./network.sh down
```

###### Tips

注意依赖包的版本可能会导致出现错误



项目基本步骤

1.书写chaincode，然后安装在channel上

2.在cloud端上 定时调用chaincode 

3.返回给client相关报告
