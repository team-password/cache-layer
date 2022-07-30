package schemas

import "reflect"

// IEntry Abstract entity type
type IEntry interface {
	TableName() string
}

// GetPKsByEntries Get the primary key of the abstract entity type
func GetPKsByEntries(entries interface{}) (PK, error) {
	pks := make(PK, 0)
	entriesElement := reflect.Indirect(reflect.ValueOf(entries))
	for i := 0; i < entriesElement.Len(); i++ {
		pk, err := GetEntryCacheKey(reflect.Indirect(entriesElement.Index(i)).Interface().(IEntry))
		if err != nil {
			return pks, err
		}
		pks = append(pks, pk)
	}
	return pks, nil
}
