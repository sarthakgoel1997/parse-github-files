package service

import (
	"net/http"
	"parse-github-files/model"
	"reflect"
	"testing"
)

func Test_prepareGitHubAPIRequest(t *testing.T) {
	type args struct {
		repository string
	}
	tests := []struct {
		name          string
		args          args
		wantGithubReq *http.Request
		wantBaseUrl   string
		wantErr       bool
	}{
		{
			name: "normal test case to generate github request and base url",
			args: args{
				repository: "https://github.com/velancio/vulnerability_scans",
			},
			wantGithubReq: func() *http.Request {
				githubReq, _ := http.NewRequest("GET", "https://api.github.com/repos/velancio/vulnerability_scans/contents/", nil)
				addhttpAuthRequestHeaders(githubReq)
				return githubReq
			}(),
			wantBaseUrl: "https://api.github.com/repos/velancio/vulnerability_scans/contents/",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGithubReq, gotBaseUrl, err := prepareGitHubAPIRequest(tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareGitHubAPIRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGithubReq, tt.wantGithubReq) {
				t.Errorf("prepareGitHubAPIRequest() gotGithubReq = %v, want %v", gotGithubReq, tt.wantGithubReq)
			}
			if gotBaseUrl != tt.wantBaseUrl {
				t.Errorf("prepareGitHubAPIRequest() gotBaseUrl = %v, want %v", gotBaseUrl, tt.wantBaseUrl)
			}
		})
	}
}

func Test_getDataFromGitHub(t *testing.T) {
	type args struct {
		baseUrl   string
		file      string
		githubReq *http.Request
	}
	tests := []struct {
		name         string
		args         args
		wantFileData model.FileData
		wantErr      bool
	}{
		{
			args: args{
				baseUrl: "https://api.github.com/repos/velancio/vulnerability_scans/contents/",
				file:    "vulnscan1011.json",
				githubReq: func() *http.Request {
					githubReq, _ := http.NewRequest("GET", "https://api.github.com/repos/velancio/vulnerability_scans/contents/", nil)
					addhttpAuthRequestHeaders(githubReq)
					return githubReq
				}(),
			},
			wantFileData: model.FileData{
				Name:     "vulnscan1011.json",
				HtmlUrl:  "https://github.com/velancio/vulnerability_scans/blob/main/vulnscan1011.json",
				Content:  "WwogIHsKICAgICJzY2FuUmVzdWx0cyI6IHsKICAgICAgInNjYW5faWQiOiAi\nVlVMTl9zY2FuXzQ1NmRlZiIsCiAgICAgICJ0aW1lc3RhbXAiOiAiMjAyNS0w\nMS0yOVQwOToxNTowMFoiLAogICAgICAic2Nhbl9zdGF0dXMiOiAiY29tcGxl\ndGVkIiwKICAgICAgInJlc291cmNlX3R5cGUiOiAiY29udGFpbmVyIiwKICAg\nICAgInJlc291cmNlX25hbWUiOiAiYXV0aC1zZXJ2aWNlOjIuMS4wIiwKICAg\nICAgInZ1bG5lcmFiaWxpdGllcyI6IFsKICAgICAgICB7CiAgICAgICAgICAi\naWQiOiAiQ1ZFLTIwMjQtMjIyMiIsCiAgICAgICAgICAic2V2ZXJpdHkiOiAi\nSElHSCIsCiAgICAgICAgICAiY3ZzcyI6IDguMiwKICAgICAgICAgICJzdGF0\ndXMiOiAiYWN0aXZlIiwKICAgICAgICAgICJwYWNrYWdlX25hbWUiOiAic3By\naW5nLXNlY3VyaXR5IiwKICAgICAgICAgICJjdXJyZW50X3ZlcnNpb24iOiAi\nNS42LjAiLAogICAgICAgICAgImZpeGVkX3ZlcnNpb24iOiAiNS42LjEiLAog\nICAgICAgICAgImRlc2NyaXB0aW9uIjogIkF1dGhlbnRpY2F0aW9uIGJ5cGFz\ncyBpbiBTcHJpbmcgU2VjdXJpdHkiLAogICAgICAgICAgInB1Ymxpc2hlZF9k\nYXRlIjogIjIwMjUtMDEtMjdUMDA6MDA6MDBaIiwKICAgICAgICAgICJsaW5r\nIjogImh0dHBzOi8vbnZkLm5pc3QuZ292L3Z1bG4vZGV0YWlsL0NWRS0yMDI0\nLTIyMjIiLAogICAgICAgICAgInJpc2tfZmFjdG9ycyI6IFsKICAgICAgICAg\nICAgIkF1dGhlbnRpY2F0aW9uIEJ5cGFzcyIsCiAgICAgICAgICAgICJIaWdo\nIENWU1MgU2NvcmUiLAogICAgICAgICAgICAiUHJvb2Ygb2YgQ29uY2VwdCBF\neHBsb2l0IEF2YWlsYWJsZSIKICAgICAgICAgIF0KICAgICAgICB9LAogICAg\nICAgIHsKICAgICAgICAgICJpZCI6ICJDVkUtMjAyNC0yMjIzIiwKICAgICAg\nICAgICJzZXZlcml0eSI6ICJNRURJVU0iLAogICAgICAgICAgImN2c3MiOiA2\nLjUsCiAgICAgICAgICAic3RhdHVzIjogImFjdGl2ZSIsCiAgICAgICAgICAi\ncGFja2FnZV9uYW1lIjogInRvbWNhdCIsCiAgICAgICAgICAiY3VycmVudF92\nZXJzaW9uIjogIjkuMC41MCIsCiAgICAgICAgICAiZml4ZWRfdmVyc2lvbiI6\nICI5LjAuNTEiLAogICAgICAgICAgImRlc2NyaXB0aW9uIjogIkluZm9ybWF0\naW9uIGRpc2Nsb3N1cmUgaW4gQXBhY2hlIFRvbWNhdCIsCiAgICAgICAgICAi\ncHVibGlzaGVkX2RhdGUiOiAiMjAyNS0wMS0yOFQwMDowMDowMFoiLAogICAg\nICAgICAgImxpbmsiOiAiaHR0cHM6Ly9udmQubmlzdC5nb3YvdnVsbi9kZXRh\naWwvQ1ZFLTIwMjQtMjIyMyIsCiAgICAgICAgICAicmlza19mYWN0b3JzIjog\nWwogICAgICAgICAgICAiSW5mb3JtYXRpb24gRGlzY2xvc3VyZSIsCiAgICAg\nICAgICAgICJNZWRpdW0gQ1ZTUyBTY29yZSIKICAgICAgICAgIF0KICAgICAg\nICB9CiAgICAgIF0sCiAgICAgICJzdW1tYXJ5IjogewogICAgICAgICJ0b3Rh\nbF92dWxuZXJhYmlsaXRpZXMiOiAyLAogICAgICAgICJzZXZlcml0eV9jb3Vu\ndHMiOiB7CiAgICAgICAgICAiQ1JJVElDQUwiOiAwLAogICAgICAgICAgIkhJ\nR0giOiAxLAogICAgICAgICAgIk1FRElVTSI6IDEsCiAgICAgICAgICAiTE9X\nIjogMAogICAgICAgIH0sCiAgICAgICAgImZpeGFibGVfY291bnQiOiAyLAog\nICAgICAgICJjb21wbGlhbnQiOiBmYWxzZQogICAgICB9LAogICAgICAic2Nh\nbl9tZXRhZGF0YSI6IHsKICAgICAgICAic2Nhbm5lcl92ZXJzaW9uIjogIjMw\nLjEuNTEiLAogICAgICAgICJwb2xpY2llc192ZXJzaW9uIjogIjIwMjUuMS4y\nOSIsCiAgICAgICAgInNjYW5uaW5nX3J1bGVzIjogWwogICAgICAgICAgInZ1\nbG5lcmFiaWxpdHkiLAogICAgICAgICAgImNvbXBsaWFuY2UiLAogICAgICAg\nICAgIm1hbHdhcmUiCiAgICAgICAgXSwKICAgICAgICAiZXhjbHVkZWRfcGF0\naHMiOiBbCiAgICAgICAgICAiL3RtcCIsCiAgICAgICAgICAiL3Zhci9sb2ci\nCiAgICAgICAgXQogICAgICB9CiAgICB9CiAgfSwKICB7CiAgICAic2NhblJl\nc3VsdHMiOiB7CiAgICAgICJzY2FuX2lkIjogIlZVTE5fc2Nhbl8xMjNhYmMi\nLAogICAgICAidGltZXN0YW1wIjogIjIwMjUtMDEtMjlUMDg6MDA6MDBaIiwK\nICAgICAgInNjYW5fc3RhdHVzIjogImNvbXBsZXRlZCIsCiAgICAgICJyZXNv\ndXJjZV90eXBlIjogImNvbnRhaW5lciIsCiAgICAgICJyZXNvdXJjZV9uYW1l\nIjogInBheW1lbnQtcHJvY2Vzc29yOjEuMC4wIiwKICAgICAgInZ1bG5lcmFi\naWxpdGllcyI6IFsKICAgICAgICB7CiAgICAgICAgICAiaWQiOiAiQ1ZFLTIw\nMjQtMTExMSIsCiAgICAgICAgICAic2V2ZXJpdHkiOiAiQ1JJVElDQUwiLAog\nICAgICAgICAgImN2c3MiOiA5LjksCiAgICAgICAgICAic3RhdHVzIjogImFj\ndGl2ZSIsCiAgICAgICAgICAicGFja2FnZV9uYW1lIjogIm9wZW5zc2wiLAog\nICAgICAgICAgImN1cnJlbnRfdmVyc2lvbiI6ICIzLjAuMCIsCiAgICAgICAg\nICAiZml4ZWRfdmVyc2lvbiI6ICIzLjAuMSIsCiAgICAgICAgICAiZGVzY3Jp\ncHRpb24iOiAiQ3JpdGljYWwgYnVmZmVyIG92ZXJmbG93IGluIE9wZW5TU0wg\nVExTIGhhbmRsaW5nIiwKICAgICAgICAgICJwdWJsaXNoZWRfZGF0ZSI6ICIy\nMDI1LTAxLTI4VDAwOjAwOjAwWiIsCiAgICAgICAgICAibGluayI6ICJodHRw\nczovL252ZC5uaXN0Lmdvdi92dWxuL2RldGFpbC9DVkUtMjAyNC0xMTExIiwK\nICAgICAgICAgICJyaXNrX2ZhY3RvcnMiOiBbCiAgICAgICAgICAgICJCdWZm\nZXIgT3ZlcmZsb3ciLAogICAgICAgICAgICAiQ3JpdGljYWwgQ1ZTUyBTY29y\nZSIsCiAgICAgICAgICAgICJQdWJsaWMgRXhwbG9pdCBBdmFpbGFibGUiLAog\nICAgICAgICAgICAiRXhwbG9pdCBpbiBXaWxkIgogICAgICAgICAgXQogICAg\nICAgIH0KICAgICAgXSwKICAgICAgInN1bW1hcnkiOiB7CiAgICAgICAgInRv\ndGFsX3Z1bG5lcmFiaWxpdGllcyI6IDEsCiAgICAgICAgInNldmVyaXR5X2Nv\ndW50cyI6IHsKICAgICAgICAgICJDUklUSUNBTCI6IDEsCiAgICAgICAgICAi\nSElHSCI6IDAsCiAgICAgICAgICAiTUVESVVNIjogMCwKICAgICAgICAgICJM\nT1ciOiAwCiAgICAgICAgfSwKICAgICAgICAiZml4YWJsZV9jb3VudCI6IDEs\nCiAgICAgICAgImNvbXBsaWFudCI6IGZhbHNlCiAgICAgIH0sCiAgICAgICJz\nY2FuX21ldGFkYXRhIjogewogICAgICAgICJzY2FubmVyX3ZlcnNpb24iOiAi\nMzAuMS41MSIsCiAgICAgICAgInBvbGljaWVzX3ZlcnNpb24iOiAiMjAyNS4x\nLjI5IiwKICAgICAgICAic2Nhbm5pbmdfcnVsZXMiOiBbCiAgICAgICAgICAi\ndnVsbmVyYWJpbGl0eSIsCiAgICAgICAgICAiY29tcGxpYW5jZSIsCiAgICAg\nICAgICAibWFsd2FyZSIKICAgICAgICBdLAogICAgICAgICJleGNsdWRlZF9w\nYXRocyI6IFsKICAgICAgICAgICIvdG1wIiwKICAgICAgICAgICIvdmFyL2xv\nZyIKICAgICAgICBdCiAgICAgIH0KICAgIH0KICB9Cl0=\n",
				Encoding: "base64",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileData, err := getDataFromGitHub(tt.args.baseUrl, tt.args.file, tt.args.githubReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDataFromGitHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFileData, tt.wantFileData) {
				t.Errorf("getDataFromGitHub() = %v, want %v", gotFileData, tt.wantFileData)
			}
		})
	}
}
