package mappers

import "fmt"

type uniqueIds map[string]struct{}

func (u uniqueIds) exists(id string, prefix string) bool {
	_id := fmt.Sprintf("%s%s", prefix, id)
	_, ok := u[_id]
	return ok
}

func (u uniqueIds) add(id string, prefix string) {
	_id := fmt.Sprintf("%s%s", prefix, id)
	u[_id] = struct{}{}
}
