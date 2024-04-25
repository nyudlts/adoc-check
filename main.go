package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"os"
	"strings"

	"github.com/nyudlts/go-aspace"
)

var (
	workOrderLocation string
	config            string
)

func init() {
	flag.StringVar(&workOrderLocation, "work-order", "", "")
	flag.StringVar(&config, "config", "", "")
}

func main() {
	flag.Parse()

	client, err := aspace.NewClient(config, "prod", 20)
	if err != nil {
		panic(err)
	}

	workOrder, _ := os.Open(workOrderLocation)
	defer workOrder.Close()
	wo := aspace.WorkOrder{}
	if err := wo.Load(workOrder); err != nil {
		panic(err)
	}

	var b bytes.Buffer
	out := csv.NewWriter(bufio.NewWriter(&b))
	out.Comma = '\t'

	for _, row := range wo.Rows {
		repoId, aoURI, err := aspace.URISplit(row.GetURI())
		if err != nil {
			panic(err)
		}

		ao, err := client.GetArchivalObject(repoId, aoURI)
		if err != nil {
			panic(err)
		}

		for _, instance := range ao.Instances {
			if instance.InstanceType == "digital_object" {
				doURI := instance.DigitalObject["ref"]
				_, doID, err := aspace.URISplit(doURI)
				if err != nil {
					panic(err)
				}

				do, err := client.GetDigitalObject(repoId, doID)
				if err != nil {
					panic(err)
				}

				if do.DigitalObjectID != row.GetComponentID() {
					out.Write([]string{row.GetURI(), do.URI, do.DigitalObjectID, "ERROR"})
				} else {
					out.Write([]string{row.GetURI(), do.Title, do.URI, do.DigitalObjectID, "OK"})
				}
				out.Flush()
			}
		}
	}
	strings.Trim(b.String(), "\n")
	os.WriteFile("adoc-check.tsv", b.Bytes(), 0777)
}
