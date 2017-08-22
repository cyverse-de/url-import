package main

import (
	"os"
	"testing"
)

func TestParseURL(t *testing.T) {
	filename = "foo.txt"
	downloadURL = "ftp://foo.com/home/bar/buzz.txt"
	scheme, host, port, user, pass, path, _ := ParseURL()
	if scheme != "ftp" {
		t.Errorf("Url Parse Error got %s, want %s", scheme, "ftp")
	}

	if host != "foo.com" {
		t.Errorf("Url Parse Error got %s, want %s", host, "foo.com")
	}

	if port != "21" {
		t.Errorf("Url Parse Error got %s, want %s", port, "21")
	}

	if user != "anonymous" {
		t.Errorf("Url Parse Error got %s, want %s", user, "anonymous")
	}

	if pass != "anonymous" {
		t.Errorf("Url Parse Error got %s, want %s", pass, "anonymous")
	}

	if path != "/home/bar/buzz.txt" {
		t.Errorf("Url Parse Error got %s, want %s", path, "/home/bar/buzz.txt")
	}

}

func TestDownloadFromURL(t *testing.T) {
	filename = "test_file.txt"
	downloadURL = "http://qa-4.cyverse.org/files/test_file.txt"

	DownloadFromURL()

	if _, err := os.Stat("test_file.txt"); os.IsNotExist(err) {
		t.Errorf("File not downloaded!")
	}

}

func TestDownloadFromFtp(t *testing.T) {
	filename = "readme.txt"
	downloadURL = "ftp://demo:password@test.rebex.net:21/pub/example/readme.txt"

	_, host, port, user, pass, path, _ := ParseURL()

	DownloadFromFtp(host, port, user, pass, path)

	if _, err := os.Stat("readme.txt"); os.IsNotExist(err) {
		t.Errorf("File not downloaded!")
	}

}
