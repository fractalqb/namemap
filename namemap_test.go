package namemap

import (
	"bytes"
	"testing"
)

func TestNameMap_Load(t *testing.T) {
	buf := bytes.NewBufferString(`[\id lang: lang:de]
	(1 foo baz)
	(2 bar quux)`)
	nm := NameMap{}
	err := nm.Load(buf)
	if err != nil {
		t.Fatal(err)
	}
	if nm.stdDom != 0 {
		t.Errorf("wrong standard domain: %d", nm.stdDom)
	}
	if nm.DomainIdx("id") != 0 {
		t.Errorf("wrong domain index for 'id': %d", nm.DomainIdx("id"))
	}
	if nm.DomainIdx("lang:") != 1 {
		t.Errorf("wrong domain index for 'lang:': %d", nm.DomainIdx("lang:"))
	}
	if nm.DomainIdx("lang:de") != 2 {
		t.Errorf("wrong domain index for 'lang:de': %d", nm.DomainIdx("lang:de"))
	}
	if mapped, ok := nm.Map(0, "2", 2); !ok {
		t.Error("mapping failed: nothing for id=2")
	} else if mapped != "quux" {
		t.Errorf("mapping failed: got %s, expect %s", mapped, "quux")
	}
}
