package code

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
)

type OrderDAO struct {
	DB *sql.DB
}

func (dao *OrderDAO) Insert(ctx context.Context, vals ...*Order) (int64, error) {
	var args = make([]interface{}, len(vals)*(3))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?)"
		args = append(args, v.OrderTime, v.OrderId, v.UserId)
	}
	sqlSen := "INSERT INTO `order`(`order_time`,`order_id`,`user_id`) VALUES" + str
	res, err := dao.DB.ExecContext(ctx, sqlSen, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (dao *OrderDAO) NewOne(row *sql.Row) (*Order, error) {
	if err := row.Err(); err != nil {
		return nil, err
	}
	var val Order
	err := row.Scan(&val.OrderTime, &val.OrderId, &val.UserId)
	return &val, err
}

func (dao *OrderDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*Order, error) {
	row := dao.DB.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *OrderDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id` FROM `order` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *OrderDAO) NewBatch(rows *sql.Rows) ([]*Order, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*Order, 0, 3)
	for rows.Next() {
		var val Order
		if err := rows.Scan(&val.OrderTime, &val.OrderId, &val.UserId); err != nil {
			return nil, err
		}
		vals = append(vals, &val)
	}
	return vals, nil
}

func (dao *OrderDAO) SelectBatchByRaw(ctx context.Context, query string, args ...any) ([]*Order, error) {
	rows, err := dao.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return dao.NewBatch(rows)
}

func (dao *OrderDAO) SelectBatchByWhere(ctx context.Context, where string, args ...any) ([]*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id` FROM `order` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}

func (dao *OrderDAO) UpdateByWhereWithNoneZero(ctx context.Context, val *Order, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + cols + " WHERE " + where
	return dao.UpdateByRaw(ctx, val, s, newArgs...)
}

func (dao *OrderDAO) UpdateNoneZero(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 3)
	var args = make([]interface{}, 0, 3)
	judge := func(x any) bool {
		return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
	}
	if judge(val.OrderTime) {
		args = append(args, val.OrderTime)
		cols = append(cols, "`order_time`")
	}

	if judge(val.OrderId) {
		args = append(args, val.OrderId)
		cols = append(cols, "`order_id`")
	}

	if judge(val.UserId) {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}

	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) UpdateByWhereWithPrimaryKey(ctx context.Context, val *Order, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateNonePrimaryKey(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + cols + " WHERE " + where
	return dao.UpdateByRaw(ctx, val, s, newArgs...)
}

func (dao *OrderDAO) UpdateNonePrimaryKey(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 3)
	var args = make([]interface{}, 0, 3)
	args = append(args, val.UserId)
	cols = append(cols, "`user_id`")

	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) UpdateByWhereWithSpecificCol(ctx context.Context, val *Order, where string, args ...any) (int64, error) {
	newArgs, cols := dao.UpdateBySpecificCol(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + cols + " WHERE " + where
	return dao.UpdateByRaw(ctx, val, s, newArgs...)
}

func (dao *OrderDAO) UpdateBySpecificCol(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 3)
	var args = make([]interface{}, 0, 3)
	args = append(args, val.OrderTime)
	cols = append(cols, "`order_time`")

	args = append(args, val.OrderId)
	cols = append(cols, "`order_id`")

	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) UpdateByRaw(ctx context.Context, val *Order, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
