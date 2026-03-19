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
	return &Mentor{engine: engine}
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
	m.err = m.engine.Context(m.ctx).Sync2(tables...)

}
