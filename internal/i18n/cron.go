package i18n

var Cron = struct {
	NotRegistered        string
	InvalidSchedule     string
	PostgresPoolError   string
	PostgresConnError   string
	PostgresLockError   string
	PostgresListError   string
	PostgresIndexError  string
	PostgresCloseError  string
	PostgresRebuildFail string
}{
	NotRegistered:        "cron.notRegistered",
	InvalidSchedule:      "cron.invalidSchedule",
	PostgresPoolError:    "cron.postgresPoolError",
	PostgresConnError:    "cron.postgresConnError",
	PostgresLockError:    "cron.postgresLockError",
	PostgresListError:    "cron.postgresListError",
	PostgresIndexError:   "cron.postgresIndexError",
	PostgresCloseError:   "cron.postgresCloseError",
	PostgresRebuildFail:  "cron.postgresRebuildFail",
}
