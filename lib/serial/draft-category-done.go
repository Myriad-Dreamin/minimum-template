package serial

// todo middleware
type Category interface {
	Path(path string) Category
	SubCate(path string, cat Category) Category
	DiveIn(path string) Category

	RawMethod(m ...Method) Category
	Method(m MethodType, descriptions ...interface{}) Category

	GetPath() string

	CreateCategoryDescription(ctx *Context) CategoryDescription
}

type category struct {
	path    string
	methods []Method
	subs    map[string]Category
}

func newCategory() *category {
	return new(category)
}

func Ink(_ ...interface{}) Category {
	return newCategory()
}

func (c *category) Path(path string) Category {
	c.path = path
	return c
}

func (c *category) SubCate(path string, cat Category) Category {
	if _, ok := c.subs[path]; ok {
		panic(ErrConflictPath)
	}
	c.subs[path] = cat
	return c
}

func (c *category) DiveIn(path string) Category {
	cat := &category{
		path: path,
	}
	c.SubCate(path, cat)
	return cat
}

func (c *category) GetPath() string {
	return c.path
}

func (c *category) RawMethod(m ...Method) Category {
	c.methods = append(c.methods, m...)
	return c
}

// todo
func (c *category) Method(m MethodType, descriptions ...interface{}) Category {
	method := newMethod(m)
	for _, description := range descriptions {
		switch desc := description.(type) {
		case string:
			method.name = desc
		case RequestObj:
			method.requests = append(method.requests, desc)
		case ReplyObj:
			method.requests = append(method.requests, desc)
		}
	}

	c.methods = append(c.methods, method)
	return c
}

func (c *category) CreateCategoryDescription(ctx *Context) CategoryDescription {
	desc := new(categoryDescription)
	for _, method := range c.methods {
		subCtx := ctx.sub()
		desc.methods = append(desc.methods, method.CreateMethodDescription(subCtx))
		desc.packages = inplaceMergePackage(desc.packages, subCtx.packages)
	}

	for _, sub := range c.subs {
		subDesc := sub.CreateCategoryDescription(ctx.sub())
		desc.subCates[subDesc.GetName()] = subDesc
	}
	return desc
}
