package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/gocarina/gocsv"
)

type Port struct {
	Ignore     string `csv:"delete_thisfield"`
	AlsoIgnore string `csv:"also_deletethisfield"`
	Protocol   string `csv:"protocol"`
	Port       int    `csv:"port_number"`
	Server     string `csv:"server"`
}

type PortsAndServersRuleset struct {
	Rule []RuleItem `yaml:"rule"`
}

type RuleItem struct {
	Server   string `yaml:"server"`
	Protocol string `yaml:"protocol"`
	TCPMin   int    `yaml:"allowed_tcp_port_min"`
	TCPMax   int    `yaml:"allowed_tcp_port_max"`
	UDPMin   int    `yaml:"allowed_udp_port_min"`
	UDPMax   int    `yaml:"allowed_udp_port_max"`
}

func main() {
	fmt.Println("Starting the review")

	p := getYaml()

	// To access a specific port / server
	//fmt.Printf("Port %v is allowed on %v\n", p[0].Rule[0].TCPMax, p[0].Rule[0].Server)

	ports := getCSV()

	for _, port := range ports {

		if port.Port == 0 {
			continue
		} else {

			switch {
			case strings.HasPrefix(port.Server, "splunk"):
				//fmt.Printf("row %v of the csv: This is a splunk server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[2].Rule[0].TCPMin, p[2].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				} else {
					if comparePort(port.Port, p[2].Rule[0].UDPMin, p[2].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				}
			case strings.HasPrefix(port.Server, "ssh"):
				//fmt.Printf("row %v of the csv: This is a ssh server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[1].Rule[0].TCPMin, p[1].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				} else {
					if comparePort(port.Port, p[1].Rule[0].UDPMin, p[1].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				}
			case strings.HasPrefix(port.Server, "lb"):
				//fmt.Printf("row %v of the csv: This is a lb server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[0].Rule[0].TCPMin, p[0].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				} else {
					if comparePort(port.Port, p[0].Rule[0].UDPMin, p[0].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				}
			case strings.HasPrefix(port.Server, "prom"):
				//fmt.Printf("row %v of the csv: This is a prom server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[3].Rule[0].TCPMin, p[3].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				} else {
					if comparePort(port.Port, p[3].Rule[0].UDPMin, p[3].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
					}
				}
			default:
				fmt.Printf("I'm not sure what this is? (%v)\n", port.Server)
			}
		}
		//fmt.Printf("row %v: port %v is open on %v\n", i+2, port.Port, port.Server)
	}
	totalPortsReviewed := len(ports)
	fmt.Printf("Review completed: total number of combinations reviewed was %v\n", totalPortsReviewed)
}

func getYaml() []PortsAndServersRuleset {
	fmt.Println("Getting the YAML file")

	var p []PortsAndServersRuleset

	yamlFile, err := ioutil.ReadFile("ports-and-servers-key.yaml")
	if err != nil {
		log.Fatal("Failed to Read the Yaml File", err)
	}
	err = yaml.Unmarshal(yamlFile, &p)
	if err != nil {
		log.Fatal("Failed to parse file ", err)
	}

	return p
}

func getCSV() []*Port {
	fmt.Println("Getting the CSV file")
	portsFile, err := os.OpenFile("ports-and-servers.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer portsFile.Close()

	ports := []*Port{}

	if err := gocsv.UnmarshalFile(portsFile, &ports); err != nil {
		panic(err)
	}
	return ports
}

func comparePort(csvPort int, keyPortMin int, keyPortMax int) bool {
	var r bool
	switch {
	case keyPortMin == keyPortMax:
		if csvPort == keyPortMin {
			r = true
		} else {
			r = false
		}
	case keyPortMin != keyPortMax:
		if (csvPort >= keyPortMin) && (csvPort <= keyPortMax) {
			r = true
		} else {
			r = false
		}
	}
	return r
}
