package Figo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/figoxu/color"
	"github.com/jinzhu/copier"
	"github.com/quexer/utee"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

func Catch(hooks ...func()) {
	if err := recover(); err != nil {
		log.Println(string(debug.Stack()))
		log.Println(err, " (recover)")
		for _, hook := range hooks {
			hook()
		}
	}
}

func Clone(to, from interface{}) {
	if reflect.TypeOf(to).Kind().String() != "ptr" {
		utee.Chk(errors.New("Parameter 'to' Should Be An PTR"))
	}
	copier.Copy(to, from)
}

func Exist(expect interface{}, objs ...interface{}) bool {
	for _, v := range objs {
		if expect == v {
			return true
		}
	}
	return false
}

func RetryExe(business func() error, times int, tips string) (int, bool) {
	err := business()
	retry := 0
	for err != nil && retry < times {
		retry++
		err = business()
	}
	success := (err == nil)
	if retry > 0 && tips != "" {
		log.Println(tips, " Execute With ", retry, " times .  @SuccessFlag:", success)
	}
	return retry, success
}

func ParseUrl(s string) (string, int, error) {
	a := strings.Split(s, ":")
	if len(a) != 2 {
		return "", 0, fmt.Errorf("bad url %s", s)
	}
	port, err := strconv.Atoi(a[1])
	return a[0], port, err
}

func NotFound(err error) bool {
	return err == mgo.ErrNotFound || err == orm.ErrNoRows
}

const (
	THEME_Black   = "black"
	THEME_Red     = "red"
	THEME_Green   = "green"
	THEME_Yellow  = "yellow"
	THEME_Blue    = "blue"
	THEME_Magenta = "megenta"
	THEME_Cyan    = "cyan"
	THEME_White   = "white"
)

func ReadInput(tips, tipTheme, inputTheme string) string {
	Print(tipTheme, tips)
	Print(inputTheme, " ")
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	return string(data)
}

func Print(theme string, v ...interface{}) {
	s := fmt.Sprint(v...)
	switch theme {
	case THEME_Black:
		color.Black(s)
	case THEME_Red:
		color.Red(s)
	case THEME_Green:
		color.Green(s)
	case THEME_Yellow:
		color.Yellow(s)
	case THEME_Blue:
		color.Blue(s)
	case THEME_Magenta:
		color.Magenta(s)
	case THEME_Cyan:
		color.Cyan(s)
	case THEME_White:
		color.White(s)
	default:
		log.Println(s)
	}
}

func Println(theme string, v ...interface{}) {
	s := fmt.Sprint(v...)
	Print(theme, s, "\n")
}

func PrintJson(prefix string, v interface{}) {
	if b, err := json.Marshal(v); err == nil {
		Print(THEME_Magenta, prefix)
		Print(THEME_Blue, string(b), "\n")
	}
}

func JsonString(obj interface{}) string {
	b, _ := json.Marshal(obj)
	s := fmt.Sprintf("%+v", string(b))
	r := strings.Replace(s, `\u003c`, "<", -1)
	r = strings.Replace(r, `\u003e`, ">", -1)
	return r
}
