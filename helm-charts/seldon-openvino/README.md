# seldon-openvino

> **:exclamation: This Helm Chart is deprecated!**

![Version: 0.1](https://img.shields.io/badge/Version-0.1-informational?style=flat-square)

Proxy integration to deploy models optimized for Intel OpenVINO in Seldon Core v1

## Source Code

* <https://github.com/SeldonIO/seldon-core>
* <https://github.com/SeldonIO/seldon-core/tree/master/helm-charts/seldon-openvino>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| engine.env.SELDON_LOG_MESSAGES_EXTERNALLY | bool | `false` |  |
| engine.env.SELDON_LOG_REQUESTS | bool | `false` |  |
| engine.env.SELDON_LOG_RESPONSES | bool | `false` |  |
| engine.resources.requests.cpu | string | `"0.1"` |  |
| openvino.image | string | `"intelaipg/openvino-model-server:0.3"` |  |
| openvino.model.env.LOG_LEVEL | string | `"DEBUG"` |  |
| openvino.model.input | string | `"data"` |  |
| openvino.model.name | string | `"squeezenet1.1"` |  |
| openvino.model.output | string | `"prob"` |  |
| openvino.model.path | string | `"/opt/ml/squeezenet"` |  |
| openvino.model.resources | object | `{}` |  |
| openvino.model.src | string | `"gs://seldon-models/openvino/squeezenet"` |  |
| openvino.model_volume | string | `"hostPath"` |  |
| openvino.port | int | `8001` |  |
| predictorLabels.fluentd | string | `"true"` |  |
| predictorLabels.version | string | `"v1"` |  |
| sdepLabels.app | string | `"seldon"` |  |
| tfserving_proxy.image | string | `"seldonio/tfserving-proxy:0.2"` |  |
