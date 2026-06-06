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

// 输出初始化完成信息（这些会显示在 docker logs 中）
print("✅ MongoDB 初始化成功！");
print("📦 数据库: InnerG");
print("📁 集合: ai_chat_sessions");
print("🔍 索引: sessionId, userId+updatedAt, status");
