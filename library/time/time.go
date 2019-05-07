package time

// 代码参考k8s Time封装

import (
	"encoding/json"
	"time"

	"github.com/google/gofuzz"
)

const (
	StampForWeb = "2006-01-02 15:04:05" // 替换原包中的RFC3339
)

// Time 针对time.Time封装，支持解析到YAML和JSON.
// 封装提供了一些工厂方法.
//
// +protobuf.options.marshal=false
// +protobuf.as=Timestamp
// +protobuf.options.(gogoproto.goproto_stringer)=false
type Time struct {
	time.Time `protobuf:"-"`
}

// DeepCopyInto 深拷贝时间值.
func (t *Time) DeepCopyInto(out *Time) {
	*out = *t
}

// String 返回时间描述.
func (t Time) String() string {
	return t.Time.Format(StampForWeb)
	// return t.Time.String()
}

// NewTime 创建一个Time
func NewTime(time time.Time) Time {
	return Time{time}
}

// Date 按描述返回Time对象
// 封装time.Date实现.
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

// Now 返回当前本地时间.
func Now() Time {
	return Time{time.Now()}
}

// IsZero 返回时间值是否为nil或0.
func (t *Time) IsZero() bool {
	if t == nil {
		return true
	}
	return t.Time.IsZero()
}

// Before 返回时间t是否在目标时间u之前.
func (t *Time) Before(u *Time) bool {
	return t.Time.Before(u.Time)
}

// Equal 比较时间t和u是否相等.
func (t *Time) Equal(u *Time) bool {
	if t == nil && u == nil {
		return true
	}
	if t != nil && u != nil {
		return t.Time.Equal(u.Time)
	}
	return false
}

// Unix 返回时间戳
// 封装time.Unix实现.
func Unix(sec int64, nsec int64) Time {
	return Time{time.Unix(sec, nsec)}
}

// Rfc3339Copy 拷贝.
func (t Time) Rfc3339Copy() Time {
	copied, _ := time.Parse(time.RFC3339, t.Format(time.RFC3339))
	return Time{copied}
}

// ForWebCopy 拷贝
func (t Time) ForWebCopy() Time {
	copied, _ := time.Parse(StampForWeb, t.Format(StampForWeb))
	return Time{copied}
}

// UnmarshalJSON 实现json.Unmarshaller接口.
func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 4 && string(b) == "null" {
		t.Time = time.Time{}
		return nil
	}

	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	pt, err := time.Parse(StampForWeb, str)
	if err != nil {
		return err
	}

	t.Time = pt.Local()
	return nil
}

// UnmarshalQueryParameter 从URL查询参数里解析
func (t *Time) UnmarshalQueryParameter(str string) error {
	if len(str) == 0 {
		t.Time = time.Time{}
		return nil
	}
	if len(str) == 4 && str == "null" {
		t.Time = time.Time{}
		return nil
	}

	pt, err := time.Parse(StampForWeb, str)
	if err != nil {
		return err
	}

	t.Time = pt.Local()
	return nil
}

// MarshalJSON 实现json.Marshaler接口.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		// Encode unset/nil objects as JSON's "null".
		return []byte("null"), nil
	}

	return json.Marshal(t.UTC().Format(StampForWeb))
}

// OpenAPISchemaType 满足OpenAPI规格.
func (_ Time) OpenAPISchemaType() []string { return []string{"string"} }

// OpenAPISchemaFormat  满足OpenAPI规格.
func (_ Time) OpenAPISchemaFormat() string { return "date-time" }

// MarshalQueryParameter  满足OpenAPI规格
func (t Time) MarshalQueryParameter() (string, error) {
	if t.IsZero() {
		// 用空字符串编码unset/nil
		return "", nil
	}

	return t.UTC().Format(StampForWeb), nil
}

// Fuzz 实现 fuzz.Interface.
func (t *Time) Fuzz(c fuzz.Continue) {
	if t == nil {
		return
	}
	t.Time = time.Unix(c.Rand.Int63n(1000*365*24*60*60), 0)
}

// 利用编译器检查Time是否实现fuzz.Interface
var _ fuzz.Interface = &Time{}
