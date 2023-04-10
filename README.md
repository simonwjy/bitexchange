# BitExchange

## Overview
搭建一个满足数字货币交易功能的Crypto Exchange，主要需要支持以下几个模块:

- 登陆验证
- 用户管理
- 订单管理
- 行情管理
- Matching Engine
- 管理后台
- Crypto Wallet
- Trading Bot
- ...

创建这个项目是为了摸索币圈交易所的实现逻辑以及Web 3相关的技术栈，有几个扩展的点是值得去探究:

- decentralised crypto wallet (去中心化钱包)
- decentralised crypto exchange (去中心化交易所)
- issue ETH based crypto currency and publish to the exchange (ETH链上发币并在交易所发布和交易)
- crypto quant trading (crypto量化交易)
- market making bot (做市机器人)

### Phase 1
搭建一个满足MVP功能的crypto exchange的框架，支持用户在上面交易基于private network的ETH

### Phase 2
在ETH private network发布基于ERC-20协议的新币，并且可以通过交易所发布和交易

### Phase 3
增加/加强exchange其他功能模块包括行情管理以及管理后台等

### Phase 4
做市交易机器人以及量化交易
