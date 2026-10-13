{{/*
Expand the name of the chart.
*/}}
{{- define "cbctf.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this.
*/}}
{{- define "cbctf.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart label.
*/}}
{{- define "cbctf.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cbctf.labels" -}}
helm.sh/chart: {{ include "cbctf.chart" . }}
{{ include "cbctf.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cbctf.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cbctf.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
ServiceAccount name
*/}}
{{- define "cbctf.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "cbctf.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
PostgreSQL fullname
*/}}
{{- define "cbctf.postgres.fullname" -}}
{{- printf "%s-postgres" (include "cbctf.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
PostgreSQL labels
*/}}
{{- define "cbctf.postgres.labels" -}}
helm.sh/chart: {{ include "cbctf.chart" . }}
{{ include "cbctf.postgres.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
PostgreSQL selector labels
*/}}
{{- define "cbctf.postgres.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cbctf.name" . }}-postgres
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Valkey resource name.

The workload itself comes from the upstream valkey chart
(https://valkey.io/valkey-helm/) declared as a dependency in Chart.yaml. This
mirrors its naming rule — `<release>-valkey`, or the release name on its own
when it already contains "valkey" — because the platform configuration has to
address the Service the subchart renders. TestBundledValkeyAddressMatchesRenderedService
in render_test.go compares this against the Service that actually shows up, so a
change upstream fails the suite instead of silently pointing at nothing.
*/}}
{{- define "cbctf.valkey.fullname" -}}
{{- if .Values.valkey.fullnameOverride }}
{{- .Values.valkey.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default "valkey" .Values.valkey.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Valkey hostname
*/}}
{{- define "cbctf.valkey.host" -}}
{{- if .Values.valkey.enabled }}
{{- printf "%s.%s" (include "cbctf.valkey.fullname" .) .Release.Namespace }}
{{- else }}
{{- required "valkey.externalHost is required when valkey.enabled is false" .Values.valkey.externalHost }}
{{- end }}
{{- end }}

{{/*
Valkey password the platform authenticates with.

An external server, and a bundled instance whose password lives in
auth.usersExistingSecret, can only be described by valkey.password. Otherwise
the inline ACL entry is the single source of truth and is read from here, so
the two can never drift apart.
*/}}
{{- define "cbctf.valkey.password" -}}
{{- if not .Values.valkey.enabled }}
{{- required "valkey.password is required when valkey.enabled is false" .Values.valkey.password }}
{{- else if .Values.valkey.auth.usersExistingSecret }}
{{- required "valkey.password is required when valkey.auth.usersExistingSecret is set" .Values.valkey.password }}
{{- else }}
{{- required "valkey.auth.aclUsers.default.password is required for the bundled Valkey" (dig "auth" "aclUsers" "default" "password" "" .Values.valkey) }}
{{- end }}
{{- end }}

{{/*
PostgreSQL service hostname
*/}}
{{- define "cbctf.postgres.host" -}}
{{- if .Values.postgres.enabled }}
{{- printf "%s.%s" (include "cbctf.postgres.fullname" .) .Release.Namespace }}
{{- else }}
{{- required "postgres.externalHost is required when postgres.enabled is false" .Values.postgres.externalHost }}
{{- end }}
{{- end }}

{{/*
Combined imagePullSecrets: merges .Values.imagePullSecrets with the auto-generated
registry secret (when imageCredentials are provided). Outputs the full
`imagePullSecrets:` block, or empty string if no secrets are configured.
*/}}
{{- define "cbctf.imagePullSecrets" -}}
{{- $secrets := .Values.imagePullSecrets | default list -}}
{{- with .Values.imageCredentials -}}
  {{- if and .registry .username .password -}}
    {{- $secrets = append $secrets (dict "name" (printf "%s-registry" (include "cbctf.fullname" $))) -}}
  {{- end -}}
{{- end -}}
{{- if $secrets -}}
imagePullSecrets:
  {{- toYaml $secrets | nindent 2 }}
{{- end -}}
{{- end }}
