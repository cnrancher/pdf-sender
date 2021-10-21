FROM alpine:3
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone && apk del tzdata
COPY bin/pdf-sender /pdf-sender
COPY package/init.yml /etc/pdf-sender.yml
ENTRYPOINT ["/pdf-sender"]
CMD [ "-f", "/etc/pdf-sender.yml", "run" ]