package service

import (
	"slices"
	"strings"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

// Runs inside the new user's registration transaction. Optional external group
// names may be unmapped; failures reading or writing the database are not skips.
func assignOAuthGroups(tx *gorm.DB, user model.User, provider model.Oauth, response map[string]any) model.RetVal {
	repo := db.InitGroupRepo(tx)
	assigned := make(map[uint]bool)
	assign := func(group model.Group) model.RetVal {
		if assigned[group.ID] {
			return model.SuccessRetVal()
		}
		ret := db.AppendUserToGroup(tx, user, group)
		if ret.OK {
			assigned[group.ID] = true
		}
		return ret
	}
	if provider.GroupsClaim != "" {
		key := strings.TrimSuffix(strings.TrimPrefix(provider.GroupsClaim, "{"), "}")
		if groups, ok := utils.GetClaimRawValue[[]string](response, key); ok {
			for _, name := range groups {
				group, ret := repo.GetByUniqueField("name", name)
				if !ret.OK {
					if ret.Msg == i18n.Model.NotFound {
						continue
					}
					return ret
				}
				if ret = assign(group); !ret.OK {
					return ret
				}
			}
			if provider.AdminGroup != "" && slices.Contains(groups, provider.AdminGroup) {
				group, ret := repo.GetByUniqueField("name", model.AdminGroupName)
				if !ret.OK {
					return ret
				}
				if ret = assign(group); !ret.OK {
					return ret
				}
			}
		}
	}
	if provider.DefaultGroup != 0 {
		group, ret := repo.GetByID(provider.DefaultGroup)
		if !ret.OK {
			return ret
		}
		return assign(group)
	}
	return model.SuccessRetVal()
}
