package migrations

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/entity"
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
	m.do("insert version", m.insertVersion)
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

func (m *Mentor) hasVersion() bool {
	exist, err := m.engine.Context(m.ctx).IsTableExist("version")
	if err != nil || !exist {
		return false
	}
	count, err := m.engine.Context(m.ctx).Count(&entity.Version{})
	if err != nil || count > 0 {
		return true
	}
	return false
}

func (m *Mentor) insertVersion() {
	if _, m.err = m.engine.Context(m.ctx).
		Exec(`DELETE FROM "version" WHERE "id" > ?`, 1); m.err != nil {
		return
	}

	_, m.err = m.engine.Context(m.ctx).Exec(
		`INSERT INTO "version" ("id", "version_number") VALUES (?, ?) ON CONFLICT ("id") DO UPDATE SET "version_number" = EXCLUDED."version_number"`,
		1,
		0,
	)
}
