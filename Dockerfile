
FROM marcelstoer/nodemcu-build 
WORKDIR /opt/nodemcu-firmware
RUN rm -rf *
RUN git clone --recurse-submodules https://github.com/nodemcu/nodemcu-firmware.git .
RUN build

FROM golang:1.13 
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-w -extldflags "-static"' -o app 
COPY --from=0 /opt/nodemcu-firmware/luac.cross .
ENTRYPOINT /app/app