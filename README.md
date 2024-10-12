# ababdonlgzj
灵敢足迹 app 的导出工具，支持查看足迹数据，导出为 csv


# 执行方法
## docker启动 （推荐）
``docker run -d -it -p 5175:80 zqzess/abandonlgzj:latest /bin/sh``

访问 `http://localhost:5175/#/upload`
## 自行编译启动
### 前台 frontend
1. `npm install`
2. `npm run dev`
### 后台 backend
参考 makefile 执行对应命令

示例:`make -f ./Makefile -C ../abandonlgzj all`

# 如何将导出的csv文件转换成其它 app 需要的格式
参考github 项目

[https://github.com/Aldenhovel/zuji](https://github.com/Aldenhovel/zuji)
# 预览图
<img width="500" alt="image" src="https://github.com/user-attachments/assets/2811ed9a-4a59-4a9d-aa6b-39df24c4628f">

<img width="500" alt="image" src="https://github.com/user-attachments/assets/d92dc959-5e47-4c73-a235-71bd0b0acc3b">

