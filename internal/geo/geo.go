package geo

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	orbjson "github.com/paulmach/orb/geojson"
)

type Feature struct {
	Id         any            `json:"id,omitempty"`
	Type       string         `json:"type"`
	Geometry   orb.Geometry   `json:"geometry"`
	Properties map[string]any `json:"properties"`
}

var (
	_ json.Marshaler   = (*Feature)(nil)
	_ json.Unmarshaler = (*Feature)(nil)
)

func (f *Feature) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":       "Feature",
		"geometry":   orbjson.NewGeometry(f.Geometry),
		"properties": f.Properties,
	}
	if f.Id != nil {
		m["id"] = f.Id
	}
	return json.Marshal(m)
}

type jsonFeature struct {
	Id         any             `json:"id,omitempty"`
	Type       string          `json:"type"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties map[string]any  `json:"properties"`
}

var rawNull = json.RawMessage([]byte("null"))

func isRawNull(raw json.RawMessage) bool {
	if len(raw) != len(rawNull) {
		return false
	}
	for i, c := range raw {
		if c != rawNull[i] {
			return false
		}
	}
	return true
}

func (f *Feature) UnmarshalJSON(data []byte) error {
	jf := &jsonFeature{}
	if err := json.Unmarshal(data, jf); err != nil {
		return err
	}

	f.Type = jf.Type
	f.Id = jf.Id
	f.Properties = jf.Properties

	if isRawNull(jf.Geometry) {
		return nil
	}
	geometry := &orbjson.Geometry{}
	if err := json.Unmarshal(jf.Geometry, geometry); err != nil {
		return err
	}

	f.Geometry = geometry.Geometry()
	return nil
}

const (
	EncodingWKB = "WKB"
	EncodingWKT = "WKT"
)

func DecodeGeometry(value any, encoding string) (*orbjson.Geometry, error) {
	if value == nil {
		return nil, nil
	}
	if encoding == "" {
		if _, ok := value.([]byte); ok {
			encoding = EncodingWKB
		} else if _, ok := value.(string); ok {
			encoding = EncodingWKT
		}
	}
	if encoding == EncodingWKB {
		data, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("expected bytes for wkb geometry, got %T", value)
		}
		g, err := wkb.Unmarshal(data)
		if err != nil {
			return nil, err
		}
		return orbjson.NewGeometry(g), nil
	}
	if encoding == EncodingWKT {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for wkt geometry, got %T", value)
		}
		g, err := wkt.Unmarshal(str)
		if err != nil {
			return nil, err
		}
		return orbjson.NewGeometry(g), nil
	}
	return nil, fmt.Errorf("unsupported encoding: %s", encoding)
}

type CollectionInfo struct {
	mutex *sync.RWMutex
	minX  float64
	maxX  float64
	minY  float64
	maxY  float64
	types map[string]bool
}

func NewCollectionInfo(concurrent bool) *CollectionInfo {
	var mutex *sync.RWMutex
	if concurrent {
		mutex = &sync.RWMutex{}
	}
	return &CollectionInfo{
		mutex: mutex,
		types: map[string]bool{},
		minX:  math.MaxFloat64,
		maxX:  -math.MaxFloat64,
		minY:  math.MaxFloat64,
		maxY:  -math.MaxFloat64,
	}
}

func (i *CollectionInfo) writeLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Lock()
}

func (i *CollectionInfo) writeUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Unlock()
}

func (i *CollectionInfo) readLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RLock()
}

func (i *CollectionInfo) readUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RUnlock()
}

func (i *CollectionInfo) AddBounds(bounds *orb.Bound) {
	i.writeLock()
	minPoint := bounds.Min
	minX := minPoint[0]
	minY := minPoint[1]
	maxPoint := bounds.Max
	maxX := maxPoint[0]
	maxY := maxPoint[1]
	i.minX = math.Min(i.minX, minX)
	i.maxX = math.Max(i.maxX, maxX)
	i.minY = math.Min(i.minY, minY)
	i.maxY = math.Max(i.maxY, maxY)
	i.writeUnlock()
}

func (i *CollectionInfo) Bounds() *orb.Bound {
	i.readLock()
	bounds := &orb.Bound{
		Min: orb.Point{i.minX, i.minY},
		Max: orb.Point{i.maxX, i.maxY},
	}
	i.readUnlock()
	return bounds
}

func (i *CollectionInfo) AddType(typ string) {
	i.writeLock()
	i.types[typ] = true
	i.writeUnlock()
}

func (i *CollectionInfo) AddTypes(types []string) {
	i.writeLock()
	for _, typ := range types {
		i.types[typ] = true
	}
	i.writeUnlock()
}

func (i *CollectionInfo) Types() []string {
	i.readLock()
	types := []string{}
	for typ, ok := range i.types {
		if ok {
			types = append(types, typ)
		}
	}
	i.readUnlock()
	return types
}

type DatasetInfo struct {
	mutex       *sync.RWMutex
	collections map[string]*CollectionInfo
}

func NewDatasetInfo(concurrent bool) *DatasetInfo {
	var mutex *sync.RWMutex
	if concurrent {
		mutex = &sync.RWMutex{}
	}
	return &DatasetInfo{
		mutex:       mutex,
		collections: map[string]*CollectionInfo{},
	}
}

func (i *DatasetInfo) writeLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Lock()
}

func (i *DatasetInfo) writeUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.Unlock()
}

func (i *DatasetInfo) readLock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RLock()
}

func (i *DatasetInfo) readUnlock() {
	if i.mutex == nil {
		return
	}
	i.mutex.RUnlock()
}

func (i *DatasetInfo) NumCollections() int {
	i.readLock()
	num := len(i.collections)
	i.readUnlock()
	return num
}

func (i *DatasetInfo) AddCollection(name string) {
	i.writeLock()
	i.collections[name] = NewCollectionInfo(i.mutex != nil)
	i.writeUnlock()
}

func (i *DatasetInfo) HasCollection(name string) bool {
	i.readLock()
	_, has := i.collections[name]
	i.readUnlock()
	return has
}

func (i *DatasetInfo) AddBounds(name string, bounds *orb.Bound) {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	collection.AddBounds(bounds)
}

func (i *DatasetInfo) Bounds(name string) *orb.Bound {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	return collection.Bounds()
}

func (i *DatasetInfo) AddTypes(name string, types []string) {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	collection.AddTypes(types)
}

func (i *DatasetInfo) Types(name string) []string {
	i.readLock()
	collection := i.collections[name]
	i.readUnlock()
	return collection.Types()
}
