package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func IsAdmin(userID int64) bool {
	adminIDsStr := os.Getenv("ADMIN_IDS")
	if adminIDsStr == "" {
		return false
	}

	ids := strings.Split(adminIDsStr, ",")
	for _, idStr := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			continue
		}
		if id == userID {
			return true
		}
	}
	return false
}

func GetAdminIDs() []int64 {
	adminIDsStr := os.Getenv("ADMIN_IDS")
	if adminIDsStr == "" {
		return nil
	}

	var ids []int64
	for _, idStr := range strings.Split(adminIDsStr, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids
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