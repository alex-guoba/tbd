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
# -> BD is a command-line productivity tool powered by Baidu Qianfan, used for simple queries, summaries, and analysis. It accepts prompts from stdin and command-line arguments, offering versatility in usage scenarios. The tool also supports multi-round conversations and syncing conversation messages to external systems.
```

Analyzing example:

```shell
 » ping -t 2 cloud.baidu.com | tbd chat "diagnostic network condition from the Ping result"
# -> Based on the provided ping results, we can make the following diagnosis of the network condition:
# -> 1. **Connectivity**: From the results, you have successfully pinged `bce.baidu.n.shifen.com` with no packet loss (0%!p(MISSING)acket loss). This indicates that the connection between your computer and the server is stable.
# -> 2. **Latency**: The minimum/average/maximum/standard deviation round-trip latency is: 11.511/11.809/12.107/0.298 ms. These values represent the response time between your computer and the server. Generally, lower latency indicates better network performance. Your latency values are relatively low, indicating good network connectivity.
# -> In summary, based on this ping result, the network connection between your computer and the `bce.baidu.n.shifen.com` server is stable, and latency is low. If you encounter issues using the server or accessing related applications, the problem may not lie with the network connection but with other factors such as issues with the application server itself or software problems.
```


### Interactive mode

TBD supports multi-round conversations, which can be used for complex tasks or clarifying questions if the user's initial question is not clear or specific enough. Use the `-i` or `--interact` parameter to enter interactive mode. Type `exit` to exit the conversation:

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
