package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	Buildstamp string
	Githash    string
)

var IP_PROVIDER = "https://echo.tinyandbeautiful.com/ip"

type Result struct {
	ID      string `json:"id"`
	ZoneID  string `json:"zone_id"`
	Content string `json:"content"`
	Name    string `json:"name"`
}

type Response struct {
	Result  json.RawMessage `json:"result"`
	Success bool            `json:"success"`
}

func parseResultList(raw json.RawMessage) ([]Result, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	if raw[0] == '[' {
		var records []Result
		if err := json.Unmarshal(raw, &records); err != nil {
			return nil, err
		}
		return records, nil
	}

	var record Result
	if err := json.Unmarshal(raw, &record); err != nil {
		return nil, err
	}
	return []Result{record}, nil
}

type DNSRecord struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	Type    string   `json:"type"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
	TTL     int      `json:"ttl"`
}

func getOwnIPv4() (string, error) {

	c := http.Client{Timeout: 10 * time.Second}

	resp, err := c.Get(IP_PROVIDER)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("empty response from IP provider")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IP provider returned status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", errors.New("empty IP returned by provider")
	}
	return ip, nil
}

func findDNSRecordByName(records []Result, name string) (Result, bool) {
	if len(records) == 0 {
		return Result{}, false
	}

	for _, value := range records {
		if strings.EqualFold(value.Name, name) {
			return value, true
		}
	}

	return Result{}, false
}

func getDomainIPv4() (string, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=A", ZONEID), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CF_TOKEN))
	c := http.Client{Timeout: 5 * time.Second}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("empty response from Cloudflare DNS API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Cloudflare DNS API returned status %s", resp.Status)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if !response.Success {
		return "", errors.New("get dns record error")
	}

	results, err := parseResultList(response.Result)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", errors.New("get dns record not equal")
	}

	if DOMAIN == "@" {
		DNSID = results[0].ID
		return results[0].Content, nil
	}

	if record, ok := findDNSRecordByName(results, DOMAIN); ok {
		DNSID = record.ID
		return record.Content, nil
	}

	return "", errors.New("get dns record not equal")
}

func putNewIP(ip string) error {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(DNSRecord{
		Content: ip,
		Name:    DOMAIN,
		Proxied: false,
		Type:    "A",
		TTL:     60,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", ZONEID, DNSID),
		&buf)

	if err != nil {
		log.Error("Error creating request:", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CF_TOKEN))
	c := http.Client{Timeout: 5 * time.Second}

	resp, err := c.Do(req)
	if err != nil {
		log.Errorf("res err %s", err)
		return err
	}
	if resp == nil {
		return errors.New("empty response from Cloudflare DNS API")
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK && response.Success {
		log.Debug("update ok")
		return nil
	}
	return fmt.Errorf("failed with HTTP status code %d", resp.StatusCode)
}

func run() {
	log.Debug("get own ip -")

	ownIP, err := getOwnIPv4()
	if err != nil {
		log.Errorf("get own ip err, %s", err)
		return
	}

	log.Debugf("get own ip: %s", ownIP)

	log.Debug("get domain ip -")

	domainIP, err := getDomainIPv4()
	if err != nil {
		log.Errorf("get domain ip err, %s", err)
		return
	}

	log.Debugf("get domain ip: %s", domainIP)

	if domainIP != ownIP {
		if err := putNewIP(ownIP); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Infof("same ip, ignore")
	}
}

// globals
var CF_TOKEN = ""
var ZONEID = ""

var DOMAIN = "@"
var DNSID = ""

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.SetLevel(log.DebugLevel)

	file, err := os.OpenFile("./dnslog", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}

	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		log.SetOutput(fileAndStdoutWriter)
	} else {
		log.Info("failed to log to file.")
	}

	// required flags
	keyPtr := flag.String("token", "", "cf Token")
	domainPtr := flag.String("domain", "@", "Your top level domain (e.g., example.com)")
	zoneidPtr := flag.String("zoneid", "", "Zone id")

	var flagversion bool
	flag.BoolVar(&flagversion, "v", false, "version")

	flag.Parse()

	if flagversion {
		fmt.Printf("Git Commit Hash: %s\n", Githash)
		fmt.Printf("Build Time : %s\n", Buildstamp)
		return
	}

	CF_TOKEN = *keyPtr
	ZONEID = *zoneidPtr
	DOMAIN = *domainPtr

	if CF_TOKEN == "" {
		log.Fatalf("You need to provide your cloudFlare TOKEN")
		return
	}

	if ZONEID == "" {
		log.Fatalf("You need to provide your cloudFlare Zone id")
		return
	}

	run()
}
