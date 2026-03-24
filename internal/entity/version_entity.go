package entity

type Version struct {
	ID            int64 `xorm:"not null pk autoincr 'id'"`
	VersionNumber int64 `xorm:"not null default 0 'version_number'"`
}

func (Version) TableName() string {
	return "version"
}
