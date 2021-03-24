package excel

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// 创建模板
func CreateExcelModel(head []string, title, sheetName string) *excelize.File {
	f := excelize.NewFile()

	index := f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", title)

	mapping := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	headStyle, _ := f.NewStyle(GetHeadCellStyle())
	for i, h := range head {
		f.SetCellValue(sheetName, mapping[i]+"2", h)
		f.SetCellStyle(sheetName, mapping[i]+"2", mapping[i]+"2", headStyle)
	}

	titleStyle, _ := f.NewStyle(GetTitleCellStyle())
	f.MergeCell(sheetName, mapping[0]+"1", mapping[len(head)-1]+"1")
	f.SetCellStyle(sheetName, mapping[0]+"1", mapping[len(head)-1]+"1", titleStyle)

	f.SetActiveSheet(index)
	return f
}

// 创建模板
func CreateFileExcelModel(f *excelize.File, head []string, title, sheetName string) {
	index := f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", title)

	mapping := [10]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	headStyle, _ := f.NewStyle(GetHeadCellStyle())
	for i, h := range head {
		f.SetCellValue(sheetName, mapping[i]+"2", h)
		f.SetCellStyle(sheetName, mapping[i]+"2", mapping[i]+"2", headStyle)
	}

	titleStyle, _ := f.NewStyle(GetTitleCellStyle())
	f.MergeCell(sheetName, mapping[0]+"1", mapping[len(head)-1]+"1")
	f.SetCellStyle(sheetName, mapping[0]+"1", mapping[len(head)-1]+"1", titleStyle)

	f.SetActiveSheet(index)
}

func GetHeadCellStyle() string {
	color := "000000"
	center := "center"
	fontFamily := "黑体"
	fontSize := 11

	var lBorder Border
	lBorder.Color = color
	lBorder.Style = 1
	lBorder.Type = "left"

	var rBorder Border
	rBorder.Color = color
	rBorder.Style = 1
	rBorder.Type = "right"

	var tBorder Border
	tBorder.Color = color
	tBorder.Style = 1
	tBorder.Type = "top"

	var bBorder Border
	bBorder.Color = color
	bBorder.Style = 1
	bBorder.Type = "bottom"

	var font Font
	font.Color = color
	font.Bold = false
	font.Family = fontFamily
	font.Size = float64(fontSize)

	var alignment Alignment
	alignment.Vertical = center
	alignment.Horizontal = center

	var style Style
	style.Border = []Border{lBorder, rBorder, tBorder, bBorder}
	style.Font = &font
	style.Alignment = &alignment

	b, err := json.Marshal(style)

	if err != nil {
		println(err.Error())
	}
	return string(b)
}

func GetTitleCellStyle() string {
	color := "000000"
	center := "center"
	fontFamily := "黑体"
	fontSize := 14

	var lBorder Border
	lBorder.Color = color
	lBorder.Style = 1
	lBorder.Type = "left"

	var rBorder Border
	rBorder.Color = color
	rBorder.Style = 1
	rBorder.Type = "right"

	var tBorder Border
	tBorder.Color = color
	tBorder.Style = 1
	tBorder.Type = "top"

	var bBorder Border
	bBorder.Color = color
	bBorder.Style = 1
	bBorder.Type = "bottom"

	var font Font
	font.Color = color
	font.Bold = false
	font.Family = fontFamily
	font.Size = float64(fontSize)

	var alignment Alignment
	alignment.Vertical = center
	alignment.Horizontal = center

	var style Style
	style.Border = []Border{lBorder, rBorder, tBorder, bBorder}
	style.Font = &font
	style.Alignment = &alignment

	b, err := json.Marshal(style)

	if err != nil {
		println(err.Error())
	}
	return string(b)
}

func GetContentCellStyle() string {
	color := "000000"
	center := "center"
	fontFamily := "黑体"
	fontSize := 10

	var lBorder Border
	lBorder.Color = color
	lBorder.Style = 1
	lBorder.Type = "left"

	var rBorder Border
	rBorder.Color = color
	rBorder.Style = 1
	rBorder.Type = "right"

	var tBorder Border
	tBorder.Color = color
	tBorder.Style = 1
	tBorder.Type = "top"

	var bBorder Border
	bBorder.Color = color
	bBorder.Style = 1
	bBorder.Type = "bottom"

	var font Font
	font.Color = color
	font.Bold = false
	font.Family = fontFamily
	font.Size = float64(fontSize)

	var alignment Alignment
	alignment.Vertical = center
	alignment.Horizontal = center

	var style Style
	style.Border = []Border{lBorder, rBorder, tBorder, bBorder}
	style.Font = &font
	style.Alignment = &alignment

	b, err := json.Marshal(style)

	if err != nil {
		println(err.Error())
	}
	return string(b)
}
