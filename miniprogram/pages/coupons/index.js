const { request } = require('../../utils/request');

Page({
  data: {
    coupons: []
  },

  onShow() {
    this.loadCoupons();
  },

  getCouponDisplay(item) {
    const now = Date.now();
    const validTo = new Date(item.validTo).getTime();
    if (item.coupon?.status && item.coupon.status !== 'active') {
      return { statusText: '已停用', statusColor: '#9ca3af', usableText: '该券已被后台停用' };
    }
    if (item.status === 'used') {
      return { statusText: '已使用', statusColor: '#60a5fa', usableText: '该券已使用' };
    }
    if (item.status === 'expired' || (validTo && validTo < now)) {
      return { statusText: '已过期', statusColor: '#9ca3af', usableText: '已超过有效期' };
    }
    return { statusText: '可使用', statusColor: '#22c55e', usableText: `满￥${item.coupon?.thresholdAmount || 0}可用` };
  },

  async loadCoupons() {
    try {
      const data = await request('/me/coupons');
      const coupons = (data.list || []).map((item) => ({
        ...item,
        display: this.getCouponDisplay(item)
      }));
      this.setData({ coupons });
    } catch (error) {
      wx.showToast({ title: '加载优惠券失败', icon: 'none' });
    }
  }
});
