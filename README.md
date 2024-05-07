# go-struct-value
library for get all primitive value of struct. support for use with pointer property.
but cannot use with nested struct

## How to install
```bash
 go get github.com/satitza/go-struct-value
```
# How to use
  ## create your struct and global variable for store custom name of property if your want to custom name
  ```bash
 var customFieldName = map[string]string{
   "test_string_ptr": "custom_test_string_ptr"
 }

 var addCustomFields = map[string]any{
   "add_custom_field": "custom_value"
 }

 var convertDateToEpoch = []string{"test_epoch_time"} // support date format 2006-01-02 15:04:05

 type CustomStruct struct {
	TestStringPtr *string                 `json:"test_string_ptr"`
	TestIntPtr    *int                    `json:"test_int_ptr"`
	TestBoolPtr   *bool                   `json:"test_bool_ptr"`
	TestString    string                  `json:"test_string"`
	TestInt       int                     `json:"test_int"`
	TestBool      bool                    `json:"test_bool"`
	TestEpochTime string                  `json:"test_epoch_time"`
	TestJson      go_struct_value.MapJson `json:"test_json"`
 } 

  ```

## create receiver function for call function in library and pass parameter
  ```bash
  func (req CustomStruct) GetAllColumnsNameString() ([]string, error) {
	 return go_struct_value.GetAllColumnsName(req, customFieldName, addCustomFields)
  }

  func (req CustomStruct) GetFieldInsertData() ([]string, []any, error) {
	 return go_struct_value.GetSqlAndDataForInsert(req, customFieldName, addCustomFields, convertDateToEpoch)
  }

  func (req CustomStruct) GetFieldUpdateMap() (map[string]any, error) {
	 return go_struct_value.GetFieldValueMap(req, customFieldName, convertDateToEpoch)
  }

  func (req CustomStruct) GetAllFieldValueMap() (map[string]any, error) {
	 return go_struct_value.GetAllFieldValueMap(req, customFieldName, convertDateToEpoch)
  }
  ```

### - GetAllColumnsNameString return all property name array and convert to snake case
### - GetFieldInsertData      return all property name array(string snake case) and array value of all property
### - GetFieldUpdateMap       return map[string]any the key is property name (string snake case) value is a value of property. but get kay and value if property has value only ( not nil or empty string)
### - GetAllFieldValueMap     return map[string]any the key is property name (string snake case) value is a value of property include nil or empty string

## Example
  ```bash
var testStringPtr = "testStringPtr"
	var testIntPtr = 1234
	var testBoolPtr = true

	var customStruct = model.CustomStruct{
		TestStringPtr: &testStringPtr,
		TestIntPtr:    &testIntPtr,
		TestBoolPtr:   &testBoolPtr,
		TestString:    "TestString",
		TestInt:       4321,
		TestBool:      false,
		TestEpochTime: "2024-02-04 13:53:05",
		TestJson: map[string]interface{}{
			"col1": "col1",
		},
	}

	var mapFieldValues, err = customStruct.GetAllFieldValueMap()
	if err != nil {
		fmt.Println(err)
	}

	for key, value := range mapFieldValues {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
  ```
## Output
```bash
  Key: test_json, Value: [123 34 99 111 108 49 34 58 34 99 111 108 49 34 125]
  Key: test_string_ptr, Value: testStringPtr
  Key: test_int_ptr, Value: 1234
  Key: test_bool_ptr, Value: true
  Key: test_string, Value: TestString
  Key: test_int, Value: 4321
  Key: test_bool, Value: false
  Key: test_epoch_time, Value: 1707029585
  ```





