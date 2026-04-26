import {
  Alert,
  App as AntdProvider,
  Button,
  Card,
  Drawer,
  DatePicker,
  Form,
  Input,
  InputNumber,
  Layout,
  Menu,
  Modal,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  Image,
} from 'antd';
import { useEffect, useMemo, useState } from 'react';
import {
  completeOrder,
  createCoupon,
  createDailyMenu,
  createGroup,
  createProduct,
  cutoffGroup,
  deleteDailyMenu as deleteDailyMenuApi,
  disableCoupon,
  fetchAdminCoupons,
  fetchAdminDailyMenus,
  fetchAdminGroups,
  fetchAdminOrders,
  fetchAdminProducts,
  fetchAdminUserCoupons,
  fetchGroupSummary,
  Coupon,
  grantCoupon,
  Group,
  Order,
  Product,
  SummaryRow,
  UserCoupon,
  markOrderDelivering,
  markOrderReadyForPickup,
  uploadAdminImage,
} from './api';

const { Header, Sider, Content } = Layout;
const { Title, Text } = Typography;

type TabKey = 'dashboard' | 'products' | 'menus' | 'groups' | 'orders' | 'coupons';

type DailyMenuItem = {
  id: number;
  date: string;
  title: string;
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

const GroupStatus = {
  Ongoing: 'ongoing',
  CutoffLocked: 'cutoff_locked',
  Completed: 'completed',
  Cancelled: 'cancelled',
} as const;

const OrderStatus = {
  Joined: 'joined',
  CutoffLocked: 'cutoff_locked',
  ReadyForPickup: 'ready_for_pickup',
  Delivering: 'delivering',
  Completed: 'completed',
  Cancelled: 'cancelled',
} as const;

const productStatusLabelMap: Record<string, string> = {
  on_sale: '上架中',
  off_sale: '已下架',
};

const groupStatusLabelMap: Record<string, string> = {
  [GroupStatus.Ongoing]: '进行中',
  [GroupStatus.CutoffLocked]: '已截单',
  [GroupStatus.Completed]: '已完成',
  [GroupStatus.Cancelled]: '已取消',
};

const fulfillmentModeLabelMap: Record<string, string> = {
  pickup: '自提',
  delivery: '配送',
  mixed: '混合履约',
};

const orderStatusLabelMap: Record<string, string> = {
  [OrderStatus.Joined]: '已下单',
  [OrderStatus.CutoffLocked]: '已截单',
  [OrderStatus.ReadyForPickup]: '待取餐',
  [OrderStatus.Delivering]: '配送中',
  [OrderStatus.Completed]: '已完成',
  [OrderStatus.Cancelled]: '已取消',
};

const couponTypeLabelMap: Record<string, string> = {
  full_reduction: '满减券',
};

const couponStatusLabelMap: Record<string, string> = {
  active: '启用',
  disabled: '停用',
};

const userCouponStatusLabelMap: Record<string, string> = {
  unused: '待使用',
  used: '已使用',
  expired: '已过期',
};

function getLabel(map: Record<string, string>, value: string) {
  return map[value] || value;
}

function getProductStatusColor(value: string) {
  return value === 'on_sale' ? 'green' : 'default';
}

function getGroupStatusColor(value: string) {
  if (value === GroupStatus.Ongoing) return 'green';
  if (value === GroupStatus.CutoffLocked) return 'orange';
  if (value === GroupStatus.Completed) return 'blue';
  return 'default';
}

function getOrderStatusColor(value: string) {
  if (value === OrderStatus.Joined) return 'cyan';
  if (value === OrderStatus.CutoffLocked) return 'orange';
  if (value === OrderStatus.ReadyForPickup) return 'gold';
  if (value === OrderStatus.Delivering) return 'blue';
  if (value === OrderStatus.Completed) return 'green';
  return 'default';
}

function getCouponStatusColor(value: string) {
  return value === 'active' ? 'green' : 'default';
}

function getUserCouponStatusColor(value: string) {
  if (value === 'unused') return 'green';
  if (value === 'used') return 'blue';
  return 'default';
}

function readFileAsDataURL(file: File) {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result || ''));
    reader.onerror = () => reject(new Error('图片读取失败'));
    reader.readAsDataURL(file);
  });
}

async function dataUrlToFile(dataUrl: string, fileName: string) {
  const response = await fetch(dataUrl);
  const blob = await response.blob();
  return new File([blob], fileName, { type: blob.type || 'image/jpeg' });
}

async function compressImage(file: File, maxWidth = 1600, quality = 0.9) {
  const dataUrl = await readFileAsDataURL(file);
  const image = await new Promise<HTMLImageElement>((resolve, reject) => {
    const img = new window.Image();
    img.onload = () => resolve(img);
    img.onerror = () => reject(new Error('图片加载失败'));
    img.src = dataUrl;
  });

  const ratio = Math.min(1, maxWidth / image.width);
  const canvas = document.createElement('canvas');
  canvas.width = Math.round(image.width * ratio);
  canvas.height = Math.round(image.height * ratio);
  const ctx = canvas.getContext('2d');
  if (!ctx) throw new Error('图片处理失败');
  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';
  ctx.drawImage(image, 0, 0, canvas.width, canvas.height);

  if (file.type === 'image/png') {
    return canvas.toDataURL('image/png');
  }
  return canvas.toDataURL('image/jpeg', quality);
}

function formatDateOnly(value: any) {
  if (!value) return '';
  if (typeof value?.format === 'function') return value.format('YYYY-MM-DD');
  if (value instanceof Date) return value.toISOString().slice(0, 10);
  return String(value).slice(0, 10);
}

export function App() {
  const { message } = AntdProvider.useApp();

  function notifyError(err: unknown, fallback: string) {
    message.error(err instanceof Error ? err.message : fallback);
  }
  const [activeTab, setActiveTab] = useState<TabKey>('dashboard');
  const [products, setProducts] = useState<Product[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [coupons, setCoupons] = useState<Coupon[]>([]);
  const [dailyMenus, setDailyMenus] = useState<DailyMenuItem[]>([]);
  const [userCoupons, setUserCoupons] = useState<UserCoupon[]>([]);
  const [couponQueryUserId, setCouponQueryUserId] = useState<number>(1);
  const [couponStatusFilter, setCouponStatusFilter] = useState<string>('all');
  const [summaryRows, setSummaryRows] = useState<SummaryRow[]>([]);
  const [summaryOpen, setSummaryOpen] = useState(false);
  const [productOpen, setProductOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const [groupOpen, setGroupOpen] = useState(false);
  const [couponOpen, setCouponOpen] = useState(false);
  const [grantOpen, setGrantOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState('');
  const [previewOpen, setPreviewOpen] = useState(false);
  const [productImageUploading, setProductImageUploading] = useState(false);
  const [groupImageUploading, setGroupImageUploading] = useState(false);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [productForm] = Form.useForm();
  const [menuForm] = Form.useForm();
  const [groupForm] = Form.useForm();
  const [couponForm] = Form.useForm();
  const [grantForm] = Form.useForm();

  useEffect(() => {
    void loadAll();
  }, []);

  const productSkuOptions = useMemo(
    () => products.flatMap((product) => (product.skus || []).map((sku) => ({ label: `${product.name} / ${sku.skuName} / ￥${sku.price}`, value: sku.id }))),
    [products],
  );

  const productSkuLabelMap = useMemo(() => {
    const map: Record<number, string> = {};
    productSkuOptions.forEach((item) => {
      map[Number(item.value)] = item.label;
    });
    return map;
  }, [productSkuOptions]);

  const menuDateOptions = useMemo(
    () =>
      [...dailyMenus]
        .sort((a, b) => a.date.localeCompare(b.date))
        .map((menu) => ({ label: `${menu.date} (${menu.title})`, value: menu.date })),
    [dailyMenus],
  );

  const menuByDateMap = useMemo(() => {
    const map: Record<string, DailyMenuItem> = {};
    dailyMenus.forEach((menu) => {
      map[menu.date] = menu;
    });
    return map;
  }, [dailyMenus]);

  const couponOptions = useMemo(
    () => coupons.map((coupon) => ({ label: `${coupon.name} / 满${coupon.thresholdAmount}减${coupon.amount}`, value: coupon.id })),
    [coupons],
  );

  const stats = useMemo(
    () => [
      { title: '商品数', value: products.length },
      { title: '每日菜单数', value: dailyMenus.length },
      { title: '团活动数', value: groups.length },
      { title: '订单数', value: orders.length },
    ],
    [dailyMenus.length, groups.length, orders.length, products.length],
  );

  async function handleImageChange(form: any, fieldName: string, file?: File, type: 'product' | 'group' = 'product') {
    if (!file) {
      form.setFieldValue(fieldName, '');
      return;
    }
    if (!file.type.startsWith('image/')) {
      throw new Error('请选择 jpg、png、webp 等图片文件');
    }
    if (file.size > 8 * 1024 * 1024) {
      throw new Error('图片原文件不能超过 8MB');
    }

    const setUploading = type === 'product' ? setProductImageUploading : setGroupImageUploading;
    setUploading(true);
    try {
      const compressed = await compressImage(file);
      const uploadFile = await dataUrlToFile(compressed, file.name || `image-${Date.now()}.jpg`);
      const uploaded = await uploadAdminImage(uploadFile);
      form.setFieldValue(fieldName, uploaded.url);
    } finally {
      setUploading(false);
    }
  }

  function clearImage(form: any, fieldName: string) {
    form.setFieldValue(fieldName, '');
  }

  function openPreview(image?: string) {
    if (!image) return;
    setPreviewImage(image);
    setPreviewOpen(true);
  }

  async function loadAll() {
    setLoading(true);
    setError('');
    try {
      const [productRes, groupRes, orderRes, couponRes, menuRes] = await Promise.all([fetchAdminProducts(), fetchAdminGroups(), fetchAdminOrders(), fetchAdminCoupons(), fetchAdminDailyMenus()]);
      setProducts(productRes.list);
      setGroups(groupRes.list);
      setOrders(orderRes.list);
      setCoupons(couponRes.list);
      setDailyMenus(
        (menuRes.list || []).map((menu) => ({
          id: menu.id,
          date: String(menu.menuDate).slice(0, 10),
          title: menu.title,
          items: (menu.items || []).map((item) => ({
            productSkuId: Number(item.productSkuId || 0),
            stockTotal: Number(item.stockTotal || 0),
            price: Number(item.price || 0),
            originalPrice: Number(item.originalPrice || 0),
            limitPerUser: Number(item.limitPerUser || 0),
            limitPerOrder: Number(item.limitPerOrder || 0),
            sortOrder: Number(item.sortOrder || 0),
          })),
        })),
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载失败');
    } finally {
      setLoading(false);
    }
  }

  async function submitProduct(values: any) {
    if (productImageUploading) {
      message.warning('商品图片还在上传中，请稍候再保存');
      return;
    }
    setSubmitting(true);
    try {
      await createProduct({
        name: values.name,
        coverImage: values.coverImage,
        categoryName: values.categoryName,
        subtitle: values.subtitle,
        skus: [{ skuName: values.skuName, skuCode: values.skuCode, price: values.price, originalPrice: values.originalPrice || 0, stockTotal: values.stockTotal, limitPerUser: 2, limitPerOrder: 4 }],
      });
      message.success('商品创建成功');
      productForm.resetFields();
      setProductOpen(false);
      await loadAll();
    } catch (err) {
      notifyError(err, '商品创建失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function submitGroup(values: any) {
    if (groupImageUploading) {
      message.warning('团活动图片还在上传中，请稍候再保存');
      return;
    }
    const linkedMenu = menuByDateMap[values.menuDate];
    if (!linkedMenu) {
      message.warning('请选择有效的每日菜单日期');
      return;
    }
    setSubmitting(true);
    try {
      await createGroup({
        title: values.title,
        coverImage: values.coverImage,
        leaderUserId: 1,
        fulfillmentMode: values.fulfillmentMode,
        startAt: values.startAt.toDate().toISOString(),
        cutoffAt: values.cutoffAt.toDate().toISOString(),
        allowModifyBeforeCutoff: true,
        showJoinList: true,
        pickupRuleDesc: '请按时取餐',
        menuDate: linkedMenu.date,
        groupNotice: `关联菜单日期：${linkedMenu.date}`,
        items: linkedMenu.items.map((item, index) => ({
          productSkuId: item.productSkuId,
          stockTotal: item.stockTotal,
          price: item.price,
          originalPrice: item.originalPrice || 0,
          limitPerUser: item.limitPerUser || 0,
          limitPerOrder: item.limitPerOrder || 0,
          sortOrder: item.sortOrder || index + 1,
        })),
      });
      message.success('团活动创建成功');
      groupForm.resetFields();
      setGroupOpen(false);
      await loadAll();
    } catch (err) {
      notifyError(err, '团活动创建失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function submitDailyMenu(values: any) {
    setSubmitting(true);
    try {
      const items = (values.items || []).map((item: any, index: number) => ({
        productSkuId: Number(item.productSkuId),
        stockTotal: Number(item.stockTotal),
        price: Number(item.price),
        originalPrice: Number(item.originalPrice || 0),
        limitPerUser: Number(item.limitPerUser || 2),
        limitPerOrder: Number(item.limitPerOrder || 4),
        sortOrder: Number(item.sortOrder || index + 1),
      }));
      if (!items.length) {
        message.warning('请至少添加一个菜单商品');
        return;
      }

      await createDailyMenu({
        menuDate: formatDateOnly(values.date),
        title: values.title,
        items,
      });
      menuForm.resetFields();
      setMenuOpen(false);
      message.success('每日菜单已保存');
      await loadAll();
    } catch (err) {
      notifyError(err, '每日菜单保存失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function deleteDailyMenu(id: number) {
    setSubmitting(true);
    try {
      await deleteDailyMenuApi(id);
      message.success('每日菜单已删除');
      await loadAll();
    } catch (err) {
      notifyError(err, '每日菜单删除失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function submitCoupon(values: any) {
    setSubmitting(true);
    try {
      await createCoupon({
        name: values.name,
        couponType: values.couponType,
        amount: values.amount,
        thresholdAmount: values.thresholdAmount,
        validDays: values.validDays,
        totalCount: values.totalCount,
        perUserLimit: values.perUserLimit,
      });
      message.success('优惠券创建成功');
      couponForm.resetFields();
      setCouponOpen(false);
      await loadAll();
    } catch (err) {
      notifyError(err, '优惠券创建失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function submitGrant(values: any) {
    setSubmitting(true);
    try {
      await grantCoupon(values.couponId, values.userId);
      message.success('发券成功');
      grantForm.resetFields();
      setGrantOpen(false);
      if (Number(values.userId) === couponQueryUserId) {
        const result = await fetchAdminUserCoupons(couponQueryUserId, couponStatusFilter);
        setUserCoupons(result.list);
      }
    } catch (err) {
      notifyError(err, '发券失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function loadUserCoupons(userId: number, status = couponStatusFilter) {
    setSubmitting(true);
    try {
      const result = await fetchAdminUserCoupons(userId, status);
      setUserCoupons(result.list);
      setCouponQueryUserId(userId);
      setCouponStatusFilter(status);
    } catch (err) {
      notifyError(err, '查询用户优惠券失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDisableCoupon(couponId: number) {
    setSubmitting(true);
    try {
      await disableCoupon(couponId);
      message.success('优惠券已停用');
      await loadAll();
    } catch (err) {
      notifyError(err, '停用失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function doCutoff(groupId: number) {
    setSubmitting(true);
    try {
      await cutoffGroup(groupId);
      message.success('截单成功');
      await loadAll();
    } catch (err) {
      notifyError(err, '截单失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function showSummary(groupId: number) {
    setSubmitting(true);
    try {
      const result = await fetchGroupSummary(groupId);
      setSummaryRows(result.bySku);
      setSummaryOpen(true);
    } catch (err) {
      notifyError(err, '获取汇总失败');
    } finally {
      setSubmitting(false);
    }
  }

  async function doOrderAction(orderId: number, action: 'pickup' | 'delivery' | 'complete') {
    setSubmitting(true);
    try {
      if (action === 'pickup') await markOrderReadyForPickup(orderId);
      if (action === 'delivery') await markOrderDelivering(orderId);
      if (action === 'complete') await completeOrder(orderId);
      message.success('操作成功');
      await loadAll();
    } catch (err) {
      notifyError(err, '操作失败');
    } finally {
      setSubmitting(false);
    }
  }

  const selectedMenuDate = Form.useWatch('menuDate', groupForm);
  const selectedMenu = selectedMenuDate ? menuByDateMap[selectedMenuDate] : undefined;

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={220}>
        <div style={{ padding: 20, fontWeight: 700, fontSize: 20 }}>selfCook</div>
        <Menu mode="inline" selectedKeys={[activeTab]} onClick={(event) => setActiveTab(event.key as TabKey)} items={[{ key: 'dashboard', label: '仪表盘' }, { key: 'products', label: '商品管理' }, { key: 'menus', label: '每日菜单管理' }, { key: 'groups', label: '团活动管理' }, { key: 'orders', label: '订单管理' }, { key: 'coupons', label: '优惠券管理' }]} />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px' }}>
          <Space direction="vertical" size={0} style={{ paddingTop: 12 }}>
            <Title level={4} style={{ margin: 0 }}>团购客饭接龙管理后台</Title>
            <Text type="secondary">支持新建商品、发团、履约、截单、汇总与优惠券运营。</Text>
          </Space>
        </Header>
        <Content style={{ padding: 24, background: '#f5f5f5' }}>
          {error ? <Alert type="error" showIcon style={{ marginBottom: 16 }} message={error} /> : null}
          {activeTab === 'dashboard' ? <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 16 }}>{stats.map((item) => <Card key={item.title} loading={loading}><Text type="secondary">{item.title}</Text><Title level={2} style={{ marginTop: 12, marginBottom: 0 }}>{item.value}</Title></Card>)}</div> : null}
          {activeTab === 'products' ? <Card title="商品管理" extra={<Button type="primary" onClick={() => setProductOpen(true)}>新建商品</Button>}><Table<Product> loading={loading} rowKey="id" dataSource={products} pagination={false} columns={[{ title: 'ID', dataIndex: 'id', width: 80 }, { title: '图片', dataIndex: 'coverImage', width: 88, render: (value?: string) => value ? <Image src={value} alt="商品图片" width={56} height={56} style={{ objectFit: 'cover', borderRadius: 8, cursor: 'pointer' }} preview={false} onClick={() => openPreview(value)} /> : '-' }, { title: '商品名称', dataIndex: 'name' }, { title: '类目', dataIndex: 'categoryName' }, { title: '状态', dataIndex: 'status', render: (value: string) => <Tag color={getProductStatusColor(value)}>{getLabel(productStatusLabelMap, value)}</Tag> }, { title: 'SKU 概览', render: (_, record) => (record.skus || []).map((sku) => `${sku.skuName} / ￥${sku.price}`).join('；') }]} /></Card> : null}
          {activeTab === 'menus' ? <Card title="每日菜单管理" extra={<Button type="primary" onClick={() => setMenuOpen(true)}>新建每日菜单</Button>}><Table<DailyMenuItem> rowKey="id" dataSource={[...dailyMenus].sort((a, b) => a.date.localeCompare(b.date))} pagination={false} columns={[{ title: '日期', dataIndex: 'date', width: 120 }, { title: '菜单标题', dataIndex: 'title', width: 160 }, { title: '菜品数', render: (_, record) => record.items.length, width: 90 }, { title: '菜单明细', render: (_, record) => record.items.map((item) => `${productSkuLabelMap[item.productSkuId] || `SKU#${item.productSkuId}`} x库存${item.stockTotal}`).join('；') }, { title: '操作', width: 100, render: (_, record) => <Button danger size="small" onClick={() => void deleteDailyMenu(record.id)}>删除</Button> }]} /></Card> : null}
          {activeTab === 'groups' ? <Card title="团活动管理" extra={<Button type="primary" onClick={() => setGroupOpen(true)}>新建团活动</Button>}><Table<Group> loading={loading} rowKey="id" dataSource={groups} pagination={false} columns={[{ title: 'ID', dataIndex: 'id', width: 80 }, { title: '图片', dataIndex: 'coverImage', width: 88, render: (value?: string) => value ? <Image src={value} alt="团活动图片" width={56} height={56} style={{ objectFit: 'cover', borderRadius: 8, cursor: 'pointer' }} preview={false} onClick={() => openPreview(value)} /> : '-' }, { title: '团名称', dataIndex: 'title' }, { title: '状态', dataIndex: 'status', render: (value: string) => <Tag color={getGroupStatusColor(value)}>{getLabel(groupStatusLabelMap, value)}</Tag> }, { title: '履约方式', dataIndex: 'fulfillmentMode', render: (value: string) => getLabel(fulfillmentModeLabelMap, value) }, { title: '截单时间', dataIndex: 'cutoffAt' }, { title: '商品数', render: (_, record) => record.items?.length || 0 }, { title: '操作', render: (_, record) => <Space><Button size="small" onClick={() => void showSummary(record.id)}>汇总</Button><Button size="small" type="primary" disabled={record.status !== GroupStatus.Ongoing} onClick={() => void doCutoff(record.id)}>截单</Button></Space> }]} /></Card> : null}
          {activeTab === 'orders' ? <Card title="订单管理"><Table<Order> loading={loading} rowKey="id" dataSource={orders} pagination={false} columns={[{ title: '订单号', dataIndex: 'orderNo' }, { title: '团ID', dataIndex: 'groupId', width: 80 }, { title: '联系人', dataIndex: 'contactName' }, { title: '手机号', dataIndex: 'contactPhone' }, { title: '履约方式', dataIndex: 'fulfillmentMode', width: 100, render: (value: string) => getLabel(fulfillmentModeLabelMap, value) }, { title: '状态', dataIndex: 'status', render: (value: string) => <Tag color={getOrderStatusColor(value)}>{getLabel(orderStatusLabelMap, value)}</Tag> }, { title: '金额', dataIndex: 'paidAmount', width: 80 }, { title: '操作', render: (_, record) => <Space wrap><Button size="small" disabled={record.status !== OrderStatus.CutoffLocked} onClick={() => void doOrderAction(record.id, 'pickup')}>待取餐</Button><Button size="small" disabled={record.status !== OrderStatus.CutoffLocked} onClick={() => void doOrderAction(record.id, 'delivery')}>配送中</Button><Button size="small" type="primary" disabled={record.status !== OrderStatus.ReadyForPickup && record.status !== OrderStatus.Delivering} onClick={() => void doOrderAction(record.id, 'complete')}>完成</Button></Space> }]} /></Card> : null}
          {activeTab === 'coupons' ? <Card title="优惠券管理" extra={<Space><Button onClick={() => void loadUserCoupons(couponQueryUserId, couponStatusFilter)}>查询用户券</Button><Button onClick={() => setGrantOpen(true)} disabled={!coupons.length}>发券</Button><Button type="primary" onClick={() => setCouponOpen(true)}>新建优惠券</Button></Space>}><Space direction="vertical" style={{ width: '100%' }} size={16}><Table<Coupon> loading={loading} rowKey="id" dataSource={coupons} pagination={false} columns={[{ title: 'ID', dataIndex: 'id', width: 80 }, { title: '券名称', dataIndex: 'name' }, { title: '类型', dataIndex: 'couponType', render: (value: string) => getLabel(couponTypeLabelMap, value) }, { title: '面额', dataIndex: 'amount' }, { title: '门槛', dataIndex: 'thresholdAmount' }, { title: '已发放', dataIndex: 'grantedCount', width: 90 }, { title: '已使用', dataIndex: 'usedCount', width: 90 }, { title: '使用率', render: (_, record) => `${record.grantedCount ? Math.round(((record.usedCount || 0) / record.grantedCount) * 100) : 0}%` }, { title: '状态', dataIndex: 'status', render: (value: string) => <Tag color={getCouponStatusColor(value)}>{getLabel(couponStatusLabelMap, value)}</Tag> }, { title: '有效期至', dataIndex: 'validTo' }, { title: '操作', render: (_, record) => <Button size="small" disabled={record.status !== 'active'} onClick={() => void handleDisableCoupon(record.id)}>停用</Button> }]} /><Card type="inner" title={`用户 ${couponQueryUserId} 的优惠券明细`} extra={<Space wrap><Text>状态</Text><Select value={couponStatusFilter} style={{ width: 140 }} options={[{ label: '全部', value: 'all' }, { label: '待使用', value: 'unused' }, { label: '已使用', value: 'used' }, { label: '已过期', value: 'expired' }]} onChange={(value) => setCouponStatusFilter(value)} /><Text>用户 ID</Text><InputNumber min={1} value={couponQueryUserId} onChange={(value) => setCouponQueryUserId(Number(value || 1))} /><Button onClick={() => void loadUserCoupons(couponQueryUserId, couponStatusFilter)}>查询</Button></Space>}><Table<UserCoupon> rowKey="id" dataSource={userCoupons} pagination={false} columns={[{ title: '记录ID', dataIndex: 'id', width: 90 }, { title: '券名称', render: (_, record) => record.coupon?.name || `券#${record.couponId}` }, { title: '状态', dataIndex: 'status', render: (value: string) => <Tag color={getUserCouponStatusColor(value)}>{getLabel(userCouponStatusLabelMap, value)}</Tag> }, { title: '领取时间', dataIndex: 'acquiredAt' }, { title: '使用时间', render: (_, record) => record.usedAt || '-' }, { title: '订单ID', render: (_, record) => record.orderId || '-' }, { title: '有效期至', dataIndex: 'validTo' }]} /></Card></Space></Card> : null}
        </Content>
      </Layout>
      <Drawer title="新建商品" open={productOpen} width={420} onClose={() => setProductOpen(false)}><Form form={productForm} layout="vertical" onFinish={(values) => void submitProduct(values)}><Form.Item name="name" label="商品名称" rules={[{ required: true }]}><Input /></Form.Item><Form.Item name="categoryName" label="类目"><Input /></Form.Item><Form.Item name="subtitle" label="副标题"><Input /></Form.Item><Form.Item name="coverImage" label="商品图片" rules={[{ required: true, message: '请上传或填写商品图片' }]}><Input.TextArea rows={3} placeholder="可粘贴图片地址，或使用下方本地上传" /></Form.Item><Form.Item label="本地上传" extra="支持 jpg/png/webp，原图最大 8MB；上传前会自动压缩优化"><Space direction="vertical" style={{ width: '100%' }} size={12}><input type="file" accept="image/jpeg,image/png,image/webp,image/*" disabled={productImageUploading} onChange={async (event) => { const file = event.target.files?.[0]; try { await handleImageChange(productForm, 'coverImage', file, 'product'); message.success(file ? '图片已上传' : '已清空图片'); } catch (err) { notifyError(err, '图片处理失败'); } finally { event.currentTarget.value = ''; } }} />{productImageUploading ? <Text type="secondary">图片上传中，请稍候...</Text> : null}{productForm.getFieldValue('coverImage') ? <Space direction="vertical" size={12}><div style={{ marginTop: 4 }}><Image src={productForm.getFieldValue('coverImage')} alt="商品预览" width={96} height={96} style={{ objectFit: 'cover', borderRadius: 8, cursor: 'pointer' }} preview={false} onClick={() => openPreview(productForm.getFieldValue('coverImage'))} /></div><Space><Button size="small" onClick={() => openPreview(productForm.getFieldValue('coverImage'))}>放大预览</Button><Button size="small" onClick={() => window.open(productForm.getFieldValue('coverImage'), '_blank')}>打开原图</Button><Button size="small" danger onClick={() => clearImage(productForm, 'coverImage')}>清空图片</Button></Space></Space> : null}</Space></Form.Item><Form.Item name="skuName" label="SKU 名称" rules={[{ required: true }]}><Input /></Form.Item><Form.Item name="skuCode" label="SKU 编码" rules={[{ required: true }]}><Input /></Form.Item><Form.Item name="price" label="售价" rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name="originalPrice" label="原价"><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name="stockTotal" label="库存" rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item><Button htmlType="submit" type="primary" loading={submitting || productImageUploading} disabled={productImageUploading} block>保存商品</Button></Form></Drawer>
      <Drawer title="新建每日菜单" open={menuOpen} width={560} onClose={() => setMenuOpen(false)}><Form form={menuForm} layout="vertical" onFinish={(values) => void submitDailyMenu(values)} initialValues={{ items: [{ limitPerUser: 2, limitPerOrder: 4 }] }}><Form.Item name="date" label="菜单日期" rules={[{ required: true, message: '请选择日期' }]}><DatePicker style={{ width: '100%' }} format="YYYY-MM-DD" /></Form.Item><Form.Item name="title" label="菜单标题"><Input placeholder="例如：周一午餐菜单" /></Form.Item><Form.List name="items">{(fields, { add, remove }) => (<Space direction="vertical" style={{ width: '100%' }} size={12}>{fields.map((field, index) => <Card key={field.key} size="small" title={`菜单商品 ${index + 1}`} extra={fields.length > 1 ? <Button danger size="small" onClick={() => remove(field.name)}>删除</Button> : null}><Form.Item name={[field.name, 'productSkuId']} label="关联商品 SKU" rules={[{ required: true, message: '请选择 SKU' }]}><Select options={productSkuOptions} showSearch optionFilterProp="label" /></Form.Item><Form.Item name={[field.name, 'stockTotal']} label="当日库存" rules={[{ required: true, message: '请输入库存' }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item><Form.Item name={[field.name, 'price']} label="当日售价" rules={[{ required: true, message: '请输入售价' }]}><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name={[field.name, 'originalPrice']} label="原价"><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name={[field.name, 'limitPerUser']} label="每人限购" initialValue={2}><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name={[field.name, 'limitPerOrder']} label="每单限购" initialValue={4}><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name={[field.name, 'sortOrder']} label="排序"><InputNumber style={{ width: '100%' }} min={1} /></Form.Item></Card>)}<Button onClick={() => add({ limitPerUser: 2, limitPerOrder: 4 })} block>新增菜单商品</Button></Space>)}</Form.List><Button htmlType="submit" type="primary" loading={submitting} block>保存每日菜单</Button></Form></Drawer>
      <Drawer title="新建团活动" open={groupOpen} width={460} onClose={() => setGroupOpen(false)}><Form form={groupForm} layout="vertical" onFinish={(values) => void submitGroup(values)}><Form.Item name="title" label="团名称" rules={[{ required: true }]}><Input /></Form.Item><Form.Item name="coverImage" label="团活动图片" rules={[{ required: true, message: '请上传或填写团活动图片' }]}><Input.TextArea rows={3} placeholder="可粘贴图片地址，或使用下方本地上传" /></Form.Item><Form.Item label="本地上传" extra="支持 jpg/png/webp，原图最大 8MB；上传前会自动压缩优化"><Space direction="vertical" style={{ width: '100%' }} size={12}><input type="file" accept="image/jpeg,image/png,image/webp,image/*" disabled={groupImageUploading} onChange={async (event) => { const file = event.target.files?.[0]; try { await handleImageChange(groupForm, 'coverImage', file, 'group'); message.success(file ? '图片已上传' : '已清空图片'); } catch (err) { notifyError(err, '图片处理失败'); } finally { event.currentTarget.value = ''; } }} />{groupImageUploading ? <Text type="secondary">图片上传中，请稍候...</Text> : null}{groupForm.getFieldValue('coverImage') ? <Space direction="vertical" size={12}><div style={{ marginTop: 4 }}><Image src={groupForm.getFieldValue('coverImage')} alt="团活动预览" width={96} height={96} style={{ objectFit: 'cover', borderRadius: 8, cursor: 'pointer' }} preview={false} onClick={() => openPreview(groupForm.getFieldValue('coverImage'))} /></div><Space><Button size="small" onClick={() => openPreview(groupForm.getFieldValue('coverImage'))}>放大预览</Button><Button size="small" onClick={() => window.open(groupForm.getFieldValue('coverImage'), '_blank')}>打开原图</Button><Button size="small" danger onClick={() => clearImage(groupForm, 'coverImage')}>清空图片</Button></Space></Space> : null}</Space></Form.Item><Form.Item name="fulfillmentMode" label="履约方式" initialValue="pickup" rules={[{ required: true }]}><Select options={[{ label: '自提', value: 'pickup' }, { label: '配送', value: 'delivery' }, { label: '混合', value: 'mixed' }]} /></Form.Item><Form.Item name="menuDate" label="关联每日菜单日期" rules={[{ required: true, message: '请选择每日菜单日期' }]} extra={dailyMenus.length ? '团活动会自动使用该日期菜单中的全部商品、库存与价格' : '请先在每日菜单管理中创建菜单'}><Select options={menuDateOptions} disabled={!dailyMenus.length} /></Form.Item>{selectedMenu ? <Card size="small" title="菜单快照" style={{ marginBottom: 12 }}><Space direction="vertical" size={4}><Text>日期：{selectedMenu.date}</Text><Text>菜品数：{selectedMenu.items.length}</Text><Text>明细：{selectedMenu.items.map((item) => `${productSkuLabelMap[item.productSkuId] || `SKU#${item.productSkuId}`} x库存${item.stockTotal}`).join('；')}</Text></Space></Card> : null}<Form.Item name="startAt" label="开始时间" rules={[{ required: true }]}><DatePicker showTime format="YYYY-MM-DD HH:mm:ss" inputReadOnly style={{ width: '100%' }} placeholder="请选择开始时间" /></Form.Item><Form.Item name="cutoffAt" label="截单时间" rules={[{ required: true }]}><DatePicker showTime format="YYYY-MM-DD HH:mm:ss" inputReadOnly style={{ width: '100%' }} placeholder="请选择截单时间" /></Form.Item><Button htmlType="submit" type="primary" loading={submitting || groupImageUploading} disabled={groupImageUploading} block>保存团活动</Button></Form></Drawer>
      <Drawer title="新建优惠券" open={couponOpen} width={420} onClose={() => setCouponOpen(false)}><Form form={couponForm} layout="vertical" onFinish={(values) => void submitCoupon(values)}><Form.Item name="name" label="券名称" rules={[{ required: true }]}><Input /></Form.Item><Form.Item name="couponType" label="券类型" initialValue="full_reduction" rules={[{ required: true }]}><Select options={[{ label: '满减券', value: 'full_reduction' }]} /></Form.Item><Form.Item name="amount" label="面额" rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name="thresholdAmount" label="使用门槛"><InputNumber style={{ width: '100%' }} min={0} /></Form.Item><Form.Item name="validDays" label="有效天数" initialValue={30}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item><Form.Item name="totalCount" label="总量" initialValue={1000}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item><Form.Item name="perUserLimit" label="每人限领" initialValue={1}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item><Button htmlType="submit" type="primary" loading={submitting} block>保存优惠券</Button></Form></Drawer>
      <Modal title="发券给用户" open={grantOpen} footer={null} onCancel={() => setGrantOpen(false)}><Form form={grantForm} layout="vertical" onFinish={(values) => void submitGrant(values)}><Form.Item name="couponId" label="选择优惠券" rules={[{ required: true }]}><Select options={couponOptions} /></Form.Item><Form.Item name="userId" label="用户 ID" initialValue={1} rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} min={1} /></Form.Item><Button htmlType="submit" type="primary" loading={submitting} block>确认发券</Button></Form></Modal>
      <Modal title="商品汇总" open={summaryOpen} footer={null} onCancel={() => setSummaryOpen(false)}><Table<SummaryRow> rowKey={(record) => `${record.productName}-${record.skuName}`} dataSource={summaryRows} pagination={false} columns={[{ title: '商品', dataIndex: 'productName' }, { title: '规格', dataIndex: 'skuName' }, { title: '数量', dataIndex: 'totalQty', width: 100 }]} /></Modal>
      <Modal open={previewOpen} footer={null} onCancel={() => setPreviewOpen(false)} closable width="auto" styles={{ body: { padding: 12, textAlign: 'center' } }}><Image src={previewImage} alt="预览图片" style={{ maxWidth: '80vw', maxHeight: '80vh', objectFit: 'contain' }} preview={false} /></Modal>
    </Layout>
  );
}
