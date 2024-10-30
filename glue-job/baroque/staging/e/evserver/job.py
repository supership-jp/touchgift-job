import sys
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext
from awsglue.context import GlueContext
from awsglue.dynamicframe import DynamicFrame
from awsglue.job import Job
from datetime import datetime, timedelta
import pytz

def apply(inputFrame, glueContext):
    frame = inputFrame.toDF()
    gc = glueContext

    # タイムゾーンをJSTに設定
    jst = pytz.timezone('Asia/Tokyo')
    now = datetime.now(pytz.utc).astimezone(jst)

    yesterday = (now - timedelta(1)).strftime('%Y%m%d')

    frame.createOrReplaceTempView("application_table")

    query = f"""
    select
        request_id as request_id,
        time as timestamp,
        visitor_uuid,
        org_code,
        mid,
        ad_id,
        view_time,
        ev,
        campaign_id,
        dt
    from application_table
    where
        dt = '{yesterday}'
        and org_code = 'baroque'
        and request_id is not null
        and request_id != '';
    """

    transformed_df = gc.sparkSession.sql(query)

    return DynamicFrame.fromDF(transformed_df, gc)

# 引数を取得
# 'mode'引数が必須なので、指定されない場合はエラーで終了
args = getResolvedOptions(sys.argv, ['JOB_NAME', 'mode'])

sc = SparkContext()
glueContext = GlueContext(sc)
spark = glueContext.spark_session

job = Job(glueContext)
job.init(args['JOB_NAME'], args)

dyf = glueContext.create_dynamic_frame.from_catalog(
    database="touchgift-datalake-evserver",
    table_name="application",
)

spark.conf.set("spark.sql.legacy.timeParserPolicy", "LEGACY")

recipe = apply(
    inputFrame=dyf,
    glueContext=glueContext)

# 'mode'が'test'でない場合はS3に書き込む
if args['mode'] != 'test':
    glueContext.write_dynamic_frame.from_options(
        frame=recipe,
        connection_type="s3",
        format="glueparquet",
        connection_options={"path": "s3://baroque-data-link-staging", "partitionKeys": ["dt", "ev"]},
        format_options={"compression": "gzip"})
else:
    print("テストが完了しました。")

job.commit()