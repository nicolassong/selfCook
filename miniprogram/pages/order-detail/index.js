const { request } = require('../../utils/request');
const { ORDER_STATUS } = require('../../utils/constants');

Page({
  data: { order: null },

  statusText(status) {
    const map = {
      [ORDER_STATUS.joined]: '待成团',
      [ORDER_STATUS.cutoffLocked]: '已截单待履约',
      [ORDER_STATUS.readyForPickup]: '待取餐',
      [ORDER_STATUS.delivering]: '配送中',
      [ORDER_STATUS.completed]: '已完成',
      [ORDER_STATUS.cancelled]: '已取消'
    };
    return map[status] || status;
  },

  fulfillmentModeText(mode) {
    const map = {
      pickup: '自提',
      delivery: '配送',
      mixed: '混合履约'
    };
    return map[mode] || mode;
  },

  onLoad(query) {
    this.orderNo = query.orderNo;
    this.loadOrder();
  },

  async loadOrder() {
    try {
      const order = await request(`/orders/no/${this.orderNo}`);
      this.setData({
        order: {
          ...order,
          statusText: this.statusText(order.status),
          fulfillmentModeText: this.fulfillmentModeText(order.fulfillmentMode)
        }
      });
    } catch (error) {
      wx.showToast({ title: '订单加载失败', icon: 'none' });
    }
  }
});
