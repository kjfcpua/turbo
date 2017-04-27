package framework

import (
	"text/template"
	"bytes"
	"os"
)

func GenerateHandler() {
	var casesStr string
	for _, v := range UrlServiceMap {
		tmpl, err := template.New("cases").Parse(cases)
		if err != nil {
			panic(err)
		}
		var casesBuf bytes.Buffer
		err = tmpl.Execute(&casesBuf, method{configs["service_name"], v[2]})
		casesStr = casesStr + casesBuf.String()
	}
	tmpl, err := template.New("handler").Parse(handlerFunc)
	if err != nil {
		panic(err)
	}
	// TODO check if dir 'gen' exist, if not, create first
	f, _ := os.Create(serviceRootPath + "/gen/handler.go")
	err = tmpl.Execute(f, handlerContent{Cases: casesStr})
	if err != nil {
		panic(err)
	}

}

type method struct {
	ServiceName string
	MethodName  string
}

type handlerContent struct {
	Cases string
}

var handlerFunc string = `package gen

import (
	"reflect"
	"net/http"
	f "zx/demo/framework"
	"fmt"
)

var Handler = func(methodName string) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		switch methodName { {{.Cases}}
		default:
			resp.Write([]byte(fmt.Sprintf("No such grpc method[%s]", methodName)))
		}
	}
}`

var cases string = `
		case "{{.MethodName}}":
			f.ParseRequestForm(req)
			request := {{.MethodName}}Request{}
			theType := reflect.TypeOf(request)
			theValue := reflect.ValueOf(&request).Elem()
			fieldNum := theType.NumField()
			for i := 0; i < fieldNum; i++ {
				fieldName := theType.Field(i).Name
				v, ok := req.Form[f.ToSnakeCase(fieldName)]
				if ok && len(v) > 0 {
					theValue.FieldByName(fieldName).SetString(v[0])
				}
			}
			params := make([]reflect.Value, 2)
			params[0] = reflect.ValueOf(req.Context())
			params[1] = reflect.ValueOf(&request)
			result := reflect.ValueOf(f.GrpcService().({{.ServiceName}}Client)).MethodByName(methodName).Call(params)

			rsp := result[0].Interface().(*{{.MethodName}}Response)
			if result[1].Interface() == nil {
				resp.Write([]byte(rsp.String() + "\n"))
			} else {
				resp.Write([]byte(result[1].Interface().(error).Error() + "\n"))
			}
			`
