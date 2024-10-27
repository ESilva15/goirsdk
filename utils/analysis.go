package utils

import "fmt"

func HexDump(buf []byte) {
	fmt.Printf("\n============ HEX DUMP ============\n")
	for k := 0; k < len(buf); k++ {
		if k%4 == 0 && k > 0 {
			fmt.Printf(" ")
		}

		if k%16 == 0 && k > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("%02X", buf[k])
	}
	fmt.Printf("\n============ HEX DUMP ============\n")
}
