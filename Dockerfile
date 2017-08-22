FROM alpine
COPY urlimport /bin
ENTRYPOINT ["urlimport"]
