package huffman

import (
    "fmt"
    "encoding/json"
)

func mapString(dict map[byte]uint32) string {
    jsonDict, err := json.MarshalIndent(dict, "", "  ")
    if err != nil {
        fmt.Println("error:", err)
    }
    return string(jsonDict)
}

func encode(data []byte) []byte {
    byteMap := make(map[byte]uint32)
    for _, b := range data {
        if val, ok := byteMap[b]; ok {
            byteMap[b] = val + 1
        } else {
            byteMap[b] = 0
        }
    }
    fmt.Println(mapString(byteMap))
    return data
}

func decode(data []byte) []byte {
    return data
}

