package string

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm"
)

func init() {
	vm.RegisterNative("stdlib.http.do", do)
	vm.RegisterNative("stdlib.http.canonicalHeaderKey", canonicalHeaderKey)
}

func do(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("http.do", 2, args...); ac != nil {
		return ac
	}

	// Argument 1 - HTTP method
	methodObj, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("http.do expected first argument to be a string, got %s", args[0].Type().String())
	}
	method := strings.ToUpper(methodObj.String())
	if method == "" {
		method = "GET"
	}

	// Argument 2 - URL
	urlObj, ok := args[1].(*object.String)
	if !ok {
		return object.NewException("http.do expected second argument to be a string, got %s", args[1].Type().String())
	}

	url := strings.TrimSpace(urlObj.String())
	if url == "" {
		return object.NewException("http.do expected a non-empty string")
	}

	// Argument 3 - Data payload
	data := ""
	if len(args) >= 3 && args[2] != object.NullConst {
		dataObj, ok := args[2].(*object.String)
		if !ok {
			return object.NewException("http.do expected third argument to be a string, got %s", args[2].Type().String())
		}
		data = dataObj.String()
	}

	// Argument 4 - Request options
	var optionsObj *object.Hash
	if len(args) >= 4 && args[3] != object.NullConst {
		dataObj, ok := args[3].(*object.Hash)
		if !ok {
			return object.NewException("http.post expected fourth argument to be a map, got %s", args[3].Type().String())
		}
		optionsObj = dataObj
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return object.NewException("error making HTTP request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if optionsObj != nil {
		headers := optionsObj.LookupKey("headers")
		if headers != nil {
			headersMap, ok := headers.(*object.Hash)
			if !ok {
				return object.NewException("headers option must be a map")
			}

			for _, pair := range headersMap.Pairs {
				if pair.Key.Type() != object.StringObj {
					continue
				}
				if pair.Value.Type() != object.StringObj {
					continue
				}
				req.Header.Set(pair.Key.(*object.String).String(), pair.Value.(*object.String).String())
			}
		}

		tlsVerify := optionsObj.LookupKey("tls_verify")
		if tlsVerify != nil {
			verifyBool, ok := tlsVerify.(*object.Boolean)
			if !ok {
				return object.NewException("tls_verify option must be a boolean")
			}

			if !verifyBool.Value {
				client.Transport = &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}
			}
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return object.NewException("error making HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	return buildReturnValue(resp)
}

func buildReturnValue(resp *http.Response) object.Object {
	body, _ := ioutil.ReadAll(resp.Body)

	headers := make(map[string]string, len(resp.Header))

	for name, values := range resp.Header {
		headers[name] = strings.Join(values, ", ")
	}

	hash := &object.Hash{
		Pairs: make(map[object.HashKey]object.HashPair, 2),
	}

	hash.SetKey("body", object.MakeStringObj(string(body)))
	hash.SetKey("headers", object.StringMapToHash(headers))

	return hash
}

func canonicalHeaderKey(interpreter object.Interpreter, env *object.Environment, args ...object.Object) object.Object {
	if ac := moduleutils.CheckMinArgs("http.canonicalHeaderKey", 1, args...); ac != nil {
		return ac
	}

	headerKey, ok := args[0].(*object.String)
	if !ok {
		return object.NewException("http.canonicalHeaderKey expected a string argument, got %s", args[0].Type().String())
	}

	return object.MakeStringObj(http.CanonicalHeaderKey(headerKey.String()))
}
