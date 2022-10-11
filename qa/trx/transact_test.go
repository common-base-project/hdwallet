package trx

import (
	"fmt"
	"github.com/lizc2003/hdwallet/trx"
	"github.com/lizc2003/hdwallet/wallet"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

/*
install:
docker pull trontools/quickstart

run:
docker run -d -it -p 9090:9090 -p 50051:50051 -p 50052:50052 --rm --name tron \
	-e "mnemonic=purse cheese cage reason cost flat jump usage hospital grit delay loan" \
	-e "defaultBalance=1000000" -e "formatJson=true" \
	trontools/quickstart

check:
curl "http://127.0.0.1:9090/admin/accounts?format=all"
*/

func TestTransaction(t *testing.T) {
	// TestTRC20Abi compile https://shasta.tronscan.org/?#/contracts/contract-compiler
	const TestTRC20Abi = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_value","type":"uint256"}],"name":"burn","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_value","type":"uint256"}],"name":"burnFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"},{"name":"_extraData","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]`
	// TestTRC20ContractBytecode base on trc20.sol
	const TestTRC20ContractBytecode = "60806040526002805460ff1916600317905534801561001d57600080fd5b50d3801561002a57600080fd5b50d2801561003757600080fd5b5060025460ff16600a0a6103e80260038190553360009081526004602081815260408084209490945583518085019094528184527f5454524300000000000000000000000000000000000000000000000000000000930192835261009b92906100e6565b506040805180820190915260048082527f545452430000000000000000000000000000000000000000000000000000000060209092019182526100e0916001916100e6565b50610181565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061012757805160ff1916838001178555610154565b82800160010185558215610154579182015b82811115610154578251825591602001919060010190610139565b50610160929150610164565b5090565b61017e91905b80821115610160576000815560010161016a565b90565b610990806101906000396000f3006080604052600436106100a05763ffffffff60e060020a60003504166306fdde0381146100a5578063095ea7b31461014957806318160ddd1461019b57806323b872dd146101dc578063313ce5671461022057806342966c681461026557806370a082311461029757806379cc6790146102d257806395d89b4114610310578063a9059cbb1461033f578063cae9ca511461037f578063dd62ed3e14610402575b600080fd5b3480156100b157600080fd5b50d380156100be57600080fd5b50d280156100cb57600080fd5b506100d4610443565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561010e5781810151838201526020016100f6565b50505050905090810190601f16801561013b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015557600080fd5b50d3801561016257600080fd5b50d2801561016f57600080fd5b50610187600160a060020a03600435166024356104d1565b604080519115158252519081900360200190f35b3480156101a757600080fd5b50d380156101b457600080fd5b50d280156101c157600080fd5b506101ca6104fe565b60408051918252519081900360200190f35b3480156101e857600080fd5b50d380156101f557600080fd5b50d2801561020257600080fd5b50610187600160a060020a0360043581169060243516604435610504565b34801561022c57600080fd5b50d3801561023957600080fd5b50d2801561024657600080fd5b5061024f610573565b6040805160ff9092168252519081900360200190f35b34801561027157600080fd5b50d3801561027e57600080fd5b50d2801561028b57600080fd5b5061018760043561057c565b3480156102a357600080fd5b50d380156102b057600080fd5b50d280156102bd57600080fd5b506101ca600160a060020a03600435166105e2565b3480156102de57600080fd5b50d380156102eb57600080fd5b50d280156102f857600080fd5b50610187600160a060020a03600435166024356105f4565b34801561031c57600080fd5b50d3801561032957600080fd5b50d2801561033657600080fd5b506100d46106b3565b34801561034b57600080fd5b50d3801561035857600080fd5b50d2801561036557600080fd5b5061037d600160a060020a036004351660243561070d565b005b34801561038b57600080fd5b50d3801561039857600080fd5b50d280156103a557600080fd5b50604080516020600460443581810135601f8101849004840285018401909552848452610187948235600160a060020a031694602480359536959460649492019190819084018382808284375094975061071c9650505050505050565b34801561040e57600080fd5b50d3801561041b57600080fd5b50d2801561042857600080fd5b506101ca600160a060020a036004358116906024351661081f565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156104c95780601f1061049e576101008083540402835291602001916104c9565b820191906000526020600020905b8154815290600101906020018083116104ac57829003601f168201915b505050505081565b336000908152600560209081526040808320600160a060020a039590951683529390529190912055600190565b60035481565b600160a060020a038316600090815260056020908152604080832033845290915281205482111561053457600080fd5b600160a060020a038416600090815260056020908152604080832033845290915290208054839003905561056984848461083c565b5060019392505050565b60025460ff1681565b3360009081526004602052604081205482111561059857600080fd5b3360008181526004602090815260409182902080548690039055600380548690039055815185815291516000805160206109458339815191529281900390910190a2506001919050565b60046020526000908152604090205481565b600160a060020a03821660009081526004602052604081205482111561061957600080fd5b600160a060020a038316600090815260056020908152604080832033845290915290205482111561064957600080fd5b600160a060020a0383166000818152600460209081526040808320805487900390556005825280832033845282529182902080548690039055600380548690039055815185815291516000805160206109458339815191529281900390910190a250600192915050565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156104c95780601f1061049e576101008083540402835291602001916104c9565b61071833838361083c565b5050565b60008361072981856104d1565b156108175760405160e060020a638f4ffcb10281523360048201818152602483018790523060448401819052608060648501908152875160848601528751600160a060020a03871695638f4ffcb195948b94938b939192909160a490910190602085019080838360005b838110156107ab578181015183820152602001610793565b50505050905090810190601f1680156107d85780820380516001836020036101000a031916815260200191505b5095505050505050600060405180830381600087803b1580156107fa57600080fd5b505af115801561080e573d6000803e3d6000fd5b50505050600191505b509392505050565b600560209081526000928352604080842090915290825290205481565b6000600160a060020a038316151561085357600080fd5b600160a060020a03841660009081526004602052604090205482111561087857600080fd5b600160a060020a038316600090815260046020526040902054828101101561089f57600080fd5b50600160a060020a038083166000818152600460209081526040808320805495891680855282852080548981039091559486905281548801909155815187815291519390950194927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a3600160a060020a0380841660009081526004602052604080822054928716825290205401811461093e57fe5b505050505600cc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5a165627a7a723058204e5a4088448d9024bd87ed8cc6e598a3aa522a31093e400297459449f5f59d810029"

	waitTime := 3 * time.Second

	mnemonic := "purse cheese cage reason cost flat jump usage hospital grit delay loan"
	hdw, err := wallet.NewHDWallet(mnemonic, "", wallet.BtcChainMainNet, wallet.ChainMainNet)
	require.NoError(t, err)

	wtmp, err := hdw.NewWallet(wallet.SymbolTrx, 0, 0, 0)
	require.NoError(t, err)
	w := wtmp.(*wallet.TrxWallet)
	fmt.Println(w.DeriveAddress())

	//client, err := trx.NewTrxClient("grpc.trongrid.io:50051", "ac0823c5-adf0-4e8d-b2b6-322036a3b6e5")
	//client, err := trx.NewTrxClient("grpc.shasta.trongrid.io:50051", "")
	client, err := trx.NewTrxClient("127.0.0.1:50051", "")
	require.NoError(t, err)

	height, err := client.GetBlockHeight()
	require.NoError(t, err)
	fmt.Println("block height:", height)

	acct, err := client.RpcClient.GetAccount(w.DeriveAddress())
	require.NoError(t, err)
	fmt.Println("Trx balance", acct.Balance, trx.SunToTrx(acct.Balance))
	fmt.Println("Energy usage", acct.AccountResource.EnergyUsage)

	{ // Trx transfer
		amount := int64(12000)
		fmt.Println("------------- transfer trx:", amount)
		_, err := trx.TransferTrx(w, client.RpcClient, "TUU9bqm9CCA1dAU9iaa6HcJF4twMBi5N86", amount)
		require.NoError(t, err)
		time.Sleep(waitTime)

		acct2, err := client.RpcClient.GetAccount(w.DeriveAddress())
		require.NoError(t, err)
		fmt.Println("Trx balance after transfer:", acct2.Balance, trx.SunToTrx(acct2.Balance))
		require.Greater(t, acct.Balance, acct2.Balance)
	}

	if acct.AccountResource.EnergyUsage < 100 {
		fmt.Println("------------- freeze balance for energy")
		_, err := trx.FreezeEnergyBalance(w, client.RpcClient, "", trx.TrxToSun(500))
		require.NoError(t, err)
		time.Sleep(waitTime)
	}

	var trc20Addr string
	var baseTrc20Supply *big.Int
	{ // deploy trc20 contract
		fmt.Println("------------- deploy trc20 token")
		txId, err := trx.DeployContract(w, client.RpcClient, "trx20test",
			TestTRC20Abi, TestTRC20ContractBytecode, 1e9, 20, 9e10)
		require.NoError(t, err)

		time.Sleep(waitTime)
		trc20Addr, err = trx.GetContractAddress(client.RpcClient, txId)
		require.NoError(t, err)
		fmt.Println("trc20 contract address:", trc20Addr)

		baseTrc20Supply, err = trx.NewErc20Contract(trc20Addr, client.RpcClient).BalanceOf(w.DeriveAddress())
		require.NoError(t, err)
		fmt.Println("baseTrc20Supply:", baseTrc20Supply.String())
	}

	{ // trc20 transfer
		fmt.Println("------------- transfer trc20 token")
		to := "TSGYZ3VAVsa2SgoYhV7mfqnde89zTU7zNh"
		sendAmount := big.NewInt(50)
		contract := trx.NewErc20Contract(trc20Addr, client.RpcClient)
		fmt.Print("Contract name: ")
		fmt.Println(contract.Name())
		fmt.Print("Contract symbol: ")
		fmt.Println(contract.Symbol())
		fmt.Print("Contract decimals: ")
		fmt.Println(contract.Decimals())

		txId, err := contract.Transfer(w, to, sendAmount, 999999)
		require.NoError(t, err)

		time.Sleep(waitTime)
		txinfo, err := client.RpcClient.GetTransactionInfoByID(txId)
		require.NoError(t, err)

		fmt.Println("Result:", txinfo.Result)
		fmt.Println("Block number:", txinfo.BlockNumber)
		fmt.Println("Fee:", txinfo.Fee)
		fmt.Println("Energy:", txinfo.Receipt.EnergyUsageTotal)
		fmt.Println("Receipt result:", txinfo.Receipt.Result)

		bal, err := contract.BalanceOf(w.DeriveAddress())
		require.NoError(t, err)
		fmt.Println("bal", bal.String())
		require.Equal(t, big.NewInt(0).Sub(baseTrc20Supply, sendAmount).Text(10), bal.Text(10))
	}
}
