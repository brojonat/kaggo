// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package dbgen

import (
	jsonb "github.com/brojonat/kaggo/server/db/jsonb"
	"github.com/jackc/pgx/v5/pgtype"
)

type InternalRandom struct {
	ID  string             `json:"id"`
	Ts  pgtype.Timestamptz `json:"ts"`
	Val int32              `json:"val"`
}

type KaggleDatasetDownload struct {
	ID        string             `json:"id"`
	Ts        pgtype.Timestamptz `json:"ts"`
	Downloads int32              `json:"downloads"`
}

type KaggleDatasetView struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Views int32              `json:"views"`
}

type KaggleDatasetVote struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Votes int32              `json:"votes"`
}

type KaggleNotebookVote struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Votes int32              `json:"votes"`
}

type Metadatum struct {
	ID          string             `json:"id"`
	RequestKind string             `json:"request_kind"`
	Data        jsonb.MetadataJSON `json:"data"`
}

type RedditCommentControversiality struct {
	ID               string             `json:"id"`
	Ts               pgtype.Timestamptz `json:"ts"`
	Controversiality float32            `json:"controversiality"`
}

type RedditCommentScore struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Score int32              `json:"score"`
}

type RedditPostRatio struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Ratio float32            `json:"ratio"`
}

type RedditPostScore struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Score int32              `json:"score"`
}

type RedditSubredditActiveUserCount struct {
	ID              string             `json:"id"`
	Ts              pgtype.Timestamptz `json:"ts"`
	ActiveUserCount int32              `json:"active_user_count"`
}

type RedditSubredditSubscriber struct {
	ID          string             `json:"id"`
	Ts          pgtype.Timestamptz `json:"ts"`
	Subscribers int32              `json:"subscribers"`
}

type User struct {
	Email string                 `json:"email"`
	Data  jsonb.UserMetadataJSON `json:"data"`
}

type UsersMetadataThrough struct {
	Email       string `json:"email"`
	ID          string `json:"id"`
	RequestKind string `json:"request_kind"`
}

type YoutubeChannelSubscriber struct {
	ID          string             `json:"id"`
	Ts          pgtype.Timestamptz `json:"ts"`
	Subscribers int32              `json:"subscribers"`
}

type YoutubeChannelVideo struct {
	ID     string             `json:"id"`
	Ts     pgtype.Timestamptz `json:"ts"`
	Videos int32              `json:"videos"`
}

type YoutubeChannelView struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Views int64              `json:"views"`
}

type YoutubeVideoComment struct {
	ID       string             `json:"id"`
	Ts       pgtype.Timestamptz `json:"ts"`
	Comments int32              `json:"comments"`
}

type YoutubeVideoLike struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Likes int32              `json:"likes"`
}

type YoutubeVideoView struct {
	ID    string             `json:"id"`
	Ts    pgtype.Timestamptz `json:"ts"`
	Views int64              `json:"views"`
}
