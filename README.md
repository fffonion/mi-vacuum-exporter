# mi-vacuum-exporter

Export Mi Home vacuum robot metrics to grafana dashboard

## Grafana dashboard

## Sample prometheus config

```yaml
# scrape vacuum devices
scrape_configs:
  - job_name: 'vacuum'
    static_configs:
    - targets:
      # IP of your vacuums
      - miio://192.168.0.233/?token=YOUR_TOKEN
    metrics_path: /scrape
    relabel_configs:
      - source_labels : [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        # don't expose your token
        regex: "miio://([^/]+)/.*"
        target_label: instance
      - target_label: __address__
        # IP of the exporter
        replacement: localhost:9234

# scrape vacuum_exporter itself
  - job_name: 'vacuum_exporter'
    static_configs:
      - targets:
        # IP of the exporter
        - localhost:9234
```

To find `token` of you device, please refer to [this guide](https://github.com/jghaanstra/com.xiaomi-miio/blob/master/docs/obtain_token.md).

## See also

- Protocol specification: [marcelrv/XiaomiRobotVacuumProtocol](https://github.com/marcelrv/XiaomiRobotVacuumProtocol)
- Python client: [rytilahti/python-miio](https://github.com/rytilahti/python-miio)
- Go client: [nickw444/miio-go](https://github.com/nickw444/miio-go)
- Go client: [vkorn/go-miio](https://github.com/vkorn/go-miio)

