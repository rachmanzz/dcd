package parse

import (
	"os"
	"testing"
)

const testInput = `[section 0]
name=userinfo
var=info
keys=username, date_field
formats=[date_field:dd-MM-yyyy]

--- BODY ---
<p>username: {{info.username}}</p>
<p>date: {{info.date_field}}</p>

[section 1]
name=address
var=addr
keys=street, city

--- BODY ---
<p>{{addr.street}}, {{addr.city}}</p>
`

func TestParse(t *testing.T) {
	f, err := os.CreateTemp("", "*.dcd")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(testInput); err != nil {
		t.Fatal(err)
	}
	f.Close()

	doc, err := Parse(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(doc.Sections))
	}

	s0 := doc.Sections[0]
	if s0.N != 0 {
		t.Errorf("section 0 N = %d, want 0", s0.N)
	}
	if s0.Props["name"] != "userinfo" {
		t.Errorf("section 0 name = %q, want userinfo", s0.Props["name"])
	}
	if s0.Props["var"] != "info" {
		t.Errorf("section 0 var = %q, want info", s0.Props["var"])
	}
	if s0.Body == "" {
		t.Error("section 0 body is empty")
	}

	s1 := doc.Sections[1]
	if s1.N != 1 {
		t.Errorf("section 1 N = %d, want 1", s1.N)
	}
	if s1.Props["name"] != "address" {
		t.Errorf("section 1 name = %q, want address", s1.Props["name"])
	}
}
