{{- define "resources" -}}
{{ if .resources -}}
resources:
{{- .resources | toYaml | nindent 2 -}}
{{ else }}
resources:
  requests:
    cpu: 0.50
    memory: "500Mi"
  limits:
    cpu: 1.0
    memory: "1000Mi"
{{- end -}}
{{- end -}}
