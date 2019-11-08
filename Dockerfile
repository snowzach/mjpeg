FROM golang:1.13-alpine3.10 as builder

RUN apk add --no-cache git gcc g++ linux-headers libc-dev make cmake curl ffmpeg ffmpeg-dev && \
    rm -rf /var/cache/apk/*

# Install GOCV
ARG OPENCV_VERSION="4.1.2"
ENV OPENCV_VERSION $OPENCV_VERSION
RUN cd /tmp && \
    curl -Lo opencv.zip https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip && \
    unzip -q opencv.zip && \
    curl -Lo opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip && \
    unzip -q opencv_contrib.zip && \
    rm opencv.zip opencv_contrib.zip && \
    cd opencv-${OPENCV_VERSION} && \
    mkdir build && cd build && \
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
    -D CMAKE_INSTALL_PREFIX=/usr/local \
    -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules \
    -D WITH_JASPER=OFF \
    -D WITH_QT=OFF \
    -D WITH_GTK=OFF \
    -D BUILD_DOCS=OFF \
    -D BUILD_EXAMPLES=OFF \
    -D BUILD_TESTS=OFF \
    -D BUILD_PERF_TESTS=OFF \
    -D BUILD_opencv_java=NO \
    -D BUILD_opencv_python=NO \
    -D BUILD_opencv_python2=NO \
    -D BUILD_opencv_python3=NO \
    -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
    make -j $(nproc --all) && \
    make preinstall && make install && \
    cd /tmp && rm -rf opencv*
ENV PKG_CONFIG_PATH /usr/local/lib64/pkgconfig

WORKDIR /build
ADD . .
RUN go build -o ./mjproxy/mjproxy ./mjproxy/...

FROM alpine:3.10
RUN apk add --no-cache live-media libstdc++ ffmpeg && \
    rm -rf /var/cache/apk/*
WORKDIR /opt/mjproxy
COPY --from=builder /build/mjproxy/mjproxy /opt/mjproxy/mjproxy
COPY --from=builder /usr/local/lib64/. /usr/local/lib/.
COPY example.yaml /opt/mjproxy/mjproxy.yaml
CMD [ "/opt/mjproxy/mjproxy", "-c", "/opt/mjproxy/mjproxy.yaml" ]