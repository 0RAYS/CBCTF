package db

import (
	"CBCTF/internal/model"
	"context"
)

// AppendUserToTeam Many2Many
func AppendUserToTeam(ctx context.Context, user model.User, team model.Team) error {
	return DB.WithContext(ctx).Model(&team).Association("Users").Append(&user)
}

// AppendUserToContest Many2Many
func AppendUserToContest(ctx context.Context, user model.User, contest model.Contest) error {
	return DB.WithContext(ctx).Model(&contest).Association("Contests").Append(&user)
}

// AppendTeamToContest HasMany
func AppendTeamToContest(ctx context.Context, team model.Team, contest model.Contest) error {
	return DB.WithContext(ctx).Model(&contest).Association("Teams").Append(&team)
}

// DeleteUserFromTeam Many2Many
func DeleteUserFromTeam(ctx context.Context, user model.User, team model.Team) error {
	return DB.WithContext(ctx).Model(&team).Association("Users").Delete(&user)
}

// DeleteUserFromContest Many2Many
func DeleteUserFromContest(ctx context.Context, user model.User, contest model.Contest) error {
	return DB.WithContext(ctx).Model(&contest).Association("Users").Delete(&user)
}

// DeleteTeamFromContest HasMany
func DeleteTeamFromContest(ctx context.Context, team model.Team, contest model.Contest) error {
	return DB.WithContext(ctx).Model(&contest).Association("Teams").Delete(&team)
}
