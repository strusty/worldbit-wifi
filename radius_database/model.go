package radius_database

type Reply struct {
	ID        uint `gorm:"primary_key"`
	Username  string
	Attribute string
	Op        string
	Value     string
}

func (Reply) TableName() string {
	return "radreply"
}

type Check Reply

func (Check) TableName() string {
	return "radcheck"
}
