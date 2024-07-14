package db

import "gosmic/db/drivers/mongodb"

type Collections struct {
	Nodes *mongodb.Collection
	Ways  *mongodb.Collection
}
