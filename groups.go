package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
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

func findGroups(groupFileName string) error {
	groupFile, err := os.Open(groupFileName)
	if err != nil {
		return err
	}
	defer groupFile.Close()
	scanner := bufio.NewScanner(groupFile)
	for scanner.Scan() {
		line := scanner.Text()
		groupFields := strings.Split(line, ":")
		// Test if config contains an entry for this group.  If it does, assign the expected hash.
		cfgSha, ok := cfg.Groups[groupFields[0]]
		// If there isn't a dictionary entry in cfg.Groups for this line in the file, move on.
		if !ok {
			log.Tracef("Unwanted group \"%s\".  Continuing", groupFields[0])
			continue
		}
		log.Debugf("Processing %s group", groupFields[0])
		group := initGroup(groupFields)
		fileSha := group.usersHash()
		// If the SHA hash defined in the configuration matches the hash
		// generated from the users in the group file, set conforms to true.
		if cfgSha == fileSha {
			group.conforms = true
		}
		prom.groupSHA.WithLabelValues(group.name, group.gid).Set(bool2Float(group.conforms))
		prom.groupUsers.WithLabelValues(group.name, group.gid).Set(group.countUsers())
	}
	return nil
}

func initGroup(groupFields []string) *etcGroupEntry {
	return &etcGroupEntry{
		name:     groupFields[0],
		password: groupFields[1],
		gid:      groupFields[2],
		users:    groupFields[3],
	}
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
