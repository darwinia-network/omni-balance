package utils

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/rand"
	"io"
	"math/big"
	"net/http"
	"omni-balance/utils/constant"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var (
	index atomic.Int64
)

func GetNextIndex(v int) int {
	if v <= 0 {
		return 0
	}
	return int(index.Add(1) % int64(v))
}

func GetNextIndexStrings(v []string) string {
	if len(v) <= 0 {
		return ""
	}
	return v[GetNextIndex(len(v))]
}

func ExtractTagFromStruct(s interface{}, tags ...string) map[string]map[string]string {
	t := reflect.TypeOf(s)
	t = t.Elem()
	var result = make(map[string]map[string]string)
	n := t.NumField()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		for _, tag := range tags {
			if _, ok := result[field.Name]; !ok {
				result[field.Name] = make(map[string]string)
			}
			result[field.Name][tag] = field.Tag.Get(tag)
		}
	}
	return result
}

func GetCurrencyPair(sourceToken, sep, targetToken string) string {
	return strings.ToUpper(fmt.Sprintf("%s%s%s", sourceToken, sep, targetToken))
}

type Element interface {
	string | int8 | int16 | int32 | int64 | int | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

func InArray[t Element](value t, arr []t) bool {
	for _, v := range arr {
		if value != v {
			continue
		}
		return true
	}
	return false
}

func InArrayFold(value string, arr []string) bool {
	for _, v := range arr {
		if !strings.EqualFold(v, value) {
			continue
		}
		return true
	}
	return false
}

func AssertEqualFold(t *testing.T, a, b string) {
	assert.Equal(t, strings.ToUpper(a), strings.ToUpper(b))
}

func ZFillByte(b []byte, n int) []byte {
	if len(b) >= n {
		return b
	}
	return append(make([]byte, n-len(b)), b...)
}

// SetBit sets the nth bit of the byte to v.
func SetBit(b byte, n byte, v bool) byte {
	if v {
		return b | (1 << n)
	}
	return b &^ (1 << n)
}

type PermitSingle struct {
	Owner    string
	Spender  string
	Value    *big.Int
	Deadline *big.Int
	Nonce    *big.Int
}

func RandomNumberIn[t constraints.Integer](min, max t) t {
	rand.New(rand.NewSource(uint64(time.Now().UnixMilli())))
	return t(rand.Intn(int(max-min+1)) + int(min))
}

func Choose[T any](arr []T) T {
	if len(arr) == 0 {
		panic("arr is empty")
	}
	if len(arr) == 1 {
		return arr[0]
	}
	return arr[RandomNumberIn(0, len(arr)-1)]
}

func Request(ctx context.Context, method string, url string, body io.Reader, dest interface{}, headers ...string) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return errors.Wrap(err, "new request")
	}
	for i := 0; i < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/json")
	}
	if req.Header.Get("accept") == "" {
		req.Header.Set("accept", "application/json")
	}
	req = req.WithContext(ctx)
	var count int
	for count < 3 {
		count++
		resp, err := new(http.Client).Do(req)
		if err != nil {
			logrus.Debugf("request failed: %s", err)
			time.Sleep(time.Second)
			continue
		}
		if dest == nil {
			_ = resp.Body.Close()
			return nil
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Debugf("read response failed: %s", err)
			_ = resp.Body.Close()
			time.Sleep(time.Second)
			continue
		}
		_ = resp.Body.Close()
		if err := json.Unmarshal(data, dest); err != nil {
			return errors.Errorf("response code is %s, body: %s", resp.Status, data)
		}
		return nil
	}
	return errors.New("request failed")
}

func String(s string) *string {
	return &s
}

func Number[t constraints.Integer | constraints.Float](n t) *t {
	return &n
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func ToMap(v interface{}) map[string]interface{} {
	data, _ := json.Marshal(v)
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	return result
}

func SetLogToCtx(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, constant.LogKeyInCtx, log)
}

func GetLogFromCtx(ctx context.Context) *logrus.Entry {
	if v, ok := ctx.Value(constant.LogKeyInCtx).(*logrus.Entry); ok {
		return v
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

func HexToString(v string) string {
	i, _ := new(big.Int).SetString(strings.TrimPrefix(v, "0x"), 16)
	return i.String()
}
