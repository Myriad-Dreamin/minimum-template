package serial

import (
	"errors"
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/sugar"
	"os"
	"reflect"
)

type VirtualService struct {
	faz      *Category
	models   []interface{}
	name string
	fileName string
}

func (v *VirtualService) GetFilePath() string {
	return v.fileName
}

func (v *VirtualService) Name(name string) *VirtualService {
	v.name = name
	return v
}

func (v *VirtualService) GetFaz() *Category {
	return v.faz
}

func (v *VirtualService) GetName() string {
	return v.name
}

func (v *VirtualService) UseModel(models ...interface{}) *VirtualService {
	v.models = append(v.models, models)
	return v
}

func (v *VirtualService) CateOf(faz *Category) *VirtualService {
	v.faz = faz
	return v
}

func (v *VirtualService) ToFile(fileName string) ProposingService {
	v.fileName = fileName
	return v
}

type ProposingService interface {
	ToFile(fileName string) ProposingService
	GetName() string
	GetFilePath() string
	GetFaz() *Category
}

type PublishingServices struct {
	rawSvc      []ProposingService
	svc         []*serviceDescription
	packageName string
}

func (c *PublishingServices) SetPackageName(packageName string) {
	c.packageName = packageName
}


func NewService(rawSvc ...ProposingService) *PublishingServices {
	return &PublishingServices{
		rawSvc:rawSvc,
		packageName: "control",
	}
}



var cateType = reflect.TypeOf(new(Category))
func (c *PublishingServices) Publish() error {
	for _, svc := range c.rawSvc {
		ctx := &context{
			svc:      svc,
		}

		desc := &serviceDescription{
			packages: make(map[string]int),
		}
		desc.name = svc.GetName()
		desc.filePath = svc.GetFilePath()
		if len(desc.filePath) == 0 {
			return errors.New("desc.filePath needed")
		}
		faz := svc.GetFaz()
		value, svcType := reflect.ValueOf(svc), reflect.TypeOf(svc)
		for svcType.Kind() == reflect.Ptr {
			value, svcType = value.Elem(), svcType.Elem()
		}
		for i := 0; i < value.NumField(); i++ {
			field, fieldType := value.Field(i), svcType.Field(i)
			if fieldType.Type == cateType {
				cate := field.Interface().(*Category)
				if cate != nil {
					desc.categories = append(desc.categories, cate.create(faz, ctx))
				}
			}
		}
		for k := range ctx.packages {
			desc.packages[k] = 1
		}
		c.svc = append(c.svc, desc)
	}
	if err := c.writeToFiles(); err != nil {
		return err
	}
	return nil
}

func (c *PublishingServices) writeToFiles() (err error) {
	for i := range c.svc {
		svc := c.svc[i]
		fmt.Println(svc.filePath)
		sugar.WithWriteFile(func(f *os.File) {
			_, err = fmt.Fprintf(f, `
package %s

import (
%s
)

%s`, c.packageName, depList(svc.packages), svc.generateObjects())
		}, svc.filePath)
	}
	return
}

func depList(pkgSet map[string]int) (res string) {
	for k := range pkgSet {
		if len(k) > 0 {
			res += `    "`+ k +`"
`
		}
	}
	return
}