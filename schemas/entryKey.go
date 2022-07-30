package schemas

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/team-password/cachelayer/tag"
)

// ServiceName Service Name
var ServiceName string

// EntryKey Cache key
type EntryKey struct {
	Name  string
	Param string
}

// EntryKeys Cache key array
type EntryKeys []EntryKey

// GetEntryCacheKey Get cache key by entity name
func (es EntryKeys) GetEntryCacheKey(entryName string) string {
	var (
		keyTemplate   = make([]string, 0)
		entryKeyNames = make([]interface{}, 0)
	)
	for _, e := range es {
		keyTemplate = append(keyTemplate, fmt.Sprintf("[%s", e.Name)+":%s]")
		entryKeyNames = append(entryKeyNames, e.Param)
	}
	return fmt.Sprintf(ServiceName+"_"+entryName+"#"+strings.Join(keyTemplate, "-"), entryKeyNames...)
}

// GetEntryCacheKey get the cache key by struct tag, if not set, find the field `id` or `key`
func GetEntryCacheKey(entry IEntry) (string, error) {
	var (
		entryKeys  EntryKeys = make([]EntryKey, 0)
		entryValue           = reflect.ValueOf(entry)
	)

	switch entryValue.Type().Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		entryValue = entryValue.Elem()
	}

	tagSortFields := tag.GetCacheTagFields(entry)

	if len(tagSortFields) == 0 {
		fieldValue := entryValue.FieldByNameFunc(func(fileName string) bool {
			if strings.ToLower(fileName) == "id" {
				return true
			}
			return false
		})
		if fieldValue != reflect.ValueOf(nil) {
			param := fieldValue.Interface()
			entryKeys = append(entryKeys, EntryKey{
				Name:  "id",
				Param: fmt.Sprint(param),
			})
			return entryKeys.GetEntryCacheKey(entryValue.Type().String()), nil
		}
		fieldValue = entryValue.FieldByNameFunc(func(fileName string) bool {
			if strings.ToLower(fileName) == "key" {
				return true
			}
			return false
		})
		if fieldValue != reflect.ValueOf(nil) {
			param := fieldValue.Interface()
			entryKeys = append(entryKeys, EntryKey{
				Name:  "key",
				Param: fmt.Sprint(param),
			})
			return entryKeys.GetEntryCacheKey(entryValue.Type().String()), nil
		}

		return "", errors.New("the field with the default value of Id and the cache tag was not found")
	}

	for _, tagField := range tagSortFields {
		fieldValue := entryValue.FieldByIndex(tagField.Index)
		entryKeys = append(entryKeys, EntryKey{
			Name:  tagField.Tag.Get(tag.GetName()),
			Param: fmt.Sprint(fieldValue.Interface()),
		})
	}

	return entryKeys.GetEntryCacheKey(entryValue.Type().String()), nil
}
