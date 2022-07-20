package main

import (
	"context"
	"douyin/v1/cmd/video/pack"
	"douyin/v1/cmd/video/service"
	"douyin/v1/kitex_gen/video"
	"douyin/v1/pkg/errno"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// GetPublishListByUser implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishListByUser(ctx context.Context, userId int64) (resp *video.PublishListResponse, err error) {
	response := new(video.PublishListResponse)
	if userId < 0 {
		response.SetBaseResp(pack.BuildBaseResp(errno.ParamErr))
		return response, nil
	}
	videos, err := service.NewQueryVideoService(ctx).GetPublishList(userId)
	if err != nil {
		response.SetBaseResp(pack.BuildBaseResp(err))
		return response, nil
	}
	response.SetBaseResp(pack.BuildBaseResp(errno.Success))
	response.SetVideoList(videos)
	return response, nil
}

// GetVideosByLastTime implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideosByLastTime(ctx context.Context, lastTime int64, userId int64) (resp *video.VideoFeedResponse, err error) {
	response := new(video.VideoFeedResponse)
	if lastTime < 0 {
		response.SetBaseResp(pack.BuildBaseResp(errno.ParamErr))
		return response, nil
	}
	videos, nextTime, err := service.NewQueryVideoService(ctx).GetVideoFeed(lastTime, userId)
	if err != nil {
		response.SetBaseResp(pack.BuildBaseResp(err))
		return response, nil
	}
	response.SetBaseResp(pack.BuildBaseResp(errno.Success))
	response.SetVideoList(videos)
	response.SetNextTime(nextTime.Unix())
	return response, nil
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, publishedVideo *video.Video) (resp *video.BaseResp, err error) {
	if publishedVideo.GetId() != 0 {
		return pack.BuildBaseResp(errno.ParamErr), nil
	}
	if err := service.NewCreateVideoService(ctx).CreateVideo(publishedVideo); err != nil {
		return pack.BuildBaseResp(errno.ServiceErr), nil
	}
	return pack.BuildBaseResp(errno.Success), nil
}

// FavoriteByUser implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) FavoriteByUser(ctx context.Context, request *video.FavoriteActionRequest) (resp *video.BaseResp, err error) {
	response := new(video.BaseResp)
	if request.UserId < 0 {
		response = pack.BuildBaseResp(errno.ParamErr)
		return response, nil
	}
	err = service.NewQueryFavoriteService(ctx).FavoriteByUser(request.UserId, request.VideoId, request.ActionType)
	if err != nil {
		response = pack.BuildBaseResp(err)
		return response, nil
	}
	response = pack.BuildBaseResp(errno.Success)
	return response, nil
}

// GetFavoriteListBYUser implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetFavoriteListBYUser(ctx context.Context, userId int64) (resp *video.FavoriteListResponse, err error) {
	response := new(video.FavoriteListResponse)

	if userId < 0 {
		response.SetBaseResp(pack.BuildBaseResp(errno.ParamErr))
		return response, nil
	}
	videos, err := service.NewQueryFavoriteService(ctx).GetFavoriteListByUser(userId)
	if err != nil {
		response.SetBaseResp(pack.BuildBaseResp(err))
		return response, nil
	}
	response.SetBaseResp(pack.BuildBaseResp(errno.Success))
	response.SetVideoList(videos)
	return response, nil
}

// GetCommentListByVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetCommentListByVideo(ctx context.Context, videoId int64) (resp *video.CommentListResponse, err error) {
	response := new(video.CommentListResponse)
	comments, err := service.NewCommentService(ctx).QueryCommentByVideoId(videoId)
	if err != nil {
		response.SetBaseResp(pack.BuildBaseResp(err))
		return response, nil
	}
	response.SetBaseResp(pack.BuildBaseResp(errno.Success))
	response.SetCommentList(comments)
	return response, nil
}

// PostComment implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PostComment(ctx context.Context, commentActionRequest *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	response := new(video.CommentActionResponse)
	if commentActionRequest.ActionType == 1 {
		comment, err := service.NewCommentService(ctx).PostComment(commentActionRequest.Comment)
		if err != nil {
			response.SetBaseResp(pack.BuildBaseResp(err))
			return response, nil
		}
		response.SetBaseResp(pack.BuildBaseResp(errno.Success))
		response.SetComment(comment)
		return response, nil
	} else {
		if err := service.NewCommentService(ctx).DeleteComment(commentActionRequest.Comment.Id); err != nil {
			response.SetBaseResp(pack.BuildBaseResp(err))
		}
		response.SetBaseResp(pack.BuildBaseResp(errno.Success))
		return response, nil
	}
}
