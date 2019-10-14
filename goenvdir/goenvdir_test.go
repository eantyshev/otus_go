package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var ExternalEnvs = []string{
	"EXTERNAL_VAR1=external value 1",
	"EXTERNAL_VAR2=external value 2",
	"EXCLUDED_VAR=excluded value",
}
var InferredEnvs = []string{
	"EXTERNAL_VAR2=modified value 2",
	"EXCLUDED_VAR=",
	"NEW_VAR=new value",
	"MULTILINE_VAR=first line   \t   \nsecond line",
	"FANCY_VAR=first\x00second    \t  ",
}
var tempDir string

func setupTempDir() {
	var err error
	if tempDir, err = ioutil.TempDir("", "goenvdir"); err != nil {
		log.Fatal(err)
	}
	log.Println("temp dir: ", tempDir)
	for _, env := range InferredEnvs {
		nameVal := strings.SplitN(env, "=", 2)
		log.Println(nameVal)
		fPath := filepath.Join(tempDir, nameVal[0])
		if err := ioutil.WriteFile(fPath, []byte(nameVal[1]), 0666); err != nil {
			log.Fatal(err)
		}
		log.Println("create file: ", fPath)
	}
}

func teardownTempDir() {
	os.RemoveAll(tempDir)
}

func TestMain(m *testing.M) {
	setupTempDir()
	code := m.Run()
	teardownTempDir()
	os.Exit(code)
}

func TestIntegration(t *testing.T) {
	cmd := exec.Command("./goenvdir", tempDir, "bash", "-c", "env")
	cmd.Env = append(os.Environ(), ExternalEnvs...)
	out, err := cmd.CombinedOutput()
	t.Log(string(out))
	if err != nil {
		t.Fatal(err)
	}
	var extVisible, extOverridden, newInferred, newMultiline bool
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "=") {
			continue
		}
		nameVal := strings.SplitN(line, "=", 2)
		name := nameVal[0]
		value := nameVal[1]
		switch name {
		case "EXTERNAL_VAR1":
			assert.Equal(t, value, "external value 1")
			extVisible = true
		case "EXTERNAL_VAR2":
			assert.Equal(t, value, "modified value 2")
			extOverridden = true
		case "EXCLUDED_VAR":
			t.Error("EXCLUDED_VAR should be masked")
		case "NEW_VAR":
			assert.Equal(t, value, "new value")
			newInferred = true
		case "MULTILINE_VAR":
			assert.Equal(t, value, "first line")
			newMultiline = true
		}
	}
	assert.True(t, extVisible)
	assert.True(t, extOverridden)
	assert.True(t, newInferred)
	assert.True(t, newMultiline)
}

func TestGetFileEnv(t *testing.T) {
	r := strings.NewReader("fancy\x00var   \t  \nsecond line")
	if result, err := GetFileEnv(r); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "fancy\nvar", result)
	}
}
