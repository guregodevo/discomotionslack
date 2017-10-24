package misc

import (
	"net/http"
	"testing"
)

func TestGetIPAdress(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.Header.Set("X-Real-Ip", "10.255.255.22, 192.168.1.22")
	r.RemoteAddr = "23.23.23.23"

	var result string

	result = GetIPAdress(r)
	if result != "23.23.23.23" {
		t.Error("IP Address didn't match: ", result)
	}

	r.Header.Set("X-Forwarded-For", "10.255.255.22, 192.168.1.22")
	r.RemoteAddr = "23.23.23.23:80"

	result = GetIPAdress(r)
	if result != "23.23.23.23" {
		t.Error("IP Address didn't match: ", result)
	}

	r.Header.Set("X-Forwarded-For", "10.255.255.22, 72.72.72.72")
	r.RemoteAddr = "23.23.23.23:80"

	result = GetIPAdress(r)
	if result != "72.72.72.72" {
		t.Error("IP Address didn't match: ", result)
	}
}
