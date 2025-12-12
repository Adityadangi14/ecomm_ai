FROM openresty/openresty:alpine

# Install build tools + luarocks
RUN apk add --no-cache --virtual .build-deps \
    make \
    gcc \
    musl-dev \
    lua5.1-dev \
    luarocks5.1 \
    && luarocks-5.1 install lua-resty-jwt \
    && luarocks-5.1 install lua-resty-prometheus \
    && apk del .build-deps

# Create directory for custom Lua modules (optional)
RUN mkdir -p /usr/local/openresty/nginx/lua/

# Copy your nginx config
COPY nginx.conf /usr/local/openresty/nginx/conf/nginx.conf

EXPOSE 80 443

CMD ["/usr/local/openresty/bin/openresty", "-g", "daemon off;"]
