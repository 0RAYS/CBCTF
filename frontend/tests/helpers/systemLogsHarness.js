import { useRef } from './pollingHookHarness.js';

export { useState, useRef, useEffect, useEffect as useLayoutEffect, toast } from './pollingHookHarness.js';
export const Button = 'button';
export const AnsiLog = 'pre';
export const IconRefresh = 'icon';
const t = (key) => key;
export const useTranslation = () => ({ t });
export default function useIpLookup() {
  return { lookup() {}, close() {}, dialogProps: {} };
}

export function useCallback(callback, deps) {
  const ref = useRef();
  if (!ref.current || deps.some((dep, index) => !Object.is(dep, ref.current.deps[index]))) {
    ref.current = { callback, deps };
  }
  return ref.current.callback;
}

export function findElement(element, predicate) {
  if (!element || typeof element !== 'object') return undefined;
  if (predicate(element)) return element;
  const children = [element.props?.children, element.props?.sentinel].flat(Infinity);
  for (const child of children) {
    const found = findElement(child, predicate);
    if (found) return found;
  }
}
