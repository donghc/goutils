package pubsub

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wjiec/gdsn"
)

func TestDsn(t *testing.T) {
	d, err := gdsn.Parse(`kafka://user:pass@10.10.10.10:9092,20.20.20.20:9092,30.30.30.30:9092/test?q=123`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if d.Address() == "" {
		t.Fatalf("parse address: %v", d.Address())
	}
	if d.User.Username() != "user" {
		t.Fatalf("parse user: %v", d.User.Username())
	}
	pass, _ := d.User.Password()
	if pass != "pass" {
		t.Fatalf("parse pass: %v", pass)
	}
	if d.Path != "/test" {
		t.Fatalf("parse path: %v", d.Path)
	}
	if d.Query().Get("q") != "123" {
		t.Fatalf("parse q: %v", d.Query().Get("q"))
	}
	fmt.Println(d.Address())
	fmt.Println(strings.Split(d.Address(), ","))
}
