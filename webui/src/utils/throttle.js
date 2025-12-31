/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} wait - 等待时间（毫秒）
 * @param {Object} options - 配置选项
 * @returns {Function} 防抖后的函数
 */
export function debounce(func, wait = 300, options = {}) {
  let timeout = null
  let lastArgs = null
  let lastThis = null
  let result = null
  
  const { leading = false, trailing = true, maxWait = null } = options
  
  let lastCallTime = 0
  let lastInvokeTime = 0
  
  function invokeFunc(time) {
    const args = lastArgs
    const thisArg = lastThis
    
    lastArgs = lastThis = null
    lastInvokeTime = time
    result = func.apply(thisArg, args)
    return result
  }
  
  function startTimer(pendingFunc, wait) {
    return setTimeout(pendingFunc, wait)
  }
  
  function cancelTimer(id) {
    clearTimeout(id)
  }
  
  function shouldInvoke(time) {
    const timeSinceLastCall = time - lastCallTime
    const timeSinceLastInvoke = time - lastInvokeTime
    
    return (
      lastCallTime === 0 ||
      timeSinceLastCall >= wait ||
      timeSinceLastCall < 0 ||
      (maxWait !== null && timeSinceLastInvoke >= maxWait)
    )
  }
  
  function timerExpired() {
    const time = Date.now()
    if (shouldInvoke(time)) {
      return trailingEdge(time)
    }
    timeout = startTimer(timerExpired, remainingWait(time))
  }
  
  function remainingWait(time) {
    const timeSinceLastCall = time - lastCallTime
    const timeSinceLastInvoke = time - lastInvokeTime
    const timeWaiting = wait - timeSinceLastCall
    
    return maxWait !== null
      ? Math.min(timeWaiting, maxWait - timeSinceLastInvoke)
      : timeWaiting
  }
  
  function trailingEdge(time) {
    timeout = null
    
    if (trailing && lastArgs) {
      return invokeFunc(time)
    }
    lastArgs = lastThis = null
    return result
  }
  
  function leadingEdge(time) {
    lastInvokeTime = time
    timeout = startTimer(timerExpired, wait)
    return leading ? invokeFunc(time) : result
  }
  
  function debounced(...args) {
    const time = Date.now()
    const isInvoking = shouldInvoke(time)
    
    lastArgs = args
    lastThis = this
    lastCallTime = time
    
    if (isInvoking) {
      if (timeout === null) {
        return leadingEdge(lastCallTime)
      }
      if (maxWait !== null) {
        timeout = startTimer(timerExpired, wait)
        return invokeFunc(lastCallTime)
      }
    }
    if (timeout === null) {
      timeout = startTimer(timerExpired, wait)
    }
    return result
  }
  
  debounced.cancel = function() {
    if (timeout !== null) {
      cancelTimer(timeout)
    }
    lastInvokeTime = 0
    lastArgs = lastCallTime = lastThis = timeout = null
  }
  
  debounced.flush = function() {
    return timeout === null ? result : trailingEdge(Date.now())
  }
  
  debounced.pending = function() {
    return timeout !== null
  }
  
  return debounced
}

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} wait - 等待时间（毫秒）
 * @param {Object} options - 配置选项
 * @returns {Function} 节流后的函数
 */
export function throttle(func, wait = 300, options = {}) {
  let timeout = null
  let previous = 0
  let result = null
  
  const { leading = true, trailing = true } = options
  
  const throttled = function(...args) {
    const now = Date.now()
    
    if (!previous && !leading) {
      previous = now
    }
    
    const remaining = wait - (now - previous)
    
    if (remaining <= 0 || remaining > wait) {
      if (timeout) {
        clearTimeout(timeout)
        timeout = null
      }
      previous = now
      result = func.apply(this, args)
    } else if (!timeout && trailing) {
      timeout = setTimeout(() => {
        previous = leading ? Date.now() : 0
        timeout = null
        result = func.apply(this, args)
      }, remaining)
    }
    
    return result
  }
  
  throttled.cancel = function() {
    if (timeout) {
      clearTimeout(timeout)
      timeout = null
    }
    previous = 0
  }
  
  return throttled
}

/**
 * 请求动画帧节流
 * @param {Function} func - 要节流的函数
 * @returns {Function} 节流后的函数
 */
export function rafThrottle(func) {
  let rafId = null
  let lastArgs = null
  
  const throttled = function(...args) {
    lastArgs = args
    
    if (rafId === null) {
      rafId = requestAnimationFrame(() => {
        func.apply(this, lastArgs)
        rafId = null
        lastArgs = null
      })
    }
  }
  
  throttled.cancel = function() {
    if (rafId !== null) {
      cancelAnimationFrame(rafId)
      rafId = null
      lastArgs = null
    }
  }
  
  return throttled
}