/**
 * 时间格式化工具函数集合
 * 统一管理所有时间相关的格式化逻辑
 */
import i18n from '../i18n';

const getLocale = () => i18n.language || 'en-US';

/**
 * 格式化时间字符串为本地化格式
 * @param {string|Date} timeStr - ISO时间字符串或Date对象
 * @param {Object} options - 格式化选项
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(timeStr, options = {}) {
  if (!timeStr) return '-';

  const defaultOptions = {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  };

  return new Date(timeStr).toLocaleString(getLocale(), {
    ...defaultOptions,
    ...options,
  });
}

/**
 * 格式化日期（不包含时间）
 * @param {string|Date} timeStr - ISO时间字符串或Date对象
 * @returns {string} 格式化后的日期字符串 (YYYY-MM-DD)
 */
export function formatDate(timeStr) {
  if (!timeStr) return '-';

  return new Date(timeStr).toLocaleString(getLocale(), {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

/**
 * 格式化剩余时间（秒数）
 * @param {number} remaining - 剩余秒数
 * @returns {string} 格式化后的剩余时间 (例如: "2h 30m 15s")
 */
export function formatRemaining(remaining) {
  if (!remaining || remaining <= 0) return i18n.t('utils.time.stopped');

  const hours = Math.floor(remaining / 3600);
  const minutes = Math.floor((remaining % 3600) / 60);
  const seconds = Math.floor(remaining % 60);

  const parts = [];
  if (hours > 0) {
    parts.push(i18n.t('utils.time.units.hour', { count: hours }));
  }
  if (minutes > 0 || hours > 0) {
    parts.push(i18n.t('utils.time.units.minute', { count: minutes }));
  }
  parts.push(i18n.t('utils.time.units.second', { count: seconds }));

  return parts.join(' ');
}

/**
 * 格式化持续时间（从纳秒转换）
 * @param {number} durationNs - 持续时间（纳秒）
 * @returns {string} 格式化后的持续时间
 */
export function formatDuration(durationNs) {
  if (!durationNs) return '-';

  // 纳秒转换为秒
  const seconds = durationNs / 1000000000;

  if (seconds < 60) {
    return i18n.t('utils.time.units.second', { count: Math.round(seconds) });
  } else if (seconds < 3600) {
    return `${i18n.t('utils.time.units.minute', { count: Math.floor(seconds / 60) })}${i18n.t(
      'utils.time.units.second',
      { count: Math.round(seconds % 60) }
    )}`;
  } else {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${i18n.t('utils.time.units.hour', { count: hours })}${minutes > 0 ? i18n.t('utils.time.units.minute', { count: minutes }) : ''}`;
  }
}

/**
 * 计算倒计时（返回对象）
 * @param {string|Date} targetTime - 目标时间
 * @returns {Object} { days, hours, minutes, seconds, isExpired }
 */
export function calculateTimeLeft(targetTime) {
  if (!targetTime) {
    return { days: 0, hours: 0, minutes: 0, seconds: 0, isExpired: true };
  }

  const now = new Date().getTime();
  const target = new Date(targetTime).getTime();
  const diff = target - now;

  if (diff <= 0) {
    return { days: 0, hours: 0, minutes: 0, seconds: 0, isExpired: true };
  }

  return {
    days: Math.floor(diff / (1000 * 60 * 60 * 24)),
    hours: Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
    minutes: Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60)),
    seconds: Math.floor((diff % (1000 * 60)) / 1000),
    isExpired: false,
  };
}

/**
 * 格式化倒计时为字符串
 * @param {string|Date} targetTime - 目标时间
 * @returns {string} 格式化后的倒计时 (例如: "2天 5小时 30分钟")
 */
export function formatTimeLeft(targetTime) {
  const { days, hours, minutes, seconds, isExpired } = calculateTimeLeft(targetTime);

  if (isExpired) {
    return i18n.t('utils.time.ended');
  }

  const parts = [];
  if (days > 0) parts.push(i18n.t('utils.time.units.day', { count: days }));
  if (hours > 0) parts.push(i18n.t('utils.time.units.hour', { count: hours }));
  if (minutes > 0) parts.push(i18n.t('utils.time.units.minute', { count: minutes }));
  if (seconds > 0 && days === 0) parts.push(i18n.t('utils.time.units.second', { count: seconds }));

  return parts.join(' ') || i18n.t('utils.time.startingSoon');
}

/**
 * 获取相对时间描述
 * @param {string|Date} timeStr - 时间字符串或Date对象
 * @returns {string} 相对时间描述 (例如: "5分钟前", "刚刚")
 */
export function getRelativeTime(timeStr) {
  if (!timeStr) return '-';

  const now = new Date().getTime();
  const time = new Date(timeStr).getTime();
  const diff = now - time;

  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (seconds < 60) return i18n.t('utils.time.justNow');
  if (minutes < 60) return i18n.t('utils.time.ago.minute', { count: minutes });
  if (hours < 24) return i18n.t('utils.time.ago.hour', { count: hours });
  if (days < 7) return i18n.t('utils.time.ago.day', { count: days });

  // 超过7天显示完整日期
  return formatDate(timeStr);
}

/**
 * 检查时间是否已过期
 * @param {string|Date} timeStr - 时间字符串或Date对象
 * @returns {boolean} 是否已过期
 */
export function isExpired(timeStr) {
  if (!timeStr) return true;
  return new Date(timeStr).getTime() < new Date().getTime();
}

/**
 * 检查时间是否在指定范围内
 * @param {string|Date} timeStr - 要检查的时间
 * @param {string|Date} startTime - 开始时间
 * @param {string|Date} endTime - 结束时间
 * @returns {boolean} 是否在范围内
 */
export function isInRange(timeStr, startTime, endTime) {
  if (!timeStr || !startTime || !endTime) return false;

  const time = new Date(timeStr).getTime();
  const start = new Date(startTime).getTime();
  const end = new Date(endTime).getTime();

  return time >= start && time <= end;
}
