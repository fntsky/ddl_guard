package migrations

import (
	"context"
	"fmt"

	"xorm.io/xorm"
)

type Mentor struct {
	ctx    context.Context
	engine *xorm.Engine
	err    error
}

func NewMentor(engine *xorm.Engine) *Mentor {
	return &Mentor{
		ctx:    context.Background(),
		engine: engine,
	}
}

func (m *Mentor) InitDB() error {
	m.do("sync tables", m.syncTables)
	return m.err
}

func (m *Mentor) do(desc string, fn func()) {
	if m.err != nil {
		return
	}
	fn()
	if m.err != nil {
		m.err = fmt.Errorf("failed to %s: %w", desc, m.err)
	}
}

func (m *Mentor) syncTables() {
	ctx := m.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	m.err = m.engine.Context(ctx).Sync2(tables...)

}
