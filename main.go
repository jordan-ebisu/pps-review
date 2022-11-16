package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/gocarina/gocsv"
)

// struct for the csv file
type Port struct {
	Ignore     string `csv:"delete_thisfield"`
	AlsoIgnore string `csv:"also_deletethisfield"`
	Port       string `csv:"port_number"`
	Server     string `csv:"server"`
}

// struct for the yaml file
type PortsAndServersRuleset struct {
	Rule []RuleItem `yaml:"rule"`
}

type RuleItem struct {
	Server   string `yaml:"server"`
	Protocol string `yaml:"protocol"`
	Min      int    `yaml:"allowed_tcp_port_min"`
	Max      int    `yaml:"allowed_tcp_port_max"`
}

func main() {
	fmt.Println("Starting the review")

	fmt.Println("Getting Yaml Key")

	var p []PortsAndServersRuleset

	yamlFile, err := ioutil.ReadFile("ports-and-servers-key.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		log.Fatal("Failed to parse file ", err)
	}
	// To access a specific port / server
	fmt.Printf("Port %v is allowed on %v\n", p[0].Rule[0].Max, p[0].Rule[0].Server)
	//fmt.Println(ports_and_server_ruleset)

	portsFile, err := os.OpenFile("ports-and-servers.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer portsFile.Close()

	ports := []*Port{}

	if err := gocsv.UnmarshalFile(portsFile, &ports); err != nil {
		panic(err)
	}

	/*for i, port := range ports {
		fmt.Printf("row %v: port %v is open on %v\n", i+2, port.Port, port.Server)
	}*/
	totalPortsReviewed := len(ports)
	fmt.Printf("The total number of combinations reviewed is %v\n", totalPortsReviewed)
}
