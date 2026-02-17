import { IconPlus, IconX } from '@tabler/icons-react';
import { Input, Button } from '../../../../common';

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
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between">
        <span className="text-xs font-mono text-neutral-400">{label}</span>
        <Button size="icon" variant="ghost" onClick={onAdd}>
          <IconPlus size={14} />
        </Button>
      </div>
      <div className="space-y-1">
        {items.map((item, index) => (
          <div key={`${label}-${index}`} className="flex items-center gap-2">
            {renderItem ? (
              renderItem(item, index)
            ) : (
              <Input size="sm" value={item ?? ''} onChange={(event) => onUpdate(index, event.target.value)} />
            )}
            <Button size="icon" variant="ghost" onClick={() => onRemove(index)}>
              <IconX size={14} />
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}
