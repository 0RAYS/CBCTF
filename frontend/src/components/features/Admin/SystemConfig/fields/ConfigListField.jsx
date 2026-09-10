import { IconPlus, IconX } from '@tabler/icons-react';
import { Input, Button } from '../../../../common';
import { useTranslation } from 'react-i18next';

/**
 * ConfigListField - Generic list field component with add/remove capabilities
 * @param {string} label - Field label
 * @param {Array} items - Array of items
 * @param {Function} onAdd - Handler for adding new item
 * @param {Function} onUpdate - Handler for updating item (index, value)
 * @param {Function} onRemove - Handler for removing item (index)
 * @param {Function} renderItem - Optional custom renderer for items (item, index) => ReactNode
 */
export function ConfigListField({ label, items = [], onAdd, onUpdate, onRemove, renderItem }) {
  const { t } = useTranslation();

  return (
    <fieldset className="min-w-0 space-y-1">
      <legend className="w-full">
        <span className="flex items-center justify-between gap-2">
          <span className="text-xs font-mono text-neutral-400">{label}</span>
          <Button size="icon" variant="ghost" aria-label={t('common.add')} onClick={onAdd}>
            <IconPlus size={14} />
          </Button>
        </span>
      </legend>
      <div className="space-y-1">
        {items.map((item, index) => (
          <div key={`${label}-${index}`} className="flex items-center gap-2">
            {renderItem ? (
              renderItem(item, index)
            ) : (
              <Input
                size="sm"
                aria-label={`${label} ${index + 1}`}
                value={item ?? ''}
                onChange={(event) => onUpdate(index, event.target.value)}
              />
            )}
            <Button size="icon" variant="ghost" aria-label={t('common.remove')} onClick={() => onRemove(index)}>
              <IconX size={14} />
            </Button>
          </div>
        ))}
      </div>
    </fieldset>
  );
}
