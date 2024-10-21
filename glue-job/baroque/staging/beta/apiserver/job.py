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

    # JSTに変換
    jst = pytz.timezone('Asia/Tokyo')
    now = datetime.now(pytz.utc).astimezone(jst)

    yesterday = (now - timedelta(1)).strftime('%Y%m%d')

    frame.createOrReplaceTempView("application_table")

    query = f"""
    select
        COALESCE(request_id, NULL) as request_id,
        COALESCE(time, NULL) as timestamp,
        COALESCE(org_code, NULL) as org_code,
        COALESCE(visitor_uuid, NULL) as visitor_uuid,
        COALESCE(mid, NULL) as mid,
        COALESCE(message, NULL) as ev,
        COALESCE(dt, NULL) as dt,
        COALESCE(adid, NULL) as adid,
        COALESCE(idfa, NULL) as idfa,
        COALESCE(fcm_token, NULL) as fcm_token,
        COALESCE(store_id, NULL) as store_id,
        COALESCE(os, NULL) as os,
        COALESCE(os_version, NULL) as os_version,
        COALESCE(device, NULL) as device,
        COALESCE(app_name, NULL) as app_name,
        COALESCE(network_type, NULL) as network_type,
        COALESCE(mcc, NULL) as mcc,
        COALESCE(mnc, NULL) as mnc,
        COALESCE(lang, NULL) as lang,
        COALESCE(campaign_id, NULL) as campaign_id,
        COALESCE(coupon_id, NULL) as coupon_id,
        COALESCE(screen_id, NULL) as screen_id
    from application_table
    where
        dt = '{yesterday}'
        and request_id is not null
        and request_id != ''
        and api = 'application'
        and (message = 'touch' or message = 'coupon_draw' or message = 'screen_imp');
    """

    transformed_df = gc.sparkSession.sql(query)
    transformed_df["ev"] = transformed_df["ev"].replace('coupon_draw', 'coupon_get_imp')

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
    database="touchgift-datalake-apiserver-beta",
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
        connection_options={"path": "s3://baroque-data-link-staging-beta", "partitionKeys": ["dt", "ev"]},
        format_options={"compression": "gzip"})
else:
    print("テストが完了しました。")

job.commit()