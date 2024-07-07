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

func TestMemberHashCount(t *testing.T) {
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
	count := group.countUsers()
	if count != 3 {
		t.Errorf("Unexpected user count. Expected=3, Got=%f", count)
	}
}

func TestFileHash(t *testing.T) {
	testFile, err := os.CreateTemp("", "testfilehash.txt")
	if err != nil {
		log.Fatalf("Unable to create test file: %v", err)
	}
	testString := `The Owl and the Pussy-cat went to sea
In a beautiful pea-green boat,
They took some honey, and plenty of money,
Wrapped up in a five-pound note.
The Owl looked up to the stars above,
And sang to a small guitar,
O lovely Pussy! O Pussy, my love,
What a beautiful Pussy you are,
You are,
You are!
What a beautiful Pussy you are!
`
	//fmt.Printf("%x\n", testString)
	_, err = testFile.WriteString(testString)
	testFile.Close()
	defer os.Remove(testFile.Name())
	if err != nil {
		log.Fatalf("Unable to write to test file: %v", err)
	}

	expectedHash := "3055540351c2f0e6ce1c4bf1a63315dda30a646bbcc9e6c12ab9cddbaecf59de"
	gotHash, err := fileHash(testFile.Name())
	if err != nil {
		log.Fatalf("File hash failed: %v", err)
	}
	if expectedHash != gotHash {
		t.Errorf("Unexpected file hash. Expected=%s, Got=%s", expectedHash, gotHash)
	}
}
