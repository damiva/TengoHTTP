package tengohttp

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/d5/tengo/v2"
)

func (s *server) read(args ...tengo.Object) (tengo.Object, error) {
	var pst bool
	switch len(args) {
	case 2:
		pst = !args[1].IsFalsy()
		fallthrough
	case 1:
		switch v := args[0].(type) {
		case *tengo.String:
			if a := v.Value; a != "" {
				if pst {
					a = s.r.PostFormValue(a)
				} else {
					a = s.r.FormValue(a)
				}
				if a != "" {
					return &tengo.String{Value: a}, nil
				} else if pst && s.r.PostForm.Has(v.Value) || !pst && s.r.Form.Has(v.Value) {
					return &tengo.String{}, nil
				}
			}
		case *tengo.Bool:
			pst = !v.IsFalsy()
			if pst {
				if s.r.PostForm == nil {
					s.r.ParseForm()
				}
				return vals2map(s.r.PostForm), nil
			} else {
				if s.r.Form == nil {
					s.r.ParseForm()
				}
				return vals2map(s.r.Form), nil
			}
		default:
			return nil, tengo.ErrInvalidArgumentType{Name: "first", Expected: "string/bool", Found: args[0].TypeName()}
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
