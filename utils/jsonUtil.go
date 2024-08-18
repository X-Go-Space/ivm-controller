package utils

import (
	"encoding/json"
	"strings"
)

// UnmarshalNestedJSON 反序列化任意级别嵌套结构的 JSON 数据
func UnmarshalNestedJSON(jsonStr string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	unmarshalNestedJSONHelper(data)

	return data, nil
}

func unmarshalNestedJSONHelper(data map[string]interface{}) {
	for key, value := range data {
		if nestedJSON, ok := value.(string); ok {
			var nestedData map[string]interface{}
			if err := json.Unmarshal([]byte(nestedJSON), &nestedData); err == nil {
				data[key] = nestedData
				unmarshalNestedJSONHelper(nestedData) // 递归处理嵌套结构
			}
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			unmarshalNestedJSONHelper(nestedMap) // 递归处理嵌套结构
		}
	}
}

// ReadNestedData 读取任意级别嵌套结构中对应的数据
func ReadNestedData(data map[string]interface{}, keyPath string) interface{} {
	keys := strings.Split(keyPath, ".")
	current := data
	for _, key := range keys {
		if val, ok := current[key]; ok {
			switch v := val.(type) {
			case map[string]interface{}:
				current = v
			default:
				return val
			}
		} else {
			return current
		}
	}

	return current
}

func SetNestedValue(data map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := data

	for i := 0; i < len(keys)-1; i++ {
		if _, ok := current[keys[i]]; !ok {
			current[keys[i]] = make(map[string]interface{})
		}
		current = current[keys[i]].(map[string]interface{})
	}

	lastKey := keys[len(keys)-1]
	current[lastKey] = value
}