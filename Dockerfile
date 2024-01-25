# FROM golang:alpine3.19 as go_builder

# FROM go_builder as backend
# RUN ["git", "clone", "https://github.com/ling734/hoppscotch-backend-go.git", "/app"]
# WORKDIR /app
# RUN ["go", "mod", "download"]
# RUN ["go", "build", "-o", "./bin/server", "./server.go"]

FROM nginx:latest
# FROM caddy:latest
WORKDIR /app
COPY template ./template/
COPY nginx.conf /etc/nginx/
# COPY aio.Caddyfile /etc/caddy/Caddyfile
COPY --chmod=755 import-meta-env .
COPY --chmod=755 server .
COPY --chmod=755 healthcheck.sh .

COPY --from=hoppscotch/hoppscotch:latest /site /site
# COPY --from=backend /app/bin/hoppscotch ./hoppscotch

RUN sed -i "s@/archive.ubuntu.com/@/mirrors.tuna.tsinghua.edu.cn/@g" /etc/apt/sources.list.d/debian.sources && apt-get update
RUN apt install -y tini curl
# RUN apt install -y nodejs npm
# RUN npm install -g @import-meta-env/cli

HEALTHCHECK --interval=2s CMD /bin/sh /app/healthcheck.sh
ENTRYPOINT [ "tini", "--" ]
# CMD printenv > build.env && npx import-meta-env -x build.env -e build.env -p "/site/**/*" && nginx && /app/server
CMD /app/import-meta-env /site & nginx & /app/server
# CMD /app/import-meta-env /site & /app/server & caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
EXPOSE 3000/tcp
EXPOSE 3100/tcp
EXPOSE 3170/tcp