// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"text/template"
)

const squidTemplate = `# uyuni-proxy-squid.service, generated by mgrpxy
# Use an uyuni-proxy-squid.service.d/local.conf file to override

[Unit]
Description=Uyuni proxy squid container service
Wants=network.target
After=network-online.target
BindsTo=uyuni-proxy-pod.service
After=uyuni-proxy-pod.service

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Environment=UYUNI_IMAGE={{ .Image }}
{{- if .HttpProxyFile }}
EnvironmentFile={{ .HttpProxyFile }}
{{- end }}
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-proxy-squid.pid %t/uyuni-proxy-squid.ctr-id

ExecStart=/usr/bin/podman run \
	--conmon-pidfile %t/uyuni-proxy-squid.pid \
	--cidfile %t/uyuni-proxy-squid.ctr-id \
	--cgroups=no-conmon \
	--pod-id-file %t/uyuni-proxy-pod.pod-id -d \
	--replace -dt \
	-v /etc/uyuni/proxy:/etc/uyuni:ro \
	{{- range $name, $path := .Volumes }}
	-v {{ $name }}:{{ $path }} \
	{{- end }}
	--name uyuni-proxy-squid \
	${UYUNI_IMAGE}

ExecStop=/usr/bin/podman stop --ignore --cidfile %t/uyuni-proxy-squid.ctr-id -t 10
ExecStopPost=/usr/bin/podman rm --ignore -f --cidfile %t/uyuni-proxy-squid.ctr-id
PIDFile=%t/uyuni-proxy-squid.pid
TimeoutStopSec=60
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

type SquidTemplateData struct {
	Volumes       map[string]string
	HttpProxyFile string
	Image         string
}

func (data SquidTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(squidTemplate))
	return t.Execute(wr, data)
}
