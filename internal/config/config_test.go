package config

import (
	"os"
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Hostname != "snugnas.local" {
		t.Errorf("Hostname = %q, want snugnas.local", c.Hostname)
	}
	if c.DashboardPort != 7777 {
		t.Errorf("DashboardPort = %d, want 7777", c.DashboardPort)
	}
	if len(c.PostgresPassword) < 32 {
		t.Errorf("PostgresPassword too short (%d) — randomHex(24) should give 48 chars", len(c.PostgresPassword))
	}
	// Twice should give different passwords
	c2 := Default()
	if c.PostgresPassword == c2.PostgresPassword {
		t.Error("two Default() calls returned the same password — generator not random")
	}
}

func TestPosixStorageRoot(t *testing.T) {
	cases := map[string]string{
		`E:\snugNAS-data`:           "E:/snugNAS-data",
		`D:\foo\bar`:                "D:/foo/bar",
		`/already/posix`:            "/already/posix",
		``:                          "",
		`C:\Users\name with space\d`: "C:/Users/name with space/d",
	}
	for in, want := range cases {
		got := (&Config{StorageRoot: in}).PosixStorageRoot()
		if got != want {
			t.Errorf("PosixStorageRoot(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APPDATA", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("HOME", tmp)

	in := &Config{
		StorageRoot:      `D:\test`,
		Hostname:         "test.local",
		DashboardPort:    1234,
		AdminEmail:       "x@y.z",
		PostgresPassword: "secret-pw",
	}
	if err := in.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, in) {
		t.Errorf("round-trip mismatch:\n got: %+v\nwant: %+v", got, in)
	}
}

func TestLoadStripsUTF8BOM(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APPDATA", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("HOME", tmp)

	// Save a default first so the directory + file exist, then rewrite with a BOM.
	if err := Default().Save(); err != nil {
		t.Fatalf("Save default: %v", err)
	}
	p, err := Path()
	if err != nil {
		t.Fatalf("Path: %v", err)
	}

	body := []byte("storage_root = 'D:\\test'\nhostname = 'bom.local'\ndashboard_port = 9999\nadmin_email = ''\npostgres_password = 'x'\n")
	withBOM := append([]byte{0xEF, 0xBB, 0xBF}, body...)
	if err := os.WriteFile(p, withBOM, 0o600); err != nil {
		t.Fatalf("write BOM: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Hostname != "bom.local" {
		t.Errorf("Hostname = %q, want bom.local — BOM not stripped?", got.Hostname)
	}
}

func TestIsNotConfigured(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APPDATA", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("HOME", tmp)

	_, err := Load()
	if err == nil {
		t.Fatal("Load on empty dir returned nil error")
	}
	if !IsNotConfigured(err) {
		t.Errorf("IsNotConfigured = false for missing-file error: %v", err)
	}
}
