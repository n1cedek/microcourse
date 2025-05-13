package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"microservices_course/week5/internal/model"
	repository "microservices_course/week5/internal/repo"
	repoMock "microservices_course/week5/internal/repo/mocks"
	"microservices_course/week5/internal/service/note"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type noteRepoMockFunc func(mc *minimock.Controller) repository.NoteRepo

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		title     = gofakeit.Name()
		content   = gofakeit.Name()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		res = &model.Note{
			ID: id,
			Info: model.NoteInfo{
				Title:   title,
				Content: content,
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         *model.Note
		err          error
		noteRepoMock noteRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			noteRepoMock: func(mc *minimock.Controller) repository.NoteRepo {
				mock := repoMock.NewNoteRepoMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		}, {
			name: "error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			noteRepoMock: func(mc *minimock.Controller) repository.NoteRepo {
				mock := repoMock.NewNoteRepoMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteRepoMock := tt.noteRepoMock(mc)
			service := notes.NewMockService(noteRepoMock)

			result, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, result)

		})
	}
}
