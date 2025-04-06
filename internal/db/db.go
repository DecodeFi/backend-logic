package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/DecodeFi/backend-logic/internal/evm_inspect"
)

type DbCreds struct {
	User     string
	Password string
	DbName   string
}

type Db struct {
	db *sql.DB
}

func NewDb(creds *DbCreds) (*Db, error) {
	credStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", creds.User, creds.Password, creds.DbName)
	db, err := sql.Open("postgres", credStr)
	if err != nil {
		return nil, err
	}

	return &Db{db: db}, nil
}

func (d *Db) Close() {
	d.db.Close()
}

type Block struct {
	BlockNumber int64
	Tag         int64
}

const BLOCK_TAG int64 = 1

func (d *Db) InsertBlock(block *Block) {
	_, err := d.db.Exec(`
		INSERT INTO blocks(block_number, tag) VALUES($1, $2) ON CONFLICT (block_number) DO UPDATE SET
			block_number=EXCLUDED.block_number,
			tag=EXCLUDED.tag
	`, block.BlockNumber, block.Tag)

	if err != nil {
		panic(err) // MOO
	}
}

var tracesSchema = []string{"trace_id", "tx_hash", "block_number", "from_addr", "to_addr", "storage_addr", "value", "action", "calldata"}

func (d *Db) InsertTraces(blockNumber string, traces []evm_inspect.Trace) {
	var schemab strings.Builder
	var updb strings.Builder

	schemab.WriteString("(")
	for i, col := range tracesSchema {
		schemab.WriteString(col)
		updb.WriteString(fmt.Sprintf("%s=EXCLUDED.%s", col, col))
		if i+1 != len(tracesSchema) {
			schemab.WriteString(", ")
			updb.WriteString(", ")
		}
	}
	schemab.WriteString(")")
	schema := schemab.String()

	upd := updb.String()

	var rowsb strings.Builder
	for r := range traces {
		var rowb strings.Builder
		rowb.WriteString("(")
		for i := range tracesSchema {
			rowb.WriteString(fmt.Sprintf("$%d", r*len(tracesSchema)+i+1))
			if i+1 != len(tracesSchema) {
				rowb.WriteString(", ")
			}
		}
		rowb.WriteString(")")

		rowsb.WriteString(rowb.String())

		if r+1 != len(traces) {
			rowsb.WriteString(", ")
		}
	}

	rows := rowsb.String()
	args := make([]interface{}, 0)

	for i, r := range traces {
		trace_id := fmt.Sprintf("%08x%s", i, r.TxHash)

		// TODO: not so nice...
		args = append(args, trace_id)
		args = append(args, r.TxHash)
		args = append(args, blockNumber)
		args = append(args, r.FromAddr)
		args = append(args, r.ToAddr)
		args = append(args, r.StorageAddr)
		args = append(args, r.Value)
		args = append(args, r.Action)
		args = append(args, r.Calldata)
	}

	query := fmt.Sprintf("INSERT INTO TRACES %s VALUES %s ON CONFLICT (trace_id) DO UPDATE SET %s", schema, rows, upd)

	_, err := d.db.Exec(query, args...)
	if err != nil {
		panic(err) // MOO
	}
}

func (d *Db) getTraces(query string, args ...interface{}) ([]evm_inspect.Trace, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var traces []evm_inspect.Trace
	for rows.Next() {
		var trace evm_inspect.Trace

		// TODO: can be not in sync with tracesSchema!
		if err := rows.Scan(&trace.TraceId, &trace.TxHash, &trace.BlockNumber, &trace.FromAddr,
			&trace.ToAddr, &trace.StorageAddr, &trace.Value, &trace.Action, &trace.Calldata); err != nil {
			return traces, err
		}
		traces = append(traces, trace)
	}
	err = rows.Err()
	return traces, err
}

func genItems() string {
	var itemsb strings.Builder
	for i, c := range tracesSchema {
		itemsb.WriteString(c)
		if i+1 < len(tracesSchema) {
			itemsb.WriteString(", ")
		}
	}

	return itemsb.String()
}

func (d *Db) GetBlockTraces(blockNumber string) ([]evm_inspect.Trace, error) {
	query := fmt.Sprintf("SELECT %s FROM traces WHERE block_number = $1", genItems())
	return d.getTraces(query, blockNumber)
}

const (
	TRACE_FROM = iota
	TRACE_TO   = iota
)

type Limit struct {
	Top    int
	Offset int
}

type BlockRange struct {
	BlockFrom int64
	BlockTo   int64
}

// TODO: the whole method is kinda messy, refactor it A LOT
func (d *Db) GetAddressTraces(address string, action string, direction int, limit *Limit, blockRange *BlockRange) ([]evm_inspect.Trace, error) {
	fromCol := "from_addr"
	toCol := "to_addr"

	if action == "delegate_call" {
		fromCol = "storage_addr"
	}

	filterCol := fromCol
	groupCol := toCol

	if direction == TRACE_TO {
		filterCol, groupCol = groupCol, filterCol
	}

	var qb strings.Builder
	qb.WriteString(fmt.Sprintf("SELECT %s, count(*) c FROM traces WHERE %s = $1 AND action = $2 ", groupCol, filterCol))

	if blockRange == nil {
		blockRange = &BlockRange{BlockFrom: 0, BlockTo: int64(1) << 60}
	}

	qb.WriteString("AND block_number >= $3 AND block_number <= $4 ")
	qb.WriteString(fmt.Sprintf("GROUP BY %s ORDER BY c DESC ", groupCol))

	if limit == nil {
		limit = &Limit{Top: 100, Offset: 0}
	}

	qb.WriteString("LIMIT $5 OFFSET $6")

	query := qb.String()

	rows, err := d.db.Query(query, address, action, blockRange.BlockFrom, blockRange.BlockTo, limit.Top, limit.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var traces []evm_inspect.Trace
	for rows.Next() {
		var groupRes string
		var counter int
		if err := rows.Scan(&groupRes, &counter); err != nil {
			return traces, err
		}

		// TODO: kinda ugly
		t := evm_inspect.Trace{
			Action: action,
		}
		if direction == TRACE_TO {
			t.ToAddr = address
			if action == "delegate_call" {
				t.StorageAddr = groupRes
			} else {
				t.FromAddr = groupRes
			}
		} else {
			t.ToAddr = groupRes
			if action == "delegate_call" {
				t.StorageAddr = address
			} else {
				t.FromAddr = address
			}
		}

		traces = append(traces, t)
	}
	err = rows.Err()
	return traces, err
}
