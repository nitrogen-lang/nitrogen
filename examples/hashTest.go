package main

import "fmt"

func main() {
	map1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range map1 {
		fmt.Println(k, v)
		if k == "key1" {
			map1["key4"] = "value4"
		}
	}
}
