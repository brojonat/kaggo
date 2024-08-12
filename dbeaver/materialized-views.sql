-- this will return the existing materialized views and continuous aggregates
SELECT view_name, format('%I.%I', materialization_hypertable_schema,
        materialization_hypertable_name) AS materialization_hypertable
FROM timescaledb_information.continuous_aggregates;
   

DROP MATERIALIZED VIEW IF EXISTS kaggle_notebook_votes_mat;
CREATE MATERIALIZED VIEW IF NOT EXISTS kaggle_notebook_votes_mat
WITH (timescaledb.continuous) AS
SELECT slug,
   time_bucket(INTERVAL '15 min', ts) AS bucket,
   MAX(val)
FROM kaggle_notebook_votes
GROUP BY slug, bucket;

SELECT add_continuous_aggregate_policy('kaggle_notebook_votes_mat',
  start_offset => INTERVAL '7 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '30 min');

-- notebook downloads continuous aggregates
DROP MATERIALIZED VIEW IF EXISTS kaggle_notebook_downloads_mat;
CREATE MATERIALIZED VIEW IF NOT EXISTS kaggle_notebook_downloads_mat
WITH (timescaledb.continuous) AS
SELECT slug,
   time_bucket(INTERVAL '15 min', ts) AS bucket,
   MAX(val)
FROM kaggle_notebook_downloads
GROUP BY slug, bucket;

SELECT add_continuous_aggregate_policy('kaggle_notebook_downloads_mat',
  start_offset => INTERVAL '7 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '30 min');
 

-- dataset downloads continuous aggregates 
DROP MATERIALIZED VIEW IF EXISTS kaggle_dataset_votes_mat;
CREATE MATERIALIZED VIEW IF NOT EXISTS kaggle_dataset_votes_mat
WITH (timescaledb.continuous) AS
SELECT slug,
   time_bucket(INTERVAL '15 min', ts) AS bucket,
   MAX(val)
FROM kaggle_dataset_votes
GROUP BY slug, bucket;

SELECT add_continuous_aggregate_policy('kaggle_dataset_votes_mat',
  start_offset => INTERVAL '7 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '30 min');

-- dataset downloads continuous aggregates
DROP MATERIALIZED VIEW IF EXISTS kaggle_dataset_downloads_mat;
CREATE MATERIALIZED VIEW IF NOT EXISTS kaggle_dataset_downloads_mat
WITH (timescaledb.continuous) AS
SELECT slug,
   time_bucket(INTERVAL '15 min', ts) AS bucket,
   MAX(val)
FROM kaggle_dataset_downloads
GROUP BY slug, bucket;

SELECT add_continuous_aggregate_policy('kaggle_dataset_downloads_mat',
  start_offset => INTERVAL '7 days',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '30 min');
 

