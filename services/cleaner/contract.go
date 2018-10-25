package cleaner

import "time"

type Cleaner interface {
	Start(period time.Duration)
}
