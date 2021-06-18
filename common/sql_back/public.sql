/*
 Navicat PostgreSQL Data Transfer

 Source Server         : 10.41.1.241
 Source Server Type    : PostgreSQL
 Source Server Version : 100010
 Source Host           : 10.41.1.241:5432
 Source Catalog        : mp_crm
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 100010
 File Encoding         : 65001

 Date: 30/04/2020 16:55:11
*/


-- ----------------------------
-- Sequence structure for servicer_count_list_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."servicer_count_list_id_seq";
CREATE SEQUENCE "public"."servicer_count_list_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS "public"."account";
CREATE TABLE "public"."account" (
  "uid" uuid NOT NULL,
  "nickname" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "account" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "password" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "use_status" int8 DEFAULT 1,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "drop_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "modify_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "update_time" timestamp(0) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "phone" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "email" varchar(512) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "master_acc" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "is_delete" int2 DEFAULT 0,
  "usd_balance" int8 DEFAULT 0,
  "khr_balance" int8 DEFAULT 0,
  "gen_key" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "is_actived" int2 DEFAULT 0,
  "head_portrait_img_no" varchar(255) COLLATE "pg_catalog"."default",
  "last_login_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "is_first_login" int2 DEFAULT 0,
  "app_lang" varchar(255) COLLATE "pg_catalog"."default",
  "pos_lang" varchar(255) COLLATE "pg_catalog"."default",
  "country_code" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."account"."uid" IS '用户id';
COMMENT ON COLUMN "public"."account"."nickname" IS '用户名';
COMMENT ON COLUMN "public"."account"."account" IS '账号';
COMMENT ON COLUMN "public"."account"."password" IS '密码';
COMMENT ON COLUMN "public"."account"."use_status" IS '使用状态：0.禁用，1.正常';
COMMENT ON COLUMN "public"."account"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."account"."modify_time" IS '最后更新时间';
COMMENT ON COLUMN "public"."account"."phone" IS '手机';
COMMENT ON COLUMN "public"."account"."email" IS '邮件';
COMMENT ON COLUMN "public"."account"."master_acc" IS '主账号';
COMMENT ON COLUMN "public"."account"."is_delete" IS '账号删除';
COMMENT ON COLUMN "public"."account"."usd_balance" IS '美金余额';
COMMENT ON COLUMN "public"."account"."khr_balance" IS '柬埔寨余额';
COMMENT ON COLUMN "public"."account"."gen_key" IS '注册生成的扫码';
COMMENT ON COLUMN "public"."account"."is_actived" IS '是否激活;0-未激活,1-激活.';
COMMENT ON COLUMN "public"."account"."head_portrait_img_no" IS '头像图片url';
COMMENT ON COLUMN "public"."account"."last_login_time" IS '最后登录时间';
COMMENT ON COLUMN "public"."account"."is_first_login" IS '是否第一次登陆,0-是;1-否(这个字段只给pos端使用)';
COMMENT ON COLUMN "public"."account"."country_code" IS '国家区号';

-- ----------------------------
-- Table structure for account_collect
-- ----------------------------
DROP TABLE IF EXISTS "public"."account_collect";
CREATE TABLE "public"."account_collect" (
  "account_no" uuid NOT NULL,
  "collect_account_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "collect_phone" int8,
  "is_delete" int2 DEFAULT 0,
  "create_time" timestamp(0),
  "modify_time" timestamp(0),
  "account_collect_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."account_collect"."account_no" IS '账号';
COMMENT ON COLUMN "public"."account_collect"."collect_account_no" IS '收集的账号Uid';
COMMENT ON COLUMN "public"."account_collect"."collect_phone" IS '收集的转账手机号';
COMMENT ON COLUMN "public"."account_collect"."modify_time" IS '最后修改时间（用于排序）';
COMMENT ON COLUMN "public"."account_collect"."account_collect_no" IS '主键';
COMMENT ON TABLE "public"."account_collect" IS '转账后，最近转账的人';

-- ----------------------------
-- Table structure for adminlog
-- ----------------------------
DROP TABLE IF EXISTS "public"."adminlog";
CREATE TABLE "public"."adminlog" (
  "log_uid" uuid NOT NULL,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "url" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "param" varchar(512) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "op" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "op_type" int4 DEFAULT 0,
  "op_acc_uid" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "during" int4 NOT NULL DEFAULT 0,
  "ip" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "status_code" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text,
  "response" text COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::text
)
;
COMMENT ON COLUMN "public"."adminlog"."log_uid" IS '日志uid';
COMMENT ON COLUMN "public"."adminlog"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."adminlog"."url" IS 'url';
COMMENT ON COLUMN "public"."adminlog"."param" IS '参数';
COMMENT ON COLUMN "public"."adminlog"."op" IS '操作';
COMMENT ON COLUMN "public"."adminlog"."op_type" IS '操作类型';
COMMENT ON COLUMN "public"."adminlog"."op_acc_uid" IS '操作账号uid';
COMMENT ON TABLE "public"."adminlog" IS '运行日志';

-- ----------------------------
-- Table structure for agreement
-- ----------------------------
DROP TABLE IF EXISTS "public"."agreement";
CREATE TABLE "public"."agreement" (
  "id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "text" text COLLATE "pg_catalog"."default",
  "lang" varchar(255) COLLATE "pg_catalog"."default",
  "type" int2,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "modify_time" timestamp(6),
  "use_status" int2 DEFAULT 1
)
;
COMMENT ON COLUMN "public"."agreement"."id" IS 'id';
COMMENT ON COLUMN "public"."agreement"."text" IS '协议内容';
COMMENT ON COLUMN "public"."agreement"."lang" IS '语言(km、zh_CN、en)';
COMMENT ON COLUMN "public"."agreement"."type" IS '0 用户协议  1隐私协议';
COMMENT ON COLUMN "public"."agreement"."is_delete" IS '1删除';
COMMENT ON COLUMN "public"."agreement"."modify_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."agreement"."use_status" IS '使用状态0禁用';

-- ----------------------------
-- Table structure for app_version
-- ----------------------------
DROP TABLE IF EXISTS "public"."app_version";
CREATE TABLE "public"."app_version" (
  "v_id" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default" DEFAULT ''::text,
  "version" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamp(0),
  "update_time" timestamp(0),
  "app_url" varchar(2000) COLLATE "pg_catalog"."default",
  "vs_code" varchar(255) COLLATE "pg_catalog"."default",
  "vs_type" int2 DEFAULT 0,
  "is_force" int2 DEFAULT 0,
  "system" int2 DEFAULT 0,
  "is_delete" int2 DEFAULT 0,
  "account_uid" uuid,
  "note" text COLLATE "pg_catalog"."default" DEFAULT ''::text,
  "status" int2 DEFAULT 1
)
;
COMMENT ON COLUMN "public"."app_version"."v_id" IS 'app版本';
COMMENT ON COLUMN "public"."app_version"."description" IS '版本描述';
COMMENT ON COLUMN "public"."app_version"."version" IS '版本';
COMMENT ON COLUMN "public"."app_version"."create_time" IS '生成时间';
COMMENT ON COLUMN "public"."app_version"."update_time" IS '更改时间';
COMMENT ON COLUMN "public"."app_version"."app_url" IS '下载地址';
COMMENT ON COLUMN "public"."app_version"."vs_code" IS '版本码,做版本比较使用';
COMMENT ON COLUMN "public"."app_version"."vs_type" IS '0-app包;1-pos包';
COMMENT ON COLUMN "public"."app_version"."is_force" IS '是否强制更新;0-否,1-是';
COMMENT ON COLUMN "public"."app_version"."system" IS '0-android系统;1-ios系统';
COMMENT ON COLUMN "public"."app_version"."is_delete" IS '1删除';
COMMENT ON COLUMN "public"."app_version"."account_uid" IS '上传版本的账户uid';
COMMENT ON COLUMN "public"."app_version"."note" IS '备注';
COMMENT ON COLUMN "public"."app_version"."status" IS '是否禁用(0禁用1启用)';

-- ----------------------------
-- Table structure for app_version_file_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."app_version_file_log";
CREATE TABLE "public"."app_version_file_log" (
  "id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "create_time" timestamp(0),
  "account_no" uuid,
  "file_name" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for billing_details_results
-- ----------------------------
DROP TABLE IF EXISTS "public"."billing_details_results";
CREATE TABLE "public"."billing_details_results" (
  "create_time" timestamp(0),
  "bill_no" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "amount" int8 DEFAULT 0,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "bill_type" int2,
  "account_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "account_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "order_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "balance" int8 DEFAULT 0,
  "order_status" int2,
  "modify_time" timestamp(0),
  "servicer_no" uuid,
  "op_acc_no" uuid
)
;
COMMENT ON COLUMN "public"."billing_details_results"."create_time" IS '记账时间';
COMMENT ON COLUMN "public"."billing_details_results"."bill_no" IS '账单流水号';
COMMENT ON COLUMN "public"."billing_details_results"."amount" IS '金额';
COMMENT ON COLUMN "public"."billing_details_results"."currency_type" IS '币种(usd,khr)';
COMMENT ON COLUMN "public"."billing_details_results"."bill_type" IS '账单类型（1:存款;2:取款;3:收益;4:充值;5-提现）';
COMMENT ON COLUMN "public"."billing_details_results"."account_no" IS '账号uid';
COMMENT ON COLUMN "public"."billing_details_results"."order_no" IS '订单号';
COMMENT ON COLUMN "public"."billing_details_results"."balance" IS '余额(暂时放这，后面确认完全无用再删)';
COMMENT ON COLUMN "public"."billing_details_results"."order_status" IS '订单状态（1-初始化;2-等待;3-已支付;4-失败(超时);5-待确认;6-取消）';
COMMENT ON COLUMN "public"."billing_details_results"."modify_time" IS '更改时间';
COMMENT ON TABLE "public"."billing_details_results" IS '服务商的总账流水-非虚拟';

-- ----------------------------
-- Table structure for card
-- ----------------------------
DROP TABLE IF EXISTS "public"."card";
CREATE TABLE "public"."card" (
  "card_no" uuid NOT NULL,
  "account_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "channel_no" uuid,
  "name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "card_number" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "balance_type" varchar(12) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "is_defalut" int2 DEFAULT 0,
  "collect_status" int2 DEFAULT 1,
  "audit_status" int2 DEFAULT 0,
  "note" text COLLATE "pg_catalog"."default" DEFAULT ''::text,
  "modify_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."card"."card_no" IS '卡uid';
COMMENT ON COLUMN "public"."card"."account_no" IS '账号uuid';
COMMENT ON COLUMN "public"."card"."channel_no" IS '渠道';
COMMENT ON COLUMN "public"."card"."name" IS '收款人名';
COMMENT ON COLUMN "public"."card"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."card"."is_delete" IS '是否删除';
COMMENT ON COLUMN "public"."card"."card_number" IS '卡号';
COMMENT ON COLUMN "public"."card"."balance_type" IS '币种（usd美金，khr瑞尔）';
COMMENT ON COLUMN "public"."card"."is_defalut" IS '1:推荐收款卡';
COMMENT ON COLUMN "public"."card"."collect_status" IS '收款状态（0禁用，1启用）';
COMMENT ON COLUMN "public"."card"."audit_status" IS '审核状态,0-未审核,1-审核通过';
COMMENT ON COLUMN "public"."card"."note" IS '备注';

-- ----------------------------
-- Table structure for cashier
-- ----------------------------
DROP TABLE IF EXISTS "public"."cashier";
CREATE TABLE "public"."cashier" (
  "uid" uuid NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "servicer_no" uuid,
  "is_delete" int2 DEFAULT 0,
  "create_time" timestamp(6),
  "op_password" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "modify_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."cashier"."uid" IS '操作员';
COMMENT ON COLUMN "public"."cashier"."name" IS '姓名';
COMMENT ON COLUMN "public"."cashier"."servicer_no" IS '服务商';
COMMENT ON COLUMN "public"."cashier"."op_password" IS '操作密码';

-- ----------------------------
-- Table structure for channel
-- ----------------------------
DROP TABLE IF EXISTS "public"."channel";
CREATE TABLE "public"."channel" (
  "channel_no" uuid NOT NULL,
  "channel_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "note" text COLLATE "pg_catalog"."default" DEFAULT ''::text,
  "idx" int8 DEFAULT 0,
  "is_recom" int2 DEFAULT 0,
  "channel_type" int2 DEFAULT 0,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "use_status" int2 DEFAULT 0,
  "logo_img_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."channel"."channel_name" IS '渠道名称';
COMMENT ON COLUMN "public"."channel"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."channel"."is_delete" IS '是否删除';
COMMENT ON COLUMN "public"."channel"."note" IS '备注';
COMMENT ON COLUMN "public"."channel"."idx" IS '序列';
COMMENT ON COLUMN "public"."channel"."is_recom" IS '是否推荐(1-推荐，0-不推荐)';
COMMENT ON COLUMN "public"."channel"."channel_type" IS '0-通用,1-用户,总部;2-pos';
COMMENT ON COLUMN "public"."channel"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."channel"."use_status" IS '状态(1正常)';
COMMENT ON COLUMN "public"."channel"."logo_img_no" IS 'logo图片id';

-- ----------------------------
-- Table structure for channel_servicer
-- ----------------------------
DROP TABLE IF EXISTS "public"."channel_servicer";
CREATE TABLE "public"."channel_servicer" (
  "channel_no" uuid NOT NULL,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "idx" int8 DEFAULT 0,
  "is_recom" int2 DEFAULT 0,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "use_status" int2 DEFAULT 1,
  "id" varchar(32) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."channel_servicer"."channel_no" IS '渠道仓库表的id';
COMMENT ON COLUMN "public"."channel_servicer"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."channel_servicer"."is_delete" IS '是否删除';
COMMENT ON COLUMN "public"."channel_servicer"."idx" IS '序列';
COMMENT ON COLUMN "public"."channel_servicer"."is_recom" IS '是否推荐(1-推荐，0-不推荐)';
COMMENT ON COLUMN "public"."channel_servicer"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."channel_servicer"."use_status" IS '状态(0禁用)';
COMMENT ON COLUMN "public"."channel_servicer"."id" IS '主键（用于删除时定位）';

-- ----------------------------
-- Table structure for collection_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."collection_order";
CREATE TABLE "public"."collection_order" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "from_vaccount_no" uuid,
  "to_vaccount_no" uuid,
  "amount" int8 DEFAULT 0,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "finish_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "order_status" int8 DEFAULT 0,
  "balance_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "payment_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT 2,
  "fees" int8 DEFAULT 0,
  "is_count" int2 DEFAULT 0,
  "modify_time" timestamp(6),
  "ip" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "lat" varchar(255) COLLATE "pg_catalog"."default",
  "lng" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."collection_order"."log_no" IS '虚拟账户收款日志主键';
COMMENT ON COLUMN "public"."collection_order"."order_status" IS '订单状态(0-交易失败，1-交易成功)';
COMMENT ON COLUMN "public"."collection_order"."balance_type" IS '金额类型（usd美金，khr瑞尔）';
COMMENT ON COLUMN "public"."collection_order"."payment_type" IS '付款方式;1-现金;2-余额';
COMMENT ON COLUMN "public"."collection_order"."fees" IS '手续费';
COMMENT ON COLUMN "public"."collection_order"."is_count" IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';

-- ----------------------------
-- Table structure for common_help
-- ----------------------------
DROP TABLE IF EXISTS "public"."common_help";
CREATE TABLE "public"."common_help" (
  "help_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "problem" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "answer" text COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "idx" int8,
  "is_delete" int2 DEFAULT 0,
  "use_status" int2 DEFAULT 1,
  "lang" varchar(255) COLLATE "pg_catalog"."default",
  "vs_type" int2,
  "modify_time" timestamp(6),
  "create_time" timestamp(6),
  "file_id" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."common_help"."help_no" IS '主键';
COMMENT ON COLUMN "public"."common_help"."problem" IS '问题';
COMMENT ON COLUMN "public"."common_help"."answer" IS '答案';
COMMENT ON COLUMN "public"."common_help"."idx" IS '排序';
COMMENT ON COLUMN "public"."common_help"."use_status" IS '使用状态0-禁用';
COMMENT ON COLUMN "public"."common_help"."lang" IS '语言(km、zh_CN、en)';
COMMENT ON COLUMN "public"."common_help"."vs_type" IS '0-app包;1-pos包';
COMMENT ON COLUMN "public"."common_help"."modify_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."common_help"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."common_help"."file_id" IS '答案的文件id';

-- ----------------------------
-- Table structure for consultation_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."consultation_config";
CREATE TABLE "public"."consultation_config" (
  "id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "use_status" int2 DEFAULT 1,
  "is_delete" int2 DEFAULT 0,
  "create_time" timestamp(0),
  "lang" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "idx" int8 DEFAULT 0,
  "logo_img_no" varchar(32) COLLATE "pg_catalog"."default",
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "text" varchar(2000) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."consultation_config"."id" IS '咨询表';
COMMENT ON COLUMN "public"."consultation_config"."use_status" IS '使用状态：0.禁用，1.正常';
COMMENT ON COLUMN "public"."consultation_config"."is_delete" IS '1删除';
COMMENT ON COLUMN "public"."consultation_config"."lang" IS '语言(km、zh_CN、en)';
COMMENT ON COLUMN "public"."consultation_config"."idx" IS '优先级（值小的再上面）';
COMMENT ON COLUMN "public"."consultation_config"."logo_img_no" IS '图片id';
COMMENT ON COLUMN "public"."consultation_config"."name" IS '主题';
COMMENT ON COLUMN "public"."consultation_config"."text" IS '内容';

-- ----------------------------
-- Table structure for cust
-- ----------------------------
DROP TABLE IF EXISTS "public"."cust";
CREATE TABLE "public"."cust" (
  "cust_no" uuid NOT NULL,
  "account_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "payment_password" varchar(255) COLLATE "pg_catalog"."default",
  "gender" int2 DEFAULT 1,
  "in_authorization" int2 DEFAULT 1,
  "out_authorization" int2 DEFAULT 1,
  "in_transfer_authorization" int2 DEFAULT 1,
  "out_transfer_authorization" int2 DEFAULT 1,
  "modify_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "def_pay_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT 'usd_balance'::character varying
)
;
COMMENT ON COLUMN "public"."cust"."cust_no" IS '用户列表主键';
COMMENT ON COLUMN "public"."cust"."account_no" IS '用户账号uuid';
COMMENT ON COLUMN "public"."cust"."payment_password" IS '用户支付密码';
COMMENT ON COLUMN "public"."cust"."gender" IS '用户性别（1男，0女）';
COMMENT ON COLUMN "public"."cust"."in_authorization" IS '充值权限(0禁用)';
COMMENT ON COLUMN "public"."cust"."out_authorization" IS '提现权限(0禁用)';
COMMENT ON COLUMN "public"."cust"."in_transfer_authorization" IS '可转账入权限(0禁用)';
COMMENT ON COLUMN "public"."cust"."out_transfer_authorization" IS '可转账出权限(0禁用)';
COMMENT ON COLUMN "public"."cust"."def_pay_no" IS '默认卡';
COMMENT ON TABLE "public"."cust" IS '主要存账户的静态信息';

-- ----------------------------
-- Table structure for dict_acc_title
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_acc_title";
CREATE TABLE "public"."dict_acc_title" (
  "title_no" int8 NOT NULL,
  "title_name" varchar(255) COLLATE "pg_catalog"."default",
  "parent_title" int8
)
;
COMMENT ON COLUMN "public"."dict_acc_title"."title_no" IS '会计科目';
COMMENT ON COLUMN "public"."dict_acc_title"."title_name" IS '会计科目名';
COMMENT ON COLUMN "public"."dict_acc_title"."parent_title" IS '父节点';
COMMENT ON TABLE "public"."dict_acc_title" IS '会计科目表';

-- ----------------------------
-- Table structure for dict_account_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_account_type";
CREATE TABLE "public"."dict_account_type" (
  "account_type" int8 NOT NULL,
  "remark" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for dict_bank_abbr
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_bank_abbr";
CREATE TABLE "public"."dict_bank_abbr" (
  "id" varchar(255) COLLATE "pg_catalog"."default",
  "bank_abbr" varchar(255) COLLATE "pg_catalog"."default",
  "bank_name" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."dict_bank_abbr"."bank_abbr" IS '缩写';
COMMENT ON COLUMN "public"."dict_bank_abbr"."bank_name" IS '银行名';

-- ----------------------------
-- Table structure for dict_bankname
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_bankname";
CREATE TABLE "public"."dict_bankname" (
  "bank_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "bank_id" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for dict_bin_bankname
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_bin_bankname";
CREATE TABLE "public"."dict_bin_bankname" (
  "bin_code" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "bank_name" varchar(255) COLLATE "pg_catalog"."default",
  "org_code" varchar(255) COLLATE "pg_catalog"."default",
  "card_name" varchar(255) COLLATE "pg_catalog"."default",
  "card_type" varchar(255) COLLATE "pg_catalog"."default",
  "card_type_no" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."dict_bin_bankname"."bank_name" IS '银行名';
COMMENT ON COLUMN "public"."dict_bin_bankname"."card_type_no" IS '0:借记卡

1:贷记卡
2:准贷记卡

3:预付费卡
';

-- ----------------------------
-- Table structure for dict_images
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_images";
CREATE TABLE "public"."dict_images" (
  "image_id" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "image_url" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "status" int2 DEFAULT 1,
  "modify_time" timestamp(6),
  "account_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "is_delete" int2 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."dict_images"."status" IS '1-正常';

-- ----------------------------
-- Table structure for dict_org_abbr
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_org_abbr";
CREATE TABLE "public"."dict_org_abbr" (
  "org_code" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "abbr" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."dict_org_abbr"."org_code" IS '机构号';
COMMENT ON COLUMN "public"."dict_org_abbr"."abbr" IS '缩写';

-- ----------------------------
-- Table structure for dict_province
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_province";
CREATE TABLE "public"."dict_province" (
  "province_code" int8 NOT NULL,
  "province_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "short_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "full_en_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "short_zh_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."dict_province"."province_code" IS '省码';
COMMENT ON COLUMN "public"."dict_province"."province_name" IS '省份名';
COMMENT ON COLUMN "public"."dict_province"."short_name" IS '简写';
COMMENT ON COLUMN "public"."dict_province"."full_en_name" IS '全拼';
COMMENT ON COLUMN "public"."dict_province"."short_zh_name" IS '中文简称';
COMMENT ON TABLE "public"."dict_province" IS '省份表';

-- ----------------------------
-- Table structure for dict_region
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_region";
CREATE TABLE "public"."dict_region" (
  "id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "code" varchar(64) COLLATE "pg_catalog"."default",
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "level" varchar(32) COLLATE "pg_catalog"."default",
  "pid" varchar(64) COLLATE "pg_catalog"."default",
  "longitude" numeric(10,4),
  "latitude" numeric(10,4),
  "is_leaf" int2,
  "pname" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."dict_region"."is_leaf" IS '1是;0不是';
COMMENT ON TABLE "public"."dict_region" IS '地区表';

-- ----------------------------
-- Table structure for dict_region_bank
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_region_bank";
CREATE TABLE "public"."dict_region_bank" (
  "code" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "province" varchar(255) COLLATE "pg_catalog"."default",
  "city" varchar(255) COLLATE "pg_catalog"."default",
  "bank_type" varchar(255) COLLATE "pg_catalog"."default",
  "province_code" varchar(255) COLLATE "pg_catalog"."default",
  "city_code" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for dict_vatype
-- ----------------------------
DROP TABLE IF EXISTS "public"."dict_vatype";
CREATE TABLE "public"."dict_vatype" (
  "va_type" int8 NOT NULL,
  "remark" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for exchange_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."exchange_order";
CREATE TABLE "public"."exchange_order" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "in_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "out_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "amount" int8 DEFAULT 0,
  "create_time" timestamp(6),
  "rate" int8 DEFAULT 0,
  "order_status" int8 DEFAULT 0,
  "finish_time" timestamp(6),
  "account_no" uuid,
  "trans_from" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "trans_amount" int8 DEFAULT 0,
  "err_reason" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "fees" int8 DEFAULT 0,
  "is_count" int2,
  "modify_time" timestamp(6),
  "ip" varchar(255) COLLATE "pg_catalog"."default",
  "lat" varchar(255) COLLATE "pg_catalog"."default",
  "lng" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."exchange_order"."in_type" IS '发起金额类型';
COMMENT ON COLUMN "public"."exchange_order"."out_type" IS '到账金额类型';
COMMENT ON COLUMN "public"."exchange_order"."amount" IS '金额';
COMMENT ON COLUMN "public"."exchange_order"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."exchange_order"."rate" IS '平台汇率';
COMMENT ON COLUMN "public"."exchange_order"."order_status" IS '订单状态';
COMMENT ON COLUMN "public"."exchange_order"."finish_time" IS '完成时间';
COMMENT ON COLUMN "public"."exchange_order"."account_no" IS '发起人';
COMMENT ON COLUMN "public"."exchange_order"."trans_from" IS 'app,trade';
COMMENT ON COLUMN "public"."exchange_order"."trans_amount" IS '转换后金额';
COMMENT ON COLUMN "public"."exchange_order"."err_reason" IS '失败原因';
COMMENT ON COLUMN "public"."exchange_order"."fees" IS '单笔手续费';
COMMENT ON COLUMN "public"."exchange_order"."is_count" IS '手续费是否统计完毕;0-否;1-是';
COMMENT ON COLUMN "public"."exchange_order"."modify_time" IS '修改时间';

-- ----------------------------
-- Table structure for func_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."func_config";
CREATE TABLE "public"."func_config" (
  "func_no" uuid NOT NULL,
  "func_name" varchar(255) COLLATE "pg_catalog"."default",
  "idx" int8,
  "use_status" int4 DEFAULT 1,
  "is_delete" int2 DEFAULT 0,
  "img" varchar(512) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "jump_url" varchar(512) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "application_type" varchar(32) COLLATE "pg_catalog"."default",
  "img_id" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."func_config"."func_no" IS '钱包功能入口配置  主键';
COMMENT ON COLUMN "public"."func_config"."use_status" IS '0禁用';
COMMENT ON COLUMN "public"."func_config"."img" IS '图片地址';
COMMENT ON COLUMN "public"."func_config"."jump_url" IS '跳转地址';
COMMENT ON COLUMN "public"."func_config"."application_type" IS '0-手机端;1-pos';
COMMENT ON COLUMN "public"."func_config"."img_id" IS '图片id';
COMMENT ON TABLE "public"."func_config" IS '钱包功能';

-- ----------------------------
-- Table structure for gen_code
-- ----------------------------
DROP TABLE IF EXISTS "public"."gen_code";
CREATE TABLE "public"."gen_code" (
  "gen_key" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "account_no" uuid,
  "amount" int8 DEFAULT 0,
  "money_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6),
  "code_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "use_status" int4 DEFAULT 1,
  "modify_time" timestamp(6),
  "sweep_account_no" uuid,
  "order_no" varchar COLLATE "pg_catalog"."default",
  "op_acc_type" int2 DEFAULT 0,
  "op_acc_no" uuid
)
;
COMMENT ON COLUMN "public"."gen_code"."account_no" IS '服务商账号id';
COMMENT ON COLUMN "public"."gen_code"."use_status" IS '1-初始化;2-已扫码;3-已支付;4-已过期;5-待确认';
COMMENT ON COLUMN "public"."gen_code"."sweep_account_no" IS '谁扫的码,accountid';
COMMENT ON COLUMN "public"."gen_code"."order_no" IS '订单ID';
COMMENT ON COLUMN "public"."gen_code"."op_acc_type" IS '这个码是谁生成的,0-用户1-服务商;2-收银员';
COMMENT ON COLUMN "public"."gen_code"."op_acc_no" IS '操作员账号,当前账号类型是服务商,就是服务商的id,如果是收银员,那就是收银员的id';

-- ----------------------------
-- Table structure for global_param
-- ----------------------------
DROP TABLE IF EXISTS "public"."global_param";
CREATE TABLE "public"."global_param" (
  "param_key" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "param_value" text COLLATE "pg_catalog"."default",
  "remark" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."global_param"."remark" IS '备注';
COMMENT ON TABLE "public"."global_param" IS '全局变量';

-- ----------------------------
-- Table structure for headquarters_profit
-- ----------------------------
DROP TABLE IF EXISTS "public"."headquarters_profit";
CREATE TABLE "public"."headquarters_profit" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "general_ledger_no" varchar COLLATE "pg_catalog"."default",
  "amount" int8,
  "create_time" timestamp(6),
  "order_status" int2 DEFAULT 0,
  "finish_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "balance_type" varchar(255) COLLATE "pg_catalog"."default",
  "profit_source" int2,
  "modify_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."headquarters_profit"."log_no" IS '平台盈利统计流水号';
COMMENT ON COLUMN "public"."headquarters_profit"."general_ledger_no" IS '订单';
COMMENT ON COLUMN "public"."headquarters_profit"."amount" IS '金额';
COMMENT ON COLUMN "public"."headquarters_profit"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."headquarters_profit"."order_status" IS '状态';
COMMENT ON COLUMN "public"."headquarters_profit"."finish_time" IS '完成时间';
COMMENT ON COLUMN "public"."headquarters_profit"."balance_type" IS '金额类型（usd美金，khr瑞尔）';
COMMENT ON COLUMN "public"."headquarters_profit"."profit_source" IS '收益来源（1-提现手续费;2-客户转账手续费;3-客户兑换手续费;4-客户收款手续费;5客户存款手续费）';

-- ----------------------------
-- Table structure for headquarters_profit_cashable
-- ----------------------------
DROP TABLE IF EXISTS "public"."headquarters_profit_cashable";
CREATE TABLE "public"."headquarters_profit_cashable" (
  "id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "cashable_balance" int8 DEFAULT 0,
  "revenue_money" int8 DEFAULT 0,
  "modify_time" timestamp(0) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."headquarters_profit_cashable"."id" IS '平台利润';
COMMENT ON COLUMN "public"."headquarters_profit_cashable"."cashable_balance" IS '可提现余额';
COMMENT ON COLUMN "public"."headquarters_profit_cashable"."revenue_money" IS '总收益的钱';
COMMENT ON COLUMN "public"."headquarters_profit_cashable"."modify_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."headquarters_profit_cashable"."currency_type" IS '币种';

-- ----------------------------
-- Table structure for headquarters_profit_withdraw
-- ----------------------------
DROP TABLE IF EXISTS "public"."headquarters_profit_withdraw";
CREATE TABLE "public"."headquarters_profit_withdraw" (
  "order_no" varchar(32) COLLATE "pg_catalog"."default" NOT NULL,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "amount" int8,
  "note" text COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "account_no" uuid
)
;
COMMENT ON COLUMN "public"."headquarters_profit_withdraw"."order_no" IS '平台盈利提现订单表主键';
COMMENT ON COLUMN "public"."headquarters_profit_withdraw"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."headquarters_profit_withdraw"."amount" IS '金额';
COMMENT ON COLUMN "public"."headquarters_profit_withdraw"."note" IS '备注';
COMMENT ON COLUMN "public"."headquarters_profit_withdraw"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."headquarters_profit_withdraw"."account_no" IS '操作人账号';

-- ----------------------------
-- Table structure for income_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."income_order";
CREATE TABLE "public"."income_order" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "act_acc_no" uuid,
  "amount" int8 DEFAULT 0,
  "servicer_no" uuid,
  "create_time" timestamp(6),
  "order_status" int8 DEFAULT 0,
  "finish_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "query_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "balance_type" varchar(255) COLLATE "pg_catalog"."default",
  "fees" int8 DEFAULT 0,
  "recv_acc_no" uuid,
  "recv_vacc" uuid,
  "op_acc_no" uuid,
  "settle_hourly_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "settle_daily_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "payment_type" int2 DEFAULT 1,
  "is_count" int2,
  "modify_time" timestamp(6),
  "ree_rate" int8,
  "real_amount" int8,
  "op_acc_type" int2
)
;
COMMENT ON COLUMN "public"."income_order"."log_no" IS '存款日志表';
COMMENT ON COLUMN "public"."income_order"."act_acc_no" IS '存款人的账号';
COMMENT ON COLUMN "public"."income_order"."amount" IS '金额';
COMMENT ON COLUMN "public"."income_order"."servicer_no" IS '收款的商户uid';
COMMENT ON COLUMN "public"."income_order"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."income_order"."order_status" IS '订单状态（1-初始化;2-等待;3-已支付;4-失败）';
COMMENT ON COLUMN "public"."income_order"."finish_time" IS '完成时间';
COMMENT ON COLUMN "public"."income_order"."query_time" IS '查询时间';
COMMENT ON COLUMN "public"."income_order"."balance_type" IS '金额类型（usd美金，khr瑞尔）';
COMMENT ON COLUMN "public"."income_order"."fees" IS '手续费';
COMMENT ON COLUMN "public"."income_order"."recv_acc_no" IS '收款人账号';
COMMENT ON COLUMN "public"."income_order"."recv_vacc" IS '收款人虚拟账号';
COMMENT ON COLUMN "public"."income_order"."op_acc_no" IS '操作员账号,当前账号类型是服务商,就是服务商的id,如果是收银员,那就是收银员的id';
COMMENT ON COLUMN "public"."income_order"."settle_hourly_log_no" IS '小时对账流水';
COMMENT ON COLUMN "public"."income_order"."settle_daily_log_no" IS '天对账流水';
COMMENT ON COLUMN "public"."income_order"."payment_type" IS '付款方式.1-现金;2-余额';
COMMENT ON COLUMN "public"."income_order"."is_count" IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';
COMMENT ON COLUMN "public"."income_order"."ree_rate" IS '手续费率';
COMMENT ON COLUMN "public"."income_order"."real_amount" IS '实际到账金额';
COMMENT ON COLUMN "public"."income_order"."op_acc_type" IS '这笔订单是谁产生的.1-服务商;2-店员';
COMMENT ON TABLE "public"."income_order" IS '收款订单';

-- ----------------------------
-- Table structure for income_ougo_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."income_ougo_config";
CREATE TABLE "public"."income_ougo_config" (
  "income_ougo_config_no" uuid NOT NULL,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "use_status" varchar(255) COLLATE "pg_catalog"."default" DEFAULT 0,
  "idx" int8,
  "config_type" int2,
  "is_delete" int2 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."income_ougo_config"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."income_ougo_config"."name" IS '名称';
COMMENT ON COLUMN "public"."income_ougo_config"."use_status" IS '状态(1正常)';
COMMENT ON COLUMN "public"."income_ougo_config"."idx" IS '排序序号';
COMMENT ON COLUMN "public"."income_ougo_config"."config_type" IS '类型（1.充值方式。2.提现方式）';

-- ----------------------------
-- Table structure for income_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."income_type";
CREATE TABLE "public"."income_type" (
  "income_type" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "income_name" varchar(255) COLLATE "pg_catalog"."default",
  "use_status" int2 DEFAULT 0,
  "idx" int8
)
;

-- ----------------------------
-- Table structure for lang
-- ----------------------------
DROP TABLE IF EXISTS "public"."lang";
CREATE TABLE "public"."lang" (
  "key" varchar(512) COLLATE "pg_catalog"."default" NOT NULL,
  "type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 1,
  "is_delete" int2 DEFAULT 0,
  "lang_km" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "lang_en" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "lang_ch" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."lang"."type" IS '类型（1-文字，2-图片, 3-错误提示）';
COMMENT ON COLUMN "public"."lang"."lang_km" IS '柬埔寨语';
COMMENT ON COLUMN "public"."lang"."lang_ch" IS '中文';

-- ----------------------------
-- Table structure for log_account
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_account";
CREATE TABLE "public"."log_account" (
  "log_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "account_uid" varchar COLLATE "pg_catalog"."default",
  "log_time" timestamp(6),
  "type" int4
)
;
COMMENT ON COLUMN "public"."log_account"."description" IS '描述';
COMMENT ON TABLE "public"."log_account" IS '账户操作日志';

-- ----------------------------
-- Table structure for log_account_web
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_account_web";
CREATE TABLE "public"."log_account_web" (
  "log_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "account_uid" varchar COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "type" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."log_account_web"."log_no" IS 'id（此表专门记WEB也就是后台管理系统的日志）';
COMMENT ON COLUMN "public"."log_account_web"."description" IS '描述';
COMMENT ON COLUMN "public"."log_account_web"."type" IS '1-账号；
';
COMMENT ON TABLE "public"."log_account_web" IS '账户操作日志';

-- ----------------------------
-- Table structure for log_app_messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_app_messages";
CREATE TABLE "public"."log_app_messages" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default",
  "order_no" varchar(32) COLLATE "pg_catalog"."default",
  "order_type" int2,
  "is_read" int2 DEFAULT 0,
  "is_push" int2 DEFAULT 0,
  "account_no" uuid,
  "create_time" timestamp(0)
)
;
COMMENT ON COLUMN "public"."log_app_messages"."order_no" IS '订单号';
COMMENT ON COLUMN "public"."log_app_messages"."order_type" IS '订单类型（1-兑换,2-充值,3-提现,4-转账,5-收款）';
COMMENT ON COLUMN "public"."log_app_messages"."is_read" IS '是否已读';
COMMENT ON COLUMN "public"."log_app_messages"."is_push" IS '是否推送';
COMMENT ON COLUMN "public"."log_app_messages"."account_no" IS '要推送的账号';
COMMENT ON TABLE "public"."log_app_messages" IS '用户消息中心';

-- ----------------------------
-- Table structure for log_card
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_card";
CREATE TABLE "public"."log_card" (
  "log_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "card_num" varchar COLLATE "pg_catalog"."default" DEFAULT 0,
  "name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT 0,
  "account_no" uuid,
  "va_type" int2 DEFAULT 0,
  "channel_no" uuid,
  "channel_type" int2,
  "create_time" timestamp(6),
  "descript" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."log_card"."card_num" IS '卡号';
COMMENT ON COLUMN "public"."log_card"."name" IS '持卡人名字';
COMMENT ON COLUMN "public"."log_card"."account_no" IS '谁创建的';
COMMENT ON COLUMN "public"."log_card"."va_type" IS '币种';
COMMENT ON COLUMN "public"."log_card"."channel_no" IS '渠道id';
COMMENT ON COLUMN "public"."log_card"."channel_type" IS '1-用户,总部;2-pos';
COMMENT ON COLUMN "public"."log_card"."descript" IS '描述';

-- ----------------------------
-- Table structure for log_exchange_rate
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_exchange_rate";
CREATE TABLE "public"."log_exchange_rate" (
  "log_time" timestamp(6),
  "usd_khr" int8 DEFAULT 0,
  "khr_usd" int8 DEFAULT 0
)
;
COMMENT ON TABLE "public"."log_exchange_rate" IS '汇率变更表';

-- ----------------------------
-- Table structure for log_login
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_login";
CREATE TABLE "public"."log_login" (
  "log_time" timestamp(6),
  "acc_no" uuid,
  "ip" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "result" int4 DEFAULT 0,
  "client" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "lat" varchar(255) COLLATE "pg_catalog"."default",
  "lng" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."log_login"."lat" IS '维度';
COMMENT ON COLUMN "public"."log_login"."lng" IS '经度';

-- ----------------------------
-- Table structure for log_to_headquarters
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_to_headquarters";
CREATE TABLE "public"."log_to_headquarters" (
  "log_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "servicer_no" uuid,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "amount" int8,
  "order_status" int2 DEFAULT 0,
  "collection_type" int2,
  "card_no" uuid,
  "create_time" timestamp(0) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "finish_time" timestamp(0) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "order_type" int2,
  "image_id" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."log_to_headquarters"."log_no" IS '转账至总部';
COMMENT ON COLUMN "public"."log_to_headquarters"."servicer_no" IS '服务商id';
COMMENT ON COLUMN "public"."log_to_headquarters"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."log_to_headquarters"."amount" IS '金额';
COMMENT ON COLUMN "public"."log_to_headquarters"."order_status" IS '订单状态(0待审核、1已完成，2已关闭)';
COMMENT ON COLUMN "public"."log_to_headquarters"."collection_type" IS '收款方式,1-支票;2-现金;3-银行转账;4-其他';
COMMENT ON COLUMN "public"."log_to_headquarters"."card_no" IS '卡uid';
COMMENT ON COLUMN "public"."log_to_headquarters"."create_time" IS '发起时间';
COMMENT ON COLUMN "public"."log_to_headquarters"."finish_time" IS '完成时间';
COMMENT ON COLUMN "public"."log_to_headquarters"."order_type" IS '订单类型（1-交易转账、2-结算转账）';

-- ----------------------------
-- Table structure for log_to_servicer
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_to_servicer";
CREATE TABLE "public"."log_to_servicer" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "servicer_no" uuid,
  "collection_type" int2,
  "card_no" uuid,
  "amount" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamp(0),
  "order_type" int2,
  "order_status" int2,
  "finish_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "motify_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone
)
;
COMMENT ON COLUMN "public"."log_to_servicer"."log_no" IS '转账至服务商流水';
COMMENT ON COLUMN "public"."log_to_servicer"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."log_to_servicer"."servicer_no" IS '服务商id';
COMMENT ON COLUMN "public"."log_to_servicer"."collection_type" IS '收款方式,1-支票;2-现金;3-银行转账;4-其他';
COMMENT ON COLUMN "public"."log_to_servicer"."card_no" IS '卡主键uid';
COMMENT ON COLUMN "public"."log_to_servicer"."amount" IS '金额';
COMMENT ON COLUMN "public"."log_to_servicer"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."log_to_servicer"."order_type" IS '订单类型（1-交易转账、2-结算转账）';
COMMENT ON COLUMN "public"."log_to_servicer"."order_status" IS '订单状态(0-待审核、1-已完成，2-已关闭)';
COMMENT ON COLUMN "public"."log_to_servicer"."finish_time" IS '完成时间';

-- ----------------------------
-- Table structure for log_vaccount
-- ----------------------------
DROP TABLE IF EXISTS "public"."log_vaccount";
CREATE TABLE "public"."log_vaccount" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "vaccount_no" uuid NOT NULL,
  "create_time" timestamp(6),
  "amount" int8 DEFAULT 0,
  "op_type" int2 DEFAULT 0,
  "frozen_balance" int8 DEFAULT 0,
  "balance" int8 DEFAULT 0,
  "reason" int8 DEFAULT 0,
  "settle_hourly_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "settle_daily_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "biz_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."log_vaccount"."log_no" IS '虚拟账户日志表';
COMMENT ON COLUMN "public"."log_vaccount"."amount" IS '变化的金额';
COMMENT ON COLUMN "public"."log_vaccount"."op_type" IS '1:+;2:-;3:冻结;4:解冻;';
COMMENT ON COLUMN "public"."log_vaccount"."frozen_balance" IS '冻结余额';
COMMENT ON COLUMN "public"."log_vaccount"."balance" IS '现有余额';
COMMENT ON COLUMN "public"."log_vaccount"."reason" IS '原因;1-兑换,2-充值,3-提现,4-转账,5-收款;6-手续费;7-pos 端取消提现';
COMMENT ON COLUMN "public"."log_vaccount"."settle_hourly_log_no" IS '每小时对账流水';
COMMENT ON COLUMN "public"."log_vaccount"."settle_daily_log_no" IS '每日对账流水';
COMMENT ON COLUMN "public"."log_vaccount"."biz_log_no" IS '业务流水';
COMMENT ON TABLE "public"."log_vaccount" IS '虚拟账户日志表';

-- ----------------------------
-- Table structure for login_token
-- ----------------------------
DROP TABLE IF EXISTS "public"."login_token";
CREATE TABLE "public"."login_token" (
  "acc_no" uuid NOT NULL,
  "routes" text COLLATE "pg_catalog"."default",
  "token" text COLLATE "pg_catalog"."default",
  "login_time" timestamp(6),
  "ip" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "last_op_time" timestamp(6),
  "imei" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
COMMENT ON COLUMN "public"."login_token"."acc_no" IS '登录账号';
COMMENT ON COLUMN "public"."login_token"."routes" IS '路径列表';
COMMENT ON COLUMN "public"."login_token"."token" IS '令牌';
COMMENT ON COLUMN "public"."login_token"."login_time" IS '登录时间';
COMMENT ON COLUMN "public"."login_token"."ip" IS '登录ip';
COMMENT ON COLUMN "public"."login_token"."last_op_time" IS '最后操作时间';
COMMENT ON COLUMN "public"."login_token"."imei" IS 'imei';

-- ----------------------------
-- Table structure for outgo_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."outgo_order";
CREATE TABLE "public"."outgo_order" (
  "log_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "vaccount_no" uuid,
  "amount" int8,
  "create_time" timestamp(6),
  "order_status" int8,
  "modify_time" timestamp(6),
  "balance_type" varchar(255) COLLATE "pg_catalog"."default",
  "fees" int8,
  "servicer_no" uuid,
  "op_acc_no" uuid,
  "settle_hourly_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "settle_daily_log_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "rate" varchar(255) COLLATE "pg_catalog"."default",
  "payment_type" int2 DEFAULT 2,
  "is_count" int2,
  "withdraw_type" int2,
  "cancel_reason" varchar(255) COLLATE "pg_catalog"."default",
  "risk_no" varchar COLLATE "pg_catalog"."default",
  "real_amount" int8,
  "op_acc_type" int2,
  "ip" varchar(255) COLLATE "pg_catalog"."default",
  "lat" varchar(255) COLLATE "pg_catalog"."default",
  "lng" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."outgo_order"."log_no" IS '取款';
COMMENT ON COLUMN "public"."outgo_order"."order_status" IS '订单状态（1-初始化;2-等待;3-已支付;4-失败(超时);5-待确认;6-取消）';
COMMENT ON COLUMN "public"."outgo_order"."fees" IS '手续费';
COMMENT ON COLUMN "public"."outgo_order"."op_acc_no" IS '操作员账号,当前账号类型是服务商,就是服务商的id,如果是收银员,那就是收银员的id';
COMMENT ON COLUMN "public"."outgo_order"."settle_hourly_log_no" IS '小时对账流水';
COMMENT ON COLUMN "public"."outgo_order"."settle_daily_log_no" IS '天对账流水';
COMMENT ON COLUMN "public"."outgo_order"."rate" IS '费率';
COMMENT ON COLUMN "public"."outgo_order"."payment_type" IS '付款方式;1-现金,2-余额';
COMMENT ON COLUMN "public"."outgo_order"."is_count" IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';
COMMENT ON COLUMN "public"."outgo_order"."withdraw_type" IS '0-手机号提现;1-扫码提现;2-全部提现';
COMMENT ON COLUMN "public"."outgo_order"."cancel_reason" IS '取消原因';
COMMENT ON COLUMN "public"."outgo_order"."risk_no" IS '风控结果主键id,pos机查风控数据使用';
COMMENT ON COLUMN "public"."outgo_order"."real_amount" IS '实际到账金额';
COMMENT ON COLUMN "public"."outgo_order"."op_acc_type" IS '这笔订单是谁产生的,1-服务商;2-收银员';
COMMENT ON COLUMN "public"."outgo_order"."lat" IS '维度';
COMMENT ON COLUMN "public"."outgo_order"."lng" IS '经度';

-- ----------------------------
-- Table structure for outgo_type
-- ----------------------------
DROP TABLE IF EXISTS "public"."outgo_type";
CREATE TABLE "public"."outgo_type" (
  "outgo_type" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "outgo_name" varchar(255) COLLATE "pg_catalog"."default",
  "use_status" int2 DEFAULT 0,
  "idx" int8
)
;
COMMENT ON TABLE "public"."outgo_type" IS '提现方式';

-- ----------------------------
-- Table structure for platform_config
-- ----------------------------
DROP TABLE IF EXISTS "public"."platform_config";
CREATE TABLE "public"."platform_config" (
  "account_uid" uuid NOT NULL,
  "top_menu_status" int8 DEFAULT 1,
  "side_menu_status" int8 DEFAULT 1
)
;
COMMENT ON COLUMN "public"."platform_config"."account_uid" IS '账户uid';
COMMENT ON COLUMN "public"."platform_config"."top_menu_status" IS '顶部菜单配置1、打开、0关闭';
COMMENT ON COLUMN "public"."platform_config"."side_menu_status" IS '侧边菜单配置1、打开、0关闭';

-- ----------------------------
-- Table structure for rela_acc_iden
-- ----------------------------
DROP TABLE IF EXISTS "public"."rela_acc_iden";
CREATE TABLE "public"."rela_acc_iden" (
  "account_no" uuid NOT NULL,
  "account_type" int8 NOT NULL,
  "iden_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid
)
;
COMMENT ON COLUMN "public"."rela_acc_iden"."account_type" IS '1: 管理员2: 运营3: 服务商4 :用户5:操作员';

-- ----------------------------
-- Table structure for rela_account_role
-- ----------------------------
DROP TABLE IF EXISTS "public"."rela_account_role";
CREATE TABLE "public"."rela_account_role" (
  "rela_uid" uuid NOT NULL,
  "account_uid" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "role_uid" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid
)
;
COMMENT ON COLUMN "public"."rela_account_role"."account_uid" IS '账户uid';
COMMENT ON COLUMN "public"."rela_account_role"."role_uid" IS '角色uid';
COMMENT ON TABLE "public"."rela_account_role" IS '账户角色关联表';

-- ----------------------------
-- Table structure for rela_imei_pubkey
-- ----------------------------
DROP TABLE IF EXISTS "public"."rela_imei_pubkey";
CREATE TABLE "public"."rela_imei_pubkey" (
  "rela_no" uuid NOT NULL,
  "imei" varchar(255) COLLATE "pg_catalog"."default",
  "pub_key" text COLLATE "pg_catalog"."default",
  "create_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."rela_imei_pubkey"."imei" IS '设备号';
COMMENT ON COLUMN "public"."rela_imei_pubkey"."pub_key" IS '客户端公钥';

-- ----------------------------
-- Table structure for rela_role_url
-- ----------------------------
DROP TABLE IF EXISTS "public"."rela_role_url";
CREATE TABLE "public"."rela_role_url" (
  "rela_uid" uuid NOT NULL,
  "url_uid" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "role_uid" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid
)
;
COMMENT ON COLUMN "public"."rela_role_url"."url_uid" IS 'url_uid';
COMMENT ON COLUMN "public"."rela_role_url"."role_uid" IS '角色uid';
COMMENT ON TABLE "public"."rela_role_url" IS '角色-权限关系';

-- ----------------------------
-- Table structure for risk_result
-- ----------------------------
DROP TABLE IF EXISTS "public"."risk_result";
CREATE TABLE "public"."risk_result" (
  "risk_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "risk_result" int2,
  "risk_threshold" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "api_type" varchar(255) COLLATE "pg_catalog"."default",
  "payer_acc_no" varchar(255) COLLATE "pg_catalog"."default",
  "action_time" varchar(255) COLLATE "pg_catalog"."default",
  "eva_execute_type" varchar(255) COLLATE "pg_catalog"."default",
  "eva_score" varchar(255) COLLATE "pg_catalog"."default",
  "money_type" varchar(255) COLLATE "pg_catalog"."default",
  "order_no" varchar(255) COLLATE "pg_catalog"."default",
  "score" int8,
  "product_type" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."risk_result"."risk_no" IS '风控结果';
COMMENT ON COLUMN "public"."risk_result"."risk_result" IS '风控结果是否风控,0-否;1-是';
COMMENT ON COLUMN "public"."risk_result"."risk_threshold" IS '风控阈值';
COMMENT ON COLUMN "public"."risk_result"."create_time" IS '生成时间';
COMMENT ON COLUMN "public"."risk_result"."api_type" IS '事件';
COMMENT ON COLUMN "public"."risk_result"."payer_acc_no" IS '操作账号';
COMMENT ON COLUMN "public"."risk_result"."action_time" IS '操作时间';
COMMENT ON COLUMN "public"."risk_result"."eva_execute_type" IS '风控执行类型';
COMMENT ON COLUMN "public"."risk_result"."eva_score" IS '风控估算的得分';
COMMENT ON COLUMN "public"."risk_result"."money_type" IS '币种';
COMMENT ON COLUMN "public"."risk_result"."order_no" IS '订单号';
COMMENT ON COLUMN "public"."risk_result"."score" IS '实际得分';
COMMENT ON COLUMN "public"."risk_result"."product_type" IS '产品类型,如转账:trancefer';

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS "public"."role";
CREATE TABLE "public"."role" (
  "role_no" uuid NOT NULL,
  "role_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "modify_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "acc_type" varchar(255) COLLATE "pg_catalog"."default",
  "def_type" int2 DEFAULT 0,
  "master_acc" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "is_delete" int2 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."role"."role_no" IS '角色uid';
COMMENT ON COLUMN "public"."role"."role_name" IS '角色名';
COMMENT ON COLUMN "public"."role"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."role"."modify_time" IS '修改时间';
COMMENT ON COLUMN "public"."role"."acc_type" IS '账号类型';
COMMENT ON COLUMN "public"."role"."def_type" IS '是否设为默认页面';
COMMENT ON COLUMN "public"."role"."master_acc" IS '主账号';
COMMENT ON TABLE "public"."role" IS '角色表';

-- ----------------------------
-- Table structure for servicer
-- ----------------------------
DROP TABLE IF EXISTS "public"."servicer";
CREATE TABLE "public"."servicer" (
  "servicer_no" uuid NOT NULL,
  "account_no" uuid,
  "addr" varchar(512) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "use_status" int4 DEFAULT 1,
  "commission_sharing" int8 DEFAULT 0,
  "income_authorization" int2 DEFAULT 1,
  "outgo_authorization" int2 DEFAULT 1,
  "open_idx" int8 DEFAULT 0,
  "contact_person" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "contact_phone" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "contact_addr" varchar(1024) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "lat" varchar(255) COLLATE "pg_catalog"."default" DEFAULT 0,
  "lng" varchar(255) COLLATE "pg_catalog"."default" DEFAULT 0,
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "modify_time" timestamp(6),
  "income_sharing" int8 DEFAULT 0,
  "scope" varchar(64) COLLATE "pg_catalog"."default" DEFAULT 0,
  "scope_off" int2 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."servicer"."servicer_no" IS '服务商表主键';
COMMENT ON COLUMN "public"."servicer"."account_no" IS '账号uid';
COMMENT ON COLUMN "public"."servicer"."addr" IS '开店地址';
COMMENT ON COLUMN "public"."servicer"."create_time" IS '开通时间';
COMMENT ON COLUMN "public"."servicer"."is_delete" IS '是否删除';
COMMENT ON COLUMN "public"."servicer"."use_status" IS '服务商状态0禁用1启用';
COMMENT ON COLUMN "public"."servicer"."commission_sharing" IS '取款手续费分成';
COMMENT ON COLUMN "public"."servicer"."income_authorization" IS '收款权限,0-禁用;1-开通';
COMMENT ON COLUMN "public"."servicer"."outgo_authorization" IS '取款权限,0-禁用;1-开通';
COMMENT ON COLUMN "public"."servicer"."open_idx" IS '开户序列';
COMMENT ON COLUMN "public"."servicer"."contact_person" IS '联系人';
COMMENT ON COLUMN "public"."servicer"."contact_phone" IS '联系电话';
COMMENT ON COLUMN "public"."servicer"."contact_addr" IS '联系人地址';
COMMENT ON COLUMN "public"."servicer"."lat" IS '纬度';
COMMENT ON COLUMN "public"."servicer"."lng" IS '经度';
COMMENT ON COLUMN "public"."servicer"."password" IS '服务商支付密码';
COMMENT ON COLUMN "public"."servicer"."income_sharing" IS '存款手续费分成';
COMMENT ON COLUMN "public"."servicer"."scope" IS '营业范围(公里)';
COMMENT ON COLUMN "public"."servicer"."scope_off" IS '围栏开关（0关闭1开启）';

-- ----------------------------
-- Table structure for servicer_count
-- ----------------------------
DROP TABLE IF EXISTS "public"."servicer_count";
CREATE TABLE "public"."servicer_count" (
  "servicer_no" uuid NOT NULL,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "in_num" int8 NOT NULL DEFAULT 0,
  "in_amount" int8 NOT NULL DEFAULT 0,
  "out_num" int8 NOT NULL DEFAULT 0,
  "out_amount" int8 NOT NULL DEFAULT 0,
  "profit_num" int8 NOT NULL DEFAULT 0,
  "modify_time" timestamp(0),
  "profit_amount" int8 NOT NULL DEFAULT 0,
  "recharge_num" int8 NOT NULL DEFAULT 0,
  "recharge_amount" int8 NOT NULL DEFAULT 0,
  "withdraw_num" int8 NOT NULL DEFAULT 0,
  "withdraw_amount" int8 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."servicer_count"."servicer_no" IS '运营商uuid';
COMMENT ON COLUMN "public"."servicer_count"."currency_type" IS '币种(usd,khr)';
COMMENT ON COLUMN "public"."servicer_count"."in_num" IS '存款-总数量';
COMMENT ON COLUMN "public"."servicer_count"."in_amount" IS '存款-总金额';
COMMENT ON COLUMN "public"."servicer_count"."out_num" IS '取款-总数量';
COMMENT ON COLUMN "public"."servicer_count"."out_amount" IS '取款-总金额';
COMMENT ON COLUMN "public"."servicer_count"."profit_num" IS '收益-总数量';
COMMENT ON COLUMN "public"."servicer_count"."modify_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."servicer_count"."profit_amount" IS '收益-总金额';
COMMENT ON COLUMN "public"."servicer_count"."recharge_num" IS '充值-总数量';
COMMENT ON COLUMN "public"."servicer_count"."recharge_amount" IS '充值-总金额';
COMMENT ON COLUMN "public"."servicer_count"."withdraw_num" IS '提现-总数量';
COMMENT ON COLUMN "public"."servicer_count"."withdraw_amount" IS '提现-总金额';
COMMENT ON TABLE "public"."servicer_count" IS '服务商累计统计';

-- ----------------------------
-- Table structure for servicer_count_list
-- ----------------------------
DROP TABLE IF EXISTS "public"."servicer_count_list";
CREATE TABLE "public"."servicer_count_list" (
  "servicer_no" uuid,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamp(0),
  "in_num" int4 NOT NULL DEFAULT 0,
  "in_amount" int8 NOT NULL DEFAULT 0,
  "out_num" int4 NOT NULL DEFAULT 0,
  "out_amount" int8 NOT NULL DEFAULT 0,
  "profit_num" int4 NOT NULL DEFAULT 0,
  "profit_amount" int8 NOT NULL DEFAULT 0,
  "recharge_num" int4 NOT NULL DEFAULT 0,
  "recharge_amount" int8 NOT NULL DEFAULT 0,
  "withdraw_num" int4 NOT NULL DEFAULT 0,
  "withdraw_amount" int8 NOT NULL DEFAULT 0,
  "id" int8 NOT NULL DEFAULT nextval('servicer_count_list_id_seq'::regclass),
  "dates" date,
  "is_counted" int2 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."servicer_count_list"."servicer_no" IS '服务商uuid';
COMMENT ON COLUMN "public"."servicer_count_list"."currency_type" IS '币种(usd,khr)';
COMMENT ON COLUMN "public"."servicer_count_list"."create_time" IS '生成时间';
COMMENT ON COLUMN "public"."servicer_count_list"."in_num" IS '存款-数量';
COMMENT ON COLUMN "public"."servicer_count_list"."in_amount" IS '存款-金额';
COMMENT ON COLUMN "public"."servicer_count_list"."out_num" IS '取款-数量';
COMMENT ON COLUMN "public"."servicer_count_list"."out_amount" IS '取款-金额';
COMMENT ON COLUMN "public"."servicer_count_list"."profit_num" IS '收益-数量';
COMMENT ON COLUMN "public"."servicer_count_list"."profit_amount" IS '收益-金额';
COMMENT ON COLUMN "public"."servicer_count_list"."recharge_num" IS '充值-数量';
COMMENT ON COLUMN "public"."servicer_count_list"."recharge_amount" IS '充值-金额';
COMMENT ON COLUMN "public"."servicer_count_list"."withdraw_num" IS '提现-数量';
COMMENT ON COLUMN "public"."servicer_count_list"."withdraw_amount" IS '提现-金额';
COMMENT ON COLUMN "public"."servicer_count_list"."id" IS '自增主键';
COMMENT ON COLUMN "public"."servicer_count_list"."dates" IS '日期';
COMMENT ON COLUMN "public"."servicer_count_list"."is_counted" IS '是否已经统计完成(1是，0否)';
COMMENT ON TABLE "public"."servicer_count_list" IS '服务商对账单表';

-- ----------------------------
-- Table structure for servicer_img
-- ----------------------------
DROP TABLE IF EXISTS "public"."servicer_img";
CREATE TABLE "public"."servicer_img" (
  "servicer_img_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "img_id" varchar(400) COLLATE "pg_catalog"."default" DEFAULT ''::text,
  "img_type" int2 DEFAULT 0,
  "create_time" timestamp(6) NOT NULL,
  "servicer_no" uuid,
  "is_delete" int2 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."servicer_img"."servicer_img_no" IS '商户照片主键';
COMMENT ON COLUMN "public"."servicer_img"."img_id" IS '照片id';
COMMENT ON COLUMN "public"."servicer_img"."img_type" IS '照片类型（1.营业执照。2.3.4 经营场所照片）';
COMMENT ON COLUMN "public"."servicer_img"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."servicer_img"."servicer_no" IS '商户uuid';

-- ----------------------------
-- Table structure for servicer_profit_ledger
-- ----------------------------
DROP TABLE IF EXISTS "public"."servicer_profit_ledger";
CREATE TABLE "public"."servicer_profit_ledger" (
  "log_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "amount_order" int8,
  "servicefee_amount_sum" int8,
  "split_proportion" int8,
  "actual_income" int8,
  "payment_time" timestamp(0),
  "servicer_no" uuid NOT NULL,
  "currency_type" varchar(255) COLLATE "pg_catalog"."default",
  "order_type" int2
)
;
COMMENT ON COLUMN "public"."servicer_profit_ledger"."log_no" IS '服务商分成利润表,订单号';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."amount_order" IS '订单金额';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."servicefee_amount_sum" IS '手续费总额';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."split_proportion" IS '分成比例';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."actual_income" IS '实际所得';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."payment_time" IS '到账时间';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."servicer_no" IS '服务商uuid';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."currency_type" IS '币种';
COMMENT ON COLUMN "public"."servicer_profit_ledger"."order_type" IS '订单类型,1-存款,2-手机号取款;3-扫码取款';

-- ----------------------------
-- Table structure for servicer_terminal
-- ----------------------------
DROP TABLE IF EXISTS "public"."servicer_terminal";
CREATE TABLE "public"."servicer_terminal" (
  "terminal_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "servicer_no" uuid,
  "terminal_number" varchar(255) COLLATE "pg_catalog"."default",
  "pos_sn" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "is_delete" int2 DEFAULT 0,
  "use_status" int2 DEFAULT 1
)
;
COMMENT ON COLUMN "public"."servicer_terminal"."terminal_no" IS '收款设备主键';
COMMENT ON COLUMN "public"."servicer_terminal"."servicer_no" IS '运营商uuid';
COMMENT ON COLUMN "public"."servicer_terminal"."terminal_number" IS '收款设备号（终端编号）';
COMMENT ON COLUMN "public"."servicer_terminal"."pos_sn" IS 'pos机码';
COMMENT ON COLUMN "public"."servicer_terminal"."use_status" IS '0禁用1启用';

-- ----------------------------
-- Table structure for settle_servicer_hourly
-- ----------------------------
DROP TABLE IF EXISTS "public"."settle_servicer_hourly";
CREATE TABLE "public"."settle_servicer_hourly" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "start_time" timestamp(6),
  "finish_time" timestamp(6),
  "run_status" int8 DEFAULT 0,
  "begin_time" timestamp(6),
  "end_time" timestamp(6),
  "sum_income_usd" int8 DEFAULT 0,
  "sum_outgo_usd" int8 DEFAULT 0,
  "balance_usd" int8 DEFAULT 0,
  "delta_amount_usd" int8 DEFAULT 0,
  "balance_khr" int8 DEFAULT 0,
  "delta_amount_khr" int8 DEFAULT 0,
  "sum_income_khr" int8 DEFAULT 0,
  "sum_outgo_khr" int8 DEFAULT 0,
  "fbalance_usd" int8 DEFAULT 0,
  "fbalance_khr" int8 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."settle_servicer_hourly"."begin_time" IS '账期开始时间';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."end_time" IS '账期结束时间';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."sum_income_usd" IS '统计的入款';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."sum_outgo_usd" IS '统计的出款';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."balance_usd" IS '已核实金额';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."delta_amount_usd" IS '未核实金额';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."balance_khr" IS '已核实金额';
COMMENT ON COLUMN "public"."settle_servicer_hourly"."delta_amount_khr" IS '未核实金额';
COMMENT ON TABLE "public"."settle_servicer_hourly" IS '服务商每小时对账';

-- ----------------------------
-- Table structure for settle_vaccount_balance_hourly
-- ----------------------------
DROP TABLE IF EXISTS "public"."settle_vaccount_balance_hourly";
CREATE TABLE "public"."settle_vaccount_balance_hourly" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "vaccount_no" uuid NOT NULL,
  "balance" int8 DEFAULT 0,
  "frozen_balance" int8 DEFAULT 0,
  "create_time" timestamp(6)
)
;
COMMENT ON TABLE "public"."settle_vaccount_balance_hourly" IS '每小时虚帐余额';

-- ----------------------------
-- Table structure for sms_send_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."sms_send_record";
CREATE TABLE "public"."sms_send_record" (
  "id" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "msgid" varchar COLLATE "pg_catalog"."default",
  "account" varchar(255) COLLATE "pg_catalog"."default",
  "business" varchar(64) COLLATE "pg_catalog"."default",
  "mobile" varchar(255) COLLATE "pg_catalog"."default",
  "msg" varchar(255) COLLATE "pg_catalog"."default",
  "status" int4,
  "created_at" timestamp(6)
)
;
COMMENT ON COLUMN "public"."sms_send_record"."status" IS '0-成功,1-失败';

-- ----------------------------
-- Table structure for transfer_order
-- ----------------------------
DROP TABLE IF EXISTS "public"."transfer_order";
CREATE TABLE "public"."transfer_order" (
  "log_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "from_vaccount_no" uuid,
  "to_vaccount_no" uuid,
  "amount" int8 DEFAULT 0,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "finish_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "order_status" int8 DEFAULT 0,
  "balance_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "exchange_type" int8 DEFAULT 0,
  "fees" int8 DEFAULT 0,
  "payment_type" int2 DEFAULT 2,
  "is_count" int2,
  "modify_time" timestamp(6),
  "ree_rate" int8,
  "real_amount" int8,
  "ip" varchar(255) COLLATE "pg_catalog"."default",
  "lat" varchar(255) COLLATE "pg_catalog"."default",
  "lng" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."transfer_order"."log_no" IS '虚拟账户转账日志主键';
COMMENT ON COLUMN "public"."transfer_order"."order_status" IS '订单状态';
COMMENT ON COLUMN "public"."transfer_order"."balance_type" IS '金额类型（usd美金，khr瑞尔）';
COMMENT ON COLUMN "public"."transfer_order"."exchange_type" IS '兑换类型 0-扫码1-支付';
COMMENT ON COLUMN "public"."transfer_order"."fees" IS '手续费';
COMMENT ON COLUMN "public"."transfer_order"."payment_type" IS '付款方式;1-现金,2-余额';
COMMENT ON COLUMN "public"."transfer_order"."is_count" IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';
COMMENT ON COLUMN "public"."transfer_order"."ree_rate" IS '手续费率';
COMMENT ON COLUMN "public"."transfer_order"."real_amount" IS '实际到账金额';
COMMENT ON COLUMN "public"."transfer_order"."lat" IS '维度';
COMMENT ON COLUMN "public"."transfer_order"."lng" IS '经度';

-- ----------------------------
-- Table structure for url
-- ----------------------------
DROP TABLE IF EXISTS "public"."url";
CREATE TABLE "public"."url" (
  "url_uid" uuid NOT NULL,
  "url_name" varchar(255) COLLATE "pg_catalog"."default",
  "url" varchar(255) COLLATE "pg_catalog"."default",
  "parent_uid" uuid,
  "title" varchar(255) COLLATE "pg_catalog"."default",
  "icon" varchar(255) COLLATE "pg_catalog"."default",
  "component_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "component_path" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "redirect" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "idx" int4,
  "is_hidden" int4,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone
)
;
COMMENT ON COLUMN "public"."url"."url_name" IS 'url名';
COMMENT ON COLUMN "public"."url"."url" IS 'url';
COMMENT ON COLUMN "public"."url"."parent_uid" IS '父链接';
COMMENT ON COLUMN "public"."url"."title" IS '标题';
COMMENT ON COLUMN "public"."url"."icon" IS '图标';
COMMENT ON COLUMN "public"."url"."component_name" IS '组件名';
COMMENT ON COLUMN "public"."url"."component_path" IS '组件路径';
COMMENT ON COLUMN "public"."url"."redirect" IS '默认跳转';
COMMENT ON COLUMN "public"."url"."idx" IS '顺序';
COMMENT ON COLUMN "public"."url"."is_hidden" IS '隐藏';
COMMENT ON COLUMN "public"."url"."create_time" IS '创建时间';

-- ----------------------------
-- Table structure for vaccount
-- ----------------------------
DROP TABLE IF EXISTS "public"."vaccount";
CREATE TABLE "public"."vaccount" (
  "vaccount_no" uuid NOT NULL,
  "account_no" uuid NOT NULL,
  "va_type" int8,
  "balance" int8 DEFAULT 0,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0,
  "use_status" int4 DEFAULT 1,
  "delete_time" timestamp(6),
  "update_time" timestamp(6),
  "frozen_balance" int8 DEFAULT 0,
  "balance_type" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "modify_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."vaccount"."vaccount_no" IS '虚拟账户';
COMMENT ON COLUMN "public"."vaccount"."account_no" IS '账号';
COMMENT ON COLUMN "public"."vaccount"."va_type" IS '虚拟账户类型';
COMMENT ON COLUMN "public"."vaccount"."balance" IS '余额';
COMMENT ON COLUMN "public"."vaccount"."use_status" IS '0禁用';
COMMENT ON COLUMN "public"."vaccount"."delete_time" IS '销户时间';
COMMENT ON COLUMN "public"."vaccount"."update_time" IS '最后修改时间';
COMMENT ON COLUMN "public"."vaccount"."frozen_balance" IS '冻结金额';
COMMENT ON COLUMN "public"."vaccount"."balance_type" IS '余额类型（usd美金，khr瑞尔）';

-- ----------------------------
-- Table structure for wf_proc_running
-- ----------------------------
DROP TABLE IF EXISTS "public"."wf_proc_running";
CREATE TABLE "public"."wf_proc_running" (
  "running_no" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "process_no" uuid,
  "current_step" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "create_time" timestamp(6),
  "run_status" int4 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."wf_proc_running"."current_step" IS '当前执行到的步骤';

-- ----------------------------
-- Table structure for wf_process
-- ----------------------------
DROP TABLE IF EXISTS "public"."wf_process";
CREATE TABLE "public"."wf_process" (
  "process_no" uuid NOT NULL,
  "process_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "execute_no" uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
  "create_time" timestamp(6),
  "execute_status" int4,
  "steps" json
)
;
COMMENT ON COLUMN "public"."wf_process"."process_no" IS '流程号';
COMMENT ON COLUMN "public"."wf_process"."process_name" IS '流程名';
COMMENT ON COLUMN "public"."wf_process"."execute_no" IS '执行时间号';
COMMENT ON COLUMN "public"."wf_process"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."wf_process"."execute_status" IS '执行状态';
COMMENT ON COLUMN "public"."wf_process"."steps" IS '步骤集合';

-- ----------------------------
-- Table structure for wf_step
-- ----------------------------
DROP TABLE IF EXISTS "public"."wf_step";
CREATE TABLE "public"."wf_step" (
  "step_no" uuid NOT NULL,
  "step_name" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "func_name" varchar(512) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p1" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p2" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p3" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p4" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p5" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p6" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p7" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p8" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p9" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "p10" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6),
  "is_delete" int2 DEFAULT 0
)
;

-- ----------------------------
-- Table structure for writeoff
-- ----------------------------
DROP TABLE IF EXISTS "public"."writeoff";
CREATE TABLE "public"."writeoff" (
  "code" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "income_order_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "outgo_order_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "create_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "finish_time" timestamp(6) DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
  "use_status" int2 DEFAULT 0,
  "transfer_order_no" varchar(255) COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "send_phone" varchar(255) COLLATE "pg_catalog"."default",
  "recv_phone" varchar(255) COLLATE "pg_catalog"."default",
  "modify_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."writeoff"."code" IS '核销码';
COMMENT ON COLUMN "public"."writeoff"."use_status" IS '状态1-初始状态.2-已使用状态';

-- ----------------------------
-- Table structure for xlsx_file_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."xlsx_file_log";
CREATE TABLE "public"."xlsx_file_log" (
  "xlsx_task_no" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "create_time" timestamp(6),
  "account_no" uuid,
  "file_type" int2,
  "query_str" varchar(255) COLLATE "pg_catalog"."default",
  "role_type" varchar(255) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."xlsx_file_log"."file_type" IS '订单类型（1兑2存3取4转5收）';
COMMENT ON COLUMN "public"."xlsx_file_log"."query_str" IS '查询条件';

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."servicer_count_list_id_seq"
OWNED BY "public"."servicer_count_list"."id";
SELECT setval('"public"."servicer_count_list_id_seq"', 194, true);

-- ----------------------------
-- Uniques structure for table account
-- ----------------------------
ALTER TABLE "public"."account" ADD CONSTRAINT "account_gen_key_key" UNIQUE ("gen_key");

-- ----------------------------
-- Primary Key structure for table account
-- ----------------------------
ALTER TABLE "public"."account" ADD CONSTRAINT "account_pkey" PRIMARY KEY ("uid");

-- ----------------------------
-- Primary Key structure for table account_collect
-- ----------------------------
ALTER TABLE "public"."account_collect" ADD CONSTRAINT "account_collect_pkey" PRIMARY KEY ("account_collect_no");

-- ----------------------------
-- Primary Key structure for table adminlog
-- ----------------------------
ALTER TABLE "public"."adminlog" ADD CONSTRAINT "adminlog_pkey" PRIMARY KEY ("log_uid");

-- ----------------------------
-- Primary Key structure for table agreement
-- ----------------------------
ALTER TABLE "public"."agreement" ADD CONSTRAINT "agreement_privacy_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table app_version
-- ----------------------------
ALTER TABLE "public"."app_version" ADD CONSTRAINT "version_vsType_system_key" UNIQUE ("version", "vs_type", "system");

-- ----------------------------
-- Primary Key structure for table app_version_file_log
-- ----------------------------
ALTER TABLE "public"."app_version_file_log" ADD CONSTRAINT "app_file_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table billing_details_results
-- ----------------------------
ALTER TABLE "public"."billing_details_results" ADD CONSTRAINT "billing_details_results_pkey" PRIMARY KEY ("bill_no");

-- ----------------------------
-- Primary Key structure for table card
-- ----------------------------
ALTER TABLE "public"."card" ADD CONSTRAINT "servicer_card_pack_pkey" PRIMARY KEY ("card_no");

-- ----------------------------
-- Primary Key structure for table cashier
-- ----------------------------
ALTER TABLE "public"."cashier" ADD CONSTRAINT "cashier_pkey" PRIMARY KEY ("uid");

-- ----------------------------
-- Primary Key structure for table channel
-- ----------------------------
ALTER TABLE "public"."channel" ADD CONSTRAINT "channel_pkey" PRIMARY KEY ("channel_no");

-- ----------------------------
-- Primary Key structure for table collection_order
-- ----------------------------
ALTER TABLE "public"."collection_order" ADD CONSTRAINT "transfer_order_copy1_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table common_help
-- ----------------------------
ALTER TABLE "public"."common_help" ADD CONSTRAINT "common_help_pkey" PRIMARY KEY ("help_no");

-- ----------------------------
-- Primary Key structure for table consultation_config
-- ----------------------------
ALTER TABLE "public"."consultation_config" ADD CONSTRAINT "consultation_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table cust
-- ----------------------------
ALTER TABLE "public"."cust" ADD CONSTRAINT "cust_pkey" PRIMARY KEY ("cust_no");

-- ----------------------------
-- Primary Key structure for table dict_acc_title
-- ----------------------------
ALTER TABLE "public"."dict_acc_title" ADD CONSTRAINT "dict_acc_title_pkey" PRIMARY KEY ("title_no");

-- ----------------------------
-- Primary Key structure for table dict_account_type
-- ----------------------------
ALTER TABLE "public"."dict_account_type" ADD CONSTRAINT "dict_account_type_pkey" PRIMARY KEY ("account_type");

-- ----------------------------
-- Primary Key structure for table dict_bankname
-- ----------------------------
ALTER TABLE "public"."dict_bankname" ADD CONSTRAINT "dict_bankname_pkey" PRIMARY KEY ("bank_name");

-- ----------------------------
-- Primary Key structure for table dict_bin_bankname
-- ----------------------------
ALTER TABLE "public"."dict_bin_bankname" ADD CONSTRAINT "dict_bin_bankname_pkey" PRIMARY KEY ("bin_code");

-- ----------------------------
-- Primary Key structure for table dict_images
-- ----------------------------
ALTER TABLE "public"."dict_images" ADD CONSTRAINT "dict_images_pkey" PRIMARY KEY ("image_id");

-- ----------------------------
-- Primary Key structure for table dict_org_abbr
-- ----------------------------
ALTER TABLE "public"."dict_org_abbr" ADD CONSTRAINT "dict_org_abbr_pkey" PRIMARY KEY ("org_code");

-- ----------------------------
-- Primary Key structure for table dict_province
-- ----------------------------
ALTER TABLE "public"."dict_province" ADD CONSTRAINT "dict_province_pkey" PRIMARY KEY ("province_code");

-- ----------------------------
-- Primary Key structure for table dict_region_bank
-- ----------------------------
ALTER TABLE "public"."dict_region_bank" ADD CONSTRAINT "dict_region_bank_copy1_pkey" PRIMARY KEY ("code");

-- ----------------------------
-- Primary Key structure for table dict_vatype
-- ----------------------------
ALTER TABLE "public"."dict_vatype" ADD CONSTRAINT "dict_vaccount_pkey" PRIMARY KEY ("va_type");

-- ----------------------------
-- Primary Key structure for table exchange_order
-- ----------------------------
ALTER TABLE "public"."exchange_order" ADD CONSTRAINT "exchange_order_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table func_config
-- ----------------------------
ALTER TABLE "public"."func_config" ADD CONSTRAINT "func_config_pkey" PRIMARY KEY ("func_no");

-- ----------------------------
-- Primary Key structure for table gen_code
-- ----------------------------
ALTER TABLE "public"."gen_code" ADD CONSTRAINT "gen_code_pkey" PRIMARY KEY ("gen_key");

-- ----------------------------
-- Primary Key structure for table global_param
-- ----------------------------
ALTER TABLE "public"."global_param" ADD CONSTRAINT "global_param_pkey" PRIMARY KEY ("param_key");

-- ----------------------------
-- Primary Key structure for table headquarters_profit
-- ----------------------------
ALTER TABLE "public"."headquarters_profit" ADD CONSTRAINT "servicefee_order_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table headquarters_profit_withdraw
-- ----------------------------
ALTER TABLE "public"."headquarters_profit_withdraw" ADD CONSTRAINT "headquarters_profit_withdraw_pkey" PRIMARY KEY ("order_no");

-- ----------------------------
-- Primary Key structure for table income_order
-- ----------------------------
ALTER TABLE "public"."income_order" ADD CONSTRAINT "income_log_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table income_type
-- ----------------------------
ALTER TABLE "public"."income_type" ADD CONSTRAINT "income_type_pkey" PRIMARY KEY ("income_type");

-- ----------------------------
-- Primary Key structure for table lang
-- ----------------------------
ALTER TABLE "public"."lang" ADD CONSTRAINT "lang_pkey" PRIMARY KEY ("key");

-- ----------------------------
-- Primary Key structure for table log_account
-- ----------------------------
ALTER TABLE "public"."log_account" ADD CONSTRAINT "log_account_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table log_account_web
-- ----------------------------
ALTER TABLE "public"."log_account_web" ADD CONSTRAINT "log_account_copy1_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table log_card
-- ----------------------------
ALTER TABLE "public"."log_card" ADD CONSTRAINT "log_card_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table log_login
-- ----------------------------
ALTER TABLE "public"."log_login" ADD CONSTRAINT "log_login_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table log_to_headquarters
-- ----------------------------
ALTER TABLE "public"."log_to_headquarters" ADD CONSTRAINT "log_to_headquarters_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table log_to_servicer
-- ----------------------------
ALTER TABLE "public"."log_to_servicer" ADD CONSTRAINT "log_to_servicer_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table log_vaccount
-- ----------------------------
ALTER TABLE "public"."log_vaccount" ADD CONSTRAINT "log_vaccount_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table login_token
-- ----------------------------
ALTER TABLE "public"."login_token" ADD CONSTRAINT "login_token_pkey" PRIMARY KEY ("acc_no");

-- ----------------------------
-- Primary Key structure for table outgo_order
-- ----------------------------
ALTER TABLE "public"."outgo_order" ADD CONSTRAINT "outgo_order_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table outgo_type
-- ----------------------------
ALTER TABLE "public"."outgo_type" ADD CONSTRAINT "income_type_copy1_pkey" PRIMARY KEY ("outgo_type");

-- ----------------------------
-- Primary Key structure for table platform_config
-- ----------------------------
ALTER TABLE "public"."platform_config" ADD CONSTRAINT "platform_config_pkey" PRIMARY KEY ("account_uid");

-- ----------------------------
-- Primary Key structure for table rela_acc_iden
-- ----------------------------
ALTER TABLE "public"."rela_acc_iden" ADD CONSTRAINT "rela_acc_iden_pkey" PRIMARY KEY ("account_no", "account_type");

-- ----------------------------
-- Primary Key structure for table rela_account_role
-- ----------------------------
ALTER TABLE "public"."rela_account_role" ADD CONSTRAINT "rela_tag_access_pkey" PRIMARY KEY ("rela_uid");

-- ----------------------------
-- Primary Key structure for table rela_imei_pubkey
-- ----------------------------
ALTER TABLE "public"."rela_imei_pubkey" ADD CONSTRAINT "rela_imei_pubkey_pkey" PRIMARY KEY ("rela_no");

-- ----------------------------
-- Primary Key structure for table rela_role_url
-- ----------------------------
ALTER TABLE "public"."rela_role_url" ADD CONSTRAINT "rela_role_url_pkey" PRIMARY KEY ("rela_uid");

-- ----------------------------
-- Primary Key structure for table role
-- ----------------------------
ALTER TABLE "public"."role" ADD CONSTRAINT "permission_tag_pkey" PRIMARY KEY ("role_no");

-- ----------------------------
-- Primary Key structure for table servicer
-- ----------------------------
ALTER TABLE "public"."servicer" ADD CONSTRAINT "merchant_pkey" PRIMARY KEY ("servicer_no");

-- ----------------------------
-- Uniques structure for table servicer_count
-- ----------------------------
ALTER TABLE "public"."servicer_count" ADD CONSTRAINT "servicer_no_balance_type_key" UNIQUE ("servicer_no", "currency_type");

-- ----------------------------
-- Uniques structure for table servicer_count_list
-- ----------------------------
ALTER TABLE "public"."servicer_count_list" ADD CONSTRAINT "servicer_date_currency" UNIQUE ("servicer_no", "currency_type", "dates");
COMMENT ON CONSTRAINT "servicer_date_currency" ON "public"."servicer_count_list" IS '服务商同一天同一币种只能有一条记录';

-- ----------------------------
-- Primary Key structure for table servicer_img
-- ----------------------------
ALTER TABLE "public"."servicer_img" ADD CONSTRAINT "merchant_img_pkey" PRIMARY KEY ("servicer_img_no");

-- ----------------------------
-- Primary Key structure for table servicer_profit_ledger
-- ----------------------------
ALTER TABLE "public"."servicer_profit_ledger" ADD CONSTRAINT "servicer_profit_ledger_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table servicer_terminal
-- ----------------------------
ALTER TABLE "public"."servicer_terminal" ADD CONSTRAINT "income_terminal_pkey" PRIMARY KEY ("terminal_no");

-- ----------------------------
-- Primary Key structure for table settle_servicer_hourly
-- ----------------------------
ALTER TABLE "public"."settle_servicer_hourly" ADD CONSTRAINT "settle_servicer_hourly_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table settle_vaccount_balance_hourly
-- ----------------------------
ALTER TABLE "public"."settle_vaccount_balance_hourly" ADD CONSTRAINT "log_settle_vaccount_balance_pkey" PRIMARY KEY ("log_no", "vaccount_no");

-- ----------------------------
-- Primary Key structure for table sms_send_record
-- ----------------------------
ALTER TABLE "public"."sms_send_record" ADD CONSTRAINT "sms_send_record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table transfer_order
-- ----------------------------
ALTER TABLE "public"."transfer_order" ADD CONSTRAINT "transfer_order_pkey" PRIMARY KEY ("log_no");

-- ----------------------------
-- Primary Key structure for table url
-- ----------------------------
ALTER TABLE "public"."url" ADD CONSTRAINT "url_pkey" PRIMARY KEY ("url_uid");

-- ----------------------------
-- Primary Key structure for table vaccount
-- ----------------------------
ALTER TABLE "public"."vaccount" ADD CONSTRAINT "vaccount_pkey" PRIMARY KEY ("vaccount_no");

-- ----------------------------
-- Primary Key structure for table wf_proc_running
-- ----------------------------
ALTER TABLE "public"."wf_proc_running" ADD CONSTRAINT "wf_proc_running_pkey" PRIMARY KEY ("running_no");

-- ----------------------------
-- Primary Key structure for table wf_process
-- ----------------------------
ALTER TABLE "public"."wf_process" ADD CONSTRAINT "wf_process_pkey" PRIMARY KEY ("process_no");

-- ----------------------------
-- Primary Key structure for table wf_step
-- ----------------------------
ALTER TABLE "public"."wf_step" ADD CONSTRAINT "wf_step_pkey" PRIMARY KEY ("step_no");

-- ----------------------------
-- Primary Key structure for table writeoff
-- ----------------------------
ALTER TABLE "public"."writeoff" ADD CONSTRAINT "writeoff_pkey" PRIMARY KEY ("code");

-- ----------------------------
-- Primary Key structure for table xlsx_file_log
-- ----------------------------
ALTER TABLE "public"."xlsx_file_log" ADD CONSTRAINT "xlsx_task_pkey" PRIMARY KEY ("xlsx_task_no");
