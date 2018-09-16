package namemap

import (
	"bytes"
	"fmt"
	"os"
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
	if nm.StdDomain != 0 {
		t.Errorf("wrong standard domain: %d", nm.StdDomain)
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

func TestNameMap_Set(t *testing.T) {
	nm := NewNameMap()
	nm.Def(map[string]string{
		"input":   "note",
		"output":  "rem",
		"l10n:EN": "remark",
	})
	nm.Def(map[string]string{
		"input":   "warn",
		"output":  "warnig",
		"l10n:EN": "warning",
		"l10n:DE": "Warnung",
	})
	nm.SetStdDomain("input")
	domIn := nm.DomainIdx("input")
	domDe := nm.DomainIdx("l10n:DE")
	mapped, mdom := nm.Map(domIn, "note", domDe)
	if mdom >= 0 {
		t.Errorf("expected empty name, got '%s' from domain '%s' (%d)",
			mapped,
			nm.DomainName(mdom),
			mdom)
	}
	if mapped != "note" {
		t.Errorf("unmappable shoud return input name, got '%s'", mapped)
	}
	nm.Set(domIn, "note", domDe, "Bemerkung")
	mapped, mdom = nm.Map(domIn, "note", domDe)
	if mdom != domDe {
		t.Errorf("result from wrong domain '%s' (%d), expected '%s' (%d)",
			nm.DomainName(mdom), mdom,
			nm.DomainName(domDe), domDe)
	}
	if mapped != "Bemerkung" {
		t.Errorf("new mapping is wrong '%s', expected 'Bemerkung'", mapped)
	}
	mapped, mdom = nm.Map(domDe, "Bemerkung", domIn)
	if mdom != domIn || mapped != "note" {
		t.Errorf("new mapping does not map reverse: %s (%d)", mapped, mdom)
	}
	mapped, mdom = nm.MapNm("l10n:DE", "Bemerkung", "output")
	if mdom != nm.DomainIdx("output") || mapped != "rem" {
		t.Errorf("new mapping does not map cross: %s (%d)", mapped, mdom)
	}
}

// TODO cannot rely on column order => might fail occasionally => rewrite as test
func ExampleNameMap_Def() {
	nm := NewNameMap()
	nm.Def(map[string]string{
		"input":   "note",
		"output":  "rem",
		"l10n:EN": "remark",
	})
	nm.Def(map[string]string{
		"input":   "warn",
		"output":  "warnig",
		"l10n:EN": "warning",
		"l10n:DE": "Warnung",
	})
	nm.SetStdDomain("input")
	nm.Save(os.Stdout, "null")
	// Unordered output:
	// [\input output l10n:EN l10n:DE]
	// (note rem remark \null)
	// (warn warnig warning Warnung)
}

func ExampleNameMap() {
	nm := NewNameMap("key", "local")
	nm.Set(0, "akey", 1, "aloc")
	err := nm.Save(os.Stdout, "null")
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// [key local]
	// (akey aloc)
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

func ExampleIgnDom() {
	nmDef := strings.NewReader(
		`[\input output l10n:EN l10n:DE]
          (note  rem    remark  \undef)
          (warn  warnig warning Warnung)`)
	nm := NameMap{}
	if err := nm.Load(nmDef); err != nil {
		panic(err)
	}
	fmt.Println(IgnDom(nm.MapNm("input", "note", "l10n:DE", "l10n:EN")))
	// Output:
	// remark
}
