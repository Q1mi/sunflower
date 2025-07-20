-- 用户基本信息表
DROP TABLE IF EXISTS `userinfo`;

CREATE TABLE `userinfo` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` CHAR(32) NOT NULL COMMENT '用户密码(MD5加密)',
    `email` VARCHAR(100) COMMENT '用户邮箱',
    `avatar` VARCHAR(256) COMMENT '用户头像',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户基本信息表';

-- 用户积分汇总表
DROP TABLE IF EXISTS `user_points`;

CREATE TABLE `user_points` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `points` BIGINT DEFAULT 0 COMMENT '当前可用积分',
    `points_total` BIGINT DEFAULT 0 COMMENT '累计获得积分',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户积分汇总表';

-- 用户积分交易明细表
DROP TABLE IF EXISTS `user_points_transactions`;
CREATE TABLE `user_points_transactions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `points_change` BIGINT NOT NULL COMMENT '积分变动值 (正数为增加，负数为扣除)',
    `current_balance` BIGINT NOT NULL COMMENT '当前余额',
    `transaction_type` TINYINT NOT NULL COMMENT '交易类型(1:签到 2:连续签到 3:补签 4:每日任务 5:福利任务)',
    `description` VARCHAR(100) COMMENT '积分变动说明',
    `ext_json` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '扩展字段',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户积分明细表';

-- 月度奖励记录表
DROP TABLE IF EXISTS `user_monthly_bonus_log`;
CREATE TABLE `user_monthly_bonus_log` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `year_month` CHAR(6) NOT NULL COMMENT '年月（YYYYMM）',
    `bonus_type` TINYINT NOT NULL COMMENT '奖励类型（1:连续签到3天 2:连续签到7天 3:连续签到15天 4:月满签）',
    `description` VARCHAR(100) COMMENT '积分变动说明',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_yearmonth_type` (`user_id`, `year_month`, `bonus_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户月度连续签到奖励记录表';

-- 签到记录表
DROP TABLE IF EXISTS `user_checkin_records`;

CREATE TABLE `user_checkin_records` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `checkin_date` DATE NOT NULL COMMENT '签到日期',
    `checkin_type` TINYINT NOT NULL DEFAULT 1 COMMENT '签到类型：1=正常签到，2=补签',
    `points_awarded_base` INT NOT NULL DEFAULT 1 COMMENT '获得积分',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_date` (`user_id`, `checkin_date`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_checkin_date` (`checkin_date`),
    KEY `idx_checkin_type` (`checkin_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '签到记录表';