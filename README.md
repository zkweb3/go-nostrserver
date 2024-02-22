![logo](https://socialify.git.ci/zkweb3/go-nostrserver/image?description=1&descriptionEditable=a%20sample%20of%20local%20testing%20server&forks=1&issues=1&language=1&name=1&pattern=Floating%20Cogs&pulls=1&stargazers=1&theme=Light)

## How to generate ca certificate using OpenSSL:

```bash
openssl genrsa -out cakey.pem 4096
```
```bash
openssl req -new -x509 -key cakey.pem -out cacert.pem -days 365
```