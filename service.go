package main

import (
	"errors"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	log "github.com/Sirupsen/logrus"
)

func connectToNeo4J(batchSize int, neoURL string) (neoutils.NeoConnection, error) {
	conf := neoutils.DefaultConnectionConfig()
	conf.BatchSize = batchSize
	db, err := neoutils.Connect(neoURL, conf)
	if err != nil {
		log.Errorf("Could not connect to neo4j, error=[%s]\n", err)
	}
	log.Printf("Connected to NEO4J successfully")
	return db, nil
}

func updateOrganisation(neoConn neoutils.NeoConnection, canonicalUUID string, uuids []string) (bool, error) {

	// ensure that canonicalUUID does not exist, and has no relationships
	// No UPPIdentifier exists for the canonical uuid (aka its V1)
	canonicalNodeIsMissing, err := nodeIsMissing(canonicalUUID, neoConn)
	if err != nil {
		return false, err
	}

	if canonicalNodeIsMissing {
		log.Info("Canonical Node is Missing (Are you only a V1 UUID): ", canonicalUUID)
		// Find V2 Node
		// get a list of the non canonical uuids
		for i, refUUID := range uuids {
			if refUUID == canonicalUUID {
				uuids = append(uuids[:i], uuids[i+1:]...)
				break
			}
		}

		log.Info("Not canonicalUUIDs: ", uuids)

		// find the only existing node - ensure that there isn't more than one of them
		existingUUID, err := findExistingNode(uuids, neoConn)

		if err != nil {
			return false, err
		}

		// Update UUID property
		// update org with setting canonical uuid
		updateToCanonicalQuery(existingUUID, canonicalUUID, neoConn)

		// Ensure UPPidentifiers are present
		// add to the existing org UPP identifiers for the other missing nodes
		for _, refUUID := range uuids {
			if refUUID != existingUUID {
				addUPPIdentifierQuery(existingUUID, refUUID, neoConn)
			}
		}
	} else {
		log.Info("Canonical Node has been found (Are you only a V2 UUID): ", canonicalUUID)
		for _, refUUID := range uuids {
			if refUUID != canonicalUUID {
				addUPPIdentifierQuery(canonicalUUID, refUUID, neoConn)
			}

		}
	}

	return true, nil
}

func findExistingNode(uuids []string, neoConn neoutils.NeoConnection) (string, error) {
	existingNode := ""
	for _, refUUID := range uuids {
		refNodeIsMissing, err := nodeIsMissing(refUUID, neoConn)
		if err != nil {
			return existingNode, err
		}
		if !refNodeIsMissing && existingNode != "" {
			return existingNode, errors.New("there are more than 1 existing nodes: a normal way of concordance should be done!")
		} else if !refNodeIsMissing {
			existingNode = refUUID
		}
	}
	log.Printf("Existing node: %v", existingNode)
	return existingNode, nil
}
