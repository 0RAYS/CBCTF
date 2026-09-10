import { useState } from 'react';
import { IconPlayerPlay } from '@tabler/icons-react';
import { Button, Modal } from '../../../common';
import { normalizeStartCount } from './generatorUtils.js';

export default function GeneratorStartDialog({ challenges, pendingOperation, onStart, onClose, text, t }) {
  const [counts, setCounts] = useState({});

  return (
    <Modal isOpen onClose={onClose} title={text('selectChallenges')}>
      <div className="flex flex-col gap-4">
        {challenges.length === 0 ? (
          <p className="text-neutral-400 text-sm py-4 text-center">{text('noDynamicChallenges')}</p>
        ) : (
          <div className="flex flex-col gap-2 max-h-72 overflow-y-auto">
            {challenges.map((challenge) => {
              const id = challenge.rand_id ?? challenge.id;
              return (
                <label
                  key={id}
                  className="flex flex-wrap items-center gap-3 px-3 py-2 rounded-md hover:bg-neutral-800 transition-colors"
                >
                  <span className="text-sm text-neutral-200 flex-1 min-w-0 break-words">
                    {challenge.title ?? challenge.name}
                  </span>
                  <span className="text-xs text-neutral-500 font-mono break-all">{id}</span>
                  <input
                    type="number"
                    min={0}
                    step={1}
                    value={counts[id] ?? 0}
                    disabled={pendingOperation === 'start'}
                    onChange={(event) =>
                      setCounts((previous) => ({ ...previous, [id]: normalizeStartCount(event.target.value) }))
                    }
                    className="w-16 bg-neutral-800 border border-neutral-700 rounded px-2 py-1 text-sm text-neutral-200 text-center focus:outline-none focus:border-geek-400"
                  />
                </label>
              );
            })}
          </div>
        )}
        <div className="flex flex-wrap justify-end gap-2 pt-2">
          <Button variant="ghost" size="sm" onClick={onClose} disabled={pendingOperation === 'start'}>
            {t('common.cancel')}
          </Button>
          <Button
            variant="primary"
            size="sm"
            onClick={() => onStart(counts)}
            loading={pendingOperation === 'start'}
            disabled={Boolean(pendingOperation) || !Object.values(counts).some((count) => count > 0)}
            icon={<IconPlayerPlay size={14} />}
          >
            {text('startSelected')}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
