FROM nginx:alpine

# Directorio donde Nginx sirve los archivos
WORKDIR /usr/share/nginx/html

# Configuración de Nginx (SPA + WASM + COOP/COEP + caché)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiamos SOLO los archivos del sitio (el build que ya tienes en la raíz)
# Usa .dockerignore para no copiar Dockerfile, nginx.conf, etc.
COPY . ./
