package cachelayer

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/team-password/cachelayer/schemas"
)

// MemoryCacheHandler memory cache handler
type MemoryCacheHandler struct {
	data map[string][]byte
}

// StoreAll Store key value
func (m MemoryCacheHandler) StoreAll(keyValues ...KV) (err error) {
	for _, keyValue := range keyValues {
		m.data[keyValue.Key] = keyValue.Value
	}
	return nil
}

// Get Get value by key
func (m MemoryCacheHandler) Get(key string) (data []byte, has bool, err error) {
	bytes, has := m.data[key]
	return bytes, has, nil
}

// GetAll Get values by keys
func (m MemoryCacheHandler) GetAll(keys schemas.PK) (data []KV, err error) {
	returnKeyValues := make([]KV, 0)
	for _, key := range keys {
		bytes, _ := m.data[key]
		returnKeyValues = append(returnKeyValues, KV{
			Key:   key,
			Value: bytes,
		})
	}
	return returnKeyValues, nil
}

// DeleteAll Delete all key caches
func (m MemoryCacheHandler) DeleteAll(keys schemas.PK) error {
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

var data = make(map[string][]byte, 0)

// NewMemoryCacheHandler Create a cache handler
func NewMemoryCacheHandler() *MemoryCacheHandler {
	return &MemoryCacheHandler{
		data: data,
	}
}

// MemoryDb Memory Database
type MemoryDb struct {
}

// NewMemoryDb Create Memory Database
func NewMemoryDb() *MemoryDb {
	return &MemoryDb{}
}

// GetEntries Get the list of entities through sql
func (m MemoryDb) GetEntries(entries interface{}, sql string) error {
	if sql == "SELECT * FROM public_relation  WHERE  relateId = 1 AND sourceId = 2 AND propertyId = 3 ;" || sql == "SELECT * FROM public_relation  WHERE (  relateId = 1 AND sourceId = 2 AND propertyId = 3  );" {
		mockEntries := make([]MockEntry, 0)
		mockEntries = append(mockEntries, MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		})
		marshal, _ := json.Marshal(mockEntries)
		json.Unmarshal(marshal, entries)
		return nil
	} else if sql == "SELECT * FROM public_relation  WHERE  relateId = 1 AND sourceId = 2;" {
		mockEntries := make([]*MockEntry, 0)
		mockEntries = append(mockEntries, &MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		})
		mockEntries = append(mockEntries, &MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 4,
		})
		marshal, _ := json.Marshal(mockEntries)
		json.Unmarshal(marshal, entries)
		return nil
	} else if sql == "SELECT * FROM public_relation  WHERE (  relateId = 1 AND sourceId = 2 AND propertyId = 3  ) OR (  relateId = 1 AND sourceId = 2 AND propertyId = 4  );" {
		mockEntries := make([]MockEntry, 0)
		mockEntries = append(mockEntries, MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 3,
		})
		mockEntries = append(mockEntries, MockEntry{
			RelateId:   1,
			SourceId:   2,
			PropertyId: 4,
		})
		marshal, _ := json.Marshal(mockEntries)
		json.Unmarshal(marshal, entries)
		return nil
	}
	return errors.New("mockEntries not found")
}

// GetEntry Get entities through sql
func (m MemoryDb) GetEntry(entry interface{}) (bool, error) {
	mockEntry := &MockEntry{
		RelateId:   1,
		SourceId:   2,
		PropertyId: 3,
	}
	marshal, _ := json.Marshal(mockEntry)
	json.Unmarshal(marshal, entry)
	return true, nil
}

// NewMemoryCache Create a memory cache
func NewMemoryCache() *CacheHandler {
	return NewCacheHandler(NewMemoryCacheHandler(), NewMemoryDb())
}

// MockEntry Mock entity
type MockEntry struct {
	RelateId   int64 `cache:"relateId"`
	SourceId   int64 `cache:"sourceId"`
	PropertyId int64 `cache:"propertyId"`
}

// TableName Table Name
func (m MockEntry) TableName() string {
	return "public_relation"
}

func TestCacheHandler_GetEntry(t *testing.T) {
	type fields struct {
		cacheHandler    ICache
		databaseHandler IDB
		serializer      Serializer
		log             Logger
	}
	type args struct {
		entry interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				serializer:      JsonSerializer{},
				log:             DefaultLogger{},
			},
			args: args{
				entry: MockEntry{
					RelateId:   1,
					SourceId:   2,
					PropertyId: 3,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CacheHandler{
				cacheHandler:    tt.fields.cacheHandler,
				databaseHandler: tt.fields.databaseHandler,
				serializer:      tt.fields.serializer,
				log:             tt.fields.log,
			}
			got, err := c.GetEntry(tt.args.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEntry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type TestLogger struct {
}

func (t TestLogger) Info(format string, a ...interface{}) {
	panic("implement me")
}

func (t TestLogger) Error(format string, a ...interface{}) {
	panic("implement me")
}

func (t TestLogger) Debug(format string, a ...interface{}) {
	panic("implement me")
}

func (t TestLogger) Warn(format string, a ...interface{}) {
	panic("implement me")
}

type TestSerializer struct {
}

func (t TestSerializer) Serialize(value interface{}) ([]byte, error) {
	panic("implement me")
}

func (t TestSerializer) Deserialize(data []byte, ptr interface{}) error {
	panic("implement me")
}

func TestNewCacheHandler(t *testing.T) {
	type args struct {
		cacheHandler    ICache
		databaseHandler IDB
		options         []OptionsFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				options:         []OptionsFunc{WithServiceName("test")},
			},
		},
		{
			name: "",
			args: args{
				cacheHandler:    NewMemoryCacheHandler(),
				databaseHandler: NewMemoryDb(),
				options:         []OptionsFunc{WithServiceName("test"), WithSerializer(TestSerializer{}), WithCacheTagName("test"), WithLogger(TestLogger{})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCacheHandler(tt.args.cacheHandler, tt.args.databaseHandler, tt.args.options...); !(schemas.ServiceName == "test") {
				t.Errorf("NewCacheHandler() = %v", got)
			}
		})
	}
}
