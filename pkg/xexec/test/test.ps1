# 输出Hello World
Write-Host "Hello World"
Write-Host "参数数量: $($args.Count)"
for ($i=0; $i -lt $args.Count; $i++) {
    Write-Host "参数 $($i+1): $($args[$i])"
}
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

