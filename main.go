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

	"github.com/andygrunwald/go-jira"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
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
	UDPOther []int  `yaml:"allowed_udp_port_other"`
}

func main() {
	p := getYaml()

	ports := getCSV()
	year, month, _ := time.Now().Date()

	deviationsFilePath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	yearString := strconv.Itoa(year)
	deviationsFilePath += "/" + "pps-" + strings.ToLower(month.String()) + "-" + yearString + ".csv"

	fmt.Println(deviationsFilePath)

	for i, port := range ports {

		if port.Port == 0 {
			continue
		} else {
			switch {
			case strings.HasPrefix(port.Server, "splunk"):
				//fmt.Printf("row %v of the csv: This is a splunk server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[2].Rule[0].TCPMin, p[2].Rule[0].TCPMax) || containsPorts(p[2].Rule[0].TCPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				} else {
					if comparePort(port.Port, p[2].Rule[0].UDPMin, p[2].Rule[0].UDPMax) || containsPorts(p[2].Rule[0].UDPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)

						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				}
			case strings.HasPrefix(port.Server, "ssh"):
				//fmt.Printf("row %v of the csv: This is a ssh server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[1].Rule[0].TCPMin, p[1].Rule[0].TCPMax) || containsPorts(p[1].Rule[0].TCPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				} else {
					if comparePort(port.Port, p[1].Rule[0].UDPMin, p[1].Rule[0].UDPMax) || containsPorts(p[1].Rule[0].UDPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				}
			case strings.HasPrefix(port.Server, "lb"):
				//fmt.Printf("row %v of the csv: This is a lb server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[0].Rule[0].TCPMin, p[0].Rule[0].TCPMax) || containsPorts(p[0].Rule[0].TCPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				} else {
					if comparePort(port.Port, p[0].Rule[0].UDPMin, p[0].Rule[0].UDPMax) || containsPorts(p[0].Rule[0].UDPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				}
			case strings.HasPrefix(port.Server, "prom"):
				//fmt.Printf("row %v of the csv: This is a prom server\n", i+2)
				if port.Protocol == "tcp" {
					if comparePort(port.Port, p[3].Rule[0].TCPMin, p[3].Rule[0].TCPMax) || containsPorts(p[3].Rule[0].TCPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}

						updateCSV(record, deviationsFilePath)
					}
				} else {
					if comparePort(port.Port, p[3].Rule[0].UDPMin, p[3].Rule[0].UDPMax) || containsPorts(p[3].Rule[0].UDPOther, port.Port) {
						continue
					} else {
						fmt.Printf("hey this port shouldn't be here! (%v %v on %v)\n", port.Protocol, port.Port, port.Server)
						record := []string{
							port.Protocol, strconv.Itoa(port.Port), port.Server,
						}
						updateCSV(record, deviationsFilePath)
					}
				}
			case i == len(ports)-1:
				fmt.Println("everything looks good")
			default:
				fmt.Printf("I'm not sure what this is? (%v)\n", port.Server)
				record := []string{
					port.Protocol, strconv.Itoa(port.Port), port.Server,
				}
				updateCSV(record, deviationsFilePath)
			}
		}
	}

	totalPortsReviewed := len(ports)
	fmt.Printf("Review completed: total number of combinations reviewed was %v\n", totalPortsReviewed)

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	jiraClient, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		fmt.Println(err)
	}

	jql := "labels=" + strings.ToLower(month.String()) + "-" + strconv.Itoa(year)
	searchResults, resp, err := jiraClient.Issue.Search(jql, nil)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	if len(searchResults) != 0 {
		//fmt.Printf("%v", issue.Key)
		fmt.Printf("looks like there's already a ticket for this month\n")
		time.Sleep(10 * time.Second)
		fmt.Println("Removing the file locally")
		e := os.Remove(deviationsFilePath)
		if e != nil {
			log.Fatal(e)
		}
		return

	} else {
		fmt.Println("no existing tickets for this month, continuing the review")

		me, resp, err := jiraClient.User.GetSelf()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		if err != nil {
			fmt.Println(err)
		}
		if err != nil {
			fmt.Println(err)
		}

		i := jira.Issue{
			Fields: &jira.IssueFields{
				Assignee:    me,
				Description: "Test Issue",
				Type: jira.IssueType{
					Name: "Task",
				},
				Project: jira.Project{
					Key: "AT",
				},
				Summary: "PPS Review " + month.String() + " " + strconv.Itoa(year),
				Labels:  []string{"pps-review", strings.ToLower(month.String()) + "-" + strconv.Itoa(year)},
			},
		}
		newIssue, resp, err := jiraClient.Issue.Create(&i)
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%v has been created\n", newIssue.Key)

		//issue, _, _ := jiraClient.Issue.Get("AT-1", nil)

		//fmt.Printf("%s: %v\n", issue.Key, issue.Fields.Summary)

		f, err := os.OpenFile(deviationsFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		//csvReader := csv.NewReader(f)

		attachmentName := "pps-" + strings.ToLower(month.String()) + "-" + strconv.Itoa(year) + "-deviations.csv"
		_, resp, err = jiraClient.Issue.PostAttachment(newIssue.Key, f, attachmentName)
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("uploading the attachment")
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(10 * time.Second)
		fmt.Println("Removing the file locally")
		e := os.Remove(deviationsFilePath)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func getYaml() []PortsAndServersRuleset {
	fmt.Println("Getting the YAML file")

	var p []PortsAndServersRuleset

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	yamlFileName := os.Getenv("KEY_YML")
	fmt.Println(yamlFileName)
	yamlFile, err := ioutil.ReadFile(os.Getenv("KEY_YML"))
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
	csvFilePath := os.Args[1]
	fmt.Println(csvFilePath)

	portsFile, err := os.OpenFile(csvFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)

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

// Checks if a port from the csv is in the TCPOther / UDPOther allowlist in the yaml
func containsPorts(p []int, po int) bool {
	for _, v := range p {
		if v == po {
			return true
		}
	}
	return false
}
