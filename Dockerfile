FROM alpine

WORKDIR /bin/eq

COPY ./out/eq .

CMD ["./eq"]
