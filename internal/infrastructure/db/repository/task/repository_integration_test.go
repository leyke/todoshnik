//go:build integration

package task_test

import (
	"context"
	"testing"

	"todoshnik/internal/infrastructure/db/repository/task"
	"todoshnik/internal/infrastructure/identity"

	apptask "todoshnik/internal/domains/task"
	taskerror "todoshnik/internal/domains/task/errors"
	testutils "todoshnik/internal/infrastructure/utils/test"

	"github.com/stretchr/testify/require"
)

func TestDBRepository_GetByID(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	tests := []struct {
		name      string
		id        int
		scope     identity.AccessScope
		want      *TaskFixture
		wantError error
	}{
		{
			name: "есть у пользователя",
			id:   fixtures[0].ID,
			scope: identity.AccessScope{
				UserID: fixtures[0].UserID,
			},
			want: &fixtures[0],
		},
		{
			name:      "нет в бд",
			id:        999999,
			wantError: taskerror.ErrNotFound,
		},
		{
			name: "чужой пользователь",
			id:   fixtures[3].ID,
			scope: identity.AccessScope{
				UserID: fixtures[0].UserID,
			},
			wantError: taskerror.ErrNotFound,
		},
		{
			name: "есть доступ у админа",
			id:   fixtures[0].ID,
			scope: identity.AccessScope{
				IsAdmin: true,
			},
			want: &fixtures[0],
		},
		{
			name: "админ видит чужую задачу",
			id:   fixtures[3].ID,
			scope: identity.AccessScope{
				IsAdmin: true,
			},
			want: &fixtures[3],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := repo.GetByID(ctx, tt.id, tt.scope)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				require.Nil(t, task)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, task)

			require.Equal(t, tt.want.ID, task.ID)
			require.Equal(t, tt.want.Title, task.Title)
			require.Equal(t, tt.want.UserID, task.UserID)
			require.Equal(t, tt.want.Done, task.Done)
		})
	}
}

func TestDBRepository_GetByID_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.GetByID(ctx, 99999, identity.AccessScope{})

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_List_User(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	result, err := repo.List(ctx, apptask.TaskFilter{
		Scope: identity.AccessScope{
			UserID: fixtures[0].UserID,
		},
	})

	require.NoError(t, err)
	require.Len(t, result, 3)

	require.Equal(t, fixtures[1].ID, result[0].ID)
	require.Equal(t, fixtures[1].Title, result[0].Title)
	require.True(t, result[0].Done)
	require.Equal(t, fixtures[1].UserID, result[0].UserID)

	require.Equal(t, fixtures[0].ID, result[1].ID)
	require.Equal(t, fixtures[0].Title, result[1].Title)
	require.False(t, result[1].Done)
	require.Equal(t, fixtures[0].UserID, result[1].UserID)

	require.Equal(t, fixtures[2].ID, result[2].ID)
	require.Equal(t, fixtures[2].Title, result[2].Title)
	require.False(t, result[2].Done)
	require.Equal(t, fixtures[2].UserID, result[2].UserID)
}

func TestDBRepository_List_Admin(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	result, err := repo.List(ctx, apptask.TaskFilter{
		Scope: identity.AccessScope{
			IsAdmin: true,
		},
	})

	require.NoError(t, err)
	require.Len(t, result, 5)

	require.Equal(t, fixtures[1].ID, result[0].ID)
	require.Equal(t, fixtures[4].ID, result[1].ID)
	require.Equal(t, fixtures[0].ID, result[2].ID)
	require.Equal(t, fixtures[2].ID, result[3].ID)
	require.Equal(t, fixtures[3].ID, result[4].ID)
}

func TestDBRepository_List_StatusFilter(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	_ = ApplyFixtures(t, db)

	tests := []struct {
		name      string
		status    apptask.Status
		wantError error
		wantLen   int
	}{
		{
			name:    "только в процессе",
			status:  apptask.StatusPending,
			wantLen: 3,
		},
		{
			name:    "только заверщенные",
			status:  apptask.StatusCompleted,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := repo.List(ctx, apptask.TaskFilter{
				Scope:  identity.AccessScope{IsAdmin: true},
				Status: tt.status,
			})

			require.NoError(t, err)
			require.Len(t, tasks, tt.wantLen)
			for _, task := range tasks {
				if tt.status == apptask.StatusCompleted {
					require.True(t, task.Done)
				} else {
					require.False(t, task.Done)
				}
			}
		})
	}
}

func TestDBRepository_List_Empty(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	result, err := repo.List(ctx, apptask.TaskFilter{})

	require.NoError(t, err)
	require.Empty(t, result)
}

func TestDBRepository_List_DbError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	result, err := repo.List(ctx, apptask.TaskFilter{})

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Create_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	tests := []struct {
		name  string
		input *apptask.Task
	}{
		{
			name: "новый готовый",
			input: &apptask.Task{
				Title:  "New Done Task",
				Done:   true,
				UserID: fixtures[0].UserID,
			},
		},
		{
			name: "новый в процессе",
			input: &apptask.Task{
				Title:  "New Pending Task",
				Done:   false,
				UserID: fixtures[0].UserID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(ctx, tt.input)

			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotZero(t, result.ID)

			saved, err := repo.GetByID(
				ctx,
				result.ID,
				identity.AccessScope{
					UserID: result.UserID,
				},
			)

			require.NoError(t, err)
			require.NotNil(t, saved)
			require.Equal(t, tt.input.Title, saved.Title)
			require.Equal(t, tt.input.Done, saved.Done)
			require.Equal(t, tt.input.UserID, saved.UserID)
		})
	}
}

func TestDBRepository_Create_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &apptask.Task{
		Title: "New Pending Task",
		Done:  false,
	}

	result, err := repo.Create(ctx, input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestDBRepository_Update_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	tests := []struct {
		name  string
		input *apptask.Task
	}{
		{
			name: "Сменить статус и название",
			input: &apptask.Task{
				ID:     fixtures[0].ID,
				Title:  "Done Task UPDed to Pending",
				Done:   false,
				UserID: fixtures[0].UserID,
			},
		},
		{
			name: "Поставить в готово",
			input: &apptask.Task{
				ID:     fixtures[1].ID,
				Title:  fixtures[1].Title,
				Done:   true,
				UserID: fixtures[1].UserID,
			},
		},
		{
			name: "Поручить",
			input: &apptask.Task{
				ID:     fixtures[2].ID,
				Title:  fixtures[2].Title,
				Done:   fixtures[2].Done,
				UserID: fixtures[3].UserID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.input)

			require.NoError(t, err)

			result, err := repo.GetByID(
				ctx,
				tt.input.ID,
				identity.AccessScope{UserID: tt.input.UserID},
			)

			require.NoError(t, err)
			require.NotNil(t, result)

			require.Equal(t, tt.input.UserID, result.UserID)
			require.Equal(t, tt.input.Title, result.Title)
			require.Equal(t, tt.input.Done, result.Done)
		})
	}
}

func TestDBRepository_Update_NotFound(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	input := &apptask.Task{
		ID:    999999,
		Title: "unknown Task UPDed to Done",
	}

	err := repo.Update(ctx, input)

	require.ErrorIs(t, err, taskerror.ErrNotFound)
}

func TestDBRepository_Update_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &apptask.Task{
		ID:    999999,
		Title: "unknown Task UPDed to Done",
	}

	err := repo.Update(ctx, input)

	require.Error(t, err)
}

func TestDBRepository_Delete_Success(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	fixtures := ApplyFixtures(t, db)

	input := &apptask.Task{
		ID:     fixtures[0].ID,
		UserID: fixtures[0].UserID,
	}

	err := repo.Delete(ctx, input)
	require.NoError(t, err)

	result, err := repo.GetByID(ctx, input.ID, identity.AccessScope{UserID: input.UserID})

	require.ErrorIs(t, err, taskerror.ErrNotFound)
	require.Nil(t, result)
}

func TestDBRepository_Delete_NotFound(t *testing.T) {
	db := testutils.TestDBWithCleanup(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	input := &apptask.Task{
		ID: 999999,
	}

	err := repo.Delete(ctx, input)

	require.ErrorIs(t, err, taskerror.ErrNotFound)
}

func TestDBRepository_Delete_DBError(t *testing.T) {
	db := testutils.TestDB(t)
	repo := task.NewRepository(db)
	ctx := context.Background()

	require.NoError(t, db.Close())

	input := &apptask.Task{
		ID: 1,
	}

	err := repo.Delete(ctx, input)

	require.Error(t, err)
}
