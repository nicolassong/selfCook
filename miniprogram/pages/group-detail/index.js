const { request } = require('../../utils/request');
const { GROUP_STATUS } = require('../../utils/constants');

Page({
  data: {
    group: null
  },

  statusText(status) {
    const map = {
      [GROUP_STATUS.ongoing]: '进行中',
      [GROUP_STATUS.cutoffLocked]: '已截单',
      [GROUP_STATUS.completed]: '已完成',
      [GROUP_STATUS.cancelled]: '已取消'
    };
    return map[status] || status;
  },

  canBuy(item, groupStatus) {
    return groupStatus === GROUP_STATUS.ongoing && Number(item.stockAvailable || 0) > 0;
  },

  buyButtonText(item, groupStatus) {
    if (groupStatus !== GROUP_STATUS.ongoing) return '已截单';
    if (Number(item.stockAvailable || 0) <= 0) return '已售罄';
    return '立即下单';
  },

  normalizeGroup(group) {
    const items = (group.items || []).map((item) => ({
      ...item,
      canBuy: this.canBuy(item, group.status),
      buyButtonText: this.buyButtonText(item, group.status)
    }));

    return {
      ...group,
      statusText: this.statusText(group.status),
      items,
      itemCount: items.length
    };
  },

  onLoad(query) {
    this.groupId = query.id;
    this.loadGroup();
  },

  async loadGroup() {
    try {
      const group = await request(`/groups/${this.groupId}`);
      this.setData({ group: this.normalizeGroup(group) });
    } catch (err) {
      void err;
      wx.showToast({ title: '活动加载失败', icon: 'none' });
    }
  },

  goCreateOrder(event) {
    const { itemid, price, disabled } = event.currentTarget.dataset;
    if (disabled) return;
    wx.navigateTo({ url: `/pages/order-create/index?groupId=${this.groupId}&groupItemId=${itemid}&price=${price}` });
  }
});
