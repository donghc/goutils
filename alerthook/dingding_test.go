package alerthook

import "testing"

func TestDingDing_SendDingDingMessage(t *testing.T) {
	type fields struct {
		Security string
		WebHook  string
		KeyWords string
	}
	type args struct {
		contentData string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &DingDing{
				Security: tt.fields.Security,
				WebHook:  tt.fields.WebHook,
				KeyWords: tt.fields.KeyWords,
			}
			if got := dd.SendDingDingMessage(tt.args.contentData); got != tt.want {
				t.Errorf("SendDingDingMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
