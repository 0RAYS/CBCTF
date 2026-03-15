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
MySQL fullname
*/}}
{{- define "cbctf.mysql.fullname" -}}
{{- printf "%s-mysql" (include "cbctf.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
MySQL labels
*/}}
{{- define "cbctf.mysql.labels" -}}
helm.sh/chart: {{ include "cbctf.chart" . }}
{{ include "cbctf.mysql.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
MySQL selector labels
*/}}
{{- define "cbctf.mysql.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cbctf.name" . }}-mysql
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Redis fullname
*/}}
{{- define "cbctf.redis.fullname" -}}
{{- printf "%s-redis" (include "cbctf.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Redis labels
*/}}
{{- define "cbctf.redis.labels" -}}
helm.sh/chart: {{ include "cbctf.chart" . }}
{{ include "cbctf.redis.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Redis selector labels
*/}}
{{- define "cbctf.redis.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cbctf.name" . }}-redis
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
MySQL service hostname
*/}}
{{- define "cbctf.mysql.host" -}}
{{- if .Values.mysql.enabled }}
{{- include "cbctf.mysql.fullname" . }}
{{- else }}
{{- required "mysql.externalHost is required when mysql.enabled is false" .Values.mysql.externalHost }}
{{- end }}
{{- end }}

{{/*
Redis service hostname
*/}}
{{- define "cbctf.redis.host" -}}
{{- if .Values.redis.enabled }}
{{- include "cbctf.redis.fullname" . }}
{{- else }}
{{- required "redis.externalHost is required when redis.enabled is false" .Values.redis.externalHost }}
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
