package database

import (
	"reflect"
	"strings"
	"unicode"
)

// ParseRDBUpdateData is the same as ParseMDBUpdateData, but it will use the struct field name with the first letter lower case,
// avoid use field names like nationalIDNo, it will be parsed as incorrect national_idno.
// if the input struct has a gorm embedded field, gorm tag embedded should be set, e.g. `gorm:"embedded"`
// this way the function knows that the  field is an embedded struct and unpack it to the result.
func ParseRDBUpdateData(x interface{}, skipFields ...string) map[string]interface{} {
	result := make(map[string]interface{})
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	// logger.Debugf("number of field:s %v", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		vField := v.Field(i)

		name := tField.Name

		// Skip the field if it's in the skipFields list
		if isSkip(skipFields, name) {
			// logger.Debugf("skip field: %s", name)
			continue
		}

		// logger.Debugf("process field: %s", name)

		// if the field is embedded, unpack the field to result map
		if strings.Contains(tField.Tag.Get("gorm"), "embedded") {
			// logger.Debugf("embedded field: %s", name)

			if vField.Kind() != reflect.Ptr {
				// unpack the embedded field to result map[string]interface{}
				embeddedData := ParseRDBUpdateData(vField.Interface(), skipFields...)

				for k, v := range embeddedData {
					result[k] = v
				}
				continue
			}

			// Check if the field is a pointer and if it's a nil
			if vField.Kind() == reflect.Ptr && vField.IsNil() {
				// logger.Debugf("embedded field is nil: %s", name)
				continue
			}

			// If the field is a pointer, get the actual value and unpack it
			if vField.Kind() == reflect.Ptr {
				vField = vField.Elem()

				// unpack the embedded field to result map[string]interface{}
				embeddedData := ParseRDBUpdateData(vField.Interface(), skipFields...)

				for k, v := range embeddedData {
					result[k] = v
				}

				// logger.Debugf("embedded field unpacked: %s", name)
			}
			continue
		}

		name = toSnakeCase(name)

		// Check if the field is a pointer and if it's a nil
		if vField.Kind() == reflect.Ptr && vField.IsNil() {
			continue
		}

		// If the field is a pointer, get the actual value
		if vField.Kind() == reflect.Ptr {
			vField = vField.Elem()
		}

		// Add the field to the map
		result[name] = vField.Interface()
	}

	return result
}

func isSkip(skip []string, name string) bool {
	for _, v := range skip {
		if v == name {
			return true
		}
	}
	return false
}

// if field name is "Name", it will return "name"
// if field name is "FirstName", it will return "first_name"
func toSnakeCase(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(rune(s[i-1])) {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
