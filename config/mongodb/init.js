// 切换到 InnerG 数据库
db = db.getSiblingDB('InnerG');

db.createUser({
    user: "InnerG",
    pwd: "InnerG",
    roles: [{ role: "readWrite", db: "InnerG" }]
});

// 创建 AI 聊天会话集合
db.createCollection('ai_chat_sessions');

// 创建必要的索引
db.ai_chat_sessions.createIndex({ "sessionId": 1 }, { unique: true });
db.ai_chat_sessions.createIndex({ "userId": 1, "updatedAt": -1 });
db.ai_chat_sessions.createIndex({ "status": 1 });

// 插入一条示例数据（可选，测试用）
// 如果不想要示例数据，可以注释掉这部分
db.ai_chat_sessions.insertOne({
    sessionId: "sess_" + ObjectId().str,
    userId: "test_user_001",
    model: "gpt-4",
    title: "欢迎使用AI助手",
    status: "active",
    messages: [
        {
            role: "system",
            message: "你是一个有帮助的AI助手",
            createdAt: new Date()
        },
        {
            role: "user",
            message: "你好，请介绍一下你自己",
            createdAt: new Date()
        }
    ],
    createdAt: new Date(),
    updatedAt: new Date()
});

// 输出初始化完成信息（这些会显示在 docker logs 中）
print("✅ MongoDB 初始化成功！");
print("📦 数据库: InnerG");
print("📁 集合: ai_chat_sessions");
print("🔍 索引: sessionId, userId+updatedAt, status");
print("✨ 示例数据: " + (db.ai_chat_sessions.count() > 0 ? "已插入" : "未插入"));