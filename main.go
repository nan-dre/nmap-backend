package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

type Data struct {
	// XMLName          xml.Name `xml:"nmaprun"`
	// Text             string   `xml:",chardata"`
	// Scanner          string   `xml:"scanner,attr"`
	Args string `xml:"args,attr"`
	// Start            string   `xml:"start,attr"`
	// Startstr         string   `xml:"startstr,attr"`
	// Version          string   `xml:"version,attr"`
	// Xmloutputversion string   `xml:"xmloutputversion,attr"`
	// Scaninfo         struct {
	// 	Text        string `xml:",chardata"`
	// 	Type        string `xml:"type,attr"`
	// 	Protocol    string `xml:"protocol,attr"`
	// 	Numservices string `xml:"numservices,attr"`
	// 	Services    string `xml:"services,attr"`
	// } `xml:"scaninfo"`
	// Verbose struct {
	// 	Text  string `xml:",chardata"`
	// 	Level string `xml:"level,attr"`
	// } `xml:"verbose"`
	// Debugging struct {
	// 	Text  string `xml:",chardata"`
	// 	Level string `xml:"level,attr"`
	// } `xml:"debugging"`
	// Hosthint struct {
	// Text   string `xml:",chardata"`
	// Status struct {
	// 	Text      string `xml:",chardata"`
	// 	State     string `xml:"state,attr"`
	// 	Reason    string `xml:"reason,attr"`
	// 	ReasonTtl string `xml:"reason_ttl,attr"`
	// } `xml:"status"`
	// Address struct {
	// 	Text     string `xml:",chardata"`
	// 	Addr     string `xml:"addr,attr"`
	// 	Addrtype string `xml:"addrtype,attr"`
	// } `xml:"address"`
	// Hostnames struct {
	// Text     string `xml:",chardata"`
	// 		Hostname struct {
	// 			Text string `xml:",chardata"`
	// 			Name string `xml:"name,attr"`
	// 			Type string `xml:"type,attr"`
	// 		} `xml:"hostname"`
	// 	} `xml:"hostnames"`
	// } `xml:"hosthint"`
	Host struct {
		// Text      string `xml:",chardata"`
		// Starttime string `xml:"starttime,attr"`
		// Endtime   string `xml:"endtime,attr"`
		// Status    struct {
		// 	Text      string `xml:",chardata"`
		// 	State     string `xml:"state,attr"`
		// 	Reason    string `xml:"reason,attr"`
		// 	ReasonTtl string `xml:"reason_ttl,attr"`
		// } `xml:"status"`
		Address struct {
			Text     string `xml:",chardata"`
			Addr     string `xml:"addr,attr"`
			Addrtype string `xml:"addrtype,attr"`
		} `xml:"address"`
		Hostnames struct {
			Text     string `xml:",chardata"`
			Hostname []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Type string `xml:"type,attr"`
			} `xml:"hostname"`
		} `xml:"hostnames"`
		Ports struct {
			// Text       string `xml:",chardata"`
			// Extraports struct {
			// 	Text         string `xml:",chardata"`
			// 	State        string `xml:"state,attr"`
			// 	Count        string `xml:"count,attr"`
			// 	// Extrareasons struct {
			// 	// 	Text   string `xml:",chardata"`
			// 	// 	Reason string `xml:"reason,attr"`
			// 	// 	Count  string `xml:"count,attr"`
			// 	// 	Proto  string `xml:"proto,attr"`
			// 	// 	Ports  string `xml:"ports,attr"`
			// 	// } `xml:"extrareasons"`
			// } `xml:"extraports"`
			Port []struct {
				Text     string `xml:",chardata"`
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					Text      string `xml:",chardata"`
					State     string `xml:"state,attr"`
					Reason    string `xml:"reason,attr"`
					ReasonTtl string `xml:"reason_ttl,attr"`
				} `xml:"state"`
				Service struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					Method string `xml:"method,attr"`
					Conf   string `xml:"conf,attr"`
				} `xml:"service"`
			} `xml:"port"`
		} `xml:"ports"`
		// Times struct {
		// 	Text   string `xml:",chardata"`
		// 	Srtt   string `xml:"srtt,attr"`
		// 	Rttvar string `xml:"rttvar,attr"`
		// 	To     string `xml:"to,attr"`
		// } `xml:"times"`
	} `xml:"host"`
	Runstats struct {
		Text     string `xml:",chardata"`
		Finished struct {
			Text    string `xml:",chardata"`
			Time    string `xml:"time,attr"`
			Timestr string `xml:"timestr,attr"`
			Summary string `xml:"summary,attr"`
			Elapsed string `xml:"elapsed,attr"`
			Exit    string `xml:"exit,attr"`
		} `xml:"finished"`
		Hosts struct {
			Text  string `xml:",chardata"`
			Up    string `xml:"up,attr"`
			Down  string `xml:"down,attr"`
			Total string `xml:"total,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

type nmapArgs struct {
	host      string `form:"host" binding:"required"`
	startPort string `form:"start-port" binding:"required"`
	endPort   string `form:"end-port" binding:"required"`
}

func setup() {
	err := os.MkdirAll("./output", os.ModePerm)
	if err != nil {
		fmt.Print("Could not create directory ./output")
		return
	}
}

func scan(path string, args nmapArgs) string {
	arg0 := args.startPort + "-" + args.endPort
	cmd := exec.Command("nmap", "-p", arg0, "-oX", path, args.host)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalln("Could not scan: ", err)
	}

	return string(stdout)
}

func getData(args nmapArgs) Data {
	currentTime := time.Now()
	path := "./output/" + currentTime.Format("2006-01-02-15-04-05") + ".xml"
	scan(path, args)
	// path := "/home/andy/projects/go/nmap-backend/output/2021-11-20-10-58-40.xml"
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatalln("Could not open scan file")
	}
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var data Data
	xml.Unmarshal(byteValue, &data)
	return data
}

func main() {
	setup()
	router := gin.Default()
	router.LoadHTMLFiles("views/main.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})
	router.POST("/scan", func(c *gin.Context) {
		var args nmapArgs
		args.host = c.PostForm("host")
		args.startPort = c.PostForm("start-port")
		args.endPort = c.PostForm("end-port")
		data := getData(args)
		c.IndentedJSON(http.StatusOK, data)
	})

	log.Fatal(router.Run("localhost:8080"))
}
