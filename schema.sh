#!/bin/bash

# 简单的 GraphQL Schema 导出脚本

echo "Vista GraphQL Schema"
echo "===================="
echo
cat graph/schema.graphqls

# 可选：导出为 schema.graphql 文件
if [ "$1" = "export" ]; then
    cp graph/schema.graphqls schema.graphql
    echo
    echo "Schema 已导出到: schema.graphql"
fi
