### protobuf说明
* 下载谷歌protoc用于生成 对应语言的代码
>> 不用下载具体语言版本，直接下载对应系统编译好的文件即可【如windows的.exe文件，Mac的编译好的protoc文件】
下载地址 https://github.com/protocolbuffers/protobuf/releases
* 定义自己需要的数据结构 .proto文件
    1. 具体语法介绍示例：https://colobu.com/2015/01/07/Protobuf-language-guide/
    2. 示例：import 引入其他的定义的.proto文件，类似于go用的import。
       package 命名空间的意思
    3. 定义好之后：使用第一步下载好的protoc文件执行生成语言代码操作.<br>
       如抖音的 protoc --go_out=plugins=grpc:. ./opentrace/proto/*.proto --proto_path=./ 
       --go_out是自定生成的具体文件类所在根目录，会依据proto定义的package【自己定义的命名空间规则】在这个目录下生成目录及文件。<br>
       --proto_path 是引入其他proto文件的目录，当多个自定义的proto文件之间有引入的时候需要指定这个目录，另外需要拷贝第一步下载的目录里的google这个目录到当前命令执行目录下，这个里边是谷歌定义的一些元 proto类型，如Any.proto可以表示是任意类型
* 使用特别提示
    1. 在使用时可以直接操作二进制流数据如<br>
    2.  cat data-danmu.txt | ~/Downloads/protoc-3.11.4-osx-x86_64/bin/protoc --decode_raw > data-danmu.json<br>
        使用 protoc --decode_raw可以解析源二进制文件为正常的数据文件，依据数据文件可以自己定义数据的结构体文件 .proto
#### 调试grpc
   可以使用 grpcui工具调试
