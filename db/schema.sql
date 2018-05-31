BEGIN TRANSACTION;

DROP TABLE IF EXISTS chain;

CREATE TABLE chain (
  height INTEGER PRIMARY KEY,
  nonce BIGINT,
  --timestamp TIMESTAMP WITH time zone,
  timestamp CHARACTER(40),
  block_hash CHARACTER(44),
  prev_hash CHARACTER(44),
  data BYTEA
);
END TRANSACTION;
