package pagination

import (
	"encoding/base64"
	"encoding/json"
)

type CursorField struct {
	Name  string
	Value string
}

type Cursor struct {
	Fields []CursorField
}

func (c *Cursor) AddField(name string, value string) {
	c.Fields = append(c.Fields, CursorField{Name: name, Value: value})
}

func (c *Cursor) B64Encode() string {
	if c == nil {
		return ""
	}
	fieldToValue := make(map[string]string)
	for _, field := range c.Fields {
		fieldToValue[field.Name] = field.Value
	}
	jsonBytes, _ := json.Marshal(fieldToValue)
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

func CursorB64Decode(encoded string) *Cursor {
	if encoded == "" {
		return nil
	}
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	fieldToValue := make(map[string]string)
	_ = json.Unmarshal(decoded, &fieldToValue)
	cursor := &Cursor{}
	for name, value := range fieldToValue {
		cursor.Fields = append(cursor.Fields, CursorField{Name: name, Value: value})
	}
	return cursor
}
