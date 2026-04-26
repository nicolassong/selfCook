const { request } = require('../../utils/request');

Page({
  data: {
    notifications: []
  },

  onShow() {
    this.loadNotifications();
  },

  async loadNotifications() {
    try {
      const data = await request('/me/notifications');
      this.setData({ notifications: data.list || [] });
    } catch (error) {
      wx.showToast({ title: '加载通知失败', icon: 'none' });
    }
  }
});
