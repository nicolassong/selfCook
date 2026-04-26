const GROUP_STATUS = {
  ongoing: 'ongoing',
  cutoffLocked: 'cutoff_locked',
  completed: 'completed',
  cancelled: 'cancelled'
};

const ORDER_STATUS = {
  joined: 'joined',
  cutoffLocked: 'cutoff_locked',
  readyForPickup: 'ready_for_pickup',
  delivering: 'delivering',
  completed: 'completed',
  cancelled: 'cancelled'
};

module.exports = {
  GROUP_STATUS,
  ORDER_STATUS
};
