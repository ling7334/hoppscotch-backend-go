package graph

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalDateTimeScalar convert a time.Time to DateTime Scalar
func MarshalDateTimeScalar(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(strconv.Quote(t.Format(time.RFC3339))))
	})
}

// UnmarshalDateTimeScalar try to convert a string or a int object to a time.Time
func UnmarshalDateTimeScalar(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(time.RFC3339, v)
	case int:
		return time.Unix(int64(v), 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("%T is not a time.Time", v)
	}
}
