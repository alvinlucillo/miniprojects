{{- define "mongodb.name" -}}
{{ .Release.Name }}-{{ .Chart.Name }}
{{- end -}}

{{- define "mongodb.fullname" -}}
{{ include "mongodb.name" . }}
{{- end -}}
