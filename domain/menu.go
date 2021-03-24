package domain

import "go-take-lessons/model"

func CvMenuToVo(menus *[]model.Menu) interface{} {
	var data []interface{}
	for _, menu := range *menus {
		data = append(data, map[string]interface{}{
			"name":   menu.Name,
			"router": menu.Router,
		})
	}
	return data
}
