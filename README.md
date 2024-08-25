![image](https://github.com/user-attachments/assets/f1a5901c-3bf8-415f-8481-e77f4e5c720e)> 声明：这篇文章是通过模仿Hexo博客的样式请进行编写的，并且是通过网上相似的Java代码通过自学的Go编写而成，如果代码有些不好的，请多多见谅并提出意见，进而使我不断地进步 谢谢
# 后端代码文件名为goBlog，后端文件里有一个文件夹为sql 拉去代码之后
> 将sql文件导入本地数据库里
> 打开application.yaml文件，修改文件里的配置
- **首先除了数据库密码，redis密码，rabbitmq密码之外还有两个比较重要的需要自己配置**
  - 修改上传图片策略相关的配置，比如我使用的是阿里云存储桶上传(https://oss.console.aliyun.com/bucket)
    ![image](https://github.com/user-attachments/assets/f3194f06-bf9a-4264-a1f7-1d92cdeb9af3)
    ![image](https://github.com/user-attachments/assets/964d20c7-5730-43b8-813e-529aac070ceb)
    将创建好的存储桶名复制
    ![image](https://github.com/user-attachments/assets/778983bc-a62c-4935-a8e2-184618ef2dd3)
    ![1724591478007](https://github.com/user-attachments/assets/038cd003-ac78-439d-a4e4-e9ae1b043aeb)
    ![image](https://github.com/user-attachments/assets/d87f7971-2cb9-4123-9586-2277724569a3)
    其中这个accessKeyId 和accessKeySecret对应yaml文件里的直接粘贴复制进去即可
  - 修改发送邮件相关的配置，比如我使用的是QQ发送邮件(https://mail.qq.com)，按照下图操作
    ![image](https://github.com/user-attachments/assets/7a77058d-b121-42f0-bfb3-1a55eb5af8ab)
    ![image](https://github.com/user-attachments/assets/d0dfca98-89f8-47b0-af30-aa6480b532cf)
    ![image](https://github.com/user-attachments/assets/5019470b-042d-4c85-928b-92e1f850174a)
    替换掉yaml文件里的配置即可
# 前端前台文件名为blog，拉取代码之后
- 具体功能有：首页、、搜索、归档、分类、标签、相册、说说、友链、关于、留言、登录
  - 首页文章 ![image](https://github.com/user-attachments/assets/af7f40ce-c823-45a6-afad-5aef7cb92799)
  ![image](https://github.com/user-attachments/assets/e7f5b535-99eb-422e-9358-0ad90ee2fa09)
  - 文章内容 ![image](https://github.com/user-attachments/assets/1d853417-544c-4cdc-b7af-beeac94b4067)
  - 搜索 ![image](https://github.com/user-attachments/assets/528cb56e-b76f-41cb-8e23-af590f14027c)
  - 友链 ![image](https://github.com/user-attachments/assets/aed4571d-0fd4-4938-a7e2-7b594a8277ca)
    等等
# 前端后台文件为admin，拉取代码之后
- 具体功能有：文章管理、消息管理、用户管理、权限管理、相册管理、等等
  - 登录页面 ![image](https://github.com/user-attachments/assets/6869e3fb-d20f-4eed-9daf-efee54808e98)
  - 首页 ![image](https://github.com/user-attachments/assets/b46bd3f4-86e2-471b-9d67-b31a50c040a3)
  - 文章管理
     - 发布文章 ![image](https://github.com/user-attachments/assets/a32c49ee-0ad4-4d79-8c7e-82e3c411a05f)
     - 文章列表 ![image](https://github.com/user-attachments/assets/313e2080-47af-46bb-ba58-b890371517dd)
     - 分类管理 ![image](https://github.com/user-attachments/assets/39a65ac9-c1a8-4bd8-93eb-272a0056b833)
     - 标签管理 ![image](https://github.com/user-attachments/assets/bb902a89-c880-4e10-9cdb-4728ea060bdd)
  - 说说管理 ![image](https://github.com/user-attachments/assets/8e860f01-598a-47ac-b6c9-4269206b1b6c)
    等等















    






