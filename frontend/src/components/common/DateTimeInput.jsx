import { useState, useMemo, useRef, useEffect, useCallback } from 'react';
import { AnimatePresence, motion } from 'motion/react';
import { IconCalendar, IconClock, IconChevronLeft, IconChevronRight } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import {
  MONTH_NAMES,
  DAY_HEADERS,
  formatToDateTimeLocal,
  parseDateTimeString,
  getCalendarGrid,
  isSameDay,
  isToday,
  generateHourOptions,
  generateMinuteOptions,
  padZero,
} from './DateTimeInput.utils';

function DateTimeInput({
  value,
  onChange,
  name,
  placeholder,
  disabled = false,
  error,
  size = 'md',
  fullWidth = true,
  className = '',
  ...rest
}) {
  const { t } = useTranslation();
  const containerRef = useRef(null);
  const hourListRef = useRef(null);
  const minuteListRef = useRef(null);

  const [inputText, setInputText] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [viewYear, setViewYear] = useState(() => new Date().getFullYear());
  const [viewMonth, setViewMonth] = useState(() => new Date().getMonth());

  const selectedDate = useMemo(() => parseDateTimeString(value), [value]);
  const selectedHour = selectedDate ? selectedDate.getHours() : null;
  const selectedMinute = selectedDate ? selectedDate.getMinutes() : null;
  const calendarGrid = useMemo(() => getCalendarGrid(viewYear, viewMonth), [viewYear, viewMonth]);

  // Sync inputText when value prop changes (adjust state during render)
  const [prevValue, setPrevValue] = useState(value);
  if (prevValue !== value) {
    setPrevValue(value);
    if (value) {
      const parsed = parseDateTimeString(value);
      if (parsed) {
        setInputText(formatToDateTimeLocal(parsed).replace('T', ' '));
      } else {
        setInputText(value);
      }
    } else {
      setInputText('');
    }
  }

  // Sync calendar view when selected date changes (adjust state during render)
  const [prevSelectedDate, setPrevSelectedDate] = useState(selectedDate);
  if (prevSelectedDate !== selectedDate) {
    setPrevSelectedDate(selectedDate);
    if (selectedDate) {
      setViewYear(selectedDate.getFullYear());
      setViewMonth(selectedDate.getMonth());
    }
  }

  // Scroll active time items into view when panel opens
  useEffect(() => {
    if (isOpen) {
      const timer = setTimeout(() => {
        [hourListRef, minuteListRef].forEach((ref) => {
          const active = ref.current?.querySelector('[data-active="true"]');
          if (active) {
            active.scrollIntoView({ block: 'center', behavior: 'instant' });
          }
        });
      }, 50);
      return () => clearTimeout(timer);
    }
  }, [isOpen]);

  const emitChange = useCallback(
    (dateTimeLocalValue) => {
      if (onChange) {
        onChange({ target: { name, value: dateTimeLocalValue } });
      }
    },
    [onChange, name]
  );

  // Click outside / Escape to close
  useEffect(() => {
    if (!isOpen) return;

    const handleClickOutside = (e) => {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setIsOpen(false);
      }
    };

    const handleEscape = (e) => {
      if (e.key === 'Escape') setIsOpen(false);
    };

    document.addEventListener('mousedown', handleClickOutside);
    document.addEventListener('keydown', handleEscape);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [isOpen]);

  const commitInput = useCallback(() => {
    const trimmed = inputText.trim();
    if (!trimmed) {
      emitChange('');
      return;
    }
    const parsed = parseDateTimeString(trimmed);
    if (parsed) {
      const formatted = formatToDateTimeLocal(parsed);
      emitChange(formatted);
      setInputText(formatted.replace('T', ' '));
    }
    // Invalid input: keep text as-is, don't emit
  }, [inputText, emitChange]);

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      commitInput();
    }
  };

  const handleDayClick = (cell) => {
    const hour = selectedHour ?? new Date().getHours();
    const minute = selectedMinute ?? 0;
    const date = new Date(cell.year, cell.month, cell.day, hour, minute);
    const formatted = formatToDateTimeLocal(date);
    emitChange(formatted);
    setInputText(formatted.replace('T', ' '));
  };

  const handleHourClick = (hour) => {
    if (!selectedDate) return;
    const date = new Date(selectedDate);
    date.setHours(hour);
    const formatted = formatToDateTimeLocal(date);
    emitChange(formatted);
    setInputText(formatted.replace('T', ' '));
  };

  const handleMinuteClick = (minute) => {
    if (!selectedDate) return;
    const date = new Date(selectedDate);
    date.setMinutes(minute);
    const formatted = formatToDateTimeLocal(date);
    emitChange(formatted);
    setInputText(formatted.replace('T', ' '));
  };

  const handleNow = () => {
    const now = new Date();
    now.setSeconds(0, 0);
    const formatted = formatToDateTimeLocal(now);
    emitChange(formatted);
    setInputText(formatted.replace('T', ' '));
    setViewYear(now.getFullYear());
    setViewMonth(now.getMonth());
  };

  const prevMonth = () => {
    if (viewMonth === 0) {
      setViewMonth(11);
      setViewYear((y) => y - 1);
    } else {
      setViewMonth((m) => m - 1);
    }
  };

  const nextMonth = () => {
    if (viewMonth === 11) {
      setViewMonth(0);
      setViewYear((y) => y + 1);
    } else {
      setViewMonth((m) => m + 1);
    }
  };

  const sizes = {
    sm: 'h-8 text-sm',
    md: 'h-10',
    lg: 'h-12 text-base',
  };

  const inputClasses = `
    bg-black/20 border rounded-md text-neutral-50 placeholder-neutral-500
    focus:outline-none transition-all duration-200 font-mono
    ${sizes[size] || sizes.md}
    pl-10 pr-10
    ${fullWidth ? 'w-full' : ''}
    ${error ? 'border-red-400 focus:border-red-400 focus:shadow-[0_0_15px_rgba(239,68,68,0.3)]' : 'border-neutral-300/30 focus:border-geek-400 focus:shadow-[0_0_15px_rgba(89,126,247,0.3)]'}
    ${disabled ? 'opacity-50 cursor-not-allowed bg-black/10' : ''}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  return (
    <div ref={containerRef} className={fullWidth ? 'w-full' : 'inline-block'}>
      <div className="relative">
        <div className="absolute left-3 top-1/2 transform -translate-y-1/2 text-neutral-400 pointer-events-none">
          <IconCalendar size={18} />
        </div>

        <input
          type="text"
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          onBlur={commitInput}
          onKeyDown={handleKeyDown}
          placeholder={placeholder || t('common.dateTimeInput.placeholder')}
          disabled={disabled}
          className={inputClasses}
          {...rest}
        />

        <div
          className="absolute right-3 top-1/2 transform -translate-y-1/2 text-neutral-400 cursor-pointer hover:text-geek-400 transition-colors"
          onClick={() => {
            if (!disabled) setIsOpen((v) => !v);
          }}
        >
          <IconClock size={18} />
        </div>
      </div>

      {error && <div className="mt-1 text-sm text-red-400">{error}</div>}

      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -4 }}
            transition={{ duration: 0.15 }}
            className="absolute z-50 mt-1 w-[320px] bg-neutral-900 border border-neutral-600 rounded-md shadow-lg overflow-hidden"
          >
            {/* Month navigation */}
            <div className="flex items-center justify-between px-3 py-2 border-b border-neutral-700">
              <button
                type="button"
                onClick={prevMonth}
                className="p-1 rounded hover:bg-neutral-700/80 text-neutral-400 hover:text-neutral-200 transition-colors"
              >
                <IconChevronLeft size={16} />
              </button>
              <span className="text-sm font-medium text-neutral-200">
                {MONTH_NAMES[viewMonth]} {viewYear}
              </span>
              <button
                type="button"
                onClick={nextMonth}
                className="p-1 rounded hover:bg-neutral-700/80 text-neutral-400 hover:text-neutral-200 transition-colors"
              >
                <IconChevronRight size={16} />
              </button>
            </div>

            {/* Day headers */}
            <div className="grid grid-cols-7 px-3 pt-2">
              {DAY_HEADERS.map((d) => (
                <div key={d} className="text-center text-xs text-neutral-500 py-1">
                  {d}
                </div>
              ))}
            </div>

            {/* Calendar grid */}
            <div className="grid grid-cols-7 px-3 pb-2">
              {calendarGrid.map((cell, i) => {
                const cellDate = new Date(cell.year, cell.month, cell.day);
                const selected = selectedDate && isSameDay(cellDate, selectedDate);
                const today = isToday(cell);

                return (
                  <button
                    key={i}
                    type="button"
                    onClick={() => handleDayClick(cell)}
                    className={`
                      text-center text-sm py-1 rounded transition-colors
                      ${!cell.isCurrentMonth ? 'text-neutral-600' : 'text-neutral-300'}
                      ${selected ? 'bg-geek-400/20 text-geek-300 border border-geek-400/30' : 'border border-transparent'}
                      ${today && !selected ? 'text-geek-400 font-semibold' : ''}
                      hover:bg-neutral-700/80
                    `
                      .trim()
                      .replace(/\s+/g, ' ')}
                  >
                    {cell.day}
                  </button>
                );
              })}
            </div>

            {/* Divider */}
            <div className="border-t border-neutral-700 mx-3" />

            {/* Time selection */}
            <div className="px-3 py-2">
              <div className="text-xs text-neutral-500 mb-1">{t('common.dateTimeInput.time')}</div>
              <div className="flex items-center gap-1">
                {/* Hours */}
                <div ref={hourListRef} className="flex-1 max-h-[120px] overflow-y-auto scrollbar-thin">
                  <div className="grid grid-cols-4 gap-0.5">
                    {generateHourOptions().map((h) => (
                      <button
                        key={h}
                        type="button"
                        data-active={selectedHour === h}
                        onClick={() => handleHourClick(h)}
                        className={`
                          text-xs py-1 rounded text-center transition-colors
                          ${selectedHour === h ? 'bg-geek-400/20 text-geek-300' : 'text-neutral-400 hover:bg-neutral-700/80'}
                          ${!selectedDate ? 'opacity-50 cursor-not-allowed' : ''}
                        `
                          .trim()
                          .replace(/\s+/g, ' ')}
                        disabled={!selectedDate}
                      >
                        {padZero(h)}
                      </button>
                    ))}
                  </div>
                </div>

                <span className="text-neutral-500 text-lg font-mono">:</span>

                {/* Minutes */}
                <div ref={minuteListRef} className="flex-1 max-h-[120px] overflow-y-auto scrollbar-thin">
                  <div className="grid grid-cols-4 gap-0.5">
                    {generateMinuteOptions().map((m) => (
                      <button
                        key={m}
                        type="button"
                        data-active={selectedMinute === m}
                        onClick={() => handleMinuteClick(m)}
                        className={`
                          text-xs py-1 rounded text-center transition-colors
                          ${selectedMinute === m ? 'bg-geek-400/20 text-geek-300' : 'text-neutral-400 hover:bg-neutral-700/80'}
                          ${!selectedDate ? 'opacity-50 cursor-not-allowed' : ''}
                        `
                          .trim()
                          .replace(/\s+/g, ' ')}
                        disabled={!selectedDate}
                      >
                        {padZero(m)}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            {/* Footer */}
            <div className="flex items-center justify-between px-3 py-2 border-t border-neutral-700">
              <button
                type="button"
                onClick={handleNow}
                className="text-xs text-geek-400 hover:text-geek-300 transition-colors"
              >
                {t('common.dateTimeInput.now')}
              </button>
              <button
                type="button"
                onClick={() => setIsOpen(false)}
                className="text-xs px-3 py-1 rounded bg-geek-400/20 text-geek-300 hover:bg-geek-400/30 transition-colors"
              >
                {t('common.dateTimeInput.ok')}
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export default DateTimeInput;
