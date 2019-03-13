package main

import (
	"crypto/rand"
	"fmt"
	"sort"
)

type T int

// randset chooses n random element from arr, and returns an array
// that contains those randomly chosen elements.
func randset(arr []T, n int) {

}

// choice returns a random element from arr.
func choice(arr []int) int {
	n := len(arr)
	i := randIndex(n - 1)
	return arr[i]
}

// n is the MAXIMUM index.  NOT the length.  n = len(arr) - 1
func randIndex(n int) int {

	nBits := log2roundup(n)
	nBytes := nBits/8 + 1
	xtraBits := nBytes*8 - nBits

	// fmt.Println("MaxIndex:", n)
	// fmt.Println("Num_Bits:", nBits)
	// fmt.Println("Num_Bytes:", nBytes)
	// fmt.Println("Extra_Bits:", xtraBits)

	// result is the randomly generated index that we want to
	// use. It has the possibility of being too high, although the
	// trimming of bits should minimize that possibilitiy.

	var result uint

	// b is an array of bytes that will contain the random
	// numbers. We rely on the fact that these bytes will be
	// totally random, as given by the crypto/rand package.

	randBytes := make([]byte, nBytes)

	// We will have to redo the random number generation if the
	// number is too high, so jump back to this point

	numberOfRetries := -1

generate:

	numberOfRetries++
	result = 0

	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	// trim extra bits from the first bytes. This will lower the
	// chance that we have to recalculate our random number.

	result += uint(randBytes[0]) >> uint(xtraBits)

	// If there are any other bytes we need to process, we shift
	// over 8 bits in our result, and add the next byte. If you
	// imagine it as a string of bits: the next byte is appened to
	// the result string.

	for i := 1; i < len(randBytes); i++ {
		result <<= 8
		result += uint(randBytes[i])
	}

	// Check to see if the index is valid. If it's not, we will
	// redo the calculations. We don't want to mess with any more
	// bit values because that would destroy the randomness.

	if result > uint(n) {
		goto generate
	}

	// fmt.Println("====================")
	// fmt.Printf("Retries: %d\n", numberOfRetries)
	// fmt.Printf("Result_Binary: %08b\n", result)
	// fmt.Printf("Result_Uint:  %d\n", result)

	return int(result)
}

// don't call this for n <= 2, it's not designed to be actual log2
func log2roundup(n int) int {
	if n <= 2 {
		return -1
	}
	result := 1
	total := 2
	for total < n {
		result++
		total *= 2
	}
	return result
}

func stringNLines(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "▒"
	}
	return s
}

func test() {
	const setSize = 2345
	const nTests = 1200100

	const nLines = 30
	const expectedValue = nTests / setSize
	const adjustDivisor = expectedValue / nLines

	fmt.Println("Histogram_of_Results:")
	fmt.Printf(
		"EXPV: %-40v %v \n",
		stringNLines(nLines),
		expectedValue,
	)

	vals := map[int]int{}

	for i := 0; i < nTests; i++ {

		index := randIndex(setSize - 1)
		_, ok := vals[index]
		if !ok {
			vals[index] = 0
		}
		vals[index]++
	}

	keys := []int{}
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v := vals[k]
		toPrint := v / adjustDivisor
		fmt.Printf(
			"%4d: %-40v %d (%+d) \n",
			k,
			stringNLines(toPrint),
			v,
			v-expectedValue,
		)
	}
}

func main() {
	test()
	// n := 300
	// ex := make([]int, n)
	// for i := 0; i < n; i++ {
	// 	ex[i] = i
	// }
	// x := choice(ex)
	// fmt.Println(x)
}
