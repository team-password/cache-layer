package cachelayer

import (
	"reflect"

	"github.com/team-password/cachelayer/schemas"
	"github.com/team-password/cachelayer/tag"
)

// KV key-value pair entity
type KV struct {
	// Key name
	Key string
	// Value value
	Value []byte
}

// ICache Cache abstraction
type ICache interface {
	// StoreAll Store all key-value pairs
	StoreAll(KVs ...KV) (err error)
	// Get value by key
	Get(key string) (value []byte, has bool, err error)
}

// IDB Database abstraction
type IDB interface {
	// GetEntry pass in an entity pointer, which will be used as the query condition.
	// After the function is executed, the result of the query will be written to the entity pointed to by the pointer.
	GetEntry(entry interface{}) (bool, error)
}

// IHandler handler abstraction
type IHandler interface {
	// GetEntry get a pointer to an entity type and return the entity
	GetEntry(entry interface{}) (bool, error)
}

var _ IHandler = Handler{}

// Handler Default handler
type Handler struct {
	cacheHandler ICache
	dbHandler    IDB
	serializer   Serializer
	log          Logger
}

// NewCacheHandler Create a Cache handler
func NewCacheHandler(cacheHandler ICache, dbHandler IDB, options ...OptionsFunc) *Handler {
	o := Options{}
	for _, option := range options {
		option(&o)
	}

	if o.cacheTagName != "" {
		tag.ConfigTag(o.cacheTagName)
	}
	if o.serializer == nil {
		o.serializer = JsonSerializer{}
	}
	if o.log == nil {
		o.log = DefaultLogger{}
	}
	schemas.ServiceName = o.serviceName

	return &Handler{
		cacheHandler: cacheHandler,
		dbHandler:    dbHandler,
		serializer:   o.serializer,
		log:          o.log,
	}
}

// GetEntry Get cached entity
func (c Handler) GetEntry(entry interface{}) (bool, error) {
	entryKey, err := schemas.GetEntryCacheKey(entry.(schemas.IEntry))
	if err != nil {
		return false, err
	}

	entryValue, has, err := c.cacheHandler.Get(entryKey)
	if err != nil {
		c.log.Error("Failed to get data from cache err:%v entryKey:%v", err.Error(), entryKey)
	}
	if has {
		err = c.serializer.Deserialize(entryValue, entry)
	}
	if !has {
		has, err = c.dbHandler.GetEntry(entry)
		if has {
			sliceValue := reflect.MakeSlice(reflect.SliceOf(reflect.Indirect(reflect.ValueOf(entry)).Type()), 0, 0)
			sliceValue = reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(entry)))
			c.storeCache(sliceValue.Interface())
		}

	}
	return has, err
}

// EntryCache Cache entity
type EntryCache struct {
	entry    interface{}
	entryKey string
}

func (c Handler) storeCache(entries interface{}) {
	entryCaches := make([]EntryCache, 0)
	entriesValue := reflect.Indirect(reflect.ValueOf(entries))
	for i := 0; i < entriesValue.Len(); i++ {
		entryKey, err := schemas.GetEntryCacheKey(entriesValue.Index(i).Interface().(schemas.IEntry))
		if err != nil {
			continue
		}
		entryCaches = append(entryCaches, EntryCache{
			entry:    entriesValue.Index(i).Interface().(schemas.IEntry),
			entryKey: entryKey,
		})
	}

	keyValues := make([]KV, 0)
	for _, entryCache := range entryCaches {
		value, err := c.serializer.Serialize(&entryCache.entry)
		if err != nil {
			c.log.Error("Failed serialize err:%v entry:%v", err, entryCache)
		}
		keyValues = append(keyValues, KV{
			Key:   entryCache.entryKey,
			Value: value,
		})
	}
	err := c.cacheHandler.StoreAll(keyValues...)
	if err != nil {
		c.log.Error("Failed StoreAll err:%v keyValues:%v", err, keyValues)
	}
}
