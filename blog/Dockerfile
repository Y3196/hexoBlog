# 所需基础镜像
FROM nginx:latest
# 作者信息
MAINTAINER Felix
# 标签信息
LABEL version=999
# 端口指定
EXPOSE 8080
# 指定工作目录
WORKDIR /usr/share/nginx/html
# 文件拷贝 将打包后的前端项目 拷贝到。。。。。
COPY dist/ /usr/share/nginx/html
#拷贝nginx 配置文件到。。。。。。。
COPY nginx.conf /etc/nginx/conf.d/nginx.conf
