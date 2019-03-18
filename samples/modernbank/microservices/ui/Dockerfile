FROM nginx:alpine
COPY build /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/conf.d/
COPY nginx-with-proxy-pass.conf.template /etc/nginx/conf.d/
COPY entrypoint.sh /
EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
