/**
 * Motion tokens — three easing tiers for CBCTF
 *
 * T1: Page/route-level entrances   — expo-out, deliberate, ~400-500ms
 * T2: Component transitions        — quint-out, purposeful, ~220-280ms
 * T3: Micro-interactions           — quart-out, instant, ~100-150ms
 */

/** expo-out: confident, decisive — use for page entrances, hero content */
export const EASE_T1 = [0.16, 1, 0.3, 1];

/** quint-out: smooth, refined — use for modals, drawers, dropdowns, accordions */
export const EASE_T2 = [0.22, 1, 0.36, 1];

/** quart-out: snappy, responsive — use for hover, tap, toggle, small state changes */
export const EASE_T3 = [0.25, 1, 0.5, 1];

/** Canonical durations in ms */
export const DUR = {
  instant: 0.1, // 100ms — button tap, toggle feedback
  fast: 0.18, // 180ms — icon swap, chip state change
  mid: 0.25, // 250ms — component transition (modal, dropdown)
  slow: 0.4, // 400ms — page-level entrance, hero content
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

/** List item stagger helper: returns delay for index */
export const staggerDelay = (index, step = 0.04) => index * step;
