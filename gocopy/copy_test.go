package main

import (
    "bytes"
    "testing"
    "io/ioutil"
    "os"
)

func CreateTmpFile(t *testing.T, content []byte) string {
    tmpfile, err := ioutil.TempFile("", "gocopy*")
    if err != nil {
        t.Fatal(err)
    }
    if len(content) > 0 {
        _, err = tmpfile.Write(content)
        if err != nil {
            t.Fatal(err)
        }
    }
    tmpfile.Close()
    return tmpfile.Name()
}

func Scenario(t *testing.T, content []byte, limit, offset int64, result []byte, fails bool) {
    var resultTest []byte
    from := CreateTmpFile(t, content)
    defer os.Remove(from)
    to := CreateTmpFile(t, []byte{})
    defer os.Remove(to)
    err := CopyData(from, to, limit, offset)
    if err != nil {
        if fails {
            t.Log(err)
            return
        }
        t.Fatal(err)
    }
    if resultTest, err = ioutil.ReadFile(to); err != nil {
        t.Fatal(err)
    }
    if !bytes.Equal(resultTest, result) {
        t.Fatalf("result: %s, expected: %s\n", resultTest, result)
    }
}

func TestSimple(t *testing.T) {
    Scenario(t, []byte("1234567890"), 0, 0, []byte("1234567890"), false)
}

func TestOnlyLimit(t *testing.T) {
    Scenario(t, []byte("1234567890"), 0, 5, []byte("12345"), false)
}

func TestOnlyOffset(t *testing.T) {
    Scenario(t, []byte("1234567890"), 5, 0, []byte("67890"), false)
}

func TestOffsetLimit(t *testing.T) {
    Scenario(t, []byte("1234567890"), 3, 4, []byte("4567"), false)
}

func TestOffsetLimit2(t *testing.T) {
    Scenario(t, []byte("1234567890"), 5, 100, []byte("67890"), false)
}

func TestNegLargeOffset (t *testing.T) {
    Scenario(t, []byte("1234567890"), 20, 0, []byte{}, true)
}

func TestNegOffsetEOF (t *testing.T) {
    Scenario(t, []byte("1234567890"), 10, 0, []byte{}, true)
}
