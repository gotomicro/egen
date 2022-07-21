package code

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
)

type UserDAO struct {
	DB *sql.DB
}

func (dao *UserDAO) UpdateByWhereWithNoneZero(ctx context.Context, val *User, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user` SET " + cols + " WHERE " + where
	return dao.UpdateByRaw(ctx, val, s, newArgs...)
}

func (dao *UserDAO) UpdateNoneZero(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 4)
	var args = make([]interface{}, 0, 4)
	judge := func(x any) bool {
		return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
	}
	if judge(val.LoginTime) {
		args = append(args, val.LoginTime)
		cols = append(cols, "`login_time`")
	}
	
	if judge(val.FirstName) {
		args = append(args, val.FirstName)
		cols = append(cols, "`first_name`")
	}
	
	if judge(val.LastName) {
		args = append(args, val.LastName)
		cols = append(cols, "`last_name`")
	}
	
	if judge(val.UserId) {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}
	
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) UpdateByWhereWithPrimaryKey(ctx context.Context, val *User, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateNonePrimaryKey(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user` SET " + cols + " WHERE " + where
	return dao.UpdateByRaw(ctx, val, s, newArgs...)
}

func (dao *UserDAO) UpdateNonePrimaryKey(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 4)
	var args = make([]interface{}, 0, 4)
	args = append(args, val.UserId)
	cols = append(cols, "`user_id`")
	
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) UpdateByWhereWithSpecificCol(ctx context.Context, val *User, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateBySpecificCol(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `user` SET " + cols + " WHERE " + where
	return dao.UpdateByRaw(ctx, val, s, newArgs...)
}

func (dao *UserDAO) UpdateBySpecificCol(val *User) ([]interface{}, string) {
	var cols = make([]string, 0, 4)
	var args = make([]interface{}, 0, 4)
	args = append(args, val.FirstName)
	cols = append(cols, "`first_name`")
	
	args = append(args, val.LastName)
	cols = append(cols, "`last_name`")
	
	return args, strings.Join(cols, "=?,")
}

func (dao *UserDAO) UpdateByRaw(ctx context.Context, val *User, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
