# tbd
A command-line productivity tool powered by Baidu Qianfan.

English | [简体中文](./README.zh-CN.md)

## Installation

```shell
go install
cp ./config.yaml $HOME/.config/tgb/
```

You'll need an API Key and Secret Key. 

1. Create application on here: https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application
2. Retrive API Key and Secret Key from the website after application approved.

Set the API Key and Secret Key in `config.yaml` stored in `./`, `$HOME/.config/tgb/` or `/etc/tgb/`.

## Usage
`tbd` has a variety of use cases, including simple queries, summerization and analyzing.

### Simple queries
We can use it as normal search engine, asking about anything:

```shell
 » tbd chat "How to install python"
To install Python on your computer, follow these steps:

1. Go to the Python website ([python.org](http://python.org)) and download the latest version of Python that matches your operating system.
2. Run the installer and follow the instructions. Make sure to select the option to add Python to your system path during the installation process.
3. After the installation is complete, open a command prompt or terminal window and type "python" to verify that Python has been successfully installed.

If you have any issues with the installation, you can refer to the Python documentation or seek help from the Python community.
```

Also in chinese: 

```shell
» tbd chat "今天深圳天气怎样？"
# -> 深圳现在气温26℃，多云，东北风3级，2023年10月30日（今日）气温22~29℃，空气质量优，空气质量指数55。
```

### Summerization and analyzing

`tbd` accepts prompt from both stdin and command line argument. This versatile feature is particularly useful when you need to pass file content or pipe output from other commands to the Baidu Qianfan models for summarization or analysis. For example, you can easily diagnostic local network condition:

Here is a summerization example: 

```shell
 » cat ./README.md | tbd chat "Briefly summarize the functions of tbd tools from README file, keep it to 40 words or less"
# -> tbd是一个命令行生产力工具，由百度千言驱动，可用于简单查询、摘要和分析。它接受来自stdin和命令行参数的提示，具有多种使用场景。该工具还支持多轮对话和同步对话消息到外部系统。
```

Analyzing example:

```shell
 » ping -t 2 cloud.baidu.com | tbd chat "diagnostic network condition from the Ping result"
# -> 根据提供的ping结果，我们可以得出以下网络状况诊断：
# -> 1. **连通性**：从结果来看，您成功地ping到了`bce.baidu.n.shifen.com`，且没有丢包（0.0%!p(MISSING)acket loss）。这意味着您的计算机与该服务器之间的网络连接是稳定的。
# -> 2. **延迟**：往返的最小/平均/最大/标准偏差延迟分别为：11.511/11.809/12.107/0.298 ms。这些数值表示您的计算机与服务器之间的响应时间。一般来说，延迟越低，网络性能越好。您的延迟数值相对较低，表明网络连接状况良好。
# -> 综上所述，根据这个ping结果，您的计算机与综上所述，根据这个ping结果，您的计算机与`bce.baidu.n.shifen.com`服务器之间的网络连接稳定，且延迟较低。如果您在使用该服务器或访问相关应用时遇到问题，那么问题可能不在于网络连接，而是其他因素，如应用服务器本身的问题、软件问题等。
```


### Interactive mode

`tbd` support multiple rounds dialogus. It can be used in complex tasks or to clarify the question, i.e. if the question raised by the user is not clear or specific enough. 

use `-i` or `--interact` args to enter into interactive mode. `exit` will finished the dialog:

```shell
» tbd chat -i
please remember my favorite number: 9
# -> Of course, I will remember your favorite number as 9. Is there anything else I can assist you with?
what would be my favorite number + 4?
# -> Your favorite number is 9, so adding 4 to it would result in 13. Therefore, your favorite number plus 4 would be 13.
exit
```

### Sync 

`tbd` support syning dialog message (request & response) to external system. like memos. See `config.yaml` for more detail.

## Reference

- [How to write prompt](https://console.bce.baidu.com/qianfan/prompt/template)
