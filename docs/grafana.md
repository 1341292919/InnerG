## 配置Grafana监控
1. 秒级别的QPS
2. 5分钟的QPS总和
```json
{
  "title": "QPS 实时监控（秒级）",
  "uid": "qps-sec-monitor",
  "version": 1,
  "timezone": "browser",
  "schemaVersion": 36,
  "refresh": "5s",
  "time": {
    "from": "now-5m",
    "to": "now"
  },
  "panels": [
    {
      "id": 1,
      "title": "实时 QPS（每秒请求数）",
      "type": "timeseries",
      "gridPos": {
        "h": 10,
        "w": 16,
        "x": 0,
        "y": 0
      },
      "targets": [
        {
          "expr": "gin_qps_current",
          "legendFormat": "QPS",
          "refId": "A",
          "interval": "1s",
          "intervalFactor": 1
        }
      ],
      "options": {
        "tooltip": {
          "mode": "single",
          "sort": "none"
        },
        "legend": {
          "displayMode": "list",
          "placement": "bottom",
          "calcs": ["last", "max", "mean"]
        }
      },
      "fieldConfig": {
        "defaults": {
          "unit": "short",
          "decimals": 0,
          "custom": {
            "lineWidth": 2,
            "fillOpacity": 10,
            "showPoints": "always",
            "pointSize": 5,
            "lineInterpolation": "step"
          }
        },
        "overrides": []
      },
      "interval": "1s"
    },
    {
      "id": 2,
      "title": "5分钟 QPS 总数",
      "type": "timeseries",
      "gridPos": {
        "h": 10,
        "w": 16,
        "x": 0,
        "y": 10
      },
      "targets": [
        {
          "expr": "gin_qps_sum_5min",
          "legendFormat": "5分钟总请求数",
          "refId": "A",
          "interval": "5m",
          "intervalFactor": 1
        }
      ],
      "options": {
        "tooltip": {
          "mode": "single",
          "sort": "none"
        },
        "legend": {
          "displayMode": "list",
          "placement": "bottom",
          "calcs": ["last", "max", "mean"]
        }
      },
      "fieldConfig": {
        "defaults": {
          "unit": "short",
          "decimals": 0,
          "custom": {
            "lineWidth": 2,
            "fillOpacity": 10,
            "showPoints": "always",
            "pointSize": 5,
            "lineInterpolation": "step"
          }
        },
        "overrides": []
      },
      "interval": "5m"
    },
    {
      "id": 3,
      "title": "当前 QPS 实时值",
      "type": "stat",
      "gridPos": {
        "h": 4,
        "w": 4,
        "x": 16,
        "y": 0
      },
      "targets": [
        {
          "expr": "gin_qps_current",
          "legendFormat": "当前 QPS",
          "refId": "A",
          "interval": "1s"
        }
      ],
      "options": {
        "reduceOptions": {
          "values": false,
          "calcs": ["last"]
        },
        "textMode": "value_and_name"
      },
      "fieldConfig": {
        "defaults": {
          "unit": "short",
          "decimals": 0,
          "color": {
            "mode": "thresholds"
          },
          "thresholds": {
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "yellow",
                "value": 500
              },
              {
                "color": "red",
                "value": 1000
              }
            ]
          }
        }
      }
    },
    {
      "id": 4,
      "title": "5分钟总 QPS",
      "type": "stat",
      "gridPos": {
        "h": 4,
        "w": 4,
        "x": 20,
        "y": 0
      },
      "targets": [
        {
          "expr": "gin_qps_sum_5min",
          "legendFormat": "5分钟总请求数",
          "refId": "A",
          "interval": "5m"
        }
      ],
      "options": {
        "reduceOptions": {
          "values": false,
          "calcs": ["last"]
        },
        "textMode": "value_and_name"
      },
      "fieldConfig": {
        "defaults": {
          "unit": "short",
          "decimals": 0,
          "color": {
            "mode": "thresholds"
          },
          "thresholds": {
            "steps": [
              {
                "color": "blue",
                "value": null
              }
            ]
          }
        }
      }
    }
  ]
}
```