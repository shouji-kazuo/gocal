{{ range $event := .Events }}
{{ formatDate $event.Start "01/02 15:04"}} ～ {{ formatDate $event.End "15:04" }}   {{ $event.Summary }}
{{ end }}