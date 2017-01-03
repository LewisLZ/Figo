package Figo

import (
	"fmt"
	"github.com/quexer/utee"
	"regexp"
	"strconv"
	"strings"
)

func SplitUTF8BOM(str string) string {
	b := []byte(str)
	if len(b) < 3 {
		return str
	}
	prefix := fmt.Sprintf("%X", b[0:3])
	if prefix == "EFBBBF" {
		return string(b[3:len(b)])
	}
	return str
}

type Parser struct {
	PrepareReg []string
	ProcessReg []string
}

func (p *Parser) Exe(content string) []string {
	prep := func(reg string, contents ...string) []string {
		var result []string
		for _, content := range contents {
			rs := regexp.MustCompile(reg).FindAllString(content, -1)
			result = append(result, rs...)
		}
		return result
	}
	proc := func(reg string, contents ...string) []string {
		var result []string
		for _, content := range contents {
			rs := regexp.MustCompile(reg).ReplaceAllString(content, "")
			result = append(result, rs)
		}
		return result
	}
	result := []string{content}
	for _, reg := range p.PrepareReg {
		result = prep(reg, result...)
	}
	for _, reg := range p.ProcessReg {
		result = proc(reg, result...)
	}
	return TrimAndClear(result...)
}

func TrimAndClear(strs ...string) []string {
	result := []string{}
	for _, v := range strs {
		v = strings.TrimSpace(v)
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func Md5Shard(key string, piece int) string {
	shardVal, err := strconv.ParseUint(utee.PlainMd5(key)[16:32], 16, 0)
	utee.Chk(err)
	shardVal = shardVal % uint64(piece)
	return fmt.Sprint(key, "_", shardVal)
}
