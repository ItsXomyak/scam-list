package postgres

import (
	"encoding/json"
	"fmt"
)

func toJSONArray(s [][]byte) ([]byte, error) {
	// быстрая валидация
	for i := range s {
		if !json.Valid(s[i]) {
			return nil, fmt.Errorf("element %d is not valid JSON", i)
		}
	}
	// Превращаем в []json.RawMessage и маршалим в единый JSON-массив
	arr := make([]json.RawMessage, len(s))
	for i := range s {
		arr[i] = json.RawMessage(s[i])
	}

	return json.Marshal(arr) // вернёт []byte вида: [ {...}, {...}, ... ]
}
