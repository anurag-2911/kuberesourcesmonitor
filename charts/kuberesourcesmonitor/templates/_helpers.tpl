{{/*
Expand the name of the chart.
*/}}
{{- define "kuberesourcesmonitor.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fullname.
*/}}
{{- define "kuberesourcesmonitor.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" $name .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Generate the name of the ServiceAccount to use.
*/}}
{{- define "kuberesourcesmonitor.serviceAccountName" -}}
{{- if .Values.serviceAccount.name -}}
{{- .Values.serviceAccount.name -}}
{{- else -}}
{{ include "kuberesourcesmonitor.fullname" . }}-sa
{{- end -}}
{{- end -}}

{{/*
Include the chart version.
*/}}
{{- define "kuberesourcesmonitor.chart" -}}
{{ .Chart.Name }}-{{ .Chart.Version }}
{{- end -}}
