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
	 return go_struct_value.GetSqlAndDataForInsert(req, nil, nil, convertDateToEpoch)
  }

  func (req CustomStruct) GetFieldUpdateMap() (map[string]any, error) {
	 return go_struct_value.GetFieldValueMap(req, nil, convertDateToEpoch)
  }

  func (req CustomStruct) GetAllFieldValueMap() (map[string]any, error) {
	 return go_struct_value.GetAllFieldValueMap(req, nil, convertDateToEpoch)
  }
  ```






