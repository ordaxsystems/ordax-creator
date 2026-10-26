package main

import (
 "encoding/json"
 "os"
 "path/filepath"
 "testing"
)

func TestCreatorProductVersion(t *testing.T) {
 for _, s := range []string{"0.6.1-dev","1.0.0","1.0.0-rc.2"} {
  if !validCreatorVersion(s) { t.Fatalf("valid version rejected: %q",s) }
 }
 for _, s := range []string{"","owner-2bad2d125d7e","../1.0.0","1.0","1.2.3\n","V1.0.0"} {
  if validCreatorVersion(s) { t.Fatalf("invalid version accepted: %q",s) }
 }
}

func TestCreatorProductVersionSingleSource(t *testing.T) {
 content, err := os.ReadFile(filepath.Join("..","..","version.json"))
 if err != nil { t.Fatal(err) }
 var data struct {
  Schema string `json:"$schema"`
  Version string `json:"version"`
 }
 if err := json.Unmarshal(content, &data); err != nil { t.Fatal(err) }
 if data.Schema != "ordax.creator-product-version/1" || !validCreatorVersion(data.Version) {
  t.Fatalf("invalid canonical Creator version metadata: schema=%q version=%q", data.Schema, data.Version)
 }
 var fields map[string]any
 if err := json.Unmarshal(content, &fields); err != nil { t.Fatal(err) }
 if len(fields) != 2 { t.Fatalf("Creator version source must have exactly schema and version, got %d keys", len(fields)) }
}
