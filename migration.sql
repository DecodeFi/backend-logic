-- TODO: remove from this repo?

-- Remove tables
DROP TABLE IF EXISTS blocks;
DROP TABLE IF EXISTS traces;

-- Blocks are stored separately and simply
CREATE TABLE blocks (
    block_number BIGINT PRIMARY KEY,
    tag BIGINT
);

-- Main traces table

CREATE TABLE traces (
    trace_id VARCHAR(74) PRIMARY KEY, -- 8 for trace no inside block +  66 for hash
    tx_hash VARCHAR(66) NOT NULL,
    block_number BIGINT NOT NULL,
    from_addr VARCHAR(42) NOT NULL,
    to_addr VARCHAR(42) NOT NULL,
    storage_addr VARCHAR(42) NOT NULL,
    value VARCHAR(32) NOT NULL,
    action VARCHAR(16) NOT NULL,
    calldata TEXT
);

-- TODO: optimize indexes
CREATE INDEX idx_traces_tx_hash ON traces (tx_hash);
CREATE INDEX idx_traces_block_number ON traces (block_number);
CREATE INDEX idx_traces_from_hash ON traces (from_addr);
CREATE INDEX idx_traces_to_hash ON traces (to_addr);
CREATE INDEX idx_traces_storage_hash ON traces (storage_addr);
