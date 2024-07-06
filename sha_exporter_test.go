package main

import (
	"log"
	"os"
	"strings"
	"testing"
)

func makeTestFile(testFileName string) os.File {
	testFile, err := os.CreateTemp("", testFileName)
	if err != nil {
		log.Fatalf("Unable to create test file: %v", err)
	}
	_, err = testFile.WriteString(`group1:x:4:user2,user3,user1
	group2:x:5:
	group3:x:6:user1`)
	if err != nil {
		os.Remove(testFile.Name())
		log.Fatalf("Unable to write to test file: %v", err)
	}
	return *testFile
}

func TestInitGroup(t *testing.T) {
	testSlice := strings.Split("group1:x:4:user2,user3,user1", ":")
	group := initGroup(testSlice)
	if group.name != "group1" {
		t.Errorf("Incorrect user name.  Expected=\"group1\", Got=\"%s\"", group.name)
	}
}

func TestMemberHash(t *testing.T) {
	testSlice := strings.Split("group1:x:4:user2,user3,user1", ":")
	group := initGroup(testSlice)
	// Hex representation of sha256("user1,user2,user3")
	cfgSha := "6b1567cecb30391d3e64d4698edc18c91cab43a088e3823803a5864b49fafada"
	fileSha := group.usersHash()
	if cfgSha != fileSha {
		t.Errorf("Unexpected SHA256 hash of group members.  Expected=%s, Got=%s", cfgSha, fileSha)
	}
	// Now test it fails when a collision doesn't happen
	group.users = "user2,user3,fail1"
	fileSha = group.usersHash()
	if cfgSha == fileSha {
		t.Errorf("Unexpected SHA256 collision using test string: %s", group.users)
	}
}
