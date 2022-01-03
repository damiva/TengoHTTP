package tengohttp

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/d5/tengo/v2"
)

func (s *server) read(args ...tengo.Object) (tengo.Object, error) {
	var a string
	switch len(args) {
	case 2:
		a, _ = tengo.ToString(args[1])
		fallthrough
	case 1:
		var q *url.Values
		if args[0].IsFalsy() {
			q = &s.r.Form
		} else {
			q = &s.r.PostForm
		}
		if a == "" {
			return vals2map(*q), nil
		} else if q.Has(a) {
			return &tengo.String{Value: q.Get(a)}, nil
		}
	case 0:
		if b, e := ioutil.ReadAll(s.r.Body); e != nil {
			return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
		} else {
			return &tengo.Bytes{Value: b}, nil
		}
	default:
		return nil, tengo.ErrWrongNumArguments
	}
	return nil, nil
}
func (s *server) write(args ...tengo.Object) (tengo.Object, error) {
	c := 0
	for n, arg := range args {
		switch a := arg.(type) {
		case *tengo.Map:
			if s.h {
				return nil, tengo.ErrInvalidArgumentType{Name: strconv.Itoa(n) + "-th", Expected: "nor map or int", Found: a.TypeName()}
			}
			for k, vs := range map2vals(a) {
				for _, v := range vs {
					s.w.Header().Add(k, v)
				}
			}
		case *tengo.Int:
			if s.h {
				return nil, tengo.ErrInvalidArgumentType{Name: strconv.Itoa(n) + "-th", Expected: "nor map or int", Found: a.TypeName()}
			}
			s.h, c = true, int(a.Value)
		default:
			var e error
			if c > 0 {
				if v, o := tengo.ToString(a); !o {
					s.w.WriteHeader(c)
				} else if c < 300 {
					s.w.WriteHeader(c)
					_, e = s.w.Write([]byte(v))
				} else if c < 400 {
					http.Redirect(s.w, s.r, v, c)
				} else {
					http.Error(s.w, v, c)
				}
			} else if v, o := a.(*tengo.Bytes); o {
				_, e = s.w.Write(v.Value)
			} else if v, o := tengo.ToString(a); o {
				_, e = s.w.Write([]byte(v))
			}
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			} else {
				c, s.h = 0, true
			}
		}
	}
	return nil, nil
}
