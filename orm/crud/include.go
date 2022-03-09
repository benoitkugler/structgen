package crud

// code included in the generated CRUD Go code
const utils = `
func loadJSON(out interface{}, src interface{}) error {
	if src == nil {
		return nil //zero value out
	}
	bs, ok := src.([]byte)
	if !ok {
		return errors.New("not a []byte")
	}
	return json.Unmarshal(bs, out)
}

func dumpJSON(s interface{}) (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return driver.Value(string(b)), nil
}

// Set is a set of IDs.
type Set map[int64]bool

func NewSet() Set {
	return map[int64]bool{}
}

// NewSetFromSlice returns a set of unique IDs
func NewSetFromSlice(keys []int64) Set {
	out := make(Set, len(keys))
	for _, key := range keys {
		out[key] = true
	}
	return out
}

// Keys return the IDs contained in the set, as a slice.
func (s Set) Keys() []int64 {
	out := make([]int64, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	return out
}

func (s Set) Has(key int64) bool {
	_, has := s[key]
	return has
}

func (s Set) Add(key int64) {
	s[key] = true
}

type IDs []int64

func (ids IDs) AsSQL() pq.Int64Array {
	return pq.Int64Array(ids)
}

func (ids IDs) AsSet() Set {
	return NewSetFromSlice(ids)
}

// ScanIDs scans the result of a query returning a
// list of IDs.
func ScanIDs(rs *sql.Rows) (IDs, error) {
	defer rs.Close()
	ints := make(IDs, 0, 16)
	var err error
	for rs.Next() {
		var s int64
		if err = rs.Scan(&s); err != nil {
			return nil, err
		}
		ints = append(ints, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return ints, nil
}
`
