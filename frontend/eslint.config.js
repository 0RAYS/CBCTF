import js from '@eslint/js';
import globals from 'globals';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import prettier from 'eslint-plugin-prettier';

export default [
    { ignores: ['dist'] },
    {
        files: ['src/**/*.{js,jsx}'],
        languageOptions: {
            ecmaVersion: 2020,
            globals: globals.browser,
            parserOptions: {
                ecmaVersion: 'latest',
                ecmaFeatures: { jsx: true },
                sourceType: 'module',
            },
        },
        settings: { react: { version: '18.3' } },
        plugins: {
            react,
            'react-hooks': reactHooks,
            'react-refresh': reactRefresh,
            prettier,
        },
        rules: {
            ...js.configs.recommended.rules,
            ...react.configs.recommended.rules,
            ...react.configs['jsx-runtime'].rules,
            ...reactHooks.configs.recommended.rules,
            'react/jsx-no-target-blank': 'off',
            'react-refresh/only-export-components': [
                'warn',
                { allowConstantExport: true },
            ],
            // 更宽松的规则
            'react/prop-types': 'off', // 关闭 prop-types 校验
            'react/jsx-pascal-case': 'off', // 允许 JSX 标签不使用 PascalCase
            'react/no-unescaped-entities': 'off', // 允许在 JSX 中直接使用 ' 和 " 字符
            'react/react-in-jsx-scope': 'off', // React 17+ 不需要再显式引入 React
            'no-console': 'warn', // 允许使用 console, 但会提示警告
            'no-debugger': 'warn', // 允许 debugger, 但会提示警告

            // Prettier 相关规则
            'prettier/prettier': [
                'warn',
                {
                    singleQuote: true, // 使用单引号
                    trailingComma: 'es5', // 对象和数组元素后添加逗号
                    semi: true, // 强制分号
                    tabWidth: 2, // 每个缩进为 4 个空格
                    printWidth: 120, // 每行最大字符数
                    arrowParens: 'always', // 总是使用箭头函数括号
                    "endOfLine": "auto",
                },
            ],
            // 如果希望 React Hooks 更宽松, 可以禁用这两个规则
            'react-hooks/rules-of-hooks': 'off',
            'react-hooks/exhaustive-deps': 'off',
            'react-hooks/set-state-in-effect': 'off',

            // 如果有使用 TypeScript 的话, 可以适当放宽这些规则
            'no-unused-vars': ['warn'], // 仅发出警告, 允许有未使用的变量
        },
        linterOptions: {
            reportUnusedDisableDirectives: true, // 将该选项移动到这里
        },
    },
];
