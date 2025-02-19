package packaging

import (
	"os"
	"testing"

	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/types"
	godeptypes "github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_packagingAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		want      *analyzer.AnalysisResult
		wantErr   string
	}{
		{
			name:      "egg zip",
			inputFile: "testdata/kitchen-1.2.6-py2.7.egg",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "testdata/kitchen-1.2.6-py2.7.egg",
						Libraries: []types.LibraryInfo{
							{
								FilePath: "testdata/kitchen-1.2.6-py2.7.egg",
								Library: godeptypes.Library{
									Name:    "kitchen",
									Version: "1.2.6",
									License: "LGPLv2+",
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "egg-info",
			inputFile: "testdata/happy.egg-info/PKG-INFO",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "testdata/happy.egg-info/PKG-INFO",
						Libraries: []types.LibraryInfo{
							{
								FilePath: "testdata/happy.egg-info/PKG-INFO",
								Library: godeptypes.Library{
									Name:    "distlib",
									Version: "0.3.1",
									License: "Python license",
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "egg-info no-license",
			inputFile: "testdata/no_license.egg-info/PKG-INFO",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "testdata/no_license.egg-info/PKG-INFO",
						Libraries: []types.LibraryInfo{
							{
								FilePath: "testdata/no_license.egg-info/PKG-INFO",
								Library: godeptypes.Library{
									Name:    "setuptools",
									Version: "51.3.3",
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "wheel",
			inputFile: "testdata/happy.dist-info/METADATA",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "testdata/happy.dist-info/METADATA",
						Libraries: []types.LibraryInfo{
							{
								FilePath: "testdata/happy.dist-info/METADATA",
								Library: godeptypes.Library{
									Name:    "distlib",
									Version: "0.3.1",
									License: "Python license",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := os.ReadFile(tt.inputFile)
			require.NoError(t, err)

			a := packagingAnalyzer{}
			got, err := a.Analyze(analyzer.AnalysisTarget{
				FilePath: tt.inputFile,
				Content:  b,
			})

			if tt.wantErr != "" {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

}

func Test_packagingAnalyzer_Required(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{
			name:     "egg",
			filePath: "python2.7/site-packages/cssutils-1.0-py2.7.egg/EGG-INFO/PKG-INFO",
			want:     true,
		},
		{
			name:     "egg-info",
			filePath: "python3.8/site-packages/wrapt-1.12.1.egg-info",
			want:     true,
		},
		{
			name:     "egg-info PKG-INFO",
			filePath: "python3.8/site-packages/wrapt-1.12.1.egg-info/PKG-INFO",
			want:     true,
		},
		{
			name:     "wheel",
			filePath: "python3.8/site-packages/wrapt-1.12.1.dist-info/METADATA",
			want:     true,
		},
		{
			name:     "sad",
			filePath: "random/PKG-INFO",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := packagingAnalyzer{}
			got := a.Required(tt.filePath, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
