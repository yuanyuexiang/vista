#!/bin/bash

# GraphQL Schema 工具脚本

case "$1" in
    "export")
        echo "Vista GraphQL Schema"
        echo "===================="
        echo
        cat graph/schema.graphqls
        cp graph/schema.graphqls schema.graphql
        echo
        echo "✅ Schema 已导出到: schema.graphql"
        ;;
    "introspection")
        echo "GraphQL Introspection 查询:"
        echo "================================================"
        echo "curl -X POST http://localhost:8080/wechat/query \\"
        echo "  -H \"Content-Type: application/json\" \\"
        echo "  -d '{\"query\": \"{ __schema { types { name kind } } }\"}'"
        echo
        echo "或访问 GraphQL Playground: http://localhost:8080/"
        ;;
    "help"|"-h"|"--help")
        echo "GraphQL Schema 工具"
        echo "==================="
        echo "用法："
        echo "  ./schema.sh              # 显示 schema 内容"
        echo "  ./schema.sh export       # 导出 schema.graphql 文件"
        echo "  ./schema.sh introspection # 显示 introspection 命令"
        echo "  ./schema.sh help         # 显示此帮助"
        ;;
    *)
        echo "Vista GraphQL Schema"
        echo "===================="
        echo
        cat graph/schema.graphqls
        echo
        echo "💡 提示: 使用 './schema.sh help' 查看更多选项"
        ;;
esac
