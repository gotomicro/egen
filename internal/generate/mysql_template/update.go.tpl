

func (dao *{{.GoName}}DAO) UpdateByWhereWithNoneZero(ctx context.Context, val *{{.GoName}}, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateByNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE {{.QuotedTableName}} SET " + cols + " WHERE " + where
	return dao.UpdateByRow(ctx, val, s, newArgs...)
}

func (dao *{{.GoName}}DAO) UpdateByNoneZero(val *{{.GoName}}) ([]interface{}, string) {
	var cols = make([]string, 0, {{len .Fields}})
	var args = make([]interface, 0, {{len .Fields}})
	judge := func (x any) bool {
    		return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
    }
	{{- range .Fields}}
	if judge(val.{{.GoName}}) {
		args = append(args, val.{{.GoName}})
		cols = append(cols, "`{{.ColName}}`")
	}
	{{end}}
	return args, strings.Join(cols, "=?,")
}

func (dao *{{.GoName}}DAO) UpdateByWhereWithPrimaryKey(ctx context.Context, val *{{.GoName}}, where string, args ...any) (int64, error) {
	newArgs, cols := UpdateByPrimaryKey()
	newArgs = append(newArgs, args...)
	s := "UPDATE {{.QuotedTableName}} SET " + cols + " WHERE " + where
	return dao.UpdateByRow(ctx, val, s, newArgs...)
}

func (dao *{{.GoName}}DAO) UpdateByPrimaryKey() ([]interface{}, string) {
	var cols = make([]string, 0, {{len .Fields}})
    var args = make([]interface, 0, {{len .Fields}})
    {{- range .Fields}}
    {{- if .IsPrimaryKey}}
    args = append(args, {{print "val." .GoName}})
    cols = append(cols, "`{{.ColName}}`")
    {{end}}
    {{- end}}
    return args, strings.Join(cols, "=?,")
}

func (dao *{{.GoName}}DAO) UpdateByWhereWithSpecificCol(ctx context.Context, val *{{.GoName}}, where string, args ...any) (int64, error) {
	newArgs, cols := UpdateBySpecificCol()
	newArgs = append(newArgs, args...)
	s := "UPDATE {{.QuotedTableName}} SET " + cols + " WHERE " + where
	return dao.UpdateByRow(ctx, val, s, newArgs...)
}

func (dao *{{.GoName}}DAO) UpdateBySpecificCol() ([]interface{}, string) {
	var cols = make([]string, 0, {{len .Fields}})
    var args = make([]interface, 0, {{len .Fields}})
    {{- range .Fields}}
    {{- if .Order}}
    args = append(args, {{print "val." .GoName}})
    cols = append(cols, "`{{.ColName}}`")
    {{end}}
    {{- end}}
    return args, strings.Join(cols, "=?,")
}

func (dao *{{.GoName}}DAO) UpdateByRaw(ctx context.Context, val *{{.GoName}}, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}