package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage: reader topic ALL|ONE|QUORUM seed_node_ip_or_name")
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

	for lastSeq, rows := 0, 0; ; {
		scanner := session.Query(
			`SELECT seq FROM ece473.prj03 WHERE topic = ? AND seq > ?`,
			topic, lastSeq).
			Iter().Scanner()
		seq := lastSeq
		for scanner.Next() {
			err := scanner.Scan(&seq)
			if err != nil {
				log.Fatalf("Cannot read after %d from table ece473.prj03: %v", lastSeq, err)
			}
			rows++
		}
		if seq == rows {
			log.Printf("%s: seq %d with %d rows", topic, seq, rows)
		} else {
			log.Printf("%s: seq %d with %d rows, missing %d", topic, seq, rows, seq-rows)
		}
		if seq != lastSeq {
			lastSeq = seq
		} else {
			log.Printf("%s: no more data, wait 10s", topic)
			time.Sleep(10 * time.Second)
		}
	}
}
