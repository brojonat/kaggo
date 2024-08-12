CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS vector;


SELECT *, 'knv' AS "metric" FROM kaggle_notebook_votes knv 
UNION ALL
SELECT *, 'knd' AS "metric" FROM kaggle_notebook_downloads knd ;

SELECT
    knv.slug AS "knv_slug",
    knv.ts AS "knv_ts",
    knv.val AS "knv_val",
    knd.slug AS "knd_slug",
    knd.ts AS "knd_ts",
    knd.val AS "knd_val"
FROM kaggle_notebook_votes knv
FULL JOIN kaggle_notebook_downloads knd ON knv.ts = knd.ts;


SELECT timescaledb_experimental.show_policies('kaggle_dataset_votes_mat');

