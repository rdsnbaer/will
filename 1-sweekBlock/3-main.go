package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

//

//注释解释的是上一行的代码
func main() {
	// url := "http://8.209.64.180:8088/coinInfo?symbol=%E5%B8%81%E7%A7%8D"
	str := "..."

	fmt.Println(url.PathEscape(str))

	b := str2gbk(str)
	fmt.Println(b)
	//get(url)

	//a := ConvertStr2GBK(str)

	//str := "！@#中国123"
	//设定一个含有中文的字符串
	var a = regexp.MustCompile("^[\u4e00-\u9fa5]$")
	//接受正则表达式的范围
	for i, v := range str {
		//golang中string的底层是byte类型，所以单纯的for输出中文会出现乱码，这里选择for-range来输出
		if a.MatchString(string(v)) {
			//判断是否为中文，如果是返回一个true，不是返回false。这俩面MatchString的参数要求是string
			//但是 for-range 返回的 value 是 rune 类型，所以需要做一个 string() 转换
			fmt.Printf("str 字符串第1 %v 个字符是中文。是2“%v”字\n", i+1, string(v))
			m := ConvertStr2GBK(string(v))
			fmt.Println("2-m:", m)

			s := str2gbk(string(v))
			fmt.Println(s)
		}
	}

	// ret, err := simplifiedchinese.GBK.NewEncoder().String(str)
	// return ret   //如果转换失败返回空字符串

	// //如果是[]byte格式的字符串，可以使用Bytes方法
	// b, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))
	// return string(b)

	// u, err := url.Parse("http://8.209.64.180:8088/coinInfo?symbol=%E5%B8%81%E7%A7%8D")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// u.Scheme = "https"
	// u.Host = "google.com"
	// q := u.Query()
	// q.Set("q", "golang")
	// u.RawQuery = q.Encode()
	// fmt.Println(u)

}

// 字符串GBK2312编码方式解码方法
func str2gbk(text string) string {
	TextBuff, err := simplifiedchinese.GBK.NewEncoder().String(text)
	if err != nil {
		fmt.Println(err)
	}
	return TextBuff
}

func ConvertStr2GBK(str string) string {
	//将utf-8编码的字符串转换为GBK编码
	ret, _ := simplifiedchinese.GBK.NewEncoder().String(str)
	return ret //如果转换失败返回空字符串

	//如果是[]byte格式的字符串，可以使用Bytes方法
	b, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))
	return string(b)
}

func get(url string) {
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("111-err:", err)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return
	}

	fmt.Println(string(body))
}

// func main() {
// 	sText := "中文"
// 	textQuoted := strconv.QuoteToASCII(sText)
// 	textUnquoted := textQuoted[1 : len(textQuoted)-1]
// 	fmt.Println(textUnquoted)

// 	sUnicodev := strings.Split(textUnquoted, "\\u")
// 	var context string
// 	for _, v := range sUnicodev {
// 		if len(v) < 1 {
// 			continue
// 		}
// 		temp, err := strconv.ParseInt(v, 16, 32)
// 		if err != nil {
// 			panic(err)
// 		}
// 		context += fmt.Sprintf("%c", temp)
// 	}
// 	fmt.Println(context)

// }
