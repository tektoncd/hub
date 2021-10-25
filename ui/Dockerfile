FROM node:16-alpine3.14 as BUILD
WORKDIR /app

# install deps and the copy the src to speed up docker build
COPY package-lock.json package.json /app/
RUN npm install

COPY tsconfig.json /app/
COPY public /app/public/
COPY src /app/src/
RUN npm run build

# Stage 2 - the production environment
FROM nginxinc/nginx-unprivileged
COPY --from=BUILD /app/build /usr/share/nginx/html
COPY image/start.sh /usr/bin/

USER root
RUN chmod ugo+rw /usr/share/nginx/html/config.js  && \
    chown nginx:nginx /usr/share/nginx/html/config.js && \
    chmod +x /usr/bin/start.sh
USER nginx

EXPOSE 8080

COPY image/nginx.conf /etc/nginx/conf.d/default.conf

CMD /usr/bin/start.sh
