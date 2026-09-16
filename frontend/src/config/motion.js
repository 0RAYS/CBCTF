/**
 * Motion tokens for CBCTF
 *
 * T2: Component transitions        — quint-out, purposeful, ~220-280ms
 * T3: Micro-interactions           — quart-out, instant, ~100-150ms
 */

/** quint-out: smooth, refined — use for modals, drawers, dropdowns, accordions */
export const EASE_T2 = [0.22, 1, 0.36, 1];

/** quart-out: snappy, responsive — use for hover, tap, toggle, small state changes */
export const EASE_T3 = [0.25, 1, 0.5, 1];

/** Transition durations in seconds */
const DUR = {
  fast: 0.18, // 180ms — icon swap, chip state change
  mid: 0.25, // 250ms — component transition (modal, dropdown)
};

/** Ready-made variant objects for AnimatePresence / motion components */

/** Backdrop fade */
export const backdropVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1, transition: { duration: DUR.mid, ease: EASE_T3 } },
  exit: { opacity: 0, transition: { duration: DUR.fast, ease: EASE_T3 } },
};

/** Panel/modal slide-up enter, slide-down exit */
export const panelVariants = {
  hidden: { opacity: 0, y: 16, scale: 0.97 },
  visible: { opacity: 1, y: 0, scale: 1, transition: { duration: DUR.mid, ease: EASE_T2 } },
  exit: { opacity: 0, y: 12, scale: 0.97, transition: { duration: DUR.fast, ease: EASE_T3 } },
};

/** Toast: slide down from top */
export const toastVariants = {
  hidden: { opacity: 0, y: -14, scale: 0.96 },
  visible: { opacity: 1, y: 0, scale: 1, transition: { duration: DUR.mid, ease: EASE_T2 } },
  exit: { opacity: 0, y: -8, scale: 0.97, transition: { duration: DUR.fast, ease: EASE_T3 } },
};
