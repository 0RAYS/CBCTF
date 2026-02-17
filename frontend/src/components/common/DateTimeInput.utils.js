export const MONTH_NAMES = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
];

export const DAY_HEADERS = ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su'];

export function padZero(num) {
  return String(num).padStart(2, '0');
}

export function formatToDateTimeLocal(date) {
  if (!(date instanceof Date) || isNaN(date.getTime())) return '';
  const year = date.getFullYear();
  const month = padZero(date.getMonth() + 1);
  const day = padZero(date.getDate());
  const hours = padZero(date.getHours());
  const minutes = padZero(date.getMinutes());
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

export function parseDateTimeString(str) {
  if (!str || typeof str !== 'string') return null;
  const trimmed = str.trim();
  if (!trimmed) return null;

  // Try full date+time: YYYY-MM-DD HH:MM, YYYY/MM/DD HH:MM, YYYY.MM.DD HH:MM, or with T separator
  const fullMatch = trimmed.match(/^(\d{4})[-/.](\d{1,2})[-/.](\d{1,2})[\sT](\d{1,2}):(\d{1,2})$/);
  if (fullMatch) {
    const [, y, m, d, h, min] = fullMatch.map(Number);
    const date = new Date(y, m - 1, d, h, min);
    if (date.getFullYear() === y && date.getMonth() === m - 1 && date.getDate() === d) {
      return date;
    }
    return null;
  }

  // Try date only: YYYY-MM-DD, YYYY/MM/DD, YYYY.MM.DD
  const dateOnly = trimmed.match(/^(\d{4})[-/.](\d{1,2})[-/.](\d{1,2})$/);
  if (dateOnly) {
    const [, y, m, d] = dateOnly.map(Number);
    const date = new Date(y, m - 1, d, 0, 0);
    if (date.getFullYear() === y && date.getMonth() === m - 1 && date.getDate() === d) {
      return date;
    }
    return null;
  }

  return null;
}

export function getDaysInMonth(year, month) {
  return new Date(year, month + 1, 0).getDate();
}

export function getCalendarGrid(year, month) {
  const firstDay = new Date(year, month, 1);
  // getDay() returns 0=Sun..6=Sat, convert to Mon=0..Sun=6
  let startDow = firstDay.getDay() - 1;
  if (startDow < 0) startDow = 6;

  const daysInMonth = getDaysInMonth(year, month);
  const prevMonth = month === 0 ? 11 : month - 1;
  const prevYear = month === 0 ? year - 1 : year;
  const daysInPrevMonth = getDaysInMonth(prevYear, prevMonth);

  const grid = [];

  // Previous month days
  for (let i = startDow - 1; i >= 0; i--) {
    grid.push({
      day: daysInPrevMonth - i,
      month: prevMonth,
      year: prevYear,
      isCurrentMonth: false,
    });
  }

  // Current month days
  for (let d = 1; d <= daysInMonth; d++) {
    grid.push({ day: d, month, year, isCurrentMonth: true });
  }

  // Next month days to fill 42 cells
  const nextMonth = month === 11 ? 0 : month + 1;
  const nextYear = month === 11 ? year + 1 : year;
  let nextDay = 1;
  while (grid.length < 42) {
    grid.push({ day: nextDay++, month: nextMonth, year: nextYear, isCurrentMonth: false });
  }

  return grid;
}

export function isSameDay(d1, d2) {
  if (!d1 || !d2) return false;
  return d1.getFullYear() === d2.getFullYear() && d1.getMonth() === d2.getMonth() && d1.getDate() === d2.getDate();
}

export function isToday(cell) {
  const now = new Date();
  return cell.year === now.getFullYear() && cell.month === now.getMonth() && cell.day === now.getDate();
}

export function generateHourOptions() {
  return Array.from({ length: 24 }, (_, i) => i);
}

export function generateMinuteOptions() {
  return Array.from({ length: 60 }, (_, i) => i);
}
