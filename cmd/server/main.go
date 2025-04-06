package main

import (
	"flag"
	"fmt"

	"github.com/DecodeFi/backend-logic/internal/alchemy"
	"github.com/DecodeFi/backend-logic/internal/db"
	"github.com/DecodeFi/backend-logic/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {

	port := flag.String("port", "", "a string")
	flag.Parse()
	if *port == "" {
		panic("MOO: port!")
	}

	alchemyCli := alchemy.NewAlchemyClient()
	db, err := db.NewDb(&db.DbCreds{
		User:     "admin",
		DbName:   "sber",
		Password: "0282347384973282qweyiu", // MOO: password in code!
	})

	if err != nil {
		panic(err) // MOO: handle!
	}

	server := server.NewHttpServer(alchemyCli, db)

	router := gin.Default()
	router.GET("/block_number", server.HandleBlockNumber)
	router.GET("/block/:id", server.HandleBlock)
	router.GET("/trace_block/:id", server.HandleGetBlockTraces)
	router.GET("/trace_address/:address", server.HandleAddressTraces)

	// TODO: remove (replace with scheduling mechanism)
	router.POST("/force_trace_block/:id", server.HandleForceTraceBlock)

	router.Run(fmt.Sprintf("localhost:%s", *port))
}
