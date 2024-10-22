package utils

const (
	UptimePattern          = `^([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	BackOffLimit           = 0
	ClusterDetailsTemplate = `Name: {{.Name}} | Namespace: {{.Namespace}} | State: {{.State}} | Runtime: {{.Uptime}} | Timezone: {{.Timezone}}`
)
