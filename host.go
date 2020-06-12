package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"os"
	"strconv"
	"time"
)

type Host struct {
	IP       string
	Port     int
	Username string
	Password string
}

func ReadCsv(fileName string) ([]Host, error) {
	var Hosts []Host
	csvFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("Failed to open csv file")
	}

	r := csv.NewReader(csvFile)

	record, err := r.ReadAll()

	if err != nil {
		return nil, errors.New("Failed to parse csv file")
	}
	for _, v := range record {
		port, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, errors.New("port must be number")
		}
		host := Host{
			IP:       v[0],
			Port:     port,
			Username: v[2],
			Password: v[3],
		}
		Hosts = append(Hosts, host)
	}

	return Hosts, nil
}

func WriteFile(content, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	if _, err := w.WriteString(content); err != nil {
		panic(err.Error())
	}

}

func (h Host) String() string {
	return fmt.Sprintf("IP: %s,port: %d, username: %s\n", h.IP, h.Port, h.Username)
}

func (h Host) Exec(cmd string) string {
	var (
		config  *ssh.ClientConfig
		client  *ssh.Client
		session *ssh.Session
		res     bytes.Buffer
		err     error
	)
	config = &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 3,
	}
	client, err = ssh.Dial("tcp", h.IP+":"+strconv.Itoa(h.Port), config)

	if err != nil {
			msg := fmt.Sprintf("Failed, ip %s, msg %s\n", h.IP, err.Error())
			return msg
	}

	session, err = client.NewSession()
	if err != nil {
		msg := fmt.Sprintf("Failed, ip %s, msg %s\n", h.IP, err.Error())
		return msg
	}

	defer session.Close()
	session.Stdout = &res

	if err := session.Run(cmd); err != nil {
		return fmt.Sprintf("Failed, ip %s, msg %s, err %s\n", h.IP, res.String(), err.Error())
	}

	return fmt.Sprintf("OK, ip %s, msg %s\n", h.IP, res.String())
}

func (h Host) Copy(src, dest, mask string) string {
	var (
		config     *ssh.ClientConfig
		client     scp.Client
		hostString string
		file       *os.File
		err        error
	)

	config = &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 3,
	}

	hostString = fmt.Sprintf("%s:%d", h.IP, h.Port)
	client = scp.NewClient(hostString, config)

	err = client.Connect()
	if err != nil {
		return fmt.Sprintf("Failed, ip %s, msg %s\n", h.IP, err.Error())
	}

	file, err = os.Open(src)
	if err != nil {
		return fmt.Sprintf("Failed, ip %s, msg %s\n", h.IP, err.Error())
	}

	defer client.Close()

	err = client.CopyFile(file, dest, mask)
	if err != nil {
		return fmt.Sprintf("Failed, ip %s, msg %s\n", h.IP, err.Error())
	}
	return fmt.Sprintf("OK, ip %s", h.IP)
}
