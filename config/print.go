package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// PrintConfig prints structure with field names and values in a clean format
func PrintConfig(cfg any) {
	fmt.Println("Configuration:")
	fmt.Println("--------------")
	printReflected(reflect.ValueOf(cfg), "", 0)
}

func printReflected(v reflect.Value, fieldName string, depth int) {
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Printf("%s%s: <nil>\n", strings.Repeat("  ", depth), fieldName)
			return
		}
		v = v.Elem()
	}

	// Handle basic types directly
	if v.Kind() != reflect.Struct {
		if fieldName != "" {
			fmt.Printf("%s%s: %v\n", strings.Repeat("  ", depth), fieldName, v.Interface())
		} else {
			fmt.Printf("%s%v\n", strings.Repeat("  ", depth), v.Interface())
		}
		return
	}

	// If this is a nested struct (with name), print only field name, not type name
	if fieldName != "" {
		fmt.Printf("%s%s:\n", strings.Repeat("  ", depth), fieldName)
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Mask sensitive fields
		nameLower := strings.ToLower(fieldType.Name)
		if strings.Contains(nameLower, "password") || strings.Contains(nameLower, "secret") || strings.Contains(nameLower, "key") {
			fmt.Printf("%s%s: ******\n", strings.Repeat("  ", depth+1), fieldType.Name)
			continue
		}

		// Handle time.Duration specially
		if field.Type().String() == "time.Duration" {
			fmt.Printf("%s%s: %v\n", strings.Repeat("  ", depth+1), fieldType.Name, field.Interface().(time.Duration))
			continue
		}

		// Recurse for nested structs or pointers to structs
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			printReflected(field, fieldType.Name, depth+1)
			continue
		}

		// Print normal value
		fmt.Printf("%s%s: %v\n", strings.Repeat("  ", depth+1), fieldType.Name, field.Interface())
	}
}
