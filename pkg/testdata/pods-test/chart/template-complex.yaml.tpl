{{ $cfg:=.Values }}
{{ range $runner := .Values.runners }}
{{ range $i, $secret := $runner.pull_secrets }}
{{ if $secret.username }}
apiVersion: v1
kind: Secret
metadata:
  name: {{ $secret.name }}
type: kubernetes.io/dockerconfigjson
stringData:
  .dockerconfigjson: >
    {
      "auths": {
        "{{ $secret.registry }}": {
          "username": "{{ $secret.username }}",
          "password": "{{ $secret.password }}",
          "email": "none@none.com",
          "auth": "{{ (print $secret.username ":" $secret.password ) | b64enc }}"
        }
      }
    }
---
{{ end }}
{{ end }}
{{ end }}
