package prettier

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewManifest(t *testing.T) {
	podYaml := `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  labels:
    key: value
spec:
  containers:
    - name: thing
      image: myimage:latest
      ports:
        - name: web
          containerPort: 80
          protocol: TCP	
`
	podTemplateYaml := `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  labels:
    key: value
spec:
  containers:
    - name: thing
      image: myimage:latest
      ports:
        - name: web
          containerPort: {{ .Values | get "port" 80 }}
          protocol: TCP
{{ .AdditionalContainers | toYaml }}
`

	type args struct {
		yml string
	}
	tests := []struct {
		name    string
		args    args
		want    Manifest
		wantErr bool
	}{
		{
			"default pod",
			args{
				yml: podYaml,
			},
			Manifest{
				Kind: "Pod",
				Metadata: struct {
					Name      string
					Namespace string
				}{
					Name:      "test-pod",
					Namespace: "",
				},
				Yaml: strings.TrimSpace(podYaml),
			},
			false,
		},
		{
			"pod with template",
			args{
				yml: podTemplateYaml,
			},
			Manifest{
				Kind: "Pod",
				Metadata: struct {
					Name      string
					Namespace string
				}{
					Name:      "test-pod",
					Namespace: "",
				},
				Yaml: strings.TrimSpace(podTemplateYaml),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManifest(tt.args.yml)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitManifests(t *testing.T) {
	podOne := `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod1
  labels:
    key: value
spec:
  containers:
    - name: thing
      image: myimage:latest
      ports:
        - name: web
          containerPort: 80
          protocol: TCP	
`
	podTwo := `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod2
  labels:
    key: value
spec:
  containers:
    - name: thing
      image: myimage:latest
      ports:
        - name: web
          containerPort: 80
          protocol: TCP
`

	type args struct {
		yml string
	}
	tests := []struct {
		name    string
		args    args
		want    []Manifest
		wantErr bool
	}{
		{
			name: "two pods",
			args: args{
				"\n---\n" + podOne + "\n---\n" + podTwo + "\n---\n",
			},
			want: []Manifest{
				{
					Kind: "Pod",
					Metadata: struct {
						Name      string
						Namespace string
					}{
						Name: "test-pod1",
					},
					Yaml: strings.TrimSpace(podOne),
				},
				{
					Kind: "Pod",
					Metadata: struct {
						Name      string
						Namespace string
					}{
						Name: "test-pod2",
					},
					Yaml: strings.TrimSpace(podTwo),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitManifests(tt.args.yml)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitManifests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitManifests() = %v, want %v", got, tt.want)
			}
		})
	}
}
