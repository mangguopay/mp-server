--
-- PostgreSQL database dump
--

-- Dumped from database version 11.6
-- Dumped by pg_dump version 11.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.xlsx_file_log DROP CONSTRAINT xlsx_task_pkey;
ALTER TABLE ONLY public.writeoff DROP CONSTRAINT writeoff_pkey;
ALTER TABLE ONLY public.wf_step DROP CONSTRAINT wf_step_pkey;
ALTER TABLE ONLY public.wf_process DROP CONSTRAINT wf_process_pkey;
ALTER TABLE ONLY public.wf_proc_running DROP CONSTRAINT wf_proc_running_pkey;
ALTER TABLE ONLY public.app_version DROP CONSTRAINT "version_vsType_system_key";
ALTER TABLE ONLY public.vaccount DROP CONSTRAINT vaccount_pkey;
ALTER TABLE ONLY public.url DROP CONSTRAINT url_pkey;
ALTER TABLE ONLY public.transfer_order DROP CONSTRAINT transfer_order_pkey;
ALTER TABLE ONLY public.collection_order DROP CONSTRAINT transfer_order_copy1_pkey;
ALTER TABLE ONLY public.sms_send_record DROP CONSTRAINT sms_send_record_pkey;
ALTER TABLE ONLY public.settle_servicer_hourly DROP CONSTRAINT settle_servicer_hourly_pkey;
ALTER TABLE ONLY public.servicer_profit_ledger DROP CONSTRAINT servicer_profit_ledger_pkey;
ALTER TABLE ONLY public.servicer_count DROP CONSTRAINT servicer_no_balance_type_key;
ALTER TABLE ONLY public.servicer_count_list DROP CONSTRAINT servicer_date_currency;
ALTER TABLE ONLY public.card DROP CONSTRAINT servicer_card_pack_pkey;
ALTER TABLE ONLY public.headquarters_profit DROP CONSTRAINT servicefee_order_pkey;
ALTER TABLE ONLY public.rela_account_role DROP CONSTRAINT rela_tag_access_pkey;
ALTER TABLE ONLY public.rela_role_url DROP CONSTRAINT rela_role_url_pkey;
ALTER TABLE ONLY public.rela_imei_pubkey DROP CONSTRAINT rela_imei_pubkey_pkey;
ALTER TABLE ONLY public.rela_acc_iden DROP CONSTRAINT rela_acc_iden_pkey;
ALTER TABLE ONLY public.platform_config DROP CONSTRAINT platform_config_pkey;
ALTER TABLE ONLY public.role DROP CONSTRAINT permission_tag_pkey;
ALTER TABLE ONLY public.outgo_order DROP CONSTRAINT outgo_order_pkey;
ALTER TABLE ONLY public.servicer DROP CONSTRAINT merchant_pkey;
ALTER TABLE ONLY public.servicer_img DROP CONSTRAINT merchant_img_pkey;
ALTER TABLE ONLY public.login_token DROP CONSTRAINT login_token_pkey;
ALTER TABLE ONLY public.log_vaccount DROP CONSTRAINT log_vaccount_pkey;
ALTER TABLE ONLY public.log_to_servicer DROP CONSTRAINT log_to_servicer_pkey;
ALTER TABLE ONLY public.log_to_headquarters DROP CONSTRAINT log_to_headquarters_pkey;
ALTER TABLE ONLY public.settle_vaccount_balance_hourly DROP CONSTRAINT log_settle_vaccount_balance_pkey;
ALTER TABLE ONLY public.log_login DROP CONSTRAINT log_login_pkey;
ALTER TABLE ONLY public.log_card DROP CONSTRAINT log_card_pkey;
ALTER TABLE ONLY public.log_account_web DROP CONSTRAINT log_account_web_pkey;
ALTER TABLE ONLY public.log_account DROP CONSTRAINT log_account_pkey;
ALTER TABLE ONLY public.lang DROP CONSTRAINT lang_pkey;
ALTER TABLE ONLY public.income_type DROP CONSTRAINT income_type_pkey;
ALTER TABLE ONLY public.outgo_type DROP CONSTRAINT income_type_copy1_pkey;
ALTER TABLE ONLY public.servicer_terminal DROP CONSTRAINT income_terminal_pkey;
ALTER TABLE ONLY public.income_order DROP CONSTRAINT income_log_pkey;
ALTER TABLE ONLY public.headquarters_profit_withdraw DROP CONSTRAINT headquarters_profit_withdraw_pkey;
ALTER TABLE ONLY public.global_param DROP CONSTRAINT global_param_pkey;
ALTER TABLE ONLY public.gen_code DROP CONSTRAINT gen_code_pkey;
ALTER TABLE ONLY public.func_config DROP CONSTRAINT func_config_pkey;
ALTER TABLE ONLY public.exchange_order DROP CONSTRAINT exchange_order_pkey;
ALTER TABLE ONLY public.dict_vatype DROP CONSTRAINT dict_vaccount_pkey;
ALTER TABLE ONLY public.dict_region_bank DROP CONSTRAINT dict_region_bank_copy1_pkey;
ALTER TABLE ONLY public.dict_province DROP CONSTRAINT dict_province_pkey;
ALTER TABLE ONLY public.dict_org_abbr DROP CONSTRAINT dict_org_abbr_pkey;
ALTER TABLE ONLY public.dict_images DROP CONSTRAINT dict_images_pkey;
ALTER TABLE ONLY public.dict_bin_bankname DROP CONSTRAINT dict_bin_bankname_pkey;
ALTER TABLE ONLY public.dict_bankname DROP CONSTRAINT dict_bankname_pkey;
ALTER TABLE ONLY public.dict_account_type DROP CONSTRAINT dict_account_type_pkey;
ALTER TABLE ONLY public.dict_acc_title DROP CONSTRAINT dict_acc_title_pkey;
ALTER TABLE ONLY public.cust DROP CONSTRAINT cust_pkey;
ALTER TABLE ONLY public.consultation_config DROP CONSTRAINT consultation_config_pkey;
ALTER TABLE ONLY public.common_help DROP CONSTRAINT common_help_pkey;
ALTER TABLE ONLY public.channel DROP CONSTRAINT channel_pkey;
ALTER TABLE ONLY public.cashier DROP CONSTRAINT cashier_pkey;
ALTER TABLE ONLY public.billing_details_results DROP CONSTRAINT billing_details_results_pkey;
ALTER TABLE ONLY public.app_version_file_log DROP CONSTRAINT app_file_log_pkey;
ALTER TABLE ONLY public.agreement DROP CONSTRAINT agreement_privacy_pkey;
ALTER TABLE ONLY public.adminlog DROP CONSTRAINT adminlog_pkey;
ALTER TABLE ONLY public.account DROP CONSTRAINT account_pkey;
ALTER TABLE ONLY public.account DROP CONSTRAINT account_gen_key_key;
ALTER TABLE ONLY public.account_collect DROP CONSTRAINT account_collect_pkey;
ALTER TABLE public.servicer_count_list ALTER COLUMN id DROP DEFAULT;
DROP TABLE public.xlsx_file_log;
DROP TABLE public.writeoff;
DROP TABLE public.wf_step;
DROP TABLE public.wf_process;
DROP TABLE public.wf_proc_running;
DROP TABLE public.vaccount;
DROP TABLE public.url;
DROP TABLE public.transfer_order;
DROP TABLE public.sms_send_record;
DROP TABLE public.settle_vaccount_balance_hourly;
DROP TABLE public.settle_servicer_hourly;
DROP TABLE public.servicer_terminal;
DROP TABLE public.servicer_profit_ledger;
DROP TABLE public.servicer_img;
DROP SEQUENCE public.servicer_count_list_id_seq;
DROP TABLE public.servicer_count_list;
DROP TABLE public.servicer_count;
DROP TABLE public.servicer;
DROP TABLE public.role;
DROP TABLE public.risk_result;
DROP TABLE public.rela_role_url;
DROP TABLE public.rela_imei_pubkey;
DROP TABLE public.rela_account_role;
DROP TABLE public.rela_acc_iden;
DROP TABLE public.platform_config;
DROP TABLE public.outgo_type;
DROP TABLE public.outgo_order;
DROP TABLE public.login_token;
DROP TABLE public.log_vaccount;
DROP TABLE public.log_to_servicer;
DROP TABLE public.log_to_headquarters;
DROP TABLE public.log_login;
DROP TABLE public.log_exchange_rate;
DROP TABLE public.log_card;
DROP TABLE public.log_app_messages;
DROP TABLE public.log_account_web;
DROP TABLE public.log_account;
DROP TABLE public.lang;
DROP TABLE public.income_type;
DROP TABLE public.income_ougo_config;
DROP TABLE public.income_order;
DROP TABLE public.headquarters_profit_withdraw;
DROP TABLE public.headquarters_profit_cashable;
DROP TABLE public.headquarters_profit;
DROP TABLE public.global_param;
DROP TABLE public.gen_code;
DROP TABLE public.func_config;
DROP TABLE public.exchange_order;
DROP TABLE public.dict_vatype;
DROP TABLE public.dict_region_bank;
DROP TABLE public.dict_region;
DROP TABLE public.dict_province;
DROP TABLE public.dict_org_abbr;
DROP TABLE public.dict_images;
DROP TABLE public.dict_bin_bankname;
DROP TABLE public.dict_bankname;
DROP TABLE public.dict_bank_abbr;
DROP TABLE public.dict_account_type;
DROP TABLE public.dict_acc_title;
DROP TABLE public.cust;
DROP TABLE public.consultation_config;
DROP TABLE public.common_help;
DROP TABLE public.collection_order;
DROP TABLE public.channel_servicer;
DROP TABLE public.channel;
DROP TABLE public.cashier;
DROP TABLE public.card;
DROP TABLE public.billing_details_results;
DROP TABLE public.app_version_file_log;
DROP TABLE public.app_version;
DROP TABLE public.agreement;
DROP TABLE public.adminlog;
DROP TABLE public.account_collect;
DROP TABLE public.account;
SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account (
    uid uuid NOT NULL,
    nickname character varying(255) DEFAULT ''::character varying,
    account character varying(255) DEFAULT ''::character varying NOT NULL,
    password character varying(255) DEFAULT ''::character varying NOT NULL,
    use_status bigint DEFAULT 1,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    drop_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    modify_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    update_time timestamp(0) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    phone character varying(255) DEFAULT ''::character varying,
    email character varying(512) DEFAULT ''::character varying,
    master_acc uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    is_delete smallint DEFAULT 0,
    usd_balance bigint DEFAULT 0,
    khr_balance bigint DEFAULT 0,
    gen_key character varying(255) DEFAULT ''::character varying,
    is_actived smallint DEFAULT 0,
    head_portrait_img_no character varying(255),
    last_login_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    is_first_login smallint DEFAULT 0,
    app_lang character varying(255),
    pos_lang character varying(255),
    country_code character varying(255),
    utm_source character varying(255)
);


ALTER TABLE public.account OWNER TO postgres;

--
-- Name: COLUMN account.uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.uid IS '用户id';


--
-- Name: COLUMN account.nickname; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.nickname IS '用户名';


--
-- Name: COLUMN account.account; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.account IS '账号';


--
-- Name: COLUMN account.password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.password IS '密码';


--
-- Name: COLUMN account.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.use_status IS '使用状态：0.禁用，1.正常';


--
-- Name: COLUMN account.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.create_time IS '创建时间';


--
-- Name: COLUMN account.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.modify_time IS '最后更新时间';


--
-- Name: COLUMN account.phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.phone IS '手机';


--
-- Name: COLUMN account.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.email IS '邮件';


--
-- Name: COLUMN account.master_acc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.master_acc IS '主账号';


--
-- Name: COLUMN account.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.is_delete IS '账号删除';


--
-- Name: COLUMN account.usd_balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.usd_balance IS '美金余额';


--
-- Name: COLUMN account.khr_balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.khr_balance IS '柬埔寨余额';


--
-- Name: COLUMN account.gen_key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.gen_key IS '注册生成的扫码';


--
-- Name: COLUMN account.is_actived; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.is_actived IS '是否激活;0-未激活,1-激活.';


--
-- Name: COLUMN account.head_portrait_img_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.head_portrait_img_no IS '头像图片url';


--
-- Name: COLUMN account.last_login_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.last_login_time IS '最后登录时间';


--
-- Name: COLUMN account.is_first_login; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.is_first_login IS '是否第一次登陆,0-是;1-否(这个字段只给pos端使用)';


--
-- Name: COLUMN account.country_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.country_code IS '手机号国家码';


--
-- Name: COLUMN account.utm_source; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account.utm_source IS '注册来源';


--
-- Name: account_collect; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.account_collect (
    account_no uuid NOT NULL,
    collect_account_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    collect_phone bigint,
    is_delete smallint DEFAULT 0,
    create_time timestamp(0) without time zone,
    modify_time timestamp(0) without time zone,
    account_collect_no character varying(255) NOT NULL
);


ALTER TABLE public.account_collect OWNER TO postgres;

--
-- Name: TABLE account_collect; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.account_collect IS '转账后，最近转账的人';


--
-- Name: COLUMN account_collect.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account_collect.account_no IS '账号';


--
-- Name: COLUMN account_collect.collect_account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account_collect.collect_account_no IS '收集的账号Uid';


--
-- Name: COLUMN account_collect.collect_phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account_collect.collect_phone IS '收集的转账手机号';


--
-- Name: COLUMN account_collect.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account_collect.modify_time IS '最后修改时间（用于排序）';


--
-- Name: COLUMN account_collect.account_collect_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.account_collect.account_collect_no IS '主键';


--
-- Name: adminlog; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.adminlog (
    log_uid uuid NOT NULL,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    url character varying(255) DEFAULT ''::character varying,
    param character varying(512) DEFAULT ''::character varying,
    op character varying(255) DEFAULT ''::character varying,
    op_type integer DEFAULT 0,
    op_acc_uid uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    during integer DEFAULT 0 NOT NULL,
    ip text DEFAULT ''::text NOT NULL,
    status_code text DEFAULT ''::text NOT NULL,
    response text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.adminlog OWNER TO postgres;

--
-- Name: TABLE adminlog; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.adminlog IS '运行日志';


--
-- Name: COLUMN adminlog.log_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.log_uid IS '日志uid';


--
-- Name: COLUMN adminlog.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.create_time IS '创建时间';


--
-- Name: COLUMN adminlog.url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.url IS 'url';


--
-- Name: COLUMN adminlog.param; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.param IS '参数';


--
-- Name: COLUMN adminlog.op; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.op IS '操作';


--
-- Name: COLUMN adminlog.op_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.op_type IS '操作类型';


--
-- Name: COLUMN adminlog.op_acc_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.adminlog.op_acc_uid IS '操作账号uid';


--
-- Name: agreement; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.agreement (
    id character varying(255) NOT NULL,
    text text,
    lang character varying(255),
    type smallint,
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    modify_time timestamp(6) without time zone,
    use_status smallint DEFAULT 1
);


ALTER TABLE public.agreement OWNER TO postgres;

--
-- Name: COLUMN agreement.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.id IS 'id';


--
-- Name: COLUMN agreement.text; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.text IS '协议内容';


--
-- Name: COLUMN agreement.lang; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.lang IS '语言(km、zh_CN、en)';


--
-- Name: COLUMN agreement.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.type IS '0 用户协议  1隐私协议';


--
-- Name: COLUMN agreement.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.is_delete IS '1删除';


--
-- Name: COLUMN agreement.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.modify_time IS '最后修改时间';


--
-- Name: COLUMN agreement.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.agreement.use_status IS '使用状态0禁用';


--
-- Name: app_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_version (
    v_id character varying NOT NULL,
    description text DEFAULT ''::text,
    version character varying(255),
    create_time timestamp(0) without time zone,
    update_time timestamp(0) without time zone,
    app_url character varying(2000),
    vs_code character varying(255),
    vs_type smallint DEFAULT 0,
    is_force smallint DEFAULT 0,
    system smallint DEFAULT 0,
    is_delete smallint DEFAULT 0,
    account_uid uuid,
    note text DEFAULT ''::text,
    status smallint DEFAULT 1
);


ALTER TABLE public.app_version OWNER TO postgres;

--
-- Name: COLUMN app_version.v_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.v_id IS 'app版本';


--
-- Name: COLUMN app_version.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.description IS '版本描述';


--
-- Name: COLUMN app_version.version; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.version IS '版本';


--
-- Name: COLUMN app_version.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.create_time IS '生成时间';


--
-- Name: COLUMN app_version.update_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.update_time IS '更改时间';


--
-- Name: COLUMN app_version.app_url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.app_url IS '下载地址';


--
-- Name: COLUMN app_version.vs_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.vs_code IS '版本码,做版本比较使用';


--
-- Name: COLUMN app_version.vs_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.vs_type IS '0-app包;1-pos包';


--
-- Name: COLUMN app_version.is_force; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.is_force IS '是否强制更新;0-否,1-是';


--
-- Name: COLUMN app_version.system; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.system IS '0-android系统;1-ios系统';


--
-- Name: COLUMN app_version.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.is_delete IS '1删除';


--
-- Name: COLUMN app_version.account_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.account_uid IS '上传版本的账户uid';


--
-- Name: COLUMN app_version.note; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.note IS '备注';


--
-- Name: COLUMN app_version.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.app_version.status IS '是否禁用(0禁用1启用)';


--
-- Name: app_version_file_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_version_file_log (
    id character varying(255) NOT NULL,
    create_time timestamp(0) without time zone,
    account_no uuid,
    file_name character varying(255)
);


ALTER TABLE public.app_version_file_log OWNER TO postgres;

--
-- Name: billing_details_results; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.billing_details_results (
    create_time timestamp(0) without time zone,
    bill_no character varying(32) NOT NULL,
    amount bigint DEFAULT 0,
    currency_type character varying(255),
    bill_type smallint,
    account_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    account_type character varying(255) DEFAULT ''::character varying,
    order_no character varying(255) DEFAULT ''::character varying,
    balance bigint DEFAULT 0,
    order_status smallint,
    modify_time timestamp(0) without time zone,
    servicer_no uuid,
    op_acc_no uuid
);


ALTER TABLE public.billing_details_results OWNER TO postgres;

--
-- Name: TABLE billing_details_results; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.billing_details_results IS '服务商的总账流水-非虚拟';


--
-- Name: COLUMN billing_details_results.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.create_time IS '记账时间';


--
-- Name: COLUMN billing_details_results.bill_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.bill_no IS '账单流水号';


--
-- Name: COLUMN billing_details_results.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.amount IS '金额';


--
-- Name: COLUMN billing_details_results.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.currency_type IS '币种(usd,khr)';


--
-- Name: COLUMN billing_details_results.bill_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.bill_type IS '账单类型（1:存款;2:取款;3:收益;4:充值;5-提现）';


--
-- Name: COLUMN billing_details_results.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.account_no IS '账号uid';


--
-- Name: COLUMN billing_details_results.order_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.order_no IS '订单号';


--
-- Name: COLUMN billing_details_results.balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.balance IS '余额(暂时放这，后面确认完全无用再删)';


--
-- Name: COLUMN billing_details_results.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.order_status IS '订单状态（1-初始化;2-等待;3-已支付;4-失败(超时);5-待确认;6-取消）';


--
-- Name: COLUMN billing_details_results.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.billing_details_results.modify_time IS '更改时间';


--
-- Name: card; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.card (
    card_no uuid NOT NULL,
    account_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    channel_no uuid,
    name character varying(255) DEFAULT ''::character varying,
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    card_number character varying(255) DEFAULT ''::character varying,
    balance_type character varying(12) DEFAULT ''::character varying,
    is_defalut smallint DEFAULT 0,
    collect_status smallint DEFAULT 1,
    audit_status smallint DEFAULT 0,
    note text DEFAULT ''::text,
    modify_time timestamp(6) without time zone
);


ALTER TABLE public.card OWNER TO postgres;

--
-- Name: COLUMN card.card_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.card_no IS '卡uid';


--
-- Name: COLUMN card.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.account_no IS '账号uuid';


--
-- Name: COLUMN card.channel_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.channel_no IS '渠道';


--
-- Name: COLUMN card.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.name IS '收款人名';


--
-- Name: COLUMN card.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.create_time IS '创建时间';


--
-- Name: COLUMN card.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.is_delete IS '是否删除';


--
-- Name: COLUMN card.card_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.card_number IS '卡号';


--
-- Name: COLUMN card.balance_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.balance_type IS '币种（usd美金，khr瑞尔）';


--
-- Name: COLUMN card.is_defalut; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.is_defalut IS '1:推荐收款卡';


--
-- Name: COLUMN card.collect_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.collect_status IS '收款状态（0禁用，1启用）';


--
-- Name: COLUMN card.audit_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.audit_status IS '审核状态,0-未审核,1-审核通过';


--
-- Name: COLUMN card.note; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.card.note IS '备注';


--
-- Name: cashier; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cashier (
    uid uuid NOT NULL,
    name character varying(255),
    servicer_no uuid,
    is_delete smallint DEFAULT 0,
    create_time timestamp(6) without time zone,
    op_password character varying(255) DEFAULT ''::character varying,
    modify_time timestamp(6) without time zone
);


ALTER TABLE public.cashier OWNER TO postgres;

--
-- Name: COLUMN cashier.uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cashier.uid IS '操作员';


--
-- Name: COLUMN cashier.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cashier.name IS '姓名';


--
-- Name: COLUMN cashier.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cashier.servicer_no IS '服务商';


--
-- Name: COLUMN cashier.op_password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cashier.op_password IS '操作密码';


--
-- Name: channel; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.channel (
    channel_no uuid NOT NULL,
    channel_name character varying(255) DEFAULT ''::character varying,
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    note text DEFAULT ''::text,
    idx bigint DEFAULT 0,
    is_recom smallint DEFAULT 0,
    channel_type smallint DEFAULT 0,
    currency_type character varying(255) DEFAULT ''::character varying,
    use_status smallint DEFAULT 0,
    logo_img_no character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.channel OWNER TO postgres;

--
-- Name: COLUMN channel.channel_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.channel_name IS '渠道名称';


--
-- Name: COLUMN channel.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.create_time IS '创建时间';


--
-- Name: COLUMN channel.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.is_delete IS '是否删除';


--
-- Name: COLUMN channel.note; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.note IS '备注';


--
-- Name: COLUMN channel.idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.idx IS '序列';


--
-- Name: COLUMN channel.is_recom; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.is_recom IS '是否推荐(1-推荐，0-不推荐)';


--
-- Name: COLUMN channel.channel_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.channel_type IS '0-通用,1-用户,总部;2-pos';


--
-- Name: COLUMN channel.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.currency_type IS '币种';


--
-- Name: COLUMN channel.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.use_status IS '状态(0正常)';


--
-- Name: COLUMN channel.logo_img_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel.logo_img_no IS 'logo图片id';


--
-- Name: channel_servicer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.channel_servicer (
    channel_no uuid NOT NULL,
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    idx bigint DEFAULT 0,
    is_recom smallint DEFAULT 0,
    currency_type character varying(255) DEFAULT ''::character varying NOT NULL,
    use_status smallint DEFAULT 1,
    id character varying(32)
);


ALTER TABLE public.channel_servicer OWNER TO postgres;

--
-- Name: COLUMN channel_servicer.channel_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.channel_no IS '渠道仓库表的id';


--
-- Name: COLUMN channel_servicer.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.create_time IS '创建时间';


--
-- Name: COLUMN channel_servicer.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.is_delete IS '是否删除';


--
-- Name: COLUMN channel_servicer.idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.idx IS '序列';


--
-- Name: COLUMN channel_servicer.is_recom; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.is_recom IS '是否推荐(1-推荐，0-不推荐)';


--
-- Name: COLUMN channel_servicer.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.currency_type IS '币种';


--
-- Name: COLUMN channel_servicer.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.use_status IS '状态(0禁用)';


--
-- Name: COLUMN channel_servicer.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.channel_servicer.id IS '主键（用于删除时定位）';


--
-- Name: collection_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.collection_order (
    log_no character varying(255) NOT NULL,
    from_vaccount_no uuid,
    to_vaccount_no uuid,
    amount bigint DEFAULT 0,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    finish_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    order_status bigint DEFAULT 0,
    balance_type character varying(255) DEFAULT ''::character varying,
    payment_type character varying(255) DEFAULT 2,
    fees bigint DEFAULT 0,
    is_count smallint DEFAULT 0,
    modify_time timestamp(6) without time zone,
    ip character varying(255),
    lat character varying(255),
    lng character varying(255)
);


ALTER TABLE public.collection_order OWNER TO postgres;

--
-- Name: COLUMN collection_order.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.log_no IS '虚拟账户收款日志主键';


--
-- Name: COLUMN collection_order.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.order_status IS '订单状态(0-交易失败，1-交易成功)';


--
-- Name: COLUMN collection_order.balance_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.balance_type IS '金额类型（usd美金，khr瑞尔）';


--
-- Name: COLUMN collection_order.payment_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.payment_type IS '付款方式;1-现金;2-余额';


--
-- Name: COLUMN collection_order.fees; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.fees IS '手续费';


--
-- Name: COLUMN collection_order.is_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.is_count IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';


--
-- Name: COLUMN collection_order.ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.ip IS 'ip';


--
-- Name: COLUMN collection_order.lat; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.lat IS '纬度';


--
-- Name: COLUMN collection_order.lng; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collection_order.lng IS '经度';


--
-- Name: common_help; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.common_help (
    help_no character varying(255) NOT NULL,
    problem character varying(255) DEFAULT ''::character varying,
    answer text DEFAULT ''::character varying,
    idx bigint,
    is_delete smallint DEFAULT 0,
    use_status smallint DEFAULT 1,
    lang character varying(255),
    vs_type smallint,
    modify_time timestamp(6) without time zone,
    create_time timestamp(6) without time zone,
    file_id character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.common_help OWNER TO postgres;

--
-- Name: COLUMN common_help.help_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.help_no IS '主键';


--
-- Name: COLUMN common_help.problem; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.problem IS '问题';


--
-- Name: COLUMN common_help.answer; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.answer IS '答案';


--
-- Name: COLUMN common_help.idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.idx IS '排序';


--
-- Name: COLUMN common_help.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.use_status IS '使用状态0-禁用';


--
-- Name: COLUMN common_help.lang; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.lang IS '语言(km、zh_CN、en)';


--
-- Name: COLUMN common_help.vs_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.vs_type IS '0-app包;1-pos包';


--
-- Name: COLUMN common_help.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.modify_time IS '最后修改时间';


--
-- Name: COLUMN common_help.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.create_time IS '创建时间';


--
-- Name: COLUMN common_help.file_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.common_help.file_id IS '答案的文件id';


--
-- Name: consultation_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.consultation_config (
    id character varying(255) NOT NULL,
    use_status smallint DEFAULT 1,
    is_delete smallint DEFAULT 0,
    create_time timestamp(0) without time zone,
    lang character varying(255) DEFAULT ''::character varying,
    idx bigint DEFAULT 0,
    logo_img_no character varying(32),
    name character varying(255),
    text character varying(2000)
);


ALTER TABLE public.consultation_config OWNER TO postgres;

--
-- Name: COLUMN consultation_config.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.id IS '咨询表';


--
-- Name: COLUMN consultation_config.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.use_status IS '使用状态：0.禁用，1.正常';


--
-- Name: COLUMN consultation_config.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.is_delete IS '1删除';


--
-- Name: COLUMN consultation_config.lang; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.lang IS '语言(km、zh_CN、en)';


--
-- Name: COLUMN consultation_config.idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.idx IS '优先级（值小的再上面）';


--
-- Name: COLUMN consultation_config.logo_img_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.logo_img_no IS '图片id';


--
-- Name: COLUMN consultation_config.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.name IS '主题';


--
-- Name: COLUMN consultation_config.text; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.consultation_config.text IS '内容';


--
-- Name: cust; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cust (
    cust_no uuid NOT NULL,
    account_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    payment_password character varying(255),
    gender smallint DEFAULT 1,
    in_authorization smallint DEFAULT 1,
    out_authorization smallint DEFAULT 1,
    in_transfer_authorization smallint DEFAULT 1,
    out_transfer_authorization smallint DEFAULT 1,
    modify_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    def_pay_no character varying(255) DEFAULT 'usd_balance'::character varying
);


ALTER TABLE public.cust OWNER TO postgres;

--
-- Name: TABLE cust; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.cust IS '主要存账户的静态信息';


--
-- Name: COLUMN cust.cust_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.cust_no IS '用户列表主键';


--
-- Name: COLUMN cust.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.account_no IS '用户账号uuid';


--
-- Name: COLUMN cust.payment_password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.payment_password IS '用户支付密码';


--
-- Name: COLUMN cust.gender; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.gender IS '用户性别（1男，0女）';


--
-- Name: COLUMN cust.in_authorization; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.in_authorization IS '充值权限(0禁用)';


--
-- Name: COLUMN cust.out_authorization; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.out_authorization IS '提现权限(0禁用)';


--
-- Name: COLUMN cust.in_transfer_authorization; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.in_transfer_authorization IS '可转账入权限(0禁用)';


--
-- Name: COLUMN cust.out_transfer_authorization; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.out_transfer_authorization IS '可转账出权限(0禁用)';


--
-- Name: COLUMN cust.def_pay_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.cust.def_pay_no IS '默认卡';


--
-- Name: dict_acc_title; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_acc_title (
    title_no bigint NOT NULL,
    title_name character varying(255),
    parent_title bigint
);


ALTER TABLE public.dict_acc_title OWNER TO postgres;

--
-- Name: TABLE dict_acc_title; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.dict_acc_title IS '会计科目表';


--
-- Name: COLUMN dict_acc_title.title_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_acc_title.title_no IS '会计科目';


--
-- Name: COLUMN dict_acc_title.title_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_acc_title.title_name IS '会计科目名';


--
-- Name: COLUMN dict_acc_title.parent_title; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_acc_title.parent_title IS '父节点';


--
-- Name: dict_account_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_account_type (
    account_type bigint NOT NULL,
    remark character varying(255)
);


ALTER TABLE public.dict_account_type OWNER TO postgres;

--
-- Name: dict_bank_abbr; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_bank_abbr (
    id character varying(255),
    bank_abbr character varying(255),
    bank_name character varying(255)
);


ALTER TABLE public.dict_bank_abbr OWNER TO postgres;

--
-- Name: COLUMN dict_bank_abbr.bank_abbr; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_bank_abbr.bank_abbr IS '缩写';


--
-- Name: COLUMN dict_bank_abbr.bank_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_bank_abbr.bank_name IS '银行名';


--
-- Name: dict_bankname; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_bankname (
    bank_name character varying(255) NOT NULL,
    bank_id character varying(255)
);


ALTER TABLE public.dict_bankname OWNER TO postgres;

--
-- Name: dict_bin_bankname; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_bin_bankname (
    bin_code character varying(255) NOT NULL,
    bank_name character varying(255),
    org_code character varying(255),
    card_name character varying(255),
    card_type character varying(255),
    card_type_no character varying(255)
);


ALTER TABLE public.dict_bin_bankname OWNER TO postgres;

--
-- Name: COLUMN dict_bin_bankname.bank_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_bin_bankname.bank_name IS '银行名';


--
-- Name: COLUMN dict_bin_bankname.card_type_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_bin_bankname.card_type_no IS '0:借记卡

1:贷记卡
2:准贷记卡

3:预付费卡
';


--
-- Name: dict_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_images (
    image_id character varying NOT NULL,
    image_url character varying(255),
    create_time timestamp(6) without time zone,
    status smallint DEFAULT 1,
    modify_time timestamp(6) without time zone,
    account_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    is_delete smallint DEFAULT 0
);


ALTER TABLE public.dict_images OWNER TO postgres;

--
-- Name: COLUMN dict_images.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_images.status IS '1-正常';


--
-- Name: dict_org_abbr; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_org_abbr (
    org_code character varying(255) NOT NULL,
    abbr character varying(255)
);


ALTER TABLE public.dict_org_abbr OWNER TO postgres;

--
-- Name: COLUMN dict_org_abbr.org_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_org_abbr.org_code IS '机构号';


--
-- Name: COLUMN dict_org_abbr.abbr; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_org_abbr.abbr IS '缩写';


--
-- Name: dict_province; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_province (
    province_code bigint NOT NULL,
    province_name character varying(255) DEFAULT ''::character varying,
    short_name character varying(255) DEFAULT ''::character varying,
    full_en_name character varying(255) DEFAULT ''::character varying,
    short_zh_name character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.dict_province OWNER TO postgres;

--
-- Name: TABLE dict_province; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.dict_province IS '省份表';


--
-- Name: COLUMN dict_province.province_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_province.province_code IS '省码';


--
-- Name: COLUMN dict_province.province_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_province.province_name IS '省份名';


--
-- Name: COLUMN dict_province.short_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_province.short_name IS '简写';


--
-- Name: COLUMN dict_province.full_en_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_province.full_en_name IS '全拼';


--
-- Name: COLUMN dict_province.short_zh_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_province.short_zh_name IS '中文简称';


--
-- Name: dict_region; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_region (
    id character varying(64) NOT NULL,
    code character varying(64),
    name character varying(255),
    level character varying(32),
    pid character varying(64),
    longitude numeric(10,4),
    latitude numeric(10,4),
    is_leaf smallint,
    pname character varying(255)
);


ALTER TABLE public.dict_region OWNER TO postgres;

--
-- Name: TABLE dict_region; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.dict_region IS '地区表';


--
-- Name: COLUMN dict_region.is_leaf; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.dict_region.is_leaf IS '1是;0不是';


--
-- Name: dict_region_bank; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_region_bank (
    code character varying(255) NOT NULL,
    name character varying(255),
    province character varying(255),
    city character varying(255),
    bank_type character varying(255),
    province_code character varying(255),
    city_code character varying(255)
);


ALTER TABLE public.dict_region_bank OWNER TO postgres;

--
-- Name: dict_vatype; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dict_vatype (
    va_type bigint NOT NULL,
    remark character varying(255)
);


ALTER TABLE public.dict_vatype OWNER TO postgres;

--
-- Name: exchange_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.exchange_order (
    log_no character varying(255) NOT NULL,
    in_type character varying(255) DEFAULT ''::character varying,
    out_type character varying(255) DEFAULT ''::character varying,
    amount bigint DEFAULT 0,
    create_time timestamp(6) without time zone,
    rate bigint DEFAULT 0,
    order_status bigint DEFAULT 0,
    finish_time timestamp(6) without time zone,
    account_no uuid,
    trans_from character varying(255) DEFAULT ''::character varying,
    trans_amount bigint DEFAULT 0,
    err_reason character varying(255) DEFAULT ''::character varying,
    fees bigint DEFAULT 0,
    is_count smallint,
    modify_time timestamp(6) without time zone,
    ip character varying(255),
    lat character varying(255),
    lng character varying(255)
);


ALTER TABLE public.exchange_order OWNER TO postgres;

--
-- Name: COLUMN exchange_order.in_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.in_type IS '发起金额类型';


--
-- Name: COLUMN exchange_order.out_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.out_type IS '到账金额类型';


--
-- Name: COLUMN exchange_order.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.amount IS '金额';


--
-- Name: COLUMN exchange_order.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.create_time IS '创建时间';


--
-- Name: COLUMN exchange_order.rate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.rate IS '平台汇率';


--
-- Name: COLUMN exchange_order.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.order_status IS '订单状态';


--
-- Name: COLUMN exchange_order.finish_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.finish_time IS '完成时间';


--
-- Name: COLUMN exchange_order.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.account_no IS '发起人';


--
-- Name: COLUMN exchange_order.trans_from; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.trans_from IS 'app,trade';


--
-- Name: COLUMN exchange_order.trans_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.trans_amount IS '转换后金额';


--
-- Name: COLUMN exchange_order.err_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.err_reason IS '失败原因';


--
-- Name: COLUMN exchange_order.fees; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.fees IS '单笔手续费';


--
-- Name: COLUMN exchange_order.is_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.is_count IS '手续费是否统计完毕;0-否;1-是';


--
-- Name: COLUMN exchange_order.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.modify_time IS '修改时间';


--
-- Name: COLUMN exchange_order.ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.ip IS 'ip';


--
-- Name: COLUMN exchange_order.lat; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.lat IS '纬度';


--
-- Name: COLUMN exchange_order.lng; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.exchange_order.lng IS '经度';


--
-- Name: func_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.func_config (
    func_no uuid NOT NULL,
    func_name character varying(255),
    idx bigint,
    use_status integer DEFAULT 1,
    is_delete smallint DEFAULT 0,
    img character varying(512) DEFAULT ''::character varying,
    jump_url character varying(512) DEFAULT ''::character varying,
    application_type character varying(32),
    img_id character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.func_config OWNER TO postgres;

--
-- Name: TABLE func_config; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.func_config IS '钱包功能';


--
-- Name: COLUMN func_config.func_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.func_config.func_no IS '钱包功能入口配置  主键';


--
-- Name: COLUMN func_config.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.func_config.use_status IS '0禁用';


--
-- Name: COLUMN func_config.img; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.func_config.img IS '图片地址';


--
-- Name: COLUMN func_config.jump_url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.func_config.jump_url IS '跳转地址';


--
-- Name: COLUMN func_config.application_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.func_config.application_type IS '0-手机端;1-pos';


--
-- Name: COLUMN func_config.img_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.func_config.img_id IS '图片id';


--
-- Name: gen_code; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.gen_code (
    gen_key character varying(255) NOT NULL,
    account_no uuid,
    amount bigint DEFAULT 0,
    money_type character varying(255) DEFAULT ''::character varying,
    create_time timestamp(6) without time zone,
    code_type character varying(255) DEFAULT ''::character varying,
    use_status integer DEFAULT 1,
    modify_time timestamp(6) without time zone,
    sweep_account_no uuid,
    order_no character varying,
    op_acc_type smallint DEFAULT 0,
    op_acc_no uuid
);


ALTER TABLE public.gen_code OWNER TO postgres;

--
-- Name: COLUMN gen_code.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.gen_code.account_no IS '服务商账号id';


--
-- Name: COLUMN gen_code.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.gen_code.use_status IS '1-初始化;2-已扫码;3-已支付;4-已过期;5-待确认';


--
-- Name: COLUMN gen_code.sweep_account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.gen_code.sweep_account_no IS '谁扫的码,accountid';


--
-- Name: COLUMN gen_code.order_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.gen_code.order_no IS '订单ID';


--
-- Name: COLUMN gen_code.op_acc_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.gen_code.op_acc_type IS '这个码是谁生成的,0-用户1-服务商;2-收银员';


--
-- Name: COLUMN gen_code.op_acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.gen_code.op_acc_no IS '操作员账号,当前账号类型是服务商,就是服务商的id,如果是收银员,那就是收银员的id';


--
-- Name: global_param; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_param (
    param_key character varying(255) NOT NULL,
    param_value text,
    remark character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.global_param OWNER TO postgres;

--
-- Name: TABLE global_param; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.global_param IS '全局变量';


--
-- Name: COLUMN global_param.remark; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.global_param.remark IS '备注';


--
-- Name: headquarters_profit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.headquarters_profit (
    log_no character varying(255) NOT NULL,
    general_ledger_no character varying,
    amount bigint,
    create_time timestamp(6) without time zone,
    order_status smallint DEFAULT 0,
    finish_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    balance_type character varying(255),
    profit_source smallint,
    modify_time timestamp(6) without time zone
);


ALTER TABLE public.headquarters_profit OWNER TO postgres;

--
-- Name: COLUMN headquarters_profit.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.log_no IS '平台盈利统计流水号';


--
-- Name: COLUMN headquarters_profit.general_ledger_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.general_ledger_no IS '订单';


--
-- Name: COLUMN headquarters_profit.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.amount IS '金额';


--
-- Name: COLUMN headquarters_profit.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.create_time IS '创建时间';


--
-- Name: COLUMN headquarters_profit.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.order_status IS '状态';


--
-- Name: COLUMN headquarters_profit.finish_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.finish_time IS '完成时间';


--
-- Name: COLUMN headquarters_profit.balance_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.balance_type IS '金额类型（usd美金，khr瑞尔）';


--
-- Name: COLUMN headquarters_profit.profit_source; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit.profit_source IS '收益来源（1-提现手续费;2-客户转账手续费;3-客户兑换手续费;4-客户收款手续费;5客户存款手续费）';


--
-- Name: headquarters_profit_cashable; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.headquarters_profit_cashable (
    id character varying(255) NOT NULL,
    cashable_balance bigint DEFAULT 0,
    revenue_money bigint DEFAULT 0,
    modify_time timestamp(0) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    currency_type character varying(255)
);


ALTER TABLE public.headquarters_profit_cashable OWNER TO postgres;

--
-- Name: COLUMN headquarters_profit_cashable.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_cashable.id IS '平台利润';


--
-- Name: COLUMN headquarters_profit_cashable.cashable_balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_cashable.cashable_balance IS '可提现余额';


--
-- Name: COLUMN headquarters_profit_cashable.revenue_money; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_cashable.revenue_money IS '总收益的钱';


--
-- Name: COLUMN headquarters_profit_cashable.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_cashable.modify_time IS '最后修改时间';


--
-- Name: COLUMN headquarters_profit_cashable.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_cashable.currency_type IS '币种';


--
-- Name: headquarters_profit_withdraw; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.headquarters_profit_withdraw (
    order_no character varying(32) NOT NULL,
    currency_type character varying(255),
    amount bigint,
    note text,
    create_time timestamp(6) without time zone,
    account_no uuid
);


ALTER TABLE public.headquarters_profit_withdraw OWNER TO postgres;

--
-- Name: COLUMN headquarters_profit_withdraw.order_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_withdraw.order_no IS '平台盈利提现订单表主键';


--
-- Name: COLUMN headquarters_profit_withdraw.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_withdraw.currency_type IS '币种';


--
-- Name: COLUMN headquarters_profit_withdraw.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_withdraw.amount IS '金额';


--
-- Name: COLUMN headquarters_profit_withdraw.note; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_withdraw.note IS '备注';


--
-- Name: COLUMN headquarters_profit_withdraw.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_withdraw.create_time IS '创建时间';


--
-- Name: COLUMN headquarters_profit_withdraw.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.headquarters_profit_withdraw.account_no IS '操作人账号';


--
-- Name: income_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_order (
    log_no character varying(255) NOT NULL,
    act_acc_no uuid,
    amount bigint DEFAULT 0,
    servicer_no uuid,
    create_time timestamp(6) without time zone,
    order_status bigint DEFAULT 0,
    finish_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    query_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    balance_type character varying(255),
    fees bigint DEFAULT 0,
    recv_acc_no uuid,
    recv_vacc uuid,
    op_acc_no uuid,
    settle_hourly_log_no character varying(255) DEFAULT ''::character varying,
    settle_daily_log_no character varying(255) DEFAULT ''::character varying,
    payment_type smallint DEFAULT 1,
    is_count smallint,
    modify_time timestamp(6) without time zone,
    ree_rate bigint,
    real_amount bigint,
    op_acc_type smallint
);


ALTER TABLE public.income_order OWNER TO postgres;

--
-- Name: TABLE income_order; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.income_order IS '收款订单';


--
-- Name: COLUMN income_order.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.log_no IS '存款日志表';


--
-- Name: COLUMN income_order.act_acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.act_acc_no IS '存款人的账号';


--
-- Name: COLUMN income_order.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.amount IS '金额';


--
-- Name: COLUMN income_order.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.servicer_no IS '收款的商户uid';


--
-- Name: COLUMN income_order.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.create_time IS '创建时间';


--
-- Name: COLUMN income_order.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.order_status IS '订单状态（1-初始化;2-等待;3-已支付;4-失败）';


--
-- Name: COLUMN income_order.finish_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.finish_time IS '完成时间';


--
-- Name: COLUMN income_order.query_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.query_time IS '查询时间';


--
-- Name: COLUMN income_order.balance_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.balance_type IS '金额类型（usd美金，khr瑞尔）';


--
-- Name: COLUMN income_order.fees; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.fees IS '手续费';


--
-- Name: COLUMN income_order.recv_acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.recv_acc_no IS '收款人账号';


--
-- Name: COLUMN income_order.recv_vacc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.recv_vacc IS '收款人虚拟账号';


--
-- Name: COLUMN income_order.op_acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.op_acc_no IS '操作员账号,当前账号类型是服务商,就是服务商的id,如果是收银员,那就是收银员的id';


--
-- Name: COLUMN income_order.settle_hourly_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.settle_hourly_log_no IS '小时对账流水';


--
-- Name: COLUMN income_order.settle_daily_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.settle_daily_log_no IS '天对账流水';


--
-- Name: COLUMN income_order.payment_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.payment_type IS '付款方式.1-现金;2-余额';


--
-- Name: COLUMN income_order.is_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.is_count IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';


--
-- Name: COLUMN income_order.ree_rate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.ree_rate IS '手续费率';


--
-- Name: COLUMN income_order.real_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.real_amount IS '实际到账金额';


--
-- Name: COLUMN income_order.op_acc_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_order.op_acc_type IS '这笔订单是谁产生的.1-服务商;2-店员';


--
-- Name: income_ougo_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_ougo_config (
    income_ougo_config_no uuid NOT NULL,
    currency_type character varying(255),
    name character varying(255),
    use_status character varying(255) DEFAULT 0,
    idx bigint,
    config_type smallint,
    is_delete smallint DEFAULT 0
);


ALTER TABLE public.income_ougo_config OWNER TO postgres;

--
-- Name: COLUMN income_ougo_config.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_ougo_config.currency_type IS '币种';


--
-- Name: COLUMN income_ougo_config.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_ougo_config.name IS '名称';


--
-- Name: COLUMN income_ougo_config.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_ougo_config.use_status IS '状态(0正常)';


--
-- Name: COLUMN income_ougo_config.idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_ougo_config.idx IS '排序序号';


--
-- Name: COLUMN income_ougo_config.config_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.income_ougo_config.config_type IS '类型（1.充值方式。2.提现方式）';


--
-- Name: income_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_type (
    income_type character varying NOT NULL,
    income_name character varying(255),
    use_status smallint DEFAULT 0,
    idx bigint
);


ALTER TABLE public.income_type OWNER TO postgres;

--
-- Name: lang; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.lang (
    key character varying(512) NOT NULL,
    type character varying(255) DEFAULT 1 NOT NULL,
    is_delete smallint DEFAULT 0,
    lang_km character varying(255) DEFAULT ''::character varying,
    lang_en character varying(255) DEFAULT ''::character varying,
    lang_ch character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.lang OWNER TO postgres;

--
-- Name: COLUMN lang.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.lang.type IS '类型（1-文字，2-图片, 3-错误提示）';


--
-- Name: COLUMN lang.lang_km; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.lang.lang_km IS '柬埔寨语';


--
-- Name: COLUMN lang.lang_ch; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.lang.lang_ch IS '中文';


--
-- Name: log_account; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_account (
    log_no character varying NOT NULL,
    description character varying(255),
    account_uid character varying,
    log_time timestamp(6) without time zone,
    type integer
);


ALTER TABLE public.log_account OWNER TO postgres;

--
-- Name: TABLE log_account; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.log_account IS '账户操作日志';


--
-- Name: COLUMN log_account.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_account.description IS '描述';


--
-- Name: log_account_web; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_account_web (
    log_no character varying NOT NULL,
    description text,
    account_uid character varying,
    create_time timestamp(6) without time zone,
    type character varying(255)
);


ALTER TABLE public.log_account_web OWNER TO postgres;

--
-- Name: TABLE log_account_web; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.log_account_web IS '账户操作日志';


--
-- Name: COLUMN log_account_web.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_account_web.log_no IS 'id（此表专门记WEB也就是后台管理系统的日志）';


--
-- Name: COLUMN log_account_web.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_account_web.description IS '描述';


--
-- Name: COLUMN log_account_web.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_account_web.type IS '1-账号；
';


--
-- Name: log_app_messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_app_messages (
    log_no character varying(255),
    order_no character varying(32),
    order_type smallint,
    is_read smallint DEFAULT 0,
    is_push smallint DEFAULT 0,
    account_no uuid,
    create_time timestamp(0) without time zone
);


ALTER TABLE public.log_app_messages OWNER TO postgres;

--
-- Name: TABLE log_app_messages; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.log_app_messages IS '用户消息中心';


--
-- Name: COLUMN log_app_messages.order_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_app_messages.order_no IS '订单号';


--
-- Name: COLUMN log_app_messages.order_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_app_messages.order_type IS '订单类型（1-兑换,2-充值,3-提现,4-转账,5-收款）';


--
-- Name: COLUMN log_app_messages.is_read; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_app_messages.is_read IS '是否已读';


--
-- Name: COLUMN log_app_messages.is_push; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_app_messages.is_push IS '是否推送';


--
-- Name: COLUMN log_app_messages.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_app_messages.account_no IS '要推送的账号';


--
-- Name: log_card; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_card (
    log_no character varying NOT NULL,
    card_num character varying DEFAULT 0,
    name character varying(255) DEFAULT 0,
    account_no uuid,
    va_type smallint DEFAULT 0,
    channel_no uuid,
    channel_type smallint,
    create_time timestamp(6) without time zone,
    descript character varying(255)
);


ALTER TABLE public.log_card OWNER TO postgres;

--
-- Name: COLUMN log_card.card_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.card_num IS '卡号';


--
-- Name: COLUMN log_card.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.name IS '持卡人名字';


--
-- Name: COLUMN log_card.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.account_no IS '谁创建的';


--
-- Name: COLUMN log_card.va_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.va_type IS '币种';


--
-- Name: COLUMN log_card.channel_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.channel_no IS '渠道id';


--
-- Name: COLUMN log_card.channel_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.channel_type IS '1-用户,总部;2-pos';


--
-- Name: COLUMN log_card.descript; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_card.descript IS '描述';


--
-- Name: log_exchange_rate; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_exchange_rate (
    log_time timestamp(6) without time zone,
    usd_khr bigint DEFAULT 0,
    khr_usd bigint DEFAULT 0
);


ALTER TABLE public.log_exchange_rate OWNER TO postgres;

--
-- Name: TABLE log_exchange_rate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.log_exchange_rate IS '汇率变更表';


--
-- Name: log_login; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_login (
    log_time timestamp(6) without time zone,
    acc_no uuid,
    ip character varying(255) DEFAULT ''::character varying,
    result integer DEFAULT 0,
    client character varying(255) DEFAULT ''::character varying,
    log_no character varying(255) NOT NULL
);


ALTER TABLE public.log_login OWNER TO postgres;

--
-- Name: log_to_headquarters; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_to_headquarters (
    log_no character varying NOT NULL,
    servicer_no uuid,
    currency_type character varying(255),
    amount bigint,
    order_status smallint DEFAULT 0,
    collection_type smallint,
    card_no uuid,
    create_time timestamp(0) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    finish_time timestamp(0) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    order_type smallint,
    image_id character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.log_to_headquarters OWNER TO postgres;

--
-- Name: COLUMN log_to_headquarters.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.log_no IS '转账至总部';


--
-- Name: COLUMN log_to_headquarters.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.servicer_no IS '服务商id';


--
-- Name: COLUMN log_to_headquarters.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.currency_type IS '币种';


--
-- Name: COLUMN log_to_headquarters.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.amount IS '金额';


--
-- Name: COLUMN log_to_headquarters.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.order_status IS '订单状态(0待审核、1已完成，2已关闭)';


--
-- Name: COLUMN log_to_headquarters.collection_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.collection_type IS '收款方式,1-支票;2-现金;3-银行转账;4-其他';


--
-- Name: COLUMN log_to_headquarters.card_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.card_no IS '卡uid';


--
-- Name: COLUMN log_to_headquarters.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.create_time IS '发起时间';


--
-- Name: COLUMN log_to_headquarters.finish_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.finish_time IS '完成时间';


--
-- Name: COLUMN log_to_headquarters.order_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_headquarters.order_type IS '订单类型（1-交易转账、2-结算转账）';


--
-- Name: log_to_servicer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_to_servicer (
    log_no character varying(255) NOT NULL,
    currency_type character varying(255),
    servicer_no uuid,
    collection_type smallint,
    card_no uuid,
    amount character varying(255),
    create_time timestamp(0) without time zone,
    order_type smallint,
    order_status smallint,
    finish_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    motify_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone
);


ALTER TABLE public.log_to_servicer OWNER TO postgres;

--
-- Name: COLUMN log_to_servicer.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.log_no IS '转账至服务商流水';


--
-- Name: COLUMN log_to_servicer.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.currency_type IS '币种';


--
-- Name: COLUMN log_to_servicer.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.servicer_no IS '服务商id';


--
-- Name: COLUMN log_to_servicer.collection_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.collection_type IS '收款方式,1-支票;2-现金;3-银行转账;4-其他';


--
-- Name: COLUMN log_to_servicer.card_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.card_no IS '卡主键uid';


--
-- Name: COLUMN log_to_servicer.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.amount IS '金额';


--
-- Name: COLUMN log_to_servicer.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.create_time IS '创建时间';


--
-- Name: COLUMN log_to_servicer.order_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.order_type IS '订单类型（1-交易转账、2-结算转账）';


--
-- Name: COLUMN log_to_servicer.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.order_status IS '订单状态(0-待审核、1-已完成，2-已关闭)';


--
-- Name: COLUMN log_to_servicer.finish_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_to_servicer.finish_time IS '完成时间';


--
-- Name: log_vaccount; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.log_vaccount (
    log_no character varying(255) NOT NULL,
    vaccount_no uuid NOT NULL,
    create_time timestamp(6) without time zone,
    amount bigint DEFAULT 0,
    op_type smallint DEFAULT 0,
    frozen_balance bigint DEFAULT 0,
    balance bigint DEFAULT 0,
    reason bigint DEFAULT 0,
    settle_hourly_log_no character varying(255) DEFAULT ''::character varying,
    settle_daily_log_no character varying(255) DEFAULT ''::character varying,
    biz_log_no character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.log_vaccount OWNER TO postgres;

--
-- Name: TABLE log_vaccount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.log_vaccount IS '虚拟账户日志表';


--
-- Name: COLUMN log_vaccount.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.log_no IS '虚拟账户日志表';


--
-- Name: COLUMN log_vaccount.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.amount IS '变化的金额';


--
-- Name: COLUMN log_vaccount.op_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.op_type IS '1:+;2:-;3:冻结;4:解冻;';


--
-- Name: COLUMN log_vaccount.frozen_balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.frozen_balance IS '冻结余额';


--
-- Name: COLUMN log_vaccount.balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.balance IS '现有余额';


--
-- Name: COLUMN log_vaccount.reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.reason IS '原因;1-兑换,2-充值,3-提现,4-转账,5-收款;6-手续费;7-pos 端取消提现';


--
-- Name: COLUMN log_vaccount.settle_hourly_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.settle_hourly_log_no IS '每小时对账流水';


--
-- Name: COLUMN log_vaccount.settle_daily_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.settle_daily_log_no IS '每日对账流水';


--
-- Name: COLUMN log_vaccount.biz_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.log_vaccount.biz_log_no IS '业务流水';


--
-- Name: login_token; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.login_token (
    acc_no uuid NOT NULL,
    routes text,
    token text,
    login_time timestamp(6) without time zone,
    ip character varying(255) DEFAULT ''::character varying,
    last_op_time timestamp(6) without time zone,
    imei character varying(255) DEFAULT ''::character varying
);


ALTER TABLE public.login_token OWNER TO postgres;

--
-- Name: COLUMN login_token.acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.acc_no IS '登录账号';


--
-- Name: COLUMN login_token.routes; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.routes IS '路径列表';


--
-- Name: COLUMN login_token.token; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.token IS '令牌';


--
-- Name: COLUMN login_token.login_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.login_time IS '登录时间';


--
-- Name: COLUMN login_token.ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.ip IS '登录ip';


--
-- Name: COLUMN login_token.last_op_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.last_op_time IS '最后操作时间';


--
-- Name: COLUMN login_token.imei; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.login_token.imei IS 'imei';


--
-- Name: outgo_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.outgo_order (
    log_no character varying NOT NULL,
    vaccount_no uuid,
    amount bigint,
    create_time timestamp(6) without time zone,
    order_status bigint,
    modify_time timestamp(6) without time zone,
    balance_type character varying(255),
    fees bigint,
    servicer_no uuid,
    op_acc_no uuid,
    settle_hourly_log_no character varying(255) DEFAULT ''::character varying,
    settle_daily_log_no character varying(255) DEFAULT ''::character varying,
    rate character varying(255),
    payment_type smallint DEFAULT 2,
    is_count smallint,
    withdraw_type smallint,
    cancel_reason character varying(255),
    risk_no character varying,
    real_amount bigint,
    op_acc_type smallint,
    ip character varying(255),
    lat character varying(255),
    lng character varying(255)
);


ALTER TABLE public.outgo_order OWNER TO postgres;

--
-- Name: COLUMN outgo_order.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.log_no IS '取款';


--
-- Name: COLUMN outgo_order.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.order_status IS '订单状态（1-初始化;2-等待;3-已支付;4-失败(超时);5-待确认;6-取消）';


--
-- Name: COLUMN outgo_order.fees; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.fees IS '手续费';


--
-- Name: COLUMN outgo_order.op_acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.op_acc_no IS '操作员账号,当前账号类型是服务商,就是服务商的id,如果是收银员,那就是收银员的id';


--
-- Name: COLUMN outgo_order.settle_hourly_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.settle_hourly_log_no IS '小时对账流水';


--
-- Name: COLUMN outgo_order.settle_daily_log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.settle_daily_log_no IS '天对账流水';


--
-- Name: COLUMN outgo_order.rate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.rate IS '费率';


--
-- Name: COLUMN outgo_order.payment_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.payment_type IS '付款方式;1-现金,2-余额';


--
-- Name: COLUMN outgo_order.is_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.is_count IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';


--
-- Name: COLUMN outgo_order.withdraw_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.withdraw_type IS '0-手机号提现;1-扫码提现;2-全部提现';


--
-- Name: COLUMN outgo_order.cancel_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.cancel_reason IS '取消原因';


--
-- Name: COLUMN outgo_order.risk_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.risk_no IS '风控结果主键id,pos机查风控数据使用';


--
-- Name: COLUMN outgo_order.real_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.real_amount IS '实际到账金额';


--
-- Name: COLUMN outgo_order.op_acc_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.op_acc_type IS '这笔订单是谁产生的,1-服务商;2-收银员';


--
-- Name: COLUMN outgo_order.ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.ip IS 'ip';


--
-- Name: COLUMN outgo_order.lat; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.lat IS '纬度';


--
-- Name: COLUMN outgo_order.lng; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.outgo_order.lng IS '经度';


--
-- Name: outgo_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.outgo_type (
    outgo_type character varying NOT NULL,
    outgo_name character varying(255),
    use_status smallint DEFAULT 0,
    idx bigint
);


ALTER TABLE public.outgo_type OWNER TO postgres;

--
-- Name: TABLE outgo_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.outgo_type IS '提现方式';


--
-- Name: platform_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.platform_config (
    account_uid uuid NOT NULL,
    top_menu_status bigint DEFAULT 1,
    side_menu_status bigint DEFAULT 1
);


ALTER TABLE public.platform_config OWNER TO postgres;

--
-- Name: COLUMN platform_config.account_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.platform_config.account_uid IS '账户uid';


--
-- Name: COLUMN platform_config.top_menu_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.platform_config.top_menu_status IS '顶部菜单配置1、打开、0关闭';


--
-- Name: COLUMN platform_config.side_menu_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.platform_config.side_menu_status IS '侧边菜单配置1、打开、0关闭';


--
-- Name: rela_acc_iden; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rela_acc_iden (
    account_no uuid NOT NULL,
    account_type bigint NOT NULL,
    iden_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid
);


ALTER TABLE public.rela_acc_iden OWNER TO postgres;

--
-- Name: COLUMN rela_acc_iden.account_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_acc_iden.account_type IS '1: 管理员2: 运营3: 服务商4 :用户5:操作员';


--
-- Name: rela_account_role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rela_account_role (
    rela_uid uuid NOT NULL,
    account_uid uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    role_uid uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid
);


ALTER TABLE public.rela_account_role OWNER TO postgres;

--
-- Name: TABLE rela_account_role; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.rela_account_role IS '账户角色关联表';


--
-- Name: COLUMN rela_account_role.account_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_account_role.account_uid IS '账户uid';


--
-- Name: COLUMN rela_account_role.role_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_account_role.role_uid IS '角色uid';


--
-- Name: rela_imei_pubkey; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rela_imei_pubkey (
    rela_no uuid NOT NULL,
    imei character varying(255),
    pub_key text,
    create_time timestamp(6) without time zone
);


ALTER TABLE public.rela_imei_pubkey OWNER TO postgres;

--
-- Name: COLUMN rela_imei_pubkey.imei; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_imei_pubkey.imei IS '设备号';


--
-- Name: COLUMN rela_imei_pubkey.pub_key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_imei_pubkey.pub_key IS '客户端公钥';


--
-- Name: rela_role_url; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rela_role_url (
    rela_uid uuid NOT NULL,
    url_uid uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    role_uid uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid
);


ALTER TABLE public.rela_role_url OWNER TO postgres;

--
-- Name: TABLE rela_role_url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.rela_role_url IS '角色-权限关系';


--
-- Name: COLUMN rela_role_url.url_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_role_url.url_uid IS 'url_uid';


--
-- Name: COLUMN rela_role_url.role_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rela_role_url.role_uid IS '角色uid';


--
-- Name: risk_result; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.risk_result (
    risk_no character varying NOT NULL,
    risk_result smallint,
    risk_threshold character varying(255),
    create_time timestamp(6) without time zone,
    api_type character varying(255),
    payer_acc_no character varying(255),
    action_time character varying(255),
    eva_execute_type character varying(255),
    eva_score character varying(255),
    money_type character varying(255),
    order_no character varying(255),
    score bigint,
    product_type character varying(255)
);


ALTER TABLE public.risk_result OWNER TO postgres;

--
-- Name: COLUMN risk_result.risk_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.risk_no IS '风控结果';


--
-- Name: COLUMN risk_result.risk_result; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.risk_result IS '风控结果是否风控,0-否;1-是';


--
-- Name: COLUMN risk_result.risk_threshold; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.risk_threshold IS '风控阈值';


--
-- Name: COLUMN risk_result.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.create_time IS '生成时间';


--
-- Name: COLUMN risk_result.api_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.api_type IS '事件';


--
-- Name: COLUMN risk_result.payer_acc_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.payer_acc_no IS '操作账号';


--
-- Name: COLUMN risk_result.action_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.action_time IS '操作时间';


--
-- Name: COLUMN risk_result.eva_execute_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.eva_execute_type IS '风控执行类型';


--
-- Name: COLUMN risk_result.eva_score; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.eva_score IS '风控估算的得分';


--
-- Name: COLUMN risk_result.money_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.money_type IS '币种';


--
-- Name: COLUMN risk_result.order_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.order_no IS '订单号';


--
-- Name: COLUMN risk_result.score; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.score IS '实际得分';


--
-- Name: COLUMN risk_result.product_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.risk_result.product_type IS '产品类型,如转账:trancefer';


--
-- Name: role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role (
    role_no uuid NOT NULL,
    role_name character varying(255) DEFAULT ''::character varying,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    modify_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    acc_type character varying(255),
    def_type smallint DEFAULT 0,
    master_acc uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    is_delete smallint DEFAULT 0
);


ALTER TABLE public.role OWNER TO postgres;

--
-- Name: TABLE role; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.role IS '角色表';


--
-- Name: COLUMN role.role_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.role_no IS '角色uid';


--
-- Name: COLUMN role.role_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.role_name IS '角色名';


--
-- Name: COLUMN role.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.create_time IS '创建时间';


--
-- Name: COLUMN role.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.modify_time IS '修改时间';


--
-- Name: COLUMN role.acc_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.acc_type IS '账号类型';


--
-- Name: COLUMN role.def_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.def_type IS '是否设为默认页面';


--
-- Name: COLUMN role.master_acc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.role.master_acc IS '主账号';


--
-- Name: servicer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servicer (
    servicer_no uuid NOT NULL,
    account_no uuid,
    addr character varying(512),
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    use_status integer DEFAULT 1,
    commission_sharing bigint DEFAULT 0,
    income_authorization smallint DEFAULT 1,
    outgo_authorization smallint DEFAULT 1,
    open_idx bigint DEFAULT 0,
    contact_person character varying(255) DEFAULT ''::character varying,
    contact_phone character varying(255) DEFAULT ''::character varying,
    contact_addr character varying(1024) DEFAULT ''::character varying,
    lat character varying(255) DEFAULT 0,
    lng character varying(255) DEFAULT 0,
    password character varying(255),
    modify_time timestamp(6) without time zone,
    income_sharing bigint DEFAULT 0,
    scope character varying(64) DEFAULT 0,
    scope_off smallint DEFAULT 0
);


ALTER TABLE public.servicer OWNER TO postgres;

--
-- Name: COLUMN servicer.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.servicer_no IS '服务商表主键';


--
-- Name: COLUMN servicer.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.account_no IS '账号uid';


--
-- Name: COLUMN servicer.addr; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.addr IS '开店地址';


--
-- Name: COLUMN servicer.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.create_time IS '开通时间';


--
-- Name: COLUMN servicer.is_delete; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.is_delete IS '是否删除';


--
-- Name: COLUMN servicer.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.use_status IS '服务商状态0禁用1启用';


--
-- Name: COLUMN servicer.commission_sharing; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.commission_sharing IS '取款手续费分成';


--
-- Name: COLUMN servicer.income_authorization; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.income_authorization IS '收款权限,0-禁用;1-开通';


--
-- Name: COLUMN servicer.outgo_authorization; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.outgo_authorization IS '取款权限,0-禁用;1-开通';


--
-- Name: COLUMN servicer.open_idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.open_idx IS '开户序列';


--
-- Name: COLUMN servicer.contact_person; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.contact_person IS '联系人';


--
-- Name: COLUMN servicer.contact_phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.contact_phone IS '联系电话';


--
-- Name: COLUMN servicer.contact_addr; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.contact_addr IS '联系人地址';


--
-- Name: COLUMN servicer.lat; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.lat IS '纬度';


--
-- Name: COLUMN servicer.lng; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.lng IS '经度';


--
-- Name: COLUMN servicer.password; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.password IS '服务商密码';


--
-- Name: COLUMN servicer.income_sharing; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.income_sharing IS '存款手续费分成';


--
-- Name: COLUMN servicer.scope; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.scope IS '营业范围(公里)';


--
-- Name: COLUMN servicer.scope_off; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer.scope_off IS '围栏开关（0关闭1开启）';


--
-- Name: servicer_count; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servicer_count (
    servicer_no uuid NOT NULL,
    currency_type character varying(255),
    in_num bigint DEFAULT 0 NOT NULL,
    in_amount bigint DEFAULT 0 NOT NULL,
    out_num bigint DEFAULT 0 NOT NULL,
    out_amount bigint DEFAULT 0 NOT NULL,
    profit_num bigint DEFAULT 0 NOT NULL,
    modify_time timestamp(0) without time zone,
    profit_amount bigint DEFAULT 0 NOT NULL,
    recharge_num bigint DEFAULT 0 NOT NULL,
    recharge_amount bigint DEFAULT 0 NOT NULL,
    withdraw_num bigint DEFAULT 0 NOT NULL,
    withdraw_amount bigint DEFAULT 0 NOT NULL
);


ALTER TABLE public.servicer_count OWNER TO postgres;

--
-- Name: TABLE servicer_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.servicer_count IS '服务商累计统计';


--
-- Name: COLUMN servicer_count.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.servicer_no IS '运营商uuid';


--
-- Name: COLUMN servicer_count.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.currency_type IS '币种(usd,khr)';


--
-- Name: COLUMN servicer_count.in_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.in_num IS '存款-总数量';


--
-- Name: COLUMN servicer_count.in_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.in_amount IS '存款-总金额';


--
-- Name: COLUMN servicer_count.out_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.out_num IS '取款-总数量';


--
-- Name: COLUMN servicer_count.out_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.out_amount IS '取款-总金额';


--
-- Name: COLUMN servicer_count.profit_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.profit_num IS '收益-总数量';


--
-- Name: COLUMN servicer_count.modify_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.modify_time IS '最后修改时间';


--
-- Name: COLUMN servicer_count.profit_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.profit_amount IS '收益-总金额';


--
-- Name: COLUMN servicer_count.recharge_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.recharge_num IS '充值-总数量';


--
-- Name: COLUMN servicer_count.recharge_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.recharge_amount IS '充值-总金额';


--
-- Name: COLUMN servicer_count.withdraw_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.withdraw_num IS '提现-总数量';


--
-- Name: COLUMN servicer_count.withdraw_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count.withdraw_amount IS '提现-总金额';


--
-- Name: servicer_count_list; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servicer_count_list (
    servicer_no uuid,
    currency_type character varying(255),
    create_time timestamp(0) without time zone,
    in_num integer DEFAULT 0 NOT NULL,
    in_amount bigint DEFAULT 0 NOT NULL,
    out_num integer DEFAULT 0 NOT NULL,
    out_amount bigint DEFAULT 0 NOT NULL,
    profit_num integer DEFAULT 0 NOT NULL,
    profit_amount bigint DEFAULT 0 NOT NULL,
    recharge_num integer DEFAULT 0 NOT NULL,
    recharge_amount bigint DEFAULT 0 NOT NULL,
    withdraw_num integer DEFAULT 0 NOT NULL,
    withdraw_amount bigint DEFAULT 0 NOT NULL,
    id bigint NOT NULL,
    dates date,
    is_counted smallint DEFAULT 0
);


ALTER TABLE public.servicer_count_list OWNER TO postgres;

--
-- Name: TABLE servicer_count_list; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.servicer_count_list IS '服务商对账单表';


--
-- Name: COLUMN servicer_count_list.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.servicer_no IS '服务商uuid';


--
-- Name: COLUMN servicer_count_list.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.currency_type IS '币种(usd,khr)';


--
-- Name: COLUMN servicer_count_list.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.create_time IS '生成时间';


--
-- Name: COLUMN servicer_count_list.in_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.in_num IS '存款-数量';


--
-- Name: COLUMN servicer_count_list.in_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.in_amount IS '存款-金额';


--
-- Name: COLUMN servicer_count_list.out_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.out_num IS '取款-数量';


--
-- Name: COLUMN servicer_count_list.out_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.out_amount IS '取款-金额';


--
-- Name: COLUMN servicer_count_list.profit_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.profit_num IS '收益-数量';


--
-- Name: COLUMN servicer_count_list.profit_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.profit_amount IS '收益-金额';


--
-- Name: COLUMN servicer_count_list.recharge_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.recharge_num IS '充值-数量';


--
-- Name: COLUMN servicer_count_list.recharge_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.recharge_amount IS '充值-金额';


--
-- Name: COLUMN servicer_count_list.withdraw_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.withdraw_num IS '提现-数量';


--
-- Name: COLUMN servicer_count_list.withdraw_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.withdraw_amount IS '提现-金额';


--
-- Name: COLUMN servicer_count_list.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.id IS '自增主键';


--
-- Name: COLUMN servicer_count_list.dates; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.dates IS '日期';


--
-- Name: COLUMN servicer_count_list.is_counted; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_count_list.is_counted IS '是否已经统计完成(1是，0否)';


--
-- Name: servicer_count_list_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.servicer_count_list_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.servicer_count_list_id_seq OWNER TO postgres;

--
-- Name: servicer_count_list_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.servicer_count_list_id_seq OWNED BY public.servicer_count_list.id;


--
-- Name: servicer_img; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servicer_img (
    servicer_img_no character varying(255) NOT NULL,
    img_id character varying(400) DEFAULT ''::text,
    img_type smallint DEFAULT 0,
    create_time timestamp(6) without time zone NOT NULL,
    servicer_no uuid,
    is_delete smallint DEFAULT 0
);


ALTER TABLE public.servicer_img OWNER TO postgres;

--
-- Name: COLUMN servicer_img.servicer_img_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_img.servicer_img_no IS '商户照片主键';


--
-- Name: COLUMN servicer_img.img_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_img.img_id IS '照片id';


--
-- Name: COLUMN servicer_img.img_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_img.img_type IS '照片类型（1.营业执照。2.3.4 经营场所照片）';


--
-- Name: COLUMN servicer_img.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_img.create_time IS '创建时间';


--
-- Name: COLUMN servicer_img.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_img.servicer_no IS '商户uuid';


--
-- Name: servicer_profit_ledger; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servicer_profit_ledger (
    log_no character varying NOT NULL,
    amount_order bigint,
    servicefee_amount_sum bigint,
    split_proportion bigint,
    actual_income bigint,
    payment_time timestamp(0) without time zone,
    servicer_no uuid NOT NULL,
    currency_type character varying(255),
    order_type smallint
);


ALTER TABLE public.servicer_profit_ledger OWNER TO postgres;

--
-- Name: COLUMN servicer_profit_ledger.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.log_no IS '服务商分成利润表,订单号';


--
-- Name: COLUMN servicer_profit_ledger.amount_order; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.amount_order IS '订单金额';


--
-- Name: COLUMN servicer_profit_ledger.servicefee_amount_sum; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.servicefee_amount_sum IS '手续费总额';


--
-- Name: COLUMN servicer_profit_ledger.split_proportion; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.split_proportion IS '分成比例';


--
-- Name: COLUMN servicer_profit_ledger.actual_income; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.actual_income IS '实际所得';


--
-- Name: COLUMN servicer_profit_ledger.payment_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.payment_time IS '到账时间';


--
-- Name: COLUMN servicer_profit_ledger.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.servicer_no IS '服务商uuid';


--
-- Name: COLUMN servicer_profit_ledger.currency_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.currency_type IS '币种';


--
-- Name: COLUMN servicer_profit_ledger.order_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_profit_ledger.order_type IS '订单类型,1-存款,2-手机号取款;3-扫码取款';


--
-- Name: servicer_terminal; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.servicer_terminal (
    terminal_no character varying NOT NULL,
    servicer_no uuid,
    terminal_number character varying(255),
    pos_sn character varying(255) DEFAULT ''::character varying,
    is_delete smallint DEFAULT 0,
    use_status smallint DEFAULT 1
);


ALTER TABLE public.servicer_terminal OWNER TO postgres;

--
-- Name: COLUMN servicer_terminal.terminal_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_terminal.terminal_no IS '收款设备主键';


--
-- Name: COLUMN servicer_terminal.servicer_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_terminal.servicer_no IS '运营商uuid';


--
-- Name: COLUMN servicer_terminal.terminal_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_terminal.terminal_number IS '收款设备号（终端编号）';


--
-- Name: COLUMN servicer_terminal.pos_sn; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_terminal.pos_sn IS 'pos机码';


--
-- Name: COLUMN servicer_terminal.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.servicer_terminal.use_status IS '0禁用1启用';


--
-- Name: settle_servicer_hourly; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.settle_servicer_hourly (
    log_no character varying(255) NOT NULL,
    start_time timestamp(6) without time zone,
    finish_time timestamp(6) without time zone,
    run_status bigint DEFAULT 0,
    begin_time timestamp(6) without time zone,
    end_time timestamp(6) without time zone,
    sum_income_usd bigint DEFAULT 0,
    sum_outgo_usd bigint DEFAULT 0,
    balance_usd bigint DEFAULT 0,
    delta_amount_usd bigint DEFAULT 0,
    balance_khr bigint DEFAULT 0,
    delta_amount_khr bigint DEFAULT 0,
    sum_income_khr bigint DEFAULT 0,
    sum_outgo_khr bigint DEFAULT 0,
    fbalance_usd bigint DEFAULT 0,
    fbalance_khr bigint DEFAULT 0
);


ALTER TABLE public.settle_servicer_hourly OWNER TO postgres;

--
-- Name: TABLE settle_servicer_hourly; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.settle_servicer_hourly IS '服务商每小时对账';


--
-- Name: COLUMN settle_servicer_hourly.begin_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.begin_time IS '账期开始时间';


--
-- Name: COLUMN settle_servicer_hourly.end_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.end_time IS '账期结束时间';


--
-- Name: COLUMN settle_servicer_hourly.sum_income_usd; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.sum_income_usd IS '统计的入款';


--
-- Name: COLUMN settle_servicer_hourly.sum_outgo_usd; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.sum_outgo_usd IS '统计的出款';


--
-- Name: COLUMN settle_servicer_hourly.balance_usd; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.balance_usd IS '已核实金额';


--
-- Name: COLUMN settle_servicer_hourly.delta_amount_usd; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.delta_amount_usd IS '未核实金额';


--
-- Name: COLUMN settle_servicer_hourly.balance_khr; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.balance_khr IS '已核实金额';


--
-- Name: COLUMN settle_servicer_hourly.delta_amount_khr; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.settle_servicer_hourly.delta_amount_khr IS '未核实金额';


--
-- Name: settle_vaccount_balance_hourly; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.settle_vaccount_balance_hourly (
    log_no character varying(255) NOT NULL,
    vaccount_no uuid NOT NULL,
    balance bigint DEFAULT 0,
    frozen_balance bigint DEFAULT 0,
    create_time timestamp(6) without time zone
);


ALTER TABLE public.settle_vaccount_balance_hourly OWNER TO postgres;

--
-- Name: TABLE settle_vaccount_balance_hourly; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.settle_vaccount_balance_hourly IS '每小时虚帐余额';


--
-- Name: sms_send_record; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sms_send_record (
    id character varying NOT NULL,
    msgid character varying,
    account character varying(255),
    business character varying(64),
    mobile character varying(255),
    msg character varying(255),
    status integer,
    created_at timestamp(6) without time zone
);


ALTER TABLE public.sms_send_record OWNER TO postgres;

--
-- Name: COLUMN sms_send_record.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sms_send_record.status IS '0-成功,1-失败';


--
-- Name: transfer_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transfer_order (
    log_no character varying(255) NOT NULL,
    from_vaccount_no uuid,
    to_vaccount_no uuid,
    amount bigint DEFAULT 0,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    finish_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    order_status bigint DEFAULT 0,
    balance_type character varying(255) DEFAULT ''::character varying,
    exchange_type bigint DEFAULT 0,
    fees bigint DEFAULT 0,
    payment_type smallint DEFAULT 2,
    is_count smallint,
    modify_time timestamp(6) without time zone,
    ree_rate bigint,
    real_amount bigint,
    ip character varying(255),
    lat character varying(255),
    lng character varying(255)
);


ALTER TABLE public.transfer_order OWNER TO postgres;

--
-- Name: COLUMN transfer_order.log_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.log_no IS '虚拟账户转账日志主键';


--
-- Name: COLUMN transfer_order.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.order_status IS '订单状态';


--
-- Name: COLUMN transfer_order.balance_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.balance_type IS '金额类型（usd美金，khr瑞尔）';


--
-- Name: COLUMN transfer_order.exchange_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.exchange_type IS '兑换类型 0-扫码1-支付';


--
-- Name: COLUMN transfer_order.fees; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.fees IS '手续费';


--
-- Name: COLUMN transfer_order.payment_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.payment_type IS '付款方式;1-现金,2-余额';


--
-- Name: COLUMN transfer_order.is_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.is_count IS '手续费是否统计完毕;0-否;1-已统计到临时账户但财务未确认;2-财务已确认';


--
-- Name: COLUMN transfer_order.ree_rate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.ree_rate IS '手续费率';


--
-- Name: COLUMN transfer_order.real_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.real_amount IS '实际到账金额';


--
-- Name: COLUMN transfer_order.ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.ip IS 'ip';


--
-- Name: COLUMN transfer_order.lat; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.lat IS '纬度';


--
-- Name: COLUMN transfer_order.lng; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transfer_order.lng IS '经度';


--
-- Name: url; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.url (
    url_uid uuid NOT NULL,
    url_name character varying(255),
    url character varying(255),
    parent_uid uuid,
    title character varying(255),
    icon character varying(255),
    component_name character varying(255) DEFAULT ''::character varying,
    component_path character varying(255) DEFAULT ''::character varying,
    redirect character varying(255) DEFAULT ''::character varying,
    idx integer,
    is_hidden integer,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone
);


ALTER TABLE public.url OWNER TO postgres;

--
-- Name: COLUMN url.url_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.url_name IS 'url名';


--
-- Name: COLUMN url.url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.url IS 'url';


--
-- Name: COLUMN url.parent_uid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.parent_uid IS '父链接';


--
-- Name: COLUMN url.title; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.title IS '标题';


--
-- Name: COLUMN url.icon; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.icon IS '图标';


--
-- Name: COLUMN url.component_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.component_name IS '组件名';


--
-- Name: COLUMN url.component_path; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.component_path IS '组件路径';


--
-- Name: COLUMN url.redirect; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.redirect IS '默认跳转';


--
-- Name: COLUMN url.idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.idx IS '顺序';


--
-- Name: COLUMN url.is_hidden; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.is_hidden IS '隐藏';


--
-- Name: COLUMN url.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.url.create_time IS '创建时间';


--
-- Name: vaccount; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.vaccount (
    vaccount_no uuid NOT NULL,
    account_no uuid NOT NULL,
    va_type bigint,
    balance bigint DEFAULT 0,
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0,
    use_status integer DEFAULT 1,
    delete_time timestamp(6) without time zone,
    update_time timestamp(6) without time zone,
    frozen_balance bigint DEFAULT 0,
    balance_type character varying(255) DEFAULT ''::character varying,
    modify_time timestamp(6) without time zone
);


ALTER TABLE public.vaccount OWNER TO postgres;

--
-- Name: COLUMN vaccount.vaccount_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.vaccount_no IS '虚拟账户';


--
-- Name: COLUMN vaccount.account_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.account_no IS '账号';


--
-- Name: COLUMN vaccount.va_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.va_type IS '虚拟账户类型';


--
-- Name: COLUMN vaccount.balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.balance IS '余额';


--
-- Name: COLUMN vaccount.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.use_status IS '0禁用';


--
-- Name: COLUMN vaccount.delete_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.delete_time IS '销户时间';


--
-- Name: COLUMN vaccount.update_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.update_time IS '最后修改时间';


--
-- Name: COLUMN vaccount.frozen_balance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.frozen_balance IS '冻结金额';


--
-- Name: COLUMN vaccount.balance_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.vaccount.balance_type IS '余额类型（usd美金，khr瑞尔）';


--
-- Name: wf_proc_running; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wf_proc_running (
    running_no character varying NOT NULL,
    process_no uuid,
    current_step uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    create_time timestamp(6) without time zone,
    run_status integer DEFAULT 0
);


ALTER TABLE public.wf_proc_running OWNER TO postgres;

--
-- Name: COLUMN wf_proc_running.current_step; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_proc_running.current_step IS '当前执行到的步骤';


--
-- Name: wf_process; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wf_process (
    process_no uuid NOT NULL,
    process_name character varying(255) DEFAULT ''::character varying,
    execute_no uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid,
    create_time timestamp(6) without time zone,
    execute_status integer,
    steps json
);


ALTER TABLE public.wf_process OWNER TO postgres;

--
-- Name: COLUMN wf_process.process_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_process.process_no IS '流程号';


--
-- Name: COLUMN wf_process.process_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_process.process_name IS '流程名';


--
-- Name: COLUMN wf_process.execute_no; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_process.execute_no IS '执行时间号';


--
-- Name: COLUMN wf_process.create_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_process.create_time IS '创建时间';


--
-- Name: COLUMN wf_process.execute_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_process.execute_status IS '执行状态';


--
-- Name: COLUMN wf_process.steps; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wf_process.steps IS '步骤集合';


--
-- Name: wf_step; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wf_step (
    step_no uuid NOT NULL,
    step_name character varying(255) DEFAULT ''::character varying,
    func_name character varying(512) DEFAULT ''::character varying,
    p1 character varying(255) DEFAULT ''::character varying,
    p2 character varying(255) DEFAULT ''::character varying,
    p3 character varying(255) DEFAULT ''::character varying,
    p4 character varying(255) DEFAULT ''::character varying,
    p5 character varying(255) DEFAULT ''::character varying,
    p6 character varying(255) DEFAULT ''::character varying,
    p7 character varying(255) DEFAULT ''::character varying,
    p8 character varying(255) DEFAULT ''::character varying,
    p9 character varying(255) DEFAULT ''::character varying,
    p10 character varying(255) DEFAULT ''::character varying,
    create_time timestamp(6) without time zone,
    is_delete smallint DEFAULT 0
);


ALTER TABLE public.wf_step OWNER TO postgres;

--
-- Name: writeoff; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.writeoff (
    code character varying(255) NOT NULL,
    income_order_no character varying(255) DEFAULT ''::character varying,
    outgo_order_no character varying(255) DEFAULT ''::character varying,
    create_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    finish_time timestamp(6) without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    use_status smallint DEFAULT 0,
    transfer_order_no character varying(255) DEFAULT ''::character varying,
    send_phone character varying(255),
    recv_phone character varying(255),
    modify_time timestamp(6) without time zone
);


ALTER TABLE public.writeoff OWNER TO postgres;

--
-- Name: COLUMN writeoff.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.writeoff.code IS '核销码';


--
-- Name: COLUMN writeoff.use_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.writeoff.use_status IS '状态1-初始状态.2-已使用状态';


--
-- Name: xlsx_file_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.xlsx_file_log (
    xlsx_task_no character varying(255) NOT NULL,
    create_time timestamp(6) without time zone,
    account_no uuid,
    file_type smallint,
    query_str character varying(255),
    role_type character varying(255)
);


ALTER TABLE public.xlsx_file_log OWNER TO postgres;

--
-- Name: COLUMN xlsx_file_log.file_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.xlsx_file_log.file_type IS '订单类型（1兑2存3取4转5收）';


--
-- Name: COLUMN xlsx_file_log.query_str; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.xlsx_file_log.query_str IS '查询条件';


--
-- Name: servicer_count_list id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer_count_list ALTER COLUMN id SET DEFAULT nextval('public.servicer_count_list_id_seq'::regclass);


--
-- Data for Name: account; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.account (uid, nickname, account, password, use_status, create_time, drop_time, modify_time, update_time, phone, email, master_acc, is_delete, usd_balance, khr_balance, gen_key, is_actived, head_portrait_img_no, last_login_time, is_first_login, app_lang, pos_lang, country_code, utm_source) FROM stdin;
4b21eb6b-96cf-439f-88c4-c023ddf6c4b1	admin2	admin2	f08587d6f466c0abc1252f15dc52a865	1	2020-04-28 14:46:02.40767	1970-01-01 00:00:00	1970-01-01 00:00:00	1970-01-01 00:00:00	123456		00000000-0000-0000-0000-000000000000	0	0	0	6608862260	1	\N	1970-01-01 00:00:00	0		\N	\N	\N
d4a4ca0f-973e-484a-80b7-c40187aeda3f	msuspg	886716519	f08587d6f466c0abc1252f15dc52a865	1	2020-05-08 10:14:48.613892	1970-01-01 00:00:00	2020-05-08 10:20:03.015926	1970-01-01 00:00:00	886716519		00000000-0000-0000-0000-000000000000	0	9990	0	4283226801	1	\N	1970-01-01 00:00:00	1		\N	855	
af8ca79c-995e-4336-a1f7-31a76613d300	admin0	admin0	43bfd7c5fedd7861054dbab63b3236d1	1	2020-04-23 12:08:30.478547	1970-01-01 00:00:00	1970-01-01 00:00:00	1970-01-01 00:00:00	123456789		00000000-0000-0000-0000-000000000000	0	0	0	3288723151	1	\N	1970-01-01 00:00:00	0		\N	\N	\N
92c0e2a2-08b3-444d-bc4b-640e1f5836cf	admin1	admin1	f08587d6f466c0abc1252f15dc52a865	1	2020-04-24 01:38:46.659063	1970-01-01 00:00:00	1970-01-01 00:00:00	1970-01-01 00:00:00	0		00000000-0000-0000-0000-000000000000	0	0	0	4645843196	1	\N	1970-01-01 00:00:00	0		\N	\N	\N
22222222-2222-2222-2222-222222222222	admin	admin	dc02fa6b1830b60c6e363e54781fb86b	1	2018-02-14 12:35:56.183946	1970-01-01 00:00:00	1970-01-01 00:00:00	1970-01-01 00:00:00	13466668888		00000000-0000-0000-0000-000000000000	0	\N	\N	0.690176161006093	1	\N	\N	0		\N	\N	\N
5e2447c5-47d4-4589-ae58-09423dd57fd7	lolo	18148769649	7dcf7a89e67edbcff0695c515aa3d53b	1	2020-05-08 13:31:26.33408	1970-01-01 00:00:00	1970-01-01 00:00:00	1970-01-01 00:00:00	18148769649		00000000-0000-0000-0000-000000000000	0	0	0	1530884939	1	\N	1970-01-01 00:00:00	0		\N	855	
30c9911c-3e77-4e5b-9f2e-5d7825b378b4	admin3	admin3	b90df30221e7600cb6b7e8ceda588eff	1	2020-04-30 11:14:28.111152	1970-01-01 00:00:00	1970-01-01 00:00:00	1970-01-01 00:00:00	123456789		00000000-0000-0000-0000-000000000000	0	0	0	5516370291	1	\N	1970-01-01 00:00:00	0		\N		\N
\.


--
-- Data for Name: account_collect; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.account_collect (account_no, collect_account_no, collect_phone, is_delete, create_time, modify_time, account_collect_no) FROM stdin;
f2961920-ef9b-4ebb-9de8-dcfcc1b21539	0e8d24af-bec7-4f95-b038-c48045f51abf	233	0	2020-04-23 12:08:11	2020-04-27 13:44:11	2020042320081105929219
c1abbc5a-4280-4996-a866-08b74a21d2fb	72ed5801-7efe-4cc0-8a43-ae84a5f8f055	1551234567	0	2020-04-27 15:42:00	2020-04-27 15:42:00	2020042716420077731532
0e8d24af-bec7-4f95-b038-c48045f51abf	2215fd8b-5b5f-40f6-b676-8a5080db660f	149	0	2020-04-28 23:00:32	2020-04-30 10:58:06	2020042900003225136321
\.


--
-- Data for Name: adminlog; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.adminlog (log_uid, create_time, url, param, op, op_type, op_acc_uid, during, ip, status_code, response) FROM stdin;
\.


--
-- Data for Name: agreement; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.agreement (id, text, lang, type, create_time, is_delete, modify_time, use_status) FROM stdin;
3	这是测试的协议内容zh_CN2	zh_CN	0	2020-02-13 12:41:16	1	2020-02-13 12:41:17	0
2020040819231790747011	2	zh_CN	0	2020-04-08 19:23:19.473544	1	2020-04-08 19:23:19.473544	0
2020040818571357616254	2	zh_CN	0	2020-04-08 18:57:14.983108	1	2020-04-08 18:57:14.983108	0
1	这是测试的协议内容km1	km	0	2020-02-12 16:52:58	1	2020-04-08 17:02:02.688562	0
2020021312465461966160	测试隐私协议内容1234	en	1	2020-02-13 12:45:06.661878	1	2020-02-13 12:59:04.302798	0
4	Modern Pay服务协议\n一、概述\n二、注册账户及管理\n三、使用规则\n四、交易风险提示\n五、服务费用\n六、用户合法使用Modern Pay服务的承诺\n七、隐私保护\n八、不可抗力、免责及责任限制\n九、知识产权的保护\n十、其他\n\n\n一、概述\nModern Pay服务协议由您与Modern Pay依据本协议约定的方式订立，用于明确在您使用Modern Pay提供服务时双方之间的权利义务关系。本协议一经订立，将在您与Modern Pay之间具有合同上的法律效力。Modern Pay郑重提示您，请您审慎阅读并充分理解本协议全部条款，特别是其中所涉及的免除或限制Modern Pay义务与责任的条款、排除或限制用户权利的条款等，以及粗体字样标识的内容。同时您还应仔细阅读并充分理解《隐私政策》。如您不同意本协议和隐私政策及后续Modern Pay做出的修订，请不用或立即停止使用Modern Pay服务。\n二、、注册账户及管理\n(一)注册\n1、您在使用Modern Pay服务前，应当注册并登录Modern Pay应用。\n2、您有义务在主项时依据Modern Pay平台要求提供自己真实、准确、完整、合法、有效的身份等信息；不得使用他人信息注册、实名认证或向平台提供其他虚假信息。\n3、在注册后，如发现您以虚假信息骗取注册，或您头像、昵称等注册信息存在违法或不良信息的，Modern Pay有权不经通知单方采取暂停使用、注销登记等措施。\n(二)使用\n1、您应当妥善保管自己的账号、密码及其他有效识别信息。因您的原因造成的账户、密码等信息被冒用、盗用或非法使用,由此引起的风险和损失需要由您自行承担。\n2、除本协议另有约定，您不应将账号、密码以及账户转让、出售或出借他人使用。\n3、基于运行和交易安全的需要,我们可能会暂停或者限制Modern Pay服务部分功能,或增加新的功能。\n4、为了维护良好的网络环境,我们有时需要了解您使用Modern Pay服务的真实背景及目的,如我们要求您提供相关信息或资料的,请您配合提供。\n5、为了您的交易安全,在使用Modern Pay服务时,请您事先自行核实交易对方的身份信息(如交易对方是否具有完全民事行为能力)并谨慎决定是否使用Modern Pay服务与对方进行交易。\n(三)注销\n在需要终止使用Modern Pay服务时,符合以下条件的,您可以申请注销您的Modern Pay账户：\n1、您仅能申请注销您本人名义注册的Modern Pay账户,并依照Modern Pay平台规则进行注销。\n2、您可以通过人工的方式申请注销用户标识或账户,但如果您使用了我们提供的安全产品,请在该安全产品环境下申请注销。\n3、您申请注销的账户处于正常状态,即您的账户的信息是最新、完整、正确的,且用户账户未被采取止付、有权机关冻结等限制措施。\n4、为了维护您和其他用户的合法利益,您申请注销的用户标识或账户,应当不存在未了结的权利义务或其他因为注销该账户会产生纠纷的情况,不存在未完结交易,没有余额。\n5、注销成功后，账户信息、个人身份信息、交易记录等将无法恢复或提供。\n三、使用规则\n1、您在使用Modern Pay服务时，必须遵守柬埔寨王国相关法律法规的规定，不得利用Modern Pay平台从事任何违法违规的活动。\n2、Modern Pay平台有权依照法律法规、监管规定及平台规则采取措施防范欺诈行为。如发现欺诈行为或其他损害注册会员及平台利益的情形时，有权及时公告并暂停或终止您在平台的相关活动。\n3、为了满足相关监管规定的要求,请您按照我们要求的时间提供您的身份信息以完成身份验证,否则您可能无法进行收款、提现、余额支付、兑换等操作,且我们可能对您的账户余额进行止付或注销您的账户。\n4、用户违反本协议或相关的服务条款的规定，导致或产生的任何Modern Pay或第三方主张的任何索赔、要求或损失，包括合理的律师费，您同意赔偿Modern Pay、合作机构，并使之免受损害。对此，Modern Pay有权视用户的行为性质，采取包括但不限于删除用户发布信息内容、暂停使用许可、终止服务、回收账户、追究法律责任等措施。\n四、交易风险提示\n1、在使用我们的服务时,若您或您的交易对方未遵从本协议或相关网站说明、交易、支付页面中的操作提示、规则,则我们有权拒绝为您与交易对方提供服务,且我们免于承担损害赔偿责任。\n2、请您特别注意,如在Modern Pay平台上以页面标明或其他方式表明相关服务系由第三方提供,您在使用该服务过程中如有疑问或发生纠纷,请您与第三方协商解决。\n五、服务费用\n在您使用我们的服务时,我们有权依照Modern Pay平台规则向您收取服务费用。我们拥有制订及调整服务费的权利,具体服务费用以您使用我们服务时Modern Pay产品页面上所列的收费方式公告或您与我们达成的其他书面协议为准。除非另有说明或约定,您同意我们有权自您委托我们代为收取或代为支付的款项或您任一Modern Pay账户余额或者其他资产中直接扣除上述服务费用。\n六、用户合法使用Modern Pay服务的承诺\n1、您需要柬埔寨王国相关法律法规及您所属、所居住或开展经营活动或其他业务的国家或地区的法律法规,不得将我们的服务用于非法目的(包括用于禁止或限制交易物品的交易),也不得以非法方式使用我们的服务。\n2、您不得利用我们的服务从事侵害他人合法权益之行为或违反国家法律法规,否则我们有权进行调查、延迟或拒绝结算或停止提供服务,且您需要自行承担相关法律责任,如因此导致我们或其他方受损的,您需要承担赔偿责任。\n3、上述1和2适用的情況包括但不限于\n(1)侵害他人名誉权、隐私权、商业秘密商标权、著作权、专利权等合法权益;\n(2)违反保密义务;\n(3)冒用他人名义使用我们的服务;\n(4)从事不法交易行为,如洗钱、恐怖融资、赌博、贩卖枪支、毒品、禁药、盗版软件、黄色淫秽物品、其他我们认为不得使用我们的服务进行交易的物品等；\n(5)提供赌博资讯或以任何方式引诱他人参与赌博;\n(6)未经授权使用他人银行卡,或利用信用卡、花呗套取现金;\n(7)进行与您或交易对方宣称的交易内容不符的交易,或不真实的易；\n(8)从事可能侵害我们的服务系统、数据的行为。\n4、账户仅限本人使用,您需充分了解并清楚知晓出租、出借、出售、购买账户的相关法律责任和惩戒措施,承诺依法依规开立和使用本人账户。\n您理解,我们的服务有赖于系统的准确运行及操作。若出现系统差错、故障、您或我们不当获利等情形,您同意我们可以采取更正差错、扣划款项等适当纠正措施。\n6、您不得对我们的系统和程序采取反向工程手段进行破解,不得对上述系统和程序(包括但不限于源程序、目标程序、技术文档、客户端至服务器端的数据、服务器数据)进行复制、修改、编译、整合或篡改,不得修改或增减我们系统的功能。\n七、隐私保护\n我们重视对您信息的保护。如果您是个人用户,关于您的个人信息依《Modern Pay私权政策》受到保护与规范,详情请参阅《Modern Pay私权政策》。如果您是企业用户,我们将依法保护您的商业秘密,非经您同意不会对外提供。\n八、不可抗力、免责及责任限制\n(一)免责条款\n因下列原因导致我们无法正常提供服务,我们免于承担责任\n1、我们的系统停机维护或升级;\n2、因台风、地震、洪水、雷电或恐怖袭击等不可抗力原因；\n3、您的电脑软硬件和通信线路、供电线路出现故障的;\n4、因您操作不当或通过非经我们授权或认可的方式使用我们服务的；\n5、因病毒、木马、恶意程序攻击、网络拥堵、系统不稳定、系统或设备故障、通讯故障、电力故障、银行原因、第三方服务瑕疵或政府行为等原因。尽管有前款约定,我们将采取合理行动积极促使服务恢复正常。\n(二)责任限制\n我们可能同时为您及您的(交易)对手方提供服务,您同意对我们可能存在的该等行为予以明确豁免任何实际或潜在的利益冲突,并不得以此来主张我们在提供服务时存在法律上的瑕疵。\n九、知识产权的保护\n1、我们及关联公司的系统及支付宝网站上的内容,包括但不限于著作、图片、档案、资讯、资料、网站架构、网站画面的安排、网页设计,均由我们或关联公司依法拥有知识产权,包括但不限于商标权、专利权、著作权商业秘密等。\n2、非经我们或关联公司书面同意,请勿擅自使用、修改、反向编译、复制、公开传播、改变、散布、发行或公开发表支付宝网站程序或\n内容。\n3、尊重知识产权是您应尽的义务,如有违反,您需要承担损害赔偿责任。\n十、其他\n1、本协议之效力、解释、变更、执行与争议解决均适用柬埔寨王国法律。因本协议产生的争议,均应依照柬埔寨王国法律予以处理。\n2、Modern Pay未行使或执行本服务协议任何权利或规定，不构成对前述权利或权利之放弃。\n3、如本协议中的任何条款无论因何种原因完全或部分无效或不具有执行力，本协议的其余条款仍应有效并且有约束力。	zh_CN	0	2020-02-13 12:44:14	0	2020-04-09 17:31:57.542872	1
2	Modern Pay Service Agreement\nI. Overview\n2. Account Registration and Management\n3. Rules of use\nFourth, transaction risk tips\nV. Service fees\n6. The user's commitment to legally use the Modern Pay service\nSeven, privacy protection\n8. Force Majeure, Disclaimer and Limitation of Liability\n9. Protection of Intellectual Property\nTen, other\n\n\nI. Overview\nThe Modern Pay service agreement is concluded between you and Modern Pay in accordance with the method stipulated in this agreement. It is used to clarify the rights and obligations between the two parties when you use Modern Pay to provide services. Once this agreement is concluded, it will have the legal effect of the contract between you and Modern Pay. Modern Pay solemnly reminds you that you should carefully read and fully understand all the terms of this agreement, especially the clauses that exempt or restrict Modern Pay ’s obligations and responsibilities, the clauses that exclude or restrict user rights, etc. content. At the same time, you should also read and fully understand the "Privacy Policy". If you do not agree to this agreement and the privacy policy and subsequent amendments made by Modern Pay, please do not use or immediately stop using the Modern Pay service.\nSecond, account registration and management\n(1) Registration\n1. Before you use the Modern Pay service, you should register and log in to the Modern Pay application.\n2. You are obliged to provide your true, accurate, complete, legal, and valid identity information according to the requirements of the Modern Pay platform at the time of the main project; you must not use other people's information to register, real-name authentication or provide other false information to the platform.\n3. After registration, if you find that you have fraudulently obtained registration with false information, or if your registration information such as avatar and nickname has illegal or bad information, Modern Pay reserves the right to take measures such as suspension of use and cancellation of registration without notice.\n(Two) use\n1. You should properly keep your account number, password and other valid identification information. Accounts, passwords, and other information caused by you are used, misappropriated, or illegally used, and you must bear the risks and losses caused thereby.\n2. Except as otherwise stipulated in this agreement, you should not transfer, sell, or lend your account number, password, and account for use by others.\n3. Based on the needs of operation and transaction security, we may suspend or restrict some functions of the Modern Pay service, or add new functions.\n4. In order to maintain a good network environment, we sometimes need to understand the real background and purpose of your use of the Modern Pay service. If we require you to provide relevant information or information, please cooperate.\n5. For the safety of your transactions, when using the Modern Pay service, please verify the identity information of the counterparty (such as whether the counterparty has full civil capacity) and carefully decide whether to use the Modern Pay service to conduct transactions with the counterparty.\n(3) Cancellation\nWhen you need to terminate the use of Modern Pay service, if you meet the following conditions, you can apply to cancel your Modern Pay account:\n1. You can only apply for cancellation of the Modern Pay account registered in your own name, and cancel in accordance with the rules of the Modern Pay platform.\n2. You can apply to cancel the user ID or account manually, but if you use the security products provided by us, please apply for cancellation under the security product environment.\n3. The account you applied for cancellation is in a normal state, that is, the information of your account is up-to-date, complete, and correct, and the user account has not been subject to restrictive measures such as stop payment or freezing by the authority.\n4. In order to protect the legal interests of you and other users, the user ID or account you have applied for cancellation should not have any unresolved rights and obligations or other circumstances that will cause disputes due to cancellation of the account. There will be no outstanding transactions and no balance.\n5. After successful cancellation, account information, personal identification information, transaction records, etc. will not be restored or provided.\n3. Rules of use\n1. When you use the Modern Pay service, you must abide by the relevant laws and regulations of the Kingdom of Cambodia, and you must not use the Modern Pay platform to engage in any illegal activities.\n2. The Modern Pay platform has the right to take measures to prevent fraud in accordance with laws and regulations, regulatory regulations and platform rules. If fraudulent behavior or other circumstances that damage the interests of registered members and the platform are found, they have the right to announce and suspend or terminate your relevant activities on the platform in a timely manner.\n3. In order to meet the requirements of the relevant regulatory regulations, please provide your identity information in accordance with the time required by us to complete the identity verification, otherwise you may not be able to collect, withdraw, balance payment, exchange and other operations, and we may The account balance can be used to stop payment or cancel your account.\n4. If the user violates the provisions of this agreement or the related service terms, any claims, requirements or losses caused or incurred by Modern Pay or a third party, including reasonable attorney fees, you agree to compensate Modern Pay, the cooperative agency, and make From damage. In this regard, Modern Pay has the right to take measures including, but not limited to, deleting the content of the user's information, suspending the use of the license, terminating the service, recovering the account, and holding legal responsibility, depending on the nature of the user's behavior.\nFourth, transaction risk tips\n1.When using our services, if you or your counterparty does not comply with this agreement or the operation tips and rules on the relevant website instructions, transactions, and payment pages, we have the right to refuse to provide services for you and the counterparty, and We are exempt from liability for damages.\n2. Please pay special attention. If you indicate on the Modern Pay platform that the related services are provided by a third party, or if you have any questions or disputes during the use of the service, please consult with the third party to resolve them.\nV. Service fees\nWhen you use our services, we have the right to charge you service fees in accordance with the rules of the Modern Pay platform. We have the right to formulate and adjust the service fee. The specific service fee is subject to the announcement of the charging method listed on the Modern Pay product page when you use our service or other written agreement between you and us. Unless otherwise stated or agreed, you agree that we have the right to deduct the above service fees directly from the money you entrust us to collect or pay on your behalf or any of your Modern Pay account balances or other assets.\n6. The user's commitment to legally use the Modern Pay service\n1. You need the relevant laws and regulations of the Kingdom of Cambodia and the laws and regulations of the country or region where you belong, live or carry out business activities or other business, and may not use our services for illegal purposes (including for prohibiting or restricting transactions of trading items) ), Nor can we use our services illegally.\n2. You must not use our services to infringe on the legitimate rights and interests of others or violate national laws and regulations, otherwise we have the right to investigate, delay or refuse settlement or stop providing services, and you need to bear the relevant legal responsibilities yourself, if this results in us If any other party is damaged, you need to be liable for compensation.\n3. The situations where 1 and 2 above apply include but are not limited to\n(1) Infringe on the legitimate rights and interests of others' reputation, privacy, trade secret trademark rights, copyrights, patent rights, etc .;\n(2) Violation of confidentiality obligations;\n(3) Using our services in the name of others;\n(4) Engaging in illegal transactions, such as money laundering, terrorist financing, gambling, trafficking in firearms, drugs, banned drugs, pirated software, pornographic obscene items, and other items that we believe are not allowed to use our services for transactions;\n(5) Provide gambling information or induce others to participate in gambling in any way;\n(6) Unauthorized use of other people's bank cards, or use of credit cards and flowers to withdraw cash;\n(7) Conduct a transaction that is inconsistent with the transaction content declared by you or the counterparty of the transaction, or is untrue;\n(8) Engage in actions that may infringe our service systems and data.\n4. The account is for personal use only. You need to fully understand and clearly understand the relevant legal responsibilities and disciplinary measures for leasing, lending, selling, and purchasing accounts, and promise to open and use your account in accordance with the law.\nYou understand that our service depends on the accurate operation and operation of the system. In the event of system errors, malfunctions, improper profits made by you or us, you agree that we can take appropriate corrective measures such as correcting errors and deducting payments.\n6. You must not use reverse engineering methods to crack our systems and programs, and must not copy the above systems and programs (including but not limited to source programs, target programs, technical documents, client-to-server data, server data) , Modify, compile, integrate or tamper, and must not modify or increase or decrease the function of our system.\nSeven, privacy protection\nWe value the protection of your information. If you are an individual user, your personal information is protected and regulated in accordance with the "Modern Pay Private Policy". For details, please refer to the "Modern Pay Private Policy". If you are an enterprise user, we will protect your business secrets according to law and will not provide it to you without your consent.\n8. Force Majeure, Disclaimer and Limitation of Liability\n(1) Disclaimer\nWe are unable to provide services normally due to the following reasons, and we are exempt from liability\n1. Our system is shut down for maintenance or upgrade;\n2. Force majeure reasons such as typhoon, earthquake, flood, lightning or terrorist attack;\n3. Your computer hardware and software, communication lines, and power supply lines fail;\n4. Due to your improper operation or use of our services by means not authorized or approved by us;\n5. Due to viruses, Trojan horses, malicious program attacks, network congestion, system instability, system or equipment failure, communication failure, power failure, bank reasons, third-party service defects or government actions and other reasons. Despite the agreement in the preceding paragraph, we will take reasonable action to actively promote service to normal.\n(2) Limitation of liability\nWe may provide services to you and your (transaction) counterparty at the same time, and you agree to expressly waive any actual or potential conflicts of interest in such actions that we may have, and shall not use this to claim that we have laws when providing services Flaws.\n9. Protection of Intellectual Property\n1.The contents of our and our affiliated companies' systems and Alipay's website, including but not limited to works, pictures, files, information, data, website structure, website screen arrangements, and web design, are owned by us or affiliated companies in accordance with law , Including but not limited to trademark rights, patent rights, copyright trade secrets, etc.\n2. Do not use, modify, decompile, copy, publicly distribute, change, distribute, publish or publicly publish the Alipay website program or without our written consent\ncontent.\n3. Respect for intellectual property rights is your obligation. If there is a violation, you need to be liable for damages.\nTen, other\n1. The validity, interpretation, modification, execution and dispute resolution of this agreement are all applicable to the laws of the Kingdom of Cambodia. Disputes arising from this agreement shall be handled in accordance with the laws of the Kingdom of Cambodia.\n2. Modern Pay's failure to exercise or enforce any rights or provisions of this service agreement does not constitute a waiver of the aforementioned rights or rights.\n3. If any clause in this agreement is completely or partially invalid or not enforceable for any reason, the remaining clauses of this agreement should still be valid and binding.	en	0	2020-02-13 12:40:45	0	2020-04-09 17:32:40.386737	1
2020040818545951540622	កិច្ចព្រមព្រៀងសេវាកម្មបង់ប្រាក់ទំនើប\nទិដ្ឋភាពទូទៅ\n2. ការចុះឈ្មោះនិងការគ្រប់គ្រងគណនី\n3. វិធាននៃការប្រើប្រាស់\nទីបួនគន្លឹះហានិភ័យនៃប្រតិបត្តិការ\nV. ថ្លៃសេវាកម្ម\n6. ការប្តេជ្ញាចិត្តរបស់អ្នកប្រើប្រាស់ក្នុងការប្រើសេវាកម្មបង់ប្រាក់ទំនើបតាមច្បាប់\nប្រាំពីរ, ការការពារភាពឯកជន\n៨- បង្ខំឧត្ដមភាពការបដិសេធនិងដែនកំណត់នៃការទទួលខុសត្រូវ\nការការពារកម្មសិទ្ធិបញ្ញា\nដប់, ផ្សេងទៀត\n\n\nទិដ្ឋភាពទូទៅ\nកិច្ចព្រមព្រៀងសេវាកម្មសេវាកម្មបង់ប្រាក់សម័យទំនើបត្រូវបានបញ្ចប់រវាងអ្នកនិង Modern Pay ស្របតាមវិធីសាស្ត្រដែលមានចែងក្នុងកិច្ចព្រមព្រៀងនេះ។ វាត្រូវបានប្រើដើម្បីបញ្ជាក់ពីសិទ្ធិនិងកាតព្វកិច្ចរវាងភាគីទាំងពីរនៅពេលដែលអ្នកប្រើ Modern Pay ដើម្បីផ្តល់សេវាកម្ម។ នៅពេលកិច្ចព្រមព្រៀងនេះត្រូវបានបញ្ចប់វានឹងមានប្រសិទ្ធិភាពស្របច្បាប់នៃកិច្ចសន្យារវាងអ្នកនិងប្រាក់ខែទំនើប។ Modern Pay រំលឹកយ៉ាងឱឡារិកថាអ្នកគួរតែអាននិងយល់ច្បាស់ពីល័ក្ខខ័ណ្ឌទាំងអស់នៃកិច្ចព្រមព្រៀងនេះជាពិសេសមាត្រាដែលលើកលែងឬកំណត់កាតព្វកិច្ចនិងការទទួលខុសត្រូវរបស់ Modern Pay ដែលជាប្រការដែលមិនរាប់បញ្ចូលឬរឹតត្បិតសិទ្ធិរបស់អ្នកប្រើប្រាស់។ ល។ ខ្លឹមសារ។ ទន្ទឹមនឹងនេះអ្នកក៏គួរអាននិងយល់ច្បាស់អំពី“ គោលការណ៍ឯកជន” ។ ប្រសិនបើអ្នកមិនយល់ព្រមលើកិច្ចព្រមព្រៀងនេះនិងគោលការណ៍ភាពឯកជននិងការកែប្រែជាបន្តបន្ទាប់ដែលធ្វើឡើងដោយ Modern Pay សូមកុំប្រើឬបញ្ឈប់ការប្រើប្រាស់សេវាកម្ម Modern Pay ភ្លាមៗ។\nទីពីរការចុះឈ្មោះនិងការគ្រប់គ្រងគណនី\n(១) ការចុះឈ្មោះ\nមុនពេលដែលអ្នកប្រើសេវាកម្មបង់ប្រាក់ទំនើបអ្នកគួរតែចុះឈ្មោះនិងចូលទៅក្នុងពាក្យសុំបង់ប្រាក់ទំនើប។\n2. អ្នកមានកាតព្វកិច្ចផ្តល់ព័ត៌មានអត្តសញ្ញាណត្រឹមត្រូវពេញលេញពេញលេញស្របច្បាប់និងមានសុពលភាពរបស់អ្នកយោងទៅតាមតំរូវការនៃវេទិកាបង់ប្រាក់ទំនើបនៅពេលគម្រោងសំខាន់អ្នកមិនត្រូវប្រើប្រាស់ព័ត៌មានរបស់អ្នកដទៃដើម្បីចុះឈ្មោះការផ្ទៀងផ្ទាត់ឈ្មោះពិតរឺផ្តល់ព័ត៌មានមិនពិតផ្សេងទៀតទៅកាន់វេទិកានោះទេ។\n3. បន្ទាប់ពីការចុះឈ្មោះប្រសិនបើអ្នកឃើញថាអ្នកបានទទួលការក្លែងបន្លំការចុះឈ្មោះជាមួយព័ត៌មានមិនពិតឬប្រសិនបើព័ត៌មានចុះឈ្មោះរបស់អ្នកដូចជាអាវ៉ាតានិងសម្មតិនាមមានព័ត៌មានមិនត្រឹមត្រូវឬអាក្រក់នោះក្រុមហ៊ុន Modern Pay សូមរក្សាសិទ្ធិក្នុងការចាត់វិធានការដូចជាការព្យួរការប្រើប្រាស់និងការលុបចោលការចុះឈ្មោះដោយមិនចាំបាច់ជូនដំណឹងជាមុន។\n(ពីរ) ប្រើ\nអ្នកគួរតែរក្សាទុកលេខគណនីលេខសម្ងាត់និងព័ត៌មានអត្តសញ្ញាណប័ណ្ណផ្សេងទៀតដែលត្រឹមត្រូវ។ គណនីលេខសម្ងាត់និងព័ត៌មានផ្សេងទៀតដែលបណ្តាលមកពីអ្នកត្រូវបានប្រើប្រាស់មិនត្រឹមត្រូវឬត្រូវបានប្រើប្រាស់ដោយខុសច្បាប់ហើយអ្នកត្រូវតែទទួលនូវហានិភ័យនិងការបាត់បង់ដែលបណ្តាលមកពី។\n2. លើកលែងតែមានចែងក្នុងកិច្ចព្រមព្រៀងនេះអ្នកមិនគួរផ្ទេរលក់រឺអោយខ្ចីលេខគណនីលេខសំងាត់និងគណនីសំរាប់ប្រើប្រាស់ដោយអ្នកដទៃឡើយ។\n3. ដោយផ្អែកលើតំរូវការនៃប្រតិបត្តិការនិងសុវត្ថិភាពប្រតិបត្តិការយើងអាចនឹងផ្អាករឺរឹតត្បិតមុខងារមួយចំនួននៃសេវាកម្មបង់ប្រាក់ទំនើបរឺបន្ថែមមុខងារថ្មី។\n4. ដើម្បីថែរក្សាបរិដ្ឋានបណ្តាញល្អពេលខ្លះយើងត្រូវស្វែងយល់អំពីសាវតានិងគោលបំណងពិតនៃការប្រើប្រាស់សេវាកម្មប្រាក់បៀវត្សរ៍ទំនើបរបស់អ្នកប្រសិនបើយើងតម្រូវឱ្យអ្នកផ្តល់នូវព័ត៌មានឬព័ត៌មានពាក់ព័ន្ធសូមសហការ។\n5. ដើម្បីសុវត្ថិភាពនៃប្រតិបត្តិការរបស់អ្នកនៅពេលប្រើសេវាកម្មបង់ប្រាក់ទំនើបសូមផ្ទៀងផ្ទាត់ព័ត៌មានអត្តសញ្ញាណរបស់អ្នកដែលមានប្រាក់ខែ (ដូចជាអ្នកមានសមត្ថភាពស៊ីវិលពេញលេញ) និងសម្រេចចិត្តដោយប្រុងប្រយ័ត្នថាតើត្រូវប្រើសេវាកម្ម Modern Pay ដើម្បីធ្វើប្រតិបត្តិការជាមួយអ្នកបង់ពន្ធឬអត់។\n(៣) ការលុបចោល\nនៅពេលដែលអ្នកត្រូវការបញ្ឈប់ការប្រើប្រាស់សេវាកម្មបង់ប្រាក់ទំនើបប្រសិនបើអ្នកបំពេញតាមល័ក្ខខ័ណ្ឌដូចខាងក្រោមអ្នកអាចដាក់ពាក្យសុំលុបចោលគណនីបង់ប្រាក់ទំនើបរបស់អ្នក៖\nអ្នកអាចដាក់ពាក្យសុំការលុបចោលគណនីបង់ប្រាក់ទំនើបដែលបានចុះឈ្មោះក្នុងឈ្មោះផ្ទាល់របស់អ្នកហើយលុបចោលដោយអនុលោមទៅតាមវិធាននៃវេទិកាបង់ប្រាក់ទំនើប។\n2. អ្នកអាចដាក់ពាក្យសុំលុបចោលលេខសម្គាល់អ្នកប្រើឬគណនីដោយដៃប៉ុន្តែប្រសិនបើអ្នកប្រើផលិតផលសុវត្ថិភាពដែលផ្តល់ដោយយើងសូមដាក់ពាក្យសុំលុបចោលក្រោមបរិយាកាសផលិតផលសុវត្ថិភាព។\n៣- គណនីដែលអ្នកបានស្នើសុំលុបចោលគឺស្ថិតក្នុងស្ថានភាពធម្មតាពោលគឺព័ត៌មាននៃគណនីរបស់អ្នកមានភាពទាន់សម័យពេញលេញនិងត្រឹមត្រូវហើយគណនីអ្នកប្រើប្រាស់មិនត្រូវបានដាក់កម្រិតដូចជាការបញ្ឈប់ការទូទាត់ឬការបង្កកដោយអាជ្ញាធរឡើយ។\n4. ដើម្បីការពារផលប្រយោជន៍ស្របច្បាប់របស់អ្នកនិងអ្នកប្រើប្រាស់ដទៃទៀតលេខសម្គាល់អ្នកប្រើឬគណនីដែលអ្នកបានដាក់ពាក្យសុំលុបចោលមិនគួរមានសិទ្ធិនិងកាតព្វកិច្ចដែលមិនទាន់បានដោះស្រាយនិងកាលៈទេសៈផ្សេងទៀតដែលនឹងបង្កឱ្យមានជម្លោះដោយសារតែការលុបចោលគណនីនោះនឹងមិនមានប្រតិបត្តិការលេចធ្លោនិងគ្មានសមតុល្យឡើយ។\n5. បន្ទាប់ពីការលុបចោលដោយជោគជ័យព័ត៌មានគណនីព័ត៌មានអត្តសញ្ញាណផ្ទាល់ខ្លួនកំណត់ត្រាប្រតិបតិ្តការ។ ល។ នឹងមិនត្រូវបានស្តារឬផ្តល់ជូនឡើយ។\n3. វិធាននៃការប្រើប្រាស់\nនៅពេលដែលអ្នកប្រើសេវាកម្មបង់ប្រាក់ទំនើបអ្នកត្រូវតែគោរពតាមច្បាប់និងបទប្បញ្ញត្តិពាក់ព័ន្ធនានានៃព្រះរាជាណាចក្រកម្ពុជាហើយអ្នកមិនត្រូវប្រើប្រាស់វេទិកាបង់ប្រាក់ទំនើបដើម្បីចូលរួមក្នុងសកម្មភាពខុសច្បាប់ណាមួយឡើយ។\n2. វេទិកាប្រាក់ខែទំនើបមានសិទ្ធិចាត់វិធានការដើម្បីការពារការក្លែងបន្លំស្របតាមច្បាប់និងបទប្បញ្ញត្តិបទប្បញ្ញត្តិនិងច្បាប់វេទិកា។ ប្រសិនបើអាកប្បកិរិយាក្លែងបន្លំឬកាលៈទេសៈផ្សេងទៀតដែលធ្វើឱ្យខូចផលប្រយោជន៍របស់សមាជិកដែលបានចុះឈ្មោះនិងវេទិកាត្រូវបានរកឃើញពួកគេមានសិទ្ធិប្រកាសនិងផ្អាកឬបញ្ឈប់សកម្មភាពពាក់ព័ន្ធរបស់អ្នកនៅលើវេទិកាឱ្យបានទាន់ពេលវេលា។\n3. ដើម្បីបំពេញតាមតម្រូវការនៃបទប្បញ្ញត្តិពាក់ព័ន្ធសូមផ្តល់ព័ត៌មានអំពីអត្តសញ្ញាណរបស់អ្នកស្របតាមពេលវេលាដែលយើងតម្រូវឱ្យបំពេញការផ្ទៀងផ្ទាត់អត្តសញ្ញាណបើមិនដូច្នេះទេអ្នកមិនអាចប្រមូលដកប្រាក់ទូទាត់ប្រាក់ផ្លាស់ប្តូរនិងប្រតិបត្តិការផ្សេងទៀតបានទេហើយយើងអាច សមតុល្យគណនីអាចត្រូវបានប្រើដើម្បីបញ្ឈប់ការទូទាត់ឬលុបចោលគណនីរបស់អ្នក។\n4. ប្រសិនបើអ្នកប្រើប្រាស់រំលោភលើបទប្បញ្ញត្តិនៃកិច្ចព្រមព្រៀងនេះឬល័ក្ខខ័ណ្ឌសេវាកម្មដែលពាក់ព័ន្ធការទាមទារការទាមទារឬការបាត់បង់ណាមួយដែលបណ្តាលមកពីឬកើតឡើងដោយប្រាក់ខែទំនើបឬតតិយជនរួមទាំងថ្លៃមេធាវីសមរម្យអ្នកយល់ព្រមទូទាត់សំណងទំនើបកម្មភ្នាក់ងារសហប្រតិបត្តិការនិងធ្វើឱ្យ ពីការខូចខាត។ ទាក់ទងនឹងបញ្ហានេះ Modern Pay មានសិទ្ធិចាត់វិធានការរាប់បញ្ចូលប៉ុន្តែមិនមានកំណត់ចំពោះការលុបមាតិកាព័ត៌មានរបស់អ្នកប្រើប្រាស់ផ្អាកការប្រើប្រាស់អាជ្ញាប័ណ្ណបញ្ឈប់សេវាកម្មដកប្រាក់គណនីនិងទទួលខុសត្រូវផ្នែកច្បាប់អាស្រ័យលើចរិតលក្ខណៈរបស់អ្នកប្រើប្រាស់។\nទីបួនគន្លឹះហានិភ័យនៃប្រតិបត្តិការ\n1. នៅពេលប្រើសេវាកម្មរបស់យើងប្រសិនបើអ្នកឬដៃគូរបស់អ្នកមិនគោរពតាមកិច្ចព្រមព្រៀងនេះឬការណែនាំប្រតិបត្តិការនិងច្បាប់ស្តីពីការណែនាំគេហទំព័រប្រតិបត្តិការនិងទំព័របង់ប្រាក់ពាក់ព័ន្ធយើងមានសិទ្ធិបដិសេធក្នុងការផ្តល់សេវាកម្មសម្រាប់អ្នកនិងអ្នកដែលសមភាគីហើយ យើងត្រូវបានរួចផុតពីការទទួលខុសត្រូវចំពោះការខូចខាត។\n2. សូមយកចិត្តទុកដាក់ជាពិសេសប្រសិនបើអ្នកចង្អុលបង្ហាញនៅលើវេទិកាបង់ប្រាក់ទំនើបថាសេវាកម្មដែលពាក់ព័ន្ធត្រូវបានផ្តល់ជូនដោយភាគីទីបីឬប្រសិនបើអ្នកមានសំណួរឬជម្លោះណាមួយក្នុងកំឡុងពេលប្រើប្រាស់សេវាកម្មសូមពិគ្រោះជាមួយភាគីទីបីដើម្បីដោះស្រាយ។\nV. ថ្លៃសេវាកម្ម\nនៅពេលអ្នកប្រើសេវាកម្មរបស់យើងយើងមានសិទ្ធិគិតថ្លៃសេវាកម្មដោយអនុលោមតាមច្បាប់នៃវេទិកាបង់ប្រាក់ទំនើប។ យើងមានសិទ្ធិក្នុងការបង្កើតនិងកែសំរួលថ្លៃសេវាកម្មថ្លៃសេវាជាក់លាក់គឺត្រូវអនុវត្តតាមវិធីសាស្រ្តនៃការគិតថ្លៃដែលបានរាយនៅលើទំព័រផលិតផលបង់ប្រាក់ទំនើបនៅពេលដែលអ្នកប្រើសេវាកម្មរបស់យើងឬកិច្ចព្រមព្រៀងជាលាយលក្ខណ៍អក្សររវាងអ្នកនិងយើង។ លើកលែងតែមានការបញ្ជាក់ឬយល់ព្រមអ្នកយល់ព្រមថាយើងមានសិទ្ធិដកថ្លៃសេវាខាងលើដោយផ្ទាល់ពីប្រាក់ដែលអ្នកបានប្រគល់ឱ្យយើងដើម្បីប្រមូលឬទូទាត់ប្រាក់ជំនួសអ្នកឬសមតុល្យគណនីបង់ប្រាក់ទំនើបឬទ្រព្យសម្បត្តិផ្សេងទៀត។\n6. ការប្តេជ្ញាចិត្តរបស់អ្នកប្រើប្រាស់ក្នុងការប្រើសេវាកម្មបង់ប្រាក់ទំនើបតាមច្បាប់\nអ្នកត្រូវការច្បាប់និងបទប្បញ្ញត្តិពាក់ព័ន្ធនៃព្រះរាជាណាចក្រកម្ពុជានិងច្បាប់និងបទបញ្ញត្តិនៃប្រទេសឬតំបន់ដែលអ្នកជាកម្មសិទ្ធិរស់នៅឬអនុវត្តសកម្មភាពអាជីវកម្មឬអាជីវកម្មផ្សេងទៀត។ ) យើងក៏មិនអាចប្រើប្រាស់សេវាកម្មរបស់យើងដោយខុសច្បាប់ដែរ។\n2. អ្នកមិនត្រូវប្រើសេវាកម្មរបស់យើងដើម្បីរំលោភលើសិទ្ធិនិងផលប្រយោជន៍ស្របច្បាប់របស់អ្នកដទៃឬរំលោភច្បាប់និងបទបញ្ញត្តិជាតិឡើយបើមិនដូច្នេះទេយើងមានសិទ្ធិស៊ើបអង្កេតពន្យារពេលឬបដិសេធការដោះស្រាយឬបញ្ឈប់ការផ្តល់សេវាកម្មហើយអ្នកត្រូវទទួលខុសត្រូវទំនួលខុសត្រូវដែលពាក់ព័ន្ធដោយខ្លួនឯងប្រសិនបើវានាំយើង ប្រសិនបើមានគណបក្សផ្សេងទៀតខូចខាតអ្នកចាំបាច់ត្រូវទទួលខុសត្រូវចំពោះសំណង។\n៣. ស្ថានភាពដែលទី ១ និងទី ២ ខាងលើអនុវត្តរួមមានប៉ុន្តែមិនមានកំណត់ទេ\n(១) រំលោភលើសិទ្ធិនិងផលប្រយោជន៍ស្របច្បាប់របស់កេរ្តិ៍ឈ្មោះអ្នកដទៃសិទ្ធិឯកជនពាណិជ្ជសញ្ញាសម្ងាត់សិទ្ធិអ្នកនិពន្ធសិទ្ធិអ្នកនិពន្ធកម្មសិទ្ធិបញ្ញា។ ល។\n(២) ការបំពានលើកាតព្វកិច្ចរក្សាការសម្ងាត់។\n(៣) ការប្រើប្រាស់សេវាកម្មរបស់យើងក្នុងនាមអ្នកដទៃ។\n(៤) ចូលរួមក្នុងប្រតិបត្តិការខុសច្បាប់ដូចជាការលាងលុយកខ្វក់ហិរញ្ញប្បទានភេរវកម្មល្បែងស៊ីសងការជួញដូរអាវុធគ្រឿងញៀនការហាមឃាត់គ្រឿងញៀនកម្មវិធីលួចចម្លងធាតុអាសអាភាសនិងរបស់ផ្សេងទៀតដែលយើងជឿថាមិនត្រូវបានអនុញ្ញាតឱ្យប្រើប្រាស់សេវាកម្មរបស់យើងសម្រាប់ប្រតិបត្តិការឡើយ។\n(៥) ផ្តល់ព័ត៌មានអំពីល្បែងស៊ីសងឬជំរុញអោយអ្នកដទៃចូលរួមលេងល្បែងតាមរបៀបណាមួយ។\n(៦) ការប្រើប្រាស់ប័ណ្ណធនាគាររបស់អ្នកដទៃដោយគ្មានការអនុញ្ញាតឬប្រើប័ណ្ណឥណទាននិងផ្កាដើម្បីដកសាច់ប្រាក់;\n(៧) ធ្វើប្រតិបត្តិការដែលមិនស៊ីគ្នានឹងមាតិកាប្រតិបត្តិការដែលបានប្រកាសដោយអ្នកឬសមភាគីនៃប្រតិបត្តិការឬមិនពិត។\n(៨) ចូលរួមសកម្មភាពដែលអាចរំលោភប្រព័ន្ធនិងទិន្នន័យសេវាកម្មរបស់យើង។\n៤- គណនីនេះគឺសម្រាប់តែការប្រើប្រាស់ផ្ទាល់ខ្លួនប៉ុណ្ណោះ។ អ្នកត្រូវយល់ដឹងនិងយល់ច្បាស់អំពីទំនួលខុសត្រូវច្បាប់ពាក់ព័ន្ធនិងវិធានការណ៍វិន័យសម្រាប់ការជួលប្រាក់កម្ចីលក់និងទិញគណនីហើយសន្យាបើកនិងប្រើប្រាស់គណនីរបស់អ្នកស្របតាមច្បាប់។\nអ្នកយល់ថាសេវាកម្មរបស់យើងពឹងផ្អែកលើប្រតិបត្តិការនិងប្រតិបត្តិការត្រឹមត្រូវនៃប្រព័ន្ធ។ ក្នុងករណីមានកំហុសប្រព័ន្ធដំណើរការមិនត្រឹមត្រូវប្រាក់ចំណេញមិនត្រឹមត្រូវដែលធ្វើឡើងដោយអ្នកឬយើងអ្នកយល់ព្រមថាយើងអាចចាត់វិធានការកែតម្រូវសមស្របដូចជាការកែកំហុសនិងកាត់ការទូទាត់។\n6. អ្នកមិនត្រូវប្រើវិធីសាស្រ្តវិស្វកម្មបញ្ច្រាសដើម្បីបំបែកប្រព័ន្ធនិងកម្មវិធីរបស់យើងនិងមិនត្រូវចម្លងប្រព័ន្ធនិងកម្មវិធីខាងលើ (រួមបញ្ចូលតែមិនកំណត់ចំពោះកម្មវិធីប្រភពកម្មវិធីគោលដៅឯកសារបច្ចេកទេសឯកសារអតិថិជន - ទិន្នន័យរបស់ម៉ាស៊ីនមេទិន្នន័យ) ។ , កែប្រែ, ចងក្រង, រួមបញ្ចូលឬរំខានហើយមិនត្រូវកែប្រែឬបង្កើនឬបន្ថយមុខងារនៃប្រព័ន្ធរបស់យើងឡើយ។\nប្រាំពីរ, ការការពារភាពឯកជន\nយើងវាយតម្លៃខ្ពស់ចំពោះការការពារព័ត៌មានរបស់អ្នក។ ប្រសិនបើអ្នកជាអ្នកប្រើប្រាស់ម្នាក់ៗព័ត៌មានផ្ទាល់ខ្លួនរបស់អ្នកត្រូវបានការពារនិងធ្វើឱ្យស្របទៅនឹងគោលការណ៍ឯកជនដែលមានប្រាក់ខែទំនើប។ សម្រាប់ព័ត៌មានលម្អិតសូមយោងលើ“ គោលការណ៍ឯកជននៃប្រាក់ខែទំនើប” ។ ប្រសិនបើអ្នកជាអ្នកប្រើប្រាស់សហគ្រាសយើងនឹងការពារអាថ៌កំបាំងអាជីវកម្មរបស់អ្នកស្របតាមច្បាប់ហើយនឹងមិនផ្តល់ជូនអ្នកដោយគ្មានការយល់ព្រមពីអ្នកឡើយ។\n៨- បង្ខំឧត្ដមភាពការបដិសេធនិងដែនកំណត់នៃការទទួលខុសត្រូវ\n(១) ការបដិសេធ\nយើងមិនអាចផ្តល់សេវាកម្មជាធម្មតាដោយសារតែហេតុផលដូចខាងក្រោមហើយយើងត្រូវបានរួចផុតពីការទទួលខុសត្រូវ\nប្រព័ន្ធរបស់យើងត្រូវបានបិទសម្រាប់ការថែរក្សាឬធ្វើឱ្យប្រសើរឡើង។\nបង្ខំឱ្យមានហេតុផលធំ ៗ ដូចជាព្យុះទីហ្វុងរញ្ជួយដីទឹកជំនន់រន្ទះឬការវាយប្រហារភេរវកម្ម។\n3. ផ្នែករឹងនិងសូហ្វវែរកុំព្យូទ័ររបស់អ្នកបណ្តាញទំនាក់ទំនងនិងបណ្តាញផ្គត់ផ្គង់ថាមពលបរាជ័យ។\nដោយសារតែប្រតិបត្តិការមិនត្រឹមត្រូវរបស់អ្នកឬការប្រើប្រាស់សេវាកម្មរបស់យើងដោយមិនត្រូវបានអនុញ្ញាតឬយល់ព្រមពីយើង។\n៥- ដោយសារវីរុសវីរុស Trojan Trojan ការវាយប្រហារកម្មវិធីព្យាបាទការកកស្ទះបណ្តាញអស្ថិរភាពប្រព័ន្ធឬបរិក្ខារការបរាជ័យការប្រាស្រ័យទាក់ទងការបរាជ័យថាមពលមូលហេតុរបស់ធនាគារបញ្ហារបស់ភាគីទីបីឬសកម្មភាពរបស់រដ្ឋាភិបាល។ ល។ ទោះបីមានការព្រមព្រៀងក្នុងកថាខណ្ឌខាងមុខក៏ដោយយើងនឹងចាត់វិធានការសមហេតុផលដើម្បីលើកកម្ពស់សេវាកម្មឱ្យមានដំណើរការធម្មតា។\n(២) ការកំណត់នៃការទទួលខុសត្រូវ\nយើងអាចនឹងផ្តល់សេវាកម្មដល់អ្នកនិងប្រតិបត្តិការរបស់អ្នកក្នុងពេលតែមួយ។ គុណវិបត្តិ។\nការការពារកម្មសិទ្ធិបញ្ញា\nមាតិកាខ្លឹមសារនៃប្រព័ន្ធក្រុមហ៊ុននិងសម្ព័ន្ធភាពរបស់យើងនិងគេហទំព័ររបស់ Alipay រួមមានមិនកំណត់ត្រឹមការងាររូបភាពឯកសារទិន្នន័យរចនាសម្ពន្ធ័គេហទំព័រការរៀបចំអេក្រង់គេហទំព័រនិងរចនាគេហទំព័រជាកម្មសិទ្ធិរបស់យើងឬក្រុមហ៊ុនដែលពាក់ព័ន្ធស្របតាមច្បាប់ , រួមបញ្ចូលតែមិនមានកំណត់ចំពោះសិទ្ធិពាណិជ្ជសញ្ញា, សិទ្ធិប៉ាតង់, អាថ៌កំបាំងនៃការរក្សាសិទ្ធិ។ ល។\n២- កុំប្រើប្រាស់កែប្រែបំបែកឯកសារចម្លងចែកចាយជាសាធារណៈផ្លាស់ប្តូរចែកចាយចែកចាយឬផ្សព្វផ្សាយកម្មវិធីគេហទំព័ររបស់ Alipay ជាសាធារណៈឬដោយគ្មានការយល់ព្រមជាលាយលក្ខណ៍អក្សររបស់យើង\nខ្លឹមសារ។\n3. ការគោរពសិទ្ធិកម្មសិទ្ធិបញ្ញាគឺជាកាតព្វកិច្ចរបស់អ្នកប្រសិនបើមានការរំលោភបំពានអ្នកត្រូវទទួលខុសត្រូវចំពោះការខូចខាត។\nដប់, ផ្សេងទៀត\n១. សុពលភាពការបកស្រាយការកែប្រែការប្រតិបត្តិនិងការដោះស្រាយវិវាទនៃកិច្ចព្រមព្រៀងនេះគឺអនុវត្តទៅនឹងច្បាប់នៃព្រះរាជាណាចក្រកម្ពុជា។ វិវាទដែលកើតចេញពីកិច្ចព្រមព្រៀងនេះនឹងត្រូវដោះស្រាយស្របតាមច្បាប់នៃព្រះរាជាណាចក្រកម្ពុជា។\n២- ការបរាជ័យរបស់ប្រាក់បៀវត្សរ៍ទំនើបក្នុងការអនុវត្តឬអនុវត្តសិទ្ធិឬបទប្បញ្ញត្តិណាមួយនៃកិច្ចព្រមព្រៀងសេវាកម្មនេះមិនបង្កើតឱ្យមានការលះបង់សិទ្ធិឬសិទ្ធិដែលបានរៀបរាប់ខាងលើឡើយ។\n3. ប្រសិនបើមានប្រការណាមួយនៅក្នុងកិច្ចព្រមព្រៀងនេះមានលក្ខណៈមិនពេញលេញឬដោយផ្នែកខ្លះឬមិនអាចអនុវត្តបានសម្រាប់ហេតុផលណាមួយមាត្រាដែលនៅសល់នៃកិច្ចព្រមព្រៀងនេះនៅតែមានសុពលភាពនិងមានកាតព្វកិច្ច។	km	0	2020-04-08 18:55:01.305093	0	2020-04-09 17:33:05.909238	1
2020021313120082759558	Modern Pay隐私政策\n\nModern Pay是一款为用户提供电子钱包的支付产品，为说明我们收集那些关于您的数据、如何收集以及与谁共享这些数据，我们将通过本指引向您阐述相关事宜。\n如果您对本隐私政策未予说明的隐私保护措施有疑问，请联系我们。\n\n目录\n一、我们可能收集的个人信息范围及时间\n二、我们收集、使用信息的目的\n三、我们可能分享、转让或披露的信息\n四、我们如何使用Cookie、Proxy、数据埋点等技术\n五、我们如何存储和保护您的个人信息安全\n六、您如何查询、修改或删除个人信息\n七、未成年人的个人信息保护\n八、本隐私政策如何更新\n九、如何联系我们\n\n一、我们可能收集的个人信息范围及时间\n我们可能收集的信息包括：\n1、您提供的信息\n（1）您在注册账户或使用我们的服务时，向我们提供的相关个人信息，包括电话号码、电子邮件、身份证号码及指纹信息。\n（2）您通过我们的服务向其他方提供的共享信息，以及您使用我们的服务时所存储的信息。\n2、我们获取的您的信息\n为使Modern Pay服务更加贴切用户的需求，了解产品的适配性、识别账号的异常状态，我们会收集、汇总、记录的信息，包括日志信息、设备信息、位置信息、相册信息以及手机常用属性信息。\n（1）日志信息：当您使用Modern Pay提供的产品或服务时，我们会自动收集您对我们服务的详细使用情况，作为有关网络日志保存。包括您的搜索查询内容、IP地址、浏览器的类型、电信运营商、使用的语言、访问日期和时间、安装应用程序及您访问的网页记录。\n（2）设备信息：我们会根据您在软件安装及使用中授予的具体权限，接收并记录您所使用的设备相关信息（包括设备型号、操作系统版本、设备设置、唯一设备标识符软硬件特征信息）、设备所在位置相关信息（包括IP地址、GPS位置以及能够提供相关信息的WiFi接入点、蓝牙和基站传感器信息）。\n（3）手机常用属性信息：包括读取外部存储、写入外部存储、读取WiFi状态以及网络请求。\n3、我们通过间接获得方式收集的您的个人信息\n为便于我们基于关联账号共同向您提供一站式服务并便于您统一进行管理，您可通过Modern Pay账号在我们提供的链接入口使用Modern Pay及其旗下各关联公司提供的产品或服务。当您通过我们产品使用上述服务时，您授权我们根据实际业务及合作需要从我们的关联公司处接收、汇总、分析来源合法的或您授权同意其向我们提供的您的个人信息或交易信息。如您拒绝提供上述信息或拒绝授权，可能无法使用Modern Pay相应产品或服务，或者无法展示相关信息。我们通过间接的方式收集到的信息，亦包括其他方使用Modern Pay服务时所提供有关您的共享信息。\n4、我们只会在您使用特定业务功能时，仅收集为了正常运行该业务功能所必须的信息，在您停止该业务功能的使用后，我们将会停止您的个人信息收集的行为。\n二、我们收集、使用信息的目的\n我们仅出于以下目的，收集和使用您的个人信息：\n1、帮助您完成注册\n为便于我们为您提供服务，您需要提供基本注册信息，包括手机号码，并创建您的用户名和密码。在部分单项服务中，如果您只需要使用浏览、搜索等基本服务，您不需要注册成为Modern Pay用户及提供上述信息。\n2、向您提供商品或服务\n我们所收集和使用的信息是为您提供Modern Pay服务的必要条件，如缺少相关信息，我们将无法为您提供Modern Pay服务的核心内容，包括：\n（1）在您购买商品或服务时，为便于向您交付商品或服务。您需提供收货人个人身份信息、姓名、联系电话、支付状态信息。如果您拒绝提供此类信息，我们将无法完成相关交付服务。如您通过Modern Pay为其他人订购商品或服务，您需要提供该实际订购人的上述信息。向我们提供该实际订购人的上述信息之前，您需确保您已经取得其授权同意。\n（2）为展示您账户的订单信息，我们会收集您在使用Modern Pay服务过程中产生的订单信息用于向您展示及便于您对订单进行管理，本隐私政策中所述的订单信息包括商品名称、商品数量、商品价格、订单金额、配送信息（配送地址、收件人姓名、收件人手机号）、设备信息，若您需要开票，还会记录您的发票开票信息。\n（3）当您与我们联系时，我们可能会保存您的通信/通话记录和内容或您留下的联系方式信息，以便与您联系或帮助您解决问题，或记录相关问题的处理方案及结果。\n（4）为确认交易状态及为您提供售后与争议解决服务，我们会通过您基于交易所选择的交易对象、支付机构、物流公司收集与交易进度相关的您的交易、支付、物流信息，或将您的交易信息共享给上述服务提供者。\n（5）当您需要通过我们合作方申请信贷服务时，应信贷服务提供方要求，需要收集您的通话记录、通讯录及联系人信息。\n（6）当您需要开通扫描功能或上传您的账户头像时，需要收集您的拍摄及相册信息。\n（7）当您需要开通话费充值功能时，若您同意我们读取您的通讯录及联系人信息，将有助于我们给您带来更好的体验，但如果您不同意，也不会影响您开通本功能。\n（8）当您需要在购买特定理财产品时设置定时购买功能时，需要收集您的日历信息。\n（9）当您需要开通指纹解锁功能时，需要收集您的指纹信息。\n（10）当您使用我们海淘服务时，为进行结售汇录入及国际收支申报、需要收集您的订单信息、个人身份信息、姓名、银行账号及联系方式。\n（11）当您登陆时，为了反欺诈及风险防范，需要收集您的手机设备IMEI/IMSI/ICCID信息，用来判断您的账户是否存在风险。\n（12）当您登陆我们的APP时，或者您需要开通并使用生活缴费功能及服务时，您同意我们收集您的地理位置信息，但如果您不同意，也不会影响您开通前述相关功能及使用相关服务。\n3、向您推送消息\n（1）为您展示和推送商品或服务。我们可能使用您的信息，您的浏览及搜索记录、设备信息、位置信息、订单信息，提取您的浏览、搜索偏好、行为习惯、位置信息特征，并基于特征标签通过电子邮件、短信或其他方式向您发送营销信息，提供或推广我们或第三方的如下商品或服务。\n（2）向您发出通知。我们可能在必须时（例如当我们由于系统维护而暂停某一单项服务、变更、终止提供某一单项服务时）向您发出与服务有关的通知。\n如您不希望继续接收我们推送的消息，您可要求我们停止推送，例如：根据短信退订指引，要求我们停止发送推广短信等；但我们依法律规定或单项服务的服务协议约定发送消息的情形除外。\n4、为您提供安全保障\n为确保您身份真实性、向您提供更好的安全保障，您可以向我们提供身份证明、面部特征、声纹信息以完成实名认证及后续身份校验。\n我们可能将您的信息用于客户服务、安全防范、诈骗监测、存档和备份等用途，确保我们向您提供的服务的安全性；我们可能使用或整合我们所收集的您的信息，以及我们的合作伙伴取得您授权或依据法律共享的信息，来综合判断您的账户及交易风险、进行身份验证、检测及防范安全事件，并依法采取必要的记录、审计、分析、处置措施。\n5、其他用途\n我们将信息用于本《隐私政策》未载明的其他用途，或者将基于特定目的收集而来的信息用于其他目的时，会事先征求您的同意。\n您了解并同意，在收集您的信息后，我们将通过技术手段对数据进行去标识化处理，去标识化处理的信息将无法识别您的身份，在此情况下我们有权使用已经去标识化的信息，对用户数据库进行分析并予以商业化的利用。\n6、征得授权同意的例外\n根据相关法律法规规定，以下情形中收集您的信息无需征得您的授权同意：\n（1）与国家安全、国防安全有关的；\n（2）与公共安全、公共卫生、重大公共利益有关的；\n（3）与犯罪侦查、起诉、审判和判决执行等有关的；\n（4）出于维护信息主体或其他个人的生命、财产等重大合法权益但又很难得到您本人同意的；\n（5）所收集的信息是您自行向社会公众公开的；\n（6）从合法公开披露的信息中收集信息的，如合法的新闻报道、政府信息公开等渠道；\n（7）根据您的要求签订合同所必需的；\n（8）用于维护Modern Pay服务的安全稳定运行所必需的，例如发现、处置产品或服务的故障；\n（9）为合法的新闻报道所必需的；\n（10）学术研究机构基于公共利益开展统计或学术研究所必要，且对外提供学术研究或描述的结果时，对结果中所包含的信息进行去标识化处理的；\n（11）法律法规规定的其他情形。\n三、我们可能共享、转让或披露的信息\n1、共享\n除以下情形外，未经您同意，Modern Pay公司不会与任何第三方共享您的个人信息。\n我们仅会出于合法、正当、必要、特定、明确的目的共享您的信息。对我们与之共享信息的公司、组织和个人，我们会与其签署严格的保密协定，要求他们按照我们的说明、本《隐私政策》以及其他任何相关的保密和安全措施来处理信息。\n（1）向您提供我们的服务：我们可能向合作伙伴及其他第三方共享您的信息，以实现您需要的核心功能或提供您需要的服务，包括：向物流服务商提供对应的订单信息。\n（2）维护和改善我们的服务：我们可能向合作伙伴及其他第三方共享您的信息，以帮助我们为您提供更有针对性、更完善的服务，包括：代表我们发出电子邮件或推送通知的通讯服务提供商等。\n（3）实现本《隐私政策》“我们收集、使用信息的目的”部分所述的内容。\n（4）履行我们在本《隐私政策》或我们与您达成的其他协议中的义务和行使我们的权利。\n（5）向委托我们进行推广的合作伙伴等第三方共享，但我们仅会向这些委托方提供推广的覆盖面和有效性的信息，而不会提供可以识别您身份的信息，包括姓名、电话号码或电子邮箱；或者我们将这些信息进行汇总，以便它不会识别您个人。比如我们可以告知该委托方有多少人看了他们的推广信息或在看到这些信息后购买了委托方的商品，或者向他们提供不能识别个人身份的统计信息，帮助他们了解其受众或顾客。\n（6）在法律法规允许的范围内，为了遵守法律、维护我们级我们的关联方或合作伙伴、您或其他Modern Pay用户或社会公众利益、财产或安全免遭损害，比如为防止欺诈等违法活动和减少信用风险，我们可能与其他公司和组织交换信息。不过，这并不包括违反本《隐私政策》中所作的承诺而为获利目的出售、出租、共享或以其它方式披露的信息。\n（7）应您合法需求，协助处理您与他人的纠纷或争议。\n（8）应您的监护人合法要求而提供您的信息。\n（9）根据与您签署的单项服务协议（包括在线签署的电子协议以及相应的平台规则）或其他的法律文件约定所提供。\n（10）基于学术研究而提供。\n（11）基于符合法律法规的社会公共利益而提供。\n2、转让\n（1）随着我们业务的持续发展，我们有可能进行合并、收购、资产转让或类似的交易，而您的信息有可能作为此类交易的一部分而被转移。我们会要求新的持有您信息的公司、组织继续受本《隐私政策》的约束，否则，我们将要求该公司、组织重新向您征求授权同意。\n（2）在获得您的明确同意后，我们会向其他方转让您的信息。\n3、披露\n我们仅会在以下情况下，且采取符合业界标准的安全防护措施的前提下，才会披露您的信息：\n（1）根据您的需求，在您明确同意的披露方式下披露您所指定的信息。\n（2）根据法律法规的要求、强制性的行政执法或司法要求所必须提供您信息的情况下，我们可能会依据所要求的信息类型和披露方式披露您的信息。在符合法律法规发前提下，当我们收到上述披露信息的请求时，我们会要求接收方必须出具与之相应的法律文件，如传票或调查函。我们坚信，对于要求我们提供的信息，应该在法律允许的范围内尽可能保持透明。我们对所有的请求都进行了慎重的审查，以确保其具备合法依据，且仅限于有权部门因特定调查目的且有合法权利获取的数据。\n四、我们如何使用Cookie、Proxy、数据埋点等技术\n1、为使您获得更轻松的访问体验，您访问Modern Pay提供的产品或服务时，我们可能会通过小型数据文件识别您的身份，这么做可帮您省去重复输入注册信息的步骤，或者帮助判断您的账户安全。这些数据文件可能是Cookie，Flash Cookie，或您的浏览器或关联应用程序提供的其他本地存储（统称“Cookie”）。请您理解，我们的某些服务只能通过使用Cookie才可得到实现。如果您的浏览器或浏览器附加服务允许，您可以修改对Cookie的接受程度或者拒绝Modern Pay的Cookie。多数浏览器上工具条中的帮助部分会告诉您怎样防止您的浏览器接受新的Cookies，怎样让您的浏览器在您收到一条新Cookies，怎样让您的浏览器在您收到一条新Cookies时通知您或者怎样彻底关闭Cookies。此外，您可以通过改变浏览器附加程序的设置，或通过访问提供商的网页，来关闭或删除浏览器附加程序使用的类似数据（诸如Flash Cookies)。但这一举动在某些情况下可能会影响您安全访问Modern Pay和使用Modern Pay提供的服务。\n2、我们网站上还可能包含一些数据埋点，我们会通过埋点收集您浏览网页活动的信息，包括您访问的页面地址、您先前访问的援引页面的位置、您的浏览环境以及显示设定。\n3、如您通过Modern Pay平台使用了由第三方而非Modern Pay提供的服务时，我们无法保证这些第三方主体会按照我们的要求采取保护措施，为尽力确保您的账号安全，使您获得更安全的访问体验，我们可能会使用专用的网络协议及代理技术(称为“专用网络通道”、“网络代理”)。使用专用网络通道，可以帮助您识别到我们已知的高风险站点，减少由此引起的钓鱼、账号泄露等风险，同时更有利于保障您和第三方的共同权益，阻止不法分子篡改您和您希望访问的第三方之间正常服务内容，例如不安全路由器、非法基站等引起的广告注入，、非法内容篡改等。在此过程中，我们也可能会获得和保存关于您计算机的相关信息，比如IP地址、硬件ID。\n五、我们如何存储和保护您的个人信息安全\n1、我们仅在本政策所述目的所必需期间和法律法规要求的时限内保留您的个人信息。在您未进行信息撤回、信息删除或账号注销等行为时,我们会保留相关信息,这是提供业务所必要的。当您进行信息撤回\n信息删除或账号注销后,我们将对您的个人信息进行立即删除或匿名化处理,例外情况:请您理解,我们还暠满足国家相关法律法规的要求,在法律法规及监管规定的时限内,必须保存规定的您的个人信息。\n2、我们承诺我们将使信息安全保护达到合理的安全水平。为保障您的信息安全,我们致力于使用各种安全技术及配套的管理体系来防止您的信息被泄露、毁损或者丢失,如通过网络安全层软件(SSL)进行加密传输、信息加密存储、严格限制数据中心的访问、使用专用网络通道及网络代理。同时我们设立了信息安全保护责任部门,建立了相关内控制度,对可能接触到您的信息的工作人员采取最小化授权原则;对工作人员处理您的信息的行为进行系统监控,不断对工作人员培训相关法律法规及隐私安全准则和安全意识强化宣导。另外,我们每年还会聘请外部独立第三方对我们的信息安全管理体系进行评估。\n3、请您理解,由于技术水平限制及可能存在的各种恶意手段,有可能因我们可控范围外的因素而出现问题。在不幸发生个人信息安全事件后,我们将按照法律法规的要求,及时向您告知:安全事件的基本情况和可能的影响、我们已采取或将要采取的处置措施您可自主防范和降低风险的建议、对您的补救措施等。我们将及时将事件相关情况以邮件、信函、电话或推送通知方式告知您,难以逐一告知个人信息主体时,我们会采取合理、有效的方式发布公告。同时我们还将按照监管部门要求,主动上报个人信息安全事件的处置情况。\n4、如我们停止运营,我们将及时停止继续收集您个人信息的活动,将停止运营的通知以逐一送达或公告的形式通知您,对所持有的个人信息进行删除或匿名化处理。\n六、您如何查询、修改或删除个人信息\n1、您可在Modern Pay客户端登录Modern Pay账户在“我的-个人信息”中,查阅您的身份信息、账户信息,或修改您的个人资料,或进行相关的隐私、安全设置;但出于安全性和身份识别的考虑或根据法律法规的强制规定,您可能无法修改注册时提供的初始注册信息。\n2、您可在Modern Pay客户端登录Modern Pay账户后,在“我的账单中删除您的某笔交易信息或在“我的-银行卡中管理绑定的银行卡。\n3、若您想要撤回某项授权,您可通过自助手机操作系统的设置界面关闭对应权限。\n4、您可通过如下方式注销Modern Pay账户:若您需要注销账户,可以联系我们的客服电话4008880999或通过“我的-我的客服在线客服”,联系人工客服协助处理,客服将在验证通过后三个工作日内为您完成账户注销。注销Modern Pay账户前,请您确认Modern Pay账户中资产、账户及交易状态,Modern Pay账户旦注销无法恢复。\n注销Modern Pay账户后,您该账户内的所有信息将被清空,我们将不会再收集、使用或共享与该账户相关的个人信息,但之前的信息如有法律法规要求,我们将按要求进行保存,且在该依法保存的时间内有权机关仍有权依法查询。\n5、如您发现我们采集、储存的您的个人信息有错误的或对采集的信息内容或方式持有疑异的,您可拨打我们的客服电话4008880999或通过在线客服要求删除或更正,我们将在15天内回复处理意见或结果。\n6、在响应您上述4、5项要求前,为保障安全,我们需要先验证您的身份和凭证资料,验证通过后将在三个工作日内答复。\n7、尽管有上述约定,但按照法律法规要求,在以下情形下,我们可能无法响应您的请求\n(1)与国家安全、国防安全有关的;\n(2)与公共安全、公共卫生、重大公共利益有关的;\n(3)与犯罪侦查、起诉和审判等有关的;\n(4)有充分证据表明您存在主观恶意或滥用权利的；\n(5)响应您的请求将导致其他个人、组织的合法权益受到严重损害的。\n七、未成年人的个人信息保护\n我们期望父母或监护人指导未成年人使用我们的服务。如您为未成年人,建议请您的父母或监护人阅读本政策,并在征得您父母或监护人同意的前提下使用我们的服务或向我们提供您的信息。如您的监护人不\n同意您按照本政策使用我们的服务或向我们提供信息,请您立即终止使用我们的服务并及时通知我们,以便我们采取相应的措施。Modern Pay将根据国家相关法律法规的规定保护未成年人的个人信息的保密性及安全性。\n八、本隐私权政策如何更新\n基于业务功能、使用规则、联络方式、或法律法规及监管要求,我们可能会适时对本政策进行修订,该等修改构成本《隐私政策》的一部分。对于重大变更,我们会提供更显著的通知,您可以选择停止使用Modern Pay服务;若您点击确认经修订的《隐私政策》即表示同意受经修订的本《隐私政策》的约束。\n您可以在Modern Pay官方网站服务中心”查看本政策。我们鼓励您在每次访问Modern Pay时都查阅我们的隐私政策。\n九、如何联系我们\nModern Pay专门设立了个人信息保护岗,该岗位主要负责Modern Pay用户个人信息保护及解答工作。如您有关于网络信息安全的投诉和举报,或您对本《隐私政策》、您的信息的相关事宜有任何问题、意见或建议,可通过以下方式联系我们：\n1、请拨打我们的客服电话4008880999,我们将会由客服人员进行记录并通知个人信息保护岗进行处理；\n2、通过在线客服进行咨询;\n3、以书面形式寄送至以下地址:\n\n本版更新日期:2020年01月08日\n本版生效日期:2020年01月08日	zh_CN	1	2020-02-13 13:10:12.358959	0	2020-04-09 17:33:31.419168	1
2020040917342380424242	គោលការណ៍ភាពឯកជននៃប្រាក់ខែទំនើប\n\nសេវាកម្មបង់ប្រាក់ទំនើបគឺជាផលិតផលទូទាត់ដែលផ្តល់ឱ្យអ្នកប្រើប្រាស់នូវកាបូបអេឡិចត្រូនិចដើម្បីពន្យល់ពីរបៀបដែលយើងប្រមូលទិន្នន័យអំពីអ្នករបៀបប្រមូលវាហើយជាមួយអ្នកណាយើងនឹងពន្យល់អំពីបញ្ហាដែលទាក់ទងនឹងអ្នកតាមរយៈសៀវភៅណែនាំនេះ។\nប្រសិនបើអ្នកមានសំណួរអំពីវិធានការការពារភាពឯកជនដែលមិនបានពិពណ៌នានៅក្នុងគោលការណ៍ឯកជននេះសូមទាក់ទងមកយើងខ្ញុំ។\n\nបញ្ជីឈ្មោះ\nជួរនិងពេលវេលានៃព័ត៌មានផ្ទាល់ខ្លួនដែលយើងអាចប្រមូលបាន\nទីពីរគោលបំណងនៃការប្រមូលនិងការប្រើប្រាស់ព័ត៌មានរបស់យើង\nព័ត៌មានដែលយើងអាចចែករំលែកផ្ទេរឬផ្សព្វផ្សាយ\nតើយើងប្រើខូឃីស៍ប្រូកស៊ីការបញ្ចូលទិន្នន័យនិងបច្ចេកវិទ្យាផ្សេងទៀតយ៉ាងដូចម្តេច\nតើយើងរក្សាទុកនិងការពារព័ត៌មានផ្ទាល់ខ្លួនរបស់អ្នកយ៉ាងដូចម្តេច?\n6. តើអ្នកសួរកែប្រែឬលុបព័ត៌មានផ្ទាល់ខ្លួនយ៉ាងដូចម្តេច\nការការពារព័ត៌មានផ្ទាល់ខ្លួនរបស់អនីតិជន\nវិធីធ្វើបច្ចុប្បន្នភាពគោលនយោបាយសិទ្ធិឯកជននេះ\nប្រាំបួន, របៀបទាក់ទងមកយើង\n\nជួរនិងពេលវេលានៃព័ត៌មានផ្ទាល់ខ្លួនដែលយើងអាចប្រមូលបាន\nព័ត៌មានដែលយើងអាចប្រមូលបានរួមមាន៖\n1. ព័ត៌មានដែលអ្នកផ្តល់\n(១) នៅពេលអ្នកចុះឈ្មោះគណនីឬប្រើប្រាស់សេវាកម្មរបស់យើងអ្នកផ្តល់ឱ្យយើងនូវព័ត៌មានផ្ទាល់ខ្លួនដែលពាក់ព័ន្ធរួមមានលេខទូរស័ព្ទអ៊ីមែលលេខអត្តសញ្ញាណប័ណ្ណនិងព័ត៌មានស្នាមម្រាមដៃ។\n(២) ព័ត៌មានដែលអ្នកចែករំលែកទៅឱ្យភាគីផ្សេងទៀតតាមរយៈសេវាកម្មរបស់យើងនិងព័ត៌មានដែលរក្សាទុកនៅពេលអ្នកប្រើសេវាកម្មរបស់យើង។\nព័ត៌មានរបស់អ្នកដែលយើងទទួលបាន\nដើម្បីធ្វើឱ្យសេវាកម្មទំនើបទាន់សម័យកាន់តែសមស្របទៅនឹងតំរូវការរបស់អ្នកប្រើប្រាស់ស្វែងយល់ពីភាពស៊ីគ្នានៃផលិតផលនិងកំណត់ពីស្ថានភាពមិនប្រក្រតីនៃគណនីយើងនឹងប្រមូលប្រមូលនិងកត់ត្រាព័ត៌មានរួមមានព័ត៌មានកំណត់ហេតុព័ត៌មានឧបករណ៍ទីតាំងព័ត៌មានអាល់ប៊ុមនិងគុណលក្ខណៈទូទៅនៃទូរស័ព្ទចល័ត។ ព័ត៌មាន។\n(១) ព័ត៌មានកំណត់ហេតុ៖ នៅពេលដែលអ្នកប្រើផលិតផលឬសេវាកម្មដែលផ្តល់ដោយ Modern Pay យើងនឹងប្រមូលការប្រើប្រាស់ព័ត៌មានលម្អិតរបស់អ្នកដោយស្វ័យប្រវត្តិហើយរក្សាទុកវាជាកំណត់ហេតុគេហទំព័រពាក់ព័ន្ធ។ រួមបញ្ចូលទាំងមាតិកាសំណួរស្វែងរករបស់អ្នកអាសយដ្ឋាន IP ប្រភេទកម្មវិធីរុករកប្រតិបត្តិករទូរគមនាគមន៍ភាសាដែលបានប្រើកាលបរិច្ឆេទនិងពេលវេលាចូលប្រើតំឡើងកម្មវិធីនិងកំណត់ត្រាគេហទំព័រដែលអ្នកបានចូលមើល។\n(២) ព័ត៌មានអំពីឧបករណ៍ៈយើងនឹងទទួលនិងកត់ត្រាព័ត៌មានអំពីឧបករណ៍ដែលអ្នកប្រើ (រាប់បញ្ចូលទាំងម៉ូដែលឧបករណ៍កំណែប្រព័ន្ធប្រតិបត្តិការការកំណត់ឧបករណ៍កម្មវិធីកំណត់អត្តសញ្ញាណឧបករណ៍និងលក្ខណៈពិសេសផ្នែករឹងតែមួយគត់យោងទៅតាមការអនុញ្ញាតជាក់លាក់ដែលអ្នកផ្តល់ក្នុងអំឡុងពេលដំឡើងនិងប្រើប្រាស់) ព័ត៌មាន) ព័ត៌មានអំពីទីតាំងរបស់ឧបករណ៍ (រាប់បញ្ចូលទាំងអាសយដ្ឋាន IP ទីតាំង GPS និងព័ត៌មានអំពីចំណុចចូលប្រើវ៉ាយហ្វាយប៊្លូធូសនិងឧបករណ៍ចាប់សញ្ញាស្ថានីយ៍មូលដ្ឋានដែលអាចផ្តល់ព័ត៌មានពាក់ព័ន្ធ) ។\n(៣) ព័ត៌មានគុណលក្ខណៈទូទៅនៃទូរស័ព្ទចល័ត៖ រួមមានការអានការផ្ទុកខាងក្រៅសរសេរទៅឧបករណ៍ផ្ទុកខាងក្រៅអានស្ថានភាពវ៉ាយហ្វាយនិងសំណើបណ្តាញ។\nព័ត៌មានផ្ទាល់ខ្លួនរបស់អ្នកដែលយើងប្រមូលតាមរយៈការទិញដោយប្រយោល\nដើម្បីជួយសម្រួលដល់ពួកយើងក្នុងការផ្តល់ឱ្យអ្នកនូវសេវាកម្មច្រកចេញចូលតែមួយដោយផ្អែកលើគណនីដែលពាក់ព័ន្ធនិងដើម្បីសម្រួលដល់ការគ្រប់គ្រងរួមរបស់អ្នកអ្នកអាចប្រើផលិតផលឬសេវាកម្មដែលផ្តល់ជូនដោយសេវាកម្មបង់ប្រាក់ទំនើបនិងក្រុមហ៊ុនដែលពាក់ព័ន្ធតាមរយៈគណនីបង់ប្រាក់ទំនើបនៅច្រកចូលភ្ជាប់ដែលផ្តល់ដោយយើង។ នៅពេលដែលអ្នកប្រើសេវាកម្មខាងលើតាមរយៈផលិតផលរបស់យើងអ្នកផ្តល់សិទ្ធិឱ្យយើងទទួលបានសង្ខេបនិងវិភាគប្រភពនៃព័ត៌មានផ្ទាល់ខ្លួនឬព័ត៌មានអំពីប្រតិបត្តិការពីក្រុមហ៊ុនដែលពាក់ព័ន្ធរបស់យើងស្របតាមតម្រូវការអាជីវកម្មនិងកិច្ចសហប្រតិបត្តិការជាក់ស្តែងឬអ្នកផ្តល់សិទ្ធិនិងយល់ព្រមផ្តល់ជូនយើង។ ប្រសិនបើអ្នកបដិសេធមិនផ្តល់ព័ត៌មានខាងលើឬបដិសេធមិនផ្តល់សិទ្ធិអ្នកមិនអាចប្រើផលិតផលឬសេវាកម្មដែលត្រូវគ្នានៃប្រាក់ខែទំនើបឬអ្នកមិនអាចបង្ហាញព័ត៌មានពាក់ព័ន្ធបានទេ។ ព័ត៌មានដែលយើងប្រមូលដោយប្រយោលក៏រួមបញ្ចូលព័ត៌មានដែលបានចែករំលែកអំពីអ្នកផងដែរនៅពេលដែលភាគីផ្សេងទៀតប្រើប្រាស់សេវាកម្មបង់ប្រាក់ទំនើប។\n៤. យើងនឹងប្រមូលតែព័ត៌មានចាំបាច់សម្រាប់ដំណើរការធម្មតានៃមុខងារអាជីវកម្មនៅពេលដែលអ្នកប្រើមុខងារអាជីវកម្មជាក់លាក់មួយ។ បន្ទាប់ពីអ្នកឈប់ប្រើប្រាស់មុខងារអាជីវកម្មយើងនឹងបញ្ឈប់សកម្មភាពប្រមូលព័ត៌មានផ្ទាល់ខ្លួនរបស់អ្នក។\nទីពីរគោលបំណងនៃការប្រមូលនិងការប្រើប្រាស់ព័ត៌មានរបស់យើង\nយើងប្រមូលនិងប្រើប្រាស់ព័ត៌មានផ្ទាល់ខ្លួនរបស់អ្នកសម្រាប់គោលបំណងដូចខាងក្រោម៖\nជួយអ្នកឱ្យបំពេញការចុះឈ្មោះ\nដើម្បីឱ្យយើងផ្តល់ជូនអ្នកនូវសេវាកម្មអ្នកត្រូវផ្តល់ព័ត៌មានចុះឈ្មោះមូលដ្ឋានរួមទាំងលេខទូរស័ព្ទនិងបង្កើតឈ្មោះអ្នកប្រើនិងពាក្យសម្ងាត់របស់អ្នក។ នៅក្នុងសេវាកម្មតែមួយប្រសិនបើអ្នកគ្រាន់តែត្រូវការប្រើសេវាកម្មមូលដ្ឋានដូចជាការរុករកនិងស្វែងរកអ្នកមិនចាំបាច់ចុះឈ្មោះជាអ្នកប្រើប្រាក់ខែទំនើបនិងផ្តល់ព័ត៌មានខាងលើទេ។\nផ្តល់ជូនអ្នកនូវទំនិញឬសេវាកម្ម\nព័ត៌មានដែលយើងប្រមូលនិងប្រើប្រាស់គឺជាលក្ខខណ្ឌចាំបាច់សម្រាប់ផ្តល់ជូនលោកអ្នកនូវសេវាកម្មបង់ប្រាក់ទំនើប។ ប្រសិនបើព័ត៌មានដែលមានការខ្វះខាតយើងនឹងមិនអាចផ្តល់ជូនលោកអ្នកនូវខ្លឹមសារសំខាន់នៃសេវាកម្មបង់ប្រាក់សម័យទំនើបរួមមាន៖\n(១) នៅពេលអ្នកទិញទំនិញឬសេវាកម្មដើម្បីសម្រួលដល់ការដឹកជញ្ជូនទំនិញឬសេវាកម្មដល់អ្នក។ អ្នកត្រូវផ្តល់ព័ត៌មានអត្តសញ្ញាណផ្ទាល់ខ្លួនរបស់អ្នកទទួលការយល់ព្រមឈ្មោះលេខទំនាក់ទំនងនិងព័ត៌មានស្ថានភាពការបង់ប្រាក់។ ប្រសិនបើអ្នកបដិសេធមិនផ្តល់ព័ត៌មានបែបនេះយើងនឹងមិនអាចបំពេញសេវាកម្មដឹកជញ្ជូនដែលពាក់ព័ន្ធបានទេ។ ប្រសិនបើអ្នកបញ្ជាទិញទំនិញឬសេវាកម្មសម្រាប់មនុស្សផ្សេងទៀតតាមរយៈ Modern Pay អ្នកត្រូវផ្តល់ព័ត៌មានខាងលើរបស់អ្នកកំណើតពិតប្រាកដ។ មុនពេលផ្តល់ឱ្យយើងនូវព័ត៌មានខាងលើនៃអ្នកបង្កើតពិតប្រាកដអ្នកត្រូវធានាថាអ្នកបានទទួលការយល់ព្រមដែលមានការអនុញ្ញាតពីពួកគេ។\n(២) ដើម្បីបង្ហាញព័ត៌មានបញ្ជាទិញនៃគណនីរបស់អ្នកយើងនឹងប្រមូលព័ត៌មានបញ្ជាទិញដែលបានបង្កើតក្នុងកំឡុងពេលប្រើប្រាស់សេវាកម្មទំនើបបង់ប្រាក់ដើម្បីបង្ហាញអ្នកនិងជួយសម្រួលដល់អ្នកក្នុងការគ្រប់គ្រងបញ្ជាទិញ។ ព័ត៌មានបញ្ជាទិញដែលបានពិពណ៌នានៅក្នុងគោលការណ៍ឯកជននេះរួមមានទំនិញ ឈ្មោះចំនួនផលិតផលតម្លៃផលិតផលចំនួនបញ្ជាទិញព័ត៌មានចែកចាយ (អាស័យដ្ឋានចែកចាយឈ្មោះអ្នកទទួលលេខទូរស័ព្ទអ្នកទទួល) ព័ត៌មានឧបករណ៍ប្រសិនបើអ្នកត្រូវការវិក្កយបត្រក៏នឹងកត់ត្រាព័ត៌មានវិក្កយបត្ររបស់អ្នកដែរ។\n(៣) នៅពេលអ្នកទាក់ទងមកយើងខ្ញុំយើងអាចរក្សាទុកការទំនាក់ទំនង / កំណត់ត្រាការហៅនិងមាតិការបស់អ្នកឬព័ត៌មានទំនាក់ទំនងដែលអ្នកទុកដើម្បីទាក់ទងអ្នកឬជួយអ្នកក្នុងការដោះស្រាយបញ្ហាឬកត់ត្រាផែនការព្យាបាលនិងលទ្ធផលនៃបញ្ហាដែលទាក់ទង។ ។\n(៤) ដើម្បីបញ្ជាក់ពីស្ថានភាពប្រតិបត្តិការនិងផ្តល់ជូនអ្នកនូវសេវាកម្មដោះស្រាយវិវាទក្រោយពេលលក់និងដោះស្រាយវិវាទយើងនឹងប្រមូលព័ត៌មានប្រតិបត្តិការទូទាត់ប្រាក់និងព័ត៌មានភ័ស្តុភារទាក់ទងនឹងវឌ្ឍនភាពនៃប្រតិបត្តិការតាមរយៈវត្ថុប្រតិបត្តិការស្ថាប័នទូទាត់ប្រាក់និងក្រុមហ៊ុនភ័ស្តុភារដែលអ្នកជ្រើសរើសផ្អែកលើការផ្លាស់ប្តូររឺ ចែករំលែកព័ត៌មានប្រតិបត្តិការរបស់អ្នកជាមួយអ្នកផ្តល់សេវាខាងលើ។\n(៥) នៅពេលដែលអ្នកត្រូវការដាក់ពាក្យសុំសេវាកម្មឥណទានតាមរយៈដៃគូរបស់យើងអ្នកត្រូវប្រមូលកំណត់ត្រាហៅចូលសៀវភៅអាសយដ្ឋាននិងព័ត៌មានទំនាក់ទំនងតាមការស្នើសុំរបស់អ្នកផ្តល់សេវាកម្មឥណទាន។\n(៦) នៅពេលដែលអ្នកត្រូវការបើកមុខងារស្កេនឬបញ្ចូលរូបភាពគណនីរបស់អ្នកអ្នកត្រូវប្រមូលព័ត៌មានថតនិងអាល់ប៊ុមរបស់អ្នក។\n(៧) នៅពេលដែលអ្នកត្រូវការបើកមុខងារបញ្ចូលទឹកប្រាក់នៃការហៅទូរស័ព្ទប្រសិនបើអ្នកយល់ព្រមឱ្យយើងអានសៀវភៅអាសយដ្ឋានរបស់អ្នកនិងព័ត៌មានទំនាក់ទំនងវានឹងជួយឱ្យយើងនាំមកនូវបទពិសោធន៍ប្រសើរជាងមុនប៉ុន្តែប្រសិនបើអ្នកមិនយល់ព្រមវានឹងមិន ប៉ះពាល់ដល់ការធ្វើឱ្យមុខងាររបស់អ្នកសកម្ម។\n(៨) នៅពេលអ្នកត្រូវការកំណត់មុខងារទិញទៀងទាត់នៅពេលទិញផលិតផលហិរញ្ញវត្ថុជាក់លាក់អ្នកត្រូវប្រមូលព័ត៌មានប្រតិទិនរបស់អ្នក។\n(៩) នៅពេលអ្នកត្រូវការបើកមុខងារដោះសោរស្នាមម្រាមដៃអ្នកត្រូវប្រមូលព័ត៌មានស្នាមម្រាមដៃរបស់អ្នក។\n(១០) នៅពេលអ្នកប្រើសេវាកម្មហៃថៅរបស់យើងអ្នកត្រូវប្រមូលព័ត៌មានបញ្ជាទិញរបស់អ្នកព័ត៌មានអត្តសញ្ញាណផ្ទាល់ខ្លួនឈ្មោះលេខគណនីធនាគារនិងព័ត៌មានទំនាក់ទំនងសម្រាប់ការទូទាត់ប្តូរប្រាក់បរទេសនិងតុល្យភាពនៃការប្រកាសចំណាយ។\n(១១) នៅពេលអ្នកចូលដើម្បីការពារការក្លែងបន្លំនិងការការពារហានិភ័យអ្នកត្រូវប្រមូលព័ត៌មាន IMEI / IMSI / ICCID នៃឧបករណ៍ចល័តរបស់អ្នកដើម្បីកំណត់ថាតើគណនីរបស់អ្នកមានគ្រោះថ្នាក់ឬអត់។\n(១២) នៅពេលដែលអ្នកចូលទៅក្នុង APP របស់យើងឬនៅពេលដែលអ្នកត្រូវការធ្វើឱ្យសកម្មនិងប្រើប្រាស់មុខងារនិងសេវាកម្មបង់ប្រាក់ជីវិតអ្នកយល់ព្រមថាយើងប្រមូលព័ត៌មានអំពីទីតាំងភូមិសាស្ត្ររបស់អ្នកប៉ុន្តែប្រសិនបើអ្នកមិនយល់ព្រមវានឹងមិនប៉ះពាល់ដល់ការធ្វើឱ្យមុខងាររបស់អ្នកទាក់ទងទៅនឹងអ្វីដែលបានរៀបរាប់ខាងលើទេ។ ប្រើសេវាកម្មដែលពាក់ព័ន្ធ។\n3. ជំរុញសារទៅអ្នក\n(១) បង្ហាញនិងជំរុញទំនិញឬសេវាកម្មសម្រាប់អ្នក។ យើងអាចប្រើព័ត៌មានការស្វែងរកនិងកំណត់ត្រាស្វែងរករបស់អ្នកព័ត៌មានឧបករណ៍ព័ត៌មានទីតាំងព័ត៌មានបញ្ជាទិញទាញយកការស្វែងរកចំណូលចិត្តឥរិយាបថលក្ខណៈព័ត៌មានទីតាំងនិងអ៊ីមែលសារខ្លីឬវិធីផ្សេងទៀតផ្អែកលើស្លាកលក្ខណៈពិសេស ផ្ញើព័ត៌មានទីផ្សារទៅអ្នកដើម្បីផ្តល់ឬផ្សព្វផ្សាយទំនិញឬសេវាកម្មដូចខាងក្រោមពីយើងឬភាគីទីបី។\n(២) ជូនដំណឹងដល់អ្នក។ យើងអាចនឹងផ្ញើការជូនដំណឹងទាក់ទងនឹងសេវាកម្មនៅពេលចាំបាច់ (ដូចជានៅពេលដែលយើងផ្អាកសេវាកម្មតែមួយផ្លាស់ប្តូរឬបញ្ឈប់សេវាកម្មតែមួយដោយសារការថែទាំប្រព័ន្ធ) ។\nប្រសិនបើអ្នកមិនចង់បន្តទទួលសារពីយើងអ្នកអាចស្នើសុំឱ្យយើងបញ្ឈប់ការរុញច្រានឧទាហរណ៍ៈយោងតាមគោលការណ៍ណែនាំការឈប់ជាវសារជាអក្សរយើងត្រូវបានស្នើសុំឱ្យបញ្ឈប់ការផ្ញើសារផ្សព្វផ្សាយ។ ល។ ប៉ុន្តែយើងផ្ញើសារតាមបទប្បញ្ញត្តិច្បាប់ឬកិច្ចព្រមព្រៀងសេវាកម្មរបស់សេវាកម្មតែមួយ។ លើកលែង។\nផ្តល់ជូនអ្នកនូវសុវត្ថិភាព\nដើម្បីធានាបាននូវភាពត្រឹមត្រូវនៃអត្តសញ្ញាណរបស់អ្នកនិងផ្តល់ឱ្យអ្នកនូវសុវត្ថិភាពកាន់តែប្រសើរអ្នកអាចផ្តល់ឱ្យយើងនូវអត្តសញ្ញាណអត្តសញ្ញាណមុខមាត់និងព័ត៌មានសំលេងដើម្បីបំពេញការផ្ទៀងផ្ទាត់ឈ្មោះពិតនិងការបញ្ជាក់អត្តសញ្ញាណជាបន្តបន្ទាប់។\nយើងអាចប្រើព័ត៌មានរបស់អ្នកសម្រាប់សេវាកម្មអតិថិជនសន្តិសុខការត្រួតពិនិត្យការបន្លំការរក្សាទុកនិងគោលបំណងបម្រុងទុកដើម្បីធានាសុវត្ថិភាពសេវាកម្មដែលយើងផ្តល់ជូនអ្នកយើងអាចប្រើឬបញ្ចូលព័ត៌មានរបស់អ្នកដែលប្រមូលបានដោយពួកយើងនិងរបស់យើង ដៃគូទទួលបានព័ត៌មានដែលអ្នកផ្តល់សិទ្ធិឬចែករំលែកស្របតាមច្បាប់ដើម្បីវិនិច្ឆ័យហានិភ័យគណនីនិងប្រតិបត្តិការរបស់អ្នកឱ្យបានទូលំទូលាយអនុវត្តការផ្ទៀងផ្ទាត់អត្តសញ្ញាណរកឃើញនិងការពារឧប្បត្តិហេតុសន្តិសុខហើយយកកំណត់ត្រាចាំបាច់សវនកម្មការវិភាគនិងវិធានការបោះចោលស្របតាមច្បាប់។\n5. ការប្រើប្រាស់ផ្សេងទៀត\nនៅពេលដែលយើងប្រើព័ត៌មានសម្រាប់គោលបំណងផ្សេងទៀតដែលមិនបានបញ្ជាក់នៅក្នុង“ គោលការណ៍ភាពឯកជន” នេះឬនៅពេលដែលព័ត៌មានដែលប្រមូលបានសម្រាប់គោលបំណងជាក់លាក់ត្រូវបានប្រើសម្រាប់គោលបំណងផ្សេងទៀតយើងនឹងស្វែងរកការព្រមព្រៀងពីអ្នកជាមុន។\nអ្នកយល់និងយល់ព្រមថាបន្ទាប់ពីប្រមូលព័ត៌មានរបស់អ្នកហើយយើងនឹងដកហូតទិន្នន័យតាមរយៈមធ្យោបាយបច្ចេកទេសហើយព័ត៌មានដែលត្រូវបានកំណត់អត្តសញ្ញាណនឹងមិនអាចកំណត់អត្តសញ្ញាណអ្នកបានទេក្នុងករណីនេះយើងមានសិទ្ធិប្រើអត្តសញ្ញាណដែលបានកំណត់ ព័ត៌មានវិភាគនិងធ្វើពាណិជ្ជកម្មមូលដ្ឋានទិន្នន័យអ្នកប្រើប្រាស់។\nការលើកលែងសម្រាប់ការទទួលបានការយល់ព្រមដែលមានការអនុញ្ញាត\nយោងទៅតាមច្បាប់និងបទប្បញ្ញត្តិពាក់ព័ន្ធការប្រមូលព័ត៌មានរបស់អ្នកក្នុងស្ថានភាពដូចខាងក្រោមមិនត្រូវការការអនុញ្ញាតពីអ្នកទេ៖\n(១) ពាក់ព័ន្ធនឹងសន្តិសុខជាតិនិងសន្តិសុខការពារជាតិ;\n(២) ទាក់ទងនឹងសុវត្ថិភាពសាធារណៈសុខភាពសាធារណៈនិងផលប្រយោជន៍សាធារណៈសំខាន់ៗ។\n(៣) ពាក់ព័ន្ធនឹងការស៊ើបអង្កេតព្រហ្មទណ្ឌការកាត់ទោសការជំនុំជម្រះក្តីនិងការអនុវត្តន៍សាលក្រម។ ល។\n(៤) ក្រៅពីការការពារសិទ្ធិនិងផលប្រយោជន៍ស្របច្បាប់សំខាន់ៗរបស់អង្គភាពព័ត៌មានឬបុគ្គលដទៃទៀតប៉ុន្តែវាពិបាកក្នុងការទទួលការយល់ព្រមពីអ្នក។\n(៥) ព័ត៌មានដែលប្រមូលបានត្រូវផ្សព្វផ្សាយជាសាធារណៈដោយខ្លួនឯង\n(៦) ប្រមូលព័ត៌មានពីព័ត៌មានដែលផ្សព្វផ្សាយជាសាធារណៈដោយស្របច្បាប់ដូចជារបាយការណ៍ព័ត៌មានស្របច្បាប់ការបង្ហាញព័ត៌មានរបស់រដ្ឋាភិបាលនិងបណ្តាញផ្សេងទៀត។\n(៧) ចាំបាច់ចុះហត្ថលេខាលើកិច្ចសន្យាស្របតាមតម្រូវការរបស់អ្នក។\n(៨) ចាំបាច់សម្រាប់ការរក្សាប្រតិបត្តិការប្រកបដោយសុវត្ថិភាពនិងស្ថេរភាពនៃសេវាកម្មបង់ប្រាក់ទំនើបដូចជាការរកឃើញនិងដោះស្រាយការបរាជ័យផលិតផលឬសេវាកម្ម។\n(៩) ចាំបាច់សម្រាប់ការរាយការណ៍ព័ត៌មានស្របច្បាប់។\n(១០) នៅពេលដែលចាំបាច់សម្រាប់ស្ថាប័នស្រាវជ្រាវសិក្សាធ្វើស្ថិតិឬការស្រាវជ្រាវសិក្សាផ្អែកលើចំណាប់អារម្មណ៍សាធារណៈនិងដើម្បីផ្តល់លទ្ធផលស្រាវជ្រាវឬលទ្ធផលពិពណ៌នាខាងក្រៅកំណត់ព័ត៌មានដែលមាននៅក្នុងលទ្ធផល\n(១១) ស្ថានភាពផ្សេងទៀតដែលមានចែងក្នុងច្បាប់និងបទប្បញ្ញត្តិ។\nព័ត៌មានដែលយើងអាចចែករំលែកផ្ទេរឬផ្សព្វផ្សាយ\nការចែករំលែក\nលើកលែងតែកាលៈទេសៈដូចខាងក្រោម Modern Pay នឹងមិនចែករំលែកព័ត៌មានផ្ទាល់ខ្លួនរបស់អ្នកជាមួយតតិយជនណាម្នាក់ដោយគ្មានការយល់ព្រមពីអ្នកឡើយ។\nយើងនឹងចែករំលែកព័ត៌មានរបស់អ្នកសម្រាប់គោលបំណងស្របច្បាប់ចាំបាច់ជាក់លាក់ជាក់លាក់និងច្បាស់លាស់។ យើងនឹងចុះហត្ថលេខាលើកិច្ចព្រមព្រៀងស្តីពីការរក្សាការសម្ងាត់តឹងរឹងជាមួយក្រុមហ៊ុនអង្គការនិងបុគ្គលដែលយើងចែករំលែកព័ត៌មានហើយ តម្រូវឲ្យ ពួកគេដំណើរការព័ត៌មានស្របតាមការណែនាំរបស់យើងគោលការណ៍ឯកជនភាពនេះនិងវិធានការរក្សាការសម្ងាត់និងសន្តិសុខដែលពាក់ព័ន្ធដទៃទៀត។\n(១) ផ្តល់សេវាកម្មរបស់យើងដល់អ្នក៖ យើងអាចចែករំលែកព័ត៌មានរបស់អ្នកជាមួយដៃគូនិងភាគីទីបីដើម្បីទទួលបានមុខងារស្នូលដែលអ្នកត្រូវការឬផ្តល់សេវាកម្មដែលអ្នកត្រូវការរួមទាំងៈការផ្តល់ព័ត៌មានបញ្ជាទិញដែលត្រូវគ្នាទៅនឹងអ្នកផ្តល់សេវាដឹកជញ្ចូន។\n(២) ថែរក្សានិងកែលម្អសេវាកម្មរបស់យើង៖ យើងអាចចែករំលែកព័ត៌មានរបស់អ្នកជាមួយដៃគូនិងភាគីទីបីដើម្បីជួយយើងក្នុងការផ្តល់ជូនអ្នកនូវសេវាកម្មដែលមានគោលដៅនិងប្រសើរជាងមុនរួមទាំងៈការផ្ញើអ៊ីមែលឬជំរុញការជូនដំណឹងក្នុងនាមយើង អ្នកផ្តល់សេវាទំនាក់ទំនង។ ល។\n(៣) ដើម្បីទទួលបានខ្លឹមសារដែលបានពិពណ៌នានៅក្នុងផ្នែក“ គោលបំណងនៃការប្រមូលនិងប្រើប្រាស់ព័ត៌មាន” នៃគោលការណ៍ឯកជននេះ។\n(៤) បំពេញកាតព្វកិច្ចរបស់យើងនៅក្នុង“ គោលការណ៍ឯកជន” នេះឬកិច្ចព្រមព្រៀងផ្សេងទៀតដែលយើងបានឈានដល់អ្នកហើយអនុវត្តសិទ្ធិរបស់យើង។\n(៥) ចែករំលែកជាមួយភាគីទីបីដូចជាដៃគូដែល ប្រគល់ឲ្យ យើងដើម្បីផ្សព្វផ្សាយប៉ុន្តែយើងនឹងផ្តល់ជូនអតិថិជនទាំងនេះនូវការគ្របដណ្តប់និងប្រសិទ្ធភាពនៃការផ្សព្វផ្សាយតែប៉ុណ្ណោះហើយនឹងមិនផ្តល់ព័ត៌មានដែលអាចកំណត់អត្តសញ្ញាណអ្នកបានទេរួមមានទាំងឈ្មោះនិងលេខទូរស័ព្ទ។ ឬអ៊ីមែលឬយើងប្រមូលព័ត៌មាននេះដើម្បីកុំអោយគេស្គាល់អត្តសញ្ញាណអ្នក។ ឧទាហរណ៍យើងអាចប្រាប់គណបក្សដែលទុកចិត្តបានថាតើមានមនុស្សប៉ុន្មាននាក់ដែលបានអានព័ត៌មានផ្សព្វផ្សាយរបស់ពួកគេឬទិញទំនិញរបស់គណបក្សដែលបានទុកចិត្តបន្ទាប់ពីបានឃើញព័ត៌មាននេះឬផ្តល់ឱ្យពួកគេនូវព័ត៌មានស្ថិតិដែលមិនអាចកំណត់អត្តសញ្ញាណពួកគេដើម្បីជួយពួកគេឱ្យយល់ពីទស្សនិកជនឬអតិថិជនរបស់ពួកគេ។\n(៦) ក្នុងកំរិតដែលត្រូវបានអនុញ្ញាតដោយច្បាប់និងបទបញ្ញត្តិដើម្បីអនុវត្តតាមច្បាប់រក្សាកំរិតសាខាឬដៃគូរបស់យើងអ្នកឬអ្នកប្រើប្រាស់បណ្តាញបង់ប្រាក់ទំនើបដទៃទៀតឬផលប្រយោជន៍សាធារណៈទ្រព្យសម្បត្តិឬសុវត្ថិភាពពីការខូចខាតដូចជាការការពារការលួចបន្លំនិងខុសច្បាប់ដទៃទៀត សកម្មភាពនិងកាត់បន្ថយហានិភ័យឥណទានយើងអាចផ្លាស់ប្តូរព័ត៌មានជាមួយក្រុមហ៊ុននិងអង្គការដទៃទៀត។ ទោះយ៉ាងណាក៏ដោយនេះមិនរាប់បញ្ចូលព័ត៌មានដែលត្រូវបានលក់ជួលចែកចាយឬត្រូវបានបង្ហាញផ្សេងទៀតសម្រាប់គោលបំណងរកប្រាក់ចំណេញដោយរំលោភលើការប្តេជ្ញាចិត្តដែលបានធ្វើនៅក្នុងគោលការណ៍ឯកជនភាពនេះ។\n(៧) ដើម្បីឆ្លើយតបទៅនឹងតម្រូវការផ្នែកច្បាប់របស់អ្នកជួយដោះស្រាយវិវាទឬជម្លោះរវាងអ្នកនិងអ្នកដទៃ។\n(៨) ផ្តល់ព័ត៌មានរបស់អ្នកតាមសំណើស្របច្បាប់របស់អាណាព្យាបាលរបស់អ្នក។\n(៩) ត្រូវបានផ្តល់ជូនស្របតាមកិច្ចព្រមព្រៀងសេវាកម្មតែមួយដែលបានចុះហត្ថលេខាជាមួយអ្នក (រាប់បញ្ចូលទាំងកិច្ចព្រមព្រៀងអេឡិចត្រូនិចដែលបានចុះហត្ថលេខាតាមអ៊ិនធរណេតនិងវិធានវេទិកាដែលត្រូវគ្នា) ឬឯកសារស្របច្បាប់ផ្សេងទៀត។\n(១០) ផ្តល់ដោយផ្អែកលើការស្រាវជ្រាវផ្នែកសិក្សា។\n(១១) ត្រូវបានផ្តល់ជូនផ្អែកលើផលប្រយោជន៍សង្គមនិងសាធារណៈស្របតាមច្បាប់និងបទប្បញ្ញត្តិ។\n2. ផ្ទេរប្រាក់\n(១) ជាមួយនឹងការអភិវឌ្ឍអាជីវកម្មរបស់យើងជាបន្តបន្ទាប់យើងអាចធ្វើការរួមបញ្ចូលគ្នាការទិញយកការផ្ទេរទ្រព្យសម្បត្តិឬប្រតិបត្តិការស្រដៀងគ្នាហើយព័ត៌មានរបស់អ្នកអាចត្រូវបានផ្ទេរជាផ្នែកនៃប្រតិបត្តិការបែបនេះ។ យើងនឹងតម្រូវឱ្យក្រុមហ៊ុននិងអង្គការថ្មីដែលកាន់ព័ត៌មានរបស់អ្នកបន្តត្រូវបានចងភ្ជាប់ដោយ“ គោលការណ៍ឯកជន” នេះបើមិនដូច្នេះទេយើងនឹងតម្រូវឱ្យក្រុមហ៊ុននិងអង្គការនានាស្នើសុំឱ្យអ្នកមានការអនុញ្ញាតម្តងទៀត។\n(២) បន្ទាប់ពីទទួលបានការព្រមព្រៀងច្បាស់លាស់របស់អ្នកហើយយើងនឹងផ្ទេរព័ត៌មានរបស់អ្នកទៅភាគីផ្សេងទៀត។\nការលាតត្រដាង\nយើងនឹងបង្ហាញព័ត៌មានរបស់អ្នកតែក្នុងកាលៈទេសៈដូចខាងក្រោមនិងក្រោមការអនុវត្តន៍វិធានការណ៍ការពារសុវត្ថិភាពដែលអនុលោមតាមស្តង់ដារឧស្សាហកម្ម៖\n(១) យោងទៅតាមតំរូវការរបស់អ្នកសូមបង្ហាញព័ត៌មានដែលអ្នកបញ្ជាក់ក្រោមវិធីសាស្ត្របង្ហាញដែលអ្នកយល់ស្រប។\n(២) ក្នុងករណីដែលព័ត៌មានរបស់អ្នកត្រូវបានទាមទារស្របតាមតម្រូវការនៃច្បាប់និងបទបញ្ជាការអនុវត្តន៍ច្បាប់រដ្ឋបាលចាំបាច់ឬតម្រូវការតុលាការយើងអាចនឹងលាតត្រដាងព័ត៌មានរបស់អ្នកដោយផ្អែកលើប្រភេទព័ត៌មានដែលត្រូវការនិងវិធីនៃការបង្ហាញ។ នៅក្រោមការអនុវត្តន៍តាមច្បាប់និងបទប្បញ្ញត្តិនៅពេលដែលយើងទទួលបានការស្នើសុំព័ត៌មានខាងលើនេះយើងនឹងតម្រូវឱ្យអ្នកទទួលចេញឯកសារច្បាប់ដែលត្រូវគ្នាដូចជាដីកាកោះហៅឬលិខិតស៊ើបអង្កេត។ យើងជឿជាក់យ៉ាងមុតមាំថាព័ត៌មានដែលស្នើសុំពីយើងគួរតែមានតម្លាភាពតាមដែលអាចធ្វើទៅបានក្នុងវិសាលភាពដែលច្បាប់បានអនុញ្ញាត។ យើងបានពិនិត្យឡើងវិញនូវរាល់សំណូមពរទាំងអស់ដោយយកចិត្តទុកដាក់ដើម្បីធានាថាពួកគេមានមូលដ្ឋានគតិយុត្តនិងមានកំណត់ចំពោះទិន្នន័យដែលទទួលបានដោយអាជ្ញាធរមានសមត្ថកិច្ចសម្រាប់គោលបំណងស៊ើបអង្កេតជាក់លាក់និងសិទ្ធិស្របច្បាប់។\nតើយើងប្រើខូឃីស៍ប្រូកស៊ីការបញ្ចូលទិន្នន័យនិងបច្ចេកវិទ្យាផ្សេងទៀតយ៉ាងដូចម្តេច\n1. ដើម្បីធ្វើឱ្យបទពិសោធន៍ប្រើប្រាស់របស់អ្នកមានភាពងាយស្រួលនៅពេលអ្នកចូលប្រើផលិតផលឬសេវាកម្មដែលផ្តល់ដោយម៉ូឌែនបេកយើងអាចស្គាល់អត្តសញ្ញាណរបស់អ្នកតាមរយៈឯកសារទិន្នន័យតូចមួយដែលអាចជួយសន្សំសំចៃជំហានរបស់អ្នកក្នុងការធ្វើម្តងទៀតនូវព័ត៌មានចុះឈ្មោះឬ ជួយវិនិច្ឆ័យសុវត្ថិភាពគណនីរបស់អ្នក។ ឯកសារទិន្នន័យទាំងនេះអាចជាខូឃីខូឃីស៍ឬឃ្លាំងផ្ទុកទិន្នន័យក្នុងតំបន់ផ្សេងទៀតដែលផ្តល់ដោយកម្មវិធីរុករកឬកម្មវិធីដែលពាក់ព័ន្ធ (សំដៅដល់ "ខូឃីស៍") ។ សូមយល់ថាសេវាកម្មរបស់យើងមួយចំនួនអាចទទួលបានតែតាមរយៈការប្រើប្រាស់ខូឃីស៍។ ប្រសិនបើកម្មវិធីរុករកឬកម្មវិធីរុករកបន្ថែមរបស់អ្នកអនុញ្ញាត	km	1	2020-04-09 17:34:23.077809	0	2020-04-09 17:34:23.077809	1
2020040917344845670627	Modern Pay Privacy Policy\n\nModern Pay is a payment product that provides users with e-wallets. To explain how we collect data about you, how to collect it, and with whom, we will explain related matters to you through this guide.\nIf you have any questions about privacy protection measures not described in this privacy policy, please contact us.\n\ntable of Contents\n1. The range and time of personal information we may collect\nSecond, the purpose of our collection and use of information\n3. Information we may share, transfer or disclose\n4. How do we use cookies, proxies, data embedding and other technologies\n5. How do we store and protect your personal information?\n6. How do you query, modify or delete personal information\n7. Protection of personal information of minors\n8. How to update this privacy policy\nNine, how to contact us\n\n1. The range and time of personal information we may collect\nThe information we may collect includes:\n1. Information you provide\n(1) When you register an account or use our services, you provide us with relevant personal information, including phone number, email, ID card number and fingerprint information.\n(2) The shared information you provide to other parties through our services, and the information stored when you use our services.\n2. Your information we obtained\nIn order to make the Modern Pay service more relevant to the needs of users, understand the adaptability of products, and identify the abnormal status of accounts, we will collect, summarize, and record information, including log information, device information, location information, album information, and common attributes of mobile phones information.\n(1) Log information: When you use the products or services provided by Modern Pay, we will automatically collect your detailed usage of our services and save them as relevant web logs. Including your search query content, IP address, browser type, telecommunications operator, language used, date and time of access, installation of applications and records of web pages you have visited.\n(2) Device information: We will receive and record information about the device you use (including device model, operating system version, device settings, unique device identifier software and hardware features according to the specific permissions you grant during software installation and use) Information), information about the location of the device (including IP address, GPS location, and information about WiFi access points, Bluetooth, and base station sensors that can provide relevant information).\n(3) Common attribute information of mobile phones: including reading external storage, writing to external storage, reading WiFi status and network requests.\n3. Your personal information that we collect through indirect acquisition\nIn order to facilitate us to provide you with a one-stop service based on associated accounts and to facilitate your unified management, you can use the products or services provided by Modern Pay and its affiliated companies through the Modern Pay account at the link entrance provided by us. When you use the above services through our products, you authorize us to receive, summarize, and analyze the source of your personal information or transaction information from our affiliated companies according to actual business and cooperation needs, or you authorize and agree to provide it to us. If you refuse to provide the above information or refuse to authorize, you may not be able to use the corresponding products or services of Modern Pay, or you may not be able to display relevant information. The information we collect in an indirect way also includes shared information about you when other parties use the Modern Pay service.\n4. We will only collect the information necessary for the normal operation of the business function when you use a specific business function. After you stop using the business function, we will stop the act of collecting your personal information.\nSecond, the purpose of our collection and use of information\nWe only collect and use your personal information for the following purposes:\n1. Help you complete the registration\nIn order for us to provide you with services, you need to provide basic registration information, including mobile phone number, and create your user name and password. In some single services, if you only need to use basic services such as browsing and searching, you do not need to register as a Modern Pay user and provide the above information.\n2. Provide you with goods or services\nThe information we collect and use is a necessary condition for providing you with Modern Pay services. If the relevant information is lacking, we will not be able to provide you with the core content of Modern Pay services, including:\n(1) When you buy goods or services, in order to facilitate the delivery of goods or services to you. You need to provide the consignee's personal identification information, name, contact number, and payment status information. If you refuse to provide such information, we will not be able to complete the relevant delivery services. If you order goods or services for other people through Modern Pay, you need to provide the above information of the actual orderer. Before providing us with the above information of the actual orderer, you need to ensure that you have obtained their authorized consent.\n(2) In order to display the order information of your account, we will collect the order information generated during the process of using the Modern Pay service to show you and facilitate you to manage the order. The order information described in this privacy policy includes goods Name, product quantity, product price, order amount, delivery information (distribution address, recipient name, recipient phone number), equipment information, if you need to invoice, will also record your invoice invoicing information.\n(3) When you contact us, we may save your communication / call records and content or the contact information you leave in order to contact you or help you solve the problem, or record the treatment plan and results of related problems .\n(4) In order to confirm the status of the transaction and provide you with after-sales and dispute resolution services, we will collect your transaction, payment, and logistics information related to the progress of the transaction through the transaction object, payment institution, and logistics company you choose based on the exchange, or Share your transaction information with the above service providers.\n(5) When you need to apply for credit service through our partner, you need to collect your call records, address book and contact information at the request of the credit service provider.\n(6) When you need to enable the scanning function or upload your account picture, you need to collect your shooting and album information.\n(7) When you need to open the call charge recharge function, if you agree to us to read your address book and contact information, it will help us bring you a better experience, but if you do not agree, it will not Affect your activation of this function.\n(8) When you need to set a regular purchase function when purchasing a specific financial product, you need to collect your calendar information.\n(9) When you need to activate the fingerprint unlock function, you need to collect your fingerprint information.\n(10) When you use our Haitao service, you need to collect your order information, personal identification information, name, bank account number and contact information for foreign exchange settlement and balance of payments declaration.\n(11) When you log in, for anti-fraud and risk prevention, you need to collect the IMEI / IMSI / ICCID information of your mobile device to determine whether your account is at risk.\n(12) When you log in to our APP, or when you need to activate and use life payment functions and services, you agree that we collect your geographic location information, but if you do not agree, it will not affect your activation of the aforementioned related functions and Use related services.\n3. Push messages to you\n(1) Display and push goods or services for you. We may use your information, your browsing and search records, device information, location information, order information, extract your browsing, search preferences, behavior habits, location information features, and email, SMS or other methods based on feature tags Send you marketing information to provide or promote the following goods or services from us or third parties.\n(2) Notice to you. We may send you service-related notices when necessary (such as when we suspend a single service, change, or terminate a single service due to system maintenance).\nIf you do not want to continue to receive messages from us, you can ask us to stop pushing, for example: according to the SMS unsubscription guidelines, we are asked to stop sending promotional SMS, etc .; but we send messages according to legal regulations or the service agreement of a single service except.\n4. Provide you with security\nIn order to ensure the authenticity of your identity and provide you with better security, you can provide us with identification, facial features, and voiceprint information to complete real-name authentication and subsequent identity verification.\nWe may use your information for customer service, security, fraud monitoring, archiving and backup purposes to ensure the security of the services we provide to you; we may use or integrate your information collected by us and our Partners obtain information that you authorize or share in accordance with the law to comprehensively judge your account and transaction risks, perform identity verification, detect and prevent security incidents, and take necessary records, audits, analysis, and disposal measures in accordance with the law.\n5. Other uses\nWhen we use the information for other purposes not specified in this "Privacy Policy", or when the information collected for a specific purpose is used for other purposes, we will seek your consent in advance.\nYou understand and agree that after collecting your information, we will de-identify the data through technical means, and the de-identified information will not be able to identify you. In this case, we have the right to use the de-identified Information, analyze and commercialize user databases.\n6. Exception for obtaining authorized consent\nAccording to relevant laws and regulations, collecting your information in the following situations does not require your authorization:\n(1) Related to national security and national defense security;\n(2) Related to public safety, public health, and major public interests;\n(3) Related to criminal investigation, prosecution, trial and execution of judgments, etc .;\n(4) Out of the protection of the important legal rights and interests of the main body of information or other individuals, but it is difficult to obtain your consent;\n(5) The collected information is disclosed to the public by yourself;\n(6) Collect information from legally publicly disclosed information, such as legal news reports, government information disclosure and other channels;\n(7) Necessary to sign a contract according to your requirements;\n(8) Necessary for maintaining the safe and stable operation of Modern Pay services, such as discovering and handling product or service failures;\n(9) Necessary for legal news reporting;\n(10) When it is necessary for academic research institutions to conduct statistics or academic research based on public interest, and to provide academic research or description results externally, de-identify the information contained in the results;\n(11) Other situations stipulated by laws and regulations.\n3. Information we may share, transfer or disclose\n1. Sharing\nExcept for the following circumstances, Modern Pay will not share your personal information with any third party without your consent.\nWe will only share your information for legal, legitimate, necessary, specific and clear purposes. We will sign strict confidentiality agreements with companies, organizations and individuals with whom we share information, and require them to process information in accordance with our instructions, this Privacy Policy, and any other relevant confidentiality and security measures.\n(1) Provide our services to you: We may share your information with partners and other third parties to achieve the core functions you need or provide the services you need, including: providing corresponding order information to logistics service providers.\n(2) Maintain and improve our services: We may share your information with partners and other third parties to help us provide you with more targeted and better services, including: sending emails or push notifications on our behalf Communication service providers, etc.\n(3) To achieve the content described in the "Purpose of Collecting and Using Information" section of this Privacy Policy.\n(4) Fulfill our obligations in this "Privacy Policy" or other agreements we have reached with you and exercise our rights.\n(5) Share with third parties such as partners who entrust us to promote, but we will only provide these clients with coverage and effectiveness of the promotion, and will not provide information that can identify you, including name and phone number Or email; or we aggregate this information so that it will not identify you personally. For example, we can tell the entrusting party how many people have read their promotional information or purchased the entrusting party's goods after seeing this information, or provide them with statistical information that cannot identify themselves to help them understand their audience or customers.\n(6) To the extent permitted by laws and regulations, in order to comply with the law, maintain our level of our affiliates or partners, you or other Modern Pay users or the public interest, property or safety from damage, such as to prevent fraud and other illegal Activities and reduce credit risk, we may exchange information with other companies and organizations. However, this does not include information that is sold, rented, shared, or otherwise disclosed for profitable purposes in violation of the commitments made in this Privacy Policy.\n(7) In response to your legal needs, assist in handling disputes or disputes between you and others.\n(8) Provide your information at the legal request of your guardian.\n(9) Provided according to the single service agreement signed with you (including the electronic agreement signed online and the corresponding platform rules) or other legal documents.\n(10) Provided based on academic research.\n(11) Provided based on social and public interests in compliance with laws and regulations.\n2. Transfer\n(1) With the continuous development of our business, we may conduct mergers, acquisitions, asset transfers or similar transactions, and your information may be transferred as part of such transactions. We will require new companies and organizations that hold your information to continue to be bound by this "Privacy Policy". Otherwise, we will require the companies and organizations to ask you for authorization again.\n(2) After obtaining your explicit consent, we will transfer your information to other parties.\n3. Disclosure\nWe will only disclose your information under the following circumstances and under the premise of adopting safety protection measures that comply with industry standards:\n(1) According to your needs, disclose the information you specify under the disclosure method that you explicitly agree to.\n(2) In the case where your information is required according to the requirements of laws and regulations, mandatory administrative law enforcement or judicial requirements, we may disclose your information based on the type of information required and the method of disclosure. Under the premise of compliance with laws and regulations, when we receive the above-mentioned request for information disclosure, we will require the recipient to issue corresponding legal documents, such as subpoenas or investigation letters. We firmly believe that the information requested from us should be as transparent as possible within the scope permitted by law. We have carefully reviewed all requests to ensure that they have a legal basis and are limited to data obtained by competent authorities for specific investigation purposes and legal rights.\n4. How do we use cookies, proxies, data embedding and other technologies\n1. In order to make your access experience easier, when you access the products or services provided by Modern Pay, we may recognize your identity through a small data file, which can save you the step of repeating the registration information, or Help to judge the security of your account. These data files may be cookies, flash cookies, or other local storage provided by your browser or associated applications (collectively referred to as "cookies"). Please understand that some of our services can only be achieved through the use of cookies. If your browser or browser add-on service allows	en	1	2020-04-09 17:34:48.703503	0	2020-04-09 17:34:48.703503	1
\.


--
-- Data for Name: app_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.app_version (v_id, description, version, create_time, update_time, app_url, vs_code, vs_type, is_force, system, is_delete, account_uid, note, status) FROM stdin;
2020042919500378153484	1	1.0.0	2020-04-29 18:50:03	2020-04-29 18:50:03	http:/download.khmodernpay.com/app/2020042919495764336143.apk	1	0	0	0	0	af8ca79c-995e-4336-a1f7-31a76613d300	1	0
2020042919552070824746	1	1.0.0	2020-04-29 18:55:21	2020-04-29 18:55:21	http:/download.khmodernpay.com/app/2020042919551822066717.apk	1	1	0	0	0	af8ca79c-995e-4336-a1f7-31a76613d300	1	0
\.


--
-- Data for Name: app_version_file_log; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.app_version_file_log (id, create_time, account_no, file_name) FROM stdin;
2020042919495764336143	2020-04-29 18:49:58	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919495764336143.apk
2020042919513012049855	2020-04-29 18:51:30	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919513012049855.apk
2020042919534469869503	2020-04-29 18:53:44	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919534469869503.apk
2020042919535554750572	2020-04-29 18:53:56	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919535554750572.apk
2020042919551822066717	2020-04-29 18:55:18	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919551822066717.apk
2020042919555554418106	2020-04-29 18:55:56	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919555554418106.apk
2020042919570561577491	2020-04-29 18:57:06	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919570561577491.apk
2020042919574084743403	2020-04-29 18:57:41	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919574084743403.apk
2020042919583569617359	2020-04-29 18:58:35	af8ca79c-995e-4336-a1f7-31a76613d300	2020042919583569617359.apk
2020042920031246529110	2020-04-29 19:03:13	af8ca79c-995e-4336-a1f7-31a76613d300	2020042920031246529110.apk
2020043009123606373579	2020-04-30 08:12:36	af8ca79c-995e-4336-a1f7-31a76613d300	2020043009123606373579.apk
2020043009251926531928	2020-04-30 08:25:20	af8ca79c-995e-4336-a1f7-31a76613d300	2020043009251926531928.apk
2020043015074041651054	2020-04-30 14:07:41	30c9911c-3e77-4e5b-9f2e-5d7825b378b4	2020043015074041651054.html
2020050615405892249095	2020-05-06 14:40:59	af8ca79c-995e-4336-a1f7-31a76613d300	2020050615405892249095.apk
2020050615582951236662	2020-05-06 14:58:29	af8ca79c-995e-4336-a1f7-31a76613d300	2020050615582951236662.html
2020050717473522884315	2020-05-07 16:47:35	af8ca79c-995e-4336-a1f7-31a76613d300	2020050717473522884315.apk
2020050717485717678343	2020-05-07 16:48:58	af8ca79c-995e-4336-a1f7-31a76613d300	2020050717485717678343.apk
2020050717544874311223	2020-05-07 16:54:48	af8ca79c-995e-4336-a1f7-31a76613d300	2020050717544874311223.apk
\.


--
-- Data for Name: billing_details_results; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.billing_details_results (create_time, bill_no, amount, currency_type, bill_type, account_no, account_type, order_no, balance, order_status, modify_time, servicer_no, op_acc_no) FROM stdin;
2020-05-08 10:22:02	2020050811220292247191	1000000	khr	4	d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	2020050811215377226609	0	3	\N	\N	\N
2020-05-08 10:22:05	2020050811220473130895	100000	usd	4	d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	2020050811213355347603	0	3	\N	\N	\N
2020-05-08 10:22:32	2020050811223118967703	10000	usd	1	d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	2020050811223180215481	0	3	\N	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d
2020-05-08 10:22:32	2020050811223193735069	0	usd	3	d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	2020050811223180215481	0	3	\N	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d
2020-05-08 10:22:55	2020050811225475269311	100	khr	1	d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	2020050811225407234823	0	3	\N	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d
2020-05-08 10:22:55	2020050811225413191086	0	khr	3	d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	2020050811225407234823	0	3	\N	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d
\.


--
-- Data for Name: card; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.card (card_no, account_no, channel_no, name, create_time, is_delete, card_number, balance_type, is_defalut, collect_status, audit_status, note, modify_time) FROM stdin;
be50f203-1aa6-4b99-a50b-5e6d2d0a593f	22222222-2222-2222-2222-222222222222	48f9d538-3ef4-4572-b23c-3f6e194d663b	11	2020-04-30 13:52:26.246118	1	11	usd	0	1	0	123	\N
bde275b9-8bbb-432f-8a59-e8a2e266f7cf	22222222-2222-2222-2222-222222222222	1350ea4d-14b7-42d5-b4d0-138a413ef769	展三年	2020-04-30 14:01:21.150574	0	76375684365865454	usd	0	1	0		\N
40bf2aad-c69f-4e3f-aa6a-775f307d314b	22222222-2222-2222-2222-222222222222	48f9d538-3ef4-4572-b23c-3f6e194d663b	李四	2020-04-30 14:01:39.607067	0	98758753285325325	usd	0	1	0		\N
110b8946-bd15-405c-bdae-eaaca9a3c3d0	22222222-2222-2222-2222-222222222222	241c3c48-dc62-44a8-a9b5-1837943dddcb	组织	2020-04-30 14:02:38.358505	0	88682538597395738535	khr	1	1	0		\N
e583cf59-6349-4928-a2df-d539d6be840a	22222222-2222-2222-2222-222222222222	e5b21ab0-aa8d-4e30-8d8d-4fc92200eead	AD	2020-04-30 14:03:05.601629	0	875837593827535	khr	0	1	0		\N
\.


--
-- Data for Name: cashier; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cashier (uid, name, servicer_no, is_delete, create_time, op_password, modify_time) FROM stdin;
\.


--
-- Data for Name: channel; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.channel (channel_no, channel_name, create_time, is_delete, note, idx, is_recom, channel_type, currency_type, use_status, logo_img_no) FROM stdin;
e5b21ab0-aa8d-4e30-8d8d-4fc92200eead	Canadia Bank	2020-04-30 10:58:18.447175	0		0	0	0		1	2020043011581874467977
48f9d538-3ef4-4572-b23c-3f6e194d663b	Foreign Trade Bank of Cambodia	2020-04-30 11:04:59.253446	0		0	0	0		1	2020043012045933785440
1350ea4d-14b7-42d5-b4d0-138a413ef769	Union Commercial Bank	2020-04-30 11:03:30.568828	0		0	0	0		1	2020043012033016237168
02a97f26-60ee-4ef9-8f15-b1357e79d4d1	Cambodian Public Bank	2020-04-30 10:59:23.613226	0		0	0	0		1	2020043011592385566309
241c3c48-dc62-44a8-a9b5-1837943dddcb	ANZ Royal	2020-04-30 10:58:52.862109	0		0	0	0		1	2020043011585294576819
\.


--
-- Data for Name: channel_servicer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.channel_servicer (channel_no, create_time, is_delete, idx, is_recom, currency_type, use_status, id) FROM stdin;
02a97f26-60ee-4ef9-8f15-b1357e79d4d1	2020-04-30 11:01:14.188374	0	0	1	usd	1	2020043012011491147436
241c3c48-dc62-44a8-a9b5-1837943dddcb	2020-04-30 11:01:21.101024	0	0	0	usd	1	2020043012012194534241
e5b21ab0-aa8d-4e30-8d8d-4fc92200eead	2020-04-30 11:01:28.439863	0	0	0	usd	1	2020043012012874022733
241c3c48-dc62-44a8-a9b5-1837943dddcb	2020-04-30 11:02:20.403443	0	0	0	khr	1	2020043012022042062424
\.


--
-- Data for Name: collection_order; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.collection_order (log_no, from_vaccount_no, to_vaccount_no, amount, create_time, finish_time, order_status, balance_type, payment_type, fees, is_count, modify_time, ip, lat, lng) FROM stdin;
\.


--
-- Data for Name: common_help; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.common_help (help_no, problem, answer, idx, is_delete, use_status, lang, vs_type, modify_time, create_time, file_id) FROM stdin;
2020043015074243016069	12		-1	1	1	zh_CN	0	2020-04-30 14:07:42.057215	2020-04-30 14:07:42.057215	2020043015074041651054
2020050615583152952067	11111		1	0	1	zh_CN	1	2020-05-06 14:58:31.442753	2020-05-06 14:58:31.442753	2020050615582951236662
\.


--
-- Data for Name: consultation_config; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.consultation_config (id, use_status, is_delete, create_time, lang, idx, logo_img_no, name, text) FROM stdin;
\.


--
-- Data for Name: cust; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cust (cust_no, account_no, payment_password, gender, in_authorization, out_authorization, in_transfer_authorization, out_transfer_authorization, modify_time, is_delete, def_pay_no) FROM stdin;
e93642a0-2d80-44d4-a3ab-9e8e22e183e3	d4a4ca0f-973e-484a-80b7-c40187aeda3f		1	1	1	1	1	\N	0	usd_balance
401a2161-2cef-469c-ab05-4e581ea4dffb	5e2447c5-47d4-4589-ae58-09423dd57fd7		1	1	1	1	1	\N	0	usd_balance
\.


--
-- Data for Name: dict_acc_title; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_acc_title (title_no, title_name, parent_title) FROM stdin;
\.


--
-- Data for Name: dict_account_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_account_type (account_type, remark) FROM stdin;
1	管理员
2	运营
3	服务商
4	用户
5	收银员
6	总部
\.


--
-- Data for Name: dict_bank_abbr; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_bank_abbr (id, bank_abbr, bank_name) FROM stdin;
\.


--
-- Data for Name: dict_bankname; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_bankname (bank_name, bank_id) FROM stdin;
\.


--
-- Data for Name: dict_bin_bankname; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_bin_bankname (bin_code, bank_name, org_code, card_name, card_type, card_type_no) FROM stdin;
\.


--
-- Data for Name: dict_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_images (image_id, image_url, create_time, status, modify_time, account_no, is_delete) FROM stdin;
2020042320084378459717	9ba44023926b9eff8f27432d35e7f11b.jpeg	2020-04-23 12:08:43.493363	1	2020-04-23 12:08:43.493363	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	0
2020042414092632968747	48fe64437710ddb09217519e3f7961f1.jpeg	2020-04-24 06:09:26.515146	1	2020-04-24 06:09:26.515146	c1abbc5a-4280-4996-a866-08b74a21d2fb	0
2020042510085259071383	62885aa087617e72761054704ab4737a.jpeg	2020-04-25 09:08:52.460782	1	2020-04-25 09:08:52.460782	da6eea63-13dc-44ba-a9eb-a74755e6dcf2	0
2020042617424570781144	44e2391e022826c203f89a586091a9df.png	2020-04-26 16:42:45.924842	1	2020-04-26 16:42:45.924842	af8ca79c-995e-4336-a1f7-31a76613d300	0
2020042721481143841569	64feb95094a8d6b7fcc4ee635840bcfe.jpeg	2020-04-27 20:48:11.623729	1	2020-04-27 20:48:11.623729	0e8d24af-bec7-4f95-b038-c48045f51abf	0
2020042721484254120099	0bb4ad75bba74fffb928b842a3f21b40.jpeg	2020-04-27 20:48:42.389848	1	2020-04-27 20:48:42.389848	0e8d24af-bec7-4f95-b038-c48045f51abf	0
2020043011581874467977	2c541f13ed5dd88d813674349cf88a07.png	2020-04-30 10:58:18.443691	1	2020-04-30 10:58:18.443691	af8ca79c-995e-4336-a1f7-31a76613d300	0
2020043011585294576819	841d21ed015eed9d342a982159c6f638.png	2020-04-30 10:58:52.859529	1	2020-04-30 10:58:52.859529	af8ca79c-995e-4336-a1f7-31a76613d300	0
2020043011592385566309	36744a730e3fbbe21e4fa7b0b6fa3897.png	2020-04-30 10:59:23.610552	1	2020-04-30 10:59:23.610552	af8ca79c-995e-4336-a1f7-31a76613d300	0
2020043012033016237168	27356f7de387401054ad0ba8867a3fad.png	2020-04-30 11:03:30.565877	1	2020-04-30 11:03:30.565877	af8ca79c-995e-4336-a1f7-31a76613d300	0
2020043012045933785440	957e0aa48aa6b1c1fc1f4295116e34f6.png	2020-04-30 11:04:59.250243	1	2020-04-30 11:04:59.250243	af8ca79c-995e-4336-a1f7-31a76613d300	0
2020050810211740917252	c84f78a10f36fae00b0cffac1862068a.jpeg	2020-05-08 09:21:17.228249	1	2020-05-08 09:21:17.228249	c1abbc5a-4280-4996-a866-08b74a21d2fb	0
\.


--
-- Data for Name: dict_org_abbr; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_org_abbr (org_code, abbr) FROM stdin;
\.


--
-- Data for Name: dict_province; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_province (province_code, province_name, short_name, full_en_name, short_zh_name) FROM stdin;
\.


--
-- Data for Name: dict_region; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_region (id, code, name, level, pid, longitude, latitude, is_leaf, pname) FROM stdin;
\.


--
-- Data for Name: dict_region_bank; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_region_bank (code, name, province, city, bank_type, province_code, city_code) FROM stdin;
\.


--
-- Data for Name: dict_vatype; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dict_vatype (va_type, remark) FROM stdin;
1	美金存款
2	瑞尔存款
3	美金冻结存款
4	瑞尔冻结存款
\.


--
-- Data for Name: exchange_order; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.exchange_order (log_no, in_type, out_type, amount, create_time, rate, order_status, finish_time, account_no, trans_from, trans_amount, err_reason, fees, is_count, modify_time, ip, lat, lng) FROM stdin;
\.


--
-- Data for Name: func_config; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.func_config (func_no, func_name, idx, use_status, is_delete, img, jump_url, application_type, img_id) FROM stdin;
1948cc84-8bb7-43ad-9cf6-4d326ab3d84e	扫一扫	1	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png		1	
cc27ad95-ddba-4cd5-aced-a697dcc9c020	存款	2	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	recharge	1	
889af60a-b0c8-485a-898e-0da28b818c69	手机号取款	3	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	takemoney	1	
6e117d9a-12d3-4dd8-a7df-bf9d34052558	扫一扫取款	4	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	scantakemoney	1	
ea37dd90-0796-4577-a58d-459440e9e421	交易明细	5	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	bill	1	
80f3310f-b978-49f7-acba-c60810ccde5a	额度管理	6	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	quota	1	
b2ae431e-4608-44a9-9b8c-4f70d02a042d	转账至总部	7	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	transbill	1	
0b1832f4-482c-4538-9e02-4a377b7e196d	总部打款记录	8	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	orgmoney	1	
2da14f03-832e-4283-b5d9-792adf2f701c	收益统计	9	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	profit	1	
8a99a126-5866-412b-8be2-d9ea1bddfa7c	对账中心	10	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	recordbill	1	
77764406-e9b2-46df-a01b-0a4929495388	设置	11	1	0	http://cdn.onlinewebfonts.com/svg/download_211361.png	setting	1	
bde6d899-a948-4d0c-8dce-a4b9a635dbde	转账	1	1	0	https://img.khmodernpay.com/trans@2x.png	trans	0	
32fc2072-d4da-47ee-851e-79bf394d4671	充值	2	1	0	https://img.khmodernpay.com/top-up@2x.png	recharge	0	
c462f9e1-8230-4177-a899-afa01dd2da05	收钱	3	1	0	https://img.khmodernpay.com/collection@2x.png	receivables	0	
c5e62532-3df9-41e3-bd78-98b877bdb221	兑换	4	1	0	https://img.khmodernpay.com/exchange@2x.png	exchange	0	
a00f5f30-8156-4804-a576-0c5350f60aec	付款	5	1	0	https://img.khmodernpay.com/payment@2x.png	payment	0	
c9565da0-9dc7-4104-83f6-8ef30a7b5668	账单	6	1	0	https://img.khmodernpay.com/bill@2x.png	bill	0	
\.


--
-- Data for Name: gen_code; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.gen_code (gen_key, account_no, amount, money_type, create_time, code_type, use_status, modify_time, sweep_account_no, order_no, op_acc_type, op_acc_no) FROM stdin;
\.


--
-- Data for Name: global_param; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.global_param (param_key, param_value, remark) FROM stdin;
active_init	F646F21990B7C12CC59FD877F8EDF1A0	
aes_key_active	GAV3BGG6ZCWK8RYT	
operation_no	eb636ecc-5760-470d-924b-d91f06b9b968	运营
sms_access_key_id	zzts	短信密钥1
sms_access_key_secret	zzts123	短信密钥2
aliyun_app_code	a599e0f1bfa6448da664c1ce72f9e401	aliyun ocr密码
mail_port	25	邮件端口
mail_subject	商户对接参数	邮件标题
mail_url	smtp.163.com	邮件服务
bill_file	0 0 1 * * ? 	生成账单文件时间
boc_package_no	fd607bd1-65db-4112-a8b9-86559100a59f	农行套餐号
settle_date	20190423	当前账期
yy_role_no	35fe2459-bc3b-416a-aa91-ba0b2e95b000	代付权限-商户角色uid
t1agent_role_no	8ea59cb3-89a3-4037-83d9-d5c8ef839a88	批量代付商户角色号
gxy_role_no	13658c9a-1e2e-4aab-b012-ec462371c85f	供销易商户角色ID
mail_sender		邮件发送者
mail_cc		邮件cc给谁
mail_authcode		邮箱密码
pcode_key	sa5d6g728ttg$%43JASHGFUIa72	个码私钥
login_aes_key	1234567890123456	登录aes key
password_salt	sa5d6g728ttg$%43JASHGFUIa72	密码盐
pcode_rand_min	10	个码随机金额下限差额
score_login_key	asjkn287tdgjfhgtq7t6q32yutq8	积分用户api登录key
score_user_def_pwd	1e4f9068dd5cf66a101a539860b5c2bc	积分用户默认密码
score_topagency	615a9a61-26d5-4c91-a5cd-d054e587187b	msc运营号
is_score_inv_reg_code	1	是否不要验证码
offline_yx_host		
offline_yx_appid		
pay_url		银联支付跳转基址
pic_base_url		图片基址
settle_file_path	D:\\\\goproj\\\\src\\\\a\\\\pay_chain_dist\\\\webadmin\\\\settle_files	对账文件路径
stat_xlsx_path	D:\\\\goproj\\\\src\\\\a\\\\pay_chain_dist\\\\webadmin/gxy_report	统计报表
store_xlsx_file_path	D:\\\\goproj\\\\src\\\\a\\\\pay_chain_dist\\\\webadmin/stores	门店列表导出
upload_path	D:\\\\goproj\\\\src\\\\a\\\\pay_chain_dist\\\\webadmin/upload	上传文件路径
web_pic_base_url		
bill_file_path	D:\\\\goproj\\\\src\\\\a\\\\pay_chain_dist\\\\settle/billfile	商户下载对账文件
score_setting_direct_rate	10	msc直推费率
score_pay_type_list	wx,ali	msc支付方式列表
coin_lv	3000	兑付阈值/分
score_setting_buy_rate	{"wx":20,"ali":60}	购买费率/万分比
score_coin_back	4123	用户发起兑付的金额/分
score_defrozen_interval	5	超时解冻时间/分
score_defrozen_times	6	冻结次数
score_amt_min	100	金额下限/分
score_amt_max	500000	金额上限/分
mail_body		邮件内容
usd_recv_rate	100	usd收款手续费
cl_sms_account	I6641122	创蓝短信账号
cl_sms_secret	LHSIe1sUPz86f7	创蓝短信密码
offline_yx_md5key		银信md5key
login_sign_key		登录sign key
cl_km_template	khrxxx	创蓝柬埔寨语言短信模板
cl_sms_url	http://intapi.253.com/send/json	创蓝短信商的url
sms_business	cl	创蓝短信商户
acc_plat	22222222-2222-2222-2222-222222222222	平台总账号
cl_zh_CN_template	【Modern Pay】您的验证码是：xxxxxx，请勿泄露。	创蓝中文短信模板
cl_en_template	usxxx	创蓝英文短信模板
zip_file_path		分账文件路径
khr_recv_rate	100	khr收款手续费
def_password_phone	1e4f9068dd5cf66a101a539860b5c2bc	默认手机端密码
def_password	1f82c942befda29b6ed487a51da199f78fce7f05	sha1默认密码
cur_settle_daily_no	2019121914	当前已执行的每小时对账号
usd_to_khr_rate	4000	授权收款额度usd转khr率
defrozen_trading_rights_time	3	冻结交易权限时间
continuous_err_password	5	连输错密码次数
friends_count_limit	2	朋友个数上限
usd_phone_single_min	1000	usd手机号取款单笔最小限额
khr_phone_min_withdraw_fee	100	khr手机号取款最低收取金额
khr_phone_free_fee_per_year	100	khr手机号取款每年免手续费额度
khr_phone_single_max	10000	khr手机号取款单笔最大限额
khr_phone_single_min	1	khr手机号取款单笔最小限额
usd_phone_withdraw_rate	10	usd手机号取款费率
usd_phone_min_withdraw_fee	10	usd手机号取款最低收取金额
usd_phone_single_max	100000	usd手机号取款单笔最大限额
khr_phone_withdraw_rate	1012	khr手机号取款费率
usd_to_khr	4100	USD兑换KHR
no_sms_reg_chk	1	检查短信验证码
is_check_sms	0	是否需要校验,0-需要,1-不需要
bill_xlsx_file_path	/app/static/cust_bills	订单xlsx文件路径
khr_to_usd	4200	KHR兑换USD
usd_face_min_withdraw_fee	10	usd面对面取款最低收取金额
usd_face_free_fee_per_year	0	usd面对面取款每年免手续费额度
usd_face_single_max	1000000	usd面对面取款单笔最大限额
khr_face_withdraw_rate	300	khr面对面取款费率
khr_face_min_withdraw_fee	200	khr面对面取款最低收取金额
khr_face_single_max	100000	khr面对面取款单笔最大限额
app_base_url	https://app-version.khmodernpay.com	版本包基址
usd_phone_free_fee_per_year	100	usd手机号取款每年免手续费额度
usd_deposit_rate	10	usd存款提现费率
usd_min_deposit_fee	10	usd存款最低收取金额
usd_deposit_single_max	2000000	usd存款单笔最大限额
usd_deposit_single_min	200	usd存款单笔最小限额
khr_deposit_rate	10	khr存款手续费
khr_min_deposit_fee	100	khr存款最低收取金额
khr_deposit_single_max	200000	khr存款单笔最大限额
khr_deposit_single_min	2	khr存款单笔最小限额
push_server	fcmpush	推送的商家,如果是jpush,就使用jpush,如果是谷歌,就使用fcmpush
image_base_url	https://img.khmodernpay.com/xxImg	图片域名
store_unauth_image_path	/app/static/img_unauth	不授权的图片路径
store_auth_image_path	/app/static/img_auth	需要授权的图片路径
app_version_file_path	/app/static/app_version	保存版本文件的路径
help_base_url	http://help.khmodernpay.com	帮助包基址
usd_transfer_rate	300	usd转账提现费率
usd_min_transfer_fee	100	usd转账最低收取金额
usd_transfer_single_max	10000	usd转账单笔最大限额
usd_transfer_single_min	0	usd转账单笔最小限额
khr_transfer_rate	300	khr转账手续费
khr_min_transfer_fee	100	khr转账最低收取金额
khr_transfer_single_max	2000	khr转账单笔最大限额
khr_transfer_single_min	2000	khr转账单笔最小限额
usd_to_khr_fee	50	USD兑换KHR单笔手续费
khr_to_usd_fee	1000	KHR兑换USD单笔手续费
usd_face_withdraw_rate	300	usd面对面取款费率
usd_face_single_min	100	usd面对面取款单笔最小限额
khr_face_free_fee_per_year	200	khr面对面取款每年免手续费额度
khr_face_single_min	100	khr面对面取款单笔最小限额
help_file_path	/app/static/help	帮助的答案html文件路径
\.


--
-- Data for Name: headquarters_profit; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.headquarters_profit (log_no, general_ledger_no, amount, create_time, order_status, finish_time, balance_type, profit_source, modify_time) FROM stdin;
2020042413534240823162	2020042413534231343012	9	2020-04-24 05:53:42.204427	3	2020-04-24 05:53:42.204427	usd	5	\N
2020042414023008844107	2020042414023058287986	5	2020-04-24 06:02:30.644448	3	2020-04-24 06:02:30.644448	usd	5	\N
2020042414041178788004	2020042414041126353596	5	2020-04-24 06:04:11.814637	3	2020-04-24 06:04:11.814637	usd	5	\N
2020042611000040027012	2020042610453739122130	5	2020-04-26 10:00:00.081239	3	2020-04-26 10:00:00.081239	usd	5	\N
2020042619001687796415	2020042619001631133543	8	2020-04-26 18:00:16.467531	3	2020-04-26 18:00:16.467531	usd	5	\N
2020042619004472022025	2020042619004475793973	8	2020-04-26 18:00:44.839506	3	2020-04-26 18:00:44.839506	usd	5	\N
2020042716420094918643	2020042716420018848193	111	2020-04-27 15:42:00.490213	3	2020-04-27 15:42:00.490213	usd	2	\N
2020042900003289614410	2020042900003201837845	100	2020-04-28 23:00:32.322929	3	2020-04-28 23:00:32.322929	usd	2	\N
2020042910370502082803	2020042910370516852515	10	2020-04-29 09:37:05.37181	3	2020-04-29 09:37:05.37181	usd	5	\N
2020042910412995463957	2020042910412975170512	9	2020-04-29 09:41:29.402963	3	2020-04-29 09:41:29.402963	usd	5	\N
2020042920004804913698	2020042920004751389451	9	2020-04-29 19:00:48.036091	3	2020-04-29 19:00:48.036091	usd	5	\N
2020043011580614491181	2020043011580619577052	100	2020-04-30 10:58:06.223675	3	2020-04-30 10:58:06.223675	usd	2	\N
2020043015064083765108	2020043015063911518368	5	2020-04-30 14:06:40.041967	3	2020-04-30 14:06:40.041967	usd	5	\N
2020043015085111510583	2020043015085147694949	5	2020-04-30 14:08:51.73457	3	2020-04-30 14:08:51.73457	usd	5	\N
2020043015152264436857	2020043015152267531776	5	2020-04-30 14:15:22.444994	3	2020-04-30 14:15:22.444994	usd	5	\N
2020043016483542384998	2020043016483587444644	9	2020-04-30 15:48:35.391473	3	2020-04-30 15:48:35.391473	usd	5	\N
2020050711022629168343	2020050711022622634565	5	2020-05-07 10:02:26.263836	3	2020-05-07 10:02:26.263836	usd	5	\N
2020050711273113915526	2020050711273131774598	5	2020-05-07 10:27:31.775537	3	2020-05-07 10:27:31.775537	usd	5	\N
2020050711275924341460	2020050711275996940484	5	2020-05-07 10:27:59.136667	3	2020-05-07 10:27:59.136667	usd	5	\N
2020050711290763342113	2020050711290709162496	50	2020-05-07 10:29:07.768233	3	2020-05-07 10:29:07.768233	khr	5	\N
2020050711293651578521	2020050711293658418063	50	2020-05-07 10:29:36.803617	3	2020-05-07 10:29:36.803617	khr	5	\N
2020050811223173930641	2020050811223180215481	10	2020-05-08 10:22:31.67248	3	2020-05-08 10:22:31.67248	usd	5	\N
2020050811225462717543	2020050811225407234823	100	2020-05-08 10:22:54.971183	3	2020-05-08 10:22:54.971183	khr	5	\N
\.


--
-- Data for Name: headquarters_profit_cashable; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.headquarters_profit_cashable (id, cashable_balance, revenue_money, modify_time, currency_type) FROM stdin;
\.


--
-- Data for Name: headquarters_profit_withdraw; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.headquarters_profit_withdraw (order_no, currency_type, amount, note, create_time, account_no) FROM stdin;
\.


--
-- Data for Name: income_order; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.income_order (log_no, act_acc_no, amount, servicer_no, create_time, order_status, finish_time, query_time, balance_type, fees, recv_acc_no, recv_vacc, op_acc_no, settle_hourly_log_no, settle_daily_log_no, payment_type, is_count, modify_time, ree_rate, real_amount, op_acc_type) FROM stdin;
2020050811223180215481	d4a4ca0f-973e-484a-80b7-c40187aeda3f	10000	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	2020-05-08 10:22:31.618546	3	2020-05-08 10:22:31.618546	1970-01-01 00:00:00	usd	10	d4a4ca0f-973e-484a-80b7-c40187aeda3f	8cde9f2a-fc1b-4aeb-a5ad-c0e718a3042c	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d			1	1	2020-05-08 10:22:31.67248	10	9990	1
2020050811225407234823	d4a4ca0f-973e-484a-80b7-c40187aeda3f	100	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	2020-05-08 10:22:54.922215	3	2020-05-08 10:22:54.922215	1970-01-01 00:00:00	khr	100	d4a4ca0f-973e-484a-80b7-c40187aeda3f	cd446382-79cf-4cd2-badd-35bb6d1d5c47	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d			1	1	2020-05-08 10:22:54.971183	10	0	1
\.


--
-- Data for Name: income_ougo_config; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.income_ougo_config (income_ougo_config_no, currency_type, name, use_status, idx, config_type, is_delete) FROM stdin;
\.


--
-- Data for Name: income_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.income_type (income_type, income_name, use_status, idx) FROM stdin;
\.


--
-- Data for Name: lang; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.lang (key, type, is_delete, lang_km, lang_en, lang_ch) FROM stdin;
兑换成功推送消息模板	1	0	你的账户于%s ,%s兑换%s成功,请及时查看！	你的账户于%s ,%s兑换%s成功,请及时查看！	你的账户于%s ,%s兑换%s成功,请及时查看！
兑换成功	1	0	兑换成功	兑换成功	兑换成功
提现成功	1	0	提现成功	提现成功	提现成功
转账成功	1	0	转账成功	转账成功	转账成功
收款成功	1	0	收款成功	收款成功	收款成功
AA100010	3	0	\N	\N	保存失败
0	3	0	\N	\N	正常
AA101002	3	0	\N	\N	上传失败
AA100001	3	0	\N	\N	数据库初始化失败
AA100002	3	0	\N	\N	数据库操作失败
AA100003	3	0	\N	\N	网络异常
提现申请成功推送消息模板	1	0	你的账户于%s 申请支出%s%s成功,请及时查看！	你的账户于%s 申请支出%s%s成功,请及时查看！	你的账户于%s 申请支出%s%s成功,请及时查看！
AA100004	3	0	\N	\N	调用api失败
AA2	3	0	\N	\N	参数有误[%v]为空
提现申请成功	1	0	提现申请成功	提现申请成功	提现申请成功
提现失败	1	0	提现失败	提现失败	提现失败
AA3	3	0	\N	\N	参数有误
AA100005	3	0	\N	\N	io错误
AA100006	3	0	\N	\N	查询失败
提现失败推送消息模板	1	0	你的账户于%s 申请支出%s%s失败,请及时查看！	你的账户于%s 申请支出%s%s失败,请及时查看！	你的账户于%s 申请支出%s%s失败,请及时查看！
AA100007	3	0	\N	\N	添加失败
AA100008	3	0	\N	\N	修改失败
AA100009	3	0	\N	\N	删除失败
AA100011	3	0	\N	\N	不存在此[%v]
AA100013	3	0	\N	\N	没有api调用权限
AA101031	3	0	\N	\N	无对应接口信息
AA1	3	0	\N	\N	不存在此操作员
AA101001	3	0	\N	\N	账户已存在
AA101004	3	0	\N	\N	密码错误
AA101005	3	0	\N	\N	验证码过期或不存在
AA101006	3	0	\N	\N	验证码发送失败
AA101008	3	0	\N	\N	登录过期
AA101009	3	0	\N	\N	图片token过期
AA101010	3	0	\N	\N	图片token太多，冲突了
AA101011	3	0	\N	\N	已经是最新版本
AA101012	3	0	\N	\N	token生成失败
AA101013	3	0	\N	\N	没有访问权限
AA101014	3	0	\N	\N	token刷新失败
AA101015	3	0	\N	\N	密码不能为空
AA101016	3	0	\N	\N	初始化钱包失败
AA101017	3	0	\N	\N	初始化账户失败
AA101018	3	0	\N	\N	角色不存在
AA101019	3	0	\N	\N	菜单url不存在
AA101020	3	0	\N	\N	请求数据无法解密
AA101021	3	0	\N	\N	验证码错误
AA200011	3	0	\N	\N	套餐已存在
AA200012	3	0	\N	\N	套餐不存在
AA200013	3	0	\N	\N	审核失败
AA101030	3	0	\N	\N	初始化机构业务钱包错误
AA101032	3	0	\N	\N	商户已绑定
AA101036	3	0	\N	\N	暂不支持此种类型登录
AA101037	3	0	\N	\N	添加渠道进件失败
AA101038	3	0	\N	\N	服务商角色不存在
AA101039	3	0	\N	\N	商户号为空
AA101040	3	0	\N	\N	账号关系不存在
AA200014	3	0	\N	\N	缺少通道组号
AA200015	3	0	\N	\N	通道组不一致
AA101041	3	0	\N	\N	联行号错误
AA101042	3	0	\N	\N	银行机构号错误
AA101043	3	0	\N	\N	已经登录
AA300001	3	0	\N	\N	产品不存在
AA300002	3	0	\N	\N	api不存在
AA300003	3	0	\N	\N	支付失败
AA300004	3	0	\N	\N	卡片已存在
AA300005	3	0	\N	\N	不支持此类卡片
AA300006	3	0	\N	\N	无对应商户信息
AA300007	3	0	\N	\N	卡片未绑定
AA300008	3	0	\N	\N	订单已发起支付
AA300009	3	0	\N	\N	找不到合适费率
AA300022	3	0	\N	\N	不支持该支付方式
AA300010	3	0	\N	\N	通道异常
AA300011	3	0	\N	\N	无法下单
AA300012	3	0	\N	\N	订单查询异常
AA300013	3	0	\N	\N	无此订单
AA300014	3	0	\N	\N	代付发起失败
AA300015	3	0	\N	\N	没有订单号
AA300016	3	0	\N	\N	退款失败
AA300017	3	0	\N	\N	未找到合适的实际支付方
AA300018	3	0	\N	\N	未找到合适的通道组
AA300019	3	0	\N	\N	未找到合适的通道路由策略
AA300020	3	0	\N	\N	未找到合适的接口组
AA300021	3	0	\N	\N	未找到结算时间类型
AA300023	3	0	\N	\N	对应钱包不足以支付
AA300024	3	0	\N	\N	冻结失败
AA300025	3	0	\N	\N	退款金额异常
AA300026	3	0	\N	\N	退款日期不是同一天
AA101033	3	0	\N	\N	解冻金额失败
AA300027	3	0	\N	\N	通道号不存在
AA300028	3	0	\N	\N	订单录入失败
AA300029	3	0	\N	\N	没有配置商户池
AA300030	3	0	\N	\N	没有配置套餐
AA300031	3	0	\N	\N	手续费计算错误
AA300032	3	0	\N	\N	分账信息格式不对
AA300033	3	0	\N	\N	分账金额不对
AA300034	3	0	\N	\N	商户已存在
AA300035	3	0	\N	\N	不允许多次退款
AA300036	3	0	\N	\N	分单金额不对
AA400002	3	0	\N	\N	platform不对
提现成功推送消息模板	1	0	你的账户于%s 支出%s%s,请及时查看！	你的账户于%s 支出%,请及时查看！	你的账户于%s 支出%s%s,请及时查看！
11	3	0	11	111	2222
khr	1	0	KHR	KHR	KHR
充值成功	1	0	ការបញ្ចូលទឹកប្រាក់ជោគជ័យ	Recharge successfully	充值成功
usd	1	0	USD	USD	USD
AA400003	3	0	\N	\N	invoice_type不对
AA400004	3	0	\N	\N	rate不对
AA400006	3	0	\N	\N	info_invoicer不对
AA400007	3	0	\N	\N	info_checker不对
AA400008	3	0	\N	\N	info_casher不对
AA400009	3	0	\N	\N	seller_bank不对
AA400010	3	0	\N	\N	seller_addr不对
AA400011	3	0	\N	\N	list_goods_name不对
AA400012	3	0	\N	\N	list_number不对
AA400013	3	0	\N	\N	amount不对
AA101034	3	0	\N	\N	不能添加同类型的第三方商户
AA101035	3	0	\N	\N	已关联其它终端
AA500001	3	0	\N	\N	获取套餐资料失败
AA500002	3	0	\N	\N	获取商户信息失败
AA500003	3	0	\N	\N	地址信息有误
AA500004	3	0	\N	\N	进件失败
AA200006	3	0	\N	\N	添加失败,规则已存在
AA200008	3	0	\N	\N	获取商户信息失败_EXISTENT
AA200007	3	0	\N	\N	渠道不存在
AA200009	3	0	\N	\N	通道不存在
AA200010	3	0	\N	\N	合作机构不存在
AA600001	3	0	\N	\N	邀请码信息有误
AA600002	3	0	\N	\N	商品分类不存在
AA101022	3	0	\N	\N	商品库存不足
AA101023	3	0	\N	\N	未知商品
AA101024	3	0	\N	\N	包含已下架商品
AA700001	3	0	\N	\N	超出可交易时间段
AA700002	3	0	\N	\N	单笔交易金额超出范围
AA700003	3	0	\N	\N	超出每日可交易总笔数
AA700004	3	0	\N	\N	超出每日可交易总金额
AA700005	3	0	\N	\N	超出频率可交易笔数
AA700006	3	0	\N	\N	超出频率可交易金额
AA700007	3	0	\N	\N	渠道已加入黑名单
AA700008	3	0	\N	\N	商户已加入黑名单
AA700009	3	0	\N	\N	账号手机号码已加入黑名单
AA700010	3	0	\N	\N	商户银行账户已加入黑名单
AA700011	3	0	\N	\N	用户银行账户已加入黑名单
AA700012	3	0	\N	\N	该信用卡已加入黑名单
AA700013	3	0	\N	\N	该储蓄卡已加入黑名单
AA700014	3	0	\N	\N	该通道已加入黑名单
AA700015	3	0	\N	\N	渠道已加入白名单
AA700016	3	0	\N	\N	商户已加入白名单
AA700017	3	0	\N	\N	账号手机号码已加入白名单
AA700018	3	0	\N	\N	商户银行账户已加入白名单
AA700019	3	0	\N	\N	用户银行账户已加入白名单
AA700020	3	0	\N	\N	该信用卡已加入白名单
AA700021	3	0	\N	\N	该储蓄卡已加入白名单
AA700022	3	0	\N	\N	该通道已加入白名单
AA700023	3	0	\N	\N	超出可交易时间段
AA700024	3	0	\N	\N	超出单笔交易额度
AA700025	3	0	\N	\N	超出交易频率可交易笔数
AA700026	3	0	\N	\N	超出交易频率可交易金额
AA800001	3	0	\N	\N	任务过多
AA800002	3	0	\N	\N	回调任务初始化失败
AA101025	3	0	\N	\N	账号不存在
AA101026	3	0	\N	\N	原支付密码错误
AA101027	3	0	\N	\N	用户钱包不存在
AA101028	3	0	\N	\N	支付密码错误
AA101029	3	0	\N	\N	金额不能为空
AA900001	3	0	\N	\N	统计失败
AA700027	3	0	\N	\N	商户已有关联
AA1000001	3	0	\N	\N	已经是最新版本
AA1000002	3	0	\N	\N	激活码已使用或不存在
AA1000003	3	0	\N	\N	无法激活该码
AA1000004	3	0	\N	\N	此类不允许在客户端登录
AA1100001	3	0	\N	\N	会计账号不存在
AA1100002	3	0	\N	\N	记账凭证号不存在
AA1100003	3	0	\N	\N	无对应外付类型
AA1100004	3	0	\N	\N	未找到终端
AA1100005	3	0	\N	\N	超出可结算时间段
AA1100006	3	0	\N	\N	超出时间范围
AA1100006_7	3	0	\N	\N	超出时间范围_7
AA1200003	3	0	\N	\N	执行失败
AA1200004	3	0	\N	\N	任务未审核
AA1200005	3	0	\N	\N	银行账户不能重复
AA1200006	3	0	\N	\N	已有关联商户,无法删除
AA1200008	3	0	\N	\N	订单批次重复
AA1200009	3	0	\N	\N	下载失败
AA1200010	3	0	\N	\N	Rdeis清理失败
AA300037	3	0	\N	\N	金额不足已支付手续费
AA300038	3	0	\N	\N	金额太小
AA300039	3	0	\N	\N	信用卡不能代付
AA100014	3	0	\N	\N	action或url不对
AA200016	3	0	\N	\N	钱包号不存在
AA101044	3	0	\N	\N	缺少运营商
AA101045	3	0	\N	\N	缺少商户名
AA101046	3	0	\N	\N	缺少账号
AA101047	3	0	\N	\N	缺少密码
AA101048	3	0	\N	\N	缺少简称
AA200017	3	0	\N	\N	文件解析错误
AA101049	3	0	\N	\N	角色已存在
AA300041	3	0	\N	\N	出金总额不等于申请金额
AA200018	3	0	\N	\N	还有子套餐，无法删除
AA200019	3	0	\N	\N	还有通道挂着该产品，无法删除
AA300043	3	0	\N	\N	没有匹配的设备
AA300042	3	0	\N	\N	获取二维码失败
AA300044	3	0	\N	\N	无法识别银行卡所属银行
AA300045	3	0	\N	\N	无法识别银行缩写
AA300046	3	0	\N	\N	审核不通过
AA200020	3	0	\N	\N	太多层级
AA101050	3	0	\N	\N	未登录
AA101051	3	0	\N	\N	设置用户状态参数有误
AA200021	3	0	\N	\N	礼包额度不足，购买失败
AA300047	3	0	\N	\N	支付方式错误
AA300048	3	0	\N	\N	币不足
AA300049	3	0	\N	\N	拿码失败
AA300050	3	0	\N	\N	已经有币在途，请先处理
AA300051	3	0	\N	\N	只支持少于%v币的兑付
AA300052	3	0	\N	\N	已有兑付申请未处理
AA300053	3	0	\N	\N	超时，请联系退款
AA100015	3	0	\N	\N	获取全局数据失败
AA200022	3	0	\N	\N	缺少通道号
AA200023	3	0	\N	\N	缺少上级套餐号
AA200024	3	0	\N	\N	该收款码已存在
AA1300001	3	0	\N	\N	用户不存在或密码错误
AA100016	3	0	\N	\N	请求包体为空
AA100017	3	0	\N	\N	请求包体不是json
AA100018	3	0	\N	\N	无法转换商户简码
AA100019	3	0	\N	\N	时间长度小于10位
AA100020	3	0	\N	\N	发起的时间戳已过期
AA100021	3	0	\N	\N	没有找到对应的apikey
AA100022	3	0	\N	\N	没有找到对应的api模式
AA200025	3	0	\N	\N	非被邀请人
AA200026	3	0	\N	\N	邀请人信息获取失败
AA100023	3	0	\N	\N	不支持的图片格式
AA300054	3	0	\N	\N	账号已冻结
AA300055	3	0	\N	\N	您有未结束的礼包订单
AA300056	3	0	\N	\N	订单已超时
AA200027	3	0	\N	\N	不是服务商
AA300057	3	0	\N	\N	取消失败
AA300058	3	0	\N	\N	订单状态错误
AA300059	3	0	\N	\N	暂时不允许此类交易
AA200028	3	0	\N	\N	费率不在合适范围
AA200029	3	0	\N	\N	兑付阈值不在合适范围
AA300060	3	0	\N	\N	订单已取消
AA200030	3	0	\N	\N	无此收款凭证
AA200031	3	0	\N	\N	已审核
AA200032	3	0	\N	\N	购买费率错误
AA300061	3	0	\N	\N	此收款方式暂时不能
AA200033	3	0	\N	\N	审核中
AA200034	3	0	\N	\N	冻结间隔获取失败
AA200035	3	0	\N	\N	冻结次数获取失败
AA200036	3	0	\N	\N	最小金额获取失败
AA200037	3	0	\N	\N	最大金额获取失败
AA200038	3	0	\N	\N	不允许同一时间开启两个同类型的码
AA200039	3	0	\N	\N	没有外部商户号
AA200040	3	0	\N	\N	没有外部机构号
AA300062	3	0	\N	\N	没有兑换费率
AA300063	3	0	\N	\N	码已作废
AA101052	3	0	\N	\N	不能频繁发送短信
AA300064	3	0	\N	\N	支付失踪
AA101054	3	0	\N	\N	短信验证码错误
AA101053	3	0	\N	\N	不能频繁发送短信
AA101055	3	0	\N	\N	原密码不正确
AA101056	3	0	\N	\N	修改密码失败
AA200041	3	0	\N	\N	已经最上层无法调整
AA200042	3	0	\N	\N	已经最下层无法调整
AA300065	3	0	\N	\N	存钱失败
AA300066	3	0	\N	\N	数据库密码为空
AA300067	3	0	\N	\N	密码错误
AA300068	3	0	\N	\N	取款金额不对
AA101057	3	0	\N	\N	修改支付密码失败
AA101058	3	0	\N	\N	修改手机号码失败
AA300069	3	0	\N	\N	取款失败
AA101059	3	0	\N	\N	修改昵称失败
AA300070	3	0	\N	\N	核销码不正确
AA300071	3	0	\N	\N	已有账户,不能使用核销码
AA101060	3	0	\N	\N	修改二维码状态失败
AA300072	3	0	\N	\N	二维码已过期
AA300073	3	0	\N	\N	二维码5分钟后才可以生成
AA300074	3	0	\N	\N	查询收银员的服务商失败
AA101061	3	0	\N	\N	获取扫一扫码的状态失败
AA300076	3	0	\N	\N	客户密码错误
AA101062	3	0	\N	\N	生成图片名失败
AA101063	3	0	\N	\N	保存图片失败
AA101064	3	0	\N	\N	收款人卡号不存在或不正确
AA101065	3	0	\N	\N	转账到总部失败
AA300077	3	0	\N	\N	币种不符合
AA300078	3	0	\N	\N	申请的金额不是10的倍数
AA101066	3	0	\N	\N	向总部请款失败
AA101067	3	0	\N	\N	绑定银行卡失败
AA101068	3	0	\N	\N	银行卡已存在
AA300079	3	0	\N	\N	账号不存在
AA101069	3	0	\N	\N	修改默认卡失败
AA101070	3	0	\N	\N	银行卡不存在
AA101071	3	0	\N	\N	当前卡是默认卡
AA101072	3	0	\N	\N	解除银行卡失败
AA101073	3	0	\N	\N	计算费率输入类型错误
AA101074	3	0	\N	\N	计算费率结果失败
AA101075	3	0	\N	\N	金额或费率为0
AA1300002	3	0	\N	\N	不是用户的账号,不能获取计算费率
AA101076	3	0	\N	\N	没有转入或转出的操作权限
AA101077	3	0	\N	\N	收款账户与收款渠道币种不一致
AA300080	3	0	\N	\N	手续费不正确
AA101078	3	0	\N	\N	操作失败
AA101079	3	0	\N	\N	订单不是待审核状态
AA300081	3	0	\N	\N	额度不足
AA101080	3	0	\N	\N	付款人虚拟账户不存在
AA101081	3	0	\N	\N	收款人虚拟账户不存在
AA101082	3	0	\N	\N	修改头像失败
AA101083	3	0	\N	\N	图片太大了
AA101084	3	0	\N	\N	自己不能给自己转账
AA300082	3	0	\N	\N	金额不能小于0
AA101085	3	0	\N	\N	当前码的状态可能已被扫
AA300083	3	0	\N	\N	没有转入权限
AA300084	3	0	\N	\N	没有转出权限
AA101086	3	0	\N	\N	图片不存在
AA101087	3	0	\N	\N	语言唯一键KEY冲突
AA101088	3	0	\N	\N	账号无对应服务商
AA101089	3	0	\N	\N	修改的昵称不能为空
AA101090	3	0	\N	\N	验证支付密码失败
AA101091	3	0	\N	\N	终端编号或posSn码为空
AA101092	3	0	\N	\N	终端编号已被添加过
AA101093	3	0	\N	\N	终端posSn码已被添加过
AA1400001	3	0	\N	\N	风控不通过
AA101094	3	0	\N	\N	无权限修改他人的pos状态
AA101095	3	0	\N	\N	创建账号非法,长度小于6位
AA300085	3	0	\N	\N	usd 金额应为整数
AA300086	3	0	\N	\N	请设置支付密码
AA101096	3	0	\N	\N	账号已存在,请通过扫码途径取款
AA101097	3	0	\N	\N	收款人手机号不正确
AA101098	3	0	\N	\N	协议使用中，不可删除或修改为不使用
AA101099	3	0	\N	\N	平台盈利可提现余额不足此次提现
AA101100	3	0	\N	\N	请输入取款金额
AA100027	3	0	\N	\N	调用api失败
AA100028	3	0	\N	\N	未授权访问
AA100029	3	0	\N	\N	解密失败
AA100030	3	0	\N	\N	未设置返回码
AA100012	3	0	\N	\N	jwt验签失败
AA101003	3	0	គណនីមិនមានឬពាក្យសម្ងាត់ខុសទេ	Account does not exist or password is wrong	账户不存在或密码错误
AA300040	3	0	តុល្យភាពមិនគ្រប់គ្រាន់	Insufficient balance	余额不足
reg	1	0	[Modern Pay] លេខកូដផ្ទៀងផ្ទាត់របស់អ្នកគឺ៖ xxxxxx សូមកុំបង្ហាញវា។	[Modern Pay] Your verification code is: xxxxxx, please do not disclose it.	【Modern Pay】您的验证码是：xxxxxx，请勿泄露。
backpwd	1	0			
paypwd	1	0			
充值成功推送消息模板	1	0	គណនីរបស់លោកអ្នកនៅថ្ងៃទី%sទទួលបាន%s%s សូមពិនិត្យមើល	Your account earns %s%s at %s, please check in time!	你的账户于%s 收入%s%s,请及时查看！
转账成功推送消息模板	1	0	គណនីរបស់លោកអ្នកនៅថ្ងៃទី%s ម៉ោង​%sទទួលបាន%s USD សូមពិនិត្យមើល	Your account earns %s%s at %s, please check in time!	你的账户于%s 收入%s%s,请及时查看！
\.


--
-- Data for Name: log_account; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_account (log_no, description, account_uid, log_time, type) FROM stdin;
2020050811200396585697	修改密码	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2020-05-08 10:20:03.021709	0
2020050811210930693244	修改支付密码	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2020-05-08 10:21:09.004157	1
\.


--
-- Data for Name: log_account_web; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_account_web (log_no, description, account_uid, create_time, type) FROM stdin;
\.


--
-- Data for Name: log_app_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_app_messages (log_no, order_no, order_type, is_read, is_push, account_no, create_time) FROM stdin;
2020042414023057024127	2020042414023058287986	2	1	0	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 06:02:31
2020042413534265710103	2020042413534231343012	2	1	0	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 05:53:42
2020042414041118462215	2020042414041126353596	2	1	0	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 06:04:12
2020042610453765831369	2020042610453739122130	2	1	0	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-26 09:45:38
2020042619004435054336	2020042619004475793973	2	1	0	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-26 18:00:45
2020042619001673698208	2020042619001631133543	2	1	0	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-26 18:00:16
2020042716420043989563	2020042716420018848193	4	0	0	72ed5801-7efe-4cc0-8a43-ae84a5f8f055	2020-04-27 15:42:00
2020042716420079150565	2020042716420018848193	4	1	0	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-27 15:42:00
2020042900003297028080	2020042900003201837845	4	0	0	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-28 23:00:32
2020042900003262721424	2020042900003201837845	4	0	0	2215fd8b-5b5f-40f6-b676-8a5080db660f	2020-04-28 23:00:32
2020042920004830548697	2020042920004751389451	2	1	0	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-29 19:00:48
2020042910412991158739	2020042910412975170512	2	1	0	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-29 09:41:29
2020043011580632745266	2020043011580619577052	4	0	0	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-30 10:58:06
2020043011580600469469	2020043011580619577052	4	0	0	2215fd8b-5b5f-40f6-b676-8a5080db660f	2020-04-30 10:58:06
2020043015064073816546	2020043015063911518368	2	0	0	c07ac194-8a49-4619-ad67-171cb689f987	2020-04-30 14:06:40
2020043015085110613512	2020043015085147694949	2	0	0	b04709d0-6572-4ff9-a3a3-1682db713846	2020-04-30 14:08:52
2020043015152256445127	2020043015152267531776	2	0	0	8472bdf8-e372-48f7-9544-a091d07d6b6f	2020-04-30 14:15:22
2020043016483519032482	2020043016483587444644	2	1	0	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-30 15:48:35
2020042910370540592933	2020042910370516852515	2	1	0	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-29 09:37:05
2020050711022607820698	2020050711022622634565	2	0	0	c07ac194-8a49-4619-ad67-171cb689f987	2020-05-07 10:02:26
2020050711273197348789	2020050711273131774598	2	0	0	8472bdf8-e372-48f7-9544-a091d07d6b6f	2020-05-07 10:27:32
2020050711275947821895	2020050711275996940484	2	0	0	b04709d0-6572-4ff9-a3a3-1682db713846	2020-05-07 10:27:59
2020050711290723311903	2020050711290709162496	2	0	0	8472bdf8-e372-48f7-9544-a091d07d6b6f	2020-05-07 10:29:08
2020050711293610327697	2020050711293658418063	2	0	0	b04709d0-6572-4ff9-a3a3-1682db713846	2020-05-07 10:29:37
2020050811223171438694	2020050811223180215481	2	0	0	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2020-05-08 10:22:32
2020050811225419981364	2020050811225407234823	2	0	0	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2020-05-08 10:22:55
\.


--
-- Data for Name: log_card; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_card (log_no, card_num, name, account_no, va_type, channel_no, channel_type, create_time, descript) FROM stdin;
\.


--
-- Data for Name: log_exchange_rate; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_exchange_rate (log_time, usd_khr, khr_usd) FROM stdin;
\.


--
-- Data for Name: log_login; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_login (log_time, acc_no, ip, result, client, log_no) FROM stdin;
2020-05-08 10:18:26.856677	4b21eb6b-96cf-439f-88c4-c023ddf6c4b1	172.68.255.31	0	1	2020050811182605271832
2020-05-08 13:01:22.122358	22222222-2222-2222-2222-222222222222	172.69.33.208	0	1	2020050814012204249967
\.


--
-- Data for Name: log_to_headquarters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_to_headquarters (log_no, servicer_no, currency_type, amount, order_status, collection_type, card_no, create_time, finish_time, order_type, image_id) FROM stdin;
2020050811215377226609	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	khr	1000000	1	2	110b8946-bd15-405c-bdae-eaaca9a3c3d0	2020-05-08 10:21:53	1970-01-01 00:00:00	1	1
2020050811213355347603	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	usd	100000	1	2	bde275b9-8bbb-432f-8a59-e8a2e266f7cf	2020-05-08 10:21:33	1970-01-01 00:00:00	1	1
\.


--
-- Data for Name: log_to_servicer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_to_servicer (log_no, currency_type, servicer_no, collection_type, card_no, amount, create_time, order_type, order_status, finish_time, motify_time) FROM stdin;
\.


--
-- Data for Name: log_vaccount; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.log_vaccount (log_no, vaccount_no, create_time, amount, op_type, frozen_balance, balance, reason, settle_hourly_log_no, settle_daily_log_no, biz_log_no) FROM stdin;
2020050811213302733793	cb34afc3-f848-4c95-910b-5a55f100f9bf	2020-05-08 10:21:33.420186	-100000	7	0	0	1			2020050811213355347603
2020050811215396247240	c5157dbc-764b-4256-8e3c-9c11e559c724	2020-05-08 10:21:53.483423	-1000000	7	0	0	1			2020050811215377226609
2020050811220253049007	c5157dbc-764b-4256-8e3c-9c11e559c724	2020-05-08 10:22:02.476122	-1000000	6	-1000000	0	1			2020050811215377226609
2020050811220449191410	cb34afc3-f848-4c95-910b-5a55f100f9bf	2020-05-08 10:22:04.803056	-100000	6	-100000	0	1			2020050811213355347603
2020050811223107934145	cb34afc3-f848-4c95-910b-5a55f100f9bf	2020-05-08 10:22:31.638952	10000	7	0	-100000	1			2020050811223180215481
2020050811223199897350	cb34afc3-f848-4c95-910b-5a55f100f9bf	2020-05-08 10:22:31.647155	10000	6	10000	-100000	1			2020050811223180215481
2020050811223188759602	8cde9f2a-fc1b-4aeb-a5ad-c0e718a3042c	2020-05-08 10:22:31.618546	9990	1	0	9990	2			2020050811223180215481
2020050811223154329670	cf581976-fe0a-4934-85f7-71729deab9be	2020-05-08 10:22:31.67248	10	1	0	428	6			2020050811223180215481
2020050811223190423525	0a9b5da3-e6a8-4526-9002-782dddd12127	2020-05-08 10:22:31.67248	0	1	0	0	6			2020050811223180215481
2020050811223199417055	cb34afc3-f848-4c95-910b-5a55f100f9bf	2020-05-08 10:22:31.67248	0	2	0	-90000	6			2020050811223180215481
2020050811225417748446	c5157dbc-764b-4256-8e3c-9c11e559c724	2020-05-08 10:22:54.939825	100	7	0	-1000000	1			2020050811225407234823
2020050811225410013048	c5157dbc-764b-4256-8e3c-9c11e559c724	2020-05-08 10:22:54.947218	100	6	100	-1000000	1			2020050811225407234823
2020050811225458084773	cd446382-79cf-4cd2-badd-35bb6d1d5c47	2020-05-08 10:22:54.922215	0	1	0	0	2			2020050811225407234823
2020050811225448668291	ec808b18-cb01-4c80-9e08-10c6ac077623	2020-05-08 10:22:54.971183	100	1	0	200	6			2020050811225407234823
2020050811225465458723	b651303e-dfbc-4c5f-bc01-dae033711144	2020-05-08 10:22:54.971183	0	1	0	0	6			2020050811225407234823
2020050811225409213233	c5157dbc-764b-4256-8e3c-9c11e559c724	2020-05-08 10:22:54.971183	0	2	0	-999900	6			2020050811225407234823
\.


--
-- Data for Name: login_token; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.login_token (acc_no, routes, token, login_time, ip, last_op_time, imei) FROM stdin;
4b21eb6b-96cf-439f-88c4-c023ddf6c4b1	[{"children":[{"children":[{"component_name":"UserManagemenList","icon":"el-icon-s-order","id":"07679b82-27d2-4811-85e6-8a6e73461613","idx":1,"name":"用户列表","path":"/cust_center/user_manage/user_management_list","title":"用户列表"},{"component_name":"CustList","hidden":true,"icon":"el-icon-s-custom","id":"850c0f8d-2235-431c-8e0a-d02baa738c58","idx":2,"name":"用户信息","path":"/cust_center/user_manage/cust_list","title":"用户信息"},{"component_name":"UserBalanceList","icon":"el-icon-s-order","id":"81d8c6fd-f427-409c-aac9-d9dac8e33f97","idx":3,"name":"用户账户明细","path":"/cust_center/user_manage/user_balance_list","title":"用户账户明细"},{"component_name":"UserBillList","hidden":true,"icon":"el-icon-s-custom","id":"e428e7fe-5311-447d-a13f-30d8bbf04653","idx":4,"name":"账单明细","path":"/cust_center/user_manage/user_bill_list","title":"账单明细"},{"component_name":"Agreement","icon":"el-icon-document","id":"b89d0f67-90fe-4ca8-9122-65eda7c53aa2","idx":5,"name":"用户协议与隐私管理","path":"/cust_center/user_manage/agreement","title":"用户协议与隐私管理"}],"component_name":"Content","component_path":"user_manage","icon":"el-icon-s-custom","id":"d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a","idx":1,"name":"用户管理","path":"/cust_center/user_manage","title":"用户管理"},{"children":[{"component_name":"ServicerManagemenList","icon":"el-icon-s-order","id":"81d5505b-df37-4cef-928c-10730baa7b10","idx":1,"name":"服务商列表","path":"/cust_center/servicer_managemen/servicer_managemen_list","title":"服务商列表"},{"component_name":"ServicerConfig","hidden":true,"icon":"el-icon-s-order","id":"3800a410-1da5-46a2-958f-163f98034d55","idx":2,"name":"服务商配置","path":"/cust_center/servicer_managemen/servicer_config","title":"服务商配置"},{"component_name":"ServicerAdd","hidden":true,"icon":"el-icon-s-order","id":"08454882-e75c-4933-b522-6687fda8cb24","idx":3,"name":"添加服务商","path":"/cust_center/servicer_managemen/servicer_add","title":"添加服务商"},{"component_name":"ServicerInfo","hidden":true,"icon":"el-icon-s-order","id":"2c2961d0-ef1e-4c1c-8f28-94db86f33b16","idx":5,"name":"服务商信息详情","path":"/cust_center/servicer_managemen/servicer_info","title":"服务商信息详情"},{"component_name":"ServicerOrderCount","hidden":true,"icon":"el-icon-s-order","id":"67d8761f-764e-4e5e-92df-bcb2bfae02ee","idx":6,"name":"指定运营商统计","path":"/cust_center/servicer_managemen/servicer_order_count","title":"指定运营商统计"},{"component_name":"ServiceGeneralLedgerList","icon":"el-icon-s-order","id":"8caf4478-ee5f-48f2-b0cb-a35b2bdc8aa0","idx":7,"name":"服务商交易查询","path":"/cust_center/servicer_managemen/service_general_ledger_list","title":"服务商交易查询"},{"component_name":"ServicerProfitLedgerList","icon":"el-icon-s-order","id":"d470cd75-937c-4229-b180-7b0ace983442","idx":8,"name":"收益明细查询","path":"/cust_center/servicer_managemen/servicer_profit_ledger_list","title":"收益明细查询"},{"component_name":"CashierList","hidden":true,"icon":"el-icon-user-solid","id":"3b618a8e-5499-49a3-a980-dab3944a50d4","idx":9,"name":"店员列表","path":"/cust_center/servicer_managemen/cashier_list","title":"店员列表"},{"component_name":"PosChannelConfig","icon":"el-icon-date","id":"87c1d68e-93a3-4ee9-b0a6-30989584c056","idx":10,"name":"服务商渠道配置","path":"/cust_center/servicer_managemen/pos_channel_config","title":"服务商渠道配置"}],"component_name":"Content","component_path":"servicer_managemen","icon":"el-icon-s-custom","id":"996af46e-f32c-4b26-8109-2b942d1aa8b8","idx":2,"name":"服务商管理","path":"/cust_center/servicer_managemen","title":"服务商管理"},{"children":[{"component_name":"CuponType","icon":"el-icon-s-data","id":"99b01f22-dee8-4f63-ab8e-4b5083d38044","idx":1,"name":"优惠卷类型","path":"/cust_center/cupon/cupon_type","title":"优惠卷类型"},{"component_name":"CuponDetails","icon":"el-icon-s-data","id":"67948dc1-eab6-4fed-ad68-7879030be523","idx":2,"name":"优惠卷明细管理","path":"/cust_center/cupon/cupon_details","title":"优惠卷明细管理"},{"component_name":"CuponAdd","hidden":true,"icon":"el-icon-s-data","id":"f4f9e4b1-e4df-418b-a88d-2bd0f00fa5e5","idx":3,"name":"添加优惠卷","path":"/cust_center/cupon/add_cupon","title":"添加优惠卷"},{"component_name":"CuponSearch","icon":"el-icon-s-data","id":"eb15be96-d809-4a03-b3e4-be03b79b55cc","idx":4,"name":"优惠卷查询","path":"/cust_center/cupon/cupon_search","title":"优惠卷查询"}],"component_name":"Content","hidden":true,"icon":"el-icon-s-order","id":"4e52b2f4-1da6-4001-8073-38f2b12d4951","idx":3,"name":"优惠卷管理","path":"/cust_center/cupon","title":"优惠卷管理"},{"children":[{"component_name":"FinancialCheckCount","hidden":true,"icon":"el-icon-s-data","id":"84a7f310-51f3-4799-a737-d933083f96a2","idx":1,"name":"服务商对账统计","path":"/cust_center/financial_manage/financial_checking_count","title":"服务商对账统计"},{"component_name":"FinancialServicerCheckCount","icon":"el-icon-s-data","id":"4d3344aa-04e1-4b24-80e6-ace57f7e3e6b","idx":2,"name":"服务商对账单","path":"/cust_center/financial_manage/financial_servicer_check_count","title":"服务商对账单"},{"component_name":"BillingDetailsResultsList","hidden":true,"icon":"el-icon-s-data","id":"d68e979e-1fb7-400f-b2d3-4f8b31a08231","idx":3,"name":"账单明细","path":"/cust_center/financial_manage/billing_details_results_list","title":"账单明细"},{"component_name":"CollectionManagementList","icon":"el-icon-s-custom","id":"f938a823-cd1f-4702-9163-45a5e016e2cf","idx":4,"name":"平台收款账户管理","path":"/cust_center/financial_manage/collection_management_list","title":"平台收款账户管理"},{"component_name":"ChannelEdit","hidden":true,"icon":"el-icon-picture-outline-round","id":"d5c1e2ba-4969-45d4-8b30-8da6cc872ae6","idx":5,"name":"新增收款账户","path":"/cust_center/financial_manage/channel_edit","title":"新增收款账户"},{"component_name":"ProfitStatistics","icon":"el-icon-date","id":"ed6cdd20-4817-464c-99b5-f7a2084f6a9c","idx":6,"name":"平台盈利统计","path":"/cust_center/financial_manage/profit_statistics","title":"平台盈利统计"},{"component_name":"ProfitTakeLog","hidden":true,"icon":"el-icon-s-fold","id":"c174cff5-6670-4001-bf92-8395eb9e2656","idx":7,"name":"提现记录","path":"/cust_center/financial_manage/profit_take_log","title":"提现记录"}],"component_name":"Content","icon":"el-icon-s-data","id":"b57a6f39-2ec1-4065-911a-29764b00ec2b","idx":4,"name":"财务管理","path":"/cust_center/financial_manage","title":"财务管理"},{"children":[{"component_name":"PhoneWithdraw","icon":"el-icon-s-claim","id":"35d43dd0-6e55-4bcc-ba62-efc1c4cdff5d","idx":1,"name":"手机号取款","path":"/cust_center/config_manage/phone_withdraw","title":"手机号取款"},{"component_name":"WithdrawConfig","icon":"el-icon-s-claim","id":"22028342-5177-4425-8465-60d41a46eb7e","idx":1,"name":"面对面取款","path":"/cust_center/config_manage/withdraw_config","title":"面对面取款"},{"component_name":"SaveMoneyConfig","icon":"el-icon-s-claim","id":"d344cbbc-eb28-446a-8b92-38bd340cf3b2","idx":2,"name":"存款配置","path":"/cust_center/config_manage/save_money_config","title":"存款配置"},{"component_name":"TransferConfig","icon":"el-icon-s-claim","id":"559d6495-5bab-428d-82d9-4bc1b30e8e06","idx":3,"name":"转账配置","path":"/cust_center/config_manage/transfer_config","title":"转账配置"},{"component_name":"ExchangeRateConfig","icon":"el-icon-s-claim","id":"d212c732-5c62-40ee-b1cc-184eddab4a61","idx":4,"name":"汇率设置","path":"/cust_center/config_manage/exchange_rate_config","title":"汇率设置"},{"component_name":"TransferSecurityConfig","icon":"el-icon-s-claim","id":"8ec4ba0d-da46-49e0-a349-cbaa7f13e26b","idx":5,"name":"交易安全配置","path":"/cust_center/config_manage/transfer_security_config","title":"交易安全配置"},{"component_name":"CollectMethodConfig","icon":"el-icon-s-claim","id":"4a0bcbc4-1eeb-4801-9779-d51641ecd982","idx":6,"name":"充值方式配置","path":"/cust_center/config_manage/collect_method_config","title":"充值方式配置"},{"component_name":"FetchMethodConfig","icon":"el-icon-s-claim","id":"f93f8fe1-b61f-4cd9-acf6-ff7d6fd5fd60","idx":7,"name":"提现方式配置","path":"/cust_center/config_manage/fetch_method_config","title":"提现方式配置"},{"component_name":"MultiLanguageList","icon":"el-icon-s-claim","id":"8d94bbc5-5eb8-4d3d-890f-3a3948b29ce9","idx":8,"name":"多语言配置","path":"/cust_center/config_manage/multi_language_list","title":"多语言配置"},{"component_name":"FuncManage","icon":"el-icon-s-claim","id":"040a0e65-262b-44ce-9a39-336f08560ad4","idx":9,"name":"钱包功能","path":"/cust_center/config_manage/func_manage","title":"钱包功能"},{"component_name":"FuncManageAdd","hidden":true,"icon":"el-icon-s-grid","id":"5b63f7e6-6bf1-406d-8895-b6f92418520c","idx":10,"name":"功能添加","path":"/cust_center/config_manage/func_manage_add","title":"功能添加"},{"component_name":"ChannelConfig","icon":"el-icon-s-claim","id":"1da7c457-afd1-484c-b24c-6d984baecdbf","idx":11,"name":"渠道仓库","path":"/cust_center/config_manage/channel_config","title":"渠道仓库"},{"component_name":"AppVersionCount","icon":"el-icon-s-home","id":"48d6de3b-4af6-4ea3-b796-aca905fa07f3","idx":12,"name":"版本管理","path":"/cust_center/config_manage/app_version_count","title":"版本管理"},{"component_name":"AppVersionList","hidden":true,"icon":"el-icon-menu","id":"8876a714-7d09-43d4-88ea-795c559a72a1","idx":13,"name":"版本列表","path":"/cust_center/config_manage/app_version_list","title":"版本列表"}],"component_name":"Content","icon":"el-icon-s-claim","id":"6264bced-71da-43f1-8fdd-5a96925164cc","idx":5,"name":"配置管理","path":"/cust_center/config_manage","title":"配置管理"},{"children":[{"component_name":"CommonHelpCount","icon":"el-icon-s-comment","id":"fcc9c518-48c0-4c0e-b96a-8d93b165d948","idx":1,"name":"帮助管理","path":"/cust_center/service_config/common_help_count","title":"帮助管理"},{"component_name":"ConsultationConfig","icon":"el-icon-s-custom","id":"439becc5-2668-4f15-8c7f-1b544989df47","idx":2,"name":"咨询管理","path":"/cust_center/service_config/consultation_config","title":"咨询管理"},{"component_name":"CommonHelp","hidden":true,"icon":"el-icon-s-comment","id":"135baf4c-f8cd-4597-bdb5-a49f4bb53e42","idx":3,"name":"帮助列表","path":"/cust_center/service_config/common_help","title":"帮助列表"}],"component_name":"Content","icon":"el-icon-s-data","id":"4a8285de-348e-4845-b0eb-142dd46002f6","idx":6,"name":"客服管理","path":"/cust_center/service_config","title":"客服管理"}],"component_name":"Default","component_path":"cust_center","icon":"el-icon-s-custom","id":"149d7564-5e2f-4472-a637-d1db61dfdc81","idx":1,"name":"用户中心","path":"/cust_center","redirect":{"name":"/cust_center"},"title":"用户中心"},{"children":[{"children":[{"component_name":"AccountList","icon":"el-icon-date","id":"b2776c17-fde6-4af5-8d34-f72d5ec88ac9","idx":1,"name":"账号列表","path":"/account_center/account_manage/account_list","title":"账号列表"},{"component_name":"AccountEdit","hidden":true,"icon":"el-icon-date","id":"622186a1-17b8-4b64-8406-dd3cd1dd2596","idx":3,"name":"账号修改","path":"/account_center/account_manage/account_edit","title":"账号修改"},{"component_name":"RoleList","icon":"el-icon-date","id":"5af52e1b-b6a1-4eca-87ec-466abacb4302","idx":4,"name":"角色列表","path":"/account_center/account_manage/role_list","title":"角色列表"},{"component_name":"ManageAuth","hidden":true,"icon":"el-icon-date","id":"8c9f5d83-fb7b-4f21-8b60-1ed3eb93ff41","idx":5,"name":"管理授权","path":"/account_center/account_manage/manage_auth","title":"管理授权"},{"component_name":"RoleEdit","hidden":true,"icon":"el-icon-date","id":"eb9ea9f9-cd44-4b11-a805-061612ff12ef","idx":5,"name":"角色修改","path":"/account_center/account_manage/role_edit","title":"角色修改"},{"component_name":"RoleAuth","hidden":true,"icon":"el-icon-date","id":"4badf803-d758-47a4-8a4b-fb005cf7167a","idx":6,"name":"角色授权","path":"/account_center/account_manage/role_auth","title":"角色授权"},{"component_name":"AccountCheck","hidden":true,"icon":"el-icon-date","id":"aded78fc-0d49-436e-a780-ebea7c8b0f4c","idx":9,"name":"账号查看","path":"/account_center/account_manage/account_check","title":"账号查看"},{"component_name":"RoleCheck","hidden":true,"icon":"el-icon-date","id":"14f49944-6337-4c1c-8ea4-dd111c7859f1","idx":10,"name":"角色查看","path":"/account_center/account_manage/role_check","title":"角色查看"},{"component_name":"MenuList","icon":"el-icon-date","id":"07bd2d37-21cd-43b3-9d38-eeafaa73bcf9","idx":20,"name":"菜单列表","path":"/account_center/account_manage/menu_list","title":"菜单列表"},{"component_name":"MenuEdit","hidden":true,"icon":"el-icon-date","id":"fa13b0c6-1e65-45de-bef8-ac534342f0f3","idx":21,"name":"菜单修改","path":"/account_center/account_manage/menu_edit","title":"菜单修改"}],"component_name":"Content","component_path":"account_manage","icon":"el-icon-s-custom","id":"222d3e77-c07d-4c30-a384-b0c493a3cb38","idx":1,"name":"账号与菜单","path":"/account_center/account_manage","redirect":{"name":"账号列表"},"title":"账号与菜单"},{"children":[{"component_name":"AuthList","hidden":true,"icon":"el-icon-date","id":"61ab5266-4740-49f3-9649-1b97566a4191","idx":1,"name":"权限列表","path":"/account_center/auth_manage/auth_list","title":"权限列表"},{"component_name":"PayApiList","icon":"el-icon-date","id":"824f43c8-ee42-47ed-bbb9-5f8de09acdbc","idx":3,"name":"api列表","path":"/account_center/auth_manage/pay_api_list","title":"api列表"},{"component_name":"PayApiEdit","hidden":true,"icon":"el-icon-date","id":"90476856-94be-40ed-8b1e-814392074d81","idx":4,"name":"api编辑","path":"/account_center/auth_manage/pay_api_edit","title":"api编辑"},{"component_name":"SetConfig","hidden":true,"icon":"el-icon-date","id":"07722dd0-088c-4309-9fd4-7097cb7a5c58","idx":4,"name":"结算配置","path":"/account_center/auth_manage/set_config","title":"结算配置"},{"component_name":"PayApiCheck","hidden":true,"icon":"el-icon-date","id":"977f77aa-139c-4d4c-8640-053eea3376dd","idx":5,"name":"api查看","path":"/account_center/auth_manage/pay_api_check","title":"api查看"}],"component_name":"Content","component_path":"auth_list","icon":"el-icon-menu","id":"1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc","idx":2,"name":"权限管理","path":"/account_center/auth_manage","redirect":{"name":"权限列表"},"title":"权限管理"},{"children":[{"component_name":"GlobalParamList","icon":"el-icon-date","id":"a168235d-63c0-4b58-8b49-3df20a9a815d","idx":1,"name":"全局参数配置","path":"/account_center/param_manage/global_param_list","title":"全局参数配置"},{"component_name":"GlobalParamCheck","hidden":true,"icon":"el-icon-date","id":"8bb9a2bc-aeea-435f-9132-73caead02416","idx":2,"name":"全局参数配置查看","path":"/account_center/param_manage/global_param_check","title":"全局参数配置查看"},{"component_name":"GlobalParamEdit","hidden":true,"icon":"el-icon-date","id":"c2592cf2-b4f8-4ab8-9ce3-54b4dbbc928a","idx":3,"name":"全局参数配置修改","path":"/account_center/param_manage/global_param_edit","title":"全局参数配置修改"},{"component_name":"RedisManage","icon":"el-icon-date","id":"f42406c3-08c8-443a-845a-996e04c55753","idx":4,"name":"redis管理","path":"/account_center/param_manage/redis_manage","title":"redis管理"},{"component_name":"TemplateManagement","icon":"el-icon-date","id":"d2e33ead-1f7a-4482-a8e0-d58deceea631","idx":5,"name":"模板管理","path":"/account_center/param_manage/template_management","title":"模板管理"}],"component_name":"Content","icon":"el-icon-menu","id":"14aa553e-5be3-404d-adb0-b4440e023c69","idx":3,"name":"参数管理","path":"/account_center/param_manage","title":"参数管理"},{"children":[{"component_name":"WalletList","icon":"el-icon-date","id":"761f6268-1ccb-4a23-b43d-cad43d4761d2","idx":0,"name":"钱包列表","path":"/account_center/my_use/wallet_list","title":"钱包列表"},{"component_name":"MessageCheck","hidden":true,"icon":"el-icon-date","id":"bda3fd00-36a5-43bd-9c2d-0a25511217f3","idx":0,"name":"消息查看","path":"/account_center/my_use/message_check","title":"消息查看"},{"component_name":"WalletFreezeLog","hidden":true,"icon":"el-icon-date","id":"cc9a6c3a-518a-458a-a84e-6353a1c22165","idx":0,"name":"钱包冻结日志","path":"/account_center/my_use/wallet_freeze_log","title":"钱包冻结日志"},{"component_name":"LiquidationData","hidden":true,"icon":"el-icon-date","id":"6e4f00e4-bfad-42af-949d-8b215dae11f2","idx":2,"name":"清分","path":"/account_center/my_use/liquidation_data","title":"清分"},{"component_name":"Wallet","icon":"el-icon-date","id":"46ea3bb2-76b6-4842-a015-33c9edc910ae","idx":3,"name":"钱包","path":"/account_center/my_use/wallet","title":"钱包"},{"component_name":"WalletLog","hidden":true,"icon":"el-icon-date","id":"7270ed76-fcee-4a5a-a37f-6645f69e68aa","idx":3,"name":"钱包日志","path":"/account_center/my_use/wallet_log","title":"钱包日志"},{"component_name":"MessageList","icon":"el-icon-date","id":"621bc69f-2dcf-4754-bcf8-71a18d17fa7e","idx":4,"name":"消息列表","path":"/account_center/my_use/message_list","title":"消息列表"},{"component_name":"MessageEdit","hidden":true,"icon":"reordrer","id":"0dc8b5f7-da88-4bdc-ba24-b5a6bc8a1044","idx":5,"name":"消息修改","path":"/account_center/my_use/message_edit","title":"消息修改"},{"component_name":"PersonalConfig","icon":"el-icon-date","id":"0e143513-3b89-4e2d-9865-8a848ad6c8db","idx":6,"name":"配置","path":"/account_center/my_use/personal_config","title":"配置"},{"component_name":"WalletPassword","hidden":true,"icon":"el-icon-date","id":"b3e49cac-fd53-4dcb-97ea-03ac3c328405","idx":6,"name":"钱包密码设置","path":"/account_center/my_use/wallet_password","title":"钱包密码设置"},{"component_name":"MercChannelWallet","icon":"el-icon-date","id":"8b470f68-065a-4b9a-a35e-0781d83910df","idx":10,"name":"商户通道钱包","path":"/account_center/my_use/merc_channel_wallet","title":"商户通道钱包"},{"component_name":"UnderAccountLogin","icon":"el-icon-date","id":"e9cb3b58-b3ed-4e27-a2f5-0c704dbbf32f","idx":13,"name":"下属账号登录","path":"/account_center/my_use/under_account_login","title":"下属账号登录"},{"component_name":"BillWalletLog","icon":"el-icon-date","id":"18460b2d-6e29-43da-9419-990fd196de7c","idx":14,"name":"结算钱包日志","path":"/account_center/my_use/bill_wallet_log","title":"结算钱包日志"},{"component_name":"MerchantBalance","icon":"el-icon-date","id":"57873320-1f3c-4d56-ad41-122f17f362a2","idx":17,"name":"商户钱包","path":"/account_center/my_use/merchant_balance","title":"商户钱包"}],"component_name":"Content","component_path":"my_use","icon":"el-icon-menu","id":"ed06f369-9d92-4d07-8b9a-6488694b4d1c","idx":4,"name":"我的使用","path":"/account_center/my_use","redirect":{"name":"公司信息"},"title":"我的使用"},{"children":[{"component_name":"TestList","component_path":"test_list","icon":"","id":"5a31e833-8621-4b3b-841a-a0b03a8575e1","idx":1,"name":"测试账号","path":"/account_center/test_manage/test_list","title":"测试账号"}],"component_name":"Content","component_path":"test_manage","icon":"el-icon-setting","id":"2fdb5752-a2e2-4705-8c9b-b4a618e6de1f","idx":5,"name":"测试管理","path":"/account_center/test_manage","title":"测试管理"}],"component_name":"Default","component_path":"account_center","icon":"el-icon-menu","id":"b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f","idx":4,"name":"账号中心","path":"/account_center","redirect":{"name":"账号管理"},"title":"账号中心"},{"children":[{"children":[{"component_name":"SubmitToHeadquarters","icon":"el-icon-date","id":"3ab532fe-274c-4755-b8c4-8055ea124209","idx":1,"name":"服务商充值","path":"/trade_center/order_manage/submit_to_headquarters","title":"服务商充值"},{"component_name":"ReqMoney","icon":"el-icon-magic-stick","id":"63593179-f9d6-4580-97ca-dd8670688272","idx":12,"name":"服务商提现","path":"/trade_center/order_manage/req_money","title":"服务商提现"},{"component_name":"ExchangeOrders","icon":"el-icon-picture-outline-round","id":"6f855990-13d0-43b1-87e6-f0647e99c360","idx":13,"name":"兑换订单","path":"/trade_center/order_manage/exchange_orders","title":"兑换订单"},{"component_name":"SaveOrders","icon":"el-icon-download","id":"ed0c5e11-eddb-48ff-84f9-6730d389f01a","idx":14,"name":"存款订单","path":"/trade_center/order_manage/save_orders","title":"存款订单"},{"component_name":"FetchOrders","icon":"el-icon-upload2","id":"a0e0bacb-c1d2-488f-aa5f-1af550bd97a8","idx":15,"name":"取款订单","path":"/trade_center/order_manage/fetch_orders","title":"取款订单"},{"component_name":"TransferOrders","icon":"el-icon-s-claim","id":"d141790b-1036-447d-b4cf-ae8b75de2190","idx":16,"name":"转账订单","path":"/trade_center/order_manage/transfer_orders","title":"转账订单"},{"component_name":"TradingAccountLog","icon":"el-icon-s-unfold","id":"2f4bda42-a624-4eef-9cc0-aa581733206a","idx":17,"name":"虚拟账户交易流水","path":"/trade_center/order_manage/trading_account_log","title":"虚拟账户交易流水"},{"component_name":"CollectionOrders","hidden":true,"icon":"el-icon-s-claim","id":"18d4ae82-68bc-435c-9e72-bfef6bc6b2e9","idx":18,"name":"收款订单","path":"/trade_center/order_manage/collection_orders","title":"收款订单"}],"component_name":"Content","icon":"el-icon-menu","id":"a8de69fd-5856-42cc-8812-fef457a99079","idx":1,"name":"交易管理","path":"/trade_center/order_manage","redirect":{"name":"交易明细"},"title":"订单管理"},{"children":[{"component_name":"DownloadFiles","icon":"el-icon-date","id":"29058c58-14a6-425d-b366-617860f6950f","idx":1,"name":"对账明细","path":"/trade_center/financial_manage/download_files","title":"对账明细"},{"component_name":"DownloadFiles","hidden":true,"icon":"el-icon-date","id":"d88a3357-b68c-406e-8951-97df93df66bb","idx":1,"name":"通道对账","path":"/trade_center/financial_manage/download_files","title":"通道对账"},{"component_name":"LiquidationList","icon":"el-icon-date","id":"46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d","idx":2,"name":"清算任务详情列表","path":"/trade_center/financial_manage/liquidation_list","title":"清算任务详情列表"},{"component_name":"LiquidationQueueList","hidden":true,"icon":"el-icon-date","id":"97a477c0-7361-421f-b042-2d1e0d9bf5fe","idx":3,"name":"清算任务队列列表","path":"/trade_center/financial_manage/liquidation_queue_list","title":"清算任务队列列表"},{"component_name":"ImmediateReconcilia","icon":"el-icon-date","id":"4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56","idx":3,"name":"即时对账","path":"/trade_center/financial_manage/immediate_reconcilia","title":"即时对账"},{"component_name":"LiquidationCheck","hidden":true,"icon":"el-icon-date","id":"a2bcc31f-515a-45e9-8083-b2f27d2c4afa","idx":4,"name":"查看清算任务详情","path":"/trade_center/financial_manage/liquidation_check","title":"查看清算任务详情"},{"component_name":"LiquidationAudit","icon":"el-icon-date","id":"60f259e8-d103-4fda-859d-b6c053cc6004","idx":7,"name":"清算任务审核列表","path":"/trade_center/financial_manage/liquidation_audit","title":"清算任务审核列表"},{"component_name":"MercReconcilia","icon":"el-icon-date","id":"b36e11e0-99e1-4af5-b83b-6fe9062554a8","idx":8,"name":"交易报表下载","path":"/trade_center/financial_manage/merc_reconcilia","title":"交易报表下载"},{"component_name":"ChannelSettleLogs","icon":"el-icon-date","id":"88c58e99-19a7-4f72-811e-92e21346207a","idx":10,"name":"通道对账分析","path":"/trade_center/financial_manage/channel_settle_logs","title":"通道对账分析"},{"component_name":"CheckDetails","icon":"el-icon-date","id":"59d15cbe-0f30-4057-ab75-09ba8776b7d6","idx":10,"name":"线上对账文件","path":"/trade_center/financial_manage/check_details","title":"线上对账文件"},{"component_name":"T1SettleReview","icon":"el-icon-date","id":"eedb53fb-3a3e-4f32-a232-ffc97a4464ae","idx":11,"name":"T1结算审核","path":"/trade_center/financial_manage/t1_settle_review","title":"T1结算审核"},{"component_name":"T1SettleDetails","hidden":true,"icon":"el-icon-date","id":"9b374efe-f7c1-44ff-a122-52981146b22a","idx":12,"name":"T1结算明细","path":"/trade_center/financial_manage/t1_settle_details","title":"T1结算明细"},{"component_name":"GenOfflineReports","icon":"el-icon-date","id":"8ec08232-29d6-418f-8560-64f307684f8e","idx":14,"name":"生成交易报表","path":"/trade_center/financial_manage/gen_offline_reports","title":"生成交易报表"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"cb332a03-aeee-435c-86de-b7d577e3817b","idx":2,"name":"日结算管理","path":"/trade_center/financial_manage","title":"日结算管理"},{"children":[{"component_name":"TestPrepay","icon":"el-icon-date","id":"4689bc27-a161-431b-bf58-456f66307fd9","idx":1,"name":"主扫测试","path":"/trade_center/pay_test/test_prepay","title":"主扫测试"},{"component_name":"TestPay","icon":"el-icon-date","id":"ef82cdff-1b9d-4396-8d2b-f7b751a50006","idx":2,"name":"被扫测试","path":"/trade_center/pay_test/test_pay","title":"被扫测试"},{"component_name":"TestAgentpay","icon":"el-icon-date","id":"357d5bae-8470-4d03-ab92-68cafedcfa16","idx":3,"name":"代付测试","path":"/trade_center/pay_test/test_agentpay","title":"代付测试"},{"component_name":"TestRefund","icon":"el-icon-date","id":"1897bc12-5b09-412d-a3c7-e001aedff81d","idx":4,"name":"退款测试","path":"/trade_center/pay_test/test_refund","title":"退款测试"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"065c9534-8c0a-408f-892a-914d172f90c7","idx":3,"name":"支付测试","path":"/trade_center/pay_test","title":"支付测试"},{"children":[{"component_name":"MonthlySettlementList","hidden":true,"icon":"el-icon-date","id":"db98b40b-91a4-4434-ba2c-712e3990b82a","idx":4,"name":"月结任务详情列表","path":"/trade_center/settlement/monthly_settlement_list","title":"月结任务详情列表"},{"component_name":"MonthlySettlementChannelList","hidden":true,"icon":"el-icon-date","id":"ad55ef62-54a6-41c0-bba9-830a06822659","idx":5,"name":"月结任务通道详情","path":"/trade_center/settlement/monthly_settlement_channel_list","title":"月结任务通道详情"},{"component_name":"MonthlyStatementPayment","icon":"el-icon-date","id":"e18c306e-1d61-4e6e-a2cf-c41ba357b406","idx":5,"name":"月结代付申请","path":"/trade_center/settlement/monthly_statement_payment","title":"月结代付申请"},{"component_name":"MonthlyLiquidationList","icon":"el-icon-date","id":"b995956c-75be-4635-9022-af1f5376be2b","idx":6,"name":"月清算任务列表","path":"/trade_center/settlement/monthly_liquidation_list","title":"月清算任务列表"},{"component_name":"MonthlyStatementFile","icon":"el-icon-date","id":"aa061751-3180-40dc-b101-f866d7e4dc2f","idx":6,"name":"月结文件","path":"/trade_center/settlement/monthly_statement_file","title":"月结文件"},{"component_name":"MonthlyLiquidationAudit","icon":"el-icon-date","id":"03ed060b-56dd-4ecb-9072-a736c81ae27f","idx":7,"name":"月清算任务审核","path":"/trade_center/settlement/monthly_liquidation_audit","title":"月清算任务审核"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"b4aee72a-2ac5-4748-b846-bd1828139858","idx":3,"name":"月结算管理","path":"/trade_center/settlement","redirect":{"name":"月结算"},"title":"月结算管理"},{"children":[{"component_name":"WalletAudit","icon":"el-icon-date","id":"742c55df-02ca-4d3e-909a-ac8c412ff735","idx":1,"name":"充值审核","path":"/trade_center/recharge_and_pay/wallet_audit","title":"充值审核"},{"component_name":"ImportRechargeFlow","icon":"el-icon-date","id":"48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5","idx":2,"name":"批量充值","path":"/trade_center/recharge_and_pay/import_recharge_flow","title":"批量充值"},{"component_name":"ContraryApplyAudit","icon":"el-icon-date","id":"90df4bc3-dbb0-4e92-ab9a-0aac895e5295","idx":3,"name":"代付申请审核","path":"/trade_center/recharge_and_pay/contrary_apply_audit","title":"代付申请审核"},{"component_name":"NoSettleAudit","icon":"el-icon-date","id":"01698813-a221-4b56-8bf5-ec3ec9956d97","idx":4,"name":"未结算审核","path":"/trade_center/recharge_and_pay/no_settle_audit","title":"未结算审核"},{"component_name":"ContraryApplyCheck","hidden":true,"icon":"el-icon-date","id":"c5e2efa4-c8c7-4c79-9823-59e829a4847e","idx":5,"name":"代付任务详情","path":"/trade_center/recharge_and_pay/contrary_apply_check","title":"代付任务详情"},{"component_name":"WithdrawalManage","icon":"el-icon-date","id":"0576302e-61ce-49b3-b265-28b7458e5c27","idx":14,"name":"出金管理","path":"/trade_center/recharge_and_pay/withdrawal_manage","title":"出金管理"},{"component_name":"WithdrawalManageEdit","hidden":true,"icon":"el-icon-date","id":"3c53a6aa-9d98-43be-bd14-4c5501e58ee5","idx":15,"name":"出金管理编辑","path":"/trade_center/recharge_and_pay/withdrawal_manage_edit","title":"出金管理编辑"},{"component_name":"WithdrawalReview","icon":"el-icon-date","id":"080f0752-21ba-4d3e-97c2-4846dfb4f648","idx":16,"name":"出金审核","path":"/trade_center/recharge_and_pay/withdrawal_review","title":"出金审核"},{"component_name":"SendApplyBalance","icon":"el-icon-date","id":"4e471ceb-766f-4575-ab73-856e07ceb1a8","idx":20,"name":"余额批量代付申请","path":"/trade_center/recharge_and_pay/send_apply_balance","title":"余额批量代付申请"},{"component_name":"ContraryApplyList","icon":"el-icon-date","id":"1c80833a-d9f9-4b31-ac77-55c282781167","idx":21,"name":"对公代付申请列表","path":"/trade_center/recharge_and_pay/contrary_apply_list","title":"对公代付申请列表"},{"component_name":"ContraryApplyEdit","hidden":true,"icon":"el-icon-date","id":"cab6e3e8-cff1-42da-959f-e12705adc760","idx":22,"name":"对公代付申请","path":"/trade_center/recharge_and_pay/contrary_apply_edit","title":"对公代付申请"},{"component_name":"UploadFee","icon":"el-icon-date","id":"3ed50a46-3620-43a1-b8ea-5a1ebc968851","idx":23,"name":"手续费上传","path":"/trade_center/recharge_and_pay/upload_fee","title":"手续费上传"},{"component_name":"FeeReview","icon":"el-icon-date","id":"23476194-48d6-447c-b4e9-2ab6ad7b1a48","idx":24,"name":"手续费审核","path":"/trade_center/recharge_and_pay/fee_review","title":"手续费审核"},{"component_name":"WithdrawalCheck","hidden":true,"icon":"el-icon-date","id":"82dbdde0-be98-462f-bef7-b86ed1a98336","idx":25,"name":"出金查看","path":"/trade_center/recharge_and_pay/withdrawal_check","title":"出金查看"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"e0be0b73-1f90-4c04-8da5-86513c148a16","idx":4,"name":"充值与代付","path":"/trade_center/recharge_and_pay","title":"充值与代付"},{"children":[{"component_name":"BankLimitList","icon":"el-icon-date","id":"a0dee210-7e5b-4309-ad4e-7ccdf1e7e82e","idx":1,"name":"银行限额列表","path":"/trade_center/risk_control_manage/bank_limit_list","title":"银行限额列表"},{"component_name":"BlacklistEdit","hidden":true,"icon":"el-icon-date","id":"72793652-a596-4b80-bbe4-73d10cca1379","idx":3,"name":"黑名单修改","path":"/trade_center/risk_control_manage/blacklist_edit","title":"黑名单修改"},{"component_name":"WhitelistEdit","hidden":true,"icon":"el-icon-date","id":"66cb5fb9-5bdc-4378-8010-6b58f8cedb5f","idx":4,"name":"白名单修改","path":"/trade_center/risk_control_manage/whitelist_edit","title":"白名单修改"},{"component_name":"WindControlEdit","hidden":true,"icon":"el-icon-date","id":"3293247e-3d3c-4749-ae91-1a991e27b8be","idx":5,"name":"风控规则修改","path":"/trade_center/risk_control_manage/wind_control_edit","title":"风控规则修改"},{"component_name":"RiskManage","icon":"el-icon-date","id":"5a8b13c6-3997-4846-b037-2b679f5c31e8","idx":10,"name":"风控管理","path":"/trade_center/risk_control_manage/risk_manage","title":"风控管理"},{"component_name":"ByWindList","icon":"el-icon-date","id":"7ea8ce89-9ce7-4a89-a332-7baad7bbe475","idx":11,"name":"被风控列表","path":"/trade_center/risk_control_manage/by_wind_list","title":"被风控列表"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"5a229d4e-beee-4f0a-ba2f-888f04a149e5","idx":5,"name":"风控管理","path":"/trade_center/risk_control_manage","title":"风控管理"},{"children":[{"component_name":"BankAbbrList","icon":"el-icon-date","id":"af0187d0-e84c-401f-9ffc-70194b099309","idx":1,"name":"银行缩写列表","path":"/trade_center/bank_params_manage/bank_abbr_list","title":"银行缩写列表"},{"component_name":"BankCodeList","icon":"el-icon-date","id":"c2ebd879-c09a-43d3-8088-6b0c2445c0b1","idx":2,"name":"银行编码列表","path":"/trade_center/bank_params_manage/bank_code_list","title":"银行编码列表"},{"component_name":"BankCardList","icon":"el-icon-date","id":"e296b896-76d1-4cf0-b092-e637a5909ec4","idx":3,"name":"银行卡列表","path":"/trade_center/bank_params_manage/bank_card_list","title":"银行卡列表"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"77788f4d-5d0e-458b-b7cb-297e907b2eff","idx":6,"name":"银行参数管理","path":"/trade_center/bank_params_manage","title":"银行参数管理"}],"component_name":"Default","component_path":"trade_center","icon":"el-icon-menu","id":"136f2190-6f80-46d4-a124-d512e508be5e","idx":5,"name":"交易中心","path":"/trade_center","redirect":{"name":"交易管理"},"title":"交易中心"}]	0826aedd76b64d23944d13b41bc8daf5	2020-05-08 10:18:26.880543	172.68.255.31	2020-05-08 10:22:05	
d4a4ca0f-973e-484a-80b7-c40187aeda3f		eb5ecee3282049929e3688b1c8cf458d	2020-05-08 10:23:18.531525	111.92.243.30	\N	
22222222-2222-2222-2222-222222222222	[{"children":[{"children":[{"component_name":"UserManagemenList","icon":"el-icon-s-order","id":"07679b82-27d2-4811-85e6-8a6e73461613","idx":1,"name":"用户列表","path":"/cust_center/user_manage/user_management_list","title":"用户列表"},{"component_name":"CustList","hidden":true,"icon":"el-icon-s-custom","id":"850c0f8d-2235-431c-8e0a-d02baa738c58","idx":2,"name":"用户信息","path":"/cust_center/user_manage/cust_list","title":"用户信息"},{"component_name":"UserBalanceList","icon":"el-icon-s-order","id":"81d8c6fd-f427-409c-aac9-d9dac8e33f97","idx":3,"name":"用户账户明细","path":"/cust_center/user_manage/user_balance_list","title":"用户账户明细"},{"component_name":"UserBillList","hidden":true,"icon":"el-icon-s-custom","id":"e428e7fe-5311-447d-a13f-30d8bbf04653","idx":4,"name":"账单明细","path":"/cust_center/user_manage/user_bill_list","title":"账单明细"},{"component_name":"Agreement","icon":"el-icon-document","id":"b89d0f67-90fe-4ca8-9122-65eda7c53aa2","idx":5,"name":"用户协议与隐私管理","path":"/cust_center/user_manage/agreement","title":"用户协议与隐私管理"}],"component_name":"Content","component_path":"user_manage","icon":"el-icon-s-custom","id":"d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a","idx":1,"name":"用户管理","path":"/cust_center/user_manage","title":"用户管理"},{"children":[{"component_name":"ServicerManagemenList","icon":"el-icon-s-order","id":"81d5505b-df37-4cef-928c-10730baa7b10","idx":1,"name":"服务商列表","path":"/cust_center/servicer_managemen/servicer_managemen_list","title":"服务商列表"},{"component_name":"ServicerConfig","hidden":true,"icon":"el-icon-s-order","id":"3800a410-1da5-46a2-958f-163f98034d55","idx":2,"name":"服务商配置","path":"/cust_center/servicer_managemen/servicer_config","title":"服务商配置"},{"component_name":"ServicerAdd","hidden":true,"icon":"el-icon-s-order","id":"08454882-e75c-4933-b522-6687fda8cb24","idx":3,"name":"添加服务商","path":"/cust_center/servicer_managemen/servicer_add","title":"添加服务商"},{"component_name":"ServicerInfo","hidden":true,"icon":"el-icon-s-order","id":"2c2961d0-ef1e-4c1c-8f28-94db86f33b16","idx":5,"name":"服务商信息详情","path":"/cust_center/servicer_managemen/servicer_info","title":"服务商信息详情"},{"component_name":"ServicerOrderCount","hidden":true,"icon":"el-icon-s-order","id":"67d8761f-764e-4e5e-92df-bcb2bfae02ee","idx":6,"name":"指定运营商统计","path":"/cust_center/servicer_managemen/servicer_order_count","title":"指定运营商统计"},{"component_name":"ServiceGeneralLedgerList","icon":"el-icon-s-order","id":"8caf4478-ee5f-48f2-b0cb-a35b2bdc8aa0","idx":7,"name":"服务商交易查询","path":"/cust_center/servicer_managemen/service_general_ledger_list","title":"服务商交易查询"},{"component_name":"ServicerProfitLedgerList","icon":"el-icon-s-order","id":"d470cd75-937c-4229-b180-7b0ace983442","idx":8,"name":"收益明细查询","path":"/cust_center/servicer_managemen/servicer_profit_ledger_list","title":"收益明细查询"},{"component_name":"CashierList","hidden":true,"icon":"el-icon-user-solid","id":"3b618a8e-5499-49a3-a980-dab3944a50d4","idx":9,"name":"店员列表","path":"/cust_center/servicer_managemen/cashier_list","title":"店员列表"},{"component_name":"PosChannelConfig","icon":"el-icon-date","id":"87c1d68e-93a3-4ee9-b0a6-30989584c056","idx":10,"name":"服务商渠道配置","path":"/cust_center/servicer_managemen/pos_channel_config","title":"服务商渠道配置"}],"component_name":"Content","component_path":"servicer_managemen","icon":"el-icon-s-custom","id":"996af46e-f32c-4b26-8109-2b942d1aa8b8","idx":2,"name":"服务商管理","path":"/cust_center/servicer_managemen","title":"服务商管理"},{"children":[{"component_name":"CuponType","icon":"el-icon-s-data","id":"99b01f22-dee8-4f63-ab8e-4b5083d38044","idx":1,"name":"优惠卷类型","path":"/cust_center/cupon/cupon_type","title":"优惠卷类型"},{"component_name":"CuponDetails","icon":"el-icon-s-data","id":"67948dc1-eab6-4fed-ad68-7879030be523","idx":2,"name":"优惠卷明细管理","path":"/cust_center/cupon/cupon_details","title":"优惠卷明细管理"},{"component_name":"CuponAdd","hidden":true,"icon":"el-icon-s-data","id":"f4f9e4b1-e4df-418b-a88d-2bd0f00fa5e5","idx":3,"name":"添加优惠卷","path":"/cust_center/cupon/add_cupon","title":"添加优惠卷"},{"component_name":"CuponSearch","icon":"el-icon-s-data","id":"eb15be96-d809-4a03-b3e4-be03b79b55cc","idx":4,"name":"优惠卷查询","path":"/cust_center/cupon/cupon_search","title":"优惠卷查询"}],"component_name":"Content","hidden":true,"icon":"el-icon-s-order","id":"4e52b2f4-1da6-4001-8073-38f2b12d4951","idx":3,"name":"优惠卷管理","path":"/cust_center/cupon","title":"优惠卷管理"},{"children":[{"component_name":"FinancialCheckCount","hidden":true,"icon":"el-icon-s-data","id":"84a7f310-51f3-4799-a737-d933083f96a2","idx":1,"name":"服务商对账统计","path":"/cust_center/financial_manage/financial_checking_count","title":"服务商对账统计"},{"component_name":"FinancialServicerCheckCount","icon":"el-icon-s-data","id":"4d3344aa-04e1-4b24-80e6-ace57f7e3e6b","idx":2,"name":"服务商对账单","path":"/cust_center/financial_manage/financial_servicer_check_count","title":"服务商对账单"},{"component_name":"BillingDetailsResultsList","hidden":true,"icon":"el-icon-s-data","id":"d68e979e-1fb7-400f-b2d3-4f8b31a08231","idx":3,"name":"账单明细","path":"/cust_center/financial_manage/billing_details_results_list","title":"账单明细"},{"component_name":"CollectionManagementList","icon":"el-icon-s-custom","id":"f938a823-cd1f-4702-9163-45a5e016e2cf","idx":4,"name":"平台收款账户管理","path":"/cust_center/financial_manage/collection_management_list","title":"平台收款账户管理"},{"component_name":"ChannelEdit","hidden":true,"icon":"el-icon-picture-outline-round","id":"d5c1e2ba-4969-45d4-8b30-8da6cc872ae6","idx":5,"name":"新增收款账户","path":"/cust_center/financial_manage/channel_edit","title":"新增收款账户"},{"component_name":"ProfitStatistics","icon":"el-icon-date","id":"ed6cdd20-4817-464c-99b5-f7a2084f6a9c","idx":6,"name":"平台盈利统计","path":"/cust_center/financial_manage/profit_statistics","title":"平台盈利统计"},{"component_name":"ProfitTakeLog","hidden":true,"icon":"el-icon-s-fold","id":"c174cff5-6670-4001-bf92-8395eb9e2656","idx":7,"name":"提现记录","path":"/cust_center/financial_manage/profit_take_log","title":"提现记录"}],"component_name":"Content","icon":"el-icon-s-data","id":"b57a6f39-2ec1-4065-911a-29764b00ec2b","idx":4,"name":"财务管理","path":"/cust_center/financial_manage","title":"财务管理"},{"children":[{"component_name":"PhoneWithdraw","icon":"el-icon-s-claim","id":"35d43dd0-6e55-4bcc-ba62-efc1c4cdff5d","idx":1,"name":"手机号取款","path":"/cust_center/config_manage/phone_withdraw","title":"手机号取款"},{"component_name":"WithdrawConfig","icon":"el-icon-s-claim","id":"22028342-5177-4425-8465-60d41a46eb7e","idx":1,"name":"面对面取款","path":"/cust_center/config_manage/withdraw_config","title":"面对面取款"},{"component_name":"SaveMoneyConfig","icon":"el-icon-s-claim","id":"d344cbbc-eb28-446a-8b92-38bd340cf3b2","idx":2,"name":"存款配置","path":"/cust_center/config_manage/save_money_config","title":"存款配置"},{"component_name":"TransferConfig","icon":"el-icon-s-claim","id":"559d6495-5bab-428d-82d9-4bc1b30e8e06","idx":3,"name":"转账配置","path":"/cust_center/config_manage/transfer_config","title":"转账配置"},{"component_name":"ExchangeRateConfig","icon":"el-icon-s-claim","id":"d212c732-5c62-40ee-b1cc-184eddab4a61","idx":4,"name":"汇率设置","path":"/cust_center/config_manage/exchange_rate_config","title":"汇率设置"},{"component_name":"TransferSecurityConfig","icon":"el-icon-s-claim","id":"8ec4ba0d-da46-49e0-a349-cbaa7f13e26b","idx":5,"name":"交易安全配置","path":"/cust_center/config_manage/transfer_security_config","title":"交易安全配置"},{"component_name":"CollectMethodConfig","icon":"el-icon-s-claim","id":"4a0bcbc4-1eeb-4801-9779-d51641ecd982","idx":6,"name":"充值方式配置","path":"/cust_center/config_manage/collect_method_config","title":"充值方式配置"},{"component_name":"FetchMethodConfig","icon":"el-icon-s-claim","id":"f93f8fe1-b61f-4cd9-acf6-ff7d6fd5fd60","idx":7,"name":"提现方式配置","path":"/cust_center/config_manage/fetch_method_config","title":"提现方式配置"},{"component_name":"MultiLanguageList","icon":"el-icon-s-claim","id":"8d94bbc5-5eb8-4d3d-890f-3a3948b29ce9","idx":8,"name":"多语言配置","path":"/cust_center/config_manage/multi_language_list","title":"多语言配置"},{"component_name":"FuncManage","icon":"el-icon-s-claim","id":"040a0e65-262b-44ce-9a39-336f08560ad4","idx":9,"name":"钱包功能","path":"/cust_center/config_manage/func_manage","title":"钱包功能"},{"component_name":"FuncManageAdd","hidden":true,"icon":"el-icon-s-grid","id":"5b63f7e6-6bf1-406d-8895-b6f92418520c","idx":10,"name":"功能添加","path":"/cust_center/config_manage/func_manage_add","title":"功能添加"},{"component_name":"ChannelConfig","icon":"el-icon-s-claim","id":"1da7c457-afd1-484c-b24c-6d984baecdbf","idx":11,"name":"渠道仓库","path":"/cust_center/config_manage/channel_config","title":"渠道仓库"},{"component_name":"AppVersionCount","icon":"el-icon-s-home","id":"48d6de3b-4af6-4ea3-b796-aca905fa07f3","idx":12,"name":"版本管理","path":"/cust_center/config_manage/app_version_count","title":"版本管理"},{"component_name":"AppVersionList","hidden":true,"icon":"el-icon-menu","id":"8876a714-7d09-43d4-88ea-795c559a72a1","idx":13,"name":"版本列表","path":"/cust_center/config_manage/app_version_list","title":"版本列表"}],"component_name":"Content","icon":"el-icon-s-claim","id":"6264bced-71da-43f1-8fdd-5a96925164cc","idx":5,"name":"配置管理","path":"/cust_center/config_manage","title":"配置管理"},{"children":[{"component_name":"CommonHelpCount","icon":"el-icon-s-comment","id":"fcc9c518-48c0-4c0e-b96a-8d93b165d948","idx":1,"name":"帮助管理","path":"/cust_center/service_config/common_help_count","title":"帮助管理"},{"component_name":"ConsultationConfig","icon":"el-icon-s-custom","id":"439becc5-2668-4f15-8c7f-1b544989df47","idx":2,"name":"咨询管理","path":"/cust_center/service_config/consultation_config","title":"咨询管理"},{"component_name":"CommonHelp","hidden":true,"icon":"el-icon-s-comment","id":"135baf4c-f8cd-4597-bdb5-a49f4bb53e42","idx":3,"name":"帮助列表","path":"/cust_center/service_config/common_help","title":"帮助列表"}],"component_name":"Content","icon":"el-icon-s-data","id":"4a8285de-348e-4845-b0eb-142dd46002f6","idx":6,"name":"客服管理","path":"/cust_center/service_config","title":"客服管理"}],"component_name":"Default","component_path":"cust_center","icon":"el-icon-s-custom","id":"149d7564-5e2f-4472-a637-d1db61dfdc81","idx":1,"name":"用户中心","path":"/cust_center","redirect":{"name":"/cust_center"},"title":"用户中心"},{"children":[{"children":[{"component_name":"AccountList","icon":"el-icon-date","id":"b2776c17-fde6-4af5-8d34-f72d5ec88ac9","idx":1,"name":"账号列表","path":"/account_center/account_manage/account_list","title":"账号列表"},{"component_name":"AccountEdit","hidden":true,"icon":"el-icon-date","id":"622186a1-17b8-4b64-8406-dd3cd1dd2596","idx":3,"name":"账号修改","path":"/account_center/account_manage/account_edit","title":"账号修改"},{"component_name":"RoleList","icon":"el-icon-date","id":"5af52e1b-b6a1-4eca-87ec-466abacb4302","idx":4,"name":"角色列表","path":"/account_center/account_manage/role_list","title":"角色列表"},{"component_name":"ManageAuth","hidden":true,"icon":"el-icon-date","id":"8c9f5d83-fb7b-4f21-8b60-1ed3eb93ff41","idx":5,"name":"管理授权","path":"/account_center/account_manage/manage_auth","title":"管理授权"},{"component_name":"RoleEdit","hidden":true,"icon":"el-icon-date","id":"eb9ea9f9-cd44-4b11-a805-061612ff12ef","idx":5,"name":"角色修改","path":"/account_center/account_manage/role_edit","title":"角色修改"},{"component_name":"RoleAuth","hidden":true,"icon":"el-icon-date","id":"4badf803-d758-47a4-8a4b-fb005cf7167a","idx":6,"name":"角色授权","path":"/account_center/account_manage/role_auth","title":"角色授权"},{"component_name":"AccountCheck","hidden":true,"icon":"el-icon-date","id":"aded78fc-0d49-436e-a780-ebea7c8b0f4c","idx":9,"name":"账号查看","path":"/account_center/account_manage/account_check","title":"账号查看"},{"component_name":"RoleCheck","hidden":true,"icon":"el-icon-date","id":"14f49944-6337-4c1c-8ea4-dd111c7859f1","idx":10,"name":"角色查看","path":"/account_center/account_manage/role_check","title":"角色查看"},{"component_name":"MenuList","icon":"el-icon-date","id":"07bd2d37-21cd-43b3-9d38-eeafaa73bcf9","idx":20,"name":"菜单列表","path":"/account_center/account_manage/menu_list","title":"菜单列表"},{"component_name":"MenuEdit","hidden":true,"icon":"el-icon-date","id":"fa13b0c6-1e65-45de-bef8-ac534342f0f3","idx":21,"name":"菜单修改","path":"/account_center/account_manage/menu_edit","title":"菜单修改"}],"component_name":"Content","component_path":"account_manage","icon":"el-icon-s-custom","id":"222d3e77-c07d-4c30-a384-b0c493a3cb38","idx":1,"name":"账号与菜单","path":"/account_center/account_manage","redirect":{"name":"账号列表"},"title":"账号与菜单"},{"children":[{"component_name":"AuthList","hidden":true,"icon":"el-icon-date","id":"61ab5266-4740-49f3-9649-1b97566a4191","idx":1,"name":"权限列表","path":"/account_center/auth_manage/auth_list","title":"权限列表"},{"component_name":"PayApiList","icon":"el-icon-date","id":"824f43c8-ee42-47ed-bbb9-5f8de09acdbc","idx":3,"name":"api列表","path":"/account_center/auth_manage/pay_api_list","title":"api列表"},{"component_name":"PayApiEdit","hidden":true,"icon":"el-icon-date","id":"90476856-94be-40ed-8b1e-814392074d81","idx":4,"name":"api编辑","path":"/account_center/auth_manage/pay_api_edit","title":"api编辑"},{"component_name":"SetConfig","hidden":true,"icon":"el-icon-date","id":"07722dd0-088c-4309-9fd4-7097cb7a5c58","idx":4,"name":"结算配置","path":"/account_center/auth_manage/set_config","title":"结算配置"},{"component_name":"PayApiCheck","hidden":true,"icon":"el-icon-date","id":"977f77aa-139c-4d4c-8640-053eea3376dd","idx":5,"name":"api查看","path":"/account_center/auth_manage/pay_api_check","title":"api查看"}],"component_name":"Content","component_path":"auth_list","icon":"el-icon-menu","id":"1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc","idx":2,"name":"权限管理","path":"/account_center/auth_manage","redirect":{"name":"权限列表"},"title":"权限管理"},{"children":[{"component_name":"GlobalParamList","icon":"el-icon-date","id":"a168235d-63c0-4b58-8b49-3df20a9a815d","idx":1,"name":"全局参数配置","path":"/account_center/param_manage/global_param_list","title":"全局参数配置"},{"component_name":"GlobalParamCheck","hidden":true,"icon":"el-icon-date","id":"8bb9a2bc-aeea-435f-9132-73caead02416","idx":2,"name":"全局参数配置查看","path":"/account_center/param_manage/global_param_check","title":"全局参数配置查看"},{"component_name":"GlobalParamEdit","hidden":true,"icon":"el-icon-date","id":"c2592cf2-b4f8-4ab8-9ce3-54b4dbbc928a","idx":3,"name":"全局参数配置修改","path":"/account_center/param_manage/global_param_edit","title":"全局参数配置修改"},{"component_name":"RedisManage","icon":"el-icon-date","id":"f42406c3-08c8-443a-845a-996e04c55753","idx":4,"name":"redis管理","path":"/account_center/param_manage/redis_manage","title":"redis管理"},{"component_name":"TemplateManagement","icon":"el-icon-date","id":"d2e33ead-1f7a-4482-a8e0-d58deceea631","idx":5,"name":"模板管理","path":"/account_center/param_manage/template_management","title":"模板管理"}],"component_name":"Content","icon":"el-icon-menu","id":"14aa553e-5be3-404d-adb0-b4440e023c69","idx":3,"name":"参数管理","path":"/account_center/param_manage","title":"参数管理"},{"children":[{"component_name":"WalletList","icon":"el-icon-date","id":"761f6268-1ccb-4a23-b43d-cad43d4761d2","idx":0,"name":"钱包列表","path":"/account_center/my_use/wallet_list","title":"钱包列表"},{"component_name":"MessageCheck","hidden":true,"icon":"el-icon-date","id":"bda3fd00-36a5-43bd-9c2d-0a25511217f3","idx":0,"name":"消息查看","path":"/account_center/my_use/message_check","title":"消息查看"},{"component_name":"WalletFreezeLog","hidden":true,"icon":"el-icon-date","id":"cc9a6c3a-518a-458a-a84e-6353a1c22165","idx":0,"name":"钱包冻结日志","path":"/account_center/my_use/wallet_freeze_log","title":"钱包冻结日志"},{"component_name":"LiquidationData","hidden":true,"icon":"el-icon-date","id":"6e4f00e4-bfad-42af-949d-8b215dae11f2","idx":2,"name":"清分","path":"/account_center/my_use/liquidation_data","title":"清分"},{"component_name":"Wallet","icon":"el-icon-date","id":"46ea3bb2-76b6-4842-a015-33c9edc910ae","idx":3,"name":"钱包","path":"/account_center/my_use/wallet","title":"钱包"},{"component_name":"WalletLog","hidden":true,"icon":"el-icon-date","id":"7270ed76-fcee-4a5a-a37f-6645f69e68aa","idx":3,"name":"钱包日志","path":"/account_center/my_use/wallet_log","title":"钱包日志"},{"component_name":"MessageList","icon":"el-icon-date","id":"621bc69f-2dcf-4754-bcf8-71a18d17fa7e","idx":4,"name":"消息列表","path":"/account_center/my_use/message_list","title":"消息列表"},{"component_name":"MessageEdit","hidden":true,"icon":"reordrer","id":"0dc8b5f7-da88-4bdc-ba24-b5a6bc8a1044","idx":5,"name":"消息修改","path":"/account_center/my_use/message_edit","title":"消息修改"},{"component_name":"PersonalConfig","icon":"el-icon-date","id":"0e143513-3b89-4e2d-9865-8a848ad6c8db","idx":6,"name":"配置","path":"/account_center/my_use/personal_config","title":"配置"},{"component_name":"WalletPassword","hidden":true,"icon":"el-icon-date","id":"b3e49cac-fd53-4dcb-97ea-03ac3c328405","idx":6,"name":"钱包密码设置","path":"/account_center/my_use/wallet_password","title":"钱包密码设置"},{"component_name":"MercChannelWallet","icon":"el-icon-date","id":"8b470f68-065a-4b9a-a35e-0781d83910df","idx":10,"name":"商户通道钱包","path":"/account_center/my_use/merc_channel_wallet","title":"商户通道钱包"},{"component_name":"UnderAccountLogin","icon":"el-icon-date","id":"e9cb3b58-b3ed-4e27-a2f5-0c704dbbf32f","idx":13,"name":"下属账号登录","path":"/account_center/my_use/under_account_login","title":"下属账号登录"},{"component_name":"BillWalletLog","icon":"el-icon-date","id":"18460b2d-6e29-43da-9419-990fd196de7c","idx":14,"name":"结算钱包日志","path":"/account_center/my_use/bill_wallet_log","title":"结算钱包日志"},{"component_name":"MerchantBalance","icon":"el-icon-date","id":"57873320-1f3c-4d56-ad41-122f17f362a2","idx":17,"name":"商户钱包","path":"/account_center/my_use/merchant_balance","title":"商户钱包"}],"component_name":"Content","component_path":"my_use","icon":"el-icon-menu","id":"ed06f369-9d92-4d07-8b9a-6488694b4d1c","idx":4,"name":"我的使用","path":"/account_center/my_use","redirect":{"name":"公司信息"},"title":"我的使用"},{"children":[{"component_name":"TestList","component_path":"test_list","icon":"","id":"5a31e833-8621-4b3b-841a-a0b03a8575e1","idx":1,"name":"测试账号","path":"/account_center/test_manage/test_list","title":"测试账号"}],"component_name":"Content","component_path":"test_manage","icon":"el-icon-setting","id":"2fdb5752-a2e2-4705-8c9b-b4a618e6de1f","idx":5,"name":"测试管理","path":"/account_center/test_manage","title":"测试管理"}],"component_name":"Default","component_path":"account_center","icon":"el-icon-menu","id":"b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f","idx":4,"name":"账号中心","path":"/account_center","redirect":{"name":"账号管理"},"title":"账号中心"},{"children":[{"children":[{"component_name":"SubmitToHeadquarters","icon":"el-icon-date","id":"3ab532fe-274c-4755-b8c4-8055ea124209","idx":1,"name":"服务商充值","path":"/trade_center/order_manage/submit_to_headquarters","title":"服务商充值"},{"component_name":"ReqMoney","icon":"el-icon-magic-stick","id":"63593179-f9d6-4580-97ca-dd8670688272","idx":12,"name":"服务商提现","path":"/trade_center/order_manage/req_money","title":"服务商提现"},{"component_name":"ExchangeOrders","icon":"el-icon-picture-outline-round","id":"6f855990-13d0-43b1-87e6-f0647e99c360","idx":13,"name":"兑换订单","path":"/trade_center/order_manage/exchange_orders","title":"兑换订单"},{"component_name":"SaveOrders","icon":"el-icon-download","id":"ed0c5e11-eddb-48ff-84f9-6730d389f01a","idx":14,"name":"存款订单","path":"/trade_center/order_manage/save_orders","title":"存款订单"},{"component_name":"FetchOrders","icon":"el-icon-upload2","id":"a0e0bacb-c1d2-488f-aa5f-1af550bd97a8","idx":15,"name":"取款订单","path":"/trade_center/order_manage/fetch_orders","title":"取款订单"},{"component_name":"TransferOrders","icon":"el-icon-s-claim","id":"d141790b-1036-447d-b4cf-ae8b75de2190","idx":16,"name":"转账订单","path":"/trade_center/order_manage/transfer_orders","title":"转账订单"},{"component_name":"TradingAccountLog","icon":"el-icon-s-unfold","id":"2f4bda42-a624-4eef-9cc0-aa581733206a","idx":17,"name":"虚拟账户交易流水","path":"/trade_center/order_manage/trading_account_log","title":"虚拟账户交易流水"},{"component_name":"CollectionOrders","hidden":true,"icon":"el-icon-s-claim","id":"18d4ae82-68bc-435c-9e72-bfef6bc6b2e9","idx":18,"name":"收款订单","path":"/trade_center/order_manage/collection_orders","title":"收款订单"}],"component_name":"Content","icon":"el-icon-menu","id":"a8de69fd-5856-42cc-8812-fef457a99079","idx":1,"name":"交易管理","path":"/trade_center/order_manage","redirect":{"name":"交易明细"},"title":"订单管理"},{"children":[{"component_name":"DownloadFiles","icon":"el-icon-date","id":"29058c58-14a6-425d-b366-617860f6950f","idx":1,"name":"对账明细","path":"/trade_center/financial_manage/download_files","title":"对账明细"},{"component_name":"DownloadFiles","hidden":true,"icon":"el-icon-date","id":"d88a3357-b68c-406e-8951-97df93df66bb","idx":1,"name":"通道对账","path":"/trade_center/financial_manage/download_files","title":"通道对账"},{"component_name":"LiquidationList","icon":"el-icon-date","id":"46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d","idx":2,"name":"清算任务详情列表","path":"/trade_center/financial_manage/liquidation_list","title":"清算任务详情列表"},{"component_name":"LiquidationQueueList","hidden":true,"icon":"el-icon-date","id":"97a477c0-7361-421f-b042-2d1e0d9bf5fe","idx":3,"name":"清算任务队列列表","path":"/trade_center/financial_manage/liquidation_queue_list","title":"清算任务队列列表"},{"component_name":"ImmediateReconcilia","icon":"el-icon-date","id":"4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56","idx":3,"name":"即时对账","path":"/trade_center/financial_manage/immediate_reconcilia","title":"即时对账"},{"component_name":"LiquidationCheck","hidden":true,"icon":"el-icon-date","id":"a2bcc31f-515a-45e9-8083-b2f27d2c4afa","idx":4,"name":"查看清算任务详情","path":"/trade_center/financial_manage/liquidation_check","title":"查看清算任务详情"},{"component_name":"LiquidationAudit","icon":"el-icon-date","id":"60f259e8-d103-4fda-859d-b6c053cc6004","idx":7,"name":"清算任务审核列表","path":"/trade_center/financial_manage/liquidation_audit","title":"清算任务审核列表"},{"component_name":"MercReconcilia","icon":"el-icon-date","id":"b36e11e0-99e1-4af5-b83b-6fe9062554a8","idx":8,"name":"交易报表下载","path":"/trade_center/financial_manage/merc_reconcilia","title":"交易报表下载"},{"component_name":"ChannelSettleLogs","icon":"el-icon-date","id":"88c58e99-19a7-4f72-811e-92e21346207a","idx":10,"name":"通道对账分析","path":"/trade_center/financial_manage/channel_settle_logs","title":"通道对账分析"},{"component_name":"CheckDetails","icon":"el-icon-date","id":"59d15cbe-0f30-4057-ab75-09ba8776b7d6","idx":10,"name":"线上对账文件","path":"/trade_center/financial_manage/check_details","title":"线上对账文件"},{"component_name":"T1SettleReview","icon":"el-icon-date","id":"eedb53fb-3a3e-4f32-a232-ffc97a4464ae","idx":11,"name":"T1结算审核","path":"/trade_center/financial_manage/t1_settle_review","title":"T1结算审核"},{"component_name":"T1SettleDetails","hidden":true,"icon":"el-icon-date","id":"9b374efe-f7c1-44ff-a122-52981146b22a","idx":12,"name":"T1结算明细","path":"/trade_center/financial_manage/t1_settle_details","title":"T1结算明细"},{"component_name":"GenOfflineReports","icon":"el-icon-date","id":"8ec08232-29d6-418f-8560-64f307684f8e","idx":14,"name":"生成交易报表","path":"/trade_center/financial_manage/gen_offline_reports","title":"生成交易报表"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"cb332a03-aeee-435c-86de-b7d577e3817b","idx":2,"name":"日结算管理","path":"/trade_center/financial_manage","title":"日结算管理"},{"children":[{"component_name":"TestPrepay","icon":"el-icon-date","id":"4689bc27-a161-431b-bf58-456f66307fd9","idx":1,"name":"主扫测试","path":"/trade_center/pay_test/test_prepay","title":"主扫测试"},{"component_name":"TestPay","icon":"el-icon-date","id":"ef82cdff-1b9d-4396-8d2b-f7b751a50006","idx":2,"name":"被扫测试","path":"/trade_center/pay_test/test_pay","title":"被扫测试"},{"component_name":"TestAgentpay","icon":"el-icon-date","id":"357d5bae-8470-4d03-ab92-68cafedcfa16","idx":3,"name":"代付测试","path":"/trade_center/pay_test/test_agentpay","title":"代付测试"},{"component_name":"TestRefund","icon":"el-icon-date","id":"1897bc12-5b09-412d-a3c7-e001aedff81d","idx":4,"name":"退款测试","path":"/trade_center/pay_test/test_refund","title":"退款测试"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"065c9534-8c0a-408f-892a-914d172f90c7","idx":3,"name":"支付测试","path":"/trade_center/pay_test","title":"支付测试"},{"children":[{"component_name":"MonthlySettlementList","hidden":true,"icon":"el-icon-date","id":"db98b40b-91a4-4434-ba2c-712e3990b82a","idx":4,"name":"月结任务详情列表","path":"/trade_center/settlement/monthly_settlement_list","title":"月结任务详情列表"},{"component_name":"MonthlySettlementChannelList","hidden":true,"icon":"el-icon-date","id":"ad55ef62-54a6-41c0-bba9-830a06822659","idx":5,"name":"月结任务通道详情","path":"/trade_center/settlement/monthly_settlement_channel_list","title":"月结任务通道详情"},{"component_name":"MonthlyStatementPayment","icon":"el-icon-date","id":"e18c306e-1d61-4e6e-a2cf-c41ba357b406","idx":5,"name":"月结代付申请","path":"/trade_center/settlement/monthly_statement_payment","title":"月结代付申请"},{"component_name":"MonthlyLiquidationList","icon":"el-icon-date","id":"b995956c-75be-4635-9022-af1f5376be2b","idx":6,"name":"月清算任务列表","path":"/trade_center/settlement/monthly_liquidation_list","title":"月清算任务列表"},{"component_name":"MonthlyStatementFile","icon":"el-icon-date","id":"aa061751-3180-40dc-b101-f866d7e4dc2f","idx":6,"name":"月结文件","path":"/trade_center/settlement/monthly_statement_file","title":"月结文件"},{"component_name":"MonthlyLiquidationAudit","icon":"el-icon-date","id":"03ed060b-56dd-4ecb-9072-a736c81ae27f","idx":7,"name":"月清算任务审核","path":"/trade_center/settlement/monthly_liquidation_audit","title":"月清算任务审核"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"b4aee72a-2ac5-4748-b846-bd1828139858","idx":3,"name":"月结算管理","path":"/trade_center/settlement","redirect":{"name":"月结算"},"title":"月结算管理"},{"children":[{"component_name":"WalletAudit","icon":"el-icon-date","id":"742c55df-02ca-4d3e-909a-ac8c412ff735","idx":1,"name":"充值审核","path":"/trade_center/recharge_and_pay/wallet_audit","title":"充值审核"},{"component_name":"ImportRechargeFlow","icon":"el-icon-date","id":"48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5","idx":2,"name":"批量充值","path":"/trade_center/recharge_and_pay/import_recharge_flow","title":"批量充值"},{"component_name":"ContraryApplyAudit","icon":"el-icon-date","id":"90df4bc3-dbb0-4e92-ab9a-0aac895e5295","idx":3,"name":"代付申请审核","path":"/trade_center/recharge_and_pay/contrary_apply_audit","title":"代付申请审核"},{"component_name":"NoSettleAudit","icon":"el-icon-date","id":"01698813-a221-4b56-8bf5-ec3ec9956d97","idx":4,"name":"未结算审核","path":"/trade_center/recharge_and_pay/no_settle_audit","title":"未结算审核"},{"component_name":"ContraryApplyCheck","hidden":true,"icon":"el-icon-date","id":"c5e2efa4-c8c7-4c79-9823-59e829a4847e","idx":5,"name":"代付任务详情","path":"/trade_center/recharge_and_pay/contrary_apply_check","title":"代付任务详情"},{"component_name":"WithdrawalManage","icon":"el-icon-date","id":"0576302e-61ce-49b3-b265-28b7458e5c27","idx":14,"name":"出金管理","path":"/trade_center/recharge_and_pay/withdrawal_manage","title":"出金管理"},{"component_name":"WithdrawalManageEdit","hidden":true,"icon":"el-icon-date","id":"3c53a6aa-9d98-43be-bd14-4c5501e58ee5","idx":15,"name":"出金管理编辑","path":"/trade_center/recharge_and_pay/withdrawal_manage_edit","title":"出金管理编辑"},{"component_name":"WithdrawalReview","icon":"el-icon-date","id":"080f0752-21ba-4d3e-97c2-4846dfb4f648","idx":16,"name":"出金审核","path":"/trade_center/recharge_and_pay/withdrawal_review","title":"出金审核"},{"component_name":"SendApplyBalance","icon":"el-icon-date","id":"4e471ceb-766f-4575-ab73-856e07ceb1a8","idx":20,"name":"余额批量代付申请","path":"/trade_center/recharge_and_pay/send_apply_balance","title":"余额批量代付申请"},{"component_name":"ContraryApplyList","icon":"el-icon-date","id":"1c80833a-d9f9-4b31-ac77-55c282781167","idx":21,"name":"对公代付申请列表","path":"/trade_center/recharge_and_pay/contrary_apply_list","title":"对公代付申请列表"},{"component_name":"ContraryApplyEdit","hidden":true,"icon":"el-icon-date","id":"cab6e3e8-cff1-42da-959f-e12705adc760","idx":22,"name":"对公代付申请","path":"/trade_center/recharge_and_pay/contrary_apply_edit","title":"对公代付申请"},{"component_name":"UploadFee","icon":"el-icon-date","id":"3ed50a46-3620-43a1-b8ea-5a1ebc968851","idx":23,"name":"手续费上传","path":"/trade_center/recharge_and_pay/upload_fee","title":"手续费上传"},{"component_name":"FeeReview","icon":"el-icon-date","id":"23476194-48d6-447c-b4e9-2ab6ad7b1a48","idx":24,"name":"手续费审核","path":"/trade_center/recharge_and_pay/fee_review","title":"手续费审核"},{"component_name":"WithdrawalCheck","hidden":true,"icon":"el-icon-date","id":"82dbdde0-be98-462f-bef7-b86ed1a98336","idx":25,"name":"出金查看","path":"/trade_center/recharge_and_pay/withdrawal_check","title":"出金查看"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"e0be0b73-1f90-4c04-8da5-86513c148a16","idx":4,"name":"充值与代付","path":"/trade_center/recharge_and_pay","title":"充值与代付"},{"children":[{"component_name":"BankLimitList","icon":"el-icon-date","id":"a0dee210-7e5b-4309-ad4e-7ccdf1e7e82e","idx":1,"name":"银行限额列表","path":"/trade_center/risk_control_manage/bank_limit_list","title":"银行限额列表"},{"component_name":"BlacklistEdit","hidden":true,"icon":"el-icon-date","id":"72793652-a596-4b80-bbe4-73d10cca1379","idx":3,"name":"黑名单修改","path":"/trade_center/risk_control_manage/blacklist_edit","title":"黑名单修改"},{"component_name":"WhitelistEdit","hidden":true,"icon":"el-icon-date","id":"66cb5fb9-5bdc-4378-8010-6b58f8cedb5f","idx":4,"name":"白名单修改","path":"/trade_center/risk_control_manage/whitelist_edit","title":"白名单修改"},{"component_name":"WindControlEdit","hidden":true,"icon":"el-icon-date","id":"3293247e-3d3c-4749-ae91-1a991e27b8be","idx":5,"name":"风控规则修改","path":"/trade_center/risk_control_manage/wind_control_edit","title":"风控规则修改"},{"component_name":"RiskManage","icon":"el-icon-date","id":"5a8b13c6-3997-4846-b037-2b679f5c31e8","idx":10,"name":"风控管理","path":"/trade_center/risk_control_manage/risk_manage","title":"风控管理"},{"component_name":"ByWindList","icon":"el-icon-date","id":"7ea8ce89-9ce7-4a89-a332-7baad7bbe475","idx":11,"name":"被风控列表","path":"/trade_center/risk_control_manage/by_wind_list","title":"被风控列表"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"5a229d4e-beee-4f0a-ba2f-888f04a149e5","idx":5,"name":"风控管理","path":"/trade_center/risk_control_manage","title":"风控管理"},{"children":[{"component_name":"BankAbbrList","icon":"el-icon-date","id":"af0187d0-e84c-401f-9ffc-70194b099309","idx":1,"name":"银行缩写列表","path":"/trade_center/bank_params_manage/bank_abbr_list","title":"银行缩写列表"},{"component_name":"BankCodeList","icon":"el-icon-date","id":"c2ebd879-c09a-43d3-8088-6b0c2445c0b1","idx":2,"name":"银行编码列表","path":"/trade_center/bank_params_manage/bank_code_list","title":"银行编码列表"},{"component_name":"BankCardList","icon":"el-icon-date","id":"e296b896-76d1-4cf0-b092-e637a5909ec4","idx":3,"name":"银行卡列表","path":"/trade_center/bank_params_manage/bank_card_list","title":"银行卡列表"}],"component_name":"Content","hidden":true,"icon":"el-icon-menu","id":"77788f4d-5d0e-458b-b7cb-297e907b2eff","idx":6,"name":"银行参数管理","path":"/trade_center/bank_params_manage","title":"银行参数管理"}],"component_name":"Default","component_path":"trade_center","icon":"el-icon-menu","id":"136f2190-6f80-46d4-a124-d512e508be5e","idx":5,"name":"交易中心","path":"/trade_center","redirect":{"name":"交易管理"},"title":"交易中心"}]	771ae87bcd314d049bf66e9824af2dd6	2020-05-08 13:01:22.146429	172.69.33.208	2020-05-08 13:01:43	
5e2447c5-47d4-4589-ae58-09423dd57fd7		e04d6f613ee74e6c989714035fa57f29	2020-05-08 13:32:30.962352	61.140.127.79	\N	
\.


--
-- Data for Name: outgo_order; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.outgo_order (log_no, vaccount_no, amount, create_time, order_status, modify_time, balance_type, fees, servicer_no, op_acc_no, settle_hourly_log_no, settle_daily_log_no, rate, payment_type, is_count, withdraw_type, cancel_reason, risk_no, real_amount, op_acc_type, ip, lat, lng) FROM stdin;
\.


--
-- Data for Name: outgo_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.outgo_type (outgo_type, outgo_name, use_status, idx) FROM stdin;
\.


--
-- Data for Name: platform_config; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.platform_config (account_uid, top_menu_status, side_menu_status) FROM stdin;
\.


--
-- Data for Name: rela_acc_iden; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rela_acc_iden (account_no, account_type, iden_no) FROM stdin;
22222222-2222-2222-2222-222222222222	1	00000000-0000-0000-0000-000000000000
92c0e2a2-08b3-444d-bc4b-640e1f5836cf	1	00000000-0000-0000-0000-000000000000
4b21eb6b-96cf-439f-88c4-c023ddf6c4b1	1	00000000-0000-0000-0000-000000000000
d4a4ca0f-973e-484a-80b7-c40187aeda3f	4	e93642a0-2d80-44d4-a3ab-9e8e22e183e3
5e2447c5-47d4-4589-ae58-09423dd57fd7	4	401a2161-2cef-469c-ab05-4e581ea4dffb
af8ca79c-995e-4336-a1f7-31a76613d300	1	00000000-0000-0000-0000-000000000000
30c9911c-3e77-4e5b-9f2e-5d7825b378b4	1	00000000-0000-0000-0000-000000000000
d4a4ca0f-973e-484a-80b7-c40187aeda3f	3	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d
\.


--
-- Data for Name: rela_account_role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rela_account_role (rela_uid, account_uid, role_uid) FROM stdin;
c97d084d-9667-439e-bbe2-47a41550e7d2	22222222-2222-2222-2222-222222222222	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
61ee9bc7-2c9e-438d-86d6-93463fc738e5	af8ca79c-995e-4336-a1f7-31a76613d300	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
10aaf527-9e12-4273-9631-a572313973f6	92c0e2a2-08b3-444d-bc4b-640e1f5836cf	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
cb4d831f-f1dc-45eb-9170-5f5acb868fa3	4b21eb6b-96cf-439f-88c4-c023ddf6c4b1	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
314eefae-98a8-47d4-8fb1-ee4c06511923	30c9911c-3e77-4e5b-9f2e-5d7825b378b4	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
514881b7-f74a-45fd-bbc3-0947ec4e6e31	d4a4ca0f-973e-484a-80b7-c40187aeda3f	fbb11569-57f0-4b75-ae29-3854e26a0585
\.


--
-- Data for Name: rela_imei_pubkey; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rela_imei_pubkey (rela_no, imei, pub_key, create_time) FROM stdin;
0b9480f6-7543-4eb0-8b3f-fd113b0eb551	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC7rjZ0Klxqb23caySYbQdA/9DLz7xUL6UL9rV0q+Oqh8notFufm1ml1VGN4raIaudn63/JLwrddwlBhsSqeBFGRbdaeDXpefMBQmna7aTj5DXSlVRfrtlEvgoxUITe2xCjNzy86RIov58/Kqb2H2IhHo1QWb4optZwDKE7dN1VFwIDAQAB	2020-04-23 12:02:53.569829
759e3bd3-dd22-46ea-abc8-d0d792564f9c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDH/GnGs46RrE9yNCItQNRXzyfPw6bSjBe42nMG7Ggc2hn8AnAiIkKPOC4DDBqSBEIxtouC7hLXa/1nVdthGVQn3Satk9Jp3Upj7TDmBFmndbM5OCljyJbht7Rp1t1hM050sJGVldz9QYq9Va1vRPgsfHE84iyRKkE5L+e6ichgZQIDAQAB	2020-04-23 12:09:10.057066
fe5e709a-3149-427b-a55a-b75bdc6d073c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC7rjZ0Klxqb23caySYbQdA/9DLz7xUL6UL9rV0q+Oqh8notFufm1ml1VGN4raIaudn63/JLwrddwlBhsSqeBFGRbdaeDXpefMBQmna7aTj5DXSlVRfrtlEvgoxUITe2xCjNzy86RIov58/Kqb2H2IhHo1QWb4optZwDKE7dN1VFwIDAQAB	2020-04-23 12:09:42.861539
06d833c6-e3c6-4ea7-ab8a-c33617ae861b	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIlrANJtxVM4mYWyxe+CvhZyiGa8DNyb5HXWCR/fXlpMJxCEymDnbaoJ4s9NlerQGJBjVs/DRmkLzpQ366iERA9qC0OoR57+HQQ5cjI+qLI/WCAQAbhdphOIyx7Pp+MPn9FU3yV2LwZsS3+HlFJJYoD+XptW4kAjTCw/mMn0/q0QIDAQAB	2020-04-23 12:15:00.467256
c5d0a354-6a5e-4550-b4d0-3a112cabb337	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCEROnTSPeulFDlpII9s3LIetoiv8wxqsAxhvy9/rPI7klbZMbqZLTtiaMuHK7/IBNUE/YzUmfrLZbDeZ15Z1p7P+QxGeAzSSREfNMt26IYEedne12tQcKVdZy75d1zUfRmX4Ho2NWck0Qypvc5yVK1t4w7fvvhQriDX7E7O6QmUwIDAQAB	2020-04-23 12:17:00.296133
5f79fe56-95d8-4306-b97b-016830593e86	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCEROnTSPeulFDlpII9s3LIetoiv8wxqsAxhvy9/rPI7klbZMbqZLTtiaMuHK7/IBNUE/YzUmfrLZbDeZ15Z1p7P+QxGeAzSSREfNMt26IYEedne12tQcKVdZy75d1zUfRmX4Ho2NWck0Qypvc5yVK1t4w7fvvhQriDX7E7O6QmUwIDAQAB	2020-04-23 12:18:18.849679
b56c6848-0e33-4b27-ab67-3633522ba7be	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCEROnTSPeulFDlpII9s3LIetoiv8wxqsAxhvy9/rPI7klbZMbqZLTtiaMuHK7/IBNUE/YzUmfrLZbDeZ15Z1p7P+QxGeAzSSREfNMt26IYEedne12tQcKVdZy75d1zUfRmX4Ho2NWck0Qypvc5yVK1t4w7fvvhQriDX7E7O6QmUwIDAQAB	2020-04-23 12:22:36.497912
14449ffd-54ea-46a6-b2cf-11853141164b	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:22:42.636774
275546f1-4ecd-4f0e-996c-d1ae5b08c24d	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:22:54.07504
b09cd508-5291-44a9-a76f-36fa03d99e6c	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:22:57.035909
bb0ab580-9ec1-460e-9c2b-7b6ba00d6bd7	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:03.808562
1a3f2bc6-46e7-4671-981b-fd57e7a4d8fc	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:05.445374
603b26d4-0816-4154-a8e9-778d5ca8821d	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:19.328775
fb119b0c-bf0f-4d71-a163-96898341bdce	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:21.352344
e251d0d4-ffa4-47e0-b6d5-861095772b7b	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:23.159574
b5262134-948c-418b-97ed-194fd3f8075e	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:24.960571
9de9f73e-41a3-4880-b6f7-a51dc09a95ab	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-23 12:23:27.466645
a4bb9c5c-cefd-4d71-bfe5-22d6906bee6c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCUuq8c6XeuyqGQamjLIjyO5DuZBKk5IHrQP9DRyeJlcunW21I5fZ6GzcLJpyPG5x1oqnIYqLL0goTr+WmN8Vsu3xsycEByH1LxtPlZMn3AsAC7CDDbMLGI2zLGGnD8/YqjcG12Nmrdwv6TovIiGqEljA9JmvWYhEptK19WMiF11wIDAQAB	2020-04-23 12:23:33.599207
77030360-e411-4869-b8ad-35a4292e05f3	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCz8o7jiIVDq6VNt3d0Tc1KGiY+2OQXUssnJ16x0d5AAydi+UMKdFlqZXn3FfDeVzefoBPDEfw1XGqNSoRZ8X/mCPuB5Nju+0F/BWBzZ/Ytx8e0YzDlBHxijqIVm84EHnVI5GdYomMmLrJYPMlGM0ZI/kvdKYifLa+1TW0+PMMVlwIDAQAB	2020-04-23 14:42:33.94844
71c0f0e1-adec-4293-a909-0c4ffe52691e	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDaZHGqJbtMumU7E+jTM6rB9E0yDv16EmCI7EmIDpiLpbfxck4sVvy2aR7PHOEeS8CS/LZtxvMPofSUAk6b6lCibJ1mX/hkoLrdFri+HnG18SlvhW4biXMTbMjwCowfFfpplJpRzjTlLF74DGmbX8j8IlGIgws1UZ6Q8plW3ixK+QIDAQAB	2020-04-24 02:34:18.857467
0c8a349f-626e-4a01-b326-075a19963aa1	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCH4xCHOSpucRxhZv1G8CqMkEbczVVwcOTqpG86ZywGggvd80W9dTcoUeq8lXiCWQYc77u4IzfElVEO5zd5b9ffCt3GpNI//KgaL4R/Y8dbakDkLArX4nrMVQBbPMfaHNCk1Tjzq9fleD0lTzWLCdWxvEC1icXEb/PixzvP8Ja8MQIDAQAB	2020-04-24 02:40:37.580361
c54749c1-ed7e-4b98-a4db-921d9ec767cc	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCUuq8c6XeuyqGQamjLIjyO5DuZBKk5IHrQP9DRyeJlcunW21I5fZ6GzcLJpyPG5x1oqnIYqLL0goTr+WmN8Vsu3xsycEByH1LxtPlZMn3AsAC7CDDbMLGI2zLGGnD8/YqjcG12Nmrdwv6TovIiGqEljA9JmvWYhEptK19WMiF11wIDAQAB	2020-04-24 02:47:08.005844
d50c40b5-925d-40fd-88e7-986bfa125b39	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCUuq8c6XeuyqGQamjLIjyO5DuZBKk5IHrQP9DRyeJlcunW21I5fZ6GzcLJpyPG5x1oqnIYqLL0goTr+WmN8Vsu3xsycEByH1LxtPlZMn3AsAC7CDDbMLGI2zLGGnD8/YqjcG12Nmrdwv6TovIiGqEljA9JmvWYhEptK19WMiF11wIDAQAB	2020-04-24 03:00:01.440174
b284cf20-5185-4db5-bd2c-50b22dce5238	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-24 03:01:33.758332
b1384eeb-536d-40fb-bb32-6753e5d118d2	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-24 03:39:48.438277
f59b15c8-618c-4ee6-884b-4f34ae97b9ac	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-24 03:39:51.467267
c5209b06-4d2d-462a-abcb-55c421d05985	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-24 03:39:57.508096
4dfdcfc2-a286-4534-b87e-b0b08da44fb2	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCn90px6F7pLD6VKTqdUXpGQPTgpg4IKFuS+TPdesTEAOmM0osndDOifzNa82ZqNVle6fQgLsjH84se2LRM4CjTF0CivXDBZdfhw19ouqHxVy+lmHm31AXdAcyZd2V24aH4J7GvKROsyz2/DLXtYgkrJYPrqZxdZ32DhCAdU/9ZEQIDAQAB	2020-04-24 03:39:59.553253
cd357775-b15e-4bb6-877d-0215187c0f49	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:41:20.52578
b120daf9-e956-4e6a-b071-e6f4e38e09f6	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:44:25.065745
1d995afd-8a9f-49f8-a02c-b4c5e8934bc6	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:44:26.780841
38573b0e-d19b-4796-9a94-5c60dd94c811	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:46:35.055763
e22c8761-e80b-4f4a-8511-bd576c5abfa9	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:48:47.189545
c4c6b2c7-ee7c-4957-9b37-bfa9f83ef980	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:56:16.648042
43c196ce-ceb6-404d-b6a9-8b56a4139c86	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:56:39.62824
838871ff-ff7a-48a1-8ba3-2f54eeef214b	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:57:03.360012
7d4aafcb-298b-4916-a4a2-0fe4a2df247d	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIjuybaDaxF//CeHYQ+DKIi+rusy2a5porQ9T9lahB3qIMxyF30hlWtymTT4Qz6feFFvq5akSzhmTCXRfAvssTijIqArq5M/I/ZV6E8MR3x6SxZEV63dZ8WRtigd74feMlRZM6lKmZCd+atw08AOVg+Gc0DIu1/lpE3Cer9XWkwIDAQAB	2020-04-24 03:59:30.105245
a4ace7e6-e2ac-4a90-a779-8c9f35f28716	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCWlJrebQOlgzDqpDXaNtYpMEY8JgKqXRCEqKbTh0+Hk6b0hv732gTogi/c4RJAYfDwHacewaNMt18EhP+IEoqy2wlNXsR6qWlm0166FwPWXQxklb6SqSl5DibHjSPTtwaTXi6/r+TTOsF3zIxWzPG+uZ55hOl/GKjdhSZ6jBSCYQIDAQAB	2020-04-24 04:00:18.980187
ac67bdc1-752c-429b-bebd-9643bc20c06d	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCWlJrebQOlgzDqpDXaNtYpMEY8JgKqXRCEqKbTh0+Hk6b0hv732gTogi/c4RJAYfDwHacewaNMt18EhP+IEoqy2wlNXsR6qWlm0166FwPWXQxklb6SqSl5DibHjSPTtwaTXi6/r+TTOsF3zIxWzPG+uZ55hOl/GKjdhSZ6jBSCYQIDAQAB	2020-04-24 04:01:39.463397
712f78fa-4528-43b5-88bf-21b18e2b04d5	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK7cATW2VBJ/oLoIDWmCgYi0YQSkq3gHIrJk+viCKDCFeg5UDENQ/1QxGIMqw9lnEzNQorKu4GOHABE/NbzSwNzONEHTXURKfrytM7LTho5K/hBjaUl8GFbo2c+YAN/UcZFCTwJCcoEOzCmCR9m0+5J7sJYjutMSpBjTjp9Yj02QIDAQAB	2020-04-24 04:28:25.329874
a589e5b3-d488-43f5-9a72-e475c9aaeda1	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 04:29:21.734482
edd3d6f1-c2d1-4a13-900a-f8b9307f5d39	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdB3Wj7XRq4ANqauU/MJ/9T8uuUvDQ9a5WpP83+BKoFS4X2sapjIReZ2RfIoJCSVxoqMJV6myHkbsj7EbyBBvDGeibYwdQyjUFnbm5shsx9bqOxGW+iwkD+edxVpS4uSSL+zDUGRK3hh6JwNjck9qeIfJUs+QbqBURY3wITYfJnwIDAQAB	2020-04-24 05:48:17.915714
2fa1b5ae-3bcf-4031-9f91-278394e40f6e	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 05:49:17.824844
91bbceb9-5bfc-4c83-9a13-4161110e223f	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 05:58:28.618299
c826d19e-078f-4dbd-86fb-a98c49ff2cb5	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 05:58:58.833701
4857767d-12eb-44a9-b999-b21aad52afe6	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdB3Wj7XRq4ANqauU/MJ/9T8uuUvDQ9a5WpP83+BKoFS4X2sapjIReZ2RfIoJCSVxoqMJV6myHkbsj7EbyBBvDGeibYwdQyjUFnbm5shsx9bqOxGW+iwkD+edxVpS4uSSL+zDUGRK3hh6JwNjck9qeIfJUs+QbqBURY3wITYfJnwIDAQAB	2020-04-24 05:59:16.546812
da538539-7d2c-4b25-bf80-7c9e92a95cdb	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 06:00:46.549819
33302c29-73d9-4d2d-b3e5-6a11ff6425b8	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdB3Wj7XRq4ANqauU/MJ/9T8uuUvDQ9a5WpP83+BKoFS4X2sapjIReZ2RfIoJCSVxoqMJV6myHkbsj7EbyBBvDGeibYwdQyjUFnbm5shsx9bqOxGW+iwkD+edxVpS4uSSL+zDUGRK3hh6JwNjck9qeIfJUs+QbqBURY3wITYfJnwIDAQAB	2020-04-24 06:03:15.459148
69a1d60a-f26f-46bc-86ab-7d79192fae57	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 08:25:30.019148
09d44745-1bf9-4f05-9061-a13fb480d901	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-24 08:25:34.076146
b90329b6-55df-43cf-bebb-aa916cc7f795	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZXp6QIrEkrEISd8fT7fd6zKgqll7Bq9bA64iMLAGkAljaFfg8Kj/LEYOKXWRU3oGLSS31ZCH0nZFgHjlaqzoreTh9Z4thw+vR7LPL+BVEbiBSp/DQ/3ObBHMqNcEVLOy5Yjexs7KS92U5KOJQOA0xWAcGvUdrOi0XQVehEyPaUQIDAQAB	2020-04-24 16:02:08.149646
7a686c35-61e8-4bde-8555-a08de0ff36bf	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCi3IGbk/x2EXYJej8zF6R2KDW4SPU7yxCLnLLxC+yk1sHOJ9ufzSxrs79T2E+cp+eLKYLUKfsbTkn8tLCN8IjR2NHpKe+4gStoY5KvqEPmtv3N6SuzCPLGE3hLMbnViScZAaND/yO8dam9KtGXxNjXJLXd736itYdyoB6vwLNGwIDAQAB	2020-04-24 16:09:45.096977
e0aa6ba5-0782-4847-b124-9ba93cc165aa	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCKI8Q1t4yLxIuHj7t9c2m/oK2xV98AIAP+Mka0ucdfk7T5esH9kMIcQoXlbyxxy4uE6blyiGkGmlgNG6WvlS0SCQ2xCbdQ09ZAxBogHOv/Onhs4lg+bho9nD/mA0ygaQHViGNE6uTVPRKr6W4WusGQcaqQD+jBUINk76oCkXYfgQIDAQAB	2020-04-24 16:35:17.942681
873291f0-2682-40db-b453-f0d7f120b941	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5Kt+uWX7hW+4r7kVVKMfhqXiinUrxpo+U9zArfnfdKJX9K9C4UCVmmB4q2BWRSAact81L9JNilcxYQB6TO+GG88FErbNtghTLSw6WnuXChyIEmZjwroAvjxDjzgkSfC5pDv3CoXNj/AhfnTI2TWp/anxZFdg8jqsdbAojuVtPIQIDAQAB	2020-04-24 16:53:14.596658
b54a7224-3d49-4327-a846-136f848d202c	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5Kt+uWX7hW+4r7kVVKMfhqXiinUrxpo+U9zArfnfdKJX9K9C4UCVmmB4q2BWRSAact81L9JNilcxYQB6TO+GG88FErbNtghTLSw6WnuXChyIEmZjwroAvjxDjzgkSfC5pDv3CoXNj/AhfnTI2TWp/anxZFdg8jqsdbAojuVtPIQIDAQAB	2020-04-24 17:02:12.872989
7a0ee700-2110-4d4a-bd78-37b99ecbd5d6	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5Kt+uWX7hW+4r7kVVKMfhqXiinUrxpo+U9zArfnfdKJX9K9C4UCVmmB4q2BWRSAact81L9JNilcxYQB6TO+GG88FErbNtghTLSw6WnuXChyIEmZjwroAvjxDjzgkSfC5pDv3CoXNj/AhfnTI2TWp/anxZFdg8jqsdbAojuVtPIQIDAQAB	2020-04-24 17:03:25.893592
eb38bf35-feb9-4976-b863-151931d99ed6	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5Kt+uWX7hW+4r7kVVKMfhqXiinUrxpo+U9zArfnfdKJX9K9C4UCVmmB4q2BWRSAact81L9JNilcxYQB6TO+GG88FErbNtghTLSw6WnuXChyIEmZjwroAvjxDjzgkSfC5pDv3CoXNj/AhfnTI2TWp/anxZFdg8jqsdbAojuVtPIQIDAQAB	2020-04-24 17:04:56.477433
18a374d5-bec6-45ac-bc58-eaec5b6cd78c	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5Kt+uWX7hW+4r7kVVKMfhqXiinUrxpo+U9zArfnfdKJX9K9C4UCVmmB4q2BWRSAact81L9JNilcxYQB6TO+GG88FErbNtghTLSw6WnuXChyIEmZjwroAvjxDjzgkSfC5pDv3CoXNj/AhfnTI2TWp/anxZFdg8jqsdbAojuVtPIQIDAQAB	2020-04-24 17:07:29.857589
5981e230-6e82-4974-83f6-670fbb390771	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK7cATW2VBJ/oLoIDWmCgYi0YQSkq3gHIrJk+viCKDCFeg5UDENQ/1QxGIMqw9lnEzNQorKu4GOHABE/NbzSwNzONEHTXURKfrytM7LTho5K/hBjaUl8GFbo2c+YAN/UcZFCTwJCcoEOzCmCR9m0+5J7sJYjutMSpBjTjp9Yj02QIDAQAB	2020-04-24 17:48:30.906834
a3953f75-df52-4457-acb4-e88eabc06f52	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC7rjZ0Klxqb23caySYbQdA/9DLz7xUL6UL9rV0q+Oqh8notFufm1ml1VGN4raIaudn63/JLwrddwlBhsSqeBFGRbdaeDXpefMBQmna7aTj5DXSlVRfrtlEvgoxUITe2xCjNzy86RIov58/Kqb2H2IhHo1QWb4optZwDKE7dN1VFwIDAQAB	2020-04-24 17:52:31.864387
0abe1601-1f9f-4d69-8fb5-29900e00495b	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK7cATW2VBJ/oLoIDWmCgYi0YQSkq3gHIrJk+viCKDCFeg5UDENQ/1QxGIMqw9lnEzNQorKu4GOHABE/NbzSwNzONEHTXURKfrytM7LTho5K/hBjaUl8GFbo2c+YAN/UcZFCTwJCcoEOzCmCR9m0+5J7sJYjutMSpBjTjp9Yj02QIDAQAB	2020-04-24 17:59:00.559892
19ca6c1f-92c2-47f2-a959-db24f419543c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC7uk86nofZ1G7lQg0sxykCvoJoMG5DbXTPyBk3ZcEHm+J8c8Um8pwbAMPHP7uH8HEvXMhrBitZcoVgvwjBpG9UNgxSzvVkXyUK8nWvcoTMyKYoFsZeZ/lnoDTdUQz46p4r5pekA35Y0vHc5qtQmZaN4ph2emRsUep0FgzUh/rd8wIDAQAB	2020-04-25 09:10:03.130438
287ee9cb-08ec-4edb-b8b0-874c3c4e1667	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSF659z9FxBShiOt2JWK1HXFwqE+X/Dx1DWibpIRdzJ2gtMQAhGbtLPWplE1ehMZLq9eVDGonbWbgAe5fyFqunU8xObmAjSWXIJHWO6zn4me9vqvaKYLRAxqjiw+GrMBHaahPkRiqo/DSTNmnyM5vOnD4r9CTmyo2/BQa9PJiBkwIDAQAB	2020-04-26 09:38:12.490085
b3bb42fc-6d2f-48ee-8918-ad495c2d6fe9	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdIn7dWVKoUz+b+HXzPnGBwsRrk64xaBaTN7shtolcb71kPxE0USyShyU21Et5n4zPi1Zqhuw/TF1pGe1PFKYdAPEN+R3/hdL8i7YuXYIXu4yWYr31W3o1PUHuTcAIHXGrTWRGoK/UV3AtXrh6JXDjnvvEwLYrILhOyxF/nhBliwIDAQAB	2020-04-26 09:43:33.791603
d6a9f0d1-6802-431c-b73b-57c412196bd4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdIn7dWVKoUz+b+HXzPnGBwsRrk64xaBaTN7shtolcb71kPxE0USyShyU21Et5n4zPi1Zqhuw/TF1pGe1PFKYdAPEN+R3/hdL8i7YuXYIXu4yWYr31W3o1PUHuTcAIHXGrTWRGoK/UV3AtXrh6JXDjnvvEwLYrILhOyxF/nhBliwIDAQAB	2020-04-26 09:46:37.333355
ebdeb0c1-541e-4005-8b48-ee8fd159f417	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCK7cATW2VBJ/oLoIDWmCgYi0YQSkq3gHIrJk+viCKDCFeg5UDENQ/1QxGIMqw9lnEzNQorKu4GOHABE/NbzSwNzONEHTXURKfrytM7LTho5K/hBjaUl8GFbo2c+YAN/UcZFCTwJCcoEOzCmCR9m0+5J7sJYjutMSpBjTjp9Yj02QIDAQAB	2020-04-26 13:39:10.191496
fd2e23ae-18d9-490e-915a-f23cffbf8a81	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCfccqM1NqGXPaPCWcK+xCJUQ+bBCMX/Fn5EZAsC19OeqDJqHfSMsTAZUGLcn/zeRrFDoDhzVI3/w4SUUS1xGHWxJ7VuI8vYCpNY61cCIvKv0pt6OnK+L3occr9scyxO0LKdyLTqgcX6vBnZgtEWCN0NIIaWn6rdvnsgldzljFh2wIDAQAB	2020-04-26 14:02:36.558239
85436cf9-58a2-4dbd-8c06-e4ac4d42f12e	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+Wt37uYlK9hf7n/1BqW9mqE2aEfumY9KaSijcDt9Eic6Zg8Q8zwJ2ELW7OkdGPU3vQLtJjyJT/GQu7yb3Ema/c3H4hUTxyZPR8w8o9pw2EPzM48utOYk3nbJ5Z2EssjOAkYSkf3Rkef6waFzMOm6LzPc0VHLe6ocIZbrt2BIDUQIDAQAB	2020-04-26 14:07:22.388187
fa3be783-19b5-4a4a-8719-45bac46534a3	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCCehacywOHN8+LnxZ8tK6XVgnNAIp0NHvBGajGxsXdBhAof13skWMWjIObsKVI01y/Ss2osnsLm4mfgKkKfkQkYwBwQVWFFdqmgaOQA5Z1Uh7p9uc51m1RIaBnpz5M6hkh2KZc22gsLjqjsGvyc9oqYFr1BZCEwUN1ohPkh1UeiQIDAQAB	2020-04-26 17:07:06.335619
bed378a5-7c29-42f9-880f-858cb5449334	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-26 17:09:03.764401
c3631b07-7583-4353-a31d-6fb90f7b718a	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCUG99ECH9r6l570ezFBDN4Z9tlHQedfLQwtdnwDGpPvaXgNZqPNFIChnzQGr7QmqCFfiwIEuAyIF1gQ67eIu2LWdVMV+qr0gI4E76koITCFiUQ/w4u3uajRrr/3hRSkpPLJZauhvCL2MFOlfjOQerdn5LF4n/LTzrD+OYH8DhXwwIDAQAB	2020-04-26 17:39:28.246871
76a7f90a-a466-4792-a71b-5fdfbbdaded7	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnCiyp/VPjqFPgTms64F/TqTrPfPvTHfa5bmvyucVk9ySnOvrliwEq9Up8Sx45lE7yMmY7NgsY1Fs+qHC1JTPwxTK5hL/PH5qSKsVKyg/w6Y8Ji0oQV3G05Yttjy3j2KAIlEeHB4hmNnvsX/Hvhm6wgaBeXAhfF8cwsQpDVZ/M4QIDAQAB	2020-04-26 17:44:34.262673
d3b2eaf9-298a-4d35-a7f6-26dcd7afb809	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnCiyp/VPjqFPgTms64F/TqTrPfPvTHfa5bmvyucVk9ySnOvrliwEq9Up8Sx45lE7yMmY7NgsY1Fs+qHC1JTPwxTK5hL/PH5qSKsVKyg/w6Y8Ji0oQV3G05Yttjy3j2KAIlEeHB4hmNnvsX/Hvhm6wgaBeXAhfF8cwsQpDVZ/M4QIDAQAB	2020-04-26 17:46:10.8622
2d931df0-90cd-488f-8f8f-cb596e803f51	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCq8yigU6sk7LdIbky2uNPMaS3YEoAPPZ+pEsU/h9IvZ8se4tZIHpebWduR7T1SGslA9ZGWlSq+bSscyYEdn9SXADS6LAsEemXWDjg8zzGfiucrEZKa7TL+pHgl6rq8wpezR15rMrGCyQ+DkZbvn3dbXgd2BCtpzc1pgrW7ieH0YwIDAQAB	2020-04-26 17:56:10.284845
d950e892-d8d4-40fa-92fb-58fd5f08c672	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCq8yigU6sk7LdIbky2uNPMaS3YEoAPPZ+pEsU/h9IvZ8se4tZIHpebWduR7T1SGslA9ZGWlSq+bSscyYEdn9SXADS6LAsEemXWDjg8zzGfiucrEZKa7TL+pHgl6rq8wpezR15rMrGCyQ+DkZbvn3dbXgd2BCtpzc1pgrW7ieH0YwIDAQAB	2020-04-26 17:58:11.621456
071aba4c-ca9d-4f1e-b016-089ea04ea062	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIfj8knVpdZokRQrT2S9zvEWMpnVCBxkjlk8jYrkDlodPpgkiOX/z382D8WwHBK0n6opftfKLzmW6Iq0Z6TCwraUMyvCh5q/IB+bIY3YYmLPCYPCpkEu9ZNr7IOcrtkpsKS0gS/tqnfb8ybhGLSl68zvlM9BC4MQw/SCG+WdJ1CQIDAQAB	2020-04-24 18:57:45.094086
351a9ccd-34d9-4bda-9690-9462f761f013	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZxnsroqGyjY0ITQ7i3BF6F9LfSE3GrpFSV107O42B8k/SQ7LQok2rEqiNQOVA2GrBDbvgU6zjlY18THNAyVzsrjFle0DYgZ/4T2oqc13v139jcsSeKSDA1U74ZNM8a1Mxxdrw24BkfZe8l/4FjiLGBQYR1ad/WmfcBdlUFodhOQIDAQAB	2020-04-24 19:01:25.458552
73661433-3000-42fc-88a7-203475dd9a5d	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCi3IGbk/x2EXYJej8zF6R2KDW4SPU7yxCLnLLxC+yk1sHOJ9ufzSxrs79T2E+cp+eLKYLUKfsbTkn8tLCN8IjR2NHpKe+4gStoY5KvqEPmtv3N6SuzCPLGE3hLMbnViScZAaND/yO8dam9KtGXxNjXJLXd736itYdyoB6vwLNGwIDAQAB	2020-04-24 20:58:22.964455
fc1a8c74-4aa4-4a75-a089-ab5f3fa9b2f7	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDeppmHgTtERIeZeJDERxpBHkaQvICGg847bhfO5wdOe7lCJOHEo8X/JYjU2hzO5FbxkBOWVCr9lFAy+pGs4RSKdc2oSPgBu3WKKrkRXGWlU1raEvu2rwixu6ZvYs0+mijGZ8OavG18cCopQxSPBUGPd8JtJjDm6HfwrMC3u0mwrQIDAQAB	2020-04-26 09:13:38.989103
640eb18b-e37b-4304-95a1-e1558f9d6dbb	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZxnsroqGyjY0ITQ7i3BF6F9LfSE3GrpFSV107O42B8k/SQ7LQok2rEqiNQOVA2GrBDbvgU6zjlY18THNAyVzsrjFle0DYgZ/4T2oqc13v139jcsSeKSDA1U74ZNM8a1Mxxdrw24BkfZe8l/4FjiLGBQYR1ad/WmfcBdlUFodhOQIDAQAB	2020-04-26 09:16:26.741371
774b008e-0e1c-4ff0-9ddf-79c075f69297	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCAFdQl3PcbdI+07LOL1t4L+xejpZyLIfOoaoZIuP/yWVI5Zv/QNPCclYBgfNmKr/c+I8Sh/r11D0DnzQuX7uqiY7b/xK0bu07W4joTDEP2DKCecbG4QzC3FUReVHAASttf9jIh0NYwVzwr/Kt+oh/Prm8JWvSu5T3AOyXNOXLrDQIDAQAB	2020-04-26 12:37:25.585794
8fb96681-7828-40f5-aa33-b7ccb3595531	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+Wt37uYlK9hf7n/1BqW9mqE2aEfumY9KaSijcDt9Eic6Zg8Q8zwJ2ELW7OkdGPU3vQLtJjyJT/GQu7yb3Ema/c3H4hUTxyZPR8w8o9pw2EPzM48utOYk3nbJ5Z2EssjOAkYSkf3Rkef6waFzMOm6LzPc0VHLe6ocIZbrt2BIDUQIDAQAB	2020-04-26 13:00:44.298656
b616692f-6131-4137-acda-abe0b54da55b	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCWIOH1vXavJ+ZCoxqXYMeJq4o1RXQunYGOKT6c+0X6rY6N7gwsIX5cqI96lUjhP4X4EPQdd1xRStJWo/DUVLKrvXy9RpRPVJWB1GxeL9944RCEA6aT6cOF7MmI5hBGdWlisWqkuC2ZbWqo84f1JX/++mxZL4WWfR09CW8bQ5rKtQIDAQAB	2020-04-26 13:18:20.921083
bd5414bd-61cf-45d2-af2d-0a4b9c9adb94	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-26 15:24:03.88317
d5c5f7f6-5f51-4eac-84ab-6889731c0201	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCeJrfirVW0oMUJphGhTLI4kFogIt0ROrp5WqodDtkP+0YbZuie5h4LDNnWs1Dkn70GPY/l7jzJNuARUmPdG2O22XzldUYp273SMyc63EMqsfbVMT2M5j8IvxmF5fOb8BnGA5d/p6lohGZ+88XhwSkNRLGgxTaynWhxYpKzwpfgKwIDAQAB	2020-04-26 15:24:39.992573
c7db07fd-d67a-4147-ac0b-7f5f9480e135	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5h6h5hHaewIlfXO/iLYlQfepKu/ieLqZ3jdBqtXCqRerOkplTfU8OlInZNbQjde8kTDWeBAVwoM3Oidv5w+QwgV6lCMpEgCWx+WU16/V1Q2+SivyYEZnXYhXo4VifwE7wVlsUvKojF8GWF6rVcoR+Yn/C+4SXq8JZvpCrQt8piQIDAQAB	2020-04-26 15:54:26.983888
fae6a922-565d-4718-8e0e-c110f870e12a	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCLINOxVfA5RJrpybhG3VIe1cf1W3nrPwBLUTuEBut+xekF41+91JaxyskZivmvGSdvQcqIQ8+7KPLmPYJlJ5w4lisDXxIIANu5Jv5uI9COAGwtQGH/2tzx0WD4L//yPOvHTlPojl16P87eqI1pymDEySWIWm2nnN4jJyttSS2oCwIDAQAB	2020-04-26 15:56:38.521714
b0171bb7-df46-4e8d-aa74-f9c6095de7f9	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-26 16:03:05.996862
71366c45-c318-444d-80c4-1ba95f1150bc	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC2dN2gYQSl9SykDGduTiQv2IR95h1KZv+GaKctjgQrDRwHn2OvJGRs5eRiR1Wc0fxxFeXGy2hjGSZoMKTSyPoN8fhG+kMMIP7zwymycdvbObAD/ps8JtbvYBi4t1+0MKZB/gxq78IKMZjbn1gjbmGsTGbwbFCYTuTrdMRavjx5ywIDAQAB	2020-04-26 16:18:07.388536
6fd420f6-4cf9-44f9-8635-1a2946f6cc15	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDWHvxVfpHXJ4LOaUOOy+IddK6Z5ofHkH+PbdlumAM2Pj/CelwtsKXrn3CSb/kApsWVouotrXlUAWVhzNP0ICyDCv/0Ew1imBYvNUAGORUZvlEYqMpeMyJwpNDq7sqR2maPDcKixoSJsnC5ZkwwyuKQluFRvjtYP2dqz/5eYr/GzwIDAQAB	2020-04-26 16:19:19.883206
f16c67fa-f95e-47c9-8ce4-28ba3c73fa29	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCCehacywOHN8+LnxZ8tK6XVgnNAIp0NHvBGajGxsXdBhAof13skWMWjIObsKVI01y/Ss2osnsLm4mfgKkKfkQkYwBwQVWFFdqmgaOQA5Z1Uh7p9uc51m1RIaBnpz5M6hkh2KZc22gsLjqjsGvyc9oqYFr1BZCEwUN1ohPkh1UeiQIDAQAB	2020-04-26 16:20:59.875483
f71ba68c-594d-409d-9a52-11702cec5078	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDWHvxVfpHXJ4LOaUOOy+IddK6Z5ofHkH+PbdlumAM2Pj/CelwtsKXrn3CSb/kApsWVouotrXlUAWVhzNP0ICyDCv/0Ew1imBYvNUAGORUZvlEYqMpeMyJwpNDq7sqR2maPDcKixoSJsnC5ZkwwyuKQluFRvjtYP2dqz/5eYr/GzwIDAQAB	2020-04-26 16:21:45.619976
2924b634-fd99-48ce-bcdd-9ce124f54c57	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSWU/qCNIFawlOOVUO33d+jXPz0oZJcbTdthJaKFkjG8vpHTeAQIdGp3GM/zuXOxEDIpUjUx2qiPKK1X8j168CdZZ3wZcfwWs1+wr1pGRSvc37+gAUzlTO6bqNyrIl2OQvArWJ2V779IfantyJVrSoxAyMZ9SYQu33tmtHPG5HbwIDAQAB	2020-04-26 16:30:53.447266
017a0479-94cb-49bb-90f9-0f70bdcef7b4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSWU/qCNIFawlOOVUO33d+jXPz0oZJcbTdthJaKFkjG8vpHTeAQIdGp3GM/zuXOxEDIpUjUx2qiPKK1X8j168CdZZ3wZcfwWs1+wr1pGRSvc37+gAUzlTO6bqNyrIl2OQvArWJ2V779IfantyJVrSoxAyMZ9SYQu33tmtHPG5HbwIDAQAB	2020-04-26 18:17:24.706587
46b62fbe-6469-41b7-8ff5-9a69a56915f0	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-26 18:20:09.100206
37f51975-4eaf-459a-8968-0cc78a05e54e	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSWU/qCNIFawlOOVUO33d+jXPz0oZJcbTdthJaKFkjG8vpHTeAQIdGp3GM/zuXOxEDIpUjUx2qiPKK1X8j168CdZZ3wZcfwWs1+wr1pGRSvc37+gAUzlTO6bqNyrIl2OQvArWJ2V779IfantyJVrSoxAyMZ9SYQu33tmtHPG5HbwIDAQAB	2020-04-26 18:20:43.560438
db109f48-c065-4650-8149-80614e49dd53	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCC3RL08WBh0pMS4weonleEv49iMKtNZU6O/fNCdoTarXXlaOfJQhY/Sjnxvkctw3Jnr9TxhpDNOtKKBdU9yxsObExEQ8FpnSrEISfNRFichCBTyckZfyXMlVH+TcWzHZ5/LtLHYgS6wgJTpgVGSRzjsDxZ6wq0pnzcG6TsIKVCPwIDAQAB	2020-04-26 18:40:04.982061
37743e28-3937-4767-b23d-eba0200ad3eb	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDauQ2YwFPWKip5HJE3eMhhxTIH6/c6++bYqCg9nzg57dvVUwO/gZI3VtM63XH4DkN9VqB8Tm3Mum1sukyicNrExt7LQjtTLABl3HYWBmzkNCLXp76ZKUDMz+fonTpUFd0qAkMx0MSj6Tz13pnhs01KBrW/wADUDpWDaxuKTWih4QIDAQAB	2020-04-26 20:56:50.846696
ae4a3305-f048-44a4-9b2b-1947e5db5acb	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDauQ2YwFPWKip5HJE3eMhhxTIH6/c6++bYqCg9nzg57dvVUwO/gZI3VtM63XH4DkN9VqB8Tm3Mum1sukyicNrExt7LQjtTLABl3HYWBmzkNCLXp76ZKUDMz+fonTpUFd0qAkMx0MSj6Tz13pnhs01KBrW/wADUDpWDaxuKTWih4QIDAQAB	2020-04-26 20:57:11.273576
f3805f59-ded0-4952-9131-ff8d9d47999d	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDauQ2YwFPWKip5HJE3eMhhxTIH6/c6++bYqCg9nzg57dvVUwO/gZI3VtM63XH4DkN9VqB8Tm3Mum1sukyicNrExt7LQjtTLABl3HYWBmzkNCLXp76ZKUDMz+fonTpUFd0qAkMx0MSj6Tz13pnhs01KBrW/wADUDpWDaxuKTWih4QIDAQAB	2020-04-26 20:58:00.862853
2bca951b-46ce-40f0-b85d-fa4c86e3721a	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-27 09:05:24.854472
9a4f108e-8238-4790-8531-481c192ec376	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDaZHGqJbtMumU7E+jTM6rB9E0yDv16EmCI7EmIDpiLpbfxck4sVvy2aR7PHOEeS8CS/LZtxvMPofSUAk6b6lCibJ1mX/hkoLrdFri+HnG18SlvhW4biXMTbMjwCowfFfpplJpRzjTlLF74DGmbX8j8IlGIgws1UZ6Q8plW3ixK+QIDAQAB	2020-04-27 09:06:36.096869
764e2d15-5639-4012-87fb-29c76865a562	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDaZHGqJbtMumU7E+jTM6rB9E0yDv16EmCI7EmIDpiLpbfxck4sVvy2aR7PHOEeS8CS/LZtxvMPofSUAk6b6lCibJ1mX/hkoLrdFri+HnG18SlvhW4biXMTbMjwCowfFfpplJpRzjTlLF74DGmbX8j8IlGIgws1UZ6Q8plW3ixK+QIDAQAB	2020-04-27 09:06:59.574675
c90f8b91-b8bb-463e-9196-4c0cebe404b5	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-27 09:08:26.767013
e14a35f2-cb1d-4c8e-bd20-f5bcb1ad1451	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNOCupdEBQ3B7XbhIFsi9DVTUNeaD/iXT+WEvmgNzUOxwE05zevRd5rSDsbNnuTmBptbIXX5l58PlTCTeQEuwdlIZeXLwDCw/ifl2BtGJPc5mjeBD+ZhXYJcxyZsG/p01/y3EaLOER5FoVSGe6xvEBKnp9qqILgmowP+RjMa2iMwIDAQAB	2020-04-27 09:43:33.233238
b4fcd8e8-55d6-46ad-9a9c-2ea5a5f48362	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZXp6QIrEkrEISd8fT7fd6zKgqll7Bq9bA64iMLAGkAljaFfg8Kj/LEYOKXWRU3oGLSS31ZCH0nZFgHjlaqzoreTh9Z4thw+vR7LPL+BVEbiBSp/DQ/3ObBHMqNcEVLOy5Yjexs7KS92U5KOJQOA0xWAcGvUdrOi0XQVehEyPaUQIDAQAB	2020-04-27 09:53:28.683923
b5e23859-09cd-42ac-a162-b209240ff031	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-27 09:57:00.439244
7bad278a-3992-4999-8781-5f85199f2990	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC8ZcMSO601Fc0B9vOzv5uBFkEoxaqFwZncYrt+5AMYk3vH5IMJiWvC0X9yaPxafNFIZzsGWl9pyLNjnAddkgoTPt+ndoRjEIG80X0EIWLA+GlKQlWs8roH3CVGR9IpGqU/X3sNumkaoPsOQEpZRjknBc3Hx8hqeQxKWKH7Z1eagQIDAQAB	2020-04-27 10:37:34.288462
a9a57ba3-046d-44c1-a5b0-e256b18eef02	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDAc6/P8hYw1qKalkQjqUrxBXKW7ElJZex03nfKJ8WKvFSsOYwq89LUgvl/qneU44sBHA4fQjd5KbN7xYEzoBhZS0EWpuAxlLFGrWLzjDYIwK3ZiV9cNflorgoK9DMREcIthWs1N8ctUwBPAVfVmAxpAsisnEcKSTlwfoJG4k+8jQIDAQAB	2020-04-27 10:40:13.133695
51de3599-1c99-416d-b450-af9b28ddedcb	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDAc6/P8hYw1qKalkQjqUrxBXKW7ElJZex03nfKJ8WKvFSsOYwq89LUgvl/qneU44sBHA4fQjd5KbN7xYEzoBhZS0EWpuAxlLFGrWLzjDYIwK3ZiV9cNflorgoK9DMREcIthWs1N8ctUwBPAVfVmAxpAsisnEcKSTlwfoJG4k+8jQIDAQAB	2020-04-27 10:40:59.36938
ec81b25b-bed8-413e-b91c-03d0ae629005	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCzCP1PEyjAmY54PXJMsBnS99VQp8pyfCnyAtqdiRr0dEsL6GmYp5M3x7qiIofCOlqsK2qTolVWA2q0I+qJRpRIbdMbI2nARsbmoSrOUeSd3Tn9dhTRJAuU/5yyefHlsd/YHmKya8UqsUZQbddvV+/3M9BNlGms7CnCsEn/dz01TQIDAQAB	2020-04-27 10:52:14.883112
18893bdc-7ad7-4706-8b4d-b030267a1efe	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-27 11:05:38.385171
e88644af-34a7-4a3c-9479-05f45d75691c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCzCP1PEyjAmY54PXJMsBnS99VQp8pyfCnyAtqdiRr0dEsL6GmYp5M3x7qiIofCOlqsK2qTolVWA2q0I+qJRpRIbdMbI2nARsbmoSrOUeSd3Tn9dhTRJAuU/5yyefHlsd/YHmKya8UqsUZQbddvV+/3M9BNlGms7CnCsEn/dz01TQIDAQAB	2020-04-27 13:00:44.042371
cdf7d11c-7c84-4540-9381-a4ccf9217c95	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCLvuhXtjM9xSt8kri3WwYI6e7VgkUrabtcUQ+q4mdEiKoITR/s96+o/UYtoM0keiAuapdJOvVbizkt9E7uEYiQfV+aKj2vAGjdxvbILu4YA+sDrjhWFTTlwXik1xj/1vwCjNJYQQkp673Au+sFmNjwMcsR1UyoMttWY2lG2gWbmQIDAQAB	2020-04-27 13:05:13.330223
fefa740a-ffbc-49fc-a587-42e26cb7e997	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCLvuhXtjM9xSt8kri3WwYI6e7VgkUrabtcUQ+q4mdEiKoITR/s96+o/UYtoM0keiAuapdJOvVbizkt9E7uEYiQfV+aKj2vAGjdxvbILu4YA+sDrjhWFTTlwXik1xj/1vwCjNJYQQkp673Au+sFmNjwMcsR1UyoMttWY2lG2gWbmQIDAQAB	2020-04-27 13:07:44.361843
261a0981-5cf7-4cdc-93c2-119e5f9fc6ce	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 13:16:02.28476
2a7d6689-c4dc-4540-bfb2-ee67be4f4ce5	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCbMKTYiCw7q9WvVvcjYkw4OnYRH+wPgmd8Knk5L+iVYaalU6/ZL/4PjCT+WIkXmTgzkzw6Z7Y/BZ0jVRLOCNimL1gz4dxyBmobdIyw1aR+czDHtVdTIDI8LDxnk3qOv1DdU9MSLb6O7MdJc6KqUgvL28nVlOxI3VgmwyOSxDaieQIDAQAB	2020-04-27 13:16:45.335485
19867a0a-d638-4f01-9202-b45e0ac98b2a	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCvvU9i1owEJYb+r6xkLUOD0OfUwLr3w/zjwWlUD+UH7bXyNfK9ELX6/RAS9+1YG5i3keo6EHsdfIaqa7b1GvCHQsBKoES81UzIxB/W89K1LWLN+ZORr34MVUURsqEK98SviGxcFlsqdaB0c2wRKk7wscR2xvJKc5P1ew0RaJ3+0QIDAQAB	2020-04-27 13:38:54.300759
07152744-ab3d-44d9-b25d-76ad1dc0ef3b	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDFHYFfoubjzR9aSAk5CgnbxlFuUGww4WPvSerUh4iDaFXOAHpcOm2BQqpAhASZN6/bc07/nCHrlyOREcTikNcvGTRoAfRu1ASaJklsPmDopxbOvsOGNkzA527wd8SoMJMA+qPZCI018sGVSxCcpNKTHn0k1Fno65R3DRudDnB8/QIDAQAB	2020-04-27 13:39:45.842962
d2fb8bc1-fae1-49d7-93f0-8c70258f89ad	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDFHYFfoubjzR9aSAk5CgnbxlFuUGww4WPvSerUh4iDaFXOAHpcOm2BQqpAhASZN6/bc07/nCHrlyOREcTikNcvGTRoAfRu1ASaJklsPmDopxbOvsOGNkzA527wd8SoMJMA+qPZCI018sGVSxCcpNKTHn0k1Fno65R3DRudDnB8/QIDAQAB	2020-04-27 13:42:11.218331
1371bb15-f33a-47ef-8713-7bc2ef4e75ff	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCs9PdPO2F5tyJCYJ6ySCuOLJztZP2Jr0uwJtuOMQbyJU64BzIN+1XVy0HpyDHd6VB/okOkzDaLOT3HK1ym6spczfTa/Aaz+H7JUu/BpOqWzhK+5j8EPNPFtlIYC+GXowfIqw//z0wwOKmtPBHiMMsleWmr6Le4vZxPqZtWLxAOxwIDAQAB	2020-04-27 14:41:25.93273
1ea773c5-dfa8-4068-adaf-2c2e444615fa	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNOCupdEBQ3B7XbhIFsi9DVTUNeaD/iXT+WEvmgNzUOxwE05zevRd5rSDsbNnuTmBptbIXX5l58PlTCTeQEuwdlIZeXLwDCw/ifl2BtGJPc5mjeBD+ZhXYJcxyZsG/p01/y3EaLOER5FoVSGe6xvEBKnp9qqILgmowP+RjMa2iMwIDAQAB	2020-04-27 15:04:15.505225
5220bf42-1da1-4ffb-bdf8-87449e4ca11f	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 15:05:46.548134
dc91b5ab-523a-4d17-92d9-fe1bb05d5b31	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 15:06:36.617126
ade30238-768b-4ff9-8464-87d977d91186	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCLvuhXtjM9xSt8kri3WwYI6e7VgkUrabtcUQ+q4mdEiKoITR/s96+o/UYtoM0keiAuapdJOvVbizkt9E7uEYiQfV+aKj2vAGjdxvbILu4YA+sDrjhWFTTlwXik1xj/1vwCjNJYQQkp673Au+sFmNjwMcsR1UyoMttWY2lG2gWbmQIDAQAB	2020-04-27 15:25:05.801578
bd781f1d-e2f5-4617-a5dc-dc5f4f095bf4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNOCupdEBQ3B7XbhIFsi9DVTUNeaD/iXT+WEvmgNzUOxwE05zevRd5rSDsbNnuTmBptbIXX5l58PlTCTeQEuwdlIZeXLwDCw/ifl2BtGJPc5mjeBD+ZhXYJcxyZsG/p01/y3EaLOER5FoVSGe6xvEBKnp9qqILgmowP+RjMa2iMwIDAQAB	2020-04-27 15:26:26.506341
076149d4-c2dc-4e96-a20c-f7d241d9354d	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCvvU9i1owEJYb+r6xkLUOD0OfUwLr3w/zjwWlUD+UH7bXyNfK9ELX6/RAS9+1YG5i3keo6EHsdfIaqa7b1GvCHQsBKoES81UzIxB/W89K1LWLN+ZORr34MVUURsqEK98SviGxcFlsqdaB0c2wRKk7wscR2xvJKc5P1ew0RaJ3+0QIDAQAB	2020-04-27 15:48:08.709831
bd82e59d-fe9b-485c-8020-655a20fdbfc2	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCLvuhXtjM9xSt8kri3WwYI6e7VgkUrabtcUQ+q4mdEiKoITR/s96+o/UYtoM0keiAuapdJOvVbizkt9E7uEYiQfV+aKj2vAGjdxvbILu4YA+sDrjhWFTTlwXik1xj/1vwCjNJYQQkp673Au+sFmNjwMcsR1UyoMttWY2lG2gWbmQIDAQAB	2020-04-27 15:48:27.232404
843a8cf6-ffae-45e9-9b32-6ee097c0bb16	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 16:36:32.589478
064e5270-50eb-4771-8480-3028ddfd7c81	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 16:44:23.435496
0e928e7b-8960-4ac0-ad3d-6d4745012739	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 17:03:13.417647
0e3ce2b0-5c84-4684-a157-2c72ebc7290d	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIo1S8XGUi/yJMuLUbsBhpmqtF0zijoC+qYxRvl7bTs9OchiioIA1+Qb3njXmvORjVYDJqcjnKOfVcn+xUY14SF11w4oksf93ziorDAx7+dL489xiFyO4p6eIBBxXkmWVbS0g0hE5LOV/s38nXqLRunCfHDg4iuqMTKbDCTq5oeQIDAQAB	2020-04-27 17:05:54.897433
4acf7729-d9e0-4acc-9e8c-02e57ef5f61e	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-27 18:27:42.430995
f4625abb-225b-4b0a-82a1-06cfdbd544fc	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCw+TQLsnncihwsIB5bKT+VHMu+gGwfabcmCDooJGDR5UgIgusQxSBuKnxhCUrX5HKF1Kiz11qT6uGYjE0gcn8y8jkq8rRhsqyQB+UL5POMgoJQ6oYreeGGDZeCTNYS2S7N39meLtW3DvDawwZJuMjJ3oGCMxW/f6SjBkSEJcyLDQIDAQAB	2020-04-27 18:33:26.350781
2b58f5ce-2296-418d-b3ef-eaabe8430731	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDauQ2YwFPWKip5HJE3eMhhxTIH6/c6++bYqCg9nzg57dvVUwO/gZI3VtM63XH4DkN9VqB8Tm3Mum1sukyicNrExt7LQjtTLABl3HYWBmzkNCLXp76ZKUDMz+fonTpUFd0qAkMx0MSj6Tz13pnhs01KBrW/wADUDpWDaxuKTWih4QIDAQAB	2020-04-27 20:47:29.777265
a62638ff-4034-4719-a020-37d62b5aba8a	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0qw7bZ/gxT2MlWZ6p/Tn7mtylCrUWx6JVNCwLFw5+WQbgGePhOfj/sujzMvTubZ1haTPZIcOifX3Y9Ju3bWiFUhSOyqK67jXyT+10pt7pCyJb9PIRzpqOPstJj89Xe8cp1xp/ZU6zP3MgtknIqyF0BcWGf1gX4KToIPSthuCaGwIDAQAB	2020-04-28 14:04:10.411032
5990b22d-05a0-4058-9b33-3f1f5d8e8981	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCw+TQLsnncihwsIB5bKT+VHMu+gGwfabcmCDooJGDR5UgIgusQxSBuKnxhCUrX5HKF1Kiz11qT6uGYjE0gcn8y8jkq8rRhsqyQB+UL5POMgoJQ6oYreeGGDZeCTNYS2S7N39meLtW3DvDawwZJuMjJ3oGCMxW/f6SjBkSEJcyLDQIDAQAB	2020-04-28 10:15:39.79292
e207cf44-e7a5-49c1-900f-cca1be81808e	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChbcPwsZFCdzoTkfGZbBaJ19zBqNuGbfrzIjFvaIzHZPwzl6vk6su8gmopyZv6ok1dJCNYD99JHvzfKzbQS90mMYhhuKxSNYNaJ6BXDL3aKsvvgYQYhNoS6w+ZsouPis7UOjeG0oF3bHlSHuvNk2JQXqop3+LZ8gb6rOgUm7G0xQIDAQAB	2020-04-28 18:12:22.619798
68c5d2d9-855f-4d78-ab70-d8192a24fad4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCw+TQLsnncihwsIB5bKT+VHMu+gGwfabcmCDooJGDR5UgIgusQxSBuKnxhCUrX5HKF1Kiz11qT6uGYjE0gcn8y8jkq8rRhsqyQB+UL5POMgoJQ6oYreeGGDZeCTNYS2S7N39meLtW3DvDawwZJuMjJ3oGCMxW/f6SjBkSEJcyLDQIDAQAB	2020-04-28 18:52:40.074748
f74d3379-983c-4900-97bb-1aee6b82f0d4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDauQ2YwFPWKip5HJE3eMhhxTIH6/c6++bYqCg9nzg57dvVUwO/gZI3VtM63XH4DkN9VqB8Tm3Mum1sukyicNrExt7LQjtTLABl3HYWBmzkNCLXp76ZKUDMz+fonTpUFd0qAkMx0MSj6Tz13pnhs01KBrW/wADUDpWDaxuKTWih4QIDAQAB	2020-04-28 22:57:16.60322
924b366d-9522-4a7c-8d75-c41e8f6b4451	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDauQ2YwFPWKip5HJE3eMhhxTIH6/c6++bYqCg9nzg57dvVUwO/gZI3VtM63XH4DkN9VqB8Tm3Mum1sukyicNrExt7LQjtTLABl3HYWBmzkNCLXp76ZKUDMz+fonTpUFd0qAkMx0MSj6Tz13pnhs01KBrW/wADUDpWDaxuKTWih4QIDAQAB	2020-04-28 23:01:53.719796
87fb058e-5674-4fe9-bcbd-80f4817e236e	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCx0U0QpKUMAUGNSOL9+R2yzo1aZXN3tnjq8c4hACkgB1t8jKT2+C3v1jHrHDg6j2RkLZjVmBiMIPzUH8RpC+osiTYTnvTyXafd6t3aLO+wboouv5PN+38oCBAanZBcOwtT3oe6ikV+8mo6d0XTNbRgAj4teEU7lVEw6wIYjuNkIQIDAQAB	2020-04-29 09:28:36.379871
1fc0808b-1ebb-4afa-a710-e4aad50398d4	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCx0U0QpKUMAUGNSOL9+R2yzo1aZXN3tnjq8c4hACkgB1t8jKT2+C3v1jHrHDg6j2RkLZjVmBiMIPzUH8RpC+osiTYTnvTyXafd6t3aLO+wboouv5PN+38oCBAanZBcOwtT3oe6ikV+8mo6d0XTNbRgAj4teEU7lVEw6wIYjuNkIQIDAQAB	2020-04-29 09:29:54.44778
89d3ce7d-b38d-46cb-9ee2-34411c46e823	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCS9bHM5WSsIAkF2viGDaGDBL7h7U0WFmUSNN0niHwpnnOJVAg4lF7Nbo3PPZZa/7AsvhUQP7Q32nW8a0+x3gbqiGD4nIqqHQBL3BgN/DCikGEpJfb0dFFfxKZMsyyUnJ588uZ7z0vpemc9SsxcHaHIhGrM+gzzHN8K9TFkbsw8PQIDAQAB	2020-04-29 09:35:13.899742
76236358-e176-4331-8d6a-32df6c7864ee	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCMfdIQn8Xtkj1yqiUZxBtyzaqQPZIcjmejQEG28+K8sA8s52sdsMW8zCUK42Wz+wBVfgYkPXOLMsQctI3Io1Z1Cl57jKy8hO8bPyXu6NSmswXbpi2RyvUXhxZTLncezNbJ6bfQUWdBRZX3a+xDH4xnOySv8I5csb/cS8PBPT/vcQIDAQAB	2020-04-29 10:08:20.478403
10437f9a-99c2-4ac4-8777-31251fadf954	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCshjRcGSM6FT4AifcYHMZStYdo+tMg4exf4Au4qcHRsv8D2atniFbB3jNzbdWHLfG9uV9pJ5bWQ+sus5uu6BFYsFq6X0RCIGEa/LyNswvlcd0HNJTVtcnwtRIq+o44m9LFviPrVGPs7wCVsfHy0pVANUSFCY2q3khpC9C4Dc/CJwIDAQAB	2020-04-29 10:10:55.570238
357dcfdc-a6d9-4e4c-8248-506f326f7229	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCS9bHM5WSsIAkF2viGDaGDBL7h7U0WFmUSNN0niHwpnnOJVAg4lF7Nbo3PPZZa/7AsvhUQP7Q32nW8a0+x3gbqiGD4nIqqHQBL3BgN/DCikGEpJfb0dFFfxKZMsyyUnJ588uZ7z0vpemc9SsxcHaHIhGrM+gzzHN8K9TFkbsw8PQIDAQAB	2020-04-29 10:28:28.689844
14b04a69-2fae-48f2-8555-14d692ebf362	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChbcPwsZFCdzoTkfGZbBaJ19zBqNuGbfrzIjFvaIzHZPwzl6vk6su8gmopyZv6ok1dJCNYD99JHvzfKzbQS90mMYhhuKxSNYNaJ6BXDL3aKsvvgYQYhNoS6w+ZsouPis7UOjeG0oF3bHlSHuvNk2JQXqop3+LZ8gb6rOgUm7G0xQIDAQAB	2020-04-29 10:47:51.145319
1042f823-4ab9-4822-ab54-cfac5bdc3b84	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0qw7bZ/gxT2MlWZ6p/Tn7mtylCrUWx6JVNCwLFw5+WQbgGePhOfj/sujzMvTubZ1haTPZIcOifX3Y9Ju3bWiFUhSOyqK67jXyT+10pt7pCyJb9PIRzpqOPstJj89Xe8cp1xp/ZU6zP3MgtknIqyF0BcWGf1gX4KToIPSthuCaGwIDAQAB	2020-04-29 13:23:30.053136
4255528e-49db-4947-a85c-527e67e471ec	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNIajMpWh7xLcqJgToh3/gQSQeR1paZgE6ppinZDMz4kF9I87nSy8+XWNjgnHobXiMaMWuzHm3Hl1fMq8qzDzND0FnVlJFKDGSG4DSYj3DAU8Q4vX/AzXn7wB/klErD4fuu01vqnShnx66/4BvI/hNhDSF/l7YiDvjucSs0EVTrwIDAQAB	2020-04-29 13:27:37.366601
db630ca6-12b2-47ee-b6a6-5b362bde9d9c	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNIajMpWh7xLcqJgToh3/gQSQeR1paZgE6ppinZDMz4kF9I87nSy8+XWNjgnHobXiMaMWuzHm3Hl1fMq8qzDzND0FnVlJFKDGSG4DSYj3DAU8Q4vX/AzXn7wB/klErD4fuu01vqnShnx66/4BvI/hNhDSF/l7YiDvjucSs0EVTrwIDAQAB	2020-04-29 16:14:28.408092
02a3bebc-cf5d-4a3b-8ed3-261bcd720c64	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDaZHGqJbtMumU7E+jTM6rB9E0yDv16EmCI7EmIDpiLpbfxck4sVvy2aR7PHOEeS8CS/LZtxvMPofSUAk6b6lCibJ1mX/hkoLrdFri+HnG18SlvhW4biXMTbMjwCowfFfpplJpRzjTlLF74DGmbX8j8IlGIgws1UZ6Q8plW3ixK+QIDAQAB	2020-04-29 16:34:25.097781
19931093-429d-4117-8a72-6cf30ce2db43	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdgfaygBR1xo9GFNHg93WJRiMcmdduZdFOQUEQdJILUoeVA+PUqeMpEQEBhko3ThailivltBUoZy4oh2OzZNeWbAuwJaopvavYnA3pXyJbXkBq7l4BQYx2Pv+8UIFG8ZbAWIl8G4rmVpzJDsHg8GxHJJjBtA2cnyQAvCxsp7nddQIDAQAB	2020-04-29 18:24:11.869434
dbb2bffd-7ecc-42ef-80f4-6b31cf3c0767	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRqaw6yCsYACXzwn7QhTJZt2XbE4Yn/TEQCEr9hBnmgFQiHjWsBDUr6+8amQr2R14rbPPLEzXlsqt+cHLjh4ViX7raMKHE7s2/XX8caTTYAQ4bR0UYbj4/yUjAg2wywpXMR+u5c9G41uBW3TaAV0LnXJHUsYRHL0x4LZ4rbGaFgwIDAQAB	2020-04-29 18:50:33.187549
88547f0d-bcd3-4bb2-953d-163cae5426b0	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDnPlxH3nimIC6ugu1GNUos7YuiY9I63hSgMf4+8kQaekcmgnPUjw3omWDuCHykWT6tTSCGYsetu7VtsMLFu0qaR/IcqnpyBQr9Dnk+JwTgVr6vAYBfXBrpBZYiQmpvOb1TfeYJq2TH/oe6ONNI4XidLD8lE+8ZclTGbnfKSO5rGQIDAQAB	2020-04-29 18:57:39.209953
ef2390ef-5726-42dd-b271-4a983534a642	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-04-29 18:59:12.505266
452cec2e-14d7-4e6e-8758-6d6c824029cb	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCw+TQLsnncihwsIB5bKT+VHMu+gGwfabcmCDooJGDR5UgIgusQxSBuKnxhCUrX5HKF1Kiz11qT6uGYjE0gcn8y8jkq8rRhsqyQB+UL5POMgoJQ6oYreeGGDZeCTNYS2S7N39meLtW3DvDawwZJuMjJ3oGCMxW/f6SjBkSEJcyLDQIDAQAB	2020-04-29 19:02:48.405804
8a486a35-ea5b-4fb3-8d92-ac452f36c33d	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-04-30 09:16:44.086237
3e69fd48-2478-4451-badc-6bbb4f894994	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-04-30 09:24:55.423469
59b1e173-041c-4f1f-977c-0f23ac981b3f	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-04-30 09:27:50.61436
5279e87f-f729-4428-a6da-157fbbd5fda2	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCPji2OIREEsYCeoEbJ4pdTYKhmrrwccapgR8WaXw+fBsAJzxQNcmYcdZ9mZM7i2HX5LLS07xID69/V/INQC6mK6+PswphZHIn9GHqrLvwOO8MyQ0+xtXxV8UCFvvvKXpof3SZ917nf8IlkccOs8ElQPyVHhjOAvl4WGifIrAdkCQIDAQAB	2020-04-30 10:30:22.353308
fa0d0e06-f54a-4677-a731-308ccbd560d5	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCOUGi8nluGvgCGHkfFZK98aTdJnfZ3dIsWmaq/QUeqtNg+fJyfwvnbPgQLSaDBFhKjhawlnVgIUXl7VwYDvmVEGPp3IrRsdyoFNER+uYKrD5woicYE78pOyO2+SlU334/NijRCk1CeHXT0SusCg4xOXn5AhvrZydgkQXMvmDkycQIDAQAB	2020-04-30 10:48:54.201244
050d4350-88d5-4056-9a13-38033d0289c1	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDnoa/hgynaBQ66wGglCC/T48rveZq+jenf/uMov++ron6V8AGgK+3CzThEru4XPguFzMgy7EuRJDdydG5l9NcD765ld0eJYHF21AwiW9u2U3nuTIM7auHdFTAWNHRFYnZV5Ykr1WD8QpztyYf4+Zadwz0EHyVqLdL7vpvOmCdP/wIDAQAB	2020-04-30 10:53:02.859885
36557785-3dc6-414c-83c7-e28afa559df0	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC2TVzeP63u1HSd/0c3+TSEq6YXvXScnzD1Qzprv0HUUeBxe/S9pp+OTg7cAblyyeEnq6/QPdF1fg2drGqdDudKKpnnKexBEh6+LPqIh4scneetVqmxpPzx7IlZQ7RVPbIHVKmzP10Fp60BfEPmuRYyrrKfQb6rb0Z38cghussk6wIDAQAB	2020-04-30 10:55:27.283388
733a30cb-6f17-4e4c-b3cd-57c1881278e1	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCR18T3jioINg+JDeE4CFo2UqxrREo4eKkRorCioIpS3bOtCO2akTXJas77yxlK/c1vO6bT1YirxrSvsuF7ovpXUJ8eG287YahQp2l3lLc1Flt6XlZYvZ5PQMPsAnQ9OOf1sT48VvlK7ongyDB06xoCmmsIOFa7ZbneK8uHxsJAfQIDAQAB	2020-04-30 10:55:41.594704
70fbd33a-9ef1-4d96-9183-b57fa534b6ad	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdcjfCHwT+WyIuCU6DMcNJO/xN+r542ItkbImofBU7+w74FQTCR+XhPPTR8BkaMeTtgSgdjjPP3SNa8XlXLXnYQVacKfuELm3yhJXElbv4TEQjLgNQoFbIn5buywCMUhBL6itGtomCNYB+xaTIqZf/OIr3cCt0gjol1R8wP7JwDwIDAQAB	2020-04-30 13:12:11.204997
a018017e-4131-49b7-981e-94fb34b7354e	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDTJAXBCnTZUNJwk9Vc5e/D+sDLAYDcG7x1eExih0dAAapSybn+RyxzUPOY6z/sPeaGYVyk5EZj81KcnRTCzqCBs594Ga5lZM3IjCE62wI21RknGQQKiDEO9lu03u735LsijyFQPPMocCvjV3LQG352kUmswMN/2jp7PG6yyvHeKwIDAQAB	2020-04-30 13:48:54.815908
e8f9de6e-662c-45eb-99f2-340a4b283802	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCOUGi8nluGvgCGHkfFZK98aTdJnfZ3dIsWmaq/QUeqtNg+fJyfwvnbPgQLSaDBFhKjhawlnVgIUXl7VwYDvmVEGPp3IrRsdyoFNER+uYKrD5woicYE78pOyO2+SlU334/NijRCk1CeHXT0SusCg4xOXn5AhvrZydgkQXMvmDkycQIDAQAB	2020-04-30 13:42:10.947247
705f3142-9296-4c9b-abfe-8fd723936364	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCOUGi8nluGvgCGHkfFZK98aTdJnfZ3dIsWmaq/QUeqtNg+fJyfwvnbPgQLSaDBFhKjhawlnVgIUXl7VwYDvmVEGPp3IrRsdyoFNER+uYKrD5woicYE78pOyO2+SlU334/NijRCk1CeHXT0SusCg4xOXn5AhvrZydgkQXMvmDkycQIDAQAB	2020-04-30 14:03:32.69471
653050c8-a236-4005-b2e6-a8cc78cb5dd2	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDTJAXBCnTZUNJwk9Vc5e/D+sDLAYDcG7x1eExih0dAAapSybn+RyxzUPOY6z/sPeaGYVyk5EZj81KcnRTCzqCBs594Ga5lZM3IjCE62wI21RknGQQKiDEO9lu03u735LsijyFQPPMocCvjV3LQG352kUmswMN/2jp7PG6yyvHeKwIDAQAB	2020-04-30 14:06:04.157454
c22b2c9f-31c5-42eb-857a-6f815a2394a1	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCpYgf+9FsQmEP+PBQoh5FOgbz+Mf1HOxJTMHrs5R9ALv4auWeycZYGrKfjB/618tlKagNEBw3eG7wjjfOau+RG3FAYkmJw5QwVmfNjzdHkNplxGIQO56HeCnY4XNdwm1tAgPOpDGwRsr7dnLMbDypbzztAAvoH4rxPnPsc34u2lwIDAQAB	2020-04-30 14:44:30.976928
5d091f50-262e-4760-9b53-5a62dd1e5b4a	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-04-30 15:15:35.126935
26c877b3-3f4b-42b5-bba0-9a2097da69f9	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-04-30 15:48:09.432133
dfdb7c6d-865a-4c0b-961c-60b88de7ec76	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCw+TQLsnncihwsIB5bKT+VHMu+gGwfabcmCDooJGDR5UgIgusQxSBuKnxhCUrX5HKF1Kiz11qT6uGYjE0gcn8y8jkq8rRhsqyQB+UL5POMgoJQ6oYreeGGDZeCTNYS2S7N39meLtW3DvDawwZJuMjJ3oGCMxW/f6SjBkSEJcyLDQIDAQAB	2020-04-30 16:40:48.58432
bc3d683b-bb07-4331-a8b6-68b28f5ac561	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdcjfCHwT+WyIuCU6DMcNJO/xN+r542ItkbImofBU7+w74FQTCR+XhPPTR8BkaMeTtgSgdjjPP3SNa8XlXLXnYQVacKfuELm3yhJXElbv4TEQjLgNQoFbIn5buywCMUhBL6itGtomCNYB+xaTIqZf/OIr3cCt0gjol1R8wP7JwDwIDAQAB	2020-04-30 16:59:40.857938
c985d9dd-e06b-48df-a6a3-86e87db98a4f	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdcjfCHwT+WyIuCU6DMcNJO/xN+r542ItkbImofBU7+w74FQTCR+XhPPTR8BkaMeTtgSgdjjPP3SNa8XlXLXnYQVacKfuELm3yhJXElbv4TEQjLgNQoFbIn5buywCMUhBL6itGtomCNYB+xaTIqZf/OIr3cCt0gjol1R8wP7JwDwIDAQAB	2020-04-30 17:01:58.928484
23f97792-0985-417c-9347-cfe62594a879	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUWjMwgBZP5apGj9lUzTtv+O/8PpSA+ea8uWL1gQce+NBGv5AYYklrF9Tygw94Sv9t2/dUhudzeGiW7PhjwC2NXzflczGB/k1HamFyji58B1OhD+XjhleVkiLI27heNRrym3kDWoB4Z7imsqDO772fAjNvhS28Rzxl5ZcKIXWiAQIDAQAB	2020-05-06 09:19:30.539639
48354877-76ad-4380-98db-1dd02d3be06f	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRN4tOewCaRbFzXkI7GWoWSmO4lJY8cOndFYmeU99APR0gVlcpvrjwMZPQSHF/lto6Fc/ycR10aMrYwTylvN5rKuFs4L9rQ+MBTW0A6gVxCUm8qjwVk654srqOjfb/JGOTtM5f0vZhiApbgebsXlJxLX8TmhAMsr9/8gKws3brDQIDAQAB	2020-05-06 10:07:54.29367
78cfae6b-a820-4e10-8988-4af23496df2c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsGFOus0TEQHZhFY5b0c8CNhTFGXFUTwpXhkSJ0SN3iSjhMz9UO9rqeau2+0I95LJVY8+7ZS0R8g/JlBrc1YJCMpZvapl/2xrUZ3B57PLoLkkfdQIZ+kFQED0HOK5Lyp/uTFZ6cDJrz00FcpMhx/56dtKVsSnFx44Ohr9mo0fZWQIDAQAB	2020-05-06 14:38:27.472742
49824feb-49f6-4f76-8f68-1b1e094ab120	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsGFOus0TEQHZhFY5b0c8CNhTFGXFUTwpXhkSJ0SN3iSjhMz9UO9rqeau2+0I95LJVY8+7ZS0R8g/JlBrc1YJCMpZvapl/2xrUZ3B57PLoLkkfdQIZ+kFQED0HOK5Lyp/uTFZ6cDJrz00FcpMhx/56dtKVsSnFx44Ohr9mo0fZWQIDAQAB	2020-05-06 14:43:46.101754
eb68947b-85b2-47f0-ab2a-1b7523b053c0	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC5Kt+uWX7hW+4r7kVVKMfhqXiinUrxpo+U9zArfnfdKJX9K9C4UCVmmB4q2BWRSAact81L9JNilcxYQB6TO+GG88FErbNtghTLSw6WnuXChyIEmZjwroAvjxDjzgkSfC5pDv3CoXNj/AhfnTI2TWp/anxZFdg8jqsdbAojuVtPIQIDAQAB	2020-05-06 14:49:28.166923
4769fba7-c653-4b67-b7aa-1aa28b0b02a7	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCEmGTe0dfT51NyKUF3BUaHDhF0QrWaLpgG28ru9OmSEzfFi53EOuTcbl3gth3zsgjGFjN+LucPF/F0Ay7zi2FjfisfKL5WhWliHnY9F+3gZe7RDAd8G5XHn0XSucNa1ZRFlsCo6xkLZvoPz1ldBmH9L4xZ+vEUytj8Pm8TKMgOxwIDAQAB	2020-05-06 14:53:03.007894
89081378-249c-4570-a6af-c07796da6a81	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRN4tOewCaRbFzXkI7GWoWSmO4lJY8cOndFYmeU99APR0gVlcpvrjwMZPQSHF/lto6Fc/ycR10aMrYwTylvN5rKuFs4L9rQ+MBTW0A6gVxCUm8qjwVk654srqOjfb/JGOTtM5f0vZhiApbgebsXlJxLX8TmhAMsr9/8gKws3brDQIDAQAB	2020-05-06 15:10:31.675278
eba30aa0-eaa9-45e8-8d2e-9c5ea6f67b07	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBOd7VjspQ3SS6mPaWW1agETH+NKe8Emt2yFqaVWS/nLDgnrgqXZipVljsNQwKGCXMJ92uMeIZZ14VFCgDhY9XiyUZJpXoAurJaqyjgBvBWPFSeqXp91Ha6e5OjqWe1iAPcC4n0NIbyjACmQEks6dvkEkjHkqE9hgUDTreOkQbCQIDAQAB	2020-05-06 15:16:11.365707
ab4d8185-c187-4496-9b1e-0d1b391a94d4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBOd7VjspQ3SS6mPaWW1agETH+NKe8Emt2yFqaVWS/nLDgnrgqXZipVljsNQwKGCXMJ92uMeIZZ14VFCgDhY9XiyUZJpXoAurJaqyjgBvBWPFSeqXp91Ha6e5OjqWe1iAPcC4n0NIbyjACmQEks6dvkEkjHkqE9hgUDTreOkQbCQIDAQAB	2020-05-06 15:48:59.611232
837842c0-40b8-4b17-9663-d7e9e202472c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBOd7VjspQ3SS6mPaWW1agETH+NKe8Emt2yFqaVWS/nLDgnrgqXZipVljsNQwKGCXMJ92uMeIZZ14VFCgDhY9XiyUZJpXoAurJaqyjgBvBWPFSeqXp91Ha6e5OjqWe1iAPcC4n0NIbyjACmQEks6dvkEkjHkqE9hgUDTreOkQbCQIDAQAB	2020-05-06 18:38:17.853522
aa3c390b-b3cd-45b9-9fc3-26bdf445692a	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBOd7VjspQ3SS6mPaWW1agETH+NKe8Emt2yFqaVWS/nLDgnrgqXZipVljsNQwKGCXMJ92uMeIZZ14VFCgDhY9XiyUZJpXoAurJaqyjgBvBWPFSeqXp91Ha6e5OjqWe1iAPcC4n0NIbyjACmQEks6dvkEkjHkqE9hgUDTreOkQbCQIDAQAB	2020-05-06 19:00:38.707413
d92c9a58-ae97-4673-8be5-03b2dd380dfe	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBOd7VjspQ3SS6mPaWW1agETH+NKe8Emt2yFqaVWS/nLDgnrgqXZipVljsNQwKGCXMJ92uMeIZZ14VFCgDhY9XiyUZJpXoAurJaqyjgBvBWPFSeqXp91Ha6e5OjqWe1iAPcC4n0NIbyjACmQEks6dvkEkjHkqE9hgUDTreOkQbCQIDAQAB	2020-05-07 09:13:46.214483
ee85f635-4cc0-47d9-ad73-c8238ef87866	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsGFOus0TEQHZhFY5b0c8CNhTFGXFUTwpXhkSJ0SN3iSjhMz9UO9rqeau2+0I95LJVY8+7ZS0R8g/JlBrc1YJCMpZvapl/2xrUZ3B57PLoLkkfdQIZ+kFQED0HOK5Lyp/uTFZ6cDJrz00FcpMhx/56dtKVsSnFx44Ohr9mo0fZWQIDAQAB	2020-05-07 12:39:16.125642
132f05b2-2a2e-43f7-b42e-63779b5a4a62	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCi3IGbk/x2EXYJej8zF6R2KDW4SPU7yxCLnLLxC+yk1sHOJ9ufzSxrs79T2E+cp+eLKYLUKfsbTkn8tLCN8IjR2NHpKe+4gStoY5KvqEPmtv3N6SuzCPLGE3hLMbnViScZAaND/yO8dam9KtGXxNjXJLXd736itYdyoB6vwLNGwIDAQAB	2020-05-07 12:59:31.94606
67d4ef23-cf9d-401e-80c8-4e1a8b55518f	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBOd7VjspQ3SS6mPaWW1agETH+NKe8Emt2yFqaVWS/nLDgnrgqXZipVljsNQwKGCXMJ92uMeIZZ14VFCgDhY9XiyUZJpXoAurJaqyjgBvBWPFSeqXp91Ha6e5OjqWe1iAPcC4n0NIbyjACmQEks6dvkEkjHkqE9hgUDTreOkQbCQIDAQAB	2020-05-07 14:26:28.19677
95c61a8a-aafe-459b-b247-c98f9bc1b7da	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDg/sfWRagSRCkjhhey9Ky4+SEgac5bbE4Gf8zHj9CH61F4BaNWQFk4KWm7bVBaPQ7nxAWj+JvcD72DZneqfvFDsUKqagLlygxgt2OWeoI0Uk++hX7uf0GIhmcP22Lh4xOOWfCbwFCtan2Zov3sEvtxUIJt83jA+Qct9zaXBbmMQIDAQAB	2020-05-07 16:13:36.500095
c974590b-a1c3-434d-9018-f75f4cd1008f	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCXonkRQV0rWOcVphCoHT6y4iR20QUlqMXkbcCb7YqdrOgd0Ihj8VSsvi0ARJ4nurz2e9I/Wr5gnfT1caYOfGtLQT+KbI3Wg0geOAilhRaNdjn+778+nhF9wTbtqN/GmTTh6ggOrAbGE4rsPQeRhkwa8gJDNqv3fURi4/uVs7qFlQIDAQAB	2020-05-07 16:17:55.269157
63e4e2ae-f401-4c2a-a99e-716f384146c0	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC03N5Dvq4afxv/eY7gOgLfr539MSZBgGrKtiPW+nqVHnAcE8OWW5P5aujarptbP/GBoD7y+G4PQ92uL1sVdV4V1qJfbEsOCyaepiLc+eHztgQQKBeAjZEfVfHKbWJdq4zSzqGPJnR6+OER6PzMjIXj+U9uaaG9otEDtOWNknUrtwIDAQAB	2020-05-07 17:41:46.238812
dd5162b1-9197-4dc6-9e20-67502413e5e5	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC03N5Dvq4afxv/eY7gOgLfr539MSZBgGrKtiPW+nqVHnAcE8OWW5P5aujarptbP/GBoD7y+G4PQ92uL1sVdV4V1qJfbEsOCyaepiLc+eHztgQQKBeAjZEfVfHKbWJdq4zSzqGPJnR6+OER6PzMjIXj+U9uaaG9otEDtOWNknUrtwIDAQAB	2020-05-08 09:14:23.63347
6502906f-7147-4062-8389-feaa3e056e32	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC03N5Dvq4afxv/eY7gOgLfr539MSZBgGrKtiPW+nqVHnAcE8OWW5P5aujarptbP/GBoD7y+G4PQ92uL1sVdV4V1qJfbEsOCyaepiLc+eHztgQQKBeAjZEfVfHKbWJdq4zSzqGPJnR6+OER6PzMjIXj+U9uaaG9otEDtOWNknUrtwIDAQAB	2020-05-08 09:18:10.484973
2088e5e4-0be5-4db4-bb2a-7c22a4c98901	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCgcCKn8zbxFSAtJibJYFGSwLt9cAcipPWZP/teB00V1AqAHqa7gGROrKmd6Q+yzQ595QxmyesuSFPWI1Uabsw0hpFrF8EvE3ILt+YNcTeTzJzry3FIEOgqwBeL6fHs3aQYGXYeioh+tOrVo5CP8Yx35MSqQLR5RAydc/+mE7kVmQIDAQAB	2020-05-08 09:19:46.08589
c969ea22-fe16-4147-a72a-156671f208e5	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqN9KKI8iyOEGCo2o1c+hJnLT/hhK24Re7aPChUCefkHLN3/jc83rkYZFllI2Blmz1aXetWp73xIKnA+1+txuHjbsZGsKa+k+fuPQT4UUwiUtsmZsXJnrxKmi8VzSb6zvx4IeIvBSiO2bO3KoAkKBM52dNg5SwdmxE1B9AODxjGwIDAQAB	2020-05-08 09:22:40.640792
6d4f86ab-2da8-4dd0-8bea-718d5aea3b11	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC03N5Dvq4afxv/eY7gOgLfr539MSZBgGrKtiPW+nqVHnAcE8OWW5P5aujarptbP/GBoD7y+G4PQ92uL1sVdV4V1qJfbEsOCyaepiLc+eHztgQQKBeAjZEfVfHKbWJdq4zSzqGPJnR6+OER6PzMjIXj+U9uaaG9otEDtOWNknUrtwIDAQAB	2020-05-08 09:29:25.765858
c0e4f6a7-6b87-4bc2-bd4a-8c393b1e30e2	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCo3yJG1+gvD0HZLXKvOSbTNAw0owfBbuQhmL0nqHseg0lZ9cEW400iXIwRB5Vdu3Nx7+mQVV/bph6bqEarUmmzG1yRhQ8RcoP7XstyR1ANqeW74ysW23KQ7BXFnjmf5vaZVyJ8k+M5erMGS716Iwl9HIJwUYYAAYn4PfvRr2lLywIDAQAB	2020-05-08 09:43:45.149699
1be519b3-5f11-4762-973d-69372f1802e9	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCpYgf+9FsQmEP+PBQoh5FOgbz+Mf1HOxJTMHrs5R9ALv4auWeycZYGrKfjB/618tlKagNEBw3eG7wjjfOau+RG3FAYkmJw5QwVmfNjzdHkNplxGIQO56HeCnY4XNdwm1tAgPOpDGwRsr7dnLMbDypbzztAAvoH4rxPnPsc34u2lwIDAQAB	2020-05-08 10:20:19.287848
02b0dfda-7cac-49dd-8628-d5fea9109125	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCpYgf+9FsQmEP+PBQoh5FOgbz+Mf1HOxJTMHrs5R9ALv4auWeycZYGrKfjB/618tlKagNEBw3eG7wjjfOau+RG3FAYkmJw5QwVmfNjzdHkNplxGIQO56HeCnY4XNdwm1tAgPOpDGwRsr7dnLMbDypbzztAAvoH4rxPnPsc34u2lwIDAQAB	2020-05-08 09:21:24.124122
fde3f258-632c-493e-846b-ff8c13497056	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCgcCKn8zbxFSAtJibJYFGSwLt9cAcipPWZP/teB00V1AqAHqa7gGROrKmd6Q+yzQ595QxmyesuSFPWI1Uabsw0hpFrF8EvE3ILt+YNcTeTzJzry3FIEOgqwBeL6fHs3aQYGXYeioh+tOrVo5CP8Yx35MSqQLR5RAydc/+mE7kVmQIDAQAB	2020-05-08 09:22:09.424962
3341b466-6f7b-4763-a6f2-46abc6236c76	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCo3yJG1+gvD0HZLXKvOSbTNAw0owfBbuQhmL0nqHseg0lZ9cEW400iXIwRB5Vdu3Nx7+mQVV/bph6bqEarUmmzG1yRhQ8RcoP7XstyR1ANqeW74ysW23KQ7BXFnjmf5vaZVyJ8k+M5erMGS716Iwl9HIJwUYYAAYn4PfvRr2lLywIDAQAB	2020-05-08 09:42:58.386189
8556808d-cebb-4560-913c-9aba02dbcb90	111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCpYgf+9FsQmEP+PBQoh5FOgbz+Mf1HOxJTMHrs5R9ALv4auWeycZYGrKfjB/618tlKagNEBw3eG7wjjfOau+RG3FAYkmJw5QwVmfNjzdHkNplxGIQO56HeCnY4XNdwm1tAgPOpDGwRsr7dnLMbDypbzztAAvoH4rxPnPsc34u2lwIDAQAB	2020-05-08 10:19:45.525035
03b9bb04-9a5e-4c88-b5f5-28b5453ea0d4	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCgcCKn8zbxFSAtJibJYFGSwLt9cAcipPWZP/teB00V1AqAHqa7gGROrKmd6Q+yzQ595QxmyesuSFPWI1Uabsw0hpFrF8EvE3ILt+YNcTeTzJzry3FIEOgqwBeL6fHs3aQYGXYeioh+tOrVo5CP8Yx35MSqQLR5RAydc/+mE7kVmQIDAQAB	2020-05-08 10:23:18.505051
3b3d20f5-c269-4322-8f9d-ac4c9b64086c	1111	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqN9KKI8iyOEGCo2o1c+hJnLT/hhK24Re7aPChUCefkHLN3/jc83rkYZFllI2Blmz1aXetWp73xIKnA+1+txuHjbsZGsKa+k+fuPQT4UUwiUtsmZsXJnrxKmi8VzSb6zvx4IeIvBSiO2bO3KoAkKBM52dNg5SwdmxE1B9AODxjGwIDAQAB	2020-05-08 13:32:30.950456
\.


--
-- Data for Name: rela_role_url; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rela_role_url (rela_uid, url_uid, role_uid) FROM stdin;
f34345a2-e288-428f-89ef-3c9851acb0d5	0e143513-3b89-4e2d-9865-8a848ad6c8db	ad726a2c-6581-4675-a770-a4b58797b813
0ec26acf-c706-4a33-bf89-bac05a7efe97	ed06f369-9d92-4d07-8b9a-6488694b4d1c	ad726a2c-6581-4675-a770-a4b58797b813
effcbd6e-7183-461a-bdf2-2492f3c15442	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	ad726a2c-6581-4675-a770-a4b58797b813
9d095dfb-5ea7-42ef-9f2d-cceb3a3be27c	58405e28-6fff-42e1-8567-01466845abfb	ad726a2c-6581-4675-a770-a4b58797b813
b8728010-817f-46af-a5a5-f25b1d64a6bd	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	ad726a2c-6581-4675-a770-a4b58797b813
d2c5d820-9ba8-4473-bb9b-0f1220c20286	68fc3f82-6076-412f-b986-20d3cb1bceaa	ad726a2c-6581-4675-a770-a4b58797b813
6d2e978b-f271-4fab-adbf-6e36597e3ed8	ccad896d-13db-4fc8-be92-574af1d80341	ad726a2c-6581-4675-a770-a4b58797b813
adb14322-62f5-4611-aa9b-5887825186db	03264ac4-7754-48ef-a824-5269d1f84e4b	ad726a2c-6581-4675-a770-a4b58797b813
49458e12-8f35-44b3-a972-5026d70ff373	51f985e0-4f2a-412a-adb5-3adf5399fc75	ad726a2c-6581-4675-a770-a4b58797b813
16310b94-a403-4480-8d34-315918af983b	77f2a333-b579-421b-bf3f-5b913e6002b5	ad726a2c-6581-4675-a770-a4b58797b813
bdad8ef5-c822-40cb-9a6f-0f9210df8047	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	ad726a2c-6581-4675-a770-a4b58797b813
d6880c26-ea8d-482e-a9c1-8c1b40a5d070	05a2e58e-0551-40c5-9bff-657768d7e689	ad726a2c-6581-4675-a770-a4b58797b813
109c7b87-5b18-4079-8248-1b72d0ca17d2	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	ad726a2c-6581-4675-a770-a4b58797b813
11e087aa-7b77-4a8c-9aea-752ac273daeb	fb9fbdf7-d05f-40ac-95e1-e10019ed8165	ad726a2c-6581-4675-a770-a4b58797b813
145ddc19-5caa-4b8a-940a-65b4dd853ef0	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	ad726a2c-6581-4675-a770-a4b58797b813
a7a3b8da-f15f-4362-8846-b376c019871c	39732090-0856-4b05-a241-df9886a313a5	ad726a2c-6581-4675-a770-a4b58797b813
2d7de6f4-10fb-45c8-8240-0631a0cd192f	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	ad726a2c-6581-4675-a770-a4b58797b813
d560e403-3807-4289-9579-28320f6ac492	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	ad726a2c-6581-4675-a770-a4b58797b813
dd1fe3d6-32e9-4008-9716-1deac51ad2bd	b48dccdc-da4f-4914-bd2e-1971de87fda1	ad726a2c-6581-4675-a770-a4b58797b813
be949ed6-55c0-4064-957c-1590aa9e4e0d	3aa46e51-bac0-4711-bf1b-bb396e71fe06	ad726a2c-6581-4675-a770-a4b58797b813
0757f422-3a3c-4e08-a0ec-5d95d94609b4	27917132-65e1-4218-bdfc-40a62550a141	ad726a2c-6581-4675-a770-a4b58797b813
f31cf135-d650-4cef-8cca-07e3fbde7246	09cfad53-de06-4a25-88c0-9eaa26345df9	ad726a2c-6581-4675-a770-a4b58797b813
694f2783-672e-4199-9953-7f32a795a831	7767088d-9354-42d4-8cc7-9d3f2c0bef88	ad726a2c-6581-4675-a770-a4b58797b813
0b74f974-6d17-49fe-86d2-a99a7623c429	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	ad726a2c-6581-4675-a770-a4b58797b813
f713cf4c-ae9a-4e1a-8813-34a694d32375	f22b7a3d-2d76-483c-af40-dd0c04defb20	ad726a2c-6581-4675-a770-a4b58797b813
a10b0606-66fd-4a77-b923-59d17d33764f	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	ad726a2c-6581-4675-a770-a4b58797b813
08d515e4-dea4-4496-880f-527988662fc5	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	ad726a2c-6581-4675-a770-a4b58797b813
9a3a06df-5e5e-49b3-a8d7-65c875d8e089	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	ad726a2c-6581-4675-a770-a4b58797b813
040af3dd-af68-4e59-91ed-6ce5e899a4f9	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	ad726a2c-6581-4675-a770-a4b58797b813
d4ac4a33-1505-427b-b6d9-2bee288ea681	c45be260-e11b-4fae-8e33-05d0803f4e6b	ad726a2c-6581-4675-a770-a4b58797b813
65be26e7-c299-4be4-8d32-2f0e54d2a5af	24127ef9-9ef3-46f3-88cd-97378e694a6c	ad726a2c-6581-4675-a770-a4b58797b813
19248a55-1f32-4580-9994-5e87ad78ec11	24a99997-51d4-446c-bc51-2f46e3ca01a6	ad726a2c-6581-4675-a770-a4b58797b813
53fd853c-8e72-4180-815e-f3d9da12ba58	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	ad726a2c-6581-4675-a770-a4b58797b813
518d5d30-00c0-4b48-875c-757b88d8cd6e	f3c2ef69-090e-4f61-a9d6-954330e86d8f	ad726a2c-6581-4675-a770-a4b58797b813
fab8f63c-5a81-437b-84cf-20b3f85c6a0e	c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	ad726a2c-6581-4675-a770-a4b58797b813
761127a1-39f9-4356-abcd-5a1bcb97f557	ebb13d64-4568-4cce-9f26-72f5e8a9bc14	ad726a2c-6581-4675-a770-a4b58797b813
b801d2b8-58f2-4afe-948a-16a651396079	288731cb-561c-4d4a-8951-d8c359a51846	ad726a2c-6581-4675-a770-a4b58797b813
0d5d4786-0370-43c7-a75d-ae9650a2a9da	9f4da9ef-7f9c-407e-91a9-034f623b810f	ad726a2c-6581-4675-a770-a4b58797b813
d58af92b-3cc6-4a3d-b8e9-6e874a1de8cf	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	ad726a2c-6581-4675-a770-a4b58797b813
1da9919c-00f4-485a-b599-56d75febca5d	653eb679-3e69-442e-b0ff-fd22e3fd8cd3	ad726a2c-6581-4675-a770-a4b58797b813
863130d2-264c-489a-9ca5-71512c806ed2	b439edb9-605f-47dc-aee2-8df45976736e	ad726a2c-6581-4675-a770-a4b58797b813
81239cd0-b526-47ca-8f2f-91f4e6a6708f	086e9acf-0405-449c-9dee-ee33d80df312	ad726a2c-6581-4675-a770-a4b58797b813
e73318e2-76e9-4187-95d6-cbc45f385da8	6f4cc747-5450-4514-9a3e-123ace9cb232	ad726a2c-6581-4675-a770-a4b58797b813
96758cbc-b10a-4f71-a37a-848517a12f9b	93421699-5c88-4dee-995c-f74222b3a1a7	ad726a2c-6581-4675-a770-a4b58797b813
30d0002a-931e-45cf-bb9d-59653eba8db9	c070d24f-76eb-4684-8459-dbb7977acc46	ad726a2c-6581-4675-a770-a4b58797b813
466a9999-3f90-4f39-9e93-7692aa5aaf0b	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	ad726a2c-6581-4675-a770-a4b58797b813
ad57b653-e773-4638-be24-d8462965d5e7	d81d630f-3d80-4168-923b-65eefae629c9	ad726a2c-6581-4675-a770-a4b58797b813
f0ac912a-cdd2-49f1-b395-655c814f55b5	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	ad726a2c-6581-4675-a770-a4b58797b813
0217405c-d8af-4c0e-8242-56dd89d9141e	ae18fd5a-e0ba-4ace-a637-902317a67773	ad726a2c-6581-4675-a770-a4b58797b813
f8007280-a727-4d25-8c3b-076e9b2b7a39	72237dba-0beb-450e-9df6-ea38d3255cab	ad726a2c-6581-4675-a770-a4b58797b813
ba5d67a9-9ca2-411b-96e5-94e2190c0f11	223c3afe-ef65-4b86-9733-6284396cd65f	ad726a2c-6581-4675-a770-a4b58797b813
ee98eb02-d935-4bf5-9c0c-a4e64b967f81	67522e75-9672-4bbc-904d-57dd942290b7	ad726a2c-6581-4675-a770-a4b58797b813
660d4883-d7e2-4283-8a46-710d3700295b	78cd88ae-d982-48c6-998c-1a6ba961b5f0	ad726a2c-6581-4675-a770-a4b58797b813
6197871e-e9b9-4d3a-b4d2-954783005c19	131cc1b8-2668-4392-9443-6c2e42a426ef	ad726a2c-6581-4675-a770-a4b58797b813
170f5a4d-244e-4e6d-ae73-b4eecfce2ade	20805d38-e63d-4117-b902-7007c1697a13	ad726a2c-6581-4675-a770-a4b58797b813
9eb117e8-2fd7-4ee1-bdc8-6fedee3f3290	f53252df-70b6-4907-8847-908554862006	ad726a2c-6581-4675-a770-a4b58797b813
e783595c-a665-45db-bff6-8727932b8b5d	99a0afca-2fef-4083-82a2-4785ac933ec0	ad726a2c-6581-4675-a770-a4b58797b813
3577592e-5020-4399-af84-c9b504f4c79f	3780a865-b70a-4a34-93f1-56c920e44100	ad726a2c-6581-4675-a770-a4b58797b813
f8897c11-815c-4cfa-8c2a-bfdc6b22e0d0	a06b41ea-3abb-4dc4-a4f7-90bfc86a9592	ad726a2c-6581-4675-a770-a4b58797b813
072263e5-a23e-4340-84e5-7a24b25c813d	0a9b8878-506b-4f76-af54-478e96f06eb9	ad726a2c-6581-4675-a770-a4b58797b813
c43ebf53-91b2-4db7-bef9-188423a6e529	25c7b3c0-efd1-4d00-a748-e37de9ca4e3c	ad726a2c-6581-4675-a770-a4b58797b813
1c2b3a90-8ef3-47a1-ab37-4a1d7ba51335	fc72276f-66bb-4609-b83a-4a84ff55100c	ad726a2c-6581-4675-a770-a4b58797b813
774c6364-52c4-477b-bb18-a2ae5608f2f2	fcea9f1c-b882-4368-a7a3-c8bbd2491684	ad726a2c-6581-4675-a770-a4b58797b813
e132e075-704e-43e9-89f8-79dc179a7833	0e143513-3b89-4e2d-9865-8a848ad6c8db	b2222924-ce4e-444f-8e1a-c3479bb99467
0dae224a-83b5-4389-a432-d51777468973	ed06f369-9d92-4d07-8b9a-6488694b4d1c	b2222924-ce4e-444f-8e1a-c3479bb99467
a0e4bde3-972c-49e0-8aea-b9a2e736d8e1	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	b2222924-ce4e-444f-8e1a-c3479bb99467
6dc68e39-391a-4643-9fff-6dd492048be6	58405e28-6fff-42e1-8567-01466845abfb	b2222924-ce4e-444f-8e1a-c3479bb99467
a1f42cce-1a49-4f01-a968-23edd503baa3	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	b2222924-ce4e-444f-8e1a-c3479bb99467
295435a3-7136-42fb-95c1-fe07e944e40e	68fc3f82-6076-412f-b986-20d3cb1bceaa	b2222924-ce4e-444f-8e1a-c3479bb99467
48e07b62-bed0-40d5-9ff4-16e8a30b65d9	ccad896d-13db-4fc8-be92-574af1d80341	b2222924-ce4e-444f-8e1a-c3479bb99467
ad72c96e-44d8-4570-b6a9-fa1308e6ea63	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	b2222924-ce4e-444f-8e1a-c3479bb99467
8a1ea9cb-744c-4388-8e58-49bdceb6e04d	05a2e58e-0551-40c5-9bff-657768d7e689	b2222924-ce4e-444f-8e1a-c3479bb99467
87a3216b-fc0c-47cc-94bf-82e4cfd9a0b7	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	b2222924-ce4e-444f-8e1a-c3479bb99467
ebbededb-cb5c-48b3-a5e0-367d9ce5ea1c	39732090-0856-4b05-a241-df9886a313a5	b2222924-ce4e-444f-8e1a-c3479bb99467
e26a77c0-e43d-4d0e-bf1f-a5fb20ced461	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	b2222924-ce4e-444f-8e1a-c3479bb99467
a46f9b0b-98e6-415a-bfdb-154a81697af6	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	b2222924-ce4e-444f-8e1a-c3479bb99467
24363e83-1165-47ea-83a8-0db187b0c79a	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	b2222924-ce4e-444f-8e1a-c3479bb99467
90d9aae9-ab13-493e-95ca-249360037d84	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	b2222924-ce4e-444f-8e1a-c3479bb99467
e89f9667-d71e-47a9-a94b-507c32a940d7	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	b2222924-ce4e-444f-8e1a-c3479bb99467
6a697d3a-4bb4-45a1-8348-24a2b379a8fa	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	b2222924-ce4e-444f-8e1a-c3479bb99467
da515b79-5c9e-48c7-a378-87a504f88836	c45be260-e11b-4fae-8e33-05d0803f4e6b	b2222924-ce4e-444f-8e1a-c3479bb99467
a1dce43d-958c-467d-a17c-68764a0ebaf1	24127ef9-9ef3-46f3-88cd-97378e694a6c	b2222924-ce4e-444f-8e1a-c3479bb99467
24e243af-f40d-4e00-901e-30b4881adfe5	6f4cc747-5450-4514-9a3e-123ace9cb232	b2222924-ce4e-444f-8e1a-c3479bb99467
a2714a5f-f69a-48c8-89ff-ce5281162953	93421699-5c88-4dee-995c-f74222b3a1a7	b2222924-ce4e-444f-8e1a-c3479bb99467
26e01127-4baf-4a0b-9c93-42c46a98920f	c070d24f-76eb-4684-8459-dbb7977acc46	b2222924-ce4e-444f-8e1a-c3479bb99467
612496dc-f66e-4e49-8c12-c87e2c320597	99b01f22-dee8-4f63-ab8e-4b5083d38044	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
f6b73f8a-6619-4558-a24d-82aa02589b98	67948dc1-eab6-4fed-ad68-7879030be523	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
ea415195-0f94-4a83-9204-b3b1515507cd	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	b2222924-ce4e-444f-8e1a-c3479bb99467
16064de6-f52a-4700-bfc0-24534ad41fe7	d81d630f-3d80-4168-923b-65eefae629c9	b2222924-ce4e-444f-8e1a-c3479bb99467
3e116966-a656-49bc-aa74-58d926890dee	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	b2222924-ce4e-444f-8e1a-c3479bb99467
eb84d1e6-cfb0-4519-9f61-d0bd286398cb	ae18fd5a-e0ba-4ace-a637-902317a67773	b2222924-ce4e-444f-8e1a-c3479bb99467
c7d969f9-64e2-44d4-a3d2-4d52fa386a83	72237dba-0beb-450e-9df6-ea38d3255cab	b2222924-ce4e-444f-8e1a-c3479bb99467
89b1d6fa-3464-4b71-9580-10e78224d26b	3ab532fe-274c-4755-b8c4-8055ea124209	75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a
61e1b45e-01aa-481c-a930-5c43bfbd4a3e	a8de69fd-5856-42cc-8812-fef457a99079	75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a
dd10b48f-cdd5-4ba4-94d8-9cf0d08c48e5	136f2190-6f80-46d4-a124-d512e508be5e	75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a
a1b647b5-0365-44ea-a591-195a0ed0735d	0ad5ba51-e406-477c-bdc2-fe35ca35b6ec	75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a
3bc97363-cd14-47b7-8ece-a2824b6ae1ed	eb15be96-d809-4a03-b3e4-be03b79b55cc	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
17b13db7-d119-4b05-8777-25ff65b03517	f4f9e4b1-e4df-418b-a88d-2bd0f00fa5e5	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
abed23dd-7a88-41ae-a7de-d56549d4f28d	4e52b2f4-1da6-4001-8073-38f2b12d4951	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
04c47dc9-354c-47ad-9e83-380aef573431	d212c732-5c62-40ee-b1cc-184eddab4a61	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
75c2e851-5ea5-4d8b-86cb-14ac68116627	797cecc2-d1f9-499a-b79d-e1f7e77c7874	75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a
6b2a0764-f31c-4d78-b44b-f4f1d876a6c8	15668ec7-6731-4646-a516-41035ddebe0d	75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a
77c82800-92f6-434e-b0d8-2e00229438f9	f50b81e9-1797-40e7-a91b-5ce8be353e8b	fbb11569-57f0-4b75-ae29-3854e26a0585
f2790454-9d7a-4154-b04a-a1ae52754bc5	fa083891-8c7b-4d3e-a769-9c765760f79e	fbb11569-57f0-4b75-ae29-3854e26a0585
2f95c2c9-e876-4234-bf12-a9ae524af62c	5e6cb690-6a8d-4fab-9234-754db94e1caa	fbb11569-57f0-4b75-ae29-3854e26a0585
0d96bd9f-fa40-4522-8e2c-685a290b8b25	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	fbb11569-57f0-4b75-ae29-3854e26a0585
44b46443-0308-431c-901c-f6e38b6c67a3	58d6c0a6-e80d-4d79-972f-0ad5f8517e97	fbb11569-57f0-4b75-ae29-3854e26a0585
93391d43-18de-48c0-8d91-2c2386b033f0	0f13980e-ffc3-4712-9c58-1726e1ec6787	fbb11569-57f0-4b75-ae29-3854e26a0585
964cffa6-db5b-4347-864a-d1ce381257cb	33bea49e-3565-469c-aca3-b524ed7232dd	fbb11569-57f0-4b75-ae29-3854e26a0585
90872318-aaaf-4a6a-823f-2a5291a759b5	d6791542-86e6-43a7-ae04-7d42c51889a5	fbb11569-57f0-4b75-ae29-3854e26a0585
40368b5e-8816-484a-9d69-e37f1f661b41	7bd65052-1ff5-4870-abdf-e398ac9a5fc3	fbb11569-57f0-4b75-ae29-3854e26a0585
a5c59f53-e874-4f51-bfc5-152804d2f301	eb9ea9f9-cd44-4b11-a805-061612ff12ef	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
eaa8e84e-3dc2-471e-99f9-f49896c53ce4	aded78fc-0d49-436e-a780-ebea7c8b0f4c	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
eb64de69-87dc-4ab1-93b8-439d931830cc	14f49944-6337-4c1c-8ea4-dd111c7859f1	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
eff8076e-c149-4a65-b2a2-c197113ce5e7	fa13b0c6-1e65-45de-bef8-ac534342f0f3	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
97cfd9cb-70a9-4225-83b5-4d40c176533b	22028342-5177-4425-8465-60d41a46eb7e	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
a27fd75a-aa44-451c-8153-e6026dd400a6	4a0bcbc4-1eeb-4801-9779-d51641ecd982	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
54d6a381-4992-4c42-a18a-cc6bd036945d	f93f8fe1-b61f-4cd9-acf6-ff7d6fd5fd60	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
940419e7-af6a-4027-8bd7-dccdb2dfd420	040a0e65-262b-44ce-9a39-336f08560ad4	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fd12eee0-615f-4ef9-9cb6-01446fac40e7	6264bced-71da-43f1-8fdd-5a96925164cc	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
9e8e418c-7b46-4187-898f-59434ed785f4	84a7f310-51f3-4799-a737-d933083f96a2	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
045681b4-9767-4f0c-a99f-467f87dda4e7	acde5628-57d0-495f-88ba-db5c6c079e87	fbb11569-57f0-4b75-ae29-3854e26a0585
bb69147a-c7d8-40b9-941d-eb4b1daa3acb	64cbf7b5-71e0-4179-a0a8-07020335901d	fbb11569-57f0-4b75-ae29-3854e26a0585
f6701603-6e0e-4fc6-a16f-0a276e2730a5	680460f5-af36-47e7-9278-cc316be1ec1f	fbb11569-57f0-4b75-ae29-3854e26a0585
a268bcf4-ac45-486f-aa5c-4d8404193b1f	a27e89ae-9bb8-4425-be64-76b5a811c137	fbb11569-57f0-4b75-ae29-3854e26a0585
9ec9f93b-0855-4019-8675-36587781ce3e	e16cdfff-e9b4-42c2-b712-018adcfee5e2	fbb11569-57f0-4b75-ae29-3854e26a0585
1fecf608-fc7f-40fa-8a30-b174fa4ab74c	1e1a0c33-63e3-45e7-aae6-5856c0add615	fbb11569-57f0-4b75-ae29-3854e26a0585
48951b5c-d01c-4979-b949-f289d01dd9ae	359d32ec-2919-4ca3-b927-232440bd886f	fbb11569-57f0-4b75-ae29-3854e26a0585
1fc216de-e290-4397-8b45-a543b9c6b72e	04032947-d9b3-4949-ad03-f6d9640a627a	fbb11569-57f0-4b75-ae29-3854e26a0585
6ed6db54-6707-4ddc-b1cc-bf55b9dcd77a	8ceae246-d9fc-4fdf-8487-c4abc76c26af	fbb11569-57f0-4b75-ae29-3854e26a0585
3f06ccea-db4c-494c-810e-aa120c9725c1	6330eaba-e5de-47f1-8c50-072120823683	fbb11569-57f0-4b75-ae29-3854e26a0585
ea8a4f5e-a819-4131-a782-ce0e8b594543	9b1d0f10-1be2-417c-a79f-a77e9dfa0c93	fbb11569-57f0-4b75-ae29-3854e26a0585
f8766be1-9a3d-413d-9f3a-4226a79a14d0	608b52b4-40a6-43f1-bca2-031b2ea2abf4	fbb11569-57f0-4b75-ae29-3854e26a0585
f80ee6a9-ef24-430a-88d2-0b69baf82edb	67591c4a-7794-421e-8c3d-cc8b94b694b5	fbb11569-57f0-4b75-ae29-3854e26a0585
8f95e391-f1dc-4d81-bfa6-8eb39576301a	a0419fd4-acd0-4e2b-8156-30066469f33e	fbb11569-57f0-4b75-ae29-3854e26a0585
a04fe268-3eec-4819-b23d-34d0eb552d79	222d3e77-c07d-4c30-a384-b0c493a3cb38	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
1242b348-f317-4660-997f-a3a0e3c96946	6b7a79ac-30dd-44ce-a6bb-82a137f6771b	fbb11569-57f0-4b75-ae29-3854e26a0585
91de2c8d-91c9-4ffe-8466-ef6a2d4859a7	5af52e1b-b6a1-4eca-87ec-466abacb4302	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
5dc957e2-013b-4ea5-ad34-8ab360b0e3d3	b2776c17-fde6-4af5-8d34-f72d5ec88ac9	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b607ed29-3d2c-4bff-8956-fceff9391d6b	4badf803-d758-47a4-8a4b-fb005cf7167a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fc050775-5301-48b8-8c21-95bedb87e457	e27c19dc-5b9c-45ef-8af3-0c024d1abf34	fbb11569-57f0-4b75-ae29-3854e26a0585
5ba8d928-5274-4436-a30d-00db491d9c9c	98a5950b-849c-4fef-8735-979aa36ed84c	fbb11569-57f0-4b75-ae29-3854e26a0585
a16ccbf3-8ee1-41fd-9308-a7e2fb25a345	6f13190f-abaa-48e1-a490-8d4406df0f34	fbb11569-57f0-4b75-ae29-3854e26a0585
75b3a148-890c-4750-8b76-d8605f4a24c6	07bd2d37-21cd-43b3-9d38-eeafaa73bcf9	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
8ee47439-de7f-4a39-b401-19468cc69072	3f3abbb0-4c6b-4b05-9e0f-4413b4adf347	fbb11569-57f0-4b75-ae29-3854e26a0585
745d848d-768d-470f-ab1e-008976f30840	189208bd-a65f-44d8-8a07-96647feaafb7	fbb11569-57f0-4b75-ae29-3854e26a0585
abfc155d-b001-4f5b-92a2-489311dedfff	9aae5b79-2f48-461e-a6b2-45e8266402dd	fbb11569-57f0-4b75-ae29-3854e26a0585
83d494d4-8380-4f43-be8b-864dda624434	b8840d1d-3908-4db7-9825-b5ce78232409	fbb11569-57f0-4b75-ae29-3854e26a0585
aa518eaa-10b0-41ad-9e59-3bdf104d059a	6be56b39-079c-410a-b41a-35b32a6a5e1f	fbb11569-57f0-4b75-ae29-3854e26a0585
502d96c9-7e4f-4d34-8f97-e77cd0e43101	28c94a48-fd22-49a2-aa57-0f72c1f635d2	fbb11569-57f0-4b75-ae29-3854e26a0585
75b3e6a3-4d50-420d-a2f4-f08fd457e5e4	05825db2-caf8-4951-9f2d-ea85ace66073	fbb11569-57f0-4b75-ae29-3854e26a0585
31b9a6ba-f210-4b6c-8f99-45bdf18d381f	e72dd18e-5845-4ef4-934c-91b3f01eddde	fbb11569-57f0-4b75-ae29-3854e26a0585
20c89af9-9ed1-4124-a318-22f5edbe00c7	8ec3bc88-f2ce-498f-8d01-a40b02804f4c	fbb11569-57f0-4b75-ae29-3854e26a0585
547b4049-ec32-480f-b891-a30239388c09	192ddf18-05b1-4d80-8b2c-923cef89df16	fbb11569-57f0-4b75-ae29-3854e26a0585
daa8e6fc-9c73-4781-8bec-b03d38bf1818	590dd762-650e-4b31-9989-268ae81a853c	fbb11569-57f0-4b75-ae29-3854e26a0585
9323177b-858e-4b01-9bf0-56552a1f7df3	19b78054-0c63-453c-9ca7-bd64c99de999	fbb11569-57f0-4b75-ae29-3854e26a0585
dd5cc7bf-fe73-4c80-949e-251c2e692a33	029a70de-2382-485b-9eb3-15d3933461d6	fbb11569-57f0-4b75-ae29-3854e26a0585
e580bcc0-0e74-48dc-a70c-917f6e0c11bd	315a3bae-c0a7-47df-9522-acbc971ff8b1	fbb11569-57f0-4b75-ae29-3854e26a0585
5b7be9ad-5e9d-4478-98ee-8abca84d7590	1e8a19be-9e51-48af-afc0-d56772e66a24	fbb11569-57f0-4b75-ae29-3854e26a0585
d407c958-614a-43c3-90de-d3d100d4db1e	ad45fb2b-38a8-4d14-8eb5-03af12a2f7ec	fbb11569-57f0-4b75-ae29-3854e26a0585
11749755-a620-437a-b39b-4fca35255d53	267939c7-d951-47f8-aad8-63f0c0d6cdfb	fbb11569-57f0-4b75-ae29-3854e26a0585
987dc327-ed0e-491f-b23d-e6423cfd9c48	a57004e3-4bac-4cc0-b2c6-8b0a060e2dbb	fbb11569-57f0-4b75-ae29-3854e26a0585
02a2cb60-2b53-447f-9fe7-3e07d299c92f	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	fbb11569-57f0-4b75-ae29-3854e26a0585
08b53b88-c3c2-46f1-af31-24dc57f82706	b39c59d0-f9d1-4901-af1c-f3a193d1c612	fbb11569-57f0-4b75-ae29-3854e26a0585
6b914a30-fc78-45ff-9798-8d042527ffc7	f76dca7f-4774-4c76-8aac-873e97f66d3c	fbb11569-57f0-4b75-ae29-3854e26a0585
c7c44f18-9b78-44c3-a7ae-7c61e88707ad	0150f1aa-0dfe-4371-a553-87e73d886b3f	fbb11569-57f0-4b75-ae29-3854e26a0585
cad1d783-6cdf-4521-86d2-623114c18f1c	e23feaf4-58ae-469c-95a3-83e6312e7f29	fbb11569-57f0-4b75-ae29-3854e26a0585
e907c8e4-59f4-4e62-82f5-19c079604b3e	665c7169-f989-436d-a611-e6effc8ddc9b	fbb11569-57f0-4b75-ae29-3854e26a0585
985e3b2a-649e-4b27-b593-e5af9a769afe	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	fbb11569-57f0-4b75-ae29-3854e26a0585
30b5d20d-2485-4717-a2ed-29f301b9570c	fa944b01-d087-4729-920d-f1f7ab36a642	fbb11569-57f0-4b75-ae29-3854e26a0585
5f4808bf-dafd-4bf5-b973-8b8a535731b0	7e62c826-7285-4b26-9235-d4a49b40ef6c	fbb11569-57f0-4b75-ae29-3854e26a0585
9addbf92-2b59-4af8-92ae-37a7888419d9	ccc79a5f-1c03-44c8-a954-2104bbdb7239	fbb11569-57f0-4b75-ae29-3854e26a0585
6a00b1dd-6057-4eaf-bd3a-ee3a16f2a4ab	8b92de65-89b3-4f4b-80d2-a2908bdca9ea	fbb11569-57f0-4b75-ae29-3854e26a0585
b56e0452-88c6-49ea-8812-ca051694eb4f	7586f8fd-b585-4555-98d1-d458115056e0	fbb11569-57f0-4b75-ae29-3854e26a0585
9b47be81-56ab-4aa5-a030-fcf34f5d9287	6419a0a5-1e69-4775-8e72-04ee288635ff	fbb11569-57f0-4b75-ae29-3854e26a0585
e8909478-b86c-47f9-9950-1c53ba0087c7	83334d84-d119-4c76-9292-bb2a34731d9d	fbb11569-57f0-4b75-ae29-3854e26a0585
0ecb7688-582b-4582-bb36-a0d7a6f7a26f	c8d5718b-6e24-43e9-9967-9a6318412dcb	fbb11569-57f0-4b75-ae29-3854e26a0585
6da93419-e5f1-46a4-8f48-a90b509d2f7f	fddf83c2-ea32-4346-8a80-b035758265ae	fbb11569-57f0-4b75-ae29-3854e26a0585
29d182eb-0a58-4d79-9b7d-a94d701fe60b	18e4c9a1-072c-4b67-83ac-8403efd66ff5	fbb11569-57f0-4b75-ae29-3854e26a0585
7063f4bc-1f0f-4158-a7c8-ee5efb29b2fa	70a95de6-a6f3-40a3-a70d-5526d23df7bf	fbb11569-57f0-4b75-ae29-3854e26a0585
9f0fe50a-ec12-4b9c-8176-ccaf93156a7d	b57a6f39-2ec1-4065-911a-29764b00ec2b	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
af41efbc-3628-43ca-908e-5a3163c8cf7b	3af18777-25a3-49aa-9e90-04a45d71a6b9	fbb11569-57f0-4b75-ae29-3854e26a0585
4feb0c3e-27b3-4bf3-9b09-c805819a4f41	3ab532fe-274c-4755-b8c4-8055ea124209	fbb11569-57f0-4b75-ae29-3854e26a0585
c781ff91-4eaa-4041-863a-405883dedeef	2ed683f9-df88-46b5-bbb2-fd7f4d09535d	fbb11569-57f0-4b75-ae29-3854e26a0585
d7f3d77e-0b46-4591-ac6f-8d9cc08bc6ba	b56fa731-9b89-40dd-966d-b92122d6bdac	fbb11569-57f0-4b75-ae29-3854e26a0585
5a65fdd6-e09a-4c68-9680-463f227ca4ed	6c8c8fcf-3ab7-4f08-8f41-5a53c3cc0ed1	fbb11569-57f0-4b75-ae29-3854e26a0585
85c0b7a5-b5ce-4aec-a2a4-bee2c381db51	c5658139-2407-4416-847f-3a5eda3b9bb7	fbb11569-57f0-4b75-ae29-3854e26a0585
6e688aa2-a5ee-4807-806c-cb53aa391ad8	63593179-f9d6-4580-97ca-dd8670688272	fbb11569-57f0-4b75-ae29-3854e26a0585
15a4daac-5ac9-446c-97e6-0b6788befa5d	559d6495-5bab-428d-82d9-4bc1b30e8e06	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
73f1860e-6643-4429-b86d-cc196da8185f	c174cff5-6670-4001-bf92-8395eb9e2656	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
c1f83930-3fc2-4b04-b512-05821843fd10	29058c58-14a6-425d-b366-617860f6950f	fbb11569-57f0-4b75-ae29-3854e26a0585
d17c3cfc-bc68-43f3-8085-9b75136fb581	97a477c0-7361-421f-b042-2d1e0d9bf5fe	fbb11569-57f0-4b75-ae29-3854e26a0585
3737dde9-c551-4baf-9671-dc18021c1ded	60f259e8-d103-4fda-859d-b6c053cc6004	fbb11569-57f0-4b75-ae29-3854e26a0585
c7d7ad1c-5f68-49e2-a7f8-573523d87079	59d15cbe-0f30-4057-ab75-09ba8776b7d6	fbb11569-57f0-4b75-ae29-3854e26a0585
88191a08-eaf6-4a7a-9445-198f4b241511	9b374efe-f7c1-44ff-a122-52981146b22a	fbb11569-57f0-4b75-ae29-3854e26a0585
4d1ae12c-7cce-44f5-a72b-b92f0504f7d8	b36e11e0-99e1-4af5-b83b-6fe9062554a8	fbb11569-57f0-4b75-ae29-3854e26a0585
00664c7e-cf5b-4e79-b38c-d52c498f3a3b	d88a3357-b68c-406e-8951-97df93df66bb	fbb11569-57f0-4b75-ae29-3854e26a0585
688f773e-7fae-4a55-ab1e-821640d1122d	1897bc12-5b09-412d-a3c7-e001aedff81d	fbb11569-57f0-4b75-ae29-3854e26a0585
90f2f01d-4f52-440e-883c-943250604cc8	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
a8df4b34-d9a3-414a-82ce-5d8b19f84f50	4689bc27-a161-431b-bf58-456f66307fd9	fbb11569-57f0-4b75-ae29-3854e26a0585
49367e84-e401-4ff6-8315-a70a3df64fc0	065c9534-8c0a-408f-892a-914d172f90c7	fbb11569-57f0-4b75-ae29-3854e26a0585
a5fa2af3-b45e-43e8-8256-8f455b75a61d	db98b40b-91a4-4434-ba2c-712e3990b82a	fbb11569-57f0-4b75-ae29-3854e26a0585
b88d2ba6-b7d8-401a-b586-4d2e51a20da2	03ed060b-56dd-4ecb-9072-a736c81ae27f	fbb11569-57f0-4b75-ae29-3854e26a0585
c54ce5e0-8989-4337-a909-0b7ec5152551	aa061751-3180-40dc-b101-f866d7e4dc2f	fbb11569-57f0-4b75-ae29-3854e26a0585
f3bc1700-8646-46b4-b71c-c77e9dd4847b	72793652-a596-4b80-bbe4-73d10cca1379	fbb11569-57f0-4b75-ae29-3854e26a0585
bbcec219-4a7d-434b-8d95-fef72608add3	66cb5fb9-5bdc-4378-8010-6b58f8cedb5f	fbb11569-57f0-4b75-ae29-3854e26a0585
904145e5-9310-4b1e-b1b3-0a6a8cc9d9e4	a0dee210-7e5b-4309-ad4e-7ccdf1e7e82e	fbb11569-57f0-4b75-ae29-3854e26a0585
de559a1e-daaa-4ade-8e15-a7e408cbc7d6	ae527d80-3b93-4ec1-a4b3-5c21ab2c80fc	fbb11569-57f0-4b75-ae29-3854e26a0585
6b6238ed-e1b2-4588-8730-22a3f96d0ee5	6b7a79ac-30dd-44ce-a6bb-82a137f67aaa	fbb11569-57f0-4b75-ae29-3854e26a0585
755acc7d-fd7f-44ce-8375-7b08c2082453	230733a9-34d8-47b9-aa97-83a85dbca11e	fbb11569-57f0-4b75-ae29-3854e26a0585
7f822ab1-7d21-429e-9744-4bd00c86e73a	805c95bb-9f47-4e74-a1f0-87c3bb12a455	fbb11569-57f0-4b75-ae29-3854e26a0585
14d0d453-10eb-4683-adc0-6349df64d516	3d74d6a0-376a-4ff2-a1fc-5150b624e769	fbb11569-57f0-4b75-ae29-3854e26a0585
ed0adaa4-39c3-471c-bf44-c5226e3d8539	e1ea71e4-fcbd-4ca5-8917-1dcab59fbd45	fbb11569-57f0-4b75-ae29-3854e26a0585
1b7aba0a-01c8-44a8-9287-a0802a7e0acd	9def3d8c-3f3d-4a1e-9020-6b382fc42310	fbb11569-57f0-4b75-ae29-3854e26a0585
91633d4a-04e4-4fe2-95b5-a8ab8588d9d0	c2731296-80f6-493f-acd6-3739949e0ad0	fbb11569-57f0-4b75-ae29-3854e26a0585
873f672f-aa41-47d7-97f4-de204d3d25ee	8caf4478-ee5f-48f2-b0cb-a35b2bdc8aa0	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
082d336b-fd92-40c4-8136-a4f6bd402b55	8ec4ba0d-da46-49e0-a349-cbaa7f13e26b	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
07a7d7db-0f25-4ed7-863d-a16b69bd76d8	4652e35d-14d1-4d5e-8c34-dd691500ce69	fbb11569-57f0-4b75-ae29-3854e26a0585
8b5c15b5-eafc-46cc-bf55-53819e28342f	5e5c4e3f-155c-48f4-9aa6-dc7dd88fa56f	fbb11569-57f0-4b75-ae29-3854e26a0585
41e9fef3-42c1-449f-80ec-5b40919712ce	23ee37f8-d777-4fdc-9d2b-8705c72ce357	fbb11569-57f0-4b75-ae29-3854e26a0585
2ae4ad89-2d65-478b-9757-5f1aa28171c0	8964e926-7f17-4370-bfe6-83cf363f5244	fbb11569-57f0-4b75-ae29-3854e26a0585
66580f76-1bd5-48cc-9149-91d6a474f4e2	6eee51ef-8a8f-4966-b19d-c9b915e681e4	fbb11569-57f0-4b75-ae29-3854e26a0585
8b0dd84a-8bae-40f1-b1aa-31c2ff9d386a	5219e94c-69eb-4baa-86ac-593d0d0255fc	fbb11569-57f0-4b75-ae29-3854e26a0585
063996ca-eb3f-4904-9ad8-ec9e5f171a87	ea6a4e30-d810-41fd-b231-a9a209adbae3	fbb11569-57f0-4b75-ae29-3854e26a0585
9bdb3f4e-f174-43cf-a460-33328a3e7cbf	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	fbb11569-57f0-4b75-ae29-3854e26a0585
1c2fa0c3-20f9-4a2e-90e9-2cbe571e4341	5e6cb690-6a8d-4fab-9234-754db94e1c2f	fbb11569-57f0-4b75-ae29-3854e26a0585
33e5871d-943b-4880-8383-a90ecf7f4daf	afe5da52-c330-4f13-b7a6-386cafd6fd0d	fbb11569-57f0-4b75-ae29-3854e26a0585
d9390789-ecea-45c6-8eae-cb29aec4b08a	9ec25419-48ea-4525-8832-2f9cbcdc3096	fbb11569-57f0-4b75-ae29-3854e26a0585
a6afb7ff-f79f-4def-b6fa-f32653cbe310	5a08d43a-1268-4b70-8514-60fccd917026	fbb11569-57f0-4b75-ae29-3854e26a0585
72356440-9331-4e09-8e59-370603460da4	a761fc46-09d5-43aa-91ad-93391ae4934a	fbb11569-57f0-4b75-ae29-3854e26a0585
4a0de66b-2981-40b9-bdff-fd4fc0efab16	7628af1d-b868-4923-8bac-a7366efdc813	fbb11569-57f0-4b75-ae29-3854e26a0585
f7cd6be5-64d9-45d5-a2d7-5684ea9d7c29	7872d920-bf48-491f-a604-f29239826d86	fbb11569-57f0-4b75-ae29-3854e26a0585
df2dc9d1-da1e-4487-90fb-e8696742b224	fa083891-8c7b-4d3e-a769-9c765710f79a	fbb11569-57f0-4b75-ae29-3854e26a0585
01c80382-2a2d-4e63-bcb3-4a8bdbb7f03f	973f37a3-9aa6-48a0-8220-66b1444e81ab	fbb11569-57f0-4b75-ae29-3854e26a0585
1ec4e87f-a00a-44a3-a7f5-a83b37695d03	43b61baf-0753-4367-9d29-7be6623d91c9	fbb11569-57f0-4b75-ae29-3854e26a0585
b31c9020-eb31-4a44-ae53-748f5f8b1960	4aacad3f-4671-4f8f-a6bc-73ce21233408	fbb11569-57f0-4b75-ae29-3854e26a0585
4034c5f8-f5ac-4240-a960-ce67d593e22a	d0f0d5f9-e375-4c11-aaac-22ae71a637ea	fbb11569-57f0-4b75-ae29-3854e26a0585
615d731d-ca7a-45d0-9543-069b3fb1f0e0	b207a77e-a3b3-4303-b757-2bd70ca8639b	fbb11569-57f0-4b75-ae29-3854e26a0585
cfefbc1e-44b9-4e70-bc3c-80dd611122e1	adcbffee-495f-4e59-90c0-f438d39ebd07	fbb11569-57f0-4b75-ae29-3854e26a0585
3d3e9af6-e390-4bfd-aab2-6140088d5249	6afbed32-4c43-43c7-b300-c1a04857d346	fbb11569-57f0-4b75-ae29-3854e26a0585
7740c6d0-5b6e-4000-ab65-4826d7a3d493	4f387a47-6b08-48f6-981c-42a8a128d1f0	fbb11569-57f0-4b75-ae29-3854e26a0585
c07e8181-1110-4554-91ce-dcff0c965bd8	a62b7765-f024-4705-877a-de0ad4878d19	fbb11569-57f0-4b75-ae29-3854e26a0585
e48db442-3d56-4be1-973a-072ec94a970c	2f4bda42-a624-4eef-9cc0-aa581733206a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
100dc0ca-11ea-4599-a4f1-7ebc7d826f8d	44a5389b-8608-4865-ab53-74be807043a2	fbb11569-57f0-4b75-ae29-3854e26a0585
c60e17c0-ee63-4a84-b5ad-7e41b823d963	6c154b5d-4689-4282-a54e-dbbeeeba5fb5	fbb11569-57f0-4b75-ae29-3854e26a0585
1facbcc4-ec50-42d7-8c70-426dc7ac7eca	a1951306-b89e-46d3-8768-3d681585e382	fbb11569-57f0-4b75-ae29-3854e26a0585
bf966f6b-277c-4ed6-b94a-95271e3d677d	55327a19-fd5c-4101-885c-ec17508e5e74	fbb11569-57f0-4b75-ae29-3854e26a0585
7606a80d-958a-43a4-9cc1-484f8dc2355d	deb75efd-1d47-4fb3-88d0-4fe07c5cbf36	fbb11569-57f0-4b75-ae29-3854e26a0585
339e4fab-7b90-45c2-85c5-97d4f12b8527	96deec3e-55a6-422c-9cb5-3f59ff7a3dd2	fbb11569-57f0-4b75-ae29-3854e26a0585
b1e59bfc-fbdd-41b7-ba36-f029ae45c17d	7e87afed-a5ab-4285-b257-b84161febe58	fbb11569-57f0-4b75-ae29-3854e26a0585
7cd2e9fc-cf3c-43f3-930f-f16692068f4e	622186a1-17b8-4b64-8406-dd3cd1dd2596	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
9b89f490-a330-4803-b5a5-28f394b1d077	d1f06d6a-a377-4419-b690-1cbe09879b1c	fbb11569-57f0-4b75-ae29-3854e26a0585
81c97751-6235-4bbe-a0d5-802a994a4ec8	ff6e92a2-a7a4-4508-9d4e-bde9f6cb0bfe	fbb11569-57f0-4b75-ae29-3854e26a0585
5a314406-073e-44f4-8ae2-32915a6fad5f	5f65f5d0-39d7-42e6-a7f5-dd21ba72abb1	fbb11569-57f0-4b75-ae29-3854e26a0585
278263f6-ad78-413d-b20b-c5e99533b14a	8feea562-c412-4546-916a-7c22e60c833c	fbb11569-57f0-4b75-ae29-3854e26a0585
803a7a41-8044-43b0-95b3-847cecfb7ac7	a35f9884-abfa-4c87-9c23-a983e04dfaf7	fbb11569-57f0-4b75-ae29-3854e26a0585
97635b4b-aeae-4ed6-b8e9-cb5f10ae8331	15362ef8-4d86-4e96-b49f-c33c5292e67f	fbb11569-57f0-4b75-ae29-3854e26a0585
fa64c2d8-948d-44a6-af7a-5936e2f3b1ec	b639edba-fe15-48c6-9e53-551c056a4426	fbb11569-57f0-4b75-ae29-3854e26a0585
b15333bf-8230-4f8f-9611-ab3fdaf20437	e3eaddf3-edb3-4930-a548-bce70c834b66	fbb11569-57f0-4b75-ae29-3854e26a0585
a7d324c7-9831-4aa7-8a94-b97782e66649	7ac101de-041e-4baa-ab7e-7471318ec81f	fbb11569-57f0-4b75-ae29-3854e26a0585
1ea8edec-1e53-4cac-b9bd-0d1b808dae50	00f047a0-8afe-40a6-89d5-74867de6cfbe	fbb11569-57f0-4b75-ae29-3854e26a0585
b350fff6-0ce3-4b0d-b7bc-4afbf0368d2e	583c9de5-4079-4eb9-b448-4fa66ac594fd	fbb11569-57f0-4b75-ae29-3854e26a0585
1a7faf06-a150-42e2-bbcb-b1cf915ca138	ab03d85a-45db-4af0-8b37-0de7e8bfe4a6	fbb11569-57f0-4b75-ae29-3854e26a0585
488e0c9b-4ab9-4b98-add4-8746006f4af6	2feb45eb-2e97-4d7f-bf43-d221b1c71928	fbb11569-57f0-4b75-ae29-3854e26a0585
4f7bf5d7-fb95-481b-9f93-e09b30a67bba	779dc403-d84a-4e4e-9aa8-3d7cbf013a69	fbb11569-57f0-4b75-ae29-3854e26a0585
a7bb39bf-f250-4ad2-a753-59db2db0dc82	115e00a4-dd87-4ebd-a471-9a69f5f3f85f	fbb11569-57f0-4b75-ae29-3854e26a0585
9f989c0d-1ff0-422a-9847-5370076337b9	44bda315-99d8-4619-99c8-877e24f259db	fbb11569-57f0-4b75-ae29-3854e26a0585
6c436433-8f74-4465-a758-aa8c51109fa9	e4db0d2c-9735-4e9a-b698-c4a395224d48	fbb11569-57f0-4b75-ae29-3854e26a0585
49ad7cdd-1837-4f62-be49-bce92a80eba9	2697141a-6da7-4cba-8237-9cc42a279742	fbb11569-57f0-4b75-ae29-3854e26a0585
8b0f785b-b88e-44b6-8160-987771d6d6ab	20ed729f-8a6b-47f0-8b79-6cd17dfa5260	fbb11569-57f0-4b75-ae29-3854e26a0585
cb32209b-2266-46b1-8949-bc407607de11	f7c21ada-b21e-4ce8-8a52-df88fa170b5d	fbb11569-57f0-4b75-ae29-3854e26a0585
6282bf93-637f-4470-9412-985a4103177f	90b65d06-171a-4975-afd6-d226f3a35910	fbb11569-57f0-4b75-ae29-3854e26a0585
5a485a7f-0fba-41bb-9ff2-9d0c8ce2e5cc	91b45a4d-262a-4a5f-a376-845bb5b8e489	fbb11569-57f0-4b75-ae29-3854e26a0585
78347247-7b7c-43eb-b209-e425c665f7c9	9392bd8f-6ef5-4f16-b87b-a36709339d5d	fbb11569-57f0-4b75-ae29-3854e26a0585
18a25930-e0bb-4569-945a-e8ce5e863f56	b92b6639-c7fb-42ff-91bc-3281971224b2	fbb11569-57f0-4b75-ae29-3854e26a0585
e3a4ab26-2397-4fe6-87a7-a5864f6f08c3	ec64aed7-c889-4d44-a3de-38b112f490ec	fbb11569-57f0-4b75-ae29-3854e26a0585
3e14a4e1-6b14-4b73-b85a-7a2c22eceffb	a4beb80a-97bf-46c4-81ab-af3735050855	fbb11569-57f0-4b75-ae29-3854e26a0585
b9291c0d-bb0c-4e71-8c56-112344e2ba54	01ba1a7b-1480-494c-8c0e-ca6b34d37d1e	fbb11569-57f0-4b75-ae29-3854e26a0585
9762d274-c51f-4f91-bdca-ab19b6051e58	3fbb2d4e-c39d-41e3-aee8-079228318790	fbb11569-57f0-4b75-ae29-3854e26a0585
511220f6-b456-458d-8b3a-336441e86350	eb467d59-fea6-4d48-afde-3011582cef2f	fbb11569-57f0-4b75-ae29-3854e26a0585
6d5c9d70-3990-4bc7-87a1-7ec1687007dc	a8de69fd-5856-42cc-8812-fef457a99079	fbb11569-57f0-4b75-ae29-3854e26a0585
3c16dc63-c0ca-437d-ae85-3e52f9f6df8d	4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56	fbb11569-57f0-4b75-ae29-3854e26a0585
ffcd6a12-0b32-4c0f-a37b-2bc3ad6ce22e	a2bcc31f-515a-45e9-8083-b2f27d2c4afa	fbb11569-57f0-4b75-ae29-3854e26a0585
8bb7b5b5-6127-4d17-a5c2-11137a323d67	46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d	fbb11569-57f0-4b75-ae29-3854e26a0585
7df62dd2-7bd8-41f9-88d3-15032540ddd5	88c58e99-19a7-4f72-811e-92e21346207a	fbb11569-57f0-4b75-ae29-3854e26a0585
c42b0706-31e1-4532-8fb2-4f8d8ca7bad9	eedb53fb-3a3e-4f32-a232-ffc97a4464ae	fbb11569-57f0-4b75-ae29-3854e26a0585
4935bab8-c90a-420c-9600-1ade9a4e301b	8ec08232-29d6-418f-8560-64f307684f8e	fbb11569-57f0-4b75-ae29-3854e26a0585
29fb6dce-277d-4d96-bf0a-a62cb2bdcad1	cb332a03-aeee-435c-86de-b7d577e3817b	fbb11569-57f0-4b75-ae29-3854e26a0585
89de4cdf-41b8-46f6-a8bf-fefaf5416ad7	357d5bae-8470-4d03-ab92-68cafedcfa16	fbb11569-57f0-4b75-ae29-3854e26a0585
1f671eab-9173-4c7d-b8f7-214d0623ce84	ef82cdff-1b9d-4396-8d2b-f7b751a50006	fbb11569-57f0-4b75-ae29-3854e26a0585
0e794f09-d089-4352-b0ee-5c1a24853b25	ad55ef62-54a6-41c0-bba9-830a06822659	fbb11569-57f0-4b75-ae29-3854e26a0585
ac8468d0-c50c-4055-b5cd-ec1677d4e3e9	b995956c-75be-4635-9022-af1f5376be2b	fbb11569-57f0-4b75-ae29-3854e26a0585
a91a987d-9ab9-4188-938a-42a470f68f00	e18c306e-1d61-4e6e-a2cf-c41ba357b406	fbb11569-57f0-4b75-ae29-3854e26a0585
a1497ee2-30d4-46bf-928a-4c0b0a3f55f9	b4aee72a-2ac5-4748-b846-bd1828139858	fbb11569-57f0-4b75-ae29-3854e26a0585
92af3203-014c-4e96-bf7e-7c3f732b8ed4	5a8b13c6-3997-4846-b037-2b679f5c31e8	fbb11569-57f0-4b75-ae29-3854e26a0585
bd495876-7d82-4137-82f4-7ca0a44b880f	3293247e-3d3c-4749-ae91-1a991e27b8be	fbb11569-57f0-4b75-ae29-3854e26a0585
2da4778e-7deb-4330-8a32-b286b8eec5d0	7ea8ce89-9ce7-4a89-a332-7baad7bbe475	fbb11569-57f0-4b75-ae29-3854e26a0585
82d2b94b-8244-474a-9835-4a6177032b3f	742c55df-02ca-4d3e-909a-ac8c412ff735	fbb11569-57f0-4b75-ae29-3854e26a0585
999a57bf-a192-49d0-9f94-5c645e8ace23	48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5	fbb11569-57f0-4b75-ae29-3854e26a0585
51da338b-00dc-44d5-8560-4126c2969435	23476194-48d6-447c-b4e9-2ab6ad7b1a48	fbb11569-57f0-4b75-ae29-3854e26a0585
080be39d-02f3-44b8-969a-62573bb8889d	cab6e3e8-cff1-42da-959f-e12705adc760	fbb11569-57f0-4b75-ae29-3854e26a0585
23181758-c56d-4869-a5ea-587b58fc129d	01698813-a221-4b56-8bf5-ec3ec9956d97	fbb11569-57f0-4b75-ae29-3854e26a0585
0ae47bcd-c61b-4c46-866e-a567fc953bad	3c53a6aa-9d98-43be-bd14-4c5501e58ee5	fbb11569-57f0-4b75-ae29-3854e26a0585
e9d46af2-48b5-4bf9-a1d6-b50b6d0d679b	080f0752-21ba-4d3e-97c2-4846dfb4f648	fbb11569-57f0-4b75-ae29-3854e26a0585
0477a05a-fd5a-4e1b-8e45-039bc2c90076	d470cd75-937c-4229-b180-7b0ace983442	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d236de36-8261-4846-83cc-51ecfbdc29e9	8d94bbc5-5eb8-4d3d-890f-3a3948b29ce9	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
8c7b51bb-5696-4354-9bb6-61666e814e0f	35d43dd0-6e55-4bcc-ba62-efc1c4cdff5d	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
4e727de7-8ecb-4d3a-9429-482ff785a079	e0be0b73-1f90-4c04-8da5-86513c148a16	fbb11569-57f0-4b75-ae29-3854e26a0585
2e4f485f-84c1-483d-abf9-ccb6d8297a63	af0187d0-e84c-401f-9ffc-70194b099309	fbb11569-57f0-4b75-ae29-3854e26a0585
822561db-738e-4910-b8e7-de9e9f9e265f	77788f4d-5d0e-458b-b7cb-297e907b2eff	fbb11569-57f0-4b75-ae29-3854e26a0585
a987d99c-7482-4282-a8af-afcfc66c8ee4	eb9ea9f9-cd44-4b11-a805-061612ff12ef	fbb11569-57f0-4b75-ae29-3854e26a0585
a7bb50d8-0aae-41a3-8b41-2dee97b6b12b	b2776c17-fde6-4af5-8d34-f72d5ec88ac9	fbb11569-57f0-4b75-ae29-3854e26a0585
2ddee019-3736-45e0-b551-c4bc9bc0c4fc	fa13b0c6-1e65-45de-bef8-ac534342f0f3	fbb11569-57f0-4b75-ae29-3854e26a0585
5e07f892-c263-457c-9aef-b4dc9bff757b	14f49944-6337-4c1c-8ea4-dd111c7859f1	fbb11569-57f0-4b75-ae29-3854e26a0585
9883107e-dafb-466b-961c-f0475ac774c2	8c9f5d83-fb7b-4f21-8b60-1ed3eb93ff41	fbb11569-57f0-4b75-ae29-3854e26a0585
aa3ac765-df9b-452b-a807-5f58206fba60	222d3e77-c07d-4c30-a384-b0c493a3cb38	fbb11569-57f0-4b75-ae29-3854e26a0585
c05f3748-b5b3-480c-a1e4-bcfdbfb6d2e0	824f43c8-ee42-47ed-bbb9-5f8de09acdbc	fbb11569-57f0-4b75-ae29-3854e26a0585
8ac26aab-de26-47b8-96ac-dfb8579d1ca4	61ab5266-4740-49f3-9649-1b97566a4191	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
7b2de00d-36ab-4fc8-a154-5aca125eed8e	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
60fa0484-4e16-4bad-a6ae-46dd65d6b250	6e4f00e4-bfad-42af-949d-8b215dae11f2	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b1e8e84b-4a77-462f-bc26-e5d06d5906c2	ed06f369-9d92-4d07-8b9a-6488694b4d1c	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
599ccb24-b885-494a-b6ab-128a030c4f28	977f77aa-139c-4d4c-8640-053eea3376dd	fbb11569-57f0-4b75-ae29-3854e26a0585
e8c6113c-5dba-4024-9470-9843a5169e3d	7270ed76-fcee-4a5a-a37f-6645f69e68aa	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
72534f7f-61f3-4e31-b971-409661b8d956	621bc69f-2dcf-4754-bcf8-71a18d17fa7e	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
a0d0d490-0255-48ce-b5ad-60b0635f0e9c	0dc8b5f7-da88-4bdc-ba24-b5a6bc8a1044	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fb5b4865-ba3a-4239-ae07-f3e6b4625ffa	bda3fd00-36a5-43bd-9c2d-0a25511217f3	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
87c21f77-bbe7-42b5-b95c-218c644c63ed	07722dd0-088c-4309-9fd4-7097cb7a5c58	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
933fa0c3-869b-4eff-a9b7-2faa54795fa2	0e143513-3b89-4e2d-9865-8a848ad6c8db	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
7c14f929-517b-43d0-8935-b56a9578f1af	07722dd0-088c-4309-9fd4-7097cb7a5c58	fbb11569-57f0-4b75-ae29-3854e26a0585
abb7eb38-ed0f-40de-a832-37c90355fc98	0dc8b5f7-da88-4bdc-ba24-b5a6bc8a1044	fbb11569-57f0-4b75-ae29-3854e26a0585
bb77de28-6de7-4ee4-9866-a3cb6f809e4e	46ea3bb2-76b6-4842-a015-33c9edc910ae	fbb11569-57f0-4b75-ae29-3854e26a0585
76ad154c-b047-448e-9bb7-70f00539f3b1	bda3fd00-36a5-43bd-9c2d-0a25511217f3	fbb11569-57f0-4b75-ae29-3854e26a0585
b0ecbad0-c1c9-408f-9556-7b07f0885293	b3e49cac-fd53-4dcb-97ea-03ac3c328405	fbb11569-57f0-4b75-ae29-3854e26a0585
cdf6ec88-dc4e-4927-9833-490ab89af20b	8b470f68-065a-4b9a-a35e-0781d83910df	fbb11569-57f0-4b75-ae29-3854e26a0585
e480efe2-5b46-4607-a005-068d87d08dc1	cc9a6c3a-518a-458a-a84e-6353a1c22165	fbb11569-57f0-4b75-ae29-3854e26a0585
1fd7e9e4-107a-41c0-a4b1-d0e56604df01	e9cb3b58-b3ed-4e27-a2f5-0c704dbbf32f	fbb11569-57f0-4b75-ae29-3854e26a0585
6ebd762a-edbd-42ab-b120-f2fb401c5866	b3f2d37b-3f4b-4fa6-b632-3a85e295fd01	fbb11569-57f0-4b75-ae29-3854e26a0585
5486d3ad-083a-4618-89c2-b7f4c7f6e443	57873320-1f3c-4d56-ad41-122f17f362a2	fbb11569-57f0-4b75-ae29-3854e26a0585
0d655bed-077f-43b8-ac97-0a80ccd87ab5	a168235d-63c0-4b58-8b49-3df20a9a815d	fbb11569-57f0-4b75-ae29-3854e26a0585
1f9e2d26-c0ef-466a-a7ef-ff1662f56180	c2592cf2-b4f8-4ab8-9ce3-54b4dbbc928a	fbb11569-57f0-4b75-ae29-3854e26a0585
8b7f74c6-725f-4b03-9ef6-ec21b36db4ee	f42406c3-08c8-443a-845a-996e04c55753	fbb11569-57f0-4b75-ae29-3854e26a0585
8a5f1f5b-42cc-4448-ae08-ed710266f65c	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	fbb11569-57f0-4b75-ae29-3854e26a0585
7a7d6db3-921f-4821-ae2a-f5dce38a3fd4	6e4e2624-bc76-49dc-86ba-6e88b0977c79	fbb11569-57f0-4b75-ae29-3854e26a0585
542a4fda-d651-48ce-925f-a57700bb3217	86f767d3-aec7-4acc-a3b3-e63b1734943f	fbb11569-57f0-4b75-ae29-3854e26a0585
46e0dc16-70bc-47dc-a19e-ef2431c0ce7b	c8b42420-731e-4736-97a0-81725e9aea49	fbb11569-57f0-4b75-ae29-3854e26a0585
75888241-7067-496a-a9ba-4746b8007065	2c1d69a3-4af1-4436-9a70-449ba88f3048	fbb11569-57f0-4b75-ae29-3854e26a0585
ac2470af-c6aa-430d-96e0-1e5626ace42e	0ad5ba51-e406-477c-bdc2-fe35ca35b6ec	fbb11569-57f0-4b75-ae29-3854e26a0585
30254681-5ac4-4110-93e5-de017c0884a5	b8e3f794-9775-4363-9bcf-d3d929f6ba27	fbb11569-57f0-4b75-ae29-3854e26a0585
ba625744-5974-494f-8a4c-d110bc8e6d16	ee775cf8-f617-4324-9220-726c7ddbe919	fbb11569-57f0-4b75-ae29-3854e26a0585
9f018889-430e-4522-9034-c5fed95f5b96	15668ec7-6731-4646-a516-41035ddebe0d	fbb11569-57f0-4b75-ae29-3854e26a0585
46b8589f-fff9-4250-98e2-42dba7954b8d	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	fbb11569-57f0-4b75-ae29-3854e26a0585
01f47e42-da06-4e32-8542-848a900a20a6	ccad896d-13db-4fc8-be92-574af1d80341	fbb11569-57f0-4b75-ae29-3854e26a0585
c7df2c0e-b5dc-4bbb-a978-75f705dffbb5	51f985e0-4f2a-412a-adb5-3adf5399fc75	fbb11569-57f0-4b75-ae29-3854e26a0585
25773714-7193-4220-86df-865450054434	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	fbb11569-57f0-4b75-ae29-3854e26a0585
455df8a7-19a5-4a83-a456-bb0da7663f40	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	fbb11569-57f0-4b75-ae29-3854e26a0585
6818ff05-fbb1-438f-93b6-12508bd2b03e	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	fbb11569-57f0-4b75-ae29-3854e26a0585
d5172f78-396f-4985-9866-c1ef62ca921d	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	fbb11569-57f0-4b75-ae29-3854e26a0585
60ae2886-c91d-4d73-9520-d308e600ef34	b48dccdc-da4f-4914-bd2e-1971de87fda1	fbb11569-57f0-4b75-ae29-3854e26a0585
a91f310d-f331-4180-ae1f-e205c309b78c	27917132-65e1-4218-bdfc-40a62550a141	fbb11569-57f0-4b75-ae29-3854e26a0585
337d14bd-79d2-4548-a99a-25d6c6163f4d	7767088d-9354-42d4-8cc7-9d3f2c0bef88	fbb11569-57f0-4b75-ae29-3854e26a0585
434b148d-474b-4291-a5d2-630b49f839e9	f22b7a3d-2d76-483c-af40-dd0c04defb20	fbb11569-57f0-4b75-ae29-3854e26a0585
ebf9dc59-9310-4d40-bf01-f12a5f72eafa	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	fbb11569-57f0-4b75-ae29-3854e26a0585
8b32d135-e9ca-4589-b2e2-b8b2b9a9d7e1	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	fbb11569-57f0-4b75-ae29-3854e26a0585
47bd66c8-7caa-477b-8574-213d7cb517f4	24127ef9-9ef3-46f3-88cd-97378e694a6c	fbb11569-57f0-4b75-ae29-3854e26a0585
29ee1f01-8ca0-48ad-9f40-bf3dfc2bf2ea	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	fbb11569-57f0-4b75-ae29-3854e26a0585
8f1e31c6-cef7-42c0-9e72-71b80e6939e9	c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	fbb11569-57f0-4b75-ae29-3854e26a0585
84e44c50-c94b-4bf1-a2dc-92bdb5228e29	288731cb-561c-4d4a-8951-d8c359a51846	fbb11569-57f0-4b75-ae29-3854e26a0585
726e455a-3857-4977-a563-e6ce4dc50949	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	fbb11569-57f0-4b75-ae29-3854e26a0585
5f78d249-8c17-40cf-aa74-af5f3ef5c6bf	b439edb9-605f-47dc-aee2-8df45976736e	fbb11569-57f0-4b75-ae29-3854e26a0585
1d1176ea-e8f8-4c72-9a70-c223203a7b09	6f4cc747-5450-4514-9a3e-123ace9cb232	fbb11569-57f0-4b75-ae29-3854e26a0585
0205ad9a-fa9d-45b3-9f53-2a33507a2270	c070d24f-76eb-4684-8459-dbb7977acc46	fbb11569-57f0-4b75-ae29-3854e26a0585
4a9c5ee1-6c91-4038-8e7a-64f5eee2eb3f	d81d630f-3d80-4168-923b-65eefae629c9	fbb11569-57f0-4b75-ae29-3854e26a0585
579a9cd1-8f85-434c-b746-8fffd983d35f	b3e49cac-fd53-4dcb-97ea-03ac3c328405	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
37dfabd4-d5d0-4b00-9fe0-6077cec6c46b	ae18fd5a-e0ba-4ace-a637-902317a67773	fbb11569-57f0-4b75-ae29-3854e26a0585
56c81a37-2894-4849-a140-93e03f933a17	223c3afe-ef65-4b86-9733-6284396cd65f	fbb11569-57f0-4b75-ae29-3854e26a0585
eba4e369-63c9-4387-9e3f-1cc91a66899f	78cd88ae-d982-48c6-998c-1a6ba961b5f0	fbb11569-57f0-4b75-ae29-3854e26a0585
15617a40-62ff-4548-a5ff-0907076471ee	20805d38-e63d-4117-b902-7007c1697a13	fbb11569-57f0-4b75-ae29-3854e26a0585
9002963c-9060-4ef0-9830-e9909b0e38b9	99a0afca-2fef-4083-82a2-4785ac933ec0	fbb11569-57f0-4b75-ae29-3854e26a0585
fe344ac6-7574-4b92-8c4c-a227006983b1	a06b41ea-3abb-4dc4-a4f7-90bfc86a9592	fbb11569-57f0-4b75-ae29-3854e26a0585
9afb2e5c-c671-4807-a43f-6e61ff791aca	25c7b3c0-efd1-4d00-a748-e37de9ca4e3c	fbb11569-57f0-4b75-ae29-3854e26a0585
92c2aef5-b3ba-4c58-bbe7-9c70497f6279	fcea9f1c-b882-4368-a7a3-c8bbd2491684	fbb11569-57f0-4b75-ae29-3854e26a0585
d47b1bf4-1397-4d84-8ef8-7c0b88071d99	6b7a79ac-30dd-44ce-a6bb-82a137f6771b	29242052-f1fd-458e-9119-3fa6a9d5ce98
7a181a67-2fc7-4be3-9e4c-3005e82c60c7	43b61baf-0753-4367-9d29-7be6623d91c9	29242052-f1fd-458e-9119-3fa6a9d5ce98
66be0eeb-586d-4249-b517-94aa66518100	6f13190f-abaa-48e1-a490-8d4406df0f34	29242052-f1fd-458e-9119-3fa6a9d5ce98
3f62812d-a2a3-414a-a775-669cd209df63	fa083891-8c7b-4d3e-a769-9c765710f79a	29242052-f1fd-458e-9119-3fa6a9d5ce98
635c5e82-fd60-409f-8ce8-fef2f8dc45e4	5a229d4e-beee-4f0a-ba2f-888f04a149e5	fbb11569-57f0-4b75-ae29-3854e26a0585
a848b0f5-a288-4dca-9b29-95e7963894c8	c5e2efa4-c8c7-4c79-9823-59e829a4847e	fbb11569-57f0-4b75-ae29-3854e26a0585
d3033ffb-ff9c-4fa8-98c1-838d9ae42f21	90df4bc3-dbb0-4e92-ab9a-0aac895e5295	fbb11569-57f0-4b75-ae29-3854e26a0585
1fbf9b21-278e-46f1-a3a1-a5caa348e8dc	4e471ceb-766f-4575-ab73-856e07ceb1a8	fbb11569-57f0-4b75-ae29-3854e26a0585
66c45b9b-480d-4c94-aebc-9f4bbc96eda4	1c80833a-d9f9-4b31-ac77-55c282781167	fbb11569-57f0-4b75-ae29-3854e26a0585
8f54c5d7-c4a7-432c-8e3d-e82a36df64e0	0576302e-61ce-49b3-b265-28b7458e5c27	fbb11569-57f0-4b75-ae29-3854e26a0585
24e253df-8a8e-4f79-b3bb-bf3fa5b72d73	3ed50a46-3620-43a1-b8ea-5a1ebc968851	fbb11569-57f0-4b75-ae29-3854e26a0585
4dcb08cf-ed2b-464b-8728-1f743b789ff1	82dbdde0-be98-462f-bef7-b86ed1a98336	fbb11569-57f0-4b75-ae29-3854e26a0585
17bd1d2b-c685-4ff9-a0cf-1115da368905	c2ebd879-c09a-43d3-8088-6b0c2445c0b1	fbb11569-57f0-4b75-ae29-3854e26a0585
21c052e0-5b78-4d47-b1ed-824fe4d6bbb1	e296b896-76d1-4cf0-b092-e637a5909ec4	fbb11569-57f0-4b75-ae29-3854e26a0585
9d67c316-6e24-404c-8b53-ed98a4ffdb2e	136f2190-6f80-46d4-a124-d512e508be5e	fbb11569-57f0-4b75-ae29-3854e26a0585
b624b113-17c3-433e-be55-f537d9db03be	5af52e1b-b6a1-4eca-87ec-466abacb4302	fbb11569-57f0-4b75-ae29-3854e26a0585
49eb9b64-8a9b-4de5-827e-e6c92b8bed13	4badf803-d758-47a4-8a4b-fb005cf7167a	fbb11569-57f0-4b75-ae29-3854e26a0585
f133ef0f-4760-4714-9aae-36afbd4390f9	aded78fc-0d49-436e-a780-ebea7c8b0f4c	fbb11569-57f0-4b75-ae29-3854e26a0585
97182a11-a25d-4d1c-b16f-b9bcb0ed1d8e	622186a1-17b8-4b64-8406-dd3cd1dd2596	fbb11569-57f0-4b75-ae29-3854e26a0585
de6ac2c8-c88a-4ac9-9242-6d4ffa79e94a	07bd2d37-21cd-43b3-9d38-eeafaa73bcf9	fbb11569-57f0-4b75-ae29-3854e26a0585
fe64573e-e223-4a74-bc30-12e1f09d9a9a	61ab5266-4740-49f3-9649-1b97566a4191	fbb11569-57f0-4b75-ae29-3854e26a0585
7dce9816-168f-4df4-bc03-a03e7a343f9c	90476856-94be-40ed-8b1e-814392074d81	fbb11569-57f0-4b75-ae29-3854e26a0585
f8f9e9d6-8249-4231-8a10-a608178187bf	7d89470c-955d-4ad7-bee4-82b65eaa07c0	fbb11569-57f0-4b75-ae29-3854e26a0585
90b61fda-967d-407d-8f9e-187e294e5cf8	4d3344aa-04e1-4b24-80e6-ace57f7e3e6b	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
89226027-6679-48ee-bb76-408034b75089	fcc9c518-48c0-4c0e-b96a-8d93b165d948	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
82a347de-ec77-4e4b-8283-40809a1cecd5	439becc5-2668-4f15-8c7f-1b544989df47	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0aefa864-8fc1-4ef1-b49a-b6acceb163af	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	fbb11569-57f0-4b75-ae29-3854e26a0585
bb395eb5-a584-42cf-bd7e-1173aa56d05b	ba6556cb-4288-428d-ac87-a67546dd70be	fbb11569-57f0-4b75-ae29-3854e26a0585
03767a08-f01f-405a-b205-b81e529abe44	621bc69f-2dcf-4754-bcf8-71a18d17fa7e	fbb11569-57f0-4b75-ae29-3854e26a0585
98252b61-781d-49ae-9782-ad353550637b	0e143513-3b89-4e2d-9865-8a848ad6c8db	fbb11569-57f0-4b75-ae29-3854e26a0585
467169bd-90c0-483d-a029-010418a9fc5d	761f6268-1ccb-4a23-b43d-cad43d4761d2	fbb11569-57f0-4b75-ae29-3854e26a0585
63a0c5e7-8cde-4b98-b9ac-28237fd3fd04	7270ed76-fcee-4a5a-a37f-6645f69e68aa	fbb11569-57f0-4b75-ae29-3854e26a0585
4ed7f200-4f93-471e-95f8-91158b88822b	6e4f00e4-bfad-42af-949d-8b215dae11f2	fbb11569-57f0-4b75-ae29-3854e26a0585
4c7ec014-deaa-4470-9eb4-b8b87a2f3019	18460b2d-6e29-43da-9419-990fd196de7c	fbb11569-57f0-4b75-ae29-3854e26a0585
1d1e6eab-36b0-4fc1-a89e-5b34572d9f7a	70ebed7b-5653-42cf-be51-902c9dce2987	fbb11569-57f0-4b75-ae29-3854e26a0585
fbb8a493-024e-40b4-b0be-cf780f2a3c4a	ed06f369-9d92-4d07-8b9a-6488694b4d1c	fbb11569-57f0-4b75-ae29-3854e26a0585
373a1bad-8b8c-4180-8d1e-4a93d61b39d9	8bb9a2bc-aeea-435f-9132-73caead02416	fbb11569-57f0-4b75-ae29-3854e26a0585
ae6e22f0-3eaf-4b3b-9145-f778251ad6d9	d2e33ead-1f7a-4482-a8e0-d58deceea631	fbb11569-57f0-4b75-ae29-3854e26a0585
f0a1606c-5097-4fc2-bdb4-d589db05bff0	14aa553e-5be3-404d-adb0-b4440e023c69	fbb11569-57f0-4b75-ae29-3854e26a0585
5330295c-b53f-4fcd-8014-26b0dccc5db0	d6344f8d-4b29-48ce-b9e0-0c255447dc33	fbb11569-57f0-4b75-ae29-3854e26a0585
56cd1dd1-6d18-433b-a75f-5e1147adbf5a	f7833e34-f3da-40a2-9bf1-a14cb550925d	fbb11569-57f0-4b75-ae29-3854e26a0585
c34f8560-e91c-420a-aaee-0745e98a8287	f58d731d-1d2e-4380-ba29-4896951518e0	fbb11569-57f0-4b75-ae29-3854e26a0585
f4969cd8-1bff-48e6-ad11-b39d3cd51cd1	f01d9cd6-4d98-43b0-8aef-c52a46964512	fbb11569-57f0-4b75-ae29-3854e26a0585
e66ba618-eff4-4fd7-bf81-f5a61f839643	e41d0e51-6eb2-4622-8b59-f76ee6e6c1e1	fbb11569-57f0-4b75-ae29-3854e26a0585
16c858e4-33b9-4164-a37a-0fc3a91f46c0	ab1aa816-baed-4d33-bd9b-80d4b83a0341	fbb11569-57f0-4b75-ae29-3854e26a0585
71d1ca64-afe7-441e-ab05-d27dbdcd0796	13f5ca18-5e46-437e-a835-4a36fd6f9955	fbb11569-57f0-4b75-ae29-3854e26a0585
a8a3f18f-0d48-48c3-87ea-5a89c3aa0b6a	797cecc2-d1f9-499a-b79d-e1f7e77c7874	fbb11569-57f0-4b75-ae29-3854e26a0585
26c503dd-13f3-4d30-88cb-a245537c9a77	58405e28-6fff-42e1-8567-01466845abfb	fbb11569-57f0-4b75-ae29-3854e26a0585
88589921-f548-4eb4-95a9-a806ea74e349	68fc3f82-6076-412f-b986-20d3cb1bceaa	fbb11569-57f0-4b75-ae29-3854e26a0585
1a52e27f-7257-4779-a6a3-f407494e9eb7	03264ac4-7754-48ef-a824-5269d1f84e4b	fbb11569-57f0-4b75-ae29-3854e26a0585
c5eb8a32-a6aa-4dd4-b21b-62511e41a7d7	77f2a333-b579-421b-bf3f-5b913e6002b5	fbb11569-57f0-4b75-ae29-3854e26a0585
4e95126b-70fd-43f3-b56d-234b59e5756f	05a2e58e-0551-40c5-9bff-657768d7e689	fbb11569-57f0-4b75-ae29-3854e26a0585
5786dd1c-0fbb-4042-ac88-1772b388970d	fb9fbdf7-d05f-40ac-95e1-e10019ed8165	fbb11569-57f0-4b75-ae29-3854e26a0585
55350984-2c7e-43cf-8815-f8d29b04f4f2	39732090-0856-4b05-a241-df9886a313a5	fbb11569-57f0-4b75-ae29-3854e26a0585
72b931f7-2a17-457d-a681-f94c71392833	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	fbb11569-57f0-4b75-ae29-3854e26a0585
7c210e03-5779-44de-862e-56e67d92bce4	3aa46e51-bac0-4711-bf1b-bb396e71fe06	fbb11569-57f0-4b75-ae29-3854e26a0585
afc656f2-5488-4179-b624-f100ebf3487e	09cfad53-de06-4a25-88c0-9eaa26345df9	fbb11569-57f0-4b75-ae29-3854e26a0585
16e593c2-5399-4b95-a0bf-0da283e34797	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	fbb11569-57f0-4b75-ae29-3854e26a0585
f000845b-c7a8-4723-bcf4-5479634eba15	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	fbb11569-57f0-4b75-ae29-3854e26a0585
cbc74a1c-114d-415f-9b52-960698c19120	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	fbb11569-57f0-4b75-ae29-3854e26a0585
02b97521-3854-403d-af8f-9818017bb030	c45be260-e11b-4fae-8e33-05d0803f4e6b	fbb11569-57f0-4b75-ae29-3854e26a0585
1a0c55cd-9504-4b98-a901-72f11bd1a55d	24a99997-51d4-446c-bc51-2f46e3ca01a6	fbb11569-57f0-4b75-ae29-3854e26a0585
c6c165ae-ecfa-4d04-ba76-042aee2ae411	f3c2ef69-090e-4f61-a9d6-954330e86d8f	fbb11569-57f0-4b75-ae29-3854e26a0585
c8f6b5f2-9bb0-49ee-a46e-be08aa460e95	ebb13d64-4568-4cce-9f26-72f5e8a9bc14	fbb11569-57f0-4b75-ae29-3854e26a0585
bdb9b86d-13fd-4402-a39b-dc1cf5dde588	9f4da9ef-7f9c-407e-91a9-034f623b810f	fbb11569-57f0-4b75-ae29-3854e26a0585
fcb54c41-539c-4c53-9438-940098db28c6	653eb679-3e69-442e-b0ff-fd22e3fd8cd3	fbb11569-57f0-4b75-ae29-3854e26a0585
bf3cbd92-bfd5-4a19-9293-c9767890f09a	086e9acf-0405-449c-9dee-ee33d80df312	fbb11569-57f0-4b75-ae29-3854e26a0585
9ab92621-b296-451e-83f1-337d74ebc69b	93421699-5c88-4dee-995c-f74222b3a1a7	fbb11569-57f0-4b75-ae29-3854e26a0585
317c3b80-41db-4657-9186-22c50da593a8	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	fbb11569-57f0-4b75-ae29-3854e26a0585
579e0249-fb8d-413c-916b-a0544780454c	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	fbb11569-57f0-4b75-ae29-3854e26a0585
d4deac6a-109f-47ec-a756-7483f8d8f856	72237dba-0beb-450e-9df6-ea38d3255cab	fbb11569-57f0-4b75-ae29-3854e26a0585
cffd9e9b-6db4-4d77-9e51-b803c9653f74	67522e75-9672-4bbc-904d-57dd942290b7	fbb11569-57f0-4b75-ae29-3854e26a0585
989a6987-3ff1-4989-ae8d-05eb25ac5561	131cc1b8-2668-4392-9443-6c2e42a426ef	fbb11569-57f0-4b75-ae29-3854e26a0585
31e6d600-c115-48c1-8683-6f2a4de3a746	f53252df-70b6-4907-8847-908554862006	fbb11569-57f0-4b75-ae29-3854e26a0585
62ed70b5-b0f5-4982-ad2b-1332c6e19b80	3780a865-b70a-4a34-93f1-56c920e44100	fbb11569-57f0-4b75-ae29-3854e26a0585
fc1ffbd5-d905-48a9-b994-962989d91dee	0a9b8878-506b-4f76-af54-478e96f06eb9	fbb11569-57f0-4b75-ae29-3854e26a0585
868a1737-8c11-4e1f-a201-ec9759553e50	fc72276f-66bb-4609-b83a-4a84ff55100c	fbb11569-57f0-4b75-ae29-3854e26a0585
f921c638-4f79-4650-afdf-7297174faf73	2ad759be-3a54-4efa-aeda-7027dc97867f	94176a90-02bb-4dea-a306-ebdcf3c02975
439e6ff4-5203-40cf-b162-8d1155266b9f	e3eaddf3-edb3-4930-a548-bce70c834b66	29242052-f1fd-458e-9119-3fa6a9d5ce98
322149aa-de51-4066-ba72-da3566aa03ed	583c9de5-4079-4eb9-b448-4fa66ac594fd	29242052-f1fd-458e-9119-3fa6a9d5ce98
e9601f9f-f876-48e1-9811-92b879457f01	e23feaf4-58ae-469c-95a3-83e6312e7f29	29242052-f1fd-458e-9119-3fa6a9d5ce98
5c14f1f3-fbfd-4e48-a0a0-5dce4e251e0f	7ac101de-041e-4baa-ab7e-7471318ec81f	29242052-f1fd-458e-9119-3fa6a9d5ce98
b3cab700-a4a0-4c06-8df8-36a24443d231	0150f1aa-0dfe-4371-a553-87e73d886b3f	29242052-f1fd-458e-9119-3fa6a9d5ce98
edb1e5a0-9d91-4b57-b45f-8426a3d27925	665c7169-f989-436d-a611-e6effc8ddc9b	29242052-f1fd-458e-9119-3fa6a9d5ce98
04b75557-f8b1-44f0-bec8-f4e3b3bf80f9	f76dca7f-4774-4c76-8aac-873e97f66d3c	29242052-f1fd-458e-9119-3fa6a9d5ce98
db132063-8755-43f8-82ef-15b3c7687b22	00f047a0-8afe-40a6-89d5-74867de6cfbe	29242052-f1fd-458e-9119-3fa6a9d5ce98
05ca84d9-680b-4986-8f47-1d8bee5f88f0	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	29242052-f1fd-458e-9119-3fa6a9d5ce98
cc996e8b-c4a6-41d8-a8c7-0c416d1c3309	4a8285de-348e-4845-b0eb-142dd46002f6	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
307ebdc9-9205-4428-bedc-6627f40ba1ad	977f77aa-139c-4d4c-8640-053eea3376dd	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
59da08f4-2d56-4fcf-a8aa-3b2b62dd0f9e	90476856-94be-40ed-8b1e-814392074d81	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
2d6b316f-9151-4ac9-b4e7-352dd5645f4c	824f43c8-ee42-47ed-bbb9-5f8de09acdbc	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
c6454292-d902-411a-a397-5cc5f48e797e	0e143513-3b89-4e2d-9865-8a848ad6c8db	8457de16-e58d-4dc0-b20c-6a5730c4ed38
3ff1fb91-90ff-41b7-ad13-d44d91704188	ed06f369-9d92-4d07-8b9a-6488694b4d1c	8457de16-e58d-4dc0-b20c-6a5730c4ed38
788359ff-811c-425d-8465-32eb675db3a0	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	8457de16-e58d-4dc0-b20c-6a5730c4ed38
04adf4ab-7fd1-45eb-ba94-763b6c3995ed	c8f239e6-0791-4bdc-95f0-8e96abb74aa3	fbb11569-57f0-4b75-ae29-3854e26a0585
b924fcbc-081f-4a92-b443-ee03b74ecc85	213afcc4-81ca-4af6-aaad-3a83a55302a4	fbb11569-57f0-4b75-ae29-3854e26a0585
e03ef52d-c986-4784-a3e5-dcd6254d79ff	805c95bb-9f47-4e74-a1f0-87c3bb12a455	29242052-f1fd-458e-9119-3fa6a9d5ce98
ed77c639-d70f-4412-a58c-ad0a738ae57d	d68e979e-1fb7-400f-b2d3-4f8b31a08231	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
47817813-8797-4c2d-af2b-9e24892c12b4	e1ea71e4-fcbd-4ca5-8917-1dcab59fbd45	29242052-f1fd-458e-9119-3fa6a9d5ce98
cecb6104-8372-48c6-b14f-f90c2b010df6	3d74d6a0-376a-4ff2-a1fc-5150b624e769	29242052-f1fd-458e-9119-3fa6a9d5ce98
eaffc02e-5c04-4423-aeb2-8a0486c353e6	0f13980e-ffc3-4712-9c58-1726e1ec6787	29242052-f1fd-458e-9119-3fa6a9d5ce98
270369c9-c650-4ad7-8bb4-c9c1abcfa904	58d6c0a6-e80d-4d79-972f-0ad5f8517e97	29242052-f1fd-458e-9119-3fa6a9d5ce98
b243bcf2-45a4-4ffa-8eb7-0460c0cdd621	33bea49e-3565-469c-aca3-b524ed7232dd	29242052-f1fd-458e-9119-3fa6a9d5ce98
eb089325-5ad1-4af7-af67-d3c67266ddad	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	8457de16-e58d-4dc0-b20c-6a5730c4ed38
da7e0939-994a-4a1c-bfaf-c431ea8e9802	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	8457de16-e58d-4dc0-b20c-6a5730c4ed38
f6829854-fe78-4e58-9311-c13a4f3f807e	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	8457de16-e58d-4dc0-b20c-6a5730c4ed38
775dca5a-8ac0-4021-9f8a-1389238722cf	05a2e58e-0551-40c5-9bff-657768d7e689	8457de16-e58d-4dc0-b20c-6a5730c4ed38
2e513e21-cb68-4699-b966-a38ea78b16b7	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	8457de16-e58d-4dc0-b20c-6a5730c4ed38
53e23172-2d14-4b79-bc94-3f7ca28fa27e	03264ac4-7754-48ef-a824-5269d1f84e4b	8457de16-e58d-4dc0-b20c-6a5730c4ed38
4f1dc8c8-31ad-4e9d-8b4d-53ccc002ed61	51f985e0-4f2a-412a-adb5-3adf5399fc75	8457de16-e58d-4dc0-b20c-6a5730c4ed38
bacd0224-b3f1-45a3-afe5-d0306d145bc2	77f2a333-b579-421b-bf3f-5b913e6002b5	8457de16-e58d-4dc0-b20c-6a5730c4ed38
44cd284b-b0c0-4b3a-be04-60e515759d69	fb9fbdf7-d05f-40ac-95e1-e10019ed8165	8457de16-e58d-4dc0-b20c-6a5730c4ed38
dbedadd4-9807-4221-9e7c-401b2aa5442d	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	8457de16-e58d-4dc0-b20c-6a5730c4ed38
c3176950-8456-4578-9098-658c3d0bfa53	c070d24f-76eb-4684-8459-dbb7977acc46	8457de16-e58d-4dc0-b20c-6a5730c4ed38
efb96bf1-9398-411b-8978-5df0e7d51c04	b92b6639-c7fb-42ff-91bc-3281971224b2	984bb718-3808-4911-ac96-925444c6a716
a023b8ec-3902-412d-874e-b0f85fb6eea4	3ab532fe-274c-4755-b8c4-8055ea124209	984bb718-3808-4911-ac96-925444c6a716
1c87ac76-bd29-447a-a2e5-71633c625d33	ec64aed7-c889-4d44-a3de-38b112f490ec	984bb718-3808-4911-ac96-925444c6a716
cc160a0b-b5c7-4199-82d2-7423cd2d7cfc	2ed683f9-df88-46b5-bbb2-fd7f4d09535d	984bb718-3808-4911-ac96-925444c6a716
af23797f-4b4e-4725-b2f6-d8b7243bc117	a4beb80a-97bf-46c4-81ab-af3735050855	984bb718-3808-4911-ac96-925444c6a716
13f35d48-537f-43fa-8967-c7341b7f78d7	b56fa731-9b89-40dd-966d-b92122d6bdac	984bb718-3808-4911-ac96-925444c6a716
74a79e4c-d0d2-47ed-80ca-2e4ef3a0ccbd	01ba1a7b-1480-494c-8c0e-ca6b34d37d1e	984bb718-3808-4911-ac96-925444c6a716
b97c8d15-365f-43b7-9490-469fb0740ad9	8bb9a2bc-aeea-435f-9132-73caead02416	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fc2f5e58-ae15-4d7b-9553-c42146314123	c2592cf2-b4f8-4ab8-9ce3-54b4dbbc928a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
c320939b-d94b-4b5e-ad06-fd218760b670	a168235d-63c0-4b58-8b49-3df20a9a815d	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
80cf518f-c7af-4ca7-9db0-1b9755f282d4	14aa553e-5be3-404d-adb0-b4440e023c69	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d7e718ec-af52-4e4b-92e2-ad44f3d88272	6c8c8fcf-3ab7-4f08-8f41-5a53c3cc0ed1	984bb718-3808-4911-ac96-925444c6a716
1f163476-0c3a-439a-a31d-39588b2f3042	3fbb2d4e-c39d-41e3-aee8-079228318790	984bb718-3808-4911-ac96-925444c6a716
6cc20313-1680-4556-8177-906a157adca5	c5658139-2407-4416-847f-3a5eda3b9bb7	984bb718-3808-4911-ac96-925444c6a716
cc934061-d643-424c-89f6-52edbb6a8aa3	eb467d59-fea6-4d48-afde-3011582cef2f	984bb718-3808-4911-ac96-925444c6a716
50540cb0-2261-4dd7-af31-f1ae4703bf39	63593179-f9d6-4580-97ca-dd8670688272	984bb718-3808-4911-ac96-925444c6a716
cccbf318-3cb2-4f0c-b9de-0cb0a11896da	a8de69fd-5856-42cc-8812-fef457a99079	984bb718-3808-4911-ac96-925444c6a716
f3306459-54a4-41e1-86c7-601daa86df23	29058c58-14a6-425d-b366-617860f6950f	984bb718-3808-4911-ac96-925444c6a716
d633a4ad-e21d-459d-aa0a-6c43ddb30a12	4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56	984bb718-3808-4911-ac96-925444c6a716
a4358e13-4340-4c25-8a9e-2ffe4f5a16ad	97a477c0-7361-421f-b042-2d1e0d9bf5fe	984bb718-3808-4911-ac96-925444c6a716
5ed0c15d-32a1-4e4e-86fc-3e6797ffa5fd	a2bcc31f-515a-45e9-8083-b2f27d2c4afa	984bb718-3808-4911-ac96-925444c6a716
671dce2d-7fe5-46b4-aa80-a25a86565798	60f259e8-d103-4fda-859d-b6c053cc6004	984bb718-3808-4911-ac96-925444c6a716
50ebf047-cad3-42a2-9e7b-86f73a67f64b	46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d	984bb718-3808-4911-ac96-925444c6a716
d974457e-dd03-4e73-b519-475941651726	59d15cbe-0f30-4057-ab75-09ba8776b7d6	984bb718-3808-4911-ac96-925444c6a716
2689be8e-c6b5-42c7-99ae-fbadb8b70e50	88c58e99-19a7-4f72-811e-92e21346207a	984bb718-3808-4911-ac96-925444c6a716
a9c90bc8-082e-485f-9f09-4b980cba63f0	9b374efe-f7c1-44ff-a122-52981146b22a	984bb718-3808-4911-ac96-925444c6a716
e9336fbf-6103-4c1e-87aa-8d3e34357b34	eedb53fb-3a3e-4f32-a232-ffc97a4464ae	984bb718-3808-4911-ac96-925444c6a716
d6db22c9-b0cd-4918-afe6-173eb4b62d17	b36e11e0-99e1-4af5-b83b-6fe9062554a8	984bb718-3808-4911-ac96-925444c6a716
89e7edd3-1297-4ae1-b72d-efa0d09ad695	8ec08232-29d6-418f-8560-64f307684f8e	984bb718-3808-4911-ac96-925444c6a716
174e2b8f-d898-4298-ab87-711d8e195575	d88a3357-b68c-406e-8951-97df93df66bb	984bb718-3808-4911-ac96-925444c6a716
5412684e-3af0-48e3-b15a-d86d730b0227	cb332a03-aeee-435c-86de-b7d577e3817b	984bb718-3808-4911-ac96-925444c6a716
2f9b4189-cbc1-47d5-9f70-f2f619c85a42	1897bc12-5b09-412d-a3c7-e001aedff81d	984bb718-3808-4911-ac96-925444c6a716
8189f0ea-acba-4e60-8812-efc9858fa545	357d5bae-8470-4d03-ab92-68cafedcfa16	984bb718-3808-4911-ac96-925444c6a716
b4fc2920-3e1a-4715-af77-1c6f6e28bb61	4689bc27-a161-431b-bf58-456f66307fd9	984bb718-3808-4911-ac96-925444c6a716
d11f559a-7475-479f-86bd-faee5eb1ee0f	ef82cdff-1b9d-4396-8d2b-f7b751a50006	984bb718-3808-4911-ac96-925444c6a716
124a378e-4259-466e-a039-e3055d8fe681	065c9534-8c0a-408f-892a-914d172f90c7	984bb718-3808-4911-ac96-925444c6a716
2a865af5-0e22-4a1e-b1cf-da139fe4d267	ad55ef62-54a6-41c0-bba9-830a06822659	984bb718-3808-4911-ac96-925444c6a716
7ed524ee-2d7d-48e4-90cb-2906853c1f33	db98b40b-91a4-4434-ba2c-712e3990b82a	984bb718-3808-4911-ac96-925444c6a716
24ab2d82-8bba-4606-8f46-93690f6ebe91	b995956c-75be-4635-9022-af1f5376be2b	984bb718-3808-4911-ac96-925444c6a716
b334683e-72af-42b2-9634-5d64c361cbbb	03ed060b-56dd-4ecb-9072-a736c81ae27f	984bb718-3808-4911-ac96-925444c6a716
7b7978e2-d2dc-414a-bcee-0934414a56b6	e18c306e-1d61-4e6e-a2cf-c41ba357b406	984bb718-3808-4911-ac96-925444c6a716
37c04259-b8aa-4404-af11-795c2d205003	aa061751-3180-40dc-b101-f866d7e4dc2f	984bb718-3808-4911-ac96-925444c6a716
6889ab0b-e671-45a9-871e-8312e8e56827	b4aee72a-2ac5-4748-b846-bd1828139858	984bb718-3808-4911-ac96-925444c6a716
854ca7c5-33fa-4741-bee6-55ad71e9bbfa	72793652-a596-4b80-bbe4-73d10cca1379	984bb718-3808-4911-ac96-925444c6a716
e7aeb1ba-6d38-44f3-86f4-3c8d9beb40dd	5a8b13c6-3997-4846-b037-2b679f5c31e8	984bb718-3808-4911-ac96-925444c6a716
bc25f91a-950c-40b3-952d-c5eeb7ea8352	66cb5fb9-5bdc-4378-8010-6b58f8cedb5f	984bb718-3808-4911-ac96-925444c6a716
38b60ada-2dc8-4c44-a71d-f9452a61c549	3293247e-3d3c-4749-ae91-1a991e27b8be	984bb718-3808-4911-ac96-925444c6a716
266b96c0-0b0f-49a0-9b6a-bacdda4caf73	a0dee210-7e5b-4309-ad4e-7ccdf1e7e82e	984bb718-3808-4911-ac96-925444c6a716
580c0dfb-ed35-443d-b4ed-670b523bef8e	7ea8ce89-9ce7-4a89-a332-7baad7bbe475	984bb718-3808-4911-ac96-925444c6a716
578fb15d-38db-455b-9a59-209e0c6b4487	5a229d4e-beee-4f0a-ba2f-888f04a149e5	984bb718-3808-4911-ac96-925444c6a716
3a22cf74-924e-4c99-91c3-dd9eca80ba1a	742c55df-02ca-4d3e-909a-ac8c412ff735	984bb718-3808-4911-ac96-925444c6a716
25ba9b74-0ec1-4581-812c-95b1d3b060f4	c5e2efa4-c8c7-4c79-9823-59e829a4847e	984bb718-3808-4911-ac96-925444c6a716
3ed12625-48d4-4eab-82bb-493643d7c7cb	48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5	984bb718-3808-4911-ac96-925444c6a716
845b0b04-2a65-4703-878c-ee44e40f264f	90df4bc3-dbb0-4e92-ab9a-0aac895e5295	984bb718-3808-4911-ac96-925444c6a716
7804ebbe-1ccd-4589-b222-8cc8aacd3462	23476194-48d6-447c-b4e9-2ab6ad7b1a48	984bb718-3808-4911-ac96-925444c6a716
62ab192a-52d2-4f85-a51c-b20cc19e7724	4e471ceb-766f-4575-ab73-856e07ceb1a8	984bb718-3808-4911-ac96-925444c6a716
4de8c022-f0d5-419c-852e-2c7c73de357c	cab6e3e8-cff1-42da-959f-e12705adc760	984bb718-3808-4911-ac96-925444c6a716
95ec1387-6055-4050-8b33-15c8217df400	1c80833a-d9f9-4b31-ac77-55c282781167	984bb718-3808-4911-ac96-925444c6a716
94f07ef3-a1cf-4a05-915c-d84db31b55a4	01698813-a221-4b56-8bf5-ec3ec9956d97	984bb718-3808-4911-ac96-925444c6a716
9a57ccf5-b9e9-4eb4-be78-83758d58a406	0576302e-61ce-49b3-b265-28b7458e5c27	984bb718-3808-4911-ac96-925444c6a716
641f17ac-baca-4134-9224-5318f9e912ab	3c53a6aa-9d98-43be-bd14-4c5501e58ee5	984bb718-3808-4911-ac96-925444c6a716
e1191912-68f9-4264-95f3-ae1d28852447	3ed50a46-3620-43a1-b8ea-5a1ebc968851	984bb718-3808-4911-ac96-925444c6a716
51117da5-2642-498a-a159-42db614cd9ca	080f0752-21ba-4d3e-97c2-4846dfb4f648	984bb718-3808-4911-ac96-925444c6a716
087ff889-6458-4984-944a-09b261ac2975	82dbdde0-be98-462f-bef7-b86ed1a98336	984bb718-3808-4911-ac96-925444c6a716
718c392f-84da-46bb-a240-c85c75f1423f	e0be0b73-1f90-4c04-8da5-86513c148a16	984bb718-3808-4911-ac96-925444c6a716
102ef145-2ef2-426e-a829-afeb7748bfb7	c2ebd879-c09a-43d3-8088-6b0c2445c0b1	984bb718-3808-4911-ac96-925444c6a716
c30f5a16-d41d-4b6e-86d7-a79000b45788	af0187d0-e84c-401f-9ffc-70194b099309	984bb718-3808-4911-ac96-925444c6a716
a0fd6ae4-a31d-4cbb-85b0-c09b5efe2d8f	e296b896-76d1-4cf0-b092-e637a5909ec4	984bb718-3808-4911-ac96-925444c6a716
6220c89f-e6a4-48fb-85fa-a39a7014dd8c	77788f4d-5d0e-458b-b7cb-297e907b2eff	984bb718-3808-4911-ac96-925444c6a716
c251b437-d2b3-4bd1-88c7-74d531bc5537	136f2190-6f80-46d4-a124-d512e508be5e	984bb718-3808-4911-ac96-925444c6a716
dfb3f0e2-2e12-4b4c-9beb-26f15369084f	2ad759be-3a54-4efa-aeda-7027dc97867f	fbb11569-57f0-4b75-ae29-3854e26a0585
b382995a-f655-46cd-b275-c03a379783e3	f938a823-cd1f-4702-9163-45a5e016e2cf	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
2124cfe2-d425-4b72-a8f8-abf396e94025	96deec3e-55a6-422c-9cb5-3f59ff7a3dd2	29242052-f1fd-458e-9119-3fa6a9d5ce98
44bd79d2-429c-4124-ad9d-3b957b67eb43	deb75efd-1d47-4fb3-88d0-4fe07c5cbf36	29242052-f1fd-458e-9119-3fa6a9d5ce98
2ad3721b-718a-4bef-9e60-0824efd4f604	315a3bae-c0a7-47df-9522-acbc971ff8b1	29242052-f1fd-458e-9119-3fa6a9d5ce98
c7259702-f0e8-4484-9f3e-192a323f5d65	590dd762-650e-4b31-9989-268ae81a853c	29242052-f1fd-458e-9119-3fa6a9d5ce98
bb685b8e-58fd-4847-9ee0-e70cf81da9af	7e87afed-a5ab-4285-b257-b84161febe58	29242052-f1fd-458e-9119-3fa6a9d5ce98
2a7b9b6a-469b-4dbc-a6c1-19cd80b8e257	029a70de-2382-485b-9eb3-15d3933461d6	29242052-f1fd-458e-9119-3fa6a9d5ce98
351902ab-b939-446e-b579-3e72b01bf736	19b78054-0c63-453c-9ca7-bd64c99de999	29242052-f1fd-458e-9119-3fa6a9d5ce98
6af09386-c050-4c40-b297-ae5c13bfffc1	d1f06d6a-a377-4419-b690-1cbe09879b1c	29242052-f1fd-458e-9119-3fa6a9d5ce98
178a09f6-b048-485c-8d63-8ff21f9fb633	6f4cc747-5450-4514-9a3e-123ace9cb232	8457de16-e58d-4dc0-b20c-6a5730c4ed38
1bda19c0-cd65-42c5-b704-98a2ec9cd317	93421699-5c88-4dee-995c-f74222b3a1a7	8457de16-e58d-4dc0-b20c-6a5730c4ed38
23684451-00ee-4b0d-9e68-866e7c8a9b60	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	8457de16-e58d-4dc0-b20c-6a5730c4ed38
fdbeb063-949e-440c-91cf-dc1010d575ee	68fc3f82-6076-412f-b986-20d3cb1bceaa	8457de16-e58d-4dc0-b20c-6a5730c4ed38
863c896f-fa40-4658-a61a-0642cda260ae	1e5eb584-f19e-4110-89e7-2d63a08d727a	8457de16-e58d-4dc0-b20c-6a5730c4ed38
13ce8a60-bf45-4ef7-a14f-9680816692bb	58405e28-6fff-42e1-8567-01466845abfb	8457de16-e58d-4dc0-b20c-6a5730c4ed38
b9aa4c90-ad8b-4f5f-840f-7e2cf61d0e6e	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	8457de16-e58d-4dc0-b20c-6a5730c4ed38
c0dfdde8-4919-4b2e-b116-95e5db1ca1c5	ccad896d-13db-4fc8-be92-574af1d80341	8457de16-e58d-4dc0-b20c-6a5730c4ed38
3e3d51ce-ae0b-4bfa-86ef-2dec400b782e	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	8457de16-e58d-4dc0-b20c-6a5730c4ed38
96df0164-03da-43fd-98ae-f653db240395	d81d630f-3d80-4168-923b-65eefae629c9	8457de16-e58d-4dc0-b20c-6a5730c4ed38
db8c0a1e-d60b-4e8b-90b7-2ba38fb6a13f	ae18fd5a-e0ba-4ace-a637-902317a67773	8457de16-e58d-4dc0-b20c-6a5730c4ed38
4a56d176-d939-4f58-aa68-8eaa2b51dc94	b48dccdc-da4f-4914-bd2e-1971de87fda1	8457de16-e58d-4dc0-b20c-6a5730c4ed38
6ea88505-92d9-4d57-9682-cacff814f5e3	7767088d-9354-42d4-8cc7-9d3f2c0bef88	8457de16-e58d-4dc0-b20c-6a5730c4ed38
49e3c837-fc87-476b-b7d0-b4a400ad601d	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	8457de16-e58d-4dc0-b20c-6a5730c4ed38
255c1d54-d153-47d5-8068-0dee27c65d1f	27917132-65e1-4218-bdfc-40a62550a141	8457de16-e58d-4dc0-b20c-6a5730c4ed38
e29a8b76-9cca-47f7-8302-68abe9987a47	09cfad53-de06-4a25-88c0-9eaa26345df9	8457de16-e58d-4dc0-b20c-6a5730c4ed38
b0935cb8-8d04-4488-8722-1ac861184052	04032947-d9b3-4949-ad03-f6d9640a627a	29242052-f1fd-458e-9119-3fa6a9d5ce98
6c6e05f2-48c6-4534-ab52-b2c52d8414f5	5e6cb690-6a8d-4fab-9234-754db94e1c2f	29242052-f1fd-458e-9119-3fa6a9d5ce98
7e227ce9-dda0-41d1-84c7-f2dca37d70d0	8ceae246-d9fc-4fdf-8487-c4abc76c26af	29242052-f1fd-458e-9119-3fa6a9d5ce98
b3747771-1466-4487-8451-2c60e76c8cff	afe5da52-c330-4f13-b7a6-386cafd6fd0d	29242052-f1fd-458e-9119-3fa6a9d5ce98
7a5238e2-ba35-4154-bc95-eb6ecda8da1d	6330eaba-e5de-47f1-8c50-072120823683	29242052-f1fd-458e-9119-3fa6a9d5ce98
e6f99a46-5f5c-4545-9048-62af7d38df43	9ec25419-48ea-4525-8832-2f9cbcdc3096	29242052-f1fd-458e-9119-3fa6a9d5ce98
3de905a5-a7a9-40af-8dbf-ef5ef9864c9b	9b1d0f10-1be2-417c-a79f-a77e9dfa0c93	29242052-f1fd-458e-9119-3fa6a9d5ce98
f12e9341-1912-4a45-b534-2e47eeab3fe2	7628af1d-b868-4923-8bac-a7366efdc813	29242052-f1fd-458e-9119-3fa6a9d5ce98
48b13e71-b77d-44ca-ae1b-b0b2fa3e9445	4aacad3f-4671-4f8f-a6bc-73ce21233408	29242052-f1fd-458e-9119-3fa6a9d5ce98
d77022e8-6064-4a3c-bd26-d5a383422e37	d0f0d5f9-e375-4c11-aaac-22ae71a637ea	29242052-f1fd-458e-9119-3fa6a9d5ce98
e46f5cd9-3834-42e3-ac22-8391e49d61c0	b207a77e-a3b3-4303-b757-2bd70ca8639b	29242052-f1fd-458e-9119-3fa6a9d5ce98
0c1da78e-8d14-42dc-a864-0d2d56877e29	9aae5b79-2f48-461e-a6b2-45e8266402dd	29242052-f1fd-458e-9119-3fa6a9d5ce98
b5049fb9-8b32-4c33-883e-5870b448b1e5	adcbffee-495f-4e59-90c0-f438d39ebd07	29242052-f1fd-458e-9119-3fa6a9d5ce98
12ae9ef8-3cf4-4748-9f4e-aaeaa419f9dd	7e62c826-7285-4b26-9235-d4a49b40ef6c	29242052-f1fd-458e-9119-3fa6a9d5ce98
f20f595c-fe47-4b69-b8b6-64e8ca34f5bc	2ad759be-3a54-4efa-aeda-7027dc97867f	29242052-f1fd-458e-9119-3fa6a9d5ce98
29d89e34-d476-4fd8-bb72-5e66a115be84	797cecc2-d1f9-499a-b79d-e1f7e77c7874	29242052-f1fd-458e-9119-3fa6a9d5ce98
74d8415c-816e-4b46-89e8-da24e66a9fb3	15668ec7-6731-4646-a516-41035ddebe0d	29242052-f1fd-458e-9119-3fa6a9d5ce98
8136b8d4-422b-4e43-b5e9-1ce5ae0eb585	3aa46e51-bac0-4711-bf1b-bb396e71fe06	8457de16-e58d-4dc0-b20c-6a5730c4ed38
522c3254-9bd7-4584-8680-bb1ee146915e	f22b7a3d-2d76-483c-af40-dd0c04defb20	8457de16-e58d-4dc0-b20c-6a5730c4ed38
45c667a5-3f37-4dae-84b5-20ffc8965558	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	8457de16-e58d-4dc0-b20c-6a5730c4ed38
acbb1c3d-1dad-4d7a-85ee-db6116d97d63	39732090-0856-4b05-a241-df9886a313a5	8457de16-e58d-4dc0-b20c-6a5730c4ed38
430888c6-faf2-430c-8e44-5d6c73b339e5	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	8457de16-e58d-4dc0-b20c-6a5730c4ed38
c43d1e5d-26a1-451f-8eb7-cbcee7619552	c45be260-e11b-4fae-8e33-05d0803f4e6b	8457de16-e58d-4dc0-b20c-6a5730c4ed38
7e58b48e-099f-4756-b517-8a29f145d452	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	8457de16-e58d-4dc0-b20c-6a5730c4ed38
50a00b66-9a42-4d10-8379-ff1355c133ed	32a0d6bb-1c39-4701-a0fb-ecd39f0f76d4	8457de16-e58d-4dc0-b20c-6a5730c4ed38
00f18f00-9488-4008-a691-889f7d86d138	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	8457de16-e58d-4dc0-b20c-6a5730c4ed38
bfddbe67-0e04-48ca-ac56-381efda4ef22	afe7beea-4242-4b06-9d4f-00dd11aae178	8457de16-e58d-4dc0-b20c-6a5730c4ed38
d52c2733-3779-4af9-8be2-0f02a3aeb63a	26b070f5-65f4-44b7-ad6d-2f8630937f3c	8457de16-e58d-4dc0-b20c-6a5730c4ed38
d4e970b3-2cda-4cc7-8e2f-3de8803e2d9a	3c15f3e4-bcab-4157-99d3-6332f2544089	8457de16-e58d-4dc0-b20c-6a5730c4ed38
ff233602-1b0e-4393-b96e-d08acc1dee84	24127ef9-9ef3-46f3-88cd-97378e694a6c	8457de16-e58d-4dc0-b20c-6a5730c4ed38
e9c0fd2d-4734-4aba-bed4-e6d9173da133	24a99997-51d4-446c-bc51-2f46e3ca01a6	8457de16-e58d-4dc0-b20c-6a5730c4ed38
6835dbf7-e98a-4fe6-a8a1-5f65ee1c7505	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	8457de16-e58d-4dc0-b20c-6a5730c4ed38
b7d97d3f-72b6-4dce-bae7-0ebca2bff246	f3c2ef69-090e-4f61-a9d6-954330e86d8f	8457de16-e58d-4dc0-b20c-6a5730c4ed38
a29d576e-4224-413c-8331-021b0f949e0b	c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	8457de16-e58d-4dc0-b20c-6a5730c4ed38
c2d32453-9162-40a4-b103-9c9b16722248	5aae3b9c-87f3-487f-b016-f7109f17d2c0	8457de16-e58d-4dc0-b20c-6a5730c4ed38
ae728bf6-d1dd-48ac-9246-7b0a3596ffec	288731cb-561c-4d4a-8951-d8c359a51846	8457de16-e58d-4dc0-b20c-6a5730c4ed38
7f10e5a2-d1c9-4176-b346-e07a3bf8cefb	ebb13d64-4568-4cce-9f26-72f5e8a9bc14	8457de16-e58d-4dc0-b20c-6a5730c4ed38
aeaddcff-2f93-40df-b668-235704d79530	9f4da9ef-7f9c-407e-91a9-034f623b810f	8457de16-e58d-4dc0-b20c-6a5730c4ed38
f2e5de68-a6ad-444a-b1da-6c8b14cd3fb0	72237dba-0beb-450e-9df6-ea38d3255cab	8457de16-e58d-4dc0-b20c-6a5730c4ed38
237f743d-854c-4baa-ac39-9a9aa72e5d56	8c9f5d83-fb7b-4f21-8b60-1ed3eb93ff41	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
33c537df-78e2-4e45-8ec0-b719d335d709	e9cb3b58-b3ed-4e27-a2f5-0c704dbbf32f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
3a6e16ac-0e89-40d3-8d6f-82bbae1a8de8	0e143513-3b89-4e2d-9865-8a848ad6c8db	85a78368-21eb-45db-9feb-07d054e5b445
b5c470dd-84e7-45ef-943c-3c8221079c38	ed06f369-9d92-4d07-8b9a-6488694b4d1c	85a78368-21eb-45db-9feb-07d054e5b445
af82bdd9-f2e6-4ba6-8c24-b98220fe2c2d	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	85a78368-21eb-45db-9feb-07d054e5b445
56dd54e8-d8f3-4eab-982e-084fbd370e34	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	85a78368-21eb-45db-9feb-07d054e5b445
d73ea6fe-5634-4f65-a345-0bf61f30f974	05a2e58e-0551-40c5-9bff-657768d7e689	85a78368-21eb-45db-9feb-07d054e5b445
d8fc9316-81cc-4a39-a7c9-dc8ead91fb06	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	85a78368-21eb-45db-9feb-07d054e5b445
42f87e37-5a0c-442f-a593-f11e2d02ac13	39732090-0856-4b05-a241-df9886a313a5	85a78368-21eb-45db-9feb-07d054e5b445
348d15e0-2038-4415-ac17-368da3c1a3e1	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	85a78368-21eb-45db-9feb-07d054e5b445
a42c9f13-24cd-4dac-a858-63a9a07a2be9	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	85a78368-21eb-45db-9feb-07d054e5b445
7ad45fef-4550-4308-9ee3-e6d181255b5a	b48dccdc-da4f-4914-bd2e-1971de87fda1	85a78368-21eb-45db-9feb-07d054e5b445
67e33689-6468-4d5d-9601-b814564abbc2	3aa46e51-bac0-4711-bf1b-bb396e71fe06	85a78368-21eb-45db-9feb-07d054e5b445
3a81e324-f08e-43a9-96ba-ee250497efdb	27917132-65e1-4218-bdfc-40a62550a141	85a78368-21eb-45db-9feb-07d054e5b445
6c7ae9d6-623e-455e-904c-5dd640f62ff8	09cfad53-de06-4a25-88c0-9eaa26345df9	85a78368-21eb-45db-9feb-07d054e5b445
2f2ffead-c0f9-4d1e-a74c-83f62e2353ff	7767088d-9354-42d4-8cc7-9d3f2c0bef88	85a78368-21eb-45db-9feb-07d054e5b445
7f30e9a6-8c22-4233-877d-2caa24629cc9	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	85a78368-21eb-45db-9feb-07d054e5b445
9cc60793-2a42-4468-84f6-db96a072a552	f22b7a3d-2d76-483c-af40-dd0c04defb20	85a78368-21eb-45db-9feb-07d054e5b445
671a84e1-3141-411f-88c1-15ccc064ccc1	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	85a78368-21eb-45db-9feb-07d054e5b445
09031dac-5466-4f68-8a0c-232d94a128de	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	85a78368-21eb-45db-9feb-07d054e5b445
82f9fe6b-d44d-4587-956d-c98fae034110	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	85a78368-21eb-45db-9feb-07d054e5b445
256dbd4d-422a-4cc7-9906-089018953dfb	f3c2ef69-090e-4f61-a9d6-954330e86d8f	85a78368-21eb-45db-9feb-07d054e5b445
ee04605b-e531-4a58-8c6b-87a2f8544a83	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	85a78368-21eb-45db-9feb-07d054e5b445
92533dc6-f132-4978-8503-565cc4d9e494	086e9acf-0405-449c-9dee-ee33d80df312	85a78368-21eb-45db-9feb-07d054e5b445
90d45f4d-9bdb-408d-a6a3-24fec8e6ba01	d81d630f-3d80-4168-923b-65eefae629c9	85a78368-21eb-45db-9feb-07d054e5b445
5ab6fd47-460d-4f8b-baab-28137a6d538b	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	85a78368-21eb-45db-9feb-07d054e5b445
a6c188ec-ae2f-4684-bead-93be443b1440	ae18fd5a-e0ba-4ace-a637-902317a67773	85a78368-21eb-45db-9feb-07d054e5b445
c38ad245-bed4-4c13-95c6-cbe471655967	72237dba-0beb-450e-9df6-ea38d3255cab	85a78368-21eb-45db-9feb-07d054e5b445
db69f53a-18f0-4a2e-b7ea-17b1e3f61155	5a31e833-8621-4b3b-841a-a0b03a8575e1	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
8793b084-1236-4000-bff1-8e90d042844f	2fdb5752-a2e2-4705-8c9b-b4a618e6de1f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
204b4872-32bf-402b-84cd-4bb62cd4347e	653eb679-3e69-442e-b0ff-fd22e3fd8cd3	8457de16-e58d-4dc0-b20c-6a5730c4ed38
a6179cda-7820-4bb0-8320-8c7574d378e2	b439edb9-605f-47dc-aee2-8df45976736e	8457de16-e58d-4dc0-b20c-6a5730c4ed38
a5d86e05-eef7-49ff-9128-07dd85a5318c	5493c560-ace3-440c-a36a-9538f618007a	8457de16-e58d-4dc0-b20c-6a5730c4ed38
0d3e28ee-9d4c-4027-8bff-a4e635d39526	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	8457de16-e58d-4dc0-b20c-6a5730c4ed38
6a9d96f9-5df0-44ef-b1fb-35b096c8c823	086e9acf-0405-449c-9dee-ee33d80df312	8457de16-e58d-4dc0-b20c-6a5730c4ed38
b488634e-36e1-4acc-a3bf-a41bb084ca71	3ab532fe-274c-4755-b8c4-8055ea124209	8457de16-e58d-4dc0-b20c-6a5730c4ed38
0616cafe-f330-492e-902c-4389c6f05e38	a8de69fd-5856-42cc-8812-fef457a99079	8457de16-e58d-4dc0-b20c-6a5730c4ed38
7b3ab9b3-a7b1-4d69-b995-21a5a1310730	136f2190-6f80-46d4-a124-d512e508be5e	8457de16-e58d-4dc0-b20c-6a5730c4ed38
ed2e37c1-bd25-41ec-8e05-09d1f24a368a	d5c1e2ba-4969-45d4-8b30-8da6cc872ae6	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
e8bc9449-9bd8-4d75-9616-8ec1b482f382	d2e33ead-1f7a-4482-a8e0-d58deceea631	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
5b8b7fff-6b93-4dc1-95f1-9d4c01961b70	d5e71527-bc03-4f52-b4e6-8460cbd3924b	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
fc926cc0-338e-489c-be1d-c890103fec1a	ccad896d-13db-4fc8-be92-574af1d80341	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
41b3f6a0-6ebc-46cd-bddf-e07fa50e2e09	cc9a6c3a-518a-458a-a84e-6353a1c22165	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
ea9ad7e1-81c2-4f34-a15c-e8a7a545d58b	761f6268-1ccb-4a23-b43d-cad43d4761d2	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
4a5f5a34-aafe-4f02-9fd3-82762c466fa5	8b470f68-065a-4b9a-a35e-0781d83910df	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fa3efa0d-2d06-49d0-918e-cba025e720f7	18460b2d-6e29-43da-9419-990fd196de7c	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
db61adac-e578-44e6-afd1-97f482ee2577	46ea3bb2-76b6-4842-a015-33c9edc910ae	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
e93cc1d4-0d32-418a-91c8-0683d1f24737	f42406c3-08c8-443a-845a-996e04c55753	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
73ee25d0-5e6d-41e1-b6bd-56eaca9cf0ba	0e143513-3b89-4e2d-9865-8a848ad6c8db	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
192b7634-63e4-479c-98c7-3cc3f6b67ffb	ed06f369-9d92-4d07-8b9a-6488694b4d1c	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
1e6e98c3-f468-48a4-ac3b-944049b49b1f	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
81b88f39-28fb-4141-94cb-d3529d126e85	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
bc69f0ef-b24a-4ed7-bb98-4035347b9234	05a2e58e-0551-40c5-9bff-657768d7e689	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
e7c8ce69-6c15-45eb-bfd7-281a5f6a324a	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
5486d322-383c-46f9-88df-ac936003d5a4	39732090-0856-4b05-a241-df9886a313a5	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
380f0f36-7d6a-4a77-97f9-e10f1e286cc0	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
d82682f1-bdf5-4d15-b43b-0faa71a31e51	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
c285421e-fdf9-4421-8197-64876da8af19	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
16191312-0b90-4dbe-a6d9-ba1b00eb4853	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
cc225ce1-50df-4a18-85b6-7ca867b29633	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
bfb1c3bc-bd05-44b2-90eb-2b88a787d72e	f3c2ef69-090e-4f61-a9d6-954330e86d8f	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
2fc97ed0-0558-4a66-bd69-74f2f9043846	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
54bc24e4-44bd-45ef-bd77-d27cb08e3ffb	086e9acf-0405-449c-9dee-ee33d80df312	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
c4a91adc-6e4f-49c1-9654-0b85661cc465	d81d630f-3d80-4168-923b-65eefae629c9	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
289c5c7c-1b63-414e-995d-4d48b8e8983e	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
5e7aeb3d-4c84-4d3b-b1ba-2027458d3148	ae18fd5a-e0ba-4ace-a637-902317a67773	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
b905fe60-db77-49ee-88fc-3b8b1ef3d737	72237dba-0beb-450e-9df6-ea38d3255cab	0ad57fbe-a74e-40ad-ab48-afa88c1dceea
89b56ef5-953c-42be-83eb-d951b9d32788	5493c560-ace3-440c-a36a-9538f618007a	fbb11569-57f0-4b75-ae29-3854e26a0585
9f77ad7b-f1bf-4b23-bf0a-8f4ae577cef6	2ad759be-3a54-4efa-aeda-7027dc97867f	8457de16-e58d-4dc0-b20c-6a5730c4ed38
6242d0bb-6890-4e41-978b-4ff061377363	797cecc2-d1f9-499a-b79d-e1f7e77c7874	8457de16-e58d-4dc0-b20c-6a5730c4ed38
349eb6b0-0876-40d0-a5ae-f6c80e5fdfe7	15668ec7-6731-4646-a516-41035ddebe0d	8457de16-e58d-4dc0-b20c-6a5730c4ed38
3bff688c-046e-4649-a8ca-d8928d14f21d	d5e71527-bc03-4f52-b4e6-8460cbd3924b	8457de16-e58d-4dc0-b20c-6a5730c4ed38
734bfdde-6518-4c1a-afb9-a6df4abf2018	0e143513-3b89-4e2d-9865-8a848ad6c8db	e920b75e-3c45-496b-9df0-f092455381eb
1c54dc6e-953b-4732-b96c-1a279b8f76a2	ed06f369-9d92-4d07-8b9a-6488694b4d1c	e920b75e-3c45-496b-9df0-f092455381eb
3d0859b5-1752-4336-8d59-3d1845307c61	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	e920b75e-3c45-496b-9df0-f092455381eb
d4664730-5104-4c45-8197-bfbb96842273	39732090-0856-4b05-a241-df9886a313a5	e920b75e-3c45-496b-9df0-f092455381eb
c1b90f54-1c92-41b1-9f2c-aab716e9d992	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	e920b75e-3c45-496b-9df0-f092455381eb
5f5989ee-7f55-4187-99fe-96e28e641acf	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	e920b75e-3c45-496b-9df0-f092455381eb
18d86ece-c075-44e5-a943-e8f9e9f238aa	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	e920b75e-3c45-496b-9df0-f092455381eb
4bba1902-52fb-44ee-a747-e21b7b68b085	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	e920b75e-3c45-496b-9df0-f092455381eb
3f51b9ba-c2da-4c28-9c52-66222991745c	d81d630f-3d80-4168-923b-65eefae629c9	e920b75e-3c45-496b-9df0-f092455381eb
591a8266-4bc7-48dc-a05c-0f98b7e83a3e	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	e920b75e-3c45-496b-9df0-f092455381eb
3e040158-23bb-40e9-a547-693751b04613	ae18fd5a-e0ba-4ace-a637-902317a67773	e920b75e-3c45-496b-9df0-f092455381eb
1e630e76-e9ea-40f3-bb40-5f02d7734041	72237dba-0beb-450e-9df6-ea38d3255cab	e920b75e-3c45-496b-9df0-f092455381eb
2c66906f-944f-4c9d-97ff-0dcd166b516b	58405e28-6fff-42e1-8567-01466845abfb	1ec06744-dd1e-4664-8666-61a1f8974b16
73bfb747-9c78-4442-9217-7b80a724edcd	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	1ec06744-dd1e-4664-8666-61a1f8974b16
2b33d5be-ac15-4fb3-922a-d261d1b26ce9	68fc3f82-6076-412f-b986-20d3cb1bceaa	1ec06744-dd1e-4664-8666-61a1f8974b16
6eaec60b-9149-4dde-86c9-a5d7f608e70a	1e5eb584-f19e-4110-89e7-2d63a08d727a	1ec06744-dd1e-4664-8666-61a1f8974b16
2ddba69a-fd8d-46ab-a010-1bd2a4c362da	d5e71527-bc03-4f52-b4e6-8460cbd3924b	1ec06744-dd1e-4664-8666-61a1f8974b16
79425901-04e4-4d56-a5e9-556ba2772397	ccad896d-13db-4fc8-be92-574af1d80341	1ec06744-dd1e-4664-8666-61a1f8974b16
acc32ad1-9058-41b5-b8ba-3ee8029708f2	39732090-0856-4b05-a241-df9886a313a5	1ec06744-dd1e-4664-8666-61a1f8974b16
89ce9f1d-e542-4d8d-9b39-01a0b211df89	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	1ec06744-dd1e-4664-8666-61a1f8974b16
e787db6a-76a1-4724-86c9-a355de9b9e79	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	1ec06744-dd1e-4664-8666-61a1f8974b16
ede2c70a-a91c-4fe3-b76c-5b0c1729211c	6f4cc747-5450-4514-9a3e-123ace9cb232	1ec06744-dd1e-4664-8666-61a1f8974b16
93e010aa-dafd-41b9-86f0-02adca943a58	93421699-5c88-4dee-995c-f74222b3a1a7	1ec06744-dd1e-4664-8666-61a1f8974b16
9008b0a8-c2c7-4c4f-a7d7-dea1a9a35227	c070d24f-76eb-4684-8459-dbb7977acc46	1ec06744-dd1e-4664-8666-61a1f8974b16
504f766e-bfbf-42c1-88d4-00e6edbbd2b7	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	1ec06744-dd1e-4664-8666-61a1f8974b16
f3036205-edc8-4905-ac99-f5648ef74102	72237dba-0beb-450e-9df6-ea38d3255cab	1ec06744-dd1e-4664-8666-61a1f8974b16
5bfeb9d6-40fc-4360-ac96-78fd96df92ae	57873320-1f3c-4d56-ad41-122f17f362a2	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
7467187f-9d32-48c1-a070-b8cd3581e586	f50b81e9-1797-40e7-a91b-5ce8be353e8b	94176a90-02bb-4dea-a306-ebdcf3c02975
a559e6f6-8901-4383-938f-3d0ab1308bfd	ae527d80-3b93-4ec1-a4b3-5c21ab2c80fc	94176a90-02bb-4dea-a306-ebdcf3c02975
f36b5169-0b05-453c-9689-94ab2dac4434	fa083891-8c7b-4d3e-a769-9c765760f79e	94176a90-02bb-4dea-a306-ebdcf3c02975
c5da6a27-207d-4718-b7da-4efdc9d113f5	6b7a79ac-30dd-44ce-a6bb-82a137f67aaa	94176a90-02bb-4dea-a306-ebdcf3c02975
7491e86b-7db6-40c9-a771-b2c7403b035f	5e6cb690-6a8d-4fab-9234-754db94e1caa	94176a90-02bb-4dea-a306-ebdcf3c02975
414558e6-90e6-45c2-9a0c-51212944d946	230733a9-34d8-47b9-aa97-83a85dbca11e	94176a90-02bb-4dea-a306-ebdcf3c02975
2bbb3531-05de-48ee-aa67-ae2dd1dbb03a	5b63f7e6-6bf1-406d-8895-b6f92418520c	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
cc866594-4cd4-40af-8df7-63e28fc6bfdd	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	94176a90-02bb-4dea-a306-ebdcf3c02975
c958f2e5-e518-4f5d-bfcb-9a00063a8f96	805c95bb-9f47-4e74-a1f0-87c3bb12a455	94176a90-02bb-4dea-a306-ebdcf3c02975
b993df60-f7d1-4acc-9b5c-a96ccc7eec87	58d6c0a6-e80d-4d79-972f-0ad5f8517e97	94176a90-02bb-4dea-a306-ebdcf3c02975
2e3271b8-e9de-4baa-a6d4-2bed7d12d21d	3d74d6a0-376a-4ff2-a1fc-5150b624e769	94176a90-02bb-4dea-a306-ebdcf3c02975
b341c0a8-1106-4ad1-9c22-54829830a24b	0f13980e-ffc3-4712-9c58-1726e1ec6787	94176a90-02bb-4dea-a306-ebdcf3c02975
3bf09db8-041a-4711-b6ea-1acaf3927717	e1ea71e4-fcbd-4ca5-8917-1dcab59fbd45	94176a90-02bb-4dea-a306-ebdcf3c02975
7b8b0f39-2409-43a9-a99d-9014148740a7	33bea49e-3565-469c-aca3-b524ed7232dd	94176a90-02bb-4dea-a306-ebdcf3c02975
82fd5a17-88fc-4add-aa15-c0657ced3e28	5e5c4e3f-155c-48f4-9aa6-dc7dd88fa56f	94176a90-02bb-4dea-a306-ebdcf3c02975
eb86ee3d-ec73-4633-8aed-e38c1e41d7b6	64cbf7b5-71e0-4179-a0a8-07020335901d	94176a90-02bb-4dea-a306-ebdcf3c02975
2193dee0-1b80-4e25-bdc4-0b05057034f6	23ee37f8-d777-4fdc-9d2b-8705c72ce357	94176a90-02bb-4dea-a306-ebdcf3c02975
307d2d8e-b7f3-4eb7-965d-3cae158bf452	6eee51ef-8a8f-4966-b19d-c9b915e681e4	94176a90-02bb-4dea-a306-ebdcf3c02975
acdf736e-8e1c-40ff-850d-bfcdbb0fc474	5219e94c-69eb-4baa-86ac-593d0d0255fc	94176a90-02bb-4dea-a306-ebdcf3c02975
965260b8-b050-42ab-8579-fcb1a58604f5	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	94176a90-02bb-4dea-a306-ebdcf3c02975
cf6f6564-b84e-4272-a0c6-4d129acb64e8	04032947-d9b3-4949-ad03-f6d9640a627a	94176a90-02bb-4dea-a306-ebdcf3c02975
473bb3ef-63f1-40e8-ad04-b2dfa35ffa28	5e6cb690-6a8d-4fab-9234-754db94e1c2f	94176a90-02bb-4dea-a306-ebdcf3c02975
37ce9780-3a15-4613-be10-6cdf9325fce4	8ceae246-d9fc-4fdf-8487-c4abc76c26af	94176a90-02bb-4dea-a306-ebdcf3c02975
ef8880e2-83ff-4629-b114-aa8ad4fa2a8f	afe5da52-c330-4f13-b7a6-386cafd6fd0d	94176a90-02bb-4dea-a306-ebdcf3c02975
c5c98482-52e7-4f76-80bb-ac42bc2422f7	6330eaba-e5de-47f1-8c50-072120823683	94176a90-02bb-4dea-a306-ebdcf3c02975
3e6b8387-5225-4efb-8e99-0f7b5065c500	9ec25419-48ea-4525-8832-2f9cbcdc3096	94176a90-02bb-4dea-a306-ebdcf3c02975
02b66e7a-edae-4c40-b750-226fc3e48f61	9b1d0f10-1be2-417c-a79f-a77e9dfa0c93	94176a90-02bb-4dea-a306-ebdcf3c02975
5ad8e4b4-86d9-429c-a8ce-d831110edcfd	5a08d43a-1268-4b70-8514-60fccd917026	94176a90-02bb-4dea-a306-ebdcf3c02975
9ee2da43-7d61-4692-b5bd-8452fcce6aff	608b52b4-40a6-43f1-bca2-031b2ea2abf4	94176a90-02bb-4dea-a306-ebdcf3c02975
950901f1-bfa3-42e3-89bd-6eea1821abc0	a761fc46-09d5-43aa-91ad-93391ae4934a	94176a90-02bb-4dea-a306-ebdcf3c02975
7010e22b-22c2-4ebe-aea0-14e43f04b6ed	67591c4a-7794-421e-8c3d-cc8b94b694b5	94176a90-02bb-4dea-a306-ebdcf3c02975
a0aed434-1f1c-4b70-960f-e2847b9a44b1	7628af1d-b868-4923-8bac-a7366efdc813	94176a90-02bb-4dea-a306-ebdcf3c02975
957bf2e8-eb56-4b4f-a5f6-61e04ace7276	a0419fd4-acd0-4e2b-8156-30066469f33e	94176a90-02bb-4dea-a306-ebdcf3c02975
f1ea960c-aa47-4fe7-8dd2-f7e95a3580bf	7872d920-bf48-491f-a604-f29239826d86	94176a90-02bb-4dea-a306-ebdcf3c02975
e2693408-b25b-408f-8fc8-2516e3c2d75b	6b7a79ac-30dd-44ce-a6bb-82a137f6771b	94176a90-02bb-4dea-a306-ebdcf3c02975
a5f49214-4b67-4a9a-8e9d-41945c3e2eaf	fa083891-8c7b-4d3e-a769-9c765710f79a	94176a90-02bb-4dea-a306-ebdcf3c02975
65ad58f8-aa32-4e2a-a881-5056b90a1c0f	e27c19dc-5b9c-45ef-8af3-0c024d1abf34	94176a90-02bb-4dea-a306-ebdcf3c02975
75a7d306-ed14-4e8d-bd85-ae4bc001571f	973f37a3-9aa6-48a0-8220-66b1444e81ab	94176a90-02bb-4dea-a306-ebdcf3c02975
b6de6f5d-da47-4bbe-86f9-010f24b22d96	98a5950b-849c-4fef-8735-979aa36ed84c	94176a90-02bb-4dea-a306-ebdcf3c02975
683e446a-eb4d-4111-87cf-2975bad4657b	43b61baf-0753-4367-9d29-7be6623d91c9	94176a90-02bb-4dea-a306-ebdcf3c02975
5e54b377-b714-4d40-8bae-809786eca049	6f13190f-abaa-48e1-a490-8d4406df0f34	94176a90-02bb-4dea-a306-ebdcf3c02975
77e6dd4c-7399-4205-97d3-e64829b8df86	4aacad3f-4671-4f8f-a6bc-73ce21233408	94176a90-02bb-4dea-a306-ebdcf3c02975
e6c3fc96-7942-4bef-9b80-874e29678fd3	3f3abbb0-4c6b-4b05-9e0f-4413b4adf347	94176a90-02bb-4dea-a306-ebdcf3c02975
e856323f-8560-4ff8-b542-fe463bb9773f	d0f0d5f9-e375-4c11-aaac-22ae71a637ea	94176a90-02bb-4dea-a306-ebdcf3c02975
189bf31d-e562-4a93-b749-0442112fd491	189208bd-a65f-44d8-8a07-96647feaafb7	94176a90-02bb-4dea-a306-ebdcf3c02975
4443c3ff-703e-40d2-a0d3-6c4c93085f8b	b207a77e-a3b3-4303-b757-2bd70ca8639b	94176a90-02bb-4dea-a306-ebdcf3c02975
9404dcc8-1376-44ff-835d-8dd6bbc8133f	9aae5b79-2f48-461e-a6b2-45e8266402dd	94176a90-02bb-4dea-a306-ebdcf3c02975
635d1635-688e-4390-854c-219de9e91191	adcbffee-495f-4e59-90c0-f438d39ebd07	94176a90-02bb-4dea-a306-ebdcf3c02975
e0f803e1-8a21-4b2c-a7d7-5317eb69c0ba	b8840d1d-3908-4db7-9825-b5ce78232409	94176a90-02bb-4dea-a306-ebdcf3c02975
6d37f8a9-49fe-461b-a78a-610685d94805	6afbed32-4c43-43c7-b300-c1a04857d346	94176a90-02bb-4dea-a306-ebdcf3c02975
4844d918-5dd7-435e-a658-c9d701161b1e	6be56b39-079c-410a-b41a-35b32a6a5e1f	94176a90-02bb-4dea-a306-ebdcf3c02975
8370af26-562c-4c20-981f-9fe774d6cb09	4f387a47-6b08-48f6-981c-42a8a128d1f0	94176a90-02bb-4dea-a306-ebdcf3c02975
7e297841-6a34-4754-a8bb-4b9be859b3b9	28c94a48-fd22-49a2-aa57-0f72c1f635d2	94176a90-02bb-4dea-a306-ebdcf3c02975
542588cc-6c09-4dbd-9c08-ba5561d74fcf	a62b7765-f024-4705-877a-de0ad4878d19	94176a90-02bb-4dea-a306-ebdcf3c02975
f55fc8b2-70b8-4f8f-8f8d-62ad6535ec4a	05825db2-caf8-4951-9f2d-ea85ace66073	94176a90-02bb-4dea-a306-ebdcf3c02975
5ce35cbe-9bfd-41f0-b24d-de4669ffc5db	44a5389b-8608-4865-ab53-74be807043a2	94176a90-02bb-4dea-a306-ebdcf3c02975
66a62f43-17d1-4cc1-98a8-800846618101	e72dd18e-5845-4ef4-934c-91b3f01eddde	94176a90-02bb-4dea-a306-ebdcf3c02975
42d94ca4-1411-4cdb-a4ec-037185aa90f6	6c154b5d-4689-4282-a54e-dbbeeeba5fb5	94176a90-02bb-4dea-a306-ebdcf3c02975
5666b0f3-f3ea-4e0a-8668-dcc25c22810c	8ec3bc88-f2ce-498f-8d01-a40b02804f4c	94176a90-02bb-4dea-a306-ebdcf3c02975
2530f462-81ae-4e7e-b95e-af3157b81335	a1951306-b89e-46d3-8768-3d681585e382	94176a90-02bb-4dea-a306-ebdcf3c02975
253f6c98-51dc-43fb-b226-4dfef22725f1	192ddf18-05b1-4d80-8b2c-923cef89df16	94176a90-02bb-4dea-a306-ebdcf3c02975
4296332f-b7f5-4d1a-b048-2715770dcf71	55327a19-fd5c-4101-885c-ec17508e5e74	94176a90-02bb-4dea-a306-ebdcf3c02975
78b6016f-5009-44cb-a47f-a030902df542	590dd762-650e-4b31-9989-268ae81a853c	94176a90-02bb-4dea-a306-ebdcf3c02975
22cdec23-f287-4b4e-9c74-8f38045e2a64	deb75efd-1d47-4fb3-88d0-4fe07c5cbf36	94176a90-02bb-4dea-a306-ebdcf3c02975
ff455e28-1412-49aa-820d-4ee5b038a2ab	19b78054-0c63-453c-9ca7-bd64c99de999	94176a90-02bb-4dea-a306-ebdcf3c02975
403b1805-8030-4645-9855-d6498aaa43c2	96deec3e-55a6-422c-9cb5-3f59ff7a3dd2	94176a90-02bb-4dea-a306-ebdcf3c02975
afc50ffb-7564-45ec-ace8-d1c0db7fc98f	029a70de-2382-485b-9eb3-15d3933461d6	94176a90-02bb-4dea-a306-ebdcf3c02975
9e4f697b-7e26-4493-911c-43e2dfb66a07	7e87afed-a5ab-4285-b257-b84161febe58	94176a90-02bb-4dea-a306-ebdcf3c02975
8edf38e6-83a8-48a1-adab-22af50821ca6	315a3bae-c0a7-47df-9522-acbc971ff8b1	94176a90-02bb-4dea-a306-ebdcf3c02975
d29ca30d-c033-413f-8f8a-62dfcfd9129b	d1f06d6a-a377-4419-b690-1cbe09879b1c	94176a90-02bb-4dea-a306-ebdcf3c02975
e3c1d9d2-9839-4e11-b86f-3eacbc7453e5	1e8a19be-9e51-48af-afc0-d56772e66a24	94176a90-02bb-4dea-a306-ebdcf3c02975
fd67d988-caf4-4ad6-9353-8c421eaeb627	ff6e92a2-a7a4-4508-9d4e-bde9f6cb0bfe	94176a90-02bb-4dea-a306-ebdcf3c02975
d2d3bad5-fa05-47c5-96dd-4f5b95c085f5	ad45fb2b-38a8-4d14-8eb5-03af12a2f7ec	94176a90-02bb-4dea-a306-ebdcf3c02975
f06db0ad-6618-4a80-88b7-0530ef258114	5f65f5d0-39d7-42e6-a7f5-dd21ba72abb1	94176a90-02bb-4dea-a306-ebdcf3c02975
3123f015-9dbf-4d59-b98c-54db74a20eec	267939c7-d951-47f8-aad8-63f0c0d6cdfb	94176a90-02bb-4dea-a306-ebdcf3c02975
22d9f8dd-86f8-4253-a3cd-136c86c58650	8feea562-c412-4546-916a-7c22e60c833c	94176a90-02bb-4dea-a306-ebdcf3c02975
a419ff27-0a88-4ffe-bcc0-13a6932defbc	a57004e3-4bac-4cc0-b2c6-8b0a060e2dbb	94176a90-02bb-4dea-a306-ebdcf3c02975
30bf5267-22fa-4e73-8126-5cd5deb101e5	a35f9884-abfa-4c87-9c23-a983e04dfaf7	94176a90-02bb-4dea-a306-ebdcf3c02975
685332c0-41cc-403a-b360-85f37a8f8928	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	94176a90-02bb-4dea-a306-ebdcf3c02975
9771eff4-23cf-415f-861f-da61fd3a2a0e	15362ef8-4d86-4e96-b49f-c33c5292e67f	94176a90-02bb-4dea-a306-ebdcf3c02975
626e1d40-378f-4e82-bbf2-c6894e7b648d	b39c59d0-f9d1-4901-af1c-f3a193d1c612	94176a90-02bb-4dea-a306-ebdcf3c02975
70ff3da8-0b28-4022-98c7-e9cdea9eefbf	b639edba-fe15-48c6-9e53-551c056a4426	94176a90-02bb-4dea-a306-ebdcf3c02975
1d7b3263-3dc3-4028-99a7-021cbb3bdf04	f76dca7f-4774-4c76-8aac-873e97f66d3c	94176a90-02bb-4dea-a306-ebdcf3c02975
b0fa89a7-a784-4b01-8702-6ea747a4de86	e3eaddf3-edb3-4930-a548-bce70c834b66	94176a90-02bb-4dea-a306-ebdcf3c02975
e5569026-5e41-4ce2-8481-af2185be8574	0150f1aa-0dfe-4371-a553-87e73d886b3f	94176a90-02bb-4dea-a306-ebdcf3c02975
11b953fd-986c-4e91-9876-00a04fe251b8	7ac101de-041e-4baa-ab7e-7471318ec81f	94176a90-02bb-4dea-a306-ebdcf3c02975
096e12cc-210b-4d35-ab22-7a10ee7ca6b5	e23feaf4-58ae-469c-95a3-83e6312e7f29	94176a90-02bb-4dea-a306-ebdcf3c02975
b50e5c33-712d-47f3-b552-acdf8240f3a2	00f047a0-8afe-40a6-89d5-74867de6cfbe	94176a90-02bb-4dea-a306-ebdcf3c02975
4febd4aa-c7a5-44bb-98ba-aec85050354b	665c7169-f989-436d-a611-e6effc8ddc9b	94176a90-02bb-4dea-a306-ebdcf3c02975
8c953ed6-a7fc-431d-abe9-82792d9b66ab	583c9de5-4079-4eb9-b448-4fa66ac594fd	94176a90-02bb-4dea-a306-ebdcf3c02975
9d524911-6ab0-4dfb-95de-507a5d434f2f	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	94176a90-02bb-4dea-a306-ebdcf3c02975
20deb501-4e81-4894-87ea-0d965b1027ef	ab03d85a-45db-4af0-8b37-0de7e8bfe4a6	94176a90-02bb-4dea-a306-ebdcf3c02975
cf838284-000d-4259-9ebf-3593fac905b4	fa944b01-d087-4729-920d-f1f7ab36a642	94176a90-02bb-4dea-a306-ebdcf3c02975
f3a782bc-863b-481f-aacc-d62b64fea003	2feb45eb-2e97-4d7f-bf43-d221b1c71928	94176a90-02bb-4dea-a306-ebdcf3c02975
4550b9bb-179a-4890-9820-1c1a5968d449	7e62c826-7285-4b26-9235-d4a49b40ef6c	94176a90-02bb-4dea-a306-ebdcf3c02975
e6410fa2-41c0-4aa3-8fdc-2fcc73951df3	b92b6639-c7fb-42ff-91bc-3281971224b2	94176a90-02bb-4dea-a306-ebdcf3c02975
dec3a68b-e819-451b-8847-81d7e6d5c334	3ab532fe-274c-4755-b8c4-8055ea124209	94176a90-02bb-4dea-a306-ebdcf3c02975
f636fbe7-5f4d-48d8-abef-1f94ee43e999	ec64aed7-c889-4d44-a3de-38b112f490ec	94176a90-02bb-4dea-a306-ebdcf3c02975
90554841-face-4212-ada5-94a7b0d564ce	2ed683f9-df88-46b5-bbb2-fd7f4d09535d	94176a90-02bb-4dea-a306-ebdcf3c02975
470897a0-ae44-4c40-8857-a1a6f8e2d44c	a4beb80a-97bf-46c4-81ab-af3735050855	94176a90-02bb-4dea-a306-ebdcf3c02975
f51386ca-a8eb-4ec8-b6e3-96edecd150c6	b56fa731-9b89-40dd-966d-b92122d6bdac	94176a90-02bb-4dea-a306-ebdcf3c02975
3d459b4f-c507-409b-a0a3-7791ff3778a4	01ba1a7b-1480-494c-8c0e-ca6b34d37d1e	94176a90-02bb-4dea-a306-ebdcf3c02975
be07ee84-17a2-4172-9c43-da72ca540bea	6c8c8fcf-3ab7-4f08-8f41-5a53c3cc0ed1	94176a90-02bb-4dea-a306-ebdcf3c02975
be8289ad-fdc2-4280-a73e-3c8c36519c4c	3fbb2d4e-c39d-41e3-aee8-079228318790	94176a90-02bb-4dea-a306-ebdcf3c02975
bbb37579-c128-49df-971a-1d7954f1ed1d	c5658139-2407-4416-847f-3a5eda3b9bb7	94176a90-02bb-4dea-a306-ebdcf3c02975
b4bd7a6d-034b-4496-beab-70696c41efce	eb467d59-fea6-4d48-afde-3011582cef2f	94176a90-02bb-4dea-a306-ebdcf3c02975
2cfbc12f-0b51-43e3-a238-1b0afe1c4d13	63593179-f9d6-4580-97ca-dd8670688272	94176a90-02bb-4dea-a306-ebdcf3c02975
e9e42f9e-bf3d-44d7-8fa0-bb3bd35307fd	a8de69fd-5856-42cc-8812-fef457a99079	94176a90-02bb-4dea-a306-ebdcf3c02975
8060a9cc-c472-4e2a-b4e9-26b3bbc07365	29058c58-14a6-425d-b366-617860f6950f	94176a90-02bb-4dea-a306-ebdcf3c02975
515d3d22-a51c-438e-a088-dbdd7c71b8c3	4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56	94176a90-02bb-4dea-a306-ebdcf3c02975
07cbe06e-9934-4dc5-8ba9-5448d482d4b0	97a477c0-7361-421f-b042-2d1e0d9bf5fe	94176a90-02bb-4dea-a306-ebdcf3c02975
1538ec9d-a834-476c-b966-002a11371702	a2bcc31f-515a-45e9-8083-b2f27d2c4afa	94176a90-02bb-4dea-a306-ebdcf3c02975
cda6d9f3-173c-42f1-9267-22c47e72bd26	60f259e8-d103-4fda-859d-b6c053cc6004	94176a90-02bb-4dea-a306-ebdcf3c02975
66d3a2cf-6a53-450b-8a64-ea4356e13902	46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d	94176a90-02bb-4dea-a306-ebdcf3c02975
d9821b80-a377-48df-9463-92a894176b8a	59d15cbe-0f30-4057-ab75-09ba8776b7d6	94176a90-02bb-4dea-a306-ebdcf3c02975
6226d044-5d71-4659-845c-4e88134352cc	88c58e99-19a7-4f72-811e-92e21346207a	94176a90-02bb-4dea-a306-ebdcf3c02975
3caa691b-4b03-473b-89cb-2d08fc67c3b0	9b374efe-f7c1-44ff-a122-52981146b22a	94176a90-02bb-4dea-a306-ebdcf3c02975
9668cd16-ce1a-442e-a83b-46e66276dd2a	eedb53fb-3a3e-4f32-a232-ffc97a4464ae	94176a90-02bb-4dea-a306-ebdcf3c02975
090fc6e9-0440-43eb-9eed-1ef2f3b0694d	b36e11e0-99e1-4af5-b83b-6fe9062554a8	94176a90-02bb-4dea-a306-ebdcf3c02975
5e7002ad-b9ad-41e6-9dae-cefc831a32c0	8ec08232-29d6-418f-8560-64f307684f8e	94176a90-02bb-4dea-a306-ebdcf3c02975
59d5b74a-14c0-472c-a33a-6f6130c087ad	d88a3357-b68c-406e-8951-97df93df66bb	94176a90-02bb-4dea-a306-ebdcf3c02975
da35043f-a545-4e24-a455-74cf7f2cfe61	cb332a03-aeee-435c-86de-b7d577e3817b	94176a90-02bb-4dea-a306-ebdcf3c02975
e81a58ab-2d42-4c31-b978-b7d90606975f	742c55df-02ca-4d3e-909a-ac8c412ff735	94176a90-02bb-4dea-a306-ebdcf3c02975
db67f016-9a48-4aa0-af48-12787ebcace7	c5e2efa4-c8c7-4c79-9823-59e829a4847e	94176a90-02bb-4dea-a306-ebdcf3c02975
682727dd-3ce1-4d8e-8785-7ec8a061334c	48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5	94176a90-02bb-4dea-a306-ebdcf3c02975
96a576df-f9f1-4e02-b2f7-9b554c1cbe94	90df4bc3-dbb0-4e92-ab9a-0aac895e5295	94176a90-02bb-4dea-a306-ebdcf3c02975
3d14cc5d-2624-498c-8665-ad052e849d72	23476194-48d6-447c-b4e9-2ab6ad7b1a48	94176a90-02bb-4dea-a306-ebdcf3c02975
74a8c7f8-9cb6-488f-9b24-3d44b3fbfc4c	4e471ceb-766f-4575-ab73-856e07ceb1a8	94176a90-02bb-4dea-a306-ebdcf3c02975
8f02f906-8af1-4fc1-8855-4ce792cf2db6	cab6e3e8-cff1-42da-959f-e12705adc760	94176a90-02bb-4dea-a306-ebdcf3c02975
16b7e308-b821-4f6d-a9d2-b21cce3ded52	1c80833a-d9f9-4b31-ac77-55c282781167	94176a90-02bb-4dea-a306-ebdcf3c02975
83553909-2e99-4f7d-a4ce-5e0952ba313a	01698813-a221-4b56-8bf5-ec3ec9956d97	94176a90-02bb-4dea-a306-ebdcf3c02975
1a16db59-f98f-4dca-8baf-04758a18157c	0576302e-61ce-49b3-b265-28b7458e5c27	94176a90-02bb-4dea-a306-ebdcf3c02975
451d363b-4fe1-4586-b1c7-3f1b799e8434	3c53a6aa-9d98-43be-bd14-4c5501e58ee5	94176a90-02bb-4dea-a306-ebdcf3c02975
895eee68-bd4d-4304-ab9d-36d267c401f7	3ed50a46-3620-43a1-b8ea-5a1ebc968851	94176a90-02bb-4dea-a306-ebdcf3c02975
dd0cffed-dadf-4033-aa47-d9ac81f0e4bf	080f0752-21ba-4d3e-97c2-4846dfb4f648	94176a90-02bb-4dea-a306-ebdcf3c02975
a862c8f1-6f62-4256-a7c1-5fe1e0bd1cbc	82dbdde0-be98-462f-bef7-b86ed1a98336	94176a90-02bb-4dea-a306-ebdcf3c02975
c36b35ed-19ec-4cc1-b9c7-b8b8ae5d26c2	e0be0b73-1f90-4c04-8da5-86513c148a16	94176a90-02bb-4dea-a306-ebdcf3c02975
ca0391d7-d23b-43ba-94c0-ec61f51f1617	c2ebd879-c09a-43d3-8088-6b0c2445c0b1	94176a90-02bb-4dea-a306-ebdcf3c02975
ca891f5c-1d0f-4f31-a837-7e71581d2164	af0187d0-e84c-401f-9ffc-70194b099309	94176a90-02bb-4dea-a306-ebdcf3c02975
6b4a2756-e974-49a9-834b-90269b5e58c5	e296b896-76d1-4cf0-b092-e637a5909ec4	94176a90-02bb-4dea-a306-ebdcf3c02975
160d4b3f-7da5-4050-bc9d-2c7f1c1d9492	77788f4d-5d0e-458b-b7cb-297e907b2eff	94176a90-02bb-4dea-a306-ebdcf3c02975
daaf56d0-10a9-420d-891b-d3c68bbe1f6a	136f2190-6f80-46d4-a124-d512e508be5e	94176a90-02bb-4dea-a306-ebdcf3c02975
fafb0b37-18d1-4088-b8d9-29fefefe3ff3	b2776c17-fde6-4af5-8d34-f72d5ec88ac9	94176a90-02bb-4dea-a306-ebdcf3c02975
ffe42c4b-50ff-4226-ad03-c5a51065b456	aded78fc-0d49-436e-a780-ebea7c8b0f4c	94176a90-02bb-4dea-a306-ebdcf3c02975
d89f8ecd-d7d1-4e85-8ba8-cbed86be2475	622186a1-17b8-4b64-8406-dd3cd1dd2596	94176a90-02bb-4dea-a306-ebdcf3c02975
a7058f88-a6b1-4c0b-a4b6-d2898573ebfa	222d3e77-c07d-4c30-a384-b0c493a3cb38	94176a90-02bb-4dea-a306-ebdcf3c02975
3523e2f3-8230-450f-9101-08c35b850185	0e143513-3b89-4e2d-9865-8a848ad6c8db	94176a90-02bb-4dea-a306-ebdcf3c02975
ddef4e59-b271-4338-8576-653ffa900791	e9cb3b58-b3ed-4e27-a2f5-0c704dbbf32f	94176a90-02bb-4dea-a306-ebdcf3c02975
4561f2bb-cb2d-4a97-ba30-81d3f9a90d4a	ed06f369-9d92-4d07-8b9a-6488694b4d1c	94176a90-02bb-4dea-a306-ebdcf3c02975
2ad233d1-93ed-448a-bd7e-c6db3950752a	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	94176a90-02bb-4dea-a306-ebdcf3c02975
68cdb536-1f01-4bf0-8315-4e8b1d335cc3	d6344f8d-4b29-48ce-b9e0-0c255447dc33	94176a90-02bb-4dea-a306-ebdcf3c02975
32ea23c3-bdea-49f2-8b14-5c2966a4b49f	6e4e2624-bc76-49dc-86ba-6e88b0977c79	94176a90-02bb-4dea-a306-ebdcf3c02975
5d188d8a-32f6-4155-9d42-6d36a1ed2d52	b8e3f794-9775-4363-9bcf-d3d929f6ba27	94176a90-02bb-4dea-a306-ebdcf3c02975
26774e15-cf34-4a8b-ab5a-723bd6259e1b	13f5ca18-5e46-437e-a835-4a36fd6f9955	94176a90-02bb-4dea-a306-ebdcf3c02975
a98ab6cd-96a6-4625-a2e7-0febe9169a05	ee775cf8-f617-4324-9220-726c7ddbe919	94176a90-02bb-4dea-a306-ebdcf3c02975
f523765e-5acc-40a4-99eb-b7ec47f7d29e	797cecc2-d1f9-499a-b79d-e1f7e77c7874	94176a90-02bb-4dea-a306-ebdcf3c02975
bd65bb33-4dd8-4a87-813f-f4c85eec0dce	15668ec7-6731-4646-a516-41035ddebe0d	94176a90-02bb-4dea-a306-ebdcf3c02975
3fb437c4-ae59-496d-915e-61bb98d24fb8	58405e28-6fff-42e1-8567-01466845abfb	94176a90-02bb-4dea-a306-ebdcf3c02975
09ac44b1-22fb-49c0-b6ba-52ca98bf01d9	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	94176a90-02bb-4dea-a306-ebdcf3c02975
b07fe993-9148-4eb0-a583-5ffb75de3420	68fc3f82-6076-412f-b986-20d3cb1bceaa	94176a90-02bb-4dea-a306-ebdcf3c02975
d3fe04ad-220f-4e2a-bb0a-e4850a7cc007	ccad896d-13db-4fc8-be92-574af1d80341	94176a90-02bb-4dea-a306-ebdcf3c02975
7818950b-9042-4dfd-bce0-519046d1463f	03264ac4-7754-48ef-a824-5269d1f84e4b	94176a90-02bb-4dea-a306-ebdcf3c02975
b279383e-5292-4b40-841c-d253a98a59a6	51f985e0-4f2a-412a-adb5-3adf5399fc75	94176a90-02bb-4dea-a306-ebdcf3c02975
0c3acfe6-73db-44a8-b8b5-8ed4b8420706	77f2a333-b579-421b-bf3f-5b913e6002b5	94176a90-02bb-4dea-a306-ebdcf3c02975
7acd4a43-fcc3-4bd4-b345-a1e4bde5cae9	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	94176a90-02bb-4dea-a306-ebdcf3c02975
6410b710-5360-4ddc-9ff8-354aaa4d3b57	05a2e58e-0551-40c5-9bff-657768d7e689	94176a90-02bb-4dea-a306-ebdcf3c02975
0a2aa786-9a84-4125-be2e-c932c670a380	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	94176a90-02bb-4dea-a306-ebdcf3c02975
7122880a-1194-43fb-bb09-0c34151552cd	fb9fbdf7-d05f-40ac-95e1-e10019ed8165	94176a90-02bb-4dea-a306-ebdcf3c02975
9c1e04f2-996b-43b3-86f5-1fa846cfbb81	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	94176a90-02bb-4dea-a306-ebdcf3c02975
3c2fbeb4-0f2f-4029-b3a4-d83299c4159e	39732090-0856-4b05-a241-df9886a313a5	94176a90-02bb-4dea-a306-ebdcf3c02975
75dff12a-dd8d-4c34-b531-007fa88a5482	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	94176a90-02bb-4dea-a306-ebdcf3c02975
fb4dd0e6-738a-4ff8-bf49-c3f7f5f587a2	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	94176a90-02bb-4dea-a306-ebdcf3c02975
bd5ec579-8083-4c08-a66d-ac45990d1212	b48dccdc-da4f-4914-bd2e-1971de87fda1	94176a90-02bb-4dea-a306-ebdcf3c02975
8bc06889-a9f8-490e-91db-1624b5e0a625	3aa46e51-bac0-4711-bf1b-bb396e71fe06	94176a90-02bb-4dea-a306-ebdcf3c02975
891fee28-5727-4595-bfc0-165fba27281f	27917132-65e1-4218-bdfc-40a62550a141	94176a90-02bb-4dea-a306-ebdcf3c02975
13630c32-8546-4196-a0b6-d0b0bfc2f392	09cfad53-de06-4a25-88c0-9eaa26345df9	94176a90-02bb-4dea-a306-ebdcf3c02975
128db896-db51-41c9-acc5-4af9b5ded3b0	7767088d-9354-42d4-8cc7-9d3f2c0bef88	94176a90-02bb-4dea-a306-ebdcf3c02975
d54ef449-7ea2-4519-9a8d-921c505cc5c7	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	94176a90-02bb-4dea-a306-ebdcf3c02975
1bd3ac12-646d-4f51-8ab4-7e8d9a025e91	f22b7a3d-2d76-483c-af40-dd0c04defb20	94176a90-02bb-4dea-a306-ebdcf3c02975
c21c29ab-1086-413f-9f98-6a7a5cd7ccf0	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	94176a90-02bb-4dea-a306-ebdcf3c02975
f3a94622-c03b-45e9-87c1-8c2cbcf4cb5f	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	94176a90-02bb-4dea-a306-ebdcf3c02975
5fadae98-463a-469a-a5cc-6d3b7639ef69	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	94176a90-02bb-4dea-a306-ebdcf3c02975
7ebed39d-ad5f-46e3-938e-61b529c8a4be	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	94176a90-02bb-4dea-a306-ebdcf3c02975
35b5fee0-3b11-44ac-92c8-1280942c1c00	c45be260-e11b-4fae-8e33-05d0803f4e6b	94176a90-02bb-4dea-a306-ebdcf3c02975
bd4f259c-d466-4a8c-838a-06c914f39884	24127ef9-9ef3-46f3-88cd-97378e694a6c	94176a90-02bb-4dea-a306-ebdcf3c02975
fe115a70-68af-4119-8532-8dec7453c45f	24a99997-51d4-446c-bc51-2f46e3ca01a6	94176a90-02bb-4dea-a306-ebdcf3c02975
aa139cf9-8ea1-442e-8b3d-7b739ce21a6d	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	94176a90-02bb-4dea-a306-ebdcf3c02975
80d39675-1280-4972-a225-af1e01fb6e39	f3c2ef69-090e-4f61-a9d6-954330e86d8f	94176a90-02bb-4dea-a306-ebdcf3c02975
b5cbd74b-04ac-4cca-ab53-72cb4c2f4432	c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	94176a90-02bb-4dea-a306-ebdcf3c02975
b56f5fcc-6eb3-4060-9f00-1ae452a2663a	ebb13d64-4568-4cce-9f26-72f5e8a9bc14	94176a90-02bb-4dea-a306-ebdcf3c02975
0aad2054-6c81-4c57-b279-01ec2074d95c	288731cb-561c-4d4a-8951-d8c359a51846	94176a90-02bb-4dea-a306-ebdcf3c02975
4459583d-ec9a-4aa3-99a3-11a542dfe87a	9f4da9ef-7f9c-407e-91a9-034f623b810f	94176a90-02bb-4dea-a306-ebdcf3c02975
913324f6-04a0-4656-871a-97cfaef172a3	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	94176a90-02bb-4dea-a306-ebdcf3c02975
fe275bd0-61d6-4dd2-9ef8-04fd00d860a1	653eb679-3e69-442e-b0ff-fd22e3fd8cd3	94176a90-02bb-4dea-a306-ebdcf3c02975
82e1d234-a353-4364-9d28-971526180f76	b439edb9-605f-47dc-aee2-8df45976736e	94176a90-02bb-4dea-a306-ebdcf3c02975
ee52f4de-b001-469e-9fd5-f4fabeac44fb	a0e0bacb-c1d2-488f-aa5f-1af550bd97a8	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
8c2a1744-364d-43c0-8215-08f10e70df5a	d141790b-1036-447d-b4cf-ae8b75de2190	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
61fee277-4632-4114-85c5-40f2725c1838	ed0c5e11-eddb-48ff-84f9-6730d389f01a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
12d4228b-fb2e-488c-85b2-b390e30fe5b5	3ab532fe-274c-4755-b8c4-8055ea124209	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d92ad458-e2a2-44af-99b3-aac45c1e28be	6f855990-13d0-43b1-87e6-f0647e99c360	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
4b2630a5-c674-488a-afac-018bc741f8e6	086e9acf-0405-449c-9dee-ee33d80df312	94176a90-02bb-4dea-a306-ebdcf3c02975
e72ebe64-3cb6-4355-9529-4d734baea9c7	6f4cc747-5450-4514-9a3e-123ace9cb232	94176a90-02bb-4dea-a306-ebdcf3c02975
ace23ab1-2e1f-401b-a7ab-900db63200c0	93421699-5c88-4dee-995c-f74222b3a1a7	94176a90-02bb-4dea-a306-ebdcf3c02975
f148ba12-5159-4f81-b9f8-fa5ffcd250c8	c070d24f-76eb-4684-8459-dbb7977acc46	94176a90-02bb-4dea-a306-ebdcf3c02975
4b4b9c38-902f-4014-89c5-2344ff6088a8	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	94176a90-02bb-4dea-a306-ebdcf3c02975
d4ac6dcd-48e1-4e76-9652-50b2446a67d1	63593179-f9d6-4580-97ca-dd8670688272	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0a61f088-59b7-47cc-a08b-4411f4bf6dd2	a8de69fd-5856-42cc-8812-fef457a99079	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
427b279a-0409-40fd-8a31-c874e2c1c16c	8ec08232-29d6-418f-8560-64f307684f8e	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
1475c945-4303-4e07-84cc-4ecd65d0708e	d81d630f-3d80-4168-923b-65eefae629c9	94176a90-02bb-4dea-a306-ebdcf3c02975
005aba2c-5964-4c36-9e74-5600aabd7f9c	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	94176a90-02bb-4dea-a306-ebdcf3c02975
2834b721-77df-419e-b413-2dc1c8c7a723	ae18fd5a-e0ba-4ace-a637-902317a67773	94176a90-02bb-4dea-a306-ebdcf3c02975
06e97ce0-6470-4b6b-9aff-4746d9585deb	72237dba-0beb-450e-9df6-ea38d3255cab	94176a90-02bb-4dea-a306-ebdcf3c02975
3b12c881-bc0e-4c4e-be6e-dd8519ad4042	223c3afe-ef65-4b86-9733-6284396cd65f	94176a90-02bb-4dea-a306-ebdcf3c02975
0e6304fc-3ddd-4050-9bfd-57a3061e410d	67522e75-9672-4bbc-904d-57dd942290b7	94176a90-02bb-4dea-a306-ebdcf3c02975
fa9d3c01-807b-40a8-8364-e4acddaa0987	78cd88ae-d982-48c6-998c-1a6ba961b5f0	94176a90-02bb-4dea-a306-ebdcf3c02975
d5a7191a-e9d2-4577-9a84-548fdc96450a	131cc1b8-2668-4392-9443-6c2e42a426ef	94176a90-02bb-4dea-a306-ebdcf3c02975
ce3ce47f-cb9d-4f41-8a0e-623414568f48	20805d38-e63d-4117-b902-7007c1697a13	94176a90-02bb-4dea-a306-ebdcf3c02975
2aea014a-c6db-49a4-872a-5da6eca876f2	f53252df-70b6-4907-8847-908554862006	94176a90-02bb-4dea-a306-ebdcf3c02975
05ee43ef-b06a-4aff-bc31-77e8adecaeaa	99a0afca-2fef-4083-82a2-4785ac933ec0	94176a90-02bb-4dea-a306-ebdcf3c02975
d7b85b53-f4c0-43b2-bfa9-8c4be592a0d8	3780a865-b70a-4a34-93f1-56c920e44100	94176a90-02bb-4dea-a306-ebdcf3c02975
e39aae82-7063-4689-b37b-6f6cb5e3e984	a06b41ea-3abb-4dc4-a4f7-90bfc86a9592	94176a90-02bb-4dea-a306-ebdcf3c02975
09d08a7d-f63b-4fd9-97a9-2bc450f09f0e	0a9b8878-506b-4f76-af54-478e96f06eb9	94176a90-02bb-4dea-a306-ebdcf3c02975
44c4ae11-b3e4-4021-abe7-fa9ef8ea805c	25c7b3c0-efd1-4d00-a748-e37de9ca4e3c	94176a90-02bb-4dea-a306-ebdcf3c02975
ca523d1b-25c2-4d6f-ab7a-8b65deaaa5c7	fc72276f-66bb-4609-b83a-4a84ff55100c	94176a90-02bb-4dea-a306-ebdcf3c02975
672aaadf-1cb0-4c41-b63c-4fb3322f4783	fcea9f1c-b882-4368-a7a3-c8bbd2491684	94176a90-02bb-4dea-a306-ebdcf3c02975
3cf47fd2-1176-486a-8401-f8b29831ff27	0e143513-3b89-4e2d-9865-8a848ad6c8db	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
9acb7cfd-0f4a-4fed-9ac7-6fcdf8a8195a	ed06f369-9d92-4d07-8b9a-6488694b4d1c	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
ea560991-fc75-471f-b63c-1c4a87ad5a6b	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
7df15b8b-6522-44f2-8ffc-b1dc67ea9370	ee775cf8-f617-4324-9220-726c7ddbe919	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
f9633b22-81fd-4509-a32f-8bcfa0635d99	2ad759be-3a54-4efa-aeda-7027dc97867f	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
974e1ca6-68bf-480e-8299-2c953917a839	797cecc2-d1f9-499a-b79d-e1f7e77c7874	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
5d493fee-8004-4b38-ba70-5f177be4ff0e	15668ec7-6731-4646-a516-41035ddebe0d	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
7d53de13-d4e6-4d5f-9149-d8c7b87a7a81	58405e28-6fff-42e1-8567-01466845abfb	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
53c50566-8b13-4613-8749-5458b22ec1bb	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
9fede614-59af-47d9-8ce4-e63b865bba06	68fc3f82-6076-412f-b986-20d3cb1bceaa	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
59037a07-bedd-44ac-a915-cd8b3f8855bc	ccad896d-13db-4fc8-be92-574af1d80341	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
c70ff6fc-7815-4a5c-978b-d908b484131d	03264ac4-7754-48ef-a824-5269d1f84e4b	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
b9fe5b18-e7c8-4dfb-8863-90bb102eb0b6	b36e11e0-99e1-4af5-b83b-6fe9062554a8	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
e1549ef7-75e3-4c25-8460-69a540ece133	a2bcc31f-515a-45e9-8083-b2f27d2c4afa	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
67f98f90-ae47-4182-8529-3d5a7b8aea73	9b374efe-f7c1-44ff-a122-52981146b22a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
c29370fc-9699-491b-a23d-121949d729ef	46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
f183af97-b5dc-4e9e-8809-34e53e069529	51f985e0-4f2a-412a-adb5-3adf5399fc75	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
2ea916f7-d564-41ec-8dd3-aa5b5cd9d868	77f2a333-b579-421b-bf3f-5b913e6002b5	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
640ecfbb-7f81-4dbf-b9ee-2bdad3cf0a2f	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
a1fd0fb0-55e8-47c1-be5e-100e6975621c	05a2e58e-0551-40c5-9bff-657768d7e689	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
6b8b55c0-4d16-4db1-ba39-91b9f4455eca	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
3183854e-7490-406d-b659-5b807ca07a35	fb9fbdf7-d05f-40ac-95e1-e10019ed8165	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
08e6cd56-732e-4ab7-bc3e-0fa5e4c32374	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
4e1ef1e1-67b3-4f1a-a984-14a5859ff274	39732090-0856-4b05-a241-df9886a313a5	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
86735524-649d-4972-b1ae-1ab3081044fa	29058c58-14a6-425d-b366-617860f6950f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
a050d22d-4fc5-458a-b74f-eeb2effc160d	97a477c0-7361-421f-b042-2d1e0d9bf5fe	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
95adebff-80d1-482c-98b9-ab9a15e8803a	88c58e99-19a7-4f72-811e-92e21346207a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
68457c56-9c1c-4dc4-be5c-764d72d22e5c	d88a3357-b68c-406e-8951-97df93df66bb	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
441b3ac2-2306-42a5-8cd2-1f293d01e726	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
43f6807f-4d42-40ce-8392-0b3c151505ec	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
c4a1abf2-2323-4cb2-95e9-9fe68f6f01f7	b48dccdc-da4f-4914-bd2e-1971de87fda1	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
abdf5e9b-d4d9-4ecc-be3e-545da0c30cba	3aa46e51-bac0-4711-bf1b-bb396e71fe06	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
7cf627f2-6875-4a6a-a529-69fe584f849d	27917132-65e1-4218-bdfc-40a62550a141	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
b4e98d15-1e4d-47e3-b36b-dcaa8073b661	09cfad53-de06-4a25-88c0-9eaa26345df9	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
ae910601-9efc-4f35-997a-693ee248cc28	7767088d-9354-42d4-8cc7-9d3f2c0bef88	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
7a604825-8ff5-41cf-9998-5c7129923044	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
df38c95a-cf3d-4c89-9e60-0f954122d7a8	f22b7a3d-2d76-483c-af40-dd0c04defb20	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
944f20b0-b8a4-4810-a199-5d6c2c9b812a	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
35aa44fa-fb6f-4bff-8369-68c11f0442b0	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
74b2afce-0c82-4dc9-a253-cb4daf8a307b	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
dc99d1d6-0086-4f5e-a368-d6723b37784a	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
19a8937d-e542-4f07-9139-378cb6cad36e	c45be260-e11b-4fae-8e33-05d0803f4e6b	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
2fd225a1-b8b8-4e44-a68a-2130a36e95fd	24127ef9-9ef3-46f3-88cd-97378e694a6c	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
cba6cbf8-f47c-4fc7-b195-119b950d31fd	24a99997-51d4-446c-bc51-2f46e3ca01a6	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
68ca011d-b985-4d91-b770-144e7382aeb0	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
e80bfafc-e993-4326-a7f4-2d74746593e3	f3c2ef69-090e-4f61-a9d6-954330e86d8f	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
95155b06-8edb-449e-a45d-8806138b47b4	c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
916a4a89-9485-48b4-b60a-42da0d53c0ad	ebb13d64-4568-4cce-9f26-72f5e8a9bc14	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
c3b64d68-369c-41bc-9331-4a595c1141ef	288731cb-561c-4d4a-8951-d8c359a51846	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
87cf5f22-762c-45b2-aa16-00bc1a33cf9d	9f4da9ef-7f9c-407e-91a9-034f623b810f	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
fb6fd230-dee3-41ea-859f-19dd3eb23214	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
747a3570-41b3-4bc3-aaae-4eb77cc33964	653eb679-3e69-442e-b0ff-fd22e3fd8cd3	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
ffd0e613-da83-41f4-a2a7-a6240600e9da	b439edb9-605f-47dc-aee2-8df45976736e	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
3768e05f-2c99-44d6-a17e-531fdd52ba3d	5493c560-ace3-440c-a36a-9538f618007a	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
5726d39f-787b-4834-8907-e74739a22de9	086e9acf-0405-449c-9dee-ee33d80df312	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
6cebdf00-a39f-4de5-af48-3e8d552c7e20	6f4cc747-5450-4514-9a3e-123ace9cb232	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
1750717b-5d51-40e6-b385-aa0766b93379	93421699-5c88-4dee-995c-f74222b3a1a7	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
d9403fc8-192c-4b3e-bb2b-b591e4a8a7b7	c070d24f-76eb-4684-8459-dbb7977acc46	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
d6a2a34b-f80a-4102-8a57-0da9e6ff50d5	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
2bbe8a4c-35bf-4300-b7ae-70e10978b760	d81d630f-3d80-4168-923b-65eefae629c9	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
1cec2a2e-29f7-49dd-b6bd-84867485225f	60f259e8-d103-4fda-859d-b6c053cc6004	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
9359ae1d-1158-4e1b-9ca7-22408ca1296f	eedb53fb-3a3e-4f32-a232-ffc97a4464ae	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
4a5fbe2a-5f51-4259-aa73-3638ef86a03e	59d15cbe-0f30-4057-ab75-09ba8776b7d6	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
2c9896b1-c4f9-461f-af70-7f954047a6cd	4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d884a2fb-5b29-4455-a429-5dab2d05249d	cb332a03-aeee-435c-86de-b7d577e3817b	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
acb611b3-475c-41e2-93bf-9a840bac856c	1897bc12-5b09-412d-a3c7-e001aedff81d	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
7608d5e9-6782-47da-a6b4-d0d0c7c0c0f2	0e143513-3b89-4e2d-9865-8a848ad6c8db	be0d9830-2392-45b9-96d1-805f9a70063b
9039b553-8a98-4864-8373-94f10ef02eeb	357d5bae-8470-4d03-ab92-68cafedcfa16	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
6a01dbf5-e252-49a7-916a-9e607158dd2d	ef82cdff-1b9d-4396-8d2b-f7b751a50006	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
3ee0ca4c-833a-4b7d-8200-3220e73263ea	4689bc27-a161-431b-bf58-456f66307fd9	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
9cbbc5c2-00f6-43d6-996f-82ebe6fe508e	065c9534-8c0a-408f-892a-914d172f90c7	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
a9b09a48-c505-4147-94e1-c9804766ec3f	af0187d0-e84c-401f-9ffc-70194b099309	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
30ecb1ba-8e54-45de-a8f3-e5579e08bd77	ed06f369-9d92-4d07-8b9a-6488694b4d1c	be0d9830-2392-45b9-96d1-805f9a70063b
b111513d-57a4-4753-a82e-2fbe4570659e	c2ebd879-c09a-43d3-8088-6b0c2445c0b1	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b25aba08-6c21-4424-9563-4695fe10da5e	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	be0d9830-2392-45b9-96d1-805f9a70063b
2b103186-5cc0-4300-b05e-113fff0ede9a	e296b896-76d1-4cf0-b092-e637a5909ec4	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b51c3591-fb16-4ebd-a92b-f77f0e99861f	77788f4d-5d0e-458b-b7cb-297e907b2eff	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
36ef810e-c37b-470f-9d32-6cc37787a1b8	82dbdde0-be98-462f-bef7-b86ed1a98336	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
af267847-498d-4aaf-98bf-f4028a5d96dd	23476194-48d6-447c-b4e9-2ab6ad7b1a48	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fde8b764-4d42-4f54-b6e3-d02c375b8ff8	0576302e-61ce-49b3-b265-28b7458e5c27	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
f0ddddf1-3a64-4727-8c32-20a30f67f5d8	58405e28-6fff-42e1-8567-01466845abfb	be0d9830-2392-45b9-96d1-805f9a70063b
14345e2e-69a9-46b7-bcdf-34644bee979f	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	be0d9830-2392-45b9-96d1-805f9a70063b
0f9a429b-3cc6-4cc7-8b87-6eac9be48744	68fc3f82-6076-412f-b986-20d3cb1bceaa	be0d9830-2392-45b9-96d1-805f9a70063b
319d3ecc-08f8-43d5-ae2c-abb9e14989cd	ccad896d-13db-4fc8-be92-574af1d80341	be0d9830-2392-45b9-96d1-805f9a70063b
e60349f0-099c-4e15-98f6-05c6c49b312e	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	be0d9830-2392-45b9-96d1-805f9a70063b
1273d8fe-fa21-48b1-9a46-70f8b7a2bd42	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	be0d9830-2392-45b9-96d1-805f9a70063b
2a0f8058-8e99-4896-8c97-654a7849782f	6f4cc747-5450-4514-9a3e-123ace9cb232	be0d9830-2392-45b9-96d1-805f9a70063b
74c2cabf-a724-4e45-8a47-8fa520e9dc88	93421699-5c88-4dee-995c-f74222b3a1a7	be0d9830-2392-45b9-96d1-805f9a70063b
71147227-1b5a-4675-b581-4307b2edf8a1	c070d24f-76eb-4684-8459-dbb7977acc46	be0d9830-2392-45b9-96d1-805f9a70063b
e1970b34-9246-486d-8478-9c3758c3b72c	3c53a6aa-9d98-43be-bd14-4c5501e58ee5	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0c488d2e-4f90-4e29-bf57-05bb8cb02da5	48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
ae7ea925-330d-4b65-8bc2-9ebbd65f92ce	742c55df-02ca-4d3e-909a-ac8c412ff735	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
3d67f9ea-02d2-4ead-a58e-2d4200c5f692	3ed50a46-3620-43a1-b8ea-5a1ebc968851	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
53bba6d2-c31d-4ed2-9559-e2a254e3b096	080f0752-21ba-4d3e-97c2-4846dfb4f648	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
4ac597f1-f84a-48cf-865d-33414e560f26	01698813-a221-4b56-8bf5-ec3ec9956d97	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
cff46718-76eb-49f3-a801-c3eb97a01389	cab6e3e8-cff1-42da-959f-e12705adc760	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
382b7e6c-9b2c-453e-81bb-a3a87a595a96	c5e2efa4-c8c7-4c79-9823-59e829a4847e	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fae39b08-a7a4-4319-86b4-f9c0791e5b43	1c80833a-d9f9-4b31-ac77-55c282781167	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
90fd1a9e-1a0f-4717-a167-c1e05fb704ff	4e471ceb-766f-4575-ab73-856e07ceb1a8	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d57f58e1-54c6-44bd-a971-8262d0bf56dd	90df4bc3-dbb0-4e92-ab9a-0aac895e5295	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b0666c0b-3ebf-4ed5-aaff-53c7049e45e5	e0be0b73-1f90-4c04-8da5-86513c148a16	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
9d91f7a7-3ea0-4859-b927-9258684b836e	db98b40b-91a4-4434-ba2c-712e3990b82a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
6282783a-de99-4185-b696-04fb7d0a081f	03ed060b-56dd-4ecb-9072-a736c81ae27f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
a7b9f8a2-ba48-4164-b837-9a88c821feff	ad55ef62-54a6-41c0-bba9-830a06822659	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
1a034b36-9c5a-4ca5-b07b-ab4db9a31a23	e18c306e-1d61-4e6e-a2cf-c41ba357b406	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
e1b45545-76aa-492c-ae57-d8be3bff2f96	b995956c-75be-4635-9022-af1f5376be2b	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
884e1ec9-7c4d-45f4-919c-1d9ef5141c14	aa061751-3180-40dc-b101-f866d7e4dc2f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
04473a58-4d78-4708-811b-d75003533963	b4aee72a-2ac5-4748-b846-bd1828139858	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
8a5a2948-65e5-4502-96e6-4cc5972409be	5a8b13c6-3997-4846-b037-2b679f5c31e8	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
cbf7e0c5-433a-4999-93d4-a27819e85f99	a0dee210-7e5b-4309-ad4e-7ccdf1e7e82e	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
2edede79-7c05-483c-9b15-7ecf0d01c733	7ea8ce89-9ce7-4a89-a332-7baad7bbe475	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0fabecc3-f439-4539-9fe5-74e2bb6a5047	72793652-a596-4b80-bbe4-73d10cca1379	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
647d6cb4-1da1-4632-a401-91ff5494b7a2	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	be0d9830-2392-45b9-96d1-805f9a70063b
a2b04500-cd07-46df-b21c-0e36011e26e3	72237dba-0beb-450e-9df6-ea38d3255cab	be0d9830-2392-45b9-96d1-805f9a70063b
d6536988-956c-4361-9493-d94df09ec0ad	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
99b649c5-e85e-45b9-9ecd-8efe3269111d	ae18fd5a-e0ba-4ace-a637-902317a67773	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
21364540-2868-458b-8a6a-6dd1e6871457	72237dba-0beb-450e-9df6-ea38d3255cab	bf9e11be-08a0-4af3-a1f8-cb494ea8ded7
a27cd900-3959-4eb7-b221-67c91986b8e3	3293247e-3d3c-4749-ae91-1a991e27b8be	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
89b292af-2944-4d08-bdac-3efa4b7b055f	66cb5fb9-5bdc-4378-8010-6b58f8cedb5f	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d2834836-e8f2-4470-8bb6-60e8209f1442	5a229d4e-beee-4f0a-ba2f-888f04a149e5	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b80ca0d7-9ce0-4bd1-9f01-33fa9ea56099	136f2190-6f80-46d4-a124-d512e508be5e	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
fbd38b40-0d4b-4a14-a09c-9a6f0d896f24	0e143513-3b89-4e2d-9865-8a848ad6c8db	be30c63e-6c12-4e47-a2d5-26d886544e46
f72b45a5-7140-4097-9e66-858d1498ce89	ed06f369-9d92-4d07-8b9a-6488694b4d1c	be30c63e-6c12-4e47-a2d5-26d886544e46
0335f3ec-8afb-49c5-9d28-7550a42eae71	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	be30c63e-6c12-4e47-a2d5-26d886544e46
c153507c-eb94-47e0-8109-a66d59cc636e	58405e28-6fff-42e1-8567-01466845abfb	be30c63e-6c12-4e47-a2d5-26d886544e46
b787c588-da62-4013-88d4-50868aa7aa99	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	be30c63e-6c12-4e47-a2d5-26d886544e46
a0d2c633-3a54-48a8-b2ee-91a12a3116cf	68fc3f82-6076-412f-b986-20d3cb1bceaa	be30c63e-6c12-4e47-a2d5-26d886544e46
e281e815-3d72-421f-be87-3acddc29807b	ccad896d-13db-4fc8-be92-574af1d80341	be30c63e-6c12-4e47-a2d5-26d886544e46
242cd9b9-fa93-439a-b42d-0db1db9b7cc0	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	be30c63e-6c12-4e47-a2d5-26d886544e46
0329dca9-013b-4e10-b6a2-1d30391e6b88	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	be30c63e-6c12-4e47-a2d5-26d886544e46
f804a075-f393-4aef-85bf-b77c30d04f10	6f4cc747-5450-4514-9a3e-123ace9cb232	be30c63e-6c12-4e47-a2d5-26d886544e46
8913bddc-5276-4b5f-be46-55557773fd72	93421699-5c88-4dee-995c-f74222b3a1a7	be30c63e-6c12-4e47-a2d5-26d886544e46
06a1057d-68fe-4424-890e-0794b809dc2f	c070d24f-76eb-4684-8459-dbb7977acc46	be30c63e-6c12-4e47-a2d5-26d886544e46
4a8a88d4-51e9-4217-a73f-f5db06509473	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	be30c63e-6c12-4e47-a2d5-26d886544e46
90ff1d21-d568-41f1-afe5-592ad8e671e0	72237dba-0beb-450e-9df6-ea38d3255cab	be30c63e-6c12-4e47-a2d5-26d886544e46
d41aebac-de81-4d38-831b-053574197cce	3ab532fe-274c-4755-b8c4-8055ea124209	29242052-f1fd-458e-9119-3fa6a9d5ce98
ca2f9078-332d-4a24-88cc-796e4cd5babf	a8de69fd-5856-42cc-8812-fef457a99079	29242052-f1fd-458e-9119-3fa6a9d5ce98
1119f850-1ddb-4d8f-9d6d-8fb10b45f03f	136f2190-6f80-46d4-a124-d512e508be5e	29242052-f1fd-458e-9119-3fa6a9d5ce98
184c2286-7db4-4c9b-a41a-c09aab7c50f7	3ab532fe-274c-4755-b8c4-8055ea124209	283fd8af-5010-48d6-9ac2-314f013b13bc
6726fc5e-7b2e-4280-bc8f-c519924adebc	a8de69fd-5856-42cc-8812-fef457a99079	283fd8af-5010-48d6-9ac2-314f013b13bc
cd533b07-0c4d-4940-bfb9-184c7d96ce44	136f2190-6f80-46d4-a124-d512e508be5e	283fd8af-5010-48d6-9ac2-314f013b13bc
030a6862-3243-45d5-8ea9-35bcf88b423b	0e143513-3b89-4e2d-9865-8a848ad6c8db	283fd8af-5010-48d6-9ac2-314f013b13bc
fe5ff1d6-f893-49d9-9a4d-0985ce1823a8	ed06f369-9d92-4d07-8b9a-6488694b4d1c	283fd8af-5010-48d6-9ac2-314f013b13bc
2cabcd92-8433-446b-bd7e-013a377d9fdb	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	283fd8af-5010-48d6-9ac2-314f013b13bc
b0de06d9-d3ff-43eb-8993-431264b12a23	58405e28-6fff-42e1-8567-01466845abfb	283fd8af-5010-48d6-9ac2-314f013b13bc
8558ead5-aa0e-45bc-8702-9d51797599dc	aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	283fd8af-5010-48d6-9ac2-314f013b13bc
2392c815-ac0f-4bee-8076-74861379d8cc	68fc3f82-6076-412f-b986-20d3cb1bceaa	283fd8af-5010-48d6-9ac2-314f013b13bc
d5d22272-b69c-425d-87ca-01c8a74b6790	ccad896d-13db-4fc8-be92-574af1d80341	283fd8af-5010-48d6-9ac2-314f013b13bc
7b333a87-4d18-4c6e-9b35-dde7a5c243f4	03264ac4-7754-48ef-a824-5269d1f84e4b	283fd8af-5010-48d6-9ac2-314f013b13bc
80d3c0e1-c365-4eb8-b77a-a17e0bd02a73	51f985e0-4f2a-412a-adb5-3adf5399fc75	283fd8af-5010-48d6-9ac2-314f013b13bc
5910d3d6-e660-4fe3-9d09-0afeb878b0bb	77f2a333-b579-421b-bf3f-5b913e6002b5	283fd8af-5010-48d6-9ac2-314f013b13bc
9f398e0f-86e8-4bda-8b89-e8d6a1f8fef6	0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	283fd8af-5010-48d6-9ac2-314f013b13bc
518d767c-0a6a-4541-8f80-e3fc9eaa34cf	05a2e58e-0551-40c5-9bff-657768d7e689	283fd8af-5010-48d6-9ac2-314f013b13bc
f37484f0-d2d2-4fbf-8237-e7e83757ecf2	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	283fd8af-5010-48d6-9ac2-314f013b13bc
a3d50862-0a6a-4f12-84f7-eb19295f2ccb	fb9fbdf7-d05f-40ac-95e1-e10019ed8165	283fd8af-5010-48d6-9ac2-314f013b13bc
9c63c133-da1e-426d-b52a-c57ec66eec4f	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	283fd8af-5010-48d6-9ac2-314f013b13bc
909e7202-0d0c-455c-8ba3-c9277144fa0a	39732090-0856-4b05-a241-df9886a313a5	283fd8af-5010-48d6-9ac2-314f013b13bc
d671b202-468b-4bbf-81f7-812e9da32178	9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	283fd8af-5010-48d6-9ac2-314f013b13bc
fe1d31cc-7421-4155-99a5-1b80932355e0	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	283fd8af-5010-48d6-9ac2-314f013b13bc
0c622256-8625-425c-a788-03fc59fe96f0	b48dccdc-da4f-4914-bd2e-1971de87fda1	283fd8af-5010-48d6-9ac2-314f013b13bc
a839c3df-622e-47d8-aeca-2bbb7bf3f34d	3aa46e51-bac0-4711-bf1b-bb396e71fe06	283fd8af-5010-48d6-9ac2-314f013b13bc
ba16c873-e892-400e-bd1f-4034865e2d88	27917132-65e1-4218-bdfc-40a62550a141	283fd8af-5010-48d6-9ac2-314f013b13bc
9618440a-58ed-4338-bf07-8e7b1f9ad044	09cfad53-de06-4a25-88c0-9eaa26345df9	283fd8af-5010-48d6-9ac2-314f013b13bc
9d72ae89-97b5-4f78-8340-27eb5a89072c	7767088d-9354-42d4-8cc7-9d3f2c0bef88	283fd8af-5010-48d6-9ac2-314f013b13bc
1b84de48-65bf-4a3d-a8b9-0b4c85bcf4de	b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	283fd8af-5010-48d6-9ac2-314f013b13bc
673a5a55-f09a-4f85-b0d3-1e989d151c9a	f22b7a3d-2d76-483c-af40-dd0c04defb20	283fd8af-5010-48d6-9ac2-314f013b13bc
c154c061-ba39-4982-9829-0c0b0bde6076	86ce4e37-3234-4bb4-bb03-b1d3742ac42e	283fd8af-5010-48d6-9ac2-314f013b13bc
c4a3067e-4816-4147-813f-b095ce0af345	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	283fd8af-5010-48d6-9ac2-314f013b13bc
888444eb-06ec-4d51-a5f2-acf35d8e5170	3566c88a-f2b5-4c1e-b8a3-68a26b49e158	283fd8af-5010-48d6-9ac2-314f013b13bc
4e19c917-cd02-4c25-b11f-a63bdf1b1a3e	9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	283fd8af-5010-48d6-9ac2-314f013b13bc
b1981be5-a746-4612-a23f-2103048c76e3	c45be260-e11b-4fae-8e33-05d0803f4e6b	283fd8af-5010-48d6-9ac2-314f013b13bc
428c7890-7e90-4a5c-8985-0be9a3ceddf4	24127ef9-9ef3-46f3-88cd-97378e694a6c	283fd8af-5010-48d6-9ac2-314f013b13bc
28241fed-6c43-4cbf-a81e-622df838f4a4	24a99997-51d4-446c-bc51-2f46e3ca01a6	283fd8af-5010-48d6-9ac2-314f013b13bc
fd34740d-387e-4a05-ac2f-0ff29bd8e5f2	d344cbbc-eb28-446a-8b92-38bd340cf3b2	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
293ea7b2-5fc5-4fdc-99be-01739e7eed4e	9fc4b56e-d391-4564-9dcd-2c47f9e436c3	283fd8af-5010-48d6-9ac2-314f013b13bc
35d1fb32-e3b8-4503-9cb6-fd23b0de7541	f3c2ef69-090e-4f61-a9d6-954330e86d8f	283fd8af-5010-48d6-9ac2-314f013b13bc
78f35a8a-fd47-4156-86fc-f900a959ad18	ed6cdd20-4817-464c-99b5-f7a2084f6a9c	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0268dc15-6360-4850-a93a-b81bfb781a3d	c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	283fd8af-5010-48d6-9ac2-314f013b13bc
0e8bdaa5-9c8d-4481-814d-005dcfc00003	653eb679-3e69-442e-b0ff-fd22e3fd8cd3	283fd8af-5010-48d6-9ac2-314f013b13bc
10e654de-c80a-4df2-afc1-f0523c0e02ee	b439edb9-605f-47dc-aee2-8df45976736e	283fd8af-5010-48d6-9ac2-314f013b13bc
c1c43a8c-a90b-4976-9644-dc1e40d278f3	5493c560-ace3-440c-a36a-9538f618007a	283fd8af-5010-48d6-9ac2-314f013b13bc
402e5663-41f0-41d0-9283-5829a45cc20f	086e9acf-0405-449c-9dee-ee33d80df312	283fd8af-5010-48d6-9ac2-314f013b13bc
1eaba89e-ce15-4908-894a-0d454980dfb3	6f4cc747-5450-4514-9a3e-123ace9cb232	283fd8af-5010-48d6-9ac2-314f013b13bc
afaf46b1-bc9c-4c96-a957-5f1ac0bf7b32	93421699-5c88-4dee-995c-f74222b3a1a7	283fd8af-5010-48d6-9ac2-314f013b13bc
c1addae2-83e0-4a93-b9aa-7c9f738c4bef	c070d24f-76eb-4684-8459-dbb7977acc46	283fd8af-5010-48d6-9ac2-314f013b13bc
ed7cfce1-9ddc-4833-bc5e-b9297907f1eb	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	283fd8af-5010-48d6-9ac2-314f013b13bc
e798d49b-64aa-4708-bff1-e2070d8c06b6	d81d630f-3d80-4168-923b-65eefae629c9	283fd8af-5010-48d6-9ac2-314f013b13bc
0afc75aa-e4cb-4672-a0e2-dd9905b7693d	4ec8150a-0cbe-48f5-9dee-974fcb3e0558	283fd8af-5010-48d6-9ac2-314f013b13bc
36a65d4d-7215-4e76-a96c-b587b74111d1	ae18fd5a-e0ba-4ace-a637-902317a67773	283fd8af-5010-48d6-9ac2-314f013b13bc
335b482f-c3eb-43de-b2e4-98af37b03752	c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	283fd8af-5010-48d6-9ac2-314f013b13bc
aec18f83-a0c0-4ed7-8827-d10002fbfbe4	ebb13d64-4568-4cce-9f26-72f5e8a9bc14	283fd8af-5010-48d6-9ac2-314f013b13bc
ae68b423-ba35-4930-8f8e-ba9d0e7692a9	288731cb-561c-4d4a-8951-d8c359a51846	283fd8af-5010-48d6-9ac2-314f013b13bc
caa5514d-c158-4d5c-9316-54301b2f5ff4	5aae3b9c-87f3-487f-b016-f7109f17d2c0	283fd8af-5010-48d6-9ac2-314f013b13bc
7032c414-e191-42ae-b9e9-f68d35d6e546	9f4da9ef-7f9c-407e-91a9-034f623b810f	283fd8af-5010-48d6-9ac2-314f013b13bc
73d48b55-bdca-45aa-99e0-74235c191a5e	72237dba-0beb-450e-9df6-ea38d3255cab	283fd8af-5010-48d6-9ac2-314f013b13bc
ec19c872-ff6b-43ba-8ef2-674f5853fed7	3ab532fe-274c-4755-b8c4-8055ea124209	f9ba64c8-eac8-4b39-8285-f86611011cfc
a36e59d1-4efc-489b-aed5-0f27d468a2d4	a8de69fd-5856-42cc-8812-fef457a99079	f9ba64c8-eac8-4b39-8285-f86611011cfc
6d6cb89a-1f9d-45bb-8e93-5c19cee09cda	136f2190-6f80-46d4-a124-d512e508be5e	f9ba64c8-eac8-4b39-8285-f86611011cfc
b479945b-edd9-4122-b521-7fe9719b69b8	ab1aa816-baed-4d33-bd9b-80d4b83a0341	f9ba64c8-eac8-4b39-8285-f86611011cfc
096ac236-ebf4-40f8-b5a3-b5b077aea340	797cecc2-d1f9-499a-b79d-e1f7e77c7874	f9ba64c8-eac8-4b39-8285-f86611011cfc
0da7d4b6-0142-4f0c-b6b4-d36d88b1d92b	15668ec7-6731-4646-a516-41035ddebe0d	f9ba64c8-eac8-4b39-8285-f86611011cfc
a5f3d578-20b8-4c27-b904-7f3e90211a74	23bcf6d0-fa48-4ae8-8036-07f675e240b7	8457de16-e58d-4dc0-b20c-6a5730c4ed38
d722c042-6fb5-4562-b5c4-7a275c53de74	1e5eb584-f19e-4110-89e7-2d63a08d727a	be0d9830-2392-45b9-96d1-805f9a70063b
41821d4b-d0a6-4580-8bfe-bf1ab5f5ca1d	23bcf6d0-fa48-4ae8-8036-07f675e240b7	be0d9830-2392-45b9-96d1-805f9a70063b
50f530cb-416d-47fb-a5f9-95935b4a276a	1e5eb584-f19e-4110-89e7-2d63a08d727a	be30c63e-6c12-4e47-a2d5-26d886544e46
a9126a37-914d-4c05-99a4-7753ff56e283	23bcf6d0-fa48-4ae8-8036-07f675e240b7	be30c63e-6c12-4e47-a2d5-26d886544e46
7618ec5f-d19f-4dee-80e8-3c9aa95e25c9	1897bc12-5b09-412d-a3c7-e001aedff81d	8457de16-e58d-4dc0-b20c-6a5730c4ed38
531ae8d6-5c43-4dcf-9587-b110f2d1aee8	357d5bae-8470-4d03-ab92-68cafedcfa16	8457de16-e58d-4dc0-b20c-6a5730c4ed38
dde35897-c227-481c-89bc-f2aa1c18e3af	ef82cdff-1b9d-4396-8d2b-f7b751a50006	8457de16-e58d-4dc0-b20c-6a5730c4ed38
435d4a9e-1821-4a52-9101-f03d12771801	4689bc27-a161-431b-bf58-456f66307fd9	8457de16-e58d-4dc0-b20c-6a5730c4ed38
b09c3142-1529-4771-928c-36cf5180ba04	065c9534-8c0a-408f-892a-914d172f90c7	8457de16-e58d-4dc0-b20c-6a5730c4ed38
d96d4078-7eab-4146-b9cb-1c7f82f8e8fd	07679b82-27d2-4811-85e6-8a6e73461613	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
e4d7052f-b6ca-45c5-8871-9d52694623b0	d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
85cb5888-87c0-431e-98e1-77cdc5412a98	149d7564-5e2f-4472-a637-d1db61dfdc81	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
879b419d-cda3-4c8e-a318-93476015d8b2	81d8c6fd-f427-409c-aac9-d9dac8e33f97	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
ae8bc083-6258-45af-97a5-330a29010955	850c0f8d-2235-431c-8e0a-d02baa738c58	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0edef4d1-16c7-4673-9c0f-df78157e6f3b	e428e7fe-5311-447d-a13f-30d8bbf04653	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
580c955f-66cc-42ed-afc2-16c05dff30a6	3800a410-1da5-46a2-958f-163f98034d55	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
46618d57-7eeb-4638-a717-14d67049f568	81d5505b-df37-4cef-928c-10730baa7b10	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
b4119030-06e3-404f-b58c-7598735060ec	08454882-e75c-4933-b522-6687fda8cb24	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
9f33c142-1156-4190-9938-3f0e30a75e30	2c2961d0-ef1e-4c1c-8f28-94db86f33b16	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
14025e40-d354-4bf4-889d-fe1dc72bbc7b	67d8761f-764e-4e5e-92df-bcb2bfae02ee	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
81b4a848-30d8-49b2-a2a9-5a6691708542	996af46e-f32c-4b26-8109-2b942d1aa8b8	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
0b2fd262-2691-4970-be4d-d56d1cf5fa47	b89d0f67-90fe-4ca8-9122-65eda7c53aa2	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
875f9350-90a1-43f5-aac8-b02c6bf1ceaf	18d4ae82-68bc-435c-9e72-bfef6bc6b2e9	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
d892c7e6-8cc3-41d8-aa4e-460bcb5d3233	48d6de3b-4af6-4ea3-b796-aca905fa07f3	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
3fea735d-a584-4527-9485-208c28e88e49	8876a714-7d09-43d4-88ea-795c559a72a1	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
94d1ffea-6445-4ea9-8f81-24003967e588	3b618a8e-5499-49a3-a980-dab3944a50d4	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
edae18e7-d744-460b-8bcc-b580b943c2c3	87c1d68e-93a3-4ee9-b0a6-30989584c056	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
403c9958-1ddb-47e1-b2bd-647a7e376993	1da7c457-afd1-484c-b24c-6d984baecdbf	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
8c17acf8-abda-46ba-87d7-ac19ff329351	135baf4c-f8cd-4597-bdb5-a49f4bb53e42	2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9
\.


--
-- Data for Name: risk_result; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.risk_result (risk_no, risk_result, risk_threshold, create_time, api_type, payer_acc_no, action_time, eva_execute_type, eva_score, money_type, order_no, score, product_type) FROM stdin;
2020042320032685519983	0	9	2020-04-23 12:03:26.04945	exchange	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-23 19:03:26.044537488 +0700 +07 m=+10577.836197671	online	10	usd	2020042320032673158468	8	exchange
2020042320033889174421	0	9	2020-04-23 12:03:38.504535	exchange	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-23 19:03:38.50137908 +0700 +07 m=+10590.293039223	online	10	usd	2020042320033872751669	8	exchange
2020042320043782864843	0	9	2020-04-23 12:04:37.158525	exchange	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-23 19:04:37.154609486 +0700 +07 m=+10648.946269629	online	10	khr	2020042320043796136626	8	exchange
2020042320274748934623	0	9	2020-04-23 12:27:47.242718	exchange	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-23 19:27:47.238055253 +0700 +07 m=+12039.029715386	online	10	usd	2020042320274738216920	8	exchange
2020042413493906888030	0	9	2020-04-24 05:49:39.301013	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 12:49:39.293473995 +0700 +07 m=+7128.297379129	online	10	usd	2020042413493958372565	8	save_money
2020042413494793523299	0	9	2020-04-24 05:49:47.358123	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 12:49:47.354353509 +0700 +07 m=+7136.358258633	online	10	usd	2020042413494756368434	8	save_money
2020042413534254661130	0	9	2020-04-24 05:53:42.164062	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 12:53:42.159177272 +0700 +07 m=+7371.163082416	online	10	usd	2020042413534231343012	8	save_money
2020042414023012137944	0	9	2020-04-24 06:02:30.607948	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 13:02:30.602504707 +0700 +07 m=+7899.606409851	online	10	usd	2020042414023058287986	8	save_money
2020042414041140442131	0	9	2020-04-24 06:04:11.781471	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-24 13:04:11.776938642 +0700 +07 m=+8000.780843766	online	10	usd	2020042414041126353596	8	save_money
2020042610450325788993	0	9	2020-04-26 09:45:03.98104	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-26 09:45:03.972702866 +0700 +07 m=+150302.957021960	online	10	usd	2020042610450398550018	8	save_money
2020042610453028097413	0	9	2020-04-26 09:45:30.644965	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-26 09:45:30.641739261 +0700 +07 m=+150329.626058325	online	10	usd	2020042610453004687366	8	save_money
2020042610453763718388	0	9	2020-04-26 09:45:37.756242	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-26 09:45:37.753120841 +0700 +07 m=+150336.737439885	online	10	usd	2020042610453739122130	8	save_money
2020042618572087887390	0	9	2020-04-26 17:57:20.398356	save_money	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-26 17:57:20.393355751 +0700 +07 m=+179839.377674805	online	10	usd	2020042618572059971004	8	save_money
2020042619001660858341	0	9	2020-04-26 18:00:16.432152	save_money	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-26 18:00:16.428034817 +0700 +07 m=+180015.412353861	online	10	usd	2020042619001631133543	8	save_money
2020042619004417580697	0	9	2020-04-26 18:00:44.806721	save_money	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-26 18:00:44.803617404 +0700 +07 m=+180043.787936458	online	10	usd	2020042619004475793973	8	save_money
2020042714430825061097	0	9	2020-04-27 13:43:08.216023	save_money	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-27 13:43:08.208480964 +0700 +07 m=+70759.881445115	online	0	usd	2020042714430860352023	8	save_money
2020042716420052265180	0	9	2020-04-27 15:42:00.46079	transfer	c1abbc5a-4280-4996-a866-08b74a21d2fb	2020-04-27 15:42:00.455043054 +0700 +07 m=+77892.128007205	online	10	usd	2020042716420018848193	8	transfer
2020042900003241589472	0	9	2020-04-28 23:00:32.293046	transfer	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-28 23:00:32.284860485 +0700 +07 m=+22680.305207674	online	10	usd	2020042900003201837845	8	transfer
2020042910365409045130	0	9	2020-04-29 09:36:54.110244	save_money	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-29 09:36:54.104581807 +0700 +07 m=+60862.124929016	online	0	usd	2020042910365455690386	8	save_money
2020042910370559618522	0	9	2020-04-29 09:37:05.335862	save_money	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-29 09:37:05.332017619 +0700 +07 m=+60873.352364828	online	10	usd	2020042910370516852515	8	save_money
2020042910412997016030	0	9	2020-04-29 09:41:29.364536	save_money	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-29 09:41:29.35947697 +0700 +07 m=+61137.379824179	online	10	usd	2020042910412975170512	8	save_money
2020042920004776934891	0	9	2020-04-29 19:00:47.982486	save_money	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-29 19:00:47.975334898 +0700 +07 m=+2304.074210383	online	10	usd	2020042920004751389451	8	save_money
2020043011580655439107	0	9	2020-04-30 10:58:06.19318	transfer	0e8d24af-bec7-4f95-b038-c48045f51abf	2020-04-30 10:58:06.185056775 +0700 +07 m=+59742.283932270	online	10	usd	2020043011580619577052	8	transfer
2020043014502404936670	0	9	2020-04-30 13:50:24.179046	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-04-30 13:50:24.173592046 +0700 +07 m=+70080.272467521	online	10	usd	2020043014502473435188	8	save_money
2020043015063939120392	0	9	2020-04-30 14:06:39.9816	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-04-30 14:06:39.976934202 +0700 +07 m=+71056.075809697	online	10	usd	2020043015063911518368	8	save_money
2020043015085186239893	0	9	2020-04-30 14:08:51.69915	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-04-30 14:08:51.695326648 +0700 +07 m=+71187.794202123	online	10	usd	2020043015085147694949	8	save_money
2020043015152293317190	0	9	2020-04-30 14:15:22.407189	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-04-30 14:15:22.403451757 +0700 +07 m=+71578.502327232	online	10	usd	2020043015152267531776	8	save_money
2020043016483508497973	0	9	2020-04-30 15:48:35.357046	save_money	f2961920-ef9b-4ebb-9de8-dcfcc1b21539	2020-04-30 15:48:35.353594361 +0700 +07 m=+77171.452469826	online	10	usd	2020043016483587444644	8	save_money
2020050711022666167294	0	9	2020-05-07 10:02:26.225049	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-05-07 10:02:26.218141395 +0700 +07 m=+661202.317016890	online	10	usd	2020050711022622634565	8	save_money
2020050711273189223415	0	9	2020-05-07 10:27:31.737545	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-05-07 10:27:31.731143179 +0700 +07 m=+662707.830018674	online	10	usd	2020050711273131774598	8	save_money
2020050711275959881881	0	9	2020-05-07 10:27:59.096503	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-05-07 10:27:59.092014274 +0700 +07 m=+662735.190889749	online	10	usd	2020050711275996940484	8	save_money
2020050711290739173419	0	9	2020-05-07 10:29:07.731072	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-05-07 10:29:07.725863984 +0700 +07 m=+662803.824739469	online	10	khr	2020050711290709162496	8	save_money
2020050711293602573628	0	9	2020-05-07 10:29:36.769395	save_money	c07ac194-8a49-4619-ad67-171cb689f987	2020-05-07 10:29:36.765435958 +0700 +07 m=+662832.864311433	online	10	khr	2020050711293658418063	8	save_money
2020050811223141619423	0	9	2020-05-08 10:22:31.636115	save_money	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2020-05-08 10:22:31.627700288 +0700 +07 m=+64263.251218498	online	10	usd	2020050811223180215481	8	save_money
2020050811225435314836	0	9	2020-05-08 10:22:54.936442	save_money	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2020-05-08 10:22:54.931829422 +0700 +07 m=+64286.555347642	online	10	khr	2020050811225407234823	8	save_money
\.


--
-- Data for Name: role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.role (role_no, role_name, create_time, modify_time, acc_type, def_type, master_acc, is_delete) FROM stdin;
2804a64e-0d6a-4553-9b45-a6f8ae3dc5b9	后台管理员	2018-06-07 07:51:46.287437	2018-06-25 02:28:44.231932	1	1	00000000-0000-0000-0000-000000000000	0
ad726a2c-6581-4675-a770-a4b58797b813	运营-明	2019-09-18 16:53:49.044707	1970-01-01 00:00:00	3	0	00000000-0000-0000-0000-000000000000	0
fbb11569-57f0-4b75-ae29-3854e26a0585	运营2-1	2019-09-18 16:53:56.81534	2019-09-18 16:55:30.341407	3	1	00000000-0000-0000-0000-000000000000	0
984bb718-3808-4911-ac96-925444c6a716	普通商户	2019-09-18 17:00:53.557317	1970-01-01 00:00:00	5	1	00000000-0000-0000-0000-000000000000	0
eaa6135b-2f8e-479c-9336-2dbff8b03ce3	商户a	2019-09-30 10:43:42.317187	1970-01-01 00:00:00	5	0	00000000-0000-0000-0000-000000000000	0
5e40ef5d-fa01-40fb-a3d3-32fa378ce287	商户bbbbbbbb	2019-09-30 10:44:19.239339	1970-01-01 00:00:00	5	0	00000000-0000-0000-0000-000000000000	1
8bcc6602-0b78-4d78-936d-fd50c3ae893f	123123123	2019-12-04 16:17:34.889386	1970-01-01 00:00:00	1	0	00000000-0000-0000-0000-000000000000	1
a36c637e-32f1-406d-82d4-6617d2fda821	123123123	2019-12-04 16:20:53.50281	1970-01-01 00:00:00	5	0	00000000-0000-0000-0000-000000000000	1
b4639be8-2b98-4cca-9692-95f305cef25d	运营1-商户	2019-10-15 11:29:49.18975	2019-12-04 18:01:04.295448	5	0	00000000-0000-0000-0000-000000000000	0
85a78368-21eb-45db-9feb-07d054e5b445	大卡商	2019-09-18 17:12:47.932414	2019-09-18 17:21:36.705208	16	0	00000000-0000-0000-0000-000000000000	1
29242052-f1fd-458e-9119-3fa6a9d5ce98	普通代理	2019-09-18 17:02:58.253367	1970-01-01 00:00:00	4	1	00000000-0000-0000-0000-000000000000	1
8457de16-e58d-4dc0-b20c-6a5730c4ed38	MSC运营	2019-09-18 16:54:05.469463	2019-09-18 17:03:38.96391	3	1	00000000-0000-0000-0000-000000000000	1
e920b75e-3c45-496b-9df0-f092455381eb	运营1-小卡商	2019-09-18 17:35:06.825055	1970-01-01 00:00:00	17	1	00000000-0000-0000-0000-000000000000	1
0ad57fbe-a74e-40ad-ab48-afa88c1dceea	运营1-卡商管理者	2019-09-18 17:30:11.826326	1970-01-01 00:00:00	16	0	00000000-0000-0000-0000-000000000000	1
be30c63e-6c12-4e47-a2d5-26d886544e46	小代理	2019-09-18 17:56:28.231667	2019-09-18 18:45:23.834373	13	0	00000000-0000-0000-0000-000000000000	1
be0d9830-2392-45b9-96d1-805f9a70063b	大代理	2019-09-18 17:54:38.151817	1970-01-01 00:00:00	13	1	00000000-0000-0000-0000-000000000000	1
1ec06744-dd1e-4664-8666-61a1f8974b16	运营2-1-服务商	2019-10-29 17:53:42.71019	2019-12-04 18:00:52.324438	13	1	00000000-0000-0000-0000-000000000000	1
283fd8af-5010-48d6-9ac2-314f013b13bc	运营-运营csa	2019-09-27 15:42:05.180201	1970-01-01 00:00:00	3	0	00000000-0000-0000-0000-000000000000	1
b2222924-ce4e-444f-8e1a-c3479bb99467	服务商-运营-明	2019-09-18 17:59:02.696273	1970-01-01 00:00:00	13	1	00000000-0000-0000-0000-000000000000	1
94176a90-02bb-4dea-a306-ebdcf3c02975	运营-运营s	2019-09-18 17:45:25.054334	1970-01-01 00:00:00	3	0	00000000-0000-0000-0000-000000000000	1
bf9e11be-08a0-4af3-a1f8-cb494ea8ded7	运营-运营cs1	2019-09-20 19:09:21.748461	1970-01-01 00:00:00	3	0	00000000-0000-0000-0000-000000000000	1
75e36809-c0cf-4ce5-a0a4-1cecbd6f1d6a	运营2-1-子商户	2019-10-08 18:18:58.951902	1970-01-01 00:00:00	6	1	00000000-0000-0000-0000-000000000000	1
f9ba64c8-eac8-4b39-8285-f86611011cfc	运营2-1-门店	2019-10-08 18:17:33.50539	1970-01-01 00:00:00	9	1	00000000-0000-0000-0000-000000000000	1
\.


--
-- Data for Name: servicer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.servicer (servicer_no, account_no, addr, create_time, is_delete, use_status, commission_sharing, income_authorization, outgo_authorization, open_idx, contact_person, contact_phone, contact_addr, lat, lng, password, modify_time, income_sharing, scope, scope_off) FROM stdin;
f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	d4a4ca0f-973e-484a-80b7-c40187aeda3f	\N	2020-05-08 10:19:42.483307	0	1	0	1	1	0				0	0	f08587d6f466c0abc1252f15dc52a865	2020-05-08 10:21:09.000415	0	0	0
\.


--
-- Data for Name: servicer_count; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.servicer_count (servicer_no, currency_type, in_num, in_amount, out_num, out_amount, profit_num, modify_time, profit_amount, recharge_num, recharge_amount, withdraw_num, withdraw_amount) FROM stdin;
dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	usd	3	12500	0	0	3	2020-04-27 03:00:00	15	0	0	0	0
dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	khr	0	0	0	0	0	2020-04-27 03:00:00	0	0	0	0	0
60f15170-c1db-41b0-bb3d-14185ab43d28	usd	4	14300	0	0	4	2020-04-30 03:00:00	4	0	0	0	0
60f15170-c1db-41b0-bb3d-14185ab43d28	khr	0	0	0	0	0	2020-04-30 03:00:00	0	0	0	0	0
99f047f6-30d0-4e76-853d-57d124c76cc0	usd	3	22200	0	0	3	2020-05-01 03:00:00	6	0	0	0	0
99f047f6-30d0-4e76-853d-57d124c76cc0	khr	0	0	0	0	0	2020-05-01 03:00:00	0	0	0	0	0
cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	6	55000	0	0	6	2020-05-08 03:00:00	30	0	0	0	0
cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	khr	2	200	0	0	2	2020-05-08 03:00:00	100	0	0	0	0
\.


--
-- Data for Name: servicer_count_list; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.servicer_count_list (servicer_no, currency_type, create_time, in_num, in_amount, out_num, out_amount, profit_num, profit_amount, recharge_num, recharge_amount, withdraw_num, withdraw_amount, id, dates, is_counted) FROM stdin;
cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	2020-05-08 02:00:00	3	25000	0	0	3	15	0	0	0	0	188	2020-05-07	1
cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	khr	2020-05-08 02:00:00	2	200	0	0	2	100	0	0	0	0	189	2020-05-07	1
60f15170-c1db-41b0-bb3d-14185ab43d28	usd	2020-04-25 02:00:00	1	1000	0	0	1	1	0	0	0	0	174	2020-04-24	1
60f15170-c1db-41b0-bb3d-14185ab43d28	khr	2020-04-25 02:00:00	0	0	0	0	0	0	0	0	0	0	175	2020-04-24	1
dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	usd	2020-04-25 02:00:00	2	2500	0	0	2	10	0	0	0	0	176	2020-04-24	1
dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	khr	2020-04-25 02:00:00	0	0	0	0	0	0	0	0	0	0	177	2020-04-24	1
99f047f6-30d0-4e76-853d-57d124c76cc0	usd	2020-04-27 02:00:00	2	11100	0	0	2	4	0	0	0	0	178	2020-04-26	1
99f047f6-30d0-4e76-853d-57d124c76cc0	khr	2020-04-27 02:00:00	0	0	0	0	0	0	0	0	0	0	179	2020-04-26	1
dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	usd	2020-04-27 02:00:00	1	10000	0	0	1	5	0	0	0	0	180	2020-04-26	1
dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	khr	2020-04-27 02:00:00	0	0	0	0	0	0	0	0	0	0	181	2020-04-26	1
60f15170-c1db-41b0-bb3d-14185ab43d28	usd	2020-04-30 02:00:00	3	13300	0	0	3	3	0	0	0	0	182	2020-04-29	1
60f15170-c1db-41b0-bb3d-14185ab43d28	khr	2020-04-30 02:00:00	0	0	0	0	0	0	0	0	0	0	183	2020-04-29	1
99f047f6-30d0-4e76-853d-57d124c76cc0	usd	2020-05-01 02:00:00	1	11100	0	0	1	2	0	0	0	0	184	2020-04-30	1
99f047f6-30d0-4e76-853d-57d124c76cc0	khr	2020-05-01 02:00:00	0	0	0	0	0	0	0	0	0	0	185	2020-04-30	1
cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	2020-05-01 02:00:00	3	30000	0	0	3	15	0	0	0	0	186	2020-04-30	1
cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	khr	2020-05-01 02:00:00	0	0	0	0	0	0	0	0	0	0	187	2020-04-30	1
\.


--
-- Data for Name: servicer_img; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.servicer_img (servicer_img_no, img_id, img_type, create_time, servicer_no, is_delete) FROM stdin;
\.


--
-- Data for Name: servicer_profit_ledger; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.servicer_profit_ledger (log_no, amount_order, servicefee_amount_sum, split_proportion, actual_income, payment_time, servicer_no, currency_type, order_type) FROM stdin;
2020042413534231343012	1000	10	1000	1	2020-04-24 05:53:42	60f15170-c1db-41b0-bb3d-14185ab43d28	usd	1
2020042414023058287986	1000	10	5000	5	2020-04-24 06:02:31	dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	usd	1
2020042414041126353596	1500	10	5000	5	2020-04-24 06:04:12	dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	usd	1
2020042610453739122130	10000	10	5000	5	2020-04-26 10:00:00	dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	usd	1
2020042619001631133543	10000	10	2000	2	2020-04-26 18:00:16	99f047f6-30d0-4e76-853d-57d124c76cc0	usd	1
2020042619004475793973	1100	10	2000	2	2020-04-26 18:00:45	99f047f6-30d0-4e76-853d-57d124c76cc0	usd	1
2020042910370516852515	11100	11	1000	1	2020-04-29 09:37:05	60f15170-c1db-41b0-bb3d-14185ab43d28	usd	1
2020042910412975170512	1100	10	1000	1	2020-04-29 09:41:29	60f15170-c1db-41b0-bb3d-14185ab43d28	usd	1
2020042920004751389451	1100	10	1000	1	2020-04-29 19:00:48	60f15170-c1db-41b0-bb3d-14185ab43d28	usd	1
2020043015063911518368	10000	10	5000	5	2020-04-30 14:06:40	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	1
2020043015085147694949	10000	10	5000	5	2020-04-30 14:08:52	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	1
2020043015152267531776	10000	10	5000	5	2020-04-30 14:15:22	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	1
2020043016483587444644	11100	11	2000	2	2020-04-30 15:48:35	99f047f6-30d0-4e76-853d-57d124c76cc0	usd	1
2020050711022622634565	5000	10	5000	5	2020-05-07 10:02:26	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	1
2020050711273131774598	10000	10	5000	5	2020-05-07 10:27:32	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	1
2020050711275996940484	10000	10	5000	5	2020-05-07 10:27:59	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	usd	1
2020050711290709162496	100	100	5000	50	2020-05-07 10:29:08	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	khr	1
2020050711293658418063	100	100	5000	50	2020-05-07 10:29:37	cf1c20ac-cd4a-45d6-9f07-5aefb184a5f8	khr	1
2020050811223180215481	10000	10	0	0	2020-05-08 10:22:32	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	usd	1
2020050811225407234823	100	100	0	0	2020-05-08 10:22:55	f226c8e8-d4c1-4a79-9e3b-ada0eadf088d	khr	1
\.


--
-- Data for Name: servicer_terminal; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.servicer_terminal (terminal_no, servicer_no, terminal_number, pos_sn, is_delete, use_status) FROM stdin;
2020042419404232320123	dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	111111	1111	1	1
2020042419405060743903	dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	111111	1111	0	1
2020042419405046197048	dd6e1d2d-9aae-42d1-8ebd-a760bfcadd1f	2222	222222	0	1
2020042812024114125703	60f15170-c1db-41b0-bb3d-14185ab43d28	19CTCA884616555	19CTCA884616555	0	1
\.


--
-- Data for Name: settle_servicer_hourly; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.settle_servicer_hourly (log_no, start_time, finish_time, run_status, begin_time, end_time, sum_income_usd, sum_outgo_usd, balance_usd, delta_amount_usd, balance_khr, delta_amount_khr, sum_income_khr, sum_outgo_khr, fbalance_usd, fbalance_khr) FROM stdin;
\.


--
-- Data for Name: settle_vaccount_balance_hourly; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.settle_vaccount_balance_hourly (log_no, vaccount_no, balance, frozen_balance, create_time) FROM stdin;
\.


--
-- Data for Name: sms_send_record; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sms_send_record (id, msgid, account, business, mobile, msg, status, created_at) FROM stdin;
fb45e28e-705b-4eb7-abc1-5f07b3e4a017	1190422169923096576	I6641122	cl	855886716519	【Modern Pay】您的验证码是：674481，请勿泄露。	0	2020-04-24 07:47:02.824267
9af11ce5-fd7a-4ee8-b9be-ca58c3fecc90	1190442261624786944	I6641122	cl	855976355466	【Modern Pay】您的验证码是：003668，请勿泄露。	0	2020-04-24 16:06:53.06194
e27c5eee-c7de-4ef3-a997-328b151d1551	1190442751058120704	I6641122	cl	8550976355466	【Modern Pay】您的验证码是：967197，请勿泄露。	0	2020-04-24 16:08:49.75097
163fed02-2a59-4b91-92c4-fd2936cd3760	1190699093224198144	I6641122	cl	855085921025	【Modern Pay】您的验证码是：776247，请勿泄露。	0	2020-04-25 09:07:26.485493
6bce5058-653a-4dc3-8016-d1a1727f3a33	1190699582661726208	I6641122	cl	855085921025	【Modern Pay】您的验证码是：975645，请勿泄露。	0	2020-04-25 09:09:23.176198
f84ec674-8e7c-476b-a4cc-06f699607e50		I6641122	cl	855566595		1	2020-04-26 17:00:52.038839
9ee1996a-04ac-4169-8900-260d9506f2dd		I6641122	cl	855566595		1	2020-04-26 17:00:52.291222
5dd070d9-30fe-4482-b244-0369a25e9aa2		I6641122	cl	855566595		1	2020-04-26 17:00:52.582264
00770c3a-fcf1-43a4-9a96-68219e9ec669		I6641122	cl	855566595		1	2020-04-26 17:00:52.952875
ade64f6f-f9e9-4a5b-8762-7ea2ebeed30e		I6641122	cl	855566595		1	2020-04-26 17:00:53.57719
ee3d768b-c505-4e68-ac94-b9b9279d0563	1191513553639837696	I6641122	cl	855150123456	【Modern Pay】您的验证码是：413710，请勿泄露。	0	2020-04-27 15:03:48.979979
28bbd4f4-6e2a-4956-b603-c9b0a293718f	1191516649707278336	I6641122	cl	8551561234567	【Modern Pay】您的验证码是：494046，请勿泄露。	0	2020-04-27 15:16:07.145036
6d94daff-1c52-46f0-8cde-bd91689e6abe	1191855918778945536	I6641122	cl	85567567423	【Modern Pay】您的验证码是：636087，请勿泄露。	0	2020-04-28 13:44:15.192496
84a13c93-1450-4425-92a1-0f5bf56333b9	1191866336402477056	I6641122	cl	85567567423	【Modern Pay】您的验证码是：090426，请勿泄露。	0	2020-04-28 14:25:38.946091
c211fd22-8226-45fc-9f8e-e391fae4fe46	1191866823575080960	I6641122	cl	85567567423	【Modern Pay】您的验证码是：032810，请勿泄露。	0	2020-04-28 14:27:35.1022
ed36c70e-26f0-4ca0-ac7f-829967aad17f		I6641122	cl	8550976355466		1	2020-04-30 10:57:56.858032
eb763cef-71d1-49ca-9844-49bc37a5a916	1192539226814681088	I6641122	cl	855886716519	【Modern Pay】您的验证码是：794705，请勿泄露。	0	2020-04-30 10:59:28.527893
641d1e42-14b8-4f71-9aab-9eaefdf9a2a8		I6641122	cl	855886716519		1	2020-04-30 11:36:29.874767
71050f10-b8f8-472f-be38-1e5b65b0e711		I6641122	cl	855886716519		1	2020-04-30 11:42:12.277647
f19a4d5b-6f0e-437b-aed3-115172b0bd42		I6641122	cl	855886716519		1	2020-04-30 11:42:12.529548
4040e7ae-99cc-4fe0-9387-00590f67d2b4		I6641122	cl	855886716519		1	2020-04-30 11:42:13.528235
339da4b0-a117-4692-90f4-e1d2ec6f89b9	1192550324729155584	I6641122	cl	855886716519	【Modern Pay】您的验证码是：717384，请勿泄露。	0	2020-04-30 11:43:34.478348
08834565-abc2-4f16-8d2a-fa1d3089b667	1192550770927734784	I6641122	cl	855886716519	【Modern Pay】您的验证码是：575537，请勿泄露。	0	2020-04-30 11:45:20.860134
9a9cead6-7d52-4ccb-8b46-bc5ea801f5f6	1192564040208224256	I6641122	cl	855150123456	【Modern Pay】您的验证码是：562789，请勿泄露。	0	2020-04-30 12:38:04.502471
191070b6-d122-4732-a845-5da128b38a74	1192564050471555072	I6641122	cl	855150123456	【Modern Pay】您的验证码是：937524，请勿泄露。	0	2020-04-30 12:38:06.949792
b65ad53f-5310-4ac8-b4ba-0dbd4a46522e	1192573270017445888	I6641122	cl	855886716519	【Modern Pay】您的验证码是：521855，请勿泄露。	0	2020-04-30 13:14:45.0599
b553ece6-f19f-4e17-8b02-80d58277b564	1192573535143464960	I6641122	cl	855886716519	【Modern Pay】您的验证码是：171250，请勿泄露。	0	2020-04-30 13:15:48.272901
2d6a000f-98ae-4c26-9e3d-20f1f0963a6a		I6641122	cl	855886716519		1	2020-04-30 13:20:42.824105
25c436da-c7fd-48bb-bb07-c558f9f8a87d		I6641122	cl	855886716519		1	2020-04-30 13:20:49.098985
c2ee29e6-fffa-4b43-851e-a9596db194f6		I6641122	cl	855886716519		1	2020-04-30 13:45:16.270945
58a16956-6f1f-4bf3-9f77-43801a0b985e		I6641122	cl	855886716519		1	2020-04-30 13:45:18.819871
bf83c752-cf19-445f-99a5-28f2f05c03a6		I6641122	cl	855150123456		1	2020-04-30 13:46:21.569707
44c6a761-b6ba-4456-884d-dc93bf285fd1		I6641122	cl	855150123456		1	2020-04-30 13:46:22.994895
9f63e4c7-71c8-41d4-aff9-dd2d5eca737f	1192581925055893504	I6641122	cl	855886716519	【Modern Pay】您的验证码是：387623，请勿泄露。	0	2020-04-30 13:49:08.579686
988747ec-c1a7-4452-aa46-228f89429c25		I6641122	cl	855886716519		1	2020-04-30 13:50:46.713549
0f27ca3e-b8e9-426b-958b-57440b756042	1192583044540796928	I6641122	cl	855886716519	【Modern Pay】您的验证码是：367174，请勿泄露。	0	2020-04-30 13:53:35.485819
9af57a0c-32f1-41ac-94ce-1eb15525cc0e	1192584055389032448	I6641122	cl	855886716519	【Modern Pay】您的验证码是：732816，请勿泄露。	0	2020-04-30 13:57:37.048542
a4bd975f-8bb6-430c-87ce-00065c93547e	1192584324118089728	I6641122	cl	85536666555	【Modern Pay】您的验证码是：084530，请勿泄露。	0	2020-04-30 13:58:41.118918
0803b117-d6c0-4d3d-a1b2-7215494e3c70	1192591345949216768	I6641122	cl	855150123456	【Modern Pay】您的验证码是：891104，请勿泄露。	0	2020-04-30 14:26:34.70741
731f4c64-4753-4385-a9a7-b682e1e3ec39	1192591360134483968	I6641122	cl	855150123456	【Modern Pay】您的验证码是：360106，请勿泄露。	0	2020-04-30 14:26:38.076343
40b440c5-39b7-4ad1-b5e9-d8f6177eb654	1192592883891441664	I6641122	cl	855886716519	【Modern Pay】您的验证码是：136449，请勿泄露。	0	2020-04-30 14:32:41.368532
75880dd8-a5cc-4e8d-8dbb-d4f1bf63b15a	1192592889423728640	I6641122	cl	855886716519	【Modern Pay】您的验证码是：130958，请勿泄露。	0	2020-04-30 14:32:42.687132
258a5702-06c6-4a50-8bd4-e251f16dbe3e	1192592890552127488	I6641122	cl	855886716519	【Modern Pay】您的验证码是：892408，请勿泄露。	0	2020-04-30 14:32:42.955837
41f06a37-acc9-4231-ad74-ee05f3ff165e	1192595695291600896	I6641122	cl	855886716519	【Modern Pay】您的验证码是：192043，请勿泄露。	0	2020-04-30 14:43:51.65572
b2d3918a-7bc1-41ff-97f8-ba72420f2e52	1194743115546759168	I6641122	cl	85513425882336	【Modern Pay】您的验证码是：398223，请勿泄露。	0	2020-05-06 12:56:56.545418
79583741-4a44-4034-8866-4d8e265f90f3	1194770107461668864	I6641122	cl	855150123456	【Modern Pay】您的验证码是：140516，请勿泄露。	0	2020-05-06 14:44:11.920598
b7778059-6e5c-4938-ba57-b7ecd4e85ad4	1194770810330550272	I6641122	cl	8551551234567	【Modern Pay】您的验证码是：865261，请勿泄露。	0	2020-05-06 14:46:59.499135
df90d372-4f17-44d9-b1ef-383dc9717979	1194786892500176896	I6641122	cl	855233	【Modern Pay】您的验证码是：835810，请勿泄露。	0	2020-05-06 15:50:53.787392
7d7cb59e-6623-46c3-847e-626ee4040a5c	1194787296260657152	I6641122	cl	855233	【Modern Pay】您的验证码是：414231，请勿泄露。	0	2020-05-06 15:52:30.050181
aac3da6a-162d-4b39-ae24-1fdbce2e776c	1195062225136783360	I6641122	cl	855886716519	【Modern Pay】您的验证码是：440751，请勿泄露。	0	2020-05-07 10:04:58.201391
d6ce2b9a-b246-40c3-b082-661c06189731	1195062439079841792	I6641122	cl	855886716519	【Modern Pay】您的验证码是：629572，请勿泄露。	0	2020-05-07 10:05:49.209395
0570e4ff-5cc8-4370-8ec6-5c992855e628	1195409105712451584	I6641122	cl	855886716519	【Modern Pay】您的验证码是：442692，请勿泄露。	0	2020-05-08 09:03:20.975605
dfdd36ed-a748-40f5-aa3a-ec256494f2d6	1195413616170635264	I6641122	cl	855233	【Modern Pay】您的验证码是：673278，请勿泄露。	0	2020-05-08 09:21:16.353488
b46edecb-e533-4df9-a37a-d95529585892	1195415095606513664	I6641122	cl	855150123456	【Modern Pay】您的验证码是：541895，请勿泄露。	0	2020-05-08 09:27:09.077527
36f8fcfb-b16e-4f7d-9788-86a362567aaf	1195416275749310464	I6641122	cl	855886716519	【Modern Pay】您的验证码是：325868，请勿泄露。	0	2020-05-08 09:31:50.445375
f7054fd8-a967-4465-a91f-8f5a0129446c	1195417016111075328	I6641122	cl	855886716519	【Modern Pay】您的验证码是：144909，请勿泄露。	0	2020-05-08 09:34:46.96077
7707d0e0-7760-404c-b7cb-58f6ce9caad6	1195418475393323008	I6641122	cl	855000000000	【Modern Pay】您的验证码是：798962，请勿泄露。	0	2020-05-08 09:40:34.881141
a7d46ea4-ea56-45bf-b3ab-99bea52bbff6	1195418732902748160	I6641122	cl	85597983564986	【Modern Pay】您的验证码是：956010，请勿泄露。	0	2020-05-08 09:41:36.277296
2a3204e7-e799-42f8-befe-a1e1fc21bd5f	1195426966476689408	I6641122	cl	855886716519	【Modern Pay】您的验证码是：431777，请勿泄露。	0	2020-05-08 10:14:19.312802
585d124f-92f3-4dad-9d61-50587e564d20	1195428563546017792	I6641122	cl	855886716519	【Modern Pay】您的验证码是：672837，请勿泄露。	0	2020-05-08 10:20:40.0842
e6356480-62f6-4185-88d2-7e934b00970c	1195476302875070464	I6641122	cl	855233	【Modern Pay】您的验证码是：734179，请勿泄露。	0	2020-05-08 13:30:22.025779
fa78be27-89cd-4366-8025-55c0208e9117	1195476489186054144	I6641122	cl	85518148769649	【Modern Pay】您的验证码是：674868，请勿泄露。	0	2020-05-08 13:31:06.445203
\.


--
-- Data for Name: transfer_order; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transfer_order (log_no, from_vaccount_no, to_vaccount_no, amount, create_time, finish_time, order_status, balance_type, exchange_type, fees, payment_type, is_count, modify_time, ree_rate, real_amount, ip, lat, lng) FROM stdin;
2020042716420018848193	6fdd9297-54f6-4c9f-b158-170abd0ae07c	ea5c1573-69a7-44df-93ca-87f2298ffb9d	1000	2020-04-27 15:42:00.442053	2020-04-27 15:42:00.442053	3	usd	1	111	2	1	2020-04-27 15:42:00.490213	1110	1000	\N	\N	\N
2020042900003201837845	ab714dd1-1384-424e-9aa3-85b7903676da	c9e9b725-31f5-44eb-8e4f-55e5bc6d69cf	100	2020-04-28 23:00:32.272726	2020-04-28 23:00:32.272726	3	usd	1	100	2	1	2020-04-28 23:00:32.322929	1110	100	113.66.218.253		
2020043011580619577052	ab714dd1-1384-424e-9aa3-85b7903676da	c9e9b725-31f5-44eb-8e4f-55e5bc6d69cf	100	2020-04-30 10:58:06.175771	2020-04-30 10:58:06.175771	3	usd	1	100	2	1	2020-04-30 10:58:06.223675	300	100	61.140.127.170	0	0
\.


--
-- Data for Name: url; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.url (url_uid, url_name, url, parent_uid, title, icon, component_name, component_path, redirect, idx, is_hidden, create_time) FROM stdin;
222d3e77-c07d-4c30-a384-b0c493a3cb38	账号与菜单	/account_manage	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	账号与菜单	el-icon-s-custom	Content	account_manage	账号列表	1	0	1970-01-01 00:00:00
04032947-d9b3-4949-ad03-f6d9640a627a	简易商户修改	/easy_merc_edit	adcbffee-495f-4e59-90c0-f438d39ebd07	简易商户修改	el-icon-date	EasyMercEdit			24	1	1970-01-01 00:00:00
5e6cb690-6a8d-4fab-9234-754db94e1c2f	商户查看	/merc_check	adcbffee-495f-4e59-90c0-f438d39ebd07	商户查看	el-icon-date	MercCheck			17	1	1970-01-01 00:00:00
42fc26ac-7638-4b3b-9ed3-4b0a3f964512	门店管理	/store_manage	7e62c826-7285-4b26-9235-d4a49b40ef6c	门店管理	el-icon-menu	Content	store_manage	门店配置	2	0	1970-01-01 00:00:00
33bea49e-3565-469c-aca3-b524ed7232dd	业务员管理	/salesman_manage	7e62c826-7285-4b26-9235-d4a49b40ef6c	业务员管理	el-icon-menu	Content		业务员列表	7	0	1970-01-01 00:00:00
1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道管理	/aisle	7e62c826-7285-4b26-9235-d4a49b40ef6c	通道管理	el-icon-menu	Content	aisle	通道列表	3	0	1970-01-01 00:00:00
8ceae246-d9fc-4fdf-8487-c4abc76c26af	下属详情	/subordinate_details	adcbffee-495f-4e59-90c0-f438d39ebd07	下属详情	el-icon-menu	SubordinateDetails			19	1	1970-01-01 00:00:00
adcbffee-495f-4e59-90c0-f438d39ebd07	机构管理	/cooperation_agency	7e62c826-7285-4b26-9235-d4a49b40ef6c	机构管理	el-icon-menu	Content	cooperation_agency	合作列表	1	0	1970-01-01 00:00:00
eb9ea9f9-cd44-4b11-a805-061612ff12ef	角色修改	/role_edit	222d3e77-c07d-4c30-a384-b0c493a3cb38	角色修改	el-icon-date	RoleEdit			5	1	1970-01-01 00:00:00
5af52e1b-b6a1-4eca-87ec-466abacb4302	角色列表	/role_list	222d3e77-c07d-4c30-a384-b0c493a3cb38	角色列表	el-icon-date	RoleList			4	0	1970-01-01 00:00:00
082a3f2b-d30f-4463-af60-88ef60c0aca0	税务管理	/tax_manage	b8b970f8-84fb-450e-91f9-f98691152b83	税务管理	el-icon-menu	Content	tax_manage	税务信息	0	0	1970-01-01 00:00:00
6419a0a5-1e69-4775-8e72-04ee288635ff	日志管理	/log_manage	3af18777-25a3-49aa-9e90-04a45d71a6b9	日志管理	el-icon-menu	Content	log_manage	请求日志列表	1	0	1970-01-01 00:00:00
e4d8d6fb-c6c3-46be-8c71-e041d11fccf5	财务管理	/financial_manage	b8b970f8-84fb-450e-91f9-f98691152b83	财务管理	el-icon-menu	Content		财务报表	0	0	1970-01-01 00:00:00
afe5da52-c330-4f13-b7a6-386cafd6fd0d	简易商户查看	/easy_merc_check	adcbffee-495f-4e59-90c0-f438d39ebd07	简易商户查看	el-icon-date	EasyMercCheck			26	1	1970-01-01 00:00:00
a62b7765-f024-4705-877a-de0ad4878d19	商户池	/merchant	7e62c826-7285-4b26-9235-d4a49b40ef6c	商户池	el-icon-menu	Content		关联商户	10	0	1970-01-01 00:00:00
55327a19-fd5c-4101-885c-ec17508e5e74	商品管理	/goods_manage	7e62c826-7285-4b26-9235-d4a49b40ef6c	商品管理	el-icon-menu	Content		商品列表	8	0	1970-01-01 00:00:00
6330eaba-e5de-47f1-8c50-072120823683	合作修改(代理)	/agent_cooperation_agency_edit	adcbffee-495f-4e59-90c0-f438d39ebd07	合作修改(代理)	el-icon-date	AgentCooperationAgencyEdit			20	1	1970-01-01 00:00:00
1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	权限管理	/auth_manage	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	权限管理	el-icon-menu	Content	auth_list	权限列表	2	0	1970-01-01 00:00:00
ed06f369-9d92-4d07-8b9a-6488694b4d1c	我的使用	/my_use	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	我的使用	el-icon-menu	Content	my_use	公司信息	4	0	1970-01-01 00:00:00
797cecc2-d1f9-499a-b79d-e1f7e77c7874	统计	/statistics	15668ec7-6731-4646-a516-41035ddebe0d	统计	el-icon-menu	Content		交易分析	7	0	1970-01-01 00:00:00
a8de69fd-5856-42cc-8812-fef457a99079	交易管理	/order_manage	136f2190-6f80-46d4-a124-d512e508be5e	订单管理	el-icon-menu	Content		交易明细	1	0	1970-01-01 00:00:00
9b1d0f10-1be2-417c-a79f-a77e9dfa0c93	简易商户审核	/easy_merc_audit	adcbffee-495f-4e59-90c0-f438d39ebd07	简易商户审核	el-icon-date	EasyMercAudit			27	1	1970-01-01 00:00:00
0dc8b5f7-da88-4bdc-ba24-b5a6bc8a1044	消息修改	/message_edit	ed06f369-9d92-4d07-8b9a-6488694b4d1c	消息修改	reordrer	MessageEdit			5	1	1970-01-01 00:00:00
d1f06d6a-a377-4419-b690-1cbe09879b1c	终端发布	/terminal_publish	7e62c826-7285-4b26-9235-d4a49b40ef6c	终端发布	el-icon-menu	Content	endpoint_config_list	终端配置列表	5	0	1970-01-01 00:00:00
20ed729f-8a6b-47f0-8b79-6cd17dfa5260	定时任务	/regular_task	3af18777-25a3-49aa-9e90-04a45d71a6b9	定时任务	el-icon-menu	Content		定时任务列表	3	0	1970-01-01 00:00:00
14aa553e-5be3-404d-adb0-b4440e023c69	参数管理	/param_manage	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	参数管理	el-icon-menu	Content			3	0	1970-01-01 00:00:00
6f775fcf-6b53-49a6-960f-5be9ea2c32f7	扩展	/expansion	7e62c826-7285-4b26-9235-d4a49b40ef6c	扩展	el-icon-menu	Content	loan_application	借款申请	13	0	1970-01-01 00:00:00
b639edba-fe15-48c6-9e53-551c056a4426	第三方参数	/third_party_params	7e62c826-7285-4b26-9235-d4a49b40ef6c	第三方参数	el-icon-menu	Content	rate_list	汇卡费率配置	11	0	1970-01-01 00:00:00
a59a1727-baa8-472a-ab4c-fe0f185fbc5e	终端管理	/terminal_manage	7e62c826-7285-4b26-9235-d4a49b40ef6c	终端管理	el-icon-menu	Content	operator_list	收银员列表	4	0	1970-01-01 00:00:00
90b65d06-171a-4975-afd6-d226f3a35910	统计任务	/statistical_work	3af18777-25a3-49aa-9e90-04a45d71a6b9	统计任务	el-icon-menu	Content			4	0	1970-01-01 00:00:00
1897bc12-5b09-412d-a3c7-e001aedff81d	退款测试	/test_refund	065c9534-8c0a-408f-892a-914d172f90c7	退款测试	el-icon-date	TestRefund			4	0	1970-01-01 00:00:00
2feb45eb-2e97-4d7f-bf43-d221b1c71928	工具	/tools	7e62c826-7285-4b26-9235-d4a49b40ef6c	工具	el-icon-menu	Content		rsa_p12解码	12	0	1970-01-01 00:00:00
d6344f8d-4b29-48ce-b9e0-0c255447dc33	支付成功率(实时)	/merc_pay_analysis_online	797cecc2-d1f9-499a-b79d-e1f7e77c7874	支付成功率(实时)	el-icon-date	MercPayAnalysisOnline			11	0	1970-01-01 00:00:00
61ab5266-4740-49f3-9649-1b97566a4191	权限列表	/auth_list	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	权限列表	el-icon-date	AuthList			1	1	1970-01-01 00:00:00
b8b970f8-84fb-450e-91f9-f98691152b83	财务中心	/financial_center	00000000-0000-0000-0000-000000000000	财务中心	el-icon-menu	Default	financial_center	财务管理	10	1	1970-01-01 00:00:00
cb332a03-aeee-435c-86de-b7d577e3817b	日结算管理	/financial_manage	136f2190-6f80-46d4-a124-d512e508be5e	日结算管理	el-icon-menu	Content			2	1	1970-01-01 00:00:00
136f2190-6f80-46d4-a124-d512e508be5e	交易中心	/trade_center	00000000-0000-0000-0000-000000000000	交易中心	el-icon-menu	Default	trade_center	交易管理	5	0	1970-01-01 00:00:00
b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	账号中心	/account_center	00000000-0000-0000-0000-000000000000	账号中心	el-icon-menu	Default	account_center	账号管理	4	0	1970-01-01 00:00:00
3af18777-25a3-49aa-9e90-04a45d71a6b9	数据中心	/data_center	00000000-0000-0000-0000-000000000000	数据中心	el-icon-menu	Default	data_center	日志管理	6	1	1970-01-01 00:00:00
15668ec7-6731-4646-a516-41035ddebe0d	统计中心	/statistics_center	00000000-0000-0000-0000-000000000000	统计中心	el-icon-menu	Default	statistics_center	统计	3	1	1970-01-01 00:00:00
7e62c826-7285-4b26-9235-d4a49b40ef6c	运营中心	/operations_center	00000000-0000-0000-0000-000000000000	运营中心	el-icon-menu	Default	operations_center	我的使用	2	1	1970-01-01 00:00:00
065c9534-8c0a-408f-892a-914d172f90c7	支付测试	/pay_test	136f2190-6f80-46d4-a124-d512e508be5e	支付测试	el-icon-menu	Content			3	1	1970-01-01 00:00:00
b4aee72a-2ac5-4748-b846-bd1828139858	月结算管理	/settlement	136f2190-6f80-46d4-a124-d512e508be5e	月结算管理	el-icon-menu	Content		月结算	3	1	1970-01-01 00:00:00
5a229d4e-beee-4f0a-ba2f-888f04a149e5	风控管理	/risk_control_manage	136f2190-6f80-46d4-a124-d512e508be5e	风控管理	el-icon-menu	Content			5	1	1970-01-01 00:00:00
f35b0f93-775b-443a-b3fe-07a9a100d1a0	财务报表	/financial_report	e4d8d6fb-c6c3-46be-8c71-e041d11fccf5	财务报表	el-icon-date	WorkBench			0	0	1970-01-01 00:00:00
296464e5-f911-457a-be84-7f85532b0c5b	税务信息	/tax_info	082a3f2b-d30f-4463-af60-88ef60c0aca0	税务信息	el-icon-date	WorkBench			0	0	1970-01-01 00:00:00
ea351efd-4942-4912-ac47-51545913d554	开票记录	/tax_log	082a3f2b-d30f-4463-af60-88ef60c0aca0	开票记录	el-icon-date	WorkBench			0	0	1970-01-01 00:00:00
233d5915-ed17-41fa-8a13-b8e02f4c510f	冲正记录	/reverse	082a3f2b-d30f-4463-af60-88ef60c0aca0	冲正记录	el-icon-date	WorkBench			0	0	1970-01-01 00:00:00
9392bd8f-6ef5-4f16-b87b-a36709339d5d	报表管理	/report_form	3af18777-25a3-49aa-9e90-04a45d71a6b9	报表管理	el-icon-menu	Content		报表列表	4	0	1970-01-01 00:00:00
a317c13a-414a-4f44-9b46-3c9fe0b735ee	科目管理	/subject_manage	274aca4a-db46-450c-99f7-8a6f44f1a6d0	科目管理	el-icon-menu	Content			1	0	1970-01-01 00:00:00
6d40ba44-ccbf-4491-bdd3-ebb3a4ad3cce	商品税收编码	/tax_code	082a3f2b-d30f-4463-af60-88ef60c0aca0	商品税收编码	el-icon-date	WorkBench			0	0	1970-01-01 00:00:00
4badf803-d758-47a4-8a4b-fb005cf7167a	角色授权	/role_auth	222d3e77-c07d-4c30-a384-b0c493a3cb38	角色授权	el-icon-date	RoleAuth			6	1	1970-01-01 00:00:00
59cfe00e-f2ad-4c84-965b-6cee95622718	退款记录	/refund_record	a051904d-9315-48c3-9ecd-1fcedaaa38e9	退款记录	el-icon-date	RefundRecord			0	1	1970-01-01 00:00:00
5eadf649-5e8a-4172-9857-ef565a5414ed	日结账单	/dayorder	a051904d-9315-48c3-9ecd-1fcedaaa38e9	日结账单	el-icon-date	DayOrder			0	1	1970-01-01 00:00:00
779dc403-d84a-4e4e-9aa8-3d7cbf013a69	请求日志列表	/log_list	6419a0a5-1e69-4775-8e72-04ee288635ff	请求日志列表	el-icon-date	LogList			0	0	1970-01-01 00:00:00
9e7d73a9-49c9-4dc1-ac35-897e0040ca51	商户编辑	/submerchant_edit	5f224689-4f42-40fc-88bb-df7ad04d0e14	商户编辑	el-icon-date	SubmerchantEdit			1	1	1970-01-01 00:00:00
ccc79a5f-1c03-44c8-a954-2104bbdb7239	后台日志列表	/backstage_list	6419a0a5-1e69-4775-8e72-04ee288635ff	后台日志列表	el-icon-date	BackstageList			0	0	1970-01-01 00:00:00
805c95bb-9f47-4e74-a1f0-87c3bb12a455	业务员修改	/salesman_edit	33bea49e-3565-469c-aca3-b524ed7232dd	业务员修改	el-icon-date	SalesmanEdit			2	1	1970-01-01 00:00:00
58d6c0a6-e80d-4d79-972f-0ad5f8517e97	业务员列表	/salesman_list	33bea49e-3565-469c-aca3-b524ed7232dd	业务员列表	el-icon-date	SalesmanList			1	0	1970-01-01 00:00:00
72793652-a596-4b80-bbe4-73d10cca1379	黑名单修改	/blacklist_edit	5a229d4e-beee-4f0a-ba2f-888f04a149e5	黑名单修改	el-icon-date	BlacklistEdit			3	1	1970-01-01 00:00:00
5a8b13c6-3997-4846-b037-2b679f5c31e8	风控管理	/risk_manage	5a229d4e-beee-4f0a-ba2f-888f04a149e5	风控管理	el-icon-date	RiskManage			10	0	1970-01-01 00:00:00
66cb5fb9-5bdc-4378-8010-6b58f8cedb5f	白名单修改	/whitelist_edit	5a229d4e-beee-4f0a-ba2f-888f04a149e5	白名单修改	el-icon-date	WhitelistEdit			4	1	1970-01-01 00:00:00
3d74d6a0-376a-4ff2-a1fc-5150b624e769	业务员查看	/salesman_check	33bea49e-3565-469c-aca3-b524ed7232dd	业务员查看	el-icon-date	SalesmanCheck			3	1	1970-01-01 00:00:00
fa13b0c6-1e65-45de-bef8-ac534342f0f3	菜单修改	/menu_edit	222d3e77-c07d-4c30-a384-b0c493a3cb38	菜单修改	el-icon-date	MenuEdit			21	1	1970-01-01 00:00:00
aded78fc-0d49-436e-a780-ebea7c8b0f4c	账号查看	/account_check	222d3e77-c07d-4c30-a384-b0c493a3cb38	账号查看	el-icon-date	AccountCheck			9	1	1970-01-01 00:00:00
14f49944-6337-4c1c-8ea4-dd111c7859f1	角色查看	/role_check	222d3e77-c07d-4c30-a384-b0c493a3cb38	角色查看	el-icon-date	RoleCheck			10	1	1970-01-01 00:00:00
05825db2-caf8-4951-9f2d-ea85ace66073	商品列表	/goods_list	55327a19-fd5c-4101-885c-ec17508e5e74	商品列表	el-icon-date	GoodsList			0	0	1970-01-01 00:00:00
44a5389b-8608-4865-ab53-74be807043a2	商品修改	/good_edit	55327a19-fd5c-4101-885c-ec17508e5e74	商品修改	el-icon-date	GoodEdit			1	1	1970-01-01 00:00:00
ba6556cb-4288-428d-ac87-a67546dd70be	公司信息	/company_info	ed06f369-9d92-4d07-8b9a-6488694b4d1c	公司信息	el-icon-date	CompanyInfo			2	0	1970-01-01 00:00:00
e72dd18e-5845-4ef4-934c-91b3f01eddde	分类列表	/class_list	55327a19-fd5c-4101-885c-ec17508e5e74	分类列表	el-icon-date	ClassList			2	0	1970-01-01 00:00:00
6c154b5d-4689-4282-a54e-dbbeeeba5fb5	销售订单	/sell_order	55327a19-fd5c-4101-885c-ec17508e5e74	销售订单	el-icon-date	SellOrder			3	0	1970-01-01 00:00:00
8ec3bc88-f2ce-498f-8d01-a40b02804f4c	商品详情	/good_detail	55327a19-fd5c-4101-885c-ec17508e5e74	商品详情	el-icon-date	GoodDetail			5	1	1970-01-01 00:00:00
a1951306-b89e-46d3-8768-3d681585e382	销售开单	/sell_make	55327a19-fd5c-4101-885c-ec17508e5e74	销售开单	el-icon-date	SellMake			6	0	1970-01-01 00:00:00
622186a1-17b8-4b64-8406-dd3cd1dd2596	账号修改	/account_edit	222d3e77-c07d-4c30-a384-b0c493a3cb38	账号修改	el-icon-date	AccountEdit			3	1	1970-01-01 00:00:00
98c83395-0888-4128-b6c9-f6daa055d5ac	报表修改	/report_form_edit	ae401632-6914-4d39-b3eb-398852202a41	报表修改	el-icon-date	ReportFormEdit			2	0	1970-01-01 00:00:00
192ddf18-05b1-4d80-8b2c-923cef89df16	销售订单详情	/sell_order_detail	55327a19-fd5c-4101-885c-ec17508e5e74	销售订单详情	el-icon-date	SellOrderDetail			4	1	1970-01-01 00:00:00
608b52b4-40a6-43f1-bca2-031b2ea2abf4	进件资料审核	/incoming_data_audit	adcbffee-495f-4e59-90c0-f438d39ebd07	进件资料审核	el-icon-date	IncomingDataAudit			6	1	1970-01-01 00:00:00
0ad6d6ee-1106-4940-b9ee-916b16fe6622	报表查看	/report_form_check	ae401632-6914-4d39-b3eb-398852202a41	报表查看	el-icon-date	ReportFormCheck			3	0	1970-01-01 00:00:00
6823c594-6616-47db-9eae-8263b0ec13f7	报表下载	/report_form_download	ae401632-6914-4d39-b3eb-398852202a41	报表下载	el-icon-date	ReportFormDownload			4	0	1970-01-01 00:00:00
6e4e2624-bc76-49dc-86ba-6e88b0977c79	交易分析	/transaction	797cecc2-d1f9-499a-b79d-e1f7e77c7874	交易分析	el-icon-date	Transaction			3	0	1970-01-01 00:00:00
f76dca7f-4774-4c76-8aac-873e97f66d3c	收银员列表	/operator_list	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	收银员列表	el-icon-date	OperatorList			1	0	1970-01-01 00:00:00
e3eaddf3-edb3-4930-a548-bce70c834b66	收银员修改	/operator_edit	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	收银员修改	el-icon-date	OperatorEdit			2	1	1970-01-01 00:00:00
9def3d8c-3f3d-4a1e-9020-6b382fc42310	接口组列表	/interface_group_list	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	接口组列表	el-icon-date	InterfaceGroupList			1	0	1970-01-01 00:00:00
590dd762-650e-4b31-9989-268ae81a853c	终端配置列表	/endpoint_config_list	d1f06d6a-a377-4419-b690-1cbe09879b1c	终端配置列表	el-icon-date	EndpointConfigList			1	0	1970-01-01 00:00:00
f50b81e9-1797-40e7-a91b-5ce8be353e8b	门店修改	/store_config_edit	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	门店修改	el-icon-date	StoreConfigEdit			6	1	1970-01-01 00:00:00
ae527d80-3b93-4ec1-a4b3-5c21ab2c80fc	门店列表	/store_list	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	门店列表	el-icon-date	StoreList			4	0	1970-01-01 00:00:00
d6791542-86e6-43a7-ae04-7d42c51889a5	接口组查看	/interface_group_check	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	接口组查看	el-icon-date	InterfaceGroupCheck			13	1	1970-01-01 00:00:00
7628af1d-b868-4923-8bac-a7366efdc813	合作列表(代理)	/agent_cooperation_agency_list	adcbffee-495f-4e59-90c0-f438d39ebd07	合作列表(代理)	el-icon-date	AgentCooperationAgencyList			21	0	1970-01-01 00:00:00
a0419fd4-acd0-4e2b-8156-30066469f33e	风控交易列表	/risk_control_transactions_list	adcbffee-495f-4e59-90c0-f438d39ebd07	风控交易列表	el-icon-date	RiskControlTransactionsList			22	1	1970-01-01 00:00:00
7872d920-bf48-491f-a604-f29239826d86	局部规则修改	/local_rule_edit	adcbffee-495f-4e59-90c0-f438d39ebd07	局部规则修改	el-icon-date	LocalRuleEdit			23	1	1970-01-01 00:00:00
1e8a19be-9e51-48af-afc0-d56772e66a24	oem设置	/oem_set	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	oem设置	el-icon-date	OemSet			8	0	1970-01-01 00:00:00
c2731296-80f6-493f-acd6-3739949e0ad0	接口组配置	/interface_group_config	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	接口组配置	el-icon-date	InterfaceGroupConfig			9	1	1970-01-01 00:00:00
b2776c17-fde6-4af5-8d34-f72d5ec88ac9	账号列表	/account_list	222d3e77-c07d-4c30-a384-b0c493a3cb38	账号列表	el-icon-date	AccountList			1	0	1970-01-01 00:00:00
e0be0b73-1f90-4c04-8da5-86513c148a16	充值与代付	/recharge_and_pay	136f2190-6f80-46d4-a124-d512e508be5e	充值与代付	el-icon-menu	Content			4	1	1970-01-01 00:00:00
7bd65052-1ff5-4870-abdf-e398ac9a5fc3	通道列表	/aisle_list	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道列表	el-icon-date	AisleList			2	0	1970-01-01 00:00:00
4652e35d-14d1-4d5e-8c34-dd691500ce69	通道修改	/aisle_edit	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道修改	el-icon-date	AisleEdit			6	1	1970-01-01 00:00:00
acde5628-57d0-495f-88ba-db5c6c079e87	通道添加	/channel_add	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道添加	el-icon-date	ChannelAdd			8	1	1970-01-01 00:00:00
5e5c4e3f-155c-48f4-9aa6-dc7dd88fa56f	套餐查看	/rate_package_check	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	套餐查看	el-icon-date	RatePackageCheck			14	1	1970-01-01 00:00:00
64cbf7b5-71e0-4179-a0a8-07020335901d	费率套餐	/rate_package	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	费率套餐	el-icon-date	RatePackage			4	0	1970-01-01 00:00:00
23ee37f8-d777-4fdc-9d2b-8705c72ce357	费率套餐修改	/rate_package_edit	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	费率套餐修改	el-icon-date	RatePackageEdit			5	1	1970-01-01 00:00:00
b8840d1d-3908-4db7-9825-b5ce78232409	商户池列表	/merchant_pool_list	a62b7765-f024-4705-877a-de0ad4878d19	商户池列表	el-icon-date	MerchantPoolList			0	0	1970-01-01 00:00:00
6b7a79ac-30dd-44ce-a6bb-82a137f6771b	商户修改	/merc_edit	adcbffee-495f-4e59-90c0-f438d39ebd07	商户修改	el-icon-date	MercEdit			8	1	1970-01-01 00:00:00
18e4c9a1-072c-4b67-83ac-8403efd66ff5	报表列表	/project_report_list	9392bd8f-6ef5-4f16-b87b-a36709339d5d	报表列表	el-icon-date	ProjectReportList			2	0	1970-01-01 00:00:00
6afbed32-4c43-43c7-b300-c1a04857d346	关联商户	/related_businesses	a62b7765-f024-4705-877a-de0ad4878d19	关联商户	el-icon-date	RelatedBusinesses			2	1	1970-01-01 00:00:00
6be56b39-079c-410a-b41a-35b32a6a5e1f	商户池查看	/merchant_pool_check	a62b7765-f024-4705-877a-de0ad4878d19	商户池查看	el-icon-date	MerchantPoolCheck			5	1	1970-01-01 00:00:00
4f387a47-6b08-48f6-981c-42a8a128d1f0	商户池修改	/merchant_pool_edit	a62b7765-f024-4705-877a-de0ad4878d19	商户池修改	el-icon-date	MerchantPoolEdit			0	1	1970-01-01 00:00:00
28c94a48-fd22-49a2-aa57-0f72c1f635d2	商户属性修改	/merc_attr_edit	a62b7765-f024-4705-877a-de0ad4878d19	商户属性修改	el-icon-date	MercAttrEdit			3	1	1970-01-01 00:00:00
46ea3bb2-76b6-4842-a015-33c9edc910ae	钱包	/wallet	ed06f369-9d92-4d07-8b9a-6488694b4d1c	钱包	el-icon-date	Wallet			3	0	1970-01-01 00:00:00
2697141a-6da7-4cba-8237-9cc42a279742	定时任务列表	/regular_task_list	20ed729f-8a6b-47f0-8b79-6cd17dfa5260	定时任务列表	el-icon-date	RegularTaskList			1	0	1970-01-01 00:00:00
83334d84-d119-4c76-9292-bb2a34731d9d	定时任务修改	/regular_task_edit	20ed729f-8a6b-47f0-8b79-6cd17dfa5260	定时任务修改	el-icon-date	RegularTaskEdit			2	1	1970-01-01 00:00:00
621bc69f-2dcf-4754-bcf8-71a18d17fa7e	消息列表	/message_list	ed06f369-9d92-4d07-8b9a-6488694b4d1c	消息列表	el-icon-date	MessageList			4	0	1970-01-01 00:00:00
bda3fd00-36a5-43bd-9c2d-0a25511217f3	消息查看	/message_check	ed06f369-9d92-4d07-8b9a-6488694b4d1c	消息查看	el-icon-date	MessageCheck			0	1	1970-01-01 00:00:00
0e143513-3b89-4e2d-9865-8a848ad6c8db	配置	/personal_config	ed06f369-9d92-4d07-8b9a-6488694b4d1c	配置	el-icon-date	PersonalConfig			6	0	1970-01-01 00:00:00
f7833e34-f3da-40a2-9bf1-a14cb550925d	当天统计	/statistics_for_day	797cecc2-d1f9-499a-b79d-e1f7e77c7874	当天统计	el-icon-date	StatisticsForDay			1	1	1970-01-01 00:00:00
115e00a4-dd87-4ebd-a471-9a69f5f3f85f	回调日志列表	/callback_list	6419a0a5-1e69-4775-8e72-04ee288635ff	回调日志列表	el-icon-date	CallbackList			0	1	1970-01-01 00:00:00
b3e49cac-fd53-4dcb-97ea-03ac3c328405	钱包密码设置	/wallet_password	ed06f369-9d92-4d07-8b9a-6488694b4d1c	钱包密码设置	el-icon-date	WalletPassword			6	1	1970-01-01 00:00:00
29058c58-14a6-425d-b366-617860f6950f	对账明细	/download_files	cb332a03-aeee-435c-86de-b7d577e3817b	对账明细	el-icon-date	DownloadFiles			1	0	1970-01-01 00:00:00
0f13980e-ffc3-4712-9c58-1726e1ec6787	业务经理列表	/sales_manager_list	33bea49e-3565-469c-aca3-b524ed7232dd	业务经理列表	el-icon-date	SalesmanagerList			0	0	1970-01-01 00:00:00
761f6268-1ccb-4a23-b43d-cad43d4761d2	钱包列表	/wallet_list	ed06f369-9d92-4d07-8b9a-6488694b4d1c	钱包列表	el-icon-date	WalletList			0	0	1970-01-01 00:00:00
8b470f68-065a-4b9a-a35e-0781d83910df	商户通道钱包	/merc_channel_wallet	ed06f369-9d92-4d07-8b9a-6488694b4d1c	商户通道钱包	el-icon-date	MercChannelWallet			10	0	1970-01-01 00:00:00
7270ed76-fcee-4a5a-a37f-6645f69e68aa	钱包日志	/wallet_log	ed06f369-9d92-4d07-8b9a-6488694b4d1c	钱包日志	el-icon-date	WalletLog			3	1	1970-01-01 00:00:00
cc9a6c3a-518a-458a-a84e-6353a1c22165	钱包冻结日志	/wallet_freeze_log	ed06f369-9d92-4d07-8b9a-6488694b4d1c	钱包冻结日志	el-icon-date	WalletFreezeLog			0	1	1970-01-01 00:00:00
e1ea71e4-fcbd-4ca5-8917-1dcab59fbd45	业务经理修改	/sales_manager_edit	33bea49e-3565-469c-aca3-b524ed7232dd	业务经理修改	el-icon-date	SalesmanagerEdit			0	1	1970-01-01 00:00:00
86f767d3-aec7-4acc-a3b3-e63b1734943f	用户分析	/user_analysis	797cecc2-d1f9-499a-b79d-e1f7e77c7874	用户分析	el-icon-date	UserAnalysis			4	1	1970-01-01 00:00:00
0150f1aa-0dfe-4371-a553-87e73d886b3f	终端列表	/terminal_list	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	终端列表	el-icon-date	TerminalList			3	0	1970-01-01 00:00:00
7ac101de-041e-4baa-ab7e-7471318ec81f	终端修改	/terminal_edit	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	终端修改	el-icon-date	TerminalEdit			4	1	1970-01-01 00:00:00
e23feaf4-58ae-469c-95a3-83e6312e7f29	第三方账号列表	/endpoint_config	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	第三方账号列表	el-icon-date	EndpointConfig			5	0	1970-01-01 00:00:00
deb75efd-1d47-4fb3-88d0-4fe07c5cbf36	终端配置修改	/endpoint_config_edit	d1f06d6a-a377-4419-b690-1cbe09879b1c	终端配置修改	el-icon-date	EndpointConfigEdit			2	1	1970-01-01 00:00:00
19b78054-0c63-453c-9ca7-bd64c99de999	终端版本列表	/version_list	d1f06d6a-a377-4419-b690-1cbe09879b1c	终端版本列表	el-icon-date	VersionList			3	0	1970-01-01 00:00:00
96deec3e-55a6-422c-9cb5-3f59ff7a3dd2	分支列表	/branch_list	d1f06d6a-a377-4419-b690-1cbe09879b1c	分支列表	el-icon-date	BranchList			5	0	1970-01-01 00:00:00
fa083891-8c7b-4d3e-a769-9c765760f79e	子商户列表	/sub_merc_list	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	商户列表	el-icon-date	SubMerchantList			1	0	1970-01-01 00:00:00
6b7a79ac-30dd-44ce-a6bb-82a137f67aaa	子商户修改	/sub_merc_edit	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	子商户修改	el-icon-date	SubMerchantEdit			3	1	1970-01-01 00:00:00
680460f5-af36-47e7-9278-cc316be1ec1f	通道组设置	/channel_group_set	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道组设置	el-icon-date	ChannelGroupSet			7	1	1970-01-01 00:00:00
5e6cb690-6a8d-4fab-9234-754db94e1caa	子商户查看	/sub_merc_check	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	子商户查看	el-icon-date	SubMerchantCheck			2	1	1970-01-01 00:00:00
6e4f00e4-bfad-42af-949d-8b215dae11f2	清分	/liquidation_data	ed06f369-9d92-4d07-8b9a-6488694b4d1c	清分	el-icon-date	LiquidationData			2	1	1970-01-01 00:00:00
43b61baf-0753-4367-9d29-7be6623d91c9	代理修改	/agency_edit	adcbffee-495f-4e59-90c0-f438d39ebd07	代理修改	el-icon-date	AgencyEdit			7	1	1970-01-01 00:00:00
8964e926-7f17-4370-bfe6-83cf363f5244	通道组查看	/channel_group_check	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道组查看	el-icon-date	ChannelGroupCheck			11	1	1970-01-01 00:00:00
a27e89ae-9bb8-4425-be64-76b5a811c137	通道查看	/channel_check	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道查看	el-icon-date	ChannelCheck			12	1	1970-01-01 00:00:00
6eee51ef-8a8f-4966-b19d-c9b915e681e4	套餐授予	/rate_package_grant	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	套餐授予	el-icon-date	RatePackageGrant			15	1	1970-01-01 00:00:00
e16cdfff-e9b4-42c2-b712-018adcfee5e2	通道组列表	/channel_group_list	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	通道组列表	el-icon-date	ChannelGroupList			3	0	1970-01-01 00:00:00
824f43c8-ee42-47ed-bbb9-5f8de09acdbc	api列表	/pay_api_list	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	api列表	el-icon-date	PayApiList			3	0	1970-01-01 00:00:00
90476856-94be-40ed-8b1e-814392074d81	api编辑	/pay_api_edit	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	api编辑	el-icon-date	PayApiEdit			4	1	1970-01-01 00:00:00
977f77aa-139c-4d4c-8640-053eea3376dd	api查看	/pay_api_check	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	api查看	el-icon-date	PayApiCheck			5	1	1970-01-01 00:00:00
3293247e-3d3c-4749-ae91-1a991e27b8be	风控规则修改	/wind_control_edit	5a229d4e-beee-4f0a-ba2f-888f04a149e5	风控规则修改	el-icon-date	WindControlEdit			5	1	1970-01-01 00:00:00
742c55df-02ca-4d3e-909a-ac8c412ff735	充值审核	/wallet_audit	e0be0b73-1f90-4c04-8da5-86513c148a16	充值审核	el-icon-date	WalletAudit			1	0	1970-01-01 00:00:00
4aacad3f-4671-4f8f-a6bc-73ce21233408	我的商户列表	/merc_configuration_list	adcbffee-495f-4e59-90c0-f438d39ebd07	我的商户列表	el-icon-date	MercConfigurationList			16	0	1970-01-01 00:00:00
f58d731d-1d2e-4380-ba29-4896951518e0	近七天交易分析	/transaction_week	797cecc2-d1f9-499a-b79d-e1f7e77c7874	近七天交易分析	el-icon-date	TransactionWeek			6	0	1970-01-01 00:00:00
c8b42420-731e-4736-97a0-81725e9aea49	近一月交易分析	/transaction_month	797cecc2-d1f9-499a-b79d-e1f7e77c7874	近一月交易分析	el-icon-date	TransactionMonth			7	0	1970-01-01 00:00:00
357d5bae-8470-4d03-ab92-68cafedcfa16	代付测试	/test_agentpay	065c9534-8c0a-408f-892a-914d172f90c7	代付测试	el-icon-date	TestAgentpay			3	0	1970-01-01 00:00:00
4cb4759c-7efc-4a7f-9b50-4d2fc1a00c56	即时对账	/immediate_reconcilia	cb332a03-aeee-435c-86de-b7d577e3817b	即时对账	el-icon-date	ImmediateReconcilia			3	0	1970-01-01 00:00:00
97a477c0-7361-421f-b042-2d1e0d9bf5fe	清算任务队列列表	/liquidation_queue_list	cb332a03-aeee-435c-86de-b7d577e3817b	清算任务队列列表	el-icon-date	LiquidationQueueList			3	1	1970-01-01 00:00:00
a2bcc31f-515a-45e9-8083-b2f27d2c4afa	查看清算任务详情	/liquidation_check	cb332a03-aeee-435c-86de-b7d577e3817b	查看清算任务详情	el-icon-date	LiquidationCheck			4	1	1970-01-01 00:00:00
60f259e8-d103-4fda-859d-b6c053cc6004	清算任务审核列表	/liquidation_audit	cb332a03-aeee-435c-86de-b7d577e3817b	清算任务审核列表	el-icon-date	LiquidationAudit			7	0	1970-01-01 00:00:00
46a9fbb8-a2d4-4bc0-aaf8-ded371f5f05d	清算任务详情列表	/liquidation_list	cb332a03-aeee-435c-86de-b7d577e3817b	清算任务详情列表	el-icon-date	LiquidationList			2	0	1970-01-01 00:00:00
f01d9cd6-4d98-43b0-8aef-c52a46964512	商户池	/merchants_pool	797cecc2-d1f9-499a-b79d-e1f7e77c7874	商户池	el-icon-date	MerchantsPool			8	0	1970-01-01 00:00:00
2c1d69a3-4af1-4436-9a70-449ba88f3048	通道交易统计	/channel_transaction_statistics	797cecc2-d1f9-499a-b79d-e1f7e77c7874	通道交易统计	el-icon-date	ChannelTransactionStatistics			8	0	1970-01-01 00:00:00
e41d0e51-6eb2-4622-8b59-f76ee6e6c1e1	总交易分析	/transaction_analysis	797cecc2-d1f9-499a-b79d-e1f7e77c7874	总交易分析	el-icon-date	TransactionAnalysis			9	0	1970-01-01 00:00:00
a168235d-63c0-4b58-8b49-3df20a9a815d	全局参数配置	/global_param_list	14aa553e-5be3-404d-adb0-b4440e023c69	全局参数配置	el-icon-date	GlobalParamList			1	0	1970-01-01 00:00:00
8bb9a2bc-aeea-435f-9132-73caead02416	全局参数配置查看	/global_param_check	14aa553e-5be3-404d-adb0-b4440e023c69	全局参数配置查看	el-icon-date	GlobalParamCheck			2	1	1970-01-01 00:00:00
c2592cf2-b4f8-4ab8-9ce3-54b4dbbc928a	全局参数配置修改	/global_param_edit	14aa553e-5be3-404d-adb0-b4440e023c69	全局参数配置修改	el-icon-date	GlobalParamEdit			3	1	1970-01-01 00:00:00
c8d5718b-6e24-43e9-9967-9a6318412dcb	统计任务列表	/statistical_work_list	90b65d06-171a-4975-afd6-d226f3a35910	统计任务列表	el-icon-date	StatisticalWorkList			1	0	1970-01-01 00:00:00
f7c21ada-b21e-4ce8-8a52-df88fa170b5d	统计任务修改	/statistical_work_edit	90b65d06-171a-4975-afd6-d226f3a35910	统计任务修改	el-icon-date	StatisticalWorkEdit			2	1	1970-01-01 00:00:00
fddf83c2-ea32-4346-8a80-b035758265ae	统计任务查看	/statistical_work_check	90b65d06-171a-4975-afd6-d226f3a35910	统计任务查看	el-icon-date	StatisticalWorkCheck			3	1	1970-01-01 00:00:00
3f3abbb0-4c6b-4b05-9e0f-4413b4adf347	上游资料提交	/upstream_files_list	adcbffee-495f-4e59-90c0-f438d39ebd07	上游资料提交	el-icon-date	UpstreamFilesList			15	1	1970-01-01 00:00:00
c5e2efa4-c8c7-4c79-9823-59e829a4847e	代付任务详情	/contrary_apply_check	e0be0b73-1f90-4c04-8da5-86513c148a16	代付任务详情	el-icon-date	ContraryApplyCheck			5	1	1970-01-01 00:00:00
ad55ef62-54a6-41c0-bba9-830a06822659	月结任务通道详情	/monthly_settlement_channel_list	b4aee72a-2ac5-4748-b846-bd1828139858	月结任务通道详情	el-icon-date	MonthlySettlementChannelList			5	1	1970-01-01 00:00:00
db98b40b-91a4-4434-ba2c-712e3990b82a	月结任务详情列表	/monthly_settlement_list	b4aee72a-2ac5-4748-b846-bd1828139858	月结任务详情列表	el-icon-date	MonthlySettlementList			4	1	1970-01-01 00:00:00
b995956c-75be-4635-9022-af1f5376be2b	月清算任务列表	/monthly_liquidation_list	b4aee72a-2ac5-4748-b846-bd1828139858	月清算任务列表	el-icon-date	MonthlyLiquidationList			6	0	1970-01-01 00:00:00
03ed060b-56dd-4ecb-9072-a736c81ae27f	月清算任务审核	/monthly_liquidation_audit	b4aee72a-2ac5-4748-b846-bd1828139858	月清算任务审核	el-icon-date	MonthlyLiquidationAudit			7	0	1970-01-01 00:00:00
15362ef8-4d86-4e96-b49f-c33c5292e67f	汇卡费率配置	/rate_list	b639edba-fe15-48c6-9e53-551c056a4426	汇卡费率配置	el-icon-date	RateList			3	0	1970-01-01 00:00:00
ff6e92a2-a7a4-4508-9d4e-bde9f6cb0bfe	贷款银行查看	/loan_bank_check	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	贷款银行查看	el-icon-date	LoanBankCheck			4	1	1970-01-01 00:00:00
4689bc27-a161-431b-bf58-456f66307fd9	主扫测试	/test_prepay	065c9534-8c0a-408f-892a-914d172f90c7	主扫测试	el-icon-date	TestPrepay			1	0	1970-01-01 00:00:00
ef82cdff-1b9d-4396-8d2b-f7b751a50006	被扫测试	/test_pay	065c9534-8c0a-408f-892a-914d172f90c7	被扫测试	el-icon-date	TestPay			2	0	1970-01-01 00:00:00
48c6f1fa-5349-4f16-9d45-fc4d1eb1fbf5	批量充值	/import_recharge_flow	e0be0b73-1f90-4c04-8da5-86513c148a16	批量充值	el-icon-date	ImportRechargeFlow			2	0	1970-01-01 00:00:00
90df4bc3-dbb0-4e92-ab9a-0aac895e5295	代付申请审核	/contrary_apply_audit	e0be0b73-1f90-4c04-8da5-86513c148a16	代付申请审核	el-icon-date	ContraryApplyAudit			3	0	1970-01-01 00:00:00
23476194-48d6-447c-b4e9-2ab6ad7b1a48	手续费审核	/fee_review	e0be0b73-1f90-4c04-8da5-86513c148a16	手续费审核	el-icon-date	FeeReview			24	0	1970-01-01 00:00:00
00f047a0-8afe-40a6-89d5-74867de6cfbe	第三方账号修改	/endpoint_config_add	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	第三方账号修改	el-icon-date	EndpointConfigAdd			6	1	1970-01-01 00:00:00
0ad5ba51-e406-477c-bdc2-fe35ca35b6ec	日结算	/day_settlement	797cecc2-d1f9-499a-b79d-e1f7e77c7874	日结算	el-icon-date	DaySettlement			2	0	1970-01-01 00:00:00
665c7169-f989-436d-a611-e6effc8ddc9b	第三方商户列表	/third_party_merc_list	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	第三方商户列表	el-icon-date	ThirdPartyMercList			7	0	1970-01-01 00:00:00
029a70de-2382-485b-9eb3-15d3933461d6	终端版本修改	/version_edit	d1f06d6a-a377-4419-b690-1cbe09879b1c	终端版本修改	el-icon-date	VersionEdit			4	1	1970-01-01 00:00:00
7e87afed-a5ab-4285-b257-b84161febe58	分支修改	/branch_edit	d1f06d6a-a377-4419-b690-1cbe09879b1c	分支修改	el-icon-date	BranchEdit			6	1	1970-01-01 00:00:00
d0f0d5f9-e375-4c11-aaac-22ae71a637ea	api授权	/pay_api_auth	adcbffee-495f-4e59-90c0-f438d39ebd07	api授权	el-icon-date	PayApiAuth			14	1	1970-01-01 00:00:00
230733a9-34d8-47b9-aa97-83a85dbca11e	门店查看	/store_check	42fc26ac-7638-4b3b-9ed3-4b0a3f964512	门店查看	el-icon-date	StoreCheck			5	1	1970-01-01 00:00:00
315a3bae-c0a7-47df-9522-acbc971ff8b1	终端配置查看	/endpoint_config_check	d1f06d6a-a377-4419-b690-1cbe09879b1c	终端配置查看	el-icon-date	EndpointConfigCheck			7	1	1970-01-01 00:00:00
189208bd-a65f-44d8-8a07-96647feaafb7	我的子商户列表	/merchant_sub_list	adcbffee-495f-4e59-90c0-f438d39ebd07	我的子商户列表	el-icon-date	MerchantSubList			2	0	1970-01-01 00:00:00
ad45fb2b-38a8-4d14-8eb5-03af12a2f7ec	贷款银行	/loan_banks	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	贷款银行	el-icon-date	LoanBanks			2	0	1970-01-01 00:00:00
5f65f5d0-39d7-42e6-a7f5-dd21ba72abb1	贷款银行修改	/loan_bank_edit	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	贷款银行修改	el-icon-date	LoanBankEdit			3	1	1970-01-01 00:00:00
267939c7-d951-47f8-aad8-63f0c0d6cdfb	借款申请	/loan_application	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	借款申请	el-icon-date	LoanApplication			1	0	1970-01-01 00:00:00
8feea562-c412-4546-916a-7c22e60c833c	已关联银行商户	/related_banks_merc	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	已关联银行商户	el-icon-date	RelatedBanksMerc			5	0	1970-01-01 00:00:00
ab1aa816-baed-4d33-bd9b-80d4b83a0341	门店日结算	/store_day_settlement	797cecc2-d1f9-499a-b79d-e1f7e77c7874	门店日结算	el-icon-date	StoreDaySettlement			0	0	1970-01-01 00:00:00
59d15cbe-0f30-4057-ab75-09ba8776b7d6	线上对账文件	/check_details	cb332a03-aeee-435c-86de-b7d577e3817b	线上对账文件	el-icon-date	CheckDetails			10	0	1970-01-01 00:00:00
e9cb3b58-b3ed-4e27-a2f5-0c704dbbf32f	下属账号登录	/under_account_login	ed06f369-9d92-4d07-8b9a-6488694b4d1c	下属账号登录	el-icon-date	UnderAccountLogin			13	0	1970-01-01 00:00:00
18460b2d-6e29-43da-9419-990fd196de7c	结算钱包日志	/bill_wallet_log	ed06f369-9d92-4d07-8b9a-6488694b4d1c	结算钱包日志	el-icon-date	BillWalletLog			14	0	1970-01-01 00:00:00
4e471ceb-766f-4575-ab73-856e07ceb1a8	余额批量代付申请	/send_apply_balance	e0be0b73-1f90-4c04-8da5-86513c148a16	余额批量代付申请	el-icon-date	SendApplyBalance			20	0	1970-01-01 00:00:00
cab6e3e8-cff1-42da-959f-e12705adc760	对公代付申请	/contrary_apply_edit	e0be0b73-1f90-4c04-8da5-86513c148a16	对公代付申请	el-icon-date	ContraryApplyEdit			22	1	1970-01-01 00:00:00
1c80833a-d9f9-4b31-ac77-55c282781167	对公代付申请列表	/contrary_apply_list	e0be0b73-1f90-4c04-8da5-86513c148a16	对公代付申请列表	el-icon-date	ContraryApplyList			21	0	1970-01-01 00:00:00
a57004e3-4bac-4cc0-b2c6-8b0a060e2dbb	点击流水	/click_bills	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	点击流水	el-icon-date	ClickBills			6	1	1970-01-01 00:00:00
8b92de65-89b3-4f4b-80d2-a2908bdca9ea	进件日志列表	/incoming_log_list	6419a0a5-1e69-4775-8e72-04ee288635ff	进件日志列表	el-icon-date	IncomingLogList			4	0	1970-01-01 00:00:00
a0dee210-7e5b-4309-ad4e-7ccdf1e7e82e	银行限额列表	/bank_limit_list	5a229d4e-beee-4f0a-ba2f-888f04a149e5	银行限额列表	el-icon-date	BankLimitList			1	0	1970-01-01 00:00:00
b3f2d37b-3f4b-4fa6-b632-3a85e295fd01	首页	/merc_index	ed06f369-9d92-4d07-8b9a-6488694b4d1c	首页	el-icon-date	MercIndex			1	0	1970-01-01 00:00:00
44bda315-99d8-4619-99c8-877e24f259db	自动代付日志	/auto_agentpays	6419a0a5-1e69-4775-8e72-04ee288635ff	自动代付日志	el-icon-date	AutoAgentpays			5	0	1970-01-01 00:00:00
7586f8fd-b585-4555-98d1-d458115056e0	代付订单日志	/auto_agentpay_logs	6419a0a5-1e69-4775-8e72-04ee288635ff	代付订单日志	el-icon-date	AutoAgentpayLogs			6	0	1970-01-01 00:00:00
e4db0d2c-9735-4e9a-b698-c4a395224d48	登录日志	/login_logs	6419a0a5-1e69-4775-8e72-04ee288635ff	登录日志	el-icon-date	LoginLogs			7	0	1970-01-01 00:00:00
5219e94c-69eb-4baa-86ac-593d0d0255fc	分账设置	/ledger_account_set	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	分账设置	el-icon-date	LedgerAccountSet			10	1	1970-01-01 00:00:00
88c58e99-19a7-4f72-811e-92e21346207a	通道对账分析	/channel_settle_logs	cb332a03-aeee-435c-86de-b7d577e3817b	通道对账分析	el-icon-date	ChannelSettleLogs			10	0	1970-01-01 00:00:00
91b45a4d-262a-4a5f-a376-845bb5b8e489	统计报表配置	/statistical_report_list	9392bd8f-6ef5-4f16-b87b-a36709339d5d	统计报表配置	el-icon-date	StatisticalReportList			8	0	1970-01-01 00:00:00
70a95de6-a6f3-40a3-a70d-5526d23df7bf	统计报表配置修改	/statistical_report_edit	9392bd8f-6ef5-4f16-b87b-a36709339d5d	统计报表配置修改	el-icon-date	StatisticalReportEdit			9	1	1970-01-01 00:00:00
b8e3f794-9775-4363-9bcf-d3d929f6ba27	商户支付统计	/merc_pay_analysis	797cecc2-d1f9-499a-b79d-e1f7e77c7874	商户支付统计	el-icon-date	MercPayAnalysis			10	0	1970-01-01 00:00:00
a35f9884-abfa-4c87-9c23-a983e04dfaf7	银行已关联商户	/loan_bank_merc_check	6f775fcf-6b53-49a6-960f-5be9ea2c32f7	银行已关联商户	el-icon-date	LoanBankMercCheck			7	1	1970-01-01 00:00:00
a5226119-0c62-4194-b097-d608a41e00c0	科目列表	/subject_list	a317c13a-414a-4f44-9b46-3c9fe0b735ee	科目列表	el-icon-date	SubjectList			1	0	1970-01-01 00:00:00
7d89470c-955d-4ad7-bee4-82b65eaa07c0	文档列表	/document_list	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	文档列表	el-icon-date	DocumentList			6	0	1970-01-01 00:00:00
7ea8ce89-9ce7-4a89-a332-7baad7bbe475	被风控列表	/by_wind_list	5a229d4e-beee-4f0a-ba2f-888f04a149e5	被风控列表	el-icon-date	ByWindList			11	0	1970-01-01 00:00:00
9b374efe-f7c1-44ff-a122-52981146b22a	T1结算明细	/t1_settle_details	cb332a03-aeee-435c-86de-b7d577e3817b	T1结算明细	el-icon-date	T1SettleDetails			12	1	1970-01-01 00:00:00
eedb53fb-3a3e-4f32-a232-ffc97a4464ae	T1结算审核	/t1_settle_review	cb332a03-aeee-435c-86de-b7d577e3817b	T1结算审核	el-icon-date	T1SettleReview			11	0	1970-01-01 00:00:00
b36e11e0-99e1-4af5-b83b-6fe9062554a8	交易报表下载	/merc_reconcilia	cb332a03-aeee-435c-86de-b7d577e3817b	交易报表下载	el-icon-date	MercReconcilia			8	0	1970-01-01 00:00:00
8ec08232-29d6-418f-8560-64f307684f8e	生成交易报表	/gen_offline_reports	cb332a03-aeee-435c-86de-b7d577e3817b	生成交易报表	el-icon-date	GenOfflineReports			14	0	1970-01-01 00:00:00
d2e33ead-1f7a-4482-a8e0-d58deceea631	模板管理	/template_management	14aa553e-5be3-404d-adb0-b4440e023c69	模板管理	el-icon-date	TemplateManagement			5	0	1970-01-01 00:00:00
e18c306e-1d61-4e6e-a2cf-c41ba357b406	月结代付申请	/monthly_statement_payment	b4aee72a-2ac5-4748-b846-bd1828139858	月结代付申请	el-icon-date	MonthlyStatementPayment			5	0	1970-01-01 00:00:00
f42406c3-08c8-443a-845a-996e04c55753	redis管理	/redis_manage	14aa553e-5be3-404d-adb0-b4440e023c69	redis管理	el-icon-date	RedisManage			4	0	1970-01-01 00:00:00
01698813-a221-4b56-8bf5-ec3ec9956d97	未结算审核	/no_settle_audit	e0be0b73-1f90-4c04-8da5-86513c148a16	未结算审核	el-icon-date	NoSettleAudit			4	0	1970-01-01 00:00:00
0576302e-61ce-49b3-b265-28b7458e5c27	出金管理	/withdrawal_manage	e0be0b73-1f90-4c04-8da5-86513c148a16	出金管理	el-icon-date	WithdrawalManage			14	0	1970-01-01 00:00:00
3c53a6aa-9d98-43be-bd14-4c5501e58ee5	出金管理编辑	/withdrawal_manage_edit	e0be0b73-1f90-4c04-8da5-86513c148a16	出金管理编辑	el-icon-date	WithdrawalManageEdit			15	1	1970-01-01 00:00:00
aa061751-3180-40dc-b101-f866d7e4dc2f	月结文件	/monthly_statement_file	b4aee72a-2ac5-4748-b846-bd1828139858	月结文件	el-icon-date	MonthlyStatementFile			6	0	1970-01-01 00:00:00
70ebed7b-5653-42cf-be51-902c9dce2987	文档下载	/document_download	ed06f369-9d92-4d07-8b9a-6488694b4d1c	文档下载	el-icon-date	DocumentDownload			17	0	1970-01-01 00:00:00
13f5ca18-5e46-437e-a835-4a36fd6f9955	套餐树	/rate_package_tree	797cecc2-d1f9-499a-b79d-e1f7e77c7874	套餐树	el-icon-date	RatePackageTree			15	0	1970-01-01 00:00:00
57873320-1f3c-4d56-ad41-122f17f362a2	商户钱包	/merchant_balance	ed06f369-9d92-4d07-8b9a-6488694b4d1c	商户钱包	el-icon-date	MerchantBalance			17	0	1970-01-01 00:00:00
b39c59d0-f9d1-4901-af1c-f3a193d1c612	农信初始化	/nx_init	b639edba-fe15-48c6-9e53-551c056a4426	农信初始化	el-icon-date	NxInit			2	0	1970-01-01 00:00:00
07bd2d37-21cd-43b3-9d38-eeafaa73bcf9	菜单列表	/menu_list	222d3e77-c07d-4c30-a384-b0c493a3cb38	菜单列表	el-icon-date	MenuList			20	0	1970-01-01 00:00:00
d88a3357-b68c-406e-8951-97df93df66bb	通道对账	/download_files	cb332a03-aeee-435c-86de-b7d577e3817b	通道对账	el-icon-date	DownloadFiles			1	1	1970-01-01 00:00:00
3ed50a46-3620-43a1-b8ea-5a1ebc968851	手续费上传	/upload_fee	e0be0b73-1f90-4c04-8da5-86513c148a16	手续费上传	el-icon-date	UploadFee			23	0	1970-01-01 00:00:00
583c9de5-4079-4eb9-b448-4fa66ac594fd	第三方商户修改	/third_party_merc_edit	a59a1727-baa8-472a-ab4c-fe0f185fbc5e	第三方商户修改	el-icon-date	ThirdPartyMercEdit			8	1	1970-01-01 00:00:00
b207a77e-a3b3-4303-b757-2bd70ca8639b	商户审核	/merc_audit	adcbffee-495f-4e59-90c0-f438d39ebd07	商户审核	el-icon-date	MercAudit			18	1	1970-01-01 00:00:00
080f0752-21ba-4d3e-97c2-4846dfb4f648	出金审核	/withdrawal_review	e0be0b73-1f90-4c04-8da5-86513c148a16	出金审核	el-icon-date	WithdrawalReview			16	0	1970-01-01 00:00:00
82dbdde0-be98-462f-bef7-b86ed1a98336	出金查看	/withdrawal_check	e0be0b73-1f90-4c04-8da5-86513c148a16	出金查看	el-icon-date	WithdrawalCheck			25	1	1970-01-01 00:00:00
ab03d85a-45db-4af0-8b37-0de7e8bfe4a6	rsa_p12解码	/rsa_priv_p12	2feb45eb-2e97-4d7f-bf43-d221b1c71928	rsa_p12解码	el-icon-date	RsaPrivP12			1	0	1970-01-01 00:00:00
fa944b01-d087-4729-920d-f1f7ab36a642	生成rsa密钥对	/gen_rsa_key_pair	2feb45eb-2e97-4d7f-bf43-d221b1c71928	生成rsa密钥对	el-icon-date	GenRsaKeyPair			2	0	1970-01-01 00:00:00
1e1a0c33-63e3-45e7-aae6-5856c0add615	产品列表	/product_list	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	产品列表	el-icon-date	ProductList			16	0	1970-01-01 00:00:00
ea6a4e30-d810-41fd-b231-a9a209adbae3	产品修改	/product_edit	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	产品修改	el-icon-date	ProductEdit			17	1	1970-01-01 00:00:00
359d32ec-2919-4ca3-b927-232440bd886f	产品查看	/product_check	1d2e808a-62bc-41cc-996a-cd1b2e2b14c4	产品查看	el-icon-date	ProductCheck			18	1	1970-01-01 00:00:00
ccad896d-13db-4fc8-be92-574af1d80341	卖出管理	/score_sell_manage	72237dba-0beb-450e-9df6-ea38d3255cab	卖出管理	el-icon-menu	Content	hl_sell_list	卖出列表	5	0	1970-01-01 00:00:00
223c3afe-ef65-4b86-9733-6284396cd65f	设备管理员列表	/device_admin_list	fc72276f-66bb-4609-b83a-4a84ff55100c	设备管理员列表	el-icon-date	DeviceAdminList			1	0	1970-01-01 00:00:00
67522e75-9672-4bbc-904d-57dd942290b7	设备管理员修改	/device_admin_edit	fc72276f-66bb-4609-b83a-4a84ff55100c	设备管理员修改	el-icon-date	DeviceAdminEdit			2	1	1970-01-01 00:00:00
78cd88ae-d982-48c6-998c-1a6ba961b5f0	设备管理员查看	/device_admin_check	fc72276f-66bb-4609-b83a-4a84ff55100c	设备管理员查看	el-icon-date	DeviceAdminCheck			3	1	1970-01-01 00:00:00
131cc1b8-2668-4392-9443-6c2e42a426ef	设备列表	/device_list	fc72276f-66bb-4609-b83a-4a84ff55100c	设备列表	el-icon-date	DeviceList			4	0	1970-01-01 00:00:00
20805d38-e63d-4117-b902-7007c1697a13	设备修改	/device_edit	fc72276f-66bb-4609-b83a-4a84ff55100c	设备修改	el-icon-date	DeviceEdit			5	1	1970-01-01 00:00:00
f53252df-70b6-4907-8847-908554862006	设备查看	/device_check	fc72276f-66bb-4609-b83a-4a84ff55100c	设备查看	el-icon-date	DeviceCheck			6	1	1970-01-01 00:00:00
6f4cc747-5450-4514-9a3e-123ace9cb232	一级服务商列表	/hl_one_agent_list	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	一级服务商列表	el-icon-date	HlOneAgentList			1	0	1970-01-01 00:00:00
58405e28-6fff-42e1-8567-01466845abfb	卖出列表	/hl_sell_list	ccad896d-13db-4fc8-be92-574af1d80341	卖出列表	el-icon-date	HlSellList			1	0	1970-01-01 00:00:00
93421699-5c88-4dee-995c-f74222b3a1a7	三级服务商列表	/hl_three_agent_list	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	三级服务商列表	el-icon-date	HlThreeAgentList			3	0	1970-01-01 00:00:00
3566c88a-f2b5-4c1e-b8a3-68a26b49e158	订单记录	/hl_order_record	24127ef9-9ef3-46f3-88cd-97378e694a6c	订单记录	el-icon-date	HlOrderRecord			4	1	1970-01-01 00:00:00
03264ac4-7754-48ef-a824-5269d1f84e4b	申诉详情	/hl_appeal_detail	77f2a333-b579-421b-bf3f-5b913e6002b5	申诉详情	el-icon-date	HlAppealDetail			2	1	1970-01-01 00:00:00
51f985e0-4f2a-412a-adb5-3adf5399fc75	申诉列表	/hl_appeal_list	77f2a333-b579-421b-bf3f-5b913e6002b5	申诉列表	el-icon-date	HlAppealList			1	0	1970-01-01 00:00:00
0cd8481f-84d0-43a8-b9d3-5b3ba2d6e2ed	兑付列表	/hl_redemption_list	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	兑付列表	el-icon-date	HlRedemptionList			1	0	1970-01-01 00:00:00
b48dccdc-da4f-4914-bd2e-1971de87fda1	礼包列表	/gift_list	f22b7a3d-2d76-483c-af40-dd0c04defb20	礼包列表	el-icon-date	GiftList			1	0	1970-01-01 00:00:00
3aa46e51-bac0-4711-bf1b-bb396e71fe06	礼包修改	/gift_edit	f22b7a3d-2d76-483c-af40-dd0c04defb20	礼包修改	el-icon-date	GiftEdit			2	1	1970-01-01 00:00:00
27917132-65e1-4218-bdfc-40a62550a141	礼包查看	/gift_check	f22b7a3d-2d76-483c-af40-dd0c04defb20	礼包查看	el-icon-date	GiftCheck			3	1	1970-01-01 00:00:00
09cfad53-de06-4a25-88c0-9eaa26345df9	礼包类型列表	/gift_type_list	f22b7a3d-2d76-483c-af40-dd0c04defb20	礼包类型列表	el-icon-date	GiftTypeList			4	0	1970-01-01 00:00:00
7767088d-9354-42d4-8cc7-9d3f2c0bef88	礼包类型修改	/gift_type_edit	f22b7a3d-2d76-483c-af40-dd0c04defb20	礼包类型修改	el-icon-date	GiftTypeEdit			5	1	1970-01-01 00:00:00
b49aa1c7-1f8e-4d0d-bdaf-09d3ec791ff4	礼包类型查看	/gift_type_check	f22b7a3d-2d76-483c-af40-dd0c04defb20	礼包类型查看	el-icon-date	GiftTypeCheck			6	1	1970-01-01 00:00:00
39732090-0856-4b05-a241-df9886a313a5	买入列表	/buy_list	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	买入列表	el-icon-date	BuyList			1	0	1970-01-01 00:00:00
05a2e58e-0551-40c5-9bff-657768d7e689	兑付详情	/hl_redemption_detail	7eddcf4a-530a-4210-8e1c-17da3d0f87c0	兑付详情	el-icon-date	HlRedemptionDetail			2	1	1970-01-01 00:00:00
fb9fbdf7-d05f-40ac-95e1-e10019ed8165	意见反馈列表	/hl_feedback_list	e31f3e0e-1158-4f37-97a4-bedbfd2a6589	意见反馈列表	el-icon-date	HlFeedbackList			1	0	1970-01-01 00:00:00
9a84a06f-73ba-4a9d-9818-dfd8bbcfb1fa	用户列表	/hl_account_list	24127ef9-9ef3-46f3-88cd-97378e694a6c	用户列表	el-icon-date	HlAccountList			1	0	1970-01-01 00:00:00
fc72276f-66bb-4609-b83a-4a84ff55100c	设备管理	/device_manage	fcea9f1c-b882-4368-a7a3-c8bbd2491684	设备管理	el-icon-menu	Content	device_manage	设备列表	1	0	1970-01-01 00:00:00
c45be260-e11b-4fae-8e33-05d0803f4e6b	统计信息	/hl_statistics	24127ef9-9ef3-46f3-88cd-97378e694a6c	统计信息	el-icon-date	HlStatistics			3	1	1970-01-01 00:00:00
d81d630f-3d80-4168-923b-65eefae629c9	银行卡列表	/card_limit_list	ae18fd5a-e0ba-4ace-a637-902317a67773	银行卡列表	el-icon-date	CardLimitList			2	0	1970-01-01 00:00:00
77f2a333-b579-421b-bf3f-5b913e6002b5	申诉管理	/score_appeal_manage	72237dba-0beb-450e-9df6-ea38d3255cab	申诉管理	el-icon-menu	Content	hl_appeal_list	申诉列表	11	0	1970-01-01 00:00:00
aca0acfd-7930-4b4d-8a5b-5b78aaaa4db1	卖出订单详情	/hl_sell_detail	ccad896d-13db-4fc8-be92-574af1d80341	卖出订单详情	el-icon-date	HlSellDetail			2	1	1970-01-01 00:00:00
c070d24f-76eb-4684-8459-dbb7977acc46	二级服务商列表	/hl_two_agent_list	0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	二级服务商列表	el-icon-date	HlTwoAgentList			2	0	1970-01-01 00:00:00
7eddcf4a-530a-4210-8e1c-17da3d0f87c0	兑付管理	/score_redemption_manage	72237dba-0beb-450e-9df6-ea38d3255cab	兑付管理	el-icon-menu	Content	hl_redemption_list	兑付列表	6	0	1970-01-01 00:00:00
e31f3e0e-1158-4f37-97a4-bedbfd2a6589	意见反馈	/score_feedback	72237dba-0beb-450e-9df6-ea38d3255cab	意见反馈	el-icon-menu	Content	hl_feedback_list	意见反馈列表	10	0	1970-01-01 00:00:00
3f5a7355-9e0e-40ea-9d5c-81103b167e2c	买入管理	/score_buy_manage	72237dba-0beb-450e-9df6-ea38d3255cab	买入管理	el-icon-menu	Content	buy_list	买入列表	4	0	1970-01-01 00:00:00
f22b7a3d-2d76-483c-af40-dd0c04defb20	商品管理	/score_good_manage	72237dba-0beb-450e-9df6-ea38d3255cab	商品管理	el-icon-menu	Content	gift_list	礼包列表	3	0	1970-01-01 00:00:00
99a0afca-2fef-4083-82a2-4785ac933ec0	设备二维码列表	/device_qrcode_list	fc72276f-66bb-4609-b83a-4a84ff55100c	设备二维码列表	el-icon-date	DeviceQrcodeList			7	0	1970-01-01 00:00:00
e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	首页	/score_my	72237dba-0beb-450e-9df6-ea38d3255cab	首页	el-icon-menu	Content	hl_index	首页	1	0	1970-01-01 00:00:00
4ec8150a-0cbe-48f5-9dee-974fcb3e0558	银行限额列表	/bank_money_limit_list	ae18fd5a-e0ba-4ace-a637-902317a67773	银行限额列表	el-icon-date	BankMoneyLimitList			1	0	1970-01-01 00:00:00
24127ef9-9ef3-46f3-88cd-97378e694a6c	用户管理	/score_account_manage	72237dba-0beb-450e-9df6-ea38d3255cab	用户管理	el-icon-menu	Content	hl_account_list	用户列表	8	0	1970-01-01 00:00:00
f3c2ef69-090e-4f61-a9d6-954330e86d8f	卡商管理	/score_cm	72237dba-0beb-450e-9df6-ea38d3255cab	卡商管理	el-icon-menu	Content			12	0	1970-01-01 00:00:00
9c2c0ede-430f-4981-bb5d-7fd23b1bbe6c	订单详情(买入)	/hl_buy_order_detail	3f5a7355-9e0e-40ea-9d5c-81103b167e2c	订单详情(买入)	el-icon-date	HlBuyOrderDetail			2	1	1970-01-01 00:00:00
c2ebd879-c09a-43d3-8088-6b0c2445c0b1	银行编码列表	/bank_code_list	77788f4d-5d0e-458b-b7cb-297e907b2eff	银行编码列表	el-icon-date	BankCodeList			2	0	1970-01-01 00:00:00
af0187d0-e84c-401f-9ffc-70194b099309	银行缩写列表	/bank_abbr_list	77788f4d-5d0e-458b-b7cb-297e907b2eff	银行缩写列表	el-icon-date	BankAbbrList			1	0	1970-01-01 00:00:00
e296b896-76d1-4cf0-b092-e637a5909ec4	银行卡列表	/bank_card_list	77788f4d-5d0e-458b-b7cb-297e907b2eff	银行卡列表	el-icon-date	BankCardList			3	0	1970-01-01 00:00:00
24a99997-51d4-446c-bc51-2f46e3ca01a6	卡商管理者列表	/hl_score_cm_manager	f3c2ef69-090e-4f61-a9d6-954330e86d8f	卡商管理者列表	el-icon-date	HlScoreCmManager			1	0	1970-01-01 00:00:00
9fc4b56e-d391-4564-9dcd-2c47f9e436c3	小卡商列表	/hl_score_cm_worker	f3c2ef69-090e-4f61-a9d6-954330e86d8f	小卡商列表	el-icon-date	HlScoreCmWorker			2	0	1970-01-01 00:00:00
07722dd0-088c-4309-9fd4-7097cb7a5c58	结算配置	/set_config	1585834c-8fdf-4d7d-a590-2bb3cf4bd6cc	结算配置	el-icon-date	SetConfig			4	1	1970-01-01 00:00:00
fe4ff53d-f25a-489b-a9e7-ea0df86fe4b5	我的管理	/user_manage	6cc30ed8-3fe5-49b5-b800-3361298c3a85	我的管理	el-icon-s-goods	Content	user_manage		1	0	1970-01-01 00:00:00
c74df4ad-f3d3-4362-8011-ed5fdfba4f4d	银行卡收款统计	/card_recv_stat	086e9acf-0405-449c-9dee-ee33d80df312	银行卡收款统计	el-icon-bank-card	CardRecvStatList	WorkBench		3	0	1970-01-01 00:00:00
086e9acf-0405-449c-9dee-ee33d80df312	财务管理	/msc_financial	72237dba-0beb-450e-9df6-ea38d3255cab	MSC财务管理	el-icon-collection	Content		赠送币统计	13	0	1970-01-01 00:00:00
0c4a97b5-5dcb-4ec6-9eaa-cc74f9b041bc	服务商管理	/score_agent_manage	72237dba-0beb-450e-9df6-ea38d3255cab	服务商管理	el-icon-menu	Content	hl_big_agent_list	一级代理列表	7	0	1970-01-01 00:00:00
ae18fd5a-e0ba-4ace-a637-902317a67773	银行卡管理	/score_bank_card_manage	72237dba-0beb-450e-9df6-ea38d3255cab	银行卡管理	el-icon-menu	Content	bank_money_limit_list	银行卡限额列表	2	0	1970-01-01 00:00:00
ee775cf8-f617-4324-9220-726c7ddbe919	成功率	/hl_success_rate	797cecc2-d1f9-499a-b79d-e1f7e77c7874	成功率	el-icon-date	HlSuccessRate			12	0	1970-01-01 00:00:00
c2dde18e-3fab-450c-b2c0-554fa1dcbf6b	收款渠道设置	/hl_receipt_channel_set	9f4da9ef-7f9c-407e-91a9-034f623b810f	收款渠道设置	el-icon-menu	HlReceiptChannelSet			1	0	1970-01-01 00:00:00
9f4da9ef-7f9c-407e-91a9-034f623b810f	基础设置	/score_channel_set	72237dba-0beb-450e-9df6-ea38d3255cab	基础设置	el-icon-menu	Content	score_channel_set	收款渠道设置	9	0	1970-01-01 00:00:00
ebb13d64-4568-4cce-9f26-72f5e8a9bc14	轮播图列表	/hl_bananas_set	9f4da9ef-7f9c-407e-91a9-034f623b810f	轮播图列表	el-icon-date	HlBananasSet			2	0	1970-01-01 00:00:00
3780a865-b70a-4a34-93f1-56c920e44100	设备池列表	/device_pool_list	fc72276f-66bb-4609-b83a-4a84ff55100c	设备池列表	el-icon-date	DevicePoolList			8	0	1970-01-01 00:00:00
a06b41ea-3abb-4dc4-a4f7-90bfc86a9592	设备池修改	/device_pool_edit	fc72276f-66bb-4609-b83a-4a84ff55100c	设备池修改	el-icon-date	DevicePoolEdit			9	1	1970-01-01 00:00:00
0a9b8878-506b-4f76-af54-478e96f06eb9	设备池查看	/device_pool_check	fc72276f-66bb-4609-b83a-4a84ff55100c	设备池查看	el-icon-date	DevicePoolCheck			10	1	1970-01-01 00:00:00
25c7b3c0-efd1-4d00-a748-e37de9ca4e3c	商户关联设备池列表	/merc_rela_pool_list	fc72276f-66bb-4609-b83a-4a84ff55100c	商户关联设备池列表	el-icon-date	MercRelaPoolList			11	0	1970-01-01 00:00:00
68fc3f82-6076-412f-b986-20d3cb1bceaa	交易记录	/hl_sell_transaction_record	ccad896d-13db-4fc8-be92-574af1d80341	交易记录	el-icon-date	HlSellTransactionRecord			3	1	1970-01-01 00:00:00
86ce4e37-3234-4bb4-bb03-b1d3742ac42e	工作台	/hl_index	e0bed7b3-d7ab-4bc6-b852-501d1b6e7ed5	工作台	el-icon-date	HlIndex			1	0	1970-01-01 00:00:00
288731cb-561c-4d4a-8951-d8c359a51846	轮播图修改	/hl_bananas_set_edit	9f4da9ef-7f9c-407e-91a9-034f623b810f	轮播图修改	el-icon-date	HlBananasSetEdit			3	1	1970-01-01 00:00:00
653eb679-3e69-442e-b0ff-fd22e3fd8cd3	赠送币统计	/reward_coin_list	086e9acf-0405-449c-9dee-ee33d80df312	赠送币统计	el-icon-toilet-paper	RewardCoinList	WorkBench		1	0	1970-01-01 00:00:00
53819810-8a9e-4741-bfaa-80bd10264a06	测试列表		fed719a1-3a39-4051-beb3-eba29a803a34	测试列表	el-icon-picture-outline-round				1	0	1970-01-01 00:00:00
5493c560-ace3-440c-a36a-9538f618007a	直推统计	/hl_straight_push_statistics	086e9acf-0405-449c-9dee-ee33d80df312	直推统计	el-icon-date	HlStraightPushStatistics			13	0	1970-01-01 00:00:00
2ad759be-3a54-4efa-aeda-7027dc97867f	交易统计	/merc_statistics	797cecc2-d1f9-499a-b79d-e1f7e77c7874	交易统计	el-icon-date	MercStatistics			12	0	1970-01-01 00:00:00
5aae3b9c-87f3-487f-b016-f7109f17d2c0	多语言列表	/multi_language_list	9f4da9ef-7f9c-407e-91a9-034f623b810f	多语言列表	el-icon-date	MultiLanguageList			4	0	1970-01-01 00:00:00
2fdb5752-a2e2-4705-8c9b-b4a618e6de1f	测试管理	/test_manage	b01086fd-b0c8-4c3f-8ca9-f05cd737ae9f	测试管理	el-icon-setting	Content	test_manage		5	0	1970-01-01 00:00:00
5a31e833-8621-4b3b-841a-a0b03a8575e1	测试账号	/test_list	2fdb5752-a2e2-4705-8c9b-b4a618e6de1f	测试账号		TestList	test_list		1	0	1970-01-01 00:00:00
c8f239e6-0791-4bdc-95f0-8e96abb74aa3	测试账号	/test_list	213afcc4-81ca-4af6-aaad-3a83a55302a4	测试账号	el-icon-s-operation	TestList	test_list		1	0	1970-01-01 00:00:00
3c15f3e4-bcab-4157-99d3-6332f2544089	app崩溃日志	/hl_crash_log_list	24127ef9-9ef3-46f3-88cd-97378e694a6c	app崩溃日志	el-icon-s-platform	HlCrashLogList			5	0	1970-01-01 00:00:00
afe7beea-4242-4b06-9d4f-00dd11aae178	app崩溃日志查看	/hl_carsh_log_check	24127ef9-9ef3-46f3-88cd-97378e694a6c	app崩溃日志查看	el-icon-date	HlCarshLogCheck			6	1	1970-01-01 00:00:00
1e5eb584-f19e-4110-89e7-2d63a08d727a	卖出交易记录列表	/hl_sell_transaction_record_list	ccad896d-13db-4fc8-be92-574af1d80341	卖出交易记录列表	el-icon-date	HlSellTransactionRecordList			4	0	1970-01-01 00:00:00
26b070f5-65f4-44b7-ad6d-2f8630937f3c	手机列表	/phone_list	24127ef9-9ef3-46f3-88cd-97378e694a6c	手机列表	el-icon-date	PhoneList			7	0	1970-01-01 00:00:00
32a0d6bb-1c39-4701-a0fb-ecd39f0f76d4	手机查看	/phone_check	24127ef9-9ef3-46f3-88cd-97378e694a6c	手机查看	el-icon-date	PhoneCheck			8	1	1970-01-01 00:00:00
b439edb9-605f-47dc-aee2-8df45976736e	服务商分润统计	/agency_reward_list	086e9acf-0405-449c-9dee-ee33d80df312	服务商分润统计	el-icon-medal	AgencyRewardCoinList	WorkBench		2	0	1970-01-01 00:00:00
d5e71527-bc03-4f52-b4e6-8460cbd3924b	交易方式列表	/hl_pay_method_list	ccad896d-13db-4fc8-be92-574af1d80341	交易方式列表	el-icon-date	HlPayMethodList			5	0	1970-01-01 00:00:00
23bcf6d0-fa48-4ae8-8036-07f675e240b7	补单审核列表	/force_callabck_list	ccad896d-13db-4fc8-be92-574af1d80341	补单审核列表	el-icon-s-flag	ForceCallbackList			6	0	1970-01-01 00:00:00
1d64225d-ad7b-4cb9-910f-49283f636cca	获取ip	/get_ip	6419a0a5-1e69-4775-8e72-04ee288635ff	获取ip	el-icon-s-grid	GetIpList			8	0	1970-01-01 00:00:00
ce3db5e2-8372-475b-b492-f727f5240569	卖出当天统计	/hl_sell_stat_daily_list	ccad896d-13db-4fc8-be92-574af1d80341	卖出当天统计	el-icon-s-custom	HlSellStatDailyList			7	0	1970-01-01 00:00:00
d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	用户管理	/user_manage	149d7564-5e2f-4472-a637-d1db61dfdc81	用户管理	el-icon-s-custom	Content	user_manage		1	0	1970-01-01 00:00:00
dd351f02-ed87-426b-952e-e1247841b8c3	用户列表	/user_management_list	47ad1c65-a39a-4c11-8d8d-3c2362907c49	用户列表	el-icon-date	UserManagemenList			1	0	1970-01-01 00:00:00
07679b82-27d2-4811-85e6-8a6e73461613	用户列表	/user_management_list	d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	用户列表	el-icon-s-order	UserManagemenList			1	0	1970-01-01 00:00:00
9ec25419-48ea-4525-8832-2f9cbcdc3096	商户列表	/easy_merc_list	adcbffee-495f-4e59-90c0-f438d39ebd07	商户列表	el-icon-date	EasyMercList			25	0	1970-01-01 00:00:00
81d5505b-df37-4cef-928c-10730baa7b10	服务商列表	/servicer_managemen_list	996af46e-f32c-4b26-8109-2b942d1aa8b8	服务商列表	el-icon-s-order	ServicerManagemenList			1	0	1970-01-01 00:00:00
3800a410-1da5-46a2-958f-163f98034d55	服务商配置	/servicer_config	996af46e-f32c-4b26-8109-2b942d1aa8b8	服务商配置	el-icon-s-order	ServicerConfig			2	1	1970-01-01 00:00:00
81d8c6fd-f427-409c-aac9-d9dac8e33f97	用户账户明细	/user_balance_list	d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	用户账户明细	el-icon-s-order	UserBalanceList			3	0	1970-01-01 00:00:00
850c0f8d-2235-431c-8e0a-d02baa738c58	用户信息	/cust_list	d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	用户信息	el-icon-s-custom	CustList			2	1	1970-01-01 00:00:00
149d7564-5e2f-4472-a637-d1db61dfdc81	用户中心	/cust_center	00000000-0000-0000-0000-000000000000	用户中心	el-icon-s-custom	Default	cust_center	/cust_center	1	0	1970-01-01 00:00:00
e428e7fe-5311-447d-a13f-30d8bbf04653	账单明细	/user_bill_list	d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	账单明细	el-icon-s-custom	UserBillList			4	1	1970-01-01 00:00:00
8caf4478-ee5f-48f2-b0cb-a35b2bdc8aa0	服务商交易查询	/service_general_ledger_list	996af46e-f32c-4b26-8109-2b942d1aa8b8	服务商交易查询	el-icon-s-order	ServiceGeneralLedgerList			7	0	1970-01-01 00:00:00
996af46e-f32c-4b26-8109-2b942d1aa8b8	服务商管理	/servicer_managemen	149d7564-5e2f-4472-a637-d1db61dfdc81	服务商管理	el-icon-s-custom	Content	servicer_managemen		2	0	1970-01-01 00:00:00
d68e979e-1fb7-400f-b2d3-4f8b31a08231	账单明细	/billing_details_results_list	b57a6f39-2ec1-4065-911a-29764b00ec2b	账单明细	el-icon-s-data	BillingDetailsResultsList			3	1	1970-01-01 00:00:00
d5c1e2ba-4969-45d4-8b30-8da6cc872ae6	新增收款账户	/channel_edit	b57a6f39-2ec1-4065-911a-29764b00ec2b	新增收款账户	el-icon-picture-outline-round	ChannelEdit			5	1	1970-01-01 00:00:00
08454882-e75c-4933-b522-6687fda8cb24	添加服务商	/servicer_add	996af46e-f32c-4b26-8109-2b942d1aa8b8	添加服务商	el-icon-s-order	ServicerAdd			3	1	1970-01-01 00:00:00
2c2961d0-ef1e-4c1c-8f28-94db86f33b16	服务商信息详情	/servicer_info	996af46e-f32c-4b26-8109-2b942d1aa8b8	服务商信息详情	el-icon-s-order	ServicerInfo			5	1	1970-01-01 00:00:00
67d8761f-764e-4e5e-92df-bcb2bfae02ee	指定运营商统计	/servicer_order_count	996af46e-f32c-4b26-8109-2b942d1aa8b8	指定运营商统计	el-icon-s-order	ServicerOrderCount			6	1	1970-01-01 00:00:00
274aca4a-db46-450c-99f7-8a6f44f1a6d0	会计中心	/accounting_center	00000000-0000-0000-0000-000000000000	会计中心	el-icon-menu	Default	accounting_center	/accounting_center/subject_manage/subject_list	7	1	1970-01-01 00:00:00
8c9f5d83-fb7b-4f21-8b60-1ed3eb93ff41	管理授权	/manage_auth	222d3e77-c07d-4c30-a384-b0c493a3cb38	管理授权	el-icon-date	ManageAuth			5	1	1970-01-01 00:00:00
99b01f22-dee8-4f63-ab8e-4b5083d38044	优惠卷类型	/cupon_type	4e52b2f4-1da6-4001-8073-38f2b12d4951	优惠卷类型	el-icon-s-data	CuponType			1	0	1970-01-01 00:00:00
67948dc1-eab6-4fed-ad68-7879030be523	优惠卷明细管理	/cupon_details	4e52b2f4-1da6-4001-8073-38f2b12d4951	优惠卷明细管理	el-icon-s-data	CuponDetails			2	0	1970-01-01 00:00:00
f4f9e4b1-e4df-418b-a88d-2bd0f00fa5e5	添加优惠卷	/add_cupon	4e52b2f4-1da6-4001-8073-38f2b12d4951	添加优惠卷	el-icon-s-data	CuponAdd			3	1	1970-01-01 00:00:00
eb15be96-d809-4a03-b3e4-be03b79b55cc	优惠卷查询	/cupon_search	4e52b2f4-1da6-4001-8073-38f2b12d4951	优惠卷查询	el-icon-s-data	CuponSearch			4	0	1970-01-01 00:00:00
b57a6f39-2ec1-4065-911a-29764b00ec2b	财务管理	/financial_manage	149d7564-5e2f-4472-a637-d1db61dfdc81	财务管理	el-icon-s-data	Content			4	0	1970-01-01 00:00:00
6264bced-71da-43f1-8fdd-5a96925164cc	配置管理	/config_manage	149d7564-5e2f-4472-a637-d1db61dfdc81	配置管理	el-icon-s-claim	Content			5	0	1970-01-01 00:00:00
77788f4d-5d0e-458b-b7cb-297e907b2eff	银行参数管理	/bank_params_manage	136f2190-6f80-46d4-a124-d512e508be5e	银行参数管理	el-icon-menu	Content			6	1	1970-01-01 00:00:00
4d3344aa-04e1-4b24-80e6-ace57f7e3e6b	服务商对账单	/financial_servicer_check_count	b57a6f39-2ec1-4065-911a-29764b00ec2b	服务商对账单	el-icon-s-data	FinancialServicerCheckCount			2	0	1970-01-01 00:00:00
4a0bcbc4-1eeb-4801-9779-d51641ecd982	充值方式配置	/collect_method_config	6264bced-71da-43f1-8fdd-5a96925164cc	充值方式配置	el-icon-s-claim	CollectMethodConfig			6	0	1970-01-01 00:00:00
8ec4ba0d-da46-49e0-a349-cbaa7f13e26b	交易安全配置	/transfer_security_config	6264bced-71da-43f1-8fdd-5a96925164cc	交易安全配置	el-icon-s-claim	TransferSecurityConfig			5	0	1970-01-01 00:00:00
d470cd75-937c-4229-b180-7b0ace983442	收益明细查询	/servicer_profit_ledger_list	996af46e-f32c-4b26-8109-2b942d1aa8b8	收益明细查询	el-icon-s-order	ServicerProfitLedgerList			8	0	1970-01-01 00:00:00
84a7f310-51f3-4799-a737-d933083f96a2	服务商对账统计	/financial_checking_count	b57a6f39-2ec1-4065-911a-29764b00ec2b	服务商对账统计	el-icon-s-data	FinancialCheckCount			1	1	1970-01-01 00:00:00
5b63f7e6-6bf1-406d-8895-b6f92418520c	功能添加	/func_manage_add	6264bced-71da-43f1-8fdd-5a96925164cc	功能添加	el-icon-s-grid	FuncManageAdd			10	1	1970-01-01 00:00:00
ed6cdd20-4817-464c-99b5-f7a2084f6a9c	平台盈利统计	/profit_statistics	b57a6f39-2ec1-4065-911a-29764b00ec2b	平台盈利统计	el-icon-date	ProfitStatistics			6	0	1970-01-01 00:00:00
c174cff5-6670-4001-bf92-8395eb9e2656	提现记录	/profit_take_log	b57a6f39-2ec1-4065-911a-29764b00ec2b	提现记录	el-icon-s-fold	ProfitTakeLog			7	1	1970-01-01 00:00:00
d344cbbc-eb28-446a-8b92-38bd340cf3b2	存款配置	/save_money_config	6264bced-71da-43f1-8fdd-5a96925164cc	存款配置	el-icon-s-claim	SaveMoneyConfig			2	0	1970-01-01 00:00:00
559d6495-5bab-428d-82d9-4bc1b30e8e06	转账配置	/transfer_config	6264bced-71da-43f1-8fdd-5a96925164cc	转账配置	el-icon-s-claim	TransferConfig			3	0	1970-01-01 00:00:00
2f4bda42-a624-4eef-9cc0-aa581733206a	虚拟账户交易流水	/trading_account_log	a8de69fd-5856-42cc-8812-fef457a99079	虚拟账户交易流水	el-icon-s-unfold	TradingAccountLog			17	0	1970-01-01 00:00:00
22028342-5177-4425-8465-60d41a46eb7e	面对面取款	/withdraw_config	6264bced-71da-43f1-8fdd-5a96925164cc	面对面取款	el-icon-s-claim	WithdrawConfig			1	0	1970-01-01 00:00:00
35d43dd0-6e55-4bcc-ba62-efc1c4cdff5d	手机号取款	/phone_withdraw	6264bced-71da-43f1-8fdd-5a96925164cc	手机号取款	el-icon-s-claim	PhoneWithdraw			1	0	1970-01-01 00:00:00
4e52b2f4-1da6-4001-8073-38f2b12d4951	优惠卷管理	/cupon	149d7564-5e2f-4472-a637-d1db61dfdc81	优惠卷管理	el-icon-s-order	Content			3	1	1970-01-01 00:00:00
f938a823-cd1f-4702-9163-45a5e016e2cf	平台收款账户管理	/collection_management_list	b57a6f39-2ec1-4065-911a-29764b00ec2b	平台收款账户管理	el-icon-s-custom	CollectionManagementList			4	0	1970-01-01 00:00:00
4a8285de-348e-4845-b0eb-142dd46002f6	客服管理	/service_config	149d7564-5e2f-4472-a637-d1db61dfdc81	客服管理	el-icon-s-data	Content			6	0	1970-01-01 00:00:00
439becc5-2668-4f15-8c7f-1b544989df47	咨询管理	/consultation_config	4a8285de-348e-4845-b0eb-142dd46002f6	咨询管理	el-icon-s-custom	ConsultationConfig			2	0	1970-01-01 00:00:00
b89d0f67-90fe-4ca8-9122-65eda7c53aa2	用户协议与隐私管理	/agreement	d9d1bfbd-1d8d-467c-b95a-5c9112b52e3a	用户协议与隐私管理	el-icon-document	Agreement			5	0	1970-01-01 00:00:00
d212c732-5c62-40ee-b1cc-184eddab4a61	汇率设置	/exchange_rate_config	6264bced-71da-43f1-8fdd-5a96925164cc	汇率设置	el-icon-s-claim	ExchangeRateConfig			4	0	1970-01-01 00:00:00
f93f8fe1-b61f-4cd9-acf6-ff7d6fd5fd60	提现方式配置	/fetch_method_config	6264bced-71da-43f1-8fdd-5a96925164cc	提现方式配置	el-icon-s-claim	FetchMethodConfig			7	0	1970-01-01 00:00:00
8d94bbc5-5eb8-4d3d-890f-3a3948b29ce9	多语言配置	/multi_language_list	6264bced-71da-43f1-8fdd-5a96925164cc	多语言配置	el-icon-s-claim	MultiLanguageList			8	0	1970-01-01 00:00:00
040a0e65-262b-44ce-9a39-336f08560ad4	钱包功能	/func_manage	6264bced-71da-43f1-8fdd-5a96925164cc	钱包功能	el-icon-s-claim	FuncManage			9	0	1970-01-01 00:00:00
3ab532fe-274c-4755-b8c4-8055ea124209	服务商充值	/submit_to_headquarters	a8de69fd-5856-42cc-8812-fef457a99079	服务商充值	el-icon-date	SubmitToHeadquarters			1	0	1970-01-01 00:00:00
63593179-f9d6-4580-97ca-dd8670688272	服务商提现	/req_money	a8de69fd-5856-42cc-8812-fef457a99079	服务商提现	el-icon-magic-stick	ReqMoney			12	0	1970-01-01 00:00:00
18d4ae82-68bc-435c-9e72-bfef6bc6b2e9	收款订单	/collection_orders	a8de69fd-5856-42cc-8812-fef457a99079	收款订单	el-icon-s-claim	CollectionOrders			18	1	1970-01-01 00:00:00
6f855990-13d0-43b1-87e6-f0647e99c360	兑换订单	/exchange_orders	a8de69fd-5856-42cc-8812-fef457a99079	兑换订单	el-icon-picture-outline-round	ExchangeOrders			13	0	1970-01-01 00:00:00
ed0c5e11-eddb-48ff-84f9-6730d389f01a	存款订单	/save_orders	a8de69fd-5856-42cc-8812-fef457a99079	存款订单	el-icon-download	SaveOrders			14	0	1970-01-01 00:00:00
d141790b-1036-447d-b4cf-ae8b75de2190	转账订单	/transfer_orders	a8de69fd-5856-42cc-8812-fef457a99079	转账订单	el-icon-s-claim	TransferOrders			16	0	1970-01-01 00:00:00
a0e0bacb-c1d2-488f-aa5f-1af550bd97a8	取款订单	/fetch_orders	a8de69fd-5856-42cc-8812-fef457a99079	取款订单	el-icon-upload2	FetchOrders			15	0	1970-01-01 00:00:00
48d6de3b-4af6-4ea3-b796-aca905fa07f3	版本管理	/app_version_count	6264bced-71da-43f1-8fdd-5a96925164cc	版本管理	el-icon-s-home	AppVersionCount			12	0	1970-01-01 00:00:00
8876a714-7d09-43d4-88ea-795c559a72a1	版本列表	/app_version_list	6264bced-71da-43f1-8fdd-5a96925164cc	版本列表	el-icon-menu	AppVersionList			13	1	1970-01-01 00:00:00
3b618a8e-5499-49a3-a980-dab3944a50d4	店员列表	/cashier_list	996af46e-f32c-4b26-8109-2b942d1aa8b8	店员列表	el-icon-user-solid	CashierList			9	1	1970-01-01 00:00:00
87c1d68e-93a3-4ee9-b0a6-30989584c056	服务商渠道配置	/pos_channel_config	996af46e-f32c-4b26-8109-2b942d1aa8b8	服务商渠道配置	el-icon-date	PosChannelConfig			10	0	1970-01-01 00:00:00
1da7c457-afd1-484c-b24c-6d984baecdbf	渠道仓库	/channel_config	6264bced-71da-43f1-8fdd-5a96925164cc	渠道仓库	el-icon-s-claim	ChannelConfig			11	0	1970-01-01 00:00:00
fcc9c518-48c0-4c0e-b96a-8d93b165d948	帮助管理	/common_help_count	4a8285de-348e-4845-b0eb-142dd46002f6	帮助管理	el-icon-s-comment	CommonHelpCount			1	0	1970-01-01 00:00:00
135baf4c-f8cd-4597-bdb5-a49f4bb53e42	帮助列表	/common_help	4a8285de-348e-4845-b0eb-142dd46002f6	帮助列表	el-icon-s-comment	CommonHelp			3	1	1970-01-01 00:00:00
\.


--
-- Data for Name: vaccount; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.vaccount (vaccount_no, account_no, va_type, balance, create_time, is_delete, use_status, delete_time, update_time, frozen_balance, balance_type, modify_time) FROM stdin;
f6d550e1-7bf5-41eb-abdf-04f8739376a7	d4a4ca0f-973e-484a-80b7-c40187aeda3f	5	0	2020-05-08 10:22:31.638952	0	1	\N	\N	0	usd	\N
8cde9f2a-fc1b-4aeb-a5ad-c0e718a3042c	d4a4ca0f-973e-484a-80b7-c40187aeda3f	1	9990	2020-05-08 10:22:31.625795	0	1	\N	\N	0	usd	\N
cf581976-fe0a-4934-85f7-71729deab9be	22222222-2222-2222-2222-222222222222	11	428	2020-04-24 05:53:42.223499	0	1	\N	\N	0	usd	\N
0a9b5da3-e6a8-4526-9002-782dddd12127	d4a4ca0f-973e-484a-80b7-c40187aeda3f	11	0	2020-05-08 10:22:31.694078	0	1	\N	\N	0	usd	\N
cb34afc3-f848-4c95-910b-5a55f100f9bf	d4a4ca0f-973e-484a-80b7-c40187aeda3f	7	-90000	2020-05-08 10:21:33.420186	0	1	\N	\N	0	usd	2020-05-08 10:22:31.647155
0c8f8a8a-3df1-4742-9525-14c61c3a655e	d4a4ca0f-973e-484a-80b7-c40187aeda3f	6	0	2020-05-08 10:22:54.939825	0	1	\N	\N	0	khr	\N
cd446382-79cf-4cd2-badd-35bb6d1d5c47	d4a4ca0f-973e-484a-80b7-c40187aeda3f	2	0	2020-05-08 10:22:54.929877	0	1	\N	\N	0	khr	\N
ec808b18-cb01-4c80-9e08-10c6ac077623	22222222-2222-2222-2222-222222222222	12	200	2020-05-07 10:29:07.773315	0	1	\N	\N	0	khr	\N
b651303e-dfbc-4c5f-bc01-dae033711144	d4a4ca0f-973e-484a-80b7-c40187aeda3f	12	0	2020-05-08 10:22:54.976284	0	1	\N	\N	0	khr	\N
c5157dbc-764b-4256-8e3c-9c11e559c724	d4a4ca0f-973e-484a-80b7-c40187aeda3f	8	-999900	2020-05-08 10:21:53.483423	0	1	\N	\N	0	khr	2020-05-08 10:22:54.947218
\.


--
-- Data for Name: wf_proc_running; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wf_proc_running (running_no, process_no, current_step, create_time, run_status) FROM stdin;
\.


--
-- Data for Name: wf_process; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wf_process (process_no, process_name, execute_no, create_time, execute_status, steps) FROM stdin;
\.


--
-- Data for Name: wf_step; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wf_step (step_no, step_name, func_name, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, create_time, is_delete) FROM stdin;
\.


--
-- Data for Name: writeoff; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.writeoff (code, income_order_no, outgo_order_no, create_time, finish_time, use_status, transfer_order_no, send_phone, recv_phone, modify_time) FROM stdin;
7856277982	2020043015152267531776		2020-04-30 14:15:22.390277	1970-01-01 00:00:00	1		886716519	85921025	\N
4636941226	2020050711273131774598		2020-05-07 10:27:31.723308	1970-01-01 00:00:00	1		886716519	85921025	\N
2953753781	2020050711290709162496		2020-05-07 10:29:07.715479	1970-01-01 00:00:00	1		886716519	85921025	\N
\.


--
-- Data for Name: xlsx_file_log; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.xlsx_file_log (xlsx_task_no, create_time, account_no, file_type, query_str, role_type) FROM stdin;
2020043012125082157732	2020-04-30 11:12:50.976964	af8ca79c-995e-4336-a1f7-31a76613d300	2		\N
\.


--
-- Name: servicer_count_list_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.servicer_count_list_id_seq', 189, true);


--
-- Name: account_collect account_collect_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account_collect
    ADD CONSTRAINT account_collect_pkey PRIMARY KEY (account_collect_no);


--
-- Name: account account_gen_key_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_gen_key_key UNIQUE (gen_key);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_pkey PRIMARY KEY (uid);


--
-- Name: adminlog adminlog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.adminlog
    ADD CONSTRAINT adminlog_pkey PRIMARY KEY (log_uid);


--
-- Name: agreement agreement_privacy_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.agreement
    ADD CONSTRAINT agreement_privacy_pkey PRIMARY KEY (id);


--
-- Name: app_version_file_log app_file_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_version_file_log
    ADD CONSTRAINT app_file_log_pkey PRIMARY KEY (id);


--
-- Name: billing_details_results billing_details_results_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.billing_details_results
    ADD CONSTRAINT billing_details_results_pkey PRIMARY KEY (bill_no);


--
-- Name: cashier cashier_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cashier
    ADD CONSTRAINT cashier_pkey PRIMARY KEY (uid);


--
-- Name: channel channel_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channel
    ADD CONSTRAINT channel_pkey PRIMARY KEY (channel_no);


--
-- Name: common_help common_help_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.common_help
    ADD CONSTRAINT common_help_pkey PRIMARY KEY (help_no);


--
-- Name: consultation_config consultation_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.consultation_config
    ADD CONSTRAINT consultation_config_pkey PRIMARY KEY (id);


--
-- Name: cust cust_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cust
    ADD CONSTRAINT cust_pkey PRIMARY KEY (cust_no);


--
-- Name: dict_acc_title dict_acc_title_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_acc_title
    ADD CONSTRAINT dict_acc_title_pkey PRIMARY KEY (title_no);


--
-- Name: dict_account_type dict_account_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_account_type
    ADD CONSTRAINT dict_account_type_pkey PRIMARY KEY (account_type);


--
-- Name: dict_bankname dict_bankname_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_bankname
    ADD CONSTRAINT dict_bankname_pkey PRIMARY KEY (bank_name);


--
-- Name: dict_bin_bankname dict_bin_bankname_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_bin_bankname
    ADD CONSTRAINT dict_bin_bankname_pkey PRIMARY KEY (bin_code);


--
-- Name: dict_images dict_images_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_images
    ADD CONSTRAINT dict_images_pkey PRIMARY KEY (image_id);


--
-- Name: dict_org_abbr dict_org_abbr_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_org_abbr
    ADD CONSTRAINT dict_org_abbr_pkey PRIMARY KEY (org_code);


--
-- Name: dict_province dict_province_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_province
    ADD CONSTRAINT dict_province_pkey PRIMARY KEY (province_code);


--
-- Name: dict_region_bank dict_region_bank_copy1_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_region_bank
    ADD CONSTRAINT dict_region_bank_copy1_pkey PRIMARY KEY (code);


--
-- Name: dict_vatype dict_vaccount_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dict_vatype
    ADD CONSTRAINT dict_vaccount_pkey PRIMARY KEY (va_type);


--
-- Name: exchange_order exchange_order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_order
    ADD CONSTRAINT exchange_order_pkey PRIMARY KEY (log_no);


--
-- Name: func_config func_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.func_config
    ADD CONSTRAINT func_config_pkey PRIMARY KEY (func_no);


--
-- Name: gen_code gen_code_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gen_code
    ADD CONSTRAINT gen_code_pkey PRIMARY KEY (gen_key);


--
-- Name: global_param global_param_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_param
    ADD CONSTRAINT global_param_pkey PRIMARY KEY (param_key);


--
-- Name: headquarters_profit_withdraw headquarters_profit_withdraw_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.headquarters_profit_withdraw
    ADD CONSTRAINT headquarters_profit_withdraw_pkey PRIMARY KEY (order_no);


--
-- Name: income_order income_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_order
    ADD CONSTRAINT income_log_pkey PRIMARY KEY (log_no);


--
-- Name: servicer_terminal income_terminal_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer_terminal
    ADD CONSTRAINT income_terminal_pkey PRIMARY KEY (terminal_no);


--
-- Name: outgo_type income_type_copy1_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.outgo_type
    ADD CONSTRAINT income_type_copy1_pkey PRIMARY KEY (outgo_type);


--
-- Name: income_type income_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_type
    ADD CONSTRAINT income_type_pkey PRIMARY KEY (income_type);


--
-- Name: lang lang_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.lang
    ADD CONSTRAINT lang_pkey PRIMARY KEY (key);


--
-- Name: log_account log_account_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_account
    ADD CONSTRAINT log_account_pkey PRIMARY KEY (log_no);


--
-- Name: log_account_web log_account_web_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_account_web
    ADD CONSTRAINT log_account_web_pkey PRIMARY KEY (log_no);


--
-- Name: log_card log_card_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_card
    ADD CONSTRAINT log_card_pkey PRIMARY KEY (log_no);


--
-- Name: log_login log_login_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_login
    ADD CONSTRAINT log_login_pkey PRIMARY KEY (log_no);


--
-- Name: settle_vaccount_balance_hourly log_settle_vaccount_balance_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.settle_vaccount_balance_hourly
    ADD CONSTRAINT log_settle_vaccount_balance_pkey PRIMARY KEY (log_no, vaccount_no);


--
-- Name: log_to_headquarters log_to_headquarters_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_to_headquarters
    ADD CONSTRAINT log_to_headquarters_pkey PRIMARY KEY (log_no);


--
-- Name: log_to_servicer log_to_servicer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_to_servicer
    ADD CONSTRAINT log_to_servicer_pkey PRIMARY KEY (log_no);


--
-- Name: log_vaccount log_vaccount_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.log_vaccount
    ADD CONSTRAINT log_vaccount_pkey PRIMARY KEY (log_no);


--
-- Name: login_token login_token_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.login_token
    ADD CONSTRAINT login_token_pkey PRIMARY KEY (acc_no);


--
-- Name: servicer_img merchant_img_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer_img
    ADD CONSTRAINT merchant_img_pkey PRIMARY KEY (servicer_img_no);


--
-- Name: servicer merchant_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer
    ADD CONSTRAINT merchant_pkey PRIMARY KEY (servicer_no);


--
-- Name: outgo_order outgo_order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.outgo_order
    ADD CONSTRAINT outgo_order_pkey PRIMARY KEY (log_no);


--
-- Name: role permission_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT permission_tag_pkey PRIMARY KEY (role_no);


--
-- Name: platform_config platform_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.platform_config
    ADD CONSTRAINT platform_config_pkey PRIMARY KEY (account_uid);


--
-- Name: rela_acc_iden rela_acc_iden_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rela_acc_iden
    ADD CONSTRAINT rela_acc_iden_pkey PRIMARY KEY (account_no, account_type);


--
-- Name: rela_imei_pubkey rela_imei_pubkey_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rela_imei_pubkey
    ADD CONSTRAINT rela_imei_pubkey_pkey PRIMARY KEY (rela_no);


--
-- Name: rela_role_url rela_role_url_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rela_role_url
    ADD CONSTRAINT rela_role_url_pkey PRIMARY KEY (rela_uid);


--
-- Name: rela_account_role rela_tag_access_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rela_account_role
    ADD CONSTRAINT rela_tag_access_pkey PRIMARY KEY (rela_uid);


--
-- Name: headquarters_profit servicefee_order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.headquarters_profit
    ADD CONSTRAINT servicefee_order_pkey PRIMARY KEY (log_no);


--
-- Name: card servicer_card_pack_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.card
    ADD CONSTRAINT servicer_card_pack_pkey PRIMARY KEY (card_no);


--
-- Name: servicer_count_list servicer_date_currency; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer_count_list
    ADD CONSTRAINT servicer_date_currency UNIQUE (servicer_no, currency_type, dates);


--
-- Name: CONSTRAINT servicer_date_currency ON servicer_count_list; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON CONSTRAINT servicer_date_currency ON public.servicer_count_list IS '服务商同一天同一币种只能有一条记录';


--
-- Name: servicer_count servicer_no_balance_type_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer_count
    ADD CONSTRAINT servicer_no_balance_type_key UNIQUE (servicer_no, currency_type);


--
-- Name: servicer_profit_ledger servicer_profit_ledger_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.servicer_profit_ledger
    ADD CONSTRAINT servicer_profit_ledger_pkey PRIMARY KEY (log_no);


--
-- Name: settle_servicer_hourly settle_servicer_hourly_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.settle_servicer_hourly
    ADD CONSTRAINT settle_servicer_hourly_pkey PRIMARY KEY (log_no);


--
-- Name: sms_send_record sms_send_record_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sms_send_record
    ADD CONSTRAINT sms_send_record_pkey PRIMARY KEY (id);


--
-- Name: collection_order transfer_order_copy1_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collection_order
    ADD CONSTRAINT transfer_order_copy1_pkey PRIMARY KEY (log_no);


--
-- Name: transfer_order transfer_order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transfer_order
    ADD CONSTRAINT transfer_order_pkey PRIMARY KEY (log_no);


--
-- Name: url url_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.url
    ADD CONSTRAINT url_pkey PRIMARY KEY (url_uid);


--
-- Name: vaccount vaccount_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.vaccount
    ADD CONSTRAINT vaccount_pkey PRIMARY KEY (vaccount_no);


--
-- Name: app_version version_vsType_system_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_version
    ADD CONSTRAINT "version_vsType_system_key" UNIQUE (version, vs_type, system);


--
-- Name: wf_proc_running wf_proc_running_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wf_proc_running
    ADD CONSTRAINT wf_proc_running_pkey PRIMARY KEY (running_no);


--
-- Name: wf_process wf_process_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wf_process
    ADD CONSTRAINT wf_process_pkey PRIMARY KEY (process_no);


--
-- Name: wf_step wf_step_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wf_step
    ADD CONSTRAINT wf_step_pkey PRIMARY KEY (step_no);


--
-- Name: writeoff writeoff_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.writeoff
    ADD CONSTRAINT writeoff_pkey PRIMARY KEY (code);


--
-- Name: xlsx_file_log xlsx_task_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.xlsx_file_log
    ADD CONSTRAINT xlsx_task_pkey PRIMARY KEY (xlsx_task_no);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM rdsadmin;
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

