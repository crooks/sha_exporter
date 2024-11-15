package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Masterminds/log-go"
)

type etcGroupEntry struct {
	name     string
	password string
	gid      string
	users    string
	conforms bool
}

func findGroups(groupFileName string) (countSuccess, countFail int, err error) {
	groupFile, err := os.Open(groupFileName)
	if err != nil {
		return
	}
	defer groupFile.Close()
	scanner := bufio.NewScanner(groupFile)
	for scanner.Scan() {
		line := scanner.Text()
		groupFields := strings.Split(line, ":")
		// Test if config contains an entry for this group.  If it does, assign the expected hash.
		cfgGroup, ok := cfg.Groups[groupFields[0]]
		// If there isn't a dictionary entry in cfg.Groups for this line in the file, move on.
		if !ok {
			log.Tracef("Unwanted group \"%s\".  Continuing", groupFields[0])
			continue
		}
		group := initGroup(groupFields)
		// fileSha is the SHA hash of the group line in /etc/group
		fileSha := group.usersHash()
		// If the SHA hash defined in the configuration matches the hash
		// generated from the users in the group file, set conforms to true.
		group.conforms = cfgGroup.Hash == fileSha
		if group.conforms {
			log.Tracef("Hash collision for %s: %s", groupFields[0], fileSha)
			countSuccess++
		} else {
			log.Debugf("No hash collision for %s: Expected=%s, Got=%s", groupFields[0], cfgGroup.Hash, fileSha)
			countFail++
		}

		prom.groupSHA.WithLabelValues(group.name, group.gid).Set(bool2Float(group.conforms))
		prom.groupUsers.WithLabelValues(group.name, group.gid).Set(group.countUsers())
	}
	return
}

// debugGroups prints the sorted group members and the associated hash
func debugGroups(groupFileName string) (err error) {
	groupFile, err := os.Open(groupFileName)
	if err != nil {
		return
	}
	defer groupFile.Close()
	scanner := bufio.NewScanner(groupFile)
	for scanner.Scan() {
		line := scanner.Text()
		groupFields := strings.Split(line, ":")
		// Test if config contains an entry for this group.  If it does, assign the expected hash.
		_, ok := cfg.Groups[groupFields[0]]
		// If there isn't a dictionary entry in cfg.Groups for this line in the file, move on.
		if ok {
			group := initGroup(groupFields)
			group.debug()
		}
	}
	return
}

func initGroup(groupFields []string) *etcGroupEntry {
	return &etcGroupEntry{
		name:     groupFields[0],
		password: groupFields[1],
		gid:      groupFields[2],
		users:    groupFields[3],
	}
}

// debug prints the sorted users string and the associated hash
func (group *etcGroupEntry) debug() {
    fmt.Printf("%s:-\n", group.name)
	fmt.Printf("  %s\n", group.sortUsers())
	fmt.Printf("  %s\n", group.usersHash())
}

// Take a users string E.g. "user3,user1,user2" and return "user1,user2,user3"
func (group *etcGroupEntry) sortUsers() string {
	users := strings.Split(group.users, ",")
	slices.Sort(users)
	return strings.Join(users, ",")
}

func (group *etcGroupEntry) countUsers() float64 {
	users := strings.Split(group.users, ",")
	return float64(len(users))
}

func (group *etcGroupEntry) usersHash() string {
	h := sha256.New()
	h.Write([]byte(group.sortUsers()))
	return hex.EncodeToString(h.Sum(nil))
}
