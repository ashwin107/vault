package framework

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// FieldData is the structure passed to the callback to handle a path
// containing the populated parameters for fields. This should be used
// instead of the raw (*vault.Request).Data to access data in a type-safe
// way.
type FieldData struct {
	Raw    map[string]interface{}
	Schema map[string]*FieldSchema
}

// Get gets the value for the given field. If the key is an invalid field,
// FieldData will panic. If you want a safer version of this method, use
// GetOk. If the field k is not set, the default value (if set) will be
// returned, otherwise the zero value will be returned.
func (d *FieldData) Get(k string) interface{} {
	schema, ok := d.Schema[k]
	if !ok {
		panic(fmt.Sprintf("field %s not in the schema", k))
	}

	value, ok := d.GetOk(k)
	if !ok {
		value = schema.DefaultOrZero()
	}

	return value
}

// GetOk gets the value for the given field. The second return value
// will be false if the key is invalid or the key is not set at all.
func (d *FieldData) GetOk(k string) (interface{}, bool) {
	schema, ok := d.Schema[k]
	if !ok {
		return nil, false
	}

	result, ok, err := d.GetOkErr(k)
	if err != nil {
		panic(fmt.Sprintf("error reading %s: %s", k, err))
	}

	if ok && result == nil {
		result = schema.DefaultOrZero()
	}

	return result, ok
}

// GetOkErr is the most conservative of all the Get methods. It returns
// whether key is set or not, but also an error value. The error value is
// non-nil if the field doesn't exist or there was an error parsing the
// field value.
func (d *FieldData) GetOkErr(k string) (interface{}, bool, error) {
	schema, ok := d.Schema[k]
	if !ok {
		return nil, false, fmt.Errorf("unknown field: %s", k)
	}

	switch schema.Type {
	case TypeBool:
		fallthrough
	case TypeInt:
		fallthrough
	case TypeMap:
		fallthrough
	case TypeString:
		return d.getPrimitive(k, schema)
	default:
		return nil, false,
			fmt.Errorf("unknown field type %s for field %s", schema.Type, k)
	}
}

func (d *FieldData) getPrimitive(
	k string, schema *FieldSchema) (interface{}, bool, error) {
	raw, ok := d.Raw[k]
	if !ok {
		return nil, false, nil
	}

	switch schema.Type {
	case TypeBool:
		var result bool
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}

		return result, true, nil
	case TypeInt:
		var result int
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}

		return result, true, nil
	case TypeString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}

		return result, true, nil
	case TypeMap:
		var result map[string]interface{}
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}

		return result, true, nil
	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}
}
