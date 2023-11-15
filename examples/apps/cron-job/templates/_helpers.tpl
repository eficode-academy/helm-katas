{{- define "resources" -}}
{{ if .resources -}}
  {{- .resources | toYaml | nindent 2 -}}
{{ else }}
  requests:
    cpu: 0.50
    memory: "500Mi"
  limits:
    cpu: 1.0
    memory: "1000Mi"
{{- end -}}
{{- end -}}

{{- define "metadata" -}}
app: {{ .Values.app }}
managed_by: {{ .Release.Service }}
chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end -}}