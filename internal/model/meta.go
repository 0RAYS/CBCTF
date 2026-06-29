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
	QueryFields   []string
	SearchFields  []string
}

var metadataRegistry = map[reflect.Type]Metadata{
	reflect.TypeFor[Challenge](): {
		Name:         "Challenge",
		Table:        "challenges",
		UniqueFields: []string{"id", "rand_id"},
		QueryFields:  []string{"id", "rand_id", "name", "description", "category", "type", "generator_image"},
		SearchFields: []string{"rand_id", "name", "description", "category", "type", "generator_image"},
	},
	reflect.TypeFor[Branding](): {
		Name:         "Branding",
		Table:        "brandings",
		UniqueFields: []string{"id", "code"},
		QueryFields:  []string{"id", "code"},
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
		QueryFields:  []string{"id", "ip", "reason", "reason_type", "type", "checked", "hash", "comment", "time", "contest_id"},
		SearchFields: []string{"ip", "reason", "reason_type", "type", "hash", "comment"},
	},
	reflect.TypeFor[Contest](): {
		Name:         "Contest",
		Table:        "contests",
		UniqueFields: []string{"id", "name"},
		QueryFields:  []string{"id", "name", "description", "prefix", "start", "duration", "hidden"},
		SearchFields: []string{"name", "description", "prefix"},
	},
	reflect.TypeFor[ContestChallenge](): {
		Name:          "ContestChallenge",
		Table:         "contest_challenges",
		UniqueFields:  []string{"id"},
		UniqueIndexes: [][]string{{"contest_id", "challenge_id"}},
		QueryFields:   []string{"id", "contest_id", "challenge_id", "name", "category", "type", "hidden"},
		SearchFields:  []string{"name", "category", "type"},
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
		QueryFields:  []string{"id", "name", "description", "schedule", "success_last", "failure_last", "success", "failure"},
	},
	reflect.TypeFor[Email](): {
		Name:         "Email",
		Table:        "emails",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "from", "to", "subject", "content", "success", "smtp_id"},
		SearchFields: []string{"from", "to", "subject", "content"},
	},
	reflect.TypeFor[Event](): {
		Name:         "Event",
		Table:        "events",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "type", "success", "ip"},
		SearchFields: []string{"type", "ip"},
	},
	reflect.TypeFor[File](): {
		Name:         "File",
		Table:        "files",
		UniqueFields: []string{"id", "rand_id"},
		QueryFields:  []string{"id", "rand_id", "model", "model_id", "filename", "size", "suffix", "hash", "type"},
		SearchFields: []string{"rand_id", "model", "filename", "suffix", "hash", "type"},
	},
	reflect.TypeFor[Generator](): {
		Name:         "Generator",
		Table:        "generators",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "challenge_id", "challenge_name", "contest_id", "success", "failure", "status"},
	},
	reflect.TypeFor[Group](): {
		Name:         "Group",
		Table:        "groups",
		UniqueFields: []string{"id", "name"},
		QueryFields:  []string{"id", "role_id", "name", "description", "default"},
		SearchFields: []string{"name", "description"},
	},
	reflect.TypeFor[Notice](): {
		Name:         "Notice",
		Table:        "notices",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "title", "content", "type", "contest_id"},
		SearchFields: []string{"title", "content", "type"},
	},
	reflect.TypeFor[Oauth](): {
		Name:         "Oauth",
		Table:        "oauths",
		UniqueFields: []string{"id", "provider"},
		QueryFields:  []string{"id", "provider", "on"},
		SearchFields: []string{"provider"},
	},
	reflect.TypeFor[Permission](): {
		Name:         "Permission",
		Table:        "permissions",
		UniqueFields: []string{"id", "name"},
		QueryFields:  []string{"id", "name", "resource", "operation", "description"},
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
		QueryFields:  []string{"id", "ip", "user_agent", "user_id", "method", "path", "status"},
		SearchFields: []string{"ip", "user_agent", "method", "path"},
	},
	reflect.TypeFor[Role](): {
		Name:         "Role",
		Table:        "roles",
		UniqueFields: []string{"id", "name"},
		QueryFields:  []string{"id", "name", "description", "default"},
		SearchFields: []string{"name", "description"},
	},
	reflect.TypeFor[Setting](): {
		Name:         "Setting",
		Table:        "settings",
		UniqueFields: []string{"id", "key"},
		QueryFields:  []string{"id", "key"},
	},
	reflect.TypeFor[Smtp](): {
		Name:         "Smtp",
		Table:        "smtps",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "address", "host", "on"},
		SearchFields: []string{"address", "host"},
	},
	reflect.TypeFor[Task](): {
		Name:         "Task",
		Table:        "tasks",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "task_id", "type", "queue", "status", "retry_count", "max_retry", "processed_at"},
		SearchFields: []string{"task_id", "type", "queue", "status", "error"},
	},
	reflect.TypeFor[Submission](): {
		Name:         "Submission",
		Table:        "submissions",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "value", "solved", "ip", "user_id", "team_id", "contest_id", "challenge_id", "contest_challenge_id"},
		SearchFields: []string{"value", "ip"},
	},
	reflect.TypeFor[Team](): {
		Name:          "Team",
		Table:         "teams",
		UniqueFields:  []string{"id"},
		UniqueIndexes: [][]string{{"contest_id", "name"}},
		QueryFields:   []string{"id", "name", "description", "banned", "hidden", "contest_id"},
		SearchFields:  []string{"name", "description"},
	},
	reflect.TypeFor[TeamFlag](): {
		Name:          "TeamFlag",
		Table:         "team_flags",
		UniqueFields:  []string{"id"},
		UniqueIndexes: [][]string{{"team_id", "contest_flag_id"}},
		SearchFields:  []string{"value"},
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
		QueryFields:   []string{"id", "name", "email", "description", "verified", "banned", "hidden", "provider"},
		SearchFields:  []string{"name", "email", "description", "provider"},
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
		QueryFields:  []string{"id", "name", "url", "on", "method"},
		SearchFields: []string{"name", "url", "method"},
	},
	reflect.TypeFor[WebhookHistory](): {
		Name:         "WebhookHistory",
		Table:        "webhook_histories",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "success", "webhook_id", "resp_code"},
		SearchFields: []string{"error"},
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

func QueryFields(m any) []string {
	return slices.Clone(MetadataOf(m).QueryFields)
}

func SearchFields(m any) []string {
	return slices.Clone(MetadataOf(m).SearchFields)
}
