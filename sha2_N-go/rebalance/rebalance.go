package rebalance

import (
	"fmt"
	"log"
	"sha-go/config"
	"sha-go/node"
	"sha-go/ring"
)

func BalanceAddition(newNode *node.Node, rightNode *node.Node, n int) {
	selectQuery := fmt.Sprintf("SELECT key,value,hash FROM files_%d", rightNode.Slot)
	rightNodeRows, err := rightNode.DB.Query(selectQuery)
	if err != nil {
		log.Println("Error in right Node table query: ",err)
		return
	}
	defer rightNodeRows.Close()

	var recordsToMove []config.FileRecord
	for rightNodeRows.Next() {
		var record config.FileRecord
		err := rightNodeRows.Scan(&record.Key, &record.Value, &record.Hash)
		if err != nil {
			log.Println("Error scanning row: ",err)
			continue
		}

		recordSlot := ring.GetSlot(record.Hash,n)
		if recordSlot <= newNode.Slot{
			recordsToMove = append(recordsToMove, record)
		}
	}

	if len(recordsToMove) > 0{
		insertQuery := fmt.Sprintf("INSERT INTO files_%d (key, value, hash) VALUES ($1, $2, $3)", newNode.Slot)
		deleteQuery := fmt.Sprintf("DELETE FROM files_%d WHERE key = $1", rightNode.Slot)

		for _,record := range recordsToMove{
			_,err := newNode.DB.Exec(insertQuery, record.Key, record.Value, record.Hash)
			if err != nil {
				log.Println("Error inserting into new node: ",err)
				continue
			}

			_,err = rightNode.DB.Exec(deleteQuery, record.Key)
			if err != nil {
				log.Println("Error deleting from right node: ",err)
			}
		}

		log.Printf("Moved %d records to new node (slot %d)",len(recordsToMove),newNode.Slot)
	}
}


func BalanceDeletion(tobeDeletedNode *node.Node, rightNode *node.Node){
	deleteAndReturnQuery := fmt.Sprintf("DELETE FROM files_%d RETURNING key,value,hash",tobeDeletedNode.Slot)
	removedRows,err := tobeDeletedNode.DB.Query(deleteAndReturnQuery) 
	if err != nil{
		log.Println("Error in querying rows from to be deleted node: ",err)
		return
	}
	defer removedRows.Close()
	var recordsToMove []config.FileRecord
	for removedRows.Next(){
		var recordToMove config.FileRecord
		err := removedRows.Scan(&recordToMove.Key,&recordToMove.Value,&recordToMove.Hash)
		if err != nil {
			log.Println("Error scanning a row of the sql result during balancing deletion: ",err)
			continue
		}
		recordsToMove = append(recordsToMove, recordToMove)
	}
	if len(recordsToMove)>0{
		insertQuery := fmt.Sprintf("INSERT INTO files_%d (key, value, hash) VALUES ($1, $2, $3)",rightNode.Slot)

		for _,record := range(recordsToMove){
			_,err := rightNode.DB.Exec(insertQuery,record.Key,record.Value,record.Hash)
			if err != nil {
				log.Printf("Error inserting into new node: %v", err)
				continue
			}
		}
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS files_%d", tobeDeletedNode.Slot)
		tobeDeletedNode.DB.Exec(dropQuery)
		log.Printf("Moved %d records to rightmost node (slot %d) and deleted node at slot %d",len(recordsToMove),rightNode.Slot,tobeDeletedNode.Slot)
	}
}