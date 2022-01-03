package tengohttp

import (
	"net/url"

	"github.com/d5/tengo/v2"
)

func (s *server) parse(args ...tengo.Object) (tengo.Object, error) {
	var a string
	switch len(args) {
	case 0:
		a = "http://" + s.r.Host + s.r.RequestURI
		fallthrough
	case 1:
		a, _ = tengo.ToString(args[0])
	default:
		return nil, tengo.ErrWrongNumArguments
	}
	if u, e := url.Parse(a); e != nil {
		return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
	} else {
		ret := &tengo.Map{Value: map[string]tengo.Object{
			"scheme":       &tengo.String{Value: u.Scheme},
			"opaque":       &tengo.String{Value: u.Opaque},
			"user":         &tengo.String{Value: u.User.Username()},
			"host":         &tengo.String{Value: u.Host},
			"path":         &tengo.String{Value: u.Path},
			"raw_path":     &tengo.String{Value: u.EscapedPath()},
			"query":        vals2map(u.Query()),
			"raw_query":    &tengo.String{Value: u.RawQuery},
			"fragment":     &tengo.String{Value: u.Fragment},
			"raw_fragment": &tengo.String{Value: u.EscapedFragment()},
		}}
		if p, o := u.User.Password(); o {
			ret.Value["pass"] = &tengo.String{Value: p}
		}
		return ret, nil
	}
}
func (s *server) resolve(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) == 1 {
		args = append(args, &tengo.String{Value: "http://" + s.r.Host + s.r.URL.Path})
	}
	if len(args) != 2 {
		err = tengo.ErrWrongNumArguments
	} else if ro, o := args[0].(*tengo.String); !o {
		err = tengo.ErrInvalidArgumentType{Name: "first", Expected: "string", Found: args[0].TypeName()}
	} else if bo, o := args[1].(*tengo.String); !o {
		err = tengo.ErrInvalidArgumentType{Name: "second", Expected: "string", Found: args[1].TypeName()}
	} else if bu, e := url.Parse(bo.Value); e != nil {
		ret = &tengo.Error{Value: &tengo.String{Value: e.Error()}}
	} else if ru, e := url.Parse(ro.Value); e != nil {
		ret = &tengo.Error{Value: &tengo.String{Value: e.Error()}}
	} else {
		ret = &tengo.String{Value: bu.ResolveReference(ru).String()}
	}
	return
}
func encode(args ...tengo.Object) (r tengo.Object, e error) {
	pth := false
	switch len(args) {
	case 2:
		pth = !args[1].IsFalsy()
		fallthrough
	case 1:
		switch a := args[0].(type) {
		case *tengo.String:
			if pth {
				r = &tengo.String{Value: url.PathEscape(a.Value)}
			} else {
				r = &tengo.String{Value: url.QueryEscape(a.Value)}
			}
		case *tengo.Map:
			r = &tengo.String{Value: url.Values(map2vals(a)).Encode()}
		default:
			e = tengo.ErrInvalidArgumentType{Name: "first", Expected: "string/map", Found: args[0].TypeName()}
		}
	default:
		e = tengo.ErrWrongNumArguments
	}
	return
}
func decode(args ...tengo.Object) (r tengo.Object, e error) {
	pth := false
	switch len(args) {
	case 2:
		pth = !args[1].IsFalsy()
		fallthrough
	case 1:
		if s, o := args[0].(*tengo.String); !o {
			e = tengo.ErrInvalidArgumentType{Name: "first", Expected: "string", Found: args[0].TypeName()}
		} else {
			var rs string
			if pth {
				rs, e = url.PathUnescape(s.Value)
			} else {
				rs, e = url.QueryUnescape(s.Value)
			}
			if e != nil {
				r, e = &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			} else {
				r = &tengo.String{Value: rs}
			}
		}
	default:
		e = tengo.ErrWrongNumArguments
	}
	return
}

func vals2map(vals map[string][]string) *tengo.Map {
	r := &tengo.Map{Value: make(map[string]tengo.Object)}
	for k, vs := range vals {
		a := &tengo.Array{}
		for _, v := range vs {
			a.Value = append(a.Value, &tengo.String{Value: v})
		}
		r.Value[k] = a
	}
	return r
}
func map2vals(m *tengo.Map) map[string][]string {
	vals := make(map[string][]string)
	for k, o := range m.Value {
		if a, i := o.(*tengo.Array); i {
			for _, v := range a.Value {
				s, _ := tengo.ToString(v)
				vals[k] = append(vals[k], s)
			}
		} else {
			s, _ := tengo.ToString(o)
			vals[k] = []string{s}
		}
	}
	return vals
}
