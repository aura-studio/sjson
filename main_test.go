package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestSjsonCmd(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		set           string
		expectCompact string
		expectIndent  string
	}{
		{
			"int修改",
			`{"a":1}`,
			"a=2",
			`{"a":2}`,
			`{
  "a": 2
}`,
		},
		{
			"嵌套路径",
			`{"a":1}`,
			"b.c=3",
			`{"a":1,"b":{"c":3}}`,
			`{
  "a": 1,
  "b": {
    "c": 3
  }
}`,
		},
		{
			"bool类型",
			`{"a":false}`,
			"a=true",
			`{"a":true}`,
			`{
  "a": true
}`,
		},
		{
			"null类型",
			`{"a":1}`,
			"a=null",
			`{"a":null}`,
			`{
  "a": null
}`,
		},
		{
			"字符串类型",
			`{"a":1}`,
			"b=hello",
			`{"a":1,"b":"hello"}`,
			`{
  "a": 1,
  "b": "hello"
}`,
		},
	}

	// 只保留 compact/indent 两种测试
	for _, tc := range tests {
		t.Run(tc.name+"_compact", func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go", "-json", tc.json, "-set", tc.set)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("运行失败: %v, 输出: %s", err, string(out))
			}
			got := strings.TrimSpace(string(out))
			if got != tc.expectCompact {
				t.Errorf("期望: %s\n实际: %s", tc.expectCompact, got)
			}
		})
		t.Run(tc.name+"_indent", func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go", "-json", tc.json, "-set", tc.set, "-indent")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("运行失败: %v, 输出: %s", err, string(out))
			}
			got := strings.TrimSpace(string(out))
			if got != tc.expectIndent {
				t.Errorf("期望: %s\n实际: %s", tc.expectIndent, got)
			}
		})
	}
}
