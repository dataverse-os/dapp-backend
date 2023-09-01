FROM ceramicnetwork/composedb:2.35.0

COPY cmd.sh /app/

RUN apk add jq bash moreutils

CMD ["bash", "/app/cmd.sh"]