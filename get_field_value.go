package go_struct_value

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
)

type MapJson map[string]interface{}

func (j *MapJson) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *MapJson) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &j)
}

// GetAllColumnsName For Select all column
func GetAllColumnsName(req interface{}, customFieldName map[string]string, addCustomFields map[string]any) ([]string, error) {

	var columns []string

	for index := 0; index < reflect.TypeOf(req).NumField(); index++ {
		field := reflect.TypeOf(req).Field(index)
		dataType := field.Type.Kind()
		if dataType != reflect.Struct &&
			dataType != reflect.Slice &&
			dataType != reflect.Interface &&
			dataType != reflect.Array &&
			dataType != reflect.Chan &&
			dataType != reflect.Func {

			if dataType == reflect.Pointer {

				if field.Type.Elem().Kind() == reflect.Map {
					if field.Type.Elem() != reflect.TypeOf(map[string]interface{}{}) &&
						field.Type.Elem() != reflect.TypeOf(MapJson{}) {
						continue
					}
				}

				if field.Type.Elem().Kind() == reflect.Struct ||
					field.Type.Elem().Kind() == reflect.Slice ||
					field.Type.Elem().Kind() == reflect.Interface ||
					field.Type.Elem().Kind() == reflect.Array ||
					field.Type.Elem().Kind() == reflect.Chan ||
					field.Type.Elem().Kind() == reflect.Func {
					continue
				}
			}

			fieldName := strcase.ToSnake(field.Name)
			if val, ok := customFieldName[fieldName]; ok {
				fieldName = val
			}

			columns = append(columns, fieldName)
		}
	}

	if len(addCustomFields) > 0 {
		for key := range addCustomFields {
			columns = append(columns, key)
		}
	}

	return columns, nil
}

// GetSqlAndDataForInsert For Insert
func GetSqlAndDataForInsert(req interface{}, customFieldName map[string]string, addCustomFields map[string]any, convertDateToEpoch []string) ([]string, []any, error) {

	var columns []string
	var data []any

	values := reflect.ValueOf(req)
	for index := 0; index < reflect.TypeOf(req).NumField(); index++ {

		field := reflect.TypeOf(req).Field(index)
		dataType := field.Type.Kind()
		if dataType != reflect.Struct &&
			dataType != reflect.Slice &&
			dataType != reflect.Interface &&
			dataType != reflect.Array &&
			dataType != reflect.Chan &&
			dataType != reflect.Func {

			if dataType == reflect.Pointer {

				if field.Type.Elem().Kind() == reflect.Map {
					if field.Type.Elem() != reflect.TypeOf(map[string]interface{}{}) &&
						field.Type.Elem() != reflect.TypeOf(MapJson{}) {
						continue
					}
				}

				if field.Type.Elem().Kind() == reflect.Struct ||
					field.Type.Elem().Kind() == reflect.Slice ||
					field.Type.Elem().Kind() == reflect.Interface ||
					field.Type.Elem().Kind() == reflect.Array ||
					field.Type.Elem().Kind() == reflect.Chan ||
					field.Type.Elem().Kind() == reflect.Func {
					continue
				}
			}

			fieldName := strcase.ToSnake(field.Name)
			if val, ok := customFieldName[fieldName]; ok {
				fieldName = val
			}

			columns = append(columns, fieldName)
			value := values.Field(index)

			if dataType == reflect.Pointer {
				if field.Type.Elem() == reflect.TypeOf(map[string]interface{}{}) ||
					field.Type.Elem() == reflect.TypeOf(MapJson{}) {
					realValue := value.Interface()
					if !reflect.ValueOf(realValue).IsNil() {
						jsonByte, err := json.Marshal(realValue)
						if err != nil {
							return nil, nil, err
						}
						data = append(data, jsonByte)
					} else {
						data = append(data, nil)
					}

				} else {
					if !reflect.ValueOf(value.Interface()).IsNil() {
						realValue := value.Elem()
						data = append(data, realValue.Interface())
					} else {
						data = append(data, nil)
					}
				}
			} else {
				var realValue any
				if dataType == reflect.Bool {
					realValue = value.Bool()
					data = append(data, realValue)
				} else {
					switch dataType {
					case reflect.String:
						if !reflect.ValueOf(value.String()).IsZero() {
							realValue = value.String()
							data = append(data, realValue)
						} else {
							data = append(data, reflect.Zero(value.Type()).Interface())
						}
					case reflect.Int,
						reflect.Int8,
						reflect.Int16,
						reflect.Int32,
						reflect.Int64:
						if !reflect.ValueOf(value.Int()).IsZero() {
							realValue = value.Int()
							data = append(data, realValue)
						} else {
							data = append(data, reflect.Zero(value.Type()).Int())
						}
					case reflect.Uint,
						reflect.Uint8,
						reflect.Uint16,
						reflect.Uint32,
						reflect.Uint64:
						if !reflect.ValueOf(value.Uint()).IsZero() {
							realValue = value.Uint()
							data = append(data, realValue)
						} else {
							data = append(data, reflect.Zero(value.Type()).Uint())
						}
					case reflect.Float32, reflect.Float64:
						if !reflect.ValueOf(value.Float()).IsZero() {
							realValue = value.Float()
							data = append(data, realValue)
						} else {
							data = append(data, reflect.Zero(value.Type()).Float())
						}
					case reflect.Complex64, reflect.Complex128:
						if !reflect.ValueOf(value.Complex()).IsZero() {
							realValue = value.Complex()
							data = append(data, realValue)
						} else {
							data = append(data, reflect.Zero(value.Type()).Complex())
						}
					case reflect.Map:
						if field.Type == reflect.TypeOf(map[string]interface{}{}) ||
							field.Type == reflect.TypeOf(MapJson{}) {
							realValue = value.Interface()
							if !reflect.ValueOf(realValue).IsZero() {
								jsonByte, err := json.Marshal(realValue)
								if err != nil {
									return nil, nil, err
								}
								data = append(data, jsonByte)
							} else {
								data = append(data, nil)
							}
						}
					default:
						data = append(data, nil)
					}
				}
			}
		}
	}

	if len(addCustomFields) > 0 {
		for key := range addCustomFields {
			columns = append(columns, key)
			if val, ok := addCustomFields[key]; ok {
				data = append(data, val)
			} else {
				data = append(data, nil)
			}
		}
	}

	for _, convertDateFieldName := range convertDateToEpoch {
		for index := 0; index < len(columns); index++ {
			if columns[index] == convertDateFieldName {
				dateString := data[index]
				if dateString != nil && dateString != "" {
					location, err := time.LoadLocation("Asia/Bangkok")
					if err != nil {
						return nil, nil, err
					}
					epoch, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s", dateString), location)
					if err != nil {
						return nil, nil, err
					}
					data[index] = epoch.Unix()
				}
			}
		}
	}

	return columns, data, nil
}

// GetFieldValueMap For Update
func GetFieldValueMap(req interface{}, customFieldName map[string]string, convertDateToEpoch []string) (map[string]any, error) {

	var mapFields = make(map[string]any)

	values := reflect.ValueOf(req)
	for index := 0; index < reflect.TypeOf(req).NumField(); index++ {

		field := reflect.TypeOf(req).Field(index)
		dataType := field.Type.Kind()

		if dataType != reflect.Struct &&
			dataType != reflect.Slice &&
			dataType != reflect.Interface &&
			dataType != reflect.Array &&
			dataType != reflect.Chan &&
			dataType != reflect.Func {

			if dataType == reflect.Pointer {

				if field.Type.Elem().Kind() == reflect.Map {
					if field.Type.Elem() != reflect.TypeOf(map[string]interface{}{}) &&
						field.Type.Elem() != reflect.TypeOf(MapJson{}) {
						continue
					}
				}

				if field.Type.Elem().Kind() == reflect.Struct ||
					field.Type.Elem().Kind() == reflect.Slice ||
					field.Type.Elem().Kind() == reflect.Interface ||
					field.Type.Elem().Kind() == reflect.Array ||
					field.Type.Elem().Kind() == reflect.Chan ||
					field.Type.Elem().Kind() == reflect.Func {
					continue
				}
			}

			fieldName := strcase.ToSnake(field.Name)
			if val, ok := customFieldName[fieldName]; ok {
				fieldName = val
			}

			value := values.Field(index)
			if dataType == reflect.Pointer {
				if field.Type.Elem() == reflect.TypeOf(map[string]interface{}{}) ||
					field.Type.Elem() == reflect.TypeOf(MapJson{}) {
					realValue := value.Interface()
					if !reflect.ValueOf(realValue).IsNil() {
						jsonByte, err := json.Marshal(realValue)
						if err != nil {
							return nil, err
						}
						mapFields[fieldName] = jsonByte
					}
				} else {
					if !reflect.ValueOf(value.Interface()).IsNil() {
						realValue := value.Elem()
						mapFields[fieldName] = realValue.Interface()
					}
				}
			} else {
				var realValue any
				if dataType == reflect.Bool {
					realValue = value.Bool()
					mapFields[fieldName] = realValue
				} else {
					switch dataType {
					case reflect.String:
						if !reflect.ValueOf(value.String()).IsZero() {
							realValue = value.String()
							mapFields[fieldName] = realValue
						}
					case reflect.Int,
						reflect.Int8,
						reflect.Int16,
						reflect.Int32,
						reflect.Int64:
						if !reflect.ValueOf(value.Int()).IsZero() {
							realValue = value.Int()
							mapFields[fieldName] = realValue
						}
					case reflect.Uint,
						reflect.Uint8,
						reflect.Uint16,
						reflect.Uint32,
						reflect.Uint64:
						if !reflect.ValueOf(value.Uint()).IsZero() {
							realValue = value.Uint()
							mapFields[fieldName] = realValue
						}
					case reflect.Float32, reflect.Float64:
						if !reflect.ValueOf(value.Int()).IsZero() {
							realValue = value.Float()
							mapFields[fieldName] = realValue
						}
					case reflect.Complex64,
						reflect.Complex128:
						if !reflect.ValueOf(value.Complex()).IsZero() {
							realValue = value.Complex()
							mapFields[fieldName] = realValue
						}
					case reflect.Map:
						if field.Type == reflect.TypeOf(map[string]interface{}{}) ||
							field.Type == reflect.TypeOf(MapJson{}) {
							realValue = value.Interface()
							if !reflect.ValueOf(realValue).IsZero() {
								jsonByte, err := json.Marshal(realValue)
								if err != nil {
									return nil, err
								}
								mapFields[fieldName] = jsonByte
							}
						}
					default:
					}
				}
			}
		}
	}

	for _, name := range convertDateToEpoch {
		if val, ok := mapFields[name]; ok {
			location, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return nil, err
			}
			epoch, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s", val), location)
			if err != nil {
				return nil, err
			}
			mapFields[name] = epoch.Unix()
		}
	}

	return mapFields, nil

}

// GetAllFieldValueMap For Update With Null Value
func GetAllFieldValueMap(req interface{}, customFieldName map[string]string, convertDateToEpoch []string) (map[string]any, error) {

	var mapFields = make(map[string]any)

	values := reflect.ValueOf(req)
	for index := 0; index < reflect.TypeOf(req).NumField(); index++ {

		field := reflect.TypeOf(req).Field(index)
		dataType := field.Type.Kind()
		if dataType != reflect.Struct &&
			dataType != reflect.Slice &&
			dataType != reflect.Interface &&
			dataType != reflect.Array &&
			dataType != reflect.Chan &&
			dataType != reflect.Func {

			if dataType == reflect.Pointer {
				if field.Type.Elem().Kind() == reflect.Map {
					if field.Type.Elem() != reflect.TypeOf(map[string]interface{}{}) &&
						field.Type.Elem() != reflect.TypeOf(MapJson{}) {
						continue
					}
				}

				if dataType == reflect.Pointer {
					if field.Type.Elem().Kind() == reflect.Struct ||
						field.Type.Elem().Kind() == reflect.Slice ||
						field.Type.Elem().Kind() == reflect.Interface ||
						field.Type.Elem().Kind() == reflect.Array ||
						field.Type.Elem().Kind() == reflect.Chan ||
						field.Type.Elem().Kind() == reflect.Func {
						continue
					}
				}
			}

			fieldName := strcase.ToSnake(field.Name)
			if val, ok := customFieldName[fieldName]; ok {
				fieldName = val
			}

			mapFields[fieldName] = nil
			value := values.Field(index)

			if dataType == reflect.Pointer {
				if field.Type.Elem() == reflect.TypeOf(map[string]interface{}{}) ||
					field.Type.Elem() == reflect.TypeOf(MapJson{}) {
					realValue := value.Interface()
					if !reflect.ValueOf(realValue).IsNil() {
						jsonByte, err := json.Marshal(realValue)
						if err != nil {
							return nil, err
						}
						mapFields[fieldName] = jsonByte
					} else {
						mapFields[fieldName] = nil
					}
				} else {
					if !reflect.ValueOf(value.Interface()).IsNil() {
						realValue := value.Elem()
						mapFields[fieldName] = realValue.Interface()
					} else {
						mapFields[fieldName] = nil
					}
				}
			} else {
				var realValue any
				if dataType == reflect.Bool {
					realValue = value.Bool()
					mapFields[fieldName] = realValue
				} else {
					switch dataType {
					case reflect.String:
						if !reflect.ValueOf(value.String()).IsZero() {
							realValue = value.String()
							mapFields[fieldName] = realValue
						} else {
							mapFields[fieldName] = reflect.Zero(value.Type()).Interface()
						}
					case reflect.Int,
						reflect.Int8,
						reflect.Int16,
						reflect.Int32,
						reflect.Int64:
						if !reflect.ValueOf(value.Int()).IsZero() {
							realValue = value.Int()
							mapFields[fieldName] = realValue
						} else {
							mapFields[fieldName] = reflect.Zero(value.Type()).Int()
						}
					case reflect.Uint,
						reflect.Uint8,
						reflect.Uint16,
						reflect.Uint32,
						reflect.Uint64:
						if !reflect.ValueOf(value.Uint()).IsZero() {
							realValue = value.Uint()
							mapFields[fieldName] = realValue
						} else {
							mapFields[fieldName] = reflect.Zero(value.Type()).Uint()
						}
					case reflect.Float32, reflect.Float64:
						if !reflect.ValueOf(value.Float()).IsZero() {
							realValue = value.Float()
							mapFields[fieldName] = realValue
						} else {
							mapFields[fieldName] = reflect.Zero(value.Type()).Float()
						}
					case reflect.Complex64,
						reflect.Complex128:
						if !reflect.ValueOf(value.Complex()).IsZero() {
							realValue = value.Complex()
							mapFields[fieldName] = realValue
						} else {
							mapFields[fieldName] = reflect.Zero(value.Type()).Complex()
						}
					case reflect.Map:
						if field.Type == reflect.TypeOf(map[string]interface{}{}) ||
							field.Type == reflect.TypeOf(MapJson{}) {
							realValue = value.Interface()
							if !reflect.ValueOf(realValue).IsZero() {
								jsonByte, err := json.Marshal(realValue)
								if err != nil {
									return nil, err
								}
								mapFields[fieldName] = jsonByte
							} else {
								mapFields[fieldName] = reflect.Zero(value.Type()).Interface()
							}
						}
					default:
						mapFields[fieldName] = nil
					}
				}
			}
		}
	}

	for _, name := range convertDateToEpoch {
		if val, ok := mapFields[name]; ok {
			if val != nil && val != "" {
				location, err := time.LoadLocation("Asia/Bangkok")
				if err != nil {
					return nil, err
				}
				epoch, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s", val), location)
				if err != nil {
					return nil, err
				}

				mapFields[name] = epoch.Unix()
			}
		}
	}

	return mapFields, nil
}


// Use for audit logs
func ConvertDateTimeStringToEpochTimeString(req interface{}, convertDateToEpoch []string) (interface{}, error) {

	values := reflect.ValueOf(req)
	for index := 0; index < reflect.TypeOf(req).NumField(); index++ {

		field := reflect.TypeOf(req).Field(index)
		dataType := field.Type.Kind()
		fieldName := strcase.ToSnake(field.Name)

		for _, fieldConfigName := range convertDateToEpoch {
			if fieldName == fieldConfigName {
				value := values.Field(index)

				if !value.IsNil() && !value.IsZero() {
					var oldDateString any
					if dataType == reflect.Pointer {
						oldDateString = value.Elem().Interface()
					} else {
						oldDateString = value.Interface()
					}

					location, err := dateutil.LoadThaiLocation()
					if err != nil {
						return nil, err
					}

					epoch, err := dateutil.DateTimeString2EpochInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s", oldDateString), location)
					if err != nil {
						return nil, err
					}

					newValue := strconv.FormatInt(epoch, 10)

					if dataType == reflect.Pointer {
						value.Elem().Set(reflect.ValueOf(newValue))
					} else {
						value.Set(reflect.ValueOf(newValue))
					}
				}
			}
		}
	}

	return req, nil
}

