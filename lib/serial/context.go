package serial

import (
	"reflect"
)

type Context struct {
	vars map[string]interface{}
	method *Method
	svc ProposingService
	sources map[uintptr]*source
	packages map[string]int
}

func (c *Context) appendPackage(pkg string) {
	if c.packages == nil {
		c.packages = make(map[string]int)
	}
	c.packages[pkg] = 1
}

func (c *Context) set(k string, v interface{}) {
	if c.vars == nil {
		c.vars = make(map[string]interface{})
	}
	c.vars[k] = v
}

func (c *Context) get(k string) (v interface{}) {
	if c.vars != nil {
		v, _ = c.vars[k]
	}
	return
}

func (c *Context) getSource(ptr uintptr) *source {
	s, _ := c.sources[ptr]
	return s
}

func (c *Context) makeSources() {
	c.sources = make(map[uintptr]*source)
	models := c.svc.GetModels()
	for _, xmodel := range models {
		v, t := reflect.ValueOf(xmodel.refer).Elem(), reflect.TypeOf(xmodel.refer).Elem()
		tt := t
		for t.Kind() == reflect.Ptr {
			v, t = v.Elem(), t.Elem()
		}
		c.appendPackage(t.PkgPath())
		for i := 0; i < t.NumField(); i++ {
			c.sources[v.Addr().Pointer() + t.Field(i).Offset] = &source{
				modelName: xmodel.name, faz: tt, fazElem: t, fieldIndex: i}
		}
	}
	//fmt.Println(c.sources)
}



