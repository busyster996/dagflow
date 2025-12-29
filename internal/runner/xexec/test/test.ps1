# 输出Hello World
Write-Host "Hello World"

# 从标准输入读取每一行（支持管道和非交互模式）
try {
    while ($null -ne ($line = [Console]::In.ReadLine())) {
        if ([string]::IsNullOrEmpty($line)) { break }
        Write-Host "Line: $line"
    }
} catch {
    # 如果没有输入流，静默处理
}
# 持续ping
ping -t 1.1.1.1

