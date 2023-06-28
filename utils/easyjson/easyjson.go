package easyjson

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Get 返回给定 JSON 字符串中键的值
func Get(jsonStr string, key string) (interface{}, error) {
	// 先将 json 字符串解析为通用接口类型
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	// 根据 / 分割要查询的 key
	keys := strings.Split(key, "/")
	// 依次使用分割后的 key 遍历接口
	for _, k := range keys {
		// 如果是一个数组索引，将字符串转换为整数
		index, err := strconv.Atoi(k)
		if err == nil {
			v := data.([]interface{})[index]
			data = v
		} else {
			v := data.(map[string]interface{})[k]
			data = v
		}
	}
	return data, nil
}
