/ 定义表结构
persons: ([] name: `symbol$(); age: (); city: `symbol$())

/ 向表中插入数据
persons: insert[`persons; (`John; 25; `London)]
persons: insert[`persons; (`Alice; 32; `New York)]
persons: insert[`persons; (`Bob; 28; `Paris)]

/ 启动KDB服务器监听端口
\p 5000