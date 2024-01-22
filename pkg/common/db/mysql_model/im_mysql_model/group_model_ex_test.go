package im_mysql_model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestGetArticleList(t *testing.T) {
	type args struct {
		groupID  string
		fromID   int64
		offset   int64
		PageSize int64
	}
	tests := []struct {
		name           string
		args           args
		wantResultData []proto.Message
		wantErr        assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResultData, err := GetArticleList(tt.args.groupID, tt.args.fromID, tt.args.offset, tt.args.PageSize)
			if !tt.wantErr(t, err, fmt.Sprintf("GetArticleList(%v, %v, %v, %v)", tt.args.groupID, tt.args.fromID, tt.args.offset, tt.args.PageSize)) {
				return
			}
			fmt.Println(gotResultData)
			assert.Equalf(t, tt.wantResultData, gotResultData, "GetArticleList(%v, %v, %v, %v)", tt.args.groupID, tt.args.fromID, tt.args.offset, tt.args.PageSize)
		})
	}
}
