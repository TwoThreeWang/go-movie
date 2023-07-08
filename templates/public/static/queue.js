// 创建一个名为Queue的构造函数
function Queue(maxLength) {
    this.maxLength = maxLength;
    this.items = JSON.parse(localStorage.getItem("queue")) || [];
}

// 向队列中添加一个元素
Queue.prototype.push = function(item) {
    this.items.unshift(item);
    if (this.items.length > this.maxLength) {
        this.items.pop();
    }
    localStorage.setItem("queue", JSON.stringify(this.items));
};

// 从队列中移除一个元素
Queue.prototype.pop = function() {
    var item = this.items.pop();
    localStorage.setItem("queue", JSON.stringify(this.items));
    return item;
};

// 获取队列中的所有元素
Queue.prototype.getItems = function() {
    return this.items;
};