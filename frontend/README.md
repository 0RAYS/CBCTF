# CBCTF-frontend
CBCTF平台的前端项目

## 环境要求

请确保您的开发环境满足以下要求：
```bash
node -v  # v18.0.0 或更高
npm -v   # v9.0.0 或更高
```

## 技术栈

- React 18
- Vite 6
- Tailwind CSS 3
- React Router DOM 6
- Axios

## 包管理器
本项目使用 pnpm 进行包管理，安装命令：
```bash
npm install -g pnpm
```

## 开始使用

1. 安装依赖
```bash
pnpm install
```

2. 开发环境运行
```bash
pnpm dev
```

3. 生产环境构建
```bash
pnpm build
```

## 项目结构

```
cbctf-frontend/
├── src/
│   ├── api/          # API 请求
│   ├── components/   # 公共组件
│   ├── pages/        # 页面组件
│   ├── App.jsx       # 应用入口
│   └── main.jsx      # 主入口
├── public/           # 静态资源
└── ...配置文件
```

## 开发规范

- 使用 ESLint + Prettier 进行代码规范和格式化
- 运行 `pnpm lint` 进行代码检查和自动修复
