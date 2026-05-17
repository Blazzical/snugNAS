package mdns

import "testing"

func TestCanonicalHostname(t *testing.T) {
	cases := map[string]string{
		"snugnas":         "snugnas",
		"snugnas.local":   "snugnas",
		"snugnas.local.":  "snugnas",
		"foo.bar.local":   "foo.bar",
		"foo.bar.local.":  "foo.bar",
		"":                "",
		".local":          "",
		".local.":         "",
	}
	for in, want := range cases {
		got := canonicalHostname(in)
		if got != want {
			t.Errorf("canonicalHostname(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLanIPv4Filters(t *testing.T) {
	ips, err := lanIPv4()
	if err != nil {
		t.Fatalf("lanIPv4: %v", err)
	}
	for _, ip := range ips {
		if ip.IsLoopback() {
			t.Errorf("got loopback IP: %v", ip)
		}
		if ip.IsLinkLocalUnicast() {
			t.Errorf("got link-local IP: %v", ip)
		}
		if ip.To4() == nil {
			t.Errorf("got non-IPv4: %v", ip)
		}
	}
}

func TestPublishRejectsEmptyHostname(t *testing.T) {
	_, err := Publish("", 8080)
	if err == nil {
		t.Error("Publish(\"\") returned nil err")
	}
}
