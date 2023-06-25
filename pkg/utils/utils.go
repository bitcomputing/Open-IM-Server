package utils

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}
func cleanUpFuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

// Get the intersection of two slices
func Intersect(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func Difference(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

// Get the intersection of two slices
func IntersectString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	inter := IntersectString(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

func RemoveFromSlice(slice1, slice2 []string) []string {
	for _, v1 := range slice1 {
		for i2, v2 := range slice2 {
			if v2 == v1 {
				if i2 != len(slice2)-1 {
					slice2 = append(slice2[:i2], slice2[i2+1:]...)
				} else {
					slice2 = append(slice2[:i2])
				}
			}
		}
	}
	return slice2
}

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func RemoveRepeatedStringInList(slc []string) []string {
	var result []string
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func Pb2String(pb proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}
	buffer, err := marshaler.Marshal(pb)
	return string(buffer), err
}

func String2Pb(s string, pb proto.Message) error {
	return protojson.Unmarshal([]byte(s), pb)

}

func Map2Pb(m map[string]string) (pb proto.Message, err error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(b, pb)
	if err != nil {
		return nil, err
	}
	return pb, nil
}
func Pb2Map(pb proto.Message) (map[string]interface{}, error) {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: false,
	}
	ret := make(map[string]interface{})
	buffer, err := marshaler.Marshal(pb)
	if err != nil {
		return ret, err
	}

	if err := json.Unmarshal(buffer, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
