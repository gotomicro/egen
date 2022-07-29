package code

import (
	"context"
	"database/sql"
	"strings"
)

type UserDAO struct {
	DB *sql.DB
}

func (dao *UserDAO) Insert(ctx context.Context, vals ...*User) (int64, error) {
	var args = make([]interface{}, len(vals)*(5))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?,?,?)"
		args = append(args, v.LoginTime, v.FirstName, v.LastName, v.UserId, v.Password)
	}
	sqlSen := "INSERT INTO `user`(`login_time`,`first_name`,`last_name`,`user_id`,`password`) VALUES" + str
	res, err := dao.DB.ExecContext(ctx, sqlSen, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *UserDAO) NewOne(row *sql.Row) (*User, error) {
	if err := row.Err(); err != nil {
		return nil, err
	}
	var val User
	err := row.Scan(&val.LoginTime, &val.FirstName, &val.LastName, &val.UserId, &val.Password)
	return &val, err
}

func (dao *UserDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*User, error) {
	row := dao.DB.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *UserDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*User, error) {
	s := "SELECT `login_time`,`first_name`,`last_name`,`user_id`,`password` FROM `user` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *UserDAO) NewBatch(rows *sql.Rows) ([]*User, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*User, 0, 5)
	for rows.Next() {
		var val User
		if err := rows.Scan(&val.LoginTime, &val.FirstName, &val.LastName, &val.UserId, &val.Password); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *UserDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*User, error) {
	rows, err := dao.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *UserDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*User, error) {
	s := "SELECT `login_time`,`first_name`,`last_name`,`user_id`,`password` FROM `user` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}

func (dao *UserDAO) UpdateColsByWhere(ctx context.Context, val *User, where string, args ...any) (int64, error) {
	newArgs, cols := dao.quotedNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user` SET " + cols + " WHERE " + where
	return dao.UpdateColByRaw(ctx, val, s, newArgs...)
}

func (dao *UserDAO) quotedNoneZero(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	if val.LoginTime != "" {
		args = append(args, val.LoginTime)
		cols = append(cols, "`login_time`")
	}
	if val.FirstName != "" {
		args = append(args, val.FirstName)
		cols = append(cols, "`first_name`")
	}
	if val.LastName != "" {
		args = append(args, val.LastName)
		cols = append(cols, "`last_name`")
	}
	if val.UserId != 0 {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}
	if len(val.Password) != 0 {
		args = append(args, val.Password)
		cols = append(cols, "`password`")
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) quotedNonePK(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	args = append(args, val.UserId)
	cols = append(cols, "`user_id`")
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) quotedSpecificCol(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	args = append(args, val.FirstName)
	cols = append(cols, "`first_name`")
	args = append(args, val.LastName)
	cols = append(cols, "`last_name`")
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) UpdateColByRaw(ctx context.Context, val *User, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
