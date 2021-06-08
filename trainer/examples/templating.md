This is an example of confusing text templating in Helm.

```TEXT
{{- define "anotherMyTemplate" }}
{{- if and "hej" ( "foo"| upper | ne "bar" ) }}
anders and
{{- end }}
{{ end }}

{{ define "myTemplate" }}
helmIsCool: true
{{ template "anotherMyTemplate" . }}
{{- end -}}
hvor står jeg henne
{{- template "myTemplate" }}
```
