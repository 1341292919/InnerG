CREATE TABLE InnerG.user
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account VARCHAR(32) NULL,
    username VARCHAR(64) NOT NULL,
    email VARCHAR(128) NULL,
    avatar VARCHAR(512) NULL,
    status TINYINT DEFAULT 1 NOT NULL,
    role_type TINYINT NOT NULL DEFAULT 1 COMMENT '角色类型：0-管理员，1-普通用户',
    password_hash VARCHAR(255) NOT NULL,
    gender TINYINT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE KEY (account),
    UNIQUE KEY (email)
) AUTO_INCREMENT = 108000;

INSERT INTO InnerG.user (id, account, username, email, avatar, gender, role_type, password_hash, status, created_at, updated_at, deleted_at)
VALUES (
           108000,
           'lbl102300218',
           '林柏林爱喝酒',
           '15759762783@163.com',
           'https://ts3.tc.mm.bing.net/th/id/OIP-C.MBqJCFa5AMwdgTr5_91G4gHaE7?rs=1&pid=ImgDetMain&o=7&rm=3',
           1,
           1,  -- role_type: 1-普通用户
           '$2a$12$sHsUsfhcJPri5.LO0wgovuG4KFozx85tc7CQivRPGCDTz.TUplc7q',
           0,
           '2026-03-18 15:02:46',
           '2026-03-18 23:03:31',
           NULL
       );
-- 插入管理员账号
INSERT INTO InnerG.user (id, account, username, email, avatar, gender, role_type, password_hash, status, created_at, updated_at, deleted_at)
VALUES (
           108001,
           'admin_zhang',
           '张知行',  -- 管理员姓名，寓意"知行合一"
           'zhangzhixing@qq.com',
           'https://ts3.tc.mm.bing.net/th/id/OIP-C.MBqJCFa5AMwdgTr5_91G4gHaE7?rs=1&pid=ImgDetMain&o=7&rm=3',  -- 头像不变
           1,  -- 性别：男
           0,  -- role_type: 0-管理员
           '$2a$12$sHsUsfhcJPri5.LO0wgovuG4KFozx85tc7CQivRPGCDTz.TUplc7q',  -- 密码不变
           1,  -- status: 1-正常
           '2026-03-20 10:00:00',  -- 创建时间
           '2026-03-20 10:00:00',  -- 更新时间
           NULL
       );