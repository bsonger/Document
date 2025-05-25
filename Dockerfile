###############################################################################
# 1️⃣ Build Stage —— 使用 Hugo 生成 static site                               #
###############################################################################
FROM klakegg/hugo:0.125.4-ext-alpine AS builder

# 项目根目录（与 config.toml / hugo.toml 同级）
WORKDIR /src

# 将本地 Hugo 源码复制进来
COPY . .

# 允许通过 --build-arg 设置站点 baseURL
ARG HUGO_BASEURL=/
ENV HUGO_BASEURL=${HUGO_BASEURL}

# 生成 public/ 目录（--gc 清理无用文件，--minify 压缩）
RUN hugo --baseURL "${HUGO_BASEURL}" --gc --minify

###############################################################################
# 2️⃣ Runtime Stage —— 使用 NGINX 提供静态文件                               #
###############################################################################
FROM nginx:1.25-alpine

# 清空 NGINX 默认文件
RUN rm -rf /usr/share/nginx/html/*

# 拷贝 Hugo 输出
COPY --from=builder /src/public /usr/share/nginx/html

# 如需自定义 NGINX 配置，取消注释下一行并提供 nginx.conf
# COPY nginx.conf /etc/nginx/nginx.conf

# 暴露 80 端口
EXPOSE 80
STOPSIGNAL SIGTERM

# 健康检查：每 30 秒访问一次首页
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost || exit 1

# 启动 NGINX（前台）
CMD ["nginx", "-g", "daemon off;"]
