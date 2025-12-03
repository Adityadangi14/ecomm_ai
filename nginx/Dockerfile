FROM openresty/openresty:alpine

RUN apk add --no-cache --virtual .build-deps \
    make \
    gcc \
    musl-dev \
    lua5.1-dev \
    luarocks5.1 \
    && luarocks-5.1 install lua-resty-jwt \
    && apk del .build-deps

# Create directory for custom Lua modules
RUN mkdir -p /usr/local/openresty/nginx/lua/

# Copy your JWT auth module
# COPY lua/jwt_auth.lua /usr/local/openresty/nginx/lua/

# Copy nginx configuration
COPY nginx.conf /usr/local/openresty/nginx/conf/nginx.conf

EXPOSE 80 443

CMD ["/usr/local/openresty/bin/openresty", "-g", "daemon off;"]
