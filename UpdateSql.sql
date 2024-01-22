--更新eventlogs 字段 设置 taskcount 为整数类型
ALTER TABLE `openIM_v2`.`event_logs`
    MODIFY COLUMN `user_taskcount` int(10) NULL DEFAULT 0 AFTER `chatwithaddress`;