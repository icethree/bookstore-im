package model

import (
	"context"
	"gorm.io/gorm"
)

const (
	DBTransKey = "DBTransManager"
)

type dbTran struct {
	s *gorm.DB
	c int64
}

func GetDB(ctx context.Context) *gorm.DB {
	return DB.Session(&gorm.Session{Context: ctx})
}

type TxCtxFunc func(ctx context.Context, s *gorm.DB) error

func newTrans(ctx context.Context) (context.Context, *dbTran) {
	m, ok := ctx.Value(DBTransKey).(*dbTran)
	if !ok {
		m = &dbTran{
			s: DB,
			c: 1,
		}
		return context.WithValue(ctx, DBTransKey, m), m
	}
	m.c++

	return context.WithValue(ctx, DBTransKey, m), m
}

func rollBack(ctx context.Context) error {
	m, ok := ctx.Value(DBTransKey).(*dbTran)
	if !ok {
		return nil
	}

	if m.c == 1 {
		return m.s.Rollback().Error
	}

	m.c--
	return nil
}

func RunNestedTx(ctx context.Context, t TxCtxFunc) error {
	nCtx, m := newTrans(ctx)
	if m.c == 1 {
		m.s.Begin()
	}

	if err := t(nCtx, m.s); err != nil {
		rollBack(nCtx)
		return err
	}

	if m.c == 1 {
		m.s.Commit()
	} else {
		m.c--
	}

	return nil
}
