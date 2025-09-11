package urltool

import "testing"

func TestGetBasePath(t *testing.T) {
	type args struct {
		targetUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				targetUrl: "https://www.liwenzhou.com/posts/Go/gRPC/",
			},
			want:    "gRPC",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				targetUrl: "/xxx/123",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBasePath(tt.args.targetUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
