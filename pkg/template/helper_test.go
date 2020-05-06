package template

import (
	"github.com/golang/mock/gomock"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/repository"
	mock_repository "github.com/pxecore/pxecore/pkg/repository/mock"
	"reflect"
	"testing"
)

func Test_mergeMaps(t *testing.T) {
	type args struct {
		m1 map[string]string
		m2 map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"M2_OVERRIDES_M1", args{
			map[string]string{"a": "m1", "b": "m1"},
			map[string]string{"a": "m2", "c": "m2"}},
			map[string]string{"a": "m2", "b": "m1", "c": "m2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeMaps(tt.args.m1, tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recursiveGroupMerge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	gr := mock_repository.NewMockGroupRepository(ctrl)
	gr.EXPECT().Get("recursive").Return(entity.Group{
		ID: "recursive", Vars: map[string]string{},
		ParentID: "recursive", TemplateID: "template1",
	}, nil).AnyTimes()
	gr.EXPECT().Get("children").Return(entity.Group{
		ID: "children", Vars: map[string]string{"a": "children"},
		ParentID: "parent", TemplateID: "children",
	}, nil).AnyTimes()
	gr.EXPECT().Get("parent").Return(entity.Group{
		ID: "parent", Vars: map[string]string{"a": "parent"},
		ParentID: "", TemplateID: "parent",
	}, nil).AnyTimes()
	ms := mock_repository.NewMockSession(ctrl)
	ms.EXPECT().Group().Return(gr).AnyTimes()

	type args struct {
		session repository.Session
		groups  []string
		groupID string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   string
		wantErr bool
	}{
		{"OK",
			args{ms, []string{}, "children"},
			map[string]string{"a": "children"}, "children", false},
		{"KO_RECURSIVE",
			args{ms, []string{}, "recursive"},
			nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := recursiveGroupMerge(tt.args.session, tt.args.groups, tt.args.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("recursiveGroupMerge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("recursiveGroupMerge() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("recursiveGroupMerge() got1 = %v, want %v", got1, tt.want1)
				}
			}
		})
	}
}
