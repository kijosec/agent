package kubectl

import "testing"

func TestParseHelmChartLabel(t *testing.T) {
	tests := []struct {
		input       string
		wantName    string
		wantVersion string
	}{
		{"ingress-nginx-4.9.0", "ingress-nginx", "4.9.0"},
		{"kube-prometheus-stack-55.0.0", "kube-prometheus-stack", "55.0.0"},
		{"cert-manager-1.13.2", "cert-manager", "1.13.2"},
		{"no-version", "no-version", ""},
		{"", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotName, gotVersion := parseHelmChartLabel(tt.input)
			if gotName != tt.wantName {
				t.Errorf("name = %q, want %q", gotName, tt.wantName)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("version = %q, want %q", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestExtractHelmMetadata(t *testing.T) {
	tests := []struct {
		name            string
		labels          map[string]string
		wantRelease     string
		wantChart       string
		wantVersion     string
	}{
		{
			name: "helm-managed resource",
			labels: map[string]string{
				"app.kubernetes.io/managed-by": "Helm",
				"app.kubernetes.io/instance":   "my-ingress",
				"helm.sh/chart":                "ingress-nginx-4.9.0",
			},
			wantRelease: "my-ingress",
			wantChart:   "ingress-nginx",
			wantVersion: "4.9.0",
		},
		{
			name: "non-helm resource",
			labels: map[string]string{
				"app": "myapp",
			},
			wantRelease: "",
			wantChart:   "",
			wantVersion: "",
		},
		{
			name: "helm without chart label",
			labels: map[string]string{
				"app.kubernetes.io/managed-by": "Helm",
				"app.kubernetes.io/instance":   "my-release",
			},
			wantRelease: "my-release",
			wantChart:   "",
			wantVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRelease, gotChart, gotVersion := extractHelmMetadata(tt.labels)
			if gotRelease != tt.wantRelease {
				t.Errorf("release = %q, want %q", gotRelease, tt.wantRelease)
			}
			if gotChart != tt.wantChart {
				t.Errorf("chart = %q, want %q", gotChart, tt.wantChart)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("version = %q, want %q", gotVersion, tt.wantVersion)
			}
		})
	}
}
