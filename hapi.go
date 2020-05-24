package hapi

import (
	"log"
	"reflect"

	"github.com/labstack/echo"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type HAPI struct {
	e *echo.Echo
}

func New(e *echo.Echo) *HAPI {
	return &HAPI{e: e}
}

func (h *HAPI) GET(path string, f interface{}, m ...echo.MiddlewareFunc) *echo.Route {
	return h.e.GET(path, handle(f), m...)
}

func (h *HAPI) POST(path string, f interface{}, m ...echo.MiddlewareFunc) *echo.Route {
	return h.e.POST(path, handle(f), m...)
}

func handle(f interface{}) echo.HandlerFunc {
	rv := reflect.ValueOf(f)
	rt := reflect.TypeOf(f)

	return func(c echo.Context) (err error) {
		in := make([]reflect.Value, rt.NumIn())

		// req := c.Request()

		// bind request to handler input variables
		for i := 0; i < rt.NumIn(); i++ {
			arg := rt.In(i)

			switch {

			case arg.String() == "echo.Context":
				in[i] = reflect.ValueOf(c)

			case arg.String() == "http.ResponseWriter":
				in[i] = reflect.ValueOf(c.Response().Writer)

			case arg.String() == "*http.Request":
				in[i] = reflect.ValueOf(c.Request())

			case arg.String() == "*multipart.FileHeader":
				f, err := c.MultipartForm()
				if err != nil {
					log.Println(err)
					continue
				}

				for _, file := range f.File {
					in[i] = reflect.ValueOf(file[0])
					break
				}

				// in[i] =
			// case arg.Kind() == reflect.Ptr:
			// 	argPtr := reflect.New(arg)
			// 	argIntf := argPtr.Interface()

			// 	if err := c.Bind(argIntf); err != nil {
			// 		return err
			// 	}

			// 	in[i] = argPtr.Elem()

			// case arg.Kind() == reflect.Struct:
			// 	argPtr := reflect.New(arg)
			// 	argIntf := argPtr.Interface()

			// 	if err := c.Bind(argIntf); err != nil {
			// 		return err
			// 	}

			// 	in[i] = argPtr.Elem()

			default:
				argPtr := reflect.New(arg)
				argIntf := argPtr.Interface()

				if err := c.Bind(argIntf); err != nil {
					return err
				}

				in[i] = argPtr.Elem()
			}
		}

		// call handler
		out := rv.Call(in)

		// if handler returns nothing
		if len(out) == 0 {
			return nil
		}

		var (
			resp interface{}
			code = 200
			e    error
		)

		// bind returned variables to response
		for _, a := range out {
			val := a.Interface()
			if val == nil {
				continue
			}

			if e, ok := val.(error); ok {
				err = e
				continue
			}

			if v, ok := val.([]byte); ok {
				resp = v
				// rw.Write(data)
				continue
			}

			if v, ok := val.(string); ok {
				resp = v
				// fmt.Fprint(rw, data)
				continue
			}

			if v, ok := val.(int); ok {
				code = v
				// rw.WriteHeader(code)
				continue
			}

			resp = val
		}

		if e != nil {
			if code <= 200 {
				code = 500
			}
			return c.JSON(code, echo.Map{"error": e})
		}

		if resp != nil {
			switch resp.(type) {
			case []byte:
				_, err = c.Response().Write(resp.([]byte))
				return err
				// _, err = rw.Write(resp.([]byte))
			case string:
				return c.String(code, resp.(string))
				// _, err = fmt.Fprint(rw, resp.(string))
			default:
				return c.JSON(code, resp)
			}
		}

		return c.JSON(code, echo.Map{"status": "ok"})
	}
}
