#!/usr/bin/env python3
import sys
import subprocess
import platform

# 输出 Hello World
print("Hello World")

# 获取并打印命令行参数
print(f"\n脚本路径: {sys.argv[0]}")
print(f"参数数量: {len(sys.argv) - 1}")
if len(sys.argv) > 1:
    print("参数列表:")
    for i, arg in enumerate(sys.argv[1:], 1):
        print(f"  参数 {i}: {arg}")
else:
    print("没有传递任何参数")
print()
# 从标准输入读取每一行（支持管道和非交互模式）
try:
    while True:
        line = sys.stdin.readline()
        # 如果读到EOF或空行则退出
        if not line or line.strip() == '':
            break
        # 去掉末尾的换行符并输出
        print(f"Line: {line.rstrip()}")
except (EOFError, KeyboardInterrupt):
    # 如果没有输入流或用户中断，静默处理
    pass

# 持续 ping
try:
    # 根据操作系统选择合适的 ping 参数
    if platform.system().lower() == 'windows':
        # Windows: -t 表示持续 ping
        subprocess.run(['ping', '-t', '1.1.1.1'])
    else:
        # Linux/Mac: 不需要 -t 参数，默认就是持续 ping
        subprocess.run(['ping', '1.1.1.1'])
except KeyboardInterrupt:
    print("\nPing 已停止")

