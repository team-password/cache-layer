package cachelayer

import (
	"encoding/json"
	"testing"

	"github.com/team-password/cachelayer/schemas"
)

// MemoryCacheHandler memory cache handler
type MemoryCacheHandler struct {
	data map[string][]byte
}

// StoreAll store key value
func (m MemoryCacheHandler) StoreAll(keyValues ...KV) (err error) {
	for _, keyValue := range keyValues {
		m.data[keyValue.Key] = keyValue.Value
	}
	return nil
}

// Get gets value by key
func (m MemoryCacheHandler) Get(key string) (data []byte, has bool, err error) {
	bytes, has := m.data[key]
	return bytes, has, nil
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
func NewMemoryCache() *Handler {
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
			c := Handler{
				cacheHandler: tt.fields.cacheHandler,
				dbHandler:    tt.fields.databaseHandler,
				serializer:   tt.fields.serializer,
				log:          tt.fields.log,
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
