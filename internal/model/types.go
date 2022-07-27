// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"strings"
)

// Model 模型定义
type Model struct {
	TableName string
	GoName    string
	Fields    []Field
}

type Field struct {
	ColName      string
	IsPrimaryKey bool
	GoName       string
	Order        bool
	GoType       string
}

func (f *Field) IsInteger() bool {
	switch f.GoType {
	case "int64", "int32", "int16", "int8", "int":
		return true
	case "uint64", "uint32", "uint16", "uint8", "uint":
		return true
	case "byte", "rune":
		return true
	}
	return false
}

func (f *Field) IsFloat() bool {
	switch f.GoType {
	case "float32", "float64":
		return true
	}
	return false
}

func (f *Field) IsString() bool {
	switch f.GoType {
	case "string":
		return true
	}
	return false
}

func (f *Field) IsBool() bool {
	if f.GoType == "bool" {
		return true
	}
	return false
}

func (f *Field) IsSlice() bool {
	if f.GoType == "slice" {
		return true
	}
	return false
}

func (f *Field) IsMap() bool {
	if f.GoType == "map" {
		return true
	}
	return false
}

func (f *Field) IsArray() bool {
	if f.GoType == "array" {
		return true
	}
	return false
}

func (f *Field) IsPtr() bool {
	if f.GoType == "ptr" {
		return true
	}
	return false
}

func (m *Model) QuotedTableName() string {
	return "`" + m.TableName + "`"
}

func (m *Model) QuotedExecArgsWithParameter(flag, owner, col string) string {
	var str []string
	for _, v := range m.Fields {
		if strings.Contains(col, v.ColName) {
			str = append(str, flag+owner+"."+v.GoName)
		}
	}
	return strings.Join(str, ", ")
}

func (m *Model) InsertWithReplaceParameter() string {
	var str strings.Builder
	for k := range m.Fields {
		if k != 0 {
			str.WriteByte(',')
		}
		str.WriteByte('?')
	}
	return str.String()
}

func (m *Model) QuotedExecArgsWithAll() string {
	var str strings.Builder
	for k, v := range m.Fields {
		if k != 0 {
			str.WriteString(", ")
		}
		str.WriteString("v." + v.GoName)
	}
	return str.String()
}

func (m *Model) QuotedAllCol() string {
	var str strings.Builder
	for k, v := range m.Fields {
		if k != 0 {
			str.WriteByte(',')
		}
		str.WriteString("`" + v.ColName + "`")
	}
	
	return str.String()
}
