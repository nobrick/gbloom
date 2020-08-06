package gbloom_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math"
	"os"

	"github.com/nobrick/gbloom"
)

func ExampleBloomFilter() {
	hashes, err := readTestData()
	if err != nil {
		fmt.Println(err)
		return
	}

	n := len(hashes)
	nAdd := len(hashes) / 2
	fmt.Printf("nAdd: %d\nnTest: %d\n\n", nAdd, n-nAdd)

	testBloomFillter(hashes, 0, nAdd)

	maxBits := nAdd
	for i := 0; i < 17; i++ {
		fmt.Println()

		testBloomFillter(hashes, maxBits, nAdd)
		maxBits <<= 1
	}

	// Output:
	// nAdd: 10598
	// nTest: 10598
	//
	// MaxBits: 100000007 (9436x)
	// State size: 12.498 mB
	// False positive: 1 (0.01%)
	//
	// MaxBits: 10598 (1x)
	// State size: 0.001 mB
	// False positive: 6755 (63.74%)
	//
	// MaxBits: 21196 (2x)
	// State size: 0.003 mB
	// False positive: 4181 (39.45%)
	//
	// MaxBits: 42392 (4x)
	// State size: 0.005 mB
	// False positive: 2265 (21.37%)
	//
	// MaxBits: 84784 (8x)
	// State size: 0.011 mB
	// False positive: 1225 (11.56%)
	//
	// MaxBits: 169568 (16x)
	// State size: 0.021 mB
	// False positive: 601 (5.67%)
	//
	// MaxBits: 339136 (32x)
	// State size: 0.042 mB
	// False positive: 301 (2.84%)
	//
	// MaxBits: 678272 (64x)
	// State size: 0.085 mB
	// False positive: 133 (1.25%)
	//
	// MaxBits: 1356544 (128x)
	// State size: 0.170 mB
	// False positive: 74 (0.70%)
	//
	// MaxBits: 2713088 (256x)
	// State size: 0.339 mB
	// False positive: 34 (0.32%)
	//
	// MaxBits: 5426176 (512x)
	// State size: 0.678 mB
	// False positive: 18 (0.17%)
	//
	// MaxBits: 10852352 (1024x)
	// State size: 1.357 mB
	// False positive: 12 (0.11%)
	//
	// MaxBits: 21704704 (2048x)
	// State size: 2.713 mB
	// False positive: 7 (0.07%)
	//
	// MaxBits: 43409408 (4096x)
	// State size: 5.425 mB
	// False positive: 3 (0.03%)
	//
	// MaxBits: 86818816 (8192x)
	// State size: 10.851 mB
	// False positive: 1 (0.01%)
	//
	// MaxBits: 173637632 (16384x)
	// State size: 21.703 mB
	// False positive: 1 (0.01%)
	//
	// MaxBits: 347275264 (32768x)
	// State size: 43.408 mB
	// False positive: 0 (0.00%)
	//
	// MaxBits: 694550528 (65536x)
	// State size: 86.817 mB
	// False positive: 0 (0.00%)
}

func testBloomFillter(hashes []string, maxBits, nAdd int) {
	f := gbloom.New(maxBits, nil)
	n := len(hashes)

	for i := 0; i < nAdd; i++ {
		hash := hexToBytes(hashes[i])
		f.AddHash(hash)

		ok := f.TestHash(hash)
		if !ok {
			fmt.Printf("Expect %s to be present in the filter\n", hashes[i])
		}
	}

	nFalsePositive := 0
	for i := nAdd; i < n; i++ {
		found := f.TestHash(hexToBytes(hashes[i]))
		if found {
			nFalsePositive++
		}
	}

	state := f.Bytes()
	m := f.MaxBits()
	fpRate := float64(nFalsePositive) / float64(n-nAdd)

	fmt.Printf("MaxBits: %d (%.0fx)\n", m, math.Round(float64(m)/float64(nAdd)))
	fmt.Printf("State size: %.03f mB\n", float64(len(state))/1_000_000.0)
	fmt.Printf("False positive: %d (%.2f%%)\n", nFalsePositive, fpRate*100)

	filterNew := gbloom.NewFromBytes(maxBits, nil, state)
	for i := 0; i < nAdd; i++ {
		ok := filterNew.TestHash(hexToBytes(hashes[i]))
		if !ok {
			fmt.Printf("Expect %s to be present in the restored filter\n", hashes[i])
		}
	}
}

func readTestData() (hashes []string, _ error) {
	file, err := os.Open("testdata/hashes")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hashes = append(hashes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hashes, nil
}

func hexToBytes(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return bytes
}
