package reflext

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// StructTag :
type StructTag struct {
	name string
	opts map[string]string
}

// Get :
func (st StructTag) Get(key string) string {
	return st.opts[key]
}

// LookUp :
func (st StructTag) LookUp(key string) (val string, exist bool) {
	val, exist = st.opts[key]
	return
}

// StructField :
type StructField struct {
	ID         string
	Index      []int
	Name       string
	Path       string
	IsNullable bool
	Zero       reflect.Value
	Tag        StructTag
	Embedded   bool
	Parent     *StructField
	Children   []*StructField
}

// Struct :
type Struct struct {
	Tree       *StructField
	Fields     Fields         // all fields belong to this struct
	NameFields []*StructField // names fields in sequence
	Indexes    map[string]*StructField
	Names      map[string]*StructField
}

// GetByTraversal :
func (s Struct) GetByTraversal(index []int) *StructField {
	if len(index) == 0 {
		return nil
	}

	tree := s.Tree
	for _, i := range index {
		if i >= len(tree.Children) || tree.Children[i] == nil {
			return nil
		}
		tree = tree.Children[i]
	}
	return tree
}

// Fields :
type Fields []*StructField

func (x Fields) Len() int { return len(x) }

func (x Fields) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x Fields) Less(i, j int) bool {
	for k, xik := range x[i].Index {
		if k >= len(x[j].Index) {
			return false
		}
		if xik != x[j].Index[k] {
			return xik < x[j].Index[k]
		}
	}
	return len(x[i].Index) < len(x[j].Index)
}

type typeQueue struct {
	t  reflect.Type
	sf *StructField
	pp string // parent path
}

func getCodec(t reflect.Type, tagName string, mapFunc MapFunc) *Struct {
	fields := make([]*StructField, 0)

	root := &StructField{}
	queue := []typeQueue{}
	queue = append(queue, typeQueue{Deref(t), root, ""})

	for len(queue) > 0 {
		q := queue[0]

		q.sf.Children = make([]*StructField, 0)

		for i := 0; i < q.t.NumField(); i++ {
			f := q.t.Field(i)

			// skip unexported fields
			if len(f.PkgPath) != 0 && !f.Anonymous {
				continue
			}

			tag := parseTag(f.Tag, tagName)
			if tag.name == "-" {
				continue
			}

			sf := &StructField{
				ID:         strings.TrimLeft(q.sf.ID+"."+strconv.Itoa(i), "."),
				Name:       f.Name,
				Path:       tag.name,
				IsNullable: q.sf.IsNullable || IsNullable(f.Type),
				Zero:       reflect.Zero(f.Type),
				Tag:        tag,
				Parent:     q.sf,
				Children:   make([]*StructField, 0),
			}

			if q.pp == "" {
				if sf.Path == "" {
					sf.Path = sf.Name
				}
			} else {
				sf.Path = q.pp + "." + sf.Name
			}

			ft := Deref(f.Type)
			q.sf.Children = append(q.sf.Children, sf)
			sf.Index = appendSlice(q.sf.Index, i)
			sf.Embedded = ft.Kind() == reflect.Struct && f.Anonymous

			if mapFunc != nil {
				if mapFunc(sf) {
					fields = append(fields, sf)
					continue
				}
			}

			if ft.Kind() == reflect.Struct {
				// embedded struct
				if f.Anonymous {
					path := sf.Path
					if sf.Tag.name == "" {
						path = q.pp
					}
					queue = append(queue, typeQueue{ft, sf, path})
				} else {
					queue = append(queue, typeQueue{ft, sf, sf.Path})
				}
			}

			fields = append(fields, sf)
		}

		queue = queue[1:]
	}

	codec := &Struct{
		Tree:       root,
		Fields:     fields,
		NameFields: make([]*StructField, 0, len(fields)),
		Indexes:    make(map[string]*StructField),
		Names:      make(map[string]*StructField),
	}

	sort.Sort(codec.Fields)

	for _, sf := range codec.Fields {
		codec.Indexes[sf.ID] = sf
		if sf.Name != "" && !sf.Embedded {
			codec.Names[sf.Path] = sf
			if len(sf.Index) > 1 && sf.Parent != nil && !sf.Parent.Embedded {
				continue
			}
			codec.NameFields = append(codec.NameFields, sf)
		}
	}

	return codec
}

func appendSlice(s []int, i int) []int {
	x := make([]int, len(s)+1)
	copy(x, s)
	x[len(x)-1] = i
	return x
}

func parseTag(tag reflect.StructTag, tagName string) (st StructTag) {
	parts := strings.Split(tag.Get(tagName), ",")
	st.name = strings.TrimSpace(parts[0])
	st.opts = make(map[string]string)
	if len(parts) > 1 {
		for _, opt := range parts[1:] {
			opt = strings.TrimSpace(opt)
			if strings.Contains(opt, ":") {
				kv := strings.SplitN(opt, ":", 2)
				k := strings.TrimSpace(strings.ToLower(kv[0]))
				st.opts[k] = strings.TrimSpace(kv[1])
				continue
			}
			opt = strings.ToLower(opt)
			st.opts[opt] = ""
		}
	}
	return
}