package server

import (
	"net/http"
	"strconv"

	"github.com/DecodeFi/backend-logic/internal/alchemy"
	"github.com/DecodeFi/backend-logic/internal/db"
	"github.com/DecodeFi/backend-logic/internal/evm_inspect"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	serverBase
}

func NewHttpServer(alchemyCli *alchemy.AlchemyClient, db *db.Db) *HttpServer {
	return &HttpServer{serverBase: serverBase{alchemyCli: alchemyCli, db: db}}
}

/////////////////////////////////////////////////////////////////////////

func (h *HttpServer) HandleBlockNumber(c *gin.Context) {
	resp, err := h.alchemyCli.BlockNumber(alchemy.NewBlockNumberRequest())
	if err != nil {
		panic(err) // MOO: handle!
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func (h *HttpServer) HandleBlock(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.alchemyCli.Block(alchemy.NewBlockRequest(id))
	if err != nil {
		panic(err) // MOO: handle!
	}

	c.IndentedJSON(http.StatusOK, resp)
}

/////////////////////////////////////////////////////////////////////////

func (h *HttpServer) HandleForceTraceBlock(c *gin.Context) {
	id := c.Param("id")

	cli := evm_inspect.NewEvmInspectClient()
	traces, err := cli.TraceBlock(id)

	if err != nil {
		panic(err) // MOO: handle!
	}

	intId, _ := strconv.ParseInt(id, 10, 64)

	h.db.InsertTraces(id, traces)
	h.db.InsertBlock(&db.Block{BlockNumber: intId, Tag: db.BLOCK_TAG})

	c.IndentedJSON(http.StatusOK, traces)
}

func (h *HttpServer) HandleGetBlockTraces(c *gin.Context) {
	id := c.Param("id")

	// TODO: lookup in blocks table to early reject
	traces, err := h.db.GetBlockTraces(id)
	if err != nil {
		panic(err) // MOO: handle
	}

	for i := range traces {
		traces[i].Calldata = ""
	}

	c.JSON(http.StatusOK, traces)
}

func (h *HttpServer) HandleAddressTraces(c *gin.Context) {
	address := c.Param("address")

	directions := []int{db.TRACE_FROM, db.TRACE_TO}
	actions := []string{"call", "delegate_call", "create", "create2"}

	traces := make([]evm_inspect.Trace, 0)

	for _, d := range directions {
		for _, a := range actions {
			top := 20
			if a == "create" || a == "create2" {
				top = 50
			}
			t, err := h.db.GetAddressTraces(address, a, d, &db.Limit{Top: top, Offset: 0}, nil)
			if err != nil {
				panic(err) // MOO: handle
			}

			traces = append(traces, t...)
		}
	}

	c.JSON(http.StatusOK, traces)
}
