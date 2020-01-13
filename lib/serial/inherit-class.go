package serial

type inheritClass struct {
	name string
	bases []interface{}
}

func (i inheritClass) String() string {
	panic("implement me")
	//return fmt.Sprintf("type %s = %s", i.name, reflect.TypeOf(i.base))
}

func (i inheritClass) GetTypeString() string {
	return i.name
}

func (i inheritClass) GetEmbedObject() []ObjectDescription {
	return nil
}

func (i inheritClass) CreateObjectDescription(ctx *Context) ObjectDescription {
	return i
}

func Inherit(name string, bases ...interface{}) *inheritClass {
	return &inheritClass{name:name, bases:bases}
}

