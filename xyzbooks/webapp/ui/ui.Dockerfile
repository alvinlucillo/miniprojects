# build stage
FROM node:lts-alpine as build-stage

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

# production stage
FROM nginx:stable-alpine as production-stage

# Install dockerize
ENV DOCKERIZE_VERSION v0.6.1
RUN apk add --no-cache openssl
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz


COPY --from=build-stage /app/dist/spa /usr/share/nginx/html
COPY --from=build-stage /app/config/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD dockerize -wait http://$SERVER_HOST:$SERVER_PORT/ready -timeout 20s && nginx -g 'daemon off;'