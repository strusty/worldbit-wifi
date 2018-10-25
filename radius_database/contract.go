package radius_database

type AccountingStore interface {
	SessionTimeSum(username string) (int64, error)
	DeleteAccountingByUsername(username string) error
}

type CheckStore interface {
	Create(check *Check) error
	SessionChecks() ([]Check, error)
	DeleteChecksByUsername(username string) error
}

type ReplyStore interface {
	Create(check *Reply) error
	DeleteRepliesByUsername(username string) error
}
