package main

import (
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
)

type countResult []struct {
	Count int `json:"c"`
}

func nodeIsMissing(uuid string, neoConn neoutils.NeoConnection) (bool, error) {

	results := countResult{}
	err := neoConn.CypherBatch([]*neoism.CypherQuery{{
		Statement: `match (uppId:UPPIdentifier{value:{uuid}})-[:IDENTIFIES]->(a:Thing) return count(a) as c`,
		Parameters: map[string]interface{}{
			"uuid": uuid,
		},
		Result: &results,
	}})

	if err != nil {
		return false, err
	}

	if results[0].Count != 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func updateToCanonicalQuery(oldUUID string, canonicalUUID string, neoConn neoutils.NeoConnection) error {
	return neoConn.CypherBatch([]*neoism.CypherQuery{{
		Statement: `match (a:Organisation)<-[:IDENTIFIES]-(i:UPPIdentifier{value:{oldUUID}})
		set a.uuid = {canonicalUUID}
		create (uppId:UPPIdentifier:Identifier{value:{canonicalUUID}})-[:IDENTIFIES]->(a)
		return a`,
		Parameters: map[string]interface{}{
			"oldUUID":       oldUUID,
			"canonicalUUID": canonicalUUID,
		},
	}})

}

func addUPPIdentifierQuery(oldUUID string, newUUID string, neoConn neoutils.NeoConnection) error {
	return neoConn.CypherBatch([]*neoism.CypherQuery{{
		Statement: `match (a:Organisation)<-[:IDENTIFIES]-(i:UPPIdentifier{value:{oldUUID}})
		create (uppId:UPPIdentifier:Identifier{value:{newUUID}})-[i:IDENTIFIES]->(a)`,
		Parameters: map[string]interface{}{
			"oldUUID": oldUUID,
			"newUUID": newUUID,
		},
	}})
}
