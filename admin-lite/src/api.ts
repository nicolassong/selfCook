export type ApiResponse<T> = {
  code: number;
  message: string;
  data: T;
};

const API_BASE = 'http://localhost:8081/api/v1';
const ASSET_BASE = 'http://localhost:8081';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options?.headers || {}),
    },
    ...options,
  });

  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`);
  }

  const result = (await response.json()) as ApiResponse<T>;
  if (result.code !== 0) {
    throw new Error(result.message || 'Request failed');
  }
  return result.data;
}

async function uploadFile(path: string, file: File): Promise<{ url: string; fileName: string; size: number }> {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    throw new Error(`Upload failed: ${response.status}`);
  }

  const result = (await response.json()) as ApiResponse<{ url: string; fileName: string; size: number }>;
  if (result.code !== 0) {
    throw new Error(result.message || 'Upload failed');
  }
  return result.data;
}

export type Product = {
  id: number;
  name: string;
  coverImage?: string;
  categoryName: string;
  status: string;
  sortOrder: number;
  skus?: { id: number; skuName: string; price: number; stockAvailable: number }[];
};

export type Group = {
  id: number;
  title: string;
  coverImage?: string;
  status: string;
  fulfillmentMode: string;
  cutoffAt: string;
  items?: { id: number; productName: string; skuName: string; stockAvailable: number; productSkuId?: number }[];
};

export type Order = {
  id: number;
  orderNo: string;
  groupId: number;
  status: string;
  fulfillmentMode: string;
  contactName: string;
  contactPhone: string;
  paidAmount: number;
  createdAt: string;
  items?: { id: number; productName: string; skuName: string; quantity: number }[];
};

export type Coupon = {
  id: number;
  name: string;
  couponType: string;
  amount: number;
  thresholdAmount: number;
  status: string;
  validFrom: string;
  validTo: string;
  totalCount: number;
  perUserLimit: number;
  grantedCount?: number;
  usedCount?: number;
};

export type UserCoupon = {
  id: number;
  couponId: number;
  userId: number;
  status: string;
  acquiredAt: string;
  usedAt?: string | null;
  orderId?: number | null;
  validFrom: string;
  validTo: string;
  coupon: Coupon;
};

export type SummaryRow = {
  productName: string;
  skuName: string;
  totalQty: number;
};

export type DailyMenu = {
  id: number;
  menuDate: string;
  title: string;
  status: string;
  items?: {
    id: number;
    productSkuId: number;
    stockTotal: number;
    stockAvailable: number;
    price: number;
    originalPrice: number;
    limitPerUser: number;
    limitPerOrder: number;
    sortOrder: number;
  }[];
};

export type DailyMenuCreateItem = {
  productSkuId: number;
  stockTotal: number;
  price: number;
  originalPrice?: number;
  limitPerUser?: number;
  limitPerOrder?: number;
  sortOrder?: number;
};

export type CreateProductPayload = {
  name: string;
  subtitle?: string;
  coverImage?: string;
  categoryName?: string;
  description?: string;
  sortOrder?: number;
  skus: {
    skuName: string;
    skuCode: string;
    price: number;
    originalPrice?: number;
    stockTotal: number;
    limitPerUser?: number;
    limitPerOrder?: number;
  }[];
};

export type CreateGroupPayload = {
  title: string;
  coverImage?: string;
  leaderUserId?: number;
  startAt: string;
  cutoffAt: string;
  fulfillmentMode: string;
  allowModifyBeforeCutoff?: boolean;
  showJoinList?: boolean;
  pickupRuleDesc?: string;
  deliveryRuleDesc?: string;
  groupNotice?: string;
  menuDate?: string;
  items: {
    productSkuId: number;
    stockTotal: number;
    price: number;
    originalPrice?: number;
    limitPerUser?: number;
    limitPerOrder?: number;
    sortOrder?: number;
  }[];
};

export type CreateCouponPayload = {
  name: string;
  couponType: string;
  amount: number;
  thresholdAmount?: number;
  applicableScope?: string;
  validDays?: number;
  totalCount?: number;
  perUserLimit?: number;
};

export type CreateDailyMenuPayload = {
  menuDate: string;
  title?: string;
  items: DailyMenuCreateItem[];
};

export async function fetchAdminProducts() {
  return request<{ list: Product[] }>('/admin/products');
}

export async function fetchAdminDailyMenus() {
  return request<{ list: DailyMenu[] }>('/admin/daily-menus');
}

export async function createDailyMenu(payload: CreateDailyMenuPayload) {
  return request<DailyMenu>('/admin/daily-menus', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function deleteDailyMenu(menuId: number) {
  return request<{ id: number }>(`/admin/daily-menus/${menuId}`, {
    method: 'DELETE',
  });
}

export async function fetchAdminGroups() {
  return request<{ list: Group[] }>('/admin/groups');
}

export async function fetchAdminOrders() {
  return request<{ list: Order[] }>('/admin/orders');
}

export async function fetchAdminCoupons() {
  return request<{ list: Coupon[] }>('/admin/coupons');
}

export async function fetchAdminUserCoupons(userId: number, status = 'all') {
  return request<{ list: UserCoupon[] }>(`/admin/user-coupons?userId=${userId}&status=${status}`);
}

export async function disableCoupon(couponId: number) {
  return request<{ id: number; status: string }>(`/admin/coupons/${couponId}/disable`, {
    method: 'POST',
  });
}

export async function createProduct(payload: CreateProductPayload) {
  return request<Product>('/admin/products', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function createGroup(payload: CreateGroupPayload) {
  return request<Group>('/admin/groups', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function createCoupon(payload: CreateCouponPayload) {
  return request<Coupon>('/admin/coupons', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function grantCoupon(couponId: number, userId: number) {
  return request<{ id: number }>('/admin/coupons/grant', {
    method: 'POST',
    body: JSON.stringify({ couponId, userId }),
  });
}

export async function cutoffGroup(groupId: number) {
  return request<{ groupId: number; status: string }>(`/leader/groups/${groupId}/cutoff`, {
    method: 'POST',
    body: JSON.stringify({ reason: '后台手动截单' }),
  });
}

export async function fetchGroupSummary(groupId: number) {
  return request<{ groupId: number; bySku: SummaryRow[] }>(`/leader/groups/${groupId}/summary`);
}

export async function markOrderReadyForPickup(orderId: number) {
  return request<{ orderId: number; status: string }>(`/admin/orders/${orderId}/ready-for-pickup`, {
    method: 'POST',
    body: JSON.stringify({ remark: '后台标记待取餐' }),
  });
}

export async function markOrderDelivering(orderId: number) {
  return request<{ orderId: number; status: string }>(`/admin/orders/${orderId}/start-delivery`, {
    method: 'POST',
    body: JSON.stringify({ remark: '后台标记配送中' }),
  });
}

export async function completeOrder(orderId: number) {
  return request<{ orderId: number; status: string }>(`/admin/orders/${orderId}/complete`, {
    method: 'POST',
    body: JSON.stringify({ remark: '后台完成订单' }),
  });
}

export async function uploadAdminImage(file: File) {
  return uploadFile('/admin/uploads/image', file);
}

export { ASSET_BASE };
