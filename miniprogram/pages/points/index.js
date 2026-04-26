const { request } = require('../../utils/request');

Page({
  data: {
    points: []
  },

  onShow() {
    this.loadPoints();
  },

  async loadPoints() {
    try {
      const data = await request('/me/points');
      this.setData({ points: data.list || [] });
    } catch (error) {
      wx.showToast({ title: '加载积分失败', icon: 'none' });
    }
  }
});
