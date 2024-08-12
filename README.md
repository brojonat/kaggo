# kaggo

This idea started out as a premium Kaggle service. The value proposition was mining Kaggle for "good" notebooks and datasets. However, it has morphed into a general purpose metric tracker for things like kaggle notebooks and datasets. Users pay for us to monitor the performance of their content.

# Kaggle API

Check out the Kaggle API here: https://github.com/Kaggle/kaggle-api/blob/main/docs/KaggleApi.md.

The official Kaggle (python) client is surprisingly great, but it only provides a point in time snapshot of the current state of the user specified input. We can do better by providing a real time dashboard of trending notebooks and datasets. Here's the API of our service:

GET kaggle/datasets -> semantic search for datasets
GET kaggle/datasets/slug -> dataset timeseries data
POST kaggle/datasets/slug {votes: N, downloads: M}

GET kaggle/kernels -> semantic search for kernels
GET kaggle/kernels/slug -> dataset timeseries data
POST kaggle/kernels/slug {votes: N, downloads: M}
