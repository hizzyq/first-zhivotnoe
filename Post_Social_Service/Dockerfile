# Этап 1: Сборка и публикация приложения
FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine AS build-env
WORKDIR /app
# Копируем файл проекта и восстанавливаем зависимости
COPY *.slnx ./
# 2. Создаем структуру папок и копируем файлы проектов (.csproj)
# Это нужно, чтобы dotnet restore закешировал зависимости
COPY Core/Core.csproj ./Core/
COPY Date/Date.csproj ./Date/
COPY vas_projects/vas_projects.csproj ./vas_projects/
RUN dotnet restore
# Копируем остальной код и собираем релизную версию
COPY . ./
RUN dotnet publish vas_projects/vas_projects.csproj -c Release -o out

# Этап 2: Запуск с использованием .NET Runtime на Alpine
FROM mcr.microsoft.com/dotnet/aspnet:10.0-alpine
WORKDIR /app
COPY --from=build-env /app/out .
EXPOSE 5000
ENTRYPOINT ["dotnet", "vas_projects.dll"]
