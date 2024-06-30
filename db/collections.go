package db

import "gosmic/db/drivers/mongodb"

type Collections struct {
	Objects *mongodb.Collection
}
