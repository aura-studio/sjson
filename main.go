package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// setJSONPath 设置 json 中的路径为指定值，支持嵌套路径 a.b.c
func setJSONPath(data map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	m := data
	for i, part := range parts {
		if i == len(parts)-1 {
			m[part] = value
			return
		}
		if _, ok := m[part]; !ok {
			m[part] = make(map[string]interface{})
		}
		if next, ok := m[part].(map[string]interface{}); ok {
			m = next
		} else {
			// 路径冲突，覆盖为 map
			m[part] = make(map[string]interface{})
			m = m[part].(map[string]interface{})
		}
	}
}

func main() {
	jsonStr := flag.String("json", "", "输入的json字符串")
	setArg := flag.String("set", "", "设置路径和值，格式为 a.b.c=value")
	indent := flag.Bool("indent", false, "是否缩进输出json")
	flag.Parse()

	if *jsonStr == "" || *setArg == "" {
		fmt.Println("用法: -json '{\"a\":1}' -set a.b.c=123 [-indent]")
		os.Exit(1)
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(*jsonStr), &data)
	if err != nil {
		fmt.Println("json解析失败:", err)
		os.Exit(1)
	}

	// 解析 -set 参数
	parts := strings.SplitN(*setArg, "=", 2)
	if len(parts) != 2 {
		fmt.Println("-set 参数格式错误，应为 a.b.c=value")
		os.Exit(1)
	}
	path, valStr := parts[0], parts[1]

	// 尝试将 valStr 解析为 int、float、bool、null 或字符串
	var val interface{}
	if valStr == "null" {
		val = nil
	} else if valStr == "true" {
		val = true
	} else if valStr == "false" {
		val = false
	} else if i, err := json.Number(valStr).Int64(); err == nil {
		val = i
	} else if f, err := json.Number(valStr).Float64(); err == nil {
		val = f
	} else {
		val = valStr
	}

	setJSONPath(data, path, val)

	var out []byte
	if *indent {
		out, err = json.MarshalIndent(data, "", "  ")
	} else {
		out, err = json.Marshal(data)
	}
	if err != nil {
		fmt.Println("json序列化失败:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
