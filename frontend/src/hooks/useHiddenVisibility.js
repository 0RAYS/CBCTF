import { useEffect, useState } from 'react';

/** 记住是否显示排行榜中被隐藏的条目，选择保存在浏览器本地。 */
export function useHiddenVisibility(storageKey) {
  const [visible, setVisible] = useState(() => window.localStorage.getItem(storageKey) === '1');

  useEffect(() => {
    window.localStorage.setItem(storageKey, visible ? '1' : '0');
  }, [storageKey, visible]);

  return [visible, setVisible];
}
