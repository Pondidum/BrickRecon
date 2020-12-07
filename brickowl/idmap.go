package brickowl

import "encoding/json"

type idMap map[string]string

func (idmap *idMap) UnmarshalJSON(b []byte) error {

	var ids []struct {
		ID   string
		Type string
	}

	if err := json.Unmarshal(b, &ids); err != nil {
		return err
	}

	m := idMap{}

	for _, pair := range ids {
		m[pair.Type] = pair.ID
	}

	*idmap = m

	return nil
}
