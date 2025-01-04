package pagination

import "fmt"

type Paginator struct {
	Order  *Order
	Filter *Filter
}

func (p *Paginator) EncodeCursor(i interface{}) string {
	return p.Order.EncodeCursor(i)
}

func (p *Paginator) DecodeCursor(encoded string) *Cursor {
	return p.Order.DecodeCursor(encoded)
}

func (p *Paginator) ContributeCursorToFilter(cursor *Cursor, isAfter bool) {
	o := p.Order
	filter := p.Filter
	if o == nil {
		return
	}
	orderFieldsByName := make(map[string]OrderField)
	for _, field := range o.Fields {
		orderFieldsByName[field.Name()] = field
	}
	for _, cursorField := range cursor.Fields {
		orderField, ok := orderFieldsByName[cursorField.Name]
		if !ok {
			continue
		}
		if isAfter {
			if orderField.Asc() {
				filter.AddFilterField(cursorField.Name, FilterOperatorGt, cursorField.Value)
			} else {
				filter.AddFilterField(cursorField.Name, FilterOperatorLt, cursorField.Value)
			}
		} else {
			if orderField.Asc() {
				filter.AddFilterField(cursorField.Name, FilterOperatorLt, cursorField.Value)
			} else {
				filter.AddFilterField(cursorField.Name, FilterOperatorGt, cursorField.Value)
			}
		}
	}
}

func (p *Paginator) ContributeCursorStringToFilter(cursorString *string, isAfter bool) error {
	if cursorString == nil {
		return nil
	}
	cursor := p.DecodeCursor(*cursorString)
	if cursor == nil {
		return fmt.Errorf("invalid cursor")
	}
	p.ContributeCursorToFilter(cursor, isAfter)
	return nil
}
