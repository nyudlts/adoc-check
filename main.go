package main

import (
	"flag"
	"fmt"
	"os"

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
					fmt.Println(do.DigitalObjectID, "ERROR")
				} else {
					fmt.Println(do.DigitalObjectID, "OK")
				}
			}
		}
	}

}
