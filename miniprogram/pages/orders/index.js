const { request } = require('../../utils/request');
const { ORDER_STATUS } = require('../../utils/constants');

Page({
  data: {
    orders: [],
    ORDER_STATUS
  },

  onShow() {
    this.loadOrders();
  },

  fulfillmentModeText(mode) {
    const map = {
      pickup: '自提',
      delivery: '配送',
      mixed: '混合履约'
    };
    return map[mode] || mode;
  },

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

  async loadOrders() {
    try {
      const data = await request('/me/orders');
      const list = (data.list || []).map((item) => ({ ...item, statusText: this.statusText(item.status), fulfillmentModeText: this.fulfillmentModeText(item.fulfillmentMode) }));
      this.setData({ orders: list });
    } catch (error) {
      wx.showToast({ title: '加载订单失败', icon: 'none' });
    }
  },

  async cancelOrder(event) {
    const { id } = event.currentTarget.dataset;
    try {
      await request(`/orders/${id}/cancel`, 'POST', { reason: '用户取消' });
      wx.showToast({ title: '已取消', icon: 'success' });
      this.loadOrders();
    } catch (error) {
      wx.showToast({ title: '取消失败', icon: 'none' });
    }
  },

  goDetail(event) {
    const { orderno } = event.currentTarget.dataset;
    wx.navigateTo({ url: `/pages/order-detail/index?orderNo=${orderno}` });
  }
});
