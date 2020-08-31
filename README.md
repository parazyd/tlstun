我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

tlstun
======

tlstun is a simple Go program that will add TLS support for your
programs that do not have it.

It simply proxies from one TLS-listening host:port to another plaintext
host:port. If TLS is not your thing, you can also proxy plain TCP
traffic.


Installation
------------

```
$ go get github.com/parazyd/tlstun
```

Make sure you generate or acquire a TLS certificate keypair to use with
tlstun.


Usage
-----

```
Usage of ./tlstun:
  -cacert string
        Path for CA certificate file (default "ca.pem")
  -cert string
        Path for Certificate file (default "server.pem")
  -forward string
        Forward address (default "127.0.0.1:72")
  -key string
        Path for Key file (default "server-key.pem")
  -listen string
        Listen address (default "127.0.0.1:7443")
  -notls
        Disable TLS and just tunnel plain TCP
  -tlsver int
        TLS version to use (11, 12, 13) (default 13)
  -verbose
        Verbose mode
  -verifyclient
        Do client verification
```

tlstun supports two different ways of multiplexing, one being normal TLS
proxying, and the other being TLS proxying with client certificate
authentication. In addition to this, tlstun can also opt-out of TLS and
proxy plain TCP without encryption by using the `-notls` flag.


### Without client verification

Start tlstun with `-cert` and `-key`, and it will simply provide a TLS
forward to its destination with the given TLS certificate.


### With client verification

With client verification, start tlstun with `-cacert`, `-cert`, `-key`,
and `-verifyclient` and it will do client certificate verification. This
means it will only allow access from clients providing a certificate
signed by the CA certificate that is being loaded/used with tlstun on
startup with `-cacert`.
