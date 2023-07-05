package handlers

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	entitiesHttp "github.com/akrillis/nakamatest/internal/entities/http"
	"github.com/akrillis/nakamatest/internal/service"
	"github.com/akrillis/nakamatest/internal/storage"
)

func Test_getContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       entitiesHttp.FileRequest
		want      string
		wantErr   bool
		setupMock func(
			svc *service.MockFiler,
			storage *storage.MockSaver,
		)
	}{
		{
			name:    "content getting error",
			req:     entitiesHttp.FileRequest{},
			want:    "",
			wantErr: true,
			setupMock: func(
				svc *service.MockFiler,
				storage *storage.MockSaver,
			) {
				svc.EXPECT().ReadWithCheck().Return("", os.ErrNotExist)
			},
		},
		{
			name: "ok with saver error",
			req: entitiesHttp.FileRequest{
				Type:    "type",
				Version: "version",
				Hash:    "hash",
			},
			want:    `{"field0": "value0", "field1": 1}`,
			wantErr: false,
			setupMock: func(
				svc *service.MockFiler,
				storage *storage.MockSaver,
			) {
				svc.EXPECT().ReadWithCheck().Return(`{"field0": "value0", "field1": 1}`, nil)

				svc.EXPECT().GetPath().Return("type/version.json").Times(2)

				storage.
					EXPECT().
					Save("type/version.json", `{"field0": "value0", "field1": 1}`).
					Return(errors.New("error"))
			},
		},
		{
			name: "ok",
			req: entitiesHttp.FileRequest{
				Type:    "type",
				Version: "version",
				Hash:    "hash",
			},
			want:    `{"field0": "value0", "field1": 1}`,
			wantErr: false,
			setupMock: func(
				svc *service.MockFiler,
				storage *storage.MockSaver,
			) {
				svc.EXPECT().ReadWithCheck().Return(`{"field0": "value0", "field1": 1}`, nil)

				svc.EXPECT().GetPath().Return("type/version.json")

				storage.
					EXPECT().
					Save("type/version.json", `{"field0": "value0", "field1": 1}`).
					Return(nil)
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			filer := service.NewMockFiler(ctrl)
			saver := storage.NewMockSaver(ctrl)

			test.setupMock(filer, saver)

			content, err := getContent(filer, saver)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.want, content)
			}
		})
	}
}
