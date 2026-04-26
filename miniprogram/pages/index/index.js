const { request } = require('../../utils/request');
const { GROUP_STATUS } = require('../../utils/constants');

Page({
  data: {
    groups: [],
    loading: false
  },

  fulfillmentModeText(mode) {
    const map = {
      pickup: '自提',
      delivery: '配送',
      mixed: '混合履约'
    };
    return map[mode] || mode || '混合履约';
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

  buildItemsSummary(items) {
    if (!Array.isArray(items) || !items.length) return '暂无菜单商品';
    return items
      .map((item) => `${item.productName || '商品'}(${item.skuName || '规格'}) x库存${item.stockAvailable}`)
      .join('；');
  },

  onShow() {
    this.loadGroups();
  },

  async loadGroups() {
    this.setData({ loading: true });
    try {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      const minCutoffAt = today.toISOString().slice(0, 10);
      const query = `status=${GROUP_STATUS.ongoing}&minCutoffAt=${encodeURIComponent(minCutoffAt)}`;
      const data = await request(`/groups?${query}`);
      const groups = (data.list || []).map((item) => ({
        ...item,
        items: item.items || [],
        itemCount: (item.items || []).length,
        itemsSummary: this.buildItemsSummary(item.items || []),
        fulfillmentModeText: this.fulfillmentModeText(item.fulfillmentMode),
        statusText: this.statusText(item.status)
      }));
      this.setData({ groups });
    } catch {
      wx.showToast({ title: '加载失败', icon: 'none' });
    } finally {
      this.setData({ loading: false });
    }
  },

  goDetail(event) {
    const { id } = event.currentTarget.dataset;
    wx.navigateTo({ url: `/pages/group-detail/index?id=${id}` });
  }
});
