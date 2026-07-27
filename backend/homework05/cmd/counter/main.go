package main

import (
	"context"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"homework05/contract"
)

func main() {
	rpcURL := os.Getenv("ETH_RPC_URL")
	pkHex := os.Getenv("ETH_PRIVATE_KEY")
	if rpcURL == "" || pkHex == "" {
		log.Fatal("ETH_RPC_URL / ETH_PRIVATE_KEY 未设置")
	}

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("获取 chainID 失败: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	// 造签名器（带 chainID）
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("创建签名器失败: %v", err)
	}
	instance := getDeployCounterInstance(ctx, client, auth)
	// instance := getCallContractInstance(client, "0x7Ae48fb752A819147b2b7c30f60711DAfD56C5fb")
	count, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("读取 count 失败: %v", err)
	}
	log.Printf("初始 count = %s", count.String())

	// 3. 调 increment（发交易，花 gas）
	incTx, err := instance.Increment(auth)
	if err != nil {
		log.Fatalf("increment 失败: %v", err)
	}
	log.Printf("increment 交易: %s", incTx.Hash().Hex())
	if _, err := bind.WaitMined(ctx, client, incTx); err != nil {
		log.Fatalf("等待 increment 上链失败: %v", err)
	}

	// 4. 再读，验证 +1
	count, err = instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("读取 count 失败: %v", err)
	}
	log.Printf("increment 后 count = %s", count.String())
}

func getDeployCounterInstance(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts) *contract.Counter {
	// 1. 部署 Counter 合约
	address, deployTx, instance, err := contract.DeployCounter(auth, client)
	if err != nil {
		log.Fatalf("部署失败: %v", err)
	}
	log.Printf("合约地址: %s | 部署交易: %s", address.Hex(), deployTx.Hash().Hex())

	// 等部署上链（合约要被打包后才能调用）
	if _, err := bind.WaitDeployed(ctx, client, deployTx); err != nil {
		log.Fatalf("等待部署上链失败: %v", err)
	}
	log.Println("合约已部署上链")
	return instance
}

func getCallContractInstance(client *ethclient.Client, address string) *contract.Counter {
	ctrAddress := common.HexToAddress(address)
	instance, err := contract.NewCounter(ctrAddress, client)
	if err != nil {
		log.Fatalf("创建合约实例失败: %v", err)
	}
	return instance
}
