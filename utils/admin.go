package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	adminIDs  []int64
	adminOnce sync.Once
)

// loadAdminIDs parses ADMIN_IDS env once and caches the result.
func loadAdminIDs() {
	adminOnce.Do(func() {
		raw := os.Getenv("ADMIN_IDS")
		if raw == "" {
			return
		}
		for _, s := range strings.Split(raw, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err != nil {
				continue
			}
			adminIDs = append(adminIDs, id)
		}
	})
}

func IsAdmin(userID int64) bool {
	loadAdminIDs()
	for _, id := range adminIDs {
		if id == userID {
			return true
		}
	}
	return false
}

func GetAdminIDs() []int64 {
	loadAdminIDs()
	return adminIDs
}

func FormatAdminList() string {
	ids := GetAdminIDs()
	if len(ids) == 0 {
		return "لا يوجد مشرفون"
	}

	var parts []string
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("`%d`", id))
	}
	return strings.Join(parts, ", ")
}
