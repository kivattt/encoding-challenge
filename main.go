package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivattt/getopt"
)

const lookup = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// Code taken from the encoding/base64 package in the Golang standard library, then edited
func encode(str, realStr string) string {
	if len(str) != len(realStr) {
		panic("Epic fail")
	}

	if len(str) == 0 {
		return ""
	}

	str += str
	str += str
	str += str

	size := 2 * ((len(str) + 2) / 3 * 4)
	var ret []byte
	for i := 0; i < size; i++ {
		ret = append(ret, ' ')
	}

	di, si := 0, 0
	for realIndex := 0; realIndex < len(realStr); realIndex++ {
		val := uint(str[si])<<16 | uint(str[si+1])<<8 | uint(str[si+2])

		ret[di+0] = lookup[val>>18&0x3F]
		ret[di+1] = lookup[val>>12&0x3F]
		ret[di+2] = lookup[realStr[realIndex]&0x3F]
		ret[di+3] = lookup[(realStr[realIndex]>>6)&3]

		si += 3
		di += 4
	}

	remain := len(realStr) % 4

	if remain == 0 {
		return strings.TrimRight(string(ret), " ")
	}

	val := uint(str[si]) << 16
	if remain == 2 {
		val |= uint(str[si+1]) << 8
	}

	ret[di] = lookup[val>>18&0x3F]
	ret[di+1] = lookup[val>>12&0x3F]

	switch remain {
	case 2:
		ret[di+2] = lookup[val>>6&0x3F]
		ret[di+3] = '='
	case 1:
		ret[di+2] = '='
		ret[di+3] = '='
	}

	return strings.TrimRight(string(ret), " ")
}

func main() {
	help := flag.Bool("help", false, "display this help and exit")
	decode := flag.Bool("decode", false, "decode instead of encode")
	data := flag.String("data", "", "data to encode/decode")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("challenge", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"d", "decode",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *help {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS]")
		fmt.Println("Encoding challenge")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	if *decode {
		for i := 2; i < len(*data); i += 4 {
			c1 := (*data)[i]
			c2 := (*data)[i+1]

			if c1 == '=' || c2 == '=' {
				break
			}

			a := strings.IndexByte(lookup, c1)
			b := strings.IndexByte(lookup, c2)

			c := (a & 0x3F) | ((b & 3) << 6)
			os.Stdout.Write([]byte{byte(c)})
		}
		os.Exit(0)
	}

	var seed int64
	for i := 0; i < len(*data); i++ {
		if i > 4 {
			break
		}

		seed |= int64((*data)[i]) << (i * 8)
	}
	r := rand.New(rand.NewSource(seed))

	var randomDataBuilder strings.Builder
	for i := 0; i < len(*data); i++ {
		randomDataBuilder.WriteByte(byte(r.Intn(255)))
	}
	fmt.Println(encode(randomDataBuilder.String(), *data))
}
