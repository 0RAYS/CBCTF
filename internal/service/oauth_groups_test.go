package service

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

func TestOAuthGroupsDistinguishUnmappedNamesFromDatabaseErrors(t *testing.T) {
	for _, failRead := range []bool{false, true} {
		database := queryFailureDB(t)
		if failRead {
			_ = database.Callback().Query().Before("gorm:query").Register("test:failed_read", func(tx *gorm.DB) { _ = tx.AddError(errors.New("read unavailable")) })
		}
		ret := assignOAuthGroups(database, model.User{ID: 1}, model.Oauth{GroupsClaim: "{groups}"}, map[string]any{"groups": []string{"unmapped"}})
		if ret.OK == failRead {
			t.Fatalf("database error confused with unmapped group: %+v", ret)
		}
	}
}

func TestOAuthGroupWriteFailurePropagates(t *testing.T) {
	database := queryFailureDB(t)
	_ = database.Callback().Query().After("gorm:query").Register("test:group", func(tx *gorm.DB) {
		if group, ok := tx.Statement.Dest.(*model.Group); ok {
			group.ID = 1
			group.Name = "mapped"
			tx.RowsAffected = 1
		}
	})
	_ = database.Callback().Create().Before("gorm:create").Register("test:membership_failure", func(tx *gorm.DB) { _ = tx.AddError(errors.New("membership write failed")) })
	ret := assignOAuthGroups(database, model.User{ID: 1}, model.Oauth{GroupsClaim: "groups"}, map[string]any{"groups": []string{"mapped"}})
	if ret.OK || ret.Msg != i18n.Model.UserGroup.CreateError {
		t.Fatalf("membership failure lost: %+v", ret)
	}
}
