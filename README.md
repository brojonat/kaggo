# kaggo

This idea started out as a premium Kaggle service. The value proposition was mining Kaggle for "good" notebooks and datasets. However, it has morphed into a general purpose metric tracker for things on the Internet. Users pay for us to monitor the performance of their content.

## For Developers

### How to add a new metric

Let's say you want to start tracking Twitch metrics. Here's how you'd go about implementing that.

- Go to the Twitch API and start browsing the reference for endpoints/resources of interest. You're only interested in endpoints/resources that are accessible with a vanilla API token; you're not going to have access to users' specific data since we're not getting their consent. We'll refer to these as "resources" below; resources are effectively endpoints that we can hit and return a response body that contain 1 or more metrics we want to track. Resources map 1:1 to `RequestKind` types.

- Here's the Twitch resources of interest that I found:

  - twitch.clip https://dev.twitch.tv/docs/api/reference/#get-clips
    - Metadata: id, url, broadcaster_name, creator_name, game_id, title, duration
    - Metrics: view_count
  - twitch.video https://dev.twitch.tv/docs/api/reference/#get-videos
    - Metadata: id, user_name, title, url,
    - Metric: view_count
  - twitch.stream https://dev.twitch.tv/docs/api/reference/#get-streams
    - Metadata: user_id, user_name
    - Metrics: viewer_count
    - NOTE: we supply a user_id and this resource only returns an entry if the user is actively streaming. So it's a cool metric BUT we need to handle the case when the user is offline to return 0. This is only potentially problematic for fetching metadata, but the metadata activity can be a no-op because the caller will need to specify the metadata anyway in order for us to find the stream.
  - twitch.user-past-dec
    - Metadata user_name
    - Metrics: viewer_count
    - This is problematic because we'll need to hit the get-videos endpoint and page through the videos. Or maybe we do sum(last 100 video views)? We can get that in a single query and it's a "relevant" metric.

- You'll need to inspect the response body for each of these. If the API documentation is good, you can do this on the docs page, otherwise you'll need to use your favorite HTTP client (e.g., curl, Bruno, etc), and manually make requests against the API to get some sample data. Then you can take the sample data and drop it into https://play.jmespath.org/ or something similar and determine the correct path to extract the quantity of interest (e.g., data[0].view_count).

- Add the payload structs to the server API.
- Open your Go code and open the `temporal/handlers_[metadata|metrics].go` files and add support to BOTH for extracting the 1) the metadata param(s) and 2) the metric(s) you want to track. Note that the metadata workflow and the long polling workflow _may_ pass different types of response bodies to their respective handlers! Typically they'll be the same, but in some cases (i.e., Twitch streaming and "recent" metrics), the responses are different.
- Add `RequestKindTwitchClips` and so on for each resource.
- Update `prepareRequest` to handle the new request kinds.
- Update `UploadMetadata` and `UploadMetrics` to handle the new request kinds.
- Write and apply migrations for the new tables
- Write and generate the queries for the new metrics (particularly insertion, and getting bucketed)
- Add the handlers to serve the metrics
- Update the schedule handler to perform the actual requests
