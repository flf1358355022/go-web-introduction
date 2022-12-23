DROP TABLE IF EXISTS `user_test`;
CREATE TABLE `user_test` (
                             `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户主键id',
                             `name` varchar(50) NOT NULL COMMENT '用户姓名',
                             `phone` varchar(50) NOT NULL COMMENT '用户手机号',
                             PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

insert into user_test(`name`, `phone`) values ("Mary", "13111110000");
insert into user_test(`name`, `phone`) values ("Mr Li", "13989890000");