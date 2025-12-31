echo "Hello World"
echo "脚本路径: $0"
echo "参数1: $1"
echo "参数2: $2"
echo "所有参数: $@"
echo "参数数量: $#"
while read line; do
  echo "Line: $line"
done
exit 123
# 持续ping
ping 1.1.1.1