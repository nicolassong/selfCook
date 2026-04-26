const { request } = require('../../utils/request');

Page({
  data: {
    contactName: '演示用户',
    contactPhone: '13800001111',
    pickupPointId: 1,
    mode: 'pickup',
    addresses: [],
    selectedAddressId: 1,
    selectedAddressText: '',
    coupons: [],
    couponOptions: [],
    selectedCouponId: null,
    selectedCouponText: '暂不使用优惠券',
    couponHint: '可选择一张满足门槛的优惠券',
    groupItems: [],
    quantityMap: {},
    totalCount: 0,
    totalAmount: 0,
    discountAmount: 0,
    payableAmount: 0,
    quantityHint: '请至少选择一件商品'
  },

  onLoad(query) {
    this.groupId = Number(query.groupId);
    this.groupItemId = Number(query.groupItemId || 0);
    this.loadGroupItems();
    this.updateAmount();
  },

  onShow() {
    this.loadAddresses();
    this.loadCoupons();
  },

  async loadGroupItems() {
    try {
      const group = await request(`/groups/${this.groupId}`);
      const items = (group.items || []).map((item) => {
        const stock = Number(item.stockAvailable || 0);
        const limit = Number(item.limitPerUser || 0);
        let maxPurchasable = stock;
        if (limit > 0) {
          maxPurchasable = Math.min(stock, limit);
        }
        return {
          ...item,
          maxPurchasable,
          purchasable: maxPurchasable > 0
        };
      });

      const quantityMap = {};
      const preferredItem = items.find((it) => Number(it.id) === this.groupItemId);
      if (preferredItem && preferredItem.purchasable) {
        quantityMap[preferredItem.id] = 1;
      }

      this.setData({ groupItems: items, quantityMap });
      this.updateAmount();
    } catch {
      wx.showToast({ title: '读取商品信息失败', icon: 'none' });
    }
  },

  async loadAddresses() {
    try {
      const data = await request('/me/addresses');
      const addresses = data.list || [];
      const defaultAddress = addresses.find((item) => item.isDefault) || addresses[0];
      this.setData({
        addresses,
        selectedAddressId: defaultAddress ? defaultAddress.id : 1,
        selectedAddressText: defaultAddress ? `${defaultAddress.province}${defaultAddress.city}${defaultAddress.district}${defaultAddress.detailAddress}` : ''
      });
    } catch {}
  },

  async loadCoupons() {
    try {
      const data = await request('/me/coupons');
      this.setData({ coupons: data.list || [] });
      this.updateAmount();
    } catch {}
  },

  getCouponState(coupon, totalAmount) {
    if (!coupon) {
      return { usable: false, discount: 0, text: '暂不使用优惠券' };
    }
    if (coupon.coupon?.status && coupon.coupon.status !== 'active') {
      return { usable: false, discount: 0, text: '券已停用' };
    }
    if (coupon.status !== 'unused') {
      return { usable: false, discount: 0, text: `不可用：${coupon.status}` };
    }
    if (coupon.coupon && totalAmount < Number(coupon.coupon.thresholdAmount || 0)) {
      return { usable: false, discount: 0, text: `未满￥${coupon.coupon.thresholdAmount}门槛` };
    }
    const discount = Number(coupon.coupon?.amount || 0);
    return { usable: true, discount, text: `可用，预计减￥${discount}` };
  },

  updateAmount() {
    const quantityMap = this.data.quantityMap || {};
    const selectedItems = (this.data.groupItems || []).filter((item) => Number(quantityMap[item.id] || 0) > 0);
    const totalCount = selectedItems.reduce((sum, item) => sum + Number(quantityMap[item.id] || 0), 0);
    const totalAmount = selectedItems.reduce((sum, item) => sum + Number(item.price || 0) * Number(quantityMap[item.id] || 0), 0);

    const couponOptions = this.data.coupons.map((item) => {
      const state = this.getCouponState(item, totalAmount);
      return {
        id: item.id,
        label: `${item.coupon?.name || `券#${item.couponId}`} / ${state.text}`,
        usable: state.usable,
        discount: state.discount
      };
    });

    const selectedCoupon = this.data.coupons.find((item) => item.id === this.data.selectedCouponId);
    const state = this.getCouponState(selectedCoupon, totalAmount);
    const discountAmount = state.usable ? Math.min(totalAmount, state.discount) : 0;
    const payableAmount = Math.max(0, totalAmount - discountAmount);

    const quantityHint = totalCount > 0
      ? `已选 ${totalCount} 份商品`
      : '请至少选择一件商品';

    this.setData({
      totalCount,
      totalAmount,
      discountAmount,
      payableAmount,
      couponOptions,
      couponHint: state.text,
      quantityHint
    });
  },

  setField(event) {
    const { field } = event.currentTarget.dataset;
    this.setData({ [field]: event.detail.value });
  },

  setMode(event) {
    const options = ['pickup', 'delivery'];
    this.setData({ mode: options[event.detail.value] });
  },

  changeItemQuantity(event) {
    const itemId = Number(event.currentTarget.dataset.itemid);
    const delta = Number(event.currentTarget.dataset.delta);
    const item = (this.data.groupItems || []).find((it) => Number(it.id) === itemId);
    if (!item || !item.purchasable) return;

    const quantityMap = { ...(this.data.quantityMap || {}) };
    const current = Number(quantityMap[itemId] || 0);
    const next = Math.max(0, Math.min(Number(item.maxPurchasable || 0), current + delta));

    if (next <= 0) {
      delete quantityMap[itemId];
    } else {
      quantityMap[itemId] = next;
    }

    this.setData({ quantityMap });
    this.updateAmount();
  },

  clearAllItems() {
    this.setData({ quantityMap: {} });
    this.updateAmount();
  },

  quickAddRecommended() {
    const quantityMap = {};
    const items = this.data.groupItems || [];
    let picked = 0;
    for (const item of items) {
      if (!item.purchasable) continue;
      quantityMap[item.id] = 1;
      picked += 1;
      if (picked >= 2) break;
    }
    this.setData({ quantityMap });
    this.updateAmount();
  },

  chooseAddress(event) {
    const index = Number(event.detail.value);
    const address = this.data.addresses[index];
    this.setData({
      selectedAddressId: Number(address.id),
      selectedAddressText: `${address.province}${address.city}${address.district}${address.detailAddress}`
    });
  },

  chooseCoupon(event) {
    const index = Number(event.detail.value);
    const coupon = this.data.coupons[index];
    const state = this.getCouponState(coupon, this.data.totalAmount);
    this.setData({
      selectedCouponId: coupon ? coupon.id : null,
      selectedCouponText: coupon ? `${coupon.coupon?.name || `券#${coupon.couponId}`} / ${state.text}` : '暂不使用优惠券'
    });
    this.updateAmount();
  },

  async submitOrder() {
    if (Number(this.data.totalCount || 0) <= 0) {
      wx.showToast({ title: '请至少选择一件商品', icon: 'none' });
      return;
    }

    const orderItems = [];
    for (const item of this.data.groupItems || []) {
      const quantity = Number((this.data.quantityMap || {})[item.id] || 0);
      if (quantity <= 0) continue;
      if (quantity > Number(item.maxPurchasable || 0)) {
        wx.showToast({ title: `${item.productName} 超出可购上限`, icon: 'none' });
        return;
      }
      orderItems.push({ groupItemId: Number(item.id), quantity, tasteRemark: '' });
    }

    const selectedCoupon = this.data.coupons.find((item) => item.id === this.data.selectedCouponId);
    const state = this.getCouponState(selectedCoupon, this.data.totalAmount);
    if (selectedCoupon && !state.usable) {
      wx.showToast({ title: state.text, icon: 'none' });
      return;
    }

    try {
      const data = await request('/orders', 'POST', {
        groupId: this.groupId,
        fulfillmentMode: this.data.mode,
        pickupPointId: this.data.mode === 'pickup' ? Number(this.data.pickupPointId) : undefined,
        addressId: this.data.mode === 'delivery' ? Number(this.data.selectedAddressId) : undefined,
        userCouponId: state.usable ? this.data.selectedCouponId : undefined,
        contactName: this.data.contactName,
        contactPhone: this.data.contactPhone,
        remark: '小程序下单',
        items: orderItems
      });
      wx.showToast({ title: '下单成功', icon: 'success' });
      setTimeout(() => {
        wx.redirectTo({ url: `/pages/order-detail/index?orderNo=${data.orderNo}` });
      }, 600);
    } catch {
      wx.showToast({ title: '下单失败', icon: 'none' });
    }
  }
});
