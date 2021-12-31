package uid

import (
	"github.com/satori/go.uuid"
	"regexp"
)

func UUid() string {
	//todo
	return uuid.NewV4().String()
}

var reg = regexp.MustCompile(`^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$`)

func IsUUID(id string) bool {
	return reg.MatchString(id)
}

