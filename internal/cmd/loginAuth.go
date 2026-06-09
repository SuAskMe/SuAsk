package cmd

// 鉴权逻辑已迁移到 middleware.SessionRequired，
// 通过路由分组直接控制，不再需要 trie 匹配和排除列表。
// 本文件保留为空，避免其他地方有残留引用时编译报错。
