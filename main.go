package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strconv"
	cons "uniswap-simulator/lib/constants"
	ppool "uniswap-simulator/lib/pool"
	strat "uniswap-simulator/lib/strategy"
	ent "uniswap-simulator/lib/transaction"
	ui "uniswap-simulator/uint256"
)

func main() {
	fmt.Println("Start")
	transactions := getTransactions()
	fmt.Println("Transactions: ", len(transactions))
	token0 := "USDC"
	token1 := "WETH"
	fee := 500
	sqrtX96big, _ := new(big.Int).SetString("1350174849792634181862360983626536", 10)
	sqrtX96, _ := ui.FromBig(sqrtX96big)

	pool := ppool.NewPool(token0, token1, fee, sqrtX96)

	startAmount0 := ui.NewInt(1_000_000)                                // 1 USDC
	startAmount1big, _ := new(big.Int).SetString("290000000000000", 10) // 290_000_000_000_000 wei ~= 1 USD worth of ETH
	startAmount1, _ := ui.FromBig(startAmount1big)

	strategy := strat.NewStrategy(startAmount0, startAmount1, pool, 4000)
	fmt.Printf("AmountBefore: %d %d \n", strategy.Amount0, strategy.Amount1)
	//
	//starttime := transactions[0].Timestamp
	//// 30 days
	//nextUpdate := starttime + (60 * 60 * 24 * 30)
	//fmt.Printf("NextUpdate: %d\n", nextUpdate)
	//// 24 hours
	//updateInterval := 60 * 60 * 24
	strategy.Rebalance()

	//amountbig , _ := new(big.Int).SetString("93924580278", 10)
	//amount, _ := ui.FromBig(amountbig)
	//strategy.Pool.MintStrategy(190880, 198880, amount)

	for _, trans := range transactions {
		switch trans.Type {
		case "Mint":
			if !trans.Amount.IsZero() {
				strategy.Pool.Mint(trans.TickLower, trans.TickUpper, trans.Amount)
			}

		case "Burn":
			if !trans.Amount.IsZero() {
				strategy.Pool.Burn(trans.TickLower, trans.TickUpper, trans.Amount)
			}

		case "Swap":

			if trans.Amount0.Sign() > 0 {
				if trans.UseX96 {
					strategy.Pool.ExactInputSwap(trans.Amount0, token0, trans.SqrtPriceX96)
				} else {
					strategy.Pool.ExactInputSwap(trans.Amount0, token0, cons.Zero)
				}
			} else if trans.Amount1.Sign() > 0 {
				if trans.UseX96 {
					strategy.Pool.ExactInputSwap(trans.Amount1, token1, trans.SqrtPriceX96)
				} else {
					strategy.Pool.ExactInputSwap(trans.Amount1, token1, cons.Zero)
				}
			}
		case "Flash":
			strategy.Pool.Flash(trans.Amount0, trans.Amount1)
		}

	}
	strategy.BurnAll()
	fmt.Printf("AmountAfter: %d %d \n", strategy.Amount0, strategy.Amount1)
	fmt.Println(strategy.Pool.SqrtRatioX96)

}

func getTransactions() []ent.Transaction {
	filename := "trans.json"
	filepath := path.Join("data", filename)
	file, err := os.Open(filepath)
	check(err)
	value, err := ioutil.ReadAll(file)
	check(err)
	var transactionsInput []ent.TransactionInput
	err = json.Unmarshal([]byte(value), &transactionsInput)
	check(err)
	var transactions []ent.Transaction
	for _, transIn := range transactionsInput {
		useX96, _ := strconv.ParseBool(transIn.UseX96)
		trans := ent.Transaction{
			transIn.Type,
			stringToUint256(transIn.Amount),
			stringToUint256(transIn.Amount0),
			stringToUint256(transIn.Amount1),
			transIn.ID,
			stringToUint256(transIn.SqrtPriceX96),
			transIn.Tick,
			transIn.TickLower,
			transIn.TickUpper,
			transIn.Timestamp,
			useX96,
		}
		transactions = append(transactions, trans)
	}
	return transactions
}

func stringToUint256(amount string) *ui.Int {
	bigint := new(big.Int)
	bigint.SetString(amount, 10)
	uint256, _ := ui.FromBig(bigint)
	return uint256
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
