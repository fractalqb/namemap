package namemap

import (
	"bytes"
	"fmt"
	"strings"
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
	if mapped, toDom := nm.Map(0, "2", 2); toDom < 0 {
		t.Error("mapping failed: nothing for id=2")
	} else if mapped != "quux" {
		t.Errorf("mapping failed: got %s, expect %s", mapped, "quux")
	}
}

func ExampleNameMap_Load() {
	nmDef := strings.NewReader(
		`[\input output l10n:EN l10n:DE]
          (note  rem    remark  \undef)
          (warn  warnig warning Warnung)`)
	nm := NameMap{}
	if err := nm.Load(nmDef); err != nil {
		panic(err)
	}
	mapped, inDom := nm.MapNm("input", "warn", "l10n:DE", "l10n:EN")
	fmt.Printf("input name 'warn' maps to '%s' in domain %s\n", mapped, nm.DomainName(inDom))
	mapped, inDom = nm.MapNm("output", "rem", "l10n:DE", "l10n:EN")
	fmt.Printf("output name 'rem' maps to '%s' in domain %s\n", mapped, nm.DomainName(inDom))
	// Output:
	// input name 'warn' maps to 'Warnung' in domain l10n:DE
	// output name 'rem' maps to 'remark' in domain l10n:EN
}
