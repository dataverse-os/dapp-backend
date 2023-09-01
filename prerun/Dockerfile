FROM ceramicnetwork/composedb:2.35.0

COPY cmd.sh /app/

RUN apt install jq -y

CMD ["bash", "/app/cmd.sh"]