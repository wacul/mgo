package txn

import (
	"bytes"
	"fmt"
	"labix.org/v2/mgo/bson"
	"os"
	"sort"
	"sync/atomic"
)

var debugEnabled bool

// SetDebug enables or disables debugging.
func SetDebug(debug bool) {
	debugEnabled = debug
}

var ErrChaos = fmt.Errorf("interrupted by chaos")

var debugId uint32

func debugPrefix() string {
	d := atomic.AddUint32(&debugId, 1)-1
	s := make([]byte, 0, 10)
	for i := uint(0); i < 8; i++ {
		s = append(s, "abcdefghijklmnop"[(d>>(4*i))&0xf])
		if d>>(4*(i+1)) == 0 {
			break
		}
	}
	s = append(s, ')', ' ')
	return string(s)
}

func debugf(format string, args ...interface{}) {
	if !debugEnabled {
		return
	}
	for i, arg := range args {
		switch v := arg.(type) {
		case bson.ObjectId:
			args[i] = v.Hex()
		case []bson.ObjectId:
			lst := make([]string, len(v))
			for j, id := range v {
				lst[j] = id.Hex()
			}
			args[i] = lst
		case map[docKey][]bson.ObjectId:
			buf := &bytes.Buffer{}
			var dkeys docKeys
			for dkey := range v {
				dkeys = append(dkeys, dkey)
			}
			sort.Sort(dkeys)
			for i, dkey := range dkeys {
				if i > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteString(fmt.Sprintf("%v: {", dkey))
				for j, id := range v[dkey] {
					if j > 0 {
						buf.WriteByte(' ')
					}
					buf.WriteString(id.Hex())
				}
				buf.WriteByte('}')
			}
			args[i] = buf.String()
		case map[docKey][]int64:
			buf := &bytes.Buffer{}
			var dkeys docKeys
			for dkey := range v {
				dkeys = append(dkeys, dkey)
			}
			sort.Sort(dkeys)
			for i, dkey := range dkeys {
				if i > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteString(fmt.Sprintf("%v: %v", dkey, v[dkey]))
			}
			args[i] = buf.String()
		}
	}
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}