package sui

import (
	"fmt"
	"reflect"

	"github.com/howjmay/sui-go/sui_types"
)

type SuiJsonArg any

// signed integers (e.g. int, int8, int16, etc) and arbitrary size 'uint' are not supported
// more about Sui RPC JSON format (aka SuiJSON) see https://docs.sui.io/references/sui-api
// TODO u128/u256 are not supported yet
func ToSuiJsonArg(input any) SuiJsonArg {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	// If the value is a pointer, dereference it
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		panic("int/int8/int16/int32/int64/uint are not supported")
	case reflect.Uint:
		panic("uint is not supported")
	case reflect.Float32, reflect.Float64:
		panic("float32/float64 are not supported")

	// Force integers to have zero padding. In this way move vm won't confuse about the arg types.
	case reflect.Uint8:
		return fmt.Sprintf("0x%02x", input)
	case reflect.Uint16:
		return fmt.Sprintf("0x%04x", input)
	case reflect.Uint32:
		return fmt.Sprintf("0x%08x", input)
	case reflect.Uint64:
		return fmt.Sprintf("0x%016x", input)
	case reflect.Bool:
		return input
	case reflect.String:
		return input
	case reflect.Array:
		typ := reflect.TypeOf(input)
		if typ == reflect.TypeOf(&sui_types.SuiAddress{}) ||
			typ == reflect.TypeOf(&sui_types.ObjectID{}) ||
			typ == reflect.TypeOf(&sui_types.PackageID{}) {
			return input.(*sui_types.SuiAddress).String()
		}
		panic("not supported type")
	case reflect.Slice:
		elemType := typ.Elem()
		if elemType == reflect.TypeOf(&sui_types.SuiAddress{}) ||
			elemType == reflect.TypeOf(&sui_types.ObjectID{}) ||
			elemType == reflect.TypeOf(&sui_types.PackageID{}) {
			// Handle slices of pointers to SuiAddress, ObjectID, and PackageID
			ret := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(&sui_types.SuiAddress{})), val.Len(), val.Len())
			for i := 0; i < val.Len(); i++ {
				innerVal := val.Index(i)
				ret.Index(i).Set(innerVal)
			}
			return ret.Interface()
		} else if elemType.Kind() == reflect.Uint8 {
			return input
		} else if elemType.Kind() == reflect.String {
			return input
		} else {
			inputType := val.Type().Elem()
			newSliceType := reflect.SliceOf(inputType)
			ret := reflect.MakeSlice(newSliceType, 0, val.Len())
			for i := 0; i < val.Len(); i++ {
				innerVal := val.Index(i)

				// Recursively convert inner elements
				convertedInnerVal := reflect.ValueOf(ToSuiJsonArg(innerVal.Interface()))

				// Ensure that the converted value is assignable to the element type of the new slice
				if !convertedInnerVal.Type().AssignableTo(inputType) {
					panic("converted value is not assignable to the slice element type")
				}

				ret = reflect.Append(ret, convertedInnerVal)
			}
			return ret.Interface()
		}
	default:
		panic("not supported type")
	}
}
