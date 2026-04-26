const { request } = require('../../utils/request');

Page({
  data: {
    user: null,
    addresses: [],
    coupons: [],
    points: [],
    notifications: []
  },

  onShow() {
    this.loadData();
  },

  async loadData() {
    try {
      const [user, addresses, coupons, points, notifications] = await Promise.all([
        request('/me'),
        request('/me/addresses'),
        request('/me/coupons'),
        request('/me/points'),
        request('/me/notifications')
      ]);
      this.setData({
        user,
        addresses: addresses.list || [],
        coupons: coupons.list || [],
        points: points.list || [],
        notifications: notifications.list || []
      });
    } catch (error) {
      wx.showToast({ title: '加载我的页面失败', icon: 'none' });
    }
  }
});
