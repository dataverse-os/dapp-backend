FROM node:16-alpine

COPY cmd.sh /app/

RUN apk add jq bash moreutils \
    && npm install --location=global @composedb/cli

CMD ["bash", "/app/cmd.sh"]