package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var pool *pgxpool.Pool

func init() {
	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("unable to connect to database", err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("init db pool")
}

func SelectWallets() (wallets []Wallet, err error) {
	rows, err := pool.Query(context.Background(), "SELECT name,address,nodes,balance,daily,\"group\" FROM wallet ORDER BY daily DESC,balance DESC")
	if err != nil {
		return nil, err
	}
	wallets, err = pgx.CollectRows(rows, pgx.RowToStructByName[Wallet])
	return wallets, nil
}

func InsertWallet(wallet *Wallet) (err error) {
	_, err = pool.Exec(context.Background(), "INSERT INTO wallet(name,address,\"group\") VALUES ($1,$2,$3)", wallet.Name, wallet.Address, wallet.Group)
	return
}

func Batch(sql string, params []string) error {
	batch := &pgx.Batch{}
	for _, p := range params {
		batch.Queue(sql, p)
	}
	br := pool.SendBatch(context.Background(), batch)
	if err := br.Close(); err != nil {
		return err
	}
	return nil
}

func SelectNodes() (nodes []Node, err error) {
	rows, err := pool.Query(context.Background(), "SELECT id,name,ip,bandwidth,traffic,price,renew,cpu,ram,disk,state,node_id,type,pool_id FROM node ORDER BY renew")
	if err != nil {
		return nil, err
	}
	nodes, err = pgx.CollectRows(rows, pgx.RowToStructByName[Node])
	return nodes, err
}

func SelectNodesCosts() (costs float32) {
	err := pool.QueryRow(context.Background(), "SELECT SUM(price) costs FROM node WHERE price > 0").Scan(&costs)
	if err != nil {
		log.Println(err)
	}
	return 0
}
