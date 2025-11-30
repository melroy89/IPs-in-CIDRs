package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Matcher struct {
	CIDRFile string
	IPFile   string
	mu       sync.RWMutex
	cidrs    []*net.IPNet
	ips      []net.IP
}

func main() {
	m := &Matcher{
		CIDRFile: "cidrs.txt",
		IPFile:   "ips.txt",
	}
	log.Println("Loading CIDRs and IPs")
	err := m.init()
	if err != nil {
		log.Fatalf("failed to initialize matcher: %v", err)
	}
	log.Println("Starting to match IPs in CIDRs")
	m.matching()
}

func (m *Matcher) init() error {
	if err := m.loadCIDRs(); err != nil {
		log.Printf("failed to load cidrs at startup: %v", err)
		return err
	}

	if err := m.loadIPs(); err != nil {
		log.Printf("failed to load ips at startup: %v", err)
		return err
	}

	return nil
}

func (m *Matcher) loadCIDRs() error {
	path := m.CIDRFile
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var list []*net.IPNet
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// support both plain IPs and CIDRs
		if !strings.Contains(line, "/") {
			// treat single IP as /32 or /128
			if strings.Contains(line, ":") {
				// ipv6
				line = line + "/128"
			} else {
				// ipv4
				line = line + "/32"
			}
		}
		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			log.Printf("invalid cidr '%s' in cidrs: %v", scanner.Text(), err)
			continue
		}
		list = append(list, ipnet)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	m.cidrs = list
	m.mu.Unlock()
	return nil
}

func (m *Matcher) loadIPs() error {
	path := m.IPFile
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var list []net.IP
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ip := net.ParseIP(line)
		if ip == nil {
			log.Printf("invalid ip '%s'", scanner.Text())
			continue
		}
		list = append(list, ip)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	m.ips = list
	m.mu.Unlock()
	return nil
}

func (m *Matcher) matching() error {
	// Keep track of uniuqe CIDR (IPNet) matches
	matchesMap := make(map[string]*net.IPNet)
	for _, ip := range m.ips {
		if cidr, ok := m.ipMatchesWatchList(ip.String()); ok {
			fmt.Printf("IP: %s matches CIDR: %s\n", ip.String(), cidr.String())
			cidrStr := cidr.String()
			if _, exists := matchesMap[cidrStr]; !exists {
				matchesMap[cidrStr] = cidr
			}
		}
	}
	fmt.Println("\n\nCIDRs with matches:\n---------------")
	for _, cidr := range matchesMap {
		fmt.Printf("%s\n", cidr.String())
	}
	return nil
}

func (m *Matcher) ipMatchesWatchList(ipStr string) (*net.IPNet, bool) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, n := range m.cidrs {
		if n.Contains(ip) {
			return n, true
		}
	}
	return nil, false
}
