package model

import (
	"fmt"
	"reflect"
	"slices"
)

type Metadata struct {
	Name          string
	Table         string
	UniqueFields  []string
	UniqueIndexes [][]string
}

var metadataRegistry = map[reflect.Type]Metadata{
	reflect.TypeFor[Challenge](): {
		Name:         "Challenge",
		Table:        "challenges",
		UniqueFields: []string{"id", "rand_id"},
	},
	reflect.TypeFor[Branding](): {
		Name:         "Branding",
		Table:        "brandings",
		UniqueFields: []string{"id", "code"},
	},
	reflect.TypeFor[ChallengeFlag](): {
		Name:         "ChallengeFlag",
		Table:        "challenge_flags",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Cheat](): {
		Name:         "Cheat",
		Table:        "cheats",
		UniqueFields: []string{"id", "hash"},
	},
	reflect.TypeFor[Contest](): {
		Name:         "Contest",
		Table:        "contests",
		UniqueFields: []string{"id", "name"},
	},
	reflect.TypeFor[ContestChallenge](): {
		Name:          "ContestChallenge",
		Table:         "contest_challenges",
		UniqueFields:  []string{"id"},
		UniqueIndexes: [][]string{{"contest_id", "challenge_id"}},
	},
	reflect.TypeFor[ContestFlag](): {
		Name:         "ContestFlag",
		Table:        "contest_flags",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[CronJob](): {
		Name:         "CronJob",
		Table:        "cron_jobs",
		UniqueFields: []string{"id", "name"},
	},
	reflect.TypeFor[Email](): {
		Name:         "Email",
		Table:        "emails",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Event](): {
		Name:         "Event",
		Table:        "events",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[File](): {
		Name:         "File",
		Table:        "files",
		UniqueFields: []string{"id", "rand_id"},
	},
	reflect.TypeFor[Generator](): {
		Name:         "Generator",
		Table:        "generators",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Group](): {
		Name:         "Group",
		Table:        "groups",
		UniqueFields: []string{"id", "name"},
	},
	reflect.TypeFor[Notice](): {
		Name:         "Notice",
		Table:        "notices",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Oauth](): {
		Name:         "Oauth",
		Table:        "oauths",
		UniqueFields: []string{"id", "provider"},
	},
	reflect.TypeFor[Permission](): {
		Name:         "Permission",
		Table:        "permissions",
		UniqueFields: []string{"id", "name"},
	},
	reflect.TypeFor[Pod](): {
		Name:         "Pod",
		Table:        "pods",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Request](): {
		Name:         "Request",
		Table:        "requests",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Role](): {
		Name:         "Role",
		Table:        "roles",
		UniqueFields: []string{"id", "name"},
	},
	reflect.TypeFor[Setting](): {
		Name:         "Setting",
		Table:        "settings",
		UniqueFields: []string{"id", "key"},
	},
	reflect.TypeFor[Smtp](): {
		Name:         "Smtp",
		Table:        "smtps",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Task](): {
		Name:         "Task",
		Table:        "tasks",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Submission](): {
		Name:         "Submission",
		Table:        "submissions",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Team](): {
		Name:          "Team",
		Table:         "teams",
		UniqueFields:  []string{"id"},
		UniqueIndexes: [][]string{{"contest_id", "name"}},
	},
	reflect.TypeFor[TeamFlag](): {
		Name:          "TeamFlag",
		Table:         "team_flags",
		UniqueFields:  []string{"id"},
		UniqueIndexes: [][]string{{"team_id", "contest_flag_id"}},
	},
	reflect.TypeFor[Traffic](): {
		Name:         "Traffic",
		Table:        "traffics",
		UniqueFields: []string{"id", "victim_id"},
	},
	reflect.TypeFor[User](): {
		Name:          "User",
		Table:         "users",
		UniqueFields:  []string{"id", "name", "email"},
		UniqueIndexes: [][]string{{"provider", "provider_user_id"}},
	},
	reflect.TypeFor[Victim](): {
		Name:         "Victim",
		Table:        "victims",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[Webhook](): {
		Name:         "Webhook",
		Table:        "webhooks",
		UniqueFields: []string{"id"},
	},
	reflect.TypeFor[WebhookHistory](): {
		Name:         "WebhookHistory",
		Table:        "webhook_histories",
		UniqueFields: []string{"id"},
	},
}

func MetadataOf(m any) Metadata {
	if m == nil {
		panic("model metadata lookup on nil value")
	}
	t := reflect.TypeOf(m)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	meta, ok := metadataRegistry[t]
	if !ok {
		panic(fmt.Sprintf("missing model metadata for %s", t.String()))
	}
	return meta
}

func Name(m any) string {
	return MetadataOf(m).Name
}

func UniqueFields(m any) []string {
	return slices.Clone(MetadataOf(m).UniqueFields)
}

func UniqueIndexes(m any) [][]string {
	indexes := MetadataOf(m).UniqueIndexes
	result := make([][]string, 0, len(indexes))
	for _, index := range indexes {
		result = append(result, slices.Clone(index))
	}
	return result
}
