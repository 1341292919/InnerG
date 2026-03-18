CREATE TABLE InnerG.user
(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    account VARCHAR(32)  NULL,
    username VARCHAR(64) NOT NULL,
    email VARCHAR(128) NULL,
    avatar VARCHAR(512) NULL,
    status TINYINT DEFAULT 1 NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    gender TINYINT DEFAULT 0 COMMENT '性别：0-未知，1-男，2-女',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE KEY (account),
    UNIQUE KEY (email)
) AUTO_INCREMENT = 108000;

INSERT INTO InnerG.user (id, account, username, email, avatar, gender, password_hash, status, created_at, updated_at, deleted_at)
VALUES (
           108000,
           'lbl102300218',
           '林柏林爱喝酒',
           '15759762783@163.com',
           'https://ts3.tc.mm.bing.net/th/id/OIP-C.MBqJCFa5AMwdgTr5_91G4gHaE7?rs=1&pid=ImgDetMain&o=7&rm=3',
           1,
           '$2a$12$vA8.XsWrL9fABqpe8bt3Q.ITabaS/efkNBBGRNrrFFkhJb6JzxfC',
           0,
           '2026-03-18 15:02:46',
           '2026-03-18 23:03:31',
           NULL
       );