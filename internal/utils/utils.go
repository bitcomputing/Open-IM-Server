package utils

import (
	"encoding/json"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func JsonDataList(resp interface{}) []map[string]interface{} {
	var list []proto.Message
	if reflect.TypeOf(resp).Kind() == reflect.Slice {
		s := reflect.ValueOf(resp)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface().(proto.Message))
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range list {
		m := ProtoToMap(v, false)
		result = append(result, m)
	}
	return result
}

func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
	}

	s, _ := marshaler.Marshal(pb)
	out := make(map[string]interface{})
	json.Unmarshal(s, &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}
