import { getContestStatus, getContestTimeRange } from '../../../../../config/contest.js';

export function createContestDraft(data) {
  if (!data)
    return {
      title: '',
      description: '',
      image: '',
      status: 'upcoming',
      startTime: '',
      endTime: '',
      participants: 0,
      rules: [],
      prizes: [],
      timeline: [],
      prefix: 'CBCTF',
      size: 4,
      hidden: false,
      captcha: '',
      blood: true,
      victims: 1,
    };
  return {
    title: data.name,
    description: data.description,
    image: data.picture,
    status: getContestStatus(data.start, data.duration),
    startTime: data.start,
    endTime: getContestTimeRange(data.start, data.duration).endTime,
    participants: data.users,
    rules: [...(data.rules || [])],
    prizes: (data.prizes || []).map(({ amount, description }) => ({ amount, description })),
    timeline: (data.timelines || []).map(({ date, title, description }) => ({ date, title, description })),
    prefix: data.prefix,
    size: data.size,
    hidden: data.hidden,
    captcha: data.captcha,
    blood: data.blood,
    victims: data.victims,
  };
}

export function contestUpdatePayload(draft) {
  return {
    name: draft.title,
    description: draft.description,
    prefix: draft.prefix,
    size: draft.size,
    hidden: draft.hidden,
    captcha: draft.captcha,
    blood: draft.blood,
    victims: draft.victims,
    start: new Date(draft.startTime).toISOString(),
    duration: Math.floor((new Date(draft.endTime) - new Date(draft.startTime)) / 1000),
    rules: [...draft.rules],
    prizes: draft.prizes.map(({ amount, description }) => ({ amount, description })),
    timelines: draft.timeline.map(({ date, title, description }) => ({
      date: date ? new Date(date).toISOString() : '',
      title,
      description,
    })),
  };
}

export function validateContestDraft(draft) {
  const errors = {};
  if (!draft.title?.trim()) errors.title = 'titleRequired';
  for (const field of ['startTime', 'endTime']) {
    if (!draft[field] || Number.isNaN(new Date(draft[field]).getTime())) errors[field] = 'invalidDate';
  }
  if (!errors.startTime && !errors.endTime && new Date(draft.startTime) >= new Date(draft.endTime)) {
    errors.endTime = 'endTimeAfterStart';
  }
  return errors;
}

export function formatDateForInput(value) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  const pad = (part) => String(part).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}
