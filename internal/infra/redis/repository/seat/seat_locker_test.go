package seatrepo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domaincache "ticket-reservation/internal/domain/cache"
	seatrepo "ticket-reservation/internal/infra/redis/repository/seat"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
	lockmanager "github.com/kittipat1413/go-common/framework/lockmanager"
	locker_mocks "github.com/kittipat1413/go-common/framework/lockmanager/mocks"
)

func TestNewSeatLockerRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLocker := locker_mocks.NewMockLockManager(ctrl)

	// Execute
	repo := seatrepo.NewSeatLockerRepository(mockLocker)

	// Assert
	assert.NotNil(t, repo)
}

func TestSeatLockerRepositoryImpl_LockSeat(t *testing.T) {
	concertID := uuid.New()
	zoneID := uuid.New()
	seatID := uuid.New()
	token := "test-token-123"
	ttl := 5 * time.Minute
	expectedKey := "seat_lock:concert:" + concertID.String() + ":zone:" + zoneID.String() + ":seat:" + seatID.String()

	tests := []struct {
		name              string
		setupMock         func(mock *locker_mocks.MockLockManager)
		expectedError     bool
		expectedErrorMsg  string
		expectedErrorType error
	}{
		{
			name: "successful lock acquisition",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Acquire(gomock.Any(), expectedKey, ttl, token).
					Return(token, nil)
			},
			expectedError: false,
		},
		{
			name: "seat already locked by another process",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Acquire(gomock.Any(), expectedKey, ttl, token).
					Return("", lockmanager.ErrLockAlreadyTaken)
			},
			expectedError:     true,
			expectedErrorMsg:  "seat already locked",
			expectedErrorType: domaincache.ErrSeatAlreadyLocked,
		},
		{
			name: "lock manager connection error",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Acquire(gomock.Any(), expectedKey, ttl, token).
					Return("", errors.New("redis connection failed"))
			},
			expectedError:     true,
			expectedErrorMsg:  "failed to lock seat",
			expectedErrorType: &errsFramework.DatabaseError{},
		},
		{
			name: "lock manager timeout error",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Acquire(gomock.Any(), expectedKey, ttl, token).
					Return("", context.DeadlineExceeded)
			},
			expectedError:     true,
			expectedErrorMsg:  "failed to lock seat",
			expectedErrorType: &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLocker := locker_mocks.NewMockLockManager(ctrl)
			tt.setupMock(mockLocker)

			repository := seatrepo.NewSeatLockerRepository(mockLocker)

			// Execute
			err := repository.LockSeat(context.Background(), concertID, zoneID, seatID, token, ttl)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[repository seat/seat_locker LockSeat]")
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)

				if tt.expectedErrorType != nil {
					assert.ErrorAs(t, err, &tt.expectedErrorType, "Expected error to be of type %T", tt.expectedErrorType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSeatLockerRepositoryImpl_UnlockSeat(t *testing.T) {
	concertID := uuid.New()
	zoneID := uuid.New()
	seatID := uuid.New()
	token := "test-token-123"
	expectedKey := "seat_lock:concert:" + concertID.String() + ":zone:" + zoneID.String() + ":seat:" + seatID.String()

	tests := []struct {
		name              string
		setupMock         func(mock *locker_mocks.MockLockManager)
		expectedError     bool
		expectedErrorMsg  string
		expectedErrorType error
	}{
		{
			name: "successful unlock",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Release(gomock.Any(), expectedKey, token).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "unlock denied - token mismatch",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Release(gomock.Any(), expectedKey, token).
					Return(lockmanager.ErrUnlockNotPermitted)
			},
			expectedError:     true,
			expectedErrorMsg:  "seat unlock denied",
			expectedErrorType: domaincache.ErrSeatUnlockDenied,
		},
		{
			name: "lock manager connection error",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Release(gomock.Any(), expectedKey, token).
					Return(errors.New("redis connection failed"))
			},
			expectedError:     true,
			expectedErrorMsg:  "failed to unlock seat",
			expectedErrorType: &errsFramework.DatabaseError{},
		},
		{
			name: "lock manager timeout error",
			setupMock: func(mock *locker_mocks.MockLockManager) {
				mock.EXPECT().
					Release(gomock.Any(), expectedKey, token).
					Return(context.DeadlineExceeded)
			},
			expectedError:     true,
			expectedErrorMsg:  "failed to unlock seat",
			expectedErrorType: &errsFramework.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLocker := locker_mocks.NewMockLockManager(ctrl)
			tt.setupMock(mockLocker)

			repository := seatrepo.NewSeatLockerRepository(mockLocker)

			// Execute
			err := repository.UnlockSeat(context.Background(), concertID, zoneID, seatID, token)

			// Assert
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "[repository seat/seat_locker UnlockSeat]")
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)

				if tt.expectedErrorType != nil {
					assert.ErrorAs(t, err, &tt.expectedErrorType, "Expected error to be of type %T", tt.expectedErrorType)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
