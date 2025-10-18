/**
 * 统一的请求发送函数 (基于 Fetch API)
 *
 * @param {string} method - HTTP 方法 (如 'GET', 'POST', 'PUT', 'DELETE')
 * @param {string} url - 请求的 URL
 * @param {object | null} data - 请求体数据 (仅 POST/PUT/DELETE 有效) 或查询参数 (GET)
 * @param {object} options - 额外的配置，如 headers, credentials 等
 * @returns {Promise<any>} 返回一个 Promise，成功时解析为解析后的 JSON 或文本数据。
 */
function sendRequest(method, url, data, options = {}) {
    method = method.toUpperCase();
    const config = {
        method: method,
        headers: options.headers || {},
        // 允许传入其他 Fetch 选项，如 mode, cache, credentials 等
        ...options
    };

    let fullUrl = url;

    // 1. 处理 GET/HEAD 请求的查询参数
    if (method === 'GET' || method === 'HEAD') {
        if (data && typeof data === 'object') {
            const query = new URLSearchParams(data).toString();
            fullUrl = query ? `${url}?${query}` : url;
        }
    }
    // 2. 处理 POST/PUT/DELETE 请求的请求体
    else if (data !== null && data !== undefined) {
        // 默认将数据视为 JSON 发送
        if (!config.headers['Content-Type']) {
            config.headers['Content-Type'] = 'application/json';
        }
        config.body = JSON.stringify(data);
    }

    // 核心：使用 Fetch API 发送请求，并返回 Promise
    return fetch(fullUrl, config) // 这里的 fetch 本身就返回一个 Promise
        .then(response => {
            // 检查 HTTP 状态码
            if (!response.ok) {
                // 如果请求失败 (例如 404, 500)，抛出错误
                // 注意：Fetch 在网络错误时才 reject，HTTP 错误时是 resolve
                throw new Error(`请求失败，状态码: ${response.status} (${response.statusText})`);
            }

            // 决定如何解析响应体
            const contentType = response.headers.get("content-type");
            if (contentType && contentType.includes("application/json")) {
                return response.json(); // Promise<JSON Object>
            }

            return response.text(); // Promise<string>
        })
        .then(resp => {
            // 在这里，resp 就是解析后的 JSON 对象或文本字符串
            // 最终的返回值是这个 Promise 解析的值
            return resp;
        });
        // 外部调用者通过 .then() 或 await 来获取这个 resp
}
function showMsg(message, type = 'info') {
    const existingAlert = document.querySelector('.alert');
    if (existingAlert) {
        existingAlert.remove();
    }
    const alert = document.createElement('div');
    alert.className = `alert ${type}`;
    alert.textContent = message;
    document.body.appendChild(alert);
    setTimeout(() => alert.classList.add('show'), 10);
    setTimeout(() => {
        alert.classList.remove('show');
        setTimeout(() => alert.remove(), 300);
    }, 3000);
}
// 滑动刷新模拟
let startY = 0;
window.addEventListener('touchstart', e => {
    startY = e.touches[0].clientY;
});
window.addEventListener('touchend', e => {
    const endY = e.changedTouches[0].clientY;
    const isAtTop = (document.documentElement.scrollTop || document.body.scrollTop) === 0;
    // 3. 判断：是否在顶部 AND 向下滑动距离超过 120px
    if (isAtTop && (endY - startY > 240)) {
        showMsg('页面刷新中...');
        setTimeout(() => {
            window.location.reload();
        }, 600);
    }
});
// localStorage播放记录操作封装
const Storage = {
  key: "moovie_history",
  // 读取所有记录
  getAll() {
    try {
      return JSON.parse(localStorage.getItem(this.key)) || [];
    } catch (e) {
      console.error("读取 localStorage 失败:", e);
      return [];
    }
  },
  // 保存完整数组
  saveAll(list) {
    localStorage.setItem(this.key, JSON.stringify(list));
  },
  // 新增或更新（按 id 匹配）
  upsert(record) {
    const list = this.getAll();
    const index = list.findIndex(item => item.id === record.id);
    if (index !== -1) {
      // 更新旧记录
      list[index] = { ...list[index], ...record, timestamp: Date.now() };
    } else {
      // 新增新记录，插入到最前
      list.unshift({ ...record, timestamp: Date.now() });
    }
    // 限制最多存 20 条
    this.saveAll(list.slice(0, 20));
  },
  // 删除某条记录
  remove(id) {
    const list = this.getAll().filter(item => item.id !== id);
    this.saveAll(list);
  },
  // 清空所有记录
  clear() {
    localStorage.removeItem(this.key);
  },
  // 查询单条记录
  getById(id) {
    return this.getAll().find(item => item.id === id) || null;
  }
};
function testWeb(url) {
    return new Promise((resolve, reject) => {
        const startTime = performance.now();
        const timeoutDuration = 2; // 2秒
        let isCompleted = false; // 标记是否已完成，防止重复处理

        const timeout = setTimeout(() => {
            if (!isCompleted) {
                isCompleted = true;
                // 连接超时，返回 -1
                reject(-1);
            }
        }, timeoutDuration);

        // 成功处理函数：计算时长并返回结果
        function handleSuccess() {
            if (isCompleted) return;
            isCompleted = true;
            clearTimeout(timeout);
            const endTime = performance.now();
            const duration = Math.floor(endTime - startTime);

            // 成功时，解决 Promise 并返回 duration
            resolve(duration);
        }

        // 失败处理函数：清除计时器并返回 -1
        function handleFailure(reason) {
            if (isCompleted) return;
            isCompleted = true;
            clearTimeout(timeout);

            // 失败时，拒绝 Promise 并返回 -1
            reject(-1);
        }

        // 核心请求
        fetch(url, {
            mode: 'no-cors',
            cache: 'no-store'
        })
        .then(handleSuccess)
        .catch((err) => {
            console.error("Fetch 捕获到错误:", err);
            handleFailure(err);
        });
    });
}