package main

import (
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gocql/gocql"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage: writer topic ALL|ONE|QUORUM seed_node_ip_or_name")
	}

	topic := os.Args[1]

	consistency := gocql.All
	switch strings.ToUpper(os.Args[2]) {
	case "ALL":
	case "ONE":
		consistency = gocql.One
	case "QUORUM":
		consistency = gocql.Quorum
	default:
		log.Fatalf("Unknown consistency level %s", os.Args[2])
	}

	seed := os.Args[3]
	log.Printf(
		"Connecting cluster at %s with consistency %s for topic %s",
		seed, consistency, topic)

	cluster := gocql.NewCluster(seed)
	cluster.Consistency = consistency
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Cannot connect to cluster at %s: %v", seed, err)
	}
	defer session.Close()

	var clusterName string
	if err := session.Query(
		"SELECT cluster_name FROM system.local").
		Scan(&clusterName); err != nil {
		log.Fatalf("Cannot query cluster: %v", err)
	}
	log.Printf("Connected to cluster %s", clusterName)

	if err := session.Query(
		`CREATE KEYSPACE IF NOT EXISTS ece473
			WITH replication = {
				'class':'SimpleStrategy',
				'replication_factor':3}`).
		Exec(); err != nil {
		log.Fatalf("Cannot create keyspace ece473: %v", err)
	}

	if err := session.Query(
		`CREATE TABLE IF NOT EXISTS ece473.prj03 (
			topic text, seq int, value double,
			PRIMARY KEY (topic, seq))`).
		Exec(); err != nil {
		log.Fatalf("Cannot create table ece473.prj03: %v", err)
	}

	for seq := 1; ; seq++ {
		value := rand.Float64()
		err := session.Query(
			`INSERT INTO ece473.prj03 (topic, seq, value) VALUES (?, ?, ?)`,
			topic, seq, value).
			Exec()
		if err != nil {
			log.Fatalf("Cannot write %d to table ece473.prj03: %v", seq, err)
		}
		if seq%1000 == 0 {
			log.Printf("%s: inserted %d rows", topic, seq)
		}
	}
}
