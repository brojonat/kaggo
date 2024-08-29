BEGIN;

DROP TABLE IF EXISTS twitch_clip_views;
DROP TABLE IF EXISTS twitch_video_views;
DROP TABLE IF EXISTS twitch_stream_views;
DROP TABLE IF EXISTS twitch_user_past_dec_avg_views;
DROP TABLE IF EXISTS twitch_user_past_dec_med_views;
DROP TABLE IF EXISTS twitch_user_past_dec_std_views;
DROP TABLE IF EXISTS twitch_user_past_dec_avg_duration;
DROP TABLE IF EXISTS twitch_user_past_dec_med_duration;
DROP TABLE IF EXISTS twitch_user_past_dec_std_duration;

COMMIT;