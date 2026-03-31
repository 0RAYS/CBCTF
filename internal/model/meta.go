package model

import (
	"fmt"
	"reflect"
	"slices"
)

type Metadata struct {
	Name         string
	Table        string
	UniqueFields []string
	QueryFields  []string
	SearchFields []string
}

var metadataRegistry = map[reflect.Type]Metadata{
	reflect.TypeFor[Challenge](): {
		Name:         "Challenge",
		Table:        "challenges",
		UniqueFields: []string{"id", "rand_id"},
		QueryFields:  []string{"id", "rand_id", "name", "description", "category", "type", "generator_image"},
		SearchFields: []string{"rand_id", "name", "description", "category", "type", "generator_image"},
	},
	reflect.TypeFor[ChallengeFlag](): {
		Name:         "ChallengeFlag",
		Table:        "challenge_flags",
		UniqueFields: []string{"id"},
		QueryFields:  []string{},
	},
	reflect.TypeFor[Cheat](): {
		Name:         "Cheat",
		Table:        "cheats",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "magic", "ip", "reason", "reason_type", "type", "checked", "hash", "comment", "time", "contest_id"},
		SearchFields: []string{"magic", "ip", "reason", "reason_type", "type", "hash", "comment"},
	},
	reflect.TypeFor[Contest](): {
		Name:         "Contest",
		Table:        "contests",
		UniqueFields: []string{"id", "name"},
		QueryFields:  []string{"id", "name", "description", "prefix", "start", "duration", "hidden"},
		SearchFields: []string{"name", "description", "prefix"},
	},
	reflect.TypeFor[ContestChallenge](): {
		Name:         "ContestChallenge",
		Table:        "contest_challenges",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "contest_id", "challenge_id", "name", "category", "type", "hidden"},
		SearchFields: []string{"name", "category", "type"},
	},
	reflect.TypeFor[ContestFlag](): {
		Name:         "ContestFlag",
		Table:        "contest_flags",
		UniqueFields: []string{"id"},
		QueryFields:  []string{},
	},
	reflect.TypeFor[CronJob](): {
		Name:         "CronJob",
		Table:        "cron_jobs",
		UniqueFields: []string{"id", "name"},
		QueryFields:  []string{"id", "name", "description", "schedule", "success_last", "failure_last", "success", "failure"},
		SearchFields: []string{"name", "description"},
	},
	reflect.TypeFor[Device](): {
		Name:         "Device",
		Table:        "devices",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "user_id", "magic"},
		SearchFields: []string{"magic"},
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
		SearchFields: []string{"challenge_name", "status"},
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
		SearchFields: []string{"name", "resource", "operation", "description"},
	},
	reflect.TypeFor[Pod](): {
		Name:         "Pod",
		Table:        "pods",
		UniqueFields: []string{"id"},
		QueryFields:  []string{},
	},
	reflect.TypeFor[Request](): {
		Name:         "Request",
		Table:        "requests",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "ip", "user_agent", "user_id", "method", "path", "status", "magic"},
		SearchFields: []string{"ip", "user_agent", "method", "path", "magic"},
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
		UniqueFields: []string{"key"},
		QueryFields:  []string{"id", "key"},
		SearchFields: []string{"key"},
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
		Name:         "Team",
		Table:        "teams",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "name", "description", "banned", "hidden", "contest_id"},
		SearchFields: []string{"name", "description"},
	},
	reflect.TypeFor[TeamFlag](): {
		Name:         "TeamFlag",
		Table:        "team_flags",
		UniqueFields: []string{"id"},
		QueryFields:  []string{},
	},
	reflect.TypeFor[Traffic](): {
		Name:         "Traffic",
		Table:        "traffics",
		UniqueFields: []string{"id"},
		QueryFields:  []string{"id", "src_ip", "dst_ip", "type", "subtype"},
		SearchFields: []string{"src_ip", "dst_ip", "type", "subtype"},
	},
	reflect.TypeFor[User](): {
		Name:         "User",
		Table:        "users",
		UniqueFields: []string{"id", "name", "email"},
		QueryFields:  []string{"id", "name", "email", "description", "verified", "banned", "hidden", "provider"},
		SearchFields: []string{"name", "email", "description", "provider"},
	},
	reflect.TypeFor[Victim](): {
		Name:         "Victim",
		Table:        "victims",
		UniqueFields: []string{"id"},
		QueryFields:  []string{},
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

func TableName(m any) string {
	return MetadataOf(m).Table
}

func ModelName(m any) string {
	return MetadataOf(m).Name
}

func UniqueFields(m any) []string {
	return slices.Clone(MetadataOf(m).UniqueFields)
}

func QueryFields(m any) []string {
	return slices.Clone(MetadataOf(m).QueryFields)
}

func SearchFields(m any) []string {
	return slices.Clone(MetadataOf(m).SearchFields)
}
