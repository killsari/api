项目结构

    common:公共
    config:初始化
    controller:控制
    model:数据
    release:部署
    service:服务层
    tool:脚本工具
    utils:工具类
    main.go:运行入口

运行配置文件路径:

    release/config/production.yaml

运行:

    go run main.go [-c release/config/production.yaml]
    注意：配置文件对应的相对路径需要正确

编译(for linux):

    ./build

编译运行:

    ./xxx [-c release/config/production.yaml]
    注意：配置文件对应的相对路径需要正确