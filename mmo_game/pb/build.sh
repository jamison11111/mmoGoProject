#编译Proto协议文件的编译指令脚本,把proto协议文件编译后输出的go文件输出到当前路径.下,编译当前路径下所有的proto文件
#!/bin/bash
protoc --go_out=. *.proto    
