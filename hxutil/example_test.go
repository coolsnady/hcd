package hxutil_test

import (
	"fmt"
	"math"

	"github.com/coolsnady/hxd/hxutil"
)

func ExampleAmount() {

	a := hxutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = hxutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = hxutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 HX
	// 100,000,000 Atoms: 1 HX
	// 100,000 Atoms: 0.001 HX
}

func ExampleNewAmount() {
	amountOne, err := hxutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := hxutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := hxutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := hxutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 HX
	// 0.01234567 HX
	// 0 HX
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := hxutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(hxutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(hxutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(hxutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(hxutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kDCR
	// Atom to Coin: 444333.222111 HX
	// Atom to MilliCoin: 444333222.111 mDCR
	// Atom to MicroCoin: 444333222111 μDCR
	// Atom to Atom: 44433322211100 Atom
}
