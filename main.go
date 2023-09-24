package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// 生成22位随机字符串 `rgHrZPfgAsTpkgXJj4qgBA`
func getEi() string {
	list := []string{`a`, `A`, `b`, `B`, `c`, `C`, `d`, `D`, `e`, `E`, `f`, `F`, `g`, `G`, `h`, `H`, `i`, `I`, `j`, `J`, `k`, `K`, `l`, `L`, `m`, `M`, `n`, `N`, `o`, `O`, `p`, `P`, `q`, `Q`, `r`, `R`, `s`, `S`, `t`, `T`, `u`, `U`, `v`, `V`, `w`, `W`, `x`, `X`, `y`, `Y`, `z`, `Z`}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
	var out string
	for i := 0; i < 22; i++ {
		out += list[r.Intn(52)]
	}
	return out
}

// 获取谷歌TTS朗读结果 ['语言','文本']['base64','错误'] `%25E7%2594%25A8%25E6%2588%25B6%25E7%2595%258C%25E9%259D%25A2`
func getTTS(lang, text string) (string, error) {
	var proxyurl *url.URL
	client := &http.Client{Timeout: 30 * time.Second}
	if proxy != `` {
		var err error
		proxyurl, err = url.Parse(proxy)
		if err != nil {
			return ``, fmt.Errorf(`代理URL解析失败: %v`, err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyurl)}
	}

	type Resp struct {
		Url []string `json:"translate_tts"`
	}
	api := `https://www.google.de/async/translate_tts?ei=` + getEi() + `&opi=89978449&yv=3&ttsp=tl:` + lang + `,txt:` + text + `,spd:1&cs=0&async=_fmt:jspb`

	req, _ := http.NewRequest(`GET`, api, nil)
	req.Header.Set(`User-Agent`, `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.188`)
	req.Header.Add(`Accept`, `application/json, text/plain, */*`)
	req.Header.Add(`Referer`, `https://www.google.de/`)
	resp, err := client.Do(req)
	if err != nil {
		return ``, fmt.Errorf(`HTTP请求失败: %v`, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ``, fmt.Errorf(`读取Body内容失败: %v`, err)
	}
	var ret Resp
	err = json.Unmarshal(bytes.TrimPrefix(body, []byte(")]}'\n")), &ret)
	if err != nil {
		fmt.Println(body)
		return ``, fmt.Errorf(`返回Json解析失败: %v`, err)
	}

	if len(ret.Url) >= 1 {
		return ret.Url[0], nil
	}

	return ``, fmt.Errorf(`未返回有效数据: %v`, ret)
}

func showHelp() {
	fmt.Print(`
暂无帮助
`)
	os.Exit(0)
}

// 再次拆分
func nlps(s []string, m string) (o []string) {
	for _, v := range s {
		o = append(o, strings.SplitAfter(v, m)...)
	}
	return
}

// 分隔字符串
func nlp(s string) (a, b string) {
	sr := strings.SplitAfter(s, "\n")
	out := nlps(nlps(nlps(sr, ` `), `。`), `，`)
	var c string
	for i, v := range out {
		c += v
		if len(c) > 550 {
			for _, d := range out[i:] {
				b += d
			}
			break
		}
		a = c
	}
	return
}

// 执行命令
func Cmd(command string) (string, error) {
	cmd := exec.Command(`/usr/bin/env`, `bash`, `-c`, command)
	var out bytes.Buffer
	var ert bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &ert
	err := cmd.Start()
	if err != nil {
		return ert.String(), err
	}
	err = cmd.Wait()
	return out.String(), err
}

// func try(f func()error) error {
// 	return f()
// }

var (
	lang   string   // 源语言
	text   string   // 待翻译文本
	usetxt bool     // 解析文本文档
	name   string   // 输出文件名
	help   bool     // 显示帮助
	proxy  string   // 代理地址
	format bool     // 混合分段音频
	nlist  []string // 文件名列表
	// i18n = []string{`zh-CN`} // 支持的语言列表

	reg, _ = regexp.Compile("\n+") // 匹配连续换行
)

func init() {
	flag.StringVar(&lang, `lang`, ``, `源语言`)
	flag.StringVar(&text, `text`, ``, `要翻译的文本`)
	flag.BoolVar(&usetxt, `usetxt`, false, `从text参数解析文本文档`)
	flag.StringVar(&name, `name`, time.Now().Format(time.RFC3339), `输出文件名`)
	flag.BoolVar(&help, `help`, false, `显示帮助`)
	flag.StringVar(&proxy, `proxy`, ``, `使用代理`)
	flag.BoolVar(&format, `format`, false, `混合分段音频 (需要FFmpeg命令)`)
	flag.Parse()
}

func main() {
	if help {
		showHelp()
	}
	if lang == `` {
		log.Fatal(`未输入源语言参数 (lang)`)
	}
	if text == `` {
		log.Fatal(`未输入待翻译文本参数 (text)`)
	}
	if usetxt {
		stat, err := os.Stat(text)
		if err != nil || stat.IsDir() {
			log.Fatal(`指定文件不存在 (` + text + `)`)
		}
		file, err := os.Open(text)
		if err != nil {
			log.Fatal(`打开文件失败: `, err)
		}
		defer file.Close()
		txt, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(`读取文件失败: `, err)
		}
		text = string(txt)
	}
	text = reg.ReplaceAllString(text, "\n")
	fmt.Printf("语言: %v, 文本内容: \n[%v]\n", lang, text)
	log.Println(`请稍候...`)
	time.Sleep(time.Millisecond * 300 * 2)

	nla, nlb := nlp(text)
	getTTS2(nla, `0`)
	for i := 0; nlb != ``; i++ {
		nla, nlb = nlp(nlb)
		getTTS2(nla, fmt.Sprint(i+1))
	}

	if format {
		msg, err := Cmd(fmt.Sprintf(`ffmpeg -i "concat:%v" -loglevel quiet -c copy "%v.mp3"`, strings.Join(nlist, `|`), name))
		if err != nil {
			log.Printf("命令执行失败，%v，输出: %v\n", err, msg)
		} else {
			log.Println(`混合分段音频成功`)
		}
	}

	fmt.Println(`Success!`)
}

func getTTS2(text, exname string) {
	out, err := getTTS(lang, url.QueryEscape(strings.ReplaceAll(text, "\n", ` `)))
	if err != nil {
		log.Fatal(`TTS请求错误: `, err)
	}
	dec, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		log.Fatal(`Base64解码错误: `, err)
	}

	fname := name + `_` + exname + `.mp3`
	nlist = append(nlist, fname)
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal(`创建文件错误: `, err)
	}
	file.Write(dec)
	file.Close()
	time.Sleep(time.Millisecond * 300)
}
