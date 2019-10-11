package util

import "reflect"
import "github.com/fatih/structtag"

func GetPayload(target interface{}, data interface{}) {
	tv := reflect.Indirect(reflect.ValueOf(target))
	t := tv.Type()
	for k, v := range data.(map[string]interface{}) {
		for i := 0; i < tv.NumField(); i++ {
			if tag, err := structtag.Parse(string(t.Field(i).Tag)); err == nil {
				if jt, er := tag.Get("json"); er == nil {
					if jt.Name == k {
						tv.FieldByName(t.Field(i).Name).Set(reflect.ValueOf(v))
					}
				}
			}
		}
	}
}
