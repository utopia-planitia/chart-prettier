package prettier

import (
	"reflect"
	"testing"

	"github.com/damoon/fstesting"
	"github.com/spf13/afero"
)

func TestChart(t *testing.T) {
	chart := &Chart{}

	appFs, err := fstesting.InMemoryCopy("testdata/pods-test/chart", "testdata/pods-test/chart")
	if err != nil {
		t.Errorf("create mock filesystem: %v", err)
		return
	}

	err = chart.LoadChart(appFs, "testdata/pods-test/chart")
	if err != nil {
		t.Errorf("loading from existing chart: %v", err)
		return
	}

	err = chart.DeleteFiles(appFs, "testdata/pods-test/chart")
	if err != nil {
		t.Errorf("cleanup preexisting manifests: %v", err)
		return
	}

	err = chart.WriteOut(appFs, "testdata/pods-test/chart")
	if err != nil {
		t.Errorf("create new manifests in chart: %v", err)
		return
	}

	goldenFs := afero.NewOsFs()

	equals, diff, err := fstesting.DiffDir(appFs, goldenFs, "testdata/pods-test/chart", "testdata/pods-test/golden")
	if err != nil {
		t.Errorf("compare test result with golden state: %v", err)
		return
	}

	if !equals {
		t.Errorf("result does not match golden state: %v", diff)
		return
	}
}

func Test_uniqueNames(t *testing.T) {
	testPod := `
apiVersion: v1
kind: Pod
metadata:
  name: test
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
	testPod2 := `
apiVersion: v1
kind: Pod
metadata:
  name: test2
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

	testPodManifest, err := NewManifest(testPod)
	if err != nil {
		t.Errorf("generate manifest object for test-pod: %v", err)
		return
	}

	testPod2Manifest, err := NewManifest(testPod2)
	if err != nil {
		t.Errorf("generate manifest object for test-pod2: %v", err)
		return
	}

	type args struct {
		manifests []Manifest
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]Manifest
		wantErr bool
	}{
		{
			name: "one pod",
			args: args{
				[]Manifest{testPodManifest},
			},
			want: map[string]Manifest{
				"pod": testPodManifest,
			},
			wantErr: false,
		},
		{
			name: "two pod",
			args: args{
				[]Manifest{testPodManifest, testPod2Manifest},
			},
			want: map[string]Manifest{
				"pod-test":  testPodManifest,
				"pod-test2": testPod2Manifest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uniqueNames(tt.args.manifests)
			if (err != nil) != tt.wantErr {
				t.Errorf("uniqueNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqueNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
