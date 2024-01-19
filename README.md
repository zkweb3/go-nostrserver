openssl genrsa -out cakey.pem 4096

openssl req -new -x509 -key cakey.pem -out cacert.pem -days 365