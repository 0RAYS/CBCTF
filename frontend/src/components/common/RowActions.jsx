import Button from './Button';

/** The column and context menu share action callbacks, visibility and disabled state. */
export default function RowActions({ actions }) {
  return (
    <div className="flex flex-wrap items-center gap-2" onClick={(event) => event.stopPropagation()}>
      {actions
        .filter((action) => action && action.inline && !action.hidden)
        .map((action) => (
          <Button
            key={action.key}
            variant={action.danger ? 'danger' : 'ghost'}
            size={action.icon ? 'icon' : 'sm'}
            className="shrink-0"
            title={action.label}
            aria-label={action.label}
            disabled={action.disabled}
            onClick={() => action.onClick()}
          >
            {action.icon || action.label}
          </Button>
        ))}
    </div>
  );
}
