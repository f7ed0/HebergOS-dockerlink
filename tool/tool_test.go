package tool

import "testing"

func TestVersionCheck(t *testing.T) {
	// Checking VeionSup
	want := VersionCheck("v0.0","v0.0",VersionSup)
	if !want {
		t.Fatalf("VersionCheck(\"v0.0\",\"v0.0\",VersionSup) should return true but returns %v",want)
	}
	want = VersionCheck("v0.1","v0.0",VersionSup)
	if !want {
		t.Fatalf("VersionCheck(\"v0.1\",\"v0.0\",VersionSup) should return true but returns %v",want)
	}
	want = VersionCheck("v0.0.1","v0.1",VersionSup)
	if want {
		t.Fatalf("VersionCheck(\"v0.0.1\",\"v0.1\",VersionSup) should return false but returns %v",want)
	}
	want = VersionCheck("v0.1.1","v0.1",VersionSup)
	if !want {
		t.Fatalf("VersionCheck(\"v0.1.1\",\"v0.1\",VersionSup) should return true but returns %v",want)
	}
	want = VersionCheck("d0.0.1","v0.1",VersionSup)
	if want {
		t.Fatalf("VersionCheck(\"d0.0.1\",\"v0.1\",VersionSup) should return false but returns %v",want)
	}
	want = VersionCheck("v0.0.1","d0.1",VersionSup)
	if want {
		t.Fatalf("VersionCheck(\"d0.0.1\",\"v0.1\",VersionSup) should return false but returns %v",want)
	}
}