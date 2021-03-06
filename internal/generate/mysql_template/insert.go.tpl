package code

import (
	"context"
	"database/sql"
)

type {{.GoName}}DAO struct {
	DB *sql.DB
}

func (dao *{{.GoName}}DAO) Insert(ctx context.Context, vals ...*{{.GoName}}) (int64, error) {
	var args = make([]interface{}, len(vals)*({{len .Fields}}))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "({{.InsertWithReplaceParameter}})"
		args = append(args, {{.QuotedExecArgsWithAll}})
	}
	sqlSen := "INSERT INTO {{.QuotedTableName}}({{.QuotedAllCol}}) VALUES" + str
	res, err := dao.DB.ExecContext(ctx, sqlSen, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}