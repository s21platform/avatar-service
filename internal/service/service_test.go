package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/model"
	"github.com/s21platform/avatar-service/pkg/avatar"
)

func TestService_SetUserAvatar(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
	ctx = context.WithValue(ctx, config.KeyUUID, "test-uuid")

	mockS3 := NewMockS3Storage(ctrl)
	mockRepo := NewMockDBRepo(ctrl)
	mockUserKafka := NewMockKafkaProducer(ctrl)
	mockSocietyKafka := NewMockKafkaProducer(ctrl)

	t.Run("set_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetUserAvatar")

		stream := NewMockAvatarService_SetUserAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetUserAvatarIn{
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
		stream.EXPECT().SendAndClose(&avatar.SetUserAvatarOut{Link: "https://s3.example.com/avatar.jpg"}).Return(nil)

		mockS3.EXPECT().PutObject(gomock.Any(), &model.AvatarContent{
			AvatarType: model.UserAvatarType,
			UUID:       "test-uuid",
			Filename:   "test.jpg",
			ImageData:  []byte{1, 2, 3},
		}).Return("https://s3.example.com/avatar.jpg", nil)

		mockRepo.EXPECT().SetUserAvatar(gomock.Any(), "test-uuid", "https://s3.example.com/avatar.jpg").Return(nil)
		mockUserKafka.EXPECT().ProduceMessage(gomock.Any(), &avatar.NewAvatarRegister{
			Uuid: "test-uuid",
			Link: "https://s3.example.com/avatar.jpg",
		}, "test-uuid").Return(nil)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetUserAvatar(stream)

		assert.NoError(t, err)
	})

	t.Run("no_uuid", func(t *testing.T) {
		ctxNoUUID := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		mockLogger.EXPECT().AddFuncName("SetUserAvatar")
		mockLogger.EXPECT().Error("uuid is required")

		stream := NewMockAvatarService_SetUserAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctxNoUUID).AnyTimes()

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetUserAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "uuid is required")
	})

	t.Run("stream_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetUserAvatar")
		mockLogger.EXPECT().Error("failed to receive data from stream: stream error")

		stream := NewMockAvatarService_SetUserAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(nil, errors.New("stream error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetUserAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to receive data from stream")
	})

	t.Run("s3_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetUserAvatar")
		mockLogger.EXPECT().Error("failed to upload file to S3: s3 error")

		stream := NewMockAvatarService_SetUserAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetUserAvatarIn{
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		mockS3.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return("", errors.New("s3 error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetUserAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to upload file to S3")
	})

	t.Run("repo_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetUserAvatar")
		mockLogger.EXPECT().Error("failed to save avatar to database: repo error")

		stream := NewMockAvatarService_SetUserAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetUserAvatarIn{
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		mockS3.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return("https://s3.example.com/avatar.jpg", nil)
		mockRepo.EXPECT().SetUserAvatar(gomock.Any(), "test-uuid", "https://s3.example.com/avatar.jpg").Return(errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetUserAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to save avatar to database")
	})

	t.Run("kafka_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetUserAvatar")
		mockLogger.EXPECT().Error("failed to produce message to user service: kafka error")

		stream := NewMockAvatarService_SetUserAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetUserAvatarIn{
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		mockS3.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return("https://s3.example.com/avatar.jpg", nil)
		mockRepo.EXPECT().SetUserAvatar(gomock.Any(), "test-uuid", "https://s3.example.com/avatar.jpg").Return(nil)
		mockUserKafka.EXPECT().ProduceMessage(gomock.Any(), gomock.Any(), "test-uuid").Return(errors.New("kafka error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetUserAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to produce message to user service")
	})
}

func TestService_GetAllUserAvatars(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
	ctx = context.WithValue(ctx, config.KeyUUID, "test-uuid")

	mockS3 := NewMockS3Storage(ctrl)
	mockRepo := NewMockDBRepo(ctrl)
	mockUserKafka := NewMockKafkaProducer(ctrl)
	mockSocietyKafka := NewMockKafkaProducer(ctrl)

	t.Run("get_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllUserAvatars")

		expectedAvatars := &model.AvatarMetadataList{
			{ID: 1, Link: "https://s3.example.com/avatar1.jpg"},
			{ID: 2, Link: "https://s3.example.com/avatar2.jpg"},
		}

		mockRepo.EXPECT().GetAllUserAvatars(gomock.Any(), "test-uuid").Return(expectedAvatars, nil)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		result, err := s.GetAllUserAvatars(ctx, &emptypb.Empty{})

		assert.NoError(t, err)
		assert.Len(t, result.AvatarList, 2)
		assert.Equal(t, int32(1), result.AvatarList[0].Id)
		assert.Equal(t, "https://s3.example.com/avatar1.jpg", result.AvatarList[0].Link)
		assert.Equal(t, int32(2), result.AvatarList[1].Id)
		assert.Equal(t, "https://s3.example.com/avatar2.jpg", result.AvatarList[1].Link)
	})

	t.Run("no_uuid", func(t *testing.T) {
		ctxNoUUID := context.WithValue(context.Background(), config.KeyLogger, mockLogger)
		mockLogger.EXPECT().AddFuncName("GetAllUserAvatars")
		mockLogger.EXPECT().Error("uuid is required")

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.GetAllUserAvatars(ctxNoUUID, &emptypb.Empty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "uuid is required")
	})

	t.Run("repo_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllUserAvatars")
		mockLogger.EXPECT().Error("failed to get all user avatars: repo error")

		mockRepo.EXPECT().GetAllUserAvatars(gomock.Any(), "test-uuid").Return(nil, errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.GetAllUserAvatars(ctx, &emptypb.Empty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to get all user avatars")
	})
}

func TestService_DeleteUserAvatar(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockS3 := NewMockS3Storage(ctrl)
	mockRepo := NewMockDBRepo(ctrl)
	mockUserKafka := NewMockKafkaProducer(ctrl)
	mockSocietyKafka := NewMockKafkaProducer(ctrl)

	t.Run("delete_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteUserAvatar")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "test-uuid",
			Link: "https://s3.example.com/avatar.jpg",
		}

		mockRepo.EXPECT().GetUserAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/avatar.jpg").Return(nil)
		mockRepo.EXPECT().DeleteUserAvatar(gomock.Any(), 1).Return(nil)
		mockRepo.EXPECT().GetLatestUserAvatar(gomock.Any(), "test-uuid").Return("https://s3.example.com/latest.jpg")
		mockUserKafka.EXPECT().ProduceMessage(gomock.Any(), &avatar.NewAvatarRegister{
			Uuid: "test-uuid",
			Link: "https://s3.example.com/latest.jpg",
		}, "test-uuid").Return(nil)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		result, err := s.DeleteUserAvatar(ctx, &avatar.DeleteUserAvatarIn{AvatarId: 1})

		assert.NoError(t, err)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "https://s3.example.com/avatar.jpg", result.Link)
	})

	t.Run("get_avatar_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteUserAvatar")
		mockLogger.EXPECT().Error("failed to get user avatar: repo error")

		mockRepo.EXPECT().GetUserAvatarData(gomock.Any(), 1).Return(nil, errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteUserAvatar(ctx, &avatar.DeleteUserAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "failed to get avatar data")
	})

	t.Run("s3_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteUserAvatar")
		mockLogger.EXPECT().Error("failed to delete user avatar: s3 error")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "test-uuid",
			Link: "https://s3.example.com/avatar.jpg",
		}

		mockRepo.EXPECT().GetUserAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/avatar.jpg").Return(errors.New("s3 error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteUserAvatar(ctx, &avatar.DeleteUserAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to delete avatar in s3")
	})

	t.Run("repo_delete_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteUserAvatar")
		mockLogger.EXPECT().Error("failed to delete user avatar: repo error")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "test-uuid",
			Link: "https://s3.example.com/avatar.jpg",
		}

		mockRepo.EXPECT().GetUserAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/avatar.jpg").Return(nil)
		mockRepo.EXPECT().DeleteUserAvatar(gomock.Any(), 1).Return(errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteUserAvatar(ctx, &avatar.DeleteUserAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to delete avatar in postgres")
	})

	t.Run("kafka_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteUserAvatar")
		mockLogger.EXPECT().Error("failed to produce avatar to user service: kafka error")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "test-uuid",
			Link: "https://s3.example.com/avatar.jpg",
		}

		mockRepo.EXPECT().GetUserAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/avatar.jpg").Return(nil)
		mockRepo.EXPECT().DeleteUserAvatar(gomock.Any(), 1).Return(nil)
		mockRepo.EXPECT().GetLatestUserAvatar(gomock.Any(), "test-uuid").Return("https://s3.example.com/latest.jpg")
		mockUserKafka.EXPECT().ProduceMessage(gomock.Any(), gomock.Any(), "test-uuid").Return(errors.New("kafka error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteUserAvatar(ctx, &avatar.DeleteUserAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to produce avatar to user service")
	})
}

func TestService_SetSocietyAvatar(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockS3 := NewMockS3Storage(ctrl)
	mockRepo := NewMockDBRepo(ctrl)
	mockUserKafka := NewMockKafkaProducer(ctrl)
	mockSocietyKafka := NewMockKafkaProducer(ctrl)

	t.Run("set_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetSocietyAvatar")

		stream := NewMockAvatarService_SetSocietyAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetSocietyAvatarIn{
			Uuid:     "society-uuid",
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
		stream.EXPECT().SendAndClose(&avatar.SetSocietyAvatarOut{Link: "https://s3.example.com/society.jpg"}).Return(nil)

		mockS3.EXPECT().PutObject(gomock.Any(), &model.AvatarContent{
			AvatarType: model.SocietyAvatarType,
			UUID:       "society-uuid",
			Filename:   "test.jpg",
			ImageData:  []byte{1, 2, 3},
		}).Return("https://s3.example.com/society.jpg", nil)

		mockRepo.EXPECT().SetSocietyAvatar(gomock.Any(), "society-uuid", "https://s3.example.com/society.jpg").Return(nil)
		mockSocietyKafka.EXPECT().ProduceMessage(gomock.Any(), &avatar.NewAvatarRegister{
			Uuid: "society-uuid",
			Link: "https://s3.example.com/society.jpg",
		}, "society-uuid").Return(nil)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetSocietyAvatar(stream)

		assert.NoError(t, err)
	})

	t.Run("no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetSocietyAvatar")
		mockLogger.EXPECT().Error("society uuid is required")

		stream := NewMockAvatarService_SetSocietyAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetSocietyAvatarIn{
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetSocietyAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "society uuid is required")
	})

	t.Run("stream_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetSocietyAvatar")
		mockLogger.EXPECT().Error("failed to receive data from stream: stream error")

		stream := NewMockAvatarService_SetSocietyAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(nil, errors.New("stream error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetSocietyAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to receive data from stream")
	})

	t.Run("s3_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetSocietyAvatar")
		mockLogger.EXPECT().Error("failed to upload file to S3: s3 error")

		stream := NewMockAvatarService_SetSocietyAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetSocietyAvatarIn{
			Uuid:     "society-uuid",
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		mockS3.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return("", errors.New("s3 error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetSocietyAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to upload file to S3")
	})

	t.Run("repo_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetSocietyAvatar")
		mockLogger.EXPECT().Error("failed to save avatar to database: repo error")

		stream := NewMockAvatarService_SetSocietyAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetSocietyAvatarIn{
			Uuid:     "society-uuid",
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		mockS3.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return("https://s3.example.com/society.jpg", nil)
		mockRepo.EXPECT().SetSocietyAvatar(gomock.Any(), "society-uuid", "https://s3.example.com/society.jpg").Return(errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetSocietyAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to save avatar to database")
	})

	t.Run("kafka_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("SetSocietyAvatar")
		mockLogger.EXPECT().Error("failed to produce message to society service: kafka error")

		stream := NewMockAvatarService_SetSocietyAvatarServer(ctrl)
		stream.EXPECT().Context().Return(ctx).AnyTimes()
		stream.EXPECT().Recv().Return(&avatar.SetSocietyAvatarIn{
			Uuid:     "society-uuid",
			Filename: "test.jpg",
			Batch:    []byte{1, 2, 3},
		}, nil).Times(1)
		stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

		mockS3.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return("https://s3.example.com/society.jpg", nil)
		mockRepo.EXPECT().SetSocietyAvatar(gomock.Any(), "society-uuid", "https://s3.example.com/society.jpg").Return(nil)
		mockSocietyKafka.EXPECT().ProduceMessage(gomock.Any(), gomock.Any(), "society-uuid").Return(errors.New("kafka error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		err := s.SetSocietyAvatar(stream)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to produce message to society service")
	})
}

func TestService_GetAllSocietyAvatars(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockS3 := NewMockS3Storage(ctrl)
	mockRepo := NewMockDBRepo(ctrl)
	mockUserKafka := NewMockKafkaProducer(ctrl)
	mockSocietyKafka := NewMockKafkaProducer(ctrl)

	t.Run("get_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllSocietyAvatars")

		expectedAvatars := &model.AvatarMetadataList{
			{ID: 1, Link: "https://s3.example.com/society1.jpg"},
			{ID: 2, Link: "https://s3.example.com/society2.jpg"},
		}

		mockRepo.EXPECT().GetAllSocietyAvatars(gomock.Any(), "society-uuid").Return(expectedAvatars, nil)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		result, err := s.GetAllSocietyAvatars(ctx, &avatar.GetAllSocietyAvatarsIn{Uuid: "society-uuid"})

		assert.NoError(t, err)
		assert.Len(t, result.AvatarList, 2)
		assert.Equal(t, int32(1), result.AvatarList[0].Id)
		assert.Equal(t, "https://s3.example.com/society1.jpg", result.AvatarList[0].Link)
		assert.Equal(t, int32(2), result.AvatarList[1].Id)
		assert.Equal(t, "https://s3.example.com/society2.jpg", result.AvatarList[1].Link)
	})

	t.Run("no_uuid", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllSocietyAvatars")
		mockLogger.EXPECT().Error("society uuid is required")

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.GetAllSocietyAvatars(ctx, &avatar.GetAllSocietyAvatarsIn{Uuid: ""})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, st.Message(), "society uuid is required")
	})

	t.Run("repo_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("GetAllSocietyAvatars")
		mockLogger.EXPECT().Error("failed to get all society avatars: repo error")

		mockRepo.EXPECT().GetAllSocietyAvatars(gomock.Any(), "society-uuid").Return(nil, errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.GetAllSocietyAvatars(ctx, &avatar.GetAllSocietyAvatarsIn{Uuid: "society-uuid"})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to get all society avatars")
	})
}

func TestService_DeleteSocietyAvatar(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

	mockS3 := NewMockS3Storage(ctrl)
	mockRepo := NewMockDBRepo(ctrl)
	mockUserKafka := NewMockKafkaProducer(ctrl)
	mockSocietyKafka := NewMockKafkaProducer(ctrl)

	t.Run("delete_ok", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteSocietyAvatar")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "society-uuid",
			Link: "https://s3.example.com/society.jpg",
		}

		mockRepo.EXPECT().GetSocietyAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/society.jpg").Return(nil)
		mockRepo.EXPECT().DeleteSocietyAvatar(gomock.Any(), 1).Return(nil)
		mockRepo.EXPECT().GetLatestSocietyAvatar(gomock.Any(), "society-uuid").Return("https://s3.example.com/latest_society.jpg")
		mockSocietyKafka.EXPECT().ProduceMessage(gomock.Any(), &avatar.NewAvatarRegister{
			Uuid: "society-uuid",
			Link: "https://s3.example.com/latest_society.jpg",
		}, "society-uuid").Return(nil)

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		result, err := s.DeleteSocietyAvatar(ctx, &avatar.DeleteSocietyAvatarIn{AvatarId: 1})

		assert.NoError(t, err)
		assert.Equal(t, int32(1), result.Id)
		assert.Equal(t, "https://s3.example.com/society.jpg", result.Link)
	})

	t.Run("get_avatar_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteSocietyAvatar")
		mockLogger.EXPECT().Error("failed to get avatar data: repo error")

		mockRepo.EXPECT().GetSocietyAvatarData(gomock.Any(), 1).Return(nil, errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteSocietyAvatar(ctx, &avatar.DeleteSocietyAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Contains(t, st.Message(), "failed to get avatar data")
	})

	t.Run("s3_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteSocietyAvatar")
		mockLogger.EXPECT().Error("failed to delete avatar in s3: s3 error")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "society-uuid",
			Link: "https://s3.example.com/society.jpg",
		}

		mockRepo.EXPECT().GetSocietyAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/society.jpg").Return(errors.New("s3 error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteSocietyAvatar(ctx, &avatar.DeleteSocietyAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to delete avatar in s3")
	})

	t.Run("repo_delete_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteSocietyAvatar")
		mockLogger.EXPECT().Error("failed to delete avatar in postgres: repo error")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "society-uuid",
			Link: "https://s3.example.com/society.jpg",
		}

		mockRepo.EXPECT().GetSocietyAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/society.jpg").Return(nil)
		mockRepo.EXPECT().DeleteSocietyAvatar(gomock.Any(), 1).Return(errors.New("repo error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteSocietyAvatar(ctx, &avatar.DeleteSocietyAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to delete avatar in postgres")
	})

	t.Run("kafka_error", func(t *testing.T) {
		mockLogger.EXPECT().AddFuncName("DeleteSocietyAvatar")
		mockLogger.EXPECT().Error("failed to produce avatar: kafka error")

		avatarInfo := &model.AvatarMetadata{
			ID:   1,
			UUID: "society-uuid",
			Link: "https://s3.example.com/society.jpg",
		}

		mockRepo.EXPECT().GetSocietyAvatarData(gomock.Any(), 1).Return(avatarInfo, nil)
		mockS3.EXPECT().RemoveObject(gomock.Any(), "https://s3.example.com/society.jpg").Return(nil)
		mockRepo.EXPECT().DeleteSocietyAvatar(gomock.Any(), 1).Return(nil)
		mockRepo.EXPECT().GetLatestSocietyAvatar(gomock.Any(), "society-uuid").Return("https://s3.example.com/latest_society.jpg")
		mockSocietyKafka.EXPECT().ProduceMessage(gomock.Any(), gomock.Any(), "society-uuid").Return(errors.New("kafka error"))

		s := New(mockS3, mockRepo, mockUserKafka, mockSocietyKafka)
		_, err := s.DeleteSocietyAvatar(ctx, &avatar.DeleteSocietyAvatarIn{AvatarId: 1})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to produce avatar")
	})
}
