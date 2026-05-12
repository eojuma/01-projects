FROM alpine:latest
LABEL name="ascii-art-web" \
      description="Golang web application that renders text as ASCII art banners using multiple font styles" \
      maintainers="Aaron Brown, Evans Juma, Clinton Odhiambo" \
      project="Zone01 Go Curriculum" \
      language="Golang 1.24.3"
COPY ascii-art-web /app/
COPY templates /app/templates
WORKDIR /app
ENTRYPOINT ["./ascii-art-web"]

