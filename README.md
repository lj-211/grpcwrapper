# grpcwrapper
对grpc进行简单封装，达到以下两个目的
    1. 尽量减少业务层对grpc原始的需求
    2. 增加中间件支持对grpc单个请求进行中间处理

注意:
    不限制使用grpc原生接口，按需使用即可；这个封装为了统一各个服务对于server和client的
    行为。
    肯定会存在需要单独设置特殊的option选项，这个时候自己选择原生API即可。
## 特性
### 拦截器
拦截器分为客户端拦截器和服务器拦截器:
服务器拦截器可以用来实现以下功能:
    1. 统一鉴权
    2. 流量控制
    3. 输出统计信息
    4. 等等...
客户端多拦截器
    1. 统一附带规范元数据
    2. 等等...
### 元数据[todo]
对部分元数据和业务关联的行为进行封装，如果需要操作元数据，请直接操作原始API。
### 基础操作
提供通用的服务搭建处理接口，如果要有不同需求，请使用原生接口。
### 客户端
    1. 提供统一超时的调用函数
    2. 提供统一keepalive参数
## 补充特性
    1. none
