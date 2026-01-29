# Используем официальный образ OpenSearch Dashboards
FROM opensearchproject/opensearch-dashboards:2.15.0

# Переключаемся на root для копирования файлов и установки (если нужно)
USER root

# Копируем собственный конфиг с Security
COPY opensearch_dashboards.yml /usr/share/opensearch-dashboards/config/

# Устанавливаем владельца обратно на opensearch-dashboards
RUN chown opensearch-dashboards:opensearch-dashboards /usr/share/opensearch-dashboards/config/opensearch_dashboards.yml

# Если нужно, можно добавить дополнительные плагины
# RUN /usr/share/opensearch-dashboards/bin/opensearch-dashboards-plugin install <plugin-url>

# Возвращаемся к пользователю Dashboards
USER opensearch-dashboards

# Стандартная команда запуска
CMD ["/usr/share/opensearch-dashboards/bin/opensearch-dashboards"]
