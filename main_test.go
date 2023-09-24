package main_test

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// func nlp(s string) (a, b string) {
// 	sr := strings.SplitAfter(s, "\n")
// 	var c string
// 	for i, v := range sr {
// 		c += v
// 		if len(c) > 1293 {
// 			for _, d := range sr[i:] {
// 				b += d
// 			}
// 			break
// 		}
// 		a = c
// 	}
// 	return
// }

var reg, _ = regexp.Compile("\n+")

// 再次拆分
func nlps(s []string, m string) (o []string) {
	for _, v := range s {
		o = append(o, strings.SplitAfter(v, m)...)
	}
	return
}

// 分隔字符串
func nlp2(s string) (a, b string) {
	sr := strings.SplitAfter(s, "\n")
	out := nlps(
		nlps(
			nlps(
				sr,
				` `),
			`。`),
		`，`)
	// var out []string
	// for _, v := range sr {
	// 	out = append(out, strings.SplitAfter(v, ` `)...)
	// }
	var c string
	for i, v := range out {
		c += v
		if len(c) > 570 {
			for _, d := range out[i:] {
				b += d
			}
			break
		}
		a = c
	}
	return
}

var text string

func init() {
	file, _ := os.Open(`test2.txt`)
	defer file.Close()
	txt, _ := io.ReadAll(file)
	text = reg.ReplaceAllString(string(txt), "\n")
}

// func TestNLP(t *testing.T) {
// 	// file, _ := os.Open(`test.txt`)
// 	// defer file.Close()
// 	// txt, _ := io.ReadAll(file)
// 	// text := string(txt)

// 	a, b := nlp2(text)
// 	t.Logf("a:\n[%v]\nb:\n[%v]\n", a, b)
// }

func TestTTS(t *testing.T) {
	nla, nlb := nlp2(text)
	fmt.Printf("[0]: %v\n", nla)
	for i := 0; nlb != ``; i++ {
		nla, nlb = nlp2(nlb)
		fmt.Printf("[%v]: %v\n", i+1, nla)
	}
}
