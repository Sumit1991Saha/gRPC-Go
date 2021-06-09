#!/bin/zsh

function generateCACertificate() {
  echo "Generating CA Certificate"
  openssl ecparam -name prime256v1 -genkey -noout -out certs/ca/cakey.key
  openssl req -x509 -new -nodes -key certs/ca/cakey.key -subj "/CN=ServerMonCA/C=SM" -days 3650 -out certs/ca/cacert.pem
}

function generateServerCertificate() {
  echo "Generating Server Certificate"
  openssl ecparam -name prime256v1 -genkey -noout -out certs/server/server.key
  generateCSRConfigForServer
  openssl req -new -key certs/server/server.key -out certs/server/server.csr -config certs/server/csrserver.conf
  openssl x509 -req -in certs/server/server.csr -CA certs/ca/cacert.pem -CAkey certs/ca/cakey.key -CAcreateserial -out certs/server/server.pem -days 3650 -extfile certs/server/csrserver.conf -extensions req_ext
}

function generateClientCertificate() {
  echo "Generating Client Certificate"
  openssl ecparam -name prime256v1 -genkey -noout -out certs/client/client.key
  generateCSRConfigForClient
  openssl req -new -key certs/client/client.key -out certs/client/client.csr -config certs/client/csrclient.conf
  openssl x509 -req -in certs/client/client.csr -CA certs/ca/cacert.pem -CAkey certs/ca/cakey.key -CAcreateserial -out certs/client/client.pem -days 3650 -extfile certs/client/csrclient.conf -extensions req_ext
}

function generateCSRConfigForServer() {
  echo "Generating CSR config for server"
cat > certs/server/csrserver.conf <<EOF
  [ req ]
  default_bits = 256
  prompt = no
  default_md = sha256
  req_extensions = req_ext
  distinguished_name = dn

  [ dn ]
  C = SM
  CN = localhost

  [ req_ext ]
  keyUsage = keyEncipherment, dataEncipherment
  extendedKeyUsage = serverAuth
  subjectAltName = @alt_names

  [ alt_names ]
  DNS.1 = localhost
  IP.1 = 127.0.0.1

EOF
}

function generateCSRConfigForClient() {
  echo "Generating CSR config for Client"
cat > certs/client/csrclient.conf <<EOF
  [ req ]
  default_bits = 256
  prompt = no
  default_md = sha256
  req_extensions = req_ext
  distinguished_name = dn

  [ dn ]
  C = SM
  CN = client

  [ req_ext ]
  keyUsage = keyEncipherment
  extendedKeyUsage = clientAuth

EOF
}

function generateCerts() {
  rm -rf certs
  mkdir certs
  mkdir certs/ca/
  generateCACertificate
  mkdir certs/server/
  generateServerCertificate
  mkdir certs/client/
  generateClientCertificate
}

generateCerts
