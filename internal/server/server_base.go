package server

import (
	"github.com/DecodeFi/backend-logic/internal/alchemy"
	"github.com/DecodeFi/backend-logic/internal/db"
)

type serverBase struct {
	alchemyCli *alchemy.AlchemyClient
	db         *db.Db
}
