# Imagen base para tu aplicación Go
FROM golang:1.21.4

# Directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar todos los archivos del proyecto al contenedor
COPY . .

# Instalar MySQL client (si se requiere para realizar migraciones)
#RUN apt-get update && apt-get install -y mysql-client

# Comando para ejecutar tu aplicación al iniciar el contenedor
CMD ["go", "run", "main.go"]
