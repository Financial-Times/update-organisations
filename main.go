package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func readUUIDS(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	uuids := strings.Split(string(content), "\n")
	log.Printf("Uuids: %v", uuids)
	return uuids, nil
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 128,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	},
}

func getCompositeOrgModel(baseURL string, id string) (organisation, error) {
	req, _ := http.NewRequest("GET", baseURL+id, nil)

	resp, err := httpClient.Do(req)
	if err != nil {
		return organisation{}, fmt.Errorf("Could not get concept with uuid: %v (%v)", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		return organisation{}, fmt.Errorf("Could not get concept %v from %v. Returned %v", id, baseURL, resp.StatusCode)
	}

	var org organisation

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&org); err != nil {
		return organisation{}, fmt.Errorf("Error decoding response: %+v", err)
	}

	resp.Body.Close()
	if err != nil {
		return organisation{}, fmt.Errorf("Could not get concept with uuid: %v (%v)", id, err)
	}
	log.Printf("Organisation found in the comp transformer: %v", org)
	return org, nil
}

func main() {

	app := cli.App("update-organisations-neo4j", "A simple app for hacking the relationship transfer between apps")
	neoURL := app.String(cli.StringOpt{
		Name:   "neo-url",
		Value:  "http://localhost:7474/db/data",
		Desc:   "neo4j endpoint URL",
		EnvVar: "NEO_URL",
	})
	batchSize := app.Int(cli.IntOpt{
		Name:   "batchSize",
		Value:  1024,
		Desc:   "Maximum number of statements to execute per batch",
		EnvVar: "BATCH_SIZE",
	})
	uuidsPath := app.String(cli.StringOpt{
		Name:   "uuids",
		Value:  "uuids.txt",
		Desc:   "Uuids separated by \n that need to be updated",
		EnvVar: "UUIDS",
	})
	transformerUrl := app.String(cli.StringOpt{
		Name:   "composite-transformer-url",
		Value:  "https://pub-pre-prod-uk-up.ft.com/__composite-orgs-transformer/transformers/organisations/",
		Desc:   "Composite org transformer",
		EnvVar: "COMPOSITE_TRANSFORMER_URL",
	})
	app.Action = func() {
		log.Printf("Params! %s | %d | %s | %s | %s", *neoURL, *batchSize, *uuidsPath, *transformerUrl)

		uuids, err := readUUIDS(*uuidsPath)
		if err != nil {
			log.Error(err)
		}

		//connect to neo4j
		db, err := connectToNeo4J(*batchSize, *neoURL)
		if err != nil {
			log.Error(err)
		}

		for _, uuid := range uuids {
			if len(uuid) > 5 {
				log.Println("**** Starting updates for: ", uuid)

				//make a call for the first uuid:
				org, err := getCompositeOrgModel(*transformerUrl, uuid)
				log.Println("     -> alternative uuids: ", org.AlternativeIdentifiers.UUIDS)

				//execute queries
				updated, err := updateOrganisation(db, org.UUID, org.AlternativeIdentifiers.UUIDS)
				if err != nil {
					log.Error(err)
				}

				if updated {
					log.Println("     The update was successful")
				} else {
					log.Println("     The update failed for %s", uuid)
				}
			}
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Println("Application started with args %s", os.Args)

	app.Run(os.Args)
}
