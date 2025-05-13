package tests

import (
	"context"
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

func TestCreate(t *testing.T) {
	t.Parallel()
	type noteRepoMockFunc func(mc *minimock.Controller) repository.NoteRepo

	type args struct {
		ctx context.Context
		req *model.NoteInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		title   = gofakeit.Name()
		content = gofakeit.Name()

		repoErr = fmt.Errorf("repo error")

		req = &model.NoteInfo{
			Title:   title,
			Content: content,
		}
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		noteRepoMock noteRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			noteRepoMock: func(mc *minimock.Controller) repository.NoteRepo {
				mock := repoMock.NewNoteRepoMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		}, {
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			noteRepoMock: func(mc *minimock.Controller) repository.NoteRepo {
				mock := repoMock.NewNoteRepoMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
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

			result, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, result)

		})
	}
}
