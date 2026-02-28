package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed fixes/*.md
var fixesFS embed.FS

func main() {
	low := flag.String("low", "", "low Hugo version, inclusive (e.g. v0.110.0)")
	high := flag.String("high", "", "high Hugo version, inclusive (e.g. v0.156.0)")
	flag.Parse()

	var lowVersion, highVersion [3]int
	var hasLow, hasHigh bool

	if *low != "" {
		v, err := parseVersion(*low)
		if err != nil {
			log.Fatalf("invalid low version %q: %v", *low, err)
		}
		lowVersion = v
		hasLow = true
	}

	if *high != "" {
		v, err := parseVersion(*high)
		if err != nil {
			log.Fatalf("invalid high version %q: %v", *high, err)
		}
		highVersion = v
		hasHigh = true
	}

	entries, err := fs.ReadDir(fixesFS, "fixes")
	if err != nil {
		log.Fatalf("reading fixes directory: %v", err)
	}

	var files []versionedFile
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		v, err := parseVersion(name)
		if err != nil {
			continue
		}
		if hasLow && compareSemver(v[0], v[1], v[2], lowVersion[0], lowVersion[1], lowVersion[2]) < 0 {
			continue
		}
		if hasHigh && compareSemver(v[0], v[1], v[2], highVersion[0], highVersion[1], highVersion[2]) > 0 {
			continue
		}
		files = append(files, versionedFile{v[0], v[1], v[2], e.Name()})
	}

	if len(files) == 0 {
		log.Fatal("no fix files found for the specified version range")
	}

	sort.Slice(files, func(i, j int) bool {
		return compareSemver(files[i].major, files[i].minor, files[i].patch,
			files[j].major, files[j].minor, files[j].patch) < 0
	})

	firstFileVersion := files[0].String()
	lastFileVersion := files[len(files)-1].String()

	header := fmt.Sprintf(`## Fixes from Hugo %s to Hugo %s

Note that these upgrades typically also requires upgrading to Hugo %s in e.g. netlify.toml.
	
`, firstFileVersion, lastFileVersion, lastFileVersion)
	os.Stdout.WriteString(header)

	for _, f := range files {
		data, err := fs.ReadFile(fixesFS, "fixes/"+f.name)
		if err != nil {
			log.Fatalf("reading %s: %v", f.name, err)
		}
		os.Stdout.Write(data)
		fmt.Println()
	}
}

func parseVersion(s string) ([3]int, error) {
	s = strings.TrimPrefix(s, "v")
	var major, minor, patch int
	n, err := fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
	if err != nil || n != 3 {
		return [3]int{}, fmt.Errorf("invalid semver: %s", s)
	}
	return [3]int{major, minor, patch}, nil
}

func compareSemver(aMaj, aMin, aPat, bMaj, bMin, bPat int) int {
	if aMaj != bMaj {
		return aMaj - bMaj
	}
	if aMin != bMin {
		return aMin - bMin
	}
	return aPat - bPat
}

type versionedFile struct {
	major, minor, patch int
	name                string
}

func (v versionedFile) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}
