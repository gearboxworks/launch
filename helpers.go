package main

import (
)

//// func foo() {
//// 	fmt.Printf("Go runs OK!\n")
//// 	fmt.Printf("PPID: %d -> PID:%d\n", os.Getppid(), os.Getpid())
//// 	fmt.Printf("Compiler: %s v%s\n", runtime.Compiler, runtime.Version())
//// 	fmt.Printf("Architecture: %s v%s\n", runtime.GOARCH, runtime.GOOS)
//// 	fmt.Printf("GOROOT: %s\n", runtime.GOROOT())
//// }
//
//
//func isInt(i interface{}) bool {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
//			return true
//		default:
//			return false
//	}
//}
//
//func isString(i interface{}) bool {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.String:
//			return true
//		default:
//			return false
//	}
//}
//
//func isSlice(i interface{}) bool {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.Slice:
//			return true
//		default:
//			return false
//	}
//}
//
//func isArray(i interface{}) bool {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.Array:
//			return true
//		default:
//			return false
//	}
//}
//
//func isMap(i interface{}) bool {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.Map:
//			return true
//		default:
//			return false
//	}
//}
//
//// ToUpper function.
//func ToUpper(i interface{}) string {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.String:
//			return strings.ToUpper(i.(string))
//		default:
//			return ""
//	}
//}
//
//// ToLower function.
//func ToLower(i interface{}) string {
//	v := reflect.ValueOf(i)
//	switch v.Kind() {
//		case reflect.String:
//			return strings.ToLower(i.(string))
//		default:
//			return ""
//	}
//}
//
//// ToString function.
//func ToString(i interface{}) string {
//	ret := ""
//	var j []byte
//	var err error
//	j, err = json.Marshal(i)
//	if err == nil {
//		ret = string(j)
//	}
//	return ret
//}
//
//// FindInMap function.
//func ReadFile(f string) string {
//	var ret string
//
//	// var err error
//	var data []byte
//	data, _ = ioutil.ReadFile(f)
//	ret = string(data)
//
//	return ret
//}
//
//// FindInMap function.
//func FindInMap(i interface{}, n string) interface{} {
//	var ret interface{}
//	n = strings.TrimPrefix(n, "\"")
//	n = strings.TrimSuffix(n, "\"")
//
//	ret, _ = findKey(i, n)
//
//	// v := reflect.ValueOf(i)
//	// switch v.Kind() {
//	// 	case reflect.Map:
//	// 		// for i := 0; i < v.Len(); i++ {
//	// 		// 	v.
//	// 		// 	fmt.Println(v.Index(i))
//	// 		// }
//	// 		//
//	// 		// for _, m := range v.MapKeys() {
//	// 		// 	if m.
//	// 		// }
//	// }
//	return ret
//}
//
//func findKey(obj interface{}, key string) (interface{}, bool) {
//
//	//if the argument is not a map, ignore it
//	mobj, ok := obj.(map[string]interface{})
//	if !ok {
//		return nil, false
//	}
//
//	for k, v := range mobj {
//		// key match, return value
//		if k == key {
//			return v, true
//		}
//
//		// if the value is a map, search recursively
//		if m, ok := v.(map[string]interface{}); ok {
//			if res, ok := findKey(m, key); ok {
//				return res, true
//			}
//		}
//
//		// if the value is an array, search recursively
//		// from each element
//		if va, ok := v.([]interface{}); ok {
//			for _, a := range va {
//				if res, ok := findKey(a, key); ok {
//					return res,true
//				}
//			}
//		}
//	}
//
//	// element not found
//	return nil,false
//}
//
//// FindInMap function.
//// func PrintEnv(ex []string) string {
//func PrintEnv() string {
//	var ret string
//
//	for range only.Once {
//		var env Environment
//		var err error
//		env, err = getEnv()
//		if err != nil {
//			break
//		}
//
//		for k, v := range env {
//			// Bit of a hack for now...
//			// Will strip out env for Docker init
//			switch {
//				case k == "MAIL":
//				case k == "HOME":
//				case k == "LOGNAME":
//				case k == "PATH":
//				case k == "PWD":
//				case k == "SHELL":
//				case k == "SHLVL":
//				case k == "USER":
//				case k == "_":
//
//				default:
//					ret += fmt.Sprintf("%s=\"%s\"; export %s\n", k, v, k)
//			}
//		}
//	}
//
//	return ret
//}
//
//func getEnv() (Environment, error) {
//	var e Environment
//	var err error
//
//	for range only.Once {
//		e = make(Environment)
//		for _, item := range os.Environ() {
//			s := strings.SplitN(item, "=", 2)
//			e[s[0]] = s[1]
//		}
//	}
//
//	return e, err
//}
