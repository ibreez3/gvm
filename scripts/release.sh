#!/bin/bash

set -e

# 确保工作区是干净的
if [[ -n $(git status --porcelain) ]]; then
    echo "Error: Working directory is not clean. Please commit or stash changes first."
    exit 1
fi

# 确保当前分支是 main
current_branch=$(git rev-parse --abbrev-ref HEAD)
if [[ "$current_branch" != "main" ]]; then
    echo "Error: You must be on the 'main' branch to release."
    exit 1
fi

# 确保本地 main 是最新的
echo "Pulling latest changes from origin/main..."
git pull origin main

# 获取最新的 tag，如果不存在则从 v0.0.0 开始
latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo "Current version: $latest_tag"

# 解析版本号 (假设格式为 vX.Y.Z)
# 去掉 'v' 前缀
version_raw=${latest_tag#v}
IFS='.' read -r major minor patch <<< "$version_raw"

# 如果解析失败（例如没有 patch 版本），给默认值
major=${major:-0}
minor=${minor:-0}
patch=${patch:-0}

# 计算建议版本
next_patch="v$major.$minor.$((patch + 1))"
next_minor="v$major.$((minor + 1)).0"
next_major="v$((major + 1)).0.0"

echo ""
echo "Select release type:"
echo "1) Patch ($next_patch) - Bug fixes / small changes"
echo "2) Minor ($next_minor) - New features (backward compatible)"
echo "3) Major ($next_major) - Breaking changes"
echo "4) Custom input"

read -p "Enter choice [1]: " choice
choice=${choice:-1}

case $choice in
    1) new_version=$next_patch ;;
    2) new_version=$next_minor ;;
    3) new_version=$next_major ;;
    4) 
        read -p "Enter version (e.g. v1.2.3): " input_version
        if [[ ! $input_version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "Error: Invalid format. Must be vX.Y.Z"
            exit 1
        fi
        new_version=$input_version
        ;;
    *) 
        echo "Invalid choice"
        exit 1 
        ;;
esac

echo ""
read -p "Prepare to release $new_version. Continue? [y/N] " confirm
if [[ $confirm != [yY] && $confirm != [yY][eE][sS] ]]; then
    echo "Aborted."
    exit 0
fi

echo ""
echo "Creating tag $new_version..."
git tag -a "$new_version" -m "Release $new_version"

echo "Pushing tag to origin..."
git push origin "$new_version"

echo ""
echo "✅ Done! Release workflow should be triggered on GitHub."
echo "Check status at: https://github.com/ibreez3/gvm/actions"
