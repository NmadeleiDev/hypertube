// Универсальный пакет для заполнения структур настроек подключаемых пакетов
// из единого конфигурационного файла
package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

var bConf []byte
var confmap map[string]interface{}

// Получаем содержимое файла конфига
func SetConfigFile(fileName string) error {
	var err error
	bConf, err = ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bConf, &confmap)
	if err != nil {
		return err
	}
	return nil
}

// Заполнение структуры конфига
func ParsePackageConfig(packStruct interface{}, blockName string) error {
	var err error

	structVal := reflect.ValueOf(packStruct).Elem()
	pack_path := structVal.Type().PkgPath()

	pack := blockName
	// Если передать пустое имя блока в конфиге
	// имя будет взято из полного наименования пакета
	if len(pack) == 0 {
		pack_slice := strings.Split(pack_path, "/")
		pack = pack_slice[len(pack_slice)-1]
	}

	if V, ok := confmap[pack]; ok {
		err = SwitchSetType(structVal, V, structVal.Type(), "", pack_path)
	} else {
		err = errors.New(fmt.Sprintf("No have block <%s> in config file", pack))
	}

	return err
}

// Рекурсивная функция заполнения полей конфига
func SwitchSetType(field reflect.Value, json_value interface{}, field_type reflect.Type, field_tag reflect.StructTag, pack_path string) error {
	var err error
	switch field_type.Kind() {
	case reflect.String:
		val := json_value.(string)
		if field.CanAddr() && field.Type().Kind() != reflect.Map {
			field.Set(reflect.ValueOf(val))
		}
		if field.Type().Kind() == reflect.Map {
			keys := field.MapKeys()
			field.SetMapIndex(keys[field.Len()-1], reflect.ValueOf(val))
		}
	case reflect.Uint:
		fl, ok := json_value.(float64)
		if !ok {
			t := field_type.String()
			err = errors.New(fmt.Sprintf("In conf file does not match type %s for package %s", t, pack_path))
			return err
		}
		val := uint(fl)
		if field.CanAddr() && field.Type().Kind() != reflect.Map {
			field.Set(reflect.ValueOf(val))
		}
		if field.Type().Kind() == reflect.Map {
			keys := field.MapKeys()
			field.SetMapIndex(keys[field.Len()-1], reflect.ValueOf(val))
		}
	case reflect.Uint64:
		fl, ok := json_value.(float64)
		if !ok {
			t := field_type.String()
			err = errors.New(fmt.Sprintf("In conf file does not match type %s for package %s", t, pack_path))
			return err
		}
		val := uint64(fl)
		if field.CanAddr() && field.Type().Kind() != reflect.Map {
			field.Set(reflect.ValueOf(val))
		}
		if field.Type().Kind() == reflect.Map {
			keys := field.MapKeys()
			field.SetMapIndex(keys[field.Len()-1], reflect.ValueOf(val))
		}
	case reflect.Int:
		fl, ok := json_value.(float64)
		if !ok {
			t := field_type.String()
			err = errors.New(fmt.Sprintf("In conf file does not match type %s for package %s", t, pack_path))
			return err
		}
		val := int(fl)
		if field.CanAddr() && field.Type().Kind() != reflect.Map {
			field.Set(reflect.ValueOf(val))
		}
		if field.Type().Kind() == reflect.Map {
			keys := field.MapKeys()
			field.SetMapIndex(keys[field.Len()-1], reflect.ValueOf(val))
		}
	case reflect.Int64:
		fl, ok := json_value.(float64)
		if !ok {
			t := field_type.String()
			err = errors.New(fmt.Sprintf("In conf file does not match type %s for package %s", t, pack_path))
			return err
		}
		var val reflect.Value
		switch field_type.PkgPath() + `.` + field_type.Name() {
		case "time.Duration":
			d_name := field_tag.Get("time")
			d := multiplyDuration(time.Duration(int64(fl)), d_name) // time.Second * time.Duration(int64(fl))
			val = reflect.ValueOf(d)
		default:
			val = reflect.ValueOf(int64(fl))
		}
		if field.CanAddr() && field.Type().Kind() != reflect.Map {
			field.Set(val)
		}
		if field.Type().Kind() == reflect.Map {
			keys := field.MapKeys()
			field.SetMapIndex(keys[field.Len()-1], reflect.ValueOf(val))
		}
	case reflect.Bool:
		b, ok := json_value.(bool)
		if !ok {
			t := field_type.String()
			err = errors.New(fmt.Sprintf("In conf file does not match type %s for package %s", t, pack_path))
			return err
		}
		val := reflect.ValueOf(b)
		if field.CanAddr() && field.Type().Kind() != reflect.Map {
			field.Set(val)
		}
		if field.Type().Kind() == reflect.Map {
			keys := field.MapKeys()
			field.SetMapIndex(keys[field.Len()-1], val)
		}
	case reflect.Slice:
		v_slice := reflect.ValueOf(json_value)
		t_slice := field_type.Elem()
		field.Set(reflect.MakeSlice(reflect.SliceOf(t_slice), v_slice.Len(), v_slice.Cap()))

		for j := 0; j < v_slice.Len(); j++ {
			field.Index(j).Set(reflect.Zero(t_slice))
			err := SwitchSetType(field.Index(j), v_slice.Index(j).Interface(), t_slice, field_tag, pack_path)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		v_map := reflect.ValueOf(json_value)
		field.Set(reflect.MakeMap(field_type))

		for _, k := range v_map.MapKeys() {
			val_json := v_map.MapIndex(k)
			n_key := reflect.ValueOf(k.Interface().(string))
			field.SetMapIndex(n_key, reflect.Zero(field_type.Elem()))
			err := SwitchSetType(field, val_json.Interface(), field_type.Elem(), field_tag, pack_path)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		var structmap map[string]interface{}
		structmap = json_value.(map[string]interface{})
		for i := 0; i < field_type.NumField(); i++ {
			tag := field_type.Field(i).Tag.Get("conf")
			// Если тэга conf нет у поля структуры либо там стоит "-",
			// то это поле пропускаем
			if tag == "" || tag == "-" {
				continue
			}
			if structmap[tag] == nil {
				return errors.New(fmt.Sprintf("Not option %s for package %s", tag, pack_path))
			}
			field_tag_child := field_type.Field(i).Tag
			field_type_child := field_type.Field(i).Type
			json_value_child := structmap[tag]
			err = SwitchSetType(field.Field(i), json_value_child, field_type_child, field_tag_child, pack_path)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// Возможность в конфиге задавать время в разных единицах.
func multiplyDuration(duration time.Duration, multiplier string) time.Duration {
	switch multiplier {
	case "Microsecond":
		return duration * time.Microsecond
	case "Millisecond":
		return duration * time.Millisecond
	case "Second":
		return duration * time.Second
	case "Minute":
		return duration * time.Minute
	case "Hour":
		return duration * time.Hour
	default:
		return duration
	}
}
