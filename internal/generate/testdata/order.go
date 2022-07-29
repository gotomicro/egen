package code

import (
	"context"
	"database/sql"
	"strings"
)

type OrderDAO struct {
	DB *sql.DB
}

func (dao *OrderDAO) Insert(ctx context.Context, vals ...*Order) (int64, error) {
	var args = make([]interface{}, len(vals)*(5))
	var str = ""
	for k, v := range vals {
		if k != 0 {
			str += ", "
		}
		str += "(?,?,?,?,?)"
		args = append(args, v.OrderTime, v.OrderId, v.UserId, v.HasBuy, v.Price)
	}
	sqlSen := "INSERT INTO `order`(`order_time`,`order_id`,`user_id`,`has_buy`,`price`) VALUES" + str
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
	err := row.Scan(&val.OrderTime, &val.OrderId, &val.UserId, &val.HasBuy, &val.Price)
	return &val, err
}

func (dao *OrderDAO) SelectByRaw(ctx context.Context, query string, args ...any) (*Order, error) {
	row := dao.DB.QueryRowContext(ctx, query, args...)
	return dao.NewOne(row)
}

func (dao *OrderDAO) SelectByWhere(ctx context.Context, where string, args ...any) (*Order, error) {
	s := "SELECT `order_time`,`order_id`,`user_id`,`has_buy`,`price` FROM `order` WHERE " + where
	return dao.SelectByRaw(ctx, s, args...)
}

func (dao *OrderDAO) NewBatch(rows *sql.Rows) ([]*Order, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var vals = make([]*Order, 0, 5)
	for rows.Next() {
		var val Order
		if err := rows.Scan(&val.OrderTime, &val.OrderId, &val.UserId, &val.HasBuy, &val.Price); err != nil {
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
	s := "SELECT `order_time`,`order_id`,`user_id`,`has_buy`,`price` FROM `order` WHERE " + where
	return dao.SelectBatchByRaw(ctx, s, args...)
}

func (dao *OrderDAO) UpdateColsByWhere(ctx context.Context, val *Order, where string, args ...any) (int64, error) {
	newArgs, cols := dao.quotedNoneZero(val)
	newArgs = append(newArgs, args...)
	s := "UPDATE `order` SET " + cols + " WHERE " + where
	return dao.UpdateColByRaw(ctx, val, s, newArgs...)
}

func (dao *OrderDAO) quotedNoneZero(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	if val.OrderTime != "" {
		args = append(args, val.OrderTime)
		cols = append(cols, "`order_time`")
	}
	if val.OrderId != 0 {
		args = append(args, val.OrderId)
		cols = append(cols, "`order_id`")
	}
	if val.UserId != 0 {
		args = append(args, val.UserId)
		cols = append(cols, "`user_id`")
	}
	if val.HasBuy != false {
		args = append(args, val.HasBuy)
		cols = append(cols, "`has_buy`")
	}
	if val.Price != 0.0 {
		args = append(args, val.Price)
		cols = append(cols, "`price`")
	}
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) quotedNonePK(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	args = append(args, val.UserId)
	cols = append(cols, "`user_id`")
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) quotedSpecificCol(val *Order) ([]interface{}, string) {
	var cols = make([]string, 0, 5)
	var args = make([]interface{}, 0, 5)
	args = append(args, val.OrderTime)
	cols = append(cols, "`order_time`")
	args = append(args, val.OrderId)
	cols = append(cols, "`order_id`")
	return args, strings.Join(cols, "=?,")
}

func (dao *OrderDAO) UpdateColByRaw(ctx context.Context, val *Order, query string, args ...any) (int64, error) {
	res, err := dao.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
