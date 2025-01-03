package pagination

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
)

type OrderField interface {
	Name() string
	Asc() bool
	getValue(object interface{}) string
}

type BaseOrderField struct {
	FieldName string
	FieldAsc  bool
}

func (b *BaseOrderField) Name() string {
	return b.FieldName
}

func (b *BaseOrderField) Asc() bool {
	return b.FieldAsc
}

type Order struct {
	Fields []OrderField
}

func (o *Order) checksum() string {
	if o == nil {
		return ""
	}
	fieldToValue := make(map[string]bool)
	for _, field := range o.Fields {
		fieldToValue[field.Name()] = field.Asc()
	}
	jsonBytes, _ := json.Marshal(fieldToValue)
	h := fnv.New32a()
	h.Write(jsonBytes)
	sum := h.Sum32()
	return fmt.Sprintf("%x", sum)
}

func (o *Order) EncodeCursor(i interface{}) string {
	if o == nil {
		return ""
	}
	if o.checksum() == "" {
		return ""
	}
	cursor := &Cursor{}
	for _, field := range o.Fields {
		cursor.AddField(field.Name(), field.getValue(i))
	}
	return o.checksum() + "|" + cursor.B64Encode()
}

func (o *Order) DecodeCursor(encoded string) *Cursor {
	if encoded == "" {
		return nil
	}
	parts := strings.SplitN(encoded, "|", 2)
	if len(parts) != 2 {
		return nil
	}
	if parts[0] != o.checksum() {
		return nil
	}
	return CursorB64Decode(parts[1])
}
