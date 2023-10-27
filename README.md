# tbd
A command-line productivity tool powered by Baidu Qianfan.

## Installation
```shell
go build
```

You'll need an API Key and secret. 

1. Create application on here: https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application
2. Then you will get API Key and Secret Key.

Set the API Key and Secret Key in `config.yaml` stored in `$HOME/.tgb/` or `/etc/tgb/`.

## Usage
`tbd` has a variety of use cases, including simple queries, shell queries.

### Simple queries
We can use it as normal search engine, asking about anything:

```shell
 » ./tbd chat "How to install python on my computer"
To install Python on your computer, follow these steps:

1. Go to the Python website ([python.org](http://python.org)) and download the latest version of Python that matches your operating system.
2. Run the installer and follow the instructions. Make sure to select the option to add Python to your system path during the installation process.
3. After the installation is complete, open a command prompt or terminal window and type "python" to verify that Python has been successfully installed.

If you have any issues with the installation, you can refer to the Python documentation or seek help from the Python community.
```

Also in chinese: 

```shell
 » ./tbd chat "今天深圳天气怎样？"                                                                                                              1 ↵
深圳现在气温25℃，晴，东北风1级，2023年10月27日（今日）气温22~29℃，空气质量优，空气质量指数27。


近几日天气信息：

* 2023-10-26：阴转多云，23~30℃，东南风<3级，空气质量优。

* **2023-10-27：阴，22~29℃，无持续风向<3级，空气质量优**。

* 2023-10-28：小雨转中雨，22~25℃，无持续风向<3级，空气质量优。

* 2023-10-29：中雨，21~25℃，东风3-4级，空气质量优。

* 2023-10-30：阴转晴，21~26℃，东北风3-4级，空气质量优。
* 2023-10-28：小雨转中雨，22~25℃，无持续风向<3级，空气质量优。
```
