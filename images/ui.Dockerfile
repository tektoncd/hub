FROM --platform=$BUILDPLATFORM node:19-alpine3.16 as BUILD
WORKDIR /app

# install deps and the copy the src to speed up docker build
COPY ui/package-lock.json ui/package.json ui/.npmrc ui/config-overrides.js /app/
RUN npm install

COPY ui/tsconfig.json /app/
COPY ui/public /app/public/
COPY ui/src /app/src/
RUN npm run build

# Stage 2 - the production environment
FROM nginxinc/nginx-unprivileged:alpine
COPY --from=BUILD /app/build /usr/share/nginx/html
COPY ui/image/start.sh /usr/bin/

USER root
RUN chmod ugo+rw /usr/share/nginx/html/config.js  && \
    chown nginx:nginx /usr/share/nginx/html/config.js && \
    chmod +x /usr/bin/start.sh
USER nginx

EXPOSE 8080

COPY ui/image/nginx.conf /etc/nginx/conf.d/default.conf
COPY ui/image/location.locations /etc/nginx/conf.d/location.locations

CMD sh /usr/bin/start.sh
