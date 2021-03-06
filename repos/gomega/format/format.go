package format

import (
    "fmt"
    "reflect"
    "strings"
)

var Indent = "    "
var longFormThreshold = 20

func Message(actual interface{}, message string, expected ...interface{}) string {
    if len(expected) == 0 {
        return fmt.Sprintf("Expected\n%s\n%s", Object(actual, 1), message)
    } else {
        return fmt.Sprintf("Expected\n%s\n%s\n%s", Object(actual, 1), message, Object(expected[0], 1))
    }
}

func Object(object interface{}, indentation uint) string {
    indent := strings.Repeat(Indent, int(indentation))
    return fmt.Sprintf("%s<%s>: %s", indent, formatType(object), formatValue(object, indentation))
}

func IndentString(s string, indentation uint) string {
    components := strings.Split(s, "\n")
    result := ""
    indent := strings.Repeat(Indent, int(indentation))
    for i, component := range components {
        result += indent + component
        if i < len(components)-1 {
            result += "\n"
        }
    }

    return result
}

func formatType(object interface{}) string {
    t := reflect.TypeOf(object)
    if t == nil {
        return "nil"
    }
    switch t.Kind() {
    case reflect.Chan:
        v := reflect.ValueOf(object)
        return fmt.Sprintf("%T | len:%d, cap:%d", object, v.Len(), v.Cap())
    case reflect.Ptr:
        return fmt.Sprintf("%T | %p", object, object)
    case reflect.Slice:
        v := reflect.ValueOf(object)
        return fmt.Sprintf("%T | len:%d, cap:%d", object, v.Len(), v.Cap())
    case reflect.Map:
        v := reflect.ValueOf(object)
        return fmt.Sprintf("%T | len:%d", object, v.Len())
    default:
        return fmt.Sprintf("%T", object)
    }
}

func formatValue(object interface{}, indentation uint) string {
    if isNil(object) {
        return "nil"
    }
    t := reflect.TypeOf(object)
    switch t.Kind() {
    case reflect.Bool:
        return fmt.Sprintf("%v", object)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return fmt.Sprintf("%v", object)
    case reflect.Uintptr:
        return fmt.Sprintf("%#v", object)
    case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
        return fmt.Sprintf("%v", object)
    case reflect.Chan:
        return fmt.Sprintf("%v", object)
    case reflect.Func:
        return fmt.Sprintf("%v", object)
    case reflect.Ptr:
        v := reflect.ValueOf(object)
        return formatValue(v.Elem().Interface(), indentation)
    case reflect.Slice:
        if t.Elem().Kind() == reflect.Uint8 {
            return formatString(object, indentation)
        }
        return formatSlice(object, indentation)
    case reflect.String:
        return formatString(object, indentation)
    case reflect.Array:
        return formatSlice(object, indentation)
    case reflect.Map:
        return formatMap(object, indentation)
    case reflect.Struct:
        return formatStruct(object, indentation)
    default:
        return fmt.Sprintf("%#v", object)
    }
}

func formatString(object interface{}, indentation uint) string {
    if indentation == 1 {
        s := fmt.Sprintf("%s", object)
        components := strings.Split(s, "\n")
        result := ""
        for i, component := range components {
            if i == 0 {
                result += component
            } else {
                result += Indent + component
            }
            if i < len(components)-1 {
                result += "\n"
            }
        }

        return fmt.Sprintf("%s", result)
    } else {
        return fmt.Sprintf("%q", object)
    }
}

func formatSlice(object interface{}, indentation uint) string {
    v := reflect.ValueOf(object)

    l := v.Len()
    result := make([]string, l)
    longest := 0
    for i := 0; i < l; i++ {
        result[i] = formatValue(v.Index(i).Interface(), indentation+1)
        if len(result[i]) > longest {
            longest = len(result[i])
        }
    }

    if longest > longFormThreshold {
        indenter := strings.Repeat(Indent, int(indentation))
        return fmt.Sprintf("[\n%s%s,\n%s]", indenter+Indent, strings.Join(result, ",\n"+indenter+Indent), indenter)
    } else {
        return fmt.Sprintf("[%s]", strings.Join(result, ", "))
    }
}

func formatMap(object interface{}, indentation uint) string {
    v := reflect.ValueOf(object)

    l := v.Len()
    result := make([]string, l)

    longest := 0
    for i, key := range v.MapKeys() {
        value := v.MapIndex(key)
        result[i] = fmt.Sprintf("%s: %s", formatValue(key.Interface(), 0), formatValue(value.Interface(), indentation+1))
        if len(result[i]) > longest {
            longest = len(result[i])
        }
    }

    if longest > longFormThreshold {
        indenter := strings.Repeat(Indent, int(indentation))
        return fmt.Sprintf("{\n%s%s,\n%s}", indenter+Indent, strings.Join(result, ",\n"+indenter+Indent), indenter)
    } else {
        return fmt.Sprintf("{%s}", strings.Join(result, ", "))
    }
}

func formatStruct(object interface{}, indentation uint) string {
    v := reflect.ValueOf(object)
    t := reflect.TypeOf(object)

    l := v.NumField()
    result := []string{}
    longest := 0
    for i := 0; i < l; i++ {
        structField := t.Field(i)
        if structField.PkgPath == "" { //implies the field is exported
            fieldEntry := v.Field(i).Interface()
            representation := fmt.Sprintf("%s: %s", structField.Name, formatValue(fieldEntry, indentation+1))
            result = append(result, representation)
            if len(representation) > longest {
                longest = len(representation)
            }
        }
    }
    if longest > longFormThreshold {
        indenter := strings.Repeat(Indent, int(indentation))
        return fmt.Sprintf("{\n%s%s,\n%s}", indenter+Indent, strings.Join(result, ",\n"+indenter+Indent), indenter)
    } else {
        return fmt.Sprintf("{%s}", strings.Join(result, ", "))
    }
}

func isNil(a interface{}) bool {
    if a == nil {
        return true
    }

    switch reflect.TypeOf(a).Kind() {
    case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
        return reflect.ValueOf(a).IsNil()
    }

    return false
}
