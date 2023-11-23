# Imagen base para tu aplicación Go
FROM golang:1.21.4

# Directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar todos los archivos del proyecto al contenedor
COPY . .

COPY wait-for-it.sh wait-for-it.sh
RUN chmod +x wait-for-it.sh

# Instalar MySQL client (si se requiere para realizar migraciones)
#RUN apt-get update && apt-get install -y mysql-client

# Comando para ejecutar tu aplicación al iniciar el contenedor
#CMD ["go", "run", "main.go"]



# Establece el punto de entrada de la aplicación
#ENTRYPOINT ["dotnet", "log_management_ms.dll"]
CMD /bin/bash -c "./wait-for-it.sh rabbitmq:5672 -- go run main.go"
