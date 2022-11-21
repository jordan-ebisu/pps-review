package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	TCPOther []int  `yaml:"allowed_tcp_port_other"`
	UDPMin   int    `yaml:"allowed_udp_port_min"`
	UDPMax   int    `yaml:"allowed_udp_port_max"`
	UDPOther []int  `yaml:"allowed_udp_port_othe"`
}

func main() {
	fmt.Println("Starting the review")

	p := getYaml()

	// To access a specific port / server
	//fmt.Printf("Port %v is allowed on %v\n", p[0].Rule[0].TCPMax, p[0].Rule[0].Server)

	ports := getCSV()
	year, month, _ := time.Now().Date()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	yearString := strconv.Itoa(year)
	path += "/" + "pps-" + month.String() + "-" + yearString + ".csv"

	fmt.Println(path)

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
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				} else {
					if comparePort(port.Port, p[2].Rule[0].UDPMin, p[2].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)

						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				}
			case strings.HasPrefix(port.Server, "ssh"):
				//fmt.Printf("row %v of the csv: This is a ssh server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[1].Rule[0].TCPMin, p[1].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				} else {
					if comparePort(port.Port, p[1].Rule[0].UDPMin, p[1].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				}
			case strings.HasPrefix(port.Server, "lb"):
				//fmt.Printf("row %v of the csv: This is a lb server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[0].Rule[0].TCPMin, p[0].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				} else {
					if comparePort(port.Port, p[0].Rule[0].UDPMin, p[0].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				}
			case strings.HasPrefix(port.Server, "prom"):
				//fmt.Printf("row %v of the csv: This is a prom server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[3].Rule[0].TCPMin, p[3].Rule[0].TCPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
					}
				} else {
					if comparePort(port.Port, p[3].Rule[0].UDPMin, p[3].Rule[0].UDPMax) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, path)
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

	// TODO create function that deletes the file
	time.Sleep(10 * time.Second)
	e := os.Remove(path)
	if e != nil {
		log.Fatal(e)
	}

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

// updateCSV checks if the file for the month exists. If it does then it appends the record. If it does not, it creates the file then writes the information
func updateCSV(record []string, path string) {
	if _, err := os.Stat(path); err == nil {
		// File exists
		f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		csvWriter := csv.NewWriter(f)

		csvWriter.Write(record)
		csvWriter.Flush()

	} else {
		// File doesn't exist
		f, err := os.Create(path)
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		defer f.Close()

		csvWriter := csv.NewWriter(f)

		csvWriter.Write(record)
		csvWriter.Flush()

	}
}
