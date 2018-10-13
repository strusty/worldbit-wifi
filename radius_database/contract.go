package radius_database

type CheckStore interface {
	Create(check *Check) error
}

type ReplyStore interface {
	Create(check *Reply) error
}
