package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

var filename, downloadURL, logfileName string
var log *os.File

const (
	defaultPort = "21"
	defaultUser = "anonymous"
	defaultPass = "anonymous"
	logfileExtn = "-url-import-error.txt"
)

func main() {
	logfileName = strings.Join(strings.Split(time.Now().Format(time.UnixDate), " "), "-") + logfileExtn
	file, err := os.Create(logfileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create error file. ", err)
		os.Exit(1)
		return
	}
	log = file
	os.Stderr = file
	defer file.Close()

	parseArgs()
	scheme, host, port, user, pass, path, exitcode := ParseURL()
	if exitcode != 0 {
		cleanup()
		defer os.Exit(exitcode)
		return
	}

	switch strings.ToLower(scheme) {
	case "http", "https":
		exitcode = DownloadFromURL()
	case "ftp":
		exitcode = DownloadFromFtp(host, port, user, pass, path)
	default:
		logMessage("Invalid Scheme.")
	}

	cleanup()
	defer os.Exit(exitcode)
}

//cleanup clean up 0 byte log file or partially downloaded output file
func cleanup() {
	// if error file is 0 bytes, then remove it
	if fi, err := os.Stat(logfileName); !os.IsNotExist(err) {
		if fi.Size() == 0 { // no errors logged
			os.Remove(logfileName)
		} else {
			os.Remove(filename) //there were some errors logged. Remove partial downloaded file.
		}
	}
}

//logMessage to a log file
func logMessage(message string) {
	log.WriteString(message) // nolint:errcheck
	log.WriteString("\n")    // nolint:errcheck
}

//parseArgs parse command line args
func parseArgs() int {
	flag.StringVar(&filename, "filename", "", "file name to use when saving the imported file")
	flag.StringVar(&downloadURL, "url", "", "Url to import the file from")
	flag.Parse()
	if len(filename) == 0 || len(downloadURL) == 0 {
		flag.PrintDefaults()
		logMessage("Invalid command line arguments!")
		return 1
	}
	return 0
}

//ParseURL parses given url to extract host, post, user, pass and path.
func ParseURL() (string, string, string, string, string, string, int) {
	var username string
	var password string

	u, err := url.Parse(downloadURL)
	if err != nil {
		logMessage("Unable to parse url. " + err.Error())
		return "", "", "", "", "", "", 1
	}
	host, port, _ := net.SplitHostPort(u.Host)

	if len(host) == 0 {
		host = u.Host
	}

	if u.User != nil {
		username = u.User.Username()
		pass, _ := u.User.Password()
		password = pass
	} else {
		username = defaultUser
		password = defaultPass
	}

	if len(port) == 0 {
		port = defaultPort
	}

	return u.Scheme, host, port, username, password, u.Path, 0
}

//DownloadFromURL downloads file from given http(s) url
func DownloadFromURL() int {
	output, err := os.Create(filename)
	if err != nil {
		logMessage("Unable to create output file. " + err.Error())
		return 1
	}
	defer output.Close()
	response, err := http.Get(downloadURL)
	if err != nil {
		logMessage("Unable to download file. " + err.Error())
		return 1
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		logMessage("Request failed. Status code:" + response.Status)
		return 1
	}

	n, err := io.Copy(output, response.Body)
	if err != nil {
		logMessage("Unable to copy contents. " + err.Error())
		return 1
	}
	fmt.Println(n, " bytes downloaded.")
	return 0
}

//DownloadFromFtp download a file from a given ftp url, port, user, pass and path
func DownloadFromFtp(host string, port string, user string, pass string, path string) int {
	fmt.Println(host, port, user, pass, path)
	s := []string{host, port}
	hostPort := strings.Join(s, ":")
	conn, err := ftp.DialTimeout(hostPort, 5*time.Second)

	if err != nil {
		logMessage("Unable to connect to ftp server. " + err.Error())
		return 1
	}
	defer conn.Quit() // nolint:errcheck

	err = conn.Login(user, pass)
	if err != nil {
		logMessage("Unable to login to FTP server. " + err.Error())
		return 1
	}

	fileSize, err := conn.FileSize(path)
	if err != nil {
		fmt.Println("Couldn't retrieve file size from FTP server, continuing anyway")
	} else {
		fmt.Printf("File size: %d\n", fileSize)
	}

	response, err := conn.Retr(path)
	if err != nil {
		logMessage("Unable to retrieve file from FTP server. " + err.Error())
		return 1
	}
	defer response.Close()

	output, err := os.Create(filename)
	if err != nil {
		logMessage("Unable to create output file. " + err.Error())
		return 1
	}
	defer output.Close()

	n, err := io.Copy(output, response)
	if err != nil {
		logMessage("Unable to copy contents. " + err.Error())
		return 1
	}
	fmt.Println(n, " bytes downloaded.")
	return 0
}
